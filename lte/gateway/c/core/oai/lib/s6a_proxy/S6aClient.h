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

#pragma once

#include <gmp.h>
#include <grpc++/grpc++.h>
#include <stdint.h>
#include <functional>
#include <memory>

#include "feg/protos/s6a_proxy.grpc.pb.h"
#include "includes/GRPCReceiver.h"
#include "s6a_messages_types.h"

extern "C" {
#include "intertask_interface.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace feg {
class AuthenticationInformationAnswer;
class PurgeUEAnswer;
class UpdateLocationAnswer;
}  // namespace feg
}  // namespace magma
}

using grpc::Status;

namespace magma {

/**
 * S6aClient is the main asynchronous client for interacting with s6a_proxy.
 * Responses will come in a queue and call the callback passed
 * To start the client and make sure it receives calls, one must call the
 * rpc_response_loop method defined in the GRPCReceiver base class
 */
class S6aClient : public GRPCReceiver {
 public:
  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void purge_ue(
      const char* imsi,
      std::function<void(Status, feg::PurgeUEAnswer)> callback);

  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void authentication_info_req(
      const s6a_auth_info_req_t* const msg,
      std::function<void(Status, feg::AuthenticationInformationAnswer)> callbk);

  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void update_location_request(
      const s6a_update_location_req_t* const msg,
      std::function<void(Status, feg::UpdateLocationAnswer)> callback);

 public:
  S6aClient(S6aClient const&) = delete;
  void operator=(S6aClient const&) = delete;

 private:
  S6aClient(bool enable_s6a_proxy_channel);
  static S6aClient& get_instance();
  static S6aClient& get_s6a_proxy_instance();
  static S6aClient& get_subscriberdb_instance();
  static S6aClient& get_client_based_on_fed_mode(const char* imsi);
  std::unique_ptr<feg::S6aProxy::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
};

// There are 3 services which can handle authentication:
// 1) Local subscriberdb
// 2) Subscriberdb in the cloud (EPS Authentication)
// 3) S6a Proxy running in the FeG
// When relay_enabled is true, then auth requests are sent to the S6a Proxy.
// Otherwise, if cloud_subscriberdb_enabled is true, then auth requests are
// sent to the EPS Authentication service.
// If neither flag is true, then a local instance of subscriberdb receives the
// auth messages.
bool get_s6a_relay_enabled(void);
bool get_cloud_subscriberdb_enabled(void);

}  // namespace magma
