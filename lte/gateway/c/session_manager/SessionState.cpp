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

#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>

#include <functional>
#include <string>
#include <unordered_set>
#include <utility>
#include <vector>

#include "CreditKey.h"
#include "DiameterCodes.h"
#include "EnumToString.h"
#include "magma_logging.h"
#include "MetricsHelpers.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "StoredState.h"
#include "Utilities.h"

namespace {
const char* UE_TRAFFIC_COUNTER_NAME = "ue_traffic";
const char* UE_DROPPED_COUNTER_NAME = "ue_dropped_usage";
const char* UE_USED_COUNTER_NAME    = "ue_reported_usage";
const char* LABEL_IMSI              = "IMSI";
const char* LABEL_APN               = "apn";
const char* LABEL_MSISDN            = "msisdn";
const char* LABEL_DIRECTION         = "direction";
const char* DIRECTION_UP            = "up";
const char* DIRECTION_DOWN          = "down";
const char* LABEL_SESSION_ID        = "session_id";
}  // namespace

using magma::service303::increment_counter;
using magma::service303::remove_counter;

namespace magma {

template<class T>
void remove_from_vec_by_value(std::vector<T>& vec, T value) {
  vec.erase(std::remove(vec.begin(), vec.end(), value), vec.end());
}

std::unique_ptr<SessionState> SessionState::unmarshal(
    const StoredSessionState& marshaled, StaticRuleStore& rule_store) {
  return std::make_unique<SessionState>(marshaled, rule_store);
}

StoredSessionState SessionState::marshal() {
  StoredSessionState marshaled{};

  marshaled.fsm_state  = curr_state_;
  marshaled.config     = config_;
  marshaled.imsi       = get_imsi();
  marshaled.session_id = session_id_;
  marshaled.local_teid = local_teid_;
  // 5G session version handling
  marshaled.current_version         = current_version_;
  marshaled.subscriber_quota_state  = subscriber_quota_state_;
  marshaled.tgpp_context            = tgpp_context_;
  marshaled.request_number          = request_number_;
  marshaled.pdp_start_time          = pdp_start_time_;
  marshaled.pdp_end_time            = pdp_end_time_;
  marshaled.pending_event_triggers  = pending_event_triggers_;
  marshaled.revalidation_time       = revalidation_time_;
  marshaled.bearer_id_by_policy     = bearer_id_by_policy_;
  marshaled.create_session_response = create_session_response_;

  marshaled.monitor_map = StoredMonitorMap();
  for (auto& monitor_pair : monitor_map_) {
    StoredMonitor monitor{};
    monitor.credit = monitor_pair.second->credit.marshal();
    monitor.level  = monitor_pair.second->level;
    marshaled.monitor_map[monitor_pair.first] = monitor;
  }
  marshaled.session_level_key = session_level_key_;

  marshaled.credit_map = StoredChargingCreditMap(4, &ccHash, &ccEqual);
  for (auto& credit_pair : credit_map_) {
    auto key                  = CreditKey();
    key.rating_group          = credit_pair.first.rating_group;
    key.service_identifier    = credit_pair.first.service_identifier;
    marshaled.credit_map[key] = credit_pair.second->marshal();
  }

  for (auto& rule_id : active_static_rules_) {
    marshaled.static_rule_ids.push_back(rule_id);
  }
  for (auto& rule : PdrList_) {
    marshaled.PdrList.push_back(rule);
  }

  std::vector<PolicyRule> dynamic_rules;
  dynamic_rules_.get_rules(dynamic_rules);
  marshaled.dynamic_rules = std::move(dynamic_rules);

  std::vector<PolicyRule> gy_dynamic_rules;
  gy_dynamic_rules_.get_rules(gy_dynamic_rules);
  marshaled.gy_dynamic_rules = std::move(gy_dynamic_rules);

  for (auto& rule_id : scheduled_static_rules_) {
    marshaled.scheduled_static_rules.insert(rule_id);
  }
  std::vector<PolicyRule> scheduled_dynamic_rules;
  scheduled_dynamic_rules_.get_rules(scheduled_dynamic_rules);
  marshaled.scheduled_dynamic_rules = std::move(scheduled_dynamic_rules);

  for (auto& it : rule_lifetimes_) {
    marshaled.rule_lifetimes[it.first] = it.second;
  }
  marshaled.policy_version_and_stats = policy_version_and_stats_;
  return marshaled;
}

SessionState::SessionState(
    const StoredSessionState& marshaled, StaticRuleStore& rule_store)
    : imsi_(marshaled.imsi),
      session_id_(marshaled.session_id),
      local_teid_(marshaled.local_teid),
      request_number_(marshaled.request_number),
      curr_state_(marshaled.fsm_state),
      config_(marshaled.config),
      pdp_start_time_(marshaled.pdp_start_time),
      pdp_end_time_(marshaled.pdp_end_time),
      // 5G session version handlimg
      current_version_(marshaled.current_version),
      subscriber_quota_state_(marshaled.subscriber_quota_state),
      tgpp_context_(marshaled.tgpp_context),
      create_session_response_(marshaled.create_session_response),
      policy_version_and_stats_(marshaled.policy_version_and_stats),
      static_rules_(rule_store),
      pending_event_triggers_(marshaled.pending_event_triggers),
      revalidation_time_(marshaled.revalidation_time),
      credit_map_(4, &ccHash, &ccEqual),
      bearer_id_by_policy_(marshaled.bearer_id_by_policy) {
  session_level_key_ = marshaled.session_level_key;
  for (auto it : marshaled.monitor_map) {
    Monitor monitor;
    monitor.credit = SessionCredit(it.second.credit);
    monitor.level  = it.second.level;

    monitor_map_[it.first] = std::make_unique<Monitor>(monitor);
  }

  for (const auto& it : marshaled.credit_map) {
    credit_map_[it.first] =
        std::make_unique<ChargingGrant>(ChargingGrant(it.second));
  }

  for (const std::string& rule_id : marshaled.static_rule_ids) {
    active_static_rules_.push_back(rule_id);
  }
  for (auto& rule : marshaled.dynamic_rules) {
    dynamic_rules_.insert_rule(rule);
  }
  for (auto& rule : marshaled.PdrList) {
    PdrList_.push_back(rule);
  }
  for (const std::string& rule_id : marshaled.scheduled_static_rules) {
    scheduled_static_rules_.insert(rule_id);
  }
  for (auto& rule : marshaled.scheduled_dynamic_rules) {
    scheduled_dynamic_rules_.insert_rule(rule);
  }
  for (auto& it : marshaled.rule_lifetimes) {
    rule_lifetimes_[it.first] = it.second;
  }
  for (auto& rule : marshaled.gy_dynamic_rules) {
    gy_dynamic_rules_.insert_rule(rule);
  }
}

SessionState::SessionState(
    const std::string& imsi, const std::string& session_id,
    const SessionConfig& cfg, StaticRuleStore& rule_store,
    const magma::lte::TgppContext& tgpp_context, uint64_t pdp_start_time,
    const CreateSessionResponse& csr)
    : imsi_(imsi),
      session_id_(session_id),
      local_teid_(0),
      // Request number set to 1, because request 0 is INIT call
      request_number_(1),
      curr_state_(SESSION_ACTIVE),
      config_(cfg),
      pdp_start_time_(pdp_start_time),
      pdp_end_time_(0),
      tgpp_context_(tgpp_context),
      create_session_response_(csr),
      static_rules_(rule_store),
      credit_map_(4, &ccHash, &ccEqual) {
  // other default initializations
  current_version_        = 0;
  session_level_key_      = "";
  subscriber_quota_state_ = SubscriberQuotaUpdate_Type_VALID_QUOTA;
}

/*For 5G which doesn't have response context*/
SessionState::SessionState(
    const std::string& imsi, const std::string& session_ctx_id,
    const SessionConfig& cfg, StaticRuleStore& rule_store)
    : imsi_(imsi),
      session_id_(session_ctx_id),
      local_teid_(0),
      // Request number set to 1, because request 0 is INIT call
      request_number_(1),
      /*current state would be CREATING and version would be 0 */
      curr_state_(CREATING),
      config_(cfg),
      current_version_(0),
      static_rules_(rule_store) {}

/* get-set methods of new messages  for 5G*/
uint32_t SessionState::get_current_version() {
  return current_version_;
}

void SessionState::set_current_version(
    int new_session_version, SessionStateUpdateCriteria& session_uc) {
  current_version_                      = new_session_version;
  session_uc.is_current_version_updated = true;
  session_uc.updated_current_version    = new_session_version;
  MLOG(MINFO) << " Current version is " << get_current_version();
}
/* Add PDR rule to this rules session list */
void SessionState::insert_pdr(SetGroupPDR* rule) {
  PdrList_.push_back(*rule);
}

void SessionState::set_remove_all_pdrs() {
  for (auto& rule : PdrList_) {
    rule.set_pdr_state(PdrState::REMOVE);
  }
}
/* Remove all Pdr, FAR rules */
void SessionState::remove_all_rules() {
  PdrList_.clear();
}

/* It gets all PDR rule list of the session */
std::vector<SetGroupPDR>& SessionState::get_all_pdr_rules() {
  return PdrList_;
}

SessionFsmState SessionState::get_state() {
  return curr_state_;
}

magma::lte::Fsm_state_FsmState SessionState::get_proto_fsm_state() {
  SessionFsmState curr_state = get_state();
  switch (curr_state) {
    case CREATING:
      return magma::lte::Fsm_state_FsmState_CREATING;
      break;
    case CREATED:
      return magma::lte::Fsm_state_FsmState_CREATED;
      break;
    case ACTIVE:
      return magma::lte::Fsm_state_FsmState_ACTIVE;
      break;
    case RELEASE:
      return magma::lte::Fsm_state_FsmState_RELEASE;
      break;
    case INACTIVE:
    default:
      return magma::lte::Fsm_state_FsmState_INACTIVE;
      break;
  }
  return magma::lte::Fsm_state_FsmState_INACTIVE;
}

/*temporary copy to be removed after upf node code completes */
void SessionState::sess_infocopy(struct SessionInfo* info) {
  // Static SessionInfo vlaue till UPF node value implementation
  // gets stablized.
  std::string imsi_num;
  // TODO we could eventually  migrate to SMF-UPF proto enum directly.
  info->state = get_proto_fsm_state();
  info->subscriber_id.assign(get_imsi());
  info->ver_no              = get_current_version();
  info->nodeId.node_id_type = SessionInfo::IPv4;
  strcpy(info->nodeId.node_id, "192.168.2.1");
  /* TODO below to be changed after UPF node association message
   * completes . Revisit
   */
}

void SessionState::set_teids(uint32_t enb_teid, uint32_t agw_teid) {
  Teids teids;
  teids.set_agw_teid(agw_teid);
  teids.set_enb_teid(enb_teid);
  set_teids(teids);
}

void SessionState::set_teids(Teids teids) {
  config_.common_context.mutable_teids()->CopyFrom(teids);
}

static UsageMonitorUpdate make_usage_monitor_update(
    const SessionCredit::Usage& usage_in, const std::string& monitoring_key,
    MonitoringLevel level) {
  UsageMonitorUpdate update;
  update.set_bytes_tx(usage_in.bytes_tx);
  update.set_bytes_rx(usage_in.bytes_rx);
  update.set_level(level);
  update.set_monitoring_key(monitoring_key);
  return update;
}

SessionCreditUpdateCriteria* SessionState::get_credit_uc(
    const CreditKey& key, SessionStateUpdateCriteria& uc) {
  if (uc.charging_credit_map.find(key) == uc.charging_credit_map.end()) {
    uc.charging_credit_map[key] = credit_map_[key]->get_update_criteria();
  }
  return &(uc.charging_credit_map[key]);
}

bool SessionState::apply_update_criteria(SessionStateUpdateCriteria& uc) {
  if (uc.is_fsm_updated) {
    curr_state_ = uc.updated_fsm_state;
  }

  if (uc.is_current_version_updated) {
    current_version_ = uc.updated_current_version;
  }

  if (uc.is_local_teid_updated) {
    local_teid_ = uc.local_teid_updated;
  }

  if (uc.is_pending_event_triggers_updated) {
    for (auto it : uc.pending_event_triggers) {
      pending_event_triggers_[it.first] = it.second;
      if (it.first == REVALIDATION_TIMEOUT) {
        revalidation_time_ = uc.revalidation_time;
      }
    }
  }
  // QoS Management
  if (uc.is_bearer_mapping_updated) {
    bearer_id_by_policy_ = uc.bearer_id_by_policy;
  }

  // Config
  if (uc.is_config_updated) {
    config_ = uc.updated_config;
  }

  // Rule versions
  if (uc.policy_version_and_stats) {
    policy_version_and_stats_ = *uc.policy_version_and_stats;
  }

  // Manually update these policy structures to avoid incrementing version
  // Static rules
  for (const auto& rule_id : uc.static_rules_to_uninstall) {
    if (is_static_rule_installed(rule_id)) {
      remove_from_vec_by_value<std::string>(active_static_rules_, rule_id);
    }
    if (is_static_rule_scheduled(rule_id)) {
      scheduled_static_rules_.erase(rule_id);
    }
    rule_lifetimes_.erase(rule_id);
  }
  for (const auto& rule_id : uc.static_rules_to_install) {
    if (!is_static_rule_installed(rule_id)) {
      active_static_rules_.push_back(rule_id);
    }
    if (uc.new_rule_lifetimes.find(rule_id) != uc.new_rule_lifetimes.end()) {
      rule_lifetimes_[rule_id] = uc.new_rule_lifetimes[rule_id];
    }
    if (is_static_rule_scheduled(rule_id)) {
      scheduled_static_rules_.erase(rule_id);
    }
  }
  for (const auto& rule_id : uc.new_scheduled_static_rules) {
    if (is_static_rule_scheduled(rule_id)) {
      continue;
    }
    if (uc.new_rule_lifetimes.find(rule_id) != uc.new_rule_lifetimes.end()) {
      rule_lifetimes_[rule_id] = uc.new_rule_lifetimes[rule_id];
    }
    scheduled_static_rules_.insert(rule_id);
  }

  // Dynamic rules
  for (const auto& rule_id : uc.dynamic_rules_to_uninstall) {
    scheduled_dynamic_rules_.remove_rule(rule_id, nullptr);
    dynamic_rules_.remove_rule(rule_id, nullptr);
    rule_lifetimes_.erase(rule_id);
  }
  for (const auto& rule : uc.dynamic_rules_to_install) {
    if (uc.new_rule_lifetimes.find(rule.id()) != uc.new_rule_lifetimes.end()) {
      rule_lifetimes_[rule.id()] = uc.new_rule_lifetimes[rule.id()];
    }
    dynamic_rules_.insert_rule(rule);
    scheduled_dynamic_rules_.remove_rule(rule.id(), nullptr);
  }
  for (const auto& rule : uc.new_scheduled_dynamic_rules) {
    if (uc.new_rule_lifetimes.find(rule.id()) != uc.new_rule_lifetimes.end()) {
      rule_lifetimes_[rule.id()] = uc.new_rule_lifetimes[rule.id()];
    }
    scheduled_dynamic_rules_.insert_rule(rule);
  }

  // Gy Dynamic rules
  for (const auto& rule : uc.gy_dynamic_rules_to_install) {
    if (uc.new_rule_lifetimes.find(rule.id()) != uc.new_rule_lifetimes.end()) {
      rule_lifetimes_[rule.id()] = uc.new_rule_lifetimes[rule.id()];
    }
    gy_dynamic_rules_.insert_rule(rule);
  }
  for (const auto& rule_id : uc.gy_dynamic_rules_to_uninstall) {
    gy_dynamic_rules_.remove_rule(rule_id, nullptr);
  }

  // Charging credit
  for (const auto& it : uc.charging_credit_map) {
    auto key           = it.first;
    auto credit_update = it.second;
    apply_charging_credit_update(key, credit_update);
  }
  for (const auto& it : uc.charging_credit_to_install) {
    auto key           = it.first;
    auto stored_credit = it.second;
    credit_map_[key]   = std::make_unique<ChargingGrant>(stored_credit);
  }

  // Monitoring credit
  if (uc.is_session_level_key_updated) {
    set_session_level_key(uc.updated_session_level_key);
  }
  for (const auto& it : uc.monitor_credit_map) {
    auto key           = it.first;
    auto credit_update = it.second;
    apply_monitor_updates(key, uc, credit_update);
  }
  for (const auto& it : uc.monitor_credit_to_install) {
    auto key            = it.first;
    auto stored_monitor = it.second;
    monitor_map_[key]   = std::make_unique<Monitor>(stored_monitor);
  }

  if (uc.updated_pdp_end_time > 0) {
    pdp_end_time_ = uc.updated_pdp_end_time;
  }

  return true;
}

void SessionState::add_rule_usage(
    const std::string& rule_id, uint64_t used_tx, uint64_t used_rx,
    uint64_t dropped_tx, uint64_t dropped_rx,
    SessionStateUpdateCriteria& update_criteria) {
  CreditKey charging_key;
  if (dynamic_rules_.get_charging_key_for_rule_id(rule_id, &charging_key) ||
      static_rules_.get_charging_key_for_rule_id(rule_id, &charging_key)) {
    MLOG(MINFO) << "Updating used charging credit for Rule=" << rule_id
                << " Rating Group=" << charging_key.rating_group
                << " Service Identifier=" << charging_key.service_identifier;
    auto it = credit_map_.find(charging_key);
    if (it != credit_map_.end()) {
      auto credit_uc = get_credit_uc(charging_key, update_criteria);
      it->second->credit.add_used_credit(used_tx, used_rx, *credit_uc);
      if (it->second->should_deactivate_service()) {
        it->second->set_service_state(SERVICE_NEEDS_DEACTIVATION, *credit_uc);
      }
    } else {
      MLOG(MDEBUG) << "Rating Group " << charging_key.rating_group
                   << " not found, not adding the usage";
    }
  }
  std::string monitoring_key;
  if (dynamic_rules_.get_monitoring_key_for_rule_id(rule_id, &monitoring_key) ||
      static_rules_.get_monitoring_key_for_rule_id(rule_id, &monitoring_key)) {
    MLOG(MINFO) << "Updating used monitoring credit for Rule=" << rule_id
                << " Monitoring Key=" << monitoring_key;
    add_to_monitor(monitoring_key, used_tx, used_rx, update_criteria);
  }
  if (session_level_key_ != "" && monitoring_key != session_level_key_) {
    // Update session level key if its different
    add_to_monitor(session_level_key_, used_tx, used_rx, update_criteria);
  }
  if (is_dynamic_rule_installed(rule_id) || is_static_rule_installed(rule_id)) {
    update_data_metrics(UE_USED_COUNTER_NAME, used_tx, used_rx);
  }
  update_data_metrics(UE_DROPPED_COUNTER_NAME, dropped_tx, dropped_rx);
}

void SessionState::apply_session_rule_set(
    const RuleSetToApply& rule_set, RulesToProcess* pending_activation,
    RulesToProcess* pending_deactivation, RulesToProcess* pending_bearer_setup,
    SessionStateUpdateCriteria& uc) {
  apply_session_static_rule_set(
      rule_set.static_rules, pending_activation, pending_deactivation,
      pending_bearer_setup, uc);
  apply_session_dynamic_rule_set(
      rule_set.dynamic_rules, pending_activation, pending_deactivation,
      pending_bearer_setup, uc);
}

void SessionState::apply_session_static_rule_set(
    const std::unordered_set<std::string> static_rules,
    RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
    RulesToProcess* pending_bearer_setup, SessionStateUpdateCriteria& uc) {
  // No activation time / deactivation support yet for rule set interface
  RuleLifetime lifetime;
  // Go through the rule set and install any rules not yet installed
  for (const auto& static_rule_id : static_rules) {
    PolicyRule rule;
    if (!static_rules_.get_rule(static_rule_id, &rule)) {
      MLOG(MERROR) << "Static rule " << static_rule_id
                   << " is not found. Skipping activation";
      continue;
    }
    if (is_static_rule_installed(static_rule_id)) {
      continue;
    }

    MLOG(MINFO) << "Installing static rule " << static_rule_id << " for "
                << session_id_;
    RuleToProcess to_process =
        activate_static_rule(static_rule_id, lifetime, uc);
    classify_policy_activation(
        to_process, STATIC, pending_activation, pending_bearer_setup);
  }
  std::vector<PolicyRule> static_pending_deactivation;

  // Go through the existing rules and uninstall any rule not in the rule set
  for (const auto static_rule_id : active_static_rules_) {
    if (static_rules.find(static_rule_id) == static_rules.end()) {
      PolicyRule rule;
      if (static_rules_.get_rule(static_rule_id, &rule)) {
        static_pending_deactivation.push_back(rule);
      }
    }
  }
  // Do the actual removal separately so we're not modifying the vector while
  // looping
  for (const PolicyRule static_rule : static_pending_deactivation) {
    MLOG(MINFO) << "Removing static rule " << static_rule.id() << " for "
                << session_id_;
    optional<RuleToProcess> op_rule_info =
        deactivate_static_rule(static_rule.id(), uc);
    if (!op_rule_info) {
      MLOG(MWARNING) << "Failed to deactivate static rule " << static_rule.id()
                     << " for " << session_id_;
    } else {
      pending_deactivation->push_back(*op_rule_info);
    }
  }
}

void SessionState::apply_session_dynamic_rule_set(
    const std::unordered_map<std::string, PolicyRule> dynamic_rules,
    RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
    RulesToProcess* pending_bearer_setup, SessionStateUpdateCriteria& uc) {
  // No activation time / deactivation support yet for rule set interface
  RuleLifetime lifetime;
  for (const auto& dynamic_rule_pair : dynamic_rules) {
    if (is_dynamic_rule_installed(dynamic_rule_pair.first)) {
      continue;
    }
    MLOG(MINFO) << "Installing dynamic rule " << dynamic_rule_pair.first
                << " for " << session_id_;
    RuleToProcess to_process =
        insert_dynamic_rule(dynamic_rule_pair.second, lifetime, uc);
    classify_policy_activation(
        to_process, DYNAMIC, pending_activation, pending_bearer_setup);
  }
  std::vector<PolicyRule> active_dynamic_rules;
  dynamic_rules_.get_rules(active_dynamic_rules);
  for (const auto& dynamic_rule : active_dynamic_rules) {
    if (dynamic_rules.find(dynamic_rule.id()) == dynamic_rules.end()) {
      MLOG(MINFO) << "Removing dynamic rule " << dynamic_rule.id() << " for "
                  << session_id_;
      pending_deactivation->push_back(
          *remove_dynamic_rule(dynamic_rule.id(), nullptr, uc));
    }
  }
}

void SessionState::set_subscriber_quota_state(
    const magma::lte::SubscriberQuotaUpdate_Type state,
    SessionStateUpdateCriteria* session_uc) {
  if (session_uc != nullptr) {
    session_uc->updated_subscriber_quota_state = state;
  }
  subscriber_quota_state_ = state;
}

bool SessionState::active_monitored_rules_exist() {
  return total_monitored_rules_count() > 0;
}

bool SessionState::is_terminating() {
  if (curr_state_ == SESSION_RELEASED || curr_state_ == SESSION_TERMINATED ||
      curr_state_ == RELEASE) {
    return true;
  }
  return false;
}

void SessionState::get_monitor_updates(
    UpdateSessionRequest& update_request_out,
    SessionStateUpdateCriteria& update_criteria) {
  for (auto& monitor_pair : monitor_map_) {
    if (!monitor_pair.second->should_send_update()) {
      continue;  // no update
    }

    auto mkey      = monitor_pair.first;
    auto& credit   = monitor_pair.second->credit;
    auto credit_uc = get_monitor_uc(mkey, update_criteria);

    if (curr_state_ == SESSION_RELEASED) {
      MLOG(MDEBUG)
          << "Session " << session_id_
          << " is in Session Released state. Not sending update to the core"
             "for monitor "
          << mkey;
      continue;  // no update
    }

    std::string last_update_message = "";
    if (credit.is_report_last_credit()) {
      // in case this is the last report, here we mark the session as to be
      // deleted
      credit_uc->deleted  = true;
      last_update_message = " (this is the last updcate sent for this monitor)";
    }

    MLOG(MDEBUG) << "Session " << session_id_ << " monitoring key " << mkey
                 << " updating due to quota exhaustion"
                 << " with request number " << request_number_
                 << last_update_message;

    auto usage = credit.get_usage_for_reporting(*credit_uc);
    auto update =
        make_usage_monitor_update(usage, mkey, monitor_pair.second->level);
    auto new_req = update_request_out.mutable_usage_monitors()->Add();

    add_common_fields_to_usage_monitor_update(new_req);
    new_req->mutable_update()->CopyFrom(update);
    new_req->set_event_trigger(USAGE_REPORT);
    request_number_++;
    update_criteria.request_number_increment++;
  }
}

void SessionState::add_common_fields_to_usage_monitor_update(
    UsageMonitoringUpdateRequest* req) {
  req->set_session_id(session_id_);
  req->set_request_number(request_number_);
  req->set_sid(get_imsi());
  req->set_ue_ipv4(config_.common_context.ue_ipv4());
  req->set_rat_type(config_.common_context.rat_type());
  fill_protos_tgpp_context(req->mutable_tgpp_ctx());
  if (config_.rat_specific_context.has_wlan_context()) {
    const auto& wlan_context = config_.rat_specific_context.wlan_context();
    req->set_hardware_addr(wlan_context.mac_addr_binary());
  } else {
    const auto& lte_context = config_.rat_specific_context.lte_context();
    req->set_charging_characteristics(lte_context.charging_characteristics());
  }
}

void SessionState::get_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& update_criteria) {
  if (curr_state_ != SESSION_ACTIVE) return;
  get_charging_updates(update_request_out, actions_out, update_criteria);
  get_monitor_updates(update_request_out, update_criteria);
  get_event_trigger_updates(update_request_out, update_criteria);
}

