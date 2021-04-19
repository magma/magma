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

#include <vector>
#include <string>

#include "Consts.h"
#include "ProtobufCreators.h"

namespace magma {

CommonSessionContext build_common_context(
    const std::string& imsi,  // assumes IMSI prefix
    const std::string& ue_ipv4, const std::string& ue_ipv6, const Teids teids,
    const std::string& apn, const std::string& msisdn, const RATType rat_type) {
  CommonSessionContext common_context;
  common_context.mutable_sid()->set_id(imsi);
  common_context.set_ue_ipv4(ue_ipv4);
  common_context.set_ue_ipv6(ue_ipv6);
  common_context.mutable_teids()->CopyFrom(teids);
  common_context.set_apn(apn);
  common_context.set_msisdn(msisdn);
  common_context.set_rat_type(rat_type);
  return common_context;
}

LTESessionContext build_lte_context(
    const std::string& spgw_ipv4, const std::string& imei,
    const std::string& plmn_id, const std::string& imsi_plmn_id,
    const std::string& user_location, uint32_t bearer_id,
    QosInformationRequest* qos_info) {
  LTESessionContext lte_context;
  lte_context.set_spgw_ipv4(spgw_ipv4);
  lte_context.set_imei(imei);
  lte_context.set_plmn_id(plmn_id);
  lte_context.set_imsi_plmn_id(imsi_plmn_id);
  lte_context.set_user_location(user_location);
  lte_context.set_bearer_id(bearer_id);
  if (qos_info != nullptr) {
    lte_context.mutable_qos_info()->CopyFrom(*qos_info);
  }
  return lte_context;
}

WLANSessionContext build_wlan_context(
    const std::string& mac_addr, const std::string& radius_session_id) {
  WLANSessionContext wlan_context;
  wlan_context.set_mac_addr(mac_addr);
  wlan_context.set_radius_session_id(radius_session_id);
  return wlan_context;
}

RuleSet create_rule_set(
    const bool apply_subscriber_wide, const std::string& apn,
    std::vector<std::string> static_rules,
    std::vector<PolicyRule> dynamic_rules) {
  RuleSet rule_set;
  rule_set.set_apply_subscriber_wide(apply_subscriber_wide);
  rule_set.set_apn(apn);
  for (const auto& rule : static_rules) {
    rule_set.mutable_static_rules()->Add()->set_rule_id(rule);
  }
  for (const auto& rule : dynamic_rules) {
    rule_set.mutable_dynamic_rules()->Add()->mutable_policy_rule()->CopyFrom(
        rule);
  }
  return rule_set;
}

void create_rule_record(
    const std::string& imsi, const std::string& rule_id, uint64_t bytes_rx,
    uint64_t bytes_tx, RuleRecord* rule_record) {
  rule_record->set_sid(imsi);
  rule_record->set_rule_id(rule_id);
  rule_record->set_bytes_rx(bytes_rx);
  rule_record->set_bytes_tx(bytes_tx);
  rule_record->set_dropped_rx(0);
  rule_record->set_dropped_tx(0);
}

void create_rule_record(
    const std::string& imsi, const uint32_t teid, const std::string& rule_id,
    uint64_t bytes_rx, uint64_t bytes_tx, RuleRecord* rule_record) {
  create_rule_record(imsi, rule_id, bytes_rx, bytes_tx, rule_record);
  rule_record->set_teid(teid);
}

void create_rule_record(
    const std::string& imsi, const uint32_t teid, const std::string& rule_id,
    uint64_t bytes_rx, uint64_t bytes_tx, uint64_t dropped_rx,
    uint64_t dropped_tx, RuleRecord* rule_record) {
  create_rule_record(imsi, teid, rule_id, bytes_rx, bytes_tx, rule_record);
  rule_record->set_dropped_rx(dropped_rx);
  rule_record->set_dropped_tx(dropped_tx);
}

void create_charging_credit(
    uint64_t volume, bool is_final, ChargingCredit* credit) {
  create_granted_units(&volume, NULL, NULL, credit->mutable_granted_units());
  credit->set_type(ChargingCredit::BYTES);
  credit->set_is_final(is_final);
}

void create_charging_credit(
    uint64_t volume, ChargingCredit_FinalAction action,
    std::string redirect_server, std::string restrict_rule,
    ChargingCredit* credit) {
  create_granted_units(&volume, NULL, NULL, credit->mutable_granted_units());
  credit->set_type(ChargingCredit::BYTES);
  credit->set_is_final(true);
  credit->set_final_action(action);
  credit->mutable_redirect_server()->set_redirect_server_address(
      redirect_server);
  credit->add_restrict_rules(restrict_rule);
}

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, CreditLimitType limit_type,
    CreditUpdateResponse* response) {
  response->set_success(true);
  response->set_sid(imsi);
  response->set_session_id(session_id);
  response->set_charging_key(charging_key);
  response->set_limit_type(limit_type);
}

