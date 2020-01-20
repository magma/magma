"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.pipelined.openflow import flows
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
import shlex
import subprocess
import logging

from ryu.lib.packet import ether_types

# TBD: Add more apps and make it dynamic
appMap = {"facebook": 2, "whatsapp": 3, "instagram": 4, "twitter": 5}

class DPIController(MagmaController):
    """
    DPI controller.

    The DPI controller is responsible for marking a flow with an App ID derived
    from DPI. The APP ID should be stored in register 3
    """

    APP_NAME = "dpi"
    APP_TYPE = ControllerType.LOGICAL
    UPDATE_INTERVAL = 10  # seconds

    def __init__(self, *args, **kwargs):
        super(DPIController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self._datapath = None
        self._dpi_enabled = kwargs['config']['dpi']['enabled']
        self._mon_port = kwargs['config']['dpi']['mon_port']
        self._mon_port_number = kwargs['config']['dpi']['mon_port_number']
        if self._dpi_enabled:
            self._create_monitor_port()

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._datapath = datapath

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def classify_flow(self, match, app):
        ryu_match = {'eth_type': ether_types.ETH_TYPE_IP}

        ryu_match['ipv4_dst'] = self._get_ip_tuple(match.ipv4_dst)
        ryu_match['ipv4_src'] = self._get_ip_tuple(match.ipv4_src)
        ryu_match['ip_proto'] = match.ip_proto
        if match.ip_proto == match.IPProto.IPPROTO_TCP:
            ryu_match['tcp_dst'] = match.tcp_dst
            ryu_match['tcp_src'] = match.tcp_src
        elif match.ip_proto == match.IPProto.IPPROTO_UDP:
            ryu_match['udp_dst'] = match.udp_dst
            ryu_match['udp_src'] = match.udp_src

        parser = self._datapath.ofproto_parser
        app_id = appMap.get(app, 1)  # 1 is returned for unknown apps
        actions = [parser.NXActionRegLoad2(dst='reg3', value=app_id)]

        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             MagmaMatch(**ryu_match), actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=self.next_table)
        return True

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next table.
        The policies for each subscriber would be added when the IP session is
        created, by reaching out to the controller/PCRF.

        Args:
            datapath: ryu datapath struct
        """
        parser = datapath.ofproto_parser
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)
        if self._dpi_enabled:
            # TODO: Do not directly add this action once there is a way to
            # add a flow with both resubmit to the next service and send to
            # another port.
            actions = [parser.OFPActionOutput(self._mon_port_number)]
        else:
            actions = []

        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             inbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             outbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def _create_monitor_port(self):
        add_cmd = "sudo ovs-vsctl add-port gtp_br0 mon1 -- set interface \
                {} ofport_request={} \
                type=internal".format(self._mon_port, self._mon_port_number)
        enable_cmd = "sudo ifconfig {} up".format(self._mon_port)

        args = shlex.split(add_cmd)
        ret = subprocess.call(args)
        logging.debug("created monitor port ret %d", ret)

        args = shlex.split(enable_cmd)
        ret = subprocess.call(args)
        logging.debug("enabled monitor port ret %d", ret)
