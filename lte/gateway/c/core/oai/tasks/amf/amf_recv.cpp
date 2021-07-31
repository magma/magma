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

#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif
#include "common_defs.h"
#include <unordered_map>
#include "amf_app_ue_context_and_proc.h"
#include "amf_authentication.h"
#include "amf_as.h"
#include "amf_recv.h"
#include "amf_identity.h"

#define AMF_CAUSE_SUCCESS (1)
namespace magma5g {
extern std::unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map;

int amf_handle_service_request(
    amf_ue_ngap_id_t ue_id, ServiceRequestMsg* msg,
    const amf_nas_message_decode_status_t decode_status) {
  int rc                         = RETURNok;
  ue_m5gmm_context_s* ue_context = nullptr;
  notify_ue_event notify_ue_event_type;
  tmsi_t tmsi_rcv;

  memcpy(
      &tmsi_rcv, &msg->m5gs_mobile_identity.mobile_identity.tmsi.m5g_tmsi,
      sizeof(tmsi_t));

  ue_context = ue_context_loopkup_by_guti(tmsi_rcv);
  if ((ue_context) && (ue_id != ue_context->amf_ue_ngap_id)) {
    ue_context_update_ue_id(ue_context, ue_id);
  }

  if (ue_context) {
    OAILOG_INFO(
        LOG_NAS_AMF,
        "TMSI matched for the UE id %d "
        " received TMSI %08X\n",
        ue_id, tmsi_rcv);
    // Calculate MAC and compare if matches send message to SMF
    if (decode_status.mac_matched) {
      OAILOG_INFO(
          LOG_NAS_AMF, "MAC in security header matched for the UE id %d ",
          ue_id);
      // Set event type as service request
      notify_ue_event_type = UE_SERVICE_REQUEST_ON_PAGING;
      // construct the proto structure and send message to SMF
      rc = amf_smf_notification_send(ue_id, ue_context, notify_ue_event_type);
    }  // MAC matched
    else {
      OAILOG_INFO(
          LOG_NAS_AMF,
          "MAC in security header not matched for the UE id %d "
          "and prepare for reject message on DL",
          ue_id);
    }
  } else {
    OAILOG_INFO(
        LOG_NAS_AMF,
        "TMSI not matched for "
        "(ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
        ue_id);

    // Send prepare and send reject message.
  }
  return rc;
}

/* Identifies 5GS Registration type and processes the Message accordingly */
int amf_handle_registration_request(
    amf_ue_ngap_id_t ue_id, tai_t* originating_tai, ecgi_t* originating_ecgi,
    RegistrationRequestMsg* msg, const bool is_initial,
    const bool is_amf_ctx_new, int amf_cause,
    const amf_nas_message_decode_status_t decode_status) {
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  int rc = RETURNok;
  // Local imsi to be put in imsi defined in 3gpp_23.003.h
  supi_as_imsi_t supi_imsi;
  amf_guti_m5g_t amf_guti;
  guti_and_amf_id_t guti_and_amf_id;
  /*
   * Handle message checking error
   */
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    rc = amf_proc_registration_reject(ue_id, amf_cause);
    OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  }
  amf_registration_request_ies_t* params =
      new (amf_registration_request_ies_t)();
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  /*
   * Message processing
   */
  /*
   * Get the 5GS Registration type
   */
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_RESERVED;
  if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_INITIAL) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
  } else if (
      msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_MOBILITY_UPDATING) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_MOBILITY_UPDATING;
  } else if (
      msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_PERIODIC_UPDATING) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_PERIODIC_UPDATING;
  } else if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_EMERGENCY) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_EMERGENCY;
  } else if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_RESERVED) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_RESERVED;
  } else {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
  }
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context == NULL) {
    OAILOG_INFO(LOG_AMF_APP, "ue_context is NULL for UE ID:%d \n", ue_id);
    return RETURNerror;
  }
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  // Save the UE Security Capability into AMF's UE Context
  memcpy(
      &(ue_context->amf_context.ue_sec_capability), &(msg->ue_sec_capability),
      sizeof(UESecurityCapabilityMsg));
  memcpy(
      &(ue_context->amf_context.originating_tai), (const void*) originating_tai,
      sizeof(tai_t));

  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  OAILOG_DEBUG(
      LOG_NAS_AMF, "m5gs_reg_type.type_val :%d", msg->m5gs_reg_type.type_val);
  if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_INITIAL) {
    OAILOG_INFO(LOG_NAS_AMF, "New REGITRATION_REQUEST processing\n");
    /*
     * Get the AMF mobile identity. For new registration
     * mobility type suppose to be SUCI
     * This is SUCI message identity type is SUPI as IMSI type
     * Extract the SUPI from SUCI directly as scheme is NULL */
    if (msg->m5gs_mobile_identity.mobile_identity.imsi.type_of_identity ==
        M5GSMobileIdentityMsg_SUCI_IMSI) {
      // Only considering protection scheme as NULL else return error.
      if (msg->m5gs_mobile_identity.mobile_identity.imsi.protect_schm_id ==
          MOBILE_IDENTITY_PROTECTION_SCHEME_NULL) {
        /*
         * Extract the SUPI or IMSI from SUCI as scheme output is not encrypted
         */
        params->imsi = new imsi_t();
        /* Copying PLMN to local supi which is imsi*/
        supi_imsi.plmn.mcc_digit1 =
            msg->m5gs_mobile_identity.mobile_identity.imsi.mcc_digit1;
        supi_imsi.plmn.mcc_digit2 =
            msg->m5gs_mobile_identity.mobile_identity.imsi.mcc_digit2;
        supi_imsi.plmn.mcc_digit3 =
            msg->m5gs_mobile_identity.mobile_identity.imsi.mcc_digit3;
        supi_imsi.plmn.mnc_digit1 =
            msg->m5gs_mobile_identity.mobile_identity.imsi.mnc_digit1;
        supi_imsi.plmn.mnc_digit2 =
            msg->m5gs_mobile_identity.mobile_identity.imsi.mnc_digit2;
        supi_imsi.plmn.mnc_digit3 =
            msg->m5gs_mobile_identity.mobile_identity.imsi.mnc_digit3;
        // copy 5 octet scheme_output to msin of supi_imsi
        memcpy(
            &supi_imsi.msin,
            &msg->m5gs_mobile_identity.mobile_identity.imsi.scheme_output,
            MSIN_MAX_LENGTH);
        // Copy entire supi_imsi to param->imsi->u.value
        memcpy(&params->imsi->u.value, &supi_imsi, IMSI_BCD8_SIZE);

        if (supi_imsi.plmn.mnc_digit3 != 0xf) {
          params->imsi->u.value[0] = ((supi_imsi.plmn.mcc_digit1 << 4) & 0xf0) |
                                     (supi_imsi.plmn.mcc_digit2 & 0xf);
          params->imsi->u.value[1] = ((supi_imsi.plmn.mcc_digit3 << 4) & 0xf0) |
                                     (supi_imsi.plmn.mnc_digit1 & 0xf);
          params->imsi->u.value[2] = ((supi_imsi.plmn.mnc_digit2 << 4) & 0xf0) |
                                     (supi_imsi.plmn.mnc_digit3 & 0xf);
        }
        OAILOG_INFO(
            LOG_AMF_APP, "Value of SUPI/IMSI from params->imsi->u.value\n");
        OAILOG_INFO(
            LOG_AMF_APP,
            "SUPI as IMSI derived : %02x%02x%02x%02x%02x%02x%02x%02x \n",
            params->imsi->u.value[0], params->imsi->u.value[1],
            params->imsi->u.value[2], params->imsi->u.value[3],
            params->imsi->u.value[4], params->imsi->u.value[5],
            params->imsi->u.value[6], params->imsi->u.value[7]);

        ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit1 =
            supi_imsi.plmn.mcc_digit1;
        ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit2 =
            supi_imsi.plmn.mcc_digit2;
        ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit3 =
            supi_imsi.plmn.mcc_digit3;
        ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit1 =
            supi_imsi.plmn.mnc_digit1;
        ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit2 =
            supi_imsi.plmn.mnc_digit2;
        ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit3 =
            supi_imsi.plmn.mnc_digit3;

        ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_SUCI_IMSI;

        amf_app_generate_guti_on_supi(&amf_guti, &supi_imsi);

        amf_ue_context_on_new_guti(ue_context, (guti_m5_t*) &amf_guti);

        ue_context->amf_context.m5_guti.m_tmsi = amf_guti.m_tmsi;

        imsi64_t imsi64                = amf_imsi_to_imsi64(params->imsi);
        guti_and_amf_id.amf_guti       = amf_guti;
        guti_and_amf_id.amf_ue_ngap_id = ue_id;
        OAILOG_DEBUG(LOG_AMF_APP, "imsi64 : " IMSI_64_FMT "\n", imsi64);
        if (amf_supi_guti_map.size() == 0) {
          // first entry.
          amf_supi_guti_map.insert(
              std::pair<imsi64_t, guti_and_amf_id_t>(imsi64, guti_and_amf_id));
        } else {
          /* already elements exist then check if same imsi already present
           * if same imsi then update/overwrite the element
           */
          std::unordered_map<imsi64_t, guti_and_amf_id_t>::iterator found_imsi =
              amf_supi_guti_map.find(imsi64);
          if (found_imsi == amf_supi_guti_map.end()) {
            // it is new entry to map
            amf_supi_guti_map.insert(std::pair<imsi64_t, guti_and_amf_id_t>(
                imsi64, guti_and_amf_id));
          } else {
            // Overwrite the second element.
            found_imsi->second = guti_and_amf_id;
          }
        }
      }
    } else if (
        msg->m5gs_mobile_identity.mobile_identity.guti.type_of_identity ==
        M5GSMobileIdentityMsg_GUTI) {
      OAILOG_INFO(LOG_NAS_AMF, "New REGITRATION_REQUEST Id is GUTI\n");
      params->guti                        = new (guti_m5_t)();
      ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_GUTI;
    }
  }  // end of AMF_REGISTRATION_TYPE_INITIAL
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");

  if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_PERIODIC_UPDATING) {
    /*
     * This request for periodic registration update
     * For registered UE, is_amf_ctx_new = False
     * and Identity type should be GUTI
     *    1> Already PLMN had been updated in amf_context
     *    2> Generate new GUTI
     *    3> Update amf_context and MAP
     *    4> call accept message API and send DL message.
     */
    OAILOG_INFO(
        LOG_NAS_AMF,
        "AMF_REGISTRATION_TYPE_PERIODIC_UPDATING processing"
        " is_amf_ctx_new = %d and identity type = %d ",
        is_amf_ctx_new,
        msg->m5gs_mobile_identity.mobile_identity.imsi.type_of_identity);
    if ((msg->m5gs_mobile_identity.mobile_identity.guti.type_of_identity ==
         M5GSMobileIdentityMsg_GUTI)) {
      /* Copying PLMN to local supi which is imsi*/
      supi_imsi.plmn.mcc_digit1 =
          msg->m5gs_mobile_identity.mobile_identity.guti.mcc_digit1;
      supi_imsi.plmn.mcc_digit2 =
          msg->m5gs_mobile_identity.mobile_identity.guti.mcc_digit2;
      supi_imsi.plmn.mcc_digit3 =
          msg->m5gs_mobile_identity.mobile_identity.guti.mcc_digit3;
      supi_imsi.plmn.mnc_digit1 =
          msg->m5gs_mobile_identity.mobile_identity.guti.mnc_digit1;
      supi_imsi.plmn.mnc_digit2 =
          msg->m5gs_mobile_identity.mobile_identity.guti.mnc_digit2;
      supi_imsi.plmn.mnc_digit3 =
          msg->m5gs_mobile_identity.mobile_identity.guti.mnc_digit3;

      amf_app_generate_guti_on_supi(&amf_guti, &supi_imsi);
      OAILOG_INFO(
          LOG_NAS_AMF,
          "In process of periodic registration update"
          " new 5G-TMSI value 0x%08" PRIx32 "\n",
          amf_guti.m_tmsi);
      /* Update this new GUTI in amf_context and map
       * unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map
       */
      amf_ue_context_on_new_guti(ue_context, (guti_m5_t*) &amf_guti);
      ue_context->amf_context.m5_guti.m_tmsi = amf_guti.m_tmsi;

      imsi64_t imsi64                = ue_context->amf_context.imsi64;
      guti_and_amf_id.amf_guti       = amf_guti;
      guti_and_amf_id.amf_ue_ngap_id = ue_id;
      // Find the respective element in map with key imsi_64
      std::unordered_map<imsi64_t, guti_and_amf_id_t>::iterator found_imsi =
          amf_supi_guti_map.find(imsi64);
      if (found_imsi != amf_supi_guti_map.end()) {
        // element found in map and update the GUTI
        found_imsi->second = guti_and_amf_id;
      }

      params->guti = new (guti_m5_t)();
      memcpy(
          params->guti, &(ue_context->amf_context.m5_guti), sizeof(guti_m5_t));

      ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_GUTI;

    } else {
      // UE context is new and/or UE identity type is not GUTI
      // add log message.
      OAILOG_INFO(
          LOG_AMF_APP,
          "UE context was not existing or UE identity type is not GUTI "
          "Periodic Registration Update failed and sending reject message\n");
      // TODO Implement Reject message
      return RETURNerror;
    }
  }  // end of AMF_REGISTRATION_TYPE_PERIODIC_UPDATING

  params->decode_status = decode_status;
  /*
   * Execute the requested new UE registration procedure
   * This will initiate identity req in DL.
   */
  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  rc = amf_proc_registration_request(ue_id, is_amf_ctx_new, params);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);

  OAILOG_DEBUG(LOG_NAS_AMF, "Processing REGITRATION_REQUEST message\n");
  return rc;
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_identity_response()                                **
 **                                                                        **
 ** Description: Processes Identity Response message                       **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      msg:       The received AMF message                               **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                                 **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/

