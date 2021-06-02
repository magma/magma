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

#include <google/protobuf/util/time_util.h>

#include <memory>
#include <string>
#include <utility>
#include <vector>

#include "EnumToString.h"
#include "GrpcMagmaUtils.h"
#include "magma_logging.h"
#include "PipelinedClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include "Types.h"

using grpc::Status;

namespace {  // anonymous
using std::experimental::optional;
// Preparation of Set Session request to UPF
magma::SessionSet create_session_set_req(
    magma::SessionState::SessionInfo info) {
  magma::SessionSet req;
  req.set_subscriber_id(info.subscriber_id);
  req.set_session_version(info.ver_no);
  req.set_local_f_teid(info.local_f_teid);
  req.mutable_node_id()->set_node_id(info.nodeId.node_id);
  req.mutable_node_id()->set_node_id_type(magma::NodeID::IPv4);
  req.mutable_state()->set_state(info.state);

  for (const auto& final_req : info.Pdr_rules_) {
    req.mutable_set_gr_pdr()->Add()->CopyFrom(final_req);
  }
  return req;
}

magma::DeactivateFlowsRequest make_deactivate_req(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const magma::Teids teids,
    const magma::RequestOriginType_OriginType origin_type,
    const bool remove_default_drop_rules) {
  magma::DeactivateFlowsRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(ip_addr);
  req.set_ipv6_addr(ipv6_addr);
  req.set_downlink_tunnel(teids.enb_teid());
  req.set_uplink_tunnel(teids.agw_teid());
  req.set_remove_default_drop_flows(remove_default_drop_rules);
  req.mutable_request_origin()->set_type(origin_type);
  return req;
}

/**
 * @brief Create a map of Teids -> DeactivateFlowsRequest
 * If to_process is empty, create one default_teids -> DeactivateFlowsRequest
 * with no rules. For each to_process item, create item.teid ->
 * DeactivateFlowsRequest with item.rules
 * @param imsi
 * @param ip_addr
 * @param ipv6_addr
 * @param default_teids this value is only used if to_process is empty
 * @param to_process
 * @param origin_type
 * @param remove_default_drop_rules
 * @return magma::DeactivateReqByTeids
 */
magma::DeactivateReqByTeids make_deactivate_req_by_teid(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const magma::Teids default_teids,
    const magma::RulesToProcess& to_process,
    const magma::RequestOriginType_OriginType origin_type,
    const bool remove_default_drop_rules) {
  magma::DeactivateReqByTeids deactivate_req_by_teids;
  if (to_process.empty()) {
    // Send an empty request with the default teid
    deactivate_req_by_teids[default_teids] = make_deactivate_req(
        imsi, ip_addr, ipv6_addr, default_teids, origin_type,
        remove_default_drop_rules);
    return deactivate_req_by_teids;
  }

  for (const magma::RuleToProcess& val : to_process) {
    const magma::Teids& dedicated_teids = val.teids;
    if (deactivate_req_by_teids.find(dedicated_teids) ==
        deactivate_req_by_teids.end()) {
      deactivate_req_by_teids[dedicated_teids] = make_deactivate_req(
          imsi, ip_addr, ipv6_addr, dedicated_teids, origin_type,
          remove_default_drop_rules);
    }
    auto versioned_policy =
        deactivate_req_by_teids[dedicated_teids].mutable_policies()->Add();
    versioned_policy->set_version(val.version);
    versioned_policy->set_rule_id(val.rule.id());
  }
  return deactivate_req_by_teids;
}

magma::ActivateFlowsRequest make_activate_req(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const magma::Teids teids,
    const std::string& msisdn,
    const optional<magma::AggregatedMaximumBitrate>& ambr,
    const magma::RequestOriginType_OriginType origin_type) {
  magma::ActivateFlowsRequest req;
  req.mutable_sid()->set_id(imsi);
  req.set_ip_addr(ip_addr);
  req.set_ipv6_addr(ipv6_addr);
  req.set_downlink_tunnel(teids.enb_teid());
  req.set_uplink_tunnel(teids.agw_teid());
  req.set_msisdn(msisdn);
  req.mutable_request_origin()->set_type(origin_type);
  if (ambr) {
    req.mutable_apn_ambr()->CopyFrom(*ambr);
  }
  return req;
}

/**
 * @brief Create a map of Teids -> ActivateFlowsRequest
 * If to_process is empty, create one default_teids -> ActivateFlowsRequest with
 * no rules. For each to_process item, create item.teid -> ActivateFlowsRequest
 * with item.rules
 * @param imsi
 * @param ip_addr
 * @param ipv6_addr
 * @param default_teids
 * @param msisdn
 * @param ambr
 * @param to_process
 * @param origin_type
 * @return magma::ActivateReqByTeids
 */
magma::ActivateReqByTeids make_activate_req_by_teid(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const magma::Teids default_teids,
    const std::string& msisdn,
    const optional<magma::AggregatedMaximumBitrate>& ambr,
    const magma::RulesToProcess& to_process,
    const magma::RequestOriginType_OriginType origin_type) {
  magma::ActivateReqByTeids activate_req_by_teids;
  if (to_process.empty()) {
    // Send an empty request with the default teid
    activate_req_by_teids[default_teids] = make_activate_req(
        imsi, ip_addr, ipv6_addr, default_teids, msisdn, ambr, origin_type);
    return activate_req_by_teids;
  }

  for (const magma::RuleToProcess& val : to_process) {
    const magma::Teids& dedicated_teids = val.teids;
    if (activate_req_by_teids.find(dedicated_teids) ==
        activate_req_by_teids.end()) {
      activate_req_by_teids[dedicated_teids] = make_activate_req(
          imsi, ip_addr, ipv6_addr, dedicated_teids, msisdn, ambr, origin_type);
    }
    auto versioned_policy =
        activate_req_by_teids[dedicated_teids].mutable_policies()->Add();
    versioned_policy->set_version(val.version);
    versioned_policy->mutable_rule()->CopyFrom(val.rule);
  }
  return activate_req_by_teids;
}

magma::UEMacFlowRequest create_add_ue_mac_flow_req(
    const magma::SubscriberID& sid, const std::string& ue_mac_addr,
    const std::string& msisdn, const std::string& ap_mac_addr,
    const std::string& ap_name, const std::uint64_t& pdp_start_time) {
  magma::UEMacFlowRequest req;
  req.mutable_sid()->CopyFrom(sid);
  req.set_mac_addr(ue_mac_addr);
  req.set_msisdn(msisdn);
  req.set_ap_mac_addr(ap_mac_addr);
  req.set_ap_name(ap_name);
  req.set_pdp_start_time(pdp_start_time);
  return req;
}

magma::UEMacFlowRequest create_delete_ue_mac_flow_req(
    const magma::SubscriberID& sid, const std::string& ue_mac_addr) {
  magma::UEMacFlowRequest req;
  req.mutable_sid()->CopyFrom(sid);
  req.set_mac_addr(ue_mac_addr);
  return req;
}

magma::SetupDefaultRequest create_setup_default_req(
    const std::uint64_t& epoch) {
  magma::SetupDefaultRequest req;
  req.set_epoch(epoch);
  return req;
}

magma::SetupPolicyRequest create_setup_policy_req(
    const std::vector<magma::SessionState::SessionInfo>& infos,
    const std::uint64_t& epoch) {
  magma::SetupPolicyRequest req;
  req.set_epoch(epoch);
  auto mut_requests = req.mutable_requests();

  for (auto it = infos.begin(); it != infos.end(); it++) {
    magma::ActivateReqByTeids gx_activate_reqs = make_activate_req_by_teid(
        it->imsi, it->ip_addr, it->ipv6_addr, it->teids, it->msisdn, it->ambr,
        it->gx_rules, magma::RequestOriginType::GX);
    for (auto& activate_pair : gx_activate_reqs) {
      mut_requests->Add()->CopyFrom(activate_pair.second);
    }

    magma::ActivateReqByTeids gy_activate_reqs = make_activate_req_by_teid(
        it->imsi, it->ip_addr, it->ipv6_addr, it->teids, it->msisdn, {},
        it->gy_dynamic_rules, magma::RequestOriginType::GY);
    for (auto& activate_pair : gy_activate_reqs) {
      mut_requests->Add()->CopyFrom(activate_pair.second);
    }
  }
  return req;
}

magma::SetupUEMacRequest create_setup_ue_mac_req(
    const std::vector<magma::SessionState::SessionInfo>& infos,
    const std::vector<std::string> ue_mac_addrs,
    const std::vector<std::string> msisdns,
    const std::vector<std::string> apn_mac_addrs,
    const std::vector<std::string> apn_names,
    const std::vector<std::uint64_t> pdp_start_times,
    const std::uint64_t& epoch) {
  magma::SetupUEMacRequest req;
  std::vector<magma::UEMacFlowRequest> activation_reqs;

  for (unsigned i = 0; i < infos.size(); i++) {
    magma::SubscriberID sid;
    sid.set_id(infos[i].imsi);
    auto activate_req = create_add_ue_mac_flow_req(
        sid, ue_mac_addrs[i], msisdns[i], apn_mac_addrs[i], apn_names[i],
        pdp_start_times[i]);
    activation_reqs.push_back(activate_req);
  }
  auto mut_requests = req.mutable_requests();
  for (const auto& act_req : activation_reqs) {
    mut_requests->Add()->CopyFrom(act_req);
  }
  req.set_epoch(epoch);
  return req;
}

magma::UpdateSubscriberQuotaStateRequest create_subscriber_quota_state_req(
    const std::vector<magma::SubscriberQuotaUpdate>& updates) {
  magma::UpdateSubscriberQuotaStateRequest req;
  auto p_updates = req.mutable_updates();
  for (const auto& update : updates) {
    p_updates->Add()->CopyFrom(update);
  }
  return req;
}

}  // namespace

