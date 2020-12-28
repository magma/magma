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

from lte.protos.pipelined_pb2 import SetupFlowsResult, ActivateFlowsRequest
from magma.pipelined.app.base import ControllerNotReadyException
from magma.pipelined.openflow import flows
from magma.policydb.rule_store import PolicyRuleDict

from magma.pipelined.policy_converters import FlowMatchError, \
    convert_ipv4_str_to_ip_proto, ovs_flow_match_to_magma_match


class RestartMixin(metaclass=ABCMeta):
    """
    RestartMixin

    Mixin class for policy enforcement apps that includes common methods
    used for rule activation/deactivation.
    """
    def __init__(self, *args, **kwargs):
        super(RestartMixin, self).__init__(*args, **kwargs)
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
        Sets up the controller after the restart
         - Add default/missing flows if needed
         - Remove stale flows (not default and not in passed requsts)
        """
        if not self._datapath:
            self.logger.error('Controller restart not ready, datapath is None')
            return SetupFlowsResult.FAILURE
        if requests is None:
            requests = []
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
        if not hasattr(self, '_tbls'):
            self._tbls = [self.tbl_num]
        try:
            startup_flows = \
                {i: self._startup_flow_controller.get_flows(i) for i
                 in self._tbls}
        except ControllerNotReadyException as err:
            self.logger.error('Setup failed: %s', err)
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

        self.logger.debug('Setting up %s default rules', self.APP_NAME)
        remaining_flows = self._install_default_flows_if_not_installed(
            self._datapath, startup_flows)

        for tbl in startup_flows:
            self.logger.debug('Startup flows before filtering -> %s',
                              [flow.match for flow in startup_flows[tbl]])
        extra_flows = self._add_missing_flows(requests, remaining_flows)

        for tbl in startup_flows:
            self.logger.debug('Startup extra flows that will be deleted -> %s',
                              [flow.match for flow in startup_flows[tbl]])
        self._remove_extra_flows(extra_flows)

        # For now just reinsert redirection rules, this is a bit of a hack but
        # redirection relies on async dns request to be setup and we can't
        # currently do this from out synchronous setup request. So just reinsert
        self._process_redirection_rules(requests)

        # TODO I don't think this is relevant here, move to specific controller
        if self.proxy_controller_fut and self.proxy_controller_fut.done():
            if not self.proxy_controller:
                self.proxy_controller = self.proxy_controller_fut.result()
        self.logger.info("Initialized proxy_controller %s",
                         self.proxy_controller)

        self.init_finished = True
        return SetupFlowsResult(result=SetupFlowsResult.SUCCESS)

    def _remove_extra_flows(self, extra_flows):
        msg_list = []
        for tbl in extra_flows:
            for flow in extra_flows[tbl]:
                match = ovs_flow_match_to_magma_match(flow)
                self.logger.debug('Sending msg for deletion -> %s',
                                  match.ryu_match)
                msg_list.append(flows.get_delete_flow_msg(
                    self._datapath, tbl, match, cookie=flow.cookie,
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

        ret = {}
        for tbl in current_flows:
            msgs_to_send, remaining_flows = \
                self._msg_hub.filter_msgs_if_not_in_flow_list(
                    self._datapath, msg_list, current_flows[tbl])
            ret[tbl] = remaining_flows
        if msgs_to_send:
            chan = self._msg_hub.send(msgs_to_send, self._datapath)
            self._wait_for_responses(chan, len(msgs_to_send))

        return ret

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
