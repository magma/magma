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
import logging
from typing import Any, List, Optional

# there's a cyclic dependency in ryu
import ryu.base.app_manager  # pylint: disable=unused-import
from magma.pipelined.metrics import DP_SEND_MSG_ERROR
from magma.pipelined.openflow.exceptions import (
    MagmaDPDisconnectedError,
    MagmaOFError,
)
from magma.pipelined.policy_converters import MATCH_ATTRIBUTES
from ryu.controller.controller import Datapath
from ryu.lib import hub
from ryu.ofproto.ofproto_parser import MsgBase

logger = logging.getLogger(__name__)

DEFAULT_TIMEOUT_SEC = 10


def send_msg(datapath, msg, retries=3):
    """
    Send a message to the given datapath, with retries on failure.
    Prefer the higher-level functions in flows.py before dropping down to
    send_msg.

    Args:
        datapath (ryu.controller.controller.Datapath): The datapath to send the
        message to
        msg: Openflow message to send
        retries: Number of retries on failure

    Raises:
        MagmaOFError: if the message fails to send in the specified number of
        attempts
    """
    for i in range(0, retries):
        try:
            ret = datapath.send_msg(msg)
            if not ret:
                raise MagmaDPDisconnectedError()
            return
        except Exception as e:  # pylint: disable=broad-except
            logger.warning(
                'Error sending message, retrying again '
                '(attempt %s/%s)', i, retries,
            )
            if i == retries - 1:    # Only propagate if all retries are up
                logging.error('Send msg error! Type: %s, Reason: %s',
                              type(e).__name__, e)
                DP_SEND_MSG_ERROR.labels(cause=type(e).__name__).inc()
                raise MagmaOFError(e)
            else:
                continue


class MsgReply(object):
    """
    Reply for a single message sent to OVS. If an exception occurs in the
    transaction, it is thrown when `result()` is called.
    """

    def __init__(self, txn_id: Any,
                 exception: Optional[Exception] = None) -> None:
        """
        Create a reply marked by the transaction id. If an error occurs,
        include an exception
        """
        self._exception = exception
        self.txn_id = txn_id

    def ok(self) -> bool:
        """
        Return true if no exception occured
        """
        return self._exception is None

    def exception(self) -> Optional[Exception]:
        return self._exception


class MsgChannel(object):
    """
    MsgChannel is a wrapper for a queue of message replies coming in
    asynchronously. MsgChannel now uses hub.Queue in order to not block any
    greenthreads.
    """
    class Timeout(Exception):
        pass

    def __init__(self) -> None:
        self._queue = hub.Queue()

    def get(self, timeout: int=DEFAULT_TIMEOUT_SEC) -> MsgReply:
        try:
            return self._queue.get(block=True, timeout=timeout)
        except hub.QueueEmpty:
            raise self.Timeout()

    def put(self, reply: MsgReply) -> None:
        self._queue.put(reply)


