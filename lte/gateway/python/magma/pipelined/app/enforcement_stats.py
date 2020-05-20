"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from typing import List
from collections import defaultdict

from lte.protos.pipelined_pb2 import RuleModResult
from lte.protos.session_manager_pb2 import RuleRecord, \
    RuleRecordTable
from ryu.controller import dpset, ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib import hub
from ryu.ofproto.ofproto_v1_4 import OFPMPF_REPLY_MORE
from ryu.ofproto.ofproto_v1_4_parser import OFPFlowStats

from magma.pipelined.app.base import MagmaController, ControllerType, \
    global_epoch
from magma.pipelined.app.policy_mixin import PolicyMixin
from magma.pipelined.openflow import messages, flows
from magma.pipelined.openflow.exceptions import MagmaOFError
from magma.pipelined.imsi import decode_imsi, encode_imsi
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MsgChannel, MessageHub

from magma.pipelined.openflow.registers import Direction, DIRECTION_REG, \
    IMSI_REG, RULE_VERSION_REG, SCRATCH_REGS


ETH_FRAME_SIZE_BYTES = 14
PROCESS_STATS = 0x0
IGNORE_STATS = 0x1


class RelayDisabledException(Exception):
    pass


class EnforcementStatsController(PolicyMixin, MagmaController):
    """
    This openflow controller installs flows for aggregating policy usage
    statistics, which are sent to sessiond for tracking.

    It periodically polls OVS for flow stats on the its table and reports the
    usage records to session manager via RPC. Flows are deleted when their
    version (reg4 match) is different from the current version of the rule for
    the subscriber maintained by the rule version mapper.
    """

    APP_NAME = 'enforcement_stats'
    APP_TYPE = ControllerType.LOGICAL
    SESSIOND_RPC_TIMEOUT = 10
    # 0xffffffffffffffff is reserved in openflow
    DEFAULT_FLOW_COOKIE = 0xfffffffffffffffe

    _CONTEXTS = {
        'dpset': dpset.DPSet,
    }

    def __init__(self, *args, **kwargs):
        super(EnforcementStatsController, self).__init__(*args, **kwargs)
        # No need to report usage if relay mode is not enabled.
        self._relay_enabled = kwargs['mconfig'].relay_enabled
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self.dpset = kwargs['dpset']
        self.loop = kwargs['loop']
        # Spawn a thread to poll for flow stats
        poll_interval = kwargs['config']['enforcement']['poll_interval']
        # Create a rpc channel to sessiond
        self.sessiond = kwargs['rpc_stubs']['sessiond']
        self._msg_hub = MessageHub(self.logger)
        self.unhandled_stats_msgs = []  # Store multi-part responses from ovs
        self.total_usage = {}  # Store total usage
        # Store last usage excluding deleted flows for calculating deltas
        self.last_usage_for_delta = {}
        self.failed_usage = {}  # Store failed usage to retry rpc to sessiond
        self._unmatched_bytes = 0  # Store bytes matched by default rule if any
        self._clean_restart = kwargs['config']['clean_restart']
        if not self._relay_enabled:
            self.logger.info('Relay mode is not enabled. '
                             'enforcement_stats will not report usage.')
            return
        self.flow_stats_thread = hub.spawn(self._monitor, poll_interval)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def cleanup_state(self):
        """
        When we remove/reinsert flows we need to remove old usage maps as new
        flows will have reset stat counters
        """
        self.unhandled_stats_msgs = []
        self.total_usage = {}
        self.last_usage_for_delta = {}
        self.failed_usage = {}
        self._unmatched_bytes = 0

    def _check_relay(func):  # pylint: disable=no-self-argument
        def wrapped(self, *args, **kwargs):
            if self._relay_enabled:  # pylint: disable=protected-access
                func(self, *args, **kwargs)  # pylint: disable=not-callable

        return wrapped

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath
        if not self._relay_enabled:
            self._install_default_flows_if_not_installed(datapath, [])

    def _install_default_flows_if_not_installed(self, datapath,
            existing_flows: List[OFPFlowStats]) -> List[OFPFlowStats]:
        """
        Install default flows(if not already installed) to forward the traffic,
        If no other flows are matched.

        Returns:
            The list of flows that remain after inserting default flows
        """
        match = MagmaMatch()
        msg = flows.get_add_resubmit_next_service_flow_msg(
            datapath, self.tbl_num, match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_table,
            cookie=self.DEFAULT_FLOW_COOKIE)

        msg, remaining_flows = \
            self._msg_hub.filter_msgs_if_not_in_flow_list([msg], existing_flows)
        if msg:
            chan = self._msg_hub.send(msg, datapath)
            self._wait_for_responses(chan, 1)

        return remaining_flows

    @_check_relay
    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        if self._clean_restart:
            self.delete_all_flows(datapath)

    def _install_flow_for_rule(self, imsi, ip_addr, rule):
        """
        Install a flow to get stats for a particular rule. Flows will match on
        IMSI, cookie (the rule num), in/out direction

        Args:
            imsi (string): subscriber to install rule for
            ip_addr (string): subscriber session ipv4 address
            rule (PolicyRule): policy rule proto
        """
        # Do not install anything if relay is disabled
        if not self._relay_enabled:
            return RuleModResult.SUCCESS

        def fail(err):
            self.logger.error(
                "Failed to install rule %s for subscriber %s: %s",
                rule.id, imsi, err)
            return RuleModResult.FAILURE

        msgs = self._get_rule_match_flow_msgs(imsi, rule)

        chan = self._msg_hub.send(msgs, self._datapath)
        for _ in range(len(msgs)):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS")
            if not result.ok():
                return fail(result.exception())

        return RuleModResult.SUCCESS

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    @_check_relay
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    @_check_relay
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def _get_rule_match_flow_msgs(self, imsi, rule):
        """
        Returns flow add messages used for rule matching.
        """
        rule_num = self._rule_mapper.get_or_create_rule_num(rule.id)
        version = self._session_rule_version_mapper.get_version(imsi, rule.id)
        self.logger.debug(
            'Installing flow for %s with rule num %s (version %s)', imsi,
            rule_num, version)
        inbound_rule_match = _generate_rule_match(imsi, rule_num, version,
                                                  Direction.IN)
        outbound_rule_match = _generate_rule_match(imsi, rule_num, version,
                                                   Direction.OUT)

        inbound_rule_match._match_kwargs[SCRATCH_REGS[1]] = PROCESS_STATS
        outbound_rule_match._match_kwargs[SCRATCH_REGS[1]] = PROCESS_STATS
        msgs = [
            flows.get_add_resubmit_next_service_flow_msg(
                self._datapath,
                self.tbl_num,
                inbound_rule_match,
                [],
                priority=flows.DEFAULT_PRIORITY,
                cookie=rule_num,
                resubmit_table=self.next_table),
            flows.get_add_resubmit_next_service_flow_msg(
                self._datapath,
                self.tbl_num,
                outbound_rule_match,
                [],
                priority=flows.DEFAULT_PRIORITY,
                cookie=rule_num,
                resubmit_table=self.next_table),
        ]

        if rule.app_name:
            inbound_rule_match._match_kwargs[SCRATCH_REGS[1]] = IGNORE_STATS
            outbound_rule_match._match_kwargs[SCRATCH_REGS[1]] = IGNORE_STATS
            msgs.extend([
                flows.get_add_resubmit_next_service_flow_msg(
                    self._datapath,
                    self.tbl_num,
                    inbound_rule_match,
                    [],
                    priority=flows.DEFAULT_PRIORITY,
                    cookie=rule_num,
                    resubmit_table=self.next_table),
                flows.get_add_resubmit_next_service_flow_msg(
                    self._datapath,
                    self.tbl_num,
                    outbound_rule_match,
                    [],
                    priority=flows.DEFAULT_PRIORITY,
                    cookie=rule_num,
                    resubmit_table=self.next_table),
            ])
        return msgs

    def _get_default_flow_msg_for_subscriber(self, _):
        return None

    def _install_redirect_flow(self, imsi, ip_addr, rule):
        pass

    def _install_default_flow_for_subscriber(self, imsi):
        pass

    def get_policy_usage(self, fut):
        if not self._relay_enabled:
            fut.set_exception(RelayDisabledException())
            return

        record_table = RuleRecordTable(
            records=self.total_usage.values(),
            epoch=global_epoch)
        fut.set_result(record_table)

    def _monitor(self, poll_interval):
        """
        Main thread that sends a stats request at the configured interval in
        seconds.
        """
        while True:
            for _, datapath in self.dpset.get_all():
                if self.init_finished:
                    self._poll_stats(datapath)
                else:
                    # Still send an empty report -> needed for pipelined setup
                    self._report_usage({})
            hub.sleep(poll_interval)

    def _poll_stats(self, datapath):
        """
        Send a FlowStatsRequest message to the datapath
        """
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser
        req = parser.OFPFlowStatsRequest(
            datapath,
            table_id=self.tbl_num,
            out_group=ofproto.OFPG_ANY,
            out_port=ofproto.OFPP_ANY,
        )
        try:
            messages.send_msg(datapath, req)
        except MagmaOFError as e:
            self.logger.warning("Couldn't poll datapath stats: %s", e)

    @set_ev_cls(ofp_event.EventOFPFlowStatsReply, MAIN_DISPATCHER)
    @_check_relay
    def _flow_stats_reply_handler(self, ev):
        """
        Schedule the flow stats handling in the main event loop, so as to
        unblock the ryu event loop
        """
        if not self.init_finished:
            self.logger.debug('Setup not finished, skipping stats reply')
            return

        self.unhandled_stats_msgs.append(ev.msg.body)
        if ev.msg.flags == OFPMPF_REPLY_MORE:
            # Wait for more multi-part responses thats received for the
            # single stats request.
            return
        self.loop.call_soon_threadsafe(
            self._handle_flow_stats, self.unhandled_stats_msgs)
        self.unhandled_stats_msgs = []

    def _handle_flow_stats(self, stats_msgs):
        """
        Aggregate flow stats by rule, and report to session manager
        """
        stat_count = sum(len(flow_stats) for flow_stats in stats_msgs)
        if stat_count == 0:
            return

        self.logger.debug("Processing %s stats responses", len(stats_msgs))
        # Aggregate flows into rule records
        current_usage = defaultdict(RuleRecord)
        for flow_stats in stats_msgs:
            self.logger.debug("Processing stats of %d flows", len(flow_stats))
            for stat in flow_stats:
                if stat.table_id != self.tbl_num:
                    # this update is not intended for policy
                    return
                current_usage = self._update_usage_from_flow_stat(
                    current_usage, stat)

        # Calculate the delta values from last stat update
        delta_usage = _delta_usage_maps(current_usage,
                                        self.last_usage_for_delta)
        self.total_usage = current_usage

        # Append any records which we couldn't send to session manager earlier
        delta_usage = _merge_usage_maps(delta_usage, self.failed_usage)
        self.failed_usage = {}

        # Send report even if usage is empty. Sessiond uses empty reports to
        # recognize when flows have ended
        self._report_usage(delta_usage)

        self._delete_old_flows(stats_msgs)

    def _report_usage(self, delta_usage):
        """
        Report usage to sessiond using rpc
        """
        record_table = RuleRecordTable(records=delta_usage.values(),
                                       epoch=global_epoch)
        future = self.sessiond.ReportRuleStats.future(
            record_table, self.SESSIOND_RPC_TIMEOUT)
        future.add_done_callback(
            lambda future: self.loop.call_soon_threadsafe(
                self._report_usage_done, future, delta_usage))

    def _report_usage_done(self, future, delta_usage):
        """
        Callback after sessiond RPC completion
        """
        err = future.exception()
        if err:
            self.logger.error('Couldnt send flow records to sessiond: %s', err)
            self.failed_usage = _merge_usage_maps(
                delta_usage, self.failed_usage)

    def _update_usage_from_flow_stat(self, current_usage, flow_stat):
        """
        Update the rule record map with the flow stat and return the
        updated map.
        """
        rule_id = self._get_rule_id(flow_stat)
        # Rule not found, must be default flow
        if rule_id == "":
            default_flow_matched = \
                flow_stat.cookie == self.DEFAULT_FLOW_COOKIE and \
                flow_stat.byte_count != 0 and \
                self._unmatched_bytes != flow_stat.byte_count
            if default_flow_matched:
                self.logger.error('%s bytes total not reported.',
                                  flow_stat.byte_count)
                self._unmatched_bytes = flow_stat.byte_count
            return current_usage
        # If this is a pass through app name flow ignore stats
        if flow_stat.match[SCRATCH_REGS[1]] == IGNORE_STATS:
            return current_usage
        sid = _get_sid(flow_stat)

        # use a compound key to separate flows for the same rule but for
        # different subscribers
        key = sid + "|" + rule_id
        record = current_usage[key]
        record.rule_id = rule_id
        record.sid = sid
        if flow_stat.match[DIRECTION_REG] == Direction.IN:
            # HACK decrement byte count for downlink packets by the length
            # of an ethernet frame. Only IP and below should be counted towards
            # a user's data. Uplink does this already because the GTP port is
            # an L3 port.
            record.bytes_rx += _get_downlink_byte_count(flow_stat)
        else:
            record.bytes_tx += flow_stat.byte_count
        current_usage[key] = record
        return current_usage

    def _delete_old_flows(self, stats_msgs):
        """
        Check if the version of any flow is older than the current version. If
        so, delete the flow and update last_usage_for_delta so we calculate the
        correct usage delta for the next poll.
        """
        deleted_flow_usage = defaultdict(RuleRecord)
        for deletable_stat in self._old_flow_stats(stats_msgs):
            stat_rule_id = self._get_rule_id(deletable_stat)
            stat_sid = _get_sid(deletable_stat)
            rule_version = _get_version(deletable_stat)

            try:
                self._delete_flow(deletable_stat, stat_sid, rule_version)
                # Only remove the usage of the deleted flow if deletion
                # is successful.
                self._update_usage_from_flow_stat(deleted_flow_usage,
                                                  deletable_stat)
            except MagmaOFError as e:
                self.logger.error(
                    'Failed to delete rule %s for subscriber %s '
                    '(version: %s): %s', stat_rule_id,
                    stat_sid, rule_version, e)

        self.last_usage_for_delta = _delta_usage_maps(self.total_usage,
                                                      deleted_flow_usage)

    def _old_flow_stats(self, stats_msgs):
        """
        Generator function to filter the flow stats that should be deleted from
        the stats messages.
        """
        for flow_stats in stats_msgs:
            for stat in flow_stats:
                if stat.table_id != self.tbl_num:
                    # this update is not intended for policy
                    return

                rule_id = self._get_rule_id(stat)
                sid = _get_sid(stat)
                rule_version = _get_version(stat)
                if rule_id == "":
                    continue

                current_ver = \
                    self._session_rule_version_mapper.get_version(sid, rule_id)
                if current_ver != rule_version:
                    yield stat

    def _delete_flow(self, flow_stat, sid, version):
        cookie, mask = (
            flow_stat.cookie, flows.OVS_COOKIE_MATCH_ALL)
        match = _generate_rule_match(
            sid, flow_stat.cookie, version,
            Direction(flow_stat.match[DIRECTION_REG]))
        flows.delete_flow(self._datapath,
                          self.tbl_num,
                          match,
                          cookie=cookie,
                          cookie_mask=mask)

    def _get_rule_id(self, flow):
        """
        Return the rule id from the rule cookie
        """
        # the default rule will have a cookie of 0
        rule_num = flow.cookie
        if rule_num == 0 or rule_num == self.DEFAULT_FLOW_COOKIE:
            return ""
        try:
            return self._rule_mapper.get_rule_id(rule_num)
        except KeyError as e:
            self.logger.error('Could not find rule id for num %d: %s',
                              rule_num, e)
            return ""


