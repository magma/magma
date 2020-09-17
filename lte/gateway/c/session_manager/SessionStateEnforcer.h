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
#include "RuleStore.h"
#include "PipelinedClient.h"
#include "SessionState.h"
#include "SessionStore.h"
#include "AmfServiceClient.h"


namespace magma
{
class SmSessionContextSendClient {};

//Object flow calss for 5G and composed from 4G LocalEnforcer.
class SessionStateEnforcer
{
   public:

      SessionStateEnforcer(
      std::shared_ptr<StaticRuleStore> rule_store,
      SessionStore& session_store,
      /*M5G specific parameter new objects to communicate UPF and response to AMF*/
      std::shared_ptr<PipelinedClient> pipelined_client,
      std::shared_ptr<AmfServiceClient> amf_srv_client,
      magma::mconfig::SessionD mconfig);
      //std::shared_ptr<SmSessionContextSendClient> SmSession_Context_Send_Client, //Revisit on sending msg back
      //std::shared_ptr<SpgwServiceClient> spgw_client,
      //long session_force_termination_timeout_ms,

      ~SessionStateEnforcer() {}

      void attachEventBase(folly::EventBase* evb);

      // starts the event base thread with loop
      //void start();//moved to global and get started for both 4G and 5G

      void stop();

      folly::EventBase& get_event_base();

      /*Member functions*/
      bool m5g_init_session_credit(SessionMap& session_map,
		      const std::string& imsi, const std::string& session_ctx_id,
		      const SessionConfig& cfg);
      /*Charging & rule related*/
      bool handle_session_init_rule_updates( SessionMap& session_map,
		      const std::string& imsi, SessionState& session_state);
   private:
    std::vector<std::string> static_rules;
    //std::vector<PolicyRule> dynamic_rules;

    ConvergedRuleStore  GlobalRuleList;
    std::unordered_multimap<std::string,uint32_t> pdr_map_;
    std::unordered_multimap<std::string,uint32_t> far_map_;

    std::shared_ptr<StaticRuleStore> rule_store_;
    //AsyncEventdClient& eventd_client_;
    SessionStore& session_store_;
    // Two new objects are responsible to communicate UPF and AMF 
    //std::shared_ptr<UpfClient> upf_client_;
    std::shared_ptr<PipelinedClient> pipelined_client_;
    std::shared_ptr<SmSessionContextSendClient> SmSession_Context_Send_Client_;
    std::shared_ptr<AmfServiceClient> amf_srv_client_;
    //std::unordered_map<std::string, std::vector<std::unique_ptr<SessionState>>>
    //                  session_map_;
    folly::EventBase* evb_;
    std::chrono::seconds retry_timeout_;
    magma::mconfig::SessionD mconfig_;  // Is this really reqd ?
    bool static_rule_init();
    /* To send response back to AMF 
     * Fill the response structure and call rpc of AmfServiceClient */
    void prepare_response_to_access(
		    const std::string& imsi,
		    SessionState& session_state,
		    //const std::vector<magma::PolicyRule>& flows,TODO
		    const magma::lte::M5GSMCause m5gsmcause);

    //long session_force_termination_timeout_ms_;
    //std::shared_ptr<PipelinedClient> pipelined_client_;
 }; //End of class SessionStateEnforcer

} //end namespace magma


