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
#include <glog/logging.h>
#include <gtest/gtest.h>

#include <future>
#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "Consts.h"
#include "ProtobufCreators.h"
#include "SessiondMocks.h"
#include "SessionState.h"

using ::testing::Test;

namespace magma {

class SessionStateTest5G : public ::testing::Test {
 protected:
  virtual void SetUp() {
    Teids teids;
    rule_store = std::make_shared<StaticRuleStore>();
    cfg = build_sm_context(IMSI1, "10.20.30.40", 5);
    session_state =
        std::make_shared<SessionState>(IMSI1, SESSION_ID_1, cfg, *rule_store);
    session_state->set_fsm_state(SESSION_ACTIVE, nullptr);
    session_state->set_create_session_response(CreateSessionResponse(),
                                               nullptr);
    update_criteria = get_default_update_criteria();
  }

  SessionConfig build_sm_context(
      const std::string& imsi,  // assumes IMSI prefix
      const std::string& ue_ipv4, uint32_t pdu_id) {
    SetSMSessionContext request;
    auto* req =
        request.mutable_rat_specific_context()->mutable_m5gsm_session_context();
    auto* reqcmn = request.mutable_common_context();
    req->set_pdu_session_id(pdu_id);
    req->set_request_type(magma::RequestType::INITIAL_REQUEST);
    req->mutable_pdu_address()->set_redirect_address_type(
        magma::RedirectServer::IPV4);
    req->set_access_type(magma::AccessType::M_3GPP_ACCESS_3GPP);
    req->mutable_pdu_address()->set_redirect_server_address(ue_ipv4);
    req->set_priority_access(magma::priorityaccess::High);
    req->set_imei("123456789012345");
    req->set_gpsi("9876543210");
    req->set_pcf_id("1357924680123456");

    reqcmn->mutable_sid()->set_id(imsi);
    reqcmn->set_apn("BLR");
    reqcmn->set_ue_ipv4("192.168.128.11");
    reqcmn->set_rat_type(magma::RATType::TGPP_NR);
    reqcmn->set_sm_session_state(magma::SMSessionFSMState::CREATING_0);

    SessionConfig cfg;
    cfg.common_context = request.common_context();
    cfg.rat_specific_context = request.rat_specific_context();
    cfg.rat_specific_context.mutable_m5gsm_session_context()->set_ssc_mode(
        SSC_MODE_3);
    return cfg;
  }

  void insert_static_rule_into_store(uint32_t rating_group,
                                     const std::string& m_key,
                                     const std::string& rule_id) {
    rule_store->insert_rule(create_policy_rule(rule_id, m_key, rating_group));
  }

  void insert_static_rule_with_qos_into_store(uint32_t rating_group,
                                              const std::string& m_key,
                                              const int qci,
                                              const std::string& rule_id) {
    PolicyRule rule =
        create_policy_rule_with_qos(rule_id, m_key, rating_group, qci);
    rule_store->insert_rule(rule);
  }

  uint32_t activate_rule(uint32_t rating_group, const std::string& m_key,
                         const std::string& rule_id, PolicyType rule_type,
                         std::time_t activation_time,
                         std::time_t deactivation_time) {
    PolicyRule rule = create_policy_rule(rule_id, m_key, rating_group);
    RuleLifetime lifetime(activation_time, deactivation_time);
    switch (rule_type) {
      case STATIC:
        rule_store->insert_rule(rule);
        return session_state
            ->activate_static_rule(rule_id, lifetime, &update_criteria)
            .version;
        break;
      case DYNAMIC:
        return session_state
            ->insert_dynamic_rule(rule, lifetime, &update_criteria)
            .version;
        break;
      default:
        break;
    }
    return 0;
  }

  uint32_t get_monitored_rule_count(const std::string& mkey) {
    std::vector<PolicyRule> rules;
    EXPECT_TRUE(session_state->get_dynamic_rules().get_rules(rules));
    uint32_t count = 0;
    for (PolicyRule& rule : rules) {
      if (rule.monitoring_key() == mkey) {
        count++;
      }
    }
    return count;
  }

  void schedule_rule(uint32_t rating_group, const std::string& m_key,
                     const std::string& rule_id, PolicyType rule_type,
                     std::time_t activation_time,
                     std::time_t deactivation_time) {
    PolicyRule rule = create_policy_rule(rule_id, m_key, rating_group);
    RuleLifetime lifetime(activation_time, deactivation_time);
    switch (rule_type) {
      case STATIC:
        // insert into list of existing rules
        rule_store->insert_rule(rule);
        // mark the rule as scheduled in the session
        session_state->schedule_static_rule(rule_id, lifetime,
                                            &update_criteria);
        break;
      case DYNAMIC:
        session_state->schedule_dynamic_rule(rule, lifetime, &update_criteria);
        break;
    }
  }

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionState> session_state;
  SessionStateUpdateCriteria update_criteria;
  SessionConfig cfg;
};

};  // namespace magma
