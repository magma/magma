/*
Copyright 2020 The Magma Authors.
This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*****************************************************************************
  Source      	SessionStateEnforcer.h
  Version     	0.1
  Date       	2020/08/08
  Product     	SessionD
  Subsystem   	5G managing & maintaining state & store of session of SessionD
                Fanout message to Access and UPF through respective client obj
  Author/Editor Sanjay Kumar Ojha
  Description 	Objects run in main thread context invoked by folly event
*****************************************************************************/

#pragma once

#include <unordered_map>
#include <map>
#include <unordered_set>
#include <vector>

#include <folly/io/async/EventBaseManager.h>
#include <lte/protos/mconfig/mconfigs.pb.h>
#include <lte/protos/policydb.pb.h>
#include "PipelinedClient.h"
#include "RuleStore.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "AmfServiceClient.h"

namespace magma {

class SessionStateEnforcer {
 public:
  SessionStateEnforcer(
      std::shared_ptr<StaticRuleStore> rule_store, SessionStore& session_store,
      /*M5G specific parameter new objects to communicate UPF and response to
         AMF*/
      std::shared_ptr<PipelinedClient> pipelined_client,
      std::shared_ptr<AmfServiceClient> amf_srv_client,
      magma::mconfig::SessionD mconfig,
      long session_force_termination_timeout_ms);

  ~SessionStateEnforcer() {}

  void attachEventBase(folly::EventBase* evb);

  void stop();

  folly::EventBase& get_event_base();

  /*Member functions*/
  bool m5g_init_session_credit(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, const SessionConfig& cfg);
  /*Charging & rule related*/
  void handle_session_init_rule_updates(
      SessionMap& session_map, const std::string& imsi,
      SessionState& session_state);

  /*Release request handle*/
  bool m5g_release_session(
      SessionMap& session_map, const std::string& imsi, const std::string& dnn,
      SessionUpdate& session_update);

  /*Handle and update respective session upon receiving message from UPF*/
  void m5g_update_session_state_to_amf(
      const std::string& imsi, uint32_t teid, uint32_t version);

  std::vector<std::string> static_rules;

  /* Get N3 ip  address of UPF */
  std::string get_upf_n3_addr() const;

  /* Get N3 ip  address of UPF */
  std::string get_upf_node_id() const;

  /* Initialize the upf node id and n3 address
   */
  bool set_upf_node(const std::string& node_id, const std::string& n3_addr);

 private:
  ConvergedRuleStore GlobalRuleList;
  std::unordered_multimap<std::string, uint32_t> pdr_map_;
  std::unordered_multimap<std::string, uint32_t> far_map_;

  std::shared_ptr<StaticRuleStore> rule_store_;
  SessionStore& session_store_;
  std::shared_ptr<PipelinedClient> pipelined_client_;
  std::shared_ptr<AmfServiceClient> amf_srv_client_;
  folly::EventBase* evb_;
  std::chrono::seconds retry_timeout_;
  magma::mconfig::SessionD mconfig_;
  std::string upf_node_id_;
  uint32_t teid_counter_;
  std::string upf_node_ip_addr_;

  // Timer used to forcefully terminate session context on time out
  long session_force_termination_timeout_ms_;
  bool static_rule_init();
  /* To send response back to AMF
   * Fill the response structure and call rpc of AmfServiceClient
   */
  void prepare_response_to_access(
      const std::string& imsi, SessionState& session_state,
      const magma::lte::M5GSMCause m5gsmcause);

  void handle_state_update_to_amf(
      const std::string& imsi, SessionState& session_state,
      const magma::lte::M5GSMCause m5gsmcause);

  /*Start processing to terminate respective session requested from AMF*/
  void m5g_start_session_termination(
      const std::string& imsi, const std::unique_ptr<SessionState>& session,
      const std::string& dnn, SessionStateUpdateCriteria& uc);

  /* Function to handle termination if UPF doesn't send required report
   * As per current implementation, upf report is not in place and
   * termination on time out will be executed forcefully
   */
  void m5g_handle_termination_on_timeout(
      const std::string& imsi, const std::string& session_id);

  /*Function will clean up all resources related to requested session*/
  void m5g_complete_termination(
      SessionMap& session_map, const std::string& imsi,
      const std::string& session_id, SessionUpdate& session_update);

  /*Function is to remove the associated rules to respective sessions*/
  void m5g_remove_rules_and_update_upf(
      const std::string& imsi, const std::unique_ptr<SessionState>& session,
      const std::string& dnn, SessionStateUpdateCriteria& uc);

};  // End of class SessionStateEnforcer

}  // end namespace magma