// defaults to not final credit
void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t volume, CreditUpdateResponse* response) {
  create_credit_update_response(
      imsi, session_id, charging_key, volume, false, response);
}

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t volume, bool is_final,
    CreditUpdateResponse* response) {
  create_charging_credit(volume, is_final, response->mutable_credit());
  response->set_success(true);
  response->set_sid(imsi);
  response->set_session_id(session_id);
  response->set_charging_key(charging_key);
}

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t volume, ChargingCredit_FinalAction action,
    std::string redirect_server, std::string restrict_rule,
    CreditUpdateResponse* response) {
  create_charging_credit(
      volume, action, redirect_server, restrict_rule,
      response->mutable_credit());
  response->set_success(true);
  response->set_sid(imsi);
  response->set_session_id(session_id);
  response->set_charging_key(charging_key);
}

void create_credit_update_response_with_error(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, bool is_final, DiameterResultCode resultCode,
    CreditUpdateResponse* response) {
  response->set_success(false);
  create_charging_credit(0, is_final, response->mutable_credit());
  response->set_sid(imsi);
  response->set_session_id(session_id);
  response->set_charging_key(charging_key);
  response->set_result_code(resultCode);
}

void create_credit_update_response_with_error(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, bool is_final, DiameterResultCode resultCode,
    ChargingCredit_FinalAction action, std::string redirect_server,
    std::string restrict_rule, CreditUpdateResponse* response) {
  response->set_success(false);
  create_charging_credit(
      0, action, redirect_server, restrict_rule, response->mutable_credit());
  response->set_sid(imsi);
  response->set_session_id(session_id);
  response->set_charging_key(charging_key);
  response->set_result_code(resultCode);
}

void create_charging_credit(
    uint64_t total_volume, uint64_t tx_volume, uint64_t rx_volume,
    bool is_final, ChargingCredit* credit) {
  create_granted_units(
      &total_volume, &tx_volume, &rx_volume, credit->mutable_granted_units());
  credit->set_type(ChargingCredit::BYTES);
  credit->set_is_final(is_final);
}

void create_credit_update_response(
    const std::string& imsi, const std::string session_id,
    uint32_t charging_key, uint64_t total_volume, uint64_t tx_volume,
    uint64_t rx_volume, bool is_final, CreditUpdateResponse* response) {
  create_charging_credit(
      total_volume, tx_volume, rx_volume, is_final, response->mutable_credit());
  response->set_success(true);
  response->set_sid(imsi);
  response->set_session_id(session_id);
  response->set_charging_key(charging_key);
}

void create_update_session_request(
    std::string imsi, std::string session_id, uint32_t ckey, std::string mkey,
    CreditUsage::UpdateType type, uint64_t bytes_rx, uint64_t bytes_tx,
    UpdateSessionRequest* usr) {
  CreditUsageUpdate* credit_update = usr->add_updates();
  create_usage_update(imsi, ckey, bytes_rx, bytes_tx, type, credit_update);

  UsageMonitoringUpdateRequest* monitor_credit_update =
      usr->add_usage_monitors();
  create_usage_monitoring_update_request(
      imsi, mkey, bytes_rx, bytes_tx, monitor_credit_update);
}

void create_usage_monitoring_update_request(
    const std::string& imsi, std::string monitoring_key, uint64_t bytes_rx,
    uint64_t bytes_tx, UsageMonitoringUpdateRequest* update) {
  auto usage = update->mutable_update();
  usage->set_monitoring_key(monitoring_key);
  usage->set_bytes_rx(bytes_rx);
  usage->set_bytes_tx(bytes_tx);
}

void create_usage_update(
    const std::string& imsi, uint32_t charging_key, uint64_t bytes_rx,
    uint64_t bytes_tx, CreditUsage::UpdateType type,
    CreditUsageUpdate* update) {
  auto usage = update->mutable_usage();
  update->mutable_common_context()->mutable_sid()->set_id(imsi);
  usage->set_charging_key(charging_key);
  usage->set_bytes_rx(bytes_rx);
  usage->set_bytes_tx(bytes_tx);
  usage->set_type(type);
}

void create_monitor_credit(
    const std::string& m_key, MonitoringLevel level, uint64_t volume,
    UsageMonitoringCredit* credit) {
  create_monitor_credit(m_key, level, volume, 0, 0, credit);
}

