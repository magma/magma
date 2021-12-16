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

  Source      Attach.c

  Version     0.1

  Date        2012/12/04

  Product     NAS stack

  Subsystem   EPS Mobility Management

  Author      Frederic Maurel, Lionel GAUTHIER

  Description Defines the attach related EMM procedure executed by the
        Non-Access Stratum.

        To get internet connectivity from the network, the network
        have to know about the UE. When the UE is switched on, it
        has to initiate the attach procedure to get initial access
        to the network and register its presence to the Evolved
        Packet Core (EPC) network in order to receive EPS services.

        As a result of a successful attach procedure, a context is
        created for the UE in the MME, and a default bearer is esta-
        blished between the UE and the PDN-GW. The UE gets the home
        agent IPv4 and IPv6 addresses and full connectivity to the
        IP network.

        The network may also initiate the activation of additional
        dedicated bearers for the support of a specific service.

*****************************************************************************/

#include <stdint.h>
#include <stdbool.h>
#include <string.h>
#include <stdlib.h>

#include "lte/gateway/c/core/oai/common/assertions.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_ies.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/mme_app_state.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.h"
#include "lte/gateway/c/core/oai/include/mme_app_ue_context.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_itti_messaging.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/include/mme_events.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_procedures.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/network/nas_message.h"
#include "lte/gateway/c/core/oai/tasks/nas/util/nas_timer.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"

#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.008.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_36.401.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_23.003.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_24.301.h"
#include "lte/gateway/c/core/oai/include/3gpp_requirements_24.301.h"

#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_proc.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_sap.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/emm_data.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/emm_cause.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_asDef.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_cnDef.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_fsm.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/sap/emm_regDef.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/esm_data.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/sap/esm_sapDef.h"
#include "lte/gateway/c/core/oai/tasks/nas/esm/sap/esm_sap.h"
#include "lte/gateway/c/core/oai/tasks/nas/nas_proc.h"

#include "lte/gateway/c/core/oai/tasks/nas/ies/AdditionalUpdateType.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/EmmCommon.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EmmCause.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/EpsNetworkFeatureSupport.h"
#include "lte/gateway/c/core/oai/include/TrackingAreaIdentity.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/TrackingAreaIdentityList.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/mme_app/mme_app_timer.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/* String representation of the EPS attach type */
static const char* emm_attach_type_str[] = {"EPS", "IMSI", "EMERGENCY",
                                            "RESERVED"};

/*
   --------------------------------------------------------------------------
        Internal data handled by the attach procedure in the MME
   --------------------------------------------------------------------------
*/

/*
   Functions that may initiate EMM common procedures
*/
static int emm_start_attach_proc_authentication(
    emm_context_t* emm_context, nas_emm_attach_proc_t* attach_proc);
static int emm_start_attach_proc_security(
    emm_context_t* emm_context, nas_emm_attach_proc_t* attach_proc);

static int emm_attach_security_a(emm_context_t* emm_context);
static int emm_attach(emm_context_t* emm_context);

static int emm_attach_success_identification_cb(emm_context_t* emm_context);
static int emm_attach_failure_identification_cb(emm_context_t* emm_context);
static int emm_attach_success_authentication_cb(emm_context_t* emm_context);
static int emm_attach_failure_authentication_cb(emm_context_t* emm_context);
static int emm_attach_success_security_cb(emm_context_t* emm_context);
static int emm_attach_failure_security_cb(emm_context_t* emm_context);
static int emm_attach_identification_after_smc_success_cb(
    emm_context_t* emm_context);

/*
   Abnormal case attach procedures
*/
static int emm_attach_release(emm_context_t* emm_context);
static int emm_attach_abort(
    struct emm_context_s* emm_context, struct nas_base_proc_s* base_proc);
static int emm_attach_run_procedure(emm_context_t* emm_context);
static int emm_send_attach_accept(emm_context_t* emm_context);

static bool emm_attach_ies_have_changed(
    mme_ue_s1ap_id_t ue_id, emm_attach_request_ies_t* const ies1,
    emm_attach_request_ies_t* const ies2);

static void emm_proc_create_procedure_attach_request(
    ue_mm_context_t* const ue_mm_context,
    STOLEN_REF emm_attach_request_ies_t* const ies);

static int emm_attach_update(
    emm_context_t* const emm_context, emm_attach_request_ies_t* const ies);

static int emm_attach_accept_retx(emm_context_t* emm_context);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
            Attach procedure executed by the MME
   --------------------------------------------------------------------------
*/
/*
 *
 * Name:    emm_proc_attach_request()
 *
 * Description: Performs the UE requested attach procedure
 *
 *              3GPP TS 24.301, section 5.5.1.2.3
 *      The network may initiate EMM common procedures, e.g. the
 *      identification, authentication and security mode control
 *      procedures during the attach procedure, depending on the
 *      information received in the ATTACH REQUEST message (e.g.
 *      IMSI, GUTI and KSI).
 *
 * Inputs:  ue_id:      UE lower layer identifier
 *      type:      Type of the requested attach
 *      ies:       Information ElementStrue if the security context is of type
 *      ctx_is_new:   Is the mm context has been newly created in the context of
 * this procedure Others:    _emm_data
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    _emm_data
 *
 */
