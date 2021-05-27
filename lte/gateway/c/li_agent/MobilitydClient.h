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
#pragma once

#include <grpc++/grpc++.h>
#include <arpa/inet.h>
#include <memory>
#include <string>

#include "lte/protos/mobilityd.grpc.pb.h"
#include "GRPCReceiver.h"

namespace grpc {
class Channel;
class ClientContext;
class Status;
}  // namespace grpc

namespace magma {
namespace lte {

/*
 * gRPC client for MobilityService
 */
class MobilityClient : public GRPCReceiver {
 public:
  virtual ~MobilityClient() = default;
  /*
   * Get the subscriber id given its allocated IPv4 address. If the address
   * isn't associated with a subscriber, then it returns an error
   * @param addr: ipv4 address of subscriber
   * @param imsi (out): contains the imsi of the associated subscriber if it
   *                    exists
   * @return 0 on success
   * @return -RPC_STATUS_NOT_FOUND if IPv4 address is not found
   */
  int GetSubscriberIDFromIP(const struct in_addr& addr, std::string* imsi);

 public:
  static MobilityClient& getInstance();

  MobilityClient(MobilityClient const&) = delete;
  void operator=(MobilityClient const&) = delete;

 private:
  MobilityClient();
  static const uint32_t RESPONSE_TIMEOUT = 10;  // seconds
  std::unique_ptr<MobilityService::Stub> stub_{};
};

}  // namespace lte
}  // namespace magma
