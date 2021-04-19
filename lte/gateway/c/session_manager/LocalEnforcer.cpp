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

#include <google/protobuf/repeated_field.h>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>
#include <grpcpp/channel.h>
#include <lte/protos/session_manager.pb.h>
#include <time.h>

#include <memory>
#include <string>
#include <tuple>
#include <utility>
#include <vector>

#include "DiameterCodes.h"
#include "EnumToString.h"
#include "LocalEnforcer.h"
#include "magma_logging.h"
#include "ServiceRegistrySingleton.h"
#include "Utilities.h"

namespace magma {
uint32_t LocalEnforcer::BEARER_CREATION_DELAY_ON_SESSION_INIT = 2000;
uint32_t LocalEnforcer::REDIRECT_FLOW_PRIORITY                = 2000;
bool LocalEnforcer::SEND_ACCESS_TIMEZONE                      = false;
bool LocalEnforcer::CLEANUP_DANGLING_FLOWS                    = true;

using google::protobuf::RepeatedPtrField;

using namespace std::placeholders;

static void handle_command_level_result_code_for_monitors(
    const std::string& imsi, const std::string& session_id,
    const uint32_t result_code,
    std::unordered_set<ImsiAndSessionID>& sessions_to_terminate);

static bool is_valid_mac_address(const char* mac);

static bool parse_apn(
    const std::string& apn, std::string& mac_addr, std::string& name);

static SubscriberQuotaUpdate make_subscriber_quota_update(
    const std::string& imsi, const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state);

LocalEnforcer::LocalEnforcer(
    std::shared_ptr<SessionReporter> reporter,
    std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
    std::shared_ptr<PipelinedClient> pipelined_client,
    std::shared_ptr<EventsReporter> events_reporter,
    std::shared_ptr<SpgwServiceClient> spgw_client,
    std::shared_ptr<aaa::AAAClient> aaa_client,
    long session_force_termination_timeout_ms,
    long quota_exhaustion_termination_on_init_ms,
    magma::mconfig::SessionD mconfig)
    : reporter_(reporter),
      rule_store_(rule_store),
      pipelined_client_(pipelined_client),
      events_reporter_(events_reporter),
      spgw_client_(spgw_client),
      aaa_client_(aaa_client),
      session_store_(session_store),
      session_force_termination_timeout_ms_(
          session_force_termination_timeout_ms),
      quota_exhaustion_termination_on_init_ms_(
          quota_exhaustion_termination_on_init_ms),
      retry_timeout_(2000),
      mconfig_(mconfig),
      access_timezone_(compute_access_timezone()) {}

void LocalEnforcer::start() {
  evb_->loopForever();
}

void LocalEnforcer::attachEventBase(folly::EventBase* evb) {
  evb_ = evb;
}

void LocalEnforcer::stop() {
  evb_->terminateLoopSoon();
}

folly::EventBase& LocalEnforcer::get_event_base() {
  return *evb_;
}

void LocalEnforcer::setup(
    SessionMap& session_map, const std::uint64_t& epoch,
    std::function<void(Status status, SetupFlowsResult)> callback) {
  std::vector<SessionState::SessionInfo> session_infos;
  std::vector<SubscriberQuotaUpdate> quota_updates;
  std::vector<std::string> msisdns;
  std::vector<std::string> ue_mac_addrs;
  std::vector<std::string> apn_mac_addrs;
  std::vector<std::string> apn_names;
  std::vector<std::uint64_t> pdp_start_times;
  bool cwf = false;
  for (auto it = session_map.begin(); it != session_map.end(); it++) {
    for (const auto& session : it->second) {
      SessionState::SessionInfo session_info;
      session->get_session_info(session_info);
      session_infos.push_back(session_info);
      const auto& config = session->get_config();
      msisdns.push_back(config.common_context.msisdn());

      std::string apn_mac_addr;
      std::string apn_name;
      auto apn = config.common_context.apn();
      if (!parse_apn(apn, apn_mac_addr, apn_name)) {
        MLOG(MWARNING) << "Failed mac/name parsing for apn " << apn;
        apn_mac_addr = "";
        apn_name     = apn;
      }
      apn_mac_addrs.push_back(apn_mac_addr);
      apn_names.push_back(apn_name);
      pdp_start_times.push_back(session->get_pdp_start_time());

      if (session->is_radius_cwf_session()) {
        cwf                      = true;
        const auto& wlan_context = config.rat_specific_context.wlan_context();
        const auto& ue_mac_addr  = wlan_context.mac_addr();
        ue_mac_addrs.push_back(ue_mac_addr);
        SubscriberQuotaUpdate update = make_subscriber_quota_update(
            session_info.imsi, ue_mac_addr,
            session->get_subscriber_quota_state());
        quota_updates.push_back(update);
      }
    }
  }
  // TODO this assumption of CWF only deployments will not be relevant for long
  if (cwf) {
    pipelined_client_->setup_cwf(
        session_infos, quota_updates, ue_mac_addrs, msisdns, apn_mac_addrs,
        apn_names, pdp_start_times, epoch, callback);
  } else {
    pipelined_client_->setup_lte(session_infos, epoch, callback);
  }
}

void LocalEnforcer::sync_sessions_on_restart(std::time_t current_time) {
  auto session_map    = session_store_.read_all_sessions();
  auto session_update = SessionStore::get_default_session_update(session_map);
  // Update the sessions so that their rules match the current timestamp
  for (auto& it : session_map) {
    const auto& imsi = it.first;
    for (auto& session : it.second) {
      const std::string session_id = session->get_session_id();
      auto& uc                     = session_update[imsi][session_id];
      // Reschedule Revalidation Timer if it was pending before
      auto triggers   = session->get_event_triggers();
      auto trigger_it = triggers.find(REVALIDATION_TIMEOUT);
      if (trigger_it != triggers.end() &&
          triggers[REVALIDATION_TIMEOUT] == PENDING) {
        // the bool value indicates whether the trigger has been triggered
        const auto revalidation_time = session->get_revalidation_time();
        schedule_revalidation(*session, revalidation_time, uc);
      }

      session->sync_rules_to_time(current_time, uc);

      for (std::string rule_id : session->get_static_rules()) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        if (lifetime.deactivation_time > current_time) {
          auto rule_install =
              session->get_static_rule_install(rule_id, lifetime);
          schedule_static_rule_deactivation(
              imsi, session_id, rule_id, lifetime.deactivation_time);
        }
      }
      // Schedule rule activations / deactivations
      for (std::string rule_id : session->get_scheduled_static_rules()) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        schedule_static_rule_activation(
            imsi, session_id, rule_id, lifetime.activation_time);
        if (lifetime.deactivation_time > current_time) {
          schedule_static_rule_deactivation(
              imsi, session_id, rule_id, lifetime.deactivation_time);
        }
      }

      std::vector<std::string> rule_ids;
      session->get_dynamic_rules().get_rule_ids(rule_ids);
      for (std::string rule_id : rule_ids) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        if (lifetime.deactivation_time > current_time) {
          auto rule_install =
              session->get_dynamic_rule_install(rule_id, lifetime);
          schedule_dynamic_rule_deactivation(
              imsi, session_id, rule_id, lifetime.deactivation_time);
        }
      }
      rule_ids.clear();
      session->get_scheduled_dynamic_rules().get_rule_ids(rule_ids);
      for (auto rule_id : rule_ids) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        auto rule_install =
            session->get_dynamic_rule_install(rule_id, lifetime);
        schedule_dynamic_rule_activation(
            imsi, session_id, rule_id, lifetime.deactivation_time);
        if (lifetime.deactivation_time > current_time) {
          schedule_dynamic_rule_deactivation(
              imsi, session_id, rule_id, lifetime.deactivation_time);
        }
      }
      // Reschedule termination if subscriber has no quota
      if (terminate_on_wallet_exhaust()) {
        handle_session_activate_subscriber_quota_state(*session);
      }
    }
  }
  bool success = session_store_.update_sessions(session_update);
  if (success) {
    MLOG(MDEBUG) << "Successfully synced sessions after restart";
  } else {
    MLOG(MERROR) << "Failed to sync sessions after restart";
  }
}

void LocalEnforcer::aggregate_records(
    SessionMap& session_map, const RuleRecordTable& records,
    SessionUpdate& session_update) {
  // Insert the IMSI+SessionID for sessions we received a rule record into a set
  // for easy access
  std::unordered_set<ImsiAndSessionID> sessions_with_reporting_flows;
  // In some failure cases, PipelineD may still hold onto flows for sessions
  // that do not exist in SessionD. In this case, send DeactivateFlowsRequest
  RuleRecordSet dead_sessions_to_cleanup;

  for (const RuleRecord& record : records.records()) {
    const std::string& imsi = record.sid();
    const uint32_t teid     = record.teid();
    SessionSearchCriteria criteria(imsi, IMSI_AND_TEID, teid);
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "Could not find session for " << imsi << " and " << teid
                   << " during record aggregation";
      dead_sessions_to_cleanup.insert(record);
      continue;
    }
    auto& session          = **session_it;
    const auto& session_id = session->get_session_id();
    if (record.bytes_tx() > 0 || record.bytes_rx() > 0) {
      sessions_with_reporting_flows.insert(ImsiAndSessionID(imsi, session_id));
      MLOG(MINFO) << session_id << " used " << record.bytes_tx()
                  << " tx bytes and " << record.bytes_rx()
                  << " rx bytes for rule " << record.rule_id();
    }
    SessionStateUpdateCriteria& uc = session_update[imsi][session_id];

    session->add_rule_usage(
        record.rule_id(), record.bytes_tx(), record.bytes_rx(),
        record.dropped_tx(), record.dropped_rx(), uc);
  }
  complete_termination_for_released_sessions(
      session_map, sessions_with_reporting_flows, session_update);
  cleanup_dead_sessions(dead_sessions_to_cleanup);
}

