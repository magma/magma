"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import shlex
import subprocess

from magma.pipelined.openflow import flows
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, DPI_REG
from magma.pipelined.policy_converters import FlowMatchError, \
    flow_match_to_magma_match


from ryu.lib.packet import ether_types

# TBD: Add more apps and make it dynamic
appMap = {"facebook_messenger": 1, "instagram": 1, "facebook": 3, "youtube": 4,
          "gmail": 5, "google": 6, "google_docs": 7, "viber": 8, "imo": 9,
          "netflix": 10, "apple": 11, "microsoft": 12}


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
        self._bridge_name = kwargs['config']['bridge_name']
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

    def add_classify_flow(self, match, app):
        try:
            match = flow_match_to_magma_match(match)
        except FlowMatchError as e:
            self.logger.error(e)
            return False

        parser = self._datapath.ofproto_parser
        app_id = appMap.get(app, 0)
        if app_id != 0:
            actions = [parser.NXActionRegLoad2(dst=DPI_REG, value=app_id)]
            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                                 match, actions,
                                                 priority=flows.DEFAULT_PRIORITY,
                                                 resubmit_table=self.next_table)
        else:
            self.logger.error("Unrecognized app name %s", app)

        return True

    def remove_classify_flow(self, match, app):
        try:
            match = flow_match_to_magma_match(match)
        except FlowMatchError as e:
            self.logger.error(e)
            return False

        flows.delete_flow(self._datapath, self.tbl_num, match)
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
        add_cmd = "ovs-vsctl add-port {} mon1 -- set interface {} \
            ofport_request={} type=internal" \
            .format(self._bridge_name, self._mon_port, self._mon_port_number)

        args = shlex.split(add_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Created monitor port ret %d", ret)

        enable_cmd = "ifconfig {} up".format(self._mon_port)
        args = shlex.split(enable_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Enabled monitor port ret %d", ret)
