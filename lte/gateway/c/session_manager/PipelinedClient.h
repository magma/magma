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
#pragma once

#include <mutex>
#include <unordered_map>
#include <queue>
#include "yaml-cpp/yaml.h"

#include <lte/protos/policydb.pb.h>
#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/pipelined.pb.h>
#include <lte/protos/subscriberdb.pb.h>

#include "GRPCReceiver.h"
#include "SessionState.h"

#define M5G_MIN_TEID (UINT32_MAX / 2)

using grpc::Status;

namespace magma {
using namespace lte;
using std::experimental::optional;

enum PipelineDRequestType {
  ACTIVATE   = 0,
  DEACTIVATE = 1,
};

struct PendingActivateRequest {
  ActivateFlowsRequest request;
  std::function<void(Status, ActivateFlowsResult)> callback_fn;
  PendingActivateRequest() {}
  PendingActivateRequest(
      ActivateFlowsRequest req,
      std::function<void(Status, ActivateFlowsResult)> cb)
      : request(req), callback_fn(cb) {}
};

struct PendingDeactivateRequest {
  DeactivateFlowsRequest request;
  std::function<void(Status, DeactivateFlowsResult)> callback_fn;
  PendingDeactivateRequest() {}
  PendingDeactivateRequest(
      DeactivateFlowsRequest req,
      std::function<void(Status, DeactivateFlowsResult)> cb)
      : request(req), callback_fn(cb) {}
};

struct PendingRequest {
  // TODO maybe use a union?
  PendingActivateRequest activate_req;
  PendingDeactivateRequest deactivate_req;
  PipelineDRequestType request_type;
  std::time_t expiry_time;

  PendingRequest() {}
  PendingRequest(
      ActivateFlowsRequest req,
      std::function<void(Status, ActivateFlowsResult)> cb, std::time_t time)
      : activate_req(req, cb), request_type(ACTIVATE), expiry_time(time) {}
  PendingRequest(
      DeactivateFlowsRequest req,
      std::function<void(Status, DeactivateFlowsResult)> cb, std::time_t time)
      : deactivate_req(req, cb), request_type(DEACTIVATE), expiry_time(time) {}
};

/**
 * PipelinedClient is the base class for managing rules and their activations.
 * The class is intended on interfacing with the data pipeline to enforce rules.
 */
class PipelinedClient {
 public:
  /**
   * Activates all rules for provided SessionInfos
   * @param infos - list of SessionInfos to setup flows for
   * @return true if the operation was successful
   */
  virtual bool setup_cwf(
      const std::vector<SessionState::SessionInfo>& infos,
      const std::vector<SubscriberQuotaUpdate>& quota_updates,
      const std::vector<std::string> ue_mac_addrs,
      const std::vector<std::string> msisdns,
      const std::vector<std::string> apn_mac_addrs,
      const std::vector<std::string> apn_names,
      const std::vector<std::uint64_t> pdp_start_times,
      const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback) = 0;

  /**
   * Activates all rules for provided SessionInfos
   * @param infos - list of SessionInfos to setup flows for
   * @return true if the operation was successful
   */
  virtual bool setup_lte(
      const std::vector<SessionState::SessionInfo>& infos,
      const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback) = 0;

  /**
   * Deactivate all flows for a subscriber's session
   * @param imsi - UE to delete all policy flows for
   * @return true if the operation was successful
   */
  virtual bool deactivate_all_flows(const std::string& imsi) = 0;

  /**
   * Deactivate all flows for the specified rules plus any drop default rule
   * added by pipelined
   * @param imsi - UE to delete flows for
   * @param rule_ids - rules to deactivate
   * @return true if the operation was successful
   */
  virtual bool deactivate_flows_for_rules_for_termination(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::vector<std::string>& rule_ids,
      const std::vector<PolicyRule>& dynamic_rules,
      const RequestOriginType_OriginType origin_type) = 0;

  /**
   * Deactivate all flows for the specified rules
   * @param imsi - UE to delete flows for
   * @param rule_ids - rules to deactivate
   * @return true if the operation was successful
   */
  virtual bool deactivate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::vector<std::string>& rule_ids,
      const std::vector<PolicyRule>& dynamic_rules,
      const RequestOriginType_OriginType origin_type) = 0;

