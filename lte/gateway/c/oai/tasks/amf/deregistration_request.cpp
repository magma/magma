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
/*****************************************************************************
  Source      deregistration_request.cpp
  Version     0.1
  Date        2020/12/01
  Product     AMF Core
  Subsystem   Access and Mobility Management Function
  Author      Sanjay Kumar Ojha
  Description Defination of deRegistration requesr response procedure
*****************************************************************************/
#ifndef DEREGISTRATION_REQUEST_SEEN
#define DEREGISTRATION_REQUEST_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "3gpp_38.401.h"
//#include "bstrlib.h"
#ifdef __cplusplus
};
#endif
//#include "deregistration_request.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_app_defs.h"
#include "amf_recv.h"
#include "amf_asDefs.h"
#include "amf_as.h"
#include "amf_sap.h"
#include "amf_app_state_manager.h"

using namespace std;

namespace magma5g {
amf_as_data_t amf_data_de_reg_sec;
amf_sap_c amf_sap_de_reg;
amf_app_msg app_msg_de_reg;
/*
 * name : amf_handle_deregistration_ue_origin_req()
 * Description: Starts processing de-registration request from UE.
 *        Request comes from AS to AMF as UL NAS message.
 *        Current scope is 3GPP connection, irrespective of
 *        switch-off or normal de-registration.
 *        re-registration required is out of mvc scope now.
 */
int amf_procedure_handler::amf_handle_deregistration_ue_origin_req(
    amf_ue_ngap_id_t ue_id, DeRegistrationRequestUEInitMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF-TEST: UE originated Deregistration procedures started\n");
  int rc = RETURNerror;

  amf_deregistration_request_ies_t params;

  if (msg->m5gs_de_reg_type.switchoff) {
    params.de_reg_type = AMF_SWITCHOFF_DEREGISTRATION;
  } else {
    params.de_reg_type = AMF_NORMAL_DEREGISTRATION;
  }
  /*value of access_type would be 1 or 2 or 3, 24-501 - 9.11.3.20 */
  switch (msg->m5gs_de_reg_type.access_type) {
    case 1:
      params.de_reg_access_type = AMF_3GPP_ACCESS;
      OAILOG_INFO(
          LOG_NAS_AMF,
          "AMF_TEST: "
          "Access type is AMF_3GPP_ACCESS for deregistration request from "
          "UE\n");
      break;
    case 2:
      params.de_reg_access_type = AMF_NONE_3GPP_ACCESS;
      OAILOG_INFO(
          LOG_NAS_AMF,
          "Access type AMF_NONE_3GPP_ACCESS for deregistration request from "
          "UE\n");
      break;
    case 3:
      params.de_reg_access_type = AMF_3GPP_ACCESS_AND_NONE_3GPP_ACCESS;
      OAILOG_INFO(
          LOG_NAS_AMF,
          "Access type AMF_3GPP_ACCESS_AND_NONE_3GPP_ACCESS for deregistration "
          "request from UE\n");
      break;
    default:
      /*Wrong access-type and nothing to do and return with error*/
      OAILOG_INFO(
          LOG_NAS_AMF, "Wrong access type received for deregistration\n");
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
      /*TODO handle error case rather returning from here*/
      break;
  }
  /*setting key set identifier as received from UE*/
  params.ksi = msg->nas_key_set_identifier.nas_key_set_identifier;
  // TODO need to recheck on this counter arguments
  increment_counter("ue_deregistration", 1, 1, "amf_cause", "ue_initiated");
  rc = amf_proc_deregistration_request(ue_id, &params);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/*
 * Add explanation TODO
 * Function : amf_proc_deregistration_request
 *
 * Logic: Because it's UE originated, it has to be deleted from AMF
 *        and other related NFs.
 *        If switch off, no need to send deregistration accept in DL
 *        Normal deregistration on 3GPP,
 *        Send accept message in DL to UE.
 *        Update UE state
 *        Sends ue context release notification to NGAP
 *        Clean up availble ue context
 *        send message to SMF to remove all PDU sessions
 *        AMF is not expecting response from SMF.
 * Note: Not considered below messages/elements,
 * Deregistration type as identity type, i.e. IMSI/TIMSI/GUTI
 * teid, pdn_context, cm_state as IDLE, purge_request, UDM/network
 * initiated deregistration, timer assosiation release
 */

int amf_proc_deregistration_request(
    amf_ue_ngap_id_t ue_id, amf_deregistration_request_ies_t* params) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_INFO(
      LOG_NAS_AMF,
      "AMF_TEST: processing deregistration UE-id = %d "
      "type = %d\n",
      ue_id, params->de_reg_type);
  int rc                 = RETURNerror;
  amf_context_t* amf_ctx = amf_context_get(ue_id);
  if (!amf_ctx) {
    OAILOG_INFO(
        LOG_NAS_AMF,
        "AMF icontext not present for UE-id = %d "
        "type = %d\n",
        ue_id, params->de_reg_type);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  amf_sap_t amf_sap;
  amf_as_data_t* amf_as = &amf_sap.u.amf_as.u.data;

  /* if switched off, directly release all resources and
   * dont send accept to UE
   * Note: TODO deregister_success_event TBD
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
    rc                = amf_sap_de_reg.amf_sap_send(&amf_sap);
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
    rc = amf_sap_de_reg.amf_sap_send(&amf_sap);
    /* Handle releasing all context related resources
     */
    rc = amf_app_handle_deregistration_req(ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/*
 * TODO add explanation
 * */
int amf_app_handle_deregistration_req(amf_ue_ngap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  // OAILOG_INFO(
  //   LOG_NAS_AMF, "AMF amf_app_handle_deregistration_req UE-id = %d \n",
  //   ue_id);
  int rc                         = RETURNerror;
  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_context) {
    // OAILOG_ERROR(
    //     LOG_AMF_APP, "UE context doesn't exist -> Nothing to do. Return \n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }
  // TODO Condition for teid of N11, active_pdn_contexts and amy more
  // will be taken care later as PDU session related info not stored
  // but proceeding to release all the resources and notify NGAP
  amf_app_desc_t* amf_app_desc_p = get_amf_nas_state(false);
  if (!amf_app_desc_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Failed to fetch amf_app_desc_p \n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }
  // UE context release notification to NGAP
  app_msg_de_reg.amf_app_ue_context_release(
      ue_context, ue_context->ue_context_rel_cause);

  // Remove stored UE context from AMF core.
  amf_remove_ue_context(&amf_app_desc_p->amf_ue_contexts, ue_context);

  // At the end remove PDU session related info and send all session release
  // message to SMF
  if (ue_context->amf_teid_n11 && (M5GCM_IDLE == ue_context->ecm_state)) {
    // check for "ue_context_p->pdn_contexts[i]" as well later TODO
    // amf_app_send_delete_session_request(ue_context,
    //		    ue_context->pdn_contexts[i]->default_ebi, i);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/*
 * Add logic description
 */
void amf_remove_ue_context(
    amf_ue_context_t* amf_ue_context_p, ue_m5gmm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_INFO(LOG_NAS_AMF, "AMF UE context remove from amf core\n");

  hashtable_rc_t hash_rc              = HASH_TABLE_OK;
  hash_table_ts_t* amf_state_ue_id_ht = get_amf_ue_state();

  if (!amf_ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid AMF UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  if (!ue_context_p) {
    OAILOG_ERROR(LOG_AMF_APP, "Invalid UE context received\n");
    OAILOG_FUNC_OUT(LOG_AMF_APP);
  }
  /* UE state had not been mainted & nothing to remove
   * clear_amf_ctxt() will be implemeted later
   * Timer had not been implemented, nothing to be removed
   * 4 hash table elements to be removed.
   */
  // for imsi_amf_ue_id_htbl removal
  if (ue_context_p->amf_context._imsi64) {
    hash_rc = hashtable_uint64_ts_remove(
        amf_ue_context_p->imsi_amf_ue_id_htbl,
        (const hash_key_t) ue_context_p->amf_context._imsi64);
    if (HASH_TABLE_OK != hash_rc) {
      OAILOG_ERROR_UE(
          LOG_AMF_APP, ue_context_p->amf_context._imsi64,
          "UE context not found for IMSI\n"
          " gnb_ue_ngap_id"
          "0x%08" PRIX64
          " amf_ue_ngap_id "
          "0x%08" PRIX64 " not in IMSI collection\n",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
    }
  }
  // tun11_ue_context_htbl removal
  if (ue_context_p->amf_teid_n11) {
    hash_rc = hashtable_uint64_ts_remove(
        amf_ue_context_p->tun11_ue_context_htbl,
        (const hash_key_t) ue_context_p->amf_teid_n11);
    if (HASH_TABLE_OK != hash_rc)
      OAILOG_ERROR_UE(
          LOG_AMF_APP, ue_context_p->amf_context._imsi64,
          "UE context not found for n11_ue_context\n"
          " gnb_ue_ngap_id "
          "0x%08" PRIX64
          " amf_ue_ngap_id "
          "0x%08" PRIX64 " \n",
          ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
  }
  // for gnb_ue_ngap_id_ue_context_htbl
  hash_rc = hashtable_uint64_ts_remove(
      amf_ue_context_p->gnb_ue_ngap_id_ue_context_htbl,
      (const hash_key_t) ue_context_p->gnb_ngap_id_key);
  OAILOG_ERROR_UE(
      LOG_AMF_APP, ue_context_p->amf_context._imsi64,
      "UE context not found for gnb_ue_ngap_id_ue_context_htbl\n"
      " gnb_ue_ngap_id "
      "0x%08" PRIX64
      " amf_ue_ngap_id "
      "0x%08" PRIX64 " \n",
      ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);

  // for guti_ue_context_htbl obj hash table
  hash_rc = obj_hashtable_uint64_ts_remove(
      amf_ue_context_p->guti_ue_context_htbl,
      (const void* const) & ue_context_p->amf_context._m5_guti,
      sizeof(ue_context_p->amf_context._m5_guti));
  OAILOG_ERROR_UE(
      LOG_AMF_APP, ue_context_p->amf_context._imsi64,
      "UE context not found for guti_ue_context_htbl\n"
      " gnb_ue_ngap_id "
      "0x%08" PRIX64
      " amf_ue_ngap_id "
      "0x%08" PRIX64 " \n",
      ue_context_p->gnb_ue_ngap_id, ue_context_p->amf_ue_ngap_id);
}

}  // end  namespace magma5g
#endif
