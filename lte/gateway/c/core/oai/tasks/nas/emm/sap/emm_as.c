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
#include <netinet/in.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "common_defs.h"
#include "3gpp_requirements_24.301.h"
#include "common_types.h"
#include "3gpp_24.007.h"
#include "mme_app_ue_context.h"
#include "nas/as_message.h"
#include "emm_cause.h"
#include "mme_app_itti_messaging.h"
#include "emm_as.h"
#include "emm_recv.h"
#include "LowerLayer.h"
#include "emm_send.h"
#include "nas_proc.h"
#include "emm_proc.h"
#include "service303.h"
#include "conversions.h"
#include "3gpp_23.003.h"
#include "3gpp_24.301.h"
#include "3gpp_36.401.h"
#include "AttachRequest.h"
#include "DetachRequest.h"
#include "TrackingAreaIdentity.h"
#include "emm_data.h"
#include "emm_msg.h"
#include "emm_msgDef.h"
#include "mme_api.h"
#include "mme_app_state.h"
#include "nas_message.h"
#include "nas_procedures.h"
#include "mme_app_defs.h"

/****************************************************************************/
/****************  E X T E R N A L    D E F I N I T I O N S  ****************/
/****************************************************************************/

/****************************************************************************/
/*******************  L O C A L    D E F I N I T I O N S  *******************/
/****************************************************************************/

/*
   String representation of EMMAS-SAP primitives
*/
static const char* emm_as_primitive_str[] = {
    "EMMAS_SECURITY_REQ",   "EMMAS_SECURITY_IND",   "EMMAS_SECURITY_RES",
    "EMMAS_SECURITY_REJ",   "EMMAS_ESTABLISH_REQ",  "EMMAS_ESTABLISH_CNF",
    "EMMAS_ESTABLISH_REJ",  "EMMAS_RELEASE_REQ",    "EMMAS_RELEASE_IND",
    "EMMAS_ERAB_SETUP_REQ", "EMMAS_ERAB_SETUP_CNF", "EMMAS_ERAB_SETUP_REJ",
    "EMMAS_DATA_REQ",       "EMMAS_DATA_IND",       "EMMAS_PAGE_IND",
    "EMMAS_STATUS_IND",     "EMMAS_ERAB_REL_CMD",
};

/*
   Functions executed to process EMM procedures upon receiving
   data from the network
*/
static int emm_as_recv(
    mme_ue_s1ap_id_t ue_id, tai_t const* originating_tai,
    ecgi_t const* originating_ecgi, bstring msg, size_t len, int* emm_cause,
    nas_message_decode_status_t* decode_status);

static int emm_as_establish_req(emm_as_establish_t* msg, int* emm_cause);
static int emm_as_data_ind(emm_as_data_t* msg, int* emm_cause);
static int emm_as_release_ind(
    const emm_as_release_t* const release, int* emm_cause);

/*
   Functions executed to send data to the network when requested
   within EMM procedure processing
*/
static EMM_msg* emm_as_set_header(
    nas_message_t* msg, const emm_as_security_data_t* security);

static int emm_as_encode(
    bstring* info, nas_message_t* msg, size_t length,
    emm_security_context_t* emm_security_context);

static int emm_as_encrypt(
    bstring* info, const nas_message_security_header_t* header,
    const unsigned char* buffer, size_t length,
    emm_security_context_t* emm_security_context);

static int emm_as_send_a(const emm_as_t* msg);
static int emm_as_security_req(
    const emm_as_security_t*, dl_info_transfer_req_t*);
static int emm_as_security_rej(
    const emm_as_security_t*, dl_info_transfer_req_t*);
static int emm_as_establish_cnf(
    const emm_as_establish_t*, nas_establish_rsp_t*);
static int emm_as_establish_rej(
    const emm_as_establish_t*, nas_establish_rsp_t*);
static int emm_as_data_req(const emm_as_data_t*, dl_info_transfer_req_t*);
static int emm_as_status_ind(const emm_as_status_t*, dl_info_transfer_req_t*);
static int emm_as_release_req(const emm_as_release_t*, nas_release_req_t*);
static int emm_as_erab_setup_req(
    const emm_as_activate_bearer_context_req_t*,
    activate_bearer_context_req_t*);
static int emm_as_erab_rel_cmd(
    const emm_as_deactivate_bearer_context_req_t* msg,
    deactivate_bearer_context_req_t* as_msg);

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

/****************************************************************************
 **                                                                        **
 ** Name:    emm_as_initialize()                                       **
 **                                                                        **
 ** Description: Initializes the EMMAS Service Access Point                **
 **                                                                        **
 ** Inputs:  None                                                      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    None                                       **
 **      Others:    NONE                                       **
 **                                                                        **
 ***************************************************************************/
void emm_as_initialize(void) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  /*
   * TODO: Initialize the EMMAS-SAP
   */
  OAILOG_FUNC_OUT(LOG_NAS_EMM);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_as_send()                                             **
 **                                                                        **
 ** Description: Processes the EMMAS Service Access Point primitive.       **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
int emm_as_send(emm_as_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc                       = RETURNok;
  int emm_cause                = EMM_CAUSE_SUCCESS;
  emm_as_primitive_t primitive = msg->primitive;
  mme_ue_s1ap_id_t ue_id       = 0;

  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - Received primitive %s (%d)\n",
      emm_as_primitive_str[primitive - _EMMAS_START - 1], primitive);

  switch (primitive) {
    case _EMMAS_DATA_IND:
      rc    = emm_as_data_ind(&msg->u.data, &emm_cause);
      ue_id = msg->u.data.ue_id;
      break;

    case _EMMAS_ESTABLISH_REQ:
      rc    = emm_as_establish_req(&msg->u.establish, &emm_cause);
      ue_id = msg->u.establish.ue_id;
      break;

    case _EMMAS_RELEASE_IND:
      rc    = emm_as_release_ind(&msg->u.release, &emm_cause);
      ue_id = msg->u.release.ue_id;
      break;

    default:
      /*
       * Other primitives are forwarded to lower layers (S1AP)
       */
      rc = emm_as_send_a(msg);

      if (rc != RETURNok) {
        OAILOG_ERROR(
            LOG_NAS_EMM,
            "EMMAS-SAP - "
            "Failed to process primitive %s (%d) for ue id " MME_UE_S1AP_ID_FMT
            "\n",
            emm_as_primitive_str[primitive - _EMMAS_START - 1], primitive,
            ue_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
      }

      break;
  }

  /*
   * Handle decoding errors
   */
  if ((emm_cause != EMM_CAUSE_SUCCESS) &&
      (emm_cause != EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW)) {
    /*
     * Ignore received message that is too short to contain a complete
     * * * * message type information element
     */
    if (rc == TLV_BUFFER_TOO_SHORT) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
    }
    /*
     * Ignore received message that contains not supported protocol
     * * * * discriminator
     */
    else if (rc == TLV_PROTOCOL_NOT_SUPPORTED) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
    } else if (rc == TLV_WRONG_MESSAGE_TYPE) {
      emm_cause = EMM_CAUSE_MESSAGE_TYPE_NOT_IMPLEMENTED;
    }

    /*
     * EMM message processing failed
     */
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMAS-SAP - Received EMM message is not valid "
        "(cause=%d) for ue id " MME_UE_S1AP_ID_FMT "\n",
        emm_cause, ue_id);
    /*
     * Return an EMM status message
     */
    rc = emm_proc_status(ue_id, emm_cause);
  }

  if (rc != RETURNok) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMAS-SAP - Failed to process primitive %s (%d) for ue "
        "id " MME_UE_S1AP_ID_FMT "\n",
        emm_as_primitive_str[primitive - _EMMAS_START - 1], primitive, ue_id);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************/
