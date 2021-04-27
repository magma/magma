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
from abc import ABCMeta, abstractmethod
from typing import List

from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import (
    ActivateFlowsRequest,
    ActivateFlowsResult,
    RuleModResult,
)
from lte.protos.policydb_pb2 import PolicyRule
from magma.pipelined.app.dpi import UNCLASSIFIED_PROTO_ID, get_app_id
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.messages import MsgChannel
from magma.pipelined.openflow.registers import (
    RULE_NUM_REG,
    RULE_VERSION_REG,
    SCRATCH_REGS,
)
from magma.pipelined.policy_converters import (
    FlowMatchError,
    convert_ipv4_str_to_ip_proto,
    convert_ipv6_bytes_to_ip_proto,
    flow_match_to_magma_match,
    get_direction_for_match,
    get_flow_ip_dst,
    ipv4_address_to_str,
)
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
        self._rule_mapper = kwargs['rule_id_mapper']
        self._setup_type = kwargs['config']['setup_type']
        self._session_rule_version_mapper = kwargs[
            'session_rule_version_mapper']
        if 'proxy' in kwargs['app_futures']:
            self.proxy_controller_fut = kwargs['app_futures']['proxy']
        else:
            self.proxy_controller_fut = None
        self.proxy_controller = None

    def activate_rules(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, policies):
        """
        Activate the flows for a subscriber based on the rules stored in Redis.
        During activation, a default flow may be installed for the subscriber.

        Args:
            imsi (string): subscriber id
            msisdn (bytes): subscriber MSISDN
            uplink_tunnel(int): Tunnel ID of the subscriber session.
            ip_addr (string): subscriber session ipv4 address
            policies (VersionedPolicies []): list of versioned policies to activate
        """
        if self._datapath is None:
            self.logger.error('Datapath not initialized for adding flows')
            return ActivateFlowsResult(
                policy_results=[RuleModResult(
                    rule_id=policy.rule.id,
                    version=policy.version,
                    result=RuleModResult.FAILURE,
                ) for policy in policies],
            )
        policy_results = []
        for policy in policies:
            res = self._install_flow_for_rule(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, policy.rule, policy.version)
            policy_results.append(RuleModResult(rule_id=policy.rule.id, version=policy.version, result=res))

        # Install a base flow for when no rule is matched.
        self._install_default_flow_for_subscriber(imsi, ip_addr, uplink_tunnel)
        return ActivateFlowsResult(
            policy_results=policy_results,
        )

    def _remove_he_flows(self, ip_addr: IPAddress, rule_id: str = "",
                         rule_num: int = -1):
        if self.proxy_controller:
            self.proxy_controller.remove_subscriber_he_flows(ip_addr, rule_id,
                                                             rule_num)

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
        flow_match = flow_match_to_magma_match(flow.match, ip_addr,
                                               uplink_tunnel)
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

    def _get_ue_specific_flow_msgs(self, requests: List[ActivateFlowsRequest]):
        msg_list = []
        for add_flow_req in requests:
            imsi = add_flow_req.sid.id
            apn_ambr = add_flow_req.apn_ambr
            policies = add_flow_req.policies
            msisdn = add_flow_req.msisdn
            uplink_tunnel = add_flow_req.uplink_tunnel

            if self._setup_type == 'CWF' or add_flow_req.ip_addr:
                ipv4 = convert_ipv4_str_to_ip_proto(add_flow_req.ip_addr)
                msgs = self._get_default_flow_msgs_for_subscriber(imsi, ipv4,
                    uplink_tunnel)
                if msgs:
                    msg_list.extend(msgs)

                for policy in policies:
                    msg_list.extend(self._get_policy_flows(imsi, msisdn, uplink_tunnel, ipv4, apn_ambr, policy))
            if add_flow_req.ipv6_addr:
                ipv6 = convert_ipv6_bytes_to_ip_proto(add_flow_req.ipv6_addr)
                msgs = self._get_default_flow_msgs_for_subscriber(imsi, ipv6,
                    uplink_tunnel)
                if msgs:
                    msg_list.extend(msgs)

                for policy in policies:
                    msg_list.extend(self._get_policy_flows(imsi, msisdn, uplink_tunnel, ipv6, apn_ambr, policy))

        return {self.tbl_num: msg_list}

    def _get_policy_flows(self, imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr,
                          policy):
        msg_list = []
        # As the versions are managed by sessiond, save state here
        self._service_manager.session_rule_version_mapper.save_version(
            imsi, uplink_tunnel, policy.rule.id, policy.version)
        try:
            if policy.rule.redirect.support == policy.rule.redirect.ENABLED:
                return msg_list
            flow_adds = self._get_rule_match_flow_msgs(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, policy.rule, policy.version)
            msg_list.extend(flow_adds)
        except FlowMatchError:
            self.logger.error("Failed to verify rule_id: %s", policy.rule.id)
        return msg_list

    def _process_redirection_rules(self, requests):
        for add_flow_req in requests:
            imsi = add_flow_req.sid.id
            ip_addr = convert_ipv4_str_to_ip_proto(add_flow_req.ip_addr)
            policies = add_flow_req.policies

            for policy in policies:
                if policy.rule.redirect.support == policy.rule.redirect.ENABLED:
                    self._install_redirect_flow(imsi, ip_addr, policy.rule, policy.version)

    def finish_init(self, requests):
        # For now just reinsert redirection rules, this is a bit of a hack but
        # redirection relies on async dns request to be setup and we can't
        # currently do this from out synchronous setup request. So just reinsert
        self._process_redirection_rules(requests)

        if self.proxy_controller_fut and self.proxy_controller_fut.done():
            if not self.proxy_controller:
                self.proxy_controller = self.proxy_controller_fut.result()
        self.logger.info("Initialized proxy_controller %s",
                         self.proxy_controller)


    @abstractmethod
    def _install_flow_for_rule(self, imsi, msisdn: bytes, uplink_tunnel: int, ip_addr, apn_ambr, rule, version):
        """
        Install a flow given a rule. Subclass should implement this.

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        raise NotImplementedError

    @abstractmethod
    def _install_default_flow_for_subscriber(self, imsi, ip_addr, uplink_tunnel):
        """
        Install a flow for the subscriber in the event no rule is matched.
        Subclass should implement this.

        Args:
            imsi (string): subscriber id
        """
        raise NotImplementedError

    @abstractmethod
    def _install_redirect_flow(self, imsi, ip_addr, rule, version):
        """
        Install a redirection flow for the subscriber.
        Subclass should implement this.

        Args:
            imsi (string): subscriber id
            ip_addr (string): subscriber ip
            rule (PolicyRule): policyrule to install
        """
        raise NotImplementedError
