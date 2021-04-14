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
  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    rc = amf_proc_registration_reject(ue_id, amf_cause);
  }
  amf_registration_request_ies_t* params = new (amf_registration_request_ies_t);
  /*
   * Message processing
   */
  /*
   * Get the 5GS Registration type
   */
  params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_RESERVED;
  if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_INITIAL) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
  } else if (
      msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_MOBILITY_UPDATING) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_MOBILITY_UPDATING;
  } else if (
      msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_PERODIC_UPDATING) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_PERODIC_UPDATING;
  } else if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_EMERGENCY) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_EMERGENCY;
  } else if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_RESERVED) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_RESERVED;
  } else {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
  }
  /* This is SUCI message identity type is SUPI as IMSI type
   * Extract the SUPI from SUCI directly as scheme is NULL */
  if (msg->m5gs_mobile_identity.mobile_identity.imsi.type_of_identity ==
      M5GSMobileIdentityMsg_IMSI) {
    // Only considering protection scheme as NULL else return error.
    if (msg->m5gs_mobile_identity.mobile_identity.imsi.protect_schm_id ==
        MOBILE_IDENTITY_PROTECTION_SCHEME_NULL) {
      /*
       * Extract the SUPI or IMSI from SUCI as scheme output is not encrypted
       */
      params->imsi = new (imsi_t);
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
      OAILOG_DEBUG(
          LOG_AMF_APP, "Value of SUPI/IMSI from params->imsi->u.value\n");
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "SUPI as IMSI derived : %02x%02x%02x%02x%02x%02x%02x%02x \n",
          params->imsi->u.value[0], params->imsi->u.value[1],
          params->imsi->u.value[2], params->imsi->u.value[3],
          params->imsi->u.value[4], params->imsi->u.value[5],
          params->imsi->u.value[6], params->imsi->u.value[7]);
    }
  }

  rc = amf_proc_registration_request(ue_id, is_amf_ctx_new, params);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
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
      M5GSMobileIdentityMsg_IMSI) {
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

    /* SUPI as IMSI retrived from SUSI. Generate GUTI based on incoming
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

/* Lookup the guti structure basd on imsi  in amf_supi_guti_map
 * params@ imsi
 * return guti structure on success ,Null on failure
 */
ue_m5gmm_context_s* lookup_ue_ctxt_by_imsi(imsi64_t imsi64) {
  //  imsi64_t imsi64 = amf_imsi_to_imsi64(imsi);

  /*Check imsi found
   *
   */
  std::unordered_map<imsi64_t, guti_and_amf_id_t>::iterator found_imsi =
      amf_supi_guti_map.find(imsi64);
  if (found_imsi == amf_supi_guti_map.end()) {
    // it is new entry to map
    OAILOG_ERROR(LOG_NAS_AMF, "UE_ID context not found in ue-context_map\n");
    return NULL;

  } else {
    // Overwrite the second element.
    return amf_ue_context_exists_amf_ue_ngap_id(
        found_imsi->second.amf_ue_ngap_id);
    // return found_imsi->second;
  }
}

}  // namespace magma5g
