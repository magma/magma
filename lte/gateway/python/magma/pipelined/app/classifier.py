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

from collections import namedtuple
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from .base import MagmaController
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.app.inout import INGRESS
from ryu.lib.packet import ether_types
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.utils import Utils

GTP_PORT_MAC = "02:00:00:00:00:01"

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
            ['gtp_port', 'mtr_ip', 'mtr_port', 'internal_sampling_port', 'internal_sampling_fwd_tbl'],
    )

    def __init__(self, *args, **kwargs):
        super(Classifier, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_table_num(INGRESS)
        self._uplink_port = OFPP_LOCAL
        self._datapath = None

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None
        
        if 'mtr_ip' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']

        return self.ClassifierConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            internal_sampling_port=
                               config_dict['ovs_internal_sampling_port_number'],
            internal_sampling_fwd_tbl=
                              config_dict['ovs_internal_sampling_fwd_tbl_number'],

        )

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self._delete_all_flows()
        self._install_default_tunnel_flows()
        self._install_internal_pkt_fwd_flow()

    def _delete_all_flows(self):
        flows.delete_all_flows_from_table(self._datapath, self.tbl_num)

    def cleanup_on_disconnect(self, datapath):
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
        priority = Utils.get_of_priority(precedence)
        # Add flow for gtp port
        match = MagmaMatch(tunnel_id=i_teid, in_port=self.config.gtp_port)

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
                   parser.OFPActionSetField(tun_ipv4_dst=enodeb_ip_addr)]
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


    def _delete_tunnel_flows(self, i_teid:int, ue_ip_adr:str):

        # Delete flow for gtp port
        match = MagmaMatch(tunnel_id=i_teid, in_port=self.config.gtp_port)

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

