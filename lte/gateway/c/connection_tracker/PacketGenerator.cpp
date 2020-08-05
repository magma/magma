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

#include "PacketGenerator.h"

namespace magma {

using namespace Tins;

PacketGenerator::PacketGenerator(std::string iface_name)
    : iface_name_(iface_name) {
  iface_ = NetworkInterface(iface_name_);
  MLOG(MINFO) << "Using interface " << iface_name_.c_str()
              << "for pkt generation";
}

bool PacketGenerator::send_packet(struct flow_information* flow) {
  PacketSender sender;

  // Random mac header for our internal packets
  EthernetII eth_("33:aa:99:33:aa:00", "55:11:44:ee:00:00");
  eth_ /= IP(IPv4Address(flow->saddr), IPv4Address(flow->daddr));

  if (flow->l4_proto == 6) {
    eth_ /= TCP(flow->dport, flow->sport);
  } else if (flow->l4_proto == 17) {
    eth_ /= UDP(flow->dport, flow->sport);
  } else {
    MLOG(MDEBUG) << "Encountered unsupported protocol, not sending pkt";
    return false;
  }

  sender.send(eth_, iface_);

  return true;
}

}  // namespace magma