namespace magma {

AsyncPipelinedClient::AsyncPipelinedClient(
    std::shared_ptr<grpc::Channel> channel)
    : stub_(Pipelined::NewStub(channel)) {
  teid = M5G_MIN_TEID;
}

AsyncPipelinedClient::AsyncPipelinedClient()
    : AsyncPipelinedClient(ServiceRegistrySingleton::Instance()->GetGrpcChannel(
          "pipelined", ServiceRegistrySingleton::LOCAL)) {}

void AsyncPipelinedClient::setup_cwf(
    const std::vector<SessionState::SessionInfo>& infos,
    const std::vector<SubscriberQuotaUpdate>& quota_updates,
    const std::vector<std::string> ue_mac_addrs,
    const std::vector<std::string> msisdns,
    const std::vector<std::string> apn_mac_addrs,
    const std::vector<std::string> apn_names,
    const std::vector<std::uint64_t> pdp_start_times,
    const std::uint64_t& epoch,
    std::function<void(Status status, SetupFlowsResult)> callback) {
  SetupDefaultRequest setup_default_req = create_setup_default_req(epoch);
  setup_default_controllers_rpc(setup_default_req, callback);
  SetupPolicyRequest setup_policy_req = create_setup_policy_req(infos, epoch);
  setup_policy_rpc(setup_policy_req, callback);

  SetupUEMacRequest setup_ue_mac_req = create_setup_ue_mac_req(
      infos, ue_mac_addrs, msisdns, apn_mac_addrs, apn_names, pdp_start_times,
      epoch);
  setup_ue_mac_rpc(setup_ue_mac_req, callback);

  update_subscriber_quota_state(quota_updates);
}

void AsyncPipelinedClient::setup_lte(
    const std::vector<SessionState::SessionInfo>& infos,
    const std::uint64_t& epoch,
    std::function<void(Status status, SetupFlowsResult)> callback) {
  SetupDefaultRequest setup_default_req = create_setup_default_req(epoch);
  setup_default_controllers_rpc(setup_default_req, callback);
  SetupPolicyRequest setup_policy_req = create_setup_policy_req(infos, epoch);
  setup_policy_rpc(setup_policy_req, callback);
}

// Method to Setup UPF Session
void AsyncPipelinedClient::set_upf_session(
    const SessionState::SessionInfo info,
    std::function<void(Status status, UPFSessionContextState)> callback) {
  SessionSet setup_session_req = create_session_set_req(info);
  set_upf_session_rpc(setup_session_req, callback);
}

void AsyncPipelinedClient::deactivate_flows_for_rules_for_termination(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const std::vector<Teids>& teids,
    const RequestOriginType_OriginType origin_type) {
  for (const Teids& t : teids) {
    MLOG(MDEBUG) << "Deactivating all rules and default drop flows "
                 << "for " << imsi << ", ipv4: " << ip_addr
                 << ", ipv6: " << ipv6_addr << ", agw teid: " << t.agw_teid()
                 << ", enb teid: " << t.enb_teid() << ", origin_type: "
                 << request_origin_type_to_str(origin_type);
    DeactivateFlowsRequest req =
        make_deactivate_req(imsi, ip_addr, ipv6_addr, t, origin_type, true);
    deactivate_flows(req);
  }
}

void AsyncPipelinedClient::deactivate_flows_for_rules(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const Teids teids,
    const RulesToProcess to_process,
    const RequestOriginType_OriginType origin_type) {
  MLOG(MDEBUG) << "Deactivating " << to_process.size()
               << " rules and for subscriber " << imsi << " IP " << ip_addr
               << " " << ipv6_addr;

  DeactivateReqByTeids reqs = make_deactivate_req_by_teid(
      imsi, ip_addr, ipv6_addr, teids, to_process, origin_type, false);
  for (auto& req_pair : reqs) {
    deactivate_flows(req_pair.second);
  }
}

void AsyncPipelinedClient::deactivate_flows(DeactivateFlowsRequest& request) {
  auto imsi = request.sid().id();
  deactivate_flows_rpc(
      request, [imsi](Status status, DeactivateFlowsResult resp) {
        if (!status.ok()) {
          MLOG(MERROR) << "Could not deactivate flows for subscriber " << imsi
                       << ": " << status.error_message();
        }
      });
}

void AsyncPipelinedClient::activate_flows_for_rules(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const Teids teids, const std::string& msisdn,
    const optional<AggregatedMaximumBitrate>& ambr,
    const RulesToProcess to_process,
    std::function<void(Status status, ActivateFlowsResult)> callback) {
  MLOG(MDEBUG) << "Activating " << to_process.size() << " rules for " << imsi
               << " msisdn " << msisdn << " and ip " << ip_addr << " "
               << ipv6_addr;
  ActivateReqByTeids reqs = make_activate_req_by_teid(
      imsi, ip_addr, ipv6_addr, teids, msisdn, ambr, to_process,
      RequestOriginType::GX);
  for (auto& activate_pair : reqs) {
    activate_flows_rpc(activate_pair.second, callback);
  }
}

void AsyncPipelinedClient::add_ue_mac_flow(
    const SubscriberID& sid, const std::string& ue_mac_addr,
    const std::string& msisdn, const std::string& ap_mac_addr,
    const std::string& ap_name,
    std::function<void(Status status, FlowResponse)> callback) {
  auto req = create_add_ue_mac_flow_req(
      sid, ue_mac_addr, msisdn, ap_mac_addr, ap_name, 0);
  add_ue_mac_flow_rpc(req, callback);
}

void AsyncPipelinedClient::update_ipfix_flow(
    const SubscriberID& sid, const std::string& ue_mac_addr,
    const std::string& msisdn, const std::string& ap_mac_addr,
    const std::string& ap_name, const uint64_t& pdp_start_time) {
  auto req = create_add_ue_mac_flow_req(
      sid, ue_mac_addr, msisdn, ap_mac_addr, ap_name, pdp_start_time);
  update_ipfix_flow_rpc(req, [ue_mac_addr](Status status, FlowResponse resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not update ipfix flow for subscriber with MAC"
                   << ue_mac_addr << ": " << status.error_message();
    }
  });
}

