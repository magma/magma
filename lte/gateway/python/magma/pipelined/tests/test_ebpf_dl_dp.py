"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import logging
import socket
import subprocess
import unittest
from typing import List

from lte.protos.mobilityd_pb2 import IPAddress
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.ebpf.ebpf_manager import EbpfManager
from scapy.all import AsyncSniffer
from scapy.layers.inet import IP

PKT_SCRIPT = "/home/vagrant/magma/lte/gateway/python/magma/pipelined/tests/script/ip-packet.py"
PY_PATH = "/home/vagrant/build/python/bin/python"
UL_HANDLER = "/home/vagrant/magma/lte/gateway/python/magma/pipelined/ebpf/ebpf_ul_handler.c"
DL_HANDLER = "/home/vagrant/magma/lte/gateway/python/magma/pipelined/ebpf/ebpf_dl_handler.c"
BPF_HEADER_PATH = "/home/vagrant/magma/orc8r/gateway/c/common/ebpf/"


# This test works when ran separately.
@unittest.skip("AsyncSniffer is not working")
class eBpfDatapathDLTest(unittest.TestCase):
    NS_NAME = 'ens1'
    gtp_veth = "enb0"
    gtp_veth_ns = "enb1"

    sgi_veth = "sgi0"
    sgi_veth1 = "sgi1"

    sgi_veth_ip = "3.3.3.3"
    inner_src_ip = '2.2.2.2'
    inner_dst_ip = '2.2.2.1'

    gtp_pkt_dst = '11.1.1.1'
    gtp_pkt_src = '11.1.1.2'

    gtp_tunnel_id = 101

    imsi = '122321231222333'

    packet_cap1: List = []
    sniffer = None
    ebpf_man = None

    @classmethod
    def setUpClass(cls):
        pass

    @classmethod
    def setUpClassDevices(cls):
        BridgeTools.delete_ns_all()

        BridgeTools.create_veth_pair(cls.gtp_veth, cls.gtp_veth_ns)
        BridgeTools.ifup_netdev(cls.gtp_veth, cls.gtp_pkt_dst + "/24")

        BridgeTools.create_veth_pair(cls.sgi_veth, cls.sgi_veth1)

        BridgeTools.create_ns_and_move_veth(cls.NS_NAME, cls.sgi_veth1, cls.inner_src_ip + "/24")

        BridgeTools.ifup_netdev(cls.sgi_veth, cls.inner_dst_ip + "/24")
        BridgeTools.ifup_netdev(cls.gtp_veth_ns, cls.gtp_pkt_src + "/24")

        gw_ip = IPAddress(version=IPAddress.IPV4, address=socket.inet_aton(cls.sgi_veth_ip))

        cls.ebpf_man = EbpfManager(cls.sgi_veth, cls.gtp_veth, gw_ip, bpf_ul_file=UL_HANDLER, bpf_dl_file=DL_HANDLER, bpf_header_path=BPF_HEADER_PATH)
        cls.ebpf_man.detach_dl_ebpf()
        cls.ebpf_man.attach_dl_ebpf()

        cls.sniffer = AsyncSniffer(
            iface='gtpu_sys_2152',
            store=False,
            prn=cls.pkt_cap_fun,
        )
        cls.sniffer.start()

    @classmethod
    def sendPacket(cls, gtp_src, gtp_dst):
        try:
            xmit_cmd = [
                "ip", "netns", "exec", cls.NS_NAME,
                PY_PATH,
                PKT_SCRIPT,
                gtp_src, gtp_dst,
                cls.sgi_veth1,
            ]
            subprocess.check_call(xmit_cmd)
            logging.debug("del ns %s", xmit_cmd)

        except subprocess.CalledProcessError as e:
            logging.debug("Error while xmit from ns: %s", e)

    @classmethod
    def tearDownClassDevices(cls):
        cls.ebpf_man.detach_ul_ebpf()
        cls.sniffer.stop()
        BridgeTools.delete_ns_all()
        BridgeTools.delete_veth(cls.gtp_veth)
        BridgeTools.delete_veth(cls.sgi_veth)

    @classmethod
    def pkt_cap_fun(cls, packet):
        print("got packet: %s", packet)
        print(packet.show(dump=True))
        cls.packet_cap1.append(packet)

    @classmethod
    def count_udp_packet(cls):
        cnt = 0
        for pkt in cls.packet_cap1:
            # print(pkt.show(dump=True))
            if IP in pkt:
                if pkt[IP].src == cls.inner_src_ip and pkt[IP].dst == cls.inner_dst_ip:
                    cnt = cnt + 1
        return cnt

    def testEbpfDlFrw1(self):
        cls = self.__class__
        cls.setUpClassDevices()
        cls.sendPacket(cls.inner_src_ip, cls.inner_dst_ip)
        self.assertEqual(len(cls.packet_cap1), 0)

        cls.ebpf_man.add_dl_entry(cls.inner_dst_ip, cls.gtp_pkt_dst, cls.gtp_tunnel_id, cls.imsi)
        cls.sendPacket(cls.inner_src_ip, cls.inner_dst_ip)

        self.assertEqual(cls.count_udp_packet(), 1)
        cls.sendPacket(cls.inner_src_ip, cls.inner_dst_ip)

        self.assertEqual(cls.count_udp_packet(), 2)

        cls.ebpf_man.del_dl_entry(cls.inner_dst_ip)
        cls.ebpf_man.print_dl_map()
        cls.sendPacket(cls.inner_src_ip, cls.inner_dst_ip)
        self.assertEqual(cls.count_udp_packet(), 2)

        cls.tearDownClassDevices()
