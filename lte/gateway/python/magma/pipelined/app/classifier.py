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
import subprocess
import ipaddress
import logging
from collections import namedtuple
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from .base import MagmaController
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.app.inout import INGRESS
from ryu.lib.packet import ether_types
from ryu.lib.ovs import bridge
from ryu import cfg
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.utils import Utils
from magma.pipelined.openflow.registers import TUN_PORT_REG

GTP_PORT_MAC = "02:00:00:00:00:01"
TUNNEL_OAM_FLAG = 1
OVSDB_PORT = 6640  # The IANA registered port for OVSDB [RFC7047]
OVSDB_MANAGER_ADDR = 'ptcp:6640'

class Classifier(MagmaController):
    """
    A controller that sets up an openflow pipeline for Magma.
    This controller is used for install the default tunnel entry,
    add tunnel flows for gtp_port/local_port/mtr_port and delete
    tunnel flows for gtp_port/local_port/mtr_port into OVS table 0.
    """
    APP_NAME = "classifier"
    APP_TYPE = ControllerType.SPECIAL
    ClassifierConfig = namedtuple(
            'ClassifierConfig',
            ['gtp_port', 'mtr_ip', 'mtr_port', 'internal_sampling_port', 'internal_sampling_fwd_tbl', 'multi_tunnel_flag'],
    )

    def __init__(self, *args, **kwargs):
        super(Classifier, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_table_num(INGRESS)
        self._uplink_port = OFPP_LOCAL
        self._datapath = None
        self._clean_restart = kwargs['config']['clean_restart']
        if self.config.multi_tunnel_flag:
            self._ovs_multi_tunnel_init()
        self.CONF = cfg.CONF
        # OVSBridge instance instantiated later
        self.ovs = None

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None
        
        if 'mtr_ip' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']

        if 'ovs_multi_tunnel' in config_dict:
            multi_tunnel_flag = config_dict['ovs_multi_tunnel']
        else:
            multi_tunnel_flag = False

        return self.ClassifierConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            internal_sampling_port=
                               config_dict['ovs_internal_sampling_port_number'],
            internal_sampling_fwd_tbl=
                              config_dict['ovs_internal_sampling_fwd_tbl_number'],
            multi_tunnel_flag=multi_tunnel_flag,

        )

    def _ovs_multi_tunnel_init(self):
        try:
            subprocess.check_output("sudo ovs-vsctl list Open_vSwitch | grep gtpu",
                                     shell=True,)
            self.ovs_gtp_type = "gtpu"
        except subprocess.CalledProcessError:
            self.ovs_gtp_type = "gtp"

        try:
            subprocess.check_output('ovs-vsctl set-manager %s' % OVSDB_MANAGER_ADDR,
                                     shell=True,)
        except subprocess.CalledProcessError as e:
            logging.warning("Error for set manager with ovsdb: %s", e)

    def _ip_addr_to_gtp_port_name(self, enodeb_ip_addr:str):
        ip_no = hex(int(ipaddress.ip_address(enodeb_ip_addr)))
        buf = "g_{}".format(ip_no[2:])
        return buf

    def _get_ovs_bridge(self, dpid):
        ovsdb_addr = 'tcp:%s:%d' % (self._datapath.address[0], OVSDB_PORT)
        if (self.ovs is not None
                and self.ovs.datapath_id == dpid
                and self.ovs.vsctl.remote == ovsdb_addr):
            return self.ovs

        try:
            self.ovs = bridge.OVSBridge(
                CONF=self.CONF,
                datapath_id=self._datapath.id,
                ovsdb_addr=ovsdb_addr)
            self.ovs.init()
        except (ValueError, KeyError) as e:
            self.logger.exception('Cannot initiate OVSDB connection: %s', e)
            return None
        return self.ovs

    def _get_ofport(self, dpid, port_name):
        ovs = self._get_ovs_bridge(dpid)
        if ovs is None:
            return None

        try:
            return ovs.get_ofport(port_name)
        except Exception as e:
            self.logger.debug('Cannot get port number for %s: %s',
                              port_name, e)
            return None

    def _add_gtp_port(self, dpid, gnb_ip):
        if not self.config.multi_tunnel_flag:
            return self.config.gtp_port

        port_name = self._ip_addr_to_gtp_port_name(gnb_ip)
        # If GTP port already exists, returns OFPort number
        gtpport = self._get_ofport(dpid, port_name)
        if gtpport is not None:
            return gtpport

        ovs = self._get_ovs_bridge(dpid)
        if ovs is None:
            return None

        port_name = self._ip_addr_to_gtp_port_name(gnb_ip)
        ovs.add_tunnel_port(port_name, self.ovs_gtp_type,
                            gnb_ip, key="flow")

        return self._get_ofport(dpid, port_name)

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        if self._clean_restart:
            self._delete_all_flows()

        self._install_default_tunnel_flows()
        self._install_internal_pkt_fwd_flow()

    def _delete_all_flows(self):
        flows.delete_all_flows_from_table(self._datapath, self.tbl_num)

    def cleanup_on_disconnect(self, datapath):
        if self._clean_restart:
            self._delete_all_flows()

    def _install_default_tunnel_flows(self):
        match = MagmaMatch()
        flows.add_flow(self._datapath,self.tbl_num, match,
                       priority=flows.MINIMUM_PRIORITY,
                       goto_table=self.next_table)

    def _install_internal_pkt_fwd_flow(self):
        match = MagmaMatch(in_port=self.config.internal_sampling_port)
        flows.add_flow(self._datapath,self.tbl_num, match,
                       priority=flows.MINIMUM_PRIORITY,
                       goto_table=self.config.internal_sampling_fwd_tbl)


    def _add_tunnel_flows(self, precedence:int, i_teid:int,
                          o_teid:int, ue_ip_adr:str,
                          enodeb_ip_addr:str, sid:int = None):

        parser = self._datapath.ofproto_parser
        dpid = self._datapath.id
        priority = Utils.get_of_priority(precedence)
        # Add flow for gtp port
        if enodeb_ip_addr:
            gtp_portno = self._add_gtp_port(dpid, enodeb_ip_addr)
        else:
            gtp_portno = self.config.gtp_port
        match = MagmaMatch(tunnel_id=i_teid, in_port=gtp_portno)

        actions = [parser.OFPActionSetField(eth_src=GTP_PORT_MAC),
                   parser.OFPActionSetField(eth_dst="ff:ff:ff:ff:ff:ff")]
        if sid:
            actions.append(parser.OFPActionSetField(metadata=sid))

        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)

        # Add flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,in_port=self._uplink_port,
                           ipv4_dst=ue_ip_adr)
        actions = [parser.OFPActionSetField(tunnel_id=o_teid),
                   parser.OFPActionSetField(tun_ipv4_dst=enodeb_ip_addr),
                   parser.OFPActionSetField(tun_flags=TUNNEL_OAM_FLAG),
                   parser.NXActionRegLoad2(dst=TUN_PORT_REG, value=gtp_portno)]
        if sid:
            actions.append(parser.OFPActionSetField(metadata=sid))

        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)

        # Add flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,
                           ipv4_dst=ue_ip_adr)

        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)
       
        # Add ARP flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                           in_port=self._uplink_port, arp_tpa=ue_ip_adr)
        actions = []
        if sid:
            actions = [parser.OFPActionSetField(metadata=sid)]

        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)

        # Add ARP flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                               in_port=self.config.mtr_port,
                               arp_tpa=ue_ip_adr)

        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)


    def _delete_tunnel_flows(self, i_teid:int, ue_ip_adr:str,
                                 enodeb_ip_addr:str = None):

        dpid = self._datapath.id
        # Delete flow for gtp port
        if enodeb_ip_addr:
            gtp_portno = self._add_gtp_port(dpid, enodeb_ip_addr)
        else:
            gtp_portno = self.config.gtp_port

        match = MagmaMatch(tunnel_id=i_teid, in_port=gtp_portno)

        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self._uplink_port, ipv4_dst=ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,ipv4_dst=ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete ARP flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                           in_port=self._uplink_port, arp_tpa=ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete ARP flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                           in_port=self.config.mtr_port, arp_tpa=ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match)


    def _discard_tunnel_flows(self, precedence:int, i_teid:int,
                              ue_ip_adr:str):
        priority = Utils.get_of_priority(precedence)
        # discard uplink Tunnel
        match = MagmaMatch(tunnel_id=i_teid, in_port=self.config.gtp_port)

        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority + 1)

        # discard downlink Tunnel for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self._uplink_port, ipv4_dst=ue_ip_adr)

        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority + 1)

        # discard downlink Tunnel for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,ipv4_dst=ue_ip_adr)
 
        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority + 1)

    def _resume_tunnel_flows(self, precedence:int, i_teid:int,
                             ue_ip_adr:str):

        priority = Utils.get_of_priority(precedence)
        # Forward flow for gtp port
        match = MagmaMatch(tunnel_id=i_teid, in_port=self.config.gtp_port)

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority + 1)

        # Forward flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self._uplink_port,ipv4_dst=ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority +1)

        # Forward flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,ipv4_dst=ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority + 1)

