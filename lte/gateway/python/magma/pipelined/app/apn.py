"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.registers import Direction, APN_TAG_REG
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.apn import encode_apn

from ryu.controller.controller import  Datapath

class APNController(MagmaController):
    """
    APN Controller.


    The APN controller is responsible for marking a flow with an encoded APN name. The APN name should be stored in
    register.

    """

    APP_NAME = "apn"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(APNController, self).__init__(*args, **kwargs)
        self._clean_start = True # get from config file
        self._datapath = None

        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self.apn_tagging_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]


    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self.delete_all_flows(datapath)
        # In case wee need to clean all existing  buggy / orphaned flows before start the controller
        # if self._clean_start
        #   self.delete_existing_flows()
        pass

    def add_apn_flow_for_ue(self, ue_ip_addr, apn):
        """ Add flow which match all IN traffic with specified UE_IP and set APN hash in to register.

        Args:
            ue_ip_addr: ip addr allocated for the UE in scope of connection to specific APN
            apn: APN UE is connected to with specified IP addr
        """
        # TODO(119vik): same IP is reused for several bearers connected to the same APN - take care about duplications

        parser = self._datapath.ofproto_parser

        # Tag all downlink traffic
        outbound_match = MagmaMatch(direction=Direction.OUT, ipv4_dst=ue_ip_addr)
        actions = [
            parser.NXActionRegLoad2(dst=APN_TAG_REG, value=encode_apn(apn)),
            parser.NXActionResubmitTable(table_id=self.apn_tagging_scratch)
        ]
        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             outbound_match , actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

        # Tag all uplink traffic
        inbound_match = MagmaMatch(direction=Direction.IN, ipv4_src=ue_ip_addr)
        actions = [
            parser.NXActionRegLoad2(dst=APN_TAG_REG, value=encode_apn(apn)),
            parser.NXActionResubmitTable(table_id=self.apn_tagging_scratch)
        ]
        flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                                             inbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def delete_apn_flow_for_ue(self, ue_ip_addr, apn):
        """ Delete flow been created in scope of add_apn_flow_for_ue.

        Args:
            ue_ip_addr: ip addr allocated for the UE in scope of connection to specific APN
            apn: APN UE is connected to with specified IP addr
        """
        # TODO(119vik): same IP is reused for several bearers connected to the same APN - take care about duplications
        # flow delete
        pass

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        """ Delete all flows which set APN register"""
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self.apn_tagging_scratch)