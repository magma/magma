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

from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from ryu.lib.packet import ether_types


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
        ['setup_type', 'ip_blocklist', 'allowed_gre_peers'],
    )

    def __init__(self, *args, **kwargs):
        super(AccessControlController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.config = self._get_config(kwargs['config'], kwargs['mconfig'])
        self._tunnel_acl_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]

    def _get_config(self, config_dict, mconfig):
        return self.AccessControlConfig(
            setup_type=config_dict['setup_type'],
            ip_blocklist=config_dict['access_control']['ip_blocklist'],
            allowed_gre_peers=mconfig.allowed_gre_peers
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
            priority=flows.MINIMUM_PRIORITY, resubmit_table=self.next_table)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, MagmaMatch(direction=Direction.OUT), [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self._tunnel_acl_scratch)
        if self.config.setup_type == 'CWF':
            flows.add_drop_flow(
                datapath, self._tunnel_acl_scratch, MagmaMatch(), [],
                priority=flows.MINIMUM_PRIORITY)
        else:
            # TODO add LTE WLAN peers
            flows.add_resubmit_next_service_flow(
                datapath, self._tunnel_acl_scratch, MagmaMatch(), [],
                priority=flows.MINIMUM_PRIORITY, resubmit_table=self.next_table)

    def _install_gre_allow_flows(self, datapath):
        """
        Adds flows to allow specific gre tunnels.

        Args:
            datapath: ryu datapath struct
        """
        for peer in self.config.allowed_gre_peers:
            self._add_gre_tun_allow_flow(datapath, peer.ip, peer.key)

    def _add_gre_tun_allow_flow(self, datapath, gre_ip, gre_key):
        #TODO how to check if protobuf field is set(only works for msgs in pr3)
        if gre_key:
            ulink_match_gre = MagmaMatch(direction=Direction.OUT,
                                         tun_ipv4_src=gre_ip, tunnel_id=gre_key)
        else:
            ulink_match_gre = MagmaMatch(direction=Direction.OUT,
                                         tun_ipv4_src=gre_ip)

        flows.add_resubmit_next_service_flow(datapath, self._tunnel_acl_scratch,
                                             ulink_match_gre, [],
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=self.next_table)

    def _install_ip_blocklist_flow(self, datapath):
        """
        Install flows to drop any packets with ip address blocks matching the
        blocklist.
        """
        for entry in self.config.ip_blocklist:
            ip_network = ipaddress.IPv4Network(entry['ip'])
            direction = entry.get('direction', None)
            if direction is not None and \
                    direction not in [self.CONFIG_INBOUND_DIRECTION,
                                      self.CONFIG_OUTBOUND_DIRECTION]:
                self.logger.error(
                    'Invalid direction found in ip blocklist: %s', direction)
                continue
            # If no direction is specified, both outbound and inbound traffic
            # will be dropped.
            if direction is None or direction == self.CONFIG_INBOUND_DIRECTION:
                match = MagmaMatch(direction=Direction.OUT,
                                   eth_type=ether_types.ETH_TYPE_IP,
                                   ipv4_dst=(ip_network.network_address,
                                             ip_network.netmask))
                flows.add_drop_flow(datapath, self.tbl_num, match, [],
                                    priority=flows.DEFAULT_PRIORITY)
            if direction is None or \
                    direction == self.CONFIG_OUTBOUND_DIRECTION:
                match = MagmaMatch(direction=Direction.IN,
                                   eth_type=ether_types.ETH_TYPE_IP,
                                   ipv4_src=(ip_network.network_address,
                                             ip_network.netmask))
                flows.add_drop_flow(datapath, self.tbl_num, match, [],
                                    priority=flows.DEFAULT_PRIORITY)
