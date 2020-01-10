"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
from typing import NamedTuple, Dict

from ryu.lib.packet import ether_types
from ryu.ofproto.inet import IPPROTO_TCP
from ryu.controller.controller import Datapath

from lte.protos.pipelined_pb2 import SubscriberQuotaUpdate
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.inout import INGRESS
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction


class CheckQuotaController(MagmaController):
    """
    Quota Check Controller

    This controller recognizes special IP addr that IMSI sends a request to and
    routes that request to a flask server to check user quota.
    """

    APP_NAME = "check_quota"
    APP_TYPE = ControllerType.LOGICAL
    CheckQuotaConfig = NamedTuple(
        'CheckQuotaConfig',
        [('bridge_ip', str), ('quota_check_ip', str),
         ('has_quota_port', int), ('no_quota_port', int)],
    )

    def __init__(self, *args, **kwargs):
        super(CheckQuotaController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.next_table = \
            self._service_manager.get_table_num(INGRESS)
        self._datapath = None

    def _get_config(self, config_dict: Dict) -> NamedTuple:
        return self.CheckQuotaConfig(
            bridge_ip=config_dict['bridge_ip_address'],
            quota_check_ip=config_dict['quota_check_ip'],
            has_quota_port=config_dict['has_quota_port'],
            no_quota_port=config_dict['no_quota_port'],
        )

    def initialize_on_connect(self, datapath: Datapath):
        self._datapath = datapath
        self._delete_all_flows(datapath)
        self._install_default_flows(datapath)

    def cleanup_on_disconnect(self, datapath: Datapath):
        self._delete_all_flows(datapath)

    def update_subscriber_quota_state(self, update: SubscriberQuotaUpdate):
        imsi = update.sid.id
        if update.update_type == SubscriberQuotaUpdate.VALID_QUOTA:
            self._add_subscriber_flow(imsi, True)
        elif update.update_type == SubscriberQuotaUpdate.NO_QUOTA:
            self._add_subscriber_flow(imsi, False)
        elif update.update_type == SubscriberQuotaUpdate.TERMINATE:
            self._remove_subscriber_flow(imsi)

    def _add_subscriber_flow(self, imsi: str, has_quota: bool):
        parser = self._datapath.ofproto_parser
        if has_quota:
            tcp_dst = self.config.has_quota_port
        else:
            tcp_dst = self.config.no_quota_port
        match = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.IN,
            ipv4_dst=self.config.quota_check_ip
        )
        actions = [
            parser.OFPActionSetField(ipv4_dst=self.config.bridge_ip),
            parser.OFPActionSetField(tcp_dst=tcp_dst)
        ]
        flows.add_resubmit_current_service_flow(
            self._datapath, self.tbl_num, match, actions,
            priority=flows.UE_FLOW_PRIORITY,
            resubmit_table=self.next_main_table)

    def _remove_subscriber_flow(self, imsi: str):
        match = MagmaMatch(
            imsi=encode_imsi(imsi), eth_type=ether_types.ETH_TYPE_IP,
            ip_proto=IPPROTO_TCP, direction=Direction.IN,
            ipv4_dst=self.config.quota_check_ip
        )
        flows.delete_flow(self._datapath, self.tbl_num, match)

    def _install_default_flows(self, datapath: Datapath):
        """
        Set the default flows to just forward to next app.

        Args:
            datapath: ryu datapath struct
        """
        parser = self._datapath.ofproto_parser

        # Default flows for non matched traffic
        inbound_match = MagmaMatch(direction=Direction.IN)
        outbound_match = MagmaMatch(direction=Direction.OUT)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, inbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, outbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        # For traffic from the check quota server(ip rewrite)
        match = MagmaMatch(
            eth_type=ether_types.ETH_TYPE_IP, ip_proto=IPPROTO_TCP,
            direction=Direction.OUT, ipv4_src=self.config.bridge_ip)
        actions = [
            parser.OFPActionSetField(ipv4_src=self.config.quota_check_ip),
            parser.OFPActionSetField(tcp_src=80)
        ]
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, match, actions,
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self.next_main_table
        )

    def _delete_all_flows(self, datapath: Datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
