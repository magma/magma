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

#include <string.h>

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/common/assertions.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/itti_types.h"
#ifdef __cplusplus
}
#endif

#include "lte/gateway/c/core/oai/include/sgs_messages_types.hpp"
#include "lte/gateway/c/core/oai/tasks/sgs/sgs_messages.hpp"

static void sgs_send_sgsap_vlr_reset_ack(void);

status_code_e handle_sgs_location_update_accept(
    const itti_sgsap_location_update_acc_t* itti_sgsap_location_update_acc_p) {
  /* Received SGS Location Update Accept from FedGW
   *send it to MME App for further processing
   */
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS,
                                                     SGSAP_LOCATION_UPDATE_ACC);
  memset((void*)&message_p->ittiMsg.sgsap_location_update_acc, 0,
         sizeof(itti_sgsap_location_update_acc_t));
  memcpy(&message_p->ittiMsg.sgsap_location_update_acc,
         itti_sgsap_location_update_acc_p,
         sizeof(itti_sgsap_location_update_acc_t));

  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS Location Update Acc message from FedGW with IMSI %s\n",
      itti_sgsap_location_update_acc_p->imsi);
  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_SGS, rc);
}

//------------------------------------------------------------------------------------------------------------
status_code_e handle_sgs_location_update_reject(
    const itti_sgsap_location_update_rej_t* itti_sgsap_loc_updt_rej_p) {
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;

  /* Received SGS Location Update Reject from FedGW
   *send it to MME App for further processing
   */
  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS,
                                                     SGSAP_LOCATION_UPDATE_REJ);
  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS Location Update Reject message from FedGW with IMSI %s\n",
      itti_sgsap_loc_updt_rej_p->imsi);
  memset((void*)&message_p->ittiMsg.sgsap_location_update_rej, 0,
         sizeof(itti_sgsap_location_update_rej_t));
  memcpy(&message_p->ittiMsg.sgsap_location_update_rej,
         itti_sgsap_loc_updt_rej_p, sizeof(itti_sgsap_location_update_rej_t));

  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_SGS, rc);
}

status_code_e handle_sgs_eps_detach_ack(
    const itti_sgsap_eps_detach_ack_t* sgsap_eps_detach_ack_p) {
  // send it to MME module for further processing
  status_code_e rc = RETURNok;
  MessageDef* message_p = NULL;
  itti_sgsap_eps_detach_ack_t* sgs_eps_detach_ack_p = NULL;

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_S6A, SGSAP_EPS_DETACH_ACK);
  sgs_eps_detach_ack_p = &message_p->ittiMsg.sgsap_eps_detach_ack;
  memset((void*)sgs_eps_detach_ack_p, 0, sizeof(itti_sgsap_eps_detach_ack_t));
  OAILOG_DEBUG(LOG_SGS, "Received SGS EPS Detach Ack message from FedGW\n");
  memcpy(sgs_eps_detach_ack_p, sgsap_eps_detach_ack_p,
         sizeof(itti_sgsap_eps_detach_ack_t));

  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_SGS, rc);
}

status_code_e handle_sgs_imsi_detach_ack(
    const itti_sgsap_imsi_detach_ack_t* sgsap_imsi_detach_ack_p) {
  // send it to MME module for further processing
  status_code_e rc = RETURNok;
  MessageDef* message_p = NULL;
  itti_sgsap_imsi_detach_ack_t* sgs_imsi_detach_ack_p = NULL;

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_S6A, SGSAP_IMSI_DETACH_ACK);
  sgs_imsi_detach_ack_p = &message_p->ittiMsg.sgsap_imsi_detach_ack;
  memset((void*)sgs_imsi_detach_ack_p, 0, sizeof(itti_sgsap_imsi_detach_ack_t));
  OAILOG_DEBUG(LOG_SGS, "Received SGS IMSI Detach Ack message from FedGW\n");
  memcpy(sgs_imsi_detach_ack_p, sgsap_imsi_detach_ack_p,
         sizeof(itti_sgsap_imsi_detach_ack_t));

  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_SGS, rc);
}

status_code_e handle_sgs_downlink_unitdata(
    const itti_sgsap_downlink_unitdata_t* sgs_dl_unitdata_p) {
  status_code_e rc = RETURNok;

  MessageDef* message_p = NULL;
  itti_sgsap_downlink_unitdata_t* sgs_dl_unit_data_p = NULL;

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_DOWNLINK_UNITDATA);
  sgs_dl_unit_data_p = &message_p->ittiMsg.sgsap_downlink_unitdata;
  memset((void*)sgs_dl_unit_data_p, 0, sizeof(itti_sgsap_downlink_unitdata_t));
  OAILOG_DEBUG(LOG_SGS, "Received SGS Downlink UnitData message from FedGW\n");
  memcpy(sgs_dl_unit_data_p, sgs_dl_unitdata_p,
         sizeof(itti_sgsap_downlink_unitdata_t));
  // send it to NAS module for further processing
  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

status_code_e handle_sgs_release_req(
    const itti_sgsap_release_req_t* sgs_release_req_p) {
  status_code_e rc = RETURNok;

  MessageDef* message_p = NULL;
  itti_sgsap_release_req_t* sgs_rel_req_p = NULL;

  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_RELEASE_REQ);
  sgs_rel_req_p = &message_p->ittiMsg.sgsap_release_req;
  memset((void*)sgs_rel_req_p, 0, sizeof(itti_sgsap_release_req_t));
  OAILOG_DEBUG(LOG_SGS, "Received SGS Release Request message from FedGW\n");
  memcpy(sgs_rel_req_p, sgs_release_req_p, sizeof(itti_sgsap_release_req_t));
  // send it to NAS module for further processing
  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
/* Fed GW calls below function, on reception of MM Information Request from
 * MSC/VLR
 */

status_code_e handle_sgs_mm_information_request(
    const itti_sgsap_mm_information_req_t* mm_information_req_pP) {
  /* Received SGS MM Information Request from FedGW
   *send it to NAS task for further processing
   */
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SGS);

  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS,
                                                     SGSAP_MM_INFORMATION_REQ);
  itti_sgsap_mm_information_req_t* mm_information_req_p =
      &message_p->ittiMsg.sgsap_mm_information_req;
  memset((void*)mm_information_req_p, 0,
         sizeof(itti_sgsap_mm_information_req_t));

  memcpy((void*)mm_information_req_p, (void*)mm_information_req_pP,
         sizeof(itti_sgsap_mm_information_req_t));
  OAILOG_DEBUG(
      LOG_SGS,
      "Received MM Information Request message from FedGW and send to NAS for "
      "Imsi :%s \n",
      mm_information_req_p->imsi);

  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_SGS, rc);
}

//------------------------------------------------------------------------------------------------------------
status_code_e handle_sgs_service_abort_req(
    const itti_sgsap_service_abort_req_t* itti_sgsap_service_abort_req_p) {
  /* Received SGS SERVICE ABORT Request from FedGW
   *send it to MME App for further processing
   */
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;

  OAILOG_FUNC_IN(LOG_SGS);
  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_SERVICE_ABORT_REQ);
  memset((void*)&message_p->ittiMsg.sgsap_service_abort_req, 0,
         sizeof(itti_sgsap_service_abort_req_t));
  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS SERVICE ABORT Req message from FedGW for IMSI %s\n",
      itti_sgsap_service_abort_req_p->imsi);

  memcpy(&message_p->ittiMsg.sgsap_service_abort_req,
         itti_sgsap_service_abort_req_p,
         sizeof(itti_sgsap_service_abort_req_t));

  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_SGS, rc);
}

/* handle_sgs_paging_request()
 * Is function invoked by FedGW on reception of Paging Request from MSC/VLR
 */

status_code_e handle_sgs_paging_request(
    const itti_sgsap_paging_request_t* const sgs_paging_req_pP) {
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SGS);

  /* Received SGS Paging Request from FedGW
   *send it to MME App for further processing
   */
  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_PAGING_REQUEST);

  itti_sgsap_paging_request_t* sgs_paging_req_p =
      &message_p->ittiMsg.sgsap_paging_request;
  memset((void*)sgs_paging_req_p, 0, sizeof(itti_sgsap_paging_request_t));

  memcpy((void*)sgs_paging_req_p, (void*)sgs_paging_req_pP,
         sizeof(itti_sgsap_paging_request_t));

  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS Paging Request message from FedGW and send Paging request "
      "to "
      "MME app"
      "for Imsi :%s \n",
      sgs_paging_req_p->imsi);
  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/* handle_sgs_vlr_reset_indication()
 * Is Function invoked by FedGW on reception of Reset Indication from MSC/VLR
 */