//------------------------------------------------------------------------------
status_code_e emm_proc_attach_request(
    mme_ue_s1ap_id_t ue_id, const bool is_mm_ctx_new,
    STOLEN_REF emm_attach_request_ies_t* const ies) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  ue_mm_context_t ue_ctx;
  emm_fsm_state_t fsm_state       = EMM_DEREGISTERED;
  bool clear_emm_ctxt             = false;
  bool is_unknown_guti            = false;
  ue_mm_context_t* ue_mm_context  = NULL;
  ue_mm_context_t* guti_ue_mm_ctx = NULL;
  ue_mm_context_t* imsi_ue_mm_ctx = NULL;
  emm_context_t* new_emm_ctx      = NULL;
  imsi64_t imsi64                 = INVALID_IMSI64;
  mme_ue_s1ap_id_t old_ue_id      = INVALID_MME_UE_S1AP_ID;

  if (ies->imsi) {
    imsi64 = imsi_to_imsi64(ies->imsi);
    OAILOG_INFO(
        LOG_NAS_EMM,
        "ATTACH REQ (ue_id = " MME_UE_S1AP_ID_FMT ") (IMSI = " IMSI_64_FMT
        ") \n",
        ue_id, imsi64);
  } else if (ies->guti) {
    OAILOG_INFO(
        LOG_NAS_EMM,
        "ATTACH REQ (ue_id = " MME_UE_S1AP_ID_FMT ") (GUTI = " GUTI_FMT ") \n",
        ue_id, GUTI_ARG(ies->guti));
  } else if (ies->imei) {
    char imei_str[16];
    IMEI_TO_STRING(ies->imei, imei_str, 16);
    OAILOG_INFO(
        LOG_NAS_EMM,
        "ATTACH REQ (ue_id = " MME_UE_S1AP_ID_FMT ") (IMEI = %s ) \n", ue_id,
        imei_str);
  }

  OAILOG_INFO(
      LOG_NAS_EMM, "EMM-PROC:  ATTACH - EPS attach type = %s (%d)\n",
      emm_attach_type_str[ies->type], ies->type);
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "is_initial request = %u\n (ue_id=" MME_UE_S1AP_ID_FMT
      ") \n(imsi = " IMSI_64_FMT ") \n",
      ies->is_initial, ue_id, imsi64);
  /*
   * Initialize the temporary UE context
   */
  memset(&ue_ctx, 0, sizeof(ue_mm_context_t));
  ue_ctx.emm_context.is_dynamic = false;
  ue_ctx.mme_ue_s1ap_id         = ue_id;

  // Check whether request if for emergency bearer service.

  /*
   * Requirement MME24.301R10_5.5.1.1_1
   * MME not configured to support attach for emergency bearer services
   * shall reject any request to attach with an attach type set to "EPS
   * emergency attach".
   */
  if (!(_emm_data.conf.eps_network_feature_support[0] &
        EPS_NETWORK_FEATURE_SUPPORT_EMERGENCY_BEARER_SERVICES_IN_S1_MODE_SUPPORTED) &&
      (EMM_ATTACH_TYPE_EMERGENCY == ies->type)) {
    REQUIREMENT_3GPP_24_301(R10_5_5_1__1);
    // TODO: update this if/when emergency attach is supported
    ue_ctx.emm_context.emm_cause = ies->imei ? EMM_CAUSE_IMEI_NOT_ACCEPTED :
                                               EMM_CAUSE_NOT_AUTHORIZED_IN_PLMN;
    /*
     * Do not accept the UE to attach for emergency services
     */
    struct nas_emm_attach_proc_s no_attach_proc = {0};
    no_attach_proc.ue_id                        = ue_id;
    no_attach_proc.emm_cause                    = ue_ctx.emm_context.emm_cause;
    no_attach_proc.esm_msg_out                  = NULL;
    OAILOG_ERROR_UE(
        LOG_NAS_EMM, ue_ctx.emm_context._imsi64,
        "EMM-PROC - Sending Attach Reject for ue_id = " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    rc = _emm_attach_reject(
        &ue_ctx.emm_context, (struct nas_base_proc_s*) &no_attach_proc);
    increment_counter(
        "ue_attach", 1, 2, "result", "failure", "cause", "emergency_attach");
    if (ies) {
      free_emm_attach_request_ies((emm_attach_request_ies_t * * const) & ies);
    }
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  /*
   * Get the UE's EMM context if it exists
   */
  mme_app_desc_t* mme_app_desc_p = get_mme_nas_state(false);
  ue_mm_context                  = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);
  if (!ue_mm_context) {
    OAILOG_ERROR_UE(
        LOG_NAS_EMM, ue_ctx.emm_context._imsi64,
        "EMM-PROC - Sending Attach Reject for ue_id = " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    struct nas_emm_attach_proc_s no_attach_proc = {0};
    no_attach_proc.ue_id                        = ue_id;
    no_attach_proc.emm_cause                    = ue_ctx.emm_context.emm_cause;
    no_attach_proc.esm_msg_out                  = NULL;
    ue_ctx.emm_context.emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
    rc                           = _emm_attach_reject(
        &ue_ctx.emm_context, (struct nas_base_proc_s*) &no_attach_proc);
    increment_counter(
        "ue_attach", 1, 2, "result", "failure", "cause",
        "ue_context_not_found");
    if (ies) {
      free_emm_attach_request_ies((emm_attach_request_ies_t * * const) & ies);
    }
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  // if is_mm_ctx_new==TRUE then ue_mm_context should always be not NULL

  // Actually uplink_nas_transport is sent from S1AP task to NAS task without
  // passing by the MME_APP task... Since now UE_MM context and NAS EMM context
  // are tied together, we may change/use/split logic across MME_APP and NAS

  // Search UE context using GUTI -
  if (ies->guti) {  // no need for  && (is_native_guti)
    guti_ue_mm_ctx =
        mme_ue_context_exists_guti(&mme_app_desc_p->mme_ue_contexts, ies->guti);
    // Allocate new context and process the new request as fresh attach
    // request
    if (guti_ue_mm_ctx) {
      create_new_attach_info(
          &guti_ue_mm_ctx->emm_context, ue_mm_context->mme_ue_s1ap_id,
          STOLEN_REF ies, is_mm_ctx_new);
      /*
       * This implies either UE or eNB has not sent S-TMSI in initial UE
       * message even though UE has old GUTI. Trigger clean up
       */
      nas_proc_implicit_detach_ue_ind(guti_ue_mm_ctx->mme_ue_s1ap_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
    }
    // Allocate new context and process the new request as fresh attach request
    clear_emm_ctxt  = true;
    is_unknown_guti = true;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "EMM-PROC - Received Attach Request with unknown GUTI for ue_id "
        "= " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
  }
  if (ies->imsi) {
    imsi_ue_mm_ctx =
        mme_ue_context_exists_imsi(&mme_app_desc_p->mme_ue_contexts, imsi64);
    if (imsi_ue_mm_ctx) {
      old_ue_id = imsi_ue_mm_ctx->mme_ue_s1ap_id;
      fsm_state = emm_fsm_get_state(&imsi_ue_mm_ctx->emm_context);

      nas_emm_attach_proc_t* attach_proc =
          get_nas_specific_procedure_attach(&imsi_ue_mm_ctx->emm_context);
      if (is_nas_common_procedure_identification_running(
              &imsi_ue_mm_ctx->emm_context)) {
        nas_emm_ident_proc_t* ident_proc =
            get_nas_common_procedure_identification(
                &imsi_ue_mm_ctx->emm_context);
        if (attach_proc) {
          if ((is_nas_attach_accept_sent(attach_proc)) ||
              (is_nas_attach_reject_sent(attach_proc))) {
            REQUIREMENT_3GPP_24_301(R10_5_4_4_6_c);  // continue
            // TODO Need to be reviewed and corrected
            increment_counter(
                "duplicate_attach_request", 1, 1, "action", "not_handled");
          } else {
            REQUIREMENT_3GPP_24_301(R10_5_4_4_6_d);
            emm_sap_t emm_sap           = {0};
            emm_sap.primitive           = EMMREG_COMMON_PROC_ABORT;
            emm_sap.u.emm_reg.ue_id     = ue_id;
            emm_sap.u.emm_reg.ctx       = &imsi_ue_mm_ctx->emm_context;
            emm_sap.u.emm_reg.notify    = false;
            emm_sap.u.emm_reg.free_proc = true;
            emm_sap.u.emm_reg.u.common.common_proc = &ident_proc->emm_com_proc;
            emm_sap.u.emm_reg.u.common.previous_emm_fsm_state =
                ident_proc->emm_com_proc.emm_proc.previous_emm_fsm_state;
            // TODO Need to be reviewed and corrected
            // trigger clean up
            memset(&emm_sap, 0, sizeof(emm_sap));
            emm_sap.primitive = EMMCN_IMPLICIT_DETACH_UE;
            emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = old_ue_id;
            rc = emm_sap_send(&emm_sap);
            // Allocate new context and process the new request as fresh attach
            // request
            clear_emm_ctxt = true;
            increment_counter(
                "duplicate_attach_request", 1, 1, "action",
                "processed_old_ctxt_cleanup");
          }
        } else {
          // TODO Need to be reviewed and corrected
          REQUIREMENT_3GPP_24_301(R10_5_4_4_6_c);  // continue
          increment_counter(
              "duplicate_attach_request", 1, 1, "action", "not_handled");
        }
      }
      if (EMM_REGISTERED == fsm_state) {
        REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_f);
        if (imsi_ue_mm_ctx->emm_context.is_attached) {
          OAILOG_INFO(
              LOG_NAS_EMM,
              "EMM-PROC  - the new ATTACH REQUEST is progressed for ue "
              "id " MME_UE_S1AP_ID_FMT "\n",
              ue_mm_context->mme_ue_s1ap_id);
          // process the new request as fresh attach request
          create_new_attach_info(
              &imsi_ue_mm_ctx->emm_context, ue_mm_context->mme_ue_s1ap_id,
              STOLEN_REF ies, is_mm_ctx_new);
          // Trigger clean up
          nas_proc_implicit_detach_ue_ind(old_ue_id);

          increment_counter(
              "duplicate_attach_request", 1, 1, "action",
              "processed_old_ctxt_cleanup");
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        }
      } else if (
          (attach_proc) &&
          (is_nas_attach_accept_sent(
              attach_proc))) {  // && (!emm_ctx->is_attach_complete_received):
                                // implicit

        imsi_ue_mm_ctx->emm_context.num_attach_request++;
        if (emm_attach_ies_have_changed(
                imsi_ue_mm_ctx->mme_ue_s1ap_id, attach_proc->ies, ies)) {
          OAILOG_WARNING(
              LOG_NAS_EMM,
              "EMM-PROC  - Attach parameters have changed for ue "
              "id " MME_UE_S1AP_ID_FMT "\n",
              imsi_ue_mm_ctx->mme_ue_s1ap_id);
          REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_d__1);
          /*
           * If one or more of the information elements in the ATTACH REQUEST
           * message differ from the ones received within the previous ATTACH
           * REQUEST message, the previously initiated attach procedure shall
           * be aborted if the ATTACH COMPLETE message has not been received
           * and the new attach procedure shall be progressed;
           */
          // After releasing of contexts of old UE, process the new request as
          // fresh attach request
          create_new_attach_info(
              &imsi_ue_mm_ctx->emm_context, ue_mm_context->mme_ue_s1ap_id,
              STOLEN_REF ies, is_mm_ctx_new);

          nas_proc_implicit_detach_ue_ind(old_ue_id);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        } else {
          imsi_ue_mm_ctx->emm_context.num_attach_request++;
          REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_d__2);
          /*
           * - if the information elements do not differ, then the ATTACH ACCEPT
           * message shall be resent and the timer T3450 shall be restarted if
           * an ATTACH COMPLETE message is expected. In that case, the
           * retransmission counter related to T3450 is not incremented.
           */
          emm_attach_accept_retx(&imsi_ue_mm_ctx->emm_context);
          increment_counter(
              "duplicate_attach_request", 1, 1, "action",
              "ignored_duplicate_req_retx_attach_accept");
          if (imsi_ue_mm_ctx->mme_ue_s1ap_id != ue_mm_context->mme_ue_s1ap_id) {
            /* Re-transmitted attach request will be sent in UL nas message
             * and it will have same mme_ue_s1ap_id, so there will not be new
             * contexts created,
             * If Attach Request comes in initial ue message, new
             * mme_ue_s1ap_id and UE contexts will be created,
             * which needs to be deleted
             */
            OAILOG_DEBUG(
                LOG_NAS_EMM,
                "EMM-PROC - Sending Detach Request message to MME APP"
                "module for ue_id =" MME_UE_S1AP_ID_FMT "\n",
                ue_id);
            mme_app_handle_detach_req(ue_mm_context->mme_ue_s1ap_id);
          }
          if (ies) {
            free_emm_attach_request_ies(
                (emm_attach_request_ies_t * * const) & ies);
          }
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        }
      } else if (
          (imsi_ue_mm_ctx) &&
          (0 < imsi_ue_mm_ctx->emm_context.num_attach_request) &&
          ((attach_proc) && ((!is_nas_attach_accept_sent(attach_proc)) &&
                             (!is_nas_attach_reject_sent(attach_proc))))) {
        if (emm_attach_ies_have_changed(
                imsi_ue_mm_ctx->mme_ue_s1ap_id, attach_proc->ies, ies)) {
          OAILOG_WARNING(
              LOG_NAS_EMM,
              "EMM-PROC  - Attach parameters have changed for ue "
              "id " MME_UE_S1AP_ID_FMT "\n",
              imsi_ue_mm_ctx->mme_ue_s1ap_id);
          REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_e__1);
          /*
           * If one or more of the information elements in the ATTACH REQUEST
           * message differs from the ones received within the previous ATTACH
           * REQUEST message, the previously initiated attach procedure shall be
           * aborted and the new attach procedure shall be executed;
           */
          // Allocate new context and process the new request as fresh attach
          // request
          increment_counter(
              "duplicate_attach_request", 1, 1, "action",
              "processed_old_ctxt_cleanup");
          create_new_attach_info(
              &imsi_ue_mm_ctx->emm_context, ue_mm_context->mme_ue_s1ap_id,
              STOLEN_REF ies, is_mm_ctx_new);

          // trigger clean up
          nas_proc_implicit_detach_ue_ind(old_ue_id);
          OAILOG_INFO(
              LOG_NAS_EMM,
              "Sent implicit detach for ue_id " MME_UE_S1AP_ID_FMT "\n",
              ue_mm_context->mme_ue_s1ap_id);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        } else {
          REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_e__2);
          /*
           * if the information elements do not differ, then the network shall
           * continue with the previous attach procedure and shall ignore the
           * second ATTACH REQUEST message.
           */
          // Clean up new UE context that was created to handle new attach
          // request
          OAILOG_DEBUG(
              LOG_NAS_EMM,
              "EMM-PROC - Sending Detach Request message to MME APP"
              "module for ue_id =" MME_UE_S1AP_ID_FMT "\n",
              ue_id);
          /* Release s1 connection only if attach req is received in the
           * initial ue message
           */
          if (ies->is_initial) {
            mme_app_handle_detach_req(ue_mm_context->mme_ue_s1ap_id);
          }

          OAILOG_WARNING(
              LOG_NAS_EMM,
              "EMM-PROC  - Received duplicated Attach Request for ue "
              "id " MME_UE_S1AP_ID_FMT "\n",
              ue_id);
          increment_counter(
              "duplicate_attach_request", 1, 1, "action", "ignored");
          if (ies) {
            free_emm_attach_request_ies(
                (emm_attach_request_ies_t * * const) & ies);
          }
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        }
      }
    }  // if imsi_emm_ctx != NULL
    // Allocate new context and process the new request as fresh attach request
    clear_emm_ctxt = true;
  }

  if (clear_emm_ctxt) {
    /*
     * Create UE's EMM context
     */
    new_emm_ctx = &ue_mm_context->emm_context;

    bdestroy(new_emm_ctx->esm_msg);
    emm_init_context(new_emm_ctx, true);

    new_emm_ctx->num_attach_request++;
    new_emm_ctx->attach_type            = ies->type;
    new_emm_ctx->additional_update_type = ies->additional_update_type;
    OAILOG_NOTICE(
        LOG_NAS_EMM,
        "EMM-PROC  - Create EMM context ue_id = " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    new_emm_ctx->is_dynamic = true;
    new_emm_ctx->emm_cause  = EMM_CAUSE_SUCCESS;
    // Store Voice Domain pref IE to be sent to MME APP
    if (ies->voicedomainpreferenceandueusagesetting) {
      memcpy(
          &new_emm_ctx->volte_params
               .voice_domain_preference_and_ue_usage_setting,
          ies->voicedomainpreferenceandueusagesetting,
          sizeof(voice_domain_preference_and_ue_usage_setting_t));
      new_emm_ctx->volte_params.presencemask |=
          VOICE_DOMAIN_PREF_UE_USAGE_SETTING;
    }
  }
  if (is_unknown_guti) {
    is_unknown_guti                = false;
    new_emm_ctx->emm_context_state = UNKNOWN_GUTI;
  }
  if (!is_nas_specific_procedure_attach_running(&ue_mm_context->emm_context)) {
    emm_proc_create_procedure_attach_request(ue_mm_context, STOLEN_REF ies);
  } else if (ies) {  // we should not be really here
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMM-PROC  - Freeing Attach Request IEs for ue_id "
        "= " MME_UE_S1AP_ID_FMT,
        ue_id);
    free_emm_attach_request_ies((emm_attach_request_ies_t * * const) & ies);
  }
  rc = emm_attach_run_procedure(&ue_mm_context->emm_context);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
/*
 *
 * Name:        emm_proc_attach_reject()
 *
 * Description: Performs the protocol error abnormal case
 *
 *              3GPP TS 24.301, section 5.5.1.2.7, case b
 *              If the ATTACH REQUEST message is received with a protocol
 *              error, the network shall return an ATTACH REJECT message.
 *
 * Inputs:  ue_id:              UE lower layer identifier
 *                  emm_cause: EMM cause code to be reported
 *                  Others:    None
 *
 * Outputs:     None
 *                  Return:    RETURNok, RETURNerror
 *                  Others:    _emm_data
 *
 */
