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
import ipaddress
from pprint import pformat

from ryu.controller import ofp_event
from ryu.lib.packet import packet
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, load_passthrough

from scapy.arch import get_if_hwaddr, get_if_addr
from scapy.data import ETHER_BROADCAST, ETH_P_ALL
from scapy.error import Scapy_Exception
from scapy.layers.l2 import ARP, Ether, Dot1Q
from scapy.layers.inet6 import IPv6, ICMPv6ND_RA, ICMPv6NDOptSrcLLAddr, \
    ICMPv6NDOptPrefixInfo, ICMPv6ND_NA
from scapy.sendrecv import srp1, sendp

from ryu.controller import dpset
from ryu.lib.packet import ether_types, icmpv6, ipv6
from ryu.ofproto.inet import IPPROTO_ICMPV6

from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL


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

    ICMPV6_RS_TYPE = 133
    ICMPV6_RA_TYPE = 134
    ICMPV6_NS_TYPE = 135
    ICMPV6_NA_TYPE = 136

    def __init__(self, *args, **kwargs):
        super(IPV6RouterSolicitationController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.setup_type = kwargs['config']['setup_type']
        self.gtp_tunnel = 'gtp_sys_2152'
        self._ipv6_src = 'fe80::24c3:d0ff:fef3:dd82'
        self._ll_addr = get_if_hwaddr(kwargs['config']['bridge_name'])
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
                              icmpv6_type=self.ICMPV6_RS_TYPE,
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
                              icmpv6_type=self.ICMPV6_NS_TYPE,
                              direction=Direction.IN)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=match_ns, actions=[],
                              priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=self.next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)

        # For testing -----
        # flows.add_output_flow(self._datapath, self.tbl_num,
        #                       match=MagmaMatch(), actions=[],
        #                       priority=flows.PASSTHROUGH_PRIORITY+1000,
        #                       output_port=ofproto.OFPP_CONTROLLER,
        #                       max_len=ofproto.OFPCML_NO_BUFFER)
        # flows.add_output_flow(datapath, self.tbl_num,
        #                       match=MagmaMatch(), actions=[],
        #                       priority=flows.PASSTHROUGH_PRIORITY,
        #                       output_port=ofproto.OFPP_CONTROLLER,
        #                       copy_table=self.next_table,
        #                       max_len=ofproto.OFPCML_NO_BUFFER)

    def _send_router_advertisement(self, prefix: str, output_port):
        try:
            port = OFPP_LOCAL
            ''' 
            gw_ip = ipaddress.ip_address(ip.address)
            self.logger.debug("sending arp via egress: %s", self.config.non_nat_arp_egress_port)
            eth_mac_src = get_if_hwaddr(self.config.non_nat_arp_egress_port)
            psrc = "0.0.0.0"
            egress_port_ip = get_if_addr(self.config.non_nat_arp_egress_port)
            if egress_port_ip:
                psrc = egress_port_ip

            pkt = Ether(dst=ETHER_BROADCAST, src=eth_mac_src)
            if vlan != "":
                pkt /= Dot1Q(vlan=int(vlan))
            pkt /= ARP(op="who-has", pdst=gw_ip, hwsrc=et3h_mac_src, psrc=psrc)
            self.logger.debug("ARP Req pkt %s", pkt.show(dump=True))
            
            '''
            pkt = Ether()#(dst=ETHER_BROADCAST, src=eth_mac_src)
            pkt /= IPv6(src=self._ipv6_src,
                        dst=self.DEVICE_MULTICAST)
            pkt /= ICMPv6ND_RA()
            pkt /= ICMPv6NDOptSrcLLAddr(lladdr=self._ll_addr)
            pkt /= ICMPv6NDOptPrefixInfo(prefixlen=64, prefix=prefix)
            res = None

            self.logger.error(pkt.show(dump=True))

            output_port = self.gtp_tunnel

            sendp(pkt,
                  type=ETH_P_ALL,
                  iface=output_port,
                  verbose=0,
                  nofilter=1,
                  promisc=0)
            '''   
            res = srp1(pkt,
                       type=ETH_P_ALL,
                       iface='testing_br',
                       timeout=1,
                       verbose=0,
                       nofilter=1,
                       promisc=0)
            '''

            if res is not None:
                self.logger.debug("ARP Res pkt %s", res.show(dump=True))
                if str(res[ARP].psrc) != str(gw_ip):
                    self.logger.warning("Unexpected ARP response. %s", res.show(dump=True))
                    return ""
                else:
                    mac = res[ARP].hwsrc
                return mac
            else:
                self.logger.debug("Got Null response")
                return ""

        except Scapy_Exception as ex:
            self.logger.warning("Error in probing Mac address: err %s", ex)
            return ""
        except ValueError:
            self.logger.warning("Invalid GW Ip address: [%s]", ip)
            return ""

    def _send_neighbor_advertisement(self, output_port):
            port = OFPP_LOCAL

            pkt = Ether()#(dst=ETHER_BROADCAST, src=eth_mac_src)
            pkt /= IPv6(src=self._ipv6_src,
                        dst=self.DEVICE_MULTICAST)
            pkt /= ICMPv6ND_NA()
            pkt /= ICMPv6NDOptSrcLLAddr(lladdr=self._ll_addr)

            self.logger.error(pkt.show(dump=True))

            output_port = self.gtp_tunnel

            sendp(pkt,
                  type=ETH_P_ALL,
                  iface=output_port,
                  verbose=0,
                  nofilter=1,
                  promisc=0)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _parse_rs(self, ev):
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
            return

        if 'reg6' in ev.msg.match:
            tunnel_id = ev.msg.match['reg6'] #Replace with tunnel reg
            self.logger.error(tunnel_id)

        pkt = packet.Packet(msg.data)
        for p in pkt.protocols:
            self.logger.error(p)

        ipv6_header = pkt.get_protocols(ipv6.ipv6)[0]

        icmpv6_header = pkt.get_protocols(icmpv6.icmpv6)[0]

        prefix = self.get_custom_prefix(ipv6_header)
        if icmpv6_header.type_ == self.ICMPV6_RS_TYPE:
            self.logger.error("Recieved router soli MSG---------------")
            self._send_router_advertisement(prefix, self.gtp_tunnel)
        elif icmpv6_header.type_ == self.ICMPV6_NS_TYPE:
            self.logger.error("Recieved neighbor soli MSG---------------")
            self._send_neighbor_advertisement(self.gtp_tunnel)

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
        ip_block = ipaddress.ip_address(ipv6_header.src)
        self.logger.error(int(ip_block) & 0xffffffffffffffff)
        ip_block2 = ipaddress.ip_address(int(ip_block) & 0xffffffffffffffff)
        self.logger.error(ip_block2)

        return str(ip_block2)