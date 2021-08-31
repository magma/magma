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

  uint32_t len = nas5g_pkt_snap.get_pdu_session_est_type1_len();

  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_est_req, nas5g_pkt_snap.pdu_session_est_req_type1, len);

  EXPECT_EQ(decode_res, true);
}

TEST(test_amf_nas5g_pkt_process, test_amf_pdu_sess_est_req_type2_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;

  uint32_t len = nas5g_pkt_snap.get_pdu_session_est_type2_len();

  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));
  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_est_req, nas5g_pkt_snap.pdu_session_est_req_type2, len);

  EXPECT_EQ(decode_res, true);
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

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  return RUN_ALL_TESTS();
}

}  // namespace magma5g