void LocalEnforcer::cleanup_dead_sessions(
    const RuleRecordSet dead_sessions_to_cleanup) {
  if (dead_sessions_to_cleanup.size() == 0) {
    return;
  }
  if (!CLEANUP_DANGLING_FLOWS) {
    MLOG(MWARNING) << "Not cleaning up " << dead_sessions_to_cleanup.size()
                   << " dangling sessions in PipelineD due to "
                      "'cleanup_all_dangling_flows: false' in sessiond.yml";
    return;
  }

  MLOG(MINFO) << "Deactivating " << dead_sessions_to_cleanup.size()
              << " dangling sessions in PipelineD";
  for (const RuleRecord& record : dead_sessions_to_cleanup) {
    Teids teids;
    teids.set_agw_teid(record.teid());
    const std::vector<Teids> teids_vec{teids};
    pipelined_client_->deactivate_flows_for_rules_for_termination(
        record.sid(), record.ue_ipv4(), record.ue_ipv6(), teids_vec,
        RequestOriginType::WILDCARD);
  }
}

void LocalEnforcer::complete_termination_for_released_sessions(
    SessionMap& session_map,
    std::unordered_set<ImsiAndSessionID> sessions_with_reporting_flows,
    SessionUpdate& session_update) {
  // Iterate through sessions and notify that report has finished. Terminate any
  // sessions that can be terminated.
  std::vector<ImsiAndSessionID> sessions_to_terminate;
  for (const auto& session_pair : session_map) {
    const std::string imsi = session_pair.first;
    for (const auto& session : session_pair.second) {
      const std::string session_id = session->get_session_id();
      // If we did not receive a rule record for the session, this means
      // PipelineD has reported all usage for the session
      auto imsi_and_session_id = ImsiAndSessionID(imsi, session_id);
      if (session->get_state() == SESSION_RELEASED &&
          sessions_with_reporting_flows.find(imsi_and_session_id) ==
              sessions_with_reporting_flows.end()) {
        sessions_to_terminate.push_back(imsi_and_session_id);
      }
    }
  }
  // Do the actual termination in a separate loop since this can modify the
  // session map structure
  for (const auto& pair : sessions_to_terminate) {
    const auto &imsi = pair.first, &session_id = pair.second;
    complete_termination(session_map, imsi, session_id, session_update);
  }
}

void LocalEnforcer::execute_actions(
    SessionMap& session_map,
    const std::vector<std::unique_ptr<ServiceAction>>& actions,
    SessionUpdate& session_update) {
  for (const auto& action_p : actions) {
    auto imsi       = action_p->get_imsi();
    auto session_id = action_p->get_session_id();
    switch (action_p->get_type()) {
      case ACTIVATE_SERVICE:
        handle_activate_service_action(action_p);
        break;
      case REDIRECT:
      case RESTRICT_ACCESS: {
        install_final_unit_action_flows(action_p);
        break;
      }
      case TERMINATE_SERVICE: {
        auto found = find_and_terminate_session(
            session_map, imsi, session_id, session_update);
        if (!found) {
          MLOG(MERROR) << "Cannot act on TERMINATE action since session "
                       << session_id << " does not exist";
        }
        break;
      }
      case CONTINUE_SERVICE:
        break;
    }
  }
}

// TODO look into whether we need to re-install all Gx rules on activation
void LocalEnforcer::handle_activate_service_action(
    const std::unique_ptr<ServiceAction>& action_p) {
  const std::string imsi      = action_p->get_imsi(),
                    msisdn    = action_p->get_msisdn(),
                    ip_addr   = action_p->get_ip_addr(),
                    ipv6_addr = action_p->get_ipv6_addr();
  const Teids teids           = action_p->get_teids();
  const auto ambr             = action_p->get_ambr();
  RulesToProcess to_process   = action_p->get_gx_rules_to_install();
  pipelined_client_->activate_flows_for_rules(
      imsi, ip_addr, ipv6_addr, teids, msisdn, ambr, to_process,
      std::bind(
          &LocalEnforcer::handle_activate_ue_flows_callback, this, imsi,
          ip_addr, ipv6_addr, teids, _1, _2));
}

// Terminates sessions that correspond to the given IMSI and session.
void LocalEnforcer::start_session_termination(
    const std::unique_ptr<SessionState>& session, bool notify_access,
    SessionStateUpdateCriteria& session_uc) {
  const std::string& imsi       = session->get_imsi();
  const std::string& session_id = session->get_session_id();
  if (session->is_terminating()) {
    // If the session is terminating already, do nothing.
    MLOG(MINFO) << "Session " << session_id << " is already terminating, "
                << "ignoring termination request";
    return;
  }
  MLOG(MINFO) << "Initiating session termination for " << session_id;
  session->set_pdp_end_time(magma::get_time_in_sec_since_epoch(), session_uc);

  remove_all_rules_for_termination(session, session_uc);
  session->set_fsm_state(SESSION_RELEASED, session_uc);
  const auto& config         = session->get_config();
  const auto& common_context = config.common_context;
  if (notify_access) {
    notify_termination_to_access_service(session_id, config);
  }
  if (common_context.rat_type() == TGPP_WLAN) {
    MLOG(MDEBUG) << "Deleting UE MAC flow for subscriber " << imsi;
    pipelined_client_->delete_ue_mac_flow(
        common_context.sid(),
        config.rat_specific_context.wlan_context().mac_addr());
  }
  if (terminate_on_wallet_exhaust()) {
    handle_subscriber_quota_state_change(
        *session, SubscriberQuotaUpdate_Type_TERMINATE, &session_uc);
  }
  // The termination should be completed when aggregated usage record no
  // longer
  // includes the imsi. If this has not occurred after the timeout, force
  // terminate the session.
  MLOG(MDEBUG) << "Scheduling a force termination timeout for " << session_id
               << " in " << session_force_termination_timeout_ms_ << "ms";
  evb_->runAfterDelay(
      [this, imsi, session_id] {
        handle_force_termination_timeout(imsi, session_id);
      },
      session_force_termination_timeout_ms_);
}

void LocalEnforcer::handle_force_termination_timeout(
    const std::string& imsi, const std::string& session_id) {
  auto session_map    = session_store_.read_sessions_for_deletion({imsi});
  auto session_update = SessionStore::get_default_session_update(session_map);
  bool needs_termination =
      session_update[imsi].find(session_id) != session_update[imsi].end();
  MLOG(MDEBUG) << "Forced termination timeout! Checking if termination has to "
               << "be forced for " << session_id << "... => "
               << (needs_termination ? "YES" : "NO");
  // If the session doesn't exist in the session_update, then the session was
  // already terminated + removed
  if (needs_termination) {
    complete_termination(session_map, imsi, session_id, session_update);
    bool end_success = session_store_.update_sessions(session_update);
    if (end_success) {
      MLOG(MDEBUG) << "Updated session termination of " << session_id
                   << " in to SessionStore";
    } else {
      MLOG(MDEBUG) << "Failed to update session termination of " << session_id
                   << " in to SessionStore";
    }
  }
}

void LocalEnforcer::remove_all_rules_for_termination(
    const std::unique_ptr<SessionState>& session,
    SessionStateUpdateCriteria& uc) {
  const std::string ip_addr = session->get_config().common_context.ue_ipv4();
  const auto ipv6_addr      = session->get_config().common_context.ue_ipv6();
  const std::vector<Teids> teids = session->get_active_teids();
  pipelined_client_->deactivate_flows_for_rules_for_termination(
      session->get_imsi(), ip_addr, ipv6_addr, teids,
      RequestOriginType::WILDCARD);
  session->remove_all_rules_for_termination(uc);
}

void LocalEnforcer::notify_termination_to_access_service(
    const std::string& session_id, const SessionConfig& config) {
  const std::string& imsi = config.get_imsi();
  auto common_context     = config.common_context;
  switch (common_context.rat_type()) {
    case TGPP_WLAN: {
      // tell AAA service to terminate radius session if necessary
      const auto& radius_session_id =
          config.rat_specific_context.wlan_context().radius_session_id();
      MLOG(MDEBUG) << "Asking AAA service to terminate session with "
                   << "Radius ID: " << radius_session_id << ", IMSI: " << imsi;
      aaa_client_->terminate_session(radius_session_id, imsi);
      break;
    }
    case TGPP_LTE: {
      // Deleting the PDN session by triggering network issued default bearer
      // deactivation
      const auto& lte_context = config.rat_specific_context.lte_context();
      spgw_client_->delete_default_bearer(
          imsi, common_context.ue_ipv4(), lte_context.bearer_id());
      break;
    }
    default:
      // Should not get here
      MLOG(MWARNING) << session_id << " has an invalid RAT Type "
                     << config.common_context.rat_type();
      return;
  }
}

void LocalEnforcer::handle_subscriber_quota_state_change(
    SessionState& session, SubscriberQuotaUpdate_Type new_state,
    SessionStateUpdateCriteria* session_uc) {
  auto config                   = session.get_config();
  const std::string& imsi       = session.get_imsi();
  const std::string& session_id = session.get_session_id();
  MLOG(MINFO) << session_id << " now has subscriber wallet status: "
              << wallet_state_to_str(new_state);
  session.set_subscriber_quota_state(new_state, session_uc);
  std::string ue_mac_addr = "";
  auto rat_specific       = config.rat_specific_context;
  if (rat_specific.has_wlan_context()) {
    ue_mac_addr = rat_specific.wlan_context().mac_addr();
  }
  report_subscriber_state_to_pipelined(imsi, ue_mac_addr, new_state);
}

void LocalEnforcer::install_final_unit_action_flows(
    const std::unique_ptr<ServiceAction>& action_p) {
  const std::string imsi       = action_p->get_imsi(),
                    msisdn     = action_p->get_msisdn(),
                    ip_addr    = action_p->get_ip_addr(),
                    ipv6_addr  = action_p->get_ipv6_addr(),
                    session_id = action_p->get_session_id();
  const Teids teids            = action_p->get_teids();
  const auto fua_type          = action_p->get_type();

  MLOG(MINFO) << "Installing final unit action "
              << service_action_type_to_str(fua_type) << " flows for "
              << session_id;
  RulesToProcess to_process = action_p->get_gy_rules_to_install();
  pipelined_client_->add_gy_final_action_flow(
      imsi, ip_addr, ipv6_addr, teids, msisdn, to_process);
}

