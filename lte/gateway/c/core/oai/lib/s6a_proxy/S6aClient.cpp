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
#include <iostream>
#include <utility>

#include "lte/protos/mconfig/mconfigs.pb.h"
#include "includes/MConfigLoader.h"
#include "S6aClient.h"
#include "includes/ServiceRegistrySingleton.h"
#include "itti_msg_to_proto_msg.h"
#include "feg/protos/s6a_proxy.pb.h"
#include "mme_config.h"
#include "common_defs.h"
#include "common_utility_funs.h"
extern "C" {
#include "log.h"
}

namespace grpc {
class Status;
}  // namespace grpc

using grpc::Status;
#define MME_SERVICE "mme"

namespace magma {
using namespace feg;

static bool read_hss_relay_enabled(void);

static bool read_mme_cloud_subscriberdb_enabled(void);

static const bool hss_relay_enabled = read_hss_relay_enabled();

static const bool cloud_subscriberdb_enabled =
    read_mme_cloud_subscriberdb_enabled();

// Relay enabled must be controlled by mconfigs in ONE place
// and mme app should be restarted on the config changes.
//
// The relay decision is encapsulated in S6a Client API, it is STATIC
// and CONSTANT during the application's lifetime and can only be queried
// (not changed)

bool get_s6a_relay_enabled(void) {
  return hss_relay_enabled;
}

bool get_cloud_subscriberdb_enabled(void) {
  return cloud_subscriberdb_enabled;
}

static bool read_hss_relay_enabled(void) {
  magma::mconfig::MME mconfig;
  magma::MConfigLoader loader;

  if (!loader.load_service_mconfig(MME_SERVICE, &mconfig)) {
    std::cout << "[INFO] Unable to load mconfig for mme, S6a relay is disabled"
              << std::endl;
    return false;  // default is - relay disabled
  }
  if (mconfig.relay_enabled()) {
    return true;
  }
  return mconfig.hss_relay_enabled();
}

static bool read_mme_cloud_subscriberdb_enabled(void) {
  magma::mconfig::MME mconfig;
  magma::MConfigLoader loader;

  if (!loader.load_service_mconfig(MME_SERVICE, &mconfig)) {
    std::cout << "[INFO] Unable to load mconfig for mme, cloud subscriberdb is "
                 "disabled"
              << std::endl;
    return false;  // default is - cloud subscriberdb disabled
  }
  return mconfig.cloud_subscriberdb_enabled();
}

S6aClient& S6aClient::get_s6a_proxy_instance() {
  static S6aClient s6a_proxy_instance(true);
  return s6a_proxy_instance;
}

S6aClient& S6aClient::get_subscriberdb_instance() {
  static S6aClient subscriberdb_instance(false);
  return subscriberdb_instance;
}

S6aClient& S6aClient::get_client_based_on_fed_mode(const char* imsi) {
  // get_client_based_on_fed_mode finds out the s6a_client (either subscribrdb
  // or FEG) based on imsi and fed map configured
  switch (match_fed_mode_map(imsi, LOG_S6A)) {
    case magma::mconfig::ModeMapItem_FederatedMode_SPGW_SUBSCRIBER:
    case magma::mconfig::ModeMapItem_FederatedMode_S8_SUBSCRIBER:
      return get_s6a_proxy_instance();
    case magma::mconfig::ModeMapItem_FederatedMode_LOCAL_SUBSCRIBER:
      return get_subscriberdb_instance();
    default:
      std::cout << "[ERROR] Unable to find appropriate fed mode for " << imsi
                << ". Using local s6a_cli" << std::endl;
      return get_subscriberdb_instance();
  }
}

S6aClient::S6aClient(bool enable_s6a_proxy_channel) {
  // Create channel based on relay_enabled, enable_s6a_proxy_channel and
  // cloud_subscriberdb_enabled flags.
  // If relay_enabled is true and enable_s6a_proxy_channel is true i.e federated
  // mode is SPGW_SUBSCRIBER or S8_SUBSCRIBER,
  // then create a channel towards the FeG.
  // Otherwise, create a channel towards either local or cloud-based
  // subscriberdb.
  if ((get_s6a_relay_enabled() == true) && (enable_s6a_proxy_channel)) {
    auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        "s6a_proxy", ServiceRegistrySingleton::CLOUD);
    // Create stub for S6aProxy gRPC service
    stub_ = S6aProxy::NewStub(channel);
  } else {
    auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        "subscriberdb", ServiceRegistrySingleton::LOCAL);
    // Create stub for subscriberdb gRPC service
    stub_ = S6aProxy::NewStub(channel);
  }

  std::thread resp_loop_thread([&]() { rpc_response_loop(); });
  resp_loop_thread.detach();
}

void S6aClient::purge_ue(
    const char* imsi, std::function<void(Status, PurgeUEAnswer)> callbk) {
  S6aClient& client = get_client_based_on_fed_mode(imsi);

  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto resp = new AsyncLocalResponse<PurgeUEAnswer>(
      std::move(callbk), RESPONSE_TIMEOUT);

  // Create a response reader for the `PurgeUE` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  PurgeUERequest puRequest;
  puRequest.set_user_name(imsi);
  auto resp_rdr = client.stub_->AsyncPurgeUE(
      resp->get_context(), puRequest, &client.queue_);

  // Set the reader for the response. This executes the `PurgeUE`
  // response using the response reader. When it is done, the callback stored
  // in `resp` will be called
  resp->set_response_reader(std::move(resp_rdr));
}

void S6aClient::authentication_info_req(
    const s6a_auth_info_req_t* const msg,
    std::function<void(Status, feg::AuthenticationInformationAnswer)> callbk) {
  S6aClient& client = get_client_based_on_fed_mode(msg->imsi);

  AuthenticationInformationRequest proto_msg =
      convert_itti_s6a_authentication_info_req_to_proto_msg(msg);
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto resp = new AsyncLocalResponse<AuthenticationInformationAnswer>(
      std::move(callbk), RESPONSE_TIMEOUT);

  // Create a response reader for the `authentication_info_req` RPC call.
  // This reader stores the client context, the request to pass in, and
  // the queue to add the response to when done
  auto resp_rdr = client.stub_->AsyncAuthenticationInformation(
      resp->get_context(), proto_msg, &client.queue_);

  // Set the reader for the response. This executes the `Auth_Info_Req`
  // response using the response reader. When it is done, the callback stored
  // in `resp` will be called
  resp->set_response_reader(std::move(resp_rdr));
}

void S6aClient::update_location_request(
    const s6a_update_location_req_t* const msg,
    std::function<void(Status, feg::UpdateLocationAnswer)> callbk) {
  S6aClient& client = get_client_based_on_fed_mode(msg->imsi);

  UpdateLocationRequest proto_msg =
      convert_itti_s6a_update_location_request_to_proto_msg(msg);
  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto resp = new AsyncLocalResponse<UpdateLocationAnswer>(
      std::move(callbk), RESPONSE_TIMEOUT);

  // Create a response reader for the `update_location_request` RPC call.
  // This reader stores the client context, the request to pass in, and
  // the queue to add the response to when done
  auto resp_rdr = client.stub_->AsyncUpdateLocation(
      resp->get_context(), proto_msg, &client.queue_);

  // Set the reader for the response. This executes the `Location Update Req`
  // response using the response reader. When it is done, the callback stored
  // in `resp` will be called
  resp->set_response_reader(std::move(resp_rdr));
}

}  // namespace magma
