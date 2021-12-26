/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#ifndef DEREGISTRATION_REQUEST_SEEN
#define DEREGISTRATION_REQUEST_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#ifdef __cplusplus
};
#endif
#include <unordered_map>
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"

using magma5g::AsyncM5GMobilityServiceClient;

namespace magma5g {
amf_as_data_t amf_data_de_reg_sec;
extern std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*> ue_context_map;

/*
 * name : amf_handle_deregistration_ue_origin_req()
 * Description: Starts processing de-registration request from UE.
 *        Request comes from AS to AMF as UL NAS message.
 *        Current scope is 3GPP connection, irrespective of
 *        switch-off or normal de-registration.
 *        re-registration required is out of mvc scope now.
 */
int amf_handle_deregistration_ue_origin_req(
    amf_ue_ngap_id_t ue_id, DeRegistrationRequestUEInitMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(
      LOG_NAS_AMF, "UE originated deregistration procedures started\n");
  int rc = RETURNerror;
  amf_deregistration_request_ies_t params;
  if (msg->m5gs_de_reg_type.switchoff) {
    params.de_reg_type = AMF_SWITCHOFF_DEREGISTRATION;
  } else {
    params.de_reg_type = AMF_NORMAL_DEREGISTRATION;
  }
  /*value of access_type would be 1 or 2 or 3, 24-501 - 9.11.3.20 */
  switch (msg->m5gs_de_reg_type.access_type) {
    case AMF_3GPP_ACCESS:
      params.de_reg_access_type = AMF_3GPP_ACCESS;
      OAILOG_DEBUG(
          LOG_NAS_AMF,
          "Access type is AMF_3GPP_ACCESS for deregistration request from "
          "UE\n");
      break;
    case NON_AMF_3GPP_ACCESS:
      params.de_reg_access_type = AMF_NONE_3GPP_ACCESS;
      OAILOG_DEBUG(
          LOG_NAS_AMF,
          "Access type AMF_NONE_3GPP_ACCESS for deregistration request from "
          "UE\n");
      break;
    case AMF_3GPP_ACCESS_AND_NONE_3GPP_ACCESS:
      params.de_reg_access_type = AMF_3GPP_ACCESS_AND_NONE_3GPP_ACCESS;
      OAILOG_DEBUG(
          LOG_NAS_AMF,
          "Access type AMF_3GPP_ACCESS_AND_NONE_3GPP_ACCESS for deregistration "
          "request from UE\n");
      break;
    default:
      OAILOG_WARNING(
          LOG_NAS_AMF, "Wrong access type received for deregistration\n");
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
      break;
  }
  /*setting key set identifier as received from UE*/
  params.ksi = msg->nas_key_set_identifier.nas_key_set_identifier;
  increment_counter("ue_deregistration", 1, 1, "amf_cause", "ue_initiated");
  rc = amf_proc_deregistration_request(ue_id, &params);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/*
 * Function : amf_proc_deregistration_request
 *
 * Description : Process the UE originated De-Registration request
 */
int amf_proc_deregistration_request(
    amf_ue_ngap_id_t ue_id, amf_deregistration_request_ies_t* params) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(
      LOG_NAS_AMF,
      "Processing deregistration UE-id = " AMF_UE_NGAP_ID_FMT " type = %d",
      ue_id, params->de_reg_type);
  int rc = RETURNerror;

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context == NULL) {
    return -1;
  }

  amf_context_t* amf_ctx = amf_context_get(ue_id);
  if (!amf_ctx) {
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "AMF icontext not present for UE-id = " AMF_UE_NGAP_ID_FMT
        " type = %d\n",
        ue_id, params->de_reg_type);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  amf_sap_t amf_sap     = {};
  amf_as_data_t* amf_as = &amf_sap.u.amf_as.u.data;

