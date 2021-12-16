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

class SessionStateTest : public ::testing::Test {
 protected:
  virtual void SetUp() {
    Teids teids;
    cfg.common_context =
        build_common_context(IMSI1, IP1, IPv6_1, teids, APN1, MSISDN, TGPP_LTE);
    auto tgpp_ctx = TgppContext();
    auto pdp_start_time = 12345;
    create_tgpp_context("gx.dest.com", "gy.dest.com", &tgpp_ctx);
    rule_store = std::make_shared<StaticRuleStore>();
    session_state = std::make_shared<SessionState>(SESSION_ID_1, cfg,
                                                   *rule_store, pdp_start_time);
    session_state->set_tgpp_context(tgpp_ctx, nullptr);
    session_state->set_fsm_state(SESSION_ACTIVE, nullptr);
    session_state->set_create_session_response(CreateSessionResponse(),
                                               nullptr);
    update_criteria = get_default_update_criteria();
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

  uint32_t insert_rule(uint32_t rating_group, const std::string& m_key,
                       const std::string& rule_id, PolicyType rule_type,
                       std::time_t activation_time,
                       std::time_t deactivation_time) {
    PolicyRule rule = create_policy_rule(rule_id, m_key, rating_group);
    RuleLifetime lifetime(activation_time, deactivation_time);
    switch (rule_type) {
      case STATIC:
        // insert into list of existing rules
        rule_store->insert_rule(rule);
        // mark the rule as active in session
        return session_state
            ->activate_static_rule(rule_id, lifetime, &update_criteria)
            .version;
      case DYNAMIC:
        return session_state
            ->insert_dynamic_rule(rule, lifetime, &update_criteria)
            .version;
        break;
    }
    return 0;
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

  void initialize_session_with_qos() {
    SessionConfig cfg;
    Teids teids;
    cfg.common_context =
        build_common_context(IMSI1, IP1, IPv6_1, teids, APN1, MSISDN, TGPP_LTE);
    QosInformationRequest qos_info;
    qos_info.set_apn_ambr_dl(32);
    qos_info.set_apn_ambr_dl(64);
    const auto& lte_context =
        build_lte_context(IP2, "", "", "", "", BEARER_ID_1, &qos_info);
    cfg.rat_specific_context.mutable_lte_context()->CopyFrom(lte_context);
    session_state->set_config(cfg, nullptr);
  }

  // TODO: make session_manager.proto and policydb.proto to use common field
  static RedirectInformation_AddressType address_type_converter(
      RedirectServer_RedirectAddressType address_type) {
    switch (address_type) {
      case RedirectServer_RedirectAddressType_IPV4:
        return RedirectInformation_AddressType_IPv4;
      case RedirectServer_RedirectAddressType_IPV6:
        return RedirectInformation_AddressType_IPv6;
      case RedirectServer_RedirectAddressType_URL:
        return RedirectInformation_AddressType_URL;
      case RedirectServer_RedirectAddressType_SIP_URI:
        return RedirectInformation_AddressType_SIP_URI;
      default:
        return RedirectInformation_AddressType_IPv4;
    }
  }

  uint32_t insert_gy_redirection_rule(const std::string& rule_id) {
    PolicyRule redirect_rule;
    redirect_rule.set_id(rule_id);
    redirect_rule.set_priority(999);

    RedirectInformation* redirect_info = redirect_rule.mutable_redirect();
    redirect_info->set_support(RedirectInformation_Support_ENABLED);

    RedirectServer redirect_server;
    redirect_server.set_redirect_address_type(RedirectServer::URL);
    redirect_server.set_redirect_server_address("http://www.example.com/");

    redirect_info->set_address_type(
        address_type_converter(redirect_server.redirect_address_type()));
    redirect_info->set_server_address(
        redirect_server.redirect_server_address());

    RuleLifetime lifetime{};
    return session_state
        ->insert_gy_rule(redirect_rule, lifetime, &update_criteria)
        .version;
  }

  void receive_credit_from_ocs(uint32_t rating_group, uint64_t volume) {
    CreditUpdateResponse charge_resp;
    create_credit_update_response(IMSI1, SESSION_ID_1, rating_group, volume,
                                  &charge_resp);
    session_state->receive_charging_credit(charge_resp, &update_criteria);
  }

  void receive_credit_from_ocs(uint32_t rating_group, uint64_t total_volume,
                               uint64_t tx_volume, uint64_t rx_volume,
                               bool is_final) {
    CreditUpdateResponse charge_resp;
    create_credit_update_response(IMSI1, SESSION_ID_1, rating_group,
                                  total_volume, tx_volume, rx_volume, is_final,
                                  &charge_resp);
    session_state->receive_charging_credit(charge_resp, &update_criteria);
  }

  void receive_credit_from_pcrf(const std::string& mkey, uint64_t volume,
                                MonitoringLevel level) {
    UsageMonitoringUpdateResponse monitor_resp;
    receive_credit_from_pcrf(mkey, volume, 0, 0, level);
  }

  void receive_credit_from_pcrf(const std::string& mkey, uint64_t total_volume,
                                uint64_t tx_volume, uint64_t rx_volume,
                                MonitoringLevel level) {
    UsageMonitoringUpdateResponse monitor_resp;
    create_monitor_update_response(IMSI1, SESSION_ID_1, mkey, level,
                                   total_volume, tx_volume, rx_volume,
                                   &monitor_resp);
    session_state->receive_monitor(monitor_resp, &update_criteria);
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

 protected:
  std::shared_ptr<StaticRuleStore> rule_store;
  std::shared_ptr<SessionState> session_state;
  SessionStateUpdateCriteria update_criteria;
  SessionConfig cfg;
};
};  // namespace magma