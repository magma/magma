/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the Apache License, Version 2.0  (the "License"); you may not use this file
 * except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
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

/*! \file sgw_task.c
  \brief
  \author Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/
#define PGW
#define PGW_TASK_C

#include <stdio.h>
#include <sys/types.h>

#include "log.h"
#include "intertask_interface.h"
#include "pgw_defs.h"
#include "pgw_handlers.h"
#include "sgw.h"
#include "common_defs.h"
#include "bstrlib.h"
#include "intertask_interface_types.h"
#include "spgw_config.h"
#include "spgw_state.h"
#include "assertions.h"

extern __pid_t g_pid;

static void pgw_exit(void);

//------------------------------------------------------------------------------
static void *pgw_intertask_interface(void *args_p)
{
  itti_mark_task_ready(TASK_PGW_APP);

  spgw_state_t *spgw_state_p;

  while (1) {
    MessageDef *received_message_p = NULL;

    itti_receive_msg(TASK_PGW_APP, &received_message_p);

    imsi64_t imsi64 = itti_get_associated_imsi(received_message_p);
    OAILOG_DEBUG(
      LOG_PGW_APP,
      "Received message with imsi: " IMSI_64_FMT,
      imsi64);

    if (ITTI_MSG_ID(received_message_p) != TERMINATE_MESSAGE) {
      spgw_state_p = get_spgw_state(false);
      AssertFatal(
        spgw_state_p != NULL, "Failed to retrieve SPGW state on PGW task");
    }

    switch (ITTI_MSG_ID(received_message_p)) {
      case TERMINATE_MESSAGE: {
        pgw_exit();
        OAI_FPRINTF_INFO("TASK_PGW terminated\n");
        itti_exit_task();
      } break;

      default: {
        OAILOG_DEBUG(
          LOG_PGW_APP,
          "Unkwnon message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p),
          ITTI_MSG_NAME(received_message_p));
      } break;
    }

    itti_free(ITTI_MSG_ORIGIN_ID(received_message_p), received_message_p);
    received_message_p = NULL;
  }

  return NULL;
}

int pgw_init(spgw_config_t *spgw_config_pP)
{
  if (itti_create_task(TASK_PGW_APP, &pgw_intertask_interface, NULL) < 0) {
    perror("pthread_create");
    OAILOG_ALERT(LOG_PGW_APP, "Initializing PGW-APP task interface: ERROR\n");
    return RETURNerror;
  }

  FILE *fp = NULL;
  bstring filename = bformat("/tmp/pgw_%d.status", g_pid);
  fp = fopen(bdata(filename), "w+");
  bdestroy(filename);
  fprintf(fp, "STARTED\n");
  fflush(fp);
  fclose(fp);

  OAILOG_DEBUG(LOG_PGW_APP, "Initializing PGW-APP task interface: DONE\n");
  return RETURNok;
}

static void pgw_exit(void)
{
  return;
}