/*********************  L O C A L    F U N C T I O N S  *********************/
/****************************************************************************/

/*
   --------------------------------------------------------------------------
   Functions executed to process EMM procedures upon receiving data from the
   network
   --------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_recv()                                            **
 **                                                                        **
 ** Description: Decodes and processes the EPS Mobility Management message **
 **      received from the Access Stratum                          **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      msg:       The EMM message to process                 **
 **      len:       The length of the EMM message              **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_cause: EMM cause code                             **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_recv(
    mme_ue_s1ap_id_t ue_id, tai_t const* originating_tai,
    ecgi_t const* originating_ecgi, bstring msg, size_t len, int* emm_cause,
    nas_message_decode_status_t* decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  nas_message_decode_status_t local_decode_status = {0};
  int decoder_rc                                  = RETURNok;
  int rc                                          = RETURNerror;
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};
  emm_security_context_t* emm_security_context =
      NULL; /* Current EPS NAS security context     */

  if (decode_status) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "EMMAS-SAP - Received EMM message (length=%lu)\n integrity protected "
        "%d \n"
        " ciphered %d \n mac matched %d \n security context %d\n",
        len, decode_status->integrity_protected_message,
        decode_status->ciphered_message, decode_status->mac_matched,
        decode_status->security_context_available);
  } else {
    OAILOG_DEBUG(
        LOG_NAS_EMM, "EMMAS-SAP - Received EMM message (length=%lu)\n", len);
  }

  memset(&nas_msg, 0, sizeof(nas_msg));
  if (!decode_status) {
    memset(&local_decode_status, 0, sizeof(local_decode_status));
    decode_status = &local_decode_status;
  }

  ue_mm_context_t* ue_mm_context = mme_ue_context_exists_mme_ue_s1ap_id(ue_id);

  emm_context_t* emm_ctx = NULL;

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
    if (emm_ctx) {
      if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
        emm_security_context = &emm_ctx->_security;
      }
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMAS-SAP - UE MM Context NULL for ue id " MME_UE_S1AP_ID_FMT "\n",
        ue_id);
  }

  /*
   * Decode the received message
   */
  decoder_rc = nas_message_decode(
      msg->data, &nas_msg, len, emm_security_context, decode_status);

  if (decoder_rc < 0) {
    OAILOG_ERROR(
        LOG_NAS_EMM,
        "EMMAS-SAP - Failed to decode NAS message for ue_id = "
        "(" MME_UE_S1AP_ID_FMT
        ")\n"
        "(err=%d)\n",
        ue_id, decoder_rc);
    *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
  }

  /*
   * Process NAS message
   */
  EMM_msg* emm_msg = &nas_msg.plain.emm;

  switch (emm_msg->header.message_type) {
    case EMM_STATUS:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }
      rc = emm_recv_status(
          ue_id, &emm_msg->emm_status, emm_cause, decode_status);
      break;

    case ATTACH_REQUEST:
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__1);  // Integrity checking of NAS signalling messages in
                            // the MME
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__2);  // Integrity checking of NAS signalling messages in
                            // the MME
      rc = emm_recv_attach_request(
          ue_id, originating_tai, originating_ecgi, &emm_msg->attach_request,
          false, false, emm_cause, decode_status);
      // Free the ESM container
      bdestroy(emm_msg->attach_request.esmmessagecontainer);
      break;

    case IDENTITY_RESPONSE:
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__1);  // Integrity checking of NAS signalling messages in
                            // the MME
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__2);  // Integrity checking of NAS signalling messages in
                            // the MME
      rc = emm_recv_identity_response(
          ue_id, &emm_msg->identity_response, emm_cause, decode_status);
      break;

    case AUTHENTICATION_RESPONSE:
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__1);  // Integrity checking of NAS signalling messages in
                            // the MME
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__2);  // Integrity checking of NAS signalling messages in
                            // the MME
      rc = emm_recv_authentication_response(
          ue_id, &emm_msg->authentication_response, emm_cause, decode_status);
      break;

    case AUTHENTICATION_FAILURE:
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__1);  // Integrity checking of NAS signalling messages in
                            // the MME
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__2);  // Integrity checking of NAS signalling messages in
                            // the MME
      rc = emm_recv_authentication_failure(
          ue_id, &emm_msg->authentication_failure, emm_cause, decode_status);
      break;

    case SECURITY_MODE_COMPLETE:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }

      rc = emm_recv_security_mode_complete(
          ue_id, &emm_msg->security_mode_complete, emm_cause, decode_status);
      break;

    case SECURITY_MODE_REJECT:
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__1);  // Integrity checking of NAS signalling messages in
                            // the MME
      REQUIREMENT_3GPP_24_301(
          R10_4_4_4_3__2);  // Integrity checking of NAS signalling messages in
                            // the MME
      rc = emm_recv_security_mode_reject(
          ue_id, &emm_msg->security_mode_reject, emm_cause, decode_status);
      break;

    case ATTACH_COMPLETE:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }

      rc = emm_recv_attach_complete(
          ue_id, &emm_msg->attach_complete, emm_cause, decode_status);
      bdestroy((bstring)(emm_msg->attach_complete.esmmessagecontainer));
      break;

    case TRACKING_AREA_UPDATE_COMPLETE:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }
      rc =
          emm_recv_tau_complete(ue_id, &emm_msg->tracking_area_update_complete);
      /* send the SGSAP TMSI Reallocation complete message towards SGS.
       * if csfb newTmsiAllocated flag is true
       * After sending set it to false
       */
      if (emm_ctx->csfbparams.newTmsiAllocated) {
        char imsi_str[IMSI_BCD_DIGITS_MAX + 1];
        IMSI_TO_STRING(&(emm_ctx->_imsi), imsi_str, IMSI_BCD_DIGITS_MAX + 1);
        mme_app_itti_sgsap_tmsi_reallocation_comp(imsi_str, strlen(imsi_str));
        emm_ctx->csfbparams.newTmsiAllocated = false;
        /* update the neaf flag to false after sending the Tmsi Reallocation
         * Complete message to SGS
         */
        mme_ue_context_update_ue_sgs_neaf(ue_id, false);
      }
      break;
    case GUTI_REALLOCATION_COMPLETE:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }
      break;
    case UPLINK_NAS_TRANSPORT:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }
      rc = emm_recv_uplink_nas_transport(
          ue_id, &emm_msg->uplink_nas_transport, emm_cause, decode_status);
      break;

    case DETACH_REQUEST:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__2);
      if ((1 == decode_status->security_context_available) &&
          (0 < emm_security_context->activated) &&
          ((0 == decode_status->integrity_protected_message) ||
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }

      rc = emm_recv_detach_request(
          ue_id, &emm_msg->detach_request, false, emm_cause, decode_status);
      break;
    case DETACH_ACCEPT:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__2);
      if ((1 == decode_status->security_context_available) &&
          (0 < emm_security_context->activated) &&
          ((0 == decode_status->integrity_protected_message) ||
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
      }

      rc = emm_recv_detach_accept(ue_id, emm_cause);
      break;

    case EXTENDED_SERVICE_REQUEST:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__2);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
        // Send Service Reject with cause "UE identity cannot be derived by the
        // network"
        rc = emm_proc_service_reject(
            ue_id, EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW);
        increment_counter(
            "extended_service_request", 1, 2, "result", "failure", "cause",
            "ue_context_not_available");
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
      }
      /* Process Extended-Service request */
      rc = emm_recv_ext_service_request(
          ue_id, &emm_msg->extended_service_request, emm_cause, decode_status);

      break;

    case TRACKING_AREA_UPDATE_REQUEST:
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__2);
      if ((0 == decode_status->security_context_available) ||
          (0 == decode_status->integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status->security_context_available) &&
           (0 == decode_status->mac_matched))) {
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
        rc         = emm_proc_tracking_area_update_reject(
            ue_id, EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW);
        increment_counter(
            "tracking_area_update", 1, 2, "result", "failure", "cause",
            "ue_context_not_available");
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
      }
      /* Process Extended-Service request */
      rc = emm_recv_tracking_area_update_request(
          ue_id, &emm_msg->tracking_area_update_request, false,
          originating_tai->tac, emm_cause, decode_status);

      break;

    default:
      OAILOG_WARNING(
          LOG_NAS_EMM,
          "EMMAS-SAP - EMM message 0x%x is not valid for ue "
          "id " MME_UE_S1AP_ID_FMT "\n",
          emm_msg->header.message_type, ue_id);
      *emm_cause = EMM_CAUSE_MESSAGE_TYPE_NOT_COMPATIBLE;
      break;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_data_ind()                                        **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP data transfer indication          **
 **      primitive                                                 **
 **                                                                        **
 ** EMMAS-SAP - AS->EMM: DATA_IND - Data transfer procedure                **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_cause: EMM cause code                             **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_data_ind(emm_as_data_t* msg, int* emm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = RETURNerror;

  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Received AS data transfer indication "
      "(ue_id=" MME_UE_S1AP_ID_FMT ", delivered=%s, length=%d)\n",
      msg->ue_id, (msg->delivered) ? "true" : "false", blength(msg->nas_msg));

  if (EMM_AS_DATA_DELIVERED_TRUE == msg->delivered) {
    if (blength(msg->nas_msg) > 0) {
      /*
       * Process the received NAS message
       */
      bstring plain_msg = bstrcpy(msg->nas_msg);

      if (plain_msg) {
        nas_message_security_header_t header = {0};
        emm_security_context_t* security =
            NULL; /* Current EPS NAS security context     */
        nas_message_decode_status_t decode_status = {0};

        /*
         * Decrypt the received security protected message
         */
        ue_mm_context_t* ue_mm_context =
            mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

        emm_context_t* emm_ctx = NULL;

        if (ue_mm_context) {
          emm_ctx = &ue_mm_context->emm_context;
          if (emm_ctx) {
            if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
              security = &emm_ctx->_security;
            }
          }
        } else {
          OAILOG_WARNING(
              LOG_NAS_EMM,
              "EMMAS-SAP - UE MM Context is NULL for ue id " MME_UE_S1AP_ID_FMT
              "\n",
              msg->ue_id);
        }

        int bytes = nas_message_decrypt(
            msg->nas_msg->data, plain_msg->data, &header, blength(msg->nas_msg),
            security, &decode_status);

        plain_msg->slen = bytes;

        if ((bytes < 0) &&
            (bytes !=
             TLV_MAC_MISMATCH)) {  // not in spec, (case identity response for
                                   // attach with unknown GUTI)
          /*
           * Failed to decrypt the message
           */
          *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, bytes);
        } else if (
            header.protocol_discriminator == EPS_MOBILITY_MANAGEMENT_MESSAGE) {
          /*
           * Process EMM data
           */
          tai_t originating_tai = {0};  // originating TAI
          memcpy(&originating_tai, msg->tai, sizeof(originating_tai));

          rc = emm_as_recv(
              msg->ue_id, &originating_tai, &msg->ecgi, plain_msg, bytes,
              emm_cause, &decode_status);
        } else if (
            header.protocol_discriminator == EPS_SESSION_MANAGEMENT_MESSAGE) {
          /*
           * Foward ESM data to EPS session management
           */
          // shrink plain_msg
          btrunc(plain_msg, bytes);
          rc = lowerlayer_data_ind(msg->ue_id, plain_msg);
        }

        bdestroy_wrapper(&plain_msg);
      }
    } else {
      /*
       * Process successfull lower layer transfer indication
       */
      rc = lowerlayer_success(msg->ue_id, &msg->nas_msg);
    }
  } else if (EMM_AS_DATA_DELIVERED_LOWER_LAYER_FAILURE == msg->delivered) {
    /*
     * Process lower layer transmission failure of NAS message
     */
    rc = lowerlayer_failure(msg->ue_id, &msg->nas_msg);
  } else {
    rc = lowerlayer_non_delivery_indication(msg->ue_id, &msg->nas_msg);
  }
  bdestroy_wrapper(&msg->nas_msg);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_establish_req()                                   **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP connection establish request      **
 **      primitive                                                 **
 **                                                                        **
 ** EMMAS-SAP - AS->EMM: ESTABLISH_REQ - NAS signalling connection         **
 **     The AS notifies the NAS that establishment of the signal-  **
 **     ling connection has been requested to tranfer initial NAS  **
 **     message from the UE.                                       **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     emm_cause: EMM cause code                             **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_establish_req(emm_as_establish_t* msg, int* emm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  struct emm_context_s* emm_ctx                = NULL;
  emm_security_context_t* emm_security_context = NULL;
  nas_message_decode_status_t decode_status    = {0};
  int decoder_rc                               = 0;
  int rc                                       = RETURNerror;
  tai_t originating_tai                        = {0};

  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - Received AS connection establish request\n");
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  ue_mm_context_t* ue_mm_context =
      mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

  if (ue_mm_context) {
    emm_ctx = &ue_mm_context->emm_context;
    if (emm_ctx) {
      OAILOG_DEBUG(LOG_NAS_EMM, "EMMAS-SAP - got context %p\n", emm_ctx);
      if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
        emm_security_context = &emm_ctx->_security;
      }
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMAS-SAP - ue mm context null for ue id " MME_UE_S1AP_ID_FMT "\n",
        msg->ue_id);
  }

  /*
   * Decode initial NAS message
   */
  OAILOG_DEBUG(
      LOG_NAS_EMM,
      "EMMAS-SAP - Decoding Initial NAS message for ue_id = (%u)\n",
      msg->ue_id);
  decoder_rc = nas_message_decode(
      msg->nas_msg->data, &nas_msg, blength(msg->nas_msg), emm_security_context,
      &decode_status);
  bdestroy_wrapper(&msg->nas_msg);

  // TODO conditional IE error
  if (decoder_rc < 0) {
    if (decoder_rc < TLV_FATAL_ERROR) {
      *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoder_rc);
    } else if (decoder_rc == TLV_MANDATORY_FIELD_NOT_PRESENT) {
      *emm_cause = EMM_CAUSE_INVALID_MANDATORY_INFO;
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_b__1);
    } else if (decoder_rc == TLV_UNEXPECTED_IEI) {
      *emm_cause = EMM_CAUSE_IE_NOT_IMPLEMENTED;
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_b__2);
    } else {
      *emm_cause = EMM_CAUSE_PROTOCOL_ERROR;
      REQUIREMENT_3GPP_24_301(R10_5_5_1_2_7_b__4);
    }
  }

  /*
   * Process initial NAS message
   */
  EMM_msg* emm_msg = &nas_msg.plain.emm;

  switch (emm_msg->header.message_type) {
    case ATTACH_REQUEST:
      memcpy(&originating_tai, msg->tai, sizeof(originating_tai));
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - Message Type = ATTACH_REQUEST(0x%x) for (ue_id "
          "= " MME_UE_S1AP_ID_FMT ")\n",
          emm_msg->header.message_type, msg->ue_id);
      rc = emm_recv_attach_request(
          msg->ue_id, &originating_tai, &msg->ecgi, &emm_msg->attach_request,
          msg->is_initial, msg->is_mm_ctx_new, emm_cause, &decode_status);
      break;

    case DETACH_REQUEST:
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - Message Type = DETACH_REQUEST(0x%x) for (ue_id "
          "= " MME_UE_S1AP_ID_FMT ")\n",
          emm_msg->header.message_type, msg->ue_id);
      if (emm_ctx == NULL) {
        /*
         * This means UE context is not present and this UE is not known in the
         * EPC. If message is Detach Request. Ignore the message.
         */
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;

        OAILOG_WARNING(
            LOG_NAS_EMM,
            "EMM-PROC- Received Detach Request from unknown "
            "UE.ue_id=" MME_UE_S1AP_ID_FMT " \n",
            msg->ue_id);
        // Clean up S1AP and MME UE Context
        mme_app_handle_detach_req(ue_mm_context->mme_ue_s1ap_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      }

      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__1);
      REQUIREMENT_3GPP_24_301(R10_4_4_4_3__2);
      if ((1 == decode_status.security_context_available) &&
          (0 < emm_security_context->activated) &&
          ((0 == decode_status.integrity_protected_message) ||
           (0 == decode_status.mac_matched))) {
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
        // Delete EMM,ESM conext, MMEAPP UE context and S1AP context
        nas_proc_implicit_detach_ue_ind(ue_mm_context->mme_ue_s1ap_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      }
      // Process Detach Request
      rc = emm_recv_detach_request(
          msg->ue_id, &emm_msg->detach_request, msg->is_initial, emm_cause,
          &decode_status);
      break;

    case TRACKING_AREA_UPDATE_REQUEST:
      increment_counter("tracking_area_update", 1, NO_LABELS);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - Message Type = TRACKING_AREA_UPDATE_REQUEST(0x%x)"
          "for (ue_id = " MME_UE_S1AP_ID_FMT ")\n",
          emm_msg->header.message_type, msg->ue_id);
      // Check for emm_ctx and integrity verification
      if ((emm_ctx == NULL) ||
          ((0 == decode_status.security_context_available) ||
           (0 == decode_status.integrity_protected_message) ||
           ((1 == decode_status.security_context_available) &&
            (0 == decode_status.mac_matched)))) {
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
        // Send Reject with cause "UE identity cannot be derived by the network"
        // to trigger fresh attach
        rc = emm_proc_tracking_area_update_reject(
            msg->ue_id, EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW);
        increment_counter(
            "tracking_area_update", 1, 2, "result", "failure", "cause",
            "ue_context_not_available");
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
      }

      // Process periodic TAU
      rc = emm_recv_tracking_area_update_request(
          msg->ue_id, &emm_msg->tracking_area_update_request, msg->is_initial,
          msg->tai->tac, emm_cause, &decode_status);
      break;

    case SERVICE_REQUEST:
      // Requirement MME24.301R10_4.4.4.3_1
      increment_counter("service_request", 1, NO_LABELS);
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - Message Type = SERVICE_REQUEST(0x%x) for (ue_id "
          "= " MME_UE_S1AP_ID_FMT ")\n",
          emm_msg->header.message_type, msg->ue_id);
      if ((emm_ctx == NULL) ||
          ((0 == decode_status.security_context_available) ||
           (0 == decode_status.integrity_protected_message) ||
           // Requirement MME24.301R10_4.4.4.3_2
           ((1 == decode_status.security_context_available) &&
            (0 == decode_status.mac_matched)))) {
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
        // Send Service Reject with cause "UE identity cannot be derived by the
        // network" to trigger fresh attach
        rc = emm_proc_service_reject(
            msg->ue_id, EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW);
        increment_counter(
            "service_request", 1, 2, "result", "failure", "cause",
            "ue_context_not_available");
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
      }
      // Process Service request
      rc = emm_recv_service_request(
          msg->ue_id, &emm_msg->service_request, msg->is_initial, emm_cause,
          &decode_status);
      break;

    case EXTENDED_SERVICE_REQUEST:
      /* Requirement MME24.301R10_4.4.4.3_1 */
      OAILOG_INFO(
          LOG_NAS_EMM,
          "EMMAS-SAP - Message Type = EXTENDED_SERVICE_REQUEST(0x%x) for "
          "(ue_id = " MME_UE_S1AP_ID_FMT ")\n",
          emm_msg->header.message_type, msg->ue_id);
      if ((0 == decode_status.security_context_available) ||
          (0 == decode_status.integrity_protected_message) ||
          // Requirement MME24.301R10_4.4.4.3_2
          ((1 == decode_status.security_context_available) &&
           (0 == decode_status.mac_matched))) {
        *emm_cause = EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW;
        // Send Service Reject with cause "UE identity cannot be derived by the
        // network" to trigger fresh attach
        rc = emm_proc_service_reject(
            msg->ue_id, EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW);
        increment_counter(
            "extended_service_request", 1, 2, "result", "failure", "cause",
            "ue_context_not_available");
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
      }
      /* Process Extended-Service request */
      rc = emm_recv_initial_ext_service_request(
          msg->ue_id, &emm_msg->extended_service_request, emm_cause,
          &decode_status);
      break;

    default:
      OAILOG_WARNING(
          LOG_NAS_EMM,
          "EMMAS-SAP - Initial NAS message 0x%x is "
          "not valid (ue_id = " MME_UE_S1AP_ID_FMT ")\n",
          emm_msg->header.message_type, msg->ue_id);
      *emm_cause = EMM_CAUSE_MESSAGE_TYPE_NOT_COMPATIBLE;
      break;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

//------------------------------------------------------------------------------
static int emm_as_release_ind(
    const emm_as_release_t* const release, int* emm_cause) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int rc = lowerlayer_release(release->ue_id, release->cause);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/*
   --------------------------------------------------------------------------
   Functions executed to send data to the network when requested within EMM
   procedure processing
   --------------------------------------------------------------------------
*/

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_set_header()                                      **
 **                                                                        **
 ** Description: Setup the security header of the given NAS message        **
 **                                                                        **
 ** Inputs:  security:  The NAS security data to use               **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     msg:       The NAS message                            **
 **      Return:    Pointer to the plain NAS message to be se- **
 **             curity protected if setting of the securi- **
 **             ty header succeed;                         **
 **             NULL pointer otherwise                     **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static EMM_msg* emm_as_set_header(
    nas_message_t* msg, const emm_as_security_data_t* security) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  msg->header.protocol_discriminator = EPS_MOBILITY_MANAGEMENT_MESSAGE;

  if (security && (security->ksi != KSI_NO_KEY_AVAILABLE)) {
    /*
     * A valid EPS security context exists
     */
    if (security->is_new) {
      /*
       * New EPS security context is taken into use
       */
      if (security->is_knas_int_present) {
        if (security->is_knas_enc_present) {
          /*
           * NAS integrity and cyphering keys are available
           */
          msg->header.security_header_type =
              SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED_NEW;
        } else {
          /*
           * NAS integrity key only is available
           */
          msg->header.security_header_type =
              SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_NEW;
        }

        OAILOG_FUNC_RETURN(LOG_NAS_EMM, &msg->security_protected.plain.emm);
      }
    } else if (security->is_knas_int_present) {
      if (security->is_knas_enc_present) {
        /*
         * NAS integrity and cyphering keys are available
         */
        msg->header.security_header_type =
            SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
      } else {
        /*
         * NAS integrity key only is available
         */
        msg->header.security_header_type =
            SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED;
      }

      OAILOG_FUNC_RETURN(LOG_NAS_EMM, &msg->security_protected.plain.emm);
    } else {
      /*
       * No valid EPS security context exists
       */
      msg->header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, &msg->plain.emm);
    }
  } else {
    /*
     * No valid EPS security context exists
     */
    msg->header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, &msg->plain.emm);
  }

  /*
   * A valid EPS security context exists but NAS integrity key
   * * * * is not available
   */
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, NULL);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_encode()                                          **
 **                                                                        **
 ** Description: Encodes NAS message into NAS information container        **
 **                                                                        **
 ** Inputs:  msg:       The NAS message to encode                  **
 **      length:    The maximum length of the NAS message      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     info:      The NAS information container              **
 **      msg:       The NAS message to encode                  **
 **      Return:    The number of bytes successfully encoded   **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_encode(
    bstring* info, nas_message_t* msg, size_t length,
    emm_security_context_t* emm_security_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int bytes = 0;

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  EMM_GET_BYTE_ALIGNED_LENGTH(length);
  if (msg->header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    emm_msg_header_t* header = &msg->security_protected.plain.emm.header;

    /*
     * Expand size of protected NAS message
     */
    length += NAS_MESSAGE_SECURITY_HEADER_SIZE;
    /*
     * Set header of plain NAS message
     */
    header->protocol_discriminator = EPS_MOBILITY_MANAGEMENT_MESSAGE;
    header->security_header_type   = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  }

  /*
   * Allocate memory to the NAS information container
   */
  *info = bfromcstralloc(length, "\0");

  if (*info) {
    /*
     * Encode the NAS message
     */
    bytes =
        nas_message_encode((*info)->data, msg, length, emm_security_context);

    if (bytes > 0) {
      (*info)->slen = bytes;
    } else {
      bdestroy_wrapper(info);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_encrypt()                                         **
 **                                                                        **
 ** Description: Encryts NAS message into NAS information container        **
 **                                                                        **
 ** Inputs:  header:    The Security header in used                **
 **      msg:       The NAS message to encrypt                 **
 **      length:    The maximum length of the NAS message      **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     info:      The NAS information container              **
 **      Return:    The number of bytes successfully encrypted **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_encrypt(
    bstring* info, const nas_message_security_header_t* header,
    const unsigned char* msg, size_t length,
    emm_security_context_t* emm_security_context) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int bytes = 0;

  if (header->security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    /*
     * Expand size of protected NAS message
     */
    length += NAS_MESSAGE_SECURITY_HEADER_SIZE;
  }

  /*
   * Allocate memory to the NAS information container
   */
  *info = bfromcstralloc(length, "\0");

  if (*info) {
    /*
     * Encrypt the NAS information message
     */
    bytes = nas_message_encrypt(
        msg, (*info)->data, header, length, emm_security_context);

    if (bytes > 0) {
      (*info)->slen = bytes;
    } else {
      bdestroy_wrapper(info);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, bytes);
}

/****************************************************************************
 **                                                                        **
 ** Name:    emm_as_send_a()                                            **
 **                                                                        **
 ** Description: Builds NAS message according to the given EMMAS Service   **
 **      Access Point primitive and sends it to the Access Stratum **
 **      sublayer                                                  **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to be sent         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     None                                                      **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_send_a(const emm_as_t* msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  as_message_t as_msg = {0};

  switch (msg->primitive) {
    case _EMMAS_DATA_REQ:
      as_msg.msg_id =
          emm_as_data_req(&msg->u.data, &as_msg.msg.dl_info_transfer_req);
      break;

    case _EMMAS_ERAB_SETUP_REQ:
      as_msg.msg_id = emm_as_erab_setup_req(
          &msg->u.activate_bearer_context_req,
          &as_msg.msg.activate_bearer_context_req);
      break;

    case _EMMAS_ERAB_REL_CMD:
      as_msg.msg_id = emm_as_erab_rel_cmd(
          &msg->u.deactivate_bearer_context_req,
          &as_msg.msg.deactivate_bearer_context_req);
      break;

    case _EMMAS_STATUS_IND:
      as_msg.msg_id =
          emm_as_status_ind(&msg->u.status, &as_msg.msg.dl_info_transfer_req);
      break;

    case _EMMAS_RELEASE_REQ:
      as_msg.msg_id =
          emm_as_release_req(&msg->u.release, &as_msg.msg.nas_release_req);
      break;

    case _EMMAS_SECURITY_REQ:
      as_msg.msg_id = emm_as_security_req(
          &msg->u.security, &as_msg.msg.dl_info_transfer_req);
      break;

    case _EMMAS_SECURITY_REJ:
      as_msg.msg_id = emm_as_security_rej(
          &msg->u.security, &as_msg.msg.dl_info_transfer_req);
      break;

    case _EMMAS_ESTABLISH_CNF:
      as_msg.msg_id = emm_as_establish_cnf(
          &msg->u.establish, &as_msg.msg.nas_establish_rsp);
      break;

    case _EMMAS_ESTABLISH_REJ:
      as_msg.msg_id = emm_as_establish_rej(
          &msg->u.establish, &as_msg.msg.nas_establish_rsp);
      break;

    default:
      as_msg.msg_id = 0;
      break;
  }

  /*
   * Send the message to the Access Stratum or S1AP in case of MME
   */
  if (as_msg.msg_id > 0) {
    OAILOG_DEBUG(
        LOG_NAS_EMM,
        "EMMAS-SAP - "
        "Sending msg with id 0x%x, primitive %s (%d) to S1AP layer for "
        "transmission\n",
        as_msg.msg_id, emm_as_primitive_str[msg->primitive - _EMMAS_START - 1],
        msg->primitive);

    switch (as_msg.msg_id) {
      case AS_DL_INFO_TRANSFER_REQ: {
        mme_app_handle_nas_dl_req(
            as_msg.msg.dl_info_transfer_req.ue_id,
            as_msg.msg.dl_info_transfer_req.nas_msg,
            as_msg.msg.dl_info_transfer_req.err_code);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      } break;

      case AS_ACTIVATE_BEARER_CONTEXT_REQ: {
        mme_app_handle_erab_setup_req(
            as_msg.msg.activate_bearer_context_req.ue_id,
            as_msg.msg.activate_bearer_context_req.ebi,
            as_msg.msg.activate_bearer_context_req.mbr_dl,
            as_msg.msg.activate_bearer_context_req.mbr_ul,
            as_msg.msg.activate_bearer_context_req.gbr_dl,
            as_msg.msg.activate_bearer_context_req.gbr_ul,
            as_msg.msg.activate_bearer_context_req.nas_msg);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      } break;

      case AS_DEACTIVATE_BEARER_CONTEXT_REQ: {
        mme_app_handle_erab_rel_cmd(
            as_msg.msg.deactivate_bearer_context_req.ue_id,
            as_msg.msg.deactivate_bearer_context_req.ebi,
            as_msg.msg.deactivate_bearer_context_req.nas_msg);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
      } break;

      case AS_NAS_ESTABLISH_RSP:
      case AS_NAS_ESTABLISH_CNF: {
        if (as_msg.msg.nas_establish_rsp.err_code != AS_SUCCESS) {
          // This flow is to release the UE context after sending the NAS
          // message.
          mme_app_handle_nas_dl_req(
              as_msg.msg.nas_establish_rsp.ue_id,
              as_msg.msg.nas_establish_rsp.nas_msg,
              as_msg.msg.nas_establish_rsp.err_code);
          as_msg.msg.nas_establish_rsp.nas_msg = NULL;
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        } else {
          OAILOG_DEBUG(
              LOG_NAS_EMM,
              "EMMAS-SAP - Sending establish_cnf to MME-APP module for UE "
              "ID: " MME_UE_S1AP_ID_FMT
              " selected eea "
              "0x%04X selected eia 0x%04X\n",
              as_msg.msg.nas_establish_rsp.ue_id,
              as_msg.msg.nas_establish_rsp.selected_encryption_algorithm,
              as_msg.msg.nas_establish_rsp.selected_integrity_algorithm);
          /*
           * Handle success case
           */
          mme_app_handle_conn_est_cnf(&as_msg.msg.nas_establish_rsp);
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        }
      } break;

      case AS_NAS_RELEASE_REQ:
        mme_app_handle_detach_req(as_msg.msg.nas_release_req.ue_id);
        OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNok);
        break;

      default:
        break;
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_data_req()                                        **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP data transfer request             **
 **      primitive                                                 **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: DATA_REQ - Data transfer procedure                **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_data_req(
    const emm_as_data_t* msg, dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size       = 0;
  int is_encoded = false;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS data transfer request\n");
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  /*
   * Setup the AS message
   */
  if (msg->guti) {
    as_msg->s_tmsi.mme_code = msg->guti->gummei.mme_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  /*
   * Setup the NAS security header
   */
  EMM_msg* emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS information message
   */
  if (emm_msg) switch (msg->nas_info) {
      case EMM_AS_NAS_DATA_ATTACH_ACCEPT:
        size = emm_send_attach_accept_dl_nas(msg, &emm_msg->attach_accept);
        break;

      case EMM_AS_NAS_DATA_DETACH_ACCEPT:
        size = emm_send_detach_accept(msg, &emm_msg->detach_accept);
        break;

      case EMM_AS_NAS_DATA_TAU:
        size = emm_send_tracking_area_update_accept_dl_nas(
            msg, &emm_msg->tracking_area_update_accept);
        break;

      case EMM_AS_NAS_EMM_INFORMATION:
        size = emm_send_emm_information(msg, &emm_msg->emm_information);
        break;

      case EMM_AS_NAS_DATA_DETACH_REQ:
        size = emm_send_nw_detach_request(msg, &emm_msg->nw_detach_request);
        break;

      case EMM_AS_NAS_DL_NAS_TRANSPORT:
        size = emm_send_dl_nas_transport(msg, &emm_msg->downlink_nas_transport);
        break;
      case EMM_AS_NAS_DATA_CS_SERVICE_NOTIFICATION:
        size = emm_send_cs_service_notification(
            msg, &emm_msg->cs_service_notification);
        break;

      case EMM_AS_NAS_DATA_INFO_SR:
        size = emm_send_service_reject(
            *(msg->emm_cause), &emm_msg->service_reject);
        break;

      default:
        /*
         * Send other NAS messages as already encoded ESM messages
         */
        size       = msg->nas_msg->slen;
        is_encoded = true;
        break;
    }

  if (size > 0) {
    int bytes                                    = 0;
    emm_security_context_t* emm_security_context = NULL;
    struct emm_context_s* emm_ctx                = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
      if (emm_ctx) {
        if (emm_msg->nw_detach_request.nw_detachtype ==
            NW_DETACH_TYPE_IMSI_DETACH) {
          emm_ctx->is_imsi_only_detach = true;
        }
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context = &emm_ctx->_security;
        }
      }
    }

    if (emm_security_context) {
      nas_msg.header.sequence_number = emm_security_context->dl_count.seq_num;
      OAILOG_DEBUG(
          LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
          nas_msg.header.sequence_number);
    } else {
      OAILOG_ERROR(
          LOG_NAS_EMM,
          "Security context is NULL for UE -> " MME_UE_S1AP_ID_FMT "\n",
          msg->ue_id);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, RETURNerror);
    }

    if (!is_encoded) {
      /*
       * Encode the NAS information message
       */
      bytes =
          emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);
    } else {
      /*
       * Encrypt the NAS information message
       */
      bytes = emm_as_encrypt(
          &as_msg->nas_msg, &nas_msg.header, msg->nas_msg->data, size,
          emm_security_context);
    }

    // Free any allocated data
    switch (msg->nas_info) {
      // emm_information message and downlink_nas_transtport is the only message
      // that has allocated data
      case EMM_AS_NAS_EMM_INFORMATION:
        emm_free_send_emm_information(&emm_msg->emm_information);
        break;
      case EMM_AS_NAS_DL_NAS_TRANSPORT:
        emm_free_send_dl_nas_transport(&emm_msg->downlink_nas_transport);
        break;
      case EMM_AS_NAS_DATA_CS_SERVICE_NOTIFICATION:
        emm_free_send_cs_service_notification(
            &emm_msg->cs_service_notification);
        break;
      case EMM_AS_NAS_DATA_ATTACH_ACCEPT:
        bdestroy_wrapper(&(emm_msg->attach_accept.esmmessagecontainer));
        break;
    }

    if (bytes > 0) {
      /*If new TMSI is allocated by MSC. We need to wait for TAU Complete
       * message. After receiving TAU Complete, we should trigger UE context
       * release
       */
      if ((msg->nas_info == EMM_AS_NAS_DATA_TAU) &&
          !(emm_ctx->csfbparams.newTmsiAllocated) &&
          !(emm_ctx->csfbparams.tau_active_flag)) {
        as_msg->err_code = AS_TERMINATED_NAS;
      } else {
        as_msg->err_code = AS_SUCCESS;
      }
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_DL_INFO_TRANSFER_REQ);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_status_ind()                                      **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP status indication primitive       **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: STATUS_IND - EMM status report procedure          **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_status_ind(
    const emm_as_status_t* msg, dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = 0;

  OAILOG_INFO(
      LOG_NAS_EMM, "EMMAS-SAP - Send AS status indication (cause=%d)\n",
      msg->emm_cause);
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  /*
   * Setup the AS message
   */
  if (msg->guti) {
    as_msg->s_tmsi.mme_code = msg->guti->gummei.mme_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  /*
   * Setup the NAS security header
   */
  EMM_msg* emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS information message
   */
  if (emm_msg) {
    size = emm_send_status(msg, &emm_msg->emm_status);
  }

  if (size > 0) {
    emm_security_context_t* emm_security_context = NULL;
    struct emm_context_s* emm_ctx                = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
      if (emm_ctx) {
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context = &emm_ctx->_security;
        }
      }
    }

    if (emm_security_context) {
      nas_msg.header.sequence_number = emm_security_context->dl_count.seq_num;
      OAILOG_DEBUG(
          LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
          nas_msg.header.sequence_number);
    }

    /*
     * Encode the NAS information message
     */
    int bytes =
        emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);

    if (bytes > 0) {
      as_msg->err_code = AS_SUCCESS;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_DL_INFO_TRANSFER_REQ);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_release_req()                                     **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP connection release request        **
 **      primitive                                                 **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: RELEASE_REQ - NAS signalling release procedure    **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_release_req(
    const emm_as_release_t* msg, nas_release_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS release request\n");

  /*
   * Setup the AS message
   */
  if (msg->guti) {
    as_msg->s_tmsi.mme_code = msg->guti->gummei.mme_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  if (msg->cause == EMM_AS_CAUSE_AUTHENTICATION) {
    as_msg->cause = AS_AUTHENTICATION_FAILURE;
  } else if (msg->cause == EMM_AS_CAUSE_DETACH) {
    as_msg->cause = AS_DETACH;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_NAS_RELEASE_REQ);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_security_req()                                    **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP security request primitive        **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: SECURITY_REQ - Security mode control procedure    **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_security_req(
    const emm_as_security_t* msg, dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = 0;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS security request\n");
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  /*
   * Setup the AS message
   */
  if (msg->guti) {
    as_msg->s_tmsi.mme_code = msg->guti->gummei.mme_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }
  /*
   * Setup the NAS security header
   */
  EMM_msg* emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS security message
   */
  if (emm_msg) switch (msg->msg_type) {
      case EMM_AS_MSG_TYPE_IDENT:
        size = emm_send_identity_request(msg, &emm_msg->identity_request);
        break;

      case EMM_AS_MSG_TYPE_AUTH:
        size = emm_send_authentication_request(
            msg, &emm_msg->authentication_request);
        break;

      case EMM_AS_MSG_TYPE_SMC:
        size = emm_send_security_mode_command(
            msg, &emm_msg->security_mode_command);
        break;

      default:
        OAILOG_WARNING(
            LOG_NAS_EMM,
            "EMMAS-SAP - Type of NAS security "
            "message 0x%.2x is not valid\n",
            msg->msg_type);
    }

  if (size > 0) {
    struct emm_context_s* emm_ctx                = NULL;
    emm_security_context_t* emm_security_context = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;

      if (emm_ctx) {
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context           = &emm_ctx->_security;
          nas_msg.header.sequence_number = emm_ctx->_security.dl_count.seq_num;
          OAILOG_DEBUG(
              LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
              nas_msg.header.sequence_number);
        }
      }
    } else {
      OAILOG_WARNING(
          LOG_NAS_EMM,
          "UE MM context NULL for ue_id = (" MME_UE_S1AP_ID_FMT ")\n",
          msg->ue_id);
    }

    /*
     * Encode the NAS security message
     */
    int bytes =
        emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);
    // Free any allocated data
    switch (msg->msg_type) {
      // authentication_request is the only message with allocated mem
      case EMM_AS_MSG_TYPE_AUTH:
        emm_free_send_authentication_request(&emm_msg->authentication_request);
        break;
    }

    if (bytes > 0) {
      as_msg->err_code = AS_SUCCESS;
      nas_emm_procedure_register_emm_message(
          msg->ue_id, msg->puid, as_msg->nas_msg);
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_DL_INFO_TRANSFER_REQ);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_security_rej()                                    **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP security reject primitive         **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: SECURITY_REJ - Security mode control procedure    **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_security_rej(
    const emm_as_security_t* msg, dl_info_transfer_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size = 0;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS security reject\n");
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  /*
   * Setup the AS message
   */
  if (msg->guti) {
    as_msg->s_tmsi.mme_code = msg->guti->gummei.mme_code;
    as_msg->s_tmsi.m_tmsi   = msg->guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  /*
   * Setup the NAS security header
   */
  EMM_msg* emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS security message
   */
  if (emm_msg) switch (msg->msg_type) {
      case EMM_AS_MSG_TYPE_AUTH:
        size = emm_send_authentication_reject(&emm_msg->authentication_reject);
        break;

      default:
        OAILOG_WARNING(
            LOG_NAS_EMM,
            "EMMAS-SAP - Type of NAS security "
            "message 0x%.2x is not valid\n",
            msg->msg_type);
    }

  if (size > 0) {
    struct emm_context_s* emm_ctx                = NULL;
    emm_security_context_t* emm_security_context = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
      if (emm_ctx) {
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context = &emm_ctx->_security;
          nas_msg.header.sequence_number =
              emm_security_context->dl_count.seq_num;
          OAILOG_DEBUG(
              LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
              nas_msg.header.sequence_number);
        } else {
          OAILOG_DEBUG(
              LOG_NAS_EMM,
              "No security context, not set nas_msg.header.sequence_number -> "
              "%u\n",
              nas_msg.header.sequence_number);
        }
      }
    }

    /*
     * Encode the NAS security message
     */
    int bytes =
        emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);

    if (bytes > 0) {
      /*
       *  Indicate to lower layer that procedure needs to be terminated after
       * sending DL NAS message.
       */
      as_msg->err_code = AS_TERMINATED_NAS;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_DL_INFO_TRANSFER_REQ);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}

//------------------------------------------------------------------------------
static int emm_as_erab_setup_req(
    const emm_as_activate_bearer_context_req_t* msg,
    activate_bearer_context_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size       = 0;
  int is_encoded = false;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS data transfer request\n");
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  /*
   * Setup the AS message
   */
  as_msg->ue_id  = msg->ue_id;
  as_msg->ebi    = msg->ebi;
  as_msg->gbr_dl = msg->gbr_dl;
  as_msg->gbr_ul = msg->gbr_ul;
  as_msg->mbr_dl = msg->mbr_dl;
  as_msg->mbr_ul = msg->mbr_ul;

  /*
   * Setup the NAS security header
   */
  EMM_msg* emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS information message
   */
  if (emm_msg) {
    size       = msg->nas_msg->slen;
    is_encoded = true;
  }

  if (size > 0) {
    int bytes                                    = 0;
    emm_security_context_t* emm_security_context = NULL;
    struct emm_context_s* emm_ctx                = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
      if (emm_ctx) {
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context = &emm_ctx->_security;
        }
      }
    }

    if (emm_security_context) {
      nas_msg.header.sequence_number = emm_security_context->dl_count.seq_num;
      OAILOG_DEBUG(
          LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
          nas_msg.header.sequence_number);
    }

    if (!is_encoded) {
      /*
       * Encode the NAS information message
       */
      bytes =
          emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);
    } else {
      /*
       * Encrypt the NAS information message
       */
      bytes = emm_as_encrypt(
          &as_msg->nas_msg, &nas_msg.header, msg->nas_msg->data, size,
          emm_security_context);
    }

    if (bytes > 0) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_ACTIVATE_BEARER_CONTEXT_REQ);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}