  /**
   * Activate all rules for the specified rules, using a normal vector
   */
  virtual bool activate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::string& msisdn, const optional<AggregatedMaximumBitrate>& ambr,
      const std::vector<std::string>& static_rules,
      const std::vector<PolicyRule>& dynamic_rules,
      std::function<void(Status status, ActivateFlowsResult)> callback) = 0;

  /**
   * Send the MAC address of UE and the subscriberID
   * for pipelined to add a flow for the subscriber by matching the MAC
   */
  virtual bool add_ue_mac_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name,
      std::function<void(Status status, FlowResponse)> callback) = 0;

  /**
   * Update the IPFIX export rule in pipeliend
   */
  virtual bool update_ipfix_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name, const uint64_t& pdp_start_time) = 0;

  /**
   * Send the MAC address of UE and the subscriberID
   * for pipelined to delete a flow for the subscriber by matching the MAC
   */
  virtual bool delete_ue_mac_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr) = 0;

  /**
   * Propagate whether a subscriber has quota / no quota / or terminated
   */
  virtual bool update_subscriber_quota_state(
      const std::vector<SubscriberQuotaUpdate>& updates) = 0;

  /**
   * Activate the GY final action policy
   */
  virtual bool add_gy_final_action_flow(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::string& msisdn, const std::vector<std::string>& static_rules,
      const std::vector<PolicyRule>& dynamic_rules) = 0;

  /**
   * Set up a Session of type SetMessage to be sent to UPF
   */
  virtual bool set_upf_session(
      const SessionState::SessionInfo info,
      std::function<void(Status status, UPFSessionContextState)> callback) = 0;

  virtual uint32_t get_next_teid()    = 0;
  virtual uint32_t get_current_teid() = 0;
};

/**
 * AsyncPipelinedClient implements PipelinedClient but sends calls
 * asynchronously to pipelined.
 */
class AsyncPipelinedClient : public GRPCReceiver, public PipelinedClient {
 public:
  AsyncPipelinedClient();

  AsyncPipelinedClient(std::shared_ptr<grpc::Channel> pipelined_channel);

  /**
   * Activates all rules for provided SessionInfos
   * @param infos - list of SessionInfos to setup flows for
   * @return true if the operation was successful
   */
  bool setup_cwf(
      const std::vector<SessionState::SessionInfo>& infos,
      const std::vector<SubscriberQuotaUpdate>& quota_updates,
      const std::vector<std::string> ue_mac_addrs,
      const std::vector<std::string> msisdns,
      const std::vector<std::string> apn_mac_addrs,
      const std::vector<std::string> apn_names,
      const std::vector<std::uint64_t> pdp_start_times,
      const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback);

  /**
   * Activates all rules for provided SessionInfos
   * @param infos - list of SessionInfos to setup flows for
   * @return true if the operation was successful
   */
  bool setup_lte(
      const std::vector<SessionState::SessionInfo>& infos,
      const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback);

  /**
   * Deactivate all flows for a subscriber's session
   * @param imsi - UE to delete all policy flows for
   * @return true if the operation was successful
   */
  bool deactivate_all_flows(const std::string& imsi);

  /**
   * Deactivate all flows related to a specific charging key plus any default
   * rule installed by pipelined. Used for session termination.
   * @param imsi - UE to delete flows for
   * @param charging_key - key to deactivate
   * @return true if the operation was successful
   */
  bool deactivate_flows_for_rules_for_termination(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::vector<std::string>& rule_ids,
      const std::vector<PolicyRule>& dynamic_rules,
      const RequestOriginType_OriginType origin_type);

