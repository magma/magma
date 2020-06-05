/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include <string>
#include <time.h>
#include <utility>
#include <vector>

#include <google/protobuf/repeated_field.h>
#include <google/protobuf/timestamp.pb.h>
#include <google/protobuf/util/time_util.h>
#include <grpcpp/channel.h>

#include "DiameterCodes.h"
#include "LocalEnforcer.h"
#include "ServiceRegistrySingleton.h"
#include "magma_logging.h"

namespace {

std::chrono::milliseconds time_difference_from_now(
    const google::protobuf::Timestamp& timestamp) {
  auto rule_time_sec =
      google::protobuf::util::TimeUtil::TimestampToSeconds(timestamp);
  auto now   = time(NULL);
  auto delta = std::max(rule_time_sec - now, 0L);
  std::chrono::seconds sec(delta);
  return std::chrono::duration_cast<std::chrono::milliseconds>(sec);
}
}  // namespace

namespace magma {

uint32_t LocalEnforcer::REDIRECT_FLOW_PRIORITY = 2000;

using google::protobuf::RepeatedPtrField;
using google::protobuf::util::TimeUtil;

// We will treat rule install/uninstall failures as all-or-nothing - that is,
// if we get a bad response from the pipelined client, we'll mark all the rules
// as failed in the response
static void mark_rule_failures(
    const bool activate_success, const bool deactivate_success,
    const PolicyReAuthRequest& request, PolicyReAuthAnswer& answer_out);
// For command level result codes, we will mark the subscriber to be terminated
// if the result code indicates a permanent failure.
static void handle_command_level_result_code(
    const std::string& imsi, const uint32_t result_code,
    std::unordered_set<std::string>& subscribers_to_terminate);
static bool is_valid_mac_address(const char* mac);
static int get_apn_split_locaion(const std::string& apn);
static bool parse_apn(
    const std::string& apn, std::string& mac_addr, std::string& name);

static SubscriberQuotaUpdate make_subscriber_quota_update(
    const std::string& imsi, const std::string& ue_mac_addr,
    const SubscriberQuotaUpdate_Type state);

LocalEnforcer::LocalEnforcer(
    std::shared_ptr<SessionReporter> reporter,
    std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
    std::shared_ptr<PipelinedClient> pipelined_client,
    std::shared_ptr<AsyncDirectorydClient> directoryd_client,
    std::shared_ptr<AsyncEventdClient> eventd_client,
    std::shared_ptr<SpgwServiceClient> spgw_client,
    std::shared_ptr<aaa::AAAClient> aaa_client,
    long session_force_termination_timeout_ms,
    long quota_exhaustion_termination_on_init_ms,
    magma::mconfig::SessionD mconfig)
    : reporter_(reporter),
      rule_store_(rule_store),
      session_store_(session_store),
      pipelined_client_(pipelined_client),
      directoryd_client_(directoryd_client),
      eventd_client_(eventd_client),
      spgw_client_(spgw_client),
      aaa_client_(aaa_client),
      session_force_termination_timeout_ms_(
          session_force_termination_timeout_ms),
      quota_exhaustion_termination_on_init_ms_(
          quota_exhaustion_termination_on_init_ms),
      retry_timeout_(1),
      mconfig_(mconfig){}

void LocalEnforcer::notify_new_report_for_sessions(
    SessionMap &session_map, SessionUpdate &session_update) {
  for (const auto &session_pair : session_map) {
    for (const auto &session : session_pair.second) {
      session->new_report(
          session_update[session_pair.first][session->get_session_id()]);
    }
  }
}

void LocalEnforcer::notify_finish_report_for_sessions(
    SessionMap& session_map, SessionUpdate& session_update) {
  // Iterate through sessions and notify that report has finished. Terminate any
  // sessions that can be terminated.
  std::vector<std::pair<std::string, std::string>> imsi_to_terminate;
  for (const auto &session_pair : session_map) {
    for (const auto &session : session_pair.second) {
      session->finish_report(
          session_update[session_pair.first][session->get_session_id()]);
      if (session->can_complete_termination()) {
        imsi_to_terminate.push_back(
            std::make_pair(session_pair.first, session->get_session_id()));
      }
    }
  }
  for (const auto& imsi_sid_pair : imsi_to_terminate) {
    SessionStateUpdateCriteria& update_criteria =
        session_update[imsi_sid_pair.first][imsi_sid_pair.second];
    complete_termination(
        session_map, imsi_sid_pair.first, imsi_sid_pair.second,
        update_criteria);
  }
}

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

bool LocalEnforcer::setup(
    SessionMap& session_map, const std::uint64_t& epoch,
    std::function<void(Status status, SetupFlowsResult)> callback) {
  std::vector<SessionState::SessionInfo> session_infos;
  std::vector<SubscriberQuotaUpdate> quota_updates;
  std::vector<std::string> msisdns;
  std::vector<std::string> ue_mac_addrs;
  std::vector<std::string> apn_mac_addrs;
  std::vector<std::string> apn_names;
  auto cwf = false;
  for (auto it = session_map.begin(); it != session_map.end(); it++) {
    for (const auto& session : it->second) {
      SessionState::SessionInfo session_info;
      session->get_session_info(session_info);
      session_infos.push_back(session_info);
      auto ue_mac_addr = session->get_config().mac_addr;
      ue_mac_addrs.push_back(ue_mac_addr);
      auto msisdn = session->get_config().msisdn;
      msisdns.push_back(msisdn);
      std::string apn_mac_addr;
      std::string apn_name;
      auto apn = session->get_config().apn;
      if (!parse_apn(apn, apn_mac_addr, apn_name)) {
        MLOG(MWARNING) << "Failed mac/name parsiong for apn " << apn;
        apn_mac_addr = "";
        apn_name     = apn;
      }
      apn_mac_addrs.push_back(apn_mac_addr);
      apn_names.push_back(apn_name);
      if (session->is_radius_cwf_session()) {
        cwf                          = true;
        SubscriberQuotaUpdate update = make_subscriber_quota_update(
            session_info.imsi, ue_mac_addr,
            session->get_subscriber_quota_state());
        quota_updates.push_back(update);
      }
    }
  }
  if (cwf) {
    return pipelined_client_->setup_cwf(
        session_infos, quota_updates, ue_mac_addrs, msisdns, apn_mac_addrs,
        apn_names, epoch, callback);
  } else {
    return pipelined_client_->setup_lte(session_infos, epoch, callback);
  }
}

void LocalEnforcer::sync_sessions_on_restart(std::time_t current_time) {
  std::unordered_set<std::string> imsis_to_terminate;

  auto session_map    = session_store_.read_all_sessions();
  auto session_update = SessionStore::get_default_session_update(session_map);
  // Update the sessions so that their rules match the current timestamp
  for (auto& it : session_map) {
    auto imsi = it.first;
    for (auto& session : it.second) {
      auto& uc = session_update[it.first][session->get_session_id()];
      // Reschedule termination if it was pending before
      if (session->get_state() == SESSION_TERMINATION_SCHEDULED) {
        imsis_to_terminate.insert(imsi);
      }

      session->sync_rules_to_time(current_time, uc);
      auto ip_addr = session->get_config().ue_ipv4;

      for (std::string rule_id : session->get_static_rules()) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        if (lifetime.deactivation_time > current_time) {
          auto rule_install = session->get_static_rule_install(rule_id, lifetime);
          schedule_static_rule_deactivation(imsi, rule_install);
        }
      }
      // Schedule rule activations / deactivations
      for (std::string rule_id : session->get_scheduled_static_rules()) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        auto rule_install = session->get_static_rule_install(rule_id, lifetime);
        schedule_static_rule_activation(imsi, ip_addr, rule_install);
        if (lifetime.deactivation_time > current_time) {
          schedule_static_rule_deactivation(imsi, rule_install);
        }
      }

      std::vector<std::string> rule_ids;
      session->get_dynamic_rules().get_rule_ids(rule_ids);
      for (std::string rule_id : rule_ids) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        if (lifetime.deactivation_time > current_time) {
          auto rule_install = session->get_dynamic_rule_install(rule_id, lifetime);
          schedule_dynamic_rule_deactivation(imsi, rule_install);
        }
      }
      session->get_scheduled_dynamic_rules().get_rule_ids(rule_ids);
      for (auto rule_id : rule_ids) {
        auto lifetime = session->get_rule_lifetime(rule_id);
        auto rule_install = session->get_dynamic_rule_install(rule_id, lifetime);
        schedule_dynamic_rule_activation(imsi, ip_addr, rule_install);
        if (lifetime.deactivation_time > current_time) {
          schedule_dynamic_rule_deactivation(imsi, rule_install);
        }
      }
    }
  }
  if (!imsis_to_terminate.empty()) {
    MLOG(MDEBUG) << "Scheduling termination for one or more IMSIs";
    schedule_termination(imsis_to_terminate);
  }
  bool success = session_store_.update_sessions(session_update);
  if (success) {
    MLOG(MDEBUG) << "Successfully synced sessions after restart";
  } else {
    MLOG(MERROR) << "Failed to sync sessions after restart";
  }
}

