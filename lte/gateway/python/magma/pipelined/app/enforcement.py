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
from lte.protos.pipelined_pb2 import RuleModResult

from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.policy_mixin import PolicyMixin
from magma.pipelined.app.restart_mixin import RestartMixin, DefaultMsgsMap

from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.policy_converters import FlowMatchError, \
    get_ue_ip_match_args, get_eth_type
from magma.pipelined.redirect import RedirectionManager, RedirectException
from magma.pipelined.qos.common import QosManager
from magma.pipelined.qos.qos_meter_impl import MeterManager

from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from magma.pipelined.utils import Utils
from magma.pipelined.openflow.exceptions import MagmaDPDisconnectedError


class EnforcementController(PolicyMixin, RestartMixin, MagmaController):
    """
    EnforcementController

    The enforcement controller installs flows for policy enforcement and
    classification. Each flow installed matches on a rule and an IMSI and then
    classifies the packet with the rule. The flow also redirects and drops
    the packet as specified in the policy.

    NOTE: Enforcement currently relies on the fact that policies do not
    overlap. In this implementation, there is the idea of a 'default rule'
    which is the catch-all. This rule is treated specially and tagged with a
    specific priority.
    """

    APP_NAME = "enforcement"
    APP_TYPE = ControllerType.LOGICAL
    DEFAULT_FLOW_COOKIE = 0xfffffffffffffffe

    def __init__(self, *args, **kwargs):
        super(EnforcementController, self).__init__(*args, **kwargs)
        self._config = kwargs['config']
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            EnforcementStatsController.APP_NAME)
        self._enforcement_stats_tbl = self._service_manager.get_table_num(
            EnforcementStatsController.APP_NAME)
        self.loop = kwargs['loop']

        self._msg_hub = MessageHub(self.logger)
        self._redirect_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]
        self._bridge_ip_address = kwargs['config']['bridge_ip_address']
        self._redirect_manager = None
        self._qos_mgr = None
        self._clean_restart = kwargs['config']['clean_restart']
        self._redirect_manager = RedirectionManager(
            self._bridge_ip_address,
            self.logger,
            self.tbl_num,
            self._enforcement_stats_tbl,
            self.next_main_table,
            self._redirect_scratch,
            self._session_rule_version_mapper)

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self._qos_mgr = QosManager.get_qos_manager(datapath, self.loop, self._config)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        if self._clean_restart:
            self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self._redirect_scratch)

    def cleanup_state(self):
        pass

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    @set_ev_cls(ofp_event.EventOFPMeterConfigStatsReply, MAIN_DISPATCHER)
    def meter_config_stats_reply_handler(self, ev):
        if not self._qos_mgr:
            return

        qos_impl = self._qos_mgr.impl
        if qos_impl and isinstance(qos_impl, MeterManager):
            qos_impl.handle_meter_config_stats(ev.msg.body)

    @set_ev_cls(ofp_event.EventOFPMeterFeaturesStatsReply, MAIN_DISPATCHER)
    def meter_features_stats_reply_handler(self, ev):
        if not self._qos_mgr:
            return

        qos_impl = self._qos_mgr.impl
        if qos_impl and isinstance(qos_impl, MeterManager):
            qos_impl.handle_meter_feature_stats(ev.msg.body)

    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        """
        Gets the default flow msg that forward to stats table(traffic will be
        dropped because stats table doesn't forward anything)

        Args:
            datapath: ryu datapath struct
        Returns:
            The list of default msgs to add
        """
        match = MagmaMatch()
        msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self._enforcement_stats_tbl,
            cookie=self.DEFAULT_FLOW_COOKIE)

        return {self.tbl_num: [msg]}

    def _get_rule_match_flow_msgs(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, rule,
                                  ng_session_id=0):
        """
        Get flow msgs to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            msisdn (bytes): subscriber MSISDN
            uplink_tunnel (int): tunnel ID of the subscriber.
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        priority = Utils.get_of_priority(rule.priority)

        flow_adds = []
        for flow in rule.flow_list:
            try:
                if ng_session_id:
                    version = self._session_rule_version_mapper.\
                                   get_ng_version_by_session_teid(ng_session_id)
                else:                   
                    version = self._session_rule_version_mapper.get_version(imsi, ip_addr,
                                                                            rule.id)

                if uplink_tunnel:     
                    self.logger.info("uplink_tunnel=%d, version=%d", uplink_tunnel, version)

                flow_adds.extend(self._get_classify_rule_flow_msgs(
                    imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, flow, rule_num, priority,
                    rule.qos, rule.hard_timeout, rule.id, rule.app_name,
                    rule.app_service_type, self.next_main_table,
                    version, self._qos_mgr, self._enforcement_stats_tbl, rule.he.urls, ng_session_id))

            except FlowMatchError as err:  # invalid match
                self.logger.error(
                    "Failed to get flow msg '%s' for subscriber %s: %s",
                    rule.id, imsi, err)
                raise err
        return flow_adds

    def _install_flow_for_rule(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, rule,
                               ng_session_id=0):
        """
        Install a flow to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            msisdn (bytes): subscriber MSISDN
            uplink_tunnel (int): tunnel ID of the subscriber.
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        if rule.redirect.support == rule.redirect.ENABLED:
            return self._install_redirect_flow(imsi, ip_addr, rule)

        if not rule.flow_list:
            self.logger.error('The flow list for imsi %s, rule.id - %s'
                              'is empty, this shoudn\'t happen', imsi, rule.id)
            return RuleModResult.FAILURE

        flow_adds = []
        try:
            flow_adds = self._get_rule_match_flow_msgs(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule,
                                                       ng_session_id)
        except FlowMatchError:
            return RuleModResult.FAILURE

        try:
            chan = self._msg_hub.send(flow_adds, self._datapath)
        except MagmaDPDisconnectedError:
            self.logger.error("Datapath disconnected, failed to install rule %s"
                              "for imsi %s", rule, imsi)
            return RuleModResult.FAILURE
        return self._wait_for_rule_responses(imsi, ip_addr, rule, chan)

    def _install_redirect_flow(self, imsi, ip_addr, rule):
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        rule_version = self._session_rule_version_mapper.get_version(imsi,
                                                                     ip_addr,
                                                                     rule.id)
        priority = Utils.get_of_priority(rule.priority)
        redirect_request = RedirectionManager.RedirectRequest(
            imsi=imsi,
            ip_addr=ip_addr.address.decode('utf-8'),
            rule=rule,
            rule_num=rule_num,
            rule_version=rule_version,
            priority=priority)
        try:
            self._redirect_manager.setup_lte_redirect(
                self._datapath, self.loop, redirect_request)
            return RuleModResult.SUCCESS
        except RedirectException as err:
            self.logger.error(
                'Redirect Exception for imsi %s, rule.id - %s : %s',
                imsi, rule.id, err
            )
            return RuleModResult.FAILURE

    def _get_default_flow_msgs_for_subscriber(self, *_):
        pass

    def _install_default_flow_for_subscriber(self, imsi, ip_addr):
        pass

    def _install_default_ng_flow_for_subscriber(self, imsi, ip_addr, ng_session_id):
        pass

    def _deactivate_flow_for_rule(self, imsi, ip_addr, rule_id):
        """
        Deactivate a specific rule using the flow cookie for a subscriber
        """
        try:
            num = self._rule_mapper.get_rule_num(rule_id)
        except KeyError:
            self.logger.error('Could not find rule id %s', rule_id)
            return
        if num is None:
            self.logger.error('Rule num is None for rule %s', rule_id)
            return
        cookie, mask = (num, flows.OVS_COOKIE_MATCH_ALL)

        ip_match_in = get_ue_ip_match_args(ip_addr, Direction.IN)
        match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                           imsi=encode_imsi(imsi), **ip_match_in)
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          cookie=cookie, cookie_mask=mask)
        ip_match_out = get_ue_ip_match_args(ip_addr, Direction.OUT)
        match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                           imsi=encode_imsi(imsi), **ip_match_out)
        flows.delete_flow(self._datapath, self.tbl_num, match,
                          cookie=cookie, cookie_mask=mask)
        self._redirect_manager.deactivate_flow_for_rule(self._datapath, imsi,
                                                        num)
        self._qos_mgr.remove_subscriber_qos(imsi, num)
        self._remove_he_flows(ip_addr, rule_id, num)

    def _deactivate_ng_flow_for_rule(self, imsi, ip_addr, rule_id, uplink_tunnel=0):
        """
        Deactivate a specific rule using the flow cookie for a subscriber
        """
        try:
            num = self._rule_mapper.get_rule_num(rule_id)
        except KeyError:
            self.logger.error('Could not find rule id %s', rule_id)
            return
        cookie, mask = (num, flows.OVS_COOKIE_MATCH_ALL)

        if uplink_tunnel:
            ip_match_out = get_ue_ip_match_args(ip_addr, Direction.OUT)
            match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                              imsi=encode_imsi(imsi), **ip_match_out)
        else:
            ip_match_in = get_ue_ip_match_args(ip_addr, Direction.IN)
            match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                               imsi=encode_imsi(imsi), **ip_match_in)

        flows.delete_flow(self._datapath, self.tbl_num, match,
                          cookie=cookie, cookie_mask=mask)

        #TODO
        #self._redirect_manager.deactivate_flow_for_rule(self._datapath, imsi,
        #                                                num)
        #self._qos_mgr.remove_subscriber_qos(imsi, num)

    def _deactivate_flows_for_subscriber(self, imsi, ip_addr):
        """ Deactivate all rules for specified subscriber session """
        ip_match_in = get_ue_ip_match_args(ip_addr, Direction.IN)
        match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                           imsi=encode_imsi(imsi), **ip_match_in)
        flows.delete_flow(self._datapath, self.tbl_num, match)
        ip_match_out = get_ue_ip_match_args(ip_addr, Direction.OUT)
        match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                           imsi=encode_imsi(imsi), **ip_match_out)
        flows.delete_flow(self._datapath, self.tbl_num, match)

        self._redirect_manager.deactivate_flows_for_subscriber(self._datapath,
                                                               imsi)
        self._qos_mgr.remove_subscriber_qos(imsi)
        self._remove_he_flows(ip_addr)


    def _deactivate_ng_flows_for_subscriber(self, imsi, ip_addr, uplink_tunnel=0):
        """ Deactivate all rules for specified subscriber session """

        if uplink_tunnel:
            ip_match_out = get_ue_ip_match_args(ip_addr, Direction.OUT)
            match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                              imsi=encode_imsi(imsi), **ip_match_out)

        else:
            ip_match_in = get_ue_ip_match_args(ip_addr, Direction.IN)
            match = MagmaMatch(eth_type=get_eth_type(ip_addr),
                               imsi=encode_imsi(imsi), **ip_match_in)

        flows.delete_flow(self._datapath, self.tbl_num, match)

         #TODO
         #self._redirect_manager.deactivate_flows_for_subscriber(self._datapath,
         #                                                      imsi)
         #self._qos_mgr.remove_subscriber_qos(imsi)


    def deactivate_rules(self, imsi, ip_addr, rule_ids, uplink_tunnel=0, ng_session_id=0):
        """
        Deactivate flows for a subscriber.
            Only imsi -> remove all rules for imsi
            imsi+ipv4 -> remove all rules for imsi session
            imsi+rule_ids -> remove specific rules for imsi (for all sessions)
            imsi+ipv4+rule_ids -> remove rules for specific imsi session

        Args:
            imsi (string): subscriber id
            ip_addr (string): subscriber ip address
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
            if ng_session_id:
                self._deactivate_ng_flows_for_subscriber(imsi, ip_addr, uplink_tunnel)
            else:
                self._deactivate_flows_for_subscriber(imsi, ip_addr)
        else:
            for rule_id in rule_ids:
                if ng_session_id:
                    self._deactivate_ng_flow_for_rule(imsi, ip_addr, rule_id, uplink_tunnel)
                else:
                    self._deactivate_flow_for_rule(imsi, ip_addr, rule_id)

    def recover_state(self, _):
        pass