  /**
   * Deactivate all flows related to a specific charging key
   * @param imsi - UE to delete flows for
   * @param charging_key - key to deactivate
   * @return true if the operation was successful
   */
  bool deactivate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::vector<std::string>& rule_ids,
      const std::vector<PolicyRule>& dynamic_rules,
      const RequestOriginType_OriginType origin_type);

  /**
   * Deactivate all flows included on the request
   * @param request
   * @return true if the operation was successful
   */
  bool deactivate_flows(DeactivateFlowsRequest& request);

  /**
   * Activate all rules for the specified rules, using a normal vector
   */
  bool activate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::string& msisdn, const optional<AggregatedMaximumBitrate>& ambr,
      const std::vector<std::string>& static_rules,
      const std::vector<PolicyRule>& dynamic_rules,
      std::function<void(Status status, ActivateFlowsResult)> callback);

  /**
   * Send the MAC address of UE and the subscriberID
   * for pipelined to add a flow for the subscriber by matching the MAC
   */
  bool add_ue_mac_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name,
      std::function<void(Status status, FlowResponse)> callback);

  /**
   * Update the IPFIX export rule in pipeliend
   */
  bool update_ipfix_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name, const uint64_t& pdp_start_time);

  /**
   * Propagate whether a subscriber has quota / no quota / or terminated
   */
  bool update_subscriber_quota_state(
      const std::vector<SubscriberQuotaUpdate>& updates);

  bool delete_ue_mac_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr);

  bool add_gy_final_action_flow(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids teids,
      const std::string& msisdn, const std::vector<std::string>& static_rules,
      const std::vector<PolicyRule>& dynamic_rules);

  bool set_upf_session(
      const SessionState::SessionInfo info,
      std::function<void(Status status, UPFSessionContextState)> callback);

  void handle_add_ue_mac_callback(
      const magma::UEMacFlowRequest req, const int retries, Status status,
      FlowResponse resp);

  void set_rate_limiting_config(const YAML::Node& config);

  uint32_t get_next_teid();
  uint32_t get_current_teid();

 private:
  static const uint32_t RESPONSE_TIMEOUT_SEC = 6;  // seconds
  std::unique_ptr<Pipelined::Stub> stub_;
  uint32_t teid;

  // State for rate-limiting
  bool ENABLE_RATE_LIMITING;
  uint32_t MAX_ACTIVE_REQS;
  uint32_t MAX_QUEUE_LENGTH;
  std::queue<PendingRequest> pending_reqs_;
  std::atomic<uint32_t> ongoing_request_count_;
  std::mutex queue_lock_;

 private:
  void setup_default_controllers_rpc(
      const SetupDefaultRequest& request,
      std::function<void(Status, SetupFlowsResult)> callback);

  void setup_policy_rpc(
      const SetupPolicyRequest& request,
      std::function<void(Status, SetupFlowsResult)> callback);

  void setup_ue_mac_rpc(
      const SetupUEMacRequest& request,
      std::function<void(Status, SetupFlowsResult)> callback);

  void deactivate_flows_rpc(
      const DeactivateFlowsRequest& request,
      std::function<void(Status, DeactivateFlowsResult)> callback);

  void send_deactivate_flows_rpc(
      const DeactivateFlowsRequest& request,
      std::function<void(Status, DeactivateFlowsResult)> callback,
      const uint32_t timeout);

  void activate_flows_rpc(
      const ActivateFlowsRequest& request,
      std::function<void(Status, ActivateFlowsResult)> callback);

  void send_activate_flows_rpc(
      const ActivateFlowsRequest& request,
      std::function<void(Status, ActivateFlowsResult)> callback,
      const uint32_t timeout);

  void add_ue_mac_flow_rpc(
      const UEMacFlowRequest& request,
      std::function<void(Status, FlowResponse)> callback);

  void update_ipfix_flow_rpc(
      const UEMacFlowRequest& request,
      std::function<void(Status, FlowResponse)> callback);

  void update_subscriber_quota_state_rpc(
      const UpdateSubscriberQuotaStateRequest& request,
      std::function<void(Status, FlowResponse)> callback);

  void delete_ue_mac_flow_rpc(
      const UEMacFlowRequest& request,
      std::function<void(Status, FlowResponse)> callback);

  void set_upf_session_rpc(
      const SessionSet& request,
      std::function<void(Status, UPFSessionContextState)> callback);

  /**
   * send_rpc_with_rate_limiting handles the GRPC call with rate-limiting
   * mechanism.
   * This mode only allows for MAX_ACTIVE_REQS ongoing requests at a time.
   * All new requests when there are max number of ongoing requests will be
   * pushed onto a queue.
   * The queue at any time will only contain up to MAX_QUEUE_LENGTH pending
   * requests. Additionally, the queue removes any pending request past its
   * expiry time every time the queue is accessed.
   * @param pending_req
   * @param callback
   */
  void send_rpc_with_rate_limiting(
      const PendingRequest& pending_req,
      std::function<void(Status, ActivateFlowsResult)> callback);

  /**
   * If the queue is non-empty, pop the front of the queue and send a GRPC
   * request for the pending request.
   * It is best to call get_pending_requests_to_cancel & cancel_requests to
   * cleanup the queue before calling this function.
   * The GRPC timeout will be off-setted by the original time when it was placed
   * on the queue.
   * Ex: If request X was inserted into the queue at time t and
   * popped at time t+2, the GRPC timeout for the actual request will be:
   * RESPONSE_TIMEOUT_SEC - (t+2  - t) = RESPONSE_TIMEOUT_SEC - 2
   */
  void pop_and_send_pending_request();

  /**
   * NOTE: this function should be called with the queue locked
   * Prune a pending request if:
   * 1. The queue has only elements with expiry_time > now
   * 2. The queue has max MAX_QUEUE_LENGTH elements
   * @return a vector of pending requests that should be cancelled
   */
  std::vector<PendingRequest> get_pending_requests_to_cancel();

  /**
   * Cancel the list of PendingRequest by faking a GRPC timeout
   * (Status::CANCELLED).
   * @param requests_to_cancel: vector of pending requests that should be
   * cancelled
   */
  void cancel_requests(const std::vector<PendingRequest>& requests_to_cancel);
};

}  // namespace magma