SubscriberQuotaUpdate_Type SessionState::get_subscriber_quota_state() const {
  return subscriber_quota_state_;
}

bool SessionState::can_complete_termination(SessionStateUpdateCriteria& uc) {
  switch (curr_state_) {
    case SESSION_ACTIVE:
      MLOG(MERROR) << "Encountered unexpected state 'ACTIVE' when "
                   << "completing termination for " << session_id_
                   << " Not terminating...";
      return false;
    case SESSION_TERMINATED:
      // session is already terminated. Do nothing.
      return false;
    default:
      // Continue termination but no logs are necessary for other states
      break;
  }
  // mark session as terminated
  set_fsm_state(SESSION_TERMINATED, uc);
  return true;
}

SessionTerminateRequest SessionState::make_termination_request(
    SessionStateUpdateCriteria& uc) {
  SessionTerminateRequest req;
  req.set_session_id(session_id_);
  req.set_request_number(request_number_);
  req.mutable_common_context()->CopyFrom(config_.common_context);

  fill_protos_tgpp_context(req.mutable_tgpp_ctx());
  if (config_.rat_specific_context.has_lte_context()) {
    const auto& lte_context = config_.rat_specific_context.lte_context();
    req.set_spgw_ipv4(lte_context.spgw_ipv4());
    req.set_imei(lte_context.imei());
    req.set_plmn_id(lte_context.plmn_id());
    req.set_imsi_plmn_id(lte_context.imsi_plmn_id());
    req.set_user_location(lte_context.user_location());
    req.set_charging_characteristics(lte_context.charging_characteristics());
  } else if (config_.rat_specific_context.has_wlan_context()) {
    const auto& wlan_context = config_.rat_specific_context.wlan_context();
    req.set_hardware_addr(wlan_context.mac_addr_binary());
  }

  // gx monitors
  for (auto& credit_pair : monitor_map_) {
    auto credit_uc = get_monitor_uc(credit_pair.first, uc);
    req.mutable_monitor_usages()->Add()->CopyFrom(make_usage_monitor_update(
        credit_pair.second->credit.get_all_unreported_usage_for_reporting(
            *credit_uc),
        credit_pair.first, credit_pair.second->level));
  }
  // gy credits
  for (auto& credit_pair : credit_map_) {
    auto credit_uc    = get_credit_uc(credit_pair.first, uc);
    auto credit_usage = credit_pair.second->get_credit_usage(
        CreditUsage::TERMINATED, *credit_uc, true);
    credit_pair.first.set_credit_usage(&credit_usage);
    req.mutable_credit_usages()->Add()->CopyFrom(credit_usage);
  }
  return req;
}

