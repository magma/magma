"""
All rights reserved.
Copyright (c) 2019-present, Facebook, Inc.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import ipaddress

from ryu.lib.packet import ether_types

from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.app.base import MagmaController


class AccessControlController(MagmaController):
    """
    Access control controller.

    The Access control controller is responsible for enforcing the ip blacklist
    and dropping any packets to any ipv4 addresses in the blacklist.
    """

    APP_NAME = "access_control"
    CONFIG_INBOUND_DIRECTION = 'inbound'
    CONFIG_OUTBOUND_DIRECTION = 'outbound'

    def __init__(self, *args, **kwargs):
        super(AccessControlController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self._ip_blacklist = \
            kwargs['config']['access_control']['ip_blacklist']

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        self._install_default_flows(datapath)
        self._install_ip_blacklist_flow(datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def _install_default_flows(self, datapath):
        """
        Default flow is to forward to next table.
        """
        flows.add_flow(datapath, self.tbl_num, MagmaMatch(), [],
                       priority=flows.MINIMUM_PRIORITY,
                       resubmit_next_service=self.next_table)

    def _install_ip_blacklist_flow(self, datapath):
        """
        Install flows to drop any packets with ip address blocks matching the
        blacklist.
        """
        for entry in self._ip_blacklist:
            ip_network = ipaddress.IPv4Network(entry['ip'])
            direction = entry.get('direction', None)
            if direction is not None and \
                    direction not in [self.CONFIG_INBOUND_DIRECTION,
                                      self.CONFIG_OUTBOUND_DIRECTION]:
                self.logger.error(
                    'Invalid direction found in ip blacklist: %s', direction)
                continue
            # If no direction is specified, both outbound and inbound traffic
            # will be dropped.
            if direction is None or direction == self.CONFIG_INBOUND_DIRECTION:
                match = MagmaMatch(direction=Direction.OUT,
                                   eth_type=ether_types.ETH_TYPE_IP,
                                   ipv4_dst=(ip_network.network_address,
                                             ip_network.netmask))
                flows.add_flow(datapath, self.tbl_num, match, [],
                               priority=flows.DEFAULT_PRIORITY)
            if direction is None or \
                    direction == self.CONFIG_OUTBOUND_DIRECTION:
                match = MagmaMatch(direction=Direction.IN,
                                   eth_type=ether_types.ETH_TYPE_IP,
                                   ipv4_src=(ip_network.network_address,
                                             ip_network.netmask))
                flows.add_flow(datapath, self.tbl_num, match, [],
                               priority=flows.DEFAULT_PRIORITY)