UpdateSessionRequest LocalEnforcer::collect_updates(
    SessionMap& session_map,
    std::vector<std::unique_ptr<ServiceAction>>& actions,
    SessionUpdate& session_update) const {
  UpdateSessionRequest request;
  for (const auto& session_pair : session_map) {
    for (const auto& session : session_pair.second) {
      std::string imsi = session_pair.first;
      std::string sid  = session->get_session_id();
      auto& uc         = session_update[imsi][sid];
      session->get_updates(request, &actions, uc);
    }
  }
  return request;
}

void LocalEnforcer::handle_update_failure(
    SessionMap& session_map, const UpdateRequestsBySession& failed_request,
    SessionUpdate& updates) {
  for (const auto& request_by_id : failed_request.requests_by_id) {
    const std::string& imsi       = request_by_id.first.first;
    const std::string& session_id = request_by_id.first.second;
    const auto& request           = request_by_id.second;
    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "Could not reset failed update for " << session_id
                   << " because it couldn't be found";
      continue;
    }
    (**session_it)->handle_update_failure(request, updates[imsi][session_id]);
  }
}

/*
 * If a rule needs to be tracked by the OCS, then it needs credit in order to
 * be activated. If it does not receive credit, it should not be installed.
 * If a rule has a monitoring key, it is not required that a usage monitor is
 * installed with quota
 */
static bool should_activate(
    const PolicyRule& rule,
    const std::unordered_set<uint32_t>& successful_credits) {
  if (rule.tracking_type() == PolicyRule::ONLY_OCS ||
      rule.tracking_type() == PolicyRule::OCS_AND_PCRF) {
    const bool exists = successful_credits.count(rule.rating_group()) > 0;
    if (!exists) {
      MLOG(MERROR) << "Not activating Gy tracked " << rule.id()
                   << " because credit w/ rating group " << rule.rating_group()
                   << " does not exist";
      return false;
    }
  }
  switch (rule.tracking_type()) {
    case PolicyRule::ONLY_PCRF:
      MLOG(MINFO) << "Activating Gx tracked rule " << rule.id()
                  << " with monitoring key " << rule.monitoring_key();
      break;
    case PolicyRule::ONLY_OCS:
      MLOG(MINFO) << "Activating Gy tracked rule " << rule.id()
                  << " with rating group " << rule.rating_group();
      break;
    case PolicyRule::OCS_AND_PCRF:
      MLOG(MINFO) << "Activating Gx+Gy tracked rule " << rule.id()
                  << " with monitoring key " << rule.monitoring_key()
                  << " with rating group " << rule.rating_group();
      break;
    case PolicyRule::NO_TRACKING:
      MLOG(MINFO) << "Activating untracked rule " << rule.id();
      break;
    default:
      MLOG(MINFO) << "Invalid rule tracking type " << rule.id();
      return false;
  }
  return true;
}

void LocalEnforcer::schedule_static_rule_activation(
    const std::string& imsi, const std::string& session_id,
    const std::string& rule_id, const std::time_t activation_time) {
  auto delta = magma::time_difference_from_now(activation_time);
  MLOG(MDEBUG) << "Scheduling " << session_id << " static rule " << rule_id
               << " activation in " << (delta.count() / 1000) << " secs";

  evb_->runAfterDelay(
      [=] {
        auto session_map = session_store_.read_sessions(SessionRead{imsi});
        auto session_update =
            session_store_.get_default_session_update(session_map);
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MWARNING) << "Could not find session  " << session_id
                         << "during installation of static rule " << rule_id;
          return;
        }
        PolicyRule rule;
        if (!rule_store_->get_rule(rule_id, &rule)) {
          MLOG(MWARNING) << "Could not find static rules definition for "
                         << rule_id;
          return;
        }

        auto& session = **session_it;
        auto& uc      = session_update[imsi][session_id];

        std::time_t current_time = time(nullptr);
        // don't install the rule if the current time is out of lifetime
        if (!session->should_rule_be_active(rule_id, current_time)) {
          session->deactivate_scheduled_static_rule(rule_id);
          session_store_.update_sessions(session_update);
          return;
        }

        auto config               = session->get_config();
        const std::string ip_addr = config.common_context.ue_ipv4();
        const auto ipv6_addr      = config.common_context.ue_ipv6();
        const Teids teids         = config.common_context.teids();
        const auto ambr           = config.get_apn_ambr();
        const std::string msisdn  = config.common_context.msisdn();

        RulesToProcess to_process;
        to_process.push_back(session->activate_static_rule(
            rule_id, session->get_rule_lifetime(rule_id), uc));

        pipelined_client_->activate_flows_for_rules(
            imsi, ip_addr, ipv6_addr, teids, msisdn, ambr, to_process,
            std::bind(
                &LocalEnforcer::handle_activate_ue_flows_callback, this, imsi,
                ip_addr, ipv6_addr, teids, _1, _2));

        session_store_.update_sessions(session_update);
      },
      delta.count());
}

void LocalEnforcer::schedule_dynamic_rule_activation(
    const std::string& imsi, const std::string& session_id,
    const std::string& rule_id, const std::time_t activation_time) {
  auto delta = magma::time_difference_from_now(activation_time);
  MLOG(MDEBUG) << "Scheduling " << session_id << " dynamic rule " << rule_id
               << " activation in " << (delta.count() / 1000) << " secs";
  evb_->runAfterDelay(
      [=] {
        auto session_map = session_store_.read_sessions(SessionRead{imsi});
        auto session_update =
            session_store_.get_default_session_update(session_map);
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MWARNING) << "Could not find session  " << session_id
                         << "during installation of dynamic rule " << rule_id;
          return;
        }
        auto& session    = **session_it;
        auto& session_uc = session_update[imsi][session_id];
        if (!session->is_dynamic_rule_scheduled(rule_id)) {
          return;
        }
        // don't install the rule if the current time is out of lifetime
        std::time_t current_time = time(nullptr);
        if (!session->should_rule_be_active(rule_id, current_time)) {
          session->remove_scheduled_dynamic_rule(rule_id, nullptr, session_uc);
          session_store_.update_sessions(session_update);
          return;
        }

        auto config          = session->get_config();
        const auto ip_addr   = config.common_context.ue_ipv4();
        const auto ipv6_addr = config.common_context.ue_ipv6();
        const Teids teids    = config.common_context.teids();
        const auto ambr      = config.get_apn_ambr();
        const auto msisdn    = config.common_context.msisdn();

        PolicyRule rule;
        if (!session->remove_scheduled_dynamic_rule(
                rule_id, &rule, session_uc)) {
          MLOG(MWARNING) << "Dynamic rule " << rule_id << " doesn't exist for "
                         << session_id;
          return;
        }
        RulesToProcess to_process;
        to_process.push_back(session->insert_dynamic_rule(
            rule, session->get_rule_lifetime(rule_id), session_uc));

        pipelined_client_->activate_flows_for_rules(
            imsi, ip_addr, ipv6_addr, teids, msisdn, ambr, to_process,
            std::bind(
                &LocalEnforcer::handle_activate_ue_flows_callback, this, imsi,
                ip_addr, ipv6_addr, teids, _1, _2));

        session_store_.update_sessions(session_update);
      },
      delta.count());
}

void LocalEnforcer::schedule_static_rule_deactivation(
    const std::string& imsi, const std::string& session_id,
    const std::string& rule_id, const std::time_t deactivation_time) {
  auto delta = magma::time_difference_from_now(deactivation_time);
  MLOG(MDEBUG) << "Scheduling " << session_id << " static rule " << rule_id
               << " deactivation in " << (delta.count() / 1000) << " secs";
  evb_->runAfterDelay(
      [=] {
        auto session_map = session_store_.read_sessions(SessionRead{imsi});
        auto session_update =
            session_store_.get_default_session_update(session_map);
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MERROR) << "Could not find session  " << session_id
                       << "during removal of static rule " << rule_id;
          return;
        }
        auto& session = **session_it;
        if (session->should_rule_be_active(rule_id, time(nullptr))) {
          return;
        }
        PolicyRule rule;
        if (!rule_store_->get_rule(rule_id, &rule)) {
          MLOG(MERROR) << "static rule " << rule_id
                       << " is not found, skipping deactivation...";
          return;
        }
        auto ip_addr      = session->get_config().common_context.ue_ipv4();
        auto ipv6_addr    = session->get_config().common_context.ue_ipv6();
        const Teids teids = session->get_config().common_context.teids();

        auto& session_uc = session_update[imsi][session_id];
        optional<RuleToProcess> op_rule_info =
            session->deactivate_static_rule(rule_id, session_uc);
        if (!op_rule_info) {
          MLOG(MERROR) << "Could not find rule " << rule_id << " for "
                       << session_id << " during static rule removal";
          return;
        }

        RulesToProcess to_process;
        to_process.push_back(*op_rule_info);
        pipelined_client_->deactivate_flows_for_rules(
            imsi, ip_addr, ipv6_addr, teids, to_process, RequestOriginType::GX);

        session_store_.update_sessions(session_update);
      },
      delta.count());
}

void LocalEnforcer::schedule_dynamic_rule_deactivation(
    const std::string& imsi, const std::string& session_id,
    const std::string& rule_id, const std::time_t deactivation_time) {
  auto delta = magma::time_difference_from_now(deactivation_time);
  MLOG(MDEBUG) << "Scheduling " << session_id << " dynamic rule " << rule_id
               << " deactivation in " << (delta.count() / 1000) << " secs";
  evb_->runAfterDelay(
      [=] {
        auto session_map = session_store_.read_sessions(SessionRead{imsi});
        auto session_update =
            session_store_.get_default_session_update(session_map);
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MWARNING) << "Could not find session " << session_id
                         << "during removal of dynamic rule " << rule_id;
          return;
        }
        auto& session = **session_it;
        if (session->should_rule_be_active(rule_id, time(nullptr))) {
          return;
        }
        auto ip_addr      = session->get_config().common_context.ue_ipv4();
        auto ipv6_addr    = session->get_config().common_context.ue_ipv6();
        const Teids teids = session->get_config().common_context.teids();

        PolicyRule policy;
        auto& uc = session_update[imsi][session_id];
        optional<RuleToProcess> remove_info =
            session->remove_dynamic_rule(policy.id(), &policy, uc);
        if (remove_info) {
          RulesToProcess to_process;
          to_process.push_back(*remove_info);
          pipelined_client_->deactivate_flows_for_rules(
              imsi, ip_addr, ipv6_addr, teids, to_process,
              RequestOriginType::GX);
        }
        session_store_.update_sessions(session_update);
      },
      delta.count());
}

