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
from magma.pipelined.app.li_mirror import LIMirrorController
from magma.pipelined.app.restart_mixin import DefaultMsgsMap, RestartMixin
from magma.pipelined.bridge_util import BridgeTools, DatapathLookupError
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.messages import MessageHub, MsgChannel
from magma.pipelined.openflow.registers import (
    PASSTHROUGH_REG_VAL,
    REG_ZERO_VAL,
    Direction,
    load_direction,
)
from ryu.controller import ofp_event
from ryu.controller.handler import MAIN_DISPATCHER, set_ev_cls
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL

INGRESS = "ingress"


class IngressController(RestartMixin, MagmaController):
    APP_NAME = "ingress"

    IngressConfig = namedtuple(
        'IngressConfig',
        [
            'uplink_port', 'mtr_ip', 'mtr_port', 'li_port_name',
            'setup_type', 'he_proxy_port',
        ],
    )

    def __init__(self, *args, **kwargs):
        super(IngressController, self).__init__(*args, **kwargs)
        self.config = self._get_config(kwargs['config'])
        self.logger.info("ingress config: %s", self.config)

        self._li_port = None
        # TODO Alex do we want this to be cofigurable from swagger?
        if self.config.mtr_ip:
            self._mtr_service_enabled = True
        else:
            self._mtr_service_enabled = False

        if (
            self._service_manager.is_app_enabled(LIMirrorController.APP_NAME)
            and self.config.li_port_name
        ):
            self._li_port = BridgeTools.get_ofport(self.config.li_port_name)
            self._li_table = self._service_manager.get_table_num(
                LIMirrorController.APP_NAME,
            )
        self._ingress_tbl_num = self._service_manager.get_table_num(INGRESS)
        # following fields are only used in Non Nat config
        self._clean_restart = kwargs['config']['clean_restart']
        self._msg_hub = MessageHub(self.logger)
        self._datapath = None
        self.tbl_num = self._ingress_tbl_num

    def _get_config(self, config_dict):
        mtr_ip = None
        mtr_port = None
        port_no = config_dict.get('uplink_port', None)

        he_proxy_port = 0
        try:
            if 'proxy_port_name' in config_dict:
                he_proxy_port = BridgeTools.get_ofport(config_dict.get('proxy_port_name'))
        except DatapathLookupError:
            # ignore it
            self.logger.debug("could not parse proxy port config")

        if 'mtr_ip' in config_dict and 'mtr_interface' in config_dict and 'ovs_mtr_port_number' in config_dict:
            self._mtr_service_enabled = True
            mtr_ip = config_dict['mtr_ip']
            mtr_port = config_dict['ovs_mtr_port_number']
        else:
            mtr_ip = None
            mtr_port = None

        li_port_name = None
        if 'li_local_iface' in config_dict:
            li_port_name = config_dict['li_local_iface']

        return self.IngressConfig(
            uplink_port=port_no,
            mtr_ip=mtr_ip,
            mtr_port=mtr_port,
            li_port_name=li_port_name,
            setup_type=config_dict.get('setup_type', None),
            he_proxy_port=he_proxy_port,
        )

    def _get_default_flow_msgs(self, datapath) -> DefaultMsgsMap:
        return {
            self._ingress_tbl_num: self._get_default_ingress_flow_msgs(datapath),
        }

    def _get_default_ingress_flow_msgs(self, dp):
        """
        Sets up the ingress table, the first step in the packet processing
        pipeline.

        This sets up flow rules to annotate packets with a metadata bit
        indicating the direction. Incoming packets are defined as packets
        originating from the LOCAL port, outgoing packets are defined as
        packets originating from the gtp port.

        All other packets bypass the pipeline.

        Note that the ingress rules do *not* install any flows that cause
        PacketIns (i.e., sends packets to the controller).

        Raises:
            MagmaOFError if any of the default flows fail to install.
        """
        parser = dp.ofproto_parser
        next_table = self._service_manager.get_next_table_num(INGRESS)
        msgs = []

        # set traffic direction bits

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=OFPP_LOCAL)
        actions = [load_direction(parser, Direction.IN)]
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp,
                self._ingress_tbl_num, match, actions=actions,
                priority=flows.DEFAULT_PRIORITY, resubmit_table=next_table,
            ),
        )

        # set a direction bit for incoming (internet -> UE) traffic.
        match = MagmaMatch(in_port=self.config.uplink_port)
        actions = [load_direction(parser, Direction.IN)]
        msgs.append(
            flows.get_add_resubmit_next_service_flow_msg(
                dp, self._ingress_tbl_num, match,
                actions=actions,
                priority=flows.DEFAULT_PRIORITY,
                resubmit_table=next_table,
            ),
        )

        # Send RADIUS requests directly to li table
        if self._li_port:
            match = MagmaMatch(in_port=self._li_port)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num,
                    match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                    resubmit_table=self._li_table,
                ),
            )

        # set a direction bit for incoming (mtr -> UE) traffic.
        if self._mtr_service_enabled:
            match = MagmaMatch(in_port=self.config.mtr_port)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num,
                    match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

        if self.config.he_proxy_port != 0:
            match = MagmaMatch(in_port=self.config.he_proxy_port)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num,
                    match, actions=actions, priority=flows.DEFAULT_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

        if self.config.setup_type == 'CWF':
            # set a direction bit for outgoing (pn -> inet) traffic for remaining traffic
            ps_match_out = MagmaMatch()
            actions = [load_direction(parser, Direction.OUT)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num, ps_match_out,
                    actions=actions,
                    priority=flows.MINIMUM_PRIORITY,
                    resubmit_table=next_table,
                ),
            )
        else:
            # set a direction bit for outgoing (pn -> inet) traffic for remaining traffic
            # Passthrough is zero for packets from eNodeB GTP tunnels
            ps_match_out = MagmaMatch(passthrough=REG_ZERO_VAL)
            actions = [load_direction(parser, Direction.OUT)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num, ps_match_out,
                    actions=actions,
                    priority=flows.MINIMUM_PRIORITY,
                    resubmit_table=next_table,
                ),
            )

            # Passthrough is one for packets from remote PGW GTP tunnels, set direction
            # flag to IN for such packets.
            ps_match_in = MagmaMatch(passthrough=PASSTHROUGH_REG_VAL)
            actions = [load_direction(parser, Direction.IN)]
            msgs.append(
                flows.get_add_resubmit_next_service_flow_msg(
                    dp, self._ingress_tbl_num, ps_match_in,
                    actions=actions,
                    priority=flows.MINIMUM_PRIORITY,
                    resubmit_table=next_table,
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
