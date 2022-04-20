"""
All rights reserved.
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
from collections import namedtuple
from enum import Enum

import netifaces
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from ryu.lib.packet import ether_types
from ryu.lib.packet.in_proto import IPPROTO_ICMP, IPPROTO_ICMPV6


class _IpVersion(Enum):
    IPV4 = 1
    IPV6 = 2


class AccessControlController(MagmaController):
    """
    Access control controller.

    The Access control controller is responsible for enforcing the ip blocklist,
    dropping any packets to any ipv4 addresses in the blocklist as well as
    enforcing a gre tunnel filter and dropping all packets that are not from
    allowed tunnels.
    """

    APP_NAME = "access_control"
    APP_TYPE = ControllerType.PHYSICAL
    CONFIG_INBOUND_DIRECTION = 'inbound'
    CONFIG_OUTBOUND_DIRECTION = 'outbound'

    AccessControlConfig = namedtuple(
        'AccessControlConfig',
        [
            'setup_type', 'ip_blocklist', 'allowed_gre_peers',
            'block_agw_local_ips', 'mtr_interface',
        ],
    )

    def __init__(self, *args, **kwargs):
        super(AccessControlController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME,
        )
        self.config = self._get_config(kwargs['config'], kwargs['mconfig'])
        self._tunnel_acl_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]

    def _get_config(self, config_dict, mconfig):
        block_agw_local_ips = config_dict['access_control'].get('block_agw_local_ips', True)
        mtr_interface = config_dict.get('mtr_interface', None)
        return self.AccessControlConfig(
            setup_type=config_dict['setup_type'],
            ip_blocklist=config_dict['access_control']['ip_blocklist'],
            allowed_gre_peers=mconfig.allowed_gre_peers,
            block_agw_local_ips=block_agw_local_ips,
            mtr_interface=mtr_interface,
        )

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._install_ip_blocklist_flow(datapath)
        self._install_local_ip_blocking_flows(datapath)
        if self.config.setup_type == 'CWF':
            self._install_gre_allow_flows(datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._tunnel_acl_scratch)

    def _install_default_flows(self, datapath):
        """
        Adds default flows for access control.
            For normal(ip blocklist table):
                Forward uplink to next table
                Forward downlink to scratch table
            For scratch table:
                Drop all unmatched traffic
        """
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, MagmaMatch(direction=Direction.IN), [],
            priority=flows.MINIMUM_PRIORITY, resubmit_table=self.next_table,
        )
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, MagmaMatch(direction=Direction.OUT), [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self._tunnel_acl_scratch,
        )
        if self.config.setup_type == 'CWF':
            flows.add_drop_flow(
                datapath, self._tunnel_acl_scratch, MagmaMatch(), [],
                priority=flows.MINIMUM_PRIORITY,
            )
        else:
            # TODO add LTE WLAN peers
            flows.add_resubmit_next_service_flow(
                datapath, self._tunnel_acl_scratch, MagmaMatch(), [],
                priority=flows.MINIMUM_PRIORITY, resubmit_table=self.next_table,
            )

    def _install_gre_allow_flows(self, datapath):
        """
        Adds flows to allow specific gre tunnels.

        Args:
            datapath: ryu datapath struct
        """
        for peer in self.config.allowed_gre_peers:
            self._add_gre_tun_allow_flow(datapath, peer.ip, peer.key)

    def _add_gre_tun_allow_flow(self, datapath, gre_ip, gre_key):
        # TODO how to check if protobuf field is set(only works for msgs in pr3)
        if gre_key:
            ulink_match_gre = MagmaMatch(
                direction=Direction.OUT,
                tun_ipv4_src=gre_ip, tunnel_id=gre_key,
            )
        else:
            ulink_match_gre = MagmaMatch(
                direction=Direction.OUT,
                tun_ipv4_src=gre_ip,
            )

        flows.add_resubmit_next_service_flow(
            datapath, self._tunnel_acl_scratch,
            ulink_match_gre, [],
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self.next_table,
        )

    def _install_local_ip_blocking_flows(self, datapath):
        if self.config.setup_type != 'LTE':
            return
        if not self.config.block_agw_local_ips:
            return
        interfaces = netifaces.interfaces()
        direction = self.CONFIG_INBOUND_DIRECTION

        for ip_version in (_IpVersion.IPV4, _IpVersion.IPV6):
            local_ipnet = AccessControlController._get_loopback_addresses(ip_version)
            self._install_ip_blocking_flow(datapath, local_ipnet, direction, ip_version)

            for iface in interfaces:
                self._install_local_ip_blocking_flows_for_iface(datapath, direction, iface, ip_version, local_ipnet)

    def _install_local_ip_blocking_flows_for_iface(self, datapath, direction, iface, ip_version, local_ipnet):
        if_addrs = AccessControlController._get_interface_ip_addresses(ip_version, iface)

        for addr in if_addrs:
            current_addr = AccessControlController._get_address_from_wrapper(ip_version, addr)
            if ipaddress.ip_address(current_addr) in local_ipnet:
                continue
            self.logger.info("Add blocking rule for: %s, iface %s", current_addr, iface)

            ip_network = AccessControlController._get_network_from_address(ip_version, current_addr)
            self._install_ip_blocking_flow(datapath, ip_network, direction, ip_version)
            # Add flow to allow ICMP for monitoring flows.
            if iface == self.config.mtr_interface:
                self._install_local_icmp_flows(datapath, ip_network, ip_version)

    @staticmethod
    def _get_interface_ip_addresses(ip_version, iface):
        if ip_version == _IpVersion.IPV4:
            return netifaces.ifaddresses(iface).get(netifaces.AF_INET, [])
        return netifaces.ifaddresses(iface).get(netifaces.AF_INET6, [])

    @staticmethod
    def _get_loopback_addresses(ip_version):
        if ip_version == _IpVersion.IPV4:
            return ipaddress.ip_network('127.0.0.0/8')
        return ipaddress.ip_network('::1')

    @staticmethod
    def _get_network_from_address(ip_version, address):
        if ip_version == _IpVersion.IPV4:
            return ipaddress.IPv4Network(address)
        return ipaddress.IPv6Network(address)

    @staticmethod
    def _get_address_from_wrapper(ip_version, addr):
        if ip_version == _IpVersion.IPV4:
            return addr['addr']
        # remove suffix %interface, e.g. ::1%eth1, if it exists
        return str(addr['addr'].split('%')[0])

    def _install_ip_blocklist_flow(self, datapath):
        """
        Install flows to drop any packets with ip address blocks matching the
        blocklist.
        """
        for entry in self.config.ip_blocklist:
            ip_network = ipaddress.IPv4Network(entry['ip'])
            direction = entry.get('direction', None)
            self._install_ip_blocking_flow(datapath, ip_network, direction, _IpVersion.IPV4)

    def _install_local_icmp_flows(self, datapath, ip_network, ip_version):
        match = AccessControlController._create_magma_match_icmp_flow(
            ip_version, ip_network,
        )
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num,
            match, [],
            priority=flows.MEDIUM_PRIORITY,
            resubmit_table=self.next_table,
        )

    @staticmethod
    def _create_magma_match_icmp_flow(ip_version, ip_network):
        if ip_version == _IpVersion.IPV4:
            return MagmaMatch(
                direction=Direction.OUT,
                eth_type=ether_types.ETH_TYPE_IP,
                ipv4_dst=(
                    ip_network.network_address,
                    ip_network.netmask,
                ),
                ip_proto=IPPROTO_ICMP,
            )
        return MagmaMatch(
            direction=Direction.OUT,
            eth_type=ether_types.ETH_TYPE_IPV6,
            ipv6_dst=(
                ip_network.network_address,
                ip_network.netmask,
            ),
            ip_proto=IPPROTO_ICMPV6,
        )

    def _install_ip_blocking_flow(self, datapath, ip_network, direction, ip_version):
        """
        Install flows to drop any packets with ip address blocks matching the
        blocklist.
        """
        if direction and direction not in [
            self.CONFIG_INBOUND_DIRECTION,
            self.CONFIG_OUTBOUND_DIRECTION,
        ]:
            self.logger.error(
                'Invalid direction found in ip blocklist: %s', direction,
            )
            return
        # If no direction is specified, both outbound and inbound traffic
        # will be dropped.
        if direction is None or direction == self.CONFIG_INBOUND_DIRECTION:
            match = AccessControlController._create_magma_match_blocking_flow_out(
                ip_version, ip_network,
            )
            flows.add_drop_flow(
                datapath, self.tbl_num, match, [],
                priority=flows.DEFAULT_PRIORITY,
            )
        if direction is None or direction == self.CONFIG_OUTBOUND_DIRECTION:
            match = AccessControlController._create_magma_match_blocking_flow_in(
                ip_version, ip_network,
            )
            flows.add_drop_flow(
                datapath, self.tbl_num, match, [],
                priority=flows.DEFAULT_PRIORITY,
            )

    @staticmethod
    def _create_magma_match_blocking_flow_out(ip_version, ip_network):
        if ip_version == _IpVersion.IPV4:
            return MagmaMatch(
                direction=Direction.OUT,
                eth_type=ether_types.ETH_TYPE_IP,
                ipv4_dst=(
                    ip_network.network_address,
                    ip_network.netmask,
                ),
            )
        return MagmaMatch(
            direction=Direction.OUT,
            eth_type=ether_types.ETH_TYPE_IPV6,
            ipv6_dst=(
                ip_network.network_address,
                ip_network.netmask,
            ),
        )

    @staticmethod
    def _create_magma_match_blocking_flow_in(ip_version, ip_network):
        if ip_version == _IpVersion.IPV4:
            return MagmaMatch(
                direction=Direction.IN,
                eth_type=ether_types.ETH_TYPE_IP,
                ipv4_src=(
                    ip_network.network_address,
                    ip_network.netmask,
                ),
            )
        return MagmaMatch(
            direction=Direction.IN,
            eth_type=ether_types.ETH_TYPE_IPV6,
            ipv6_src=(
                ip_network.network_address,
                ip_network.netmask,
            ),
        )
