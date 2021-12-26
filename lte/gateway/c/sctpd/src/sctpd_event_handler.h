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

#include "lte/gateway/c/sctpd/src/sctp_connection.h"

#include "lte/gateway/c/sctpd/src/sctpd_uplink_client.h"

namespace magma {
namespace sctpd {

// Sctp handler that relays events to MME/AMF over GRPC
class SctpdEventHandler : public SctpEventHandler {
 public:
  // Construct SctpdEventHandler that communicates to MME/AMF over client
  explicit SctpdEventHandler(SctpdUplinkClient& client);

  // Relay new assocation to MME/AMF over GRPC
  int HandleNewAssoc(uint32_t ppid, uint32_t assoc_id, uint32_t instreams,
                     uint32_t outstreams, std::string& ran_cp_ipaddr) override;

  // Relay close assocation to MME/AMF over GRPC
  void HandleCloseAssoc(uint32_t ppid, uint32_t assoc_id, bool reset) override;

  // Relay new message to MME over GRPC
  void HandleRecv(uint32_t ppid, uint32_t assoc_id, uint32_t stream,
                  const std::string& payload) override;

 private:
  SctpdUplinkClient& _client;
};

}  // namespace sctpd
}  // namespace magma
