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

#include "lte/gateway/c/core/oai/lib/message_utils/service303_message_utils.h"

#include <stddef.h>

#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#include "lte/gateway/c/core/oai/common/log.h"

int send_app_health_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id, bool healthy) {
  MessageDef* message_p;
  if (healthy) {
    message_p = DEPRECATEDitti_alloc_new_message_fatal(
        origin_id, APPLICATION_HEALTHY_MSG);
  } else {
    message_p = DEPRECATEDitti_alloc_new_message_fatal(
        origin_id, APPLICATION_UNHEALTHY_MSG);
  }
  return send_msg_to_task(task_zmq_ctx_p, TASK_SERVICE303, message_p);
}

int send_mme_app_stats_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id,
    application_mme_app_stats_msg_t* stats_msg) {
  MessageDef* message_p =
      itti_alloc_new_message(origin_id, APPLICATION_MME_APP_STATS_MSG);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unable to allocate memory");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  message_p->ittiMsg.application_mme_app_stats_msg.nb_ue_attached =
      stats_msg->nb_ue_attached;
  message_p->ittiMsg.application_mme_app_stats_msg.nb_ue_connected =
      stats_msg->nb_ue_connected;
  message_p->ittiMsg.application_mme_app_stats_msg.nb_default_eps_bearers =
      stats_msg->nb_default_eps_bearers;
  message_p->ittiMsg.application_mme_app_stats_msg.nb_s1u_bearers =
      stats_msg->nb_s1u_bearers;
  message_p->ittiMsg.application_mme_app_stats_msg.nb_mme_app_last_msg_latency =
      stats_msg->nb_mme_app_last_msg_latency;
  return send_msg_to_task(task_zmq_ctx_p, TASK_SERVICE303, message_p);
}

int send_s1ap_stats_to_service303(
    task_zmq_ctx_t* task_zmq_ctx_p, task_id_t origin_id,
    application_s1ap_stats_msg_t* stats_msg) {
  MessageDef* message_p =
      itti_alloc_new_message(origin_id, APPLICATION_S1AP_STATS_MSG);
  if (message_p == NULL) {
    OAILOG_ERROR(LOG_MME_APP, "Unable to allocate memory");
    OAILOG_FUNC_RETURN(LOG_MME_APP, RETURNerror);
  }
  message_p->ittiMsg.application_s1ap_stats_msg.nb_enb_connected =
      stats_msg->nb_enb_connected;
  message_p->ittiMsg.application_s1ap_stats_msg.nb_s1ap_last_msg_latency =
      stats_msg->nb_s1ap_last_msg_latency;
  return send_msg_to_task(task_zmq_ctx_p, TASK_SERVICE303, message_p);
}
