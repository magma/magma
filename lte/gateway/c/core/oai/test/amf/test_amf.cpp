/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include "util_nas5g_pkt.h"
#include "include/amf_session_manager_pco.h"
#include <gtest/gtest.h>
#include "intertask_interface.h"
#include "../../tasks/amf/amf_app_ue_context_and_proc.h"

extern "C" {
#include "dynamic_memory_check.h"
#define CHECK_PROTOTYPE_ONLY
#include "intertask_interface_init.h"
#undef CHECK_PROTOTYPE_ONLY
#include "intertask_interface.h"
#include "intertask_interface_types.h"
#include "itti_free_defined_msg.h"
}

const task_info_t tasks_info[] = {
    {THREAD_NULL, "TASK_UNKNOWN", "ipc://IPC_TASK_UNKNOWN"},
#define TASK_DEF(tHREADiD)                                                     \
  {THREAD_##tHREADiD, #tHREADiD, "ipc://IPC_" #tHREADiD},
#include <tasks_def.h>
#undef TASK_DEF
};

task_zmq_ctx_t grpc_service_task_zmq_ctx;

using ::testing::Test;

namespace magma5g {

uint8_t NAS5GPktSnapShot::reg_req_buffer[38] = {
    0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x09, 0xf1, 0x07,
    0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x10,
    0x01, 0x00, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0, 0x2f, 0x05,
    0x04, 0x01, 0x00, 0x00, 0x01, 0x53, 0x01, 0x00};

uint8_t NAS5GPktSnapShot::reg_resync_buffer[20] = {
    0x7e, 0x00, 0x59, 0x15, 0x30, 0x0e, 0xdc, 0xd5, 0xbb, 0x86,
    0xd4, 0xf0, 0xfb, 0xa9, 0xdc, 0x46, 0x8b, 0x8c, 0xdd, 0x67};

uint8_t NAS5GPktSnapShot::guti_based_registration[91] = {
    0x7e, 0x00, 0x41, 0x01, 0x00, 0x0b, 0xf2, 0x22, 0xf2, 0x54, 0x00, 0x00,
    0x00, 0x74, 0x20, 0x32, 0x00, 0x2e, 0x04, 0x80, 0xe0, 0x80, 0xe0, 0x71,
    0x00, 0x41, 0x7e, 0x00, 0x41, 0x01, 0x00, 0x0b, 0xf2, 0x22, 0xf2, 0x54,
    0x00, 0x00, 0x00, 0x74, 0x20, 0x32, 0x00, 0x10, 0x01, 0x03, 0x2e, 0x04,
    0x80, 0xe0, 0x80, 0xe0, 0x2f, 0x02, 0x01, 0x01, 0x52, 0x22, 0x62, 0x54,
    0x00, 0x00, 0x01, 0x17, 0x07, 0x80, 0xe0, 0xe0, 0x60, 0x00, 0x1c, 0x30,
    0x18, 0x01, 0x00, 0x74, 0x00, 0x0a, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65,
    0x72, 0x6e, 0x65, 0x74, 0x53, 0x01, 0x01};

uint8_t NAS5GPktSnapShot::pdu_session_est_req_type1[131] = {
    0x7e, 0x00, 0x67, 0x01, 0x00, 0x6c, 0x2e, 0x05, 0x01, 0xc1, 0x00, 0x00,
    0x91, 0x7b, 0x00, 0x62, 0x80, 0xc2, 0x23, 0x23, 0x01, 0x01, 0x00, 0x23,
    0x10, 0xec, 0xa3, 0x90, 0x00, 0x3e, 0xdb, 0xf9, 0x17, 0xbe, 0xcf, 0xa8,
    0x14, 0x8a, 0xcd, 0xde, 0x56, 0x55, 0x4d, 0x54, 0x53, 0x5f, 0x43, 0x48,
    0x41, 0x50, 0x5f, 0x53, 0x52, 0x56, 0x52, 0xc2, 0x23, 0x15, 0x02, 0x01,
    0x00, 0x15, 0x10, 0xb6, 0xfa, 0xad, 0xc5, 0x6a, 0x43, 0x6b, 0x2f, 0x0f,
    0x9f, 0x82, 0x35, 0x6e, 0x07, 0xd9, 0xd9, 0x80, 0x21, 0x1c, 0x01, 0x00,
    0x00, 0x1c, 0x81, 0x06, 0x00, 0x00, 0x00, 0x00, 0x82, 0x06, 0x00, 0x00,
    0x00, 0x00, 0x83, 0x06, 0x00, 0x00, 0x00, 0x00, 0x84, 0x06, 0x00, 0x00,
    0x00, 0x00, 0x00, 0x1a, 0x01, 0x05, 0x12, 0x05, 0x81, 0x22, 0x01, 0x01,
    0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

uint8_t NAS5GPktSnapShot::pdu_session_est_req_type2[47] = {
    0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff, 0xff,
    0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00, 0x0a, 0x00,
    0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x04, 0x01, 0x00, 0x00, 0x01,
    0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

uint8_t NAS5GPktSnapShot::pdu_session_release_complete[12] = {
    0x7e, 0x00, 0x67, 0x01, 0x00, 0x04, 0x2e, 0x05, 0x01, 0xd4, 0x12, 0x05};

uint8_t NAS5GPktSnapShot::deregistrarion_request[17] = {
    0x7e, 0x00, 0x45, 0x01, 0x00, 0x0b, 0xf2, 0x22, 0xf2,
    0x54, 0x00, 0x00, 0x00, 0x18, 0x5d, 0x2e, 0x00};

uint8_t NAS5GPktSnapShot::service_request[37] = {
    0x7e, 0x00, 0x4c, 0x10, 0x00, 0x07, 0xf4, 0x00, 0x00, 0xe4,
    0x2c, 0x6c, 0x68, 0x71, 0x00, 0x15, 0x7e, 0x00, 0x4c, 0x10,
    0x00, 0x07, 0xf4, 0x00, 0x00, 0xe4, 0x2c, 0x6c, 0x68, 0x40,
    0x02, 0x20, 0x00, 0x50, 0x02, 0x20, 0x00};

uint8_t NAS5GPktSnapShot::registration_reject[4] = {0x00, 0x00, 0x00, 0x00};

uint8_t NAS5GPktSnapShot::security_mode_reject[4] = {0x7e, 0x00, 0x5f, 0x24};

TEST(test_amf_nas5g_pkt_process, test_amf_ue_register_req_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  RegistrationRequestMsg reg_request;
  bool decode_res = false;

  uint32_t len = nas5g_pkt_snap.get_reg_req_buffer_len();

  memset(&reg_request, 0, sizeof(RegistrationRequestMsg));

  decode_res = decode_registration_request_msg(
      &reg_request, nas5g_pkt_snap.reg_req_buffer, len);

  EXPECT_EQ(decode_res, true);

  EXPECT_EQ(
      reg_request.extended_protocol_discriminator.extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);

  //  Type is registration Request
  EXPECT_EQ(reg_request.message_type.msg_type, REG_REQUEST);

  //  Registraiton Type is Initial Registration
  EXPECT_EQ(reg_request.m5gs_reg_type.type_val, 1);

  //  5GS Mobile Identity SUPI FORMAT
  EXPECT_EQ(
      reg_request.m5gs_mobile_identity.mobile_identity.imsi.type_of_identity,
      M5GSMobileIdentityMsg_SUCI_IMSI);

  //  5GS Mobile mms digit2
  EXPECT_EQ(
      reg_request.m5gs_mobile_identity.mobile_identity.imsi.mcc_digit1, 0x09);

  EXPECT_EQ(
      reg_request.m5gs_mobile_identity.mobile_identity.imsi.mcc_digit2, 0x00);

  EXPECT_EQ(
      reg_request.m5gs_mobile_identity.mobile_identity.imsi.mcc_digit3, 0x01);

  EXPECT_EQ(
      reg_request.m5gs_mobile_identity.mobile_identity.imsi.mnc_digit1, 0x07);

  EXPECT_EQ(
      reg_request.m5gs_mobile_identity.mobile_identity.imsi.mcc_digit2, 0x0);
}

TEST(test_amf_nas5g_pkt_process, test_amf_ue_guti_register_req_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  RegistrationRequestMsg reg_request;
  bool decode_res = false;

  uint32_t len = nas5g_pkt_snap.get_guti_based_registration_len();

  memset(&reg_request, 0, sizeof(RegistrationRequestMsg));

  decode_res = decode_registration_request_msg(
      &reg_request, nas5g_pkt_snap.guti_based_registration, len);

  EXPECT_EQ(decode_res, true);
}

TEST(test_amf_nas5g_pkt_process, test_amf_auth_sync_fail_res_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  AuthenticationFailureMsg auth_sync_fail;
  bool decode_res = false;

  uint32_t len = nas5g_pkt_snap.get_reg_resync_buffer_len();

  memset(&auth_sync_fail, 0, sizeof(AuthenticationFailureMsg));

  decode_res = decode_auth_failure_decode_msg(
      &auth_sync_fail, nas5g_pkt_snap.reg_resync_buffer, len);

  bdestroy(auth_sync_fail.auth_failure_ie.authentication_failure_info);
  EXPECT_EQ(decode_res, true);
}

TEST(test_amf_nas5g_pkt_process, test_amf_pdu_sess_est_req_type1_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  protocol_configuration_options_t* pco;

  uint32_t len = nas5g_pkt_snap.get_pdu_session_est_type1_len();

  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_est_req, nas5g_pkt_snap.pdu_session_est_req_type1, len);

  pco = &(pdu_sess_est_req.payload_container.smf_msg.msg
              .pdu_session_estab_request.protocolconfigurationoptions.pco);

  for (uint8_t i = 0; i < pco->num_protocol_or_container_id; i++) {
    if (pco->protocol_or_container_ids[i].contents) {
      bdestroy_wrapper(&pco->protocol_or_container_ids[i].contents);
    }
  }

  EXPECT_EQ(decode_res, true);
}

TEST(test_amf_nas5g_pkt_gen, test_amf_pdu_sess_accept_pco_msg) {
  uint8_t buffer[1024] = {};
  uint16_t buf_len     = 1024;
  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  uint32_t len    = nas5g_pkt_snap.get_pdu_session_est_type1_len();

  /* Initialize primary and secondary dns */
  inet_pton(AF_INET, "192.168.1.100", &(amf_config.ipv4.default_dns));
  inet_pton(AF_INET, "8.8.8.8", &(amf_config.ipv4.default_dns_sec));

  /* Decode the packet */
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));
  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_est_req, nas5g_pkt_snap.pdu_session_est_req_type1, len);
  EXPECT_EQ(decode_res, true);

  protocol_configuration_options_t* pco_req =
      &(pdu_sess_est_req.payload_container.smf_msg.msg.pdu_session_estab_request
            .protocolconfigurationoptions.pco);

  ProtocolConfigurationOptions protocolconfigruartionoption;
  protocol_configuration_options_t* pco_resp =
      &(protocolconfigruartionoption.pco);

  uint8_t ipcp_pattern_match[] = {0x7b, 0x0,  0x14, 0x80, 0x80, 0x21, 0x10, 0x3,
                                  0x0,  0x0,  0x10, 0x81, 0x6,  0xc0, 0xa8, 0x1,
                                  0x64, 0x83, 0x6,  0x8,  0x8,  0x8,  0x8};
  int cmp_res                  = 0;
  int pco_len                  = 0;

  sm_process_pco_request(pco_req, pco_resp);

  pco_len = protocolconfigruartionoption.EncodeProtocolConfigurationOptions(
      &protocolconfigruartionoption,
      REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE, buffer, buf_len);

  EXPECT_EQ(pco_len, 23);

  cmp_res = memcmp(
      buffer, ipcp_pattern_match, sizeof(ipcp_pattern_match) / sizeof(uint8_t));

  EXPECT_EQ(cmp_res, 0);
  sm_free_protocol_configuration_options(&pco_req);
  sm_free_protocol_configuration_options(&pco_resp);
}

