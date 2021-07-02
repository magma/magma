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

/*****************************************************************************

  Source      emm_cn.c

  Version     0.1

  Date        2013/12/05

  Product     NAS stack

  Subsystem   EPS Core Network

  Author      Sebastien Roux, Lionel GAUTHIER

  Description

*****************************************************************************/

#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <stdlib.h>

#include "bstrlib.h"
#include "log.h"
#include "common_types.h"
#include "common_defs.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_29.274.h"
#include "3gpp_24.301.h"
#include "mme_app_ue_context.h"
#include "dynamic_memory_check.h"
#include "emm_cn.h"
#include "emm_sap.h"
#include "emm_proc.h"
#include "esm_send.h"
#include "esm_proc.h"
#include "esm_cause.h"
#include "emm_data.h"
#include "conversions.h"
#include "esm_sap.h"
#include "service303.h"
#include "mme_app_itti_messaging.h"
#include "mme_app_apn_selection.h"
#include "mme_config.h"
#include "3gpp_23.003.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "EpsQualityOfService.h"
#include "EpsUpdateType.h"
#include "EsmCause.h"
#include "MobileIdentity.h"
#include "emm_fsm.h"
#include "emm_regDef.h"
#include "esm_data.h"
#include "esm_msg.h"
#include "esm_sapDef.h"
#include "mme_app_defs.h"
#include "mme_app_state.h"
#include "mme_app_messages_types.h"
#include "mme_app_sgs_fsm.h"
#include "nas_procedures.h"
#include "nas/networkDef.h"
#include "security_types.h"

extern int emm_proc_tracking_area_update_accept(
    nas_emm_tau_proc_t* const tau_proc);
/*
   Internal data used for attach procedure
*/
typedef struct {
  unsigned int ue_id; /* UE identifier        */
#define ATTACH_COUNTER_MAX 5
  unsigned int retransmission_count; /* Retransmission counter   */
  bstring esm_msg;                   /* ESM message to be sent within
                                      * the Attach Accept message    */
} attach_data_t;

extern int emm_cn_wrapper_attach_accept(emm_context_t* emm_context);

static int emm_cn_authentication_res(emm_cn_auth_res_t* const msg);
static int emm_cn_authentication_fail(const emm_cn_auth_fail_t* msg);
static int emm_cn_ula_success(emm_cn_ula_success_t* msg_pP);
static int emm_cn_cs_response_success(emm_cn_cs_response_success_t* msg_pP);

/*
   String representation of EMMCN-SAP primitives
*/
static const char* emm_cn_primitive_str[] = {
    "EMM_CN_AUTHENTICATION_PARAM_RES",
    "EMM_CN_AUTHENTICATION_PARAM_FAIL",
    "EMM_CN_ULA_SUCCESS",
    "EMM_CN_CS_RESPONSE_SUCCESS",
    "EMM_CN_ULA_OR_CSRSP_FAIL",
    "EMM_CN_ACTIVATE_DEDICATED_BEARER_REQ",
    "EMMCN_IMPLICIT_DETACH_UE",
    "EMMCN_SMC_PROC_FAIL",
    "EMMCN_NW_INITIATED_DETACH_UE",
    "EMMCN_CS_DOMAIN_LOCATION_UPDT_ACC",
    "EMMCN_CS_DOMAIN_LOCATION_UPDT_FAIL",
    "EMMCN_MM_INFORMATION_REQUEST",
    "EMMCN_DEACTIVATE_BEARER_REQ",
    "EMMCN_PDN_DISCONNECT_RES",
};