void LocalEnforcer::aggregate_records(
    SessionMap& session_map,
    const RuleRecordTable& records,
    SessionUpdate& session_update) {
  // unmark all credits
  notify_new_report_for_sessions(session_map, session_update);
  for (const RuleRecord &record : records.records()) {
    auto it = session_map.find(record.sid());
    if (it == session_map.end()) {
      MLOG(MERROR) << "Could not find session for " << record.sid()
                   << " during record aggregation";
      continue;
    }
    if (record.bytes_tx() > 0 || record.bytes_rx() > 0) {
      MLOG(MINFO) << "";
      MLOG(MINFO) << record.sid() << " used "
                  << record.bytes_tx() << " tx bytes and " << record.bytes_rx()
                  << " rx bytes for rule " << record.rule_id();
    }
    // Update sessions
    for (const auto& session : it->second) {
      SessionStateUpdateCriteria& uc =
          session_update[record.sid()][session->get_session_id()];
      session->add_used_credit(
          record.rule_id(), record.bytes_tx(), record.bytes_rx(), uc);
    }
  }
  notify_finish_report_for_sessions(session_map, session_update);
}

void LocalEnforcer::execute_actions(
    SessionMap& session_map,
    const std::vector<std::unique_ptr<ServiceAction>>& actions,
    SessionUpdate& session_update) {
  for (const auto& action_p : actions) {
    if (action_p->get_type() == TERMINATE_SERVICE) {
      terminate_service(
          session_map, action_p->get_imsi(), action_p->get_rule_ids(),
          action_p->get_rule_definitions(), session_update);
    } else if (action_p->get_type() == ACTIVATE_SERVICE) {
      pipelined_client_->activate_flows_for_rules(
          action_p->get_imsi(), action_p->get_ip_addr(),
          action_p->get_rule_ids(), action_p->get_rule_definitions());
    } else if (action_p->get_type() == REDIRECT) {
      // This is GY based REDIRECT, GX redirect will come in as a regular rule
      install_redirect_flow(action_p, session_update);
    } else if (action_p->get_type() == RESTRICT_ACCESS) {
      MLOG(MWARNING) << "RESTRICT_ACCESS mode is unsupported"
                     << ", will just terminate the service.";
      terminate_service(
          session_map, action_p->get_imsi(), action_p->get_rule_ids(),
          action_p->get_rule_definitions(), session_update);
    }
  }
}

// Terminates sessions that correspond to the given IMSI.
// (For session termination triggered by sessiond)
void LocalEnforcer::terminate_service(
    SessionMap &session_map, const std::string &imsi,
    const std::vector<std::string> &rule_ids,
    const std::vector<PolicyRule> &dynamic_rules,
    SessionUpdate &session_update) {
  pipelined_client_->deactivate_flows_for_rules(imsi, rule_ids, dynamic_rules,
                                                RequestOriginType::GX);
  MLOG(MINFO) << "Initiating session termination for " << imsi;

  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    MLOG(MWARNING) << "Could not find session with IMSI " << imsi
                   << " to terminate";
    return;
  }

  for (const auto &session : it->second) {
    if (session->is_terminating()) {
      // If the session is terminating already, do nothing.
      continue;
    }

    auto& update_criteria = session_update[imsi][session->get_session_id()];
    session->start_termination(update_criteria);

    // tell AAA service to terminate radius session if necessary
    auto config = session->get_config();
    if (session->is_radius_cwf_session()) {
      auto radius_session_id = config.radius_session_id;
      MLOG(MDEBUG) << "Asking AAA service to terminate session with "
                   << "Radius ID: " << radius_session_id << ", IMSI: " << imsi;
      aaa_client_->terminate_session(radius_session_id, imsi);

      MLOG(MDEBUG) << "Deleting UE MAC flow for subscriber " << imsi;
      auto mac_addr = config.mac_addr;
      SubscriberID sid;
      sid.set_id(imsi);
      bool delete_ue_mac_flow_success =
          pipelined_client_->delete_ue_mac_flow(sid, mac_addr);
      if (!delete_ue_mac_flow_success) {
        MLOG(MERROR) << "Failed to delete UE MAC flow for subscriber " << imsi;
      }
      MLOG(MDEBUG) << "Setting subscriber quota state as TERMINATE "
                   << "for subscriber " << imsi;
      session->set_subscriber_quota_state(
          SubscriberQuotaUpdate_Type_TERMINATE, update_criteria);
      report_subscriber_state_to_pipelined(
          imsi, mac_addr, SubscriberQuotaUpdate_Type_TERMINATE);
    } else {
      // Deleting the PDN session by triggering network issued default bearer
      // deactivation
      spgw_client_->delete_default_bearer(
        imsi, config.ue_ipv4, config.bearer_id);
    }

    std::string session_id = session->get_session_id();
    // The termination should be completed when aggregated usage record no
    // longer
    // includes the imsi. If this has not occurred after the timeout, force
    // terminate the session.
    evb_->runAfterDelay(
        [this, imsi, session_id] {
          MLOG(MDEBUG) << "Checking if termination has to be forced for "
                       << imsi;
          SessionRead req = {imsi};
          auto session_map = session_store_.read_sessions_for_deletion(req);
          auto session_update =
              SessionStore::get_default_session_update(session_map);
          if (session_update[imsi].find(session_id) !=
              session_update[imsi].end()) {
            auto& update_criteria = session_update[imsi][session_id];
            complete_termination(
                session_map, imsi, session_id, update_criteria);
            bool end_success = session_store_.update_sessions(session_update);
            if (end_success) {
              MLOG(MDEBUG) << "Ended session " << imsi
                           << " with session_id: " << session_id;
            } else {
              MLOG(MERROR)
                  << "Failed to update SessionStore with ended session " << imsi
                  << " and session_id: " << session_id;
            }
          } else {
            MLOG(MDEBUG) << "Not forcing termination for session " << imsi
                         << " and session_id: " << session_id
                         << " as it has already terminated.";
          }
        },
        session_force_termination_timeout_ms_);
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
  }
}

static PolicyRule create_redirect_rule(
    const std::unique_ptr<ServiceAction>& action) {
  PolicyRule redirect_rule;
  redirect_rule.set_id("redirect");
  redirect_rule.set_priority(LocalEnforcer::REDIRECT_FLOW_PRIORITY);
  action->get_credit_key().set_rule(&redirect_rule);

  RedirectInformation* redirect_info = redirect_rule.mutable_redirect();
  redirect_info->set_support(RedirectInformation_Support_ENABLED);

  auto redirect_server = action->get_redirect_server();
  redirect_info->set_address_type(
      address_type_converter(redirect_server.redirect_address_type()));
  redirect_info->set_server_address(redirect_server.redirect_server_address());

  return redirect_rule;
}

