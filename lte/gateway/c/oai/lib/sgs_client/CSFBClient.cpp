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

#include "CSFBClient.h"
#include "itti_msg_to_proto_msg.h"
#include "includes/ServiceRegistrySingleton.h"
#include "feg/protos/csfb.pb.h"
#include "orc8r/protos/common.pb.h"

namespace grpc {
class Status;
}  // namespace grpc

namespace magma {

CSFBClient& CSFBClient::get_instance() {
  static CSFBClient client_instance;
  return client_instance;
}

CSFBClient::CSFBClient() {
  // Create channel
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "csfb", ServiceRegistrySingleton::CLOUD);
  // Create stub for LocalSessionManager gRPC service
  stub_ = CSFBFedGWService::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

void CSFBClient::location_update_request(
    const itti_sgsap_location_update_req_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  LocationUpdateRequest proto_msg =
      convert_itti_sgsap_location_update_req_to_proto_msg(msg);
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  // Create a response reader for the `CreateSession` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  auto response_reader = client.stub_->AsyncLocationUpdateReq(
      local_response->get_context(), proto_msg, &client.queue_);
  // Set the reader for the local response. This executes the `CreateSession`
  // response using the response reader. When it is done, the callback stored in
  // `local_response` will be called
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::alert_ack(
    const itti_sgsap_alert_ack_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  AlertAck proto_msg = convert_itti_sgsap_alert_ack_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncAlertAc(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::alert_reject(
    const itti_sgsap_alert_reject_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client    = get_instance();
  AlertReject proto_msg = convert_itti_sgsap_alert_reject_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncAlertRej(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::tmsi_reallocation_complete(
    const itti_sgsap_tmsi_reallocation_comp_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  TMSIReallocationComplete proto_msg =
      convert_itti_sgsap_tmsi_reallocation_comp_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncTMSIReallocationComp(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::eps_detach_indication(
    const itti_sgsap_eps_detach_ind_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  EPSDetachIndication proto_msg =
      convert_itti_sgsap_eps_detach_ind_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncEPSDetachInd(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::imsi_detach_indication(
    const itti_sgsap_imsi_detach_ind_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  IMSIDetachIndication proto_msg =
      convert_itti_sgsap_imsi_detach_ind_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncIMSIDetachInd(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::paging_reject(
    const itti_sgsap_paging_reject_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client     = get_instance();
  PagingReject proto_msg = convert_itti_sgsap_paging_reject_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncPagingRej(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::service_request(
    const itti_sgsap_service_request_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  ServiceRequest proto_msg =
      convert_itti_sgsap_service_request_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncServiceReq(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::ue_activity_indication(
    const itti_sgsap_ue_activity_ind_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  UEActivityIndication proto_msg =
      convert_itti_sgsap_ue_activity_indication_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncUEActivityInd(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::ue_unreachable(
    const itti_sgsap_ue_unreachable_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client      = get_instance();
  UEUnreachable proto_msg = convert_itti_sgsap_ue_unreachable_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncUEUnreach(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

void CSFBClient::send_uplink_unitdata(
    const itti_sgsap_uplink_unitdata_t* msg,
    std::function<void(grpc::Status, Void)> callback) {
  CSFBClient& client = get_instance();
  UplinkUnitdata proto_msg =
      convert_itti_sgsap_uplink_unitdata_to_proto_msg(msg);
  auto local_response =
      new AsyncLocalResponse<Void>(std::move(callback), RESPONSE_TIMEOUT);
  auto response_reader = client.stub_->AsyncUplink(
      local_response->get_context(), proto_msg, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
}

}  // namespace magma