ChargingCreditSummaries SessionState::get_charging_credit_summaries() {
  ChargingCreditSummaries charging_credit_summaries(
      credit_map_.size(), &ccHash, &ccEqual);
  for (auto& credit_pair : credit_map_) {
    auto summary = credit_pair.second->credit.get_credit_summary();
    charging_credit_summaries[credit_pair.first] = summary;
  }
  return charging_credit_summaries;
}

SessionCredit::TotalCreditUsage SessionState::get_total_credit_usage() {
  // Collate unique charging/monitoring keys used by rules
  std::unordered_set<CreditKey, decltype(&ccHash), decltype(&ccEqual)>
      used_charging_keys(4, ccHash, ccEqual);
  std::unordered_set<std::string> used_monitoring_keys;

  std::vector<std::reference_wrapper<PolicyRuleBiMap>> bimaps{static_rules_,
                                                              dynamic_rules_};

  for (auto bimap : bimaps) {
    PolicyRuleBiMap& rules = bimap;
    std::vector<std::string> rule_ids{};
    std::vector<std::string>& rule_ids_ptr = rule_ids;
    rules.get_rule_ids(rule_ids_ptr);

    for (auto rule_id : rule_ids) {
      CreditKey charging_key;
      bool should_track_charging_key =
          rules.get_charging_key_for_rule_id(rule_id, &charging_key);
      std::string monitoring_key;
      bool should_track_monitoring_key =
          rules.get_monitoring_key_for_rule_id(rule_id, &monitoring_key);

      if (should_track_charging_key) {
        used_charging_keys.insert(charging_key);
      }
      if (should_track_monitoring_key) {
        used_monitoring_keys.insert(monitoring_key);
      }
    }
  }

  // Sum up usage
  SessionCredit::TotalCreditUsage usage{
      .monitoring_tx = 0,
      .monitoring_rx = 0,
      .charging_tx   = 0,
      .charging_rx   = 0,
  };
  for (auto monitoring_key : used_monitoring_keys) {
    usage.monitoring_tx += get_monitor(monitoring_key, USED_TX);
    usage.monitoring_rx += get_monitor(monitoring_key, USED_RX);
  }
  for (auto charging_key : used_charging_keys) {
    auto it = credit_map_.find(charging_key);
    if (it != credit_map_.end()) {
      usage.charging_tx += it->second->credit.get_credit(USED_TX);
      usage.charging_rx += it->second->credit.get_credit(USED_RX);
    }
  }
  return usage;
}

uint32_t SessionState::get_local_teid() const {
  return local_teid_;
}

void SessionState::set_local_teid(
    uint32_t teid, SessionStateUpdateCriteria& uc) {
  local_teid_              = teid;
  uc.is_local_teid_updated = true;
  uc.local_teid_updated    = teid;
  return;
}

void SessionState::set_config(const SessionConfig& config) {
  config_ = config;
}

bool SessionState::is_radius_cwf_session() const {
  return (config_.common_context.rat_type() == RATType::TGPP_WLAN);
}

void SessionState::get_session_info(SessionState::SessionInfo& info) {
  info.imsi      = get_imsi();
  info.ip_addr   = config_.common_context.ue_ipv4();
  info.ipv6_addr = config_.common_context.ue_ipv6();
  info.teids     = config_.common_context.teids();
  info.msisdn    = config_.common_context.msisdn();
  info.ambr      = config_.get_apn_ambr();

  std::vector<PolicyRule> gx_dynamic_rules, gy_dynamic_rules;
  dynamic_rules_.get_rules(gx_dynamic_rules);
  gy_dynamic_rules_.get_rules(gy_dynamic_rules);

  // Set versions
  for (const PolicyRule rule : gx_dynamic_rules) {
    info.gx_rules.push_back(make_rule_to_process(rule));
  }
  for (const PolicyRule rule : gy_dynamic_rules) {
    info.gy_dynamic_rules.push_back(make_rule_to_process(rule));
  }

  for (const std::string& rule_id : active_static_rules_) {
    PolicyRule rule;
    if (static_rules_.get_rule(rule_id, &rule)) {
      info.gx_rules.push_back(make_rule_to_process(rule));
    }
  }
}

std::vector<PolicyRule> SessionState::get_all_active_policies() {
  std::vector<PolicyRule> policies;
  for (auto& rule_id : active_static_rules_) {
    PolicyRule policy;
    if (static_rules_.get_rule(rule_id, &policy)) {
      policies.push_back(policy);
    }
  }
  dynamic_rules_.get_rules(policies);
  gy_dynamic_rules_.get_rules(policies);
  return policies;
}

void SessionState::remove_all_rules_for_termination(
    SessionStateUpdateCriteria& session_uc) {
  std::vector<PolicyRule> gx_dynamic_rules, gy_dynamic_rules,
      scheduled_dynamic_rules;
  dynamic_rules_.get_rules(gx_dynamic_rules);
  for (PolicyRule& policy : gx_dynamic_rules) {
    remove_dynamic_rule(policy.id(), nullptr, session_uc);
  }
  gy_dynamic_rules_.get_rules(gy_dynamic_rules);
  for (PolicyRule& policy : gy_dynamic_rules) {
    remove_gy_rule(policy.id(), nullptr, session_uc);
  }
  for (const std::string& rule_id : active_static_rules_) {
    deactivate_static_rule(rule_id, session_uc);
  }

  // remove scheduled rules
  for (const std::string& rule_id : scheduled_static_rules_) {
    deactivate_scheduled_static_rule(rule_id);
  }
  scheduled_dynamic_rules_.get_rules(scheduled_dynamic_rules);
  for (PolicyRule& policy : scheduled_dynamic_rules) {
    remove_scheduled_dynamic_rule(policy.id(), nullptr, session_uc);
  }
}

void SessionState::set_tgpp_context(
    const magma::lte::TgppContext& tgpp_context,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.updated_tgpp_context = tgpp_context;
  tgpp_context_                        = tgpp_context;
}

void SessionState::fill_protos_tgpp_context(
    magma::lte::TgppContext* tgpp_context) const {
  *tgpp_context = tgpp_context_;
}

uint32_t SessionState::get_request_number() {
  return request_number_;
}

uint64_t SessionState::get_pdp_start_time() {
  return pdp_start_time_;
}

uint64_t SessionState::get_pdp_end_time() {
  return pdp_end_time_;
}

uint64_t SessionState::get_active_duration_in_seconds() {
  if (pdp_end_time_ > 0) {  // session has ended
    return pdp_end_time_ - pdp_start_time_;
  }
  // session is still active
  return magma::get_time_in_sec_since_epoch() - pdp_start_time_;
}

void SessionState::set_pdp_end_time(
    uint64_t epoch, SessionStateUpdateCriteria& session_uc) {
  pdp_end_time_                   = epoch;
  session_uc.updated_pdp_end_time = epoch;
}

void SessionState::increment_request_number(uint32_t incr) {
  request_number_ += incr;
}

bool SessionState::is_dynamic_rule_scheduled(const std::string& rule_id) {
  return scheduled_dynamic_rules_.get_rule(rule_id, NULL);
}

bool SessionState::is_static_rule_scheduled(const std::string& rule_id) {
  return scheduled_static_rules_.count(rule_id) == 1;
}

bool SessionState::is_dynamic_rule_installed(const std::string& rule_id) {
  return dynamic_rules_.get_rule(rule_id, NULL);
}

bool SessionState::is_gy_dynamic_rule_installed(const std::string& rule_id) {
  return gy_dynamic_rules_.get_rule(rule_id, NULL);
}

bool SessionState::is_static_rule_installed(const std::string& rule_id) {
  return std::find(
             active_static_rules_.begin(), active_static_rules_.end(),
             rule_id) != active_static_rules_.end();
}

