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

  Source      amf_recv.cpp

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "conversions.h"
#ifdef __cplusplus
}
#endif
#include <unordered_map>
#include "M5GSMobileIdentity.h"
#include "M5GRegistrationAccept.h"
#include "amf_common_defs.h"
#include "amf_data.h"
#include "amf_fsm.h"
//#include "amf_nas5g_proc.h"
#include "amf_app_ue_context_and_proc.h"
#include "amf_asDefs.h"
#include "amf_as.h"
#include "amf_sap.h"
#include "5gs_registration_type.h"
#include "amf_recv.h"
#include "amf_identity.h"
#include "common_types.h"
using namespace std;
#define AMF_CAUSE_SUCCESS (1)
namespace magma5g {

amf_identity_msg amf_identity_rcv;
m5g_authentication m5g_authentication_rcv;
extern std::unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map;

int amf_procedure_handler::amf_handle_registration_request(
    amf_ue_ngap_id_t ue_id, tai_t* originating_tai, ecgi_t* originating_ecgi,
    RegistrationRequestMsg* msg, const bool is_initial,
    const bool is_amf_ctx_new, int amf_cause,
    const amf_nas_message_decode_status_t decode_status) {
  OAILOG_INFO(
      LOG_NAS_AMF, "AMF_TEST: Processing REGITRATION_REQUEST message\n");
  amf_registration_procedure amf_reg_proc;
  int rc = RETURNok;
  // Local imsi to be put in imsi defined in 3gpp_23.003.h
  supi_as_imsi_t supi_imsi;
#ifdef HANDLE_POST_MVC
  /*
   * Message checking
   */
  if (msg->uenetworkcapability.spare != 0b000) {
    /*
     * Spare bits shall be coded as zero
     */
    amf_cause = AMF_CAUSE_PROTOCOL_ERROR;
  }
#endif
  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    rc = amf_reg_proc.amf_proc_registration_reject(ue_id, amf_cause);
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

  /*
   * Get the AMF mobile identity
   */
#if 0   // TODO -  NEED-RECHECK
  if (msg->m5gs_mobile_identity.mobile_identity.guti.type_of_identity ==
      M5GSMobileIdentityMsg_GUTI) {
    // params->m5gsmobileidentity.m5gs_mobile_identity_t.mobileidentity.guti =
    // new(guti_t); // need to define guti_t like below in 3gpp 23003.h
    /*!< \brief  Globally Unique MME Identity gummei_t gummei;            */
    /*!< \brief  M-Temporary Mobile Subscriber Identity tmsi_t m_tmsi;  */

/* below need to update after header file of 5gsmobileidenty.h completion*/
                params->m5gsmobileidentity.mobileidentity.guti.mcc_digit1 = msg->oldgutiorimsi.guti.mcc_digit1;
                params->m5gsmobileidentity.mobileidentity.guti.mcc_digit2 = msg->oldgutiorimsi.guti.mcc_digit2;
                params->m5gsmobileidentity.mobileidentity.guti.mcc_digit3 = msg->oldgutiorimsi.guti.mcc_digit3;
                params->m5gsmobileidentity.mobileidentity.guti.mnc_digit1 = msg->oldgutiorimsi.guti.mnc_digit1;
                params->m5gsmobileidentity.mobileidentity.guti.plmn.mnc_digit2 = msg->oldgutiorimsi.guti.mnc_digit2;
                params->m5gsmobileidentity.mobileidentity.guti.plmn.mnc_digit3 = msg->oldgutiorimsi.guti.mnc_digit3;
                //params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti.mme_gid = msg->oldgutiorimsi.guti.mme_group_id;
                //params->m5gsmobileidentity.m5gs_mobile_identity_t.m5gguti.mme_code = msg->oldgutiorimsi.guti.mme_code;
                params->guti->m_tmsi = msg->oldgutiorimsi.guti.m_tmsi;
#endif  // TODO -  NEED-RECHECK

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
      // params->imsi->u.num.digit1 =
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

      OAILOG_INFO(
          LOG_AMF_APP,
          "AMF_TEST : printing SUPI/IMSI from params->imsi->u.value\n");
      // OAILOG_INFO(LOG_AMF_APP, "SUPI as IMSI derived : %02x %02x %02x %02x
      // %02x %02x %02x %02x",
      OAILOG_INFO(
          LOG_AMF_APP,
          "SUPI as IMSI derived : %02x%02x%02x%02x%02x%02x%02x%02x \n",
          params->imsi->u.value[0], params->imsi->u.value[1],
          params->imsi->u.value[2], params->imsi->u.value[3],
          params->imsi->u.value[4], params->imsi->u.value[5],
          params->imsi->u.value[6], params->imsi->u.value[7]);

#if 0  // remove after UT    
    params->imsi->u.num.digit7  = 1;
    params->imsi->u.num.digit8  = 0;
    params->imsi->u.num.digit9  = 0;
    params->imsi->u.num.digit10 = 0;
    params->imsi->u.num.digit11 = 0;
    params->imsi->u.num.digit12 = 0;
    params->imsi->u.num.digit13 = 0;
    params->imsi->u.num.digit14 = 0;
    params->imsi->u.num.digit15 = 1;
#endif
    }
  }

