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
from typing import List
from abc import ABCMeta, abstractmethod

from ryu.ofproto.ofproto_v1_4_parser import OFPFlowStats

from lte.protos.pipelined_pb2 import RuleModResult, SetupFlowsResult, \
    ActivateFlowsResult, ActivateFlowsRequest
from magma.pipelined.app.base import ControllerNotReadyException
from magma.pipelined.openflow import flows
from magma.policydb.rule_store import PolicyRuleDict
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, IMSI_REG, \
    DIRECTION_REG, SCRATCH_REGS, RULE_VERSION_REG, RULE_NUM_REG
from magma.pipelined.openflow.messages import MsgChannel

from lte.protos.policydb_pb2 import PolicyRule
from magma.pipelined.app.dpi import UNCLASSIFIED_PROTO_ID, get_app_id
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.policy_converters import FlowMatchError, \
    flow_match_to_magma_match, convert_ipv4_str_to_ip_proto, \
    get_flow_ip_dst, ipv4_address_to_str, get_direction_for_match
from lte.protos.mobilityd_pb2 import IPAddress

from magma.pipelined.qos.types import QosInfo
from magma.pipelined.utils import Utils

PROCESS_STATS = 0x0
IGNORE_STATS = 0x1
DROP_FLOW_STATS = 0x2