RuleToProcess SessionState::insert_dynamic_rule(
    const PolicyRule& rule, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& session_uc) {
  rule_lifetimes_[rule.id()] = lifetime;
  dynamic_rules_.insert_rule(rule);
  session_uc.dynamic_rules_to_install.push_back(rule);
  session_uc.new_rule_lifetimes[rule.id()] = lifetime;
  increment_rule_stats(rule.id(), session_uc);

  return make_rule_to_process(rule);
}

RuleToProcess SessionState::insert_gy_rule(
    const PolicyRule& rule, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& session_uc) {
  rule_lifetimes_[rule.id()] = lifetime;
  gy_dynamic_rules_.insert_rule(rule);
  session_uc.gy_dynamic_rules_to_install.push_back(rule);
  session_uc.new_rule_lifetimes[rule.id()] = lifetime;
  increment_rule_stats(rule.id(), session_uc);

  return make_rule_to_process(rule);
}

RuleToProcess SessionState::activate_static_rule(
    const std::string& rule_id, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& session_uc) {
  RuleToProcess to_process;
  PolicyRule rule;
  static_rules_.get_rule(rule_id, &rule);

  rule_lifetimes_[rule_id] = lifetime;
  if (!is_static_rule_installed(rule_id)) {
    active_static_rules_.push_back(rule_id);
  }
  session_uc.static_rules_to_install.insert(rule_id);
  session_uc.new_rule_lifetimes[rule_id] = lifetime;
  increment_rule_stats(rule_id, session_uc);

  return make_rule_to_process(rule);
}

optional<RuleToProcess> SessionState::remove_dynamic_rule(
    const std::string& rule_id, PolicyRule* rule_out,
    SessionStateUpdateCriteria& session_uc) {
  PolicyRule rule;
  bool removed = dynamic_rules_.remove_rule(rule_id, &rule);
  if (!removed) {
    return {};
  }
  if (rule_out != nullptr) {
    *rule_out = rule;
  }

  session_uc.dynamic_rules_to_uninstall.insert(rule_id);
  increment_rule_stats(rule_id, session_uc);

  return make_rule_to_process(rule);
}

bool SessionState::remove_scheduled_dynamic_rule(
    const std::string& rule_id, PolicyRule* rule_out,
    SessionStateUpdateCriteria& update_criteria) {
  bool removed = scheduled_dynamic_rules_.remove_rule(rule_id, rule_out);
  if (removed) {
    update_criteria.dynamic_rules_to_uninstall.insert(rule_id);
  }
  return removed;
}

optional<RuleToProcess> SessionState::remove_gy_rule(
    const std::string& rule_id, PolicyRule* rule_out,
    SessionStateUpdateCriteria& session_uc) {
  PolicyRule rule;
  bool removed = gy_dynamic_rules_.remove_rule(rule_id, &rule);
  if (!removed) {
    return {};
  }
  if (rule_out != nullptr) {
    *rule_out = rule;
  }
  session_uc.gy_dynamic_rules_to_uninstall.insert(rule_id);

  increment_rule_stats(rule_id, session_uc);
  return make_rule_to_process(rule);
}

optional<RuleToProcess> SessionState::deactivate_static_rule(
    const std::string& rule_id, SessionStateUpdateCriteria& session_uc) {
  auto it = std::find(
      active_static_rules_.begin(), active_static_rules_.end(), rule_id);
  if (it == active_static_rules_.end()) {
    return {};
  }

  session_uc.static_rules_to_uninstall.insert(rule_id);
  active_static_rules_.erase(it);

  increment_rule_stats(rule_id, session_uc);

  PolicyRule rule;
  if (!static_rules_.get_rule(rule_id, &rule)) {
    rule.set_id(rule_id);
  }
  return make_rule_to_process(rule);
}

RuleToProcess SessionState::make_rule_to_process(const PolicyRule& rule) {
  RuleToProcess to_process;
  to_process.version = get_current_rule_version(rule.id());
  to_process.rule    = rule;

  // At this point, we know the rule exists, so just check if it exists in the
  // static rule store or not
  bool is_static    = static_rules_.get_rule(rule.id(), nullptr);
  PolicyType p_type = is_static ? STATIC : DYNAMIC;

  // If there is a dedicated bearer TEID already in map, use it
  PolicyID policy_id = PolicyID(p_type, rule.id());
  bool dedicated_bearer_exists =
      bearer_id_by_policy_.find(policy_id) != bearer_id_by_policy_.end();
  if (dedicated_bearer_exists) {
    to_process.teids = bearer_id_by_policy_[policy_id].teids;
  } else {
    to_process.teids = config_.common_context.teids();
  }

  return to_process;
}

bool SessionState::deactivate_scheduled_static_rule(
    const std::string& rule_id) {
  if (scheduled_static_rules_.count(rule_id) == 0) {
    return false;
  }
  scheduled_static_rules_.erase(rule_id);
  return true;
}

void SessionState::classify_policy_activation(
    const RuleToProcess& to_process, const PolicyType p_type,
    RulesToProcess* pending_activation, RulesToProcess* pending_bearer_setup) {
  if (policy_needs_bearer_creation(p_type, to_process.rule.id())) {
    pending_bearer_setup->push_back(to_process);
  } else {
    pending_activation->push_back(to_process);
  }
}

void SessionState::process_rules_to_install(
    const std::vector<StaticRuleInstall>& static_rule_installs,
    const std::vector<DynamicRuleInstall>& dynamic_rule_installs,
    RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
    RulesToProcess* pending_bearer_setup, RulesToSchedule* pending_scheduling,
    SessionStateUpdateCriteria* session_uc) {
  process_static_rule_installs(
      static_rule_installs, pending_activation, pending_deactivation,
      pending_bearer_setup, pending_scheduling, session_uc);
  process_dynamic_rule_installs(
      dynamic_rule_installs, pending_activation, pending_deactivation,
      pending_bearer_setup, pending_scheduling, session_uc);
}

void SessionState::process_static_rule_installs(
    const std::vector<StaticRuleInstall>& rule_installs,
    RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
    RulesToProcess* pending_bearer_setup, RulesToSchedule* pending_scheduling,
    SessionStateUpdateCriteria* session_uc) {
  std::time_t current_time = std::time(nullptr);
  for (const StaticRuleInstall& rule_install : rule_installs) {
    const std::string& rule_id = rule_install.rule_id();
    RuleLifetime lifetime(rule_install);
    if (is_static_rule_installed(rule_id)) {
      // Session proxy may ask for duplicate rule installs.
      // Ignore them here.
      MLOG(MWARNING) << "Ignoring static rule install for " << session_id_
                     << " for rule " << rule_id
                     << " since it is alreday installed";
      continue;
    }
    PolicyRule static_rule;
    if (!static_rules_.get_rule(rule_id, &static_rule)) {
      MLOG(MERROR) << "static rule " << rule_id
                   << " is not found, skipping install...";
      continue;
    }

    // If the rule should be deactivated already, deactivate just in case and
    // continue
    if (lifetime.exceeded_lifetime(current_time)) {
      optional<RuleToProcess> op_remove_info =
          deactivate_static_rule(rule_id, *session_uc);
      if (op_remove_info) {
        pending_deactivation->push_back(*op_remove_info);
      }
      continue;
    }
    // If the rule should be active now, install
    if (lifetime.is_within_lifetime(current_time)) {
      RuleToProcess to_process =
          activate_static_rule(rule_id, lifetime, *session_uc);
      classify_policy_activation(
          to_process, STATIC, pending_activation, pending_bearer_setup);
    }
    // If the rule is for future activation, schedule
    if (lifetime.before_lifetime(current_time)) {
      schedule_static_rule(rule_id, lifetime, *session_uc);
      pending_scheduling->push_back(
          RuleToSchedule(STATIC, rule_id, ACTIVATE, lifetime.activation_time));
    }
    // Schedule deactivation time in the future
    if (lifetime.should_schedule_deactivation(current_time)) {
      pending_scheduling->push_back(RuleToSchedule(
          STATIC, rule_id, DEACTIVATE, lifetime.deactivation_time));
    }
  }
}

void SessionState::process_dynamic_rule_installs(
    const std::vector<DynamicRuleInstall>& rule_installs,
    RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
    RulesToProcess* pending_bearer_setup, RulesToSchedule* pending_scheduling,
    SessionStateUpdateCriteria* session_uc) {
  std::time_t current_time = std::time(nullptr);

  for (const DynamicRuleInstall& rule_install : rule_installs) {
    const PolicyRule& dynamic_rule = rule_install.policy_rule();
    const std::string& rule_id     = dynamic_rule.id();
    RuleLifetime lifetime(rule_install);

    // If the rule should be deactivated already, deactivate just in case and
    // continue
    if (lifetime.exceeded_lifetime(current_time)) {
      optional<RuleToProcess> op_remove_info =
          remove_dynamic_rule(rule_id, nullptr, *session_uc);
      if (op_remove_info) {
        pending_deactivation->push_back(*op_remove_info);
      }
      continue;
    }
    // If the rule should be active now, install
    if (lifetime.is_within_lifetime(current_time)) {
      RuleToProcess to_process =
          insert_dynamic_rule(dynamic_rule, lifetime, *session_uc);
      classify_policy_activation(
          to_process, DYNAMIC, pending_activation, pending_bearer_setup);
    }
    // If the rule is for future activation, schedule
    if (lifetime.before_lifetime(current_time)) {
      schedule_dynamic_rule(dynamic_rule, lifetime, *session_uc);
      pending_scheduling->push_back(
          RuleToSchedule(DYNAMIC, rule_id, ACTIVATE, lifetime.activation_time));
    }
    // Schedule deactivation time in the future
    if (lifetime.should_schedule_deactivation(current_time)) {
      pending_scheduling->push_back(RuleToSchedule(
          DYNAMIC, rule_id, DEACTIVATE, lifetime.deactivation_time));
    }
  }
}

void SessionState::process_rules_to_remove(
    const google::protobuf::RepeatedPtrField<std::basic_string<char>>
        rules_to_remove,
    RulesToProcess* pending_deactivation,
    SessionStateUpdateCriteria* session_uc) {
  for (const auto& rule_id : rules_to_remove) {
    optional<PolicyType> p_type = get_policy_type(rule_id);
    if (!p_type) {
      MLOG(MWARNING) << "Could not find rule " << rule_id << " for "
                     << session_id_ << " during static rule removal";
      continue;
    }
    optional<RuleToProcess> remove_info = {};
    PolicyRule rule;
    switch (*p_type) {
      case DYNAMIC: {
        remove_info = remove_dynamic_rule(rule_id, &rule, *session_uc);
        break;
      }
      case STATIC: {
        if (static_rules_.get_rule(rule_id, &rule)) {
          remove_info = deactivate_static_rule(rule_id, *session_uc);
        }
        break;
      }
      default:
        break;
    }
    if (!remove_info) {
      MLOG(MERROR) << "Failed to remove " << rule_id << " for " << session_id_;
      continue;
    }
    pending_deactivation->push_back(*remove_info);
  }
}

void SessionState::sync_rules_to_time(
    std::time_t current_time, SessionStateUpdateCriteria& session_uc) {
  // Update active static rules
  for (const std::string& rule_id : active_static_rules_) {
    if (should_rule_be_deactivated(rule_id, current_time)) {
      deactivate_static_rule(rule_id, session_uc);
    }
  }
  // Update scheduled static rules
  std::set<std::string> scheduled_rule_ids = scheduled_static_rules_;
  for (const std::string& rule_id : scheduled_rule_ids) {
    if (should_rule_be_active(rule_id, current_time)) {
      scheduled_static_rules_.erase(rule_id);
      activate_static_rule(rule_id, rule_lifetimes_[rule_id], session_uc);
    } else if (should_rule_be_deactivated(rule_id, current_time)) {
      scheduled_static_rules_.erase(rule_id);
      deactivate_static_rule(rule_id, session_uc);
    }
  }
  // Update active dynamic rules
  std::vector<std::string> dynamic_rule_ids;
  dynamic_rules_.get_rule_ids(dynamic_rule_ids);
  for (const std::string& rule_id : dynamic_rule_ids) {
    if (should_rule_be_deactivated(rule_id, current_time)) {
      remove_dynamic_rule(rule_id, NULL, session_uc);
    }
  }
  // Update scheduled dynamic rules
  dynamic_rule_ids.clear();
  scheduled_dynamic_rules_.get_rule_ids(dynamic_rule_ids);
  for (const std::string& rule_id : dynamic_rule_ids) {
    if (should_rule_be_active(rule_id, current_time)) {
      PolicyRule dy_rule;
      remove_scheduled_dynamic_rule(rule_id, &dy_rule, session_uc);
      insert_dynamic_rule(dy_rule, rule_lifetimes_[rule_id], session_uc);
    } else if (should_rule_be_deactivated(rule_id, current_time)) {
      remove_scheduled_dynamic_rule(rule_id, NULL, session_uc);
    }
  }
}

