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
#include <chrono>
#include <thread>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.h"

extern "C" {
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#define CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_init.h"
#undef CHECK_PROTOTYPE_ONLY
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/common/itti_free_defined_msg.h"
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
}

#include "lte/gateway/c/core/oai/test/amf/util_nas5g_pkt.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_session_manager_pco.h"
#include <gtest/gtest.h>
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/include/mme_config.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_authentication.h"
#include "lte/gateway/c/core/oai/test/amf/util_s6a_update_location.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_as.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.h"
#include "lte/gateway/c/core/oai/test/amf/amf_app_test_util.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_packet_handler.h"

using ::testing::Test;
task_zmq_ctx_t grpc_service_task_zmq_ctx;

namespace magma5g {
extern task_zmq_ctx_s amf_app_task_zmq_ctx;
extern std::unordered_map<imsi64_t, guti_and_amf_id_t> amf_supi_guti_map;
extern std::unordered_map<amf_ue_ngap_id_t, ue_m5gmm_context_s*> ue_context_map;

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

uint8_t NAS5GPktSnapShot::pdu_session_est_req_type3[34] = {
    0x7e, 0x00, 0x67, 0x01, 0x00, 0x0e, 0x2e, 0x05, 0x01, 0xc1, 0xff, 0xff,
    0x91, 0xa4, 0x28, 0x01, 0x01, 0x55, 0x02, 0x20, 0x12, 0x05, 0x81, 0x25,
    0x09, 0x08, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x45, 0x54};

uint8_t NAS5GPktSnapShot::pdu_session_release_complete[12] = {
    0x7e, 0x00, 0x67, 0x01, 0x00, 0x04, 0x2e, 0x05, 0x01, 0xd4, 0x12, 0x05};

uint8_t NAS5GPktSnapShot::deregistrarion_request[17] = {
    0x7e, 0x00, 0x45, 0x01, 0x00, 0x0b, 0xf2, 0x22, 0xf2,
    0x54, 0x00, 0x00, 0x00, 0x18, 0x5d, 0x2e, 0x00};

// service request with service type data
uint8_t NAS5GPktSnapShot::service_request[37] = {
    0x7e, 0x00, 0x4c, 0x10, 0x00, 0x07, 0xf4, 0x00, 0x00, 0xe4,
    0x2c, 0x6c, 0x68, 0x71, 0x00, 0x15, 0x7e, 0x00, 0x4c, 0x10,
    0x00, 0x07, 0xf4, 0x00, 0x00, 0xe4, 0x2c, 0x6c, 0x68, 0x40,
    0x02, 0x20, 0x00, 0x50, 0x02, 0x20, 0x00};

// service request with service type signaling
uint8_t NAS5GPktSnapShot::service_req_signaling[13] = {
    0x7e, 0x00, 0x4c, 0x00, 0x00, 0x07, 0xf4,
    0x00, 0x40, 0x21, 0x2e, 0x50, 0x25};

// service request with service type data and without IE uplink
// data status
uint8_t service_request_without_uplink_status[17] = {
    0x7e, 0x00, 0x4c, 0x1b, 0x00, 0x07, 0xf4, 0x01, 0x00,
    0x17, 0xd7, 0xb7, 0x33, 0x50, 0x02, 0x20, 0x00};

uint8_t NAS5GPktSnapShot::registration_reject[4] = {0x00, 0x00, 0x00, 0x00};

uint8_t NAS5GPktSnapShot::security_mode_reject[4] = {0x7e, 0x00, 0x5f, 0x24};

class AmfNas5GTest : public ::testing::Test {
 protected:
  NAS5GPktSnapShot nas5g_pkt_snap;
  RegistrationRequestMsg reg_request = {};
  bool decode_res;
  virtual void SetUp() { decode_res = false; }
  virtual void TearDown() {}
};

TEST_F(AmfNas5GTest, test_amf_ue_register_req_msg) {
  uint32_t len = nas5g_pkt_snap.get_reg_req_buffer_len();

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

TEST_F(AmfNas5GTest, test_amf_ue_guti_register_req_msg) {
  uint32_t len = nas5g_pkt_snap.get_guti_based_registration_len();

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
  memset(
      &protocolconfigruartionoption, 0, sizeof(ProtocolConfigurationOptions));
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

TEST(test_amf_nas5g_pkt_process, test_amf_pdu_sess_est_req_type3_msg) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ULNASTransportMsg pdu_sess_est_req;
  PDUSessionEstablishmentRequestMsg* pduSessEstReq = nullptr;
  bool decode_res                                  = false;
  uint8_t buffer[1024]                             = {};
  uint16_t buf_len                                 = 1024;

  uint32_t len = nas5g_pkt_snap.get_pdu_session_est_type3_len();

  /* Check if uplink pdu packet is parsed properly */
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));
  decode_res = decode_ul_nas_transport_msg(
      &pdu_sess_est_req, nas5g_pkt_snap.pdu_session_est_req_type3, len);