class MessageHub(object):
    """
    MessageHub can send flow modifications and and returns a channel
    to synchronously wait for any results (in the same vein as a Go channel).
    It does this by sending messages and barrier requests and responding with
    any errors that occured before the barrier response. The app can call
    `send` to send any message synchronously and wait for the result.

    IMPORTANT: To use the message hub, the app that contains it must call
    `handle_barrier` from all ofp_event.EventOFPBarrierReply events and
    `handle_error` from all ofp_event.EventOFPErrorMsg events.
    """
    def __init__(self, msg_hub_logger):
        self._switches = {}
        self.logger = msg_hub_logger

    def send(self,
             msg_list: List[MsgBase],
             datapath: Datapath,
             txn_id: Any=None,
             timeout: int=DEFAULT_TIMEOUT_SEC,
             channel: Optional[MsgChannel]=None) -> MsgChannel:
        """
        Send a message to OVS and track the result asynchronously. Multiple
        messages can be tracked using a transaction id (txn_id).

        Args:
            msg_list: list of messages to send
            datapath: datapath representing switch to send to
            txn_id: some kind of marker to track the messages sent
                (does not have to be unique)
            timeout: time before ignoring request
            channel: optional channel to use for the result. If it's not
                specified, one is created
        Returns:
            The channel passed in or the one created if it wasn't passed
        """
        switch = self._switches.get(datapath.id, None)

        if switch is None:
            # new switch to track
            switch = self._SwitchInfo()
            self._switches[datapath.id] = switch

        if channel is None:
            channel = MsgChannel()

        # set xids in all msgs
        msg_xids = [datapath.set_xid(msg) for msg in msg_list]
        req = self._MsgRequest(txn_id, msg_xids, channel)

        barrier = datapath.ofproto_parser.OFPBarrierRequest(datapath)
        datapath.set_xid(barrier)  # sets xid in the barrier

        switch.requests_by_barrier[barrier.xid] = req
        for msg in msg_list:
            switch.results_by_msg[msg.xid] = None
            ret = datapath.send_msg(msg)
            if not ret:
                raise MagmaDPDisconnectedError()
        datapath.send_msg(barrier)
        req.set_timeout(timeout, switch, barrier.xid, msg_xids)
        return channel

    def filter_msgs_if_not_in_flow_list(self,
                                        dp: Datapath,
                                        msg_list: List[MsgBase],
                                        flow_list):
        """
        Returns a list of messages not found in the provided flow_list, also
        returns a list of remaining flows(not found in the msg_list)
        """
        msgs_to_send = []
        remaining_flows = flow_list.copy()
        for msg in msg_list:
            index = self._get_msg_index_in_flow_list(dp, msg, remaining_flows)
            if index >= 0:
                remaining_flows.pop(index)
            else:
                msgs_to_send.append(msg)
        return msgs_to_send, remaining_flows

    @staticmethod
    def _respond(request, reply):
        if request.channel is not None:
            request.channel.put(reply)

    def handle_barrier(self, ev):
        """
        Barrier means that all of the messages before this barrier has been
        process. This function should be attached to the EventOFPBarrierReply
        event
        """
        msg = ev.msg
        datapath = msg.datapath
        switch = self._switches.get(datapath.id, None)
        if switch is None:
            # could be from a different application
            return
        req = switch.requests_by_barrier.pop(msg.xid, None)
        if req is None:
            # could be from a different application
            return
        req.cancel_timeout()
        for xid in req.msg_xids:
            e = switch.results_by_msg.pop(xid, None)
            MessageHub._respond(
                req,
                MsgReply(txn_id=req.txn_id, exception=e),
            )

    def handle_error(self, ev):
        """
        Error means that a message sent to OVS is unsuccessful. This function
        should be attached to the EventOFPErrorMsg event
        """
        msg = ev.msg
        datapath = msg.datapath
        switch = self._switches.get(datapath.id, None)
        if switch is None:
            self.logger.error('unknown dpid %s', datapath.id)
            return
        if msg.xid not in switch.results_by_msg:
            return
        # for now, result is unused. Just return if there's an exception
        switch.results_by_msg[msg.xid] = MagmaOFError(ev.msg)

    def _flow_matches_flowmsg(self, dp, flow, msg):
        """
        Compare the flow and flow message based on match/instructions
        """
        reg_loads_match = True
        resubmits_match = True
        outputs_match = True
        if len(flow.instructions) != len(msg.instructions):
            return False
        for j in range(0, len(flow.instructions)):
            # TODO add support for OFPInstructionMeter and others
            if type(flow.instructions[j]) != dp.ofproto_parser.OFPInstructionActions:
                continue
            if type(msg.instructions[j]) != dp.ofproto_parser.OFPInstructionActions:
                continue
            # Strip _nxm to handle nicira as eth_dst_nxm is same as eth_dst
            reg_loads_flow = {i.dst.replace('_nxm', ''): i.value for i in flow.instructions[j].actions
                              if type(i) == dp.ofproto_parser.NXActionRegLoad2}
            reg_loads_msg = {i.dst.replace('_nxm', ''): i.value for i in msg.instructions[j].actions
                             if type(i) == dp.ofproto_parser.NXActionRegLoad2}

            reg_loads_match = reg_loads_msg == reg_loads_flow

            resubmits_flow = [i.table_id for i in flow.instructions[j].actions
                              if type(i) == dp.ofproto_parser.NXActionResubmitTable]
            resubmits_msg = [i.table_id for i in msg.instructions[j].actions
                             if type(i) == dp.ofproto_parser.NXActionResubmitTable]
            resubmits_match = sorted(resubmits_flow) == sorted(resubmits_msg)

            outputs_flow = [i.port for i in flow.instructions[j].actions
                            if type(i) == dp.ofproto_parser.OFPActionOutput]
            outputs_msg = [i.port for i in msg.instructions[j].actions
                           if type(i) == dp.ofproto_parser.OFPActionOutput]
            outputs_match = sorted(outputs_flow) == sorted(outputs_msg)

        match_flow = {key: flow.match.get(key) for key in MATCH_ATTRIBUTES
                      if key in flow.match}
        match_msg = {key: msg.match.get(key, None) for key in MATCH_ATTRIBUTES
                     if key in msg.match}

        def strip_common(match_dict):
            """
            Filter out args that break equality
                - ('0.0.0.0', '0.0.0.0') is the same as unset
            """
            for key in list(match_dict.keys()):
                if match_dict[key] == ('0.0.0.0', '0.0.0.0'):
                    del match_dict[key]
            return match_dict
        flow_match = strip_common(match_flow) == strip_common(match_msg)

        return flow_match and reg_loads_match and resubmits_match and \
               outputs_match

    def _get_msg_index_in_flow_list(self, dp, msg, flow_list):
        for i in range(len(flow_list)):
            if self._flow_matches_flowmsg(dp, msg, flow_list[i]):
                return i
        return -1

    class _MsgRequest(object):
        def __init__(self, txn_id, msg_xids, channel=None):
            self.txn_id = txn_id
            self.msg_xids = msg_xids
            self.channel = channel
            self._timeout_thread = None

        def set_timeout(self, timeout_sec, switch, barrier_xid, msg_xids):
            """
            Spawn a timeout handler after timeout_sec seconds to clear up any
            associated state with the request
            """
            def _handle_timeout():
                return self._handle_timeout(switch, barrier_xid, msg_xids)
            # spawn timeout func to ensure cleanup occurs
            self._timeout_thread = hub.spawn_after(timeout_sec,
                                                   _handle_timeout)

        def cancel_timeout(self):
            """
            If a request is received, stop the timeout handler from running
            """
            if self._timeout_thread is not None:
                self._timeout_thread.cancel()

        def _handle_timeout(self, switch, barrier_xid, msg_xids):
            switch.requests_by_barrier.pop(barrier_xid, None)
            for xid in msg_xids:
                switch.results_by_msg.pop(xid, None)

    class _SwitchInfo(object):
        def __init__(self):
            self.requests_by_barrier = {}
            self.results_by_msg = {}