std::vector<std::string>& SessionState::get_static_rules() {
  return active_static_rules_;
}

std::set<std::string>& SessionState::get_scheduled_static_rules() {
  return scheduled_static_rules_;
}

DynamicRuleStore& SessionState::get_dynamic_rules() {
  return dynamic_rules_;
}

DynamicRuleStore& SessionState::get_scheduled_dynamic_rules() {
  return scheduled_dynamic_rules_;
}

RuleLifetime& SessionState::get_rule_lifetime(const std::string& rule_id) {
  return rule_lifetimes_[rule_id];
}

DynamicRuleStore& SessionState::get_gy_dynamic_rules() {
  return gy_dynamic_rules_;
}

uint32_t SessionState::total_monitored_rules_count() {
  uint32_t monitored_dynamic_rules = dynamic_rules_.monitored_rules_count();
  uint32_t monitored_static_rules  = 0;
  for (auto& rule_id : active_static_rules_) {
    if (static_rules_.get_monitoring_key_for_rule_id(rule_id, nullptr)) {
      monitored_static_rules++;
    }
  }
  return monitored_dynamic_rules + monitored_static_rules;
}

void SessionState::schedule_dynamic_rule(
    const PolicyRule& rule, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.new_rule_lifetimes[rule.id()] = lifetime;
  update_criteria.new_scheduled_dynamic_rules.push_back(rule);
  rule_lifetimes_[rule.id()] = lifetime;
  scheduled_dynamic_rules_.insert_rule(rule);
}

void SessionState::schedule_static_rule(
    const std::string& rule_id, RuleLifetime& lifetime,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.new_rule_lifetimes[rule_id] = lifetime;
  update_criteria.new_scheduled_static_rules.insert(rule_id);
  rule_lifetimes_[rule_id] = lifetime;
  scheduled_static_rules_.insert(rule_id);
}

uint32_t SessionState::get_credit_key_count() {
  return credit_map_.size() + monitor_map_.size();
}

bool SessionState::is_active() {
  return curr_state_ == SESSION_ACTIVE;
}

void SessionState::set_fsm_state(
    SessionFsmState new_state, SessionStateUpdateCriteria& uc) {
  // Only log and reflect change into update criteria if the state is new
  if (curr_state_ != new_state) {
    MLOG(MDEBUG) << "Session " << session_id_ << " Teid " << local_teid_
                 << " FSM state change from "
                 << session_fsm_state_to_str(curr_state_) << " to "
                 << session_fsm_state_to_str(new_state);
    curr_state_          = new_state;
    uc.is_fsm_updated    = true;
    uc.updated_fsm_state = new_state;
  }
}

// Suspend the service due to all the remaining credits are transient.
// Use the rg to trigger redirection
void SessionState::suspend_service_if_needed_for_credit(
    CreditKey ckey, SessionStateUpdateCriteria& update_criteria) {
  uint suspended_count = 0;

  auto it = credit_map_.find(ckey);
  if (it == credit_map_.end()) {
    MLOG(MDEBUG) << "Could not find RG " << ckey
                 << " Not suspending service for " << session_id_;
  }
  for (const auto& credit : credit_map_) {
    if (credit.second->suspended) {
      suspended_count++;
    }
  }
  if (credit_map_.size() > 0 && suspended_count == credit_map_.size()) {
    auto credit_uc = get_credit_uc(ckey, update_criteria);
    it->second->set_service_state(SERVICE_NEEDS_SUSPENSION, *credit_uc);
  }
}

bool SessionState::should_rule_be_active(
    const std::string& rule_id, std::time_t time) {
  return rule_lifetimes_[rule_id].is_within_lifetime(time);
}

bool SessionState::should_rule_be_deactivated(
    const std::string& rule_id, std::time_t time) {
  return rule_lifetimes_[rule_id].exceeded_lifetime(time);
}

StaticRuleInstall SessionState::get_static_rule_install(
    const std::string& rule_id, const RuleLifetime& lifetime) {
  StaticRuleInstall rule_install{};
  rule_install.set_rule_id(rule_id);
  rule_install.mutable_activation_time()->set_seconds(lifetime.activation_time);
  rule_install.mutable_deactivation_time()->set_seconds(
      lifetime.deactivation_time);
  return rule_install;
}

DynamicRuleInstall SessionState::get_dynamic_rule_install(
    const std::string& rule_id, const RuleLifetime& lifetime) {
  DynamicRuleInstall rule_install{};
  PolicyRule* policy_rule = rule_install.mutable_policy_rule();
  if (!dynamic_rules_.get_rule(rule_id, policy_rule)) {
    scheduled_dynamic_rules_.get_rule(rule_id, policy_rule);
  }
  rule_install.mutable_activation_time()->set_seconds(lifetime.activation_time);
  rule_install.mutable_deactivation_time()->set_seconds(
      lifetime.deactivation_time);
  return rule_install;
}

// Charging Credits
static FinalActionInfo get_final_action_info(
    const magma::lte::ChargingCredit& credit) {
  FinalActionInfo final_action_info;
  if (credit.is_final()) {
    final_action_info.final_action = credit.final_action();
    switch (final_action_info.final_action) {
      case ChargingCredit_FinalAction_REDIRECT:
        final_action_info.redirect_server = credit.redirect_server();
        break;
      case ChargingCredit_FinalAction_RESTRICT_ACCESS:
        for (auto rule : credit.restrict_rules()) {
          final_action_info.restrict_rules.push_back(rule);
        }
        break;
      default:  // do nothing;
        break;
    }
  }
  return final_action_info;
}

std::vector<PolicyRule> SessionState::get_all_final_unit_rules() {
  std::vector<PolicyRule> rules;
  for (auto& credit_pair : credit_map_) {
    auto& grant = credit_pair.second;
    if (grant->service_state != SERVICE_RESTRICTED) {
      continue;
    }
    for (const std::string rule_id : grant->final_action_info.restrict_rules) {
      PolicyRule rule;
      if (static_rules_.get_rule(rule_id, &rule)) {
        rules.push_back(rule);
      }
    }
  }
  gy_dynamic_rules_.get_rules(rules);
  return rules;
}

void SessionState::handle_update_failure(
    const UpdateRequests& failed_requests,
    SessionStateUpdateCriteria& session_uc) {
  MLOG(MDEBUG) << "Rolling back changes due to failed updates ("
               << failed_requests.charging_requests.size()
               << " charging requests and "
               << failed_requests.monitor_requests.size()
               << " monitor requests) for " << session_id_;
  for (const auto& failed_charging : failed_requests.charging_requests) {
    const auto key = failed_charging.usage().charging_key();
    if (credit_map_.find(key) == credit_map_.end()) {
      MLOG(MERROR) << "Could not find credit RG:" << key << " to reset for "
                   << session_id_;
      continue;
    }
    credit_map_[key]->reset_reporting_grant(get_credit_uc(key, session_uc));
  }
  for (const auto& failed_monitor : failed_requests.monitor_requests) {
    const auto key = failed_monitor.update().monitoring_key();
    if (monitor_map_.find(key) == monitor_map_.end()) {
      MLOG(MERROR) << "Could not find monitor:" << key << " to reset for "
                   << session_id_;
      continue;
    }
    monitor_map_[key]->credit.reset_reporting_credit(
        get_monitor_uc(key, session_uc));
  }
}

bool SessionState::receive_charging_credit(
    const CreditUpdateResponse& update,
    SessionStateUpdateCriteria& session_uc) {
  auto key = CreditKey(update);

  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    // new credit
    return init_charging_credit(update, session_uc);
  }
  auto& grant          = it->second;
  auto credit_uc       = get_credit_uc(key, session_uc);
  auto credit_validity = ChargingGrant::is_valid_credit_response(update);
  if (credit_validity == INVALID_CREDIT) {
    // update unsuccessful, reset credit and return
    grant->credit.mark_failure(update.result_code(), credit_uc);
    if (grant->should_deactivate_service()) {
      grant->set_service_state(SERVICE_NEEDS_DEACTIVATION, *credit_uc);
    }
    return false;
  }
  if (credit_validity == TRANSIENT_ERROR) {
    // for transient errors, try to install the credit
    // but clear the reported credit
    grant->credit.mark_failure(update.result_code(), credit_uc);
  }
  grant->receive_charging_grant(update, credit_uc);

  if (grant->reauth_state == REAUTH_PROCESSING) {
    grant->set_reauth_state(REAUTH_NOT_NEEDED, *credit_uc);
  }
  if (!grant->credit.is_quota_exhausted(1) &&
      grant->service_state != SERVICE_ENABLED) {
    // if quota no longer exhausted, re-enable services as needed
    MLOG(MINFO) << "New quota now available. Service is in state: "
                << service_state_to_str(grant->service_state)
                << " Activating service RG: " << key << " for " << session_id_;
    grant->set_service_state(SERVICE_NEEDS_ACTIVATION, *credit_uc);
  }
  return true;
}

bool SessionState::init_charging_credit(
    const CreditUpdateResponse& update,
    SessionStateUpdateCriteria& session_uc) {
  const uint32_t key = update.charging_key();
  if (ChargingGrant::is_valid_credit_response(update) == INVALID_CREDIT) {
    // init failed, don't track key
    return false;
  }
  ChargingGrant charging_grant;
  charging_grant.credit = SessionCredit(SERVICE_ENABLED, update.limit_type());
  charging_grant.receive_charging_grant(update);
  session_uc.charging_credit_to_install[CreditKey(update)] =
      charging_grant.marshal();
  credit_map_[CreditKey(update)] =
      std::make_unique<ChargingGrant>(charging_grant);
  MLOG(MINFO) << "Initialized a new credit RG:" << key << " for "
              << session_id_;
  return true;
}

void SessionState::set_suspend_credit(
    const CreditKey& charging_key, bool new_suspended,
    SessionStateUpdateCriteria& update_criteria) {
  auto it = credit_map_.find(charging_key);
  if (it != credit_map_.end()) {
    auto credit_uc = get_credit_uc(charging_key, update_criteria);
    auto& grant    = it->second;
    grant->set_suspended(new_suspended, credit_uc);
  }
}

bool SessionState::is_credit_suspended(const CreditKey& charging_key) {
  auto it = credit_map_.find(charging_key);
  if (it != credit_map_.end()) {
    auto& grant = it->second;
    return grant->get_suspended();
  }
  return false;
}

void SessionState::get_rules_per_credit_key(
    const CreditKey& charging_key, RulesToProcess* to_process,
    SessionStateUpdateCriteria* session_uc) {
  std::vector<PolicyRule> static_rules, dynamic_rules;
  static_rules_.get_rule_definitions_for_charging_key(
      charging_key, static_rules);
  for (PolicyRule rule : static_rules) {
    // Since the static rule store is shared across sessions, we should check
    // that the rule is activated for the session
    bool is_installed = is_static_rule_installed(rule.id());
    if (is_installed) {
      increment_rule_stats(rule.id(), *session_uc);
      to_process->push_back(make_rule_to_process(rule));
    }
  }
  dynamic_rules_.get_rule_definitions_for_charging_key(
      charging_key, dynamic_rules);
  for (PolicyRule rule : dynamic_rules) {
    increment_rule_stats(rule.id(), *session_uc);
    to_process->push_back(make_rule_to_process(rule));
  }
}

