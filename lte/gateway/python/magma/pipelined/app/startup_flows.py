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

from magma.pipelined.app.base import (
    ControllerNotReadyException,
    ControllerType,
    MagmaController,
)
from magma.pipelined.openflow import messages
from magma.pipelined.openflow.exceptions import MagmaOFError
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib import hub
from ryu.ofproto.ofproto_v1_4 import OFPMPF_REPLY_MORE


class StartupFlows(MagmaController):
    """
    StartupFlows

    Factory class for querying startup flows from OVS. This controller is used
    to retrieve initial flows for all tables, other controllers can use this
    to properly initialize themselves. This is done to prevent each table from
    requesting ovs flows and flooding everything with messages.

    The StartupFlows contoller spawns a thread to poll the startup flows from
    all tables. Once all flows were received, the poll thread terminates.
    """

    APP_NAME = "startup_flows"
    APP_TYPE = ControllerType.SPECIAL
    POLL_INTERVAL = 1

    def __init__(self, *args, **kwargs):
        super(StartupFlows, self).__init__(*args, **kwargs)
        self._msg_xid = None
        self._datapath = None
        self._startup_flows = []
        self._table_flows = {}
        self._flows_received = False
        self._clean_restart = kwargs['config']['clean_restart']
        if self._clean_restart:
            self.logger.info('Clean restart enabled, startup flows will not '
                             'query flows.')
            self._flows_received = True
            return
        self._flow_stats_thread = hub.spawn(self._poll_startup_flows,
                                            self.POLL_INTERVAL)

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self._datapath = datapath

    def get_flows(self, tbl_num: int):
        """
        Get a list of flows for the specified table number

        Args:
            tbl_num: int
        """
        if not self._flows_received:
            raise ControllerNotReadyException("Setup flows not received yet")

        if tbl_num in self._table_flows:
            return self._table_flows[tbl_num]
        else:
            return []

    def delete_all_flows(self, datapath):
        pass

    @set_ev_cls(ofp_event.EventOFPFlowStatsReply, MAIN_DISPATCHER)
    def handle_startup_flows(self, ev):
        # This is not the message we're looking for
        if ev.msg.xid != self._msg_xid:
            return

        for resp in ev.msg.body:
            if resp.table_id not in self._table_flows:
                self._table_flows[resp.table_id] = []
            self._table_flows[resp.table_id].append(resp)

        # There will be more stats, we have to wait
        if ev.msg.flags == OFPMPF_REPLY_MORE:
            return

        self._flows_received = True

    def _poll_startup_flows(self, poll_interval):
        """
        Query stats until we receive startup flows
        """
        while not self._flows_received:
            if self._datapath:
                self._poll_all_tables(self._datapath)
            hub.sleep(poll_interval)

    def _poll_all_tables(self, datapath):
        """
        Send a FlowStatsRequest message to the datapath
        """
        ofproto, parser = datapath.ofproto, datapath.ofproto_parser
        req = parser.OFPFlowStatsRequest(
            datapath,
            out_group=ofproto.OFPG_ANY,
            out_port=ofproto.OFPP_ANY,
        )
        self._msg_xid = datapath.set_xid(req)
        try:
            messages.send_msg(datapath, req)
        except MagmaOFError as e:
            self.logger.warning("Couldn't poll datapath stats: %s", e)
