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
#include <gtest/gtest.h>
#include <chrono>
#include <thread>

#include "lte/gateway/c/core/oai/test/mock_tasks/mock_tasks.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/common/common_types.h"
#include "lte/gateway/c/core/oai/include/amf_config.hpp"
}
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_state_manager.hpp"
#include "lte/gateway/c/core/oai/test/amf/amf_app_test_util.h"
#include "lte/gateway/c/core/oai/lib/secu/secu_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_identity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GNasEnums.h"
#include "lte/gateway/c/core/oai/test/amf/util_s6a_update_location.hpp"

using ::testing::Test;

namespace magma5g {

extern task_zmq_ctx_s amf_app_task_zmq_ctx;

TEST(TestAMFStateConverter, TestGutiToString) {
  guti_m5_t guti1, guti2;
  guti1.guamfi.plmn.mcc_digit1 = 2;
  guti1.guamfi.plmn.mcc_digit2 = 2;
  guti1.guamfi.plmn.mcc_digit3 = 2;
  guti1.guamfi.plmn.mnc_digit1 = 4;
  guti1.guamfi.plmn.mnc_digit2 = 5;
  guti1.guamfi.plmn.mnc_digit3 = 6;
  guti1.guamfi.amf_regionid = 1;
  guti1.guamfi.amf_set_id = 1;
  guti1.guamfi.amf_pointer = 0;
  guti1.m_tmsi = 0X212e5025;

  std::string guti1_str =
      AmfNasStateConverter::amf_app_convert_guti_m5_to_string(guti1);

  AmfNasStateConverter::amf_app_convert_string_to_guti_m5(guti1_str, &guti2);

  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit1, guti2.guamfi.plmn.mcc_digit1);
  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit2, guti2.guamfi.plmn.mcc_digit2);
  EXPECT_EQ(guti1.guamfi.plmn.mcc_digit3, guti2.guamfi.plmn.mcc_digit3);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit1, guti2.guamfi.plmn.mnc_digit1);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit2, guti2.guamfi.plmn.mnc_digit2);
  EXPECT_EQ(guti1.guamfi.plmn.mnc_digit3, guti2.guamfi.plmn.mnc_digit3);
  EXPECT_EQ(guti1.guamfi.amf_regionid, guti2.guamfi.amf_regionid);
  EXPECT_EQ(guti1.guamfi.amf_set_id, guti2.guamfi.amf_set_id);
  EXPECT_EQ(guti1.guamfi.amf_pointer, guti2.guamfi.amf_pointer);
  EXPECT_EQ(guti1.m_tmsi, guti2.m_tmsi);
}

TEST(TestAMFStateConverter, TestStateToProto) {
  // Guti setup
  guti_m5_t guti1;
  memset(&guti1, 0, sizeof(guti1));

  guti1.guamfi.plmn.mcc_digit1 = 2;
  guti1.guamfi.plmn.mcc_digit2 = 2;
  guti1.guamfi.plmn.mcc_digit3 = 2;
  guti1.guamfi.plmn.mnc_digit1 = 4;
  guti1.guamfi.plmn.mnc_digit2 = 5;
  guti1.guamfi.plmn.mnc_digit3 = 6;
  guti1.guamfi.amf_regionid = 1;
  guti1.guamfi.amf_set_id = 1;
  guti1.guamfi.amf_pointer = 0;
  guti1.m_tmsi = 556683301;

  amf_app_desc_t amf_app_desc1 = {}, amf_app_desc2 = {};
  magma::lte::oai::MmeNasState state_proto = magma::lte::oai::MmeNasState();
  uint64_t data = 0;

  amf_app_desc1.amf_app_ue_ngap_id_generator = 0x05;
  amf_app_desc1.amf_ue_contexts.imsi_amf_ue_id_htbl.insert(1, 10);
  amf_app_desc1.amf_ue_contexts.tun11_ue_context_htbl.insert(2, 20);
  amf_app_desc1.amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.insert(3, 30);
  amf_app_desc1.amf_ue_contexts.guti_ue_context_htbl.insert(guti1, 40);

  AmfNasStateConverter::state_to_proto(&amf_app_desc1, &state_proto);

  AmfNasStateConverter::proto_to_state(state_proto, &amf_app_desc2);

  EXPECT_EQ(amf_app_desc1.amf_app_ue_ngap_id_generator,
            amf_app_desc2.amf_app_ue_ngap_id_generator);

  EXPECT_EQ(amf_app_desc2.amf_ue_contexts.imsi_amf_ue_id_htbl.get(1, &data),
            magma::MAP_OK);
  EXPECT_EQ(data, 10);
  data = 0;

  EXPECT_EQ(amf_app_desc2.amf_ue_contexts.tun11_ue_context_htbl.get(2, &data),
            magma::MAP_OK);
  EXPECT_EQ(data, 20);
  data = 0;

  EXPECT_EQ(amf_app_desc2.amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.get(
                3, &data),
            magma::MAP_OK);
  EXPECT_EQ(data, 30);
  data = 0;

  EXPECT_EQ(
      amf_app_desc2.amf_ue_contexts.guti_ue_context_htbl.get(guti1, &data),
      magma::MAP_OK);
  EXPECT_EQ(data, 40);
}

