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

#include <lte/protos/pipelined.grpc.pb.h>
#include <lte/protos/pipelined.pb.h>
#include <lte/protos/policydb.pb.h>
#include <lte/protos/subscriberdb.pb.h>
#include <stdint.h>
#include <experimental/optional>
#include <functional>
#include <memory>
#include <mutex>
#include <string>
#include <unordered_map>
#include <vector>

#include "SessionState.h"
#include "Types.h"
#include "includes/GRPCReceiver.h"
#include "lte/protos/abort_session.pb.h"

namespace grpc {
class Channel;
class Status;
}  // namespace grpc
namespace magma {
namespace lte {
class AggregatedMaximumBitrate;
class RuleRecordTable;
class SubscriberID;
class Teids;
}  // namespace lte
}  // namespace magma

#define M5G_MIN_TEID (UINT32_MAX / 2)

using grpc::Status;

namespace magma {
using namespace lte;
using std::experimental::optional;

/**
 * PipelinedClient is the base class for managing rules and their activations.
 * The class is intended on interfacing with the data pipeline to enforce rules.
 */
class PipelinedClient {
 public:
  virtual ~PipelinedClient() = default;

  /**
   * @brief Activates all rules for provided SessionInfos
   *
   * @param infos - list of SessionInfos to setup flows
   * @param quota_updates
   * @param ue_mac_addrs
   * @param msisdns
   * @param apn_mac_addrs
   * @param apn_names
   * @param pdp_start_times
   * @param epoch
   * @param callback
   */
  virtual void setup_cwf(
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
   * @brief Set the up lte object
   *
   * @param infos
   * @param epoch
   * @param callback
   */
  virtual void setup_lte(
      const std::vector<SessionState::SessionInfo>& infos,
      const std::uint64_t& epoch,
      std::function<void(Status status, SetupFlowsResult)> callback) = 0;

  /**
   * @brief Send a DeactivateFlowsRequest for each of the teid pair.
   * The rules field will not be set to indicate that all flows should be
   * removed
   * @param imsi
   * @param ip_addr
   * @param ipv6_addr
   * @param teids
   * @param origin_type
   */
  virtual void deactivate_flows_for_rules_for_termination(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const std::vector<Teids>& teids,
      const RequestOriginType_OriginType origin_type) = 0;

  /**
   * @brief Send DeactivateFlowRequests to PipelineD
   *
   * @param imsi
   * @param ip_addr
   * @param ipv6_addr
   * @param default_teids teids corresponding to the default bearer
   * @param to_process vector of RuleToProcess. If this value is empty, the
   * function will send an empty request with the default teid value
   * @param origin_type
   */
  virtual void deactivate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids default_teids,
      const RulesToProcess to_process,
      const RequestOriginType_OriginType origin_type) = 0;

  /**
   * @brief Send ActivateFlowRequests to PipelineD
   *
   * @param imsi
   * @param ip_addr
   * @param ipv6_addr
   * @param default_teids teids corresponding to the default bearer
   * @param msisdn
   * @param ambr
   * @param to_process vector of RuleToProcess. If this value is empty, the
   * function will send an empty request with the default teid value
   * @param callback
   */
  virtual void activate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids default_teids,
      const std::string& msisdn, const optional<AggregatedMaximumBitrate>& ambr,
      const RulesToProcess to_process,
      std::function<void(Status status, ActivateFlowsResult)> callback) = 0;

  /**
   * @brief Send the MAC address of UE and the subscriberID
   * for pipelined to add a flow for the subscriber by matching the MAC
   *
   * @param sid
   * @param ue_mac_addr
   * @param msisdn
   * @param ap_mac_addr
   * @param ap_name
   * @param callback
   */
  virtual void add_ue_mac_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name,
      std::function<void(Status status, FlowResponse)> callback) = 0;

  /**
   * @brief Update the IPFIX export rule in pipeliend
   *
   * @param sid
   * @param ue_mac_addr
   * @param msisdn
   * @param ap_mac_addr
   * @param ap_name
   * @param pdp_start_time
   */
  virtual void update_ipfix_flow(const SubscriberID& sid,
                                 const std::string& ue_mac_addr,
                                 const std::string& msisdn,
                                 const std::string& ap_mac_addr,
                                 const std::string& ap_name,
                                 const uint64_t& pdp_start_time) = 0;

  /**
   * @brief Send the MAC address of UE and the subscriberID
   * for pipelined to delete a flow for the subscriber by matching the MAC
   *
   * @param sid
   * @param ue_mac_addr
   */
  virtual void delete_ue_mac_flow(const SubscriberID& sid,
                                  const std::string& ue_mac_addr) = 0;

  /**
   * @brief Propagate whether a subscriber has quota / no quota / or terminated
   *
   * @param updates
   */
  virtual void update_subscriber_quota_state(
      const std::vector<SubscriberQuotaUpdate>& updates) = 0;

  /**
   * @brief Activate the GY final action policies
   *
   * @param imsi
   * @param ip_addr
   * @param ipv6_addr
   * @param default_teids
   * @param msisdn
   * @param to_process
   */
  virtual void add_gy_final_action_flow(const std::string& imsi,
                                        const std::string& ip_addr,
                                        const std::string& ipv6_addr,
                                        const Teids default_teids,
                                        const std::string& msisdn,
                                        const RulesToProcess to_process) = 0;

  /**
   * @brief Set up a Session of type SetMessage to be sent to UPF
   *
   * @param info
   * @param callback
   */
  virtual void set_upf_session(
      const SessionState::SessionInfo info,
      const magma::RulesToProcess to_activate_process,
      const magma::RulesToProcess to_deactivate_process,
      std::function<void(Status status, UPFSessionContextState)> callback) = 0;

  virtual uint32_t get_next_teid() = 0;
  virtual uint32_t get_current_teid() = 0;

  virtual void poll_stats(
      int cookie, int cookie_mask,
      std::function<void(Status, RuleRecordTable)> callback) = 0;
};

