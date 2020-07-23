"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from ryu.ofproto.nicira_ext import ofs_nbits
from ryu.lib.packet import ether_types
from ryu.ofproto.inet import IPPROTO_TCP

from .base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction


class ConntrackController(MagmaController):
    """
    A controller that sets up tunnel/ue learn flows based on uplink UE traffic
    to properly route downlink packets back to the UE (through the correct GRE
    flow tunnel).

    This is an optional controller and will only be used for setups with flow
    based GRE tunnels.

    conntrack flags 0 is nothing, 1 is commit

    CT state reference tuple (x,y):
    x:
     0: -
     1: +

    y:
      0x01: new
      0x02: est
      0x04: rel
      0x08: rpl
      0x10: inv
      0x20: trk
    """

    APP_NAME = "conntrack"
    APP_TYPE = ControllerType.LOGICAL
    CT_NEW = 0x01
    CT_EST = 0x02
    CT_REL = 0x04
    CT_RPL = 0x08
    CT_INV = 0x10
    CT_TRK = 0x20

    def __init__(self, *args, **kwargs):
        super(ConntrackController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self.conntrack_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]
        self.connection_event_table = \
            self._service_manager.INTERNAL_IPFIX_SAMPLE_TABLE_NUM
        self._datapath = None

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_flows(self._datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self.conntrack_scratch)

    def _install_default_flows(self, datapath):
        parser = datapath.ofproto_parser

        match =  MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                            ct_state=(0x0, self.CT_TRK))
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.conntrack_scratch,
            alg=0,
            actions=[]
        )]
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             match, actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=self.next_table)

        # Match all new connections
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ct_state=(self.CT_NEW | self.CT_TRK,
                                     self.CT_NEW | self.CT_TRK))
        actions = [parser.NXActionCT(
            flags=0x1,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY)

        # Match tcp terminations (fin)
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ip_proto=IPPROTO_TCP,
                           tcp_flags=(0x1,0x1),
                           ct_state=(self.CT_EST | self.CT_TRK,
                                     self.CT_EST | self.CT_TRK))
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY)
        # match tcp fin
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ip_proto=IPPROTO_TCP,
                           tcp_flags=(0x1, 0x1),
                           ct_state=(self.CT_TRK | self.CT_INV, self.CT_TRK | self.CT_INV))
        # flags 0 is nothing, 1 is commit
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY)

        # match tcp rst
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ip_proto=IPPROTO_TCP,
                           tcp_flags=(0x4, 0x4),
                           ct_state=(self.CT_EST | self.CT_TRK, self.CT_EST | self.CT_TRK))
        # flags 0 is nothing, 1 is commit
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY)

        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, inbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_table)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, outbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_table)


        # TODO Currently for testing, will nuke later
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ct_state=(self.CT_EST | self.CT_TRK, self.CT_EST | self.CT_TRK))
        # flags 0 is nothing, 1 is commit
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY-5)
        # match tcp fin
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ip_proto=IPPROTO_TCP,
                           tcp_flags=(0x1, 0x1),
                           ct_state=(self.CT_TRK, self.CT_TRK))
        # flags 0 is nothing, 1 is commit
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY-1)
        # match tcp fin
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ip_proto=IPPROTO_TCP,
                           tcp_flags=(0x1,0x1))
        # flags 0 is nothing, 1 is commit
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=0,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY-1)