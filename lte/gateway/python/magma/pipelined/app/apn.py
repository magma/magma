"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.registers import Direction, APN_TAG_REG, APN_TAG_MAP
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.apn import encode_apn, split_apn
from magma.pipelined.openflow.messages import MsgChannel, MessageHub

from ryu.controller.controller import  Datapath
from ryu.ofproto import nicira_ext
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.controller import ofp_event


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
        self._msg_hub = MessageHub(self.logger)

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        # In case wee need to clean all existing  buggy / orphaned flows before start the controller
        if self._clean_start:
          self.delete_all_flows(datapath)

    def add_apn_flow_for_ue(self, imsi, ue_ip_addr, apn):
        """ Add flow which match all IN traffic with specified UE_IP and set APN hash in to register.

        Args:
            imsi: user's IMSI
            ue_ip_addr: ip addr allocated for the UE in scope of connection to specific APN
            apn: APN UE is connected to with specified IP addr
        """
        #TODO(119vik): same IP is reused for several bearers connected to the same APN - take care about duplications
        parser = self._datapath.ofproto_parser

        encoded_apn = encode_apn(apn)
        encoded_apn_registers = split_apn(encoded_apn)
        flow_adds = []

        # Tag all Uplink traffic
        outbound_match = MagmaMatch(direction=Direction.OUT, ipv4_src=ue_ip_addr, eth_type=2048)
        actions = None
        flow_adds.append(flows.get_add_resubmit_next_service_flow_msg(
            self._datapath, 
            self.tbl_num,
            outbound_match,
            actions,
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self.apn_tagging_scratch))
        
        # four low lever 4bytes registers represent single 16 bytes double extended register 
        actions = [
            parser.NXActionRegLoad2(
                dst=APN_TAG_MAP[APN_TAG_REG][reg_num], 
                value=int(encoded_apn_registers[reg_num], base=16)
            )
            for reg_num in range(4)
        ]
        
        flow_adds.append(flows.get_add_resubmit_next_service_flow_msg(
            self._datapath, 
            self.apn_tagging_scratch,
            outbound_match,
            actions,
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self.next_table))

        # Tag all downlink traffic
        inbound_match = MagmaMatch(direction=Direction.IN, ipv4_dst=ue_ip_addr)
        actions = None
        flow_adds.append(flows.get_add_resubmit_next_service_flow_msg(
            self._datapath, 
            self.tbl_num,
            inbound_match,
            actions,
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self.apn_tagging_scratch))
        
        # four low lever 4bytes registers represent single 16 bytes double extended register 
        actions = [
            parser.NXActionRegLoad2(
                dst=APN_TAG_MAP[APN_TAG_REG][reg_num], 
                value=int(encoded_apn_registers[reg_num], base=16)
            )
            for reg_num in range(4)
        ]
        
        flow_adds.append(flows.get_add_resubmit_next_service_flow_msg(
            self._datapath, 
            self.apn_tagging_scratch,
            inbound_match,
            actions,
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self.next_table))
        logging.info("Flows to add {}, datapath {}".format(flow_adds, self._datapath))
        chan = self._msg_hub.send(flow_adds, self._datapath)

        return self._wait_for_flow_responses(imsi, flow_adds, chan)

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def _wait_for_flow_responses(self, imsi, flow_adds, chan):
        def fail(err):
            self.logger.error(
                "Failed to install flow for subscriber %s: %s",
                imsi, err)
            return False

        for _ in range(len(flow_adds)):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS")
            if not result.ok():
                return fail(result.exception())
        return True


    def delete_apn_flow_for_ue(self, imsi, ue_ip_addr, apn):
        """ Delete flow been created in scope of add_apn_flow_for_ue.

        Args:
            imsi: user's IMSI
            ue_ip_addr: ip addr allocated for the UE in scope of connection to specific APN
            apn: APN UE is connected to with specified IP addr
        """
        # TODO(119vik): Add Flow deletion implementation
        pass

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        """Delete all flows which set APN register"""
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self.apn_tagging_scratch)