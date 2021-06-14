/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */

/*! \file async_system.c
   \brief
   \author  Lionel GAUTHIER
   \date 2017
   \email: lionel.gauthier@eurecom.fr
*/
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>
#include <stdarg.h>

#include "bstrlib.h"
#include "intertask_interface.h"
#include "log.h"
#include "async_system.h"
#include "assertions.h"
#include "dynamic_memory_check.h"
#include "itti_free_defined_msg.h"
#include "common_defs.h"
#include "async_system_messages_types.h"
#include "intertask_interface_types.h"
#include "itti_types.h"

static void async_system_exit(void);

task_zmq_ctx_t async_system_task_zmq_ctx;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  int rc                         = 0;
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case ASYNC_SYSTEM_COMMAND: {
      OAILOG_DEBUG(
          LOG_ASYNC_SYSTEM, "C system() call: %s\n",
          bdata(ASYNC_SYSTEM_COMMAND(received_message_p).system_command));
      rc = system(
          bdata(ASYNC_SYSTEM_COMMAND(received_message_p).system_command));

      if (rc) {
        OAILOG_ERROR(
            LOG_ASYNC_SYSTEM, "ERROR in system command %s: %d\n",
            bdata(ASYNC_SYSTEM_COMMAND(received_message_p).system_command), rc);
        if (ASYNC_SYSTEM_COMMAND(received_message_p).is_abort_on_error) {
          bdestroy_wrapper(
              &ASYNC_SYSTEM_COMMAND(received_message_p).system_command);
          exit(-1);  // may be not exit
        }
        bdestroy_wrapper(
            &ASYNC_SYSTEM_COMMAND(received_message_p).system_command);
      }
    } break;

    case TERMINATE_MESSAGE: {
      itti_free_msg_content(received_message_p);
      free(received_message_p);
      async_system_exit();
    } break;

    default: {
      OAILOG_DEBUG(
          LOG_ASYNC_SYSTEM, "Unknown message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* async_system_thread(__attribute__((unused)) void* args_p) {
  itti_mark_task_ready(TASK_ASYNC_SYSTEM);
  init_task_context(
      TASK_ASYNC_SYSTEM, (task_id_t[]){}, 0, handle_message,
      &async_system_task_zmq_ctx);

  zloop_start(async_system_task_zmq_ctx.event_loop);
  async_system_exit();
  return NULL;
}

//------------------------------------------------------------------------------
status_code_e async_system_init(void) {
  OAI_FPRINTF_INFO("Initializing ASYNC_SYSTEM\n");
  if (itti_create_task(TASK_ASYNC_SYSTEM, &async_system_thread, NULL) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(
        LOG_ASYNC_SYSTEM, "Initializing ASYNC_SYSTEM task interface: ERROR\n");
    return RETURNerror;
  }
  OAI_FPRINTF_INFO("Initializing ASYNC_SYSTEM Done\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
status_code_e async_system_command(
    int sender_itti_task, bool is_abort_on_error, char* format, ...) {
  va_list args;
  int rv       = 0;
  bstring bstr = NULL;
  va_start(args, format);
  bstr = bfromcstralloc(1024, " ");
  btrunc(bstr, 0);
  rv = bvcformata(bstr, 1024, format, args);  // big number, see bvcformata
  va_end(args);

  if (NULL == bstr || BSTR_OK != rv) {
    OAILOG_ERROR(LOG_ASYNC_SYSTEM, "Error while formatting system command");
    return RETURNerror;
  }
  MessageDef* message_p = NULL;
  message_p = itti_alloc_new_message(sender_itti_task, ASYNC_SYSTEM_COMMAND);
  ASYNC_SYSTEM_COMMAND(message_p).system_command    = bstr;
  ASYNC_SYSTEM_COMMAND(message_p).is_abort_on_error = is_abort_on_error;
  status_code_e result                              = send_msg_to_task(
      &async_system_task_zmq_ctx, TASK_ASYNC_SYSTEM, message_p);
  return result;
}

//------------------------------------------------------------------------------
void async_system_exit(void) {
  destroy_task_context(&async_system_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_ASYNC_SYSTEM terminated");
  pthread_exit(NULL);
}
