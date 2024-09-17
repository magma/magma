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

#include "lte/gateway/c/connection_tracker/src/PacketGenerator.hpp"

#include <glog/logging.h>
#include <tins/ethernetII.h>
#include <tins/ip.h>
#include <tins/ip_address.h>
#include <tins/packet_sender.h>
#include <tins/pdu.h>
#include <tins/tcp.h>
#include <tins/udp.h>
#include <iostream>
#include <string>

#include "orc8r/gateway/c/common/logging/magma_logging.hpp"

namespace magma {
namespace lte {

using Tins::EthernetII;
using Tins::IP;
using Tins::IPv4Address;
using Tins::NetworkInterface;
using Tins::PacketSender;
using Tins::TCP;
using Tins::UDP;

PacketGenerator::PacketGenerator(const std::string& iface_name,
                                 const std::string& pkt_dst_mac,
                                 const std::string& pkt_src_mac)
    : iface_name_(iface_name),
      pkt_dst_mac_(pkt_dst_mac),
      pkt_src_mac_(pkt_src_mac) {
  iface_ = NetworkInterface(iface_name_);
  MLOG(MINFO) << "Using interface " << iface_name_.c_str()
              << " for pkt generation";
}

bool PacketGenerator::send_packet(struct flow_information* flow) {
  PacketSender sender;

  // Random mac header for our internal packets
  EthernetII eth_(pkt_dst_mac_, pkt_src_mac_);
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

}  // namespace lte
}  // namespace magma