void LocalEnforcer::filter_rule_installs(
    std::vector<StaticRuleInstall>& static_installs,
    std::vector<DynamicRuleInstall>& dynamic_installs,
    const std::unordered_set<uint32_t>& successful_credits) {
  // Filter out static rules that we will not install nor schedule
  auto end_of_valid_st_rules = std::remove_if(
      static_installs.begin(), static_installs.end(),
      [&](StaticRuleInstall& rule_install) {
        auto& id = rule_install.rule_id();
        PolicyRule rule;
        if (!rule_store_->get_rule(id, &rule)) {
          LOG(ERROR) << "Not activating rule " << id
                     << " because it could not be found";
          return true;
        }
        return !should_activate(rule, successful_credits);
      });
  static_installs.erase(end_of_valid_st_rules, static_installs.end());

  // Filter out dynamic rules that we will not install nor schedule
  auto end_of_valid_dy_rules = std::remove_if(
      dynamic_installs.begin(), dynamic_installs.end(),
      [&](DynamicRuleInstall& rule_install) {
        return !should_activate(rule_install.policy_rule(), successful_credits);
      });
  dynamic_installs.erase(end_of_valid_dy_rules, dynamic_installs.end());
}

void LocalEnforcer::handle_session_activate_rule_updates(
    SessionState& session, const CreateSessionResponse& response,
    std::unordered_set<uint32_t>& charging_credits_received) {
  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;

  std::vector<StaticRuleInstall> static_rule_installs =
      to_vec(response.static_rules());
  std::vector<DynamicRuleInstall> dynamic_rule_installs =
      to_vec(response.dynamic_rules());
  filter_rule_installs(
      static_rule_installs, dynamic_rule_installs, charging_credits_received);

  SessionStateUpdateCriteria uc;  // TODO remove unused UC
  session.process_rules_to_install(
      static_rule_installs, dynamic_rule_installs, &pending_activation,
      &pending_deactivation, &pending_bearer_setup, &pending_scheduling, &uc);

  const std::string& imsi       = session.get_imsi();
  const std::string& session_id = session.get_session_id();
  handle_rule_scheduling(imsi, session_id, pending_scheduling);

  // activate_flows_for_rules() should be called even if there is no rule
  // to activate, because pipelined activates a "drop all packet" rule
  // when no rule is provided as the parameter.
  const SessionConfig& config = session.get_config();
  propagate_rule_updates_to_pipelined(
      config, pending_activation, pending_deactivation, true);

  if (config.common_context.rat_type() == TGPP_LTE) {
    BearerUpdate bearer_updates = session.get_dedicated_bearer_updates(
        pending_bearer_setup, pending_deactivation, &uc);
    if (bearer_updates.needs_creation) {
      // If a bearer creation is needed, we need to delay this by a few seconds
      // so that the attach fully completes before.
      schedule_session_init_dedicated_bearer_creations(
          imsi, session_id, bearer_updates);
    }
  }
}

void LocalEnforcer::schedule_session_init_dedicated_bearer_creations(
    const std::string& imsi, const std::string& session_id,
    BearerUpdate& bearer_updates) {
  MLOG(MINFO)
      << "Scheduling a dedicated bearer creation request for newly created "
      << session_id << " in "
      << LocalEnforcer::BEARER_CREATION_DELAY_ON_SESSION_INIT << " ms";
  evb_->runAfterDelay(
      [this, imsi, session_id, bearer_updates]() mutable {
        auto session_map = session_store_.read_sessions({imsi});
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MWARNING) << "Ignoring dedicated bearer creations from session "
                            "creation for "
                         << session_id << " since it no longer exists";
          return;
        }
        auto& session = **session_it;

        // Skip bearer update if it is no longer ACTIVE
        if (session->get_state() != SESSION_ACTIVE) {
          MLOG(MWARNING) << "Ignoring dedicated bearer create request from"
                         << " session creation for " << session_id
                         << " since the session is no longer active";
          return;
        }
        // Check that the policies are still installed and needs a
        // bearer
        auto rules = bearer_updates.create_req.mutable_policy_rules();
        auto it    = rules->begin();
        while (it != rules->end()) {
          auto policy_type = session->get_policy_type(it->id());
          if (!policy_type) {
            MLOG(MWARNING) << "Ignoring dedicated bearer create request from"
                           << " session creation for " << session_id
                           << " policy ID: " << it->id()
                           << " since the policy is no longer active in the "
                              "session";
            it = rules->erase(it);
            continue;
          }
          if (!session->policy_needs_bearer_creation(*policy_type, it->id())) {
            MLOG(MWARNING) << "Ignoring dedicated bearer create request from "
                           << "session creation for " << session_id
                           << " and policy ID: " << it->id()
                           << " since the policy no longer needs a bearer";
            it = rules->erase(it);
            continue;
          }
          ++it;
        }
        if (rules->size() > 0) {
          propagate_bearer_updates_to_mme(bearer_updates);
        }
        return;
        // No need to update session store for bearer creations.
        // SessionStore will be updated once the bearer binding
        // completes/fails.
      },
      LocalEnforcer::BEARER_CREATION_DELAY_ON_SESSION_INIT);
}

void LocalEnforcer::init_session(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, const SessionConfig& cfg,
    const CreateSessionResponse& response) {
  const auto time_since_epoch = magma::get_time_in_sec_since_epoch();
  auto session                = std::make_unique<SessionState>(
      imsi, session_id, cfg, *rule_store_, response.tgpp_ctx(),
      time_since_epoch, response);
  session_map[imsi].push_back(std::move(session));
}

bool LocalEnforcer::update_tunnel_ids(
    SessionMap& session_map, const UpdateTunnelIdsRequest& request) {
  const auto imsi             = request.sid().id();
  const uint32_t bearer       = request.bearer_id();
  const auto time_since_epoch = magma::get_time_in_sec_since_epoch();

  SessionSearchCriteria criteria(imsi, IMSI_AND_BEARER, bearer);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "Could not find session for " << imsi << " and bearer "
                 << bearer << " during update tunnelIds ";
    return false;
  }
  auto& session = **session_it;

  session->set_teids(request.enb_teid(), request.agw_teid());
  const CreateSessionResponse& csr = session->get_create_session_response();

  std::unordered_set<uint32_t> charging_credits_received;
  for (const auto& credit : csr.credits()) {
    // TODO this uc is not doing anything here, modify interface
    auto uc = get_default_update_criteria();
    if (session->receive_charging_credit(credit, uc)) {
      charging_credits_received.insert(credit.charging_key());
    }
  }
  // We don't have to check 'success' field for monitors because command level
  // errors are handled in session proxy for the init exchange
  for (const auto& monitor : csr.usage_monitors()) {
    // TODO this uc is not doing anything here, modify interface
    auto uc = get_default_update_criteria();
    session->receive_monitor(monitor, uc);
  }

  handle_session_activate_rule_updates(
      *session, csr, charging_credits_received);

  update_ipfix_flow(imsi, session->get_config(), time_since_epoch);
  if (terminate_on_wallet_exhaust()) {
    handle_session_activate_subscriber_quota_state(*session);
  }

  if (revalidation_required(csr.event_triggers())) {
    // TODO This might not work since the session is not initialized properly
    // at this point
    auto _ = get_default_update_criteria();
    schedule_revalidation(*session, csr.revalidation_time(), _);
  }

  // handle transient errors during first init
  auto session_id = session->get_session_id();
  for (const auto& credit : csr.credits()) {
    if (!credit.success() &&
        DiameterCodeHandler::is_transient_failure(credit.result_code())) {
      auto _    = get_default_update_criteria();
      auto ckey = credit.charging_key();

      // set state here because during init update criteria is not used.
      session->set_suspend_credit(ckey, true, _);
      session->suspend_service_if_needed_for_credit(ckey, _);

      // schedule the removal of rules to avoid problems with install-unistall
      // order
      evb_->runAfterDelay(
          [this, imsi, session_id, credit] {
            auto session_map = session_store_.read_sessions({imsi});
            SessionSearchCriteria criteria(
                imsi, IMSI_AND_SESSION_ID, session_id);
            auto session_it =
                session_store_.find_session(session_map, criteria);
            if (!session_it) {
              MLOG(MWARNING) << "Not suspending credit for " << session_id
                             << " since it no longer exists";
              return;
            }
            auto& session = **session_it;
            auto _        = get_default_update_criteria();
            remove_rules_for_suspended_credit(
                session, credit.charging_key(), _);
          },
          500);
    }
  }

  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    // First time a session is created for IMSI
    MLOG(MDEBUG) << "First session for " << imsi << " with SessionID "
                 << session_id;
    session_map[imsi] = SessionVector();
  }

  // once the session is activated, just get rid of the create session response
  session->clear_create_session_response();

  events_reporter_->session_created(
      imsi, session_id, session->get_config(), session);
  return true;
}

bool LocalEnforcer::terminate_on_wallet_exhaust() {
  return mconfig_.has_wallet_exhaust_detection() &&
         mconfig_.wallet_exhaust_detection().terminate_on_exhaust();
}

bool LocalEnforcer::is_wallet_exhausted(SessionState& session) {
  switch (mconfig_.wallet_exhaust_detection().method()) {
    case magma::mconfig::WalletExhaustDetection_Method_GxTrackedRules:
      return !session.active_monitored_rules_exist();
    default:
      MLOG(MWARNING) << "This method is not yet supported...";
      return false;
  }
}

