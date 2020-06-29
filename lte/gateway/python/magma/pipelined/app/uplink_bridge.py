"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import subprocess
from collections import namedtuple

from magma.pipelined.app.base import MagmaController, ControllerType


class UplinkBridgeController(MagmaController):
    """
    This controller manages uplink bridge flows
    These flows are used in Non NAT configuration.
    """

    APP_NAME = "uplink_bridge"
    APP_TYPE = ControllerType.SPECIAL
    UPLINK_DHCP_PORT_NAME = 'dhcp0'
    UPLINK_PATCH_PORT_NAME = 'patch-agw'
    UPLINK_OVS_BRIDGE_NAME = 'uplink_br0'
    DEFAULT_UPLINK_PORT_NANE = 'eth3'
    DEFAULT_UPLINK_MAC = '11:22:33:44:55:66'

    UplinkConfig = namedtuple(
        'UplinkBridgeConfig',
        ['uplink_bridge', 'uplink_eth_port_name', 'uplink_patch',
         'non_nat', 'virtual_mac', 'dhcp_port'],
    )

    def __init__(self, *args, **kwargs):
        super(UplinkBridgeController, self).__init__(*args, **kwargs)

        self.config = self._get_config(kwargs['config'])
        self.logger.info("uplink bridge app config: %s", self.config)

    def _get_config(self, config_dict) -> namedtuple:

        non_nat = config_dict.get('non_nat', False)
        bridge_name = config_dict.get('uplink_bridge',
                                      self.UPLINK_OVS_BRIDGE_NAME)
        dhcp_port = config_dict.get('uplink_dhcp_port',
                                    self.UPLINK_DHCP_PORT_NAME)
        uplink_patch = config_dict.get('uplink_patch',
                                       self.UPLINK_PATCH_PORT_NAME)

        uplink_eth_port_name = config_dict.get('uplink_eth_port_name',
                                               self.DEFAULT_UPLINK_PORT_NANE)
        virtual_mac = config_dict.get('virtual_mac',
                                      self.DEFAULT_UPLINK_MAC)

        return self.UplinkConfig(
            non_nat=non_nat,
            uplink_bridge=bridge_name,
            uplink_eth_port_name=uplink_eth_port_name,
            virtual_mac=virtual_mac,
            uplink_patch=uplink_patch,
            dhcp_port=dhcp_port,
        )

    def initialize_on_connect(self, datapath):
        if self.config.non_nat is False:
            self._delete_all_flows()
            return

        self._delete_all_flows()
        self._add_eth_port()
        # flows to forward traffic between patch port to eth port

        # 1. DHCP traffic
        match = "in_port=%s,ip,udp,tp_dst=68" % self.config.uplink_eth_port_name
        actions = "output:%s,output:%s" % (self.config.dhcp_port,
                                           self.config.uplink_patch)
        self._install_flow(2000, match, actions)

        # 2.a. all egress traffic
        match = "in_port=%s,ip" % self.config.uplink_patch
        actions = "mod_dl_src=%s, output:%s" % (self.config.virtual_mac,
                                                self.config.uplink_eth_port_name)
        self._install_flow(1000, match, actions)

        # 2.b. All ingress IP traffic for UE mac
        match = "in_port=%s,ip, dl_dst=%s" % (self.config.uplink_eth_port_name,
                                              self.config.virtual_mac)
        actions = "output:%s" % self.config.uplink_patch
        self._install_flow(1000, match, actions)

        # everything else:
        self._install_flow(100, "", "NORMAL")

    def cleanup_on_disconnect(self, datapath):
        self._del_eth_port()
        self._delete_all_flows()

    def delete_all_flows(self, datapath):
        self._delete_all_flows()

    def _delete_all_flows(self):
        if self.config.uplink_bridge is None:
            return
        del_flows = "ovs-ofctl del-flows %s" % self.config.uplink_bridge
        self.logger.info("Delete all flows: %s", del_flows)
        try:
            subprocess.Popen(del_flows, shell=True).wait()
        except subprocess.CalledProcessError as ex:
            raise Exception('Error: %s failed with: %s' % (del_flows, ex))

    def _install_flow(self, priority: int, flow_match: str, flow_action: str):
        if self.config.non_nat is False:
            return
        flow_cmd = "ovs-ofctl add-flow %s \"priority=%s,%s, actions=%s\"" % (
            self.config.uplink_bridge, priority,
            flow_match, flow_action)

        self.logger.info("Create flow %s", flow_cmd)

        try:
            subprocess.Popen(flow_cmd, shell=True).wait()
        except subprocess.CalledProcessError as ex:
            raise Exception('Error: %s failed with: %s' % (flow_cmd, ex))

    def _add_eth_port(self):
        if self.config.non_nat is False or \
                self.config.uplink_eth_port_name is None:
            return

        ovs_add_port = "ovs-vsctl --may-exist add-port %s %s" \
                       % (self.config.uplink_bridge, self.config.uplink_eth_port_name)
        self.logger.info("Add uplink port: %s", ovs_add_port)
        try:
            subprocess.Popen(ovs_add_port, shell=True).wait()
        except subprocess.CalledProcessError as ex:
            raise Exception('Error: %s failed with: %s' % (ovs_add_port, ex))

    def _del_eth_port(self):
        if self.config.non_nat is False or \
                self.config.uplink_eth_port_name is None:
            return

        ovs_rem_port = "ovs-vsctl --if-exists del-port %s %s" \
                       % (self.config.uplink_bridge, self.config.uplink_eth_port_name)
        self.logger.info("Remove ovs uplink port: %s", ovs_rem_port)
        try:
            subprocess.Popen(ovs_rem_port, shell=True).wait()
        except subprocess.CalledProcessError as ex:
            raise Exception('Error: %s failed with: %s' % (ovs_rem_port, ex))