//------------------------------------------------------------------------------
status_code_e emm_proc_attach_reject(
    mme_ue_s1ap_id_t ue_id, emm_cause_t emm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  emm_context_t* emm_ctx = emm_context_get(&_emm_data, ue_id);
  if (emm_ctx) {
    if (is_nas_specific_procedure_attach_running(emm_ctx)) {
      nas_emm_attach_proc_t* attach_proc =
          (nas_emm_attach_proc_t*) (emm_ctx->emm_procedures->emm_specific_proc);
      attach_proc->emm_cause = emm_cause;

      // TODO could be in callback of attach procedure triggered by
      // EMMREG_ATTACH_REJ
      rc = _emm_attach_reject(emm_ctx, (struct nas_base_proc_s*) attach_proc);
      emm_sap_t emm_sap               = {0};
      emm_sap.primitive               = EMMREG_ATTACH_REJ;
      emm_sap.u.emm_reg.ue_id         = ue_id;
      emm_sap.u.emm_reg.ctx           = emm_ctx;
      emm_sap.u.emm_reg.notify        = false;
      emm_sap.u.emm_reg.free_proc     = true;
      emm_sap.u.emm_reg.u.attach.proc = attach_proc;
      rc                              = emm_sap_send(&emm_sap);
    } else {
      nas_emm_attach_proc_t no_attach_proc = {0};
      no_attach_proc.ue_id                 = ue_id;
      no_attach_proc.emm_cause             = emm_cause;
      no_attach_proc.esm_msg_out           = NULL;
      rc                                   = _emm_attach_reject(
          emm_ctx, (struct nas_base_proc_s*) &no_attach_proc);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
 *
 * Name:    emm_proc_attach_complete()
 *
 * Description: Terminates the attach procedure upon receiving Attach
 *      Complete message from the UE.
 *
 *              3GPP TS 24.301, section 5.5.1.2.4
 *      Upon receiving an ATTACH COMPLETE message, the MME shall
 *      stop timer T3450, enter state EMM-REGISTERED and consider
 *      the GUTI sent in the ATTACH ACCEPT message as valid.
 *
 * Inputs:  ue_id:      UE lower layer identifier
 *      esm_msg_pP:   Activate default EPS bearer context accept
 *             ESM message
 *      Others:    _emm_data
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    _emm_data, T3450
 *
 */
//------------------------------------------------------------------------------
status_code_e emm_proc_attach_complete(
    mme_ue_s1ap_id_t ue_id, const_bstring esm_msg_pP, int emm_cause,
    const nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_mm_context_t* ue_mm_context     = NULL;
  nas_emm_attach_proc_t* attach_proc = NULL;
  int rc                             = RETURNerror;
  emm_sap_t emm_sap                  = {0};
  esm_sap_t esm_sap                  = {0};
  emm_context_t* emm_ctx             = NULL;

  /*
   * Get the UE context
   */
  ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);

  if (ue_mm_context) {
    if (is_nas_specific_procedure_attach_running(&ue_mm_context->emm_context)) {
      attach_proc =
          (nas_emm_attach_proc_t*)
              ue_mm_context->emm_context.emm_procedures->emm_specific_proc;

      /* Process attach complete msg only if T3450 timer is running
       * If its not running it means that implicit detach is in progress
       */
      if (attach_proc->T3450.id == NAS_TIMER_INACTIVE_ID) {
        OAILOG_WARNING_UE(
            LOG_NAS_EMM, ue_mm_context->emm_context._imsi64,
            "Discarding attach complete as T3450 timer is not active for "
            "ueid " MME_UE_S1AP_ID_FMT "\n",
            ue_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      }
      emm_ctx = &ue_mm_context->emm_context;
      /*
       * Upon receiving an ATTACH COMPLETE message, the MME shall enter state
       * EMM-REGISTERED and consider the GUTI sent in the ATTACH ACCEPT message
       * as valid.
       */
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__20);
      emm_ctx_set_attribute_valid(emm_ctx, EMM_CTXT_MEMBER_GUTI);
      // TODO LG REMOVE emm_context_add_guti(&_emm_data,
      // &ue_mm_context->emm_context);
      emm_ctx_clear_old_guti(emm_ctx);

      /*
       * send the SGSAP TMSI Reallocation complete message towards SGS.
       * if csfb newTmsiAllocated flag is true
       * After sending set it to false
       */
      if (emm_ctx->csfbparams.newTmsiAllocated) {
        OAILOG_DEBUG(
            LOG_NAS_EMM, " CSFB newTmsiAllocated = (%d) true!\n",
            emm_ctx->csfbparams.newTmsiAllocated);
        char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
        IMSI_TO_STRING(&(emm_ctx->_imsi), imsi_str, IMSI_BCD_DIGITS_MAX + 1);

        OAILOG_INFO(
            LOG_NAS_EMM,
            " Sending SGSAP TMSI REALLOCATION COMPLETE to SGS for ue_id = "
            "(%u)\n",
            ue_id);
        mme_app_itti_sgsap_tmsi_reallocation_comp(imsi_str, strlen(imsi_str));
        emm_ctx->csfbparams.newTmsiAllocated = false;
        /* update the neaf flag to false after sending the Tmsi Reallocation
         * Complete message to SGS */
        mme_ue_context_update_ue_sgs_neaf(ue_id, false);
      }

      /*
       * Forward the Activate Default EPS Bearer Context Accept message
       * to the EPS session management sublayer
       */
      /*currently by default Activate Default Bearer Context Accept message was
       * sent in Attach complete Now, modified the code to send the message
       * received in Uplink/esmContainer.
       * third byte of esm message container is a message_type*/
      switch (esm_msg_pP->data[2]) {
        case ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT:
          esm_sap.primitive = ESM_DEFAULT_EPS_BEARER_CONTEXT_ACTIVATE_CNF;
          break;
        case ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REJECT:
          esm_sap.primitive = ESM_DEFAULT_EPS_BEARER_CONTEXT_ACTIVATE_REJ;
          break;
        case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_ACCEPT:
          esm_sap.primitive = ESM_DEDICATED_EPS_BEARER_CONTEXT_ACTIVATE_CNF;
          break;
        case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REJECT:
          esm_sap.primitive = ESM_DEDICATED_EPS_BEARER_CONTEXT_ACTIVATE_REJ;
          break;
        default:
          OAILOG_ERROR(
              LOG_NAS_EMM, "Invalid ESM Message type, value = [%x] \n",
              esm_msg_pP->data[2]);
          break;
      }
      esm_sap.is_standalone = false;
      esm_sap.ue_id         = ue_id;
      esm_sap.recv          = esm_msg_pP;
      esm_sap.ctx           = &ue_mm_context->emm_context;
      rc                    = esm_sap_send(&esm_sap);
    } else {
      NOT_REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__20);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT
          " ATTACH COMPLETE discarded (EMM procedure not found)\n",
          ue_id);
    }
  } else {
    NOT_REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__20);
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT
        " ATTACH COMPLETE discarded (context not found)\n",
        ue_id);
  }

  if ((rc != RETURNerror) && (esm_sap.err == ESM_SAP_SUCCESS)) {
    /*
     * Set the network attachment indicator
     */
    ue_mm_context->emm_context.is_attached = true;
    /*
     * Notify EMM that attach procedure has successfully completed
     */
    emm_sap.primitive               = EMMREG_ATTACH_CNF;
    emm_sap.u.emm_reg.ue_id         = ue_id;
    emm_sap.u.emm_reg.ctx           = &ue_mm_context->emm_context;
    emm_sap.u.emm_reg.notify        = true;
    emm_sap.u.emm_reg.free_proc     = true;
    emm_sap.u.emm_reg.u.attach.proc = attach_proc;
    rc                              = emm_sap_send(&emm_sap);
    if (rc == RETURNok) {
      /*
       * Send EMM Information after handling Attach Complete message
       * */
      OAILOG_INFO(
          LOG_NAS_EMM,
          " Sending EMM INFORMATION for ue_id = " MME_UE_S1AP_ID_FMT "\n",
          ue_id);
      emm_proc_emm_information(ue_mm_context);
      increment_counter("ue_attach", 1, 1, "result", "attach_proc_successful");
      attach_success_event(ue_mm_context->emm_context._imsi64);
    }
  } else if (esm_sap.err != ESM_SAP_DISCARDED) {
    /*
     * Notify EMM that attach procedure failed
     */
    emm_sap.primitive               = EMMREG_ATTACH_REJ;
    emm_sap.u.emm_reg.ue_id         = ue_id;
    emm_sap.u.emm_reg.ctx           = &ue_mm_context->emm_context;
    emm_sap.u.emm_reg.notify        = true;
    emm_sap.u.emm_reg.free_proc     = true;
    emm_sap.u.emm_reg.u.attach.proc = attach_proc;
    rc                              = emm_sap_send(&emm_sap);
  } else {
    /*
     * ESM procedure failed and, received message has been discarded or
     * Status message has been returned; ignore ESM procedure failure
     */
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "Ignore ESM procedure failure/received message has been discarded for "
        "ue_id = " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
    rc = RETURNok;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/**
 * When the NAS EMM procedures are restored from data store, the references to
 * callback functions need to be re-populated with the local scope. The function
 * below set these callbacks for attach, authentication, scurity mode control
 * and other procedures.
 * The memory for the EMM procedure is allocated by the caller
 */

void set_callbacks_for_attach_proc(nas_emm_attach_proc_t* attach_proc) {
  ((nas_base_proc_t*) attach_proc)->abort   = emm_attach_abort;
  ((nas_base_proc_t*) attach_proc)->fail_in = NULL;
  ((nas_base_proc_t*) attach_proc)->time_out =
      mme_app_handle_emm_attach_t3450_expiry;
  ((nas_base_proc_t*) attach_proc)->fail_out = _emm_attach_reject;
}

void set_notif_callbacks_for_auth_proc(nas_emm_auth_proc_t* auth_proc) {
  auth_proc->emm_com_proc.emm_proc.base_proc.success_notif =
      emm_attach_success_authentication_cb;
  auth_proc->emm_com_proc.emm_proc.base_proc.failure_notif =
      emm_attach_failure_authentication_cb;
}

void set_notif_callbacks_for_smc_proc(nas_emm_smc_proc_t* smc_proc) {
  smc_proc->emm_com_proc.emm_proc.base_proc.success_notif =
      emm_attach_success_security_cb;
  smc_proc->emm_com_proc.emm_proc.base_proc.failure_notif =
      emm_attach_failure_security_cb;
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

static void emm_proc_create_procedure_attach_request(
    ue_mm_context_t* const ue_mm_context,
    STOLEN_REF emm_attach_request_ies_t* const ies) {
  nas_emm_attach_proc_t* attach_proc =
      nas_new_attach_procedure(&ue_mm_context->emm_context);
  AssertFatal(attach_proc, "TODO Handle this");
  if ((attach_proc)) {
    attach_proc->ies                          = ies;
    attach_proc->ue_id                        = ue_mm_context->mme_ue_s1ap_id;
    ((nas_base_proc_t*) attach_proc)->abort   = emm_attach_abort;
    ((nas_base_proc_t*) attach_proc)->fail_in = NULL;  // No parent procedure
    ((nas_base_proc_t*) attach_proc)->time_out =
        mme_app_handle_emm_attach_t3450_expiry;
    ((nas_base_proc_t*) attach_proc)->fail_out = _emm_attach_reject;
  }
}
/*
 * --------------------------------------------------------------------------
 * Timer handlers
 * --------------------------------------------------------------------------
 */

/*
 *
 * Name:    mme_app_handle_emm_attach_t3450_expiry
 *
 * Description: T3450 timeout handler
 *
 *              3GPP TS 24.301, section 5.5.1.2.7, case c
 *      On the first expiry of the timer T3450, the network shall
 *      retransmit the ATTACH ACCEPT message and shall reset and
 *      restart timer T3450. This retransmission is repeated four
 *      times, i.e. on the fifth expiry of timer T3450, the at-
 *      tach procedure shall be aborted and the MME enters state
 *      EMM-DEREGISTERED.
 *
 * Inputs:  args:      handler parameters
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    None
 *      Others:    None
 *
 */
status_code_e mme_app_handle_emm_attach_t3450_expiry(
    zloop_t* loop, int timer_id, void* args) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  mme_ue_s1ap_id_t mme_ue_s1ap_id = 0;
  if (!mme_pop_timer_arg_ue_id(timer_id, &mme_ue_s1ap_id)) {
    OAILOG_WARNING(
        LOG_NAS_EMM, "Invalid Timer Id expiration, Timer Id: %u\n", timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  struct ue_mm_context_s* ue_context_p = mme_app_get_ue_context_for_timer(
      mme_ue_s1ap_id, "Attach Procedure T3450 Timer");
  if (ue_context_p == NULL) {
    OAILOG_ERROR(
        LOG_MME_APP,
        "Invalid UE context received, MME UE S1AP Id: " MME_UE_S1AP_ID_FMT "\n",
        mme_ue_s1ap_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
  }

  emm_context_t* emm_context = &ue_context_p->emm_context;

  if (is_nas_specific_procedure_attach_running(emm_context)) {
    nas_emm_attach_proc_t* attach_proc =
        get_nas_specific_procedure_attach(emm_context);

    attach_proc->T3450.id = NAS_TIMER_INACTIVE_ID;

    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMM-PROC  - T3450 timer expired, retransmission "
        "counter = %d\n",
        attach_proc->attach_accept_sent);
    if (attach_proc->attach_accept_sent < ATTACH_COUNTER_MAX) {
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_c__1);
      /*
       * On the first expiry of the timer, the network shall retransmit the
       * ATTACH ACCEPT message and shall reset and restart timer T3450.
       */
      emm_attach_accept_retx(emm_context);
      attach_proc->attach_accept_sent++;
    } else {
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_c__2);
      /*
       * Abort the attach procedure
       */
      emm_sap_t emm_sap               = {0};
      emm_sap.primitive               = EMMREG_ATTACH_ABORT;
      emm_sap.u.emm_reg.ue_id         = attach_proc->ue_id;
      emm_sap.u.emm_reg.ctx           = emm_context;
      emm_sap.u.emm_reg.notify        = true;
      emm_sap.u.emm_reg.free_proc     = true;
      emm_sap.u.emm_reg.u.attach.proc = attach_proc;
      emm_sap_send(&emm_sap);
      increment_counter("nas_attach_accept_timer_expired", 1, NO_LABELS);
      increment_counter(
          "ue_attach", 1, 2, "result", "failure", "cause",
          "no_response_for_attach_accept");
    }
    // TODO REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_c__3) not coded
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

//------------------------------------------------------------------------------
static int emm_attach_release(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (emm_context) {
    mme_ue_s1ap_id_t ue_id =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id;
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMM-PROC  - Release UE context data (ue_id=" MME_UE_S1AP_ID_FMT ")\n",
        ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
 *
 * Name:    _emm_attach_reject()
 *
 * Description: Performs the attach procedure not accepted by the network.
 *
 *              3GPP TS 24.301, section 5.5.1.2.5
 *      If the attach request cannot be accepted by the network,
 *      the MME shall send an ATTACH REJECT message to the UE in-
 *      including an appropriate EMM cause value.
 *
 * Inputs:  args:      UE context data
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    None
 *
 */
status_code_e _emm_attach_reject(
    emm_context_t* emm_context, struct nas_base_proc_s* nas_base_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  emm_sap_t emm_sap = {0};
  struct nas_emm_attach_proc_s* attach_proc =
      (struct nas_emm_attach_proc_s*) nas_base_proc;

  OAILOG_WARNING(
      LOG_NAS_EMM,
      "EMM-PROC  - EMM attach procedure not accepted "
      "by the network (ue_id=" MME_UE_S1AP_ID_FMT ", cause=%d)\n",
      attach_proc->ue_id, attach_proc->emm_cause);
  /*
   * Notify EMM-AS SAP that Attach Reject message has to be sent
   * onto the network
   */
  emm_sap.primitive                        = EMMAS_ESTABLISH_REJ;
  emm_sap.u.emm_as.u.establish.ue_id       = attach_proc->ue_id;
  emm_sap.u.emm_as.u.establish.eps_id.guti = NULL;

  emm_sap.u.emm_as.u.establish.emm_cause = attach_proc->emm_cause;
  emm_sap.u.emm_as.u.establish.nas_info  = EMM_AS_NAS_INFO_ATTACH;

  if (attach_proc->emm_cause != EMM_CAUSE_ESM_FAILURE) {
    emm_sap.u.emm_as.u.establish.nas_msg = NULL;
  } else if (attach_proc->esm_msg_out) {
    emm_sap.u.emm_as.u.establish.nas_msg = attach_proc->esm_msg_out;
  } else {
    OAILOG_ERROR(LOG_NAS_EMM, "EMM-PROC  - ESM message is missing\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  /*
   * Setup EPS NAS security data
   */
  if (emm_context) {
    emm_as_set_security_data(
        &emm_sap.u.emm_as.u.establish.sctx, &emm_context->_security, false,
        false);
  } else {
    emm_as_set_security_data(
        &emm_sap.u.emm_as.u.establish.sctx, NULL, false, false);
  }
  rc                              = emm_sap_send(&emm_sap);
  attach_proc->attach_reject_sent = true;
  increment_counter("ue_attach", 1, 1, "action", "attach_reject_sent");
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
 *
 * Name:    _emm_attach_abort()
 *
 * Description: Aborts the attach procedure
 *
 * Inputs:  args:      Attach procedure data to be released
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    T3450
 *
 */
//------------------------------------------------------------------------------
static int emm_attach_abort(
    struct emm_context_s* emm_context, struct nas_base_proc_s* base_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);
  if (attach_proc) {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMM-PROC  - Abort the attach procedure (ue_id=" MME_UE_S1AP_ID_FMT
        ")\n",
        attach_proc->ue_id);

    // Trigger clean up
    emm_sap_t emm_sap                               = {0};
    emm_sap.primitive                               = EMMCN_IMPLICIT_DETACH_UE;
    emm_sap.u.emm_cn.u.emm_cn_implicit_detach.ue_id = attach_proc->ue_id;
    rc                                              = emm_sap_send(&emm_sap);
    increment_counter("ue_attach", 1, 1, "action", "attach_abort");
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
 * --------------------------------------------------------------------------
 * Functions that may initiate EMM common procedures
 * --------------------------------------------------------------------------
 */

//------------------------------------------------------------------------------
static int emm_attach_run_procedure(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_3__1);

    if (attach_proc->ies->last_visited_registered_tai)
      emm_ctx_set_valid_lvr_tai(
          emm_context, attach_proc->ies->last_visited_registered_tai);
    emm_ctx_set_valid_ue_nw_cap(
        emm_context, &attach_proc->ies->ue_network_capability);
    if (attach_proc->ies->ms_network_capability) {
      emm_ctx_set_valid_ms_nw_cap(
          emm_context, attach_proc->ies->ms_network_capability);
    }
    emm_context->originating_tai = *attach_proc->ies->originating_tai;

    if (attach_proc->ies->mob_st_clsMark2) {
      emm_ctx_set_mobile_station_clsMark2(
          emm_context, attach_proc->ies->mob_st_clsMark2);
    }
    // temporary choice to clear security context if it exist
    emm_ctx_clear_security(emm_context);

    if (attach_proc->ies->ueadditionalsecuritycapability) {
      emm_ctx_set_ue_additional_security_capability(
          emm_context, attach_proc->ies->ueadditionalsecuritycapability);
    }

    if (attach_proc->ies->imsi) {
      if ((attach_proc->ies->decode_status.mac_matched) ||
          !(attach_proc->ies->decode_status.integrity_protected_message)) {
        // force authentication, even if not necessary
        imsi64_t imsi64 = imsi_to_imsi64(attach_proc->ies->imsi);
        emm_ctx_set_valid_imsi(emm_context, attach_proc->ies->imsi, imsi64);
        emm_context_upsert_imsi(&_emm_data, emm_context);
        rc = emm_start_attach_proc_authentication(emm_context, attach_proc);
        if (rc != RETURNok) {
          OAILOG_ERROR_UE(
              LOG_NAS_EMM, imsi64,
              "Failed to start attach authentication procedure!\n");
        }
      } else {
        // force identification, even if not necessary
        rc = emm_proc_identification(
            emm_context, (nas_emm_proc_t*) attach_proc, IDENTITY_TYPE_2_IMSI,
            emm_attach_success_identification_cb,
            emm_attach_failure_identification_cb);
      }
    } else if (attach_proc->ies->guti) {
      rc = emm_proc_identification(
          emm_context, (nas_emm_proc_t*) attach_proc, IDENTITY_TYPE_2_IMSI,
          emm_attach_success_identification_cb,
          emm_attach_failure_identification_cb);
    } else if (attach_proc->ies->imei) {
      // Emergency attach is not supported
      OAILOG_ERROR(LOG_NAS_EMM, "Emergency attach is not supported");
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_success_identification_cb(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (!emm_context) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "emm_context is NULL in ATTACH - Identification success procedure!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }
  OAILOG_INFO_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Identification procedure success!\n");
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_3__1);
    rc = emm_start_attach_proc_authentication(
        emm_context,
        attach_proc);  //, IDENTITY_TYPE_2_IMSI, _emm_attach_authentified,
                       //_emm_attach_release);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_failure_identification_cb(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (!emm_context) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "emm_context is NULL in ATTACH - Identification failure procedure!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  OAILOG_ERROR_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Identification procedure failed!\n");

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_start_attach_proc_authentication(
    emm_context_t* emm_context, nas_emm_attach_proc_t* attach_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if ((emm_context) && (attach_proc)) {
    rc = emm_proc_authentication(
        emm_context, &attach_proc->emm_spec_proc,
        emm_attach_success_authentication_cb,
        emm_attach_failure_authentication_cb);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_success_authentication_cb(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (!emm_context) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "emm_context is NULL in ATTACH - Authentication success procedure!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  OAILOG_INFO_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Authentication procedure success!\n");
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_3__1);
    rc = emm_start_attach_proc_security(emm_context, attach_proc);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_failure_authentication_cb(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (!emm_context) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "emm_context is NULL in ATTACH - Authentication failure procedure!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  OAILOG_ERROR_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Authentication procedure failed!\n");
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    attach_proc->emm_cause = emm_context->emm_cause;

    emm_sap_t emm_sap               = {0};
    emm_sap.primitive               = EMMREG_ATTACH_REJ;
    emm_sap.u.emm_reg.ue_id         = attach_proc->ue_id;
    emm_sap.u.emm_reg.ctx           = emm_context;
    emm_sap.u.emm_reg.notify        = true;
    emm_sap.u.emm_reg.free_proc     = true;
    emm_sap.u.emm_reg.u.attach.proc = attach_proc;
    // dont' care emm_sap.u.emm_reg.u.attach.is_emergency = false;
    rc = emm_sap_send(&emm_sap);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_start_attach_proc_security(
    emm_context_t* emm_context, nas_emm_attach_proc_t* attach_proc) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if ((emm_context) && (attach_proc)) {
    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_3__1);
    mme_ue_s1ap_id_t ue_id =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id;
    /*
     * Create new NAS security context
     */
    emm_ctx_clear_security(emm_context);
    rc = emm_proc_security_mode_control(
        emm_context, &attach_proc->emm_spec_proc, attach_proc->ksi,
        emm_attach_success_security_cb, emm_attach_failure_security_cb);
    if (rc != RETURNok) {
      /*
       * Failed to initiate the security mode control procedure
       */
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          "EMM-PROC  - Failed to initiate security mode control procedure\n",
          ue_id);
      attach_proc->emm_cause = EMM_CAUSE_ILLEGAL_UE;
      /*
       * Do not accept the UE to attach to the network
       */
      emm_sap_t emm_sap               = {0};
      emm_sap.primitive               = EMMREG_ATTACH_REJ;
      emm_sap.u.emm_reg.ue_id         = ue_id;
      emm_sap.u.emm_reg.ctx           = emm_context;
      emm_sap.u.emm_reg.notify        = true;
      emm_sap.u.emm_reg.free_proc     = true;
      emm_sap.u.emm_reg.u.attach.proc = attach_proc;
      // dont care emm_sap.u.emm_reg.u.attach.is_emergency = false;
      rc = emm_sap_send(&emm_sap);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_success_security_cb(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (!emm_context) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "emm_context is NULL in ATTACH - Security success procedure!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  OAILOG_INFO_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Security procedure success!\n");
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);
  if (!attach_proc) {
    OAILOG_ERROR_UE(
        LOG_NAS_EMM, emm_context->_imsi64,
        "EMM-PROC  - attach_proc is NULL \n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
  }

  if (emm_context->initiate_identity_after_smc) {
    emm_context->initiate_identity_after_smc = false;
    OAILOG_DEBUG_UE(
        LOG_NAS_EMM, emm_context->_imsi64, "Trigger identity procedure\n");
    rc = emm_proc_identification(
        emm_context, (nas_emm_proc_t*) attach_proc, IDENTITY_TYPE_2_IMEISV,
        emm_attach_identification_after_smc_success_cb,
        emm_attach_failure_identification_cb);

    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  rc = emm_attach(emm_context);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_identification_after_smc_success_cb(
    emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  if (!emm_context) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "emm_context is NULL in ATTACH - Identity procedure after smc "
        "procedure success!\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  OAILOG_INFO_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Identity procedure after smc procedure success!\n");
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    rc = emm_attach(emm_context);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_attach_failure_security_cb(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  OAILOG_ERROR_UE(
      LOG_NAS_EMM, emm_context->_imsi64,
      "ATTACH - Security procedure failed!\n");
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    emm_attach_release(emm_context);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//
//  rc = _emm_start_attach_proc_authentication (emm_context, attach_proc);//,
//  IDENTITY_TYPE_2_IMSI, _emm_attach_authentified, _emm_attach_release);
//
//  if ((emm_context) && (attach_proc)) {
//    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_3__1);
//    mme_ue_s1ap_id_t                        ue_id = PARENT_STRUCT(emm_context,
//    struct ue_mm_context_s, emm_context)->mme_ue_s1ap_id; OAILOG_INFO
//    (LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " EMM-PROC  - Setup NAS
//    security\n", ue_id);
//
//    attach_proc->emm_spec_proc.emm_proc.base_proc.success_notif =
//    _emm_attach_success_authentication_cb;
//    attach_proc->emm_spec_proc.emm_proc.base_proc.failure_notif =
//    _emm_attach_failure_authentication_cb;
//    /*
//     * Create new NAS security context
//     */
//    emm_ctx_clear_security(emm_context);
//
//    /*
//     * Initialize the security mode control procedure
//     */
//    rc = emm_proc_security_mode_control (ue_id, emm_context->auth_ksi,
//                                         _emm_attach, _emm_attach_release);
//
//    if (rc != RETURNok) {
//      /*
//       * Failed to initiate the security mode control procedure
//       */
//      OAILOG_WARNING (LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT "EMM-PROC  -
//      Failed to initiate security mode control procedure\n", ue_id);
//      attach_proc->emm_cause = EMM_CAUSE_ILLEGAL_UE;
//      /*
//       * Do not accept the UE to attach to the network
//       */
//      emm_sap_t emm_sap                      = {0};
//      emm_sap.primitive                      = EMMREG_ATTACH_REJ;
//      emm_sap.u.emm_reg.ue_id                = ue_id;
//      emm_sap.u.emm_reg.ctx                  = emm_context;
//      emm_sap.u.emm_reg.notify               = true;
//      emm_sap.u.emm_reg.free_proc            = true;
//      emm_sap.u.emm_reg.u.attach.attach_proc = attach_proc;
//      // dont care emm_sap.u.emm_reg.u.attach.is_emergency = false;
//      rc = emm_sap_send (&emm_sap);
//    }
//  }
//  OAILOG_FUNC_RETURN (LOG_NAS_EMM, rc);
//}
/*
 *
 * Name:        emm_attach_security_a()
 *
 * Description: Initiates security mode control EMM common procedure.
 *
 * Inputs:          args:      security argument parameters
 *                  Others:    None
 *
 * Outputs:     None
 *                  Return:    RETURNok, RETURNerror
 *                  Others:    _emm_data
 *
 */
//------------------------------------------------------------------------------
status_code_e emm_attach_security(struct emm_context_s* emm_context) {
  return emm_attach_security_a(emm_context);
}

//------------------------------------------------------------------------------
static int emm_attach_security_a(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_3__1);
    mme_ue_s1ap_id_t ue_id =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id;
    OAILOG_INFO(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT " EMM-PROC  - Setup NAS security\n", ue_id);

    /*
     * Create new NAS security context
     */
    emm_ctx_clear_security(emm_context);
    /*
     * Initialize the security mode control procedure
     */
    rc = emm_proc_security_mode_control(
        emm_context, &attach_proc->emm_spec_proc, attach_proc->ksi, emm_attach,
        emm_attach_release);

    if (rc != RETURNok) {
      /*
       * Failed to initiate the security mode control procedure
       */
      OAILOG_WARNING(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          "EMM-PROC  - Failed to initiate security mode control procedure\n",
          ue_id);
      attach_proc->emm_cause = EMM_CAUSE_ILLEGAL_UE;
      /*
       * Do not accept the UE to attach to the network
       */
      emm_sap_t emm_sap               = {0};
      emm_sap.primitive               = EMMREG_ATTACH_REJ;
      emm_sap.u.emm_reg.ue_id         = ue_id;
      emm_sap.u.emm_reg.ctx           = emm_context;
      emm_sap.u.emm_reg.notify        = true;
      emm_sap.u.emm_reg.free_proc     = true;
      emm_sap.u.emm_reg.u.attach.proc = attach_proc;
      // dont care emm_sap.u.emm_reg.u.attach.is_emergency = false;
      rc = emm_sap_send(&emm_sap);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
   --------------------------------------------------------------------------
                MME specific local functions
   --------------------------------------------------------------------------
*/

/*
 *
 * Name:    _emm_attach()
 *
 * Description: Performs the attach signalling procedure while a context
 *      exists for the incoming UE in the network.
 *
 *              3GPP TS 24.301, section 5.5.1.2.4
 *      Upon receiving the ATTACH REQUEST message, the MME shall
 *      send an ATTACH ACCEPT message to the UE and start timer
 *      T3450.
 *
 * Inputs:  args:      attach argument parameters
 *      Others:    None
 *
 * Outputs:     None
 *      Return:    RETURNok, RETURNerror
 *      Others:    _emm_data
 *
 */
//------------------------------------------------------------------------------
static int emm_attach(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;
  mme_ue_s1ap_id_t ue_id =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
          ->mme_ue_s1ap_id;

  OAILOG_INFO(
      LOG_NAS_EMM, "ue_id=" MME_UE_S1AP_ID_FMT " EMM-PROC  - Attach UE \n",
      ue_id);

  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    if (attach_proc->ies->esm_msg) {
      esm_sap_t esm_sap     = {0};
      esm_sap.primitive     = ESM_UNITDATA_IND;
      esm_sap.is_standalone = false;
      esm_sap.ue_id         = ue_id;
      esm_sap.ctx           = emm_context;
      esm_sap.recv          = attach_proc->ies->esm_msg;
      rc                    = esm_sap_send(&esm_sap);
      if ((rc != RETURNerror) && (esm_sap.err == ESM_SAP_SUCCESS)) {
        rc = RETURNok;
      } else if (esm_sap.err != ESM_SAP_DISCARDED) {
        /*
         * The attach procedure failed due to an ESM procedure failure
         */
        attach_proc->emm_cause = EMM_CAUSE_ESM_FAILURE;

        /*
         * Setup the ESM message container to include PDN Connectivity Reject
         * message within the Attach Reject message
         */
        bdestroy_wrapper(&attach_proc->ies->esm_msg);
        attach_proc->esm_msg_out = esm_sap.send;
        OAILOG_ERROR(
            LOG_NAS_EMM,
            "Sending Attach Reject to UE for ue_id = " MME_UE_S1AP_ID_FMT
            ", emm_cause = (%d)\n",
            ue_id, attach_proc->emm_cause);
        rc = _emm_attach_reject(
            emm_context, &attach_proc->emm_spec_proc.emm_proc.base_proc);
      } else {
        /*
         * ESM procedure failed and, received message has been discarded or
         * Status message has been returned; ignore ESM procedure failure
         */
        OAILOG_WARNING(
            LOG_NAS_EMM,
            "Ignore ESM procedure failure &"
            "received message has been discarded for ue_id "
            "= " MME_UE_S1AP_ID_FMT "\n",
            ue_id);
        rc = RETURNok;
      }
    } else {
      rc = emm_send_attach_accept(emm_context);
    }
  }

  if (rc != RETURNok) {
    /*
     * The attach procedure failed
     */
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT
        " EMM-PROC  - Failed to respond to Attach Request\n",
        ue_id);
    attach_proc->emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
    /*
     * Do not accept the UE to attach to the network
     */
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "Sending Attach Reject to UE ue_id = " MME_UE_S1AP_ID_FMT
        ", emm_cause = (%d)\n",
        ue_id, attach_proc->emm_cause);
    rc = _emm_attach_reject(
        emm_context, &attach_proc->emm_spec_proc.emm_proc.base_proc);
    increment_counter(
        "ue_attach", 1, 2, "result", "failure", "cause", "protocol_error");
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:   _encode_csfb_parameters_attach_accept                          **
 **                                                                        **
 ** Description: Encode CSFB parameters to send in ATTACH ACCEPT           **
 **                                                                        **
 ** Inputs:  data:      EMM data context, emm_as_establish_t               **
 **      Others:    None                                                   **
 ** Outputs:     None                                                      **
 **      Return:    NONE                                                   **
 **                                                                        **
 ***************************************************************************/

static void encode_csfb_parameters_attach_accept(
    emm_context_t* emm_ctx, emm_as_establish_t* establish_p) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  ue_mm_context_t* ue_mm_context_p =
      PARENT_STRUCT(emm_ctx, struct ue_mm_context_s, emm_context);
  OAILOG_DEBUG(LOG_NAS_EMM, "Encoding CSFB parameters\n");

  char* non_eps_service_control = bdata(mme_config.non_eps_service_control);
  if ((emm_ctx->attach_type == EMM_ATTACH_TYPE_COMBINED_EPS_IMSI) &&
      ((!(strcmp(non_eps_service_control, "SMS")) ||
        !(strcmp(non_eps_service_control, "CSFB_SMS"))))) {
    // CSFB - check if Network Access Mode is Packet only received from HSS in
    // ULA message
    if (is_mme_ue_context_network_access_mode_packet_only(ue_mm_context_p)) {
      establish_p->emm_cause = EMM_CAUSE_CS_SERVICE_NOT_AVAILABLE;
    } else if (
        emm_ctx->csfbparams.sgs_loc_updt_status ==
        SUCCESS) {  // CSFB - Check if SGS Location update procedure is
                    // successful
      if (emm_ctx->csfbparams.presencemask & LAI_CSFB) {
        establish_p->location_area_identification = &emm_ctx->csfbparams.lai;
      }
      // CSFB-Encode Mobile Identity
      if (emm_ctx->csfbparams.presencemask & MOBILE_IDENTITY) {
        establish_p->ms_identity = &emm_ctx->csfbparams.mobileid;
        OAILOG_DEBUG(
            LOG_NAS_EMM,
            "TMSI  digit1 %d\n"
            "TMSI  digit2 %d\n"
            "TMSI  digit3 %d\n"
            "TMSI  digit4 %d\n",
            establish_p->ms_identity->tmsi.tmsi[0],
            establish_p->ms_identity->tmsi.tmsi[1],
            establish_p->ms_identity->tmsi.tmsi[2],
            establish_p->ms_identity->tmsi.tmsi[3]);
      }
    } else if (emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) {
      establish_p->emm_cause = emm_ctx->emm_cause;
    }
    /* Adding Additional Update Result if we have received
    additional_update_type in Attach Request
    or if MME is configures to support SMS only*/
    if ((emm_ctx->additional_update_type == SMS_ONLY) ||
        (emm_ctx->csfbparams.additional_updt_res ==
         ADDITONAL_UPDT_RES_SMS_ONLY)) {
      establish_p->additional_update_result =
          &emm_ctx->csfbparams.additional_updt_res;
    }
    OAILOG_DEBUG(
        LOG_NAS_EMM, "Additional update type = (%u)\n",
        emm_ctx->additional_update_type);
  }
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

//------------------------------------------------------------------------------
status_code_e emm_cn_wrapper_attach_accept(emm_context_t* emm_context) {
  return emm_send_attach_accept(emm_context);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_send_attach_accept()                                      **
 **                                                                        **
 ** Description: Sends ATTACH ACCEPT message and start timer T3450         **
 **                                                                        **
 ** Inputs:  data:      Attach accept retransmission data          **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    T3450                                      **
 **                                                                        **
 ***************************************************************************/
static int emm_send_attach_accept(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  // may be caused by timer not stopped when deleted context
  if (emm_context) {
    emm_sap_t emm_sap = {0};
    nas_emm_attach_proc_t* attach_proc =
        get_nas_specific_procedure_attach(emm_context);
    ue_mm_context_t* ue_mm_context_p =
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
    mme_ue_s1ap_id_t ue_id = ue_mm_context_p->mme_ue_s1ap_id;

    if (attach_proc) {
      emm_attach_update(emm_context, attach_proc->ies);
      /*
       * Notify EMM-AS SAP that Attach Accept message together with an Activate
       * Default EPS Bearer Context Request message has to be sent to the UE
       */
      emm_sap.primitive = EMMAS_ESTABLISH_CNF;
      emm_sap.u.emm_as.u.establish.puid =
          attach_proc->emm_spec_proc.emm_proc.base_proc.nas_puid;
      emm_sap.u.emm_as.u.establish.ue_id    = ue_id;
      emm_sap.u.emm_as.u.establish.nas_info = EMM_AS_NAS_INFO_ATTACH;

      NO_REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__3);
      bdestroy_wrapper(&ue_mm_context_p->ue_radio_capability);
      //----------------------------------------
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__4);
      emm_ctx_set_attribute_valid(
          emm_context, EMM_CTXT_MEMBER_UE_NETWORK_CAPABILITY_IE);
      emm_ctx_set_attribute_valid(
          emm_context, EMM_CTXT_MEMBER_MS_NETWORK_CAPABILITY_IE);
      //----------------------------------------
      if (attach_proc->ies->drx_parameter) {
        REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__5);
        emm_ctx_set_valid_drx_parameter(
            emm_context, attach_proc->ies->drx_parameter);
      }
      //----------------------------------------
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__9);
      // the set of emm_sap.u.emm_as.u.establish.new_guti is for including the
      // GUTI in the attach accept message
      // ONLY ONE MME NOW NO S10
      if (!IS_EMM_CTXT_PRESENT_GUTI(emm_context)) {
        // Sure it is an unknown GUTI in this MME
        guti_t old_guti = emm_context->_old_guti;
        guti_t guti     = {.gummei.plmn     = {0},
                       .gummei.mme_gid  = 0,
                       .gummei.mme_code = 0,
                       .m_tmsi          = INVALID_M_TMSI};
        clear_guti(&guti);

        rc = mme_api_new_guti(
            &emm_context->_imsi, &old_guti, &guti,
            &emm_context->originating_tai, &emm_context->_tai_list);
        if (RETURNok == rc) {
          emm_ctx_set_guti(emm_context, &guti);
          emm_ctx_set_attribute_valid(emm_context, EMM_CTXT_MEMBER_TAI_LIST);
          //----------------------------------------
          REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__6);
          REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__10);
          memcpy(
              &emm_sap.u.emm_as.u.establish.tai_list, &emm_context->_tai_list,
              sizeof(tai_list_t));
        } else {
          OAILOG_ERROR(
              LOG_NAS_EMM,
              "Failed to assign mme api new guti for ue_id "
              "= " MME_UE_S1AP_ID_FMT "\n",
              ue_id);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
        }
      } else {
        // Set the TAI attributes from the stored context for resends.
        memcpy(
            &emm_sap.u.emm_as.u.establish.tai_list, &emm_context->_tai_list,
            sizeof(tai_list_t));
      }
    }

    emm_sap.u.emm_as.u.establish.eps_id.guti = &emm_context->_guti;

    if (!IS_EMM_CTXT_VALID_GUTI(emm_context) &&
        IS_EMM_CTXT_PRESENT_GUTI(emm_context) &&
        IS_EMM_CTXT_PRESENT_OLD_GUTI(emm_context)) {
      /*
       * Implicit GUTI reallocation;
       * include the new assigned GUTI in the Attach Accept message
       */
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          " EMM-PROC  - Implicit GUTI reallocation, include the new assigned "
          "GUTI in the Attach Accept message\n",
          ue_id);
      emm_sap.u.emm_as.u.establish.new_guti = &emm_context->_guti;
    } else if (
        !IS_EMM_CTXT_VALID_GUTI(emm_context) &&
        IS_EMM_CTXT_PRESENT_GUTI(emm_context)) {
      /*
       * include the new assigned GUTI in the Attach Accept message
       */
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          " EMM-PROC  - Include the new assigned GUTI in the Attach Accept "
          "message\n",
          ue_id);
      emm_sap.u.emm_as.u.establish.new_guti = &emm_context->_guti;
    } else {  // IS_EMM_CTXT_VALID_GUTI(ue_mm_context) is true
      emm_sap.u.emm_as.u.establish.new_guti = NULL;
    }
    //----------------------------------------
    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__14);
    emm_sap.u.emm_as.u.establish.eps_network_feature_support =
        (eps_network_feature_support_t*) &_emm_data.conf
            .eps_network_feature_support;

    /*
     * Delete any preexisting UE radio capabilities, pursuant to
     * GPP 24.310:5.5.1.2.4
     */
    // Note: this is safe from double-free errors because it sets to NULL
    // after freeing, which free treats as a no-op.
    bdestroy_wrapper(&ue_mm_context_p->ue_radio_capability);

    /*
     * Setup EPS NAS security data
     */
    emm_as_set_security_data(
        &emm_sap.u.emm_as.u.establish.sctx, &emm_context->_security, false,
        true);
    emm_sap.u.emm_as.u.establish.encryption =
        emm_context->_security.selected_algorithms.encryption;
    emm_sap.u.emm_as.u.establish.integrity =
        emm_context->_security.selected_algorithms.integrity;
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT " EMM-PROC  - encryption = 0x%X (0x%X)\n",
        ue_id, emm_sap.u.emm_as.u.establish.encryption,
        emm_context->_security.selected_algorithms.encryption);
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT " EMM-PROC  - integrity  = 0x%X (0x%X)\n",
        ue_id, emm_sap.u.emm_as.u.establish.integrity,
        emm_context->_security.selected_algorithms.integrity);
    /*
     * Get the activate default EPS bearer context request message to
     * transfer within the ESM container of the attach accept message
     */
    emm_sap.u.emm_as.u.establish.nas_msg = attach_proc->esm_msg_out;
    OAILOG_TRACE(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT
        " EMM-PROC  - nas_msg  src size = %d nas_msg  dst size = %d \n",
        ue_id, blength(attach_proc->esm_msg_out),
        blength(emm_sap.u.emm_as.u.establish.nas_msg));

    // Send T3402
    emm_sap.u.emm_as.u.establish.t3402 = &mme_config.nas_config.t3402_min;

    // Encode CSFB parameters
    encode_csfb_parameters_attach_accept(
        emm_context, &emm_sap.u.emm_as.u.establish);

    REQUIREMENT_3GPP_24_301(R10_5_5_1_2_4__2);
    rc = emm_sap_send(&emm_sap);

    if (RETURNerror != rc) {
      /*
       * Start T3450 timer
       */
      nas_stop_T3450(attach_proc->ue_id, &attach_proc->T3450);
      nas_start_T3450(
          attach_proc->ue_id, &attach_proc->T3450,
          attach_proc->emm_spec_proc.emm_proc.base_proc.time_out);
      attach_proc->attach_accept_sent++;
    }
  } else {
    OAILOG_WARNING(LOG_NAS_EMM, "ue_mm_context NULL\n");
  }

  increment_counter("ue_attach", 1, 1, "action", "attach_accept_sent");
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:   _encode_csfb_parameters_attach_accept_retx()                   **
 **                                                                        **
 ** Description: Encode CSFB parameters to retransmit in ATTACH ACCEPT     **
 **                                                                        **
 ** Inputs:  data:      EMM data context EMM as data                       **
 **      Others:    None                                                   **
 ** Outputs:     None                                                      **
 **      Return:    NONE                                                   **
 **                                                                        **
 ***************************************************************************/

static void encode_csfb_parameters_attach_accept_retx(
    emm_context_t* emm_ctx, emm_as_data_t* data_p) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  ue_mm_context_t* ue_mm_context_p =
      PARENT_STRUCT(emm_ctx, struct ue_mm_context_s, emm_context);

  if ((emm_ctx->attach_type == EMM_ATTACH_TYPE_COMBINED_EPS_IMSI) &&
      ((_esm_data.conf.features & MME_API_CSFB_SMS_SUPPORTED) ||
       (_esm_data.conf.features & MME_API_SMS_SUPPORTED))) {
    // CSFB - Check if SGS Location update procedure is successful
    if (emm_ctx->csfbparams.sgs_loc_updt_status == SUCCESS) {
      if (emm_ctx->csfbparams.presencemask & LAI_CSFB) {
        data_p->location_area_identification = &emm_ctx->csfbparams.lai;
      }
      // CSFB-Encode Mobile Identity
      if (emm_ctx->csfbparams.presencemask & MOBILE_IDENTITY) {
        data_p->ms_identity = &emm_ctx->csfbparams.mobileid;
      }
    } else if (
        (emm_ctx->csfbparams.sgs_loc_updt_status == FAILURE) ||
        is_mme_ue_context_network_access_mode_packet_only(ue_mm_context_p)) {
      data_p->emm_cause = (uint32_t*) &emm_ctx->emm_cause;
    }
    if (emm_ctx->csfbparams.additional_updt_res == SMS_ONLY) {
      data_p->additional_update_result =
          &emm_ctx->csfbparams.additional_updt_res;
    }
  }
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_attach_accept_retx()                                     **
 **                                                                        **
 ** Description: Retransmit ATTACH ACCEPT message and restart timer T3450  **
 **                                                                        **
 ** Inputs:  data:      Attach accept retransmission data                  **
 **      Others:    None                                                   **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    T3450                                                  **
 **                                                                        **
 ***************************************************************************/
static int emm_attach_accept_retx(emm_context_t* emm_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  emm_sap_t emm_sap = {0};
  int rc            = RETURNerror;

  if (!emm_context) {
    OAILOG_WARNING(LOG_NAS_EMM, "emm_ctx NULL\n");
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
  }

  ue_mm_context_t* ue_mm_context_p =
      PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context);
  mme_ue_s1ap_id_t ue_id = ue_mm_context_p->mme_ue_s1ap_id;
  nas_emm_attach_proc_t* attach_proc =
      get_nas_specific_procedure_attach(emm_context);

  if (attach_proc) {
    if (!IS_EMM_CTXT_PRESENT_GUTI(emm_context)) {
      OAILOG_WARNING(
          LOG_NAS_EMM,
          " No GUTI present in emm_ctx. Abormal case. Skipping Retx of Attach "
          "Accept NULL for " MME_UE_S1AP_ID_FMT "\n",
          ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
    }
    /*
     * Notify EMM-AS SAP that Attach Accept message together with an Activate
     * Default EPS Bearer Context Request message has to be sent to the UE.
     * Retx of Attach Accept needs to be done via DL NAS Transport S1AP message
     */
    emm_sap.primitive                = EMMAS_DATA_REQ;
    emm_sap.u.emm_as.u.data.ue_id    = ue_id;
    emm_sap.u.emm_as.u.data.nas_info = EMM_AS_NAS_DATA_ATTACH_ACCEPT;
    memcpy(
        &emm_sap.u.emm_as.u.data.tai_list, &emm_context->_tai_list,
        sizeof(tai_list_t));
    emm_sap.u.emm_as.u.data.eps_id.guti = &emm_context->_guti;
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT
        " EMM-PROC  - Include the same GUTI in the Attach Accept Retx "
        "message\n",
        ue_id);
    emm_sap.u.emm_as.u.establish.eps_network_feature_support =
        (eps_network_feature_support_t*) &_emm_data.conf
            .eps_network_feature_support;
    emm_sap.u.emm_as.u.data.new_guti = &emm_context->_guti;

    /*
     * Setup EPS NAS security data
     */
    emm_as_set_security_data(
        &emm_sap.u.emm_as.u.data.sctx, &emm_context->_security, false, true);
    emm_sap.u.emm_as.u.data.encryption =
        emm_context->_security.selected_algorithms.encryption;
    emm_sap.u.emm_as.u.data.integrity =
        emm_context->_security.selected_algorithms.integrity;
    /*
     * Get the activate default EPS bearer context request message to
     * transfer within the ESM container of the attach accept message
     */
    emm_sap.u.emm_as.u.data.nas_msg = attach_proc->esm_msg_out;
    OAILOG_TRACE(
        LOG_NAS_EMM,
        "ue_id=" MME_UE_S1AP_ID_FMT
        " EMM-PROC  - nas_msg  src size = %d nas_msg  dst size = %d \n",
        ue_id, blength(attach_proc->esm_msg_out),
        blength(emm_sap.u.emm_as.u.data.nas_msg));

    // Encode CSFB parameters
    encode_csfb_parameters_attach_accept_retx(
        emm_context, &emm_sap.u.emm_as.u.data);

    rc = emm_sap_send(&emm_sap);

    if (RETURNerror != rc) {
      OAILOG_INFO(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          " EMM-PROC  -Sent Retx Attach Accept message\n",
          ue_id);
      /*
       * Re-start T3450 timer
       */
      nas_stop_T3450(ue_id, &attach_proc->T3450);
      nas_start_T3450(
          ue_id, &attach_proc->T3450,
          attach_proc->emm_spec_proc.emm_proc.base_proc.time_out);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          " EMM-PROC  T3450"
          " restarted\n",
          attach_proc->ue_id);
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Timer T3450 %ld expires in %u seconds\n",
          attach_proc->ue_id, attach_proc->T3450.id, attach_proc->T3450.msec);
    } else {
      OAILOG_WARNING(
          LOG_NAS_EMM,
          "ue_id=" MME_UE_S1AP_ID_FMT
          " EMM-PROC  - Send failed- Retx Attach Accept message\n",
          ue_id);
    }
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
 * Description: Check whether the given attach parameters differs from
 *      those previously stored when the attach procedure has
 *      been initiated.
 *
 * Outputs:     None
 *      Return:    true if at least one of the parameters
 *             differs; false otherwise.
 *      Others:    None
 *
 */
//-----------------------------------------------------------------------------
static bool emm_attach_ies_have_changed(
    mme_ue_s1ap_id_t ue_id, emm_attach_request_ies_t* const ies1,
    emm_attach_request_ies_t* const ies2) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  if (ies1->type != ies2->type) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: type EMM_ATTACH_TYPE\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }
  if (ies1->is_native_sc != ies2->is_native_sc) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT
        " Attach IEs changed: Is native securitty context\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }
  if (ies1->ksi != ies2->ksi) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: KSI %d -> %d \n", ue_id,
        ies1->ksi, ies2->ksi);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  /*
   * The GUTI if provided by the UE
   */
  if (ies1->is_native_guti != ies2->is_native_guti) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: Native GUTI %d -> %d \n",
        ue_id, ies1->is_native_guti, ies2->is_native_guti);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }
  if ((ies1->guti) && (!ies2->guti)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed:  GUTI " GUTI_FMT
        " -> None\n",
        ue_id, GUTI_ARG(ies1->guti));
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->guti) && (ies2->guti)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed:  GUTI None ->  " GUTI_FMT
        "\n",
        ue_id, GUTI_ARG(ies2->guti));
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->guti) && (ies2->guti)) {
    if (memcmp(ies1->guti, ies2->guti, sizeof(*(ies1->guti)))) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed:  guti/tmsi " GUTI_FMT
          " -> " GUTI_FMT "\n",
          ue_id, GUTI_ARG(ies1->guti), GUTI_ARG(ies2->guti));
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }

  /*
   * The IMSI if provided by the UE
   */
  if ((ies1->imsi) && (!ies2->imsi)) {
    imsi64_t imsi641 = imsi_to_imsi64(ies1->imsi);
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed:  IMSI " IMSI_64_FMT
        " -> None\n",
        ue_id, imsi641);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->imsi) && (ies2->imsi)) {
    imsi64_t imsi642 = imsi_to_imsi64(ies2->imsi);
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT
        " Attach IEs changed:  IMSI None ->  " IMSI_64_FMT "\n",
        ue_id, imsi642);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->guti) && (ies2->guti)) {
    imsi64_t imsi641 = imsi_to_imsi64(ies1->imsi);
    imsi64_t imsi642 = imsi_to_imsi64(ies2->imsi);
    if (memcmp(ies1->guti, ies2->guti, sizeof(*(ies1->guti)))) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed:  IMSI " IMSI_64_FMT
          " -> " IMSI_64_FMT "\n",
          ue_id, imsi641, imsi642);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }

  /*
   * The IMEI if provided by the UE
   */
  if ((ies1->imei) && (!ies2->imei)) {
    char imei_str[16];

    IMEI_TO_STRING(ies1->imei, imei_str, 16);
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: imei %s/NULL (ctxt)\n",
        ue_id, imei_str);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->imei) && (ies2->imei)) {
    char imei_str[16];

    IMEI_TO_STRING(ies2->imei, imei_str, 16);
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: imei NULL/%s (ctxt)\n",
        ue_id, imei_str);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->imei) && (ies2->imei)) {
    if (memcmp(ies1->imei, ies2->imei, sizeof(*(ies2->imei))) != 0) {
      char imei_str[16];
      char imei2_str[16];

      IMEI_TO_STRING(ies1->imei, imei_str, 16);
      IMEI_TO_STRING(ies2->imei, imei2_str, 16);
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: imei %s/%s (ctxt)\n",
          ue_id, imei_str, imei2_str);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }

  /*
   * The Last visited registered TAI if provided by the UE
   */
  if ((ies1->last_visited_registered_tai) &&
      (!ies2->last_visited_registered_tai)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: LVR TAI " TAI_FMT
        "/NULL\n",
        ue_id, TAI_ARG(ies1->last_visited_registered_tai));
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->last_visited_registered_tai) &&
      (ies2->last_visited_registered_tai)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: LVR TAI NULL/" TAI_FMT
        "\n",
        ue_id, TAI_ARG(ies2->last_visited_registered_tai));
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->last_visited_registered_tai) &&
      (ies2->last_visited_registered_tai)) {
    if (memcmp(
            ies1->last_visited_registered_tai,
            ies2->last_visited_registered_tai,
            sizeof(*(ies2->last_visited_registered_tai))) != 0) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: LVR TAI " TAI_FMT
          "/" TAI_FMT "\n",
          ue_id, TAI_ARG(ies1->last_visited_registered_tai),
          TAI_ARG(ies2->last_visited_registered_tai));
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }

  /*
   * Originating TAI
   */
  if ((ies1->originating_tai) && (!ies2->originating_tai)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: orig TAI " TAI_FMT
        "/NULL\n",
        ue_id, TAI_ARG(ies1->originating_tai));
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->originating_tai) && (ies2->originating_tai)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: orig TAI NULL/" TAI_FMT
        "\n",
        ue_id, TAI_ARG(ies2->originating_tai));
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->originating_tai) && (ies2->originating_tai)) {
    if (memcmp(
            ies1->originating_tai, ies2->originating_tai,
            sizeof(*(ies2->originating_tai))) != 0) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: orig TAI " TAI_FMT
          "/" TAI_FMT "\n",
          ue_id, TAI_ARG(ies1->originating_tai),
          TAI_ARG(ies2->originating_tai));
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }

  /*
   * Originating ECGI
   */
  if ((ies1->originating_ecgi) && (!ies2->originating_ecgi)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: orig ECGI\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->originating_ecgi) && (ies2->originating_ecgi)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: orig ECGI\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->originating_ecgi) && (ies2->originating_ecgi)) {
    if (memcmp(
            ies1->originating_ecgi, ies2->originating_ecgi,
            sizeof(*(ies2->originating_ecgi))) != 0) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: orig ECGI\n", ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }

  /*
   * UE network capability
   */
  if (memcmp(
          &ies1->ue_network_capability, &ies2->ue_network_capability,
          sizeof(ies1->ue_network_capability))) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: UE network capability\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  /*
   * MS network capability
   */
  if ((ies1->ms_network_capability) && (!ies2->ms_network_capability)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: MS network capability\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->ms_network_capability) && (ies2->ms_network_capability)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT " Attach IEs changed: MS network capability\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->ms_network_capability) && (ies2->ms_network_capability)) {
    if (memcmp(
            ies1->ms_network_capability, ies2->ms_network_capability,
            sizeof(*(ies2->ms_network_capability))) != 0) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT
          " Attach IEs changed: MS network capability\n",
          ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }
  /*
   * UE Additional Security Capability
   */
  if ((ies1->ueadditionalsecuritycapability) &&
      (!ies2->ueadditionalsecuritycapability)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT
        " Attach IEs changed: UE additional security capability\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((!ies1->ueadditionalsecuritycapability) &&
      (ies2->ueadditionalsecuritycapability)) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "UE " MME_UE_S1AP_ID_FMT
        " Attach IEs changed: UE additional security capability\n",
        ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
  }

  if ((ies1->ueadditionalsecuritycapability) &&
      (ies2->ueadditionalsecuritycapability)) {
    if (memcmp(
            ies1->ueadditionalsecuritycapability,
            ies2->ueadditionalsecuritycapability,
            sizeof(*(ies2->ueadditionalsecuritycapability))) != 0) {
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "UE " MME_UE_S1AP_ID_FMT
          " Attach IEs changed: UE additional security capability\n",
          ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, true);
    }
  }
  // TODO ESM MSG ?

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, false);
}

