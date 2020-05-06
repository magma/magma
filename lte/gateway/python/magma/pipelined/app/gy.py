"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import List

from ryu.controller import ofp_event
from ryu.lib.packet import ether_types
from ryu.ofproto.ofproto_v1_4_parser import OFPFlowStats
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls

from lte.protos.pipelined_pb2 import RuleModResult
from magma.pipelined.openflow.messages import MessageHub
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.inout import EGRESS
from magma.pipelined.app.policy_mixin import PolicyMixin
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.redirect import RedirectionManager, RedirectException


class GYController(PolicyMixin, MagmaController):
    """
    GYController

    The GY controller installs flows for enforcement of GY final actions, this
    includes redirection and QoS(currently not supported)
    """

    APP_NAME = "gy"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(GYController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.loop = kwargs['loop']
        self._msg_hub = MessageHub(self.logger)
        self._internal_ip_allocator = kwargs['internal_ip_allocator']
        tbls = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 2)
        self._redirect_scratch = tbls[0]
        self._mac_rewr = 210#tbls[1]
        self._bridge_ip_address = kwargs['config']['bridge_ip_address']
        self._clean_restart = kwargs['config']['clean_restart']
        self._redirect_manager = \
            RedirectionManager(
                self._bridge_ip_address,
                self.logger,
                self.tbl_num,
                self._service_manager.get_table_num(EGRESS),
                self._redirect_scratch,
                self._session_rule_version_mapper
            ).set_cwf_args(
                internal_ip_allocator=kwargs['internal_ip_allocator'],
                arp=kwargs['app_futures']['arpd'],
                mac_rewrite=self._mac_rewr,
                bridge_name=kwargs['config']['bridge_name']
            )

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self._delete_all_flows(datapath)
        self._install_default_flows(datapath)

    def deactivate_rules(self, imsi, rule_ids):
        """
        Deactivate flows for a subscriber. If only imsi is present, delete all
        rule flows for a subscriber (i.e. end its session). If rule_ids are
        present, delete the rule flows for that subscriber.

        Args:
            imsi (string): subscriber id
            rule_ids (list of strings): policy rule ids
        """
        if not self.init_finished:
            self.logger.error('Pipelined is not initialized')
            return RuleModResult.FAILURE

        if self._datapath is None:
            self.logger.error('Datapath not initialized')
            return

        if not imsi:
            self.logger.error('No subscriber specified')
            return

        if not rule_ids:
            self._deactivate_flows_for_subscriber(imsi)
        else:
            for rule_id in rule_ids:
                self._deactivate_flow_for_rule(imsi, rule_id)

    def cleanup_state(self):
        pass

    def _deactivate_flows_for_subscriber(self, imsi):
        """ Deactivate all rules for a subscriber, ending any enforcement """
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match)
        self._redirect_manager.deactivate_flows_for_subscriber(self._datapath,
                                                               imsi)

    def _deactivate_flow_for_rule(self, imsi, rule_id):
        """
        Deactivate a specific rule using the flow cookie for a subscriber
        """
        try:
            num = self._rule_mapper.get_rule_num(rule_id)
        except KeyError:
            self.logger.error('Could not find rule id %s', rule_id)
            return
        cookie, mask = (num, flows.OVS_COOKIE_MATCH_ALL)
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          cookie=cookie, cookie_mask=mask)
        self._redirect_manager.deactivate_flow_for_rule(self._datapath, imsi,
                                                        num)

    def _install_flow_for_rule(self, imsi, ip_addr, rule):
        if rule.redirect.support == rule.redirect.ENABLED:
            self._install_redirect_flow(imsi, ip_addr, rule)
            return RuleModResult.SUCCESS
        else:
            # TODO: Add support once sessiond implements restrict access QOS
            self.logger.error('GY only supports FINAL action redirect, other'
                              'final actions are not supported')
            return RuleModResult.FAILURE

    def _install_default_flow_for_subscriber(self, imsi):
        pass

    def _delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._redirect_scratch)
        flows.delete_all_flows_from_table(datapath, self._mac_rewr)

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next app.
        The enforcement flows for each subscriber would be added when the
        IP session is created, by reaching out to the controller/PCRF.

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

    def _install_redirect_flow(self, imsi, ip_addr, rule):
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        priority = rule.priority
        # TODO currently if redirection is enabled we ignore other flows
        # from rule.flow_list, confirm that this is the expected behaviour
        redirect_request = RedirectionManager.RedirectRequest(
            imsi=imsi,
            ip_addr=ip_addr,
            rule=rule,
            rule_num=rule_num,
            priority=priority)
        try:
            self._redirect_manager.setup_cwf_redirect(
                self._datapath, self.loop, redirect_request)
            return RuleModResult.SUCCESS
        except RedirectException as err:
            self.logger.error(
                'Redirect Exception for imsi %s, rule.id - %s : %s',
                imsi, rule.id, err
            )
            return RuleModResult.FAILURE

    def _install_default_flows_if_not_installed(self, datapath,
            existing_flows: List[OFPFlowStats]) -> List[OFPFlowStats]:
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)

        inbound_msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, inbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        outbound_msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, outbound_match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        msgs, remaining_flows = self._msg_hub \
            .filter_msgs_if_not_in_flow_list([inbound_msg, outbound_msg],
                                             existing_flows)
        if msgs:
            chan = self._msg_hub.send(msgs, datapath)
            self._wait_for_responses(chan, len(msgs))

        return remaining_flows

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)
