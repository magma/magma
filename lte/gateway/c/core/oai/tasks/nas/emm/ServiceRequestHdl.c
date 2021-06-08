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

#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <stdlib.h>

#include "log.h"
#include "common_types.h"
#include "common_defs.h"
#include "mme_app_ue_context.h"
#include "emm_proc.h"
#include "emm_data.h"
#include "emm_sap.h"
#include "emm_cause.h"
#include "mme_app_itti_messaging.h"
#include "service303.h"
#include "conversions.h"
#include "3gpp_23.003.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "DetachRequest.h"
#include "ExtendedServiceRequest.h"
#include "ServiceType.h"
#include "emm_asDef.h"
#include "esm_data.h"
#include "mme_api.h"
#include "nas_message.h"
#include "mme_app_defs.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/
static int emm_service_reject(mme_ue_s1ap_id_t ue_id, uint8_t emm_cause);

static int check_paging_received_without_lai(mme_ue_s1ap_id_t ue_id);
/*
   --------------------------------------------------------------------------
    Internal data handled by the service request procedure in the UE
   --------------------------------------------------------------------------
*/

/*
   --------------------------------------------------------------------------
    Internal data handled by the service request procedure in the MME
   --------------------------------------------------------------------------
*/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/
int emm_proc_service_reject(
    const mme_ue_s1ap_id_t ue_id, const uint8_t emm_cause) {
  int rc = RETURNok;
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  rc = emm_service_reject(ue_id, emm_cause);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
/** \fn void _emm_service_reject(void *args);
    \brief Performs the  SR  procedure not accepted by the network.
     @param [in]args UE EMM context data
     @returns status of operation
*/
static int emm_service_reject(mme_ue_s1ap_id_t ue_id, uint8_t emm_cause)

{
  int rc = RETURNerror;
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  emm_context_t* emm_ctx = emm_context_get(&_emm_data, ue_id);
  emm_sap_t emm_sap      = {0};

  OAILOG_WARNING(
      LOG_NAS_EMM,
      "EMM-PROC- Sending Service Reject. ue_id=" MME_UE_S1AP_ID_FMT
      ", cause=%d)\n",
      ue_id, emm_cause);
  /*
   * Notify EMM-AS SAP that Service Reject message has to be sent
   * onto the network
   */
  emm_sap.primitive                        = EMMAS_ESTABLISH_REJ;
  emm_sap.u.emm_as.u.establish.ue_id       = ue_id;
  emm_sap.u.emm_as.u.establish.eps_id.guti = NULL;

  emm_sap.u.emm_as.u.establish.emm_cause = emm_cause;
  emm_sap.u.emm_as.u.establish.nas_info  = EMM_AS_NAS_INFO_SR;
  emm_sap.u.emm_as.u.establish.nas_msg   = NULL;
  /*
   * Setup EPS NAS security data
   */
  if (emm_ctx) {
    emm_as_set_security_data(
        &emm_sap.u.emm_as.u.establish.sctx, &emm_ctx->_security, false, false);
  } else {
    emm_as_set_security_data(
        &emm_sap.u.emm_as.u.establish.sctx, NULL, false, false);
  }
  rc = emm_sap_send(&emm_sap);

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_proc_extended_service_request()                           **
 **                                                                        **
 ** Description: process the extended service request message received in  **
 **                Uplink nas message                       ,              **
 **              check validation for ue attach type                       **
 **              and EPC supports CS or SMS                                **
 **              send the extended service req notification to mme_app     **
 ***************************************************************************/
int emm_proc_extended_service_request(
    const mme_ue_s1ap_id_t ue_id, const extended_service_request_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                 = RETURNok;
  emm_context_t* emm_ctx = NULL;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMM-PROC- Extended Service Request for the UE (ue_id=" MME_UE_S1AP_ID_FMT
      ") \n",
      ue_id);
  /*
   * Get the UE context
   */
  emm_ctx = emm_context_get(&_emm_data, ue_id);

  if (!emm_ctx) {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "No EMM context exists for the UE (ue_id=" MME_UE_S1AP_ID_FMT ") \n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  /*
   * if CSFB Response is recieved for MT CSFB as accepted by ue,
   * and if neaf flag is true then send the itti message to SGS
   * For triggering SGS ue activity indication message towards MSC.
   * In case of Csfb as rejected by ue ,then SGS paging reject shall be sent.
   */

  if (msg->servicetype == MO_CS_FB) {
    /* neaf flag is true*/
    /* send the SGSAP Ue activity indication to MSC/VLR */
    if (mme_ue_context_get_ue_sgs_neaf(ue_id) == true) {
      char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
      IMSI_TO_STRING(&(emm_ctx->_imsi), imsi_str, IMSI_BCD_DIGITS_MAX + 1);
      mme_app_itti_sgsap_ue_activity_ind(imsi_str, strlen(imsi_str));
      mme_ue_context_update_ue_sgs_neaf(ue_id, false);
    }
  }

  /* check for the ue attach type  and epc supported feature type */
  if ((emm_ctx->attach_type != EMM_ATTACH_TYPE_COMBINED_EPS_IMSI) ||
      !(_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED)) {
    /* send the service reject to UE */
    rc = emm_proc_service_reject(ue_id, EMM_CAUSE_CONGESTION);
    increment_counter(
        "extended_service_request", 1, 2, "result", "failure", "cause",
        "emm_cause_congestion");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  // Handle extended service request received in ue connected mode
  mme_app_handle_nas_extended_service_req(
      ue_id, msg->servicetype, msg->csfbresponse);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_recv_initial_ext_service_request                          **
 **                                                                        **
 ** Description: process the extended service request message received in  **
 **              Initial UE message                         ,              **
 **              check validation for ue attach type                       **
 **              and EPC supports CS call or SMS                           **
 **              send the extended service req notification to mme_app     **
 ***************************************************************************/
int emm_recv_initial_ext_service_request(
    const mme_ue_s1ap_id_t ue_id, const extended_service_request_msg* msg,
    int* emm_cause, const nas_message_decode_status_t* decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                 = RETURNok;
  emm_context_t* emm_ctx = NULL;
  emm_sap_t emm_sap      = {0};

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMM-PROC- Extended Service Request for the UE (ue_id=" MME_UE_S1AP_ID_FMT
      ") \n",
      ue_id);
  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Received Extended Service Request message, Security context "
      "%s"
      "Integrity protected %s MAC matched %s Ciphered %s\n",
      (decode_status->security_context_available) ? "yes" : "no",
      (decode_status->integrity_protected_message) ? "yes" : "no",
      (decode_status->mac_matched) ? "yes" : "no",
      (decode_status->ciphered_message) ? "yes" : "no");

  /*
   * Get the UE context
   */
  emm_ctx = emm_context_get(&_emm_data, ue_id);

  if (!emm_ctx) {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "No EMM context exists for the UE (ue_id=" MME_UE_S1AP_ID_FMT ") \n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  if (msg->servicetype == MO_CS_FB) {
    /* neaf flag is true*/
    /* send the SGSAP Ue activity indication to MSC/VLR */
    if (mme_ue_context_get_ue_sgs_neaf(ue_id) == true) {
      char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
      IMSI_TO_STRING(&(emm_ctx->_imsi), imsi_str, IMSI_BCD_DIGITS_MAX + 1);
      mme_app_itti_sgsap_ue_activity_ind(imsi_str, strlen(imsi_str));
      mme_ue_context_update_ue_sgs_neaf(ue_id, false);
    }
  }

  /* check for the ue attach type  and epc supported feature type */
  if ((emm_ctx->attach_type != EMM_ATTACH_TYPE_COMBINED_EPS_IMSI) ||
      !(_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED)) {
    /* send the service reject to UE */
    rc = emm_proc_service_reject(ue_id, EMM_CAUSE_CONGESTION);
    increment_counter(
        "extended_service_request", 1, 2, "result", "failure", "cause",
        "emm_cause_congestion");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  *emm_cause = EMM_CAUSE_SUCCESS;
  if ((msg->servicetype == MT_CS_FB)) {
    if (!(EMM_CSFB_RSP_PRESENT & msg->presencemask)) {
      /* CSFB Resp Missing*/
      /*send the service reject to UE*/
      rc = emm_proc_service_reject(ue_id, EMM_CAUSE_CONDITIONAL_IE_ERROR);
      increment_counter(
          "extended_service_request", 1, 2, "result", "failure", "cause",
          "ue_csfb_response_missing");
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
    }
  }
  /* Check cs paging procedure initiated without LAI
   * If cs paging initiated without LAI, On reception of Extended service
   * request procedure send Detach Request ( IMSI Detach) to UE */
  if (check_paging_received_without_lai(ue_id)) {
    emm_sap.primitive = EMMCN_NW_INITIATED_DETACH_UE;
    emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.ue_id = ue_id;
    emm_sap.u.emm_cn.u.emm_cn_nw_initiated_detach.detach_type =
        NW_DETACH_TYPE_IMSI_DETACH;
    rc = emm_sap_send(&emm_sap);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
    emm_ctx->_security.kenb_ul_count = emm_ctx->_security.ul_count;
  }

  emm_sap.primitive                     = EMMAS_ESTABLISH_CNF;
  emm_sap.u.emm_as.u.establish.ue_id    = ue_id;
  emm_sap.u.emm_as.u.establish.nas_info = EMM_AS_NAS_INFO_NONE;
  emm_sap.u.emm_as.u.establish.encryption =
      emm_ctx->_security.selected_algorithms.encryption;
  emm_sap.u.emm_as.u.establish.integrity =
      emm_ctx->_security.selected_algorithms.integrity;
  emm_sap.u.emm_as.u.establish.nas_msg       = NULL;
  emm_sap.u.emm_as.u.establish.eps_id.guti   = &emm_ctx->_guti;
  emm_sap.u.emm_as.u.establish.csfb_response = msg->csfbresponse;
  emm_sap.u.emm_as.u.establish.presencemask |= SERVICE_TYPE_PRESENT;
  emm_sap.u.emm_as.u.establish.service_type = msg->servicetype;
  rc                                        = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

static int check_paging_received_without_lai(mme_ue_s1ap_id_t ue_id) {
  ue_mm_context_t* ue_context = NULL;
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_context = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (ue_context) {
    if ((ue_context->sgs_context) &&
        (ue_context->sgs_context->csfb_service_type ==
         CSFB_SERVICE_MT_CALL_OR_SMS_WITHOUT_LAI)) {
      ue_context->sgs_context->csfb_service_type = CSFB_SERVICE_NONE;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, false);
}

int emm_send_service_reject_in_dl_nas(
    const mme_ue_s1ap_id_t ue_id, const uint8_t emm_cause) {
  int rc                 = RETURNok;
  emm_sap_t emm_sap      = {0};
  emm_context_t* emm_ctx = emm_context_get(&_emm_data, ue_id);
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (!emm_ctx) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "Failed to find emm context for ue_id :" MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  emm_ctx->emm_cause                = emm_cause;
  emm_sap.primitive                 = EMMAS_DATA_REQ;
  emm_sap.u.emm_as.u.data.emm_cause = (uint32_t*) &emm_ctx->emm_cause;
  emm_sap.u.emm_as.u.data.ue_id     = ue_id;
  emm_sap.u.emm_as.u.data.nas_info  = EMM_AS_NAS_DATA_INFO_SR;
  emm_sap.u.emm_as.u.data.nas_msg   = NULL;  // No ESM container
  /*
   * Setup EPS NAS security data
   */
  emm_as_set_security_data(
      &emm_sap.u.emm_as.u.data.sctx, &emm_ctx->_security, false, true);

  rc = emm_sap_send(&emm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