TEST(TestAMFStateConverter, TestUEm5gmmContextToProto) {
  ue_m5gmm_context_t ue_m5gmm_context1 = {}, ue_m5gmm_context2 = {};
  magma::lte::oai::UeContext ue_context_proto = magma::lte::oai::UeContext();

  ue_m5gmm_context1.ue_context_rel_cause = NGAP_INVALID_CAUSE;
  ue_m5gmm_context1.cm_state = M5GCM_CONNECTED;
  ue_m5gmm_context1.mm_state = REGISTERED_IDLE;

  ue_m5gmm_context1.sctp_assoc_id_key = 1;
  ue_m5gmm_context1.gnb_ue_ngap_id = 0x09;
  ue_m5gmm_context1.gnb_ngap_id_key = INVALID_GNB_UE_NGAP_ID_KEY;

  ue_m5gmm_context1.amf_context.apn_config_profile.nb_apns = 1;
  strncpy(ue_m5gmm_context1.amf_context.apn_config_profile.apn_configuration[0]
              .service_selection,
          "internet", 8);
  ue_m5gmm_context1.amf_context.apn_config_profile.apn_configuration[0]
      .service_selection_length = 8;

  ue_m5gmm_context1.amf_teid_n11 = 0;

  ue_m5gmm_context1.amf_context.subscribed_ue_ambr.br_unit = KBPS;
  ue_m5gmm_context1.amf_context.subscribed_ue_ambr.br_ul = 1000;
  ue_m5gmm_context1.amf_context.subscribed_ue_ambr.br_dl = 10000;

  ue_m5gmm_context1.paging_context.paging_retx_count = 0;

  AmfNasStateConverter::ue_m5gmm_context_to_proto(&ue_m5gmm_context1,
                                                  &ue_context_proto);

  AmfNasStateConverter::proto_to_ue_m5gmm_context(ue_context_proto,
                                                  &ue_m5gmm_context2);

  EXPECT_EQ(ue_m5gmm_context1.ue_context_rel_cause,
            ue_m5gmm_context2.ue_context_rel_cause);
  EXPECT_EQ(ue_m5gmm_context1.cm_state, ue_m5gmm_context2.cm_state);
  EXPECT_EQ(ue_m5gmm_context1.mm_state, ue_m5gmm_context2.mm_state);

  EXPECT_EQ(ue_m5gmm_context1.sctp_assoc_id_key,
            ue_m5gmm_context2.sctp_assoc_id_key);
  EXPECT_EQ(ue_m5gmm_context1.gnb_ue_ngap_id, ue_m5gmm_context2.gnb_ue_ngap_id);
  EXPECT_EQ(ue_m5gmm_context1.gnb_ngap_id_key,
            ue_m5gmm_context2.gnb_ngap_id_key);

  EXPECT_EQ(ue_m5gmm_context1.amf_context.apn_config_profile.nb_apns,
            ue_m5gmm_context2.amf_context.apn_config_profile.nb_apns);

  std::string str_in =
      ue_m5gmm_context1.amf_context.apn_config_profile.apn_configuration[0]
          .service_selection;
  std::string str_out =
      ue_m5gmm_context2.amf_context.apn_config_profile.apn_configuration[0]
          .service_selection;
  EXPECT_EQ(str_in, str_out);

  EXPECT_EQ(ue_m5gmm_context1.amf_teid_n11, ue_m5gmm_context2.amf_teid_n11);

  EXPECT_EQ(ue_m5gmm_context1.amf_context.subscribed_ue_ambr.br_unit,
            ue_m5gmm_context2.amf_context.subscribed_ue_ambr.br_unit);
  EXPECT_EQ(ue_m5gmm_context1.amf_context.subscribed_ue_ambr.br_ul,
            ue_m5gmm_context2.amf_context.subscribed_ue_ambr.br_ul);
  EXPECT_EQ(ue_m5gmm_context1.amf_context.subscribed_ue_ambr.br_dl,
            ue_m5gmm_context2.amf_context.subscribed_ue_ambr.br_dl);

  EXPECT_EQ(ue_m5gmm_context1.paging_context.paging_retx_count,
            ue_m5gmm_context2.paging_context.paging_retx_count);
}