//------------------------------------------------------------------------------
void free_emm_attach_request_ies(emm_attach_request_ies_t** const ies) {
  if ((*ies)->guti) {
    free_wrapper((void**) &(*ies)->guti);
  }
  if ((*ies)->imsi) {
    free_wrapper((void**) &(*ies)->imsi);
  }
  if ((*ies)->imei) {
    free_wrapper((void**) &(*ies)->imei);
  }
  if ((*ies)->last_visited_registered_tai) {
    free_wrapper((void**) &(*ies)->last_visited_registered_tai);
  }
  if ((*ies)->originating_tai) {
    free_wrapper((void**) &(*ies)->originating_tai);
  }
  if ((*ies)->originating_ecgi) {
    free_wrapper((void**) &(*ies)->originating_ecgi);
  }
  if ((*ies)->ms_network_capability) {
    free_wrapper((void**) &(*ies)->ms_network_capability);
  }
  if ((*ies)->esm_msg) {
    bdestroy_wrapper(&(*ies)->esm_msg);
  }
  if ((*ies)->drx_parameter) {
    free_wrapper((void**) &(*ies)->drx_parameter);
  }
  if ((*ies)->mob_st_clsMark2) {
    free_wrapper((void**) &(*ies)->mob_st_clsMark2);
  }
  if ((*ies)->voicedomainpreferenceandueusagesetting) {
    free_wrapper((void**) &(*ies)->voicedomainpreferenceandueusagesetting);
  }
  if ((*ies)->ueadditionalsecuritycapability) {
    free_wrapper((void**) &(*ies)->ueadditionalsecuritycapability);
  }
  free_wrapper((void**) ies);
}