TEST(test_amf_nas5g_pkt_process, test_amf_pdu_sess_est_req_type2_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  protocol_configuration_options_t* pco_req;
  uint8_t buffer[1024]        = {};
  uint16_t buf_len            = 1024;
  int cmp_res                 = 0;
  int pco_len                 = 0;
  uint8_t dns_pattern_match[] = {0x7b, 0x0,  0x8,  0x80, 0x0, 0xd,
                                 0x4,  0xc0, 0xa8, 0x1,  0x64};

  /* Encoded Message */
  ProtocolConfigurationOptions protocolconfigruartionoption;
  protocol_configuration_options_t* pco_resp =
      &(protocolconfigruartionoption.pco);

  /* Initialize primary and secondary dns */
  inet_pton(AF_INET, "192.168.1.100", &(amf_config.ipv4.default_dns));

  uint32_t len = nas5g_pkt_snap.get_pdu_session_est_type2_len();

  /* Check if uplink pdu packet is parsed properly */
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));
  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_est_req, nas5g_pkt_snap.pdu_session_est_req_type2, len);

  EXPECT_EQ(decode_res, true);

  pco_req = &(pdu_sess_est_req.payload_container.smf_msg.msg
                  .pdu_session_estab_request.protocolconfigurationoptions.pco);

  /* Check whether the PCO field is decoded properly */
  EXPECT_EQ(pco_req->protocol_or_container_ids[0].id, 10);
  EXPECT_EQ(pco_req->protocol_or_container_ids[1].id, 13);

  sm_process_pco_request(pco_req, pco_resp);

  pco_len = protocolconfigruartionoption.EncodeProtocolConfigurationOptions(
      &protocolconfigruartionoption,
      REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE, buffer, buf_len);

  EXPECT_EQ(pco_len, 11);

  cmp_res = memcmp(
      buffer, dns_pattern_match, sizeof(dns_pattern_match) / sizeof(uint8_t));

  EXPECT_EQ(cmp_res, 0);

  sm_free_protocol_configuration_options(&pco_req);
  sm_free_protocol_configuration_options(&pco_resp);
}