TEST(TestAMFStateConverter, TestAmfContextStateToProto) {
#define AMF_CAUSE_SUCCESS 1
  amf_context_t amf_ctx1 = {}, amf_ctx2 = {};
  magma::lte::oai::EmmContext emm_context_proto = magma::lte::oai::EmmContext();

  amf_ctx1.imsi64 = 222456000000101;
  amf_ctx1.imsi.u.num.digit1 = 2;
  amf_ctx1.imsi.u.num.digit2 = 2;
  amf_ctx1.imsi.u.num.digit3 = 2;
  amf_ctx1.imsi.u.num.digit4 = 4;
  amf_ctx1.imsi.u.num.digit5 = 5;
  amf_ctx1.imsi.u.num.digit6 = 6;
  amf_ctx1.imsi.u.num.digit7 = 0;
  amf_ctx1.imsi.u.num.digit8 = 0;
  amf_ctx1.imsi.u.num.digit9 = 0;
  amf_ctx1.imsi.u.num.digit10 = 0;
  amf_ctx1.imsi.u.num.digit11 = 0;
  amf_ctx1.imsi.u.num.digit12 = 0;
  amf_ctx1.imsi.u.num.digit13 = 1;
  amf_ctx1.imsi.u.num.digit14 = 0;
  amf_ctx1.imsi.u.num.digit15 = 1;
  amf_ctx1.saved_imsi64 = 222456000000101;

  // imei
  amf_ctx1.imei.length = 10;
  amf_ctx1.imei.u.num.tac2 = 2;
  amf_ctx1.imei.u.num.tac1 = 1;
  amf_ctx1.imei.u.num.tac3 = 3;
  amf_ctx1.imei.u.num.tac4 = 4;
  amf_ctx1.imei.u.num.tac5 = 5;
  amf_ctx1.imei.u.num.tac6 = 6;
  amf_ctx1.imei.u.num.tac7 = 7;
  amf_ctx1.imei.u.num.tac8 = 8;
  amf_ctx1.imei.u.num.snr1 = 1;
  amf_ctx1.imei.u.num.snr2 = 2;
  amf_ctx1.imei.u.num.snr3 = 3;
  amf_ctx1.imei.u.num.snr4 = 4;
  amf_ctx1.imei.u.num.snr5 = 5;
  amf_ctx1.imei.u.num.snr6 = 6;
  amf_ctx1.imei.u.num.parity = 1;
  amf_ctx1.imei.u.num.cdsd = 8;
  for (int i = 0; i < IMEI_BCD8_SIZE; i++) {
    amf_ctx1.imei.u.value[i] = i;
  }

  // imeisv
  amf_ctx1.imeisv.length = 10;
  amf_ctx1.imeisv.u.num.tac2 = 2;
  amf_ctx1.imeisv.u.num.tac1 = 1;
  amf_ctx1.imeisv.u.num.tac3 = 3;
  amf_ctx1.imeisv.u.num.tac4 = 4;
  amf_ctx1.imeisv.u.num.tac5 = 5;
  amf_ctx1.imeisv.u.num.tac6 = 6;
  amf_ctx1.imeisv.u.num.tac7 = 7;
  amf_ctx1.imeisv.u.num.tac8 = 8;
  amf_ctx1.imeisv.u.num.snr1 = 1;
  amf_ctx1.imeisv.u.num.snr2 = 2;
  amf_ctx1.imeisv.u.num.snr3 = 3;
  amf_ctx1.imeisv.u.num.snr4 = 4;
  amf_ctx1.imeisv.u.num.snr5 = 5;
  amf_ctx1.imeisv.u.num.snr6 = 6;
  amf_ctx1.imeisv.u.num.parity = 1;
  for (int i = 0; i < IMEISV_BCD8_SIZE; i++) {
    amf_ctx1.imeisv.u.value[i] = i;
  }
  amf_ctx1.amf_cause = AMF_CAUSE_SUCCESS;
  amf_ctx1.amf_fsm_state = AMF_DEREGISTERED;

  amf_ctx1.m5gsregistrationtype = AMF_REGISTRATION_TYPE_INITIAL;
  amf_ctx1.member_present_mask |= AMF_CTXT_MEMBER_SECURITY;
  amf_ctx1.member_valid_mask |= AMF_CTXT_MEMBER_SECURITY;
  amf_ctx1.is_dynamic = true;
  amf_ctx1.is_registered = true;
  amf_ctx1.is_initial_identity_imsi = true;
  amf_ctx1.is_guti_based_registered = true;
  amf_ctx1.is_imsi_only_detach = false;

  // originating_tai
  amf_ctx1.originating_tai.plmn.mcc_digit1 = 2;
  amf_ctx1.originating_tai.plmn.mcc_digit2 = 2;
  amf_ctx1.originating_tai.plmn.mcc_digit3 = 2;
  amf_ctx1.originating_tai.plmn.mnc_digit3 = 6;
  amf_ctx1.originating_tai.plmn.mnc_digit2 = 5;
  amf_ctx1.originating_tai.plmn.mnc_digit1 = 4;
  amf_ctx1.originating_tai.tac = 1;

  amf_ctx1.ksi = 0x06;

  uint8_t pdu_sess_id = 1;
  smf_context_t smf_ctx = {};
  smf_ctx.pdu_session_state = ACTIVE;
  amf_ctx1.smf_ctxt_map[pdu_sess_id] = std::make_shared<smf_context_t>(smf_ctx);

  AmfNasStateConverter::amf_context_to_proto(&amf_ctx1, &emm_context_proto);
  AmfNasStateConverter::proto_to_amf_context(emm_context_proto, &amf_ctx2);

  EXPECT_EQ(amf_ctx1.imsi64, amf_ctx2.imsi64);
  EXPECT_EQ(amf_ctx1.saved_imsi64, amf_ctx2.saved_imsi64);
  EXPECT_EQ(amf_ctx1.amf_cause, amf_ctx2.amf_cause);
  EXPECT_EQ(amf_ctx1.m5gsregistrationtype, amf_ctx2.m5gsregistrationtype);
  EXPECT_EQ(amf_ctx1.member_present_mask, amf_ctx2.member_present_mask);
  EXPECT_EQ(amf_ctx1.member_valid_mask, amf_ctx2.member_valid_mask);
  EXPECT_EQ(amf_ctx1.is_dynamic, amf_ctx2.is_dynamic);
  EXPECT_EQ(amf_ctx1.is_registered, amf_ctx2.is_registered);
  EXPECT_EQ(amf_ctx1.is_initial_identity_imsi,
            amf_ctx2.is_initial_identity_imsi);
  EXPECT_EQ(amf_ctx1.is_guti_based_registered,
            amf_ctx2.is_guti_based_registered);
  EXPECT_EQ(amf_ctx1.is_imsi_only_detach, amf_ctx2.is_imsi_only_detach);
  EXPECT_EQ(memcmp(&amf_ctx1.imsi, &amf_ctx2.imsi, sizeof(amf_ctx1.imsi)), 0);
  EXPECT_EQ(amf_ctx1.imsi.u.num.digit1, amf_ctx2.imsi.u.num.digit1);
  EXPECT_EQ(amf_ctx1.amf_fsm_state, amf_ctx2.amf_fsm_state);
  EXPECT_EQ(memcmp(&amf_ctx1.imei, &amf_ctx2.imei, sizeof(amf_ctx1.imei)), 0);
  EXPECT_EQ(memcmp(&amf_ctx1.imeisv, &amf_ctx2.imeisv, sizeof(amf_ctx1.imeisv)),
            0);
  EXPECT_EQ(amf_ctx1.ksi, amf_ctx2.ksi);
  EXPECT_EQ(memcmp(&amf_ctx1.originating_tai, &amf_ctx2.originating_tai,
                   sizeof(amf_ctx1.originating_tai)),
            0);
  EXPECT_EQ(amf_ctx1.smf_ctxt_map.size(), amf_ctx2.smf_ctxt_map.size());
  auto map1 = amf_ctx1.smf_ctxt_map.find(pdu_sess_id);
  auto map2 = amf_ctx2.smf_ctxt_map.find(pdu_sess_id);
  EXPECT_EQ(map1->second.get()->pdu_session_state,
            map2->second.get()->pdu_session_state);
}