  EXPECT_EQ(decode_res, true);
  pduSessEstReq =
      &pdu_sess_est_req.payload_container.smf_msg.msg.pdu_session_estab_request;
  EXPECT_EQ(
      pduSessEstReq->extended_protocol_discriminator
          .extended_proto_discriminator,
      M5G_SESSION_MANAGEMENT_MESSAGES);
  EXPECT_EQ(pduSessEstReq->pdu_session_identity.pdu_session_id, 0x05);
  EXPECT_EQ(pduSessEstReq->pti.pti, 0x01);
  EXPECT_EQ(
      pduSessEstReq->message_type.msg_type, PDU_SESSION_ESTABLISHMENT_REQUEST);
  EXPECT_EQ(pduSessEstReq->integrity_prot_max_data_rate.max_uplink, 0xFF);
  EXPECT_EQ(pduSessEstReq->integrity_prot_max_data_rate.max_downlink, 0xFF);
  EXPECT_EQ(pduSessEstReq->pdu_session_type.type_val, 0x01);
  EXPECT_EQ(pduSessEstReq->ssc_mode.mode_val, 0x04);
  EXPECT_EQ(
      pduSessEstReq->maxNumOfSuppPacketFilters.iei,
      MAXIMUM_NUMBER_OF_SUPPORTED_PACKET_FILTERS_TYPE);
  EXPECT_EQ(
      pduSessEstReq->maxNumOfSuppPacketFilters.maxNumOfSuppPktFilters, 0x0220);
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

class AmfUeContextTest : public ::testing::Test {
 protected:
  ue_m5gmm_context_s* ue_context;

