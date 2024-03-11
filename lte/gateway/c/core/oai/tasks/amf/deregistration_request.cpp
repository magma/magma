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
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_asDefs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_session_context.hpp"
#include "orc8r/gateway/c/common/service303/MetricsHelpers.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_app_statistics.hpp"

namespace magma5g {
amf_as_data_t amf_data_de_reg_sec;

/*
 * name : amf_handle_deregistration_ue_origin_req()
 * Description: Starts processing de-registration request from UE.
 *        Request comes from AS to AMF as UL NAS message.
 *        Current scope is 3GPP connection, irrespective of
 *        switch-off or normal de-registration.
 *        re-registration required is out of mvc scope now.
 */
status_code_e amf_handle_deregistration_ue_origin_req(
    amf_ue_ngap_id_t ue_id, DeRegistrationRequestUEInitMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(LOG_NAS_AMF,
               "UE originated deregistration procedures started\n");
  status_code_e rc = RETURNerror;
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
      OAILOG_WARNING(LOG_NAS_AMF,
                     "Wrong access type received for deregistration\n");
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
status_code_e amf_proc_deregistration_request(
    amf_ue_ngap_id_t ue_id, amf_deregistration_request_ies_t* params) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(LOG_NAS_AMF,
               "Processing deregistration UE-id = " AMF_UE_NGAP_ID_FMT
               " type = %d",
               ue_id, params->de_reg_type);
  status_code_e rc = RETURNerror;

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context == NULL) {
    return RETURNerror;
  }