void AsyncPipelinedClient::delete_ue_mac_flow(
    const SubscriberID& sid, const std::string& ue_mac_addr) {
  auto req = create_delete_ue_mac_flow_req(sid, ue_mac_addr);
  delete_ue_mac_flow_rpc(req, [ue_mac_addr](Status status, FlowResponse resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not delete flow for subscriber with UE MAC"
                   << ue_mac_addr << ": " << status.error_message();
    }
  });
}

void AsyncPipelinedClient::update_subscriber_quota_state(
    const std::vector<SubscriberQuotaUpdate>& updates) {
  auto req = create_subscriber_quota_state_req(updates);
  update_subscriber_quota_state_rpc(req, [](Status status, FlowResponse resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could send quota update " << status.error_message();
    }
  });
}

void AsyncPipelinedClient::add_gy_final_action_flow(
    const std::string& imsi, const std::string& ip_addr,
    const std::string& ipv6_addr, const Teids teids, const std::string& msisdn,
    const RulesToProcess to_process) {
  MLOG(MDEBUG) << "Activating GY final action for subscriber " << imsi;
  ActivateReqByTeids reqs = make_activate_req_by_teid(
      imsi, ip_addr, ipv6_addr, teids, msisdn, {}, to_process,
      RequestOriginType::GY);
  auto cb = [imsi](Status status, ActivateFlowsResult resp) {
    if (!status.ok()) {
      MLOG(MERROR) << "Could not activate GY flows through pipelined for UE "
                   << imsi << ": " << status.error_message();
    }
  };

  for (auto& activate_pair : reqs) {
    activate_flows_rpc(activate_pair.second, cb);
  }
}