TEST(test_amf_nas5g_pkt_process, test_amf_pdu_sess_release_complete_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_rel_complete_req;
  bool decode_res = false;

  uint32_t len = nas5g_pkt_snap.get_pdu_session_release_complete_len();

  memset(&pdu_sess_rel_complete_req, 0, sizeof(ULNASTransportMsg));
  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_rel_complete_req, nas5g_pkt_snap.pdu_session_release_complete,
      len);

  EXPECT_EQ(decode_res, true);
}

TEST(test_amf_nas5g_pkt_process, test_amf_deregistration_request_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  DeRegistrationRequestUEInitMsg dereg_req;
  bool decode_res = false;

  uint32_t len = nas5g_pkt_snap.get_deregistrarion_request_len();

  memset(&dereg_req, 0, sizeof(DeRegistrationRequestUEInitMsg));
  decode_res = decode_ul_nas_deregister_request_msg(
      &dereg_req, nas5g_pkt_snap.deregistrarion_request, len);

  EXPECT_EQ(decode_res, true);
}

/* Test for service type Data */
TEST(test_amf_nas5g_pkt_process, test_amf_service_request_messagetype_data) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ServiceRequestMsg service_request;
  bool decode_res = 0;

  uint32_t len = nas5g_pkt_snap.get_service_request_len();

  memset(&service_request, 0, sizeof(ServiceRequestMsg));

  decode_res = decode_service_request_msg(
      &service_request, nas5g_pkt_snap.service_request, len);
  EXPECT_EQ(decode_res, true);
  EXPECT_EQ(
      service_request.extended_protocol_discriminator
          .extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);
  EXPECT_EQ(service_request.sec_header_type.sec_hdr, (uint8_t) 0x00);
  EXPECT_EQ(service_request.message_type.msg_type, M5G_SERVICE_REQUEST);
  EXPECT_EQ(service_request.nas_key_set_identifier.nas_key_set_identifier, 1);
  EXPECT_EQ(service_request.service_type.service_type_value, SERVICE_TYPE_DATA);
  EXPECT_EQ(service_request.uplink_data_status.iei, UP_LINK_DATA_STATUS);
  EXPECT_EQ(service_request.uplink_data_status.len, 0x02);
  EXPECT_EQ(service_request.uplink_data_status.uplinkDataStatus, 0x0020);
  EXPECT_EQ(service_request.pdu_session_status.iei, PDU_SESSION_STATUS);
  EXPECT_EQ(service_request.pdu_session_status.len, 0x02);
  EXPECT_EQ(service_request.pdu_session_status.pduSessionStatus, 0x0020);
}