def _generate_rule_match(imsi, rule_num, version, direction):
    """
    Return a MagmaMatch that matches on the rule num and the version.
    """
    return MagmaMatch(imsi=encode_imsi(imsi), direction=direction,
                      reg2=rule_num, rule_version=version)


def _delta_usage_maps(current_usage, last_usage):
    """
    Calculate the delta between the 2 usage maps and returns a new
    usage map.
    """
    if len(last_usage) == 0:
        return current_usage
    new_usage = {}
    for key, current in current_usage.items():
        last = last_usage.get(key, None)
        if last is not None:
            rec = RuleRecord()
            rec.MergeFrom(current)  # copy metadata
            rec.bytes_rx = current.bytes_rx - last.bytes_rx
            rec.bytes_tx = current.bytes_tx - last.bytes_tx
            new_usage[key] = rec
        else:
            new_usage[key] = current
    return new_usage


def _merge_usage_maps(current_usage, last_usage):
    """
    Merge the usage records from 2 map into a single map
    """
    if len(last_usage) == 0:
        return current_usage
    new_usage = {}
    for key, current in current_usage.items():
        last = last_usage.get(key, None)
        if last is not None:
            rec = RuleRecord()
            rec.MergeFrom(current)  # copy metadata
            rec.bytes_rx = current.bytes_rx + last.bytes_rx
            rec.bytes_tx = current.bytes_tx + last.bytes_tx
            new_usage[key] = rec
        else:
            new_usage[key] = current
    return new_usage


def _get_sid(flow):
    if IMSI_REG not in flow.match:
        return None
    return decode_imsi(flow.match[IMSI_REG])


def _get_version(flow):
    if RULE_VERSION_REG not in flow.match:
        return None
    return flow.match[RULE_VERSION_REG]


def _get_downlink_byte_count(flow_stat):
    total_bytes = flow_stat.byte_count
    packet_count = flow_stat.packet_count
    return total_bytes - ETH_FRAME_SIZE_BYTES * packet_count