uint64_t SessionState::get_charging_credit(
    const CreditKey& key, Bucket bucket) const {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return 0;
  }
  return it->second->credit.get_credit(bucket);
}

bool SessionState::set_credit_reporting(
    const CreditKey& key, bool reporting,
    SessionStateUpdateCriteria* session_uc) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    MLOG(MWARNING) << "Did not set reporting flag for RG:" << key << " for "
                   << session_id_;
    return false;
  }

  it->second->credit.set_reporting(reporting);
  if (session_uc != NULL) {
    auto credit_uc       = get_credit_uc(key, *session_uc);
    credit_uc->reporting = reporting;
  }
  return true;
}

ReAuthResult SessionState::reauth_key(
    const CreditKey& charging_key,
    SessionStateUpdateCriteria& update_criteria) {
  auto it = credit_map_.find(charging_key);
  if (it != credit_map_.end()) {
    // if credit is already reporting, don't initiate update
    auto& grant = it->second;
    if (grant->credit.is_reporting()) {
      return ReAuthResult::UPDATE_NOT_NEEDED;
    }
    auto uc = grant->get_update_criteria();
    grant->set_reauth_state(REAUTH_REQUIRED, uc);
    update_criteria.charging_credit_map[charging_key] = uc;
    return ReAuthResult::UPDATE_INITIATED;
  }
  // charging_key cannot be found, initialize credit and engage reauth
  auto grant           = std::make_unique<ChargingGrant>();
  grant->credit        = SessionCredit(SERVICE_DISABLED);
  grant->reauth_state  = REAUTH_REQUIRED;
  grant->service_state = SERVICE_DISABLED;
  update_criteria.charging_credit_to_install[charging_key] = grant->marshal();
  credit_map_[charging_key]                                = std::move(grant);
  return ReAuthResult::UPDATE_INITIATED;
}

ReAuthResult SessionState::reauth_all(
    SessionStateUpdateCriteria& update_criteria) {
  auto res = ReAuthResult::UPDATE_NOT_NEEDED;
  for (auto& credit_pair : credit_map_) {
    auto key    = credit_pair.first;
    auto& grant = credit_pair.second;
    // Only update credits that aren't reporting
    if (!grant->credit.is_reporting()) {
      update_criteria.charging_credit_map[key] = grant->get_update_criteria();
      grant->set_reauth_state(
          REAUTH_REQUIRED, update_criteria.charging_credit_map[key]);
      res = ReAuthResult::UPDATE_INITIATED;
    }
  }
  return res;
}

void SessionState::apply_charging_credit_update(
    const CreditKey& key, SessionCreditUpdateCriteria& credit_uc) {
  auto it = credit_map_.find(key);
  if (it == credit_map_.end()) {
    return;
  }

  if (credit_uc.deleted) {
    credit_map_.erase(key);
    MLOG(MINFO) << session_id_ << " Erasing RG " << key;
    return;
  }

  auto& charging_grant = it->second;
  auto& credit         = charging_grant->credit;

  // Credit merging
  credit.merge(credit_uc);

  // set charging grant
  charging_grant->is_final_grant    = credit_uc.is_final;
  charging_grant->final_action_info = credit_uc.final_action_info;
  charging_grant->expiry_time       = credit_uc.expiry_time;
  charging_grant->reauth_state      = credit_uc.reauth_state;
  charging_grant->service_state     = credit_uc.service_state;
  charging_grant->suspended         = credit_uc.suspended;
}

void SessionState::set_charging_credit(
    const CreditKey& key, ChargingGrant charging_grant,
    SessionStateUpdateCriteria& uc) {
  credit_map_[key] = std::make_unique<ChargingGrant>(charging_grant);
  uc.charging_credit_to_install[key] = credit_map_[key]->marshal();
}

CreditUsageUpdate SessionState::make_credit_usage_update_req(
    CreditUsage& usage) const {
  CreditUsageUpdate req;
  req.set_session_id(session_id_);
  req.set_request_number(request_number_);
  fill_protos_tgpp_context(req.mutable_tgpp_ctx());
  req.mutable_common_context()->CopyFrom(config_.common_context);

  // TODO keep RAT specific fields separate for now as we may not always want
  // to send the entire context
  if (config_.rat_specific_context.has_lte_context()) {
    const auto& lte_context = config_.rat_specific_context.lte_context();
    req.set_spgw_ipv4(lte_context.spgw_ipv4());
    req.set_imei(lte_context.imei());
    req.set_plmn_id(lte_context.plmn_id());
    req.set_imsi_plmn_id(lte_context.imsi_plmn_id());
    req.set_user_location(lte_context.user_location());
    req.set_charging_characteristics(lte_context.charging_characteristics());
  } else if (config_.rat_specific_context.has_wlan_context()) {
    const auto& wlan_context = config_.rat_specific_context.wlan_context();
    req.set_hardware_addr(wlan_context.mac_addr_binary());
  }
  req.mutable_usage()->CopyFrom(usage);
  return req;
}

void SessionState::get_charging_updates(
    UpdateSessionRequest& update_request_out,
    std::vector<std::unique_ptr<ServiceAction>>* actions_out,
    SessionStateUpdateCriteria& uc) {
  for (auto& credit_pair : credit_map_) {
    auto& key      = credit_pair.first;
    auto& grant    = credit_pair.second;
    auto credit_uc = get_credit_uc(key, uc);

    auto action_type = grant->get_action(*credit_uc);
    auto action      = std::make_unique<ServiceAction>(action_type);
    switch (action_type) {
      case CONTINUE_SERVICE: {
        optional<CreditUsageUpdate> op_update =
            get_update_for_continue_service(key, grant, uc);
        if (!op_update) {
          // no update
          break;
        }
        update_request_out.mutable_updates()->Add()->CopyFrom(*op_update);
      } break;
      case REDIRECT: {
        if (grant->service_state == SERVICE_REDIRECTED) {
          MLOG(MDEBUG) << "Redirection already activated for " << session_id_;
          continue;
        }
        grant->set_service_state(SERVICE_REDIRECTED, *credit_uc);

        PolicyRule redirect_rule = make_redirect_rule(grant);
        if (!is_gy_dynamic_rule_installed(redirect_rule.id())) {
          fill_service_action_for_redirect(
              action, key, grant, redirect_rule, uc);
          actions_out->push_back(std::move(action));
        }

        break;
      }
      case RESTRICT_ACCESS: {
        if (grant->service_state == SERVICE_RESTRICTED) {
          MLOG(MDEBUG) << "Restriction already activated for " << session_id_;
          continue;
        }
        grant->set_service_state(SERVICE_RESTRICTED, *credit_uc);

        fill_service_action_for_restrict(action, key, grant, uc);
        actions_out->push_back(std::move(action));
        break;
      }
      case ACTIVATE_SERVICE:
        fill_service_action_for_activate(action, key, uc);
        actions_out->push_back(std::move(action));
        grant->set_suspended(false, credit_uc);
        break;
      case TERMINATE_SERVICE:
        fill_service_action_with_context(action, action_type, key);
        actions_out->push_back(std::move(action));
        break;
      default:
        MLOG(MWARNING) << "Unexpected action type "
                       << service_action_type_to_str(action_type) << " for "
                       << session_id_;
        break;
    }
  }
}

optional<CreditUsageUpdate> SessionState::get_update_for_continue_service(
    const CreditKey& key, std::unique_ptr<ChargingGrant>& grant,
    SessionStateUpdateCriteria& session_uc) {
  CreditUsage::UpdateType update_type;
  if (!grant->get_update_type(&update_type)) {
    return {};  // no update
  }
  if (curr_state_ == SESSION_RELEASED) {
    MLOG(MDEBUG) << "Session " << session_id_
                 << " is in Released state. Not sending update to the core"
                    "for rating group "
                 << key;
    return {};  // no update
  }
  if (grant->suspended && update_type == CreditUsage::QUOTA_EXHAUSTED) {
    MLOG(MDEBUG) << "Credit " << key << " for " << session_id_
                 << " is suspended. Not sending update to the core";
    return {};  // no update
  }

  // Create Update struct
  MLOG(MDEBUG) << "Subscriber " << get_imsi() << " rating group " << key
               << " updating due to type "
               << credit_update_type_to_str(update_type)
               << " with request number " << request_number_;

  auto credit_uc = get_credit_uc(key, session_uc);
  if (update_type == CreditUsage::REAUTH_REQUIRED) {
    grant->set_reauth_state(REAUTH_PROCESSING, *credit_uc);
  }
  CreditUsage usage = grant->get_credit_usage(update_type, *credit_uc, false);
  key.set_credit_usage(&usage);

  auto request = make_credit_usage_update_req(usage);
  request_number_++;
  session_uc.request_number_increment++;
  return request;
}

void SessionState::fill_service_action_for_activate(
    std::unique_ptr<ServiceAction>& action_p, const CreditKey& key,
    SessionStateUpdateCriteria& session_uc) {
  std::vector<PolicyRule> static_rules, dynamic_rules;
  fill_service_action_with_context(action_p, ACTIVATE_SERVICE, key);
  static_rules_.get_rules_by_ids(active_static_rules_, static_rules);
  dynamic_rules_.get_rule_definitions_for_charging_key(key, dynamic_rules);

  RulesToProcess* to_install = action_p->get_mutable_gx_rules_to_install();
  for (PolicyRule rule : static_rules) {
    RuleLifetime lifetime;
    to_install->push_back(
        activate_static_rule(rule.id(), lifetime, session_uc));
  }
  for (PolicyRule rule : dynamic_rules) {
    RuleLifetime lifetime;
    to_install->push_back(insert_dynamic_rule(rule, lifetime, session_uc));
  }
}

void SessionState::fill_service_action_for_restrict(
    std::unique_ptr<ServiceAction>& action_p, const CreditKey& key,
    std::unique_ptr<ChargingGrant>& grant,
    SessionStateUpdateCriteria& session_uc) {
  fill_service_action_with_context(action_p, RESTRICT_ACCESS, key);

  RulesToProcess* gy_to_install = action_p->get_mutable_gy_rules_to_install();
  for (auto& rule_id : grant->final_action_info.restrict_rules) {
    PolicyRule rule;
    if (!static_rules_.get_rule(rule_id, &rule)) {
      MLOG(MWARNING) << "Static rule " << rule_id
                     << " requested as a restrict rule is not found.";
      continue;
    }
    RuleLifetime lifetime;
    gy_to_install->push_back(insert_gy_rule(rule, lifetime, session_uc));
  }
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
      MLOG(MERROR) << "Unknown redirect address type!";
      return RedirectInformation_AddressType_IPv4;
  }
}

PolicyRule SessionState::make_redirect_rule(
    std::unique_ptr<ChargingGrant>& grant) {
  PolicyRule redirect_rule;
  redirect_rule.set_id("redirect");
  redirect_rule.set_priority(SessionState::REDIRECT_FLOW_PRIORITY);
  RedirectInformation* redirect_info = redirect_rule.mutable_redirect();
  redirect_info->set_support(RedirectInformation_Support_ENABLED);

  auto redirect_server = grant->final_action_info.redirect_server;
  redirect_info->set_address_type(
      address_type_converter(redirect_server.redirect_address_type()));
  redirect_info->set_server_address(redirect_server.redirect_server_address());
  return redirect_rule;
}

void SessionState::fill_service_action_for_redirect(
    std::unique_ptr<ServiceAction>& action_p, const CreditKey& key,
    std::unique_ptr<ChargingGrant>& grant, PolicyRule redirect_rule,
    SessionStateUpdateCriteria& session_uc) {
  fill_service_action_with_context(action_p, REDIRECT, key);

  RulesToProcess* gy_to_install = action_p->get_mutable_gy_rules_to_install();
  RuleLifetime lifetime;
  gy_to_install->push_back(insert_gy_rule(redirect_rule, lifetime, session_uc));
}