void LocalEnforcer::handle_session_activate_subscriber_quota_state(
    SessionState& session) {
  if (is_wallet_exhausted(session)) {
    handle_subscriber_quota_state_change(
        session, SubscriberQuotaUpdate_Type_NO_QUOTA, nullptr);
    // Schedule a session termination for a configured number of seconds after
    // session create
    const std::string& imsi       = session.get_imsi();
    const std::string& session_id = session.get_session_id();
    MLOG(MINFO) << "Scheduling session for session " << session_id
                << " to be terminated in "
                << quota_exhaustion_termination_on_init_ms_ << " ms";
    auto sessions_to_terminate = std::unordered_set<ImsiAndSessionID>{
        ImsiAndSessionID(imsi, session_id)};
    schedule_termination(sessions_to_terminate);
    return;
  }

  // Valid Quota
  handle_subscriber_quota_state_change(
      session, SubscriberQuotaUpdate_Type_VALID_QUOTA, nullptr);
  return;
}

void LocalEnforcer::schedule_termination(
    std::unordered_set<ImsiAndSessionID>& sessions) {
  evb_->runAfterDelay(
      [this, sessions] {
        SessionRead req;
        for (auto& imsi_and_session_id : sessions) {
          req.insert(imsi_and_session_id.first);
        }
        auto session_map = session_store_.read_sessions_for_deletion(req);

        SessionUpdate session_update =
            SessionStore::get_default_session_update(session_map);

        terminate_multiple_sessions(session_map, sessions, session_update);
        bool end_success = session_store_.update_sessions(session_update);
        if (end_success) {
          MLOG(MDEBUG) << "Succeeded in updating session store with "
                       << "termination initialization";
        } else {
          MLOG(MDEBUG) << "Failed in updating session store with "
                       << "termination initialization";
        }
      },
      quota_exhaustion_termination_on_init_ms_);
}

void LocalEnforcer::report_subscriber_state_to_pipelined(
    const std::string& imsi, const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state) {
  auto update = make_subscriber_quota_update(imsi, ue_mac_addr, state);
  pipelined_client_->update_subscriber_quota_state(
      std::vector<SubscriberQuotaUpdate>{update});
}

void LocalEnforcer::complete_termination(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, SessionUpdate& session_updates) {
  // If the session cannot be found in session_map, or a new session has
  // already begun, do nothing.
  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    // Session is already deleted, or new session already began, ignore.
    MLOG(MDEBUG) << "Could not find session for IMSI " << imsi
                 << " and session ID " << session_id
                 << ". Skipping termination.";
  }
  auto& session    = **session_it;
  auto& session_uc = session_updates[imsi][session_id];
  if (!session->can_complete_termination(session_uc)) {
    return;  // error is logged in SessionState's complete_termination
  }
  auto termination_req = session->make_termination_request(session_uc);
  auto logging_cb = SessionReporter::get_terminate_logging_cb(termination_req);
  reporter_->report_terminate_session(termination_req, logging_cb);
  events_reporter_->session_terminated(imsi, session);

  // clear all metrics associated with this session
  session->clear_session_metrics();

  // Delete the session from SessionMap
  session_uc.is_session_ended = true;
  session_map[imsi].erase(*session_it);
  MLOG(MINFO) << session_id << " deleted from SessionMap";
  if (session_map[imsi].size() == 0) {
    session_map.erase(imsi);
    MLOG(MDEBUG) << "All sessions terminated for " << imsi;
  }
}

void LocalEnforcer::terminate_multiple_sessions(
    SessionMap& session_map,
    const std::unordered_set<ImsiAndSessionID>& sessions,
    SessionUpdate& session_update) {
  for (const auto& imsi_and_session_id : sessions) {
    const auto &imsi       = imsi_and_session_id.first,
               &session_id = imsi_and_session_id.second;
    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);

    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MWARNING) << "Session " << session_id
                     << " not found for termination";
      return;
    }
    auto& session = **session_it;
    auto& uc      = session_update[imsi][session_id];
    start_session_termination(session, true, uc);
  }
}

void LocalEnforcer::remove_rules_for_multiple_suspended_credit(
    SessionMap& session_map,
    std::unordered_set<
        ImsiSessionIDAndCreditkey, ImsiSessionIDAndCreditkey_hash>&
        suspended_credits,
    SessionUpdate& session_update) {
  for (const auto& suspended_credit : suspended_credits) {
    const auto& imsi       = suspended_credit.imsi;
    const auto& session_id = suspended_credit.session_id;
    const CreditKey& ckey  = suspended_credit.cKey;

    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MWARNING) << "Session " << session_id << " not found for suspension";
      return;
    }
    auto& session    = **session_it;
    auto& session_uc = session_update[imsi][session_id];
    remove_rules_for_suspended_credit(session, ckey, session_uc);
  }
}

void LocalEnforcer::remove_rules_for_suspended_credit(
    const std::unique_ptr<SessionState>& session, const CreditKey& ckey,
    SessionStateUpdateCriteria& session_uc) {
  MLOG(MWARNING) << "Suspending RG " << ckey << " for "
                 << session->get_session_id();
  // suspend this specific credit
  session->set_suspend_credit(ckey, true, session_uc);

  // check if this is the last credit. That will trigger Final Unit Actions
  session->suspend_service_if_needed_for_credit(ckey, session_uc);

  // Remove pipelined rules
  RulesToProcess rules_to_remove;
  session->get_rules_per_credit_key(ckey, &rules_to_remove, &session_uc);
  propagate_rule_updates_to_pipelined(
      session->get_config(), RulesToProcess{}, rules_to_remove, false);
}

void LocalEnforcer::add_rules_for_multiple_unsuspended_credit(
    SessionMap& session_map,
    std::unordered_set<
        ImsiSessionIDAndCreditkey, ImsiSessionIDAndCreditkey_hash>&
        unsuspended_credits,
    SessionUpdate& session_update) {
  for (const auto& suspended_credit : unsuspended_credits) {
    const auto& imsi       = suspended_credit.imsi;
    const auto& session_id = suspended_credit.session_id;
    const CreditKey& ckey  = suspended_credit.cKey;
    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MWARNING) << "Session " << session_id
                     << " not found for un-suspension";
      return;
    }
    auto& session    = **session_it;
    auto& session_uc = session_update[imsi][session_id];
    add_rules_for_unsuspended_credit(session, ckey, session_uc);
  }
}

void LocalEnforcer::add_rules_for_unsuspended_credit(
    const std::unique_ptr<SessionState>& session, const CreditKey& ckey,
    SessionStateUpdateCriteria& session_uc) {
  MLOG(MWARNING) << "Un-suspending RG " << ckey << " for "
                 << session->get_session_id();
  // unsuspend this credit
  session->set_suspend_credit(ckey, false, session_uc);

  //  add pipelined rules
  RulesToProcess rules_to_add;
  session->get_rules_per_credit_key(ckey, &rules_to_add, &session_uc);
  propagate_rule_updates_to_pipelined(
      session->get_config(), rules_to_add, RulesToProcess{}, false);
}

void LocalEnforcer::update_charging_credits(
    SessionMap& session_map, const UpdateSessionResponse& response,
    UpdateChargingCreditActions& actions, SessionUpdate& session_update) {
  for (const auto& credit_update_resp : response.responses()) {
    const std::string& imsi       = credit_update_resp.sid();
    const std::string& session_id = credit_update_resp.session_id();
    const auto ckey               = credit_update_resp.charging_key();
    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "Could not find session " << session_id
                   << " during charging update for RG " << ckey;
      continue;
    }
    auto& session    = **session_it;
    auto& uc         = session_update[imsi][session_id];
    auto result_code = credit_update_resp.result_code();

    // handle permanent failure -> terminate
    if (!credit_update_resp.success() &&
        DiameterCodeHandler::is_permanent_failure(result_code)) {
      MLOG(MERROR) << session_id << " Received permanent failure result code "
                   << result_code << " during update" << session_id;
      actions.sessions_to_terminate.insert(ImsiAndSessionID(imsi, session_id));
      continue;
    }

    const auto& credit_key(credit_update_resp);
    // We need to retrieve restrict_rules and is_final_action_state
    // prior to receiving charging credit as they will be updated.
    optional<FinalActionInfo> final_action_info =
        session->get_final_action_if_final_unit_state(credit_key);
    bool valid_credit =
        session->receive_charging_credit(credit_update_resp, uc);
    session->set_tgpp_context(credit_update_resp.tgpp_ctx(), uc);

    // handle suspended credits -> unsuspend
    if (valid_credit && session->is_credit_suspended(credit_key)) {
      actions.unsuspended_credits.insert(
          ImsiSessionIDAndCreditkey{imsi, session_id, ckey});
    }

    // handle other failures
    if (!credit_update_resp.success()) {
      // handle transient failure -> suspend
      if (DiameterCodeHandler::is_transient_failure(
              credit_update_resp.result_code())) {
        MLOG(MERROR) << session_id << " Received transient failure result code "
                     << result_code << " during update" << session_id;
        actions.suspended_credits.insert(
            ImsiSessionIDAndCreditkey{imsi, session_id, ckey});
      } else {
        // TODO: deal with other error codes
        MLOG(MERROR) << "Received Unimplemented result code " << result_code
                     << " for " << session_id
                     << " during update. Not action taken";
      }
    }

    // TODO: move it to actions vector
    if (final_action_info) {
      RulesToProcess gy_rules =
          session->remove_all_final_action_rules(*final_action_info, uc);
      if (!gy_rules.empty()) {
        auto config = session->get_config().common_context;
        // We need to cancel final unit action flows installed in pipelined here
        // following the reception of new charging credit.
        pipelined_client_->deactivate_flows_for_rules(
            session->get_imsi(), config.ue_ipv4(), config.ue_ipv6(),
            config.teids(), gy_rules, RequestOriginType::GY);
      }
    }
  }
}