  virtual void SetUp() { ue_context = amf_create_new_ue_context(); }
  virtual void TearDown() { delete ue_context; }
};

TEST_F(AmfUeContextTest, test_ue_context_creation) {
  EXPECT_TRUE(nullptr != ue_context);
  EXPECT_TRUE(0 == ue_context->amf_teid_n11);
  EXPECT_TRUE(0 == ue_context->paging_context.paging_retx_count);
}

TEST_F(AmfUeContextTest, test_smf_context_creation) {
  std::shared_ptr<smf_context_t> smf_context;
  uint8_t pdu_session_id = 10;
  smf_context            = amf_insert_smf_context(ue_context, pdu_session_id);
  EXPECT_TRUE(0 == smf_context->n_active_pdus);
  EXPECT_TRUE(0 == smf_context->pdu_session_version);
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

TEST(test_dl_msg, test_amf_pdu_session_establish_reject_message_data) {
  uint8_t sequence_number  = 0;
  bool is_security_enabled = false;
  amf_nas_message_t msg    = {};
  uint8_t cause            = 27;
  uint8_t pti              = 1;
  uint8_t session_id       = 1;
  bstring buffer;
  uint32_t bytes = 0;

  int len = construct_pdu_session_reject_dl_req(
      sequence_number, session_id, pti, cause, is_security_enabled, &msg);
  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(buffer->data, &msg, len, nullptr);
  EXPECT_GT(bytes, 0);

  amf_nas_message_t decode_msg                  = {0};
  amf_nas_message_decode_status_t decode_status = {};
  int status                                    = RETURNerror;
  status                                        = nas5g_message_decode(
      buffer->data, &decode_msg, bytes, nullptr, &decode_status);

  DLNASTransportMsg* dlmsg = &decode_msg.plain.amf.msg.downlinknas5gtransport;

  EXPECT_EQ(dlmsg->pdu_session_identity.pdu_session_id, session_id);
  SmfMsg& pdu_sess_est_reject = dlmsg->payload_container.smf_msg;

  EXPECT_EQ(pdu_sess_est_reject.header.pdu_session_id, session_id);
  EXPECT_EQ(pdu_sess_est_reject.header.procedure_transaction_id, pti);
  EXPECT_EQ(pdu_sess_est_reject.msg.pdu_session_estab_reject.pti.pti, pti);
  EXPECT_EQ(
      pdu_sess_est_reject.msg.pdu_session_estab_reject.m5gsm_cause.cause_value,
      cause);

  bdestroy_wrapper(&buffer);
}
/* Test for delete_wrapper */
TEST(test_delete_wrapper, test_delete_wrapper) {
  amf_registration_request_ies_t* req_ies =
      new (amf_registration_request_ies_t)();
  uint32_t* generic_type = new uint32_t;

  delete_wrapper(&req_ies);
  EXPECT_EQ(req_ies, nullptr);

  delete_wrapper(&generic_type);
  EXPECT_EQ(generic_type, nullptr);
}
/*Test for delete specific procedure : Registration Procedure*/
TEST(test_delete_registration_proc, test_delete_registration_proc) {
  ue_m5gmm_context_s* ue_ctxt = amf_create_new_ue_context();
  EXPECT_TRUE(ue_ctxt != nullptr);

  // Specific procedure: Registration Procedure
  nas_amf_registration_proc_t* reg_proc =
      nas_new_registration_procedure(ue_ctxt);
  EXPECT_TRUE(reg_proc != NULL);

  // Child procedures: Authentication Procedure
  nas5g_amf_auth_proc_t* auth_proc =
      nas5g_new_authentication_procedure(&ue_ctxt->amf_context);
  EXPECT_TRUE(auth_proc != NULL);
  (reinterpret_cast<nas5g_base_proc_t*>(auth_proc))->parent =
      reinterpret_cast<nas5g_base_proc_t*>(reg_proc);

  // Child procedures: Identity Procedure
  nas_amf_ident_proc_t* ident_proc =
      nas5g_new_identification_procedure(&ue_ctxt->amf_context);
  EXPECT_TRUE(ident_proc != NULL);
  (reinterpret_cast<nas5g_base_proc_t*>(ident_proc))->parent =
      reinterpret_cast<nas5g_base_proc_t*>(reg_proc);

  // Child procedures: SMC Procedure
  nas_amf_smc_proc_t* smc_proc = nas5g_new_smc_procedure(&ue_ctxt->amf_context);
  EXPECT_TRUE(smc_proc != NULL);
  (reinterpret_cast<nas5g_base_proc_t*>(smc_proc))->parent =
      reinterpret_cast<nas5g_base_proc_t*>(reg_proc);

  amf_delete_registration_proc(&ue_ctxt->amf_context);
  EXPECT_EQ(
      get_nas_specific_procedure_registration(&ue_ctxt->amf_context), nullptr);

  delete_wrapper(&ue_ctxt->amf_context.amf_procedures);
  delete ue_ctxt;
}

TEST(test_optional_dnn_pdu, test_pdu_session_establish_optional) {
  uint32_t bytes         = 0;
  uint32_t container_len = 0;
  bstring buffer;
  amf_nas_message_t msg = {};

  // build uplinknastransport
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
  // SSC mode check
  EXPECT_EQ(
      pdu_sess_est_req.payload_container.smf_msg.msg.pdu_session_estab_request
          .ssc_mode.mode_val,
      1);
  EXPECT_EQ(pdu_sess_est_req.nssai.sst, 1);
  uint8_t dnn[9] = {0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};
  EXPECT_EQ(memcmp(pdu_sess_est_req.dnn.dnn, dnn, pdu_sess_est_req.dnn.len), 0);

  buffer = bfromcstralloc(len, "\0");
  bytes  = pdu_sess_est_req.EncodeULNASTransportMsg(
      &pdu_sess_est_req, buffer->data, len);
  EXPECT_GT(bytes, 0);
  ULNASTransportMsg decode_pdu_sess_est_req = {};
  decode_res = decode_ul_nas_transport_msg(&decode_pdu_sess_est_req, pdu, len);
  EXPECT_EQ(decode_res, true);
  // SSC mode Check
  EXPECT_EQ(
      decode_pdu_sess_est_req.payload_container.smf_msg.msg
          .pdu_session_estab_request.ssc_mode.mode_val,
      1);
  EXPECT_EQ(decode_pdu_sess_est_req.nssai.sst, 1);
  EXPECT_EQ(memcmp(pdu_sess_est_req.dnn.dnn, dnn, pdu_sess_est_req.dnn.len), 0);

  bdestroy(buffer);
}

TEST(test_dnn, test_amf_handle_s6a_update_location_ans) {
  // creating ue_context
  ue_m5gmm_context_s* ue_context = amf_create_new_ue_context();

  // Building s6a_update_location_ans_t
  std::string imsi = "901700000000001";
  s6a_update_location_ans_t ula_ans;
  ula_ans = amf_send_s6a_ula(imsi);

  // Building key value pair for amf_supi_guti_map and ue_context_map
  uint64_t imsi_64 = 901700000000001;
  guti_and_amf_id_t guti_amf;
  guti_amf.amf_guti.m_tmsi = 0x2bfb815f;
  guti_amf.amf_ue_ngap_id  = 0x01;

  // Inserting into amf_supi_guti_map, ue_context_map
  amf_supi_guti_map.insert(
      std::pair<imsi64_t, guti_and_amf_id_t>(imsi_64, guti_amf));
  ue_context_map.insert(std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(
      guti_amf.amf_ue_ngap_id, ue_context));

  int rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_TRUE(rc == RETURNok);

  // Clearing the map and deleting ue_context
  amf_supi_guti_map.clear();
  ue_context_map.clear();
  delete ue_context;
}

TEST(test_dnn, test_amf_validate_dnn) {
  // uplink nas transport(pdu session request)
  uint8_t pdu[44] = {0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01,
                     0xc1, 0xff, 0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b,
                     0x00, 0x07, 0x80, 0x00, 0x0a, 0x00, 0x00, 0x0d, 0x00,
                     0x12, 0x01, 0x81, 0x22, 0x01, 0x01, 0x25, 0x09, 0x08,
                     0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};
  uint32_t len    = sizeof(pdu) / sizeof(uint8_t);

  ULNASTransportMsg msg;
  bool decode_res = false;
  memset(&msg, 0, sizeof(ULNASTransportMsg));
  std::string dnn_string(reinterpret_cast<char*>(msg.dnn.dnn), msg.dnn.len);
  int idx          = 0;
  bool ue_sent_dnn = true;
  // decoding uplink uplink nas transport(pdu session request)
  decode_res = decode_ul_nas_transport_msg(&msg, pdu, len);
  EXPECT_EQ(decode_res, true);

  amf_context_s amf_ctx = {};
  std::string imsi      = "901700000000001";
  s6a_update_location_ans_t ula_ans;

  // mock handling ans received from s6a_update_location_request
  ula_ans = amf_send_s6a_ula(imsi);
  memcpy(
      &amf_ctx.apn_config_profile,
      &ula_ans.subscription_data.apn_config_profile,
      sizeof(apn_config_profile_t));

  // validating dnn against s6a update location ans
  int rc = amf_validate_dnn(&amf_ctx, dnn_string, &idx, ue_sent_dnn);
  EXPECT_TRUE(rc == RETURNok);
}

class AmfUeContextTestServiceRequestProc : public ::testing::Test {
 protected:
#define MCC_DIGIT1 2
#define MCC_DIGIT2 2
#define MCC_DIGIT3 2
#define MNC_DIGIT1 4
#define MNC_DIGIT2 5
#define MNC_DIGIT3 6
#define IMSI64 222456000000101
#define M_TMSI 0X212e5025
#define AMF_SET_ID 1
#define AMF_POINTER 0
#define AMF_REGION_ID 1
#define AMF_TAC 0x03

  ue_m5gmm_context_s* ue_context;
  amf_app_desc_t* amf_app_desc_p;
  guti_m5_t guti;
  tai_t tai;
  const amf_ue_ngap_id_t AMF_UE_NGAP_ID = 0x05;
  const gnb_ue_ngap_id_t gNB_UE_NGAP_ID = 0x09;
  const uint32_t gnb_id                 = 0x01;

  virtual void SetUp() {
    itti_init(
        TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info, NULL,
        NULL);
    amf_config_init(&amf_config);
    amf_nas_state_init(&amf_config);

    ue_context     = amf_create_new_ue_context();
    amf_app_desc_p = get_amf_nas_state(false);

    // insert ue context
    if (ue_context) {
      amf_insert_ue_context(AMF_UE_NGAP_ID, ue_context);
    }

    // imsi64
    ue_context->amf_context.imsi64 = IMSI64;
    // ue security context
    ue_context->amf_context.member_present_mask |= AMF_CTXT_MEMBER_SECURITY;
    // ueContextReq
    ue_context->ue_context_request = M5G_UEContextRequest_requested;
    // ue state
    ue_context->mm_state = REGISTERED_IDLE;
    // 5G TMSI
    ue_context->amf_context.m5_guti.m_tmsi = M_TMSI;
    guti.m_tmsi = ue_context->amf_context.m5_guti.m_tmsi;

    amf_config_read_lock(&amf_config);
    // AMF GUAMI CONFIG
    amf_config.guamfi.guamfi[0].plmn.mcc_digit1 = MCC_DIGIT1;
    amf_config.guamfi.guamfi[0].plmn.mcc_digit2 = MCC_DIGIT2;
    amf_config.guamfi.guamfi[0].plmn.mcc_digit3 = MCC_DIGIT3;
    amf_config.guamfi.guamfi[0].plmn.mnc_digit1 = MNC_DIGIT1;
    amf_config.guamfi.guamfi[0].plmn.mnc_digit2 = MNC_DIGIT2;
    amf_config.guamfi.guamfi[0].plmn.mnc_digit3 = MNC_DIGIT3;
    amf_config.guamfi.guamfi[0].amf_set_id      = AMF_SET_ID;
    amf_config.guamfi.guamfi[0].amf_pointer     = AMF_POINTER;
    amf_config.guamfi.guamfi[0].amf_regionid    = AMF_REGION_ID;
    memcpy(
        &ue_context->amf_context.m5_guti.guamfi, &amf_config.guamfi.guamfi[0],
        sizeof(guamfi_t));
    memcpy(&guti.guamfi, &amf_config.guamfi.guamfi[0], sizeof(guamfi_t));
    amf_config_unlock(&amf_config);
    ue_context->amf_ue_ngap_id = AMF_UE_NGAP_ID;
    // insert ue context based on new guti
    amf_ue_context_on_new_guti(ue_context, &guti);
    // tai
    tai.plmn = guti.guamfi.plmn;
    tai.tac  = AMF_TAC;
  }
  virtual void TearDown() {
    delete ue_context;
    clear_amf_nas_state();
    itti_free_desc_threads();
    amf_config_free(&amf_config);
  }
};

TEST_F(AmfUeContextTestServiceRequestProc, test_amf_service_accept_message) {
  ServiceAcceptMsg service_accept;
  uint8_t buffer[50] = {0};
  amf_sap_t amf_sap;
  amf_as_message_t as_msg;

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

  // Verify nas encoding is successful
  EXPECT_NE(
      service_accept.EncodeServiceAcceptMsg(&service_accept, buffer, 0), 0);

  amf_sap.primitive                     = AMFAS_ESTABLISH_CNF;
  amf_sap.u.amf_as.u.establish.ue_id    = AMF_UE_NGAP_ID;
  amf_sap.u.amf_as.u.establish.nas_info = AMF_AS_NAS_INFO_SR;

  // Verify nas encoding is successful
  EXPECT_EQ(
      AS_NAS_ESTABLISH_CNF_,
      amf_as_establish_cnf(
          &amf_sap.u.amf_as.u.establish, &as_msg.msg.nas_establish_rsp));
}

/* Test for service type signaling */
TEST_F(
    AmfUeContextTestServiceRequestProc,
    test_amf_service_type_signaling_sunny_day) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ServiceRequestMsg service_request;
  bool decode_res                               = 0;
  amf_nas_message_decode_status_t decode_status = {0};
  MessageDef* message_p                         = NULL;
  amf_app_desc_t* amf_app_desc_p                = get_amf_nas_state(false);

  uint32_t len = nas5g_pkt_snap.get_service_request_signaling_len();

  memset(&service_request, 0, sizeof(ServiceRequestMsg));

  decode_res = decode_service_request_msg(
      &service_request, nas5g_pkt_snap.service_req_signaling, len);
  // Verify service request is decoded
  EXPECT_EQ(decode_res, true);
  // veriy service request NAS IE's
  EXPECT_EQ(
      service_request.extended_protocol_discriminator
          .extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);
  EXPECT_EQ(service_request.sec_header_type.sec_hdr, (uint8_t) 0x00);
  EXPECT_EQ(service_request.message_type.msg_type, M5G_SERVICE_REQUEST);
  EXPECT_EQ(service_request.nas_key_set_identifier.nas_key_set_identifier, 0);
  EXPECT_EQ(
      service_request.service_type.service_type_value, SERVICE_TYPE_SIGNALING);

  // Verify service request is handled
  EXPECT_EQ(
      RETURNok, amf_handle_service_request(
                    AMF_UE_NGAP_ID, &service_request, decode_status));
  // Verify UE moved to REGISTERED
  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);
  // Forcing UE state IDLE to handle initial ue message
  ue_context->mm_state = REGISTERED_IDLE;