void create_monitor_credit(
    const std::string& m_key, MonitoringLevel level, uint64_t total_volume,
    uint64_t tx_volume, uint64_t rx_volume, UsageMonitoringCredit* credit) {
  credit->mutable_granted_units()->mutable_total()->set_volume(total_volume);
  credit->mutable_granted_units()->mutable_total()->set_is_valid(true);
  credit->mutable_granted_units()->mutable_tx()->set_volume(tx_volume);
  credit->mutable_granted_units()->mutable_tx()->set_is_valid(true);
  credit->mutable_granted_units()->mutable_rx()->set_volume(rx_volume);
  credit->mutable_granted_units()->mutable_rx()->set_is_valid(true);
  credit->set_level(level);
  credit->set_monitoring_key(m_key);
}

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t total_volume,
    uint64_t tx_volume, uint64_t rx_volume,
    UsageMonitoringUpdateResponse* response) {
  std::vector<EventTrigger> event_triggers;
  create_monitor_update_response(
      imsi, session_id, m_key, level, total_volume, tx_volume, rx_volume,
      event_triggers, 0, response);
}

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t volume,
    UsageMonitoringUpdateResponse* response) {
  std::vector<EventTrigger> event_triggers;
  create_monitor_update_response(
      imsi, session_id, m_key, level, volume, event_triggers, 0, response);
}

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t volume,
    const std::vector<EventTrigger>& event_triggers,
    const uint64_t revalidation_time_unix_ts,
    UsageMonitoringUpdateResponse* response) {
  create_monitor_update_response(
      imsi, session_id, m_key, level, volume, 0, 0, event_triggers,
      revalidation_time_unix_ts, response);
}

void create_monitor_update_response(
    const std::string& imsi, const std::string session_id,
    const std::string& m_key, MonitoringLevel level, uint64_t total_volume,
    uint64_t tx_volume, uint64_t rx_volume,
    const std::vector<EventTrigger>& event_triggers,
    const uint64_t revalidation_time_unix_ts,
    UsageMonitoringUpdateResponse* response) {
  create_monitor_credit(
      m_key, level, total_volume, tx_volume, rx_volume,
      response->mutable_credit());
  response->set_success(true);
  response->set_sid(imsi);
  response->set_session_id(session_id);
  for (const auto& event_trigger : event_triggers) {
    response->add_event_triggers(event_trigger);
  }
  response->mutable_revalidation_time()->set_seconds(revalidation_time_unix_ts);
}

void create_policy_reauth_request(
    const std::string& session_id, const std::string& imsi,
    const std::vector<std::string>& rules_to_remove,
    const std::vector<StaticRuleInstall>& rules_to_install,
    const std::vector<DynamicRuleInstall>& dynamic_rules_to_install,
    const std::vector<EventTrigger>& event_triggers,
    const uint64_t revalidation_time_unix_ts,
    const std::vector<UsageMonitoringCredit>& usage_monitoring_credits,
    PolicyReAuthRequest* request) {
  request->set_session_id(session_id);
  request->set_imsi(imsi);
  for (const auto& rule_id : rules_to_remove) {
    request->add_rules_to_remove(rule_id);
  }
  auto req_rules_to_install = request->mutable_rules_to_install();
  for (const auto& static_rule_to_install : rules_to_install) {
    req_rules_to_install->Add()->CopyFrom(static_rule_to_install);
  }
  auto req_dynamic_rules_to_install =
      request->mutable_dynamic_rules_to_install();
  for (const auto& dynamic_rule_to_install : dynamic_rules_to_install) {
    req_dynamic_rules_to_install->Add()->CopyFrom(dynamic_rule_to_install);
  }
  for (const auto& event_trigger : event_triggers) {
    request->add_event_triggers(event_trigger);
  }
  request->mutable_revalidation_time()->set_seconds(revalidation_time_unix_ts);
  auto req_credits = request->mutable_usage_monitoring_credits();
  for (const auto& credit : usage_monitoring_credits) {
    req_credits->Add()->CopyFrom(credit);
  }
}

void create_tgpp_context(
    const std::string& gx_dest_host, const std::string& gy_dest_host,
    TgppContext* context) {
  context->set_gx_dest_host(gx_dest_host);
  context->set_gy_dest_host(gy_dest_host);
}

void create_subscriber_quota_update(
    const std::string& imsi, const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state, SubscriberQuotaUpdate* update) {
  auto sid = update->mutable_sid();
  sid->set_id(imsi);
  update->set_mac_addr(ue_mac_addr);
  update->set_update_type(state);
}

