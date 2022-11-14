"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from collections import namedtuple

from magma.pipelined.app.base import MagmaController
from magma.pipelined.app.egress import EGRESS
from magma.pipelined.app.restart_mixin import DefaultMsgsMap, RestartMixin
from magma.pipelined.ifaces import get_mac_address_from_iface
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub, MsgChannel
from magma.pipelined.openflow.registers import PASSTHROUGH_REG_VAL, Direction
from magma.pipelined.vlan_utils import get_vlan_egress_flow_msgs
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.lib.packet import ether_types

PHYSICAL_TO_LOGICAL = "middle"


class MiddleController(RestartMixin, MagmaController):
    APP_NAME = "middle"

    MiddleConfig = namedtuple(
        'MiddleConfig',
        [
            'mtr_ip', 'mtr_port', 'li_port_name', 'setup_type',
            'mtr_mac',
        ],
    )

    def __init__(self, *args, **kwargs):
        super(MiddleController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.logger.info("middle config: %s", self.config)

        # TODO Alex do we want this to be cofigurable from swagger?
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False

        self._midle_tbl_num = self._service_manager.get_table_num(PHYSICAL_TO_LOGICAL)
        self._egress_tbl_num = self._service_manager.get_table_num(EGRESS)
        # following fields are only used in Non Nat config
        self._clean_restart = kwargs['config']['clean_restart']
        self._msg_hub = MessageHub(self.logger)
        self._datapath = None
        self.tbl_num = self._midle_tbl_num

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None

        if 'mtr_ip' in config_dict and 'mtr_interface' in config_dict and 'ovs_mtr_port_number' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']
            mtr_mac = get_mac_address_from_iface(config_dict['mtr_interface'])
        else:
            mtr_ip = None
            mtr_mac = None
            mtr_port = None

        li_port_name = None
        if 'li_local_iface' in config_dict:
            li_port_name = config_dict['li_local_iface']

        return self.MiddleConfig(
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            li_port_name=li_port_name,
            setup_type=config_dict.get('setup_type', None),
            mtr_mac=mtr_mac,
        )

    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        return {
            self._midle_tbl_num: self._get_default_middle_flow_msgs(datapath),
        }

    def _get_default_middle_flow_msgs(self, dp):
        """
        Egress table is the last table that a packet touches in the pipeline.
        Output downlink traffic to gtp port, uplink traffic to LOCAL

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        msgs = []
        next_tbl = self._service_manager.get_next_table_num(PHYSICAL_TO_LOGICAL)

        # Allow passthrough pkts(skip enforcement and send to egress table)
        ps_match = MagmaMatch(passthrough=PASSTHROUGH_REG_VAL)
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp,
                self._midle_tbl_num, ps_match, actions=[],
                priority=flows.PASSTHROUGH_PRIORITY,
                resubmit_table=self._egress_tbl_num,
            ),
        )

        match = MagmaMatch()
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp,
                self._midle_tbl_num, match, actions=[],
                priority=flows.DEFAULT_PRIORITY, resubmit_table=next_tbl,
            ),
        )

        if self._mtr_service_enabled:
            msgs.extend(
                get_vlan_egress_flow_msgs(
                    dp,
                    self._midle_tbl_num,
                    ether_types.ETH_TYPE_IP,
                    self.config.mtr_ip,
                    self.config.mtr_port,
                    priority=flows.UE_FLOW_PRIORITY,
                    direction=Direction.OUT,
                    dst_mac=self.config.mtr_mac,
                ),
            )
        return msgs

    def _wait_for_responses(self, chan, response_count):
        def fail(err):
            self.logger.error("Failed to install rule with error: %s", err)

        for _ in range(response_count):
            try:
                result = chan.get()
            except MsgChannel.Timeout:
                return fail("No response from OVS msg channel")
            if not result.ok():
                return fail(result.exception())

    def cleanup_on_disconnect(self, datapath):
        if self._clean_restart:
            self.delete_all_flows(datapath)

    @set_ev_cls(ofp_event.EventOFPBarrierReply, MAIN_DISPATCHER)
    def _handle_barrier(self, ev):
        self._msg_hub.handle_barrier(ev)

    @set_ev_cls(ofp_event.EventOFPErrorMsg, MAIN_DISPATCHER)
    def _handle_error(self, ev):
        self._msg_hub.handle_error(ev)

    def cleanup_state(self):
        pass

    def _get_ue_specific_flow_msgs(self, _):
        return {}

    def finish_init(self, _):
        pass

    def initialize_on_connect(self, datapath):
        self._datapath = datapath

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