/*
  Name:    _emm_attach_update()

  Description: Update the EMM context with the given attach procedure
       parameters.

  Inputs:  ue_id:      UE lower layer identifier
       type:      Type of the requested attach
       ksi:       Security ket sey identifier
       guti:      The GUTI provided by the UE
       imsi:      The IMSI provided by the UE
       imei:      The IMEI provided by the UE
       eea:       Supported EPS encryption algorithms
       originating_tai Originating TAI (from eNB TAI)
       eia:       Supported EPS integrity algorithms
       esm_msg_pP:   ESM message contained with the attach re-
              quest
       Others:    None
  Outputs:     ctx:       EMM context of the UE in the network
       Return:    RETURNok, RETURNerror
       Others:    None

 */
//------------------------------------------------------------------------------
static int emm_attach_update(
    emm_context_t* const emm_context, emm_attach_request_ies_t* const ies) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * Emergency bearer services indicator
   */
  emm_context->is_emergency = (ies->type == EMM_ATTACH_TYPE_EMERGENCY);
  /*
   * Security key set identifier
   */
  if (emm_context->ksi != ies->ksi) {
    OAILOG_TRACE(
        LOG_NAS_EMM,
        "UE id " MME_UE_S1AP_ID_FMT
        " Update ue ksi %d "
        "-> %d\n",
        PARENT_STRUCT(emm_context, struct ue_mm_context_s, emm_context)
            ->mme_ue_s1ap_id,
        emm_context->ksi, ies->ksi);
    emm_context->ksi = ies->ksi;
  }
  /*
   * Supported EPS encryption algorithms
   */
  emm_ctx_set_valid_ue_nw_cap(emm_context, &ies->ue_network_capability);

  if (ies->ms_network_capability) {
    emm_ctx_set_valid_ms_nw_cap(emm_context, ies->ms_network_capability);
  } else {
    // optional IE
    emm_ctx_clear_ms_nw_cap(emm_context);
  }

  emm_context->originating_tai      = *ies->originating_tai;
  emm_context->is_guti_based_attach = false;

  /*
   * The GUTI if provided by the UE. Trigger UE Identity Procedure to fetch IMSI
   */
  if (ies->guti) {
    emm_context->is_guti_based_attach = true;
  }
  /*
   * The IMSI if provided by the UE
   */
  if (ies->imsi) {
    imsi64_t new_imsi64 = imsi_to_imsi64(ies->imsi);
    if (new_imsi64 != emm_context->_imsi64) {
      emm_ctx_set_valid_imsi(emm_context, ies->imsi, new_imsi64);
    }
  }

  /*
   * The IMEI if provided by the UE
   */
  if (ies->imei) {
    emm_ctx_set_valid_imei(emm_context, ies->imei);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
}