// RPC definition to Send Set Session request to UPF
void AsyncPipelinedClient::set_upf_session_rpc(
    const SessionSet& request,
    std::function<void(Status, UPFSessionContextState)> callback) {
  auto local_resp = new AsyncLocalResponse<UPFSessionContextState>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(
      stub_->AsyncSetSMFSessions(local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::setup_default_controllers_rpc(
    const SetupDefaultRequest& request,
    std::function<void(Status, SetupFlowsResult)> callback) {
  auto local_resp = new AsyncLocalResponse<SetupFlowsResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(stub_->AsyncSetupDefaultControllers(
      local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::setup_policy_rpc(
    const SetupPolicyRequest& request,
    std::function<void(Status, SetupFlowsResult)> callback) {
  auto local_resp = new AsyncLocalResponse<SetupFlowsResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(stub_->AsyncSetupPolicyFlows(
      local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::setup_ue_mac_rpc(
    const SetupUEMacRequest& request,
    std::function<void(Status, SetupFlowsResult)> callback) {
  auto local_resp = new AsyncLocalResponse<SetupFlowsResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(stub_->AsyncSetupUEMacFlows(
      local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::deactivate_flows_rpc(
    const DeactivateFlowsRequest& request,
    std::function<void(Status, DeactivateFlowsResult)> callback) {
  auto local_resp = new AsyncLocalResponse<DeactivateFlowsResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(stub_->AsyncDeactivateFlows(
      local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::activate_flows_rpc(
    const ActivateFlowsRequest& request,
    std::function<void(Status, ActivateFlowsResult)> callback) {
  auto local_resp = new AsyncLocalResponse<ActivateFlowsResult>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(
      stub_->AsyncActivateFlows(local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::add_ue_mac_flow_rpc(
    const UEMacFlowRequest& request,
    std::function<void(Status, FlowResponse)> callback) {
  auto local_resp = new AsyncLocalResponse<FlowResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(
      stub_->AsyncAddUEMacFlow(local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::update_ipfix_flow_rpc(
    const UEMacFlowRequest& request,
    std::function<void(Status, FlowResponse)> callback) {
  auto local_resp = new AsyncLocalResponse<FlowResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(stub_->AsyncUpdateIPFIXFlow(
      local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::delete_ue_mac_flow_rpc(
    const UEMacFlowRequest& request,
    std::function<void(Status, FlowResponse)> callback) {
  auto local_resp = new AsyncLocalResponse<FlowResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(std::move(stub_->AsyncDeleteUEMacFlow(
      local_resp->get_context(), request, &queue_)));
}

void AsyncPipelinedClient::update_subscriber_quota_state_rpc(
    const UpdateSubscriberQuotaStateRequest& request,
    std::function<void(Status, FlowResponse)> callback) {
  auto local_resp = new AsyncLocalResponse<FlowResponse>(
      std::move(callback), RESPONSE_TIMEOUT);
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  local_resp->set_response_reader(
      std::move(stub_->AsyncUpdateSubscriberQuotaState(
          local_resp->get_context(), request, &queue_)));
}

uint32_t AsyncPipelinedClient::get_next_teid() {
  /* For now TEID we use current no, increment for next, later we plan to
     maintain  release/alloc table for reu sing */
  uint32_t allocated_teid = teid++;
  return allocated_teid;
}

uint32_t AsyncPipelinedClient::get_current_teid() {
  /* For now TEID we use current no, increment for next, later we plan to
     maintain  release/alloc table for reu sing */
  return teid;
}

}  // namespace magma