void SessionState::fill_service_action_with_context(
    std::unique_ptr<ServiceAction>& action, ServiceActionType action_type,
    const CreditKey& key) {
  MLOG(MDEBUG) << "Subscriber " << get_imsi() << " rating group " << key
               << " action type " << service_action_type_to_str(action_type);
  action->set_credit_key(key);
  action->set_imsi(get_imsi());
  action->set_ambr(config_.get_apn_ambr());
  action->set_ip_addr(config_.common_context.ue_ipv4());
  action->set_ipv6_addr(config_.common_context.ue_ipv6());
  action->set_teids(config_.common_context.teids());
  action->set_msisdn(config_.common_context.msisdn());
  action->set_session_id(session_id_);
}

// Monitors
bool SessionState::receive_monitor(
    const UsageMonitoringUpdateResponse& update,
    SessionStateUpdateCriteria& session_uc) {
  if (!update.has_credit()) {
    // We are overloading UsageMonitoringUpdateResponse/Request with other
    // EventTriggered requests, so we could receive updates that don't affect
    // UsageMonitors.
    MLOG(MINFO) << "Received a UsageMonitoringUpdateResponse without a monitor"
                << ", not creating a monitor.";
    return true;
  }
  if (update.success() &&
      update.credit().level() == MonitoringLevel::SESSION_LEVEL) {
    update_session_level_key(update, session_uc);
  }
  auto mkey = update.credit().monitoring_key();
  auto it   = monitor_map_.find(mkey);

  if (session_uc.monitor_credit_map.find(mkey) !=
          session_uc.monitor_credit_map.end() &&
      session_uc.monitor_credit_map[mkey].deleted) {
    // This will only happen if the PCRF responds back with more credit when
    // the monitor has already been set to be terminated
    MLOG(MDEBUG) << session_id_ << "Ignoring  update for monitor " << mkey
                 << " because it has been set for deletion";
    return false;
  }

  if (it == monitor_map_.end()) {
    // new credit
    return init_new_monitor(update, session_uc);
  }
  auto credit_uc = get_monitor_uc(mkey, session_uc);
  if (!update.success()) {
    it->second->credit.mark_failure(update.result_code(), credit_uc);
    return false;
  }

  if (update.credit().action() == UsageMonitoringCredit::FORCE) {
    MLOG(MINFO) << session_id_
                << " Received UsageMonitoringCredit::FORCE "
                   "(`AVP: Usage-Monitoring-Report`) instruction "
                   "not implemented. Will just continue for monitor  "
                << mkey;
  }

  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    // Disable mpnitor, no matter if we still have credit
    MLOG(MINFO) << session_id_ << " Received Disabled action for monitor "
                << mkey << ". Will remove monitor after update is sent";
    // seting last update will deleted monitor after the update is sent.
    it->second->credit.set_report_last_credit(true, *credit_uc);

  } else {
    MLOG(MINFO) << session_id_ << " Received monitor credit for " << mkey;
    const auto& gsu = update.credit().granted_units();
    it->second->credit.receive_credit(gsu, credit_uc);
  }

  return true;
}

void SessionState::apply_monitor_updates(
    const std::string& key, SessionStateUpdateCriteria& session_uc,
    SessionCreditUpdateCriteria& credit_uc) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return;
  }

  // Actual deletion of monitor
  if (credit_uc.deleted) {
    if (it->second->level == MonitoringLevel::SESSION_LEVEL) {
      // session level change
      MLOG(MINFO) << "Removing Session Level monitor " << key;
      session_uc.is_session_level_key_updated = true;
      session_uc.updated_session_level_key    = "";
    }
    MLOG(MINFO) << session_id_ << " Erasing monitor " << key;
    monitor_map_.erase(key);
    return;
  }

  auto& charging_grant = it->second;
  auto& credit         = charging_grant->credit;

  // Credit merging
  credit.merge(credit_uc);
}

uint64_t SessionState::get_monitor(
    const std::string& key, Bucket bucket) const {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    return 0;
  }
  return it->second->credit.get_credit(bucket);
}

bool SessionState::set_monitor_reporting(
    const std::string& key, bool reporting,
    SessionStateUpdateCriteria* update_criteria) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    MLOG(MWARNING) << "Didn't set reporting flag for monitor key " << key;
    return false;
  }

  it->second->credit.set_reporting(reporting);

  if (update_criteria != NULL) {
    auto mon_credit_uc       = get_monitor_uc(key, *update_criteria);
    mon_credit_uc->reporting = reporting;
  }
  return true;
}

bool SessionState::add_to_monitor(
    const std::string& key, uint64_t used_tx, uint64_t used_rx,
    SessionStateUpdateCriteria& uc) {
  auto it = monitor_map_.find(key);
  if (it == monitor_map_.end()) {
    MLOG(MDEBUG) << "Monitoring Key " << key
                 << " not found, not adding the usage";
    return false;
  }

  auto credit_uc = get_monitor_uc(key, uc);

  it->second->credit.add_used_credit(used_tx, used_rx, *credit_uc);

  // after adding usage we check if monitor is exhausted
  if (it->second->should_delete_monitor()) {
    MLOG(MINFO) << "Quota exhausted for monitor " << key
                << ". Will remove monitor after update is sent";
    it->second->credit.set_report_last_credit(true, *credit_uc);
  }
  return true;
}

void SessionState::set_monitor(
    const std::string& key, Monitor monitor,
    SessionStateUpdateCriteria& update_criteria) {
  update_criteria.monitor_credit_to_install[key] = monitor.marshal();
  monitor_map_[key] = std::make_unique<Monitor>(monitor);
}

bool SessionState::init_new_monitor(
    const UsageMonitoringUpdateResponse& update,
    SessionStateUpdateCriteria& update_criteria) {
  if (!update.success()) {
    MLOG(MERROR) << "Monitoring init failed for imsi " << get_imsi()
                 << " and monitoring key " << update.credit().monitoring_key();
    return false;
  }
  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    MLOG(MWARNING) << "Monitoring init has action disabled for subscriber "
                   << get_imsi() << " and monitoring key "
                   << update.credit().monitoring_key();
    return false;
  }
  MLOG(MDEBUG) << session_id_ << " Initialized a monitoring credit for mkey "
               << update.credit().monitoring_key();
  auto monitor   = std::make_unique<Monitor>();
  monitor->level = update.credit().level();
  // validity time and final units not used for monitors
  auto _   = SessionCreditUpdateCriteria{};
  auto gsu = update.credit().granted_units();
  monitor->credit.receive_credit(gsu, NULL);

  update_criteria.monitor_credit_to_install[update.credit().monitoring_key()] =
      monitor->marshal();
  monitor_map_[update.credit().monitoring_key()] = std::move(monitor);
  return true;
}

void SessionState::update_session_level_key(
    const UsageMonitoringUpdateResponse& update,
    SessionStateUpdateCriteria& uc) {
  const auto& new_key = update.credit().monitoring_key();
  if (session_level_key_ != "" && session_level_key_ != new_key) {
    MLOG(MINFO) << "Session level monitoring key is updated from "
                << session_level_key_ << " to " << new_key;
  }
  if (update.credit().action() == UsageMonitoringCredit::DISABLE) {
    session_level_key_ = "";
  } else {
    session_level_key_ = new_key;
  }
  uc.is_session_level_key_updated = true;
  uc.updated_session_level_key    = session_level_key_;
}

void SessionState::set_session_level_key(const std::string new_key) {
  session_level_key_ = new_key;
}

BearerUpdate SessionState::get_dedicated_bearer_updates(
    const RulesToProcess& pending_activation,
    const RulesToProcess& pending_deactivation,
    SessionStateUpdateCriteria* session_uc) {
  BearerUpdate update;
  // Rule Installs
  for (const auto& to_process : pending_activation) {
    const auto& rule_id = to_process.rule.id();
    PolicyType p_type;
    if (static_rules_.get_rule(rule_id, nullptr)) {
      p_type = STATIC;
    } else {
      p_type = DYNAMIC;
    }
    update_bearer_creation_req(p_type, rule_id, &update);
  }

  // Rule Removals
  for (const auto& to_process : pending_deactivation) {
    const auto& rule_id = to_process.rule.id();
    PolicyType p_type;
    if (static_rules_.get_rule(rule_id, nullptr)) {
      p_type = STATIC;
    } else {
      p_type = DYNAMIC;
    }
    update_bearer_deletion_req(p_type, rule_id, update, *session_uc);
  }
  return update;
}

void SessionState::bind_policy_to_bearer(
    const PolicyBearerBindingRequest& request, SessionStateUpdateCriteria& uc) {
  const std::string& rule_id = request.policy_rule_id();
  auto policy_type           = get_policy_type(rule_id);
  if (!policy_type) {
    MLOG(MDEBUG) << "Policy " << rule_id
                 << " not found, when trying to bind to bearerID "
                 << request.bearer_id();
    return;
  }
  MLOG(MINFO) << session_id_ << " now has policy " << rule_id
              << " tied to bearerID " << request.bearer_id();
  BearerIDAndTeid brearer_id_and_teid;
  brearer_id_and_teid.bearer_id                         = request.bearer_id();
  brearer_id_and_teid.teids                             = request.teids();
  bearer_id_by_policy_[PolicyID(*policy_type, rule_id)] = brearer_id_and_teid;
  uc.is_bearer_mapping_updated                          = true;
  uc.bearer_id_by_policy                                = bearer_id_by_policy_;
}

optional<PolicyType> SessionState::get_policy_type(const std::string& rule_id) {
  if (is_static_rule_installed(rule_id)) {
    return STATIC;
  } else if (is_dynamic_rule_installed(rule_id)) {
    return DYNAMIC;
  } else {
    return {};
  }
}

optional<PolicyRule> SessionState::get_policy_definition(
    const std::string rule_id) {
  optional<PolicyType> policy_type = get_policy_type(rule_id);
  if (!policy_type) {
    return {};
  }
  PolicyRule rule;
  if (*policy_type == STATIC && static_rules_.get_rule(rule_id, &rule)) {
    return rule;
  }
  if (*policy_type == DYNAMIC && dynamic_rules_.get_rule(rule_id, &rule)) {
    return rule;
  }
  return {};
}

SessionCreditUpdateCriteria* SessionState::get_monitor_uc(
    const std::string& key, SessionStateUpdateCriteria& uc) {
  if (uc.monitor_credit_map.find(key) == uc.monitor_credit_map.end()) {
    uc.monitor_credit_map[key] =
        monitor_map_[key]->credit.get_update_criteria();
  }
  return &(uc.monitor_credit_map[key]);
}

// Event Triggers
void SessionState::get_event_trigger_updates(
    UpdateSessionRequest& update_request_out,
    SessionStateUpdateCriteria& update_criteria) {
  // todo We should also handle other event triggers here too
  auto it = pending_event_triggers_.find(REVALIDATION_TIMEOUT);
  if (it != pending_event_triggers_.end() && it->second == READY) {
    MLOG(MDEBUG) << "Session " << session_id_
                 << " updating due to EventTrigger: REVALIDATION_TIMEOUT"
                 << " with request number " << request_number_;
    auto new_req = update_request_out.mutable_usage_monitors()->Add();
    add_common_fields_to_usage_monitor_update(new_req);
    new_req->set_event_trigger(REVALIDATION_TIMEOUT);
    request_number_++;
    update_criteria.request_number_increment++;
    // todo we might want to make sure that the update went successfully
    // before clearing here
    remove_event_trigger(REVALIDATION_TIMEOUT, update_criteria);
  }
}

void SessionState::add_new_event_trigger(
    magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria) {
  MLOG(MINFO) << "Event Trigger " << trigger << " is pending for "
              << session_id_;
  set_event_trigger(trigger, PENDING, update_criteria);
}

void SessionState::mark_event_trigger_as_triggered(
    magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria) {
  auto it = pending_event_triggers_.find(trigger);
  if (it == pending_event_triggers_.end() ||
      pending_event_triggers_[trigger] != PENDING) {
    MLOG(MWARNING) << "Event Trigger " << trigger << " requested to be "
                   << "triggered is not pending for " << session_id_;
  }
  MLOG(MINFO) << "Event Trigger " << trigger << " is ready to update for "
              << session_id_;
  set_event_trigger(trigger, READY, update_criteria);
}

