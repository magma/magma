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
#include <glog/logging.h>
#include "lte/gateway/c/core/oai/include/mme_config.hpp"

#include "lte/protos/session_manager.pb.h"
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.hpp"
#include "lte/gateway/c/core/oai/tasks/grpc_service/AmfServiceImpl.hpp"

using ::testing::Test;

struct mme_config_s mme_config;
task_zmq_ctx_t grpc_service_task_zmq_ctx;

namespace magma {
namespace lte {

TEST(test_create_sm_pdu_session_v4, create_sm_pdu_session_v4) {
  SetSMSessionContext request;

  std::string imsi("901700000000001");
  std::string apn("magmacore.com");
  uint32_t pdu_session_id = 0x5;
  uint32_t pdu_session_type = 3;
  uint32_t gnb_gtp_teid = 1;
  uint8_t pti = 10;
  uint8_t gnb_gtp_teid_ip_addr[4] = {0};  //("10.20.30.40")
  gnb_gtp_teid_ip_addr[0] = 0xA;
  gnb_gtp_teid_ip_addr[1] = 0x14;
  gnb_gtp_teid_ip_addr[2] = 0x1E;
  gnb_gtp_teid_ip_addr[3] = 0x28;
  eps_subscribed_qos_profile_t qos_profile = {0};
  std::string gnb_ip_addr;
  for (int i = 0; i < 4; ++i) {
    gnb_ip_addr += std::to_string(gnb_gtp_teid_ip_addr[i]);
    if (i != 3) {
      gnb_ip_addr += ".";
    }
  }

  std::string ue_ipv4_addr("10.20.30.44");
  std::string ue_ipv6_addr;

  uint32_t version = 0;

  ambr_t default_ambr;

  // qos profile
  qos_profile.qci = 5;
  qos_profile.allocation_retention_priority.priority_level = 15;

  request = magma5g::create_sm_pdu_session(
      imsi, (uint8_t*)apn.c_str(), pdu_session_id, pdu_session_type,
      gnb_gtp_teid, pti, gnb_gtp_teid_ip_addr, ue_ipv4_addr, ue_ipv6_addr,
      default_ambr, version, qos_profile);

  auto* rat_req =
      request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
  auto* req_cmn = request.mutable_common_context();

  EXPECT_EQ(imsi, req_cmn->sid().id().substr(4));
  EXPECT_EQ(magma::lte::SubscriberID_IDType::SubscriberID_IDType_IMSI,
            req_cmn->sid().type());
  EXPECT_EQ(apn, req_cmn->apn());
  EXPECT_EQ(magma::lte::RATType::TGPP_NR, req_cmn->rat_type());
  EXPECT_EQ(magma::lte::SMSessionFSMState::CREATING_0,
            req_cmn->sm_session_state());
  EXPECT_EQ(0, req_cmn->sm_session_version());
  EXPECT_EQ(pdu_session_id, rat_req->pdu_session_id());
  EXPECT_EQ(magma::lte::RequestType::INITIAL_REQUEST, rat_req->request_type());

  EXPECT_EQ(magma::lte::PduSessionType::IPV4, rat_req->pdu_session_type());
  EXPECT_EQ(1, rat_req->mutable_gnode_endpoint()->teid());
  EXPECT_EQ(gnb_ip_addr, rat_req->mutable_gnode_endpoint()->end_ipv4_addr());
  uint8_t pti_decoded = (uint8_t)rat_req->procedure_trans_identity();
  EXPECT_EQ(pti, pti_decoded);
  EXPECT_EQ(ue_ipv4_addr, req_cmn->ue_ipv4());
}

TEST(test_create_sm_pdu_session_response, create_sm_pdu_session_response) {
  SetSMSessionContextAccess request;
  AmfServiceImpl amfservice;

  itti_n11_create_pdu_session_response_t itti_msg;

  auto* req_common = request.mutable_common_context();
  auto* req_m5g =
      request.mutable_rat_specific_context()->mutable_m5g_session_context_rsp();

  req_common->mutable_sid()->set_id("IMSI901700000000001");
  req_common->set_ue_ipv4("192.168.128.12");
  req_common->set_apn("internet");

  req_m5g->set_pdu_session_id(1);
  req_m5g->set_pdu_session_type(magma::lte::PduSessionType::IPV4);
  req_m5g->set_selected_ssc_mode(magma::lte::SscMode::SSC_MODE_1);
  req_m5g->set_allowed_ssc_mode(magma::lte::SscMode::SSC_MODE_2);
  req_m5g->mutable_subscribed_qos()->set_apn_ambr_dl(4000);
  req_m5g->mutable_subscribed_qos()->set_apn_ambr_ul(3000);
  req_m5g->mutable_subscribed_qos()->set_priority_level(1);
  req_m5g->mutable_subscribed_qos()->set_qos_class_id(magma::lte::QCI::QCI_9);
  req_m5g->mutable_subscribed_qos()->set_preemption_capability(
      magma::lte::prem_capab::SHALL_NOT_TRIGGER_PRE_EMPTION);
  req_m5g->mutable_subscribed_qos()->set_preemption_vulnerability(
      magma::lte::prem_vuner::PRE_EMPTABLE);
  req_m5g->set_m5gsm_cause(magma::lte::M5GSMCause::OPERATION_SUCCESS);
  req_m5g->set_m5g_sm_congestion_reattempt_indicator(true);
  req_m5g->set_procedure_trans_identity(1);
  req_m5g->mutable_upf_endpoint()->set_teid(2147483647);
  req_m5g->mutable_upf_endpoint()->set_end_ipv4_addr("192.168.60.142");
  req_m5g->set_always_on_pdu_session_indication(true);

  amfservice.SetSmfSessionContext_itti(&request, &itti_msg);

  EXPECT_EQ(itti_msg.pdu_session_id, req_m5g->pdu_session_id());
  EXPECT_EQ(itti_msg.pdu_session_type, req_m5g->pdu_session_type());
  EXPECT_EQ(itti_msg.selected_ssc_mode, req_m5g->selected_ssc_mode());
  EXPECT_EQ(itti_msg.allowed_ssc_mode, req_m5g->allowed_ssc_mode());
  EXPECT_EQ(itti_msg.m5gsm_cause, req_m5g->m5gsm_cause());
  EXPECT_EQ(itti_msg.m5gsm_congetion_re_attempt_indicator,
            req_m5g->m5g_sm_congestion_reattempt_indicator());
  EXPECT_EQ(itti_msg.always_on_pdu_session_indication,
            req_m5g->always_on_pdu_session_indication());
  EXPECT_EQ(itti_msg.procedure_trans_identity,
            req_m5g->procedure_trans_identity());
  EXPECT_EQ(itti_msg.session_ambr.uplink_units,
            req_m5g->mutable_subscribed_qos()->apn_ambr_ul());
  EXPECT_EQ(itti_msg.session_ambr.downlink_units,
            req_m5g->mutable_subscribed_qos()->apn_ambr_dl());
  EXPECT_EQ(
      itti_msg.qos_flow_list.item[0].qos_flow_req_item.qos_flow_identifier,
      req_m5g->mutable_subscribed_qos()->qos_class_id());
  EXPECT_EQ(itti_msg.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_descriptor.qos_flow_identifier,
            req_m5g->mutable_subscribed_qos()->qos_class_id());
  EXPECT_EQ(itti_msg.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_descriptor.fiveQi,
            req_m5g->mutable_subscribed_qos()->qos_class_id());
  EXPECT_EQ(itti_msg.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
                .non_dynamic_5QI_desc.fiveQI,
            req_m5g->mutable_subscribed_qos()->qos_class_id());
  EXPECT_EQ(itti_msg.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .priority_level,
            req_m5g->mutable_subscribed_qos()->priority_level());
  EXPECT_EQ(itti_msg.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_cap,
            req_m5g->mutable_subscribed_qos()->preemption_capability());
  EXPECT_EQ(itti_msg.qos_flow_list.item[0]
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_vul,
            req_m5g->mutable_subscribed_qos()->preemption_vulnerability());
}

TEST(test_received_notification_t, create_received_notification_t) {
  SetSmNotificationContext notif;
  AmfServiceImpl amfservice;

  itti_n11_received_notification_t itti_msg;
  auto* notify_common = notif.mutable_common_context();
  auto* req_m5g = notif.mutable_rat_specific_notification();

  // CommonSessionContext
  notify_common->mutable_sid()->set_id("IMSI901700000000001");
  notify_common->set_sm_session_state(
      magma::lte::SMSessionFSMState::CREATING_0);
  notify_common->set_sm_session_version(1);

  // RatSpecificContextAccess
  req_m5g->set_pdu_session_id(1);
  req_m5g->set_request_type(magma::lte::RequestType::INITIAL_REQUEST);
  req_m5g->set_pdu_session_type(magma::lte::PduSessionType::IPV4);
  req_m5g->mutable_m5g_sm_capability()->set_reflective_qos(1);
  req_m5g->mutable_m5g_sm_capability()->set_multi_homed_ipv6_pdu_session(2);
  req_m5g->set_m5gsm_cause(magma::lte::M5GSMCause::OPERATION_SUCCESS);
  req_m5g->set_notify_ue_event(PDU_SESSION_STATE_NOTIFY);

  amfservice.SetAmfNotification_itti(&notif, &itti_msg);

  EXPECT_EQ(itti_msg.pdu_session_id, req_m5g->pdu_session_id());
  EXPECT_EQ(itti_msg.pdu_session_type, req_m5g->pdu_session_type());
  EXPECT_EQ(itti_msg.sm_session_version, notify_common->sm_session_version());
  EXPECT_EQ(itti_msg.sm_session_fsm_state, notify_common->sm_session_state());
  EXPECT_EQ(itti_msg.request_type, req_m5g->request_type());
  EXPECT_EQ(itti_msg.m5g_sm_capability.reflective_qos,
            req_m5g->mutable_m5g_sm_capability()->reflective_qos());
  EXPECT_EQ(
      itti_msg.m5g_sm_capability.multi_homed_ipv6_pdu_session,
      req_m5g->mutable_m5g_sm_capability()->multi_homed_ipv6_pdu_session());
  EXPECT_EQ(itti_msg.m5gsm_cause, req_m5g->m5gsm_cause());
  EXPECT_EQ(itti_msg.notify_ue_evnt, req_m5g->notify_ue_event());
}
}  // namespace lte
}  // namespace magma
