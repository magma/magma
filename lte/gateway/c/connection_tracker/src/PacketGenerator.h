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

#include <tins/tins.h>

#include "orc8r/gateway/c/common/logging/magma_logging.h"

struct flow_information {
  uint32_t saddr;    /* Source address */
  uint32_t daddr;    /* Destination address */
  uint32_t l4_proto; /* Layer4 Proto ID */
  uint16_t sport;    /* Source port */
  uint16_t dport;    /* Destination port */
};

namespace magma {
namespace lte {

class PacketGenerator {
 public:
  PacketGenerator(const std::string& iface_name, const std::string& pkt_dst_mac,
                  const std::string& pkt_src_mac);
  /**
   * Send packet based on provided flow information
   * @param flow_information - flow_information
   * @return true if the operation was successful
   */
  bool send_packet(struct flow_information* flow);

 private:
  std::string iface_name_;
  std::string pkt_dst_mac_;
  std::string pkt_src_mac_;
  Tins::NetworkInterface iface_;
};

}  // namespace lte
}  // namespace magma