TEST(test_amf_nas5g_pkt_process, test_amf_service_accept_message) {
  ServiceAcceptMsg service_accept;
  uint8_t buffer[50] = {0};

  service_accept.extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;

  service_accept.sec_header_type.sec_hdr = 0;
  service_accept.spare_half_octet.spare  = 0;

  service_accept.message_type.msg_type               = M5G_SERVICE_ACCEPT;
  service_accept.pdu_session_status.iei              = PDU_SESSION_STATUS;
  service_accept.pdu_session_status.len              = 0x02;
  service_accept.pdu_session_status.pduSessionStatus = 0x05;
  service_accept.pdu_session_status.iei = PDU_SESSION_REACTIVATION_RESULT;
  service_accept.pdu_session_status.len = 0x02;
  service_accept.pdu_session_status.pduSessionStatus = 0x05;

  EXPECT_NE(
      service_accept.EncodeServiceAcceptMsg(&service_accept, buffer, 0), 0);
}

TEST(test_amf_nas5g_pkt_process, test_amf_service_accept) {
#define PDU_SESSION_ID 0x0005

  amf_as_establish_t svc_accpt_message = {0};
  amf_nas_message_t nas_msg            = {0};

  svc_accpt_message.pdu_session_status_ie |= AMF_AS_PDU_SESSION_STATUS;
  svc_accpt_message.pdu_session_status = PDU_SESSION_ID;
  svc_accpt_message.pdu_session_status_ie |=
      AMF_AS_PDU_SESSION_REACTIVATION_STATUS;
  svc_accpt_message.pdu_session_reactivation_status = PDU_SESSION_ID;

  int result = amf_service_acceptmsg(&svc_accpt_message, &nas_msg);

  EXPECT_GT(result, 0);
  EXPECT_EQ(
      nas_msg.security_protected.plain.amf.msg.service_accept.pdu_session_status
          .pduSessionStatus,
      PDU_SESSION_ID);
  EXPECT_EQ(
      nas_msg.security_protected.plain.amf.msg.service_accept
          .pdu_re_activation_status.pduSessionReActivationResult,
      PDU_SESSION_ID);
}
TEST(test_amf_data_struct, test_ue_context_creation) {
  ue_m5gmm_context_s* ue_context = nullptr;

  ue_context = amf_create_new_ue_context();
  EXPECT_TRUE(nullptr != ue_context);
  EXPECT_TRUE(0 == ue_context->amf_teid_n11);
  EXPECT_TRUE(0 == ue_context->paging_context.paging_retx_count);
  delete ue_context;
}

