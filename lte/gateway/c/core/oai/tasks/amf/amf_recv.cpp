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
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include <unordered_map>
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.h"
#include "orc8r/gateway/c/common/service303/includes/MetricsHelpers.h"

#define AMF_CAUSE_SUCCESS (1)
#define AMF_CAUSE_UE_SEC_CAP_MISSMATCH (23)
namespace magma5g {
extern std::unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map;

int amf_handle_service_request(
    amf_ue_ngap_id_t ue_id, ServiceRequestMsg* msg,
    const amf_nas_message_decode_status_t decode_status) {
  int rc                         = RETURNok;
  ue_m5gmm_context_s* ue_context = nullptr;
  notify_ue_event notify_ue_event_type;
  amf_sap_t amf_sap;
  tmsi_t tmsi_rcv;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  char ip_str[INET_ADDRSTRLEN];
  uint16_t pdu_session_status = 0;
  uint32_t tmsi_stored;
  paging_context_t* paging_ctx = nullptr;
  guti_and_amf_id_t guti_and_amf_id;

  OAILOG_DEBUG(
      LOG_AMF_APP, "Received TMSI in message : %02x%02x%02x%02x",
      msg->m5gs_mobile_identity.mobile_identity.tmsi.m5g_tmsi[0],
      msg->m5gs_mobile_identity.mobile_identity.tmsi.m5g_tmsi[1],
      msg->m5gs_mobile_identity.mobile_identity.tmsi.m5g_tmsi[2],
      msg->m5gs_mobile_identity.mobile_identity.tmsi.m5g_tmsi[3]);
  memcpy(
      &tmsi_rcv, &msg->m5gs_mobile_identity.mobile_identity.tmsi.m5g_tmsi,
      sizeof(tmsi_t));
  memset(&amf_sap, 0, sizeof(amf_sap_s));
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  tmsi_rcv   = ntohl(tmsi_rcv);

  if (ue_context == NULL) {
    OAILOG_DEBUG(LOG_AMF_APP, "UE context is NULL\n");
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  tmsi_stored = ue_context->amf_context.m5_guti.m_tmsi;
  // Check TMSI and then check MAC
  OAILOG_DEBUG(
      LOG_NAS_AMF, " TMSI stored in AMF CONTEXT %08" PRIx32 "\n", tmsi_stored);
  OAILOG_DEBUG(LOG_NAS_AMF, " TMSI received %08" PRIx32 "\n", tmsi_rcv);

  if (ue_context && (tmsi_rcv == tmsi_stored)) {
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "TMSI matched for UE ID " AMF_UE_NGAP_ID_FMT
        " receved TMSI %08X stored TMSI %08X \n",
        ue_id, tmsi_rcv, tmsi_stored);

    paging_ctx = &ue_context->paging_context;

    if ((paging_ctx) &&
        (NAS5G_TIMER_INACTIVE_ID != paging_ctx->m5_paging_response_timer.id) &&
        (0 != paging_ctx->m5_paging_response_timer.id)) {
      amf_app_stop_timer(paging_ctx->m5_paging_response_timer.id);
      paging_ctx->m5_paging_response_timer.id = NAS5G_TIMER_INACTIVE_ID;
      paging_ctx->paging_retx_count           = 0;
      // Fill the itti msg based on context info produced in amf core
      OAILOG_DEBUG(LOG_AMF_APP, "T3513: After stopping PAGING Timer\n");
    }

    imsi64_t imsi64                 = ue_context->amf_context.imsi64;
    guti_and_amf_id.amf_guti.m_tmsi = ue_context->amf_context.m5_guti.m_tmsi;
    guti_and_amf_id.amf_guti.guamfi = ue_context->amf_context.m5_guti.guamfi;
    guti_and_amf_id.amf_ue_ngap_id  = ue_id;
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

    if (msg->service_type.service_type_value == SERVICE_TYPE_SIGNALING) {
      OAILOG_DEBUG(LOG_NAS_AMF, "Service request type is signalling \n");
      amf_sap.primitive = AMFAS_ESTABLISH_CNF;

      amf_sap.u.amf_as.u.establish.ue_id    = ue_id;
      amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_SR;
      rc                                    = amf_sap_send(&amf_sap);
      ue_context->mm_state                  = REGISTERED_CONNECTED;
    } else if (
        (msg->service_type.service_type_value == SERVICE_TYPE_DATA) ||
        (msg->service_type.service_type_value ==
         SERVICE_TYPE_HIGH_PRIORITY_ACCESS)) {
      if ((msg->service_type.service_type_value == SERVICE_TYPE_DATA) &&
          !(msg->uplink_data_status.uplinkDataStatus)) {
        // prepare and send reject message.
        OAILOG_INFO(
            LOG_NAS_AMF,
            "Sending service reject with cuase condtional IE missing\n");
        amf_sap.primitive                     = AMFAS_ESTABLISH_REJ;
        amf_sap.u.amf_as.u.establish.ue_id    = ue_id;
        amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_SR;
        if (msg->pdu_session_status.iei) {
          amf_sap.u.amf_as.u.establish.pdu_session_status_ie =
              AMF_AS_PDU_SESSION_STATUS;
          amf_sap.u.amf_as.u.establish.pdu_session_status =
              msg->pdu_session_status.pduSessionStatus;
        }
        amf_sap.u.amf_as.u.establish.amf_cause =
            AMF_CAUSE_CONDITIONAL_IE_MISSING;
        rc = amf_sap_send(&amf_sap);
      } else {
        OAILOG_DEBUG(LOG_NAS_AMF, "Service request type is Data \n");
        for (uint16_t session_id = 1; session_id < (sizeof(session_id) * 8);
             session_id++) {
          if (msg->uplink_data_status.uplinkDataStatus & (1 << session_id)) {
            std::shared_ptr<smf_context_t> smf_context =
                amf_get_smf_context_by_pdu_session_id(ue_context, session_id);
            if (smf_context) {
              pdu_session_status |= (1 << session_id);
              IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
              if (smf_context->pdu_address.pdn_type == IPv4) {
                inet_ntop(
                    AF_INET, &(smf_context->pdu_address.ipv4_address.s_addr),
                    ip_str, INET_ADDRSTRLEN);
              }

              OAILOG_DEBUG(
                  LOG_NAS_AMF,
                  "Sending session request to SMF on service request for "
                  "sessiond %u\n",
                  session_id);
              notify_ue_event_type = UE_SERVICE_REQUEST_ON_PAGING;
              // construct the proto structure and send message to SMF
              amf_smf_notification_send(
                  ue_id, ue_context, notify_ue_event_type);
            }
          }
        }
      }
    } else if (
        msg->service_type.service_type_value ==
        SERVICE_TYPE_MOBILE_TERMINATED_SERVICES) {
      notify_ue_event_type = UE_SERVICE_REQUEST_ON_PAGING;
      amf_smf_notification_send(ue_id, ue_context, notify_ue_event_type);
    }
  } else {
    OAILOG_WARNING(
        LOG_NAS_AMF, "TMSI not matched for ue_id=" AMF_UE_NGAP_ID_FMT ")\n",
        ue_id);

    // Send prepare and send reject message.
    amf_sap.primitive                     = AMFAS_ESTABLISH_REJ;
    amf_sap.u.amf_as.u.establish.ue_id    = ue_id;
    amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_SR;

    if (msg->pdu_session_status.iei) {
      amf_sap.u.amf_as.u.establish.pdu_session_status_ie =
          AMF_AS_PDU_SESSION_STATUS;
      amf_sap.u.amf_as.u.establish.pdu_session_status =
          msg->pdu_session_status.pduSessionStatus;
    }
    amf_sap.u.amf_as.u.establish.amf_cause = AMF_CAUSE_UE_ID_CAN_NOT_BE_DERIVED;
    rc                                     = amf_sap_send(&amf_sap);
  }

  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

void amf_copy_plmn_to_supi(
    const ImsiM5GSMobileIdentity& imsi, supi_as_imsi_t& supi_imsi) {
  supi_imsi.plmn.mcc_digit1 = imsi.mcc_digit1;
  supi_imsi.plmn.mcc_digit2 = imsi.mcc_digit2;
  supi_imsi.plmn.mcc_digit3 = imsi.mcc_digit3;

  supi_imsi.plmn.mnc_digit1 = imsi.mnc_digit1;
  supi_imsi.plmn.mnc_digit2 = imsi.mnc_digit2;
  supi_imsi.plmn.mnc_digit3 = imsi.mnc_digit3;
}

int amf_copy_plmn_to_context(
    const ImsiM5GSMobileIdentity& imsi, ue_m5gmm_context_s* ue_context) {
  if (ue_context == NULL) {
    OAILOG_ERROR(LOG_AMF_APP, "UE context is null");
    return RETURNerror;
  }

  ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit1 = imsi.mcc_digit1;
  ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit2 = imsi.mcc_digit2;
  ue_context->amf_context.m5_guti.guamfi.plmn.mcc_digit3 = imsi.mcc_digit3;
  ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit1 = imsi.mnc_digit1;
  ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit2 = imsi.mnc_digit2;
  ue_context->amf_context.m5_guti.guamfi.plmn.mnc_digit3 = imsi.mnc_digit3;
  return RETURNok;
}

/* Identifies 5GS Registration type and processes the Message accordingly */
int amf_handle_registration_request(
    amf_ue_ngap_id_t ue_id, tai_t* originating_tai, ecgi_t* originating_ecgi,
    RegistrationRequestMsg* msg, const bool is_initial,
    const bool is_amf_ctx_new, int amf_cause,
    const amf_nas_message_decode_status_t decode_status) {
  int rc = RETURNok;
  // Local imsi to be put in imsi defined in 3gpp_23.003.h
  supi_as_imsi_t supi_imsi;
  amf_guti_m5g_t amf_guti;
  guti_and_amf_id_t guti_and_amf_id;
  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    amf_cause = AMF_CAUSE_SUCCESS;
    rc        = amf_proc_registration_reject(ue_id, amf_cause);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
  }
  amf_registration_request_ies_t* params =
      new (amf_registration_request_ies_t)();
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
      msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_PERIODIC_UPDATING) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_PERIODIC_UPDATING;
  } else if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_EMERGENCY) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_EMERGENCY;
  } else if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_RESERVED) {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_RESERVED;
  } else {
    params->m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
  }

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context == NULL) {
    OAILOG_ERROR(
        LOG_AMF_APP, "UE context is null for UE ID: " AMF_UE_NGAP_ID_FMT,
        ue_id);
    return RETURNerror;
  }
  // Save the UE Security Capability into AMF's UE Context
  memcpy(
      &(ue_context->amf_context.ue_sec_capability), &(msg->ue_sec_capability),
      sizeof(UESecurityCapabilityMsg));
  memcpy(
      &(ue_context->amf_context.originating_tai), (const void*) originating_tai,
      sizeof(tai_t));

  ue_context->amf_context.decode_status = decode_status;

  if (msg->m5gs_reg_type.type_val == AMF_REGISTRATION_TYPE_INITIAL) {
    OAILOG_DEBUG(LOG_NAS_AMF, "New REGISTRATION_REQUEST processing\n");
    // Check integrity and ciphering algorithm bits
    // If all bits are zero it means integrity and ciphering algorithms are not
    // valid, AMF should reject the initial registration. Note : amf_cause is
    // upto network provider for invalid algorithms, here we considering
    // CONDITIONAL_IE_ERROR as amf cause.
    if (ue_context->amf_context.ue_sec_capability.ia == 0 ||
        ue_context->amf_context.ue_sec_capability.ea == 0) {
      amf_cause = AMF_CAUSE_UE_SEC_CAP_MISSMATCH;
      OAILOG_ERROR(
          LOG_NAS_AMF,
          "UE is not supporting any algorithms, AMF rejecting the initial "
          "registration with cause : %d for UE ID: " AMF_UE_NGAP_ID_FMT,
          amf_cause, ue_id);
      rc = amf_proc_registration_reject(ue_id, amf_cause);
      OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    } else {
      // AMF supporting integrity algorinths IA0 to IA2 and ciphering algorithms
      // EA0 to EA2 checking UE supporting algorithms which are supported by AMF
      // or not
      uint8_t supported_ia = ue_context->amf_context.ue_sec_capability.ia0 +
                             ue_context->amf_context.ue_sec_capability.ia1 +
                             ue_context->amf_context.ue_sec_capability.ia2;
      uint8_t supported_ea = ue_context->amf_context.ue_sec_capability.ea0 +
                             ue_context->amf_context.ue_sec_capability.ea1 +
                             ue_context->amf_context.ue_sec_capability.ea2;

      if (supported_ia == 0 || supported_ea == 0) {
        amf_cause = AMF_CAUSE_UE_SEC_CAP_MISSMATCH;
        OAILOG_ERROR(
            LOG_NAS_AMF,
            "UE is not supporting the algorithms IA0,IA1,IA2 and EA0,EA1,EA2, "
            "AMF rejecting the initial registration with cause : %d for UE "
            "ID: " AMF_UE_NGAP_ID_FMT,
            amf_cause, ue_id);
        rc = amf_proc_registration_reject(ue_id, amf_cause);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
      }
    }
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
        amf_copy_plmn_to_supi(
            msg->m5gs_mobile_identity.mobile_identity.imsi, supi_imsi);
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
        OAILOG_DEBUG(
            LOG_AMF_APP,
            "SUPI as IMSI derived : %02x%02x%02x%02x%02x%02x%02x%02x \n",
            params->imsi->u.value[0], params->imsi->u.value[1],
            params->imsi->u.value[2], params->imsi->u.value[3],
            params->imsi->u.value[4], params->imsi->u.value[5],
            params->imsi->u.value[6], params->imsi->u.value[7]);

        amf_copy_plmn_to_context(
            msg->m5gs_mobile_identity.mobile_identity.imsi, ue_context);

        ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_SUCI_IMSI;

        amf_app_generate_guti_on_supi(&amf_guti, &supi_imsi);

        amf_ue_context_on_new_guti(ue_context, (guti_m5_t*) &amf_guti);

        ue_context->amf_context.m5_guti.m_tmsi = amf_guti.m_tmsi;
        ue_context->amf_context.m5_guti.guamfi = amf_guti.guamfi;
        imsi64_t imsi64                = amf_imsi_to_imsi64(params->imsi);
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
            amf_supi_guti_map.insert(std::pair<imsi64_t, guti_and_amf_id_t>(
                imsi64, guti_and_amf_id));
          } else {
            // Overwrite the second element.
            found_imsi->second = guti_and_amf_id;
          }
        }
      } else {
        /*
         * Call subscriberdb to decode the SUPI or IMSI from SUCI as scheme
         * output is encrypted
         */
        delete params;
        amf_copy_plmn_to_context(
            msg->m5gs_mobile_identity.mobile_identity.imsi, ue_context);

        ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_SUCI_IMSI;

        std::string empheral_public_key = reinterpret_cast<char*>(
            msg->m5gs_mobile_identity.mobile_identity.imsi.empheral_public_key);
        std::string ciphertext = reinterpret_cast<char*>(
            msg->m5gs_mobile_identity.mobile_identity.imsi.ciphertext);
        std::string mac_tag = reinterpret_cast<char*>(
            msg->m5gs_mobile_identity.mobile_identity.imsi.mac_tag);

        get_decrypt_imsi_suci_extension(
            &ue_context->amf_context,
            msg->m5gs_mobile_identity.mobile_identity.imsi.home_nw_id,
            empheral_public_key, ciphertext, mac_tag);
        OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
      }
    } else if (
        msg->m5gs_mobile_identity.mobile_identity.guti.type_of_identity ==
        M5GSMobileIdentityMsg_GUTI) {
      params->guti                        = new (guti_m5_t)();
      ue_context->amf_context.reg_id_type = M5GSMobileIdentityMsg_GUTI;
    }
  }  // end of AMF_REGISTRATION_TYPE_INITIAL

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
    OAILOG_DEBUG(
        LOG_NAS_AMF,
        "amf_registration_type_periodic_updating processing"
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
      OAILOG_DEBUG(
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
      OAILOG_ERROR(
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
  rc = amf_proc_registration_request(ue_id, is_amf_ctx_new, params);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);

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

/****************************************************************************
 **                                                                        **
 ** Name:    amf_handle_security_mode_reject()                             **
 **                                                                        **
 ** Description: Processes Security Mode Reject message                    **
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
int amf_handle_security_mode_reject(
    amf_ue_ngap_id_t ue_id, SecurityModeRejectMsg* msg, int amf_cause,
    const amf_nas_message_decode_status_t status) {
  OAILOG_FUNC_IN(LOG_NAS_AMF);
  int rc = RETURNok;

  OAILOG_WARNING(
      LOG_NAS_AMF,
      "AMFAS-SAP - Received Security Mode Reject message "
      "(cause=%d)\n",
      msg->m5gmm_cause.m5gmm_cause);

  /*
   * Message checking
   */
  if (msg->m5gmm_cause.m5gmm_cause == AMF_CAUSE_SUCCESS) {
    amf_cause = AMF_CAUSE_INVALID_MANDATORY_INFO;
  }

  /*
   * Handle message checking error
   */
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNerror);
  }

  if (msg->m5gmm_cause.m5gmm_cause == AMF_CAUSE_UE_SEC_CAP_MISSMATCH) {
    increment_counter(
        "security_mode_reject_received", 1, 1, "cause", "ue_sec_cap_mismatch");
  } else {
    increment_counter(
        "security_mode_reject_received", 1, 1, "cause", "unspecified");
  }

  /*
   * Execute the NAS security mode command not accepted by the UE
   */
  rc = amf_proc_security_mode_reject(ue_id);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

}  // namespace magma5g