void SessionState::remove_event_trigger(
    magma::lte::EventTrigger trigger,
    SessionStateUpdateCriteria& update_criteria) {
  MLOG(MINFO) << "Event Trigger " << trigger << " is removed for "
              << session_id_;
  pending_event_triggers_.erase(trigger);
  set_event_trigger(trigger, CLEARED, update_criteria);
}

void SessionState::set_event_trigger(
    magma::lte::EventTrigger trigger, const EventTriggerState value,
    SessionStateUpdateCriteria& update_criteria) {
  pending_event_triggers_[trigger]                  = value;
  update_criteria.is_pending_event_triggers_updated = true;
  update_criteria.pending_event_triggers[trigger]   = value;
}

void SessionState::set_revalidation_time(
    const google::protobuf::Timestamp& time,
    SessionStateUpdateCriteria& update_criteria) {
  revalidation_time_                = time;
  update_criteria.revalidation_time = time;
}

optional<FinalActionInfo> SessionState::get_final_action_if_final_unit_state(
    const CreditKey& charging_key) const {
  auto it = credit_map_.find(charging_key);
  if (it == credit_map_.end()) {
    return {};
  }
  if (it->second->service_state != SERVICE_REDIRECTED &&
      it->second->service_state != SERVICE_RESTRICTED) {
    return {};
  }
  return it->second->final_action_info;
}

RulesToProcess SessionState::remove_all_final_action_rules(
    const FinalActionInfo& final_action_info,
    SessionStateUpdateCriteria& session_uc) {
  RulesToProcess to_process;
  to_process = std::vector<RuleToProcess>{};
  switch (final_action_info.final_action) {
    case ChargingCredit_FinalAction_REDIRECT: {
      PolicyRule rule;
      optional<RuleToProcess> op_rule_info =
          remove_gy_rule("redirect", &rule, session_uc);
      if (op_rule_info) {
        to_process.push_back(*op_rule_info);
      }
    } break;
    case ChargingCredit_FinalAction_RESTRICT_ACCESS:
      for (std::string rule_id : final_action_info.restrict_rules) {
        PolicyRule rule;
        optional<RuleToProcess> op_rule_info =
            remove_gy_rule(rule_id, &rule, session_uc);
        if (op_rule_info) {
          to_process.push_back(*op_rule_info);
        }
      }
      break;
    default:
      break;
  }
  return to_process;
}

// QoS/Bearer Management
bool SessionState::policy_has_qos(
    const PolicyType policy_type, const std::string& rule_id,
    PolicyRule* rule_out) {
  if (policy_type == STATIC) {
    bool exists = static_rules_.get_rule(rule_id, rule_out);
    return exists && rule_out->has_qos();
  }
  if (policy_type == DYNAMIC) {
    bool exists = dynamic_rules_.get_rule(rule_id, rule_out);
    return exists && rule_out->has_qos();
  }
  return false;
}

optional<PolicyRule> SessionState::policy_needs_bearer_creation(
    const PolicyType policy_type, const std::string& id) {
  if (!config_.rat_specific_context.has_lte_context()) {
    return {};
  }
  if (bearer_id_by_policy_.find(PolicyID(policy_type, id)) !=
      bearer_id_by_policy_.end()) {
    // Policy already has a bearer
    return {};
  }
  PolicyRule policy;
  if (!policy_has_qos(policy_type, id, &policy)) {
    // Only create a bearer for policies with QoS
    return {};
  }
  auto default_qci = FlowQos_Qci(
      config_.rat_specific_context.lte_context().qos_info().qos_class_id());
  if (policy.qos().qci() == default_qci) {
    // This QCI is already covered by the default bearer
    return {};
  }
  return policy;
}

void SessionState::update_bearer_creation_req(
    const PolicyType policy_type, const std::string& rule_id,
    BearerUpdate* update) {
  auto policy = policy_needs_bearer_creation(policy_type, rule_id);
  if (!policy) {
    return;
  }

  // If it is first time filling in the CreationReq, fill in other info
  if (!update->needs_creation) {
    update->needs_creation = true;
    update->create_req.mutable_sid()->CopyFrom(config_.common_context.sid());
    update->create_req.set_ip_addr(config_.common_context.ue_ipv4());
    // TODO ipv6 add to the bearer request or remove ipv4
    update->create_req.set_link_bearer_id(
        config_.rat_specific_context.lte_context().bearer_id());
  }
  update->create_req.mutable_policy_rules()->Add()->CopyFrom(*policy);
  // We will add the new policyID to bearerID association, once we receive a
  // message from SGW.
}

void SessionState::update_bearer_deletion_req(
    const PolicyType policy_type, const std::string& rule_id,
    BearerUpdate& update, SessionStateUpdateCriteria& uc) {
  if (!config_.rat_specific_context.has_lte_context()) {
    return;
  }
  if (bearer_id_by_policy_.find(PolicyID(policy_type, rule_id)) ==
      bearer_id_by_policy_.end()) {
    return;
  }
  // map change needs to be propagated to the store
  const BearerIDAndTeid bearer_id_to_delete =
      bearer_id_by_policy_[PolicyID(policy_type, rule_id)];
  bearer_id_by_policy_.erase(PolicyID(policy_type, rule_id));
  uc.is_bearer_mapping_updated = true;
  uc.bearer_id_by_policy       = bearer_id_by_policy_;

  // If it is first time filling in the DeletionReq, fill in other info
  if (!update.needs_deletion) {
    update.needs_deletion = true;
    auto& req             = update.delete_req;
    req.mutable_sid()->CopyFrom(config_.common_context.sid());
    req.set_ip_addr(config_.common_context.ue_ipv4());
    // TODO ipv6 add to the bearer request or remove ipv4
    req.set_link_bearer_id(
        config_.rat_specific_context.lte_context().bearer_id());
  }
  update.delete_req.mutable_eps_bearer_ids()->Add(
      bearer_id_to_delete.bearer_id);
}

std::vector<Teids> SessionState::get_active_teids() {
  std::vector<Teids> teids = {config_.common_context.teids()};
  for (auto bearer_pair : bearer_id_by_policy_) {
    teids.push_back(bearer_pair.second.teids);
  }
  return teids;
}

RuleSetToApply::RuleSetToApply(const magma::lte::RuleSet& rule_set) {
  for (const auto& static_rule_install : rule_set.static_rules()) {
    static_rules.insert(static_rule_install.rule_id());
  }
  for (const auto& dynamic_rule_install : rule_set.dynamic_rules()) {
    dynamic_rules[dynamic_rule_install.policy_rule().id()] =
        dynamic_rule_install.policy_rule();
  }
}

void RuleSetToApply::combine_rule_set(const RuleSetToApply& other) {
  for (const auto& static_rule : other.static_rules) {
    static_rules.insert(static_rule);
  }
  for (const auto& dynamic_pair : other.dynamic_rules) {
    dynamic_rules[dynamic_pair.first] = dynamic_pair.second;
  }
}

RuleSetBySubscriber::RuleSetBySubscriber(
    const RulesPerSubscriber& rules_per_subscriber) {
  imsi = rules_per_subscriber.imsi();
  for (const auto& rule_set : rules_per_subscriber.rule_set()) {
    if (rule_set.apply_subscriber_wide()) {
      subscriber_wide_rule_set = RuleSetToApply(rule_set);
    } else {
      subscriber_wide_rule_set        = {};
      rule_set_by_apn[rule_set.apn()] = RuleSetToApply(rule_set);
    }
  }
}

optional<RuleSetToApply> RuleSetBySubscriber::get_combined_rule_set_for_apn(
    const std::string& apn) {
  const bool apn_rule_set_exists =
      rule_set_by_apn.find(apn) != rule_set_by_apn.end();
  // Apply subscriber wide rule set if it exists. Also apply per-APN rule
  // set if it exists.
  if (apn_rule_set_exists && subscriber_wide_rule_set) {
    auto rule_set_to_apply = rule_set_by_apn[apn];
    rule_set_to_apply.combine_rule_set(*subscriber_wide_rule_set);
    return rule_set_to_apply;
  }
  if (subscriber_wide_rule_set) {
    return subscriber_wide_rule_set;
  }
  if (apn_rule_set_exists) {
    return rule_set_by_apn[apn];
  }
  return {};
}

void SessionState::update_data_metrics(
    const char* counter_name, uint64_t bytes_tx, uint64_t bytes_rx) {
  const auto sid    = get_config().common_context.sid().id();
  const auto msisdn = get_config().common_context.msisdn();
  const auto apn    = get_config().common_context.apn();
  increment_counter(
      counter_name, bytes_tx, size_t(4), LABEL_IMSI, sid.c_str(), LABEL_APN,
      apn.c_str(), LABEL_MSISDN, msisdn.c_str(), LABEL_DIRECTION, DIRECTION_UP);
  increment_counter(
      counter_name, bytes_rx, size_t(4), LABEL_IMSI, sid.c_str(), LABEL_APN,
      apn.c_str(), LABEL_MSISDN, msisdn.c_str(), LABEL_DIRECTION,
      DIRECTION_DOWN);
}

void SessionState::clear_session_metrics() {
  const auto imsi   = get_config().common_context.sid().id();
  const auto msisdn = get_config().common_context.msisdn();
  const auto apn    = get_config().common_context.apn();
  remove_counter(
      UE_USED_COUNTER_NAME, size_t(4), LABEL_IMSI, imsi.c_str(), LABEL_APN,
      apn.c_str(), LABEL_MSISDN, msisdn.c_str(), LABEL_DIRECTION, DIRECTION_UP);
  remove_counter(
      UE_USED_COUNTER_NAME, size_t(4), LABEL_IMSI, imsi.c_str(), LABEL_APN,
      apn.c_str(), LABEL_MSISDN, msisdn.c_str(), LABEL_DIRECTION,
      DIRECTION_DOWN);

  remove_counter(
      UE_DROPPED_COUNTER_NAME, size_t(4), LABEL_IMSI, imsi.c_str(), LABEL_APN,
      apn.c_str(), LABEL_MSISDN, msisdn.c_str(), LABEL_DIRECTION, DIRECTION_UP);
  remove_counter(
      UE_DROPPED_COUNTER_NAME, size_t(4), LABEL_IMSI, imsi.c_str(), LABEL_APN,
      apn.c_str(), LABEL_MSISDN, msisdn.c_str(), LABEL_DIRECTION,
      DIRECTION_DOWN);

  remove_counter(
      UE_TRAFFIC_COUNTER_NAME, size_t(3), LABEL_IMSI, imsi.c_str(),
      LABEL_SESSION_ID, session_id_.c_str(), LABEL_DIRECTION, DIRECTION_UP);
  remove_counter(
      UE_TRAFFIC_COUNTER_NAME, size_t(3), LABEL_IMSI, imsi.c_str(),
      LABEL_SESSION_ID, session_id_.c_str(), LABEL_DIRECTION, DIRECTION_DOWN);
}

CreateSessionResponse SessionState::get_create_session_response() {
  return create_session_response_;
}

void SessionState::clear_create_session_response() {
  create_session_response_ = CreateSessionResponse();
}

uint32_t SessionState::get_current_rule_version(const std::string& rule_id) {
  if (policy_version_and_stats_.find(rule_id) ==
      policy_version_and_stats_.end()) {
    MLOG(MWARNING) << "RuleID " << rule_id
                   << " doesn't have a version registered for " << session_id_
                   << ", this is unexpected";
    return 0;
  }
  return policy_version_and_stats_[rule_id].current_version;
}

void SessionState::increment_rule_stats(
    const std::string& rule_id, SessionStateUpdateCriteria& session_uc) {
  if (policy_version_and_stats_.find(rule_id) ==
      policy_version_and_stats_.end()) {
    policy_version_and_stats_[rule_id]                       = StatsPerPolicy();
    policy_version_and_stats_[rule_id].current_version       = 0;
    policy_version_and_stats_[rule_id].last_reported_version = 0;
  }
  policy_version_and_stats_[rule_id].current_version++;

  if (!session_uc.policy_version_and_stats) {
    session_uc.policy_version_and_stats = policy_version_and_stats_;
  }
}

bool operator==(const Teids& lhs, const Teids& rhs) {
  return lhs.enb_teid() == rhs.enb_teid() && lhs.agw_teid() == rhs.agw_teid();
}

}  // namespace magma