void LocalEnforcer::update_monitoring_credits_and_rules(
    SessionMap& session_map, const UpdateSessionResponse& response,
    UpdateChargingCreditActions& actions, SessionUpdate& session_update) {
  // Since revalidation timer is session wide, we will only schedule one for
  // the entire session. The expectation is that if event triggers should be
  // included in all monitors or none.
  // To keep track of which timer is already tracked, we will have a set of
  // IMSIs that have pending re-validations
  std::unordered_set<ImsiAndSessionID> sessions_with_revalidation;
  for (const auto& usage_monitor_resp : response.usage_monitor_responses()) {
    const std::string& imsi       = usage_monitor_resp.sid();
    const std::string& session_id = usage_monitor_resp.session_id();
    const std::string& mkey = usage_monitor_resp.credit().monitoring_key();

    SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
    auto session_it = session_store_.find_session(session_map, criteria);
    if (!session_it) {
      MLOG(MERROR) << "Could not find session for " << session_id
                   << " during monitoring update for mkey " << mkey;
      continue;
    }
    auto& session               = **session_it;
    const SessionConfig& config = session->get_config();
    auto& uc                    = session_update[imsi][session_id];

    if (!usage_monitor_resp.success()) {
      handle_command_level_result_code_for_monitors(
          imsi, session_id, usage_monitor_resp.result_code(),
          actions.sessions_to_terminate);
      continue;
    }

    session->receive_monitor(usage_monitor_resp, uc);
    session->set_tgpp_context(usage_monitor_resp.tgpp_ctx(), uc);

    RulesToProcess pending_activation, pending_deactivation,
        pending_bearer_setup;
    RulesToSchedule pending_scheduling;

    session->process_rules_to_remove(
        usage_monitor_resp.rules_to_remove(), &pending_deactivation, &uc);

    session->process_rules_to_install(
        to_vec(usage_monitor_resp.static_rules_to_install()),
        to_vec(usage_monitor_resp.dynamic_rules_to_install()),
        &pending_activation, &pending_deactivation, &pending_bearer_setup,
        &pending_scheduling, &uc);

    handle_rule_scheduling(imsi, session_id, pending_scheduling);

    propagate_rule_updates_to_pipelined(
        config, pending_activation, pending_deactivation, false);

    if (terminate_on_wallet_exhaust() && is_wallet_exhausted(*session)) {
      actions.sessions_to_terminate.insert(ImsiAndSessionID(imsi, session_id));
    }

    auto imsi_and_session_id = ImsiAndSessionID(imsi, session_id);
    if (revalidation_required(usage_monitor_resp.event_triggers()) &&
        sessions_with_revalidation.count(imsi_and_session_id) == 0) {
      // All usage monitors under the same session will have the same event
      // trigger. See proto message / FeG for why. We will modify this input
      // logic later (Move event trigger out of UsageMonitorResponse), but
      // here we use a set to indicate whether a timer is already accounted
      // for.
      // Only schedule if no other revalidation timer was scheduled for
      // this session
      auto revalidation_time = usage_monitor_resp.revalidation_time();
      sessions_with_revalidation.insert(imsi_and_session_id);
      schedule_revalidation(*session, revalidation_time, uc);
    }

    if (config.common_context.rat_type() == TGPP_LTE) {
      const BearerUpdate update = session->get_dedicated_bearer_updates(
          pending_bearer_setup, pending_deactivation, &uc);
      propagate_bearer_updates_to_mme(update);
    }
  }
}

void LocalEnforcer::update_session_credits_and_rules(
    SessionMap& session_map, const UpdateSessionResponse& response,
    SessionUpdate& session_update) {
  // These subscribers will include any subscriber that received a permanent
  // diameter error code. Additionally, it will also include CWF sessions that
  // have run out of monitoring quota.
  UpdateChargingCreditActions actions;

  update_charging_credits(session_map, response, actions, session_update);

  update_monitoring_credits_and_rules(
      session_map, response, actions, session_update);

  remove_rules_for_multiple_suspended_credit(
      session_map, actions.suspended_credits, session_update);

  add_rules_for_multiple_unsuspended_credit(
      session_map, actions.unsuspended_credits, session_update);

  terminate_multiple_sessions(
      session_map, actions.sessions_to_terminate, session_update);
}

// handle_termination_from_access terminates the session that is
// associated with the given imsi and apn
bool LocalEnforcer::handle_termination_from_access(
    SessionMap& session_map, const std::string& imsi, const std::string& apn,
    SessionUpdate& session_updates) {
  SessionSearchCriteria criteria(imsi, IMSI_AND_APN, apn);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    return false;
  }
  auto& session                 = **session_it;
  const std::string& session_id = session->get_session_id();
  start_session_termination(session, false, session_updates[imsi][session_id]);
  return true;
}

bool LocalEnforcer::find_and_terminate_session(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, SessionUpdate& session_updates) {
  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    return false;
  }
  auto& session = **session_it;
  start_session_termination(session, true, session_updates[imsi][session_id]);
  return true;
}

bool LocalEnforcer::handle_abort_session(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, SessionUpdate& session_updates) {
  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    return false;
  }
  auto& session    = **session_it;
  auto& session_uc = session_updates[imsi][session_id];
  // Propagate rule removals to PipelineD and notify Access
  start_session_termination(session, true, session_uc);
  // ASRs do not require a CCR-T, this means we can immediately terminate
  // without waiting for final usage reports.
  events_reporter_->session_terminated(imsi, session);

  // Delete the session from SessionMap
  session_uc.is_session_ended = true;
  session_map[imsi].erase(*session_it);
  MLOG(MINFO) << session_id << " deleted from SessionMap";
  if (session_map[imsi].size() == 0) {
    session_map.erase(imsi);
    MLOG(MDEBUG) << "All sessions terminated for " << imsi;
  }
  return true;
}

void LocalEnforcer::handle_set_session_rules(
    SessionMap& session_map, const SessionRules& rules,
    SessionUpdate& session_update) {
  for (const auto& rules_per_sub : rules.rules_per_subscriber()) {
    const auto& imsi = rules_per_sub.imsi();
    auto session_it  = session_map.find(imsi);
    if (session_it == session_map.end()) {
      MLOG(MERROR) << "Could not find session for subscriber " << imsi
                   << " during set session rule update";
      return;
    }
    // Convert proto into a more convenient structure
    RuleSetBySubscriber rule_set_by_sub(rules_per_sub);

    for (const auto& session : session_it->second) {
      auto& uc           = session_update[imsi][session->get_session_id()];
      const auto& config = session->get_config();

      const auto& apn = config.common_context.apn();
      auto rule_set   = rule_set_by_sub.get_combined_rule_set_for_apn(apn);
      if (!rule_set) {
        // No rule change needed for this APN
        continue;
      }

      // Process the rule sets and get rules that need to be
      // activated/deactivated
      RulesToProcess pending_activation, pending_deactivation,
          pending_bearer_setup;
      session->apply_session_rule_set(
          *rule_set, &pending_activation, &pending_deactivation,
          &pending_bearer_setup, uc);

      // Propagate these rule changes to PipelineD and MME (if 4G)
      propagate_rule_updates_to_pipelined(
          config, pending_activation, pending_deactivation, false);
      if (config.common_context.rat_type() == TGPP_LTE) {
        const BearerUpdate update = session->get_dedicated_bearer_updates(
            pending_bearer_setup, pending_deactivation, &uc);
        propagate_bearer_updates_to_mme(update);
      }
    }
  }
}

ReAuthResult LocalEnforcer::init_charging_reauth(
    SessionMap& session_map, ChargingReAuthRequest request,
    SessionUpdate& session_update) {
  const std::string &imsi = request.sid(), &session_id = request.session_id();
  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "Could not find session " << session_id
                 << " for charging ReAuth";
    return ReAuthResult::SESSION_NOT_FOUND;
  }
  auto& session                  = **session_it;
  SessionStateUpdateCriteria& uc = session_update[imsi][session_id];
  switch (request.type()) {
    case ChargingReAuthRequest::SINGLE_SERVICE: {
      MLOG(MDEBUG) << "Initiating ReAuth of RG " << request.charging_key()
                   << " for session " << session_id;
      return session->reauth_key(CreditKey(request), uc);
    }
    case ChargingReAuthRequest::ENTIRE_SESSION: {
      MLOG(MDEBUG) << "Initiating ReAuth of all RGs for session " << session_id;
      return session->reauth_all(uc);
    }
    default:
      MLOG(MDEBUG) << "Received ChargingReAuthType " << request.type()
                   << " for " << session_id;
      return ReAuthResult::OTHER_FAILURE;
  }
}

void LocalEnforcer::init_policy_reauth(
    SessionMap& session_map, PolicyReAuthRequest request,
    PolicyReAuthAnswer& answer_out, SessionUpdate& session_update) {
  const std::string &imsi = request.imsi(), &session_id = request.session_id();
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    MLOG(MERROR) << "Could not find subscriber " << imsi
                 << " during policy ReAuth";
    answer_out.set_result(ReAuthResult::SESSION_NOT_FOUND);
    return;
  }
  // For empty session_id, apply changes to all sessions of subscriber
  // Changes are applied on a best-effort basis, so failures for one session
  // won't stop changes from being applied for subsequent sessions.
  if (session_id == "") {
    for (const auto& session : it->second) {
      init_policy_reauth_for_session(request, session, session_update);
    }
    answer_out.set_result(ReAuthResult::UPDATE_INITIATED);
    return;
  }

  SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "Could not find session " << session_id
                 << " during policy ReAuth";
    answer_out.set_result(ReAuthResult::SESSION_NOT_FOUND);
    return;
  }
  auto& session = **session_it;
  init_policy_reauth_for_session(request, session, session_update);
  answer_out.set_result(ReAuthResult::UPDATE_INITIATED);
}

void LocalEnforcer::init_policy_reauth_for_session(
    const PolicyReAuthRequest& request,
    const std::unique_ptr<SessionState>& session,
    SessionUpdate& session_update) {
  std::string imsi = request.imsi();
  SessionStateUpdateCriteria& uc =
      session_update[imsi][session->get_session_id()];

  receive_monitoring_credit_from_rar(request, session, uc);

  RulesToProcess pending_activation, pending_deactivation, pending_bearer_setup;
  RulesToSchedule pending_scheduling;

  MLOG(MDEBUG) << "Processing policy reauth for subscriber " << request.imsi();
  if (revalidation_required(request.event_triggers())) {
    schedule_revalidation(*session, request.revalidation_time(), uc);
  }

  session->process_rules_to_remove(
      request.rules_to_remove(), &pending_deactivation, &uc);

  session->process_rules_to_install(
      to_vec(request.rules_to_install()),
      to_vec(request.dynamic_rules_to_install()), &pending_activation,
      &pending_deactivation, &pending_bearer_setup, &pending_scheduling, &uc);

  handle_rule_scheduling(imsi, session->get_session_id(), pending_scheduling);

  propagate_rule_updates_to_pipelined(
      session->get_config(), pending_activation, pending_deactivation, false);

  if (terminate_on_wallet_exhaust() && is_wallet_exhausted(*session)) {
    start_session_termination(session, true, uc);
    return;
  }
  if (session->get_config().common_context.rat_type() == TGPP_LTE) {
    create_bearer(session, request, pending_bearer_setup);
  }
}

