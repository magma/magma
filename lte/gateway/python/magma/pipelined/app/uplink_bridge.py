"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
import os
import subprocess
from collections import namedtuple

import netaddr
import netifaces
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow import flows
from ryu.lib import hub

UPLINK_OVS_BRIDGE_NAME = 'uplink_br0'


class UplinkBridgeController(MagmaController):
    """
    This controller manages uplink bridge flows
    These flows are used in Non NAT configuration.
    """

    APP_NAME = "uplink_bridge"
    APP_TYPE = ControllerType.SPECIAL
    UPLINK_DHCP_PORT_NAME = 'dhcp0'
    UPLINK_PATCH_PORT_NAME = 'patch-agw'
    DEFAULT_UPLINK_PORT_NANE = 'eth3'
    DEFAULT_UPLINK_MAC = '11:22:33:44:55:66'
    DEFAULT_DEV_VLAN_IN = 'vlan_pop_in'
    DEFAULT_DEV_VLAN_OUT = 'vlan_pop_out'
    SGI_INGRESS_FLOW_UPDATE_FREQ = 60

    UplinkConfig = namedtuple(
        'UplinkBridgeConfig',
        ['uplink_bridge', 'uplink_eth_port_name', 'uplink_patch',
         'enable_nat', 'virtual_mac', 'dhcp_port',
         'sgi_management_iface_vlan', 'sgi_management_iface_ip_addr',
         'dev_vlan_in', 'dev_vlan_out', 'ovs_vlan_workaround',
         'sgi_management_iface_gw'],
    )

    def __init__(self, *args, **kwargs):
        super(UplinkBridgeController, self).__init__(*args, **kwargs)

        self.config = self._get_config(kwargs['config'])
        self.logger.info("uplink bridge app config: %s", self.config)

    def _get_config(self, config_dict) -> namedtuple:

        enable_nat = config_dict.get('enable_nat', True)
        bridge_name = config_dict.get('uplink_bridge', UPLINK_OVS_BRIDGE_NAME)
        dhcp_port = config_dict.get('uplink_dhcp_port',
                                    self.UPLINK_DHCP_PORT_NAME)
        uplink_patch = config_dict.get('uplink_patch',
                                       self.UPLINK_PATCH_PORT_NAME)

        uplink_eth_port_name = config_dict.get('uplink_eth_port_name',
                                               self.DEFAULT_UPLINK_PORT_NANE)
        if uplink_eth_port_name not in netifaces.interfaces():
            uplink_eth_port_name = None

        virtual_mac = config_dict.get('virtual_mac',
                                      self.DEFAULT_UPLINK_MAC)
        sgi_management_iface_vlan = config_dict.get('sgi_management_iface_vlan', "")
        sgi_management_iface_ip_addr = config_dict.get('sgi_management_iface_ip_addr', "")
        dev_vlan_in = config_dict.get('dev_vlan_in', self.DEFAULT_DEV_VLAN_IN)
        dev_vlan_out = config_dict.get('dev_vlan_out', self.DEFAULT_DEV_VLAN_OUT)
        ovs_vlan_workaround = config_dict.get('ovs_vlan_workaround', True)
        sgi_management_iface_gw = config_dict.get('sgi_management_iface_gw', "")
        return self.UplinkConfig(
            enable_nat=enable_nat,
            uplink_bridge=bridge_name,
            uplink_eth_port_name=uplink_eth_port_name,
            virtual_mac=virtual_mac,
            uplink_patch=uplink_patch,
            dhcp_port=dhcp_port,
            sgi_management_iface_vlan=sgi_management_iface_vlan,
            sgi_management_iface_ip_addr=sgi_management_iface_ip_addr,
            dev_vlan_in=dev_vlan_in,
            dev_vlan_out=dev_vlan_out,
            ovs_vlan_workaround=ovs_vlan_workaround,
            sgi_management_iface_gw=sgi_management_iface_gw
        )

    def initialize_on_connect(self, datapath):
        if self.config.enable_nat is True:
            self._delete_all_flows()
            self._del_eth_port()
            return

        self._delete_all_flows()
        self._add_eth_port()
        self._setup_vlan_pop_dev()

        # flows to forward traffic between patch port to eth port
        # 1. Setup SGi management iface flows
        if self.config.sgi_management_iface_vlan:
            # 1.a. Ingress
            match = "in_port=%s,vlan_vid=%s/0x1fff" % (self.config.uplink_eth_port_name,
                                                       hex(0x1000 | int(self.config.sgi_management_iface_vlan)))
            actions = "strip_vlan,output:LOCAL"
            self._install_flow(flows.MAXIMUM_PRIORITY, match, actions)

            # 1.b. Egress
            match = "in_port=LOCAL"
            actions = "push_vlan:0x8100,mod_vlan_vid=%s,output:%s" % (self.config.sgi_management_iface_vlan,
                                                                      self.config.uplink_eth_port_name)
            self._install_flow(flows.MAXIMUM_PRIORITY, match, actions)
        else:
            # 1.a. Egress
            match = "in_port=LOCAL"
            actions = "output:%s" % self.config.uplink_eth_port_name
            self._install_flow(flows.MINIMUM_PRIORITY, match, actions)

        # 2.a Ingress: DHCP reply flows
        match = "in_port=%s,ip,udp,tp_dst=68" % self.config.uplink_eth_port_name
        actions = "output:%s,output:LOCAL" % self.config.dhcp_port
        self._install_flow(flows.MAXIMUM_PRIORITY - 1, match, actions)
        # 2.b. Egress: DHCP Req traffic
        match = "in_port=%s" % self.config.dhcp_port
        actions = "output:%s" % self.config.uplink_eth_port_name
        self._install_flow(flows.MAXIMUM_PRIORITY - 1, match, actions)

        # 3. UE egress traffic
        match = "in_port=%s" % self.config.uplink_patch
        actions = "mod_dl_src=%s, output:%s" % (self.config.virtual_mac,
                                                self.config.uplink_eth_port_name)
        self._install_flow(flows.MEDIUM_PRIORITY, match, actions)

        # 4. Remaining Ingress traffic

        if self.config.ovs_vlan_workaround:
            # 4.a. All ingress IP traffic for UE mac
            match = "in_port=%s,dl_dst=%s, vlan_tci=0x0000/0x1000" % \
                    (self.config.uplink_eth_port_name,
                     self.config.virtual_mac)
            actions = "output:%s" % self.config.uplink_patch
            self._install_ip_v4_v6_flows(flows.MEDIUM_PRIORITY, match, actions)

            match = "in_port=%s,dl_dst=%s, vlan_tci=0x1000/0x1000" % \
                    (self.config.uplink_eth_port_name,
                     self.config.virtual_mac)
            actions = "strip_vlan,output:%s" % self.config.dev_vlan_in
            self._install_ip_v4_v6_flows(flows.MEDIUM_PRIORITY, match, actions)

            # 4.b. redirect all vlan-out traffic to patch port
            match = "in_port=%s,dl_dst=%s" % \
                    (self.config.dev_vlan_out,
                     self.config.virtual_mac)
            actions = "output:%s" % self.config.uplink_patch
            self._install_ip_v4_v6_flows(flows.MEDIUM_PRIORITY, match, actions)
        else:
            # 4.a. All ingress IP traffic for UE mac
            match = "in_port=%s, dl_dst=%s" % \
                    (self.config.uplink_eth_port_name,
                     self.config.virtual_mac)
            actions = "output:%s" % self.config.uplink_patch
            self._install_ip_v4_v6_flows(flows.MEDIUM_PRIORITY, match, actions)

        # 5. Handle ARP from eth0
        match = "in_port=%s,arp" % self.config.uplink_eth_port_name
        actions = "output:%s,output:%s,output:LOCAL" % (self.config.dhcp_port,
                                                        self.config.uplink_patch)
        self._install_flow(flows.MINIMUM_PRIORITY, match, actions)

        # config interfaces:
        self._kill_dhclient(self.config.uplink_eth_port_name)
        self._flush_ip(self.config.uplink_eth_port_name)

        self._set_sgi_ip_addr(self.config.uplink_bridge)
        self._set_sgi_gw(self.config.uplink_bridge)
        self._set_arp_ignore('all', '1')

        # 6. After setting IP, setup SGi interface Ingress flow
        self._set_sgi_interface_ingress_flows()

    def cleanup_on_disconnect(self, datapath):
        self._del_eth_port()
        self._delete_all_flows()

    def delete_all_flows(self, datapath):
        self._delete_all_flows()

    def _set_sgi_interface_ingress_flows(self):
        if_addrs = netifaces.ifaddresses(self.config.uplink_bridge).get(netifaces.AF_INET, [])
        for addr in if_addrs:
            addr = addr['addr']

            match = "in_port=%s,ip,ip_dst=%s" % (self.config.uplink_eth_port_name,
                                                 addr)
            actions = "output:LOCAL"
            self._install_flow(flows.MEDIUM_PRIORITY + 1, match, actions)

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
        if self.config.enable_nat is True:
            return
        flow_cmd = "ovs-ofctl add-flow -Oopenflow13 %s \"priority=%s,%s, actions=%s\"" % (
            self.config.uplink_bridge, priority,
            flow_match, flow_action)
        self.logger.info("Create flow %s", flow_cmd)

        try:
            subprocess.Popen(flow_cmd, shell=True).wait()
        except subprocess.CalledProcessError as ex:
            raise Exception('Error: %s failed with: %s' % (flow_cmd, ex))

    def _install_ip_v4_v6_flows(self, priority: int, flow_match: str, flow_action: str):
        if self.config.enable_nat is True:
            return

        self._install_flow(priority, flow_match + ", ip", flow_action)
        self._install_flow(priority, flow_match + ", ipv6", flow_action)

    def _add_eth_port(self):
        if self.config.enable_nat is True or \
                self.config.uplink_eth_port_name is None:
            return
        if BridgeTools.port_is_in_bridge(self.config.uplink_bridge,
                                         self.config.uplink_eth_port_name):
            return
        self._cleanup_if(self.config.uplink_eth_port_name, True)
        # Add eth interface to OVS.
        ovs_add_port = "ovs-vsctl --may-exist add-port %s %s" \
                       % (self.config.uplink_bridge, self.config.uplink_eth_port_name)
        try:
            subprocess.Popen(ovs_add_port, shell=True).wait()
        except subprocess.CalledProcessError as ex:
            raise Exception('Error: %s failed with: %s' % (ovs_add_port, ex))

        self.logger.info("Add uplink port: %s", ovs_add_port)
        # sometimes the mac address changes after port addition, so restart the service.
        self.logger.info('OVS uplink bridge reconfigured, restarting to get new config')
        os._exit(0) # pylint: disable=protected-access

    def _del_eth_port(self):
        if BridgeTools.port_is_in_bridge(self.config.uplink_bridge,
                                             self.config.uplink_eth_port_name):
            self._cleanup_if(self.config.uplink_bridge, True)
            if self.config.uplink_eth_port_name is None:
                return

            ovs_rem_port = "ovs-vsctl --if-exists del-port %s %s" \
                           % (self.config.uplink_bridge, self.config.uplink_eth_port_name)
            try:
                subprocess.Popen(ovs_rem_port, shell=True).wait()
                self.logger.info("Remove ovs uplink port: %s", ovs_rem_port)
            except subprocess.CalledProcessError as ex:
                self.logger.debug("ignore port del error: %s ", ex)
                return

        if self.config.uplink_eth_port_name:
            self._set_sgi_ip_addr(self.config.uplink_eth_port_name)
            self._set_sgi_gw(self.config.uplink_eth_port_name)

    def _set_sgi_gw(self, if_name: str):
        self.logger.debug('self.config.sgi_management_iface_gw %s',
                          self.config.sgi_management_iface_gw)

        if self.config.sgi_management_iface_gw is None or \
                self.config.sgi_management_iface_gw == "":
            return

        try:
            set_gw_command = ["ip",
                              "route", "replace", "default", "via",
                              self.config.sgi_management_iface_gw,
                              "metric", "100", "dev",
                              if_name]
            subprocess.check_call(set_gw_command)
            self.logger.debug("SGi GW config: [%s]", set_gw_command)
        except subprocess.SubprocessError as e:
            self.logger.warning("Error while setting SGi GW: %s", e)

    def _set_sgi_ip_addr(self, if_name: str):
        self.logger.debug("self.config.sgi_management_iface_ip_addr %s",
                          self.config.sgi_management_iface_ip_addr)
        if self.config.sgi_management_iface_ip_addr is None or \
                self.config.sgi_management_iface_ip_addr == "":
            if if_name == self.config.uplink_bridge:
                self._restart_dhclient(if_name)

            else:
                if_addrs = netifaces.ifaddresses(if_name).get(netifaces.AF_INET, [])
                if len(if_addrs) != 0:
                    self.logger.info("SGi has valid IP, skip reconfiguration %s", if_addrs)
                    return

                # for system port, use networking config
                try:
                    self._flush_ip(if_name)
                except subprocess.CalledProcessError as ex:
                    self.logger.info("could not flush ip addr: %s, %s",
                                     if_name, ex)

                if_up_cmd = ["ifup", if_name, "--force"]
                try:
                    subprocess.check_call(if_up_cmd)
                except subprocess.CalledProcessError as ex:
                    self.logger.info("could not bring up if: %s, %s",
                                     if_up_cmd, ex)
            return

        try:
            self._kill_dhclient(if_name)

            if self._is_iface_ip_set(if_name,
                                     self.config.sgi_management_iface_ip_addr):
                self.logger.info("ip addr %s already set for iface %s",
                                 self.config.sgi_management_iface_ip_addr,
                                 if_name)
                return

            self._flush_ip(if_name)

            set_ip_cmd = ["ip",
                          "addr", "add",
                          self.config.sgi_management_iface_ip_addr,
                          "dev",
                          if_name]
            subprocess.check_call(set_ip_cmd)
            self.logger.debug("SGi ip address config: [%s]", set_ip_cmd)
        except subprocess.SubprocessError as e:
            self.logger.warning("Error while setting SGi IP: %s", e)

    def _is_iface_ip_set(self, if_name, ip_addr):
        ip_addr = netaddr.IPNetwork(ip_addr)
        if_addrs = netifaces.ifaddresses(if_name).get(netifaces.AF_INET, [])

        for addr in if_addrs:
            addr = netaddr.IPNetwork("/".join((addr['addr'], addr['netmask'])))
            if ip_addr == addr:
                return True
        return False

    def _flush_ip(self, if_name):
        flush_ip = ["ip", "addr", "flush", "dev", if_name]
        subprocess.check_call(flush_ip)

    def _kill_dhclient(self, if_name):
        # Kill dhclient if running.
        pgrep_out = subprocess.Popen(["pgrep", "-f", "dhclient.*" + if_name],
                                     stdout=subprocess.PIPE)
        for pid in pgrep_out.stdout.readlines():
            subprocess.check_call(["kill", pid.strip()])

    def _restart_dhclient(self, if_name):
        # restart DHCP client can take loooong time, process it in separate thread:
        hub.spawn(self._restart_dhclient_if(if_name))

    def _setup_vlan_pop_dev(self):
        if self.config.ovs_vlan_workaround:
            # Create device
            BridgeTools.create_veth_pair(self.config.dev_vlan_in,
                                         self.config.dev_vlan_out)
            # Add to OVS,
            # OFP requested port (70 and 71) no are for test validation,
            # its not used anywhere else.
            BridgeTools.add_ovs_port(self.config.uplink_bridge,
                                     self.config.dev_vlan_in, "70")
            BridgeTools.add_ovs_port(self.config.uplink_bridge,
                                     self.config.dev_vlan_out, "71")

    def _cleanup_if(self, if_name, flush: bool):
        # Release eth IP first.
        release_eth_ip = ["dhclient", "-r", if_name]
        try:
            subprocess.check_call(release_eth_ip)
        except subprocess.CalledProcessError as ex:
            self.logger.info("could not release dhcp lease: %s, %s",
                             release_eth_ip, ex)

        if not flush:
            return
        try:
            self._flush_ip(if_name)
        except subprocess.CalledProcessError as ex:
            self.logger.info("could not flush ip addr: %s, %s", if_name, ex)

        self.logger.info("SGi DHCP: port [%s] ip removed", if_name)


    def _restart_dhclient_if(self, if_name):
        self._cleanup_if(if_name, False)

        setup_dhclient = ["dhclient", if_name]
        try:
            subprocess.check_call(setup_dhclient)
        except subprocess.CalledProcessError as ex:
            self.logger.info("could not release dhcp lease: %s, %s",
                             setup_dhclient, ex)

        self.logger.info("SGi DHCP: restart for %s done", if_name)
        while True:
            # keep updating flow to handle IP address change.
            self._set_sgi_interface_ingress_flows()
            hub.sleep(self.SGI_INGRESS_FLOW_UPDATE_FREQ)

    def _set_arp_ignore(self, if_name: str, val: str):
        sysctl_setting = 'net.ipv4.conf.' + if_name + '.arp_ignore=' + val
        subprocess.check_call(['sysctl', sysctl_setting])