TEST(test_smf_context_struct, test_smf_context_creation) {
  ue_m5gmm_context_s* ue_context = nullptr;
  smf_context_t* smf_context     = nullptr;

  ue_context             = amf_create_new_ue_context();
  uint8_t pdu_session_id = 10;
  smf_context            = amf_insert_smf_context(ue_context, pdu_session_id);
  EXPECT_TRUE(0 == smf_context->n_active_pdus);
  EXPECT_TRUE(0 == smf_context->pdu_session_version);
  delete ue_context;
}

/* Test for registration reject */
TEST(test_amf_nas5g_pkt_process, test_amf_registration_reject_msg) {
  uint8_t buffer[4] = {0};
  // registration reject message
  RegistrationRejectMsg reg_rej;
  RegistrationRejectMsg decode_reg_rej;
  reg_rej.extended_protocol_discriminator.extended_proto_discriminator = 0x7e;
  reg_rej.sec_header_type.sec_hdr                                      = 0;
  reg_rej.spare_half_octet.spare                                       = 0;
  reg_rej.message_type.msg_type                                        = 0x44;
  reg_rej.m5gmm_cause.m5gmm_cause                                      = 23;

  bool encode_res = false;
  bool decode_res = false;

  uint32_t len = 4;

  encode_res = encode_registration_reject_msg(&reg_rej, buffer, len);

  decode_res = decode_registration_reject_msg(&decode_reg_rej, buffer, len);

  EXPECT_EQ(encode_res, true);
  EXPECT_EQ(decode_res, true);

  EXPECT_TRUE(
      reg_rej.extended_protocol_discriminator.extended_proto_discriminator ==
      decode_reg_rej.extended_protocol_discriminator
          .extended_proto_discriminator);
  EXPECT_TRUE(
      reg_rej.sec_header_type.sec_hdr ==
      decode_reg_rej.sec_header_type.sec_hdr);
  EXPECT_TRUE(
      reg_rej.spare_half_octet.spare == decode_reg_rej.spare_half_octet.spare);
  EXPECT_TRUE(
      reg_rej.message_type.msg_type == decode_reg_rej.message_type.msg_type);
  EXPECT_TRUE(
      reg_rej.m5gmm_cause.m5gmm_cause ==
      decode_reg_rej.m5gmm_cause.m5gmm_cause);
}