void LocalEnforcer::propagate_rule_updates_to_pipelined(
    const SessionConfig& config, const RulesToProcess& pending_activation,
    const RulesToProcess& pending_deactivation, bool always_send_activate) {
  const std::string& imsi = config.get_imsi();
  const auto ip_addr      = config.common_context.ue_ipv4();
  const auto ipv6_addr    = config.common_context.ue_ipv6();
  const Teids teids       = config.common_context.teids();
  // deactivate_flows_for_rules() should not be called when there is no rule
  // to deactivate, because pipelined deactivates all rules
  // when no rule is provided as the parameter
  if (!pending_deactivation.empty()) {
    pipelined_client_->deactivate_flows_for_rules(
        imsi, ip_addr, ipv6_addr, teids, pending_deactivation,
        RequestOriginType::GX);
  }
  if (always_send_activate || !pending_activation.empty()) {
    const auto ambr   = config.get_apn_ambr();
    const auto msisdn = config.common_context.msisdn();
    pipelined_client_->activate_flows_for_rules(
        imsi, ip_addr, ipv6_addr, teids, msisdn, ambr, pending_activation,
        std::bind(
            &LocalEnforcer::handle_activate_ue_flows_callback, this, imsi,
            ip_addr, ipv6_addr, teids, _1, _2));
  }
}

void LocalEnforcer::receive_monitoring_credit_from_rar(
    const PolicyReAuthRequest& request,
    const std::unique_ptr<SessionState>& session,
    SessionStateUpdateCriteria& uc) {
  UsageMonitoringUpdateResponse monitoring_credit;
  monitoring_credit.set_session_id(request.session_id());
  monitoring_credit.set_sid("IMSI" + request.session_id());
  monitoring_credit.set_success(true);
  UsageMonitoringCredit* credit = monitoring_credit.mutable_credit();

  for (const auto& usage_monitoring_credit :
       request.usage_monitoring_credits()) {
    credit->CopyFrom(usage_monitoring_credit);
    session->receive_monitor(monitoring_credit, uc);
  }
}

std::vector<StaticRuleInstall> LocalEnforcer::to_vec(
    const google::protobuf::RepeatedPtrField<magma::lte::StaticRuleInstall>
        static_rule_installs) {
  std::vector<StaticRuleInstall> out;
  for (const auto& install : static_rule_installs) {
    out.push_back(install);
  }
  return out;
}

std::vector<DynamicRuleInstall> LocalEnforcer::to_vec(
    const google::protobuf::RepeatedPtrField<magma::lte::DynamicRuleInstall>
        dynamic_rule_installs) {
  std::vector<DynamicRuleInstall> out;
  for (const auto& install : dynamic_rule_installs) {
    out.push_back(install);
  }
  return out;
}

void LocalEnforcer::process_rules_to_install(
    SessionState& session,
    const std::vector<StaticRuleInstall>& static_rule_installs,
    const std::vector<DynamicRuleInstall>& dynamic_rule_installs,
    RulesToProcess* pending_activation, RulesToProcess* pending_deactivation,
    RulesToProcess* pending_bearer_setup, RulesToSchedule* pending_scheduling,
    SessionStateUpdateCriteria* session_uc) {
  session.process_static_rule_installs(
      static_rule_installs, pending_activation, pending_deactivation,
      pending_bearer_setup, pending_scheduling, session_uc);
  session.process_dynamic_rule_installs(
      dynamic_rule_installs, pending_activation, pending_deactivation,
      pending_bearer_setup, pending_scheduling, session_uc);
}

void LocalEnforcer::handle_rule_scheduling(
    const std::string& imsi, const std::string& session_id,
    const RulesToSchedule& pending_scheduling) {
  for (const RuleToSchedule& pending_scheduling : pending_scheduling) {
    const PolicyType& p_type      = pending_scheduling.p_type;
    const PolicyAction& p_action  = pending_scheduling.p_action;
    const std::string& rule_id    = pending_scheduling.rule_id;
    const std::time_t& sched_time = pending_scheduling.scheduled_time;
    if (p_type == STATIC && p_action == ACTIVATE) {
      schedule_static_rule_activation(imsi, session_id, rule_id, sched_time);
    } else if (p_type == STATIC && p_action == DEACTIVATE) {
      schedule_static_rule_deactivation(imsi, session_id, rule_id, sched_time);
    } else if (p_type == DYNAMIC && p_action == ACTIVATE) {
      schedule_dynamic_rule_activation(imsi, session_id, rule_id, sched_time);
    } else if (p_type == DYNAMIC && p_action == ACTIVATE) {
      schedule_dynamic_rule_deactivation(imsi, session_id, rule_id, sched_time);
    }
  }
}

bool LocalEnforcer::revalidation_required(
    const google::protobuf::RepeatedField<int>& event_triggers) {
  auto it = std::find(
      event_triggers.begin(), event_triggers.end(), REVALIDATION_TIMEOUT);
  return it != event_triggers.end();
}

void LocalEnforcer::schedule_revalidation(
    SessionState& session, const google::protobuf::Timestamp& revalidation_time,
    SessionStateUpdateCriteria& uc) {
  // Add revalidation info to session and mark as pending
  session.add_new_event_trigger(REVALIDATION_TIMEOUT, uc);
  session.set_revalidation_time(revalidation_time, uc);
  const std::string& session_id = session.get_session_id();
  const std::string& imsi       = session.get_imsi();
  SessionRead req               = {imsi};
  auto delta = magma::time_difference_from_now(revalidation_time);
  MLOG(MINFO) << "Scheduling revalidation in " << delta.count() << "ms for "
              << session_id;
  evb_->runAfterDelay(
      [=] {
        MLOG(MINFO) << "Revalidation timeout! for " << session_id;
        auto session_map = session_store_.read_sessions(req);
        SessionSearchCriteria criteria(imsi, IMSI_AND_SESSION_ID, session_id);
        auto session_it = session_store_.find_session(session_map, criteria);
        if (!session_it) {
          MLOG(MERROR) << session_id << " not found for revalidation";
          return;
        }
        auto& session = **session_it;
        SessionUpdate update =
            SessionStore::get_default_session_update(session_map);
        auto& uc = update[imsi][session_id];
        session->mark_event_trigger_as_triggered(REVALIDATION_TIMEOUT, uc);
        session_store_.update_sessions(update);
      },
      delta.count());
}  // namespace magma

void LocalEnforcer::handle_activate_ue_flows_callback(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const Teids teids, Status status,
    ActivateFlowsResult resp) {
  if (status.ok()) {
    MLOG(MDEBUG) << "Pipelined add ue enf flow succeeded for " << imsi;
    return;
  }

  MLOG(MERROR) << "Could not activate rules for " << imsi
               << ", rpc failed: " << status.error_message()
               << ", terminating session...";

  SessionSearchCriteria criteria(imsi, IMSI_AND_UE_IPV4_OR_IPV6, ip_addr);
  auto session_map = session_store_.read_sessions({imsi});
  auto session_it  = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "Could not find session for " << imsi << " and " << ip_addr
                 << " when trying to terminate after PipelineD GRPC failure";
    return;
  }
  auto& session                = **session_it;
  const std::string session_id = session->get_session_id();

  // start_session_termination
  evb_->runInEventBaseThread([this, imsi, session_id] {
    auto session_map = session_store_.read_sessions({imsi});
    auto update      = SessionStore::get_default_session_update(session_map);
    bool success =
        find_and_terminate_session(session_map, imsi, session_id, update);
    if (!success) {
      MLOG(MERROR) << "Failed to start termination for " << session_id
                   << " after PipelineD GRPC failure";
    }
    session_store_.update_sessions(update);
  });
}  // namespace magma

void LocalEnforcer::handle_add_ue_mac_flow_callback(
    const SubscriberID& sid, const std::string& ue_mac_addr,
    const std::string& msisdn, const std::string& apn_mac_addr,
    const std::string& apn_name, Status status, FlowResponse resp) {
  if (status.ok() && resp.result() == resp.SUCCESS) {
    MLOG(MDEBUG) << "Pipelined add ue mac flow succeeded for " << ue_mac_addr;
    return;
  }

  if (!status.ok()) {
    MLOG(MERROR) << "Could not add ue mac flow, rpc failed with: "
                 << status.error_message() << ", retrying...";
  } else if (resp.result() == resp.FAILURE) {
    MLOG(MWARNING) << "Pipelined add ue mac flow failed, retrying...";
  }

  evb_->runAfterDelay(
      [=] {
        MLOG(MERROR) << "Could not activate ue mac flows for subscriber "
                     << sid.id() << ": " << status.error_message()
                     << ", retrying...";
        pipelined_client_->add_ue_mac_flow(
            sid, ue_mac_addr, msisdn, apn_mac_addr, apn_name,
            [ue_mac_addr](Status status, FlowResponse resp) {
              if (!status.ok()) {
                MLOG(MERROR) << "Could not activate flows for UE "
                             << ue_mac_addr << ": " << status.error_message();
              }
            });
      },
      retry_timeout_.count());
}

