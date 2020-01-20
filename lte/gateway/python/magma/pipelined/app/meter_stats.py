"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from collections import defaultdict, namedtuple

from lte.protos.meteringd_pb2 import FlowRecord, \
    FlowTable
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.app.meter import MeterController
from magma.pipelined.imsi import decode_imsi
from magma.pipelined.openflow import messages
from magma.pipelined.openflow.exceptions import MagmaOFError
from magma.pipelined.openflow.registers import DIRECTION_REG, Direction, \
    IMSI_REG
from ryu.controller import dpset, ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib import hub


class UsageRecord(object):
    """
    A structure to capture usage stats of a sid.
    """

    def __init__(self):
        self.uuid = None
        self.sid = None
        self.bytes_tx = 0
        self.bytes_rx = 0
        self.pkts_tx = 0
        self.pkts_rx = 0

    def __str__(self):
        return ("bytes_tx: %d bytes_rx: %d pkts_tx: %d pkts_rx: %d"
                % (self.bytes_tx, self.bytes_rx, self.pkts_tx, self.pkts_rx))


class MeterStatsController(MagmaController):
    """
    This openflow controller periodically polls OVS for flow stats on the
    metering table and syncs the usage records to cloud via RPC.

    This controller is an entirely read-only controller, meaning that it will
    never push any flows.
    """

    APP_NAME = 'meter_stats'
    APP_TYPE = ControllerType.LOGICAL
    CLOUD_RPC_TIMEOUT = 10
    _CONTEXTS = {
        'dpset': dpset.DPSet,
    }

    MeterStatsConfig = namedtuple('MeterStatsConfig', ['poll_interval',
                                                       'enabled'])

    def __init__(self, *args, **kwargs):
        super(MeterStatsController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)

        self.dpset = kwargs['dpset']  # type: dpset.DPSet
        self.loop = kwargs['loop']
        self.config = self._get_config(kwargs['config'])

        # Cache of last reported flows
        self._last_reported_flows = {}

        # If metering isn't enabled, don't process any flow stats
        if not self.config.enabled:
            return

        # Spawn a thread to poll for flow stats
        self.flow_stats_thread = hub.spawn(self._monitor)

        self.meteringd_records = kwargs['rpc_stubs']['metering_cloud']

    def _get_config(self, config_dict):
        return self.MeterStatsConfig(
            poll_interval=config_dict['meter']['poll_interval'],
            enabled=config_dict['meter']['enabled']
        )

    def _monitor(self):
        """
        Main thread that sends a stats request at the configured interval in
        seconds.
        """
        # Don't enable polling when interval is negative
        if self.config.poll_interval < 0:
            return
        while True:
            for _, datapath in self.dpset.get_all():
                self._poll_stats(datapath)
            hub.sleep(self.config.poll_interval)

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
    def _flow_stats_reply(self, ev):
        """
        Schedule the flow stats handling in the main event loop, so as to
        unblock the ryu event loop
        """
        self.loop.call_soon_threadsafe(self._handle_flow_stats, ev.msg.body)

    def _handle_flow_stats(self, flow_stats):
        """
        Aggregate usage information when a datapath responds to a FlowStats
        request and sync over RPC with cloud metering service.
        """
        incoming, outgoing = 0, 1
        usage_by_sid = defaultdict(UsageRecord)

        def get_direction(match):
            inout_bit = match.get(DIRECTION_REG)
            if inout_bit == Direction.IN:
                return incoming
            elif inout_bit == Direction.OUT:
                return outgoing
            else:
                raise ValueError('No inout value found in match.')

        def get_flow_id(flow):
            actions = flow.instructions[0].actions
            if not hasattr(actions[0], 'note'):
                raise ValueError('No note found in flow stat actions.')
            # Filter out 0-bytes padded to the end of note
            return bytes(filter(lambda b: b, actions[0].note)).decode()

        def get_sid(match):
            encoded_imsi = match.get(IMSI_REG)
            if encoded_imsi is None or encoded_imsi == 0:
                raise ValueError('IMSI could not be parsed from match')
            return decode_imsi(encoded_imsi)

        for flow_stat in flow_stats:
            if flow_stat.table_id != self.tbl_num:
                # There is one global event handler for all tables. Ignore
                # requests from other tables or for default flows
                return
            if not self._should_report_flow(flow_stat):
                continue
            try:
                direction = get_direction(flow_stat.match)
                sid = get_sid(flow_stat.match)
                flow_id = get_flow_id(flow_stat)

                usage_by_sid[sid].id = flow_id
                usage_by_sid[sid].sid = sid
                if direction == incoming:
                    usage_by_sid[sid].bytes_rx = flow_stat.byte_count
                    usage_by_sid[sid].pkts_rx = flow_stat.packet_count
                else:
                    usage_by_sid[sid].bytes_tx = flow_stat.byte_count
                    usage_by_sid[sid].pkts_tx = flow_stat.packet_count
            except ValueError:
                continue
        self._sync_stats(usage_by_sid)

    def _should_report_flow(self, flow_stat):
        # Don't report the default flows
        return flow_stat.cookie != MeterController.DEFAULT_FLOW_COOKIE

    def _sync_stats(self, usage_by_sid):
        if len(usage_by_sid) == 0:
            self._last_reported_flows = {}
            return
        # Sync stats
        flow_records = _get_flow_records_from_usage(usage_by_sid.values())
        self._last_reported_flows = flow_records
        self.logger.debug('Syncing the following flow records to cloud:\n%s',
                          flow_records)
        future = self.meteringd_records.UpdateFlows.future(
            FlowTable(flows=flow_records),
            self.CLOUD_RPC_TIMEOUT)
        future.add_done_callback(
            lambda future: self.loop.call_soon_threadsafe(
                self._sync_stats_done, future))

    def _sync_stats_done(self, future):
        err = future.exception()
        if err:
            self.logger.error('Couldnt send flow records to cloud: %s', err)

    def get_subscriber_metering_flows(self, fut):
        fut.set_result(FlowTable(flows=self._last_reported_flows))


def _get_flow_records_from_usage(usage_records):
    """
    Transform a collection of UsageRecord objects into FlowRecord protobufs.

    Returns:
        [FlowRecord]:
            List of FlowRecord protobuf messages corresponding to the
            usages
    """

    def usage_to_flow_record(usage):
        return FlowRecord(
            id=FlowRecord.ID(id=usage.id),
            sid=usage.sid,
            bytes_tx=usage.bytes_tx,
            pkts_tx=usage.pkts_tx,
            bytes_rx=usage.bytes_rx,
            pkts_rx=usage.pkts_rx,
        )

    return list(map(usage_to_flow_record, usage_records))