  // Allocate initial UE message
  message_p = itti_alloc_new_message(TASK_NGAP, NGAP_INITIAL_UE_MESSAGE);
  NGAP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 1;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_id         = gnb_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).nas =
      blk2bstr(nas5g_pkt_snap.service_req_signaling, len);
  NGAP_INITIAL_UE_MESSAGE(message_p).m5g_rrc_establishment_cause =
      M5G_MO_SIGNALLING;
  NGAP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid        = true;
  NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.amf_set_id  = 1;
  NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.amf_pointer = 0;
  NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.m_tmsi      = guti.m_tmsi;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id         = gNB_UE_NGAP_ID;
  NGAP_INITIAL_UE_MESSAGE(message_p).tai                    = tai;
  NGAP_INITIAL_UE_MESSAGE(message_p).ue_context_request =
      M5G_UEContextRequest_requested;

  // verify initial ue message is handled
  EXPECT_EQ(
      ue_context->amf_context.imsi64,
      amf_app_handle_initial_ue_message(
          amf_app_desc_p, &NGAP_INITIAL_UE_MESSAGE(message_p)));

  // Verify UE moved to REGISTERED
  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);
  itti_free_msg_content(message_p);
  free(message_p);
}

/* Test for service type signaling */
TEST_F(
    AmfUeContextTestServiceRequestProc,
    test_amf_service_type_signaling_rainy_day) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ServiceRequestMsg service_request;
  bool decode_res                               = 0;
  amf_nas_message_decode_status_t decode_status = {0};

  uint32_t len = nas5g_pkt_snap.get_service_request_signaling_len();

  memset(&service_request, 0, sizeof(ServiceRequestMsg));

  decode_res = decode_service_request_msg(
      &service_request, nas5g_pkt_snap.service_req_signaling, len);
  // Verify service request is decoded
  EXPECT_EQ(decode_res, true);
  // veriy service request NAS IE's
  EXPECT_EQ(
      service_request.extended_protocol_discriminator
          .extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);
  EXPECT_EQ(service_request.sec_header_type.sec_hdr, (uint8_t) 0x00);
  EXPECT_EQ(service_request.message_type.msg_type, M5G_SERVICE_REQUEST);
  EXPECT_EQ(service_request.nas_key_set_identifier.nas_key_set_identifier, 0);
  EXPECT_EQ(
      service_request.service_type.service_type_value, SERVICE_TYPE_SIGNALING);

  ue_context->amf_context.m5_guti.m_tmsi = 0X25502e22;
  // Verify service request is not handled as TMSI not matching
  EXPECT_EQ(
      RETURNok, amf_handle_service_request(
                    AMF_UE_NGAP_ID, &service_request, decode_status));
  // Verify UE still remains in IDLE state
  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);
}

