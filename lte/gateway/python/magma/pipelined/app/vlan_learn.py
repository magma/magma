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

from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import (
    DIRECTION_REG,
    IMSI_REG,
    VLAN_TAG_REG,
    Direction,
)
from ryu.ofproto import ether


class VlanLearnController(MagmaController):
    """
    A controller that sets up the vlan header for packets.

    This is an optional controller and will only be used for setups that need
    vlan headers. The incoming vlan id is learned on uplink and loaded back\
    into the packet on downlink.
    """

    APP_NAME = "vlan_learn"
    APP_TYPE = ControllerType.PHYSICAL

    LOAD_VLAN = 0x1
    PASSTHROUGH = 0x2

    def __init__(self, *args, **kwargs):
        super(VlanLearnController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        scratch_tables = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 2)
        self.vlan_header_scratch = scratch_tables[0]
        self.vlan_id_scratch = scratch_tables[1]
        self._datapath = None

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_flows(self._datapath)
        self._install_set_vlan_id_flows(self._datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self.vlan_id_scratch)
        flows.delete_all_flows_from_table(datapath, self.vlan_header_scratch)

    def remove_subscriber_flow(self, imsi: str):
        match = MagmaMatch(imsi=encode_imsi(imsi))
        flows.delete_flow(self._datapath, self.vlan_id_scratch, match)
        flows.delete_flow(self._datapath, self.vlan_header_scratch, match)

    def _install_default_flows(self, dp):
        """
        For UPLINK packets that have a vlan header learn the vlan id and load
        it back with the learn action.
        For UPLINK packets without a vlan header don't alter the corresponding
        DOWNLINK packets.
        """
        parser = self._datapath.ofproto_parser

        match = MagmaMatch(
            direction=Direction.OUT,
            vlan_vid=(0x1000, 0x1000),
        )
        actions = [
            parser.NXActionLearn(
                table_id=self.vlan_header_scratch,
                priority=flows.DEFAULT_PRIORITY,
                specs=[
                    parser.NXFlowSpecMatch(
                        src=(IMSI_REG, 0),
                        dst=(IMSI_REG, 0),
                        n_bits=64
                    ),
                    parser.NXFlowSpecMatch(
                        src=Direction.IN,
                        dst=(DIRECTION_REG, 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=self.LOAD_VLAN,
                        dst=(VLAN_TAG_REG, 0),
                        n_bits=16
                    ),
                ]
            ),
            parser.NXActionLearn(
                table_id=self.vlan_id_scratch,
                priority=flows.DEFAULT_PRIORITY,
                specs=[
                    parser.NXFlowSpecMatch(
                        src=(IMSI_REG, 0),
                        dst=(IMSI_REG, 0),
                        n_bits=64
                    ),
                    parser.NXFlowSpecMatch(
                        src=Direction.IN,
                        dst=(DIRECTION_REG, 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=('vlan_tci', 0),
                        dst=('vlan_tci', 0),
                        # 12 should match only vlan_vid
                        n_bits=12
                    ),
                ]
            )
        ]

        flows.add_resubmit_next_service_flow(
            dp, self.tbl_num, match, actions,
            priority=flows.UE_FLOW_PRIORITY,
            resubmit_table=self.next_table
        )

        # For non vlan traffic
        match = MagmaMatch(direction=Direction.OUT)
        actions = [
            parser.NXActionLearn(
                table_id=self.vlan_header_scratch,
                priority=flows.DEFAULT_PRIORITY,
                specs=[
                    parser.NXFlowSpecMatch(
                        src=(IMSI_REG, 0),
                        dst=(IMSI_REG, 0),
                        n_bits=64
                    ),
                    parser.NXFlowSpecMatch(
                        src=Direction.IN,
                        dst=(DIRECTION_REG, 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=self.PASSTHROUGH,
                        dst=(VLAN_TAG_REG, 0),
                        n_bits=16
                    ),
                ]
            ),
        ]
        flows.add_resubmit_next_service_flow(
            dp, self.tbl_num, match, actions,
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_table
        )

        # For downlink
        match = MagmaMatch(direction=Direction.IN)
        actions = [
            parser.NXActionResubmitTable(table_id=self.vlan_header_scratch),
            parser.NXActionResubmitTable(table_id=self.vlan_id_scratch)]
        flows.add_resubmit_next_service_flow(
            dp, self.tbl_num, match, actions,
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_table
        )

    def _install_set_vlan_id_flows(self, dp):
        parser = self._datapath.ofproto_parser

        # Load VLAN header (needs to be done before settubg vlan tag)
        # Can't be done through the learn action because ryu doesn't support it
        match = MagmaMatch(direction=Direction.IN,
                           vlan_tag=self.LOAD_VLAN)
        actions = [
            parser.OFPActionPushVlan(ether.ETH_TYPE_8021Q),
            parser.NXActionRegLoad2(dst=VLAN_TAG_REG,
                                    value=0)]
        flows.add_resubmit_next_service_flow(
            dp, self.vlan_id_scratch, match, actions,
            priority=flows.UE_FLOW_PRIORITY,
            resubmit_table=self.vlan_id_scratch
        )