/**
 * AsyncPipelinedClient implements PipelinedClient but sends calls
 * asynchronously to pipelined.
 * Please refer to PipelinedClient for documentation
 */
class AsyncPipelinedClient : public GRPCReceiver, public PipelinedClient {
 public:
  AsyncPipelinedClient();

  explicit AsyncPipelinedClient(
      std::shared_ptr<grpc::Channel> pipelined_channel);

  void setup_cwf(const std::vector<SessionState::SessionInfo>& infos,
                 const std::vector<SubscriberQuotaUpdate>& quota_updates,
                 const std::vector<std::string> ue_mac_addrs,
                 const std::vector<std::string> msisdns,
                 const std::vector<std::string> apn_mac_addrs,
                 const std::vector<std::string> apn_names,
                 const std::vector<std::uint64_t> pdp_start_times,
                 const std::uint64_t& epoch,
                 std::function<void(Status status, SetupFlowsResult)> callback);

  void setup_lte(const std::vector<SessionState::SessionInfo>& infos,
                 const std::uint64_t& epoch,
                 std::function<void(Status status, SetupFlowsResult)> callback);

  void deactivate_flows_for_rules_for_termination(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const std::vector<Teids>& teids,
      const RequestOriginType_OriginType origin_type);

  void deactivate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids default_teids,
      const RulesToProcess to_process,
      const RequestOriginType_OriginType origin_type);

  void deactivate_flows(DeactivateFlowsRequest& request);

  void activate_flows_for_rules(
      const std::string& imsi, const std::string& ip_addr,
      const std::string& ipv6_addr, const Teids default_teids,
      const std::string& msisdn, const optional<AggregatedMaximumBitrate>& ambr,
      const RulesToProcess to_process,
      std::function<void(Status status, ActivateFlowsResult)> callback);

  void add_ue_mac_flow(
      const SubscriberID& sid, const std::string& ue_mac_addr,
      const std::string& msisdn, const std::string& ap_mac_addr,
      const std::string& ap_name,
      std::function<void(Status status, FlowResponse)> callback);

  void update_ipfix_flow(const SubscriberID& sid,
                         const std::string& ue_mac_addr,
                         const std::string& msisdn,
                         const std::string& ap_mac_addr,
                         const std::string& ap_name,
                         const uint64_t& pdp_start_time);

  void update_subscriber_quota_state(
      const std::vector<SubscriberQuotaUpdate>& updates);

  void delete_ue_mac_flow(const SubscriberID& sid,
                          const std::string& ue_mac_addr);

  void add_gy_final_action_flow(const std::string& imsi,
                                const std::string& ip_addr,
                                const std::string& ipv6_addr,
                                const Teids default_teids,
                                const std::string& msisdn,
                                const RulesToProcess to_process);

  void set_upf_session(
      const SessionState::SessionInfo info,
      const magma::RulesToProcess to_activate_process,
      const magma::RulesToProcess to_deactivate_process,
      std::function<void(Status status, UPFSessionContextState)> callback);

  void handle_add_ue_mac_callback(const magma::UEMacFlowRequest req,
                                  const int retries, Status status,
                                  FlowResponse resp);

  /**
   * @brief Retrieves relevant records from Pipelined stats enforcements table
   * based on cookie and cookie mask
   *
   * @param cookie require matching entries to contain the cookie value
   * @param cookie_mask mask that restricts the cookie bits that must match
   */
  void poll_stats(int cookie, int cookie_mask,
                  std::function<void(Status, RuleRecordTable)> callback);

  uint32_t get_next_teid();
  uint32_t get_current_teid();

 private:
  static const uint32_t RESPONSE_TIMEOUT = 6;  // seconds
  std::unique_ptr<Pipelined::Stub> stub_;
  uint32_t teid;

 private:
  void setup_default_controllers_rpc(
      const SetupDefaultRequest& request,
      std::function<void(Status, SetupFlowsResult)> callback);

  void setup_policy_rpc(const SetupPolicyRequest& request,
                        std::function<void(Status, SetupFlowsResult)> callback);

  void setup_ue_mac_rpc(const SetupUEMacRequest& request,
                        std::function<void(Status, SetupFlowsResult)> callback);

  void deactivate_flows_rpc(
      const DeactivateFlowsRequest& request,
      std::function<void(Status, DeactivateFlowsResult)> callback);

  void activate_flows_rpc(
      const ActivateFlowsRequest& request,
      std::function<void(Status, ActivateFlowsResult)> callback);

  void add_ue_mac_flow_rpc(const UEMacFlowRequest& request,
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

  void poll_stats_rpc(const GetStatsRequest& request,
                      std::function<void(Status, RuleRecordTable)> callback);
};

}  // namespace magma