void create_session_create_response(
    const std::string& imsi, const std::string session_id,
    const std::string& monitoring_key, std::vector<std::string>& static_rules,
    CreateSessionResponse* response) {
  create_monitor_update_response(
      imsi, session_id, monitoring_key, MonitoringLevel::PCC_RULE_LEVEL, 2048,
      response->mutable_usage_monitors()->Add());

  for (auto& rule_id : static_rules) {
    // insert into create session response
    StaticRuleInstall rule_install;
    rule_install.set_rule_id(rule_id);
    response->mutable_static_rules()->Add()->CopyFrom(rule_install);
  }
}

PolicyRule create_policy_rule(
    const std::string& rule_id, const std::string& m_key, const uint32_t rg) {
  PolicyRule rule;
  rule.set_id(rule_id);
  rule.set_rating_group(rg);
  rule.set_monitoring_key(m_key);
  if (rg == 0 && m_key.length() > 0) {
    rule.set_tracking_type(PolicyRule::ONLY_PCRF);
  } else if (rg > 0 && m_key.length() == 0) {
    rule.set_tracking_type(PolicyRule::ONLY_OCS);
  } else if (rg > 0 && m_key.length() > 0) {
    rule.set_tracking_type(PolicyRule::OCS_AND_PCRF);
  } else {
    rule.set_tracking_type(PolicyRule::NO_TRACKING);
  }
  return rule;
}

PolicyRule create_policy_rule_with_qos(
    const std::string& rule_id, const std::string& m_key, const uint32_t rg,
    const int qci) {
  PolicyRule rule = create_policy_rule(rule_id, m_key, rg);
  rule.mutable_qos()->set_qci(static_cast<magma::lte::FlowQos_Qci>(qci));
  return rule;
}

void create_granted_units(
    uint64_t* total, uint64_t* tx, uint64_t* rx, GrantedUnits* gsu) {
  if (total != NULL) {
    gsu->mutable_total()->set_is_valid(true);
    gsu->mutable_total()->set_volume(*total);
  }
  if (tx != NULL) {
    gsu->mutable_tx()->set_is_valid(true);
    gsu->mutable_tx()->set_volume(*tx);
  }
  if (rx != NULL) {
    gsu->mutable_rx()->set_is_valid(true);
    gsu->mutable_rx()->set_volume(*rx);
  }
}

magma::mconfig::SessionD get_default_mconfig() {
  magma::mconfig::SessionD mconfig;
  mconfig.set_log_level(magma::orc8r::LogLevel::INFO);
  mconfig.set_relay_enabled(false);
  mconfig.set_gx_gy_relay_enabled(false);
  auto wallet_config = mconfig.mutable_wallet_exhaust_detection();
  wallet_config->set_terminate_on_exhaust(false);
  return mconfig;
}

PolicyBearerBindingRequest create_policy_bearer_bind_req(
    const std::string& imsi, const uint32_t linked_bearer_id,
    const std::string& rule_id, const uint32_t bearer_id,
    const uint32_t agw_teid, const uint32_t enb_teid) {
  PolicyBearerBindingRequest bearer_bind_req;
  bearer_bind_req.mutable_sid()->set_id(imsi);
  bearer_bind_req.set_linked_bearer_id(linked_bearer_id);
  bearer_bind_req.set_policy_rule_id(rule_id);
  bearer_bind_req.set_bearer_id(bearer_id);
  bearer_bind_req.mutable_teids()->set_agw_teid(agw_teid);
  bearer_bind_req.mutable_teids()->set_enb_teid(enb_teid);
  return bearer_bind_req;
}

UpdateTunnelIdsRequest create_update_tunnel_ids_request(
    const std::string& imsi, const uint32_t bearer_id, const Teids teids) {
  return create_update_tunnel_ids_request(
      imsi, bearer_id, teids.agw_teid(), teids.enb_teid());
}

UpdateTunnelIdsRequest create_update_tunnel_ids_request(
    const std::string& imsi, const uint32_t bearer_id, const uint32_t agw_teid,
    const uint32_t enb_teid) {
  UpdateTunnelIdsRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_bearer_id(bearer_id);
  req.set_agw_teid(agw_teid);
  req.set_enb_teid(enb_teid);
  return req;
}

StaticRuleInstall create_static_rule_install(const std::string& rule_id) {
  StaticRuleInstall rule_install;
  rule_install.set_rule_id(rule_id);
  return rule_install;
}

DynamicRuleInstall create_dynamic_rule_install(const PolicyRule& rule) {
  DynamicRuleInstall rule_install;
  rule_install.mutable_policy_rule()->CopyFrom(rule);
  return rule_install;
}

}  // namespace magma