class PolicyMixin(metaclass=ABCMeta):
    """
    PolicyMixin

    Mixin class for policy enforcement apps that includes common methods
    used for rule activation/deactivation.
    """
    def __init__(self, *args, **kwargs):
        super(PolicyMixin, self).__init__(*args, **kwargs)
        self._datapath = None
        self._policy_dict = PolicyRuleDict()
        self._rule_mapper = kwargs['rule_id_mapper']
        self._session_rule_version_mapper = kwargs[
            'session_rule_version_mapper']
        if 'proxy' in kwargs['app_futures']:
            self.proxy_controller_fut = kwargs['app_futures']['proxy']
        else:
            self.proxy_controller_fut = None
        self.proxy_controller = None

    def handle_restart(self,
                       requests: List[ActivateFlowsRequest]
                       ) -> SetupFlowsResult:
        """
        Setup the policy flows for subscribers, this is used when
        the controller restarts.
        """
        if self._clean_restart:
            self.delete_all_flows(self._datapath)
            self.cleanup_state()
            self.logger.info('Controller is in clean restart mode, remaining '
                              'flows were removed, continuing with setup.')

        if self._startup_flow_controller is None:
            if (self._startup_flows_fut.done()):
                self._startup_flow_controller = self._startup_flows_fut.result()
            else:
                self.logger.error('Flow Startup controller is not ready')
                return SetupFlowsResult.FAILURE
        try:
            startup_flows = \
                self._startup_flow_controller.get_flows(self.tbl_num)
        except ControllerNotReadyException as err:
            self.logger.error('Setup failed: %s', err)
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

        self.logger.debug('Setting up %s default rules', self.APP_NAME)
        remaining_flows = self._install_default_flows_if_not_installed(
            self._datapath, startup_flows)

        self.logger.debug('Startup flows before filstering -> %s',
            [flow.match for flow in startup_flows])
        extra_flows = self._add_missing_flows(requests, remaining_flows)

        self.logger.debug('Startup flows after filtering will be deleted -> %s',
            [flow.match for flow in startup_flows])
        self._remove_extra_flows(extra_flows)

        # For now just reinsert redirection rules, this is a bit of a hack but
        # redirection relies on async dns request to be setup and we can't
        # currently do this from out synchronous setup request. So just reinsert
        self._process_redirection_rules(requests)

        if self.proxy_controller_fut and self.proxy_controller_fut.done():
            if not self.proxy_controller:
                self.proxy_controller = self.proxy_controller_fut.result()

        self.logger.info("Initialized proxy_controller %s", self.proxy_controller)
        self.init_finished = True
        return SetupFlowsResult(result=SetupFlowsResult.SUCCESS)

    def _remove_extra_flows(self, extra_flows):
        msg_list = []
        for flow in extra_flows:
            if DIRECTION_REG in flow.match:
                direction = Direction(flow.match.get(DIRECTION_REG, None))
            else:
                direction = None
            match = MagmaMatch(imsi=flow.match.get(IMSI_REG, None),
                direction=direction)
            self.logger.debug('Sending msg for deletion -> %s',
                flow.match.get('reg1', None))
            msg_list.append(flows.get_delete_flow_msg(
                self._datapath, self.tbl_num, match, cookie=flow.cookie,
                cookie_mask=flows.OVS_COOKIE_MATCH_ALL))
        if msg_list:
            chan = self._msg_hub.send(msg_list, self._datapath)
            self._wait_for_responses(chan, len(msg_list))

    def _add_missing_flows(self, requests, current_flows):
        msg_list = []
        for add_flow_req in requests:
            imsi = add_flow_req.sid.id
            ip_addr = convert_ipv4_str_to_ip_proto(add_flow_req.ip_addr)
            apn_ambr = add_flow_req.apn_ambr
            static_rule_ids = add_flow_req.rule_ids
            dynamic_rules = add_flow_req.dynamic_rules
            msisdn = add_flow_req.msisdn
            uplink_tunnel = add_flow_req.uplink_tunnel

            msgs = self._get_default_flow_msgs_for_subscriber(imsi, ip_addr)
            if msgs:
                msg_list.extend(msgs)

            for rule_id in static_rule_ids:
                rule = self._policy_dict[rule_id]
                if rule is None:
                    self.logger.error("Could not find rule for rule_id: %s",
                        rule_id)
                    continue
                try:
                    if rule.redirect.support == rule.redirect.ENABLED:
                        continue
                    flow_adds = self._get_rule_match_flow_msgs(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule)
                    msg_list.extend(flow_adds)
                except FlowMatchError:
                    self.logger.error("Failed to verify rule_id: %s", rule_id)

            for rule in dynamic_rules:
                try:
                    if rule.redirect.support == rule.redirect.ENABLED:
                        continue
                    flow_adds = self._get_rule_match_flow_msgs(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule)
                    msg_list.extend(flow_adds)
                except FlowMatchError:
                    self.logger.error("Failed to verify rule_id: %s", rule.id)

        msgs_to_send, remaining_flows = \
            self._msg_hub.filter_msgs_if_not_in_flow_list(msg_list,
                                                          current_flows)
        if msgs_to_send:
            chan = self._msg_hub.send(msgs_to_send, self._datapath)
            self._wait_for_responses(chan, len(msgs_to_send))

        return remaining_flows

    def _process_redirection_rules(self, requests):
        for add_flow_req in requests:
            imsi = add_flow_req.sid.id
            ip_addr = convert_ipv4_str_to_ip_proto(add_flow_req.ip_addr)
            static_rule_ids = add_flow_req.rule_ids
            dynamic_rules = add_flow_req.dynamic_rules
            for rule_id in static_rule_ids:
                rule = self._policy_dict[rule_id]
                if rule is None:
                    self.logger.error("Could not find rule for rule_id: %s",
                        rule_id)
                    continue
                if rule.redirect.support == rule.redirect.ENABLED:
                    self._install_redirect_flow(imsi, ip_addr, rule)

            for rule in dynamic_rules:
                if rule.redirect.support == rule.redirect.ENABLED:
                    self._install_redirect_flow(imsi, ip_addr, rule)

    def activate_rules(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, static_rule_ids, dynamic_rules):
        """
        Activate the flows for a subscriber based on the rules stored in Redis.
        During activation, a default flow may be installed for the subscriber.

        Args:
            imsi (string): subscriber id
            msisdn (bytes): subscriber MSISDN
            uplink_tunnel(int): Tunnel ID of the subscriber session.
            ip_addr (string): subscriber session ipv4 address
            static_rule_ids (string []): list of static rules to activate
            dynamic_rules (PolicyRule []): list of dynamic rules to activate
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized for adding flows')
            return ActivateFlowsResult(
                static_rule_results=[RuleModResult(
                    rule_id=rule_id,
                    result=RuleModResult.FAILURE,
                ) for rule_id in static_rule_ids],
                dynamic_rule_results=[RuleModResult(
                    rule_id=rule.id,
                    result=RuleModResult.FAILURE,
                ) for rule in dynamic_rules],
            )
        static_results = []
        for rule_id in static_rule_ids:
            res = self._install_flow_for_static_rule(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule_id)
            static_results.append(RuleModResult(rule_id=rule_id, result=res))
        dyn_results = []
        for rule in dynamic_rules:
            res = self._install_flow_for_rule(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule)
            dyn_results.append(RuleModResult(rule_id=rule.id, result=res))

        # Install a base flow for when no rule is matched.
        self._install_default_flow_for_subscriber(imsi, ip_addr)
        return ActivateFlowsResult(
            static_rule_results=static_results,
            dynamic_rule_results=dyn_results,
        )

    def _remove_he_flows(self, ip_addr: IPAddress, rule_id: str = "",
                         rule_num: int = -1):
        if self.proxy_controller:
            self.proxy_controller.remove_subscriber_he_flows(ip_addr, rule_id,
                                                             rule_num)

    def _install_flow_for_static_rule(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, rule_id):
        """
        Install a flow to get stats for a particular static rule id. The rule
        will be loaded from Redis and installed.

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule_id (string): policy rule id
        """
        rule = self._policy_dict[rule_id]
        if rule is None:
            self.logger.error("Could not find rule for rule_id: %s", rule_id)
            return RuleModResult.FAILURE
        return self._install_flow_for_rule(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule)

    def _wait_for_rule_responses(self, imsi, ip_addr, rule, chan):
        def fail(err):
            self.logger.error(
                "Failed to install rule %s for subscriber %s: %s",
                rule.id, imsi, err)
            self._deactivate_flow_for_rule(imsi, ip_addr, rule.id)
            return RuleModResult.FAILURE

        for _ in range(len(rule.flow_list)):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS")
            if not result.ok():
                return fail(result.exception())
        return RuleModResult.SUCCESS

    def _wait_for_responses(self, chan, response_count):
        def fail(err):
            #TODO need to rework setup to return all rule specific success/fails
            self.logger.error("Failed to install rule for subscriber: %s", err)

        for _ in range(response_count):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS policy mixin")
            if not result.ok():
                return fail(result.exception())

    def _get_classify_rule_flow_msgs(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, flow, rule_num,
                                     priority, qos, hard_timeout, rule_id, app_name,
                                     app_service_type, next_table, version, qos_mgr,
                                     copy_table, urls:List[str] = None):
        """
        Install a flow from a rule. If the flow action is DENY, then the flow
        will drop the packet. Otherwise, the flow classifies the packet with
        its matched rule and injects the rule num into the packet's register.
        """
        parser = self._datapath.ofproto_parser
        flow_match = flow_match_to_magma_match(flow.match, ip_addr)
        flow_match.imsi = encode_imsi(imsi)
        flow_match_actions, instructions = self._get_action_for_rule(
            flow, rule_num, imsi, ip_addr, apn_ambr, qos, rule_id, version, qos_mgr)
        msgs = []
        if app_name:
            # We have to allow initial traffic to pass through, before it gets
            # classified by DPI, flow match set app_id to unclassified
            flow_match.app_id = UNCLASSIFIED_PROTO_ID
            passthrough_actions = flow_match_actions + \
                [parser.NXActionRegLoad2(dst=SCRATCH_REGS[1],
                                         value=IGNORE_STATS)]
            msgs.append(
                flows.get_add_resubmit_current_service_flow_msg(
                    self._datapath,
                    self.tbl_num,
                    flow_match,
                    passthrough_actions,
                    hard_timeout=hard_timeout,
                    priority=Utils.UNCLASSIFIED_ALLOW_PRIORITY,
                    cookie=rule_num,
                    copy_table=copy_table,
                    resubmit_table=next_table)
            )
            flow_match.app_id = get_app_id(
                PolicyRule.AppName.Name(app_name),
                PolicyRule.AppServiceType.Name(app_service_type),
            )

        # For DROP flow just send to stats table, it'll get dropped there
        if flow.action == flow.DENY:
            flow_match_actions = flow_match_actions + \
                [parser.NXActionRegLoad2(dst=SCRATCH_REGS[1],
                                         value=DROP_FLOW_STATS)]
            msgs.append(flows.get_add_resubmit_current_service_flow_msg(
                self._datapath,
                self.tbl_num,
                flow_match,
                flow_match_actions,
                hard_timeout=hard_timeout,
                priority=priority,
                cookie=rule_num,
                resubmit_table=copy_table)
            )
        else:
            msgs.append(flows.get_add_resubmit_current_service_flow_msg(
                self._datapath,
                self.tbl_num,
                flow_match,
                flow_match_actions,
                instructions=instructions,
                hard_timeout=hard_timeout,
                priority=priority,
                cookie=rule_num,
                copy_table=copy_table,
                resubmit_table=next_table)
            )

        if self.proxy_controller:
            ue_ip = ipv4_address_to_str(ip_addr)
            ip_dst = get_flow_ip_dst(flow.match)
            direction = get_direction_for_match(flow.match)

            proxy_msgs = self.proxy_controller.get_subscriber_he_flows(
                rule_id, direction, ue_ip, uplink_tunnel, ip_dst, rule_num,
                urls, imsi, msisdn)
            msgs.extend(proxy_msgs)
        return msgs

    def _get_action_for_rule(self, flow, rule_num, imsi, ip_addr,
                             apn_ambr, qos, rule_id, version, qos_mgr):
        """
        Returns an action instructions list to be applied for a specific flow.
        If qos or apn_ambr are set, the appropriate action is returned based
        on the implementation used (tc vs ovs_meter)
        """
        parser = self._datapath.ofproto_parser
        instructions = []

        # encode the rule id in hex
        of_note = parser.NXActionNote(list(rule_id.encode()))
        actions = [of_note]
        if flow.action == flow.DENY:
            return actions, instructions

        mbr_ul = qos.max_req_bw_ul
        mbr_dl = qos.max_req_bw_dl
        qos_info = None
        ambr = None
        d = flow.match.direction
        if d == flow.match.UPLINK:
            if apn_ambr:
                ambr = apn_ambr.max_bandwidth_ul
            if mbr_ul != 0:
                qos_info = QosInfo(gbr=qos.gbr_ul, mbr=mbr_ul)

        if d == flow.match.DOWNLINK:
            if apn_ambr:
                ambr = apn_ambr.max_bandwidth_dl
            if mbr_dl != 0:
                qos_info = QosInfo(gbr=qos.gbr_dl, mbr=mbr_dl)

        if qos_info or ambr:
            action, inst = qos_mgr.add_subscriber_qos(
                imsi, ip_addr.address.decode('utf8'), ambr, rule_num, d, qos_info)

            self.logger.debug("adding Actions %s instruction %s ", action, inst)
            if action:
                actions.append(action)

            if inst:
                instructions.append(inst)

        actions.extend(
            [parser.NXActionRegLoad2(dst=RULE_NUM_REG, value=rule_num),
             parser.NXActionRegLoad2(dst=RULE_VERSION_REG, value=version)
             ])
        return actions, instructions

    @abstractmethod
    def _install_flow_for_rule(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, rule):
        """
        Install a flow given a rule. Subclass should implement this.

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        raise NotImplementedError

    @abstractmethod
    def _install_default_flow_for_subscriber(self, imsi, ip_addr):
        """
        Install a flow for the subscriber in the event no rule is matched.
        Subclass should implement this.

        Args:
            imsi (string): subscriber id
        """
        raise NotImplementedError

    @abstractmethod
    def _install_redirect_flow(self, imsi, ip_addr, rule):
        """
        Install a redirection flow for the subscriber.
        Subclass should implement this.

        Args:
            imsi (string): subscriber id
            ip_addr (string): subscriber ip
            rule (PolicyRule): policyrule to install
        """
        raise NotImplementedError

    @abstractmethod
    def _install_default_flows_if_not_installed(self, datapath,
            existing_flows: List[OFPFlowStats]) -> List[OFPFlowStats]:
        """
        Install default flows(if not already installed), if no other flows are
        matched.

        Returns:
            The list of flows that remain after inserting default flows
        """
        raise NotImplementedError