//------------------------------------------------------------------------------
static int emm_cn_authentication_res(emm_cn_auth_res_t* const msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_context_t* emm_ctx = NULL;
  int rc                 = RETURNerror;

  /*
   * We received security vector from HSS. Try to setup security with UE
   */
  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
    nas_auth_info_proc_t* auth_info_proc =
        get_nas_cn_procedure_auth_info(emm_ctx);

    if (auth_info_proc) {
      for (int i = 0; i < msg->nb_vectors; i++) {
        auth_info_proc->vector[i] = msg->vector[i];
        msg->vector[i]            = NULL;
      }
      auth_info_proc->nb_vectors = msg->nb_vectors;
      rc                         = (*auth_info_proc->success_notif)(emm_ctx);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-PROC  - "
          "Failed to find Auth_info procedure associated to UE "
          "id " MME_UE_S1AP_ID_FMT "...\n",
          msg->ue_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_cn_authentication_fail(const emm_cn_auth_fail_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_context_t* emm_ctx = NULL;
  int rc                 = RETURNerror;

  /*
   * We received security vector from HSS. Try to setup security with UE
   */
  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
    nas_auth_info_proc_t* auth_info_proc =
        get_nas_cn_procedure_auth_info(emm_ctx);

    if (auth_info_proc) {
      auth_info_proc->nas_cause = msg->cause;
      rc                        = (*auth_info_proc->failure_notif)(emm_ctx);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-PROC  - "
          "Failed to find Auth_info procedure associated to UE "
          "id " MME_UE_S1AP_ID_FMT "...\n",
          msg->ue_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_cn_smc_fail(const emm_cn_smc_fail_t* msg) {
  int rc = RETURNerror;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  rc = emm_proc_attach_reject(msg->ue_id, msg->emm_cause);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
void handle_apn_mismatch(ue_mm_context_t* const ue_context, int esm_cause) {
  ESM_msg esm_msg               = {.header = {0}};
  struct emm_context_s* emm_ctx = NULL;
  uint8_t emm_cn_sap_buffer[EMM_CN_SAP_BUFFER_SIZE];

  if (ue_context == NULL) {
    OAILOG_FUNC_OUT(LOG_MME_APP);
  }

  /*Setup ESM message*/
  emm_ctx = &ue_context->emm_context;
  esm_send_pdn_connectivity_reject(
      emm_ctx->esm_ctx.esm_proc_data->pti, &esm_msg.pdn_connectivity_reject,
      esm_cause);

  int size =
      esm_msg_encode(&esm_msg, emm_cn_sap_buffer, EMM_CN_SAP_BUFFER_SIZE);
  OAILOG_DEBUG(LOG_NAS_EMM, "ESM encoded MSG size %d\n", size);

  if (size > 0) {
    nas_emm_attach_proc_t* attach_proc =
        get_nas_specific_procedure_attach(emm_ctx);
    /*
     * Setup the ESM message container
     */
    attach_proc->esm_msg_out = blk2bstr(emm_cn_sap_buffer, size);
    increment_counter(
        "ue_attach", 1, 2, "result", "failure", "cause",
        "pdn_connection_estb_failed");
    emm_proc_attach_reject(ue_context->mme_ue_s1ap_id, EMM_CAUSE_ESM_FAILURE);
  }

  OAILOG_FUNC_OUT(LOG_MME_APP);
}

//------------------------------------------------------------------------------
static int emm_cn_ula_success(emm_cn_ula_success_t* msg_pP) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                            = RETURNerror;
  struct emm_context_s* emm_ctx     = NULL;
  esm_cause_t esm_cause             = ESM_CAUSE_SUCCESS;
  pdn_cid_t pdn_cid                 = 0;
  ebi_t new_ebi                     = 0;
  bool is_pdn_context_exist_for_apn = false;

  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg_pP->ue_id);

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMCN-SAP  - ue_mm_context Null for ue_id (" MME_UE_S1AP_ID_FMT ")\n",
        msg_pP->ue_id);
  }

  if (emm_ctx == NULL) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMCN-SAP  - "
        "Failed to find UE associated to id " MME_UE_S1AP_ID_FMT "...\n",
        msg_pP->ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  //----------------------------------------------------------------------------
  // PDN selection here
  // Because NAS knows APN selected by UE if any
  // default APN selection
  struct apn_configuration_s* apn_config =
      mme_app_select_apn(ue_mm_context, &esm_cause);

  if (!apn_config) {
    /*
     * Unfortunately we didn't find our default APN...
     */
    OAILOG_ERROR(
        LOG_NAS_ESM, "No suitable APN found ue_id=" MME_UE_S1AP_ID_FMT ")\n",
        ue_mm_context->mme_ue_s1ap_id);
    /*
     * If there is a mismatch between the APN sent by UE and the APN
     * provided by HSS or if we fail to select the APN provided
     * by HSS,send Attach Reject to UE
     */
    handle_apn_mismatch(ue_mm_context, esm_cause);
    return RETURNerror;
  }

  // search for an already set PDN context
  for (pdn_cid = 0; pdn_cid < MAX_APN_PER_UE; pdn_cid++) {
    if ((ue_mm_context->pdn_contexts[pdn_cid]) &&
        (ue_mm_context->pdn_contexts[pdn_cid]->context_identifier ==
         apn_config->context_identifier)) {
      is_pdn_context_exist_for_apn = true;
      break;
    }
  }

  if (pdn_cid >= MAX_APN_PER_UE) {
    /*
     * Search for an available PDN connection entry
     */
    for (pdn_cid = 0; pdn_cid < MAX_APN_PER_UE; pdn_cid++) {
      if (!ue_mm_context->pdn_contexts[pdn_cid]) break;
    }
  }
  if (pdn_cid < MAX_APN_PER_UE) {
    /*
     * Execute the PDN connectivity procedure requested by the UE
     */
    emm_ctx->esm_ctx.esm_proc_data->pdn_cid = pdn_cid;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.qci =
        apn_config->subscribed_qos.qci;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.pci =
        apn_config->subscribed_qos.allocation_retention_priority
            .pre_emp_capability;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.pl =
        apn_config->subscribed_qos.allocation_retention_priority.priority_level;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.pvi =
        apn_config->subscribed_qos.allocation_retention_priority
            .pre_emp_vulnerability;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.gbr.br_ul = 0;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.gbr.br_dl = 0;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.mbr.br_ul = 0;
    emm_ctx->esm_ctx.esm_proc_data->bearer_qos.mbr.br_dl = 0;
    // If UE has not sent APN, use the default APN sent by HSS
    if (emm_ctx->esm_ctx.esm_proc_data->apn == NULL) {
      emm_ctx->esm_ctx.esm_proc_data->apn =
          bfromcstr((const char*) apn_config->service_selection);
    }
    // TODO  "Better to throw emm_ctx->esm_ctx.esm_proc_data as a parameter or
    // as a hidden parameter ?"
    rc = esm_proc_pdn_connectivity_request(
        emm_ctx, emm_ctx->esm_ctx.esm_proc_data->pti,
        emm_ctx->esm_ctx.esm_proc_data->pdn_cid, apn_config->context_identifier,
        emm_ctx->esm_ctx.esm_proc_data->request_type,
        emm_ctx->esm_ctx.esm_proc_data->apn, apn_config->pdn_type,
        emm_ctx->esm_ctx.esm_proc_data->pdn_addr,
        &emm_ctx->esm_ctx.esm_proc_data->bearer_qos,
        (emm_ctx->esm_ctx.esm_proc_data->pco.num_protocol_or_container_id) ?
            &emm_ctx->esm_ctx.esm_proc_data->pco :
            NULL,
        &esm_cause);

    if (rc != RETURNerror) {
      /*
       * Create local default EPS bearer context
       */
      if ((!is_pdn_context_exist_for_apn) ||
          ((is_pdn_context_exist_for_apn) &&
           (EPS_BEARER_IDENTITY_UNASSIGNED ==
            ue_mm_context->pdn_contexts[pdn_cid]->default_ebi))) {
        rc = esm_proc_default_eps_bearer_context(
            emm_ctx, emm_ctx->esm_ctx.esm_proc_data->pti, pdn_cid, &new_ebi,
            emm_ctx->esm_ctx.esm_proc_data->bearer_qos.qci, &esm_cause);
        if (rc < 0) {
          OAILOG_WARNING(
              LOG_NAS_ESM,
              "Failed to Allocate resources required for activation"
              " of a default EPS bearer context for (ue_id =" MME_UE_S1AP_ID_FMT
              ")\n",
              ue_mm_context->mme_ue_s1ap_id);
        }
      }

      if (rc != RETURNerror) {
        esm_cause = ESM_CAUSE_SUCCESS;
      }
    } else {
      OAILOG_ERROR(
          LOG_NAS_ESM,
          "Failed to Perform PDN connectivity procedure requested by ue"
          "for (ue_id =" MME_UE_S1AP_ID_FMT ")\n",
          ue_mm_context->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
    }
    if (!is_pdn_context_exist_for_apn) {
      mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
      if ((mme_app_send_s11_create_session_req(
              mme_app_desc_p, ue_mm_context, pdn_cid)) == RETURNok) {
        increment_counter("mme_spgw_create_session_req", 1, NO_LABELS);
      }
    }
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
}

//------------------------------------------------------------------------------
static int emm_cn_implicit_detach_ue(const uint32_t ue_id) {
  int rc                          = RETURNok;
  struct emm_context_s* emm_ctx_p = NULL;

  OAILOG_FUNC_IN(LOG_NAS_EMM);

  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMM-PROC Implicit Detach UE" MME_UE_S1AP_ID_FMT "\n",
      ue_id);

  emm_detach_request_ies_t params = {0};
  // params.decode_status
  // params.guti = NULL;
  // params.imei = NULL;
  // params.imsi = NULL;
  params.is_native_sc = true;
  params.ksi          = 0;
  params.switch_off   = true;
  params.type         = EMM_DETACH_TYPE_EPS;

  emm_ctx_p = emm_context_get(&_emm_data, ue_id);
  if (emm_ctx_p &&
      (emm_ctx_p->attach_type == EMM_ATTACH_TYPE_COMBINED_EPS_IMSI)) {
    rc = emm_proc_sgs_detach_request(
        ue_id, EMM_SGS_NW_INITIATED_IMPLICIT_NONEPS_DETACH);
  } else {
    rc = emm_proc_sgs_detach_request(ue_id, EMM_SGS_NW_INITIATED_EPS_DETACH);
  }

  emm_proc_detach_request(ue_id, &params);
  increment_counter("ue_detach", 1, 1, "cause", "implicit_detach");
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_cn_nw_initiated_detach_ue(
    const uint32_t ue_id, uint8_t detach_type) {
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "EMM-PROC NW Initiated Detach UE" MME_UE_S1AP_ID_FMT "\n",
      ue_id);
  emm_proc_nw_initiated_detach_request(ue_id, detach_type);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_proc_combined_attach_req(
    struct emm_context_s* emm_ctx_p, bstring esm_data) {
  int rc = RETURNerror;
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_mm_context_t* ue_mm_context_p =
      PARENT_STRUCT(emm_ctx_p, struct ue_mm_context_s, emm_context);

  if (!ue_mm_context_p) {
    OAILOG_ERROR(LOG_NAS_EMM, "Failed to get ue context from emm context \n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  char* non_eps_service_control = bdata(mme_config.non_eps_service_control);

  if (emm_ctx_p->attach_type == EMM_ATTACH_TYPE_COMBINED_EPS_IMSI) {
    // Only send LUR for SMS and CSFB_SMS, not SMS_ORC8R
    if (!(strcmp(non_eps_service_control, "SMS")) ||
        !(strcmp(non_eps_service_control, "CSFB_SMS"))) {
      if (is_mme_ue_context_network_access_mode_packet_only(ue_mm_context_p)) {
        emm_ctx_p->emm_cause = EMM_CAUSE_CS_SERVICE_NOT_AVAILABLE;
        rc                   = RETURNerror;
      } else {
        rc = mme_app_handle_nas_cs_domain_location_update_req(
            ue_mm_context_p, ATTACH_REQUEST);
        /* Store ESM message, Activate Default EPS bearer Context Req to be
         * sent in Attach Accept triggered after receiving
         * SGS-Location Update Accept
         */
        if (rc != RETURNok) {
          OAILOG_ERROR(
              LOG_MME_APP,
              "Failed to send SGS Location Update Request to MSC for"
              "ue_id" MME_UE_S1AP_ID_FMT "\n",
              ue_mm_context_p->mme_ue_s1ap_id);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
        }
        emm_ctx_p->csfbparams.esm_data = esm_data;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
      }
    }
  }
  if (rc != RETURNok) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "Failed to send SGS Location Update Req for ue_id " MME_UE_S1AP_ID_FMT
        "\n",
        ue_mm_context_p->mme_ue_s1ap_id);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
//------------------------------------------------------------------------------
static int emm_cn_cs_response_success(emm_cn_cs_response_success_t* msg_pP) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                           = RETURNerror;
  struct emm_context_s* emm_ctx    = NULL;
  esm_proc_pdn_type_t esm_pdn_type = ESM_PDN_TYPE_IPV4;
  ESM_msg esm_msg                  = {.header = {0}};
  EpsQualityOfService qos          = {0};
  bstring rsp                      = NULL;
  bool is_standalone               = false;  // warning hardcoded
  bool triggered_by_ue             = true;   // warning hardcoded

  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg_pP->ue_id);

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMCN-SAP  - ue mm context null  for ue id " MME_UE_S1AP_ID_FMT "..\n",
        msg_pP->ue_id);
  }

  if (emm_ctx == NULL) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMCN-SAP  - "
        "Failed to find UE associated to id " MME_UE_S1AP_ID_FMT "...\n",
        msg_pP->ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  memset(&esm_msg, 0, sizeof(ESM_msg));

  is_standalone = (emm_ctx->esm_ctx.pending_standalone > 0 ? true : false);
  /* msg_pP->pdn_type enum is as per 29.272( S6a spec).Convert the PDN type
   * to be aligned to 3gpp 24.301 before sending the NAS message to UE
   */
  switch (msg_pP->pdn_type) {
    case IPv4:
      OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  esm_pdn_type = ESM_PDN_TYPE_IPV4\n");
      esm_pdn_type = ESM_PDN_TYPE_IPV4;
      break;

    case IPv6:
      OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  esm_pdn_type = ESM_PDN_TYPE_IPV6\n");
      esm_pdn_type = ESM_PDN_TYPE_IPV6;
      break;

    case IPv4_AND_v6:
      OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  esm_pdn_type = ESM_PDN_TYPE_IPV4V6\n");
      esm_pdn_type = ESM_PDN_TYPE_IPV4V6;
      break;

    default:
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "EMM  -  esm_pdn_type = ESM_PDN_TYPE_IPV4 (forced to default)\n");
      esm_pdn_type = ESM_PDN_TYPE_IPV4;
  }

  OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  qci       = %u \n", msg_pP->qci);
  OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  qos.qci   = %u \n", msg_pP->qos.qci);
  OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  qos.mbrUL = %u \n", msg_pP->qos.mbrUL);
  OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  qos.mbrDL = %u \n", msg_pP->qos.mbrDL);
  OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  qos.gbrUL = %u \n", msg_pP->qos.gbrUL);
  OAILOG_DEBUG(LOG_NAS_EMM, "EMM  -  qos.gbrDL = %u \n", msg_pP->qos.gbrDL);
  qos.bitRatesPresent    = 0;
  qos.bitRatesExtPresent = 0;
  //#pragma message "Some work to do here about qos"
  qos.qci                          = msg_pP->qci;
  qos.bitRates.maxBitRateForUL     = msg_pP->qos.mbrUL;
  qos.bitRates.maxBitRateForDL     = msg_pP->qos.mbrDL;
  qos.bitRates.guarBitRateForUL    = msg_pP->qos.gbrUL;
  qos.bitRates.guarBitRateForDL    = msg_pP->qos.gbrDL;
  qos.bitRatesExt.maxBitRateForUL  = 0;
  qos.bitRatesExt.maxBitRateForDL  = 0;
  qos.bitRatesExt.guarBitRateForUL = 0;
  qos.bitRatesExt.guarBitRateForDL = 0;

  int def_bearer_index = EBI_TO_INDEX(msg_pP->ebi);
  pdn_cid_t pdn_cid =
      ue_mm_context->bearer_contexts[def_bearer_index]->pdn_cx_id;

  if (ue_mm_context->pdn_contexts[pdn_cid] == NULL) {
    OAILOG_ERROR_UE(
        LOG_MME_APP, ue_mm_context->emm_context._imsi64,
        "pdn_contexts is NULL for "
        "MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "ebi-%u\n",
        ue_mm_context->mme_ue_s1ap_id, msg_pP->ebi);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  // Return default EPS bearer context request message
  rc = esm_send_activate_default_eps_bearer_context_request(
      msg_pP->pti, msg_pP->ebi,
      &esm_msg.activate_default_eps_bearer_context_request,
      ue_mm_context->pdn_contexts[pdn_cid], &msg_pP->pco, esm_pdn_type,
      msg_pP->pdn_addr, &qos,
      ue_mm_context->pdn_contexts[pdn_cid]->esm_data.esm_cause);
  clear_protocol_configuration_options(&msg_pP->pco);
  if (rc != RETURNerror) {
    // Encode the returned ESM response message
    char emm_cn_sap_buffer[EMM_CN_SAP_BUFFER_SIZE];
    int size = esm_msg_encode(
        &esm_msg, (uint8_t*) emm_cn_sap_buffer, EMM_CN_SAP_BUFFER_SIZE);

    OAILOG_DEBUG(LOG_NAS_EMM, "ESM encoded MSG size %d\n", size);

    bdestroy_wrapper(&msg_pP->pdn_addr);
    if (size > 0) {
      rsp = blk2bstr(emm_cn_sap_buffer, size);
    }

    // Complete the relevant ESM procedure
    rc = esm_proc_default_eps_bearer_context_request(
        is_standalone, emm_ctx, msg_pP->ebi, &rsp, triggered_by_ue);

    // Free protocol configuration options and its contents
    clear_protocol_configuration_options(
        &esm_msg.activate_default_eps_bearer_context_request
             .protocolconfigurationoptions);

    if ((rc != RETURNok) || (is_standalone)) {
      // Return indication that ESM procedure failed
      bdestroy_wrapper(&rsp);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
    }
  } else {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "ESM send activate_default_eps_bearer_context_request "
        "failed " MME_UE_S1AP_ID_FMT "\n",
        ue_mm_context->mme_ue_s1ap_id);
  }

  /* Notify SGS Location update request to MME App
   * if attach_type == EMM_ATTACH_TYPE_COMBINED_EPS_IMSI
   */
  if (emm_proc_combined_attach_req(emm_ctx, rsp) == RETURNok) {
    bdestroy_wrapper(&rsp);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  // OAILOG_INFO (LOG_NAS_EMM, "EMM - APN = %s\n", (char *)bdata(msg_pP->apn));
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_ctx);

  if (attach_proc) {
    // Setup the ESM message container
    attach_proc->esm_msg_out = rsp;

    // Send attach accept message to the UE
    rc = emm_cn_wrapper_attach_accept(emm_ctx);

    if (rc != RETURNerror) {
      if (IS_EMM_CTXT_PRESENT_OLD_GUTI(emm_ctx) &&
          (memcmp(
              &emm_ctx->_old_guti, &emm_ctx->_guti, sizeof(emm_ctx->_guti)))) {
        /*
         * Implicit GUTI reallocation;
         * Notify EMM that common procedure has been initiated
         * LG: TODO check this, seems very suspicious
         */
        emm_sap_t emm_sap = {0};

        emm_sap.primitive       = EMMREG_COMMON_PROC_REQ;
        emm_sap.u.emm_reg.ue_id = msg_pP->ue_id;
        emm_sap.u.emm_reg.ctx   = emm_ctx;

        rc = emm_sap_send(&emm_sap);
      }
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_cn_ula_or_csrsp_fail(const emm_cn_ula_or_csrsp_fail_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                          = RETURNok;
  struct emm_context_s* emm_ctx_p = NULL;
  ESM_msg esm_msg                 = {.header = {0}};
  int esm_cause;
  emm_ctx_p = emm_context_get(&_emm_data, msg->ue_id);
  if (emm_ctx_p == NULL) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMCN-SAP  - "
        "Failed to find UE associated to id " MME_UE_S1AP_ID_FMT "...\n",
        msg->ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  memset(&esm_msg, 0, sizeof(ESM_msg));

  // Map S11 cause to ESM cause
  switch (msg->cause) {
    case CAUSE_CONTEXT_NOT_FOUND:
      esm_cause = ESM_CAUSE_REQUEST_REJECTED_BY_GW;
      break;
    case CAUSE_INVALID_MESSAGE_FORMAT:
      esm_cause = ESM_CAUSE_REQUEST_REJECTED_BY_GW;
      break;
    case CAUSE_SERVICE_NOT_SUPPORTED:
      esm_cause = ESM_CAUSE_SERVICE_OPTION_NOT_SUPPORTED;
      break;
    case CAUSE_SYSTEM_FAILURE:
      esm_cause = ESM_CAUSE_NETWORK_FAILURE;
      break;
    case CAUSE_NO_RESOURCES_AVAILABLE:
      esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
      increment_counter(
          "ue_pdn_connectivity_req", 1, 2, "result", "failure", "cause",
          "no_resources_available");
      break;
    case CAUSE_ALL_DYNAMIC_ADDRESSES_OCCUPIED:
      esm_cause = ESM_CAUSE_INSUFFICIENT_RESOURCES;
      increment_counter(
          "ue_pdn_connectivity_req", 1, 2, "result", "failure", "cause",
          "all_dynamic_resources_occupied");
      break;
    default:
      esm_cause = ESM_CAUSE_REQUEST_REJECTED_BY_GW;
      break;
  }

  rc = esm_send_pdn_connectivity_reject(
      msg->pti, &esm_msg.pdn_connectivity_reject, esm_cause);
  /*
   * Encode the returned ESM response message
   */
  uint8_t emm_cn_sap_buffer[EMM_CN_SAP_BUFFER_SIZE];
  int size =
      esm_msg_encode(&esm_msg, emm_cn_sap_buffer, EMM_CN_SAP_BUFFER_SIZE);
  OAILOG_INFO(LOG_NAS_EMM, "ESM encoded MSG size %d\n", size);

  if (size > 0) {
    nas_emm_attach_proc_t* attach_proc =
        get_nas_specific_procedure_attach(emm_ctx_p);
    if (attach_proc) {
      /*
       * Setup the ESM message container
       */
      attach_proc->esm_msg_out = blk2bstr(emm_cn_sap_buffer, size);
    }
    increment_counter(
        "ue_attach", 1, 2, "result", "failure", "cause",
        "pdn_connection_estb_failed");
    rc = emm_proc_attach_reject(msg->ue_id, EMM_CAUSE_ESM_FAILURE);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_cn_activate_dedicated_bearer_req(
    emm_cn_activate_dedicated_bearer_req_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNok;
  // forward to ESM
  esm_sap_t esm_sap = {0};

  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

  esm_sap.primitive     = ESM_DEDICATED_EPS_BEARER_CONTEXT_ACTIVATE_REQ;
  esm_sap.ctx           = &ue_mm_context->emm_context;
  esm_sap.is_standalone = true;
  esm_sap.ue_id         = msg->ue_id;
  esm_sap.data.eps_dedicated_bearer_context_activate.cid = msg->cid;
  esm_sap.data.eps_dedicated_bearer_context_activate.ebi = msg->ebi;
  esm_sap.data.eps_dedicated_bearer_context_activate.linked_ebi =
      msg->linked_ebi;
  esm_sap.data.eps_dedicated_bearer_context_activate.tft = msg->tft;
  esm_sap.data.eps_dedicated_bearer_context_activate.qci = msg->bearer_qos.qci;
  esm_sap.data.eps_dedicated_bearer_context_activate.gbr_ul =
      msg->bearer_qos.gbr.br_ul;
  esm_sap.data.eps_dedicated_bearer_context_activate.gbr_dl =
      msg->bearer_qos.gbr.br_dl;
  esm_sap.data.eps_dedicated_bearer_context_activate.mbr_ul =
      msg->bearer_qos.mbr.br_ul;
  esm_sap.data.eps_dedicated_bearer_context_activate.mbr_dl =
      msg->bearer_qos.mbr.br_dl;
  esm_sap.data.eps_dedicated_bearer_context_activate.pco = msg->pco;
  memcpy(
      &esm_sap.data.eps_dedicated_bearer_context_activate.sgw_fteid,
      &msg->sgw_fteid, sizeof(fteid_t));

  rc = esm_sap_send(&esm_sap);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_cn_deactivate_dedicated_bearer_req(
    emm_cn_deactivate_dedicated_bearer_req_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNok;
  // forward to ESM
  esm_sap_t esm_sap = {0};

  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

  esm_sap.primitive     = ESM_EPS_BEARER_CONTEXT_DEACTIVATE_REQ;
  esm_sap.ctx           = &ue_mm_context->emm_context;
  esm_sap.is_standalone = true;
  esm_sap.ue_id         = msg->ue_id;
  esm_sap.data.eps_bearer_context_deactivate.is_pcrf_initiated = true;
  /*Currently we only support deactivation of a single bearer at NAS*/
  esm_sap.data.eps_bearer_context_deactivate.no_of_bearers = msg->no_of_bearers;
  memcpy(
      esm_sap.data.eps_bearer_context_deactivate.ebi, msg->ebi, sizeof(ebi_t));

  rc = esm_sap_send(&esm_sap);

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int send_attach_accept(struct emm_context_s* emm_ctx_p) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNok;

  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_ctx_p);

  /*
   * Setup the ESM message container
   */
  if (attach_proc == NULL) {
    OAILOG_ERROR(LOG_NAS_EMM, "attach_proc is NULL in _send_attach_accept\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  attach_proc->esm_msg_out = emm_ctx_p->csfbparams.esm_data;
  /*
   * Send attach accept message to the UE
   */
  rc = emm_cn_wrapper_attach_accept(emm_ctx_p);

  if (rc != RETURNerror) {
    if (IS_EMM_CTXT_PRESENT_OLD_GUTI(emm_ctx_p) &&
        (memcmp(
            &emm_ctx_p->_old_guti, &emm_ctx_p->_guti,
            sizeof(emm_ctx_p->_guti)))) {
      /*
       * Implicit GUTI reallocation;
       * * * * Notify EMM that common procedure has been initiated
       */
      emm_sap_t emm_sap = {0};

      emm_sap.primitive       = EMMREG_COMMON_PROC_REQ;
      emm_sap.u.emm_reg.ue_id = attach_proc->ue_id;
      emm_sap.u.emm_reg.ctx   = emm_ctx_p;
      rc                      = emm_sap_send(&emm_sap);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int emm_send_cs_domain_attach_or_tau_accept(
    struct ue_mm_context_s* ue_context_p) {
  int rc                    = RETURNok;
  emm_fsm_state_t fsm_state = EMM_DEREGISTERED;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*If Location Updt Accept is received for TAU,trigger TAU accept
   * If Location Updt Accept is received for Attach, trigger Attach accept
   */
  fsm_state = emm_fsm_get_state(&ue_context_p->emm_context);
  /* If fsm_state is REGISTERED, Location Update Accept is received for
   * TAU procedure
   * If fsm_state is DE-REGISTERED,Location Update Accept is received for
   * Attach procedure
   */
  if (fsm_state == EMM_REGISTERED) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMCN-SAP  - ue_context_p->emm_context.csfbparams.presencemask %d\n",
        ue_context_p->emm_context.csfbparams.presencemask);
    // Trigger Sending of TAU Accept
    nas_emm_tau_proc_t* tau_proc =
        get_nas_specific_procedure_tau(&ue_context_p->emm_context);
    if (tau_proc) {
      if ((emm_proc_tracking_area_update_accept(tau_proc)) == RETURNerror) {
        OAILOG_ERROR(
            LOG_NAS_EMM,
            "EMMCN-SAP  - "
            "Failed to send TAU accept for UE id " MME_UE_S1AP_ID_FMT "...\n",
            ue_context_p->mme_ue_s1ap_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
      }
    }
  } else {
    if (send_attach_accept(&ue_context_p->emm_context) == RETURNerror) {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMMCN-SAP  - "
          "Failed to send Attach accept for UE id " MME_UE_S1AP_ID_FMT "...\n",
          ue_context_p->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
//------------------------------------------------------------------------------
static int emm_cn_cs_domain_loc_updt_fail(
    emm_cn_cs_domain_location_updt_fail_t emm_cn_sgs_location_updt_fail) {
  int rc = RETURNok;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMM-PROC SGS LOCATION UPDATE FAILURE" MME_UE_S1AP_ID_FMT "\n",
      emm_cn_sgs_location_updt_fail.ue_id);

  emm_fsm_state_t fsm_state       = EMM_DEREGISTERED;
  struct emm_context_s* emm_ctx_p = NULL;
  emm_ctx_p = emm_context_get(&_emm_data, emm_cn_sgs_location_updt_fail.ue_id);
  if (emm_ctx_p == NULL) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMCN-SAP  - "
        "Failed to find UE associated to id " MME_UE_S1AP_ID_FMT "...\n",
        emm_cn_sgs_location_updt_fail.ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  // Store LAI to be sent in Attach Accept/TAU Accept
  if (emm_cn_sgs_location_updt_fail.presencemask & LAI) {
    emm_ctx_p->csfbparams.lai.mccdigit2 =
        emm_cn_sgs_location_updt_fail.laicsfb.mccdigit2;
    emm_ctx_p->csfbparams.lai.mccdigit1 =
        emm_cn_sgs_location_updt_fail.laicsfb.mccdigit1;
    emm_ctx_p->csfbparams.lai.mncdigit3 =
        emm_cn_sgs_location_updt_fail.laicsfb.mncdigit3;
    emm_ctx_p->csfbparams.lai.mccdigit3 =
        emm_cn_sgs_location_updt_fail.laicsfb.mccdigit3;
    emm_ctx_p->csfbparams.lai.mncdigit2 =
        emm_cn_sgs_location_updt_fail.laicsfb.mncdigit2;
    emm_ctx_p->csfbparams.lai.mncdigit1 =
        emm_cn_sgs_location_updt_fail.laicsfb.mncdigit1;
    emm_ctx_p->csfbparams.lai.lac = emm_cn_sgs_location_updt_fail.laicsfb.lac;
    emm_ctx_p->csfbparams.presencemask |= LAI_CSFB;
  }
  // Store SGS Reject Cause to be sent in Attach Accept
  emm_ctx_p->emm_cause = emm_cn_sgs_location_updt_fail.reject_cause;

  // Store the status of Location Update procedure(success/failure) to send
  // appropriate cause in Attach Accept
  emm_ctx_p->csfbparams.sgs_loc_updt_status = FAILURE;
  fsm_state                                 = emm_fsm_get_state(emm_ctx_p);

  if (fsm_state == EMM_DEREGISTERED) {
    nas_emm_attach_proc_t* attach_proc =
        get_nas_specific_procedure_attach(emm_ctx_p);
    /*
     * Setup the ESM message container
     */
    attach_proc->esm_msg_out = emm_ctx_p->csfbparams.esm_data;
    /*
     * Send attach accept message to the UE
     */
    rc = emm_cn_wrapper_attach_accept(emm_ctx_p);

    if (rc != RETURNerror) {
      if (IS_EMM_CTXT_PRESENT_OLD_GUTI(emm_ctx_p) &&
          (memcmp(
              &emm_ctx_p->_old_guti, &emm_ctx_p->_guti,
              sizeof(emm_ctx_p->_guti)))) {
        /*
         * Implicit GUTI reallocation;
         * * * * Notify EMM that common procedure has been initiated
         */
        emm_sap_t emm_sap       = {0};
        emm_sap.primitive       = EMMREG_COMMON_PROC_REQ;
        emm_sap.u.emm_reg.ue_id = emm_cn_sgs_location_updt_fail.ue_id;
        emm_sap.u.emm_reg.ctx   = emm_ctx_p;
        rc                      = emm_sap_send(&emm_sap);
      }
    }
  } else if (EMM_REGISTERED == fsm_state) {
    /* If TAU Update type is not TA_LA_UPDATING_WITH_IMSI_ATTACH and MSC/VLR
     * rejected Location Update Req send TAU Accept with appropriate emm_cause
     */
    if (emm_ctx_p->tau_updt_type !=
        EPS_UPDATE_TYPE_COMBINED_TA_LA_UPDATING_WITH_IMSI_ATTACH) {
      // Trigger Sending of TAU Accept
      nas_emm_tau_proc_t* tau_proc = get_nas_specific_procedure_tau(emm_ctx_p);
      if (tau_proc) rc = emm_proc_tracking_area_update_accept(tau_proc);
    } else {
      // Send TAU Reject to UE and initiate IMSI detach towards MSC
      rc = emm_proc_tracking_area_update_reject(
          emm_cn_sgs_location_updt_fail.ue_id, emm_ctx_p->emm_cause);
    }
  }
  // Store SGS Reject Cause to be sent in Attach/TAU
  emm_ctx_p->emm_cause = emm_cn_sgs_location_updt_fail.reject_cause;
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

// handle CS-Domain MM-Information Request from SGS task
static int emm_cn_cs_domain_mm_information_req(
    emm_cn_cs_domain_mm_information_req_t* mm_information_req_pP) {
  int rc                        = RETURNerror;
  imsi64_t imsi64               = INVALID_IMSI64;
  emm_context_t* ctxt           = NULL;
  ue_mm_context_t* ue_context_p = NULL;
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (mm_information_req_pP == NULL) {
    OAILOG_WARNING(LOG_NAS_EMM, "Received mm_information_req_pP is NULL \n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  IMSI_STRING_TO_IMSI64((char*) mm_information_req_pP->imsi, &imsi64);
  OAILOG_DEBUG(
      LOG_NAS_EMM, "Received MM Information Request for IMSI " IMSI_64_FMT "\n",
      imsi64);

  ctxt = emm_context_get_by_imsi(&_emm_data, imsi64);
  if (!(ctxt)) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "That's embarrassing as we don't know this IMSI " IMSI_64_FMT "\n",
        imsi64);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  /* Referring to 29118-d70, section 5.10.3 currently implementing with option-2
   to discard the received MM information and send EMM Information with contents
   generated locally by MME */
  ue_context_p = PARENT_STRUCT(ctxt, struct ue_mm_context_s, emm_context);

  if ((ue_context_p) && (ue_context_p->ecm_state == ECM_CONNECTED)) {
    if (ue_context_p->sgs_context == NULL) {
      OAILOG_WARNING(
          LOG_NAS_EMM, " Invalid SGS context for IMSI" IMSI_64_FMT "\n",
          imsi64);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
    }

    if (ue_context_p->sgs_context->sgs_state == SGS_ASSOCIATED) {
      if ((rc = emm_proc_emm_informtion(ue_context_p)) != RETURNok) {
        OAILOG_WARNING(
            LOG_NAS_EMM,
            "Failed to send Emm Information Reqest for " IMSI_64_FMT "\n",
            imsi64);
      } else {
        OAILOG_DEBUG(
            LOG_NAS_EMM, "Sent Emm Information Reqest for " IMSI_64_FMT "\n",
            imsi64);
      }
    }
  } else {
    // TODO - CSFB Handle MM Information Request in idle mode
    OAILOG_INFO(
        LOG_NAS_EMM,
        "Received MM Information Request while UE is in idle mode for "
        "imsi" IMSI_64_FMT "\n",
        imsi64);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

static int emm_cn_pdn_disconnect_rsp(emm_cn_pdn_disconnect_rsp_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNok;
  proc_tid_t pti;
#define ESM_SAP_BUFFER_SIZE 4096
  uint8_t esm_sap_buffer[ESM_SAP_BUFFER_SIZE];
  esm_cause_t esm_cause = ESM_CAUSE_SUCCESS;
  bstring rsp           = NULL;

  ESM_msg esm_msg = {0};

  ue_mm_context_t* ue_mm_context_p =
      mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

  pti = ue_mm_context_p->emm_context.esm_ctx.esm_proc_data->pti;
  /*
   * Build the ESM message to be sent to UE
   */
  rc = esm_send_deactivate_eps_bearer_context_request(
      pti, msg->lbi, &esm_msg.deactivate_eps_bearer_context_request,
      ESM_CAUSE_REGULAR_DEACTIVATION);

  if (rc != RETURNok) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "Building of deactivate_eps_bearer_context_request failed for bearer"
        "%u for ue id " MME_UE_S1AP_ID_FMT "\n",
        msg->lbi, ue_mm_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }
  /*
   * Encode the ESM message
   */
  int size =
      esm_msg_encode(&esm_msg, (uint8_t*) esm_sap_buffer, ESM_SAP_BUFFER_SIZE);

  if (size > 0) {
    rsp = blk2bstr(esm_sap_buffer, size);
  } else {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "Message encoding failed for bearer %u for ue id " MME_UE_S1AP_ID_FMT
        "\n",
        msg->lbi, ue_mm_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  /*
   * Send the ESM message
   */
  rc = esm_proc_eps_bearer_context_deactivate_request(
      true /*standalone*/, &ue_mm_context_p->emm_context, msg->lbi, &rsp,
      true /*UE triggered*/);

  bdestroy_wrapper(&rsp);
  if (rc != RETURNok) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "Processing eps_bearer_context_deactivate_request failed for bearer"
        "%u for ue id " MME_UE_S1AP_ID_FMT "\n",
        msg->lbi, ue_mm_context_p->mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  /*
   * Execute the PDN disconnect procedure requested by the UE
   */
  // Fetch the pdn conn id using the linked bearer id
  int pid = ue_mm_context_p->bearer_contexts[EBI_TO_INDEX(msg->lbi)]->pdn_cx_id;

  if (pid < MAX_APN_PER_UE) {
    /*
     * Release the associated default EPS bearer context
     */
    int bid = 0;
    rc      = esm_proc_eps_bearer_context_deactivate(
        &ue_mm_context_p->emm_context, false /* Release locally */, msg->lbi,
        &pid, &bid, &esm_cause);
    if (rc != RETURNok) {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "Processing bearer deactivation failed for ebi %u for ue "
          "id " MME_UE_S1AP_ID_FMT "\n",
          msg->lbi, ue_mm_context_p->mme_ue_s1ap_id);
    }
  } else {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "ESM-PROC  - No PDN connection found (lbi=%u) for ue "
        "id " MME_UE_S1AP_ID_FMT "\n",
        msg->lbi, ue_mm_context_p->mme_ue_s1ap_id);
    rc = RETURNerror;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
int emm_cn_send(const emm_cn_t* msg) {
  int rc                       = RETURNerror;
  emm_cn_primitive_t primitive = msg->primitive;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_INFO(
      LOG_NAS_EMM, "EMMCN-SAP - Received primitive %s (%d)\n",
      emm_cn_primitive_str[primitive - _EMMCN_START - 1], primitive);

  switch (primitive) {
    case _EMMCN_AUTHENTICATION_PARAM_RES:
      rc = emm_cn_authentication_res(msg->u.auth_res);
      break;

    case _EMMCN_AUTHENTICATION_PARAM_FAIL:
      rc = emm_cn_authentication_fail(msg->u.auth_fail);
      break;

    case EMMCN_ULA_SUCCESS:
      rc = emm_cn_ula_success(msg->u.emm_cn_ula_success);
      break;

    case EMMCN_CS_RESPONSE_SUCCESS:
      rc = emm_cn_cs_response_success(msg->u.emm_cn_cs_response_success);
      break;

    case EMMCN_ULA_OR_CSRSP_FAIL:
      rc = emm_cn_ula_or_csrsp_fail(msg->u.emm_cn_ula_or_csrsp_fail);
      break;

    case EMMCN_ACTIVATE_DEDICATED_BEARER_REQ:
      rc = emm_cn_activate_dedicated_bearer_req(
          msg->u.activate_dedicated_bearer_req);
      break;

    case EMMCN_IMPLICIT_DETACH_UE:
      rc = emm_cn_implicit_detach_ue(msg->u.emm_cn_implicit_detach.ue_id);
      break;

    case EMMCN_SMC_PROC_FAIL:
      rc = emm_cn_smc_fail(msg->u.smc_fail);
      break;
    case EMMCN_NW_INITIATED_DETACH_UE:
      rc = emm_cn_nw_initiated_detach_ue(
          msg->u.emm_cn_nw_initiated_detach.ue_id,
          msg->u.emm_cn_nw_initiated_detach.detach_type);
      break;

    case EMMCN_CS_DOMAIN_LOCATION_UPDT_FAIL:
      rc = emm_cn_cs_domain_loc_updt_fail(
          msg->u.emm_cn_cs_domain_location_updt_fail);
      break;

    case EMMCN_CS_DOMAIN_MM_INFORMATION_REQ:
      rc = emm_cn_cs_domain_mm_information_req(
          msg->u.emm_cn_cs_domain_mm_information_req);
      break;

    case EMMCN_DEACTIVATE_BEARER_REQ:
      rc = emm_cn_deactivate_dedicated_bearer_req(
          msg->u.deactivate_dedicated_bearer_req);
      break;

    case EMMCN_PDN_DISCONNECT_RES:
      rc = emm_cn_pdn_disconnect_rsp(msg->u.emm_cn_pdn_disconnect_rsp);
      break;

    default:
      /*
       * Other primitives are forwarded to the Access Stratum
       */
      rc = RETURNerror;
      break;
  }

  if (rc != RETURNok) {
    OAILOG_ERROR(
        LOG_NAS_EMM, "EMMCN-SAP - Failed to process primitive %s (%d)\n",
        emm_cn_primitive_str[primitive - _EMMCN_START - 1], primitive);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
