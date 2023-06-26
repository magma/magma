#!/usr/bin/env python3

"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import binascii
import socket
import sys
import time

from dpkt import ethernet, icmp6, ip, ip6
from magma.pipelined.ifaces import get_mac_address_from_iface

dst = sys.argv[1]

print("dst: %s" % dst)

src_ip = socket.inet_pton(socket.AF_INET6, "2001::2")
dst_ip = socket.inet_pton(socket.AF_INET6, dst)
src_mac = get_mac_address_from_iface("eth0")
src_mac_as_bytes = binascii.unhexlify(src_mac.replace(':', ''))

pkt = icmp6.ICMP6(type=icmp6.ICMP6_ECHO_REQUEST, data=icmp6.ICMP6.Echo())
pkt = ip6.IP6(
    src=src_ip, dst=dst_ip,
    plen=8, nxt=ip.IP_PROTO_ICMP6,
    hlim=64, data=bytes(pkt),
)
pkt = ethernet.Ethernet(
    src=src_mac_as_bytes, dst=b'\xff' * 6,
    type=ethernet.ETH_TYPE_IP6, data=bytes(pkt),
)

with socket.socket(socket.AF_PACKET, socket.SOCK_RAW) as sock:
    sock.bind(("gtp_br0", 0))
    sock.send(bytes(pkt))
    time.sleep(0.5)