/* Test for Initial Ue message in connected mode */
TEST_F(
    AmfUeContextTestServiceRequestProc,
    test_amf_initial_ue_message_connected_mode_sunny_day) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ServiceRequestMsg service_request;
  memset(&service_request, 0, sizeof(service_request));
  bool decode_res                               = 0;
  amf_nas_message_decode_status_t decode_status = {0};
  MessageDef* message_p                         = NULL;
  amf_app_desc_t* amf_app_desc_p                = get_amf_nas_state(false);
  gnb_ngap_id_key_t gnb_ngap_id_key             = INVALID_GNB_UE_NGAP_ID_KEY;

  uint32_t len = nas5g_pkt_snap.get_service_request_signaling_len();

  decode_res = decode_service_request_msg(
      &service_request, nas5g_pkt_snap.service_req_signaling, len);
  // Verify service request is decoded
  EXPECT_EQ(decode_res, true);
  ue_context->mm_state = REGISTERED_CONNECTED;

  // Allocate initial UE message
  message_p = itti_alloc_new_message(TASK_NGAP, NGAP_INITIAL_UE_MESSAGE);
  NGAP_INITIAL_UE_MESSAGE(message_p).sctp_assoc_id  = 1;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_id         = gnb_id;
  NGAP_INITIAL_UE_MESSAGE(message_p).nas =
      blk2bstr(nas5g_pkt_snap.service_req_signaling, len);
  NGAP_INITIAL_UE_MESSAGE(message_p).m5g_rrc_establishment_cause =
      M5G_MO_SIGNALLING;
  NGAP_INITIAL_UE_MESSAGE(message_p).is_s_tmsi_valid        = true;
  NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.amf_set_id  = 1;
  NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.amf_pointer = 0;
  NGAP_INITIAL_UE_MESSAGE(message_p).opt_s_tmsi.m_tmsi      = guti.m_tmsi;
  NGAP_INITIAL_UE_MESSAGE(message_p).tai                    = tai;
  NGAP_INITIAL_UE_MESSAGE(message_p).ue_context_request =
      M5G_UEContextRequest_requested;
  AMF_APP_GNB_NGAP_ID_KEY(
      gnb_ngap_id_key, NGAP_INITIAL_UE_MESSAGE(message_p).gnb_id,
      NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id);
  ue_context->gnb_ngap_id_key = gnb_ngap_id_key;
  // change gnb_ud_ngap_id and gnb_id to generate new gnb_ngap_key
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_ue_ngap_id = gNB_UE_NGAP_ID + 1;
  NGAP_INITIAL_UE_MESSAGE(message_p).gnb_id         = gnb_id + 1;

  // verify initial ue message is handled
  EXPECT_EQ(
      ue_context->amf_context.imsi64,
      amf_app_handle_initial_ue_message(
          amf_app_desc_p, &NGAP_INITIAL_UE_MESSAGE(message_p)));

  // Verify new gnb_ngap_id_key got generated
  EXPECT_NE(gnb_ngap_id_key, ue_context->gnb_ngap_id_key);

  // Verify UE still in CONNECTED MODE though initial ue message is received
  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);
  itti_free_msg_content(message_p);
  free(message_p);
}

