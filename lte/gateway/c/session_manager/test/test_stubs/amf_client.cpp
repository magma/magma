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

#include <grpcpp/channel.h>
#include <grpcpp/impl/codegen/async_unary_call.h>
#include <thread>
#include <iostream>
#include <string>
#include <utility>

//#include "orc8r/protos/mconfig/mconfigs.pb.h"
#include "amf_client.h"
#include "ServiceRegistrySingleton.h"
#include "MConfigLoader.h"
#include "lte/protos/session_manager.pb.h"

#if 0
#include "MagmaService.h"
#include "magma_logging.h"
#include "ServiceRegistrySingleton.h"
#include "SetMessageManagerHandler.h"
#include "GrpcMagmaUtils.h"
#endif

//Convered core development new files
#define SESSIOND_SERVICE "sessiond"
#define SESSION_PROXY_SERVICE "session_proxy"
#define POLICYDB_SERVICE "policydb"
#define SESSIOND_VERSION "1.0"
#define MIN_USAGE_REPORTING_THRESHOLD 0.4
#define MAX_USAGE_REPORTING_THRESHOLD 1.1
#define DEFAULT_USAGE_REPORTING_THRESHOLD 0.8
#define DEFAULT_QUOTA_EXHAUSTION_TERMINATION_MS 30000  // 30sec

#define MAGMAD_SERVICE "magmad"

using grpc::Status;

namespace magma {

AMFClient& AMFClient::get_instance() {
  static AMFClient client_instance;
  return client_instance;
}


AMFClient::AMFClient() {
  std::shared_ptr<Channel> channel;
  channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "sessiond", ServiceRegistrySingleton::LOCAL);

stub_ = AmfPduSessionSmContext::NewStub(channel);
std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

void AMFClient::amf_create_session_final(
    const SetSMSessionContext& request,
    std::function<void(Status, SmContextVoid)> callback) {
  AMFClient& client = get_instance();
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
    auto local_response = new AsyncLocalResponse<SmContextVoid>(
    std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `CreateSession` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
#if 0
    PrintGrpcMessage(
      static_cast<const google::protobuf::Message&>(request));
#endif
    auto response_reader = client.stub_->AsyncSetAmfSessionContext(
    local_response->get_context(), request, &client.queue_);
  // Set the reader for the local response. This executes the `CreateSession`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
    local_response->set_response_reader(std::move(response_reader));
}
}
/*forward declaration of void NULL function for call back*/
//void call_back_void(magma::SmContextVoid response) 
void call_back_void(grpc::Status, magma::SmContextVoid response) 
{
  //do nothinf but to only passing call back
  //cout <<" Only for testing call back" << endl;
}
 extern "C" {}

std::function<void(grpc::Status, magma::SmContextVoid)>callback  = call_back_void;

void amf_create_session() {
    magma::SetSMSessionContext sreq;
    auto *req =  sreq.mutable_m5g_rat_specific_context()->mutable_m5gsm_session_context();
    req->set_pdu_session_id({0x5});
    req->set_rquest_type(magma::RequestType::INITIAL_REQUEST);
    req->mutable_pdu_address()->set_redirect_address_type(magma::RedirectServer::IPV4);
    req->mutable_pdu_address()->set_redirect_server_address("10.20.30.40");
    req->set_priority_access(magma::priorityaccess::High);
    req->set_access_type(magma::AccessType::M_3GPP_ACCESS_3GPP);
    req->set_fsm_session_state(magma::SMSessionFSMState::CREATING_0);
    req->set_imei("123456789012345");
    req->set_gpsi("9876543210");
    req->set_pcf_id("1357924680123456");
    grpc::Status status;
    //magma::SmContextVoid response;

    magma::AMFClient::amf_create_session_final(sreq, callback);
    //magma::AMFClient::amf_create_session_final(sreq,grpc::Status status, magma::SmContextVoid response)
}
