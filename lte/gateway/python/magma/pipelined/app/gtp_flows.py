"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import logging
from collections import namedtuple
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

from .base import MagmaController
from magma.pipelined.openflow import flows
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import load_direction, Direction, \
    PASSTHROUGH_REG_VAL

from ryu.lib.packet import ether_types

# ingress and egress service names -- used by other controllers
INGRESS = "ingress"
EGRESS = "egress"
PHYSICAL_TO_LOGICAL = "middle"
GTP_PORT_MAC = "02:00:00:00:00:01";
EnodeB_IP = "192.168.60.141"

class GtpFlows(MagmaController):
    """
    A controller that sets up an openflow pipeline for Magma.

    This controller is used for table 0 for gtp default entry and 
    gtp tunnel.
    """
    APP_NAME = "gtp_flows"
    #APP_TYPE = ControllerType.SPECIAL
    GtpConfig = namedtuple(
            'GtpConfig',
            ['gtp_port', 'uplink_port_name', 'mtr_ip', 'mtr_port'],
    )
    GTP_PORT_MAC = "02:00:00:00:00:01"
    EnodeB_IP = "192.168.60.141"
    def __init__(self, *args, **kwargs):
        super(GtpFlows, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self._uplink_port = OFPP_LOCAL
        self.tbl_num = 0
        self.i_teid = 1
        self.next_table = self._service_manager.get_table_num(INGRESS)
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False
        if (self.config.uplink_port_name):
            self._uplink_port = BridgeTools.get_ofport(self.config.uplink_port_name)
        
        #Temporary Field for UT
        self.UE_ADDR ="None"
        logging.info("GTPFlows uplink:%d ", self._uplink_port)        

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
        self._delete_all_flows(self._datapath)
        self._install_default_tunnel_flows(self._datapath)
        logging.info ("GTPflows connect")

    def _delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)        

    def _install_default_tunnel_flows(self, datapath):
        match = MagmaMatch()
        logging.info ("prabin tbl:%d ntbl:%d", self.tbl_num, self.next_table)
        
        flows.add_flow(datapath,match,[],priority=flows.MINIMUM_PRIORITY,
                        table_id=self.tbl_num,goto_table=self.next_table)
                                        
    def _add_gtp_tunnel_flows(self, ue_addr):
        
        parser = self._datapath.ofproto_parser
        ofproto = self._datapath.ofproto

        logging.info("prabina t:%d nntbl:%d", self.i_teid, self.next_table)
        #Temporary Value set for UT
        self.UE_ADDR = ue_addr
        # Add flow for gtp port
        match = MagmaMatch(tunnel_id=self.i_teid, in_port=self.config.gtp_port,
                           direction=Direction.IN)
        
        actions = [parser.OFPActionSetField(eth_src=GTP_PORT_MAC),
                   parser.OFPActionSetField(eth_dst="ff:ff:ff:ff:ff:ff"),
                   parser.OFPActionSetField(metadata=0x7594587a00d)]

        flows.add_flow(self._datapath, match, actions=actions,
                       priority=flows.DEFAULT_PRIORITY,
                       table_id=self.tbl_num, goto_table=self.next_table)
        logging.info("firstone is done")
        # Add flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,in_port=OFPP_LOCAL,
                           ipv4_dst=ue_addr , direction=Direction.IN)
        
        actions = [parser.OFPActionSetField(tunnel_id=0xa000128),
                   parser.OFPActionSetField(tun_ipv4_dst=EnodeB_IP),
                   parser.OFPActionSetField(metadata=0x7594587a00d)]
        
        flows.add_flow(self._datapath, match, actions=actions,
                       priority=flows.DEFAULT_PRIORITY,
                       table_id=self.tbl_num,goto_table=self.next_table)
        # Add flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,ipv4_dst=ue_addr,
                           direction=Direction.IN)

        actions = [parser.OFPActionSetField(tunnel_id=0xa000128),
                   parser.OFPActionSetField(tun_ipv4_dst=EnodeB_IP),
                   parser.OFPActionSetField(metadata=0x7594587a00d)]
        

        flows.add_flow(self._datapath, match, actions=actions,
                        priority=flows.DEFAULT_PRIORITY,
                        table_id=self.tbl_num,goto_table=self.next_table)

        self.i_teid = self.i_teid + 1

    def _delete_gtp_tunnel_flows(self):
        
        # Delete flow for gtp port
        match = MagmaMatch(tunnel_id=self.i_teid, in_port=self.config.gtp_port)
        flows.delete_flow(self._datapath, self.tbl_num, match)
        
        # Delete flow for LOCAL port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,in_port=OFPP_LOCAL,
                           ipv4_dst=self.UE_ADDR)
        flows.delete_flow(self._datapath, self.tbl_num, match)
        
        # Delete flow for mtr port
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           in_port=self.config.mtr_port,ipv4_dst=self.UE_ADDR)
        flows.delete_flow(self._datapath, self.tbl_num, match)

    def _add_discard_data_gtp_tunnel_flows(self, ue_addr):
        
        parser = self._datapath.ofproto_parser
        ofproto = self._datapath.ofproto
        
        # discard uplink Tunnel
        match = MagmaMatch(tunnel_id=self.i_teid, in_port=self.config.gtp_port)
        cookie = self._cookie
        cookie_mask = self._cookie
        flows.add_flow_data_gtp(self._datapath, match,priority=flows.DEFAULT_PRIORITY + 1,
                                table_id=self.tbl_num,
                                cookie=cookie, cookie_mask=cookie_mask)
        
        # discard downlink Tunnel for LOCAL port
        match = MagmaMatch(in_port=OFPP_LOCAL, ipv4_dst=ue_addr)
        cookie = self._cookie + 1
        cookie_mask = self._cookie + 1
                           
        flows.add_flow_data_gtp(self._datapath, match,priority=flows.DEFAULT_PRIORITY + 1,
                                table_id=self.tbl_num,
                                cookie=cookie, cookie_mask=cookie_mask)
        # discard downlink Tunnel for mtr port
        match = MagmaMatch(in_port=self.config.mtr_port,ipv4_dst=ue_addr)
        cookie = self._cookie + 1
        cookie_mask = self._cookie + 1
                           
        flows.add_flow_data_gtp(self._datapath, match,priority=flows.DEFAULT_PRIORITY + 1,
                                table_id=self.tbl_num,
                                cookie=cookie, cookie_mask=cookie_mask)
        
        self.i_teid = self.i_teid + 1

    def _add_forward_data_gtp_tunnel_flows(self, ue_addr):
        
        # Forward flow for gtp port
        match = MagmaMatch(tunnel_id=self.i_teid, in_port=self.config.gtp_port)
        cookie = self._cookie
        cookie_mask = self._cookie
        
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=flows.DEFAULT_PRIORITY + 1,
                          cookie=cookie, cookie_mask=cookie_mask)
        
        # Forward flow for LOCAL port
        match = MagmaMatch(in_port=OFPP_LOCAL,ipv4_dst=ue_addr)
        cookie = self._cookie + 1
        cookie_mask = self._cookie + 1
        
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=flows.DEFAULT_PRIORITY + 1,
                          cookie=cookie, cookie_mask=cookie_mask)
        
        # Forward flow for mtr port
        match = MagmaMatch(in_port=self.config.mtr_port,ipv4_dst=ue_addr)
        cookie = self._cookie + 1
        cookie_mask = self._cookie + 1
        
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          priority=flows.DEFAULT_PRIORITY + 1,
                          cookie=cookie, cookie_mask=cookie_mask)
        
