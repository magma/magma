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
#include <string>
#include <functional>
#include <memory>

#include "orc8r/gateway/c/common/async_grpc/GRPCReceiver.hpp"
#include "lte/protos/sms_orc8r.grpc.pb.h"
#include "lte/gateway/c/core/oai/include/sgs_messages_types.hpp"

extern "C" {
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"

namespace grpc {
class Status;
}  // namespace grpc
namespace magma {
namespace orc8r {
class Void;
}  // namespace orc8r
}  // namespace magma
}

namespace magma {

using lte::SMSOrc8rService;
using orc8r::Void;

/**
 * SMSOrc8rClient is the main client for sending message to orc8r.
 */
class SMSOrc8rClient : public GRPCReceiver {
 public:
  /**
   * SGsAP-UPLINK-UNITDATA
   */
  static void send_uplink_unitdata(
      const itti_sgsap_uplink_unitdata_t* msg,
      std::function<void(grpc::Status, Void)> callback);

 public:
  SMSOrc8rClient(SMSOrc8rClient const&) = delete;
  void operator=(SMSOrc8rClient const&) = delete;

 private:
  SMSOrc8rClient();
  static SMSOrc8rClient& get_instance();
  std::unique_ptr<SMSOrc8rService::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 3;  // seconds
};

}  // namespace magma