status_code_e handle_sgs_vlr_reset_indication(
    const itti_sgsap_vlr_reset_indication_t* const sgs_vlr_reset_ind_pP) {
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SGS);

  /* Received SGS VLR Reset Indication from FedGW
   * send it to MME App for further processing
   */
  message_p = DEPRECATEDitti_alloc_new_message_fatal(
      TASK_SGS, SGSAP_VLR_RESET_INDICATION);

  itti_sgsap_vlr_reset_indication_t* sgs_vlr_reset_ind_p =
      &message_p->ittiMsg.sgsap_vlr_reset_indication;
  memset((void*)sgs_vlr_reset_ind_p, 0,
         sizeof(itti_sgsap_vlr_reset_indication_t));

  memcpy((void*)sgs_vlr_reset_ind_p, (void*)sgs_vlr_reset_ind_pP,
         sizeof(itti_sgsap_vlr_reset_indication_t));

  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS Reset Indication message from FedGW and send Reset "
      "Indication to MME app"
      "for vlr name :%s \n",
      sgs_vlr_reset_ind_p->vlr_name);
  if ((rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p)) !=
      RETURNok) {
    OAILOG_ERROR(
        LOG_SGS,
        "Failed to send SGSAP_VLR_RESET_INDICATION for vlr_name :%s \n",
        sgs_vlr_reset_ind_p->vlr_name);
  }
  /* Send SGSAP Reset Ack to VLR */
  sgs_send_sgsap_vlr_reset_ack();
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/* handle_sgs_status_message()
 * Is function invoked by FedGW on reception of Sgs Status message from MSC/VLR
 */

status_code_e handle_sgs_status_message(
    const itti_sgsap_status_t* sgs_status_pP) {
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SGS);

  /* Received SGS status message from FedGW
   * send it to MME App for further processing
   */
  message_p = DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_STATUS);

  itti_sgsap_status_t* sgs_status_p = &message_p->ittiMsg.sgsap_status;
  memset((void*)sgs_status_p, 0, sizeof(itti_sgsap_status_t));
  memcpy((void*)sgs_status_p, (void*)sgs_status_pP,
         sizeof(itti_sgsap_status_t));
  OAILOG_DEBUG(LOG_SGS,
               "Received SGS Status message from FedGW "
               "and send sgs status message to MME app for Imsi :%s \n",
               sgs_status_p->imsi);
  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/* _sgs_send_sgsap_vlr_reset_ack()
 * Send VLR Reset Ack In response to Reset Indication from VLR
 */

static void sgs_send_sgsap_vlr_reset_ack(void) {
  MessageDef* message_p = NULL;
  itti_sgsap_vlr_reset_ack_t* sgsap_reset_ack_pP = NULL;

  OAILOG_FUNC_IN(LOG_SGS);
  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_VLR_RESET_ACK);
  sgsap_reset_ack_pP = &message_p->ittiMsg.sgsap_vlr_reset_ack;
  memset((void*)sgsap_reset_ack_pP, 0, sizeof(itti_sgsap_vlr_reset_ack_t));

  /* Should  fill mme_name in sgs_service */
  OAILOG_INFO(LOG_SGS, "Send SGSAP-Reset Ack to SGS Service \n");
  /* send Reset Ack message to FeG */
  /* Below line should be un-commented, once fed GW or MSC supports VLR failure
   * procedure */
  // send_reset_ack(sgsap_reset_ack_pP);

  OAILOG_FUNC_OUT(LOG_SGS);
}

status_code_e handle_sgsap_alert_request(
    const itti_sgsap_alert_request_t* const sgsap_alert_request) {
  MessageDef* message_p = NULL;
  status_code_e rc = RETURNok;
  OAILOG_FUNC_IN(LOG_SGS);

  /* Received SGS Alert Req from FedGW
   *send it to MME App for further processing
   */
  message_p =
      DEPRECATEDitti_alloc_new_message_fatal(TASK_SGS, SGSAP_ALERT_REQUEST);

  memset((void*)&message_p->ittiMsg.sgsap_alert_request, 0,
         sizeof(itti_sgsap_alert_request_t));

  memcpy((void*)&message_p->ittiMsg.sgsap_alert_request,
         (void*)sgsap_alert_request, sizeof(itti_sgsap_alert_request_t));

  OAILOG_DEBUG(
      LOG_SGS,
      "Received SGS Alert Request message from FedGW and send Alert request to "
      "MME app"
      "for Imsi :%s \n",
      sgsap_alert_request->imsi);
  rc = send_msg_to_task(&sgs_task_zmq_ctx, TASK_MME_APP, message_p);

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