void LocalEnforcer::install_redirect_flow(
    const std::unique_ptr<ServiceAction>& action,
    SessionUpdate &session_update) {
  std::vector<std::string> static_rules;
  std::vector<PolicyRule> dynamic_rules{create_redirect_rule(action)};
  const std::string &imsi = action->get_imsi();

  auto request = directoryd_client_->get_directoryd_ip_field(
      imsi, [this, imsi, static_rules, dynamic_rules]
        (Status status, DirectoryField resp) {

        if (!status.ok()) {
          MLOG(MERROR) << "Could not fetch subscriber " << imsi << " ip, "
                       << "redirection fails, error: "
                       << status.error_message();
        } else {
          auto session_map = session_store_.read_sessions(SessionRead{imsi});
          auto it = session_map.find(imsi);
          if (it == session_map.end()) {
            MLOG(MDEBUG) << "Session for IMSI " << imsi << " not found";
            return;
          }

          auto session_update =
            session_store_.get_default_session_update(session_map);
          pipelined_client_->add_gy_final_action_flow(
            imsi, resp.value(), static_rules, dynamic_rules);

          for (const auto &session : it->second) {
            auto &uc = session_update[imsi][session->get_session_id()];
            session->insert_gy_dynamic_rule(dynamic_rules.front(), uc);
          }
          auto update_success =
            session_store_.update_sessions(session_update);
        }
      });
}

UpdateSessionRequest LocalEnforcer::collect_updates(
    SessionMap& session_map,
    std::vector<std::unique_ptr<ServiceAction>>& actions,
    SessionUpdate& session_update, const bool force_update) const {
  UpdateSessionRequest request;
  for (const auto& session_pair : session_map) {
    for (const auto& session : session_pair.second) {
      std::string imsi     = session_pair.first;
      std::string sid      = session->get_session_id();
      auto& update_criteria = session_update[imsi][sid];
      session->get_updates(request, &actions, update_criteria, force_update);
    }
  }
  return request;
}

void LocalEnforcer::reset_updates(
    SessionMap& session_map, const UpdateSessionRequest& failed_request) {
  for (const auto& update : failed_request.updates()) {
    auto it = session_map.find(update.sid());
    if (it == session_map.end()) {
      MLOG(MERROR) << "Could not reset credit for IMSI " << update.sid()
                   << " because it couldn't be found";
      return;
    }

    for (const auto& session : it->second) {
      // When updates are reset, they aren't written back into SessionStore,
      // so we can just put in a default UpdateCriteria
      auto uc = get_default_update_criteria();
      session->get_charging_pool().reset_reporting_credit(
          CreditKey(update.usage()), uc);
    }
  }
  for (const auto& update : failed_request.usage_monitors()) {
    auto it = session_map.find(update.sid());
    if (it == session_map.end()) {
      MLOG(MERROR) << "Could not reset credit for IMSI " << update.sid()
                   << " because it couldn't be found";
      return;
    }

    for (const auto& session : it->second) {
      // When updates are reset, they aren't written back into SessionStore,
      // so we can just put in a default UpdateCriteria
      auto uc = get_default_update_criteria();
      session->get_monitor_pool().reset_reporting_credit(
          update.update().monitoring_key(), uc);
    }
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
  }
  return true;
}

void LocalEnforcer::schedule_static_rule_activation(
    const std::string& imsi, const std::string& ip_addr,
    const StaticRuleInstall& static_rule) {
  std::vector<std::string> static_rules{static_rule.rule_id()};
  std::vector<PolicyRule> dynamic_rules;

  auto delta = time_difference_from_now(static_rule.activation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " static rule "
               << static_rule.rule_id() << " activation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
        std::move([=] {
          auto session_map = session_store_.read_sessions(SessionRead{imsi});
          auto session_update =
              session_store_.get_default_session_update(session_map);
          pipelined_client_->activate_flows_for_rules(
              imsi, ip_addr, static_rules, dynamic_rules);
          auto it = session_map.find(imsi);
          if (it == session_map.end()) {
            MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                           << "during installation of static rule "
                           << static_rule.rule_id();
          } else {
            for (const auto& session : it->second) {
              if (session->get_config().ue_ipv4 == ip_addr) {
                auto& uc = session_update[imsi][session->get_session_id()];
                session->install_scheduled_static_rule(
                    static_rule.rule_id(), uc);
              }
            }
            auto update_success =
                session_store_.update_sessions(session_update);
          }
        }),
        delta);
  });
}

void LocalEnforcer::schedule_dynamic_rule_activation(
    const std::string& imsi, const std::string& ip_addr,
    const DynamicRuleInstall& dynamic_rule) {
  std::vector<std::string> static_rules;
  std::vector<PolicyRule> dynamic_rules{dynamic_rule.policy_rule()};

  auto delta = time_difference_from_now(dynamic_rule.activation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " dynamic rule "
               << dynamic_rule.policy_rule().id() << " activation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
        std::move([=] {
          auto session_map = session_store_.read_sessions(SessionRead{imsi});
          auto session_update =
              session_store_.get_default_session_update(session_map);
          pipelined_client_->activate_flows_for_rules(
              imsi, ip_addr, static_rules, dynamic_rules);
          auto it = session_map.find(imsi);
          if (it == session_map.end()) {
            MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                           << "during installation of dynamic rule "
                           << dynamic_rule.policy_rule().id();
          } else {
            for (const auto& session : it->second) {
              if (session->get_config().ue_ipv4 == ip_addr) {
                auto& uc = session_update[imsi][session->get_session_id()];
                session->install_scheduled_dynamic_rule(
                    dynamic_rule.policy_rule().id(), uc);
              }
            }
            auto update_success =
                session_store_.update_sessions(session_update);
          }
        }),
        delta);
  });
}

void LocalEnforcer::schedule_static_rule_deactivation(
    const std::string& imsi, const StaticRuleInstall& static_rule) {
  std::vector<std::string> static_rules{static_rule.rule_id()};
  std::vector<PolicyRule> dynamic_rules;

  auto delta = time_difference_from_now(static_rule.deactivation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " static rule "
               << static_rule.rule_id() << " deactivation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
        std::move([=] {
          auto session_map = session_store_.read_sessions(SessionRead{imsi});
          auto session_update =
              session_store_.get_default_session_update(session_map);
          pipelined_client_->deactivate_flows_for_rules(
              imsi, static_rules, dynamic_rules, RequestOriginType::GX);
          auto it = session_map.find(imsi);
          if (it == session_map.end()) {
            MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                           << "during removal of static rule "
                           << static_rule.rule_id();
          } else {
            for (const auto& session : it->second) {
              auto& uc = session_update[imsi][session->get_session_id()];
              if (!session->deactivate_static_rule(static_rule.rule_id(), uc))
                MLOG(MWARNING)
                    << "Could not find rule " << static_rule.rule_id()
                    << "for IMSI " << imsi << " during static rule removal";
            }
            auto update_success =
                session_store_.update_sessions(session_update);
          }
        }),
        delta);
  });
}

void LocalEnforcer::schedule_dynamic_rule_deactivation(
    const std::string& imsi, DynamicRuleInstall& dynamic_rule) {
  std::vector<std::string> static_rules;
  PolicyRule* policy = dynamic_rule.release_policy_rule();
  std::vector<PolicyRule> dynamic_rules{*policy};

  auto delta = time_difference_from_now(dynamic_rule.deactivation_time());
  MLOG(MDEBUG) << "Scheduling subscriber " << imsi << " dynamic rule "
               << dynamic_rule.policy_rule().id() << " deactivation in "
               << (delta.count() / 1000) << " secs";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
        std::move([=] {
          auto session_map = session_store_.read_sessions(SessionRead{imsi});
          auto session_update =
              session_store_.get_default_session_update(session_map);
          pipelined_client_->deactivate_flows_for_rules(
              imsi, static_rules, dynamic_rules, RequestOriginType::GX);
          auto it = session_map.find(imsi);
          if (it == session_map.end()) {
            MLOG(MWARNING) << "Could not find session for IMSI " << imsi
                           << "during removal of dynamic rule "
                           << dynamic_rule.policy_rule().id();
          } else {
            PolicyRule rule_dont_care;
            for (const auto& session : it->second) {
              auto& uc = session_update[imsi][session->get_session_id()];
              session->remove_dynamic_rule(
                  dynamic_rule.policy_rule().id(), &rule_dont_care, uc);
            }
            auto update_success =
                session_store_.update_sessions(session_update);
          }
        }),
        delta);
  });
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

// return true if any credit unit is valid and has non-zero volume
static bool contains_credit(const GrantedUnits& gsu) {
  return (gsu.total().is_valid() && gsu.total().volume() > 0) ||
         (gsu.tx().is_valid() && gsu.tx().volume() > 0) ||
         (gsu.rx().is_valid() && gsu.rx().volume() > 0);
}

