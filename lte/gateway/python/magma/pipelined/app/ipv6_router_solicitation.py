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
import netifaces
import ipaddress
from pprint import pformat

from ryu.controller import ofp_event
from ryu.lib.packet import packet
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, load_passthrough, \
    TUN_PORT_REG

from scapy.arch import get_if_hwaddr, get_if_addr6
from scapy.data import ETHER_BROADCAST, ETH_P_ALL
from scapy.error import Scapy_Exception
from scapy.layers.l2 import ARP, Ether, Dot1Q
from scapy.layers.inet6 import IPv6, ICMPv6ND_RA, ICMPv6NDOptSrcLLAddr, \
    ICMPv6NDOptPrefixInfo, ICMPv6ND_NA
from scapy.sendrecv import srp1, sendp
from scapy.all import raw

from ryu.controller import dpset
from ryu.lib.packet import ether_types, icmpv6, ipv6
from ryu.ofproto.inet import IPPROTO_ICMPV6

from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL, OFPP_CONTROLLER
from ryu.lib.packet import packet
from ryu.lib.packet import ethernet
from ryu.lib.packet import ether_types
from ryu.lib.packet import ipv4
from ryu.lib.packet import in_proto
from ryu.lib.packet import tcp


class IPV6RouterSolicitationController(MagmaController):
    """
    IPV6RouterSolicitationController responds to ipv6 router solicitation
    messages

    (1) Listens to flows with IPv6 src address prefixed with ""fe80".
    (2) Extracts interface ID (lower 64 bits) from the Router Solicitation
        message.
    (3) Performs a look up to find the IPv6 prefix that corresponds to the
        interface ID. The look up can be done using a local look up table that
        is updated during session creation where the full 128 bit IPv6 address
        assigned to UE is provided.
    (4) Generates a router advertisement message targeting the GTP tunnel.
    """
    APP_NAME = 'ipv6_router_solicitation'
    APP_TYPE = ControllerType.PHYSICAL

    # Inherited from app_manager.RyuApp
    _CONTEXTS = {
        'dpset': dpset.DPSet,
    }

    DEVICE_MULTICAST = 'ff02::1'
    ROUTER_MULTICAST = 'ff02::2'

    MAC_MULTICAST = '33:33:00:00:00:01'

    ICMPV6_RS_TYPE = 133
    ICMPV6_NS_TYPE = 135

    def __init__(self, *args, **kwargs):
        super(IPV6RouterSolicitationController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.setup_type = kwargs['config']['setup_type']
        addrs = netifaces.ifaddresses(kwargs['config']['bridge_name'])
        self._ll_addr = addrs[netifaces.AF_LINK][0]['addr']
        if '%' in addrs[netifaces.AF_INET6][0]['addr']:
            ipv6_str, _ = addrs[netifaces.AF_INET6][0]['addr'].split('%')
            self._ipv6_src = ipv6_str

        self.logger.error(addrs)
        self.logger.error(self._ll_addr)
        self.logger.error(self._ipv6_src)

        self._prefix_len = 64
        self._datapath = None

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._install_default_ipv6_flows(datapath)

    def _install_default_flows(self, datapath):
        """
        Add low priority flow to forward to next app
        """

        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             match=MagmaMatch(), actions=[],
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def _install_default_ipv6_flows(self, datapath):
        """
        Install flows that match on RS/NS and trigger packet in message, that
        will respond with RA/NA.
        """
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser

        match_rs = MagmaMatch(eth_type=ether_types.ETH_TYPE_IPV6,
                              ipv6_src='fe80::/10',
                              ip_proto=IPPROTO_ICMPV6,
                              icmpv6_type=icmpv6.ND_ROUTER_SOLICIT,
                              direction=Direction.IN)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=match_rs, actions=[],
                              priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=self.next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)

        match_ns = MagmaMatch(eth_type=ether_types.ETH_TYPE_IPV6,
                              ipv6_src='fe80::/10',
                              ip_proto=IPPROTO_ICMPV6,
                              icmpv6_type=icmpv6.ND_NEIGHBOR_SOLICIT,
                              direction=Direction.IN)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=match_ns, actions=[],
                              priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=self.next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)

    def _send_router_advertisement(self, prefix: str, output_port):
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser

        pkt = packet.Packet()
        pkt.add_protocol(
            ethernet.ethernet(
                dst=self.MAC_MULTICAST,
                src=self._ll_addr,
                ethertype=ether_types.ETH_TYPE_IPV6,
            )
        )
        pkt.add_protocol(
            ipv6.ipv6(
                dst=self.DEVICE_MULTICAST,
                src=self._ipv6_src,
                nxt=in_proto.IPPROTO_ICMPV6,
            )
        )
        pkt.add_protocol(
            icmpv6.icmpv6(
                type_=icmpv6.ND_ROUTER_ADVERT,
                data=icmpv6.nd_router_advert(
                    options=[
                        icmpv6.nd_option_sla(
                            hw_src=self._ll_addr,
                        ),
                        icmpv6.nd_option_pi(
                            pl=self._prefix_len,
                            prefix=prefix,
                        )
                    ]
                ),
            )
        )
        pkt.serialize()

        actions_out = [
            #parser.NXActionResubmitTable(table_id=99),
            parser.OFPActionOutput(port=output_port)]
        out = parser.OFPPacketOut(datapath=self._datapath,
                                  buffer_id=ofproto.OFP_NO_BUFFER,
                                  in_port=ofproto.OFPP_CONTROLLER,
                                  actions=actions_out,
                                  data=pkt.data)
        self._datapath.send_msg(out)

    def _send_neighbor_advertisement(self, output_port):
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser

        pkt = packet.Packet()
        pkt.add_protocol(
            ethernet.ethernet(
                dst=self.MAC_MULTICAST,
                src=self._ll_addr,
                ethertype=ether_types.ETH_TYPE_IPV6,
            )
        )
        pkt.add_protocol(
            ipv6.ipv6(
                dst=self.DEVICE_MULTICAST,
                src=self._ipv6_src,
                nxt=in_proto.IPPROTO_ICMPV6,
            )
        )
        pkt.add_protocol(
            icmpv6.icmpv6(
                type_=icmpv6.ND_NEIGHBOR_ADVERT,
                data=icmpv6.nd_router_advert(
                    options=[
                        icmpv6.nd_option_sla(
                            hw_src=self._ll_addr,
                        )
                    ]
                ),
            )
        )
        pkt.serialize()

        actions_out = [
            parser.OFPActionOutput(port=output_port)]
        out = parser.OFPPacketOut(datapath=self._datapath,
                                  buffer_id=ofproto.OFP_NO_BUFFER,
                                  in_port=ofproto.OFPP_CONTROLLER,
                                  actions=actions_out,
                                  data=pkt.data)
        self._datapath.send_msg(out)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _parse_pkt_in(self, ev):
        """
        Learn action to process PacketIn DHCP packets, dhcp ack packets will
        be used to learn the ARP entry for the UE to install rules in the arp
        table. The DHCP packets will then be sent thorugh the pipeline.
        """
        msg = ev.msg

        self.logger.error("______ PKT ______")
        self.logger.error("______ PKT ______")
        self.logger.error("______ PKT ______")
        self.logger.error("______ PKT ______")
        self.logger.error(ev.msg)
        self.logger.error(pformat(ev.msg))
        if self.tbl_num != msg.table_id:
            # Intended for other application
            pkt = packet.Packet(msg.data)
            for p in pkt.protocols:
                self.logger.error(p)
            return

        in_port = ev.msg.match['in_port']

        pkt = packet.Packet(msg.data)
        for p in pkt.protocols:
            self.logger.error(p)

        ipv6_header = pkt.get_protocols(ipv6.ipv6)[0]

        icmpv6_header = pkt.get_protocols(icmpv6.icmpv6)[0]

        prefix = self.get_custom_prefix(ipv6_header)
        if icmpv6_header.type_ == self.ICMPV6_RS_TYPE:
            self.logger.error("Recieved router soli MSG---------------")
            self._send_router_advertisement(prefix, in_port)
        elif icmpv6_header.type_ == self.ICMPV6_NS_TYPE:
            self.logger.error("Recieved neighbor soli MSG---------------")
            self._send_neighbor_advertisement(in_port)

        self.logger.error("______ PKT ______")
        self.logger.error("______ PKT ______")
        self.logger.error("______ PKT ______")
        self.logger.error("______ PKT ______")

    def handle_restart(self):
        pass

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def get_custom_prefix(self, ipv6_header):
        """
        Parse the ipv6 header and retrieve the custom UE prefix by matching on
        unique interface ID

        (config prefix part(x) + unique session prefix part(64-x) + unique interface part(64) )
        """
        ip_block = ipaddress.ip_address(ipv6_header.src)
        self.logger.error(int(ip_block) & 0xffffffffffffffff)
        ip_block2 = ipaddress.ip_address(int(ip_block) & 0xffffffffffffffff)
        self.logger.error(ip_block2)

        return str(ip_block2)