/* Test for service request without NGAP IE ueContextRequest */
TEST_F(
    AmfUeContextTestServiceRequestProc,
    test_amf_without_ueContextRequest_sunny_day) {
  NAS5GPktSnapShot nas5g_pkt_snap;
  ServiceRequestMsg service_request;
  bool decode_res                               = 0;
  amf_nas_message_decode_status_t decode_status = {0};
  MessageDef* message_p                         = NULL;
  amf_app_desc_t* amf_app_desc_p                = get_amf_nas_state(false);

  uint32_t len = nas5g_pkt_snap.get_service_request_signaling_len();

  memset(&service_request, 0, sizeof(ServiceRequestMsg));

  decode_res = decode_service_request_msg(
      &service_request, nas5g_pkt_snap.service_req_signaling, len);
  // Verify service request is decoded
  EXPECT_EQ(decode_res, true);
  // veriy service request NAS IE's
  EXPECT_EQ(
      service_request.extended_protocol_discriminator
          .extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);
  EXPECT_EQ(service_request.sec_header_type.sec_hdr, (uint8_t) 0x00);
  EXPECT_EQ(service_request.message_type.msg_type, M5G_SERVICE_REQUEST);
  EXPECT_EQ(service_request.nas_key_set_identifier.nas_key_set_identifier, 0);
  EXPECT_EQ(
      service_request.service_type.service_type_value, SERVICE_TYPE_SIGNALING);

  // making ue_context_request IE NULL
  ue_context->ue_context_request = (m5g_uecontextrequest_t) 0;
  // Verify service request is handled
  EXPECT_EQ(
      RETURNok, amf_handle_service_request(
                    AMF_UE_NGAP_ID, &service_request, decode_status));
  // Verify UE moved to REGISTERED
  EXPECT_EQ(REGISTERED_CONNECTED, ue_context->mm_state);
}

