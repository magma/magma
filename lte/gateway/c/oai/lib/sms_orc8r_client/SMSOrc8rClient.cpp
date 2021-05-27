/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <grpcpp/impl/codegen/async_unary_call.h>
#include <thread>  // std::thread
#include <utility>

#include "SMSOrc8rClient.h"
#include "itti_msg_to_proto_msg.h"
#include "includes/ServiceRegistrySingleton.h"
#include "lte/protos/sms_orc8r.pb.h"
#include "orc8r/protos/common.pb.h"

namespace grpc {
class Status;
}  // namespace grpc

namespace magma {

SMSOrc8rClient& SMSOrc8rClient::get_instance() {
  static SMSOrc8rClient client_instance;
  return client_instance;
}

SMSOrc8rClient::SMSOrc8rClient() {
  // Create channel
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "smsd", ServiceRegistrySingleton::LOCAL);
  // Create stub for LocalSessionManager gRPC service
  stub_ = SMSOrc8rService::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

void SMSOrc8rClient::send_uplink_unitdata(
    const itti_sgsap_uplink_unitdata_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  SMSOrc8rClient& client = get_instance();
  SMOUplinkUnitdata proto_msg =
      convert_itti_sgsap_uplink_unitdata_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncSMOUplink(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

}  // namespace magma