bool LocalEnforcer::handle_session_init_rule_updates(
    SessionMap& session_map, const std::string& imsi,
    SessionState& session_state, const CreateSessionResponse& response,
    std::unordered_set<uint32_t>& charging_credits_received) {
  auto ip_addr = session_state.get_config().ue_ipv4;

  RulesToProcess rules_to_activate;
  RulesToProcess rules_to_deactivate;

  // Can use a default UpdateCriteria since SessionStore's create and update
  // methods are separate.
  std::vector<StaticRuleInstall> static_rule_installs =
      to_vec(response.static_rules());
  std::vector<DynamicRuleInstall> dynamic_rule_installs =
      to_vec(response.dynamic_rules());
  filter_rule_installs(
      static_rule_installs, dynamic_rule_installs, charging_credits_received);

  auto uc = get_default_update_criteria();
  process_rules_to_install(
      session_state, imsi, static_rule_installs, dynamic_rule_installs,
      rules_to_activate, rules_to_deactivate, uc);

  // activate_flows_for_rules() should be called even if there is no rule to
  // activate, because pipelined activates a "drop all packet" rule
  // when no rule is provided as the parameter.
  bool activate_success = pipelined_client_->activate_flows_for_rules(
      imsi, ip_addr, rules_to_activate.static_rules,
      rules_to_activate.dynamic_rules);

  // deactivate_flows_for_rules() should not be called when there is no rule
  // to deactivate, because pipelined deactivates all rules
  // when no rule is provided as the parameter
  bool deactivate_success = true;
  if (rules_to_process_is_not_empty(rules_to_deactivate)) {
    deactivate_success = pipelined_client_->deactivate_flows_for_rules(
        imsi, rules_to_deactivate.static_rules,
        rules_to_deactivate.dynamic_rules, RequestOriginType::GX);
  }

  return activate_success && deactivate_success;
}

bool LocalEnforcer::init_session_credit(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id, const SessionConfig& cfg,
    const CreateSessionResponse& response) {
  auto session_state = new SessionState(
      imsi, session_id, response.session_id(), cfg, *rule_store_,
      response.tgpp_ctx());

  std::unordered_set<uint32_t> charging_credits_received;
  for (const auto& credit : response.credits()) {
    auto uc = get_default_update_criteria();
    session_state->get_charging_pool().receive_credit(credit, uc);
    if (credit.success() && contains_credit(credit.credit().granted_units())) {
      charging_credits_received.insert(credit.charging_key());
    }
  }
  // We don't have to check 'success' field for monitors because command level
  // errors are handled in session proxy for the init exchange
  // Since revalidation timer is session wide, we will only schedule one for
  // the entire session

  // TODO [REMOVE] We will log error in case we get an unexpected input here
  // The expectation is that if event triggers should be included in all
  // monitors or none
  auto revalidation_scheduled = false;
  google::protobuf::Timestamp revalidation_time;

  for (const auto& monitor : response.usage_monitors()) {
    if (!revalidation_scheduled && revalidation_required(monitor.event_triggers())) {
      revalidation_scheduled = true;
      revalidation_time = monitor.revalidation_time();
    } else if (revalidation_scheduled &&
          !revalidation_required(monitor.event_triggers())) {
        // TODO [Remove This Clause]
        auto time_from_now_ms = time_difference_from_now(revalidation_time);
        MLOG(MWARNING) << imsi << " received some usage monitors with a "
                       << "revalidation timer in " << time_from_now_ms.count()
                       << " and some without";
      }
    auto uc = get_default_update_criteria();
    session_state->get_monitor_pool().receive_credit(monitor, uc);
  }

  auto rule_update_success = handle_session_init_rule_updates(
      session_map, imsi, *session_state, response, charging_credits_received);

  if (session_state->is_radius_cwf_session()) {
    using namespace std::placeholders;
    MLOG(MERROR) << "Adding UE MAC flow for subscriber " << imsi;
    SubscriberID sid;
    sid.set_id(imsi);
    std::string apn_mac_addr;
    std::string apn_name;
    if (!parse_apn(cfg.apn, apn_mac_addr, apn_name)) {
      MLOG(MWARNING) << "Failed mac/name parsing for apn " << cfg.apn;
      apn_mac_addr = "";
      apn_name     = cfg.apn;
    }
    auto ue_mac_addr             = session_state->get_config().mac_addr;
    bool add_ue_mac_flow_success = pipelined_client_->add_ue_mac_flow(
        sid, ue_mac_addr, cfg.msisdn, apn_mac_addr, apn_name,
        std::bind(
          &LocalEnforcer::handle_add_ue_mac_flow_callback,
          this, sid, ue_mac_addr, cfg.msisdn, apn_mac_addr, apn_name, _1, _2));
    if (!add_ue_mac_flow_success) {
      MLOG(MERROR) << "Failed to add UE MAC flow for subscriber " << imsi;
    }
    if (terminate_on_wallet_exhaust()) {
      handle_session_init_subscriber_quota_state(
        session_map, imsi, *session_state);
    }
  }

  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    // First time a session is created for IMSI
    MLOG(MDEBUG) << "First session for IMSI " << imsi << " with session ID "
                 << session_id;
    session_map[imsi] = std::vector<std::unique_ptr<SessionState>>();
  }
  session_map[imsi].push_back(
      std::move(std::unique_ptr<SessionState>(session_state)));

  if (session_state->is_radius_cwf_session() == false) {
    session_events::session_created(eventd_client_, imsi, session_id);
  }
  if (revalidation_scheduled) {
    // TODO This might not work since the session is not initialized properly
    // at this point
    schedule_revalidation(imsi, revalidation_time);
  }

  return rule_update_success;
}

bool LocalEnforcer::terminate_on_wallet_exhaust() {
  return mconfig_.has_wallet_exhaust_detection()
        && mconfig_.wallet_exhaust_detection().terminate_on_exhaust();
}

bool LocalEnforcer::is_wallet_exhausted(SessionState& session_state) {
  switch(mconfig_.wallet_exhaust_detection().method()) {
    case magma::mconfig::WalletExhaustDetection_Method_GxTrackedRules:
      return !session_state.active_monitored_rules_exist();
    default:
      MLOG(MWARNING) << "This method is not yet supported...";
      return false;
  }
}

void LocalEnforcer::handle_session_init_subscriber_quota_state(
    SessionMap& session_map, const std::string& imsi,
    SessionState& session_state) {
  auto ue_mac_addr = session_state.get_config().mac_addr;
  bool is_exhausted = is_wallet_exhausted(session_state);
  
  // This method only used for session creation and not updates, so
  // UpdateCriteria is unused.
  auto _ = get_default_update_criteria();
  if (is_exhausted) {
    MLOG(MINFO) << imsi << " Subscriber wallet is exhausted, setting subscriber"
                << " quota state as NO_QUOTA";
    session_state.set_subscriber_quota_state(
        SubscriberQuotaUpdate_Type_NO_QUOTA, _);
    report_subscriber_state_to_pipelined(
        imsi, ue_mac_addr, SubscriberQuotaUpdate_Type_NO_QUOTA);
    // Schedule a session termination for a configured number of seconds after
    // session create
    session_state.mark_as_awaiting_termination(_);
    MLOG(MINFO) << imsi << " Scheduling session for subscriber "
                << "to be terminated in "
                << quota_exhaustion_termination_on_init_ms_ << " ms";
    auto imsi_set = std::unordered_set<std::string>{imsi};
    schedule_termination(imsi_set);
    return;
  }

  // Valid Quota
  MLOG(MINFO) << imsi << " Setting subscriber quota state as VALID "
              << "for subscriber";
  session_state.set_subscriber_quota_state(
    SubscriberQuotaUpdate_Type_VALID_QUOTA, _);
  report_subscriber_state_to_pipelined(
      imsi, ue_mac_addr, SubscriberQuotaUpdate_Type_VALID_QUOTA);
  return;
}