/* service request without IE UplinkDataStatus */
TEST_F(
    AmfUeContextTestServiceRequestProc,
    test_amf_service_request_without_uplinkDataStatus_RainyDay) {
  ServiceRequestMsg service_request;
  memset(&service_request, 0, sizeof(service_request));
  bool decode_res                               = 0;
  amf_nas_message_decode_status_t decode_status = {0};
  MessageDef* message_p                         = NULL;
  amf_app_desc_t* amf_app_desc_p                = get_amf_nas_state(false);
  gnb_ngap_id_key_t gnb_ngap_id_key             = INVALID_GNB_UE_NGAP_ID_KEY;

  uint32_t len =
      sizeof(service_request_without_uplink_status) / sizeof(uint8_t);

  decode_res = decode_service_request_msg(
      &service_request, service_request_without_uplink_status, len);
  // Verify service request is decoded
  EXPECT_EQ(decode_res, true);
  EXPECT_EQ(
      service_request.extended_protocol_discriminator
          .extended_proto_discriminator,
      M5G_MOBILITY_MANAGEMENT_MESSAGES);
  EXPECT_EQ(service_request.sec_header_type.sec_hdr, (uint8_t) 0x00);
  EXPECT_EQ(service_request.message_type.msg_type, M5G_SERVICE_REQUEST);
  EXPECT_EQ(service_request.nas_key_set_identifier.nas_key_set_identifier, 1);
  EXPECT_EQ(service_request.service_type.service_type_value, SERVICE_TYPE_DATA);
  // Verify UP_LINK_DATA_STATUS is not present
  EXPECT_NE(service_request.uplink_data_status.iei, UP_LINK_DATA_STATUS);
  EXPECT_NE(service_request.uplink_data_status.len, 0x02);
  EXPECT_NE(service_request.uplink_data_status.uplinkDataStatus, 0x0020);
  EXPECT_EQ(service_request.pdu_session_status.iei, PDU_SESSION_STATUS);
  EXPECT_EQ(service_request.pdu_session_status.len, 0x02);
  EXPECT_EQ(service_request.pdu_session_status.pduSessionStatus, 0x0020);
  // Verify service request is rejected as
  // conditional IE Uplink Status is not present
  EXPECT_EQ(
      RETURNok, amf_handle_service_request(
                    AMF_UE_NGAP_ID, &service_request, decode_status));

  // Verify UE still in CONNECTED MODE though initial ue message is received
  EXPECT_EQ(REGISTERED_IDLE, ue_context->mm_state);
}