int amf_handle_identity_response(
    amf_ue_ngap_id_t ue_id, M5GSMobileIdentityMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_DEBUG(LOG_NAS_AMF, "Received IDENTITY_RESPONSE message\n");
  int rc = RETURNerror;
  /*
   * Message processing
   */
  /*
   * Get the mobile identity
   */
  imsi_t imsi = {0}, *p_imsi = NULL;
  imei_t* p_imei     = NULL;
  imeisv_t* p_imeisv = NULL;
  tmsi_t* p_tmsi     = NULL;
  supi_as_imsi_t supi_imsi;
  amf_guti_m5g_t amf_guti;
  guti_and_amf_id_t guti_and_amf_id;
  guti_m5_t* amf_ctx_guti = NULL;

  /* This is SUCI message identity type is SUPI as IMSI type
   * Extract the SUPI from SUCI directly as scheme is NULL */
  if (msg->mobile_identity.imsi.type_of_identity ==
      M5GSMobileIdentityMsg_SUCI_IMSI) {
    // Only considering protection scheme as NULL else return error.
    if (msg->mobile_identity.imsi.protect_schm_id ==
        MOBILE_IDENTITY_PROTECTION_SCHEME_NULL) {
      /*
       * Extract the SUPI or IMSI from SUCI as scheme output is not encrypted
       */
      p_imsi                    = &imsi;
      supi_imsi.plmn.mcc_digit1 = msg->mobile_identity.imsi.mcc_digit1;
      supi_imsi.plmn.mcc_digit2 = msg->mobile_identity.imsi.mcc_digit2;
      supi_imsi.plmn.mcc_digit3 = msg->mobile_identity.imsi.mcc_digit3;
      supi_imsi.plmn.mnc_digit1 = msg->mobile_identity.imsi.mnc_digit1;
      supi_imsi.plmn.mnc_digit2 = msg->mobile_identity.imsi.mnc_digit2;
      supi_imsi.plmn.mnc_digit3 = msg->mobile_identity.imsi.mnc_digit3;
      // copy 5 octet scheme_output to msin of supi_imsi
      memcpy(
          &supi_imsi.msin, &msg->mobile_identity.imsi.scheme_output,
          MSIN_MAX_LENGTH);
      // Copy entire supi_imsi to imsi.u.value which is 8 bytes
      memcpy(&imsi.u.value, &supi_imsi, IMSI_BCD8_SIZE);

      if (supi_imsi.plmn.mnc_digit3 != 0xf) {
        imsi.u.value[0] = ((supi_imsi.plmn.mcc_digit1 << 4) & 0xf0) |
                          (supi_imsi.plmn.mcc_digit2 & 0xf);
        imsi.u.value[1] = ((supi_imsi.plmn.mcc_digit3 << 4) & 0xf0) |
                          (supi_imsi.plmn.mnc_digit1 & 0xf);
        imsi.u.value[2] = ((supi_imsi.plmn.mnc_digit2 << 4) & 0xf0) |
                          (supi_imsi.plmn.mnc_digit3 & 0xf);
      }

    } else {
      /* Mobile identity is SUPI type IMSI but Protection scheme is not NULL
       * which is not valid message from UE. Return from here after
       * printing error message
       */
      OAILOG_ERROR(
          LOG_AMF_APP,
          "Invalid protection scheme  received "
          " in identity response from UE \n");
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
    }

    /* SUPI as IMSI retrieved from SUSI. Generate GUTI based on incoming
     * PLMN and amf_config file and fill the SUPI-GUTI MAP
     * Note: GUTI generation supposed to be done after verifying
     * subscriber data with AUSF/UDM. But currently we are not supporting
     * AUSF/UDM and directly generating GUTI.
     * If authentication request rejected by UE, the MAP has to be cleared
     */
    amf_app_generate_guti_on_supi(&amf_guti, &supi_imsi);
    OAILOG_DEBUG(LOG_NAS_AMF, "5G-TMSI as 0x%08" PRIx32 "\n", amf_guti.m_tmsi);

    /* Need to store guti in amf_ctx as well for quick access
     * which will be used to send in DL message during registration
     * accept message
     * TODO Note:currently adapting the way
     */
    amf_ctx_guti = (guti_m5_t*) &amf_guti;

    /* store this GUTI in
     * unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map
     */
    imsi64_t imsi64                = amf_imsi_to_imsi64(&imsi);
    guti_and_amf_id.amf_guti       = amf_guti;
    guti_and_amf_id.amf_ue_ngap_id = ue_id;

    if (amf_supi_guti_map.size() == 0) {
      // first entry.
      amf_supi_guti_map.insert(
          std::pair<imsi64_t, guti_and_amf_id_t>(imsi64, guti_and_amf_id));
    } else {
      /* already elements exist then check if same imsi already present
       * if same imsi then update/overwrite the element
       */
      std::unordered_map<imsi64_t, guti_and_amf_id_t>::iterator found_imsi =
          amf_supi_guti_map.find(imsi64);
      if (found_imsi == amf_supi_guti_map.end()) {
        // it is new entry to map
        amf_supi_guti_map.insert(
            std::pair<imsi64_t, guti_and_amf_id_t>(imsi64, guti_and_amf_id));
      } else {
        // Overwrite the second element.
        found_imsi->second = guti_and_amf_id;
      }
    }

    ue_m5gmm_context_s* ue_context =
        amf_ue_context_exists_amf_ue_ngap_id(ue_id);
    if (ue_context) {
      ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_SUCI_IMSI;
      amf_ue_context_on_new_guti(ue_context, (guti_m5_t*) &amf_guti);
    }
  }
  /*
   * Execute the identification completion procedure
   */
  rc = amf_proc_identification_complete(
      ue_id, p_imsi, p_imei, p_imeisv, (uint32_t*) (p_tmsi), amf_ctx_guti);
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_authentication_response()                          **
 **                                                                        **
 ** Description: Processes Authentication Response message                 **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      msg:       The received AMF message                               **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                                 **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_handle_authentication_response(
    amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;
  OAILOG_DEBUG(LOG_NAS_AMF, "Received AUTHENTICATION_RESPONSE message\n");
  /*
   * Message checking
   */
  if (msg->autn_response_parameter.response_parameter == NULL) {
    /*
     * RES parameter shall not be null
     */
    OAILOG_DEBUG(LOG_AMF_APP, "Response parameter is null");
    amf_cause = AMF_CAUSE_INVALID_MANDATORY_INFO;
  }
  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }
  /*
   * Execute the authentication completion procedure
   */
  rc = amf_proc_authentication_complete(
      ue_id, msg, AMF_CAUSE_SUCCESS,
      msg->autn_response_parameter.response_parameter);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_authentication_failure()                           **
 **                                                                        **
 ** Description: Processes Authentication failure  message                 **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                         **
 **      msg:       The received AMF message                               **
 **      Others:    None                                                   **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                                 **
 **      Return:    RETURNok, RETURNerror                                  **
 **      Others:    None                                                   **
 **                                                                        **
 ***************************************************************************/
int amf_handle_authentication_failure(
    amf_ue_ngap_id_t ue_id, AuthenticationFailureMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;
  OAILOG_DEBUG(LOG_NAS_AMF, "Received AUTHENTICATION_FAILURE message\n");

  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  /*
   * Execute the authentication failure procedure
   */
  rc = amf_proc_authentication_failure(ue_id, msg, AMF_CAUSE_SUCCESS);

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name:    lookup_ue_ctxt_by_imsi()                                      **
 **                                                                        **
 ** Description: Lookup the guti structure based on imsi in                **
 **              amf_supi_guti_map                                         **
 **                                                                        **
 ** Inputs:  imsi64: imsi value                                            **
 **                                                                        **
 ** Outputs: ue_m5gmm_context_s: pointer to ue context                     **
 **                                                                        **
 ***************************************************************************/
ue_m5gmm_context_s* lookup_ue_ctxt_by_imsi(imsi64_t imsi64) {
  /*Check imsi found
   *
   */
  std::unordered_map<imsi64_t, guti_and_amf_id_t>::iterator found_imsi =
      amf_supi_guti_map.find(imsi64);
  if (found_imsi == amf_supi_guti_map.end()) {
    return NULL;
  } else {
    return amf_ue_context_exists_amf_ue_ngap_id(
        found_imsi->second.amf_ue_ngap_id);
  }
}

}  // namespace magma5g