void LocalEnforcer::schedule_termination(std::unordered_set<std::string>& imsis) {
    evb_->runAfterDelay(
      [this, imsis] {
        SessionRead req;
        req.insert(imsis.begin(), imsis.end());
        auto session_map = session_store_.read_sessions_for_deletion(req);

        SessionUpdate session_update =
            SessionStore::get_default_session_update(session_map);

        terminate_multiple_services(session_map, imsis, session_update);
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
  bool add_subscriber_quota_state_success =
      pipelined_client_->update_subscriber_quota_state(
          std::vector<SubscriberQuotaUpdate>{update});
  if (!add_subscriber_quota_state_success) {
    MLOG(MERROR) << "Failed to update subscriber's quota state to " << state
                 << " for subscriber " << imsi;
  }
}

void LocalEnforcer::complete_termination(
    SessionMap& session_map, const std::string& imsi,
    const std::string& session_id,
    SessionStateUpdateCriteria& update_criteria) {
  // If the session cannot be found in session_map, or a new session has
  // already begun, do nothing.
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    // Session is already deleted, or new session already began, ignore.
    MLOG(MDEBUG) << "Could not find session for IMSI " << imsi
                 << " and session ID " << session_id
                 << ". Skipping termination.";
    return;
  }

  for (auto session_it = it->second.begin(); session_it != it->second.end();
       ++session_it) {
    if ((*session_it)->get_session_id() == session_id) {
      // Complete session termination and remove session from session_map.
      (*session_it)->complete_termination(*reporter_, update_criteria);
      // Send to eventd
      if ((*session_it)->is_radius_cwf_session() == false) {
        session_events::session_terminated(eventd_client_, *session_it);
      }
      // We break the loop below, but for extra code safety in case
      // someone removes the break in the future, adjust the iterator
      // after erasing the element
      update_criteria.is_session_ended = true;
      it->second.erase(session_it--);
      MLOG(MDEBUG) << "Successfully terminated session for " << imsi
                   << "session ID " << session_id;
      // No session left for this IMSI
      if (it->second.size() == 0) {
        session_map.erase(imsi);
        MLOG(MDEBUG) << "All sessions terminated for " << imsi;
      }
      break;
    }
  }
}

bool LocalEnforcer::rules_to_process_is_not_empty(
    const RulesToProcess& rules_to_process) {
  return rules_to_process.static_rules.size() > 0 ||
         rules_to_process.dynamic_rules.size() > 0;
}

void LocalEnforcer::terminate_multiple_services(
    SessionMap& session_map, const std::unordered_set<std::string>& imsis,
    SessionUpdate& session_update) {
  for (const auto& imsi : imsis) {
    auto it = session_map.find(imsi);
    if (it == session_map.end()) {
      continue;
    }
    for (const auto& session : it->second) {
      RulesToProcess rules;
      populate_rules_from_session_to_remove(imsi, session, rules);
      terminate_service(
          session_map, imsi, rules.static_rules, rules.dynamic_rules,
          session_update);
    }
  }
}

void LocalEnforcer::update_charging_credits(
    SessionMap& session_map, const UpdateSessionResponse& response,
    std::unordered_set<std::string>& subscribers_to_terminate,
    SessionUpdate& session_update) {
  for (const auto& credit_update_resp : response.responses()) {
    const std::string& imsi = credit_update_resp.sid();

    if (!credit_update_resp.success()) {
      handle_command_level_result_code(
          imsi, credit_update_resp.result_code(), subscribers_to_terminate);
      continue;
    }

    auto it = session_map.find(imsi);
    if (it == session_map.end()) {
      MLOG(MERROR) << "Could not find session for IMSI "
                   << credit_update_resp.sid() << " during update";
      continue;
    }
    for (const auto& session : it->second) {
      std::string sid                             = session->get_session_id();
      SessionStateUpdateCriteria& update_criteria = session_update[imsi][sid];
      session->get_charging_pool().receive_credit(
          credit_update_resp, update_criteria);
      session->set_tgpp_context(credit_update_resp.tgpp_ctx(), update_criteria);
      SessionState::SessionInfo info;
      std::vector<PolicyRule> gy_rules_to_deactivate;
      session->get_session_info(info);
      for (const auto &rule : info.gy_dynamic_rules) {
        PolicyRule dy_rule;
        auto &uc = session_update[imsi][session->get_session_id()];
        bool is_dynamic =
            session->remove_gy_dynamic_rule(rule.id(), &dy_rule, uc);
        if (is_dynamic) {
          gy_rules_to_deactivate.push_back(dy_rule);
        }
      }

      if (!gy_rules_to_deactivate.empty()) {
        std::vector<std::string> static_rules;
        bool deactivate_success = pipelined_client_->deactivate_flows_for_rules(
            imsi, static_rules, gy_rules_to_deactivate, RequestOriginType::GY);
      }
    }
  }
}

void LocalEnforcer::update_monitoring_credits_and_rules(
    SessionMap& session_map, const UpdateSessionResponse& response,
    std::unordered_set<std::string>& subscribers_to_terminate,
    SessionUpdate& session_update) {
  // Since revalidation timer is session wide, we will only schedule one for
  // the entire session. The expectation is that if event triggers should be
  // included in all monitors or none.
  // To keep track of which timer is already tracked, we will have a set of
  // IMSIs that have pending re-validations
  std::unordered_set<std::string> imsis_with_revalidation;
  for (const auto& usage_monitor_resp : response.usage_monitor_responses()) {
    const std::string& imsi = usage_monitor_resp.sid();

    if (!usage_monitor_resp.success()) {
      handle_command_level_result_code(
          imsi, usage_monitor_resp.result_code(), subscribers_to_terminate);
      continue;
    }

    auto it = session_map.find(imsi);
    if (it == session_map.end()) {
      MLOG(MERROR) << "Could not find session for IMSI " << imsi
                   << " during update";
      continue;
    }

    for (const auto& session : it->second) {
      auto& update_criteria = session_update[imsi][session->get_session_id()];
      session->get_monitor_pool().receive_credit(
          usage_monitor_resp, update_criteria);
      session->set_tgpp_context(usage_monitor_resp.tgpp_ctx(), update_criteria);

      RulesToProcess rules_to_activate;
      RulesToProcess rules_to_deactivate;

      process_rules_to_remove(
          imsi, session, usage_monitor_resp.rules_to_remove(),
          rules_to_deactivate, update_criteria);

      process_rules_to_install(
          *session, imsi, to_vec(usage_monitor_resp.static_rules_to_install()),
          to_vec(usage_monitor_resp.dynamic_rules_to_install()),
          rules_to_activate, rules_to_deactivate, update_criteria);

      auto ip_addr            = session->get_config().ue_ipv4;
      bool deactivate_success = true;
      bool activate_success   = true;

      if (rules_to_process_is_not_empty(rules_to_deactivate)) {
        // TODO: modify the SessionUpdate
        deactivate_success = pipelined_client_->deactivate_flows_for_rules(
            imsi, rules_to_deactivate.static_rules,
            rules_to_deactivate.dynamic_rules, RequestOriginType::GX);
      }

      if (rules_to_process_is_not_empty(rules_to_activate)) {
        // TODO: modify the SessionUpdate
        activate_success = pipelined_client_->activate_flows_for_rules(
            imsi, ip_addr, rules_to_activate.static_rules,
            rules_to_activate.dynamic_rules);
      }

      // TODO If either deactivating/activating rules fail, sessiond should
      // manage the failed states. In the meantime, we will just log error for
      // now.
      if (!deactivate_success) {
        MLOG(MERROR) << "Could not deactivate flows for IMSI " << imsi
                     << "during update";
      }

      if (!activate_success) {
        MLOG(MERROR) << "Could not activate flows for IMSI " << imsi
                     << "during update";
      }

      if (terminate_on_wallet_exhaust() && is_wallet_exhausted(*session)) {
        subscribers_to_terminate.insert(imsi);
      }

      if (revalidation_required(usage_monitor_resp.event_triggers()) &&
            imsis_with_revalidation.count(imsi) == 0) {
        // All usage monitors under the same session will have the same event
        // trigger. See proto message / FeG for why. We will modify this input
        // logic later (Move event trigger out of UsageMonitorResponse), but
        // here we use a set to indicate whether a timer is already accounted
        // for.
        // Only schedule if no other revalidation timer was scheduled for
        // this IMSI
        auto revalidation_time = usage_monitor_resp.revalidation_time();
        imsis_with_revalidation.insert(imsi);
        schedule_revalidation(imsi, revalidation_time);
      }
    }
  }
}

void LocalEnforcer::update_session_credits_and_rules(
    SessionMap& session_map, const UpdateSessionResponse& response,
    SessionUpdate& session_update) {
  // These subscribers will include any subscriber that received a permanent
  // diameter error code. Additionally, it will also include CWF sessions that
  // have run out of monitoring quota.
  std::unordered_set<std::string> subscribers_to_terminate;

  update_charging_credits(
      session_map, response, subscribers_to_terminate, session_update);
  update_monitoring_credits_and_rules(
      session_map, response, subscribers_to_terminate, session_update);

  terminate_multiple_services(
      session_map, subscribers_to_terminate, session_update);
}

// terminate_subscriber (for externally triggered EndSession)
// terminates the session that is associated with the given imsi and apn
void LocalEnforcer::terminate_subscriber(
    SessionMap& session_map, const std::string& imsi, const std::string& apn,
    SessionUpdate& session_update) {
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    MLOG(MERROR) << "Could not find session for IMSI " << imsi
                 << " during termination";
    throw SessionNotFound();
  }

  for (const auto& session : it->second) {
    auto config = session->get_config();
    if (config.apn == apn) {
      SessionStateUpdateCriteria& update_criteria =
          session_update[imsi][session->get_session_id()];
      RulesToProcess rules_to_deactivate;
      // The assumption here is that
      // mutually exclusive rule names are used for different apns
      populate_rules_from_session_to_remove(imsi, session, rules_to_deactivate);
      bool deactivate_success = true;
      for (const std::string& static_rule : rules_to_deactivate.static_rules) {
        update_criteria.static_rules_to_uninstall.insert(static_rule);
      }
      for (const PolicyRule& dynamic_rule : rules_to_deactivate.dynamic_rules) {
        update_criteria.dynamic_rules_to_uninstall.insert(dynamic_rule.id());
      }
      deactivate_success = pipelined_client_->deactivate_flows_for_rules(
          imsi, rules_to_deactivate.static_rules,
          rules_to_deactivate.dynamic_rules, RequestOriginType::GX);
      if (!deactivate_success) {
        MLOG(MERROR) << "Could not deactivate flows for IMSI " << imsi
                     << " and session " << session->get_session_id()
                     << " during termination";
      }

      if (session->is_radius_cwf_session()) {
        MLOG(MDEBUG) << "Deleting UE MAC flow for subscriber " << imsi;
        SubscriberID sid;
        sid.set_id(imsi);
        auto mac_addr = config.mac_addr;
        bool delete_ue_mac_flow_success =
            pipelined_client_->delete_ue_mac_flow(sid, mac_addr);
        if (!delete_ue_mac_flow_success) {
          MLOG(MERROR) << "Failed to delete UE MAC flow for subscriber "
                       << imsi;
        }
        session->set_subscriber_quota_state(
            SubscriberQuotaUpdate_Type_TERMINATE, update_criteria);
        report_subscriber_state_to_pipelined(
            imsi, mac_addr, SubscriberQuotaUpdate_Type_TERMINATE);
      }

      session->start_termination(update_criteria);
      std::string session_id = session->get_session_id();
      // The termination should be completed when aggregated usage record no
      // longer includes the imsi. If this has not occurred after the timeout,
      // force terminate the session.
      evb_->runAfterDelay(
          [this, imsi, session_id] {
            SessionRead req  = {imsi};
            auto session_map = session_store_.read_sessions_for_deletion(req);
            auto session_update =
                SessionStore::get_default_session_update(session_map);
            SessionStateUpdateCriteria& update_criteria =
                session_update[imsi][session_id];
            MLOG(MDEBUG) << "Completing forced termination for IMSI " << imsi;
            complete_termination(
                session_map, imsi, session_id, update_criteria);

            bool end_success = session_store_.update_sessions(session_update);
            if (end_success) {
              MLOG(MDEBUG) << "Ended session " << imsi
                           << " with session_id: " << session_id;
            } else {
              MLOG(MERROR)
                  << "Failed to update SessionStore with ended session " << imsi
                  << " and session_id: " << session_id;
            }
          },
          session_force_termination_timeout_ms_);
    }
  }
}