TEST(TestAMFStateConverter, TestAMFSecurityContextToProto) {
  amf_security_context_t state_amf_security_context_1 = {};
  amf_security_context_t state_amf_security_context_2 = {};
  // EmmSecurityProto
  magma::lte::oai::EmmSecurityContext emm_security_context_proto =
      magma::lte::oai::EmmSecurityContext();
  // amf_security_context setup
  state_amf_security_context_1.sc_type = SECURITY_CTX_TYPE_NOT_AVAILABLE;
  state_amf_security_context_1.eksi = 1;
  state_amf_security_context_1.vector_index = 1;
  state_amf_security_context_1.dl_count.overflow = 2;
  state_amf_security_context_1.dl_count.seq_num = 1;
  state_amf_security_context_1.ul_count.overflow = 1;
  state_amf_security_context_1.ul_count.seq_num = 2;
  state_amf_security_context_1.kenb_ul_count.overflow = 1;
  state_amf_security_context_1.kenb_ul_count.seq_num = 1;
  state_amf_security_context_1.direction_decode = SECU_DIRECTION_UPLINK;
  state_amf_security_context_1.direction_encode = SECU_DIRECTION_DOWNLINK;

  AmfNasStateConverter::amf_security_context_to_proto(
      &state_amf_security_context_1, &emm_security_context_proto);
  AmfNasStateConverter::proto_to_amf_security_context(
      emm_security_context_proto, &state_amf_security_context_2);

  EXPECT_EQ(state_amf_security_context_1.sc_type,
            state_amf_security_context_2.sc_type);
  EXPECT_EQ(state_amf_security_context_1.eksi,
            state_amf_security_context_2.eksi);
  EXPECT_EQ(state_amf_security_context_1.vector_index,
            state_amf_security_context_2.vector_index);

  // Count values
  EXPECT_EQ(state_amf_security_context_1.dl_count.overflow,
            state_amf_security_context_2.dl_count.overflow);
  EXPECT_EQ(state_amf_security_context_1.dl_count.seq_num,
            state_amf_security_context_2.dl_count.seq_num);
  EXPECT_EQ(state_amf_security_context_1.ul_count.overflow,
            state_amf_security_context_2.ul_count.overflow);
  EXPECT_EQ(state_amf_security_context_1.ul_count.seq_num,
            state_amf_security_context_2.ul_count.seq_num);
  EXPECT_EQ(state_amf_security_context_1.kenb_ul_count.overflow,
            state_amf_security_context_2.kenb_ul_count.overflow);
  EXPECT_EQ(state_amf_security_context_1.kenb_ul_count.seq_num,
            state_amf_security_context_2.kenb_ul_count.seq_num);

  // Security algorithm
  EXPECT_EQ(state_amf_security_context_1.direction_decode,
            state_amf_security_context_2.direction_decode);
  EXPECT_EQ(state_amf_security_context_1.direction_encode,
            state_amf_security_context_2.direction_encode);
}

