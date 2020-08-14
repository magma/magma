"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from collections import namedtuple
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from .base import MagmaController
from magma.pipelined.openflow import flows
#from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.app.inout import INGRESS
from ryu.lib.packet import ether_types
from magma.pipelined.app.base import MagmaController, ControllerType
#from lte.protos.smfupfif_pb2 import SetGroupPDR, SetGroupFAR

GTP_PORT_MAC = "02:00:00:00:00:01"
GTP = "gtp"

class GtpFlows(MagmaController):
    """
    A controller that sets up an openflow pipeline for Magma.

    This controller is used for table 0 for gtp default entry and 
    gtp tunnel.
    """
    APP_NAME = "gtp_flows"
    APP_TYPE = ControllerType.SPECIAL
    GtpConfig = namedtuple(
            'GtpConfig',
            ['gtp_port', 'uplink_port_name', 'mtr_ip', 'mtr_port'],
    )

    def __init__(self, *args, **kwargs):
        super(GtpFlows, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_table_num(INGRESS)
        self._uplink_port = OFPP_LOCAL
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False
        #if (self.config.uplink_port_name):
         #   self._uplink_port = BridgeTools.get_ofport(self.config.uplink_port_name)
        self._cookie = 1
        self._datapath = None

    def _get_config(self, config_dict):
        port_name = None
        mtr_ip = None
        mtr_port = None
        
        if 'ovs_uplink_port_name' in config_dict:
            port_name = config_dict['ovs_uplink_port_name']

        if 'mtr_ip' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']
        return self.GtpConfig(
            gtp_port=config_dict['ovs_gtp_port_number'],
            uplink_port_name=port_name,
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
        )

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self._delete_all_flows()
        self._install_default_tunnel_flows()

    def _delete_all_flows(self):
        flows.delete_all_flows_from_table(self._datapath, self.tbl_num)

    def cleanup_on_disconnect(self, datapath):
        self._delete_all_flows()

    def convert_precedence_to_priority(self, precedence):
        if precedence < flows.MAXIMUM_PRIORITY:
            priority = flows.MAXIMUM_PRIORITY - precedence
        else:
            priority = 0

        if priority < flows.DEFAULT_PRIORITY:
            priority = flows.DEFAULT_PRIORITY
        return priority

    def _install_default_tunnel_flows(self):
        match = MagmaMatch()
        flows.add_flow(self._datapath,self.tbl_num, match,
                       priority=flows.MINIMUM_PRIORITY,
                       goto_table=self.next_table)

    def _add_gtp_tunnel_flows(self, SetPDR, SetFAR, seid):
        
        parser = self._datapath.ofproto_parser
        priority = self.convert_precedence_to_priority(SetPDR.precedence)
        # Add flow for gtp port
        match = MagmaMatch(tunnel_id=SetPDR.pdi.local_f_teid.teid,
                          in_port=self.config.gtp_port)

        actions = [parser.OFPActionSetField(eth_src=GTP_PORT_MAC),
                   parser.OFPActionSetField(eth_dst="ff:ff:ff:ff:ff:ff"),
                   parser.OFPActionSetField(metadata=seid)]
        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)

        # Add flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,in_port=self._uplink_port,
                           ipv4_dst=SetPDR.pdi.ue_ip_adr)

        actions = [parser.OFPActionSetField(tunnel_id=SetFAR.fwd_parm.ohdrcr.o_teid),
                   parser.OFPActionSetField(tun_ipv4_dst=SetFAR.fwd_parm.ohdrcr.ipv4_adr),
                   parser.OFPActionSetField(metadata=seid)]

        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                       priority=priority, goto_table=self.next_table)

        # Add flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                       in_port=self.config.mtr_port,
                       ipv4_dst=SetPDR.pdi.ue_ip_adr)

        actions = [parser.OFPActionSetField(tunnel_id=SetFAR.fwd_parm.ohdrcr.o_teid),
                   parser.OFPActionSetField(tun_ipv4_dst=SetFAR.fwd_parm.ohdrcr.ipv4_adr),
                   parser.OFPActionSetField(metadata=seid)]
        
        flows.add_flow(self._datapath, self.tbl_num, match, actions=actions,
                        priority=priority, goto_table=self.next_table)
       
        # Add ARP flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,in_port=self._uplink_port,
                           arp_tpa=SetPDR.pdi.ue_ip_adr)
        
        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority, goto_table=self.next_table)

        # Add ARP flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                       in_port=self.config.mtr_port,
                       arp_tpa=SetPDR.pdi.ue_ip_adr)

        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority, goto_table=self.next_table)

    def _delete_gtp_tunnel_flows(self, SetPDR):

        # Delete flow for gtp port
        match = MagmaMatch(tunnel_id=SetPDR.pdi.local_f_teid.teid,
                           in_port=self.config.gtp_port)
        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,in_port=self._uplink_port,
                           ipv4_dst=SetPDR.pdi.ue_ip_adr)
        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,ipv4_dst=SetPDR.pdi.ue_ip_adr)
        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete ARP flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,in_port=self._uplink_port,
                           arp_tpa=SetPDR.pdi.ue_ip_adr)
        
        flows.delete_flow(self._datapath, self.tbl_num, match)

        # Delete ARP flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_ARP,
                           in_port=self.config.mtr_port,
                           arp_tpa=SetPDR.pdi.ue_ip_adr)

        flows.delete_flow(self._datapath, self.tbl_num, match)


    def _add_discard_data_gtp_tunnel_flows(self, SetPDR):

        priority = self.convert_precedence_to_priority(SetPDR.precedence)
        # discard uplink Tunnel
        match = MagmaMatch(tunnel_id=SetPDR.pdi.local_f_teid.teid,
                           in_port=self.config.gtp_port)
        cookie = self._cookie
        #cookie_mask = self._cookie
        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority + 1,cookie=cookie)

        # discard downlink Tunnel for LOCAL port
        match = MagmaMatch(in_port=self._uplink_port, ipv4_dst=SetPDR.pdi.ue_ip_adr)
        cookie = self._cookie + 1
        #cookie_mask = self._cookie + 1

        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority + 1,cookie=cookie)

        # discard downlink Tunnel for mtr port
        match = MagmaMatch(in_port=self.config.mtr_port,ipv4_dst=SetPDR.pdi.ue_ip_adr)
        cookie = self._cookie + 1
        #cookie_mask = self._cookie + 1
                           
        flows.add_flow(self._datapath, self.tbl_num, match,
                       priority=priority + 1,cookie=cookie)

    def _add_forward_data_gtp_tunnel_flows(self, SetPDR):
        
        priority = self.convert_precedence_to_priority(SetPDR.precedence)
        # Forward flow for gtp port
        match = MagmaMatch(tunnel_id=SetPDR.pdi.local_f_teid.teid, in_port=self.config.gtp_port)
        cookie = self._cookie
        cookie_mask = self._cookie

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority + 1,
                          cookie=cookie, cookie_mask=cookie_mask)

        # Forward flow for LOCAL port
        match = MagmaMatch(in_port=self._uplink_port,ipv4_dst=SetPDR.pdi.ue_ip_adr)
        cookie = self._cookie + 1
        cookie_mask = self._cookie + 1

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority +1,
                          cookie=cookie, cookie_mask=cookie_mask)

        # Forward flow for mtr port
        match = MagmaMatch(in_port=self.config.mtr_port,ipv4_dst=SetPDR.pdi.ue_ip_adr)
        cookie = self._cookie + 1
        cookie_mask = self._cookie + 1

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=priority + 1,
                          cookie=cookie, cookie_mask=cookie_mask)
