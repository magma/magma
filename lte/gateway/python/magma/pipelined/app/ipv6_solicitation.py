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
from collections import namedtuple

from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.ipv6_prefix_store import get_ipv6_interface_id
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import DIRECTION_REG, Direction
from ryu.controller import dpset, ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import ether_types, ethernet, icmpv6, in_proto, ipv6, packet
from ryu.ofproto.inet import IPPROTO_ICMPV6


class IPV6SolicitationController(MagmaController):
    """
    IPV6SolicitationController responds to ipv6 router solicitation
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
    APP_NAME = 'ipv6_solicitation'
    APP_TYPE = ControllerType.PHYSICAL

    # Inherited from app_manager.RyuApp
    _CONTEXTS = {
        'dpset': dpset.DPSet,
    }

    MAC_MULTICAST = '33:33:00:00:00:01'
    DEVICE_MULTICAST = 'ff02::1'
    LINK_LOCAL_PREFIX = 'fe80::/10'
    SOLICITED_NODE_MULTICAST = 'ff02::1:ff00:0/104'
    DEFAULT_PREFIX_LEN = 64

    IPv6RouterConfig = namedtuple(
        'IPv6RouterConfig',
        ['ipv6_src', 'll_addr', 'prefix_len'],
    )

    def __init__(self, *args, **kwargs):
        super(IPV6SolicitationController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.config = self._get_config(kwargs['config'])
        self._prefix_mapper = kwargs['interface_to_prefix_mapper']
        self._tunnel_id_mapper = kwargs['tunnel_id_mapper']
        self._datapath = None
        self.logger.info("IPv6 Router using ll_addr %s, and src ip %s",
                         self.config.ll_addr, self.config.ipv6_src)

    def _get_config(self, config_dict):
        return self.IPv6RouterConfig(
            ipv6_src=config_dict['ipv6_router_addr'],
            ll_addr=config_dict['virtual_mac'],
            prefix_len=self.DEFAULT_PREFIX_LEN)

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
        ofproto = datapath.ofproto

        match_rs = MagmaMatch(eth_type=ether_types.ETH_TYPE_IPV6,
                              ipv6_src=self.LINK_LOCAL_PREFIX,
                              ip_proto=IPPROTO_ICMPV6,
                              icmpv6_type=icmpv6.ND_ROUTER_SOLICIT,
                              direction=Direction.OUT)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=match_rs, actions=[],
                              priority=flows.DEFAULT_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              max_len=ofproto.OFPCML_NO_BUFFER)

        match_ns_ue = MagmaMatch(eth_type=ether_types.ETH_TYPE_IPV6,
                                 ipv6_src=self.LINK_LOCAL_PREFIX,
                                 ip_proto=IPPROTO_ICMPV6,
                                 icmpv6_type=icmpv6.ND_NEIGHBOR_SOLICIT,
                                 direction=Direction.OUT)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=match_ns_ue, actions=[],
                              priority=flows.DEFAULT_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              max_len=ofproto.OFPCML_NO_BUFFER)

        match_ns_sgi = MagmaMatch(eth_type=ether_types.ETH_TYPE_IPV6,
                                  ipv6_dst=self.SOLICITED_NODE_MULTICAST,
                                  ip_proto=IPPROTO_ICMPV6,
                                  icmpv6_type=icmpv6.ND_NEIGHBOR_SOLICIT,
                                  direction=Direction.IN)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=match_ns_sgi, actions=[],
                              priority=flows.DEFAULT_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              max_len=ofproto.OFPCML_NO_BUFFER)

    def _send_router_advertisement(self, ipv6_src: str, tun_id: int,
                                   tun_ipv4_dst: str, output_port):
        """
        Generates the Router Advertisement response packet
        """
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser

        if not tun_id or not tun_ipv4_dst:
            self.logger.error("Packet missing tunnel information, can't reply")
            return

        prefix = self.get_custom_prefix(ipv6_src)
        if not prefix:
            self.logger.debug("Can't reply to RS for UE ip %s", ipv6_src)
            return

        pkt = packet.Packet()
        pkt.add_protocol(
            ethernet.ethernet(
                dst=self.MAC_MULTICAST,
                src=self.config.ll_addr,
                ethertype=ether_types.ETH_TYPE_IPV6,
            )
        )
        pkt.add_protocol(
            ipv6.ipv6(
                dst=self.DEVICE_MULTICAST,
                src=self.config.ipv6_src,
                nxt=in_proto.IPPROTO_ICMPV6,
            )
        )
        pkt.add_protocol(
            icmpv6.icmpv6(
                type_=icmpv6.ND_ROUTER_ADVERT,
                data=icmpv6.nd_router_advert(
                    options=[
                        icmpv6.nd_option_sla(
                            hw_src=self.config.ll_addr,
                        ),
                        icmpv6.nd_option_pi(
                            pl=self.config.prefix_len,
                            prefix=prefix,
                        )
                    ]
                ),
            )
        )
        pkt.serialize()

        actions_out = [
            parser.NXActionSetTunnel(tun_id=tun_id),
            parser.NXActionRegLoad2(dst='tun_ipv4_dst', value=tun_ipv4_dst),
            parser.OFPActionOutput(port=output_port)]
        out = parser.OFPPacketOut(datapath=self._datapath,
                                  buffer_id=ofproto.OFP_NO_BUFFER,
                                  in_port=ofproto.OFPP_CONTROLLER,
                                  actions=actions_out,
                                  data=pkt.data)
        ret = self._datapath.send_msg(out)
        if not ret:
            self.logger.error("Datapath disconnected, couldn't send RA")

    def _send_neighbor_advertisement(self, target_ipv6: str, tun_id: int,
                                     tun_ipv4_dst: str, output_port, direction):
        """
        Generates the Neighbor Advertisement response packet
        """
        ofproto, parser = self._datapath.ofproto, self._datapath.ofproto_parser

        # Only check direction OUT because direction IN doesn't need tunn info
        if direction == Direction.OUT and (not tun_id or not tun_ipv4_dst):
            self.logger.error("Packet missing tunnel information, can't reply")
            return

        prefix = self.get_custom_prefix(target_ipv6)
        if not prefix:
            self.logger.debug("Can't reply to NS for UE ip %s", target_ipv6)
            return

        pkt = packet.Packet()
        pkt.add_protocol(
            ethernet.ethernet(
                dst=self.MAC_MULTICAST,
                src=self.config.ll_addr,
                ethertype=ether_types.ETH_TYPE_IPV6,
            )
        )
        pkt.add_protocol(
            ipv6.ipv6(
                dst=self.DEVICE_MULTICAST,
                src=self.config.ipv6_src,
                nxt=in_proto.IPPROTO_ICMPV6,
            )
        )
        pkt.add_protocol(
            icmpv6.icmpv6(
                type_=icmpv6.ND_NEIGHBOR_ADVERT,
                data=icmpv6.nd_neighbor(
                    dst=target_ipv6,
                    option=icmpv6.nd_option_tla(hw_src=self.config.ll_addr)
                ),
            )
        )
        pkt.serialize()

        # For NS from SGI response doesn't need tunnel information
        actions_out = []
        if direction == Direction.OUT:
            actions_out.extend([
                parser.NXActionSetTunnel(tun_id=tun_id),
                parser.NXActionRegLoad2(dst='tun_ipv4_dst', value=tun_ipv4_dst),
            ])
        actions_out.append(parser.OFPActionOutput(port=output_port))
        out = parser.OFPPacketOut(datapath=self._datapath,
                                  buffer_id=ofproto.OFP_NO_BUFFER,
                                  in_port=ofproto.OFPP_CONTROLLER,
                                  actions=actions_out,
                                  data=pkt.data)
        ret = self._datapath.send_msg(out)
        if not ret:
            self.logger.error("Datapath disconnected, couldn't send NA")

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _parse_pkt_in(self, ev):
        """
        Process the packet in message, reply with RA/NA packets
        """
        msg = ev.msg

        if self.tbl_num != msg.table_id:
            # Intended for other application
            return

        in_port = ev.msg.match['in_port']
        tun_id = None
        tun_ipv4_src = None
        if DIRECTION_REG not in ev.msg.match:
            self.logger.error("Packet missing direction reg, can't reply")
            return
        direction = ev.msg.match[DIRECTION_REG]
        if 'tunnel_id' in ev.msg.match:
            tun_id = ev.msg.match['tunnel_id']

        if 'tun_ipv4_src' in ev.msg.match:
            tun_ipv4_src = ev.msg.match['tun_ipv4_src']

        pkt = packet.Packet(msg.data)
        self.logger.debug("Received PKT ->")
        for p in pkt.protocols:
            self.logger.debug(p)

        ipv6_header = pkt.get_protocols(ipv6.ipv6)[0]
        icmpv6_header = pkt.get_protocols(icmpv6.icmpv6)[0]

        tun_id_dst = self._tunnel_id_mapper.get_tunnel(tun_id)

        if icmpv6_header.type_ == icmpv6.ND_ROUTER_SOLICIT:
            self.logger.debug("Received router solicitation MSG")
            self._send_router_advertisement(ipv6_header.src, tun_id_dst,
                                            tun_ipv4_src, in_port)
        elif icmpv6_header.type_ == icmpv6.ND_NEIGHBOR_SOLICIT:
            self.logger.debug("Received neighbor solicitation MSG")
            self._send_neighbor_advertisement(icmpv6_header.data.dst,
                                              tun_id_dst, tun_ipv4_src,
                                              in_port, direction)

    def handle_restart(self):
        pass

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def get_custom_prefix(self, ue_ll_ipv6: str) -> str:
        """
        Retrieve the custom prefix by extracting the interface id out of the
        packet
        """

        interface_id = get_ipv6_interface_id(ue_ll_ipv6)
        return self._prefix_mapper.get_prefix(interface_id)
