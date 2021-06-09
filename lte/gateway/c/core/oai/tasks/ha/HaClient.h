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

#include "lte/protos/ha_orc8r.grpc.pb.h"
#include "includes/GRPCReceiver.h"

namespace magma {

/**
 * HaClient is the main asynchronous client for interacting with Ha Daemon
 * at the orchestrator.
 */
class HaClient : public GRPCReceiver {
 public:
  /**
   * Proxy a purge gRPC call to s6a_proxy
   */
  static void get_eNB_offload_state(
      std::function<void(grpc::Status, lte::GetEnodebOffloadStateResponse)>
          callback);

  HaClient(HaClient const&) = delete;
  void operator=(HaClient const&) = delete;

 private:
  HaClient();
  static HaClient& get_instance();
  std::unique_ptr<lte::Ha::Stub> stub_;
  static const uint32_t RESPONSE_TIMEOUT = 3;  // seconds
};

}  // namespace magma