  amf_context_t* amf_ctx = amf_context_get(ue_id);
  if (!amf_ctx) {
    OAILOG_DEBUG(LOG_NAS_AMF,
                 "AMF icontext not present for UE-id = " AMF_UE_NGAP_ID_FMT
                 " type = %d\n",
                 ue_id, params->de_reg_type);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  amf_sap_t amf_sap = {};
  amf_as_data_t* amf_as = &amf_sap.u.amf_as.u.data;

  /* if switched off, directly release all resources and
   * dont send accept to UE
   */
  if (params->de_reg_type == AMF_SWITCHOFF_DEREGISTRATION) {
    increment_counter("ue_deregister", 1, 1, "result", "success");
    increment_counter("ue_deregister", 1, 1, "action",
                      "deregistration_accept_not_sent");
    rc = RETURNok;
  } else {
    /* AMF_NORMAL_DEREGISTRATION case where 3GPP getting deregistered
     * first send accept message and then release respective
     * resources
     */
    amf_as->ue_id = ue_id;
    amf_as->nas_info = AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT;
    amf_as->nas_msg = {0};
    /*setup NAS sequrity data to send accept message in DL req*/
    amf_data_de_reg_sec.amf_as_set_security_data(
        &amf_as->sctx, &amf_ctx->_security, false, true);
    /*
     * Send AMF-AS SAP Deregistration Accept message to NGAP
     * on AMF_AS_NAS_DATA_DEREGISTRATION_ACCEPT
     */
    amf_sap.primitive = AMFAS_DATA_REQ;
    rc = amf_sap_send(&amf_sap);
    increment_counter("ue_deregister", 1, 1, "result", "success");
    increment_counter("ue_deregister", 1, 1, "action",
                      "deregister_accept_sent");
  }
  /* start releasing UE related context and hash tables*/
  if (rc != RETURNerror) {
    amf_as->ue_id = ue_id;
    amf_sap.primitive = AMFREG_DEREGISTRATION_REQ;
    amf_sap.u.amf_reg.ue_id = ue_id;
    amf_sap.u.amf_reg.ctx = amf_ctx;
    /* send to update respective state UE machine*/
    rc = amf_sap_send(&amf_sap);
    /* Handle releasing all context related resources
     */

    ue_context->ue_context_rel_cause = NGAP_NAS_DEREGISTER;
    rc = ue_state_handle_message_dereg(ue_context->mm_state,
                                       STATE_EVENT_DEREGISTER, SESSION_NULL,
                                       ue_context, ue_id);
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
status_code_e amf_app_handle_deregistration_req(amf_ue_ngap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  status_code_e rc = RETURNerror;
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  ue_m5gmm_context_s* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP,
                 "ue context not found for the "
                 "ue_id = " AMF_UE_NGAP_ID_FMT "\n",
                 ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }

  // Clean up all the sessions.
  amf_smf_context_cleanup_pdu_session(ue_context_p);

  if (ue_context_p->amf_context.new_registration_info) {
    nas_delete_all_amf_procedures(&ue_context_p->amf_context);
    proc_new_registration_req(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  // UE context release notification to NGAP
  ue_context_p->mm_state = DEREGISTERED;

  if (M5GCM_IDLE == ue_context_p->cm_state) {
    ue_context_p->ue_context_rel_cause = NGAP_IMPLICIT_CONTEXT_RELEASE;
    // Notify NGAP to release NGAP UE context locally.
    amf_app_itti_ue_context_release(ue_context_p,
                                    ue_context_p->ue_context_rel_cause);

    amf_remove_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
  } else {
    if (ue_context_p->ue_context_rel_cause == NGAP_INVALID_CAUSE) {
      ue_context_p->ue_context_rel_cause = NGAP_NAS_DEREGISTER;
    }

    // Notify NGAP to send UE Context Release Command to eNB.
    amf_app_itti_ue_context_release(ue_context_p,
                                    ue_context_p->ue_context_rel_cause);
    if (ue_context_p->ue_context_rel_cause == NGAP_SCTP_SHUTDOWN_OR_RESET) {
      amf_remove_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context_p);
    } else {
      ue_context_p->ue_context_rel_cause = NGAP_INVALID_CAUSE;
    }
  }
  update_amf_app_stats_connected_ue_sub();
  update_amf_app_stats_registered_ue_sub();
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
  OAILOG_FUNC_IN(LOG_AMF_APP);
  memset(&smf_message, 0, sizeof(amf_smf_release_t));

  for (auto& it : ue_context->amf_context.smf_ctxt_map) {
    IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);

    std::shared_ptr<smf_context_t> i = it.second;
    smf_message.pdu_session_id = i->smf_proc_data.pdu_session_id;

    smf_message.pti = i->smf_proc_data.pti;

    release_session_gprc_req(&smf_message, imsi);

    if ((i->pdu_address.pdn_type == IPv4) ||
        (i->pdu_address.pdn_type == IPv4_AND_v6)) {
      AMFClientServicer::getInstance().release_ipv4_address(
          imsi, i->dnn.c_str(), &(i->pdu_address.ipv4_address));
    }

    if ((i->pdu_address.pdn_type == IPv6) ||
        (i->pdu_address.pdn_type == IPv4_AND_v6)) {
      AMFClientServicer::getInstance().release_ipv6_address(
          imsi, i->dnn.c_str(), &(i->pdu_address.ipv6_address));
    }

    OAILOG_INFO(LOG_AMF_APP,
                "Deleting Pdu Session id = %d for ue_id = " AMF_UE_NGAP_ID_FMT
                "\n",
                smf_message.pdu_session_id, ue_context->amf_ue_ngap_id);
    update_amf_app_stats_pdusessions_ue_sub();
  }

  ue_context->amf_context.smf_ctxt_map.clear();
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

//------------------------------------------------------------------------------
void amf_app_ue_context_free_content(ue_m5gmm_context_s* const ue_context_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  // Stop Mobile reachability timer,if running

  // Stop Implicit deregistration timer,if running

  // Stop Initial context setup process guard timer,if running

  ue_context_p->ue_context_rel_cause = NGAP_INVALID_CAUSE;

  amf_smf_context_cleanup_pdu_session(ue_context_p);
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}

void clear_amf_ctxt(amf_context_t* amf_context) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  if (!amf_context) {
    return;
  }

  nas_delete_all_amf_procedures(amf_context);

  ue_m5gmm_context_s* ue_context_p =
      PARENT_STRUCT(amf_context, ue_m5gmm_context_s, amf_context);

  ue_context_p->mm_state = DEREGISTERED;
  amf_ctx_clear_auth_vectors(amf_context);
  OAILOG_FUNC_OUT(LOG_NAS_AMF);
}

/***************************************************************************
**                                                                        **
** Name:    amf_remove_ue_context()                                       **
**                                                                        **
** Description: Function to remove UE Context                             **
**                                                                        **
**                                                                        **
***************************************************************************/
void amf_remove_ue_context(amf_ue_context_t* const amf_ue_context_p,
                           ue_m5gmm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_AMF_APP);

  magma::map_rc_t m_rc = magma::MAP_OK;
  map_uint64_ue_context_t* amf_state_ue_id_ht = get_amf_ue_state();

  if (!amf_ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }

  // TODO: Need clean up in redis database

  delete_amf_ue_state(ue_context_p->amf_context.imsi64);
  amf_app_ue_context_free_content(ue_context_p);

  // IMSI
  if (ue_context_p->amf_context.imsi64) {
    m_rc = amf_ue_context_p->imsi_amf_ue_id_htbl.remove(
        ue_context_p->amf_context.imsi64);

    if (m_rc != magma::MAP_OK) {
      OAILOG_ERROR_UE(
          LOG_AMF_APP, ue_context_p->amf_context.imsi64,
          "UE context not found!\n"
          " gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT
          " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT " not in IMSI collection\n",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
    }
  }

  // filled guti
  if ((ue_context_p->amf_context.m5_guti.guamfi.amf_regionid) ||
      (ue_context_p->amf_context.m5_guti.guamfi.amf_set_id) ||
      (ue_context_p->amf_context.m5_guti.guamfi.amf_pointer) ||
      (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit1) ||
      (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit2) ||
      (ue_context_p->amf_context.m5_guti.guamfi.plmn.mcc_digit3)) {
    m_rc = amf_ue_context_p->guti_ue_context_htbl.remove(
        ue_context_p->amf_context.m5_guti);
    if (m_rc != magma::MAP_OK)
      OAILOG_ERROR(LOG_AMF_APP,
                   "UE Context not found!\n"
                   " gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT
                   " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT
                   ", GUTI  not in GUTI collection\n",
                   ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
  }

  clear_amf_ctxt(&ue_context_p->amf_context);

  // gNB UE NGAP UE ID
  m_rc = amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl.remove(
      ue_context_p->gnb_ngap_id_key);
  if (m_rc != magma::MAP_OK)
    OAILOG_ERROR(LOG_AMF_APP,
                 "UE context not found!\n"
                 " gnb_ue_ngap_ue_id " GNB_UE_NGAP_ID_FMT
                 " amf_ue_ngap_id " AMF_UE_NGAP_ID_FMT,
                 ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);

  m_rc = amf_ue_context_p->tun11_ue_context_htbl.remove(
      ue_context_p->amf_teid_n11);

  // filled NAS UE ID/ MME UE S1AP ID
  if (ue_context_p->amf_ue_ngap_id != INVALID_AMF_UE_NGAP_ID) {
    m_rc = amf_state_ue_id_ht->remove(ue_context_p->amf_ue_ngap_id);
    if (m_rc != magma::MAP_OK)
      OAILOG_TRACE(LOG_AMF_APP, "Error Could not remove this ue context \n");
    ue_context_p->amf_ue_ngap_id = INVALID_AMF_UE_NGAP_ID;
  }

  delete ue_context_p;
  ue_context_p = NULL;
  OAILOG_FUNC_OUT(LOG_AMF_APP);
}
}  // end  namespace magma5g
#endif