void proc_new_attach_req(
    mme_ue_context_t* const mme_ue_context_p,
    struct ue_mm_context_s* ue_context_p) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);

  hashtable_rc_t hash_rc = HASH_TABLE_OK;
  OAILOG_INFO(
      LOG_NAS_EMM,
      "Process new Attach Request for ue_id " MME_UE_S1AP_ID_FMT "\n",
      ue_context_p->mme_ue_s1ap_id);
  new_attach_info_t attach_info = {0};
  memcpy(
      &attach_info, ue_context_p->emm_context.new_attach_info,
      sizeof(new_attach_info_t));
  free_wrapper((void**) &ue_context_p->emm_context.new_attach_info);
  /* The new Attach Request is received in s1ap initial ue message,
   * So release previous Attach Request's contexts
   */
  if (attach_info.is_mm_ctx_new) {
    if (ue_context_p->ecm_state == ECM_IDLE) {
      OAILOG_INFO_UE(
          LOG_NAS_EMM, ue_context_p->emm_context._imsi64,
          "Remove UE context for ue_id " MME_UE_S1AP_ID_FMT
          " as ue is in idle mode \n",
          ue_context_p->mme_ue_s1ap_id);
      mme_remove_ue_context(mme_ue_context_p, ue_context_p);
    } else {
      ue_context_p->ue_context_rel_cause = S1AP_NAS_DETACH;
      /* In case of Ue initiated explicit IMSI Detach or Combined EPS/IMSI
       *  detach Do not send UE Context Release Command to eNB before receiving
       *  SGs IMSI Detach Ack from MSC/VLR
       */
      if (ue_context_p->sgs_context != NULL) {
        if ((ue_context_p->sgs_detach_type ==
             SGS_EXPLICIT_UE_INITIATED_IMSI_DETACH_FROM_NONEPS) ||
            (ue_context_p->sgs_detach_type ==
             SGS_COMBINED_UE_INITIATED_IMSI_DETACH_FROM_EPS_N_NONEPS)) {
          OAILOG_FUNC_OUT(LOG_NAS_EMM);
        } else if (
            ue_context_p->sgs_context->ts9_timer.id ==
            MME_APP_TIMER_INACTIVE_ID) {
          /* Notify S1AP to send UE Context Release Command to eNB or free
           * s1 context locally.
           */
          mme_app_itti_ue_context_release(
              ue_context_p, ue_context_p->ue_context_rel_cause);
        }
      } else {
        // Notify S1AP to send UE Context Release Command to eNB or free s1
        // context locally.
        mme_app_itti_ue_context_release(
            ue_context_p, ue_context_p->ue_context_rel_cause);
      }
      ue_context_p->ue_context_rel_cause = S1AP_INVALID_CAUSE;
    }
  } else {
    uint64_t mme_ue_s1ap_id64 = 0;

    hash_rc = obj_hashtable_uint64_ts_get(
        mme_ue_context_p->guti_ue_context_htbl,
        (const void*) &ue_context_p->emm_context._guti, sizeof(guti_t),
        &mme_ue_s1ap_id64);

    if (HASH_TABLE_OK == hash_rc) {
      // While processing new attach req, remove GUTI from hashtable
      if ((ue_context_p->emm_context._guti.gummei.mme_code) ||
          (ue_context_p->emm_context._guti.gummei.mme_gid) ||
          (ue_context_p->emm_context._guti.m_tmsi) ||
          (ue_context_p->emm_context._guti.gummei.plmn.mcc_digit1) ||
          (ue_context_p->emm_context._guti.gummei.plmn.mcc_digit2) ||
          (ue_context_p->emm_context._guti.gummei.plmn.mcc_digit3)) {
        hash_rc = obj_hashtable_uint64_ts_remove(
            mme_ue_context_p->guti_ue_context_htbl,
            (const void* const) & ue_context_p->emm_context._guti,
            sizeof(ue_context_p->emm_context._guti));
        if (HASH_TABLE_OK != hash_rc)
          OAILOG_ERROR_UE(
              LOG_MME_APP, ue_context_p->emm_context._imsi64,
              "UE Context not found for GUTI " GUTI_FMT " \n",
              GUTI_ARG(&(ue_context_p->emm_context._guti)));
      }
    }
  }

  // Proceed with new attach request
  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(attach_info.mme_ue_s1ap_id);
  emm_context_t* new_emm_ctx = &ue_mm_context->emm_context;
  /* In case of GUTI attach with unknown GUTI, attach procedure is already
     created and identification procedure is also completed.
     So invoke authentication procedure
  */
  if (new_emm_ctx &&
      (new_emm_ctx->emm_context_state == NEW_EMM_CONTEXT_CREATED)) {
    nas_emm_attach_proc_t* attach_proc =
        get_nas_specific_procedure_attach(new_emm_ctx);
    if (attach_proc) {
      /* Upsert IMSI stored in emm context into the hashtable
       * as it will be deleted during implicit detach
       */
      emm_context_upsert_imsi(&_emm_data, new_emm_ctx);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMM-PROC  - Triggering authentication for ue_id "
          "= " MME_UE_S1AP_ID_FMT "\n",
          ue_mm_context->mme_ue_s1ap_id);
      if (emm_start_attach_proc_authentication(new_emm_ctx, attach_proc) !=
          RETURNok) {
        OAILOG_ERROR(
            LOG_NAS_EMM,
            "EMM-PROC  - Failed to start authentication procedure for ue_id "
            "= " MME_UE_S1AP_ID_FMT "\n",
            ue_mm_context->mme_ue_s1ap_id);
      }
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "EMM-PROC  - Attach procedure does not exist for ue_id "
          "= " MME_UE_S1AP_ID_FMT "\n",
          ue_mm_context->mme_ue_s1ap_id);
    }
    new_emm_ctx->emm_context_state = NEW_EMM_CONTEXT_NOT_CREATED;
    OAILOG_FUNC_OUT(LOG_NAS_EMM);
  }

  bdestroy(new_emm_ctx->esm_msg);
  emm_init_context(new_emm_ctx, true);

  new_emm_ctx->num_attach_request++;
  new_emm_ctx->attach_type            = attach_info.ies->type;
  new_emm_ctx->additional_update_type = attach_info.ies->additional_update_type;
  OAILOG_NOTICE(
      LOG_NAS_EMM,
      "EMM-PROC  - Create EMM context ue_id = " MME_UE_S1AP_ID_FMT "\n",
      ue_mm_context->mme_ue_s1ap_id);
  new_emm_ctx->is_dynamic = true;
  new_emm_ctx->emm_cause  = EMM_CAUSE_SUCCESS;
  // Store Voice Domain pref IE to be sent to MME APP
  if (attach_info.ies->voicedomainpreferenceandueusagesetting) {
    memcpy(
        &new_emm_ctx->volte_params.voice_domain_preference_and_ue_usage_setting,
        attach_info.ies->voicedomainpreferenceandueusagesetting,
        sizeof(voice_domain_preference_and_ue_usage_setting_t));
    new_emm_ctx->volte_params.presencemask |=
        VOICE_DOMAIN_PREF_UE_USAGE_SETTING;
  }
  if (!is_nas_specific_procedure_attach_running(&ue_mm_context->emm_context)) {
    emm_proc_create_procedure_attach_request(
        ue_mm_context, STOLEN_REF attach_info.ies);
  }
  emm_attach_run_procedure(&ue_mm_context->emm_context);
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}
