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
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.inout import EGRESS
from magma.pipelined.app.policy_mixin import PolicyMixin
from magma.pipelined.app.restart_mixin import DefaultMsgsMap, RestartMixin
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub
from magma.pipelined.policy_converters import FlowMatchError
from magma.pipelined.qos.common import QosManager
from magma.pipelined.qos.qos_meter_impl import MeterManager
from magma.pipelined.redirect import RedirectException, RedirectionManager
from magma.pipelined.utils import Utils
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls


class GYController(PolicyMixin, RestartMixin, MagmaController):
    """
    GYController

    The GY controller installs flows for enforcement of GY final actions, this
    includes redirection and QoS(currently not supported)
    """
    APP_NAME = "gy"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(GYController, self).__init__(*args, **kwargs)
        self._config = kwargs['config']
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_main_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.next_service_table = self._service_manager.get_next_table_num(
            EnforcementStatsController.APP_NAME)
        self._enforcement_stats_tbl = self._service_manager.get_table_num(
            EnforcementStatsController.APP_NAME)
        self.loop = kwargs['loop']
        self._msg_hub = MessageHub(self.logger)
        self._internal_ip_allocator = kwargs['internal_ip_allocator']
        self._redirect_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 2)[0]
        self._mac_rewr = \
            self._service_manager.INTERNAL_MAC_IP_REWRITE_TBL_NUM
        self._bridge_ip_address = kwargs['config']['bridge_ip_address']
        self._clean_restart = kwargs['config']['clean_restart']
        self._qos_mgr = None
        self._setup_type = self._config['setup_type']
        self._redirect_manager = \
            RedirectionManager(
                self._bridge_ip_address,
                self.logger,
                self.tbl_num,
                self._enforcement_stats_tbl,
                self._service_manager.get_table_num(EGRESS),
                self._redirect_scratch,
                self._session_rule_version_mapper
            )
        if self._setup_type == 'CWF':
            self._redirect_manager.set_cwf_args(
                internal_ip_allocator=kwargs['internal_ip_allocator'],
                arp=kwargs['app_futures']['arpd'],
                mac_rewrite=self._mac_rewr,
                bridge_name=kwargs['config']['bridge_name'],
                egress_table=self._service_manager.get_table_num(EGRESS)
            )

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        self._qos_mgr = QosManager.get_qos_manager(datapath, self.loop, self._config)

    def deactivate_rules(self, imsi, ip_addr, uplink_tunnel, rule_ids):
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
            self._deactivate_flows_for_subscriber(imsi, ip_addr, uplink_tunnel)
        else:
            for rule_id in rule_ids:
                self._deactivate_flow_for_rule(imsi, ip_addr, uplink_tunnel,
                                               rule_id)

    def cleanup_state(self):
        pass

    # pylint: disable=unused-argument
    def _deactivate_flows_for_subscriber(self, imsi, ip_addr, uplink_tunnel):
        """
        Deactivate all rules for a subscriber, ending any enforcement

        Args:
            imsi (string): subscriber id
            ip_addr(IPAddress): session IP address
        """
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.tbl_num, match)
        self._redirect_manager.deactivate_flows_for_subscriber(self._datapath,
                                                               imsi)
        self._qos_mgr.remove_subscriber_qos(imsi)
        self._remove_he_flows(ip_addr, None)

    # pylint: disable=unused-argument
    def _deactivate_flow_for_rule(self, imsi, ip_addr, uplink_tunnel,
                                  rule_id):
        """
        Deactivate a specific rule using the flow cookie for a subscriber

        Args:
            imsi (string): subscriber id
            rule_id (string): policy rule id
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
        self._qos_mgr.remove_subscriber_qos(imsi, num)
        self._remove_he_flows(ip_addr, rule_id)

    def _install_flow_for_rule(self, imsi, msisdn:bytes, uplink_tunnel: int,
                               ip_addr, apn_ambr, rule, version):
        """
        Install a flow to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            apn_ambr (integer): maximum bandwidth for non-GBR EPS bearers
            rule (PolicyRule): policy rule proto
        """
        if rule.redirect.support == rule.redirect.ENABLED:
            self._install_redirect_flow(imsi, ip_addr, rule, version)
            return RuleModResult.SUCCESS

        if not rule.flow_list:
            self.logger.error('The flow list for imsi %s, rule.id - %s'
                              'is empty, this shoudn\'t happen', imsi, rule.id)
            return RuleModResult.FAILURE

        flow_adds = []
        try:
            flow_adds = self._get_rule_match_flow_msgs(imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, rule, version)
        except FlowMatchError:
            return RuleModResult.FAILURE

        chan = self._msg_hub.send(flow_adds, self._datapath)
        return self._wait_for_rule_responses(imsi, ip_addr, rule, chan)

    def _get_default_flow_msgs_for_subscriber(self, *_):
        return None

    def _install_default_flow_for_subscriber(self, *_):
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
        match = MagmaMatch()
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

    def _install_redirect_flow(self, imsi, ip_addr, rule, version):
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        # CWF generates an internal IP for redirection so ip_addr is not needed
        if self._setup_type == 'CWF':
            ip_addr_str = None
        elif ip_addr and ip_addr.address:
            ip_addr_str = ip_addr.address.decode('utf-8')
        priority = rule.priority
        # TODO currently if redirection is enabled we ignore other flows
        # from rule.flow_list, confirm that this is the expected behaviour
        redirect_request = RedirectionManager.RedirectRequest(
            imsi=imsi,
            ip_addr=ip_addr_str,
            rule=rule,
            rule_num=rule_num,
            rule_version=version,
            priority=priority)
        try:
            if self._setup_type == 'CWF':
                self._redirect_manager.setup_cwf_redirect(
                    self._datapath, self.loop, redirect_request)
            else:
                self._redirect_manager.setup_lte_redirect(
                    self._datapath, self.loop, redirect_request)
            return RuleModResult.SUCCESS
        except RedirectException as err:
            self.logger.error(
                'Redirect Exception for imsi %s, rule.id - %s : %s',
                imsi, rule.id, err
            )
            return RuleModResult.FAILURE

    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        """
        Gets the default flow msg that forwards to next service

        Args:
            datapath: ryu datapath struct
        Returns:
            The list of default msgs to add
        """
        match = MagmaMatch()
        msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_main_table)

        return {self.tbl_num: [msg]}

    def _get_rule_match_flow_msgs(self, imsi, msisdn: bytes, uplink_tunnel: int,
                                  ip_addr, apn_ambr, rule, version):
        """
        Get flow msgs to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            msisdn (bytes): subscriber ISDN
            ip_addr (string): subscriber session ipv4 address
            apn_ambr (integer): maximum bandwidth for non-GBR EPS bearers
            rule (PolicyRule): policy rule proto
        """
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        priority = Utils.get_of_priority(rule.priority)

        flow_adds = []
        for flow in rule.flow_list:
            try:
                flow_adds.extend(self._get_classify_rule_flow_msgs(
                    imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, flow, rule_num, priority,
                    rule.qos, rule.hard_timeout, rule.id, rule.app_name,
                    rule.app_service_type, self.next_service_table,
                    version, self._qos_mgr, self._enforcement_stats_tbl))

            except FlowMatchError as err:  # invalid match
                self.logger.error(
                    "Failed to get flow msg '%s' for subscriber %s: %s",
                    rule.id, imsi, err)
                raise err
        return flow_adds

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

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def recover_state(self, _):
        pass
