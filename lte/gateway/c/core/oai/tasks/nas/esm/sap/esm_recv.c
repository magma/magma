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

#include <stdbool.h>
#include <stdlib.h>

#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.007.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/sap/esm_recv.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_pt.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_ebr.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/msg/esm_cause.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/NasRequestType.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/PdnType.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.h"
#include "lte/gateway/c/core/oai/include/mme_app_desc.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_apn_selection.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_itti_messaging.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/
extern int send_modify_bearer_req(mme_ue_s1ap_id_t ue_id, ebi_t ebi);

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Functions executed by both the UE and the MME upon receiving ESM messages
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_status()                                         **
 **                                                                        **
 ** Description: Processes ESM status message                              **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/

esm_cause_t esm_recv_status(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const esm_status_msg* msg) {
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  int rc                = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_ESM);
  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received ESM status message (pti=%d, ebi=%d)\n", pti, ebi);
  /*
   * Message processing
   */
  /*
   * Get the ESM cause
   */
  esm_cause = msg->esmcause;
  /*
   * Execute the ESM status procedure
   */
  rc = esm_proc_status_ind(emm_context, pti, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/*
   --------------------------------------------------------------------------
   Functions executed by the MME upon receiving ESM message from the UE
   --------------------------------------------------------------------------
*/
/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_pdn_connectivity_request()                       **
 **                                                                        **
 ** Description: Processes PDN connectivity request message                **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     new_ebi:   New assigned EPS bearer identity           **
 **      data:      PDN connection and EPS bearer context data **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_pdn_connectivity_request(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const pdn_connectivity_request_msg* msg, ebi_t* new_ebi,
    bool is_standalone) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int rc            = RETURNerror;
  int esm_cause     = ESM_CAUSE_SUCCESS;
  pdn_cid_t pdn_cid = 0;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received PDN Connectivity Request message "
      "(ue_id= " MME_UE_S1AP_ID_FMT ", pti=%u, ebi=%u)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if ((pti == ESM_PT_UNASSIGNED) || esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case a
     * * * * Reserved or unassigned PTI value
     */
    OAILOG_ERROR_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d) for (ue_id "
        "= " MME_UE_S1AP_ID_FMT ") \n",
        pti, ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case a
     * * * * Reserved or assigned EPS bearer identity value
     */
    OAILOG_ERROR_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d) for (ue_id "
        "= " MME_UE_S1AP_ID_FMT ")\n",
        ebi, ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Get PDN connection and EPS bearer context data structure to setup
   */
  if (!emm_context->esm_ctx.esm_proc_data) {
    emm_context->esm_ctx.esm_proc_data = (esm_proc_data_t*) calloc(
        1, sizeof(*emm_context->esm_ctx.esm_proc_data));
  }

  struct esm_proc_data_s* esm_data = emm_context->esm_ctx.esm_proc_data;

  esm_data->pti = pti;
  /*
   * Get the PDN connectivity request type
   */
  OAILOG_DEBUG_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - PDN Connectivity Request Type = (%d) for (ue_id = %u)\n ",
      msg->requesttype, ue_id);

  if (msg->requesttype == REQUEST_TYPE_INITIAL_REQUEST) {
    esm_data->request_type = ESM_PDN_REQUEST_INITIAL;
  } else if (msg->requesttype == REQUEST_TYPE_HANDOVER) {
    esm_data->request_type = ESM_PDN_REQUEST_HANDOVER;
  } else if (msg->requesttype == REQUEST_TYPE_EMERGENCY) {
    esm_data->request_type = ESM_PDN_REQUEST_EMERGENCY;
  } else {
    /*
     * Unkown PDN request type
     */
    esm_data->request_type = -1;
    OAILOG_ERROR_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PDN request type (INITIAL/HANDOVER/EMERGENCY) for "
        "(ue_id = " MME_UE_S1AP_ID_FMT ")\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_MANDATORY_INFO);
  }

  /*
   * Get the value of the PDN type indicator
   */
  OAILOG_DEBUG_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - PDN Type = (%d) for (ue_id = %u)\n ", msg->pdntype, ue_id);
  if (msg->pdntype == PDN_TYPE_IPV4) {
    esm_data->pdn_type = ESM_PDN_TYPE_IPV4;
  } else if (msg->pdntype == PDN_TYPE_IPV6) {
    esm_data->pdn_type = ESM_PDN_TYPE_IPV6;
  } else if (msg->pdntype == PDN_TYPE_IPV4V6) {
    esm_data->pdn_type = ESM_PDN_TYPE_IPV4V6;
  } else {
    /*
     * Unkown PDN type
     */
    OAILOG_ERROR_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PDN type for (ue_id = " MME_UE_S1AP_ID_FMT ")\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_UNKNOWN_PDN_TYPE);
  }

  /*
   * Get the Access Point Name, if provided
   */
  if (msg->presencemask & PDN_CONNECTIVITY_REQUEST_ACCESS_POINT_NAME_PRESENT) {
    if (esm_data->apn) bdestroy_wrapper(&esm_data->apn);
    if (mme_config.nas_config.enable_apn_correction) {
      esm_data->apn = mme_app_process_apn_correction(
          &(emm_context->_imsi), msg->accesspointname);
      OAILOG_INFO_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "ESM-SAP   - APN CORRECTION (apn = %s) for ue id " MME_UE_S1AP_ID_FMT
          "\n",
          (const char*) bdata(esm_data->apn), ue_id);
    } else {
      esm_data->apn = msg->accesspointname;
    }
  }

  if (msg->presencemask &
      PDN_CONNECTIVITY_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT) {
    if (esm_data->pco.num_protocol_or_container_id)
      clear_protocol_configuration_options(&esm_data->pco);
    copy_protocol_configuration_options(
        &esm_data->pco, &msg->protocolconfigurationoptions);
  }
  /*
   * Get the ESM information transfer flag
   */
  if (msg->presencemask &
      PDN_CONNECTIVITY_REQUEST_ESM_INFORMATION_TRANSFER_FLAG_PRESENT) {
    /*
     * 3GPP TS 24.301, sections 6.5.1.2, 6.5.1.3
     * * * * ESM information, i.e. protocol configuration options, APN, or both,
     * * * * has to be sent after the NAS signalling security has been activated
     * * * * between the UE and the MME.
     * * * *
     * * * * The MME then at a later stage in the PDN connectivity procedure
     * * * * initiates the ESM information request procedure in which the UE
     * * * * can provide the MME with protocol configuration options or APN
     * * * * or both.
     * * * * The MME waits for completion of the ESM information request
     * * * * procedure before proceeding with the PDN connectivity procedure.
     */
    if (!mme_config.nas_config.disable_esm_information) {
      esm_proc_esm_information_request(emm_context, pti);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
    }
  }

  OAILOG_DEBUG_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-PROC  - _esm_data.conf.features %08x, esm pdn type = %d\n",
      _esm_data.conf.features, esm_data->pdn_type);
  emm_context->emm_cause = ESM_CAUSE_SUCCESS;

  if (is_standalone) {
    ue_mm_context_t* ue_mm_context_p =
        mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
    // Select APN
    struct apn_configuration_s* apn_config =
        mme_app_select_apn(ue_mm_context_p, &esm_cause);

    /*
     * Execute the PDN connectivity procedure requested by the UE
     */
    if (!apn_config) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "ESM-PROC  - Cannot select APN for ue id" MME_UE_S1AP_ID_FMT "\n",
          ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
    }

    /* Check if a session already exists for this APN. Only 1 session
     * is supported per APN
     */
    for (uint8_t itr = 0; itr < MAX_APN_PER_UE; itr++) {
      if (ue_mm_context_p->pdn_contexts[itr]) {
        if (!(strcmp(
                (const char*) ue_mm_context_p->pdn_contexts[itr]
                    ->apn_subscribed->data,
                apn_config->service_selection)) &&
            (ue_mm_context_p->pdn_contexts[itr]->is_active)) {
          OAILOG_FUNC_RETURN(
              LOG_NAS_ESM, ESM_CAUSE_MULTIPLE_PDN_CONNECTIONS_NOT_ALLOWED);
        }
      }
    }

    /* Find a free PDN Connection ID*/
    for (pdn_cid = 0; pdn_cid < MAX_APN_PER_UE; pdn_cid++) {
      if (!ue_mm_context_p->pdn_contexts[pdn_cid]) break;
    }

    if (pdn_cid >= MAX_APN_PER_UE) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "ESM-PROC  - Cannot find free pdn_cid for ue id" MME_UE_S1AP_ID_FMT
          "\n",
          ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INSUFFICIENT_RESOURCES);
    }
    // Update pdn connection id
    emm_context->esm_ctx.esm_proc_data->pdn_cid = pdn_cid;

    // Update qci
    esm_data->bearer_qos.qci = apn_config->subscribed_qos.qci;
    rc                       = esm_proc_pdn_connectivity_request(
        emm_context, pti, pdn_cid, apn_config->context_identifier,
        emm_context->esm_ctx.esm_proc_data->request_type, esm_data->apn,
        apn_config->pdn_type, esm_data->pdn_addr, &esm_data->bearer_qos,
        (emm_context->esm_ctx.esm_proc_data->pco.num_protocol_or_container_id) ?
            &emm_context->esm_ctx.esm_proc_data->pco :
            NULL,
        &esm_cause);

    if (rc != RETURNerror) {
      /*
       * Create local default EPS bearer context
       */
      rc = esm_proc_default_eps_bearer_context(
          emm_context, pti, pdn_cid, new_ebi, esm_data->bearer_qos.qci,
          &esm_cause);

      if (rc != RETURNerror) {
        esm_cause = ESM_CAUSE_SUCCESS;
      }
    }
    // Send PDN Connectivity req
    OAILOG_INFO_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-PROC - Sending pdn_connectivity_req to MME APP for "
        "ue " MME_UE_S1AP_ID_FMT "\n",
        ue_id);

    emm_context->esm_ctx.pending_standalone += 1;

    mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
    mme_app_send_s11_create_session_req(
        mme_app_desc_p, ue_mm_context_p, pdn_cid);
  } else {
    mme_app_send_s6a_update_location_req(
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context));
    esm_cause = ESM_CAUSE_SUCCESS;
  }
  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_pdn_disconnect_request()                         **
 **                                                                        **
 ** Description: Processes PDN disconnect request message                  **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **      pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     linked_ebi:    Linked EPS bearer identity of the default  **
 **             bearer associated with the PDN to discon-  **
 **             nect from                                  **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_pdn_disconnect_request(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const pdn_disconnect_request_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  pdn_cid_t pid                    = MAX_APN_PER_UE;
  esm_cause_t esm_cause            = ESM_CAUSE_SUCCESS;
  ue_mm_context_t* ue_mm_context_p = NULL;
  ue_mm_context_p =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);

  if (!ue_mm_context_p) {
    OAILOG_WARNING(
        LOG_NAS_ESM, "Failed to find ue context from emm context \n");
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNerror);
  }

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received PDN Disconnect Request message for "
      "ue_id " MME_UE_S1AP_ID_FMT ", pti=%u, ebi=%u)\n",
      ue_mm_context_p->mme_ue_s1ap_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if ((pti == ESM_PT_UNASSIGNED) || esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case b
     * * * * Reserved or unassigned PTI value
     */
    OAILOG_WARNING(
        LOG_NAS_ESM, "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case b
     * * * * Reserved or assigned EPS bearer identity value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /* Send PDN disconnect reject if there is only one PDN connection*/
  if (ue_mm_context_p->nb_active_pdn_contexts == 1) {
    OAILOG_FUNC_RETURN(
        LOG_NAS_ESM, ESM_CAUSE_LAST_PDN_DISCONNECTION_NOT_ALLOWED);
  }
  /*
   * Message processing
   */
  /*
   * Execute the PDN disconnect procedure requested by the UE
   */
  struct esm_proc_data_s* esm_data = emm_context->esm_ctx.esm_proc_data;
  esm_data->pti                    = pti;

  if (ue_mm_context_p
          ->bearer_contexts[EBI_TO_INDEX(msg->linkedepsbeareridentity)]) {
    pid = ue_mm_context_p
              ->bearer_contexts[EBI_TO_INDEX(msg->linkedepsbeareridentity)]
              ->pdn_cx_id;
    if (pid >= MAX_APN_PER_UE) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "ESM-PROC  - No PDN connection found (lbi=%u) for ue "
          "id " MME_UE_S1AP_ID_FMT "\n",
          msg->linkedepsbeareridentity, ue_mm_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_PROTOCOL_ERROR);
    }

    if (ue_mm_context_p->pdn_contexts[pid] == NULL) {
      OAILOG_ERROR_UE(
          LOG_MME_APP, ue_mm_context_p->emm_context._imsi64,
          "pdn_contexts is NULL for "
          "MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "ebi-%u\n",
          ue_mm_context_p->mme_ue_s1ap_id, ebi);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_PDN_CONNECTION_DOES_NOT_EXIST);
    }

    // Check if the LBI received matches with the default bearer ID
    if (msg->linkedepsbeareridentity !=
        ue_mm_context_p->pdn_contexts[pid]->default_ebi) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "ESM-PROC  - Cannot perform PDN disconnect for dedicated bearer "
          "(lbi=%u) " MME_UE_S1AP_ID_FMT "\n",
          msg->linkedepsbeareridentity, ue_mm_context_p->mme_ue_s1ap_id);

      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
    }
  } else {
    OAILOG_ERROR_UE(
        LOG_NAS_ESM, ue_mm_context_p->emm_context._imsi64,
        "ESM-PROC  - No bearer context found, invalid bearer id (lbi=%u) for "
        "ue id " MME_UE_S1AP_ID_FMT "\n",
        msg->linkedepsbeareridentity, ue_mm_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  bool no_delete_gtpv2c_tunnel = true;  // Due to check on line 470
  OAILOG_INFO_UE(
      LOG_NAS_ESM, ue_mm_context_p->emm_context._imsi64,
      "ESM-SAP   - Sending Delete session req message "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", pid=%d, ebi=%d)\n",
      ue_mm_context_p->mme_ue_s1ap_id, pid, msg->linkedepsbeareridentity);
  mme_app_send_delete_session_request(
      ue_mm_context_p, msg->linkedepsbeareridentity, pid,
      no_delete_gtpv2c_tunnel);

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

//------------------------------------------------------------------------------
esm_cause_t esm_recv_information_response(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const esm_information_response_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received ESM Information response message "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d, ebi=%d)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if ((pti == ESM_PT_UNASSIGNED) || esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case b
     * * * * Reserved or unassigned PTI value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (ebi != ESM_EBI_UNASSIGNED) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case b
     * * * * Reserved or assigned EPS bearer identity value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  bstring apn = msg->accesspointname;
  if (mme_config.nas_config.enable_apn_correction) {
    apn = mme_app_process_apn_correction(
        &(emm_context->_imsi), msg->accesspointname);
    OAILOG_INFO_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - APN CORRECTION (apn = %s) for ue id " MME_UE_S1AP_ID_FMT
        "\n",
        (const char*) bdata(apn), ue_id);
  }

  /*
   * Message processing
   */
  /*
   * Execute the PDN disconnect procedure requested by the UE
   */
  int pid = esm_proc_esm_information_response(
      emm_context, pti, apn, &msg->protocolconfigurationoptions, &esm_cause);

  bdestroy_wrapper((bstring*) &msg->accesspointname);
  if (pid != RETURNerror) {
    // Continue with S6a Update Location Request
    mme_app_send_s6a_update_location_req(
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context));
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    erab_setup_rsp_tmr_exp_handler()                             **
 **                                                                        **
 ** Description: Handles Erab setup rsp timer expiry                       **
 **                                                                        **
 ** Inputs:                                                                **
 **      imsi64:     IMSI                                                  **
 **      args:       timer data                                            **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:  None                                                     **
 **      Others:  None                                                     **
 **                                                                        **
 ***************************************************************************/

status_code_e erab_setup_rsp_tmr_exp_handler(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);

  timer_arg_t timer_args;
  if (!mme_pop_timer_arg(timer_id, &timer_args)) {
    OAILOG_WARNING(
        LOG_NAS_EMM, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }
  mme_ue_s1ap_id_t ue_id = timer_args.ue_id;

  ue_mm_context_t* ue_mm_context = mme_app_get_ue_context_for_timer(
      ue_id, "EPS BEARER DEACTIVATE T3495 Timer");
  if (ue_mm_context == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  ebi_t ebi = timer_args.ebi;
  int rc;
  int bid = EBI_TO_INDEX(ebi);

  bearer_context_t* bearer_context = ue_mm_context->bearer_contexts[bid];
  if (bearer_context == NULL) {
    OAILOG_ERROR_UE(
        LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
        "Bearer context is NULL for (ebi=%u) for ue id " MME_UE_S1AP_ID_FMT
        "\n",
        ebi, ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
  }

  esm_ebr_context_t* ebr_ctx = &(bearer_context->esm_ebr_context);

  if (ebr_ctx && ebr_ctx->args) {
    // Get retransmission timer parameters data
    esm_ebr_timer_data_t* esm_ebr_timer_data =
        (esm_ebr_timer_data_t*) (ebr_ctx->args);
    // Increment the retransmission counter
    esm_ebr_timer_data->count += 1;
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
        "ESM-PROC  - erab_setup_rsp timer expired (ue_id=" MME_UE_S1AP_ID_FMT
        ", ebi=%d), "
        "retransmission counter = %d\n",
        esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi,
        esm_ebr_timer_data->count);

    if (!bearer_context->enb_fteid_s1u.teid) {
      if (esm_ebr_timer_data->count < ERAB_SETUP_RSP_COUNTER_MAX) {
        // Restart the timer
        rc = esm_ebr_start_timer(
            esm_ebr_timer_data->ctx, esm_ebr_timer_data->ebi, NULL,
            1000 * ERAB_SETUP_RSP_TMR, erab_setup_rsp_tmr_exp_handler);
        if (rc != RETURNerror) {
          OAILOG_INFO_UE(
              LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
              "ESM-PROC  - Started ERAB_SETUP_RSP_TMR for "
              "ue_id=" MME_UE_S1AP_ID_FMT
              "ebi (%u)"
              "\n",
              esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi);
        }
      } else {
        OAILOG_WARNING_UE(
            LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
            "ESM-PROC  - ERAB_SETUP_RSP_COUNTER_MAX reached for ERAB_SETUP_RSP "
            "ue_id= " MME_UE_S1AP_ID_FMT
            " ebi (%u)"
            "\n",
            esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi);
        if (bearer_context->esm_ebr_context.timer.id != NAS_TIMER_INACTIVE_ID) {
          bearer_context->esm_ebr_context.timer.id = NAS_TIMER_INACTIVE_ID;
        }
        if (esm_ebr_timer_data) {
          free_wrapper((void**) &esm_ebr_timer_data);
        }
      }
    } else {
      rc = send_modify_bearer_req(
          esm_ebr_timer_data->ue_id, esm_ebr_timer_data->ebi);
      if (rc != RETURNok) {
        OAILOG_ERROR_UE(
            LOG_NAS_ESM, ue_mm_context->emm_context._imsi64,
            "ESM-SAP - Sending Modify bearer req failed for "
            "(ebi=%u)" MME_UE_S1AP_ID_FMT "\n",
            esm_ebr_timer_data->ebi, esm_ebr_timer_data->ue_id);
      }
      if (bearer_context->esm_ebr_context.timer.id != NAS_TIMER_INACTIVE_ID) {
        bearer_context->esm_ebr_context.timer.id = NAS_TIMER_INACTIVE_ID;
      }
      if (esm_ebr_timer_data) {
        free_wrapper((void**) &esm_ebr_timer_data);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, RETURNok);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_default_eps_bearer_context_accept()     **
 **                                                                        **
 ** Description: Processes Activate Default EPS Bearer Context Accept      **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_default_eps_bearer_context_accept(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_default_eps_bearer_context_accept_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  ue_mm_context_t* ue_context_p =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
  mme_ue_s1ap_id_t ue_id = ue_context_p->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received Activate Default EPS Bearer Context "
      "Accept message (ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d, ebi=%d)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
      esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the default EPS bearer context activation procedure accepted
   * * * * by the UE
   */
  int rc =
      esm_proc_default_eps_bearer_context_accept(emm_context, ebi, &esm_cause);

  if (rc != RETURNok) {
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_PROTOCOL_ERROR);
  }
  /* If activate default EPS bearer context accept message is received for a
   * new standalone PDN connection, send modify bearer request to sgw
   */
  if (emm_context->esm_ctx.pending_standalone > 0) {
    emm_context->esm_ctx.pending_standalone -= 1;
    bearer_context_t* bearer_ctx =
        mme_app_get_bearer_context(ue_context_p, ebi);
    if (!bearer_ctx) {
      OAILOG_ERROR_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "Bearer context is NULL for (ebi=%u) for ue id " MME_UE_S1AP_ID_FMT
          "\n",
          ebi, ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
    }
    // set pdn_context to active after receiving "activate default eps bearer
    // context accept from ue
    if (ue_context_p->pdn_contexts[bearer_ctx->pdn_cx_id]) {
      ue_context_p->pdn_contexts[bearer_ctx->pdn_cx_id]->is_active = true;
    }
    /* Send MBR only after receiving ERAB_SETUP_RSP.
     * bearer_ctx->enb_fteid_s1u.teid gets updated after receiving
     * ERAB_SETUP_RSP.*/
    if (bearer_ctx->enb_fteid_s1u.teid) {
      rc = send_modify_bearer_req(ue_id, ebi);
      if (rc != RETURNok) {
        OAILOG_ERROR_UE(
            LOG_NAS_ESM, emm_context->_imsi64,
            "ESM-SAP - Sending Modify bearer req failed for "
            "(ebi=%u)" MME_UE_S1AP_ID_FMT "\n",
            ebi, ue_id);
        OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_PROTOCOL_ERROR);
      }
      OAILOG_DEBUG_UE(
          LOG_NAS_ESM, emm_context->_imsi64,
          "ESM-PROC  - Sending Modify bearer req ue_id=" MME_UE_S1AP_ID_FMT
          "ebi (%u), enb_fteid_s1u.teid %x"
          "\n",
          ue_id, ebi, bearer_ctx->enb_fteid_s1u.teid);

    } else {
      // Wait for ERAB SETUP RSP.Start a timer for 5 secs
      rc = esm_ebr_start_timer(
          emm_context, ebi, NULL, 1000 * ERAB_SETUP_RSP_TMR,
          erab_setup_rsp_tmr_exp_handler);
      if (rc != RETURNerror) {
        OAILOG_DEBUG_UE(
            LOG_NAS_ESM, emm_context->_imsi64,
            "ESM-PROC  - Started ERAB_SETUP_RSP_TMR for "
            "ue_id=" MME_UE_S1AP_ID_FMT
            "ebi (%u)"
            "\n",
            ue_id, ebi);
      }
    }
  }
  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_default_eps_bearer_context_reject()     **
 **                                                                        **
 ** Description: Processes Activate Default EPS Bearer Context Reject      **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fail                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_default_eps_bearer_context_reject(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_default_eps_bearer_context_reject_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_REQUEST_REJECTED_UNSPECIFIED;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received Activate Default EPS Bearer Context "
      "Reject message (ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d, ebi=%d)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
      esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the default EPS bearer context activation procedure not accepted
   * * * * by the UE
   */
  esm_proc_default_eps_bearer_context_reject(emm_context, ebi, &esm_cause);

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_dedicated_eps_bearer_context_accept()   **
 **                                                                        **
 ** Description: Processes Activate Dedicated EPS Bearer Context Accept    **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_dedicated_eps_bearer_context_accept(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_dedicated_eps_bearer_context_accept_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received Activate Dedicated EPS Bearer "
      "Context Accept message (ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d, ebi=%d)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
      esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the dedicated EPS bearer context activation procedure accepted
   * * * * by the UE
   */
  int rc = esm_proc_dedicated_eps_bearer_context_accept(
      emm_context, ebi, &esm_cause);

  if (rc != RETURNerror) {
    esm_cause = ESM_CAUSE_SUCCESS;
  }

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_activate_dedicated_eps_bearer_context_reject()   **
 **                                                                        **
 ** Description: Processes Activate Dedicated EPS Bearer Context Reject    **
 **      message                                                   **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fail                       **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_activate_dedicated_eps_bearer_context_reject(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const activate_dedicated_eps_bearer_context_reject_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received Activate Dedicated EPS Bearer "
      "Context Reject message (ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d, ebi=%d)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
      esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the dedicated EPS bearer context activation procedure not
   * * * *  accepted by the UE
   */
  int rc = esm_proc_dedicated_eps_bearer_context_reject(emm_context, ebi);

  if (rc != RETURNok) {
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_PROTOCOL_ERROR);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_SUCCESS);
}

/****************************************************************************
 **                                                                        **
 ** Name:    esm_recv_deactivate_eps_bearer_context_accept()           **
 **                                                                        **
 ** Description: Processes Deactivate EPS Bearer Context Accept message    **
 **                                                                        **
 ** Inputs:  ue_id:      UE local identifier                        **
 **          pti:       Procedure transaction identity             **
 **      ebi:       EPS bearer identity                        **
 **      msg:       The received ESM message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    ESM cause code whenever the processing of  **
 **             the ESM message fails                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
esm_cause_t esm_recv_deactivate_eps_bearer_context_accept(
    emm_context_t* emm_context, proc_tid_t pti, ebi_t ebi,
    const deactivate_eps_bearer_context_accept_msg* msg) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO_UE(
      LOG_NAS_ESM, emm_context->_imsi64,
      "ESM-SAP   - Received Deactivate EPS Bearer Context "
      "Accept message (ue_id=" MME_UE_S1AP_ID_FMT ", pti=%d, ebi=%d)\n",
      ue_id, pti, ebi);

  /*
   * Procedure transaction identity checking
   */
  if (esm_pt_is_reserved(pti)) {
    /*
     * 3GPP TS 24.301, section 7.3.1, case f
     * * * * Reserved PTI value
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid PTI value (pti=%d)\n", pti);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_PTI_VALUE);
  }
  /*
   * EPS bearer identity checking
   */
  else if (
      esm_ebr_is_reserved(ebi) || esm_ebr_is_not_in_use(emm_context, ebi)) {
    /*
     * 3GPP TS 24.301, section 7.3.2, case f
     * * * * Reserved or assigned value that does not match an existing EPS
     * * * * bearer context
     */
    OAILOG_WARNING_UE(
        LOG_NAS_ESM, emm_context->_imsi64,
        "ESM-SAP   - Invalid EPS bearer identity (ebi=%d)\n", ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, ESM_CAUSE_INVALID_EPS_BEARER_IDENTITY);
  }

  /*
   * Message processing
   */
  /*
   * Execute the dedicated EPS bearer context deactivation procedure accepted
   * * * * by the UE
   */
  esm_proc_eps_bearer_context_deactivate_accept(emm_context, ebi, &esm_cause);

  /*
   * Return the ESM cause value
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, esm_cause);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/