  /* if switched off, directly release all resources and
   * dont send accept to UE
   */
  if (params->de_reg_type == AMF_SWITCHOFF_DEREGISTRATION) {
    increment_counter("ue_deregister", 1, 1, "result", "success");
    increment_counter(
        "ue_deregister", 1, 1, "action", "deregistration_accept_not_sent");
    rc = RETURNok;
  } else {
    /* AMF_NORMAL_DEREGISTRATION case where 3GPP getting deregistered
     * first send accept message and then release respective
     * resources
     */
    amf_as->ue_id    = ue_id;
    amf_as->nas_info = AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT;
    amf_as->nas_msg  = {0};
    /*setup NAS sequrity data to send accept message in DL req*/
    amf_data_de_reg_sec.amf_as_set_security_data(
        &amf_as->sctx, &amf_ctx->_security, false, true);
    /*
     * Send AMF-AS SAP Deregistration Accept message to NGAP
     * on AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT
     */
    amf_sap.primitive = AMFAS_DATA_REQ;
    rc                = amf_sap_send(&amf_sap);
    increment_counter("ue_deregister", 1, 1, "result", "success");
    increment_counter(
        "ue_deregister", 1, 1, "action", "deregister_accept_sent");
  }
  /* start releasing UE related context and hash tables*/
  if (rc != RETURNerror) {
    amf_as->ue_id           = ue_id;
    amf_sap.primitive       = AMFREG_DEREGISTRATION_REQ;
    amf_sap.u.amf_reg.ue_id = ue_id;
    amf_sap.u.amf_reg.ctx   = amf_ctx;
    /* send to update respective state UE machine*/
    rc = amf_sap_send(&amf_sap);
    /* Handle releasing all context related resources
     */

    ue_context->ue_context_rel_cause.ngapCause_u.nas = ngap_CauseNas_deregister;
    rc = ue_state_handle_message_dereg(
        ue_context->mm_state, STATE_EVENT_DEREGISTER, SESSION_NULL, ue_context,
        ue_id);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/***************************************************************************
**                                                                        **
** Name:    amf_app_handle_deregistration_req()                           **
**                                                                        **
** Description: Processes Deregistration Request                          **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_app_handle_deregistration_req(amf_ue_ngap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc                         = RETURNerror;
  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_context) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "ue context not found for the "
        "ue_id = " AMF_UE_NGAP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }
  // TODO: will be taken care later as PDU session related info not stored
  // but proceeding to release all the resources and notify NGAP
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (!amf_app_desc_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Failed to fetch amf_app_desc_p \n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }
  // UE context release notification to NGAP
  if (ue_context->ue_context_rel_cause.ngapCause_u.nas ==
      ngap_CauseNas_normal_release) {
    ue_context->ue_context_rel_cause.ngapCause_u.nas = ngap_CauseNas_deregister;
  }
  amf_app_ue_context_release(ue_context, ue_context->ue_context_rel_cause);

  // Clean up all the sessions.
  amf_smf_context_cleanup_pdu_session(ue_context);

  amf_free_ue_context(ue_context);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/***************************************************************************
**                                                                        **
** Name:   amf_smf_context_cleanup_pdu_session()                          **
**                                                                        **
** Description: Function to remove UE Context                             **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_smf_context_cleanup_pdu_session(ue_m5gmm_context_s* ue_context) {
  amf_smf_release_t smf_message;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];

  memset(&smf_message, 0, sizeof(amf_smf_release_t));

  for (auto& it : ue_context->amf_context.smf_ctxt_map) {
    IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);

    std::shared_ptr<smf_context_t> i = it.second;
    smf_message.pdu_session_id =
        i->smf_proc_data.pdu_session_identity.pdu_session_id;

    smf_message.pti = i->smf_proc_data.pti.pti;

    release_session_gprc_req(&smf_message, imsi);

    if (i->pdu_address.pdn_type == IPv4) {
      AsyncM5GMobilityServiceClient::getInstance().release_ipv4_address(
          imsi, i->dnn.c_str(), &(i->pdu_address.ipv4_address));
    }
  }
}

/***************************************************************************
**                                                                        **
** Name:    amf_remove_ue_context()                                       **
**                                                                        **
** Description: Function to remove UE Context                             **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_remove_ue_context(ue_m5gmm_context_s* ue_context_p) {
  std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*>::iterator
      found_ue_id = ue_context_map.find(ue_context_p->amf_ue_ngap_id);

  if (found_ue_id != ue_context_map.end()) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "Removed ue id = " AMF_UE_NGAP_ID_FMT
        " entry from the ue context map\n",
        ue_context_p->amf_ue_ngap_id);
    ue_context_map.erase(found_ue_id);
  }
}
}  // end  namespace magma5g
#endif