uint64_t LocalEnforcer::get_charging_credit(
    SessionMap& session_map, const std::string& imsi,
    const CreditKey& charging_key, Bucket bucket) const {
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    return 0;
  }
  for (const auto& session : it->second) {
    uint64_t credit =
        session->get_charging_pool().get_credit(charging_key, bucket);
    if (credit > 0) {
      return credit;
    }
  }
  return 0;
}

uint64_t LocalEnforcer::get_monitor_credit(
    SessionMap& session_map, const std::string& imsi, const std::string& mkey,
    Bucket bucket) const {
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    return 0;
  }
  for (const auto& session : it->second) {
    uint64_t credit = session->get_monitor_pool().get_credit(mkey, bucket);
    if (credit > 0) {
      return credit;
    }
  }
  return 0;
}

ReAuthResult
LocalEnforcer::init_charging_reauth(SessionMap &session_map,
                                    ChargingReAuthRequest request,
                                    SessionUpdate &session_update) {
  auto it = session_map.find(request.sid());
  if (it == session_map.end()) {
    MLOG(MERROR) << "Could not find session for subscriber " << request.sid()
                 << " during reauth";
    return ReAuthResult::SESSION_NOT_FOUND;
  }
  SessionStateUpdateCriteria& update_criteria =
      session_update[request.sid()][request.session_id()];
  if (request.type() == ChargingReAuthRequest::SINGLE_SERVICE) {
    MLOG(MDEBUG) << "Initiating reauth of key " << request.charging_key()
                 << " for subscriber " << request.sid() << " for session "
                 << request.session_id();
    for (const auto& session : it->second) {
      if (session->get_session_id() == request.session_id()) {
        return session->get_charging_pool().reauth_key(
            CreditKey(request), update_criteria);
      }
    }
    MLOG(MERROR) << "Could not find session for subscriber " << request.sid()
                 << " during reauth";
    return ReAuthResult::SESSION_NOT_FOUND;
  }
  MLOG(MDEBUG) << "Initiating reauth of all keys for subscriber "
               << request.sid() << " for session" << request.session_id();
  for (const auto& session : it->second) {
    if (session->get_session_id() == request.session_id()) {
      return session->get_charging_pool().reauth_all(update_criteria);
    }
  }
  MLOG(MERROR) << "Could not find session for subscriber " << request.sid()
               << " during reauth";

  return ReAuthResult::SESSION_NOT_FOUND;
}

void LocalEnforcer::init_policy_reauth(
    SessionMap& session_map, PolicyReAuthRequest request,
    PolicyReAuthAnswer& answer_out, SessionUpdate& session_update) {
  auto it = session_map.find(request.imsi());
  if (it == session_map.end()) {
    MLOG(MERROR) << "Could not find session for subscriber " << request.imsi()
                 << " during policy reauth";
    answer_out.set_result(ReAuthResult::SESSION_NOT_FOUND);
    return;
  }

  bool deactivate_success = true;
  bool activate_success   = true;
  // For empty session_id, apply changes to all sessions of subscriber
  // Changes are applied on a best-effort basis, so failures for one session
  // won't stop changes from being applied for subsequent sessions.
  if (request.session_id() == "") {
    bool all_activated   = true;
    bool all_deactivated = true;
    for (const auto& session : it->second) {
      init_policy_reauth_for_session(
          session_map, request, session, activate_success, deactivate_success,
          session_update);
      all_activated &= activate_success;
      all_deactivated &= deactivate_success;
    }
    // Treat activate/deactivate as all-or-nothing when reporting rule failures
    mark_rule_failures(all_activated, all_deactivated, request, answer_out);
  } else {
    bool session_id_valid = false;
    for (const auto& session : it->second) {
      if (session->get_session_id() == request.session_id()) {
        session_id_valid = true;
        init_policy_reauth_for_session(
            session_map, request, session, activate_success, deactivate_success,
            session_update);
      }
    }
    if (!session_id_valid) {
      MLOG(MERROR) << "Found a matching IMSI " << request.imsi()
                   << ", but no matching session ID " << request.session_id()
                   << " during policy reauth";
      answer_out.set_result(ReAuthResult::SESSION_NOT_FOUND);
      return;
    }
    mark_rule_failures(
        activate_success, deactivate_success, request, answer_out);
  }
  answer_out.set_result(ReAuthResult::UPDATE_INITIATED);
}