TEST(TestAMFStateConverter, TestSMFContextToProto) {
  smf_context_t smf_context1 = {}, smf_context2 = {};
  magma::lte::oai::SmfContext state_smf_proto = magma::lte::oai::SmfContext();
  smf_context1.pdu_session_state = ACTIVE;
  smf_context1.pdu_session_version = 0;
  smf_context1.n_active_pdus = 0;
  smf_context1.is_emergency = false;

  // selected ambr
  smf_context1.selected_ambr.dl_ambr_unit = M5GSessionAmbrUnit::MULTIPLES_1KBPS;
  smf_context1.selected_ambr.dl_session_ambr = 10000;
  smf_context1.selected_ambr.ul_ambr_unit = M5GSessionAmbrUnit::MULTIPLES_1KBPS;
  smf_context1.selected_ambr.ul_session_ambr = 1000;

  // gtp_tunnel_id
  // gnb
  smf_context1.gtp_tunnel_id.gnb_gtp_teid = 1;
  smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[0] = 0xc0;
  smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[1] = 0xa8;
  smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[2] = 0x3c;
  smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[3] = 0x96;
  // upf
  smf_context1.gtp_tunnel_id.upf_gtp_teid[0] = 0x0;
  smf_context1.gtp_tunnel_id.upf_gtp_teid[1] = 0x0;
  smf_context1.gtp_tunnel_id.upf_gtp_teid[2] = 0x0;
  smf_context1.gtp_tunnel_id.upf_gtp_teid[3] = 0x1;

  smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[0] = 0xc0;
  smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[1] = 0xa8;
  smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[2] = 0x3c;
  smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[3] = 0xad;

  // pdu address
  smf_context1.pdu_address.pdn_type = IPv4;
  smf_context1.pdu_address.ipv4_address.s_addr = 0x0441a8c0;

  // apn_ambr
  smf_context1.apn_ambr.br_dl = 10000;
  smf_context1.apn_ambr.br_ul = 1000;
  smf_context1.apn_ambr.br_unit = KBPS;

  // smf_proc_data
  smf_context1.smf_proc_data.pdu_session_id = 1;
  smf_context1.smf_proc_data.pdu_session_type = M5GPduSessionType::IPV4;
  smf_context1.smf_proc_data.pti = 0x01;
  smf_context1.smf_proc_data.ssc_mode = SSC_MODE_3;
  smf_context1.smf_proc_data.max_uplink = 0xFF;
  smf_context1.smf_proc_data.max_downlink = 0xFF;

  smf_context1.retransmission_count = 1;

  // PCO
  smf_context1.pco.num_protocol_or_container_id = 2;
  smf_context1.pco.protocol_or_container_ids[0].id =
      PCO_CI_P_CSCF_IPV6_ADDRESS_REQUEST;
  bstring test_string1 = bfromcstr("teststring");
  smf_context1.pco.protocol_or_container_ids[0].contents = test_string1;
  smf_context1.pco.protocol_or_container_ids[0].length = blength(test_string1);
  smf_context1.pco.protocol_or_container_ids[1].id =
      PCO_CI_DSMIPV6_IPV4_HOME_AGENT_ADDRESS;
  bstring test_string2 = bfromcstr("longer.test.string");
  smf_context1.pco.protocol_or_container_ids[1].contents = test_string2;
  smf_context1.pco.protocol_or_container_ids[1].length = blength(test_string2);

  // dnn
  smf_context1.dnn = "internet";

  // nssai
  smf_context1.requested_nssai.sd[0] = 0x03;
  smf_context1.requested_nssai.sd[1] = 0x06;
  smf_context1.requested_nssai.sd[2] = 0x09;
  smf_context1.requested_nssai.sst = 1;

  // Qos
  smf_context1.smf_proc_data.qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_identifier = 9;
  smf_context1.smf_proc_data.qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
      .non_dynamic_5QI_desc.fiveQI = 9;
  smf_context1.smf_proc_data.qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .priority_level = 1;
  smf_context1.smf_proc_data.qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_cap = SHALL_NOT_TRIGGER_PRE_EMPTION;
  smf_context1.smf_proc_data.qos_flow_list.item[0]
      .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
      .pre_emption_vul = NOT_PREEMPTABLE;

  AmfNasStateConverter::smf_context_to_proto(&smf_context1, &state_smf_proto);
  AmfNasStateConverter::proto_to_smf_context(state_smf_proto, &smf_context2);

  EXPECT_EQ(smf_context1.pdu_session_state, smf_context2.pdu_session_state);
  EXPECT_EQ(smf_context1.pdu_session_version, smf_context2.pdu_session_version);
  EXPECT_EQ(smf_context1.n_active_pdus, smf_context2.n_active_pdus);
  EXPECT_EQ(smf_context1.is_emergency, smf_context2.is_emergency);

  EXPECT_EQ(smf_context1.selected_ambr.dl_ambr_unit,
            smf_context2.selected_ambr.dl_ambr_unit);
  EXPECT_EQ(smf_context1.selected_ambr.dl_session_ambr,
            smf_context2.selected_ambr.dl_session_ambr);
  EXPECT_EQ(smf_context1.selected_ambr.ul_ambr_unit,
            smf_context2.selected_ambr.ul_ambr_unit);
  EXPECT_EQ(smf_context1.selected_ambr.ul_session_ambr,
            smf_context2.selected_ambr.ul_session_ambr);

  EXPECT_EQ(smf_context1.gtp_tunnel_id.gnb_gtp_teid,
            smf_context2.gtp_tunnel_id.gnb_gtp_teid);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[0],
            smf_context2.gtp_tunnel_id.gnb_gtp_teid_ip_addr[0]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[1],
            smf_context2.gtp_tunnel_id.gnb_gtp_teid_ip_addr[1]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[2],
            smf_context2.gtp_tunnel_id.gnb_gtp_teid_ip_addr[2]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.gnb_gtp_teid_ip_addr[3],
            smf_context2.gtp_tunnel_id.gnb_gtp_teid_ip_addr[3]);

  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid[0],
            smf_context2.gtp_tunnel_id.upf_gtp_teid[0]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid[1],
            smf_context2.gtp_tunnel_id.upf_gtp_teid[1]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid[2],
            smf_context2.gtp_tunnel_id.upf_gtp_teid[2]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid[3],
            smf_context2.gtp_tunnel_id.upf_gtp_teid[3]);

  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[0],
            smf_context2.gtp_tunnel_id.upf_gtp_teid_ip_addr[0]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[1],
            smf_context2.gtp_tunnel_id.upf_gtp_teid_ip_addr[1]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[2],
            smf_context2.gtp_tunnel_id.upf_gtp_teid_ip_addr[2]);
  EXPECT_EQ(smf_context1.gtp_tunnel_id.upf_gtp_teid_ip_addr[3],
            smf_context2.gtp_tunnel_id.upf_gtp_teid_ip_addr[3]);

  EXPECT_EQ(smf_context1.pdu_address.pdn_type,
            smf_context2.pdu_address.pdn_type);
  EXPECT_EQ(smf_context1.pdu_address.ipv4_address.s_addr,
            smf_context2.pdu_address.ipv4_address.s_addr);

  EXPECT_EQ(smf_context1.apn_ambr.br_dl, smf_context2.apn_ambr.br_dl);
  EXPECT_EQ(smf_context1.apn_ambr.br_ul, smf_context1.apn_ambr.br_ul);
  EXPECT_EQ(smf_context1.apn_ambr.br_unit, smf_context1.apn_ambr.br_unit);

  EXPECT_EQ(smf_context1.smf_proc_data.pdu_session_id,
            smf_context2.smf_proc_data.pdu_session_id);
  EXPECT_EQ(smf_context1.smf_proc_data.pdu_session_type,
            smf_context2.smf_proc_data.pdu_session_type);
  EXPECT_EQ(smf_context1.smf_proc_data.pti, smf_context2.smf_proc_data.pti);
  EXPECT_EQ(smf_context1.smf_proc_data.ssc_mode,
            smf_context2.smf_proc_data.ssc_mode);
  EXPECT_EQ(smf_context1.smf_proc_data.max_uplink,
            smf_context2.smf_proc_data.max_uplink);
  EXPECT_EQ(smf_context1.smf_proc_data.max_downlink,
            smf_context2.smf_proc_data.max_downlink);

  EXPECT_EQ(smf_context1.retransmission_count,
            smf_context2.retransmission_count);

  EXPECT_EQ(smf_context1.pco.num_protocol_or_container_id,
            smf_context2.pco.num_protocol_or_container_id);
  EXPECT_EQ(smf_context1.pco.protocol_or_container_ids[0].id,
            smf_context2.pco.protocol_or_container_ids[0].id);

  std::string contents;
  BSTRING_TO_STRING(smf_context2.pco.protocol_or_container_ids[0].contents,
                    &contents);
  EXPECT_EQ(contents, "teststring");
  EXPECT_EQ(smf_context1.pco.protocol_or_container_ids[0].length,
            smf_context2.pco.protocol_or_container_ids[0].length);

  contents = {};

  EXPECT_EQ(smf_context1.pco.protocol_or_container_ids[1].id,
            smf_context2.pco.protocol_or_container_ids[1].id);
  BSTRING_TO_STRING(smf_context2.pco.protocol_or_container_ids[1].contents,
                    &contents);
  EXPECT_EQ(contents, "longer.test.string");
  EXPECT_EQ(smf_context1.pco.protocol_or_container_ids[1].length,
            smf_context2.pco.protocol_or_container_ids[1].length);

  EXPECT_EQ(smf_context1.dnn, smf_context2.dnn);

  EXPECT_EQ(smf_context1.requested_nssai.sd[0],
            smf_context2.requested_nssai.sd[0]);
  EXPECT_EQ(smf_context1.requested_nssai.sd[1],
            smf_context2.requested_nssai.sd[1]);
  EXPECT_EQ(smf_context1.requested_nssai.sd[2],
            smf_context2.requested_nssai.sd[2]);
  EXPECT_EQ(smf_context1.requested_nssai.sst, smf_context2.requested_nssai.sst);

  EXPECT_EQ(smf_context1.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_identifier,
            smf_context2.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_identifier);
  EXPECT_EQ(smf_context1.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
                .non_dynamic_5QI_desc.fiveQI,
            smf_context2.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
                .non_dynamic_5QI_desc.fiveQI);

  EXPECT_EQ(smf_context1.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .priority_level,
            smf_context2.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .priority_level);

  EXPECT_EQ(smf_context1.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_cap,
            smf_context2.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_cap);

  EXPECT_EQ(smf_context1.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_vul,
            smf_context2.smf_proc_data.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_vul);

  bdestroy(smf_context2.pco.protocol_or_container_ids[0].contents);
  bdestroy(smf_context2.pco.protocol_or_container_ids[1].contents);
  bdestroy(test_string1);
  bdestroy(test_string2);
}

class AMFAppStatelessTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    itti_init(TASK_MAX, THREAD_MAX, MESSAGES_ID_MAX, tasks_info, messages_info,
              NULL, NULL);

    // initialize amf config
    amf_config_init(&amf_config);
    amf_config.use_stateless = true;
    amf_nas_state_init(&amf_config);
    create_state_matrix();
    amf_config.guamfi.nb = 1;
    amf_config.guamfi.guamfi[0].plmn = {.mcc_digit2 = 2,
                                        .mcc_digit1 = 2,
                                        .mnc_digit3 = 6,
                                        .mcc_digit3 = 2,
                                        .mnc_digit2 = 5,
                                        .mnc_digit1 = 4};

    init_task_context(TASK_MAIN, nullptr, 0, NULL, &amf_app_task_zmq_ctx);

    amf_app_desc_p = get_amf_nas_state(false);
  }

  virtual void TearDown() {
    clear_amf_nas_state();
    clear_amf_config(&amf_config);
    destroy_task_context(&amf_app_task_zmq_ctx);
    itti_free_desc_threads();
    AMFClientServicer::getInstance().map_table_key_proto_str.clear();
    AMFClientServicer::getInstance().map_imsi_ue_proto_str.clear();
  }

  // This Function mocks AMF task stop.
  void pseudo_amf_stop() {
    clear_amf_nas_state();
    clear_amf_config(&amf_config);
    destroy_task_context(&amf_app_task_zmq_ctx);
    itti_free_desc_threads();
  }

  amf_app_desc_t* amf_app_desc_p;
  std::string imsi = "222456000000001";
  plmn_t plmn = {.mcc_digit2 = 0,
                 .mcc_digit1 = 0,
                 .mnc_digit3 = 0x0f,
                 .mcc_digit3 = 1,
                 .mnc_digit2 = 1,
                 .mnc_digit1 = 0};

  const uint8_t initial_ue_message_hexbuf[25] = {
      0x7e, 0x00, 0x41, 0x79, 0x00, 0x0d, 0x01, 0x22, 0x62,
      0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
      0x01, 0x2e, 0x04, 0xf0, 0xf0, 0xf0, 0xf0};

  const uint8_t ue_auth_response_hexbuf[21] = {
      0x7e, 0x0,  0x57, 0x2d, 0x10, 0x25, 0x70, 0x6f, 0x9a, 0x5b, 0x90,
      0xb6, 0xc9, 0x57, 0x50, 0x6c, 0x88, 0x3d, 0x76, 0xcc, 0x63};

  const uint8_t ue_smc_response_hexbuf[60] = {
      0x7e, 0x4,  0x54, 0xf6, 0xe1, 0x2a, 0x0,  0x7e, 0x0,  0x5e, 0x77, 0x0,
      0x9,  0x45, 0x73, 0x80, 0x61, 0x21, 0x85, 0x61, 0x51, 0xf1, 0x71, 0x0,
      0x23, 0x7e, 0x0,  0x41, 0x79, 0x0,  0xd,  0x1,  0x22, 0x62, 0x54, 0x0,
      0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0x0,  0xf1, 0x10, 0x1,  0x0,  0x2e,
      0x4,  0xf0, 0xf0, 0xf0, 0xf0, 0x2f, 0x2,  0x1,  0x1,  0x53, 0x1,  0x0};

  const uint8_t ue_registration_complete_hexbuf[10] = {
      0x7e, 0x02, 0x5d, 0x5f, 0x49, 0x18, 0x01, 0x7e, 0x00, 0x43};

  const uint8_t ue_pdu_session_est_req_hexbuf[44] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x15, 0x2e, 0x01, 0x01, 0xc1, 0xff,
      0xff, 0x91, 0xa1, 0x28, 0x01, 0x00, 0x7b, 0x00, 0x07, 0x80, 0x00,
      0x0a, 0x00, 0x00, 0x0d, 0x00, 0x12, 0x01, 0x81, 0x22, 0x01, 0x01,
      0x25, 0x09, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74};

  const uint8_t pdu_sess_release_hexbuf[14] = {0x7e, 0x00, 0x67, 0x01, 0x00,
                                               0x06, 0x2e, 0x01, 0x01, 0xd1,
                                               0x59, 0x24, 0x12, 0x01};

  const uint8_t pdu_sess_release_complete_hexbuf[12] = {
      0x7e, 0x00, 0x67, 0x01, 0x00, 0x04, 0x2e, 0x01, 0x01, 0xd4, 0x12, 0x01};

  uint8_t ue_initiated_dereg_hexbuf[24] = {
      0x7e, 0x01, 0x41, 0x21, 0xe6, 0xe2, 0x03, 0x7e, 0x00, 0x45, 0x01, 0x00,
      0x0b, 0xf2, 0x22, 0x62, 0x54, 0x01, 0x00, 0x40, 0x0,  0x0,  0x0,  0x0};
};

