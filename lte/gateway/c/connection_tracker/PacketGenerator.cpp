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

using namespace Tins;

int send_packet(struct flow_information *flow) {
    NetworkInterface iface = NetworkInterface("mtr0");

    /* Retrieve this structure which holds the interface's IP,
     * broadcast, hardware address and the network mask.
     */
    NetworkInterface::Info info = iface.addresses();

    /* Create an Ethernet II PDU which will be sent to
     * 77:22:33:11:ad:ad using the default interface's hardware
     * address as the sender.
     */
    EthernetII eth("77:22:33:11:ad:ad", info.hw_addr);

    /* Create an IP PDU, with 192.168.0.1 as the destination address
     * and the default interface's IP address as the sender.
     */
    IPv4Address ip_s = IPv4Address(flow->saddr);
    IPv4Address ip_d = IPv4Address(flow->daddr);
    eth /= IP(ip_s, ip_d);

    /* Create a TCP PDU using 13 as the destination port, and 15
     * as the source port.
     */
    if (flow->l4_proto==6) {
        eth /= TCP(flow->dport, flow->sport);
    } else if (flow->l4_proto==17) {
        eth /= UDP(flow->dport, flow->sport);
    } else {
        // unsupported
        return -1;
    }

    PacketSender sender;
    sender.send(eth, iface);
    
    return 0;
}
