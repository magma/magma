"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from collections import namedtuple
import uuid

from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import ether_types

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.exceptions import MagmaOFError
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, IMSI_REG


class MeterController(MagmaController):
    """
    Metering controller.

    The metering controller installs flows that are relevant to metering in the
    Metering table. Currently, this means we track volume and packet counts per
    subscriber. We install a flow rule for every IP we see in the datapath, and
    then poll for usage. We then map from the IP back to the SID.
    """

    APP_NAME = "meter"
    APP_TYPE = ControllerType.LOGICAL
    DEFAULT_FLOW_COOKIE = 0x1
    DEFAULT_IDLE_TIMEOUT_SEC = 60

    MeterConfig = namedtuple('MeterConfig', ['enabled', 'idle_timeout'])

    def __init__(self, *args, **kwargs):
        super(MeterController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.config = self._get_config(kwargs['config'])

    def _get_config(self, config_dict):
        return self.MeterConfig(
            enabled=config_dict['meter']['enabled'],
            idle_timeout=config_dict['meter'].get(
                'idle_timeout',
                self.DEFAULT_IDLE_TIMEOUT_SEC,
            ),
        )

    def initialize_on_connect(self, datapath):
        self.delete_all_flows(datapath)
        if self.config.enabled:
            self._install_default_flows(datapath)
        else:
            self._install_forward_flow(datapath)

    def _install_forward_flow(self, datapath):
        """
        Set a simple forward flow for when metering is disabled
        """
        match = MagmaMatch()
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num, match, [],
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def _install_default_flows(self, datapath):
        """
        For every UE IP block, this adds a pair of  0-priority flow-miss rules
        for incoming and outgoing traffic which trigger PACKET IN.
        """
        ofproto = datapath.ofproto
        imsi_match = (0x1, 0x1)  # match on the last bit set
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN,
                                   imsi=imsi_match)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT,
                                    imsi=imsi_match)
        flows.add_output_flow(datapath, self.tbl_num, inbound_match, [],
                              priority=flows.MINIMUM_PRIORITY,
                              cookie=self.DEFAULT_FLOW_COOKIE,
                              output_port=ofproto.OFPP_CONTROLLER,
                              max_len=ofproto.OFPCML_NO_BUFFER)
        flows.add_output_flow(datapath, self.tbl_num, outbound_match, [],
                              priority=flows.MINIMUM_PRIORITY,
                              cookie=self.DEFAULT_FLOW_COOKIE,
                              output_port=ofproto.OFPP_CONTROLLER,
                              max_len=ofproto.OFPCML_NO_BUFFER)

    @set_ev_cls(ofp_event.EventOFPPacketIn, MAIN_DISPATCHER)
    def _install_new_ingress_egress_flows(self, ev):
        """
        For every packet not already matched by a flow rule, install a pair of
        flows to track all packets to/from the corresponding IMSI.
        """

        msg = ev.msg
        datapath = msg.datapath
        parser = datapath.ofproto_parser

        if not self._matches_table(msg):
            # Intended for other application
            return

        try:
            # no need to decode the IMSI. The OFPMatch will
            # give the already-encoded IMSI value, and we can match on that
            imsi = _get_encoded_imsi_from_packetin(msg)
        except MagmaOFError as e:
            # No packet direction, but intended for this table
            self.logger.error("Error obtaining IMSI from pkt-in: %s", e)
            return

        # Set inbound/outbound tracking flows
        flow_id_note = list(bytes(str(uuid.uuid4()), 'utf-8'))
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.IN,
                                    imsi=imsi)
        outbound_actions = [parser.NXActionNote(note=flow_id_note)]

        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.OUT,
                                   imsi=imsi)
        inbound_actions = [parser.NXActionNote(note=flow_id_note)]

        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, inbound_match, inbound_actions,
            priority=flows.DEFAULT_PRIORITY,
            idle_timeout=self.config.idle_timeout,
            resubmit_table=self.next_table)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, outbound_match, outbound_actions,
            priority=flows.DEFAULT_PRIORITY,
            idle_timeout=self.config.idle_timeout,
            resubmit_table=self.next_table)

    def _matches_table(self, msg):
        return self.tbl_num == msg.table_id


def _get_encoded_imsi_from_packetin(msg):
    """
    Retrieve encoded imsi from the Packet-In message, or raise an exception if
    it doesn't exist.
    """
    imsi = msg.match.get(IMSI_REG)
    if imsi is None:
        raise MagmaOFError('IMSI not found in OFPMatch')
    return imsi