void LocalEnforcer::init_policy_reauth_for_session(
    SessionMap& session_map, const PolicyReAuthRequest& request,
    const std::unique_ptr<SessionState>& session, bool& activate_success,
    bool& deactivate_success, SessionUpdate& session_update) {
  std::string imsi = request.imsi();
  SessionStateUpdateCriteria& update_criteria =
      session_update[imsi][session->get_session_id()];

  activate_success   = true;
  deactivate_success = true;
  receive_monitoring_credit_from_rar(request, session, update_criteria);

  RulesToProcess rules_to_activate;
  RulesToProcess rules_to_deactivate;

  MLOG(MDEBUG) << "Processing policy reauth for subscriber " << request.imsi();
  if (revalidation_required(request.event_triggers())) {
    schedule_revalidation(imsi, request.revalidation_time());
  }

  process_rules_to_remove(
      imsi, session, request.rules_to_remove(), rules_to_deactivate,
      update_criteria);

  process_rules_to_install(
      *session, imsi, to_vec(request.rules_to_install()),
      to_vec(request.dynamic_rules_to_install()), rules_to_activate,
      rules_to_deactivate, update_criteria);

  auto ip_addr = session->get_config().ue_ipv4;
  if (rules_to_process_is_not_empty(rules_to_deactivate)) {
    deactivate_success = pipelined_client_->deactivate_flows_for_rules(
        request.imsi(), rules_to_deactivate.static_rules,
        rules_to_deactivate.dynamic_rules, RequestOriginType::GX);
  }
  if (rules_to_process_is_not_empty(rules_to_activate)) {
    activate_success = pipelined_client_->activate_flows_for_rules(
        request.imsi(), ip_addr, rules_to_activate.static_rules,
        rules_to_activate.dynamic_rules);
  }

  if (terminate_on_wallet_exhaust() && is_wallet_exhausted(*session)) {
    RulesToProcess rules;
    populate_rules_from_session_to_remove(imsi, session, rules);
    terminate_service(
        session_map, imsi, rules.static_rules, rules.dynamic_rules,
        session_update);
    return;
  }
  if (!session->is_radius_cwf_session()) {
    create_bearer(
        activate_success, session, request, rules_to_activate.dynamic_rules);
  }
}

void LocalEnforcer::receive_monitoring_credit_from_rar(
    const PolicyReAuthRequest& request,
    const std::unique_ptr<SessionState>& session,
    SessionStateUpdateCriteria& update_criteria) {
  UsageMonitoringUpdateResponse monitoring_credit;
  monitoring_credit.set_session_id(request.session_id());
  monitoring_credit.set_sid("IMSI" + request.session_id());
  monitoring_credit.set_success(true);
  UsageMonitoringCredit* credit = monitoring_credit.mutable_credit();

  for (const auto& usage_monitoring_credit :
       request.usage_monitoring_credits()) {
    credit->CopyFrom(usage_monitoring_credit);
    session->get_monitor_pool().receive_credit(
        monitoring_credit, update_criteria);
  }
}

void LocalEnforcer::process_rules_to_remove(
    const std::string& imsi, const std::unique_ptr<SessionState>& session,
    const google::protobuf::RepeatedPtrField<std::basic_string<char>>
        rules_to_remove,
    RulesToProcess& rules_to_deactivate,
    SessionStateUpdateCriteria& update_criteria) {
  for (const auto& rule_id : rules_to_remove) {
    // Try to remove as dynamic rule first
    PolicyRule dy_rule;
    bool is_dynamic =
        session->remove_dynamic_rule(rule_id, &dy_rule, update_criteria);
    if (is_dynamic) {
      rules_to_deactivate.dynamic_rules.push_back(dy_rule);
    } else {
      if (!session->deactivate_static_rule(rule_id, update_criteria))
        MLOG(MWARNING) << "Could not find rule " << rule_id << "for IMSI "
                       << imsi << " during static rule removal";
      rules_to_deactivate.static_rules.push_back(rule_id);
    }
  }
}

void LocalEnforcer::populate_rules_from_session_to_remove(
    const std::string& imsi, const std::unique_ptr<SessionState>& session,
    RulesToProcess& rules_to_deactivate) {
  SessionState::SessionInfo info;
  session->get_session_info(info);
  for (const auto& policyrule : info.dynamic_rules) {
    rules_to_deactivate.dynamic_rules.push_back(policyrule);
  }
  for (const auto& staticrule : info.static_rules) {
    rules_to_deactivate.static_rules.push_back(staticrule);
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
    SessionState& session, const std::string& imsi,
    std::vector<StaticRuleInstall> static_rule_installs,
    std::vector<DynamicRuleInstall> dynamic_rule_installs,
    RulesToProcess& rules_to_activate, RulesToProcess& rules_to_deactivate,
    SessionStateUpdateCriteria& update_criteria) {
  std::time_t current_time = time(NULL);
  std::string ip_addr      = session.get_config().ue_ipv4;
  for (const auto& rule_install : static_rule_installs) {
    const auto& id = rule_install.rule_id();
    auto activation_time =
        TimeUtil::TimestampToSeconds(rule_install.activation_time());
    auto deactivation_time =
        TimeUtil::TimestampToSeconds(rule_install.deactivation_time());
    RuleLifetime lifetime{
        // TODO: check if we're building the time correctly
        .activation_time   = std::time_t(activation_time),
        .deactivation_time = std::time_t(deactivation_time),
    };
    if (activation_time > current_time) {
      session.schedule_static_rule(id, lifetime, update_criteria);
      schedule_static_rule_activation(imsi, ip_addr, rule_install);
    } else {
      session.activate_static_rule(id, lifetime, update_criteria);
      rules_to_activate.static_rules.push_back(id);
    }

    if (deactivation_time > current_time) {
      schedule_static_rule_deactivation(imsi, rule_install);
    } else if (deactivation_time > 0) {  // 0: never scheduled to deactivate
      if (!session.deactivate_static_rule(id, update_criteria)) {
        MLOG(MWARNING) << "Could not find rule " << id << "for IMSI " << imsi
                       << " during static rule removal";
      }
      rules_to_deactivate.static_rules.push_back(id);
    }
  }

  for (auto& rule_install : dynamic_rule_installs) {
    auto activation_time =
        TimeUtil::TimestampToSeconds(rule_install.activation_time());
    auto deactivation_time =
        TimeUtil::TimestampToSeconds(rule_install.deactivation_time());
    RuleLifetime lifetime{
        // TODO: check if we're building the time correctly
        .activation_time   = std::time_t(activation_time),
        .deactivation_time = std::time_t(deactivation_time),
    };
    if (activation_time > current_time) {
      session.schedule_dynamic_rule(
          rule_install.policy_rule(), lifetime, update_criteria);
      schedule_dynamic_rule_activation(imsi, ip_addr, rule_install);
    } else {
      session.insert_dynamic_rule(
          rule_install.policy_rule(), lifetime, update_criteria);
      rules_to_activate.dynamic_rules.push_back(rule_install.policy_rule());
    }
    if (deactivation_time > current_time) {
      schedule_dynamic_rule_deactivation(imsi, rule_install);
    } else if (deactivation_time > 0) {
      PolicyRule rule_dont_care;
      session.remove_dynamic_rule(
          rule_install.policy_rule().id(), &rule_dont_care, update_criteria);
      rules_to_deactivate.dynamic_rules.push_back(rule_install.policy_rule());
    }
  }
}

bool LocalEnforcer::revalidation_required(
    const google::protobuf::RepeatedField<int>& event_triggers) {
  auto it = std::find(
      event_triggers.begin(), event_triggers.end(), REVALIDATION_TIMEOUT);
  return it != event_triggers.end();
}

// Todo support scheduling revalidation for different sessions for a IMSI
void LocalEnforcer::schedule_revalidation(
    const std::string& imsi,
    const google::protobuf::Timestamp& revalidation_time) {
  SessionRead req = {imsi};
  auto delta = time_difference_from_now(revalidation_time);
  MLOG(MINFO) << imsi << " Scheduling revalidation in "
              << delta.count() << "ms";
  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
        std::move([=] {
          MLOG(MINFO) << imsi << " Revalidation timeout!";
          auto session_map = session_store_.read_sessions_for_reporting(req);
          SessionUpdate update =
              SessionStore::get_default_session_update(session_map);
          check_usage_for_reporting(session_map, update, true);
        }),
        delta);
  });
}