// 1.Stateless triggered after Registration Request Complete
TEST_F(AMFAppStatelessTest, TestAfterRegistrationComplete) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));
  AMFClientServicer::getInstance().map_table_key_proto_str.clear();
  EXPECT_TRUE(
      AMFClientServicer::getInstance().map_table_key_proto_str.isEmpty());
  /* Writes the state to the data store */
  put_amf_nas_state();
  EXPECT_FALSE(
      AMFClientServicer::getInstance().map_table_key_proto_str.isEmpty());

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id((amf_ue_ngap_id_t)ue_id);
  ASSERT_NE(ue_context_p, nullptr);

  map_uint64_ue_context_t* amf_state_ue_id_ht =
      AmfNasStateManager::getInstance().get_ue_state_map();
  EXPECT_EQ(amf_state_ue_id_ht->size(), 1);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.size(), 1);

  put_amf_ue_state(amf_app_desc_p, imsi64, false);
  EXPECT_EQ(AMFClientServicer::getInstance().map_imsi_ue_proto_str.size(), 1);

  // Calling pseudo_amf_stop() and SetUp() simulates a service restart.
  AMFAppStatelessTest::pseudo_amf_stop();
  // Check if state data is cleared in AMF after pseudo_amf_stop()
  EXPECT_TRUE(AmfNasStateManager::getInstance().get_ue_state_map()->isEmpty());
  EXPECT_EQ(get_amf_nas_state(false), nullptr);

  // Internally reads back the state
  AMFAppStatelessTest::SetUp();
  amf_state_ue_id_ht = AmfNasStateManager::getInstance().get_ue_state_map();
  EXPECT_EQ(amf_state_ue_id_ht->size(), 1);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.size(), 1);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send ip address response from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_EQ(rc, RETURNok);

  /* Send pdu session setup response from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_EQ(rc, RETURNok);

  /* Send pdu resource setup response from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_EQ(rc, RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_EQ(rc, RETURNok);

  // Following cleanup is for ue_context after pseudo restart
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  // Following clean up is for stored ue_context before restart.
  delete ue_context_p;

  // TODO : CLEANUP_STATELESS do we need map_imsi_ue_proto_str
  // EXPECT_TRUE(AMFClientServicer::getInstance().map_imsi_ue_proto_str.isEmpty());
}

// 2.Stateless triggered after PDU Session Establishment Request
TEST_F(AMFAppStatelessTest, TestAfterPDUSessionEstReq) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;

  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));

  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  AMFClientServicer::getInstance().map_table_key_proto_str.clear();
  EXPECT_TRUE(
      AMFClientServicer::getInstance().map_table_key_proto_str.isEmpty());
  /* Writes the state to the data store */
  put_amf_nas_state();
  EXPECT_FALSE(
      AMFClientServicer::getInstance().map_table_key_proto_str.isEmpty());

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id((amf_ue_ngap_id_t)ue_id);
  ASSERT_NE(ue_context_p, nullptr);

  map_uint64_ue_context_t* amf_state_ue_id_ht =
      AmfNasStateManager::getInstance().get_ue_state_map();
  EXPECT_EQ(amf_state_ue_id_ht->size(), 1);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.size(), 1);

  EXPECT_EQ(ue_context_p->amf_context.smf_ctxt_map.size(), 1);
  put_amf_ue_state(amf_app_desc_p, imsi64, false);
  EXPECT_EQ(AMFClientServicer::getInstance().map_imsi_ue_proto_str.size(), 1);

  // Calling pseudo_amf_stop() and SetUp() simulates a service restart.
  AMFAppStatelessTest::pseudo_amf_stop();
  // Check if state data is cleared in AMF after pseudo_amf_stop()
  EXPECT_TRUE(AmfNasStateManager::getInstance().get_ue_state_map()->isEmpty());
  EXPECT_EQ(get_amf_nas_state(false), nullptr);

  // Internally reads back the state
  AMFAppStatelessTest::SetUp();
  amf_state_ue_id_ht = AmfNasStateManager::getInstance().get_ue_state_map();
  EXPECT_EQ(amf_state_ue_id_ht->size(), 1);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.size(), 1);

  ue_m5gmm_context_t* new_ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id((amf_ue_ngap_id_t)ue_id);
  EXPECT_EQ(new_ue_context_p->amf_context.smf_ctxt_map.size(), 1);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_EQ(rc, RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_EQ(rc, RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_EQ(rc, RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_EQ(rc, RETURNok);

  // Following cleanup is for ue_context after pseudo restart
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  // Following clean up is for stored ue_context before restart.
  delete ue_context_p;

  // TODO : CLEANUP_STATELESS do we need map_imsi_ue_proto_str
  // EXPECT_TRUE(AMFClientServicer::getInstance().map_imsi_ue_proto_str.isEmpty());
}

// 3.Stateless triggered after pdu session release complete
TEST_F(AMFAppStatelessTest, TestAfterPDUSessionReleaseComplete) {
  int rc = RETURNerror;
  amf_ue_ngap_id_t ue_id = 0;

  /* Send the initial UE message */
  imsi64_t imsi64 = 0;
  imsi64 = send_initial_ue_message_no_tmsi(amf_app_desc_p, 36, 1, 1, 0, plmn,
                                           initial_ue_message_hexbuf,
                                           sizeof(initial_ue_message_hexbuf));
  AMFClientServicer::getInstance().map_table_key_proto_str.clear();
  EXPECT_TRUE(
      AMFClientServicer::getInstance().map_table_key_proto_str.isEmpty());
  /* Writes the state to the data store */
  put_amf_nas_state();
  EXPECT_FALSE(
      AMFClientServicer::getInstance().map_table_key_proto_str.isEmpty());

  /* Check if UE Context is created with correct imsi */
  EXPECT_TRUE(get_ue_id_from_imsi(amf_app_desc_p, imsi64, &ue_id));

  /* Send the authentication response message from subscriberdb */
  rc = send_proc_authentication_info_answer(imsi, ue_id, true);
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for auth response from UE */
  rc = send_uplink_nas_message_ue_auth_response(
      amf_app_desc_p, ue_id, plmn, ue_auth_response_hexbuf,
      sizeof(ue_auth_response_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for security mode complete response from UE */
  rc = send_uplink_nas_message_ue_smc_response(amf_app_desc_p, ue_id, plmn,
                                               ue_smc_response_hexbuf,
                                               sizeof(ue_smc_response_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  s6a_update_location_ans_t ula_ans = util_amf_send_s6a_ula(imsi);
  rc = amf_handle_s6a_update_location_ans(&ula_ans);
  EXPECT_EQ(rc, RETURNok);

  send_initial_context_response(amf_app_desc_p, ue_id);

  /* Send uplink nas message for registration complete response from UE */
  rc = send_uplink_nas_registration_complete(
      amf_app_desc_p, ue_id, plmn, ue_registration_complete_hexbuf,
      sizeof(ue_registration_complete_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session establishment request from UE */
  rc = send_uplink_nas_pdu_session_establishment_request(
      amf_app_desc_p, ue_id, plmn, ue_pdu_session_est_req_hexbuf,
      sizeof(ue_pdu_session_est_req_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send ip address response  from pipelined */
  rc = send_ip_address_response_itti(IPv4);
  EXPECT_EQ(rc, RETURNok);

  /* Send pdu session setup response  from smf */
  rc = send_pdu_session_response_itti(IPv4);
  EXPECT_EQ(rc, RETURNok);

  /* Send pdu resource setup response  from UE */
  rc = send_pdu_resource_setup_response(ue_id);
  EXPECT_EQ(rc, RETURNok);

  rc = send_pdu_notification_response();
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session release request from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_hexbuf,
      sizeof(pdu_sess_release_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for pdu session release complete from UE */
  rc = send_uplink_nas_pdu_session_release_message(
      amf_app_desc_p, ue_id, plmn, pdu_sess_release_complete_hexbuf,
      sizeof(pdu_sess_release_complete_hexbuf));
  EXPECT_EQ(rc, RETURNok);

  ue_m5gmm_context_t* ue_context_p =
      amf_ue_context_exists_amf_ue_ngap_id((amf_ue_ngap_id_t)ue_id);
  ASSERT_NE(ue_context_p, nullptr);

  map_uint64_ue_context_t* amf_state_ue_id_ht =
      AmfNasStateManager::getInstance().get_ue_state_map();
  EXPECT_EQ(amf_state_ue_id_ht->size(), 1);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.size(), 1);

  put_amf_ue_state(amf_app_desc_p, imsi64, false);
  EXPECT_EQ(AMFClientServicer::getInstance().map_imsi_ue_proto_str.size(), 1);

  // Calling pseudo_amf_stop() and SetUp() simulates a service restart.
  AMFAppStatelessTest::pseudo_amf_stop();
  // Check if state data is cleared in AMF after pseudo_amf_stop()
  EXPECT_TRUE(AmfNasStateManager::getInstance().get_ue_state_map()->isEmpty());
  EXPECT_EQ(get_amf_nas_state(false), nullptr);

  // Internally reads back the state
  AMFAppStatelessTest::SetUp();
  amf_state_ue_id_ht = AmfNasStateManager::getInstance().get_ue_state_map();
  EXPECT_EQ(amf_state_ue_id_ht->size(), 1);
  EXPECT_EQ(
      amf_app_desc_p->amf_ue_contexts.gnb_ue_ngap_id_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.imsi_amf_ue_id_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.tun11_ue_context_htbl.size(), 1);
  EXPECT_EQ(amf_app_desc_p->amf_ue_contexts.guti_ue_context_htbl.size(), 1);

  rc = send_pdu_notification_response();
  EXPECT_EQ(rc, RETURNok);

  /* Send uplink nas message for deregistration complete response from UE */
  rc = send_uplink_nas_ue_deregistration_request(
      amf_app_desc_p, ue_id, plmn, ue_initiated_dereg_hexbuf,
      sizeof(ue_initiated_dereg_hexbuf));

  EXPECT_EQ(rc, RETURNok);

  // Following cleanup is for ue_context after pseudo restart
  send_ue_context_release_complete_message(amf_app_desc_p, 1, 1, ue_id);

  // Following clean up is for stored ue_context before restart.
  delete ue_context_p;

  // TODO : CLEANUP_STATELESS do we need map_imsi_ue_proto_str
  // EXPECT_TRUE(AMFClientServicer::getInstance().map_imsi_ue_proto_str.isEmpty());
}
}  // namespace magma5g
