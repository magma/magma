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

from magma.pipelined.app.base import MagmaController
from magma.pipelined.openflow import flows, messages
from magma.pipelined.openflow.exceptions import MagmaOFError
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.ofctl_v1_4 import to_instructions


class TestingController(MagmaController):
    """
    TestingController

    The testing controller installs flows necessary for running unittests. It
    has access to all tables and is only used for testing purposes
    """

    APP_NAME = "testing"

    def __init__(self, *args, **kwargs):
        super(TestingController, self).__init__(*args, **kwargs)
        self._datapath = None
        self._stats_queue = None
        self._agr_stats_queue = None

    def initialize_on_connect(self, datapath):
        """
        Saves the datapath for future use

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath

    def cleanup_on_disconnect(self, datapath):
        """
        Flows should be deleted in test files, testing controller doesn't own
        a table so we simply pass.

        Args:
            datapath: ryu datapath struct
        """
        pass

    def insert_flow(self, ryu_req):
        """ Parse the ryu request and add flow to the ovs table """
        actions = to_instructions(self._datapath, ryu_req["instructions"])

        flows.add_drop_flow(
            self._datapath, ryu_req["table_id"], ryu_req["match"],
            instructions=actions,
            priority=ryu_req["priority"]
        )

    def delete_flow(self, ryu_req):
        """ Parse the ryu request and delete flow from ovs table """
        actions = to_instructions(self._datapath, ryu_req["instructions"])

        flows.delete_flow(
            self._datapath, ryu_req["table_id"], ryu_req["match"],
            instructions=actions,
            priority=ryu_req["priority"]
        )

    def ryu_query_lookup(self, ryu, stats_queue):
        """
        Send a FlowStatsRequest message to the datapath
        """
        self._stats_queue = stats_queue
        parser = self._datapath.ofproto_parser
        match = parser.OFPMatch(**ryu["match"].ryu_match) \
            if ryu["match"] is not None else None
        if "cookie" not in ryu:
            # If cookie is not set in the parameter, then do not match on
            # cookie.
            req = parser.OFPFlowStatsRequest(
                self._datapath, table_id=ryu["table_id"], match=match
            )
        else:
            req = parser.OFPFlowStatsRequest(
                self._datapath, table_id=ryu["table_id"], match=match,
                cookie=ryu["cookie"], cookie_mask=flows.OVS_COOKIE_MATCH_ALL
            )
        try:
            messages.send_msg(self._datapath, req)
        except MagmaOFError as e:
            self.logger.warning("Couldn't poll datapath stats: %s", e)

    def table_stats_lookup(self, queue):
        """
        Send a TableStatsRequest message to the datapath
        """
        self._agr_stats_queue = queue

        parser = self._datapath.ofproto_parser
        req = parser.OFPTableStatsRequest(self._datapath, 0)
        try:
            messages.send_msg(self._datapath, req)
        except MagmaOFError as e:
            self.logger.warning("Couldn't poll datapath stats: %s", e)

    @set_ev_cls(ofp_event.EventOFPFlowStatsReply, MAIN_DISPATCHER)
    def _flow_stats_reply_handler(self, ev):
        """ Save stats in the queue val to be accessible in the test file """
        if self._stats_queue is None:
            return
        flow_stats = ev.msg.body
        self._stats_queue.put(flow_stats)

    @set_ev_cls(ofp_event.EventOFPTableStatsReply, MAIN_DISPATCHER)
    def _table_stats_reply_handler(self, ev):
        """ Save stats in the queue val to be accessible in the test file """
        if self._agr_stats_queue is None:
            return
        flow_stats = ev.msg.body
        self._agr_stats_queue.put(flow_stats)
