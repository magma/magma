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
#include <lte/protos/session_manager.grpc.pb.h>
#include <memory>
#include "Consts.h"
#include "LocalEnforcer.h"
#include "ProtobufCreators.h"
#include "SessiondMocks.h"
#include "SessionStore.h"

using grpc::ServerContext;
using grpc::Status;
using ::testing::InSequence;
using ::testing::Test;

namespace magma {

Teids teids0;
Teids teids1;
Teids teids2;

class LocalEnforcerStatsPollerTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    reporter = std::make_shared<MockSessionReporter>();
    pipelined_mock = std::make_shared<MockPipelined>();
    rule_store = std::make_shared<StaticRuleStore>();
    session_store = std::make_shared<SessionStore>(
        rule_store, std::make_shared<MeteringReporter>());
    pipelined_client = std::make_shared<MockPipelinedClient>();
    spgw_client = std::make_shared<MockSpgwServiceClient>();
    aaa_client = std::make_shared<MockAAAClient>();
    events_reporter = std::make_shared<MockEventsReporter>();
    auto default_mconfig = get_default_mconfig();
    auto shard_tracker = std::make_shared<ShardTracker>();
    local_enforcer = std::make_unique<LocalEnforcer>(
        reporter, rule_store, *session_store, pipelined_client, events_reporter,
        spgw_client, aaa_client, shard_tracker, 0, 0, default_mconfig);
    session_map = SessionMap{};
    test_cfg_ = get_default_config("");

    teids0.set_agw_teid(0);
    teids0.set_enb_teid(0);
    teids1.set_agw_teid(TEID_1_UL);
    teids1.set_enb_teid(TEID_1_DL);
    teids2.set_agw_teid(TEID_2_UL);
    teids2.set_enb_teid(TEID_2_DL);
  }

  SessionConfig get_default_config(const std::string& imsi) {
    SessionConfig cfg;
    cfg.common_context =
        build_common_context(imsi, IP1, IPv6_1, teids1, APN1, MSISDN, TGPP_LTE);
    QosInformationRequest qos_info;
    qos_info.set_apn_ambr_dl(32);
    qos_info.set_apn_ambr_dl(64);
    qos_info.set_br_unit(QosInformationRequest_BitrateUnitsAMBR_KBPS);
    const auto& lte_context =
        build_lte_context(IP2, "", "", "", "", BEARER_ID_1, &qos_info);
    cfg.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);
    return cfg;
  }

  void insert_static_rule(uint32_t rating_group, const std::string& m_key,
                          const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, m_key, rating_group));
  }

  void initialize_session(SessionMap& session_map,
                          const std::string& session_id,
                          const SessionConfig& cfg,
                          const CreateSessionResponse& response) {
    const std::string imsi = cfg.get_imsi();
    auto session = local_enforcer->create_initializing_session(session_id, cfg);
    local_enforcer->update_session_with_policy_response(session, response,
                                                        nullptr);
    session_map[imsi].push_back(std::move(session));
  }

 protected:
  std::shared_ptr<MockSessionReporter> reporter;
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionStore> session_store;
  std::unique_ptr<LocalEnforcer> local_enforcer;
  std::shared_ptr<MockPipelinedClient> pipelined_client;
  std::shared_ptr<MockSpgwServiceClient> spgw_client;
  std::shared_ptr<MockAAAClient> aaa_client;
  std::shared_ptr<MockEventsReporter> events_reporter;
  std::shared_ptr<MockPipelined> pipelined_mock;
  SessionMap session_map;
  SessionConfig test_cfg_;
};

TEST_F(LocalEnforcerStatsPollerTest, test_poll_stats) {
  // insert some rules to retrieve
  insert_static_rule(1, "", "rule1");
  insert_static_rule(1, "", "rule2");
  insert_static_rule(1, "", "rule3");
  insert_static_rule(1, "", "rule4");

  CreateSessionResponse response;
  initialize_session(session_map, SESSION_ID_1, get_default_config(IMSI1),
                     response);
  local_enforcer->update_tunnel_ids(
      session_map,
      create_update_tunnel_ids_request(IMSI1, BEARER_ID_1, teids1));
  RuleRecordTable table;

  auto record_list = table.mutable_records();
  auto ue_ipv4 = test_cfg_.common_context.ue_ipv4();
  create_rule_record(IMSI1, ue_ipv4, "rule1", 10, 20, record_list->Add());
  create_rule_record(IMSI1, ue_ipv4, "rule2", 15, 35, record_list->Add());
  create_rule_record(IMSI1, ue_ipv4, "rule3", 100, 150, record_list->Add());
  create_rule_record(IMSI1, ue_ipv4, "rule4", 200, 300, record_list->Add());

  auto update = SessionStore::get_default_session_update(session_map);

  local_enforcer->aggregate_records(session_map, table, update);

  int cookie = 0;
  int cookie_mask = 0;
  EXPECT_CALL(*pipelined_client, poll_stats(cookie, cookie_mask, testing::_))
      .Times(1);
  local_enforcer->poll_stats_enforcer(cookie, cookie_mask);

  cookie = 1;
  cookie_mask = 0;
  EXPECT_CALL(*pipelined_client, poll_stats(cookie, cookie_mask, testing::_))
      .Times(1);
  local_enforcer->poll_stats_enforcer(cookie, cookie_mask);

  cookie = 0;
  cookie_mask = 1;
  EXPECT_CALL(*pipelined_client, poll_stats(cookie, cookie_mask, testing::_))
      .Times(1);
  local_enforcer->poll_stats_enforcer(cookie, cookie_mask);
}

int main(int argc, char** argv) {
  ::testing::InitGoogleTest(&argc, argv);
  FLAGS_logtostderr = 1;
  FLAGS_v = 10;
  return RUN_ALL_TESTS();
}

}  // namespace magma