TEST(test_pdu_negative, test_unknown_pdu_session_type) {
  amf_nas_message_t msg = {};

  // build uplinknastransport //
  // uplink nas transport(pdu session request)
  uint8_t pdu[44] = {0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01,
                     0xc1, 0xff, 0xff, 0x95, 0xa1, 0x28, 0x01, 0x00, 0x7b,
                     0x00, 0x07, 0x80, 0x00, 0x0a, 0x00, 0x00, 0x0d, 0x00,
                     0x12, 0x01, 0x81, 0x22, 0x01, 0x01, 0x25, 0x09, 0x08,
                     0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  uint32_t len = sizeof(pdu) / sizeof(uint8_t);

  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(&pdu_sess_est_req, pdu, len);

  EXPECT_EQ(decode_res, true);

  amf_ue_ngap_id_t ue_id = 1;

  // creating ue_context
  ue_m5gmm_context_s* ue_context = amf_create_new_ue_context();
  ue_context_map.insert(
      std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(ue_id, ue_context));

  M5GSmCause cause = amf_smf_get_smcause(ue_id, &pdu_sess_est_req);
  EXPECT_EQ(cause, M5GSmCause::UNKNOWN_PDU_SESSION_TYPE);

  ue_context_map.clear();
  delete ue_context;
}

TEST(test_pdu_negative, test_pdu_unknown_dnn_missing_dnn) {
  amf_nas_message_t msg = {};

  // build uplinknastransport //
  // uplink nas transport(pdu session request)
  uint8_t pdu[33] = {0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01,
                     0xc1, 0xff, 0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b,
                     0x00, 0x07, 0x80, 0x00, 0x0a, 0x00, 0x00, 0x0d, 0x00,
                     0x12, 0x01, 0x81, 0x22, 0x01, 0x01};
  uint32_t len    = sizeof(pdu) / sizeof(uint8_t);

  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(&pdu_sess_est_req, pdu, len);

  EXPECT_EQ(decode_res, true);

  amf_ue_ngap_id_t ue_id = 1;

  ue_m5gmm_context_s* ue_context = amf_create_new_ue_context();
  ue_context_map.insert(
      std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(ue_id, ue_context));
  M5GSmCause cause = amf_smf_get_smcause(ue_id, &pdu_sess_est_req);
  EXPECT_EQ(cause, M5GSmCause::MISSING_OR_UNKNOWN_DNN);

  ue_context_map.clear();
  delete ue_context;
}

TEST(test_pdu_negative, test_pdu_invalid_pdu_identity) {
  amf_nas_message_t msg = {};

  // build uplinknastransport //
  // uplink nas transport(pdu session request)
  uint8_t pdu[44] = {0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01,
                     0xc1, 0xff, 0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b,
                     0x00, 0x07, 0x80, 0x00, 0x0a, 0x00, 0x00, 0x0d, 0x00,
                     0x12, 0x01, 0x81, 0x22, 0x01, 0x01, 0x25, 0x09, 0x08,
                     0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  uint32_t len = sizeof(pdu) / sizeof(uint8_t);

  ULNASTransportMsg pdu_sess_est_req;
  bool decode_res = false;
  memset(&pdu_sess_est_req, 0, sizeof(ULNASTransportMsg));

  decode_res = decode_ul_nas_transport_msg(&pdu_sess_est_req, pdu, len);

  EXPECT_EQ(decode_res, true);

  amf_ue_ngap_id_t ue_id = 1;
  uint8_t pdu_session_id = 1;

  // creating ue_context
  ue_m5gmm_context_s* ue_context = amf_create_new_ue_context();
  ue_context_map.insert(
      std::pair<amf_ue_ngap_id_t, ue_m5gmm_context_s*>(ue_id, ue_context));
  std::shared_ptr<smf_context_t> smf_ctx =
      amf_insert_smf_context(ue_context, pdu_session_id);
  smf_ctx->pdu_session_state = ACTIVE;

  for (int req_cnt = 0;
       req_cnt < MAX_UE_INITIAL_PDU_SESSION_ESTABLISHMENT_REQ_ALLOWED;
       req_cnt++) {
    amf_smf_get_smcause(ue_id, &pdu_sess_est_req);
  }
  M5GSmCause cause_dup = amf_smf_get_smcause(ue_id, &pdu_sess_est_req);
  EXPECT_EQ(cause_dup, M5GSmCause::INVALID_PDU_SESSION_IDENTITY);

  ue_context_map.clear();
  delete ue_context;
}

TEST(test_optional_pdu, test_pdu_session_accept_optional) {
  uint32_t bytes         = 0;
  uint32_t container_len = 0;
  bstring buffer;
  amf_nas_message_t msg                                   = {};
  protocol_configuration_options_t* msg_accept_pco        = nullptr;
  protocol_configuration_options_t* decode_msg_accept_pco = nullptr;

  // build downlinknastransport
  // downlink nas transport(pdu session accept)
  uint8_t pdu[82] = {
      0x7e, 0x00, 0x68, 0x01, 0x00, 0x4a, 0x2e, 0x01, 0x01, 0xc2, 0x11, 0x00,
      0x09, 0x02, 0x00, 0x06, 0x31, 0x31, 0x01, 0x01, 0x02, 0x09, 0x06, 0x0a,
      0x00, 0x01, 0x0a, 0x00, 0x01, 0x29, 0x05, 0x01, 0x05, 0x05, 0x05, 0x1e,
      0x22, 0x04, 0x03, 0x03, 0x06, 0x09, 0x79, 0x00, 0x06, 0x09, 0x20, 0x41,
      0x01, 0x01, 0x09, 0x7b, 0x00, 0x0f, 0x80, 0x00, 0x0d, 0x04, 0x08, 0x08,
      0x08, 0x08, 0x00, 0x0c, 0x04, 0xc0, 0xa8, 0x78, 0x0d, 0x25, 0x09, 0x08,
      0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x45, 0x54, 0x12, 0x01};

  uint32_t len = sizeof(pdu) / sizeof(uint8_t);

  NAS5GPktSnapShot nas5g_pkt_snap;
  DLNASTransportMsg pdu_sess_accept;
  int decode_res = 0;
  memset(&pdu_sess_accept, 0, sizeof(DLNASTransportMsg));
  SmfMsg* smf_msg = &pdu_sess_accept.payload_container.smf_msg;

  msg_accept_pco =
      &(smf_msg->msg.pdu_session_estab_accept.protocolconfigurationoptions.pco);
  decode_res =
      pdu_sess_accept.DecodeDLNASTransportMsg(&pdu_sess_accept, pdu, len);

  EXPECT_GT(decode_res, 0);

  // PDU Session type : IPv4 (pdu_address.type_val = 1)
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.pdu_address.type_val, 1);
  // SSC mode check
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val, 1);
  // NSSAI
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sst, 3);
  uint8_t sd[3] = {0x03, 0x06, 0x09};
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sd[0], sd[0]);
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sd[1], sd[1]);
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sd[2], sd[2]);
  // DNN
  uint8_t dnn[9] = {0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x45, 0x54};
  EXPECT_EQ(
      memcmp(
          smf_msg->msg.pdu_session_estab_accept.dnn.dnn, dnn,
          smf_msg->msg.pdu_session_estab_accept.dnn.len),
      0);

  buffer = bfromcstralloc(len, "\0");
  bytes  = pdu_sess_accept.EncodeDLNASTransportMsg(
      &pdu_sess_accept, buffer->data, len);
  EXPECT_GT(bytes, 0);
  DLNASTransportMsg decode_pdu_sess_accept;
  memset(&decode_pdu_sess_accept, 0, sizeof(DLNASTransportMsg));
  decode_res = decode_pdu_sess_accept.DecodeDLNASTransportMsg(
      &decode_pdu_sess_accept, pdu, len);

  smf_msg = &decode_pdu_sess_accept.payload_container.smf_msg;
  EXPECT_GT(decode_res, 0);

  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.pdu_address.type_val, 1);
  // SSC mode check
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.ssc_mode.mode_val, 1);
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sst, 3);
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sd[0], sd[0]);
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sd[1], sd[1]);
  EXPECT_EQ(smf_msg->msg.pdu_session_estab_accept.nssai.sd[2], sd[2]);
  EXPECT_EQ(
      memcmp(
          smf_msg->msg.pdu_session_estab_accept.dnn.dnn, dnn,
          smf_msg->msg.pdu_session_estab_accept.dnn.len),
      0);

  bdestroy(buffer);
  decode_msg_accept_pco =
      &(smf_msg->msg.pdu_session_estab_accept.protocolconfigurationoptions.pco);

  // Clean up the PCO contents
  sm_free_protocol_configuration_options(&decode_msg_accept_pco);
  // Clean up the PCO contents
  sm_free_protocol_configuration_options(&msg_accept_pco);
}

}  // namespace magma5g
