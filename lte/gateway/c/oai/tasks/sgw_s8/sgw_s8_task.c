/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*! \file sgw_s8_task.c
  \brief
  \author
  \company
  \email:
*/
#define SGW_S8
#define SGW_S8_TASK_C

#include <stdio.h>

#include "log.h"
#include "assertions.h"
#include "common_defs.h"
#include "itti_free_defined_msg.h"
#include "sgw_s8_defs.h"

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg);
static void sgw_s8_exit(void);
task_zmq_ctx_t sgw_s8_task_zmq_ctx;

static void* sgw_s8_thread(void* args) {
  itti_mark_task_ready(TASK_SGW_S8);
  init_task_context(
      TASK_SGW_S8, (task_id_t[]){TASK_MME_APP}, 1, handle_message,
      &sgw_s8_task_zmq_ctx);

  zloop_start(sgw_s8_task_zmq_ctx.event_loop);
  sgw_s8_exit();
  return NULL;
}

int sgw_s8_init(void) {
  OAILOG_DEBUG(LOG_SGW_S8, "Initializing SGW-S8 interface\n");

  if (itti_create_task(TASK_SGW_S8, &sgw_s8_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_SGW_S8, "Failed to create sgw_s8 task\n");
    return RETURNerror;
  }
  OAILOG_DEBUG(LOG_SGW_S8, "Done initialization of SGW_S8 interface\n");
  return RETURNok;
}

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  zframe_t* msg_frame = zframe_recv(reader);
  assert(msg_frame);
  MessageDef* received_message_p = (MessageDef*) zframe_data(msg_frame);

  switch (ITTI_MSG_ID(received_message_p)) {
    case MESSAGE_TEST: {
      OAI_FPRINTF_INFO("TASK_SGW_S8 received MESSAGE_TEST\n");
    } break;
    default: {
      OAILOG_DEBUG(
          LOG_S6A, "Unkwnon message ID %d: %s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  itti_free_msg_content(received_message_p);
  zframe_destroy(&msg_frame);
  return 0;
}

//------------------------------------------------------------------------------
static void sgw_s8_exit(void) {
  destroy_task_context(&sgw_s8_task_zmq_ctx);
  OAILOG_DEBUG(LOG_SGW_S8, "Finished cleaning up SGW_S8 task \n");
  OAI_FPRINTF_INFO("TASK_SGW_S8 terminated\n");
  pthread_exit(NULL);
}
