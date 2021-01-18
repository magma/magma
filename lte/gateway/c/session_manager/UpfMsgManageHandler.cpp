/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
/*****************************************************************************
  Source      	UpfMsgManagerHandler.cpp
  Version     	1.0
  Date       	2020/01/17
  Product     	SessionD
  Subsystem   	UPF converged Landing object in SessionD
  Author/Editor Venu Kumar Gurrapu
  Description 	Acts as 5G Landing object in SessionD & start 5G related flow
*****************************************************************************/
#include <chrono>
#include <thread>

#include <google/protobuf/util/time_util.h>
#include "lte/protos/session_manager.pb.h"
#include "UpfMsgManageHandler.h"
#include "magma_logging.h"
#include "GrpcMagmaUtils.h"

using grpc::Status;

namespace magma
{
/**
 * SetInterfaceForUserPlaneHandler processes gRPC requests for the sessionD
 * This composites the all the request that comes from UPF 
 */

UpfMsgManageHandler::UpfMsgManageHandler(
    std::shared_ptr<SessionStateEnforcer> enforcer, SessionStore& session_store)
    : session_store_(session_store),conv_enforcer_(enforcer) {}
  
/**
   * Node level GRPC message received from UPF 
   * during startup
   */

void UpfMsgManageHandler::SetUPFNodeState(
      ServerContext* context, const UPFNodeState* node_request,
      std::function<void(Status, SmContextVoid)> response_callback) {

  auto&  request  =*node_request;
  //Print the message from UPF
  PrintGrpcMessage(static_cast<const google::protobuf::Message&>(request));
  MLOG(MDEBUG) << "Node UPF details :";
  std::string ipv4_addr;
  std::string upf_id = node_request->upf_id();
  UPFAssociationState Assostate = node_request->associaton_state();
  //auto state_ver = Assostate.state_version();
  //auto  assoc_state = Assostate.assoc_state();
  auto recovery_time = Assostate.recovery_time_stamp();
  auto feature_set = Assostate.feature_set();
  //   auto grace_time = Assostate.graceful_relase_period();
   // For now get User Plan IPv4 resource at index '0' only
  ipv4_addr = Assostate.ip_resource_schema(0).ipv4_address();
  conv_enforcer_->set_upf_node_addr(ipv4_addr);
#if 0
  conv_enforcer_->get_event_base().runInEventBaseThread (
      [ipv4_addr,conv_enforcer_]() {
       //Set the IPv4 address
       conv_enforcer_->set_upf_node_addr(ipv4_addr);
  });
#endif
  //TODO  
  //send the same UPF Node association response back to UPF
}

  /**
   * Periodic messages about UPF session config 
   * 
   */
void UpfMsgManageHandler::SetUPFSessionsConfig(
      ServerContext* context, const UPFSessionConfigState* sess_config,
      std::function<void(Status, SmContextVoid)> response_callback) {
}
}//end namespace magma