void LocalEnforcer::create_bearer(
    const std::unique_ptr<SessionState>& session,
    const PolicyReAuthRequest& request,
    const std::vector<RuleToProcess>& dynamic_rules) {
  const auto& config = session->get_config();
  if (!config.rat_specific_context.has_lte_context()) {
    MLOG(MWARNING) << "No LTE Session Context is specified for session";
    return;
  }
  const auto& lte_context = config.rat_specific_context.lte_context();
  if (!lte_context.has_qos_info() || !request.has_qos_info()) {
    MLOG(MDEBUG) << "Not creating bearer";
    return;
  }
  auto default_qci = QCI(lte_context.qos_info().qos_class_id());
  if (request.qos_info().qci() != default_qci) {
    MLOG(MDEBUG) << "QCI sent in RAR is different from default QCI";
    CreateBearerRequest req;
    req.mutable_sid()->CopyFrom(config.common_context.sid());
    req.set_ip_addr(config.common_context.ue_ipv4());
    // TODO ipv6_addrs missing address
    req.set_link_bearer_id(lte_context.bearer_id());

    auto req_policy_rules = req.mutable_policy_rules();
    for (const auto& to_process : dynamic_rules) {
      optional<PolicyType> p_type =
          session->get_policy_type(to_process.rule.id());
      if (p_type && *p_type == DYNAMIC) {
        req_policy_rules->Add()->CopyFrom(to_process.rule);
      }
    }
    spgw_client_->create_dedicated_bearer(req);
  }
  return;
}

void LocalEnforcer::update_ipfix_flow(
    const std::string& imsi, const SessionConfig& config,
    const uint64_t pdp_start_time) {
  MLOG(MDEBUG) << "Updating IPFIX flow for subscriber " << imsi;
  SubscriberID sid;
  sid.set_id(imsi);
  std::string apn_mac_addr;
  std::string apn_name;
  if (!parse_apn(config.common_context.apn(), apn_mac_addr, apn_name)) {
    MLOG(MWARNING) << "Failed mac/name parsiong for apn "
                   << config.common_context.apn();
    apn_mac_addr = "";
    apn_name     = config.common_context.apn();
  }

  // MacAddr is only relevant for WLAN
  const auto& rat_specific = config.rat_specific_context;
  std::string ue_mac_addr  = "11:11:11:11:11:11";
  if (rat_specific.has_wlan_context()) {
    ue_mac_addr = rat_specific.wlan_context().mac_addr();
  }
  pipelined_client_->update_ipfix_flow(
      sid, ue_mac_addr, config.common_context.msisdn(), apn_mac_addr, apn_name,
      pdp_start_time);
}

void LocalEnforcer::propagate_bearer_updates_to_mme(
    const BearerUpdate& updates) {
  // Order matters!!
  // First send delete requests and then create requests to
  // ensure that the final state is the desired one.
  if (updates.needs_deletion) {
    spgw_client_->delete_dedicated_bearer(updates.delete_req);
  }
  if (updates.needs_creation) {
    spgw_client_->create_dedicated_bearer(updates.create_req);
  }
}

void LocalEnforcer::handle_cwf_roaming(
    SessionMap& session_map, const std::string& imsi,
    const SessionConfig& config, SessionUpdate& session_update) {
  auto it = session_map.find(imsi);
  if (it != session_map.end()) {
    for (const auto& session : it->second) {
      auto& uc = session_update[imsi][session->get_session_id()];
      session->set_config(config);
      uc.is_config_updated = true;
      uc.updated_config    = session->get_config();
      // TODO Check for event triggers and send updates to the core if needed
      update_ipfix_flow(imsi, config, session->get_pdp_start_time());
    }
  }
}

bool LocalEnforcer::bind_policy_to_bearer(
    SessionMap& session_map, const PolicyBearerBindingRequest& request,
    SessionUpdate& session_update) {
  const std::string& imsi = request.sid().id();
  const auto bearer_id    = request.linked_bearer_id();

  SessionSearchCriteria criteria(imsi, IMSI_AND_BEARER, bearer_id);
  auto session_it = session_store_.find_session(session_map, criteria);
  if (!session_it) {
    MLOG(MERROR) << "Could not find session for " << imsi << " and bearer "
                 << bearer_id << " during update tunnelIds ";
    return false;
  }
  auto& session = **session_it;

  const auto& config = session->get_config();
  if (!config.rat_specific_context.has_lte_context()) {
    return false;  // not LTE
  }

  auto& uc = session_update[imsi][session->get_session_id()];

  if (!request.bearer_id()) {
    // if bearer_id is 0, the rule needs to be removed since we cannot honor the
    // QoS request
    remove_rule_due_to_bearer_creation_failure(
        *session, request.policy_rule_id(), uc);
    return false;
  }

  session->bind_policy_to_bearer(request, uc);
  install_rule_after_bearer_creation(*session, request);
  return true;
}

void LocalEnforcer::install_rule_after_bearer_creation(
    SessionState& session, const PolicyBearerBindingRequest& request) {
  optional<PolicyRule> op_rule =
      session.get_policy_definition(request.policy_rule_id());
  if (!op_rule) {
    MLOG(MWARNING) << "Policy " << request.policy_rule_id() << " not found for "
                   << session.get_session_id();
    return;
  }
  MLOG(MINFO) << "Installing " << request.policy_rule_id()
              << " into PipelineD for " << session.get_session_id()
              << " after allocating a dedicated bearer";
  RulesToProcess pending_activation{session.make_rule_to_process(*op_rule)};
  propagate_rule_updates_to_pipelined(
      session.get_config(), pending_activation, {}, false);
}

void LocalEnforcer::remove_rule_due_to_bearer_creation_failure(
    SessionState& session, const std::string& rule_id,
    SessionStateUpdateCriteria& uc) {
  MLOG(MINFO) << "Removing " << rule_id
              << " since we failed to create a dedicated bearer for it";
  auto policy_type = session.get_policy_type(rule_id);
  if (!policy_type) {
    MLOG(MERROR) << "Unable to remove rule " << rule_id
                 << " since it is not found";
    return;
  }

  optional<RuleToProcess> remove_info = {};
  switch (*policy_type) {
    case STATIC:
      if (rule_store_->get_rule(rule_id, nullptr)) {
        remove_info = session.deactivate_static_rule(rule_id, uc);
      }
      break;
    case DYNAMIC: {
      remove_info = session.remove_dynamic_rule(rule_id, nullptr, uc);
      break;
    }
  }
  if (remove_info) {
    RulesToProcess pending_deactivation{*remove_info};
    propagate_rule_updates_to_pipelined(
        session.get_config(), {}, pending_deactivation, false);
  }
}

std::unique_ptr<Timezone> LocalEnforcer::compute_access_timezone() {
  if (!SEND_ACCESS_TIMEZONE) {
    MLOG(MWARNING) << "send_access_timezone field is not set, not computing "
                      "timezone information";
    return nullptr;
  }
  Timezone timezone;
  time_t now           = std::time(NULL);
  auto const localtime = *std::localtime(&now);
  std::ostringstream os;
  os << std::put_time(&localtime, "%z");
  std::string s = os.str();
  // s is in ISO 8601 format: "±HHMM"
  if (s.size() < 5) {
    MLOG(MERROR) << "Failed to set access timezone";
    return nullptr;
  }
  int hours            = std::stoi(s.substr(0, 3), nullptr, 10);
  int minutes          = std::stoi(s[0] + s.substr(3), nullptr, 10);
  int total_offset_min = hours * 60 + minutes;
  MLOG(MINFO) << "Access timezone is UTC " << (total_offset_min >= 0 ? "+" : "")
              << total_offset_min << " minutes";
  timezone.set_offset_minutes(total_offset_min);
  return std::make_unique<Timezone>(timezone);
}

static void handle_command_level_result_code_for_monitors(
    const std::string& imsi, const std::string& session_id,
    const uint32_t result_code,
    std::unordered_set<ImsiAndSessionID>& sessions_to_terminate) {
  if (DiameterCodeHandler::is_permanent_failure(result_code)) {
    MLOG(MERROR) << session_id << " Received permanent failure result code "
                 << result_code << " during update" << session_id;
    sessions_to_terminate.insert(ImsiAndSessionID(imsi, session_id));
    return;
  }

  if (DiameterCodeHandler::is_transient_failure(result_code)) {
    MLOG(MERROR) << session_id << " Received transient failure result code "
                 << result_code << " during update" << session_id
                 << " Not action taken";
  }

  // TODO: deal with other error codes
  MLOG(MERROR) << "Received Unimplemented result code " << result_code
               << " for " << session_id << " during update. Not action taken";
}

static bool is_valid_mac_address(const char* mac) {
  int i = 0;
  int s = 0;

  while (*mac) {
    if (isxdigit(*mac)) {
      i++;
    } else if (*mac == '-') {
      if (i == 0 || i / 2 - 1 != s) {
        break;
      }
      ++s;
    } else {
      s = -1;
    }
    ++mac;
  }
  return (i == 12 && s == 5);
}

static bool parse_apn(
    const std::string& apn, std::string& mac_addr, std::string& name) {
  // Format is mac:name, if format check fails return failure
  // Format example - 1C-B9-C4-36-04-F0:Wifi-Offload-hotspot20
  if (apn.empty()) {
    return false;
  }
  auto split_location = apn.find(":");
  if (split_location <= 0) {
    return false;
  }
  auto mac = apn.substr(0, split_location);
  if (!is_valid_mac_address(mac.c_str())) {
    return false;
  }
  mac_addr = mac;
  // Allow empty name, spec is unclear on this
  name = apn.substr(split_location + 1, apn.size());
  return true;
}

static SubscriberQuotaUpdate make_subscriber_quota_update(
    const std::string& imsi, const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state) {
  SubscriberQuotaUpdate update;
  auto sid = update.mutable_sid();
  sid->set_id(imsi);
  update.set_mac_addr(ue_mac_addr);
  update.set_update_type(state);
  return update;
}

UpdateRequestsBySession::UpdateRequestsBySession(
    const UpdateSessionRequest& request) {
  for (const auto& charging_request : request.updates()) {
    const auto imsi = charging_request.common_context().sid().id();
    const auto id   = ImsiAndSessionID(imsi, charging_request.session_id());
    requests_by_id[id].charging_requests.push_back(charging_request);
  }
  for (const auto& monitor_request : request.usage_monitors()) {
    const auto id =
        ImsiAndSessionID(monitor_request.sid(), monitor_request.session_id());
    requests_by_id[id].monitor_requests.push_back(monitor_request);
  }
}

}  // namespace magma