void LocalEnforcer::handle_add_ue_mac_flow_callback(
    const SubscriberID& sid,
    const std::string& ue_mac_addr,
    const std::string& msisdn,
    const std::string& apn_mac_addr,
    const std::string& apn_name,
    Status status, FlowResponse resp) {
  using namespace std::placeholders;
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

  evb_->runInEventBaseThread([=] {
    evb_->timer().scheduleTimeoutFn(
        std::move([=] {
        MLOG(MERROR) << "Could not activate ue mac flows for subscriber "
            << sid.id() << ": " << status.error_message() << ", retrying...";
        pipelined_client_->add_ue_mac_flow(
        sid, ue_mac_addr, msisdn, apn_mac_addr, apn_name,
        std::bind(
          &LocalEnforcer::handle_add_ue_mac_flow_callback,
          this, sid, ue_mac_addr, msisdn, apn_mac_addr, apn_name, _1, _2));
        }),
        retry_timeout_);
  });
}

void LocalEnforcer::create_bearer(
    const bool activate_success, const std::unique_ptr<SessionState>& session,
    const PolicyReAuthRequest& request,
    const std::vector<PolicyRule>& dynamic_rules) {
  auto config = session->get_config();
  if (!activate_success || !config.qos_info.enabled || !request.has_qos_info()) {
    MLOG(MDEBUG) << "Not creating bearer";
    return;
  }
  auto default_qci = QCI(config.qos_info.qci);
  if (request.qos_info().qci() != default_qci) {
    MLOG(MDEBUG) << "QCI sent in RAR is different from default QCI";
    spgw_client_->create_dedicated_bearer(
        request.imsi(), config.ue_ipv4,
        config.bearer_id, dynamic_rules);
  }
  return;
}

void LocalEnforcer::check_usage_for_reporting(
    SessionMap& session_map, SessionUpdate& session_update,
    const bool force_update) {
  std::vector<std::unique_ptr<ServiceAction>> actions;
  auto request =
      collect_updates(session_map, actions, session_update, force_update);
  execute_actions(session_map, actions, session_update);
  if (request.updates_size() == 0 && request.usage_monitors_size() == 0) {
    return;  // nothing to report
  }
  MLOG(MINFO) << "Sending " << request.updates_size()
              << " charging updates and " << request.usage_monitors_size()
              << " monitor updates to OCS and PCRF";

  // report to cloud
  (*reporter_)
      .report_updates(
          request,
          // For the context capture, pass by value, or create a new pointer so
          // that the value is not lost
          [this,
           request,
           session_map_ptr =
               std::make_shared<SessionMap>(std::move(session_map)),
           session_update]
               (Status status, UpdateSessionResponse response) mutable {
            if (!status.ok()) {
              MLOG(MERROR) << "Update of size " << request.updates_size()
                           << " to OCS and PCRF failed entirely: "
                           << status.error_message();
            } else {
              MLOG(MDEBUG) << "Received updated responses from OCS and PCRF";
              update_session_credits_and_rules(
                  *session_map_ptr, response, session_update);
              session_store_.update_sessions(session_update);
            }
          });
}

bool LocalEnforcer::session_with_imsi_exists(
    SessionMap& session_map, const std::string& imsi) const {
  if (session_map.find(imsi) != session_map.end()) {
    return session_map[imsi].size() > 0;
  }
  return false;
}

bool LocalEnforcer::session_with_apn_exists(
    SessionMap& session_map, const std::string& imsi,
    const std::string& apn) const {
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    return false;
  }
  for (const auto& session : it->second) {
    if (session->get_config().apn == apn) {
      return true;
    }
  }
  return false;
}

bool LocalEnforcer::is_session_active(
    SessionMap& session_map, const std::string& imsi,
    const std::string& core_session_id) const {
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    return false;
  }
  for (const auto& session : it->second) {
    if (session->get_core_session_id() == core_session_id) {
      return session->is_active();
    }
  }
  return false;
}

bool LocalEnforcer::get_core_sid_of_active_session(
    SessionMap &session_map, const std::string &imsi,
    std::string *core_session_id) const {
  auto it = session_map.find(imsi);
  if (it == session_map.end()) {
    return false;
  }
  for (const auto& session : it->second) {
    if (session->is_active()) {
      *core_session_id = session->get_core_session_id();
      return true;
    }
  }
  return false;
}

bool LocalEnforcer::get_core_sid_of_session_with_same_config(
    SessionMap &session_map, const std::string &imsi,
    const SessionConfig &config, std::string *core_session_id) const {
  auto it = session_map.find(imsi);
  if (it != session_map.end()) {
    for (const auto& session : it->second) {
      if (session->is_same_config(config)) {
        *core_session_id = session->get_core_session_id();
        return true;
      }
    }
  }
  return false;
}

void LocalEnforcer::handle_cwf_roaming(
    SessionMap& session_map, const std::string& imsi,
    const SessionConfig& config, SessionUpdate& session_update) {
  auto it = session_map.find(imsi);
  if (it != session_map.end()) {
    for (const auto& session : it->second) {
      auto& update_criteria = session_update[imsi][session->get_session_id()];
      session->set_config(config);
      update_criteria.is_config_updated = true;
      update_criteria.updated_config = session->get_config();
      // TODO Check for event triggers and send updates to the core if needed
      MLOG(MDEBUG) << "Updating IPFIX flow for subscriber " << imsi;
      SubscriberID sid;
      sid.set_id(imsi);
      std::string apn_mac_addr;
      std::string apn_name;
      if (!parse_apn(config.apn, apn_mac_addr, apn_name)) {
        MLOG(MWARNING) << "Failed mac/name parsiong for apn " << config.apn;
        apn_mac_addr = "";
        apn_name     = config.apn;
      }
      auto ue_mac_addr             = session->get_config().mac_addr;
      bool add_ue_mac_flow_success = pipelined_client_->update_ipfix_flow(
          sid, ue_mac_addr, config.msisdn, apn_mac_addr, apn_name);
      if (!add_ue_mac_flow_success) {
        MLOG(MERROR) << "Failed to update IPFIX flow for subscriber " << imsi;
      }
    }
  }
}

static void handle_command_level_result_code(
    const std::string& imsi, const uint32_t result_code,
    std::unordered_set<std::string>& subscribers_to_terminate) {
  const bool is_permanent_failure =
      DiameterCodeHandler::is_permanent_failure(result_code);
  if (is_permanent_failure) {
    MLOG(MERROR) << imsi << " Received permanent failure result code: "
                 << result_code
                 << " during update. Terminating Subscriber.";
    subscribers_to_terminate.insert(imsi);
  } else {
    // only log transient errors for now
    MLOG(MERROR) << "Received result code: " << result_code << "for IMSI "
                 << imsi << "during update";
  }
}

static void mark_rule_failures(
    const bool activate_success, const bool deactivate_success,
    const PolicyReAuthRequest& request, PolicyReAuthAnswer& answer_out) {
  auto failed_rules = *answer_out.mutable_failed_rules();
  if (!deactivate_success) {
    for (const std::string& rule_id : request.rules_to_remove()) {
      failed_rules[rule_id] = PolicyReAuthAnswer::GW_PCEF_MALFUNCTION;
    }
  }
  if (!activate_success) {
    for (const StaticRuleInstall rule : request.rules_to_install()) {
      failed_rules[rule.rule_id()] = PolicyReAuthAnswer::GW_PCEF_MALFUNCTION;
    }
    for (const DynamicRuleInstall& d_rule :
         request.dynamic_rules_to_install()) {
      failed_rules[d_rule.policy_rule().id()] =
          PolicyReAuthAnswer::GW_PCEF_MALFUNCTION;
    }
  }
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
}  // namespace magma
