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
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from ryu.lib.packet import ether_types


class ConntrackController(MagmaController):
    """
    A controller that sets up conntrack flows for the UEs.

    This is an optional controller and will only be used for setups that need
    connection tracking

    Conntrack flags 0 is no action, 1 is commit

    CT state reference tuple (x,y):
    x:
     0: -
     1: +

    y:
      0x01: new
      0x02: est
      0x04: rel
      0x08: rpl
      0x10: inv
      0x20: trk
    """

    APP_NAME = "conntrack"
    APP_TYPE = ControllerType.LOGICAL
    CT_NEW = 0x01
    CT_EST = 0x02
    CT_REL = 0x04
    CT_RPL = 0x08
    CT_INV = 0x10
    CT_TRK = 0x20

    def __init__(self, *args, **kwargs):
        super(ConntrackController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self.conntrack_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]
        self.connection_event_table = \
            self._service_manager.INTERNAL_IPFIX_SAMPLE_TABLE_NUM
        self.zone = kwargs['config']['conntrackd'].get('zone', 897)
        self._datapath = None

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_flows(self._datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self.conntrack_scratch)

    def _install_default_flows(self, datapath):
        parser = datapath.ofproto_parser

        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ct_state=(0x0, self.CT_TRK))
        actions = [parser.NXActionCT(
            flags=0x0,
            zone_src=None,
            zone_ofs_nbits=self.zone,
            recirc_table=self.conntrack_scratch,
            alg=0,
            actions=[]
        )]
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             match, actions,
                                             priority=flows.DEFAULT_PRIORITY,
                                             resubmit_table=self.next_table)

        # Match all new connections on scratch table
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                           ct_zone=self.zone,
                           ct_state=(self.CT_NEW | self.CT_TRK,
                                     self.CT_NEW | self.CT_TRK))
        actions = [parser.NXActionCT(
            flags=0x1,
            zone_src=None,
            zone_ofs_nbits=self.zone,
            recirc_table=self.connection_event_table,
            alg=0,
            actions=[]
        )]
        flows.add_drop_flow(datapath, self.conntrack_scratch,
                            match, actions,
                            priority=flows.DEFAULT_PRIORITY)

        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP)
        flows.add_resubmit_next_service_flow(
            datapath, self.tbl_num, match, [],
            priority=flows.MINIMUM_PRIORITY,
            resubmit_table=self.next_table)
