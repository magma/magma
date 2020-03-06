"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import shlex
import subprocess
from typing import NamedTuple, Dict

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from ryu.controller.controller import Datapath
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.imsi import encode_imsi
from ryu.lib.packet import ether_types


class IPFIXController(MagmaController):
    """
    IPFIXController

    The IPFIX controller installs sample flows for exporting IPFIX statistics
    to the controller. Each imsi gets a sample flow. After sampling traffic is
    forwarded to the next table.
    """

    APP_NAME = "ipfix"
    APP_TYPE = ControllerType.LOGICAL

    IPFIXConfig = NamedTuple(
        'IPFIXConfig',
        [('enabled', bool), ('collector_ip', str), ('collector_port', int),
         ('probability', int), ('collector_set_id', int),
         ('obs_domain_id', int), ('obs_point_id', int), ('gtp_port', int)],
    )

    def __init__(self, *args, **kwargs):
        super(IPFIXController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.ipfix_config = self._get_ipfix_config(kwargs['config'])
        self._bridge_name = kwargs['config']['bridge_name']
        self._datapath = None

    def _get_ipfix_config(self, config_dict: Dict) -> NamedTuple:
        if 'ipfix' not in config_dict or not config_dict['ipfix']['enabled']:
            return self.IPFIXConfig(enabled=False, probability=0,
                collector_ip='', collector_port=0, collector_set_id=0,
                obs_domain_id=0, obs_point_id=0, gtp_port=0)

        return self.IPFIXConfig(
            enabled=config_dict['ipfix']['enabled'],
            collector_ip=config_dict['ipfix']['collector_ip'],
            collector_port=config_dict['ipfix']['collector_port'],
            probability=config_dict['ipfix']['probability'],
            collector_set_id=config_dict['ipfix']['collector_set_id'],
            obs_domain_id=config_dict['ipfix']['obs_domain_id'],
            obs_point_id=config_dict['ipfix']['obs_point_id'],
            gtp_port=config_dict['ovs_gtp_port_number']
        )

    def initialize_on_connect(self, datapath: Datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self._delete_all_flows(datapath)
        self._install_default_flows(datapath)

        rm_cmd = "ovs-vsctl destroy Flow_Sample_Collector_Set {}" \
            .format(self.ipfix_config.collector_set_id)

        args = shlex.split(rm_cmd)
        ret = subprocess.call(args)
        self.logger.debug("Removed old Flow_Sample_Collector_Set ret %d", ret)

        action_str = (
            'ovs-vsctl -- --id=@{} get Bridge {} -- --id=@cs create '
            'Flow_Sample_Collector_Set id={} bridge=@{} ipfix=@i -- --id=@i '
            'create IPFIX targets=\"{}\\:{}\" obs_domain_id={} obs_point_id={}'
        ).format(
            self._bridge_name, self._bridge_name,
            self.ipfix_config.collector_set_id, self._bridge_name,
            self.ipfix_config.collector_ip, self.ipfix_config.collector_port,
            self.ipfix_config.obs_domain_id, self.ipfix_config.obs_point_id
        )
        try:
            p = subprocess.Popen(action_str, shell=True,
                                 stdout=subprocess.PIPE, stderr=subprocess.PIPE)
            output, err = p.communicate()
            err_str = err.decode('utf-8')
            self.logger.error(err_str)
        except subprocess.CalledProcessError as e:
            raise Exception('Error: {} failed with: {}'.format(action_str, e))

    def cleanup_on_disconnect(self, datapath: Datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self._delete_all_flows(datapath)

    def _delete_all_flows(self, datapath: Datapath) -> None:
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def _install_default_flows(self, datapath: Datapath) -> None:
        """
        For each direction set the default flows to just forward to next app.

        Args:
            datapath: ryu datapath struct
        """
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, inbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, outbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

    def add_ue_sample_flow(self, imsi: str, msisdn: str,
                           apn_mac_addr: str, apn_name: str) -> None:
        """
        Install a flow to sample packets for IPFIX for specific imsi

        Args:
            imsi (string): subscriber to install rule for
            msisdn (string): msisdn string
            apn_mac_addr (string): AP mac address string
            apn_name (string): AP name
        """
        imsi_hex = hex(encode_imsi(imsi))
        apn_name = apn_name.replace(' ', '_')
        action_str = (
            'ovs-ofctl add-flow {} "table={},priority={},metadata={},'
            'actions=sample(probability={},collector_set_id={},'
            'obs_domain_id={},obs_point_id={},apn_mac_addr={},msisdn={},'
            'apn_name={},sampling_port={}),resubmit(,{})"'
        ).format(
            self._bridge_name, self.tbl_num, flows.UE_FLOW_PRIORITY, imsi_hex,
            self.ipfix_config.probability, self.ipfix_config.collector_set_id,
            self.ipfix_config.obs_domain_id, self.ipfix_config.obs_point_id,
            apn_mac_addr.replace("-", ":"), msisdn, apn_name,
            self.ipfix_config.gtp_port, self.next_main_table
        )
        try:
            subprocess.Popen(action_str, shell=True).wait()
        except subprocess.CalledProcessError as e:
            raise Exception('Error: {} failed with: {}'.format(action_str, e))

    def delete_ue_sample_flow(self, imsi: str) -> None:
        """
        Delete a flow to sample packets for IPFIX for specific imsi

        Args:
            imsi (string): subscriber to install rule for
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized')
            return

        if not imsi:
            self.logger.error('No subscriber specified')
            return

        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match)

    def deactivate_rules(self, imsi: str) -> None:
        """
        Deactivate flows for a subscriber.

        Args:
            imsi (string): subscriber id
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized')
            return

        if not imsi:
            self.logger.error('No subscriber specified')
            return

        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match)