TEST(test_amf_nas5g_pkt_process, test_amf_service_reject_message) {
  ServiceRejectMsg service_reject, decoded_service_rej;
  uint8_t buffer[50] = {0};
  uint8_t len        = 8;

  int encode_res = 0, decode_res = 0;

  service_reject.extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;

  service_reject.sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  service_reject.spare_half_octet.spare  = 0;

  service_reject.message_type.msg_type               = M5G_SERVICE_REJECT;
  service_reject.pdu_session_status.iei              = PDU_SESSION_STATUS;
  service_reject.pdu_session_status.len              = 0x02;
  service_reject.pdu_session_status.pduSessionStatus = 0x05;
  service_reject.cause.iei         = static_cast<uint8_t>(M5GIei::M5GMM_CAUSE);
  service_reject.cause.m5gmm_cause = 9;
  service_reject.t3346Value.iei    = GPRS_TIMER2;
  service_reject.t3346Value.len    = 1;
  service_reject.t3346Value.timervalue = 60;

  encode_res =
      service_reject.EncodeServiceRejectMsg(&service_reject, buffer, len);

  EXPECT_EQ(encode_res, len);

  decode_res = decoded_service_rej.DecodeServiceRejectMsg(
      &decoded_service_rej, buffer, len);

  EXPECT_EQ(decode_res, len);

  EXPECT_EQ(
      service_reject.sec_header_type.sec_hdr,
      decoded_service_rej.sec_header_type.sec_hdr);
  EXPECT_EQ(
      service_reject.spare_half_octet.spare,
      decoded_service_rej.spare_half_octet.spare);
  EXPECT_EQ(
      service_reject.message_type.msg_type,
      decoded_service_rej.message_type.msg_type);
  EXPECT_EQ(
      service_reject.pdu_session_status.iei,
      decoded_service_rej.pdu_session_status.iei);
  EXPECT_EQ(
      service_reject.pdu_session_status.len,
      decoded_service_rej.pdu_session_status.len);
  EXPECT_EQ(
      service_reject.pdu_session_status.pduSessionStatus,
      decoded_service_rej.pdu_session_status.pduSessionStatus);
  EXPECT_EQ(
      service_reject.cause.m5gmm_cause, decoded_service_rej.cause.m5gmm_cause);
}

TEST(test_dlnastransport, test_dlnastransport) {
  DLNASTransportMsg* dlmsg = nullptr;
  SmfMsg* smf_msg          = nullptr;
  uint32_t bytes           = 0;
  uint32_t container_len   = 0;
  bstring buffer;
  amf_nas_message_t msg = {};

  /* build uplinknastransport */
  // uplink nas transport(pdu session request)
  uint8_t pdu[44] = {0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01,
                     0xc1, 0xff, 0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b,
                     0x00, 0x07, 0x80, 0x00, 0x0a, 0x00, 0x00, 0x0d, 0x00,
                     0x12, 0x01, 0x81, 0x22, 0x01, 0x01, 0x25, 0x09, 0x08,
                     0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};
  uint32_t len    = sizeof(pdu) / sizeof(uint8_t);

  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(&pdu_sess_est_req, pdu, len);

  EXPECT_EQ(decode_res, true);
  /* build uplinknastransport */

  ULNASTransportMsg* ulmsg = &pdu_sess_est_req;

  // Message construction for PDU Establishment Reject
  // NAS-5GS (NAS) PDU
  msg.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.plain.amf.header.message_type          = DLNASTRANSPORT;
  msg.header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  // SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.message_type                    = DLNASTRANSPORT;
  msg.header.sequence_number                 = 1;

  dlmsg = &msg.plain.amf.msg.downlinknas5gtransport;

  // AmfHeader
  dlmsg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  len++;
  dlmsg->spare_half_octet.spare  = 0x00;
  dlmsg->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  len++;
  dlmsg->message_type.msg_type = DLNASTRANSPORT;
  len++;
  dlmsg->payload_container.iei = PAYLOAD_CONTAINER;

  // SmfMsg
  dlmsg->payload_container_type.iei      = 0;
  dlmsg->payload_container_type.type_val = N1_SM_INFO;
  len++;
  dlmsg->pdu_session_identity.iei =
      static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2);
  len++;
  dlmsg->pdu_session_identity.pdu_session_id =
      ulmsg->payload_container.smf_msg.header.pdu_session_id;
  len++;

  dlmsg->m5gmm_cause.iei = static_cast<uint8_t>(M5GIei::M5GMM_CAUSE);
  dlmsg->m5gmm_cause.m5gmm_cause =
      static_cast<uint8_t>(M5GMmCause::MAX_PDU_SESSIONS_REACHED);
  len += 2;

  // Payload container IE from ulmsg
  dlmsg->payload_container.copy(ulmsg->payload_container);

  len += 2;  // 2 bytes for container.len
  len += dlmsg->payload_container.len;

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(len);

  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(buffer->data, &msg, len, nullptr);
  EXPECT_GT(bytes, 0);

  amf_nas_message_t decode_msg                  = {0};
  amf_nas_message_decode_status_t decode_status = {};
  int status                                    = RETURNerror;
  status                                        = nas5g_message_decode(
      buffer->data, &decode_msg, bytes, nullptr, &decode_status);

  EXPECT_EQ(true, dlmsg->payload_container.isEqual(ulmsg->payload_container));
  EXPECT_EQ(
      dlmsg->m5gmm_cause.m5gmm_cause,
      static_cast<uint8_t>(M5GMmCause::MAX_PDU_SESSIONS_REACHED));
  bdestroy(buffer);
}