//------------------------------------------------------------------------------
static int emm_as_erab_rel_cmd(
    const emm_as_deactivate_bearer_context_req_t* msg,
    deactivate_bearer_context_req_t* as_msg) {
  OAILOG_FUNC_IN(LOG_NAS_EMM);
  int size       = 0;
  int is_encoded = false;

  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS data transfer request\n");
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  /*
   * Setup the AS message
   */
  as_msg->ue_id = msg->ue_id;
  as_msg->ebi   = msg->ebi;

  /*
   * Setup the NAS security header
   */
  EMM_msg* emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS information message
   */
  if (emm_msg) {
    size       = msg->nas_msg->slen;
    is_encoded = true;
  }

  if (size > 0) {
    int bytes                                    = 0;
    emm_security_context_t* emm_security_context = NULL;
    struct emm_context_s* emm_ctx                = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
      if (emm_ctx) {
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context = &emm_ctx->_security;
        }
      }
    }

    if (emm_security_context) {
      nas_msg.header.sequence_number = emm_security_context->dl_count.seq_num;
      OAILOG_DEBUG(
          LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
          nas_msg.header.sequence_number);
    }

    if (!is_encoded) {
      /*
       * Encode the NAS information message
       */
      bytes =
          emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);
    } else {
      /*
       * Encrypt the NAS information message
       */
      bytes = emm_as_encrypt(
          &as_msg->nas_msg, &nas_msg.header, msg->nas_msg->data, size,
          emm_security_context);
    }

    if (bytes > 0) {
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_DEACTIVATE_BEARER_CONTEXT_REQ);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_establish_cnf()                                   **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP connection establish confirm      **
 **      primitive                                                 **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: ESTABLISH_CNF - NAS signalling connection         **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_establish_cnf(
    const emm_as_establish_t* msg, nas_establish_rsp_t* as_msg) {
  EMM_msg* emm_msg = NULL;
  int size         = 0;
  int ret_val      = 0;

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_INFO(
      LOG_NAS_EMM,
      "EMMAS-SAP - Send AS connection establish confirmation for (ue_id "
      "= " MME_UE_S1AP_ID_FMT ")\n",
      msg->ue_id);
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};
  /*
   * Setup the AS message
   */
  as_msg->ue_id = msg->ue_id;

  if (msg->eps_id.guti == NULL) {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMAS-SAP - GUTI is NULL for ue id " MME_UE_S1AP_ID_FMT "\n",
        msg->ue_id);
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, ret_val);
  }

  as_msg->s_tmsi.mme_code = msg->eps_id.guti->gummei.mme_code;
  as_msg->s_tmsi.m_tmsi   = msg->eps_id.guti->m_tmsi;
  as_msg->nas_msg         = msg->nas_msg;
  as_msg->csfb_response   = msg->csfb_response;
  as_msg->presencemask    = msg->presencemask;
  as_msg->service_type    = msg->service_type;

  struct emm_context_s* emm_ctx                = NULL;
  emm_security_context_t* emm_security_context = NULL;
  emm_ctx = emm_context_get(&_emm_data, msg->ue_id);
  if (emm_ctx) {
    if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
      emm_security_context                  = &emm_ctx->_security;
      as_msg->selected_encryption_algorithm = (uint16_t) htons(
          0x10000 >> emm_security_context->selected_algorithms.encryption);
      as_msg->selected_integrity_algorithm = (uint16_t) htons(
          0x10000 >> emm_security_context->selected_algorithms.integrity);
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "Set nas_msg.selected_encryption_algorithm -> NBO: 0x%04X (%u)\n",
          as_msg->selected_encryption_algorithm,
          emm_security_context->selected_algorithms.encryption);
      OAILOG_DEBUG(
          LOG_NAS_EMM,
          "Set nas_msg.selected_integrity_algorithm -> NBO: 0x%04X (%u)\n",
          as_msg->selected_integrity_algorithm,
          emm_security_context->selected_algorithms.integrity);
      as_msg->nas_ul_count = 0x00000000 |
                             (emm_security_context->ul_count.overflow << 8) |
                             emm_security_context->ul_count
                                 .seq_num;  // This is sent to calculate KeNB
      OAILOG_DEBUG(
          LOG_NAS_EMM, "EMMAS-SAP - NAS UL COUNT %8x\n", as_msg->nas_ul_count);
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_EMM,
        "EMMAS-SAP - EMM Context is NULL...! for ue id " MME_UE_S1AP_ID_FMT
        "\n",
        msg->ue_id);
  }
  switch (msg->nas_info) {
    case EMM_AS_NAS_INFO_ATTACH:
      /*
       * Setup the NAS security header
       */
      emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);
      if (emm_msg) {
        OAILOG_TRACE(
            LOG_NAS_EMM, "EMMAS-SAP - emm_as_establish.nasMSG.length=%d\n",
            msg->nas_msg->slen);
        size = emm_send_attach_accept(msg, &emm_msg->attach_accept);
      }
      break;

    case EMM_AS_NAS_INFO_TAU:
      /*
       * Setup the NAS security header
       */
      emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);
      if (emm_msg) {
        size = emm_send_tracking_area_update_accept(
            msg, &emm_msg->tracking_area_update_accept);
      }
      break;
    case EMM_AS_NAS_INFO_NONE:  // Response to SR
      as_msg->err_code = AS_SUCCESS;
      ret_val          = AS_NAS_ESTABLISH_CNF;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, ret_val);
    default:
      OAILOG_WARNING(
          LOG_NAS_EMM,
          "EMMAS-SAP - Type of initial NAS "
          "message 0x%.2x is not valid\n",
          msg->nas_info);
      break;
  }

  if (size > 0) {
    nas_msg.header.sequence_number = emm_security_context->dl_count.seq_num;
    OAILOG_DEBUG(
        LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
        nas_msg.header.sequence_number);
  } else {
    OAILOG_FUNC_RETURN(LOG_NAS_EMM, ret_val);
  }
  /*
   * Encode the initial NAS information message
   */
  int bytes =
      emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);

  // Free any allocated data
  if (msg->nas_info == EMM_AS_NAS_INFO_ATTACH) {
    bdestroy_wrapper(&(emm_msg->attach_accept.esmmessagecontainer));
  }

  if (bytes > 0) {
    as_msg->err_code = AS_SUCCESS;
    ret_val          = AS_NAS_ESTABLISH_CNF;
  }
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, ret_val);
}

