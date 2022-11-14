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

#include "lte/gateway/c/core/oai/lib/s6a_proxy/S6aClient.hpp"

#include <grpcpp/impl/codegen/async_unary_call.h>
#include <thread>  // std::thread
#include <iostream>
#include <utility>

#include "feg/protos/s6a_proxy.pb.h"
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/common/common_utility_funs.hpp"
#include "lte/gateway/c/core/oai/include/mme_config.hpp"
#include "lte/gateway/c/core/oai/lib/s6a_proxy/itti_msg_to_proto_msg.hpp"
#include "lte/protos/mconfig/mconfigs.pb.h"
#include "orc8r/gateway/c/common/config/MConfigLoader.hpp"
#include "orc8r/gateway/c/common/service_registry/ServiceRegistrySingleton.hpp"
extern "C" {
#include "lte/gateway/c/core/oai/common/log.h"
}

namespace grpc {
class Status;
}  // namespace grpc

using grpc::Status;
#define MME_SERVICE "mme"

namespace magma {
using namespace feg;

// initialize with -1: not yet calculated/cached
static int hss_relay_enabled = -1;
static int mme_cloud_subscriberdb_enabled = -1;

// Relay enabled must be controlled by mconfigs in ONE place
// and mme app should be restarted on the config changes.
//
// The relay decision is encapsulated in S6a Client API, it is STATIC
// and CONSTANT during the application's lifetime and can only be queried
// (not changed)

bool S6aClient::get_s6a_relay_enabled() {
  if (hss_relay_enabled == -1) {
    // cache result
    hss_relay_enabled = (S6aClient::read_hss_relay_enabled() ? 1 : 0);
  }
  return hss_relay_enabled == 1;
}

bool S6aClient::get_cloud_subscriberdb_enabled() {
  if (mme_cloud_subscriberdb_enabled == -1) {
    // cache result
    mme_cloud_subscriberdb_enabled =
        (S6aClient::read_mme_cloud_subscriberdb_enabled() ? 1 : 0);
  }
  return mme_cloud_subscriberdb_enabled == 1;
}

bool S6aClient::read_hss_relay_enabled() {
  magma::mconfig::MME mconfig;

  if (!magma::load_service_mconfig_from_file(MME_SERVICE, &mconfig)) {
    std::cout << "[INFO] Unable to load mconfig for mme, S6a relay is disabled"
              << std::endl;
    return false;  // default is - relay disabled
  }
  if (mconfig.relay_enabled()) {
    return true;
  }
  return mconfig.hss_relay_enabled();
}

bool S6aClient::read_mme_cloud_subscriberdb_enabled() {
  magma::mconfig::MME mconfig;

  if (!magma::load_service_mconfig_from_file(MME_SERVICE, &mconfig)) {
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
  } else if (get_cloud_subscriberdb_enabled()) {
    auto channel = ServiceRegistrySingleton::Instance()->GetGrpcChannel(
        "eps_authentication", ServiceRegistrySingleton::CLOUD);
    // Create S6aProxy stub for eps_authentication gRPC service
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

void S6aClient::purge_ue(const char* imsi,
                         std::function<void(Status, PurgeUEAnswer)> callbk) {
  S6aClient& client = get_client_based_on_fed_mode(imsi);

  // Create a raw response pointer that stores a callback to be called when the
  // gRPC call is answered
  auto resp = new AsyncLocalResponse<PurgeUEAnswer>(std::move(callbk),
                                                    RESPONSE_TIMEOUT);

  // Create a response reader for the `PurgeUE` RPC call. This reader
  // stores the client context, the request to pass in, and the queue to add
  // the response to when done
  PurgeUERequest puRequest;
  puRequest.set_user_name(imsi);
  auto resp_rdr = client.stub_->AsyncPurgeUE(resp->get_context(), puRequest,
                                             &client.queue_);

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
  auto resp = new AsyncLocalResponse<UpdateLocationAnswer>(std::move(callbk),
                                                           RESPONSE_TIMEOUT);

  // Create a response reader for the `update_location_request` RPC call.
  // This reader stores the client context, the request to pass in, and
  // the queue to add the response to when done
  auto resp_rdr = client.stub_->AsyncUpdateLocation(resp->get_context(),
                                                    proto_msg, &client.queue_);

  // Set the reader for the response. This executes the `Location Update Req`
  // response using the response reader. When it is done, the callback stored
  // in `resp` will be called
  resp->set_response_reader(std::move(resp_rdr));
}

void S6aClient::convert_ula_to_subscriber_data(
    feg::UpdateLocationAnswer response, magma::lte::SubscriberData* sub_data) {
  if (response.apn_size() < 1) {
    std::cout << "No APN configurations received" << std::endl;
    return;
  }
  std::cout << "Converting ULA TO Subscriber Data object" << std::endl;
  for (int i = 0; i < response.apn_size(); i++) {
    auto apn = response.apn(i);
    auto sub_apn_config = sub_data->mutable_non_3gpp()->add_apn_config();
    if (apn.context_id() != 0) {
      sub_apn_config->set_context_id(apn.context_id());
    }

    if (apn.service_selection().size() > 0) {
      sub_apn_config->set_service_selection(apn.service_selection());
    }

    if (apn.has_qos_profile()) {
      auto qos_profile = sub_apn_config->mutable_qos_profile();
      if (apn.qos_profile().class_id()) {
        qos_profile->set_class_id(apn.qos_profile().class_id());
      }
      if (apn.qos_profile().priority_level()) {
        qos_profile->set_priority_level(apn.qos_profile().priority_level());
      }
      if (apn.qos_profile().preemption_capability()) {
        qos_profile->set_preemption_capability(
            apn.qos_profile().preemption_capability());
      }
      if (apn.qos_profile().preemption_vulnerability()) {
        qos_profile->set_preemption_vulnerability(
            apn.qos_profile().preemption_vulnerability());
      }
    }

    if (apn.has_ambr()) {
      auto ambr = sub_apn_config->mutable_ambr();
      if (apn.ambr().max_bandwidth_dl() != 0) {
        ambr->set_max_bandwidth_dl(apn.ambr().max_bandwidth_dl());
      }
      if (apn.ambr().max_bandwidth_ul() != 0) {
        ambr->set_max_bandwidth_ul(apn.ambr().max_bandwidth_ul());
      }

      ambr->set_br_unit(
          (magma::lte::AggregatedMaximumBitrate_BitrateUnitsAMBR)apn.ambr()
              .unit());
    }

    sub_apn_config->set_pdn((magma::lte::APNConfiguration_PDNType)apn.pdn());

    // Only the first IP is assigned to the subscriber in the current
    // implementation
    if (apn.served_party_ip_address_size() > 0) {
      sub_apn_config->set_assigned_static_ip(apn.served_party_ip_address(0));
    }

    if (apn.has_resource()) {
      auto resource = sub_apn_config->mutable_resource();
      if (apn.resource().apn_name().size() > 0) {
        resource->set_apn_name(apn.resource().apn_name());
      }
      if (apn.resource().gateway_ip().size() > 0) {
        resource->set_gateway_ip(apn.resource().gateway_ip());
      }
      if (apn.resource().gateway_mac().size() > 0) {
        resource->set_gateway_mac(apn.resource().gateway_mac());
      }
      if (apn.resource().vlan_id() != 0) {
        resource->set_vlan_id(apn.resource().vlan_id());
      }
    }
  }
}

}  // namespace magma