/* Test for security mode reject Data */
TEST(test_amf_nas5g_pkt_process, test_amf_security_mode_reject_message_data) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  SecurityModeRejectMsg sm_reject;
  bool decode_res = 0;

  uint32_t len = nas5g_pkt_snap.get_security_mode_reject_len();

  memset(&sm_reject, 0, sizeof(SecurityModeRejectMsg));

  decode_res = decode_security_mode_reject_msg(
      &sm_reject, nas5g_pkt_snap.security_mode_reject, len);
  EXPECT_EQ(decode_res, true);
  EXPECT_EQ(
      sm_reject.extended_protocol_discriminator.extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);
  EXPECT_EQ(sm_reject.sec_header_type.sec_hdr, (uint8_t) 0x00);
  EXPECT_EQ(sm_reject.message_type.msg_type, SEC_MODE_REJECT);
  EXPECT_EQ(sm_reject.m5gmm_cause.m5gmm_cause, 0x24);
}

TEST(test_optional_pdu, test_pdu_establish) {
  DLNASTransportMsg* dlmsg = nullptr;
  uint32_t bytes           = 0;
  uint32_t container_len   = 0;
  bstring buffer;
  amf_nas_message_t msg = {};

  /* build uplinknastransport */
  // uplink nas transport(pdu session request)
  uint8_t pdu[44] = {0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01,
                     0xc1, 0xff, 0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b,
                     0x00, 0x07, 0x80, 0x00, 0x0a, 0x00, 0x00, 0x0d, 0x00,
                     0x12, 0x01, 0x81, 0x22, 0x01, 0x01, 0x25, 0x09, 0x08,
                     0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};
  uint32_t len    = sizeof(pdu) / sizeof(uint8_t);

  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(&pdu_sess_est_req, pdu, len);

  EXPECT_EQ(decode_res, true);
  /* build uplinknastransport */
  EXPECT_EQ(1, pdu_sess_est_req.nssai.sst);
  std::string dnn("internet");
  EXPECT_EQ(dnn, pdu_sess_est_req.dnn.dnn);

  ULNASTransportMsg* ulmsg = &pdu_sess_est_req;

  // Message construction for PDU Establishment Reject
  // NAS-5GS (NAS) PDU
  msg.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.plain.amf.header.message_type          = DLNASTRANSPORT;
  msg.header.security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.message_type                    = DLNASTRANSPORT;
  msg.header.sequence_number                 = 1;

  dlmsg = &msg.plain.amf.msg.downlinknas5gtransport;

  // AmfHeader
  dlmsg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  dlmsg->spare_half_octet.spare  = 0x00;
  dlmsg->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  dlmsg->message_type.msg_type = DLNASTRANSPORT;
  dlmsg->payload_container.iei = PAYLOAD_CONTAINER;

  // SmfMsg
  dlmsg->payload_container_type.iei      = 0;
  dlmsg->payload_container_type.type_val = N1_SM_INFO;

  dlmsg->pdu_session_identity.iei =
      static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2);
  dlmsg->pdu_session_identity.pdu_session_id =
      ulmsg->payload_container.smf_msg.header.pdu_session_id;

  //----------
  SmfMsg* smf_msg    = &dlmsg->payload_container.smf_msg;

  smf_msg->header.extended_protocol_discriminator =      M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->header.pdu_session_id           = 1;
  smf_msg->header.message_type             = PDU_SESSION_ESTABLISHMENT_ACCEPT;
  smf_msg->header.procedure_transaction_id = 10;
  smf_msg->msg.pdu_session_estab_accept.extended_protocol_discriminator
      .extended_proto_discriminator = M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->msg.pdu_session_estab_accept.pdu_session_identity.pdu_session_id =
      1;
  smf_msg->msg.pdu_session_estab_accept.pti.pti =
      10;
  smf_msg->msg.pdu_session_estab_accept.message_type.msg_type =
      PDU_SESSION_ESTABLISHMENT_ACCEPT;
  smf_msg->msg.pdu_session_estab_accept.pdu_session_type.type_val = 1;
  smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val         = 1;

  memset(
      &(smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info), 0, 12);
  uint8_t address_info[12] = {};
  for (int i = 0; i < PDU_ADDR_IPV4_LEN; i++) {
    smf_msg->msg.pdu_session_estab_accept.pdu_address.address_info[i] =
        address_info[i];
  }
  smf_msg->msg.pdu_session_estab_accept.pdu_address.type_val = PDU_ADDR_TYPE;

  /* QOSrules are hardcoded as it is not exchanged in AMF-SMF
   * gRPC calls as of now, handled in upcoming PR
   * TODO: get the rules for the session from SMF and use it here
   */
  smf_msg->msg.pdu_session_estab_accept.qos_rules.length = 0x9;
  QOSRule qos_rule;
  qos_rule.qos_rule_id         = 0x1;
  qos_rule.len                 = 0x6;
  qos_rule.rule_oper_code      = 0x1;
  qos_rule.dqr_bit             = 0x1;
  qos_rule.no_of_pkt_filters   = 0x1;
  qos_rule.qos_rule_precedence = 0xff;
  qos_rule.spare               = 0x0;
  qos_rule.segregation         = 0x0;
  qos_rule.qfi = 0x5;
  NewQOSRulePktFilter new_qos_rule_pkt_filter;
  new_qos_rule_pkt_filter.spare          = 0x0;
  new_qos_rule_pkt_filter.pkt_filter_dir = 0x3;
  new_qos_rule_pkt_filter.pkt_filter_id  = 0x1;
  new_qos_rule_pkt_filter.len            = 0x1;
  uint8_t contents                       = 0x1;
  memcpy(
      new_qos_rule_pkt_filter.contents, &contents, new_qos_rule_pkt_filter.len);
  memcpy(
      qos_rule.new_qos_rule_pkt_filter, &new_qos_rule_pkt_filter,
      1 * sizeof(NewQOSRulePktFilter));
  memcpy(
      smf_msg->msg.pdu_session_estab_accept.qos_rules.qos_rule, &qos_rule,
      1 * sizeof(QOSRule));
  smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_unit =  64;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_unit =  64;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.dl_session_ambr = 1;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.ul_session_ambr = 1;
  smf_msg->msg.pdu_session_estab_accept.session_ambr.length = AMBR_LEN;
  
  //----------
 
  smf_msg->msg.pdu_session_estab_accept.nssai.iei = static_cast<uint8_t>(M5GIei::S_NSSA);
	  smf_msg->msg.pdu_session_estab_accept.nssai.len = 1;
  smf_msg->msg.pdu_session_estab_accept.nssai.sst = 1;
  //3
  smf_msg->msg.pdu_session_estab_accept.dnn.iei = static_cast<uint8_t>(M5GIei::DNN);
	  smf_msg->msg.pdu_session_estab_accept.dnn.len = 9;
  smf_msg->msg.pdu_session_estab_accept.dnn.dnn = "internet";

  dlmsg->payload_container.len =
      PDU_ESTAB_ACCPET_PAYLOAD_CONTAINER_LEN + 14 ;
  len = PDU_ESTAB_ACCEPT_NAS_PDU_LEN + 14;

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(len);

  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(buffer->data, &msg, len, nullptr);
  EXPECT_GT(bytes, 0);

#if 0
  amf_nas_message_t decode_msg                  = {0};
  amf_nas_message_decode_status_t decode_status = {};
  int status                                    = RETURNerror;
  status                                        = nas5g_message_decode(
      buffer->data, &decode_msg, bytes, nullptr, &decode_status);

  EXPECT_EQ(true, dlmsg->payload_container.isEqual(ulmsg->payload_container));
  EXPECT_EQ(
      dlmsg->m5gmm_cause.m5gmm_cause,
      static_cast<uint8_t>(M5GMmCause::MAX_PDU_SESSIONS_REACHED));
#endif
  bdestroy(buffer);
}
int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma5g
