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

#include "PipelinedServiceClient.h"

#include <utility>
#include <cassert>
#include <grpcpp/impl/codegen/client_context.h>
#include <grpcpp/impl/codegen/status.h>
#include <netinet/in.h>
#include <cstring>
#include <iostream>
#include <memory>
#include <string>
#include <thread>

#include <grpcpp/impl/codegen/async_unary_call.h>

#include "lte/protos/pipelined.grpc.pb.h"
#include "lte/protos/pipelined.pb.h"
#include "lte/protos/mobilityd.pb.h"
#include "orc8r/protos/common.pb.h"
#include "lte/protos/subscriberdb.pb.h"
#include "includes/ServiceRegistrySingleton.h"
#include "common_defs.h"
#include "proto_converters.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc

using grpc::Channel;
using grpc::ChannelCredentials;
using grpc::ClientContext;
using grpc::InsecureChannelCredentials;
using grpc::Status;

namespace magma {
namespace lte {

PipelinedServiceClient& PipelinedServiceClient::get_instance() {
  static PipelinedServiceClient client_instance;
  return client_instance;
}

PipelinedServiceClient::PipelinedServiceClient() {
  auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
      "pipelined", ServiceRegistrySingleton::LOCAL);
  stub_ = Pipelined::NewStub(channel);
  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

//------------------- TUNNEL ADD -------------------

//                    ADD : v4
//--------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSet(
    const struct in_addr& ue_ipv4_addr, int vlan, struct in_addr& enb_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const std::string& imsi,
    uint32_t flow_precedence, const std::string& apn, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_add_update_request_ipv4(
      ue_ipv4_addr, vlan, enb_ipv4_addr, in_teid, out_teid, imsi,
      flow_precedence, apn, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    ADD : v4 with flow_dl
//--------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
    const struct in_addr& ue_ipv4_addr, int vlan, struct in_addr& enb_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const std::string& imsi,
    const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
    const std::string& apn, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_add_update_request_ipv4_flow_dl(
      ue_ipv4_addr, vlan, enb_ipv4_addr, in_teid, out_teid, imsi, flow_dl,
      flow_precedence, apn, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    ADD : v4v6
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr, int vlan,
    struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
    const std::string& imsi, uint32_t flow_precedence, const std::string& apn,
    uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_add_update_request_ipv4v6(
      ue_ipv4_addr, ue_ipv6_addr, vlan, enb_ipv4_addr, in_teid, out_teid, imsi,
      flow_precedence, apn, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));
  return RETURNok;
}

//                    ADD : v4v6 with flow dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr, int vlan,
    struct in_addr& enb_ipv4_addr, uint32_t in_teid, uint32_t out_teid,
    const std::string& imsi, const struct ip_flow_dl& flow_dl,
    uint32_t flow_precedence, const std::string& apn, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_add_update_request_ipv4v6_flow_dl(
      ue_ipv4_addr, ue_ipv6_addr, vlan, enb_ipv4_addr, in_teid, out_teid, imsi,
      flow_dl, flow_precedence, apn, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);

  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}  // namespace lte

//------------------- TUNNEL DEL -------------------

//                    DEL : v4
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSet(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_del_update_request_ipv4(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);

  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}  // namespace magma

//                    DEL : v4 with flow_dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    uint32_t in_teid, uint32_t out_teid, const struct ip_flow_dl& flow_dl,
    uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_del_update_request_ipv4_flow_dl(
      enb_ipv4_addr, ue_ipv4_addr, in_teid, out_teid, flow_dl, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    DEL : v4v6
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
    uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_del_update_request_ipv4v6(
      enb_ipv4_addr, ue_ipv4_addr, ue_ipv6_addr, in_teid, out_teid, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    DEL : v4v6 with flow_dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
    struct in_addr& enb_ipv4_addr, const struct in_addr& ue_ipv4_addr,
    struct in6_addr& ue_ipv6_addr, uint32_t in_teid, uint32_t out_teid,
    const struct ip_flow_dl& flow_dl, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_del_update_request_ipv4v6_flow_dl(
      enb_ipv4_addr, ue_ipv4_addr, ue_ipv6_addr, in_teid, out_teid, flow_dl,
      ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//------------------- DISCARDING DATA on TUNNEL -------------------
//                    DISCARD : v4
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSet(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid, uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request =
      create_discard_data_update_request_ipv4(ue_ipv4_addr, in_teid, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    DISCARD : v4 with flow_dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    const struct ip_flow_dl& flow_dl, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_discard_data_update_request_ipv4_flow_dl(
      ue_ipv4_addr, in_teid, flow_dl, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);

  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    DISCARD : v4v6
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_discard_data_update_request_ipv4v6(
      ue_ipv4_addr, ue_ipv6_addr, in_teid, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    DISCARD : v4v6 with flow_dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, const struct ip_flow_dl& flow_dl, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_discard_data_update_request_ipv4v6_flow_dl(
      ue_ipv4_addr, ue_ipv6_addr, in_teid, flow_dl, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);

  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//------------------- FORWARDING DATA on TUNNEL -------------------
//                    FORWARD : v4
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSet(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    uint32_t flow_precedence, uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_forwarding_data_update_request_ipv4(
      ue_ipv4_addr, in_teid, flow_precedence, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    FORWARD : v4 with flow_dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSetWithFlowdl(
    const struct in_addr& ue_ipv4_addr, uint32_t in_teid,
    const struct ip_flow_dl& flow_dl, uint32_t flow_precedence,
    uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_forwarding_data_update_request_ipv4_flow_dl(
      ue_ipv4_addr, in_teid, flow_dl, flow_precedence, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    FORWARD : v4v6
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSet(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, uint32_t flow_precedence, uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_forwarding_data_update_request_ipv4v6(
      ue_ipv4_addr, ue_ipv6_addr, in_teid, flow_precedence, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//                    FORWARD : v4v6 with flow_dl
//-------------------------------------------------------------------------
int PipelinedServiceClient::UpdateUEIPv4v6SessionSetWithFlowdl(
    const struct in_addr& ue_ipv4_addr, struct in6_addr& ue_ipv6_addr,
    uint32_t in_teid, const struct ip_flow_dl& flow_dl,
    uint32_t flow_precedence, uint32_t ue_state,
    std::function<void(Status, UESessionContextResponse)> callback) {
  UESessionSet request = create_forwarding_data_update_request_ipv4v6_flow_dl(
      ue_ipv4_addr, ue_ipv6_addr, in_teid, flow_dl, flow_precedence, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

//------------------- PAGING DATA on TUNNEL -------------------
int PipelinedServiceClient::UpdateUEIPv4SessionSet(
    const struct in_addr& ue_ipv4_addr, uint32_t ue_state,
    std::function<void(grpc::Status, UESessionContextResponse)> callback) {
  UESessionSet request =
      create_paging_update_request_ipv4(ue_ipv4_addr, ue_state);

  UESessionContextResponse response;

  PipelinedServiceClient& client = get_instance();

  auto local_response = new AsyncLocalResponse<UESessionContextResponse>(
      std::move(callback), RESPONSE_TIMEOUT);

  auto response_reader = client.stub_->AsyncUpdateUEState(
      local_response->get_context(), request, &client.queue_);
  local_response->set_response_reader(std::move(response_reader));

  return RETURNok;
}

}  // namespace lte
}  // namespace magma
