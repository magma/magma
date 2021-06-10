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
from typing import Dict, List

from lte.protos.pipelined_pb2 import SetupFlowsResult
from magma.pipelined.app.base import ControllerNotReadyException
from magma.pipelined.openflow import flows
from magma.pipelined.policy_converters import ovs_flow_match_to_magma_match
from ryu.ofproto.ofproto_v1_4_parser import OFPFlowStats

DefaultMsgsMap = Dict[int, List[OFPFlowStats]]


class RestartMixin(metaclass=ABCMeta):
    """
    RestartMixin

    Mixin class for controller restart handling
    """

    def handle_restart(self, requests) -> SetupFlowsResult:
        """
        Sets up the controller after the restart
         - Add default/missing flows if needed
         - Remove stale flows (not default and not in passed requsts)
         requests argument is controller specific
        """
        if not self._datapath:
            self.logger.error('Controller restart not ready, datapath is None')
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)
        dp = self._datapath

        if requests is None:
            requests = []

        if self._clean_restart:
            self.delete_all_flows(dp)
            self.cleanup_state()
            self.logger.info('Controller is in clean restart mode, remaining '
                             'flows were removed, continuing with setup.')

        if self._startup_flow_controller is None:
            if (self._startup_flows_fut.done()):
                self._startup_flow_controller = self._startup_flows_fut.result()
            else:
                self.logger.error('Flow Startup controller is not ready')
                return SetupFlowsResult(result=SetupFlowsResult.FAILURE)
        # Workaround for controllers with multiple tables
        if not hasattr(self, '_tbls'):
            self._tbls = [self.tbl_num]
        try:
            startup_flows_map = \
                {i: self._startup_flow_controller.get_flows(i) for i
                 in self._tbls}
        except ControllerNotReadyException as err:
            self.logger.error('Setup failed: %s', err)
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

        for tbl in startup_flows_map:
            self.logger.debug('Startup flows before filtering: tbl %d-> %s',
                              tbl,
                              [flow.match for flow in startup_flows_map[tbl]])

        default_msgs = self._get_default_flow_msgs(dp)
        for table, msgs_to_install in default_msgs.items():
            msgs, remaining_flows = self._msg_hub \
                .filter_msgs_if_not_in_flow_list(dp, msgs_to_install,
                                                 startup_flows_map[table])
            if msgs:
                chan = self._msg_hub.send(msgs, dp)
                self._wait_for_responses(chan, len(msgs))
            startup_flows_map[table] = remaining_flows

        ue_msgs = self._get_ue_specific_flow_msgs(requests)
        for table, msgs_to_install in ue_msgs.items():
            msgs, remaining_flows = self._msg_hub \
                .filter_msgs_if_not_in_flow_list(dp, msgs_to_install,
                                                 startup_flows_map[table])
            if msgs:
                chan = self._msg_hub.send(msgs, dp)
                self._wait_for_responses(chan, len(msgs))
            startup_flows_map[table] = remaining_flows

        for tbl in startup_flows_map:
            self.logger.debug('Startup flows to be deleted: tbl %d -> %s',
                              tbl,
                              [flow.match for flow in startup_flows_map[tbl]])
        self._remove_extra_flows(startup_flows_map)

        self.finish_init(requests)
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

    @abstractmethod
    def _get_ue_specific_flow_msgs(self, requests):
        """
        Gets ue flow messages for controller

        Args:
            requests: Controller specific setup information
        """
        raise NotImplementedError

    @abstractmethod
    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        """
        Gets default flow messages for controller

        Args:
            datapath (Datapath): RYU datapath
        """
        raise NotImplementedError