/****************************************************************************
 **                                                                        **
 ** Name:    _emm_as_establish_rej()                                   **
 **                                                                        **
 ** Description: Processes the EMMAS-SAP connection establish reject       **
 **      primitive                                                 **
 **                                                                        **
 ** EMMAS-SAP - EMM->AS: ESTABLISH_REJ - NAS signalling connection         **
 **                                                                        **
 ** Inputs:  msg:       The EMMAS-SAP primitive to process         **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     as_msg:    The message to send to the AS              **
 **      Return:    The identifier of the AS message           **
 **      Others:    None                                       **
 **                                                                        **
 ***************************************************************************/
static int emm_as_establish_rej(
    const emm_as_establish_t* msg, nas_establish_rsp_t* as_msg) {
  EMM_msg* emm_msg      = NULL;
  int size              = 0;
  nas_message_t nas_msg = {.security_protected.header           = {0},
                           .security_protected.plain.emm.header = {0},
                           .security_protected.plain.esm.header = {0}};

  OAILOG_FUNC_IN(LOG_NAS_EMM);
  OAILOG_INFO(LOG_NAS_EMM, "EMMAS-SAP - Send AS connection establish reject\n");

  /*
   * Setup the AS message
   */
  if (msg->eps_id.guti) {
    as_msg->s_tmsi.mme_code = msg->eps_id.guti->gummei.mme_code;
    as_msg->s_tmsi.m_tmsi   = msg->eps_id.guti->m_tmsi;
  } else {
    as_msg->ue_id = msg->ue_id;
  }

  /*
   * Setup the NAS security header
   */
  emm_msg = emm_as_set_header(&nas_msg, &msg->sctx);

  /*
   * Setup the NAS information messag
   */
  if (emm_msg) {
    switch (msg->nas_info) {
      case EMM_AS_NAS_INFO_ATTACH:
        size = emm_send_attach_reject(msg, &emm_msg->attach_reject);
        break;

      case EMM_AS_NAS_INFO_TAU:
        size = emm_send_tracking_area_update_reject(
            msg, &emm_msg->tracking_area_update_reject);
        break;

      case EMM_AS_NAS_INFO_SR:
        size =
            emm_send_service_reject(msg->emm_cause, &emm_msg->service_reject);
        break;

      default:
        OAILOG_WARNING(
            LOG_NAS_EMM,
            "EMMAS-SAP - Type of initial NAS "
            "message 0x%.2x is not valid\n",
            msg->nas_info);
    }
  }

  if (size > 0) {
    struct emm_context_s* emm_ctx                = NULL;
    emm_security_context_t* emm_security_context = NULL;
    ue_mm_context_t* ue_mm_context =
        mme_ue_context_exists_mme_ue_s1ap_id(msg->ue_id);

    if (ue_mm_context) {
      emm_ctx = &ue_mm_context->emm_context;
      if (emm_ctx) {
        if (IS_EMM_CTXT_PRESENT_SECURITY(emm_ctx)) {
          emm_security_context = &emm_ctx->_security;
          nas_msg.header.sequence_number =
              emm_security_context->dl_count.seq_num;
          OAILOG_DEBUG(
              LOG_NAS_EMM, "Set nas_msg.header.sequence_number -> %u\n",
              nas_msg.header.sequence_number);
        }
      }
    }

    /*
     * Encode the initial NAS information message
     */
    int bytes =
        emm_as_encode(&as_msg->nas_msg, &nas_msg, size, emm_security_context);
    if (bytes > 0) {
      // This is to indicate MME-APP to release the S1AP UE context after
      // sending the message.
      as_msg->err_code = AS_TERMINATED_NAS;
      OAILOG_FUNC_RETURN(LOG_NAS_EMM, AS_NAS_ESTABLISH_RSP);
    }
  }

  OAILOG_FUNC_RETURN(LOG_NAS_EMM, 0);
}
