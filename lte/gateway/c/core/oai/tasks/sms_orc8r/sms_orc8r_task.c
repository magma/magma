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

/*! \file sms_orc8r_task.c
  \brief
  \author
  \company
  \email:
*/
#define SMS_ORC8R
#define SMS_ORC8R_TASK_C

#include <stdio.h>

#include "log.h"
#include "intertask_interface.h"
#include "mme_config.h"
#include "sgs_messages_types.h"
#include "sms_orc8r_client_api.h"
#include "common_defs.h"
#include "intertask_interface_types.h"

static void sms_orc8r_exit(void);

task_zmq_ctx_t sms_orc8r_task_zmq_ctx;

static int handle_message(zloop_t* loop, zsock_t* reader, void* arg) {
  MessageDef* received_message_p = receive_msg(reader);

  switch (ITTI_MSG_ID(received_message_p)) {
    case SGSAP_UPLINK_UNITDATA: {
      /*
       * We received a SGs uplink unitdata message from NAS
       * * * * procedures might be:
       * * * *      Mobile origination SMS - Uplink Nas Transport message
       * * * *      Mobile terminating SMS - Uplink Nas Transport message
       */
      OAILOG_DEBUG(LOG_SMS_ORC8R, "Received SGSAP_UPLINK_UNITDATA message \n");
      send_smo_uplink_unitdata(
          &received_message_p->ittiMsg.sgsap_uplink_unitdata);
    } break;

    case TERMINATE_MESSAGE: {
      free(received_message_p);
      sms_orc8r_exit();
    } break;

    default: {
      OAILOG_DEBUG(
          LOG_SMS_ORC8R, "Unknown message ID %d:%s\n",
          ITTI_MSG_ID(received_message_p), ITTI_MSG_NAME(received_message_p));
    } break;
  }

  free(received_message_p);
  return 0;
}

//------------------------------------------------------------------------------
static void* sms_orc8r_thread(__attribute__((unused)) void* args_p) {
  task_zmq_ctx_t* task_zmq_ctx_p = &sms_orc8r_task_zmq_ctx;

  itti_mark_task_ready(TASK_SMS_ORC8R);
  init_task_context(
      TASK_SMS_ORC8R, (task_id_t[]){TASK_MME_APP}, 1, handle_message,
      task_zmq_ctx_p);

  zloop_start(task_zmq_ctx_p->event_loop);
  sms_orc8r_exit();
  return NULL;
}

//------------------------------------------------------------------------------
int sms_orc8r_init(const mme_config_t* mme_config_p) {
  OAILOG_DEBUG(LOG_SMS_ORC8R, "Initializing SMS_ORC8R task interface\n");

  if (itti_create_task(TASK_SMS_ORC8R, &sms_orc8r_thread, NULL) < 0) {
    OAILOG_ERROR(LOG_SMS_ORC8R, "sms_orc8r create task\n");
    return RETURNerror;
  }
  OAILOG_DEBUG(LOG_SMS_ORC8R, "Initializing SMS_ORC8R task interface: DONE\n");
  return RETURNok;
}

//------------------------------------------------------------------------------
static void sms_orc8r_exit(void) {
  destroy_task_context(&sms_orc8r_task_zmq_ctx);
  OAI_FPRINTF_INFO("TASK_SMS_ORC8R terminated\n");
  pthread_exit(NULL);
}
