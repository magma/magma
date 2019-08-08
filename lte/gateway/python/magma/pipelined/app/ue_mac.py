"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""

from collections import namedtuple

from .base import MagmaController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import IMSI_REG


class UEMacAddressController(MagmaController):
    """
    UE MAC Address Controller

    This controller controls table 0 which is the first table every packet
    touches. It matches on UE MAC address and sets IMSI metadata
    """

    APP_NAME = "ue_mac"
    UEMacConfig = namedtuple(
        'UEMacConfig',
        ['gre_tunnel_port'],
    )

    def __init__(self, *args, **kwargs):
        super(UEMacAddressController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self._datapath = None

    def _get_config(self, config_dict):
        return self.UEMacConfig(
            # TODO: rename port number to a tunneling protocol agnostic name
            gre_tunnel_port=config_dict['ovs_gtp_port_number'],
        )

    def initialize_on_connect(self, datapath):
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              self.APP_NAME))
        self._datapath = datapath

    def cleanup_on_disconnect(self, datapath):
        flows.delete_all_flows_from_table(datapath,
                                          self._service_manager.get_table_num(
                                              self.APP_NAME))

    def add_ue_mac_flow(self, sid, mac_addr):
        uplink_match = MagmaMatch(in_port=self.config.gre_tunnel_port,
                                  eth_src=mac_addr)
        self._add_resubmit_flow(sid, uplink_match)

        downlink_match = MagmaMatch(in_port=self._datapath.ofproto.OFPP_LOCAL,
                                    eth_dst=mac_addr)
        self._add_resubmit_flow(sid, downlink_match)

    def delete_ue_mac_flow(self, sid, mac_addr):
        uplink_match = MagmaMatch(in_port=self.config.gre_tunnel_port,
                                  eth_src=mac_addr)
        self._delete_resubmit_flow(sid, uplink_match)

        downlink_match = MagmaMatch(in_port=self._datapath.ofproto.OFPP_LOCAL,
                                    eth_dst=mac_addr)
        self._delete_resubmit_flow(sid, downlink_match)

    def _add_resubmit_flow(self, sid, match):
        parser = self._datapath.ofproto_parser
        tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        next_table = self._service_manager.get_next_table_num(self.APP_NAME)

        # Add IMSI metadata
        actions = [
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid))]

        flows.add_resubmit_next_service_flow(self._datapath, tbl_num, match,
                                             actions=actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=next_table)

    def _delete_resubmit_flow(self, sid, match):
        parser = self._datapath.ofproto_parser
        tbl_num = self._service_manager.get_table_num(self.APP_NAME)

        # Add IMSI metadata
        actions = [
            parser.NXActionRegLoad2(dst=IMSI_REG, value=encode_imsi(sid))]

        flows.delete_flow(self._datapath, tbl_num, match, actions=actions)