  // TODO -  other registration procedures are to be taken care later
  /*else if (msg->m5gsmobileidentity.mobileidentity.imei.typeofidentity ==
    M5GSMobileIdentityMsg_IMEI)
    {
        //assign IMEI value
    }
    else if (msg->m5gsmobileidentity.mobileidentity.m5gstmsi.typeofidentity ==
    M5GSMobileIdentityMsg_TMSI)
    {
        //assign m5gstmsi value
    }
    else if
    (msg->m5gsmobileidentity.m5gs_mobile_identity_t.imeisv.typeofidentity ==
    M5GS_Mobile_Identity_IMEISV)
    {
        //assign imeisv value
    REGISTRATION_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT
    }
  */
  /*
   * Execute the requested UE registration procedure
   */
  rc = amf_registration_procedure::amf_proc_registration_request(
      ue_id, is_amf_ctx_new, params);
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

int amf_procedure_handler::amf_handle_identity_response(
    amf_ue_ngap_id_t ue_id, M5GSMobileIdentityMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t decode_status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  OAILOG_INFO(LOG_NAS_AMF, "AMF_TEST: Received IDENTITY_RESPONSE message\n");
  int rc = RETURNerror;
  /*
   * Message processing
   */
  /*
   * Get the mobile identity
   */
  imsi_t imsi = {0}, *p_imsi = NULL;
  imei_t imei = {0}, *p_imei = NULL;
  imeisv_t imeisv = {0}, *p_imeisv = NULL;
  tmsi_t tmsi = 0, *p_tmsi = NULL;
  supi_as_imsi_t supi_imsi;
  amf_guti_m5g_t amf_guti;
  guti_and_amf_id_t guti_and_amf_id;
  guti_m5_t* amf_ctx_guti;

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
      OAILOG_INFO(
          LOG_AMF_APP, "AMF_TEST : Identity RSP SUPI/IMSI from imsi.u.value");
      OAILOG_INFO(
          LOG_AMF_APP,
          "SUPI as IMSI derived : %02x%02x%02x%02x%02x%02x%02x%02x",
          imsi.u.value[0], imsi.u.value[1], imsi.u.value[2], imsi.u.value[3],
          imsi.u.value[4], imsi.u.value[5], imsi.u.value[6], imsi.u.value[7]);

#if 0
       OAILOG_INFO(
           LOG_AMF_APP, "scheme_len: %d\n", msg->mobile_identity.imsi.scheme_len);
       for (int k = 0; k < 9; k++) {
         OAILOG_INFO(
             LOG_AMF_APP, "scheme output[k]:%x ",
             msg->mobile_identity.imsi.scheme_output[k]);
       }
       imsi.u.num.digit7  = 1;
       imsi.u.num.digit8  = 0;
       imsi.u.num.digit9  = 0;
       imsi.u.num.digit10 = 0;
       imsi.u.num.digit11 = 0;
       imsi.u.num.digit12 = 0;
       imsi.u.num.digit13 = 0;
       imsi.u.num.digit14 = 0;
       imsi.u.num.digit15 = 1;
#endif
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
    // amf_app_generate_guti_on_supi(&amf_guti, supi_imsi.plmn);
    amf_app_generate_guti_on_supi(&amf_guti, &supi_imsi);

    uint8_t* offset_plmn;  // Only testing
    offset_plmn = (uint8_t*) &amf_guti.guamfi.plmn;
    uint8_t* octet1;
    uint8_t* octet2;
    uint8_t* octet3;
    uint8_t* octet4;
    octet1 = offset_plmn;
    octet2 = offset_plmn++;
    octet3 = offset_plmn++;
    octet4 = offset_plmn++;
    uint16_t* set_id_pointer;
    set_id_pointer = (uint16_t*) ++offset_plmn;
    // set_id_pointer = (uint16_t*) &amf_guti.guamfi.amf_set_id;
    OAILOG_INFO(LOG_NAS_AMF, "AMF_TEST: Generated GUTI as below\n");
    OAILOG_INFO(
        LOG_NAS_AMF, "AMF_TEST: PLMN as 0x%02x%02x%02x\n", *octet1, *octet2,
        *octet3);
    OAILOG_INFO(LOG_NAS_AMF, "AMF_TEST: Region ID as 0x%02x\n", octet4);
    // OAILOG_INFO(LOG_NAS_AMF, "AMF_TEST: set_id and pointer as %04" PRIu16
    // "\n", *set_id_pointer);
    OAILOG_INFO(
        LOG_NAS_AMF, "AMF_TEST: set_id and pointer as 0x%02x\n",
        *set_id_pointer);
    OAILOG_INFO(
        LOG_NAS_AMF, "AMF_TEST: 5G-TMSI as 0x%08" PRIx32 "\n", amf_guti.m_tmsi);
    //===== Only debugging purpose and will be deleted later=======
    uint8_t plmn3ch[3];
    memset(plmn3ch, 0, 3);
    memcpy(&plmn3ch, &amf_guti, 3);
    OAILOG_INFO(
        LOG_AMF_APP,
        "AMF_TEST: PLMN from GUTI generated function 3 octets  0x%02x%02x%02x",
        plmn3ch[0], plmn3ch[1], plmn3ch[2]);

    /* Need to store guti in amf_ctx as well for quick access
     * which will be used to send in DL message during registration
     * accept message
     * TODO Note:currently adapting the way
     */
    amf_ctx_guti = (guti_m5_t*) &amf_guti;

    /* store this GUTI in
     * unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map
     */
    imsi64_t imsi64 = amf_imsi_to_imsi64(&imsi);

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

#if 0
   else if (
      msg->mobile_identity.imei.type_of_identity == MOBILE_IDENTITY_IMEI) {
    /*
     * Get the IMEI
  
   */
  } else if (
      msg->mobile_identity.tmsi.type_of_identity == MOBILE_IDENTITY_TMSI) {
    /*
     * Get the TMSI
     */

  }
#endif
  /*
   * Execute the identification completion procedure
   */
  rc = amf_identity_rcv.amf_proc_identification_complete(
      ue_id, p_imsi, p_imei, p_imeisv, (uint32_t*) (p_tmsi), amf_ctx_guti);
  // ue_id, p_imsi, p_imei, p_imeisv, (uint32_t*) (p_tmsi));
  OAILOG_FUNC_RETURN(LOG_NAS_EMM, rc);
}
#if 1
/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_authentication_response()                        **
 **                                                                        **
 ** Description: Processes Authentication Response message                 **
 **                                                                        **
 ** Inputs:  ue_id:      UE lower layer identifier                  **
 **      msg:       The received AMF message                   **
 **      Others:    None                                       **
 **                                                                        **
 ** Outputs:     amf_cause: AMF cause code                             **
 **      Return:    RETURNok, RETURNerror                      **
 **      Others:    None                                       **
 **                                                                      **
 ***************************************************************************/
int amf_procedure_handler::amf_handle_authentication_response(
    amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
    amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;

  OAILOG_INFO(
      LOG_NAS_AMF, "AMF_TEST: Received AUTHENTICATION_RESPONSE message\n");

  /*
   * Message checking
   */
  if (NULL == msg->autn_response_parameter.response_parameter) {
    /*
     * RES parameter shall not be null
     */
    OAILOG_INFO(
        LOG_AMF_APP, "AMF-TEST: RES parameter null, from %s\n", __FUNCTION__);
    amf_cause = AMF_CAUSE_INVALID_MANDATORY_INFO;
  }
  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    OAILOG_INFO(
        LOG_AMF_APP, "AMF-TEST: != AMF_CAUSE_SUCCESS, from %s\n", __FUNCTION__);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  /*
   * Execute the authentication completion procedure
   */
  rc = m5g_authentication_rcv.amf_proc_authentication_complete(
      ue_id, msg, AMF_CAUSE_SUCCESS,
      msg->autn_response_parameter.response_parameter);
  /*
   * Free authenticationresponseparameter IE
   */
  // bdestroy(msg->autn_response_parameter);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}
#endif  // TODO AMF-TEST, deal with later
}  // namespace magma5g
