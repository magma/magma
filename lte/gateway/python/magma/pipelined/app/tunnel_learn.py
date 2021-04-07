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
from magma.pipelined.openflow.registers import DIRECTION_REG, Direction


class TunnelLearnController(MagmaController):
    """
    A controller that sets up tunnel/ue learn flows based on uplink UE traffic
    to properly route downlink packets back to the UE (through the correct GRE
    flow tunnel).

    This is an optional controller and will only be used for setups with flow
    based GRE tunnels.
    """

    APP_NAME = "tunnel_learn"
    APP_TYPE = ControllerType.PHYSICAL

    def __init__(self, *args, **kwargs):
        super(TunnelLearnController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = \
            self._service_manager.get_next_table_num(self.APP_NAME)
        self.tunnel_learn_scratch = \
            self._service_manager.allocate_scratch_tables(self.APP_NAME, 1)[0]
        self._datapath = None

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_tunnel_classify_flows(self._datapath)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
        flows.delete_all_flows_from_table(datapath, self.tunnel_learn_scratch)

    def remove_subscriber_flow(self, mac_addr: str):
        match = MagmaMatch(eth_dst=mac_addr)
        flows.delete_flow(self._datapath, self.tbl_num, match)
        flows.delete_flow(self._datapath, self.tunnel_learn_scratch, match)

    def _install_default_tunnel_classify_flows(self, dp):
        """
        For direction OUT:
            add flow with learn action(which will add a rule matching on UE
            mac_addr in a scratch table, that will load the necessary
            gre infomration for the incoming(direction IN) flow)
        For direction IN:
            Will get forwarded to the scratch table and matched on the flow
                from the learn action
            Finally will get forwarded to the next table
        """
        parser = dp.ofproto_parser

        # Add a learn action that will match on UE mac, and:
        #   load gre tun_id, swap and load gre tun src and gre tun dst mac
        # Example learned flow:
        #   reg1=0x10,dl_dst=aa:29:3e:95:64:40
        #   actions=load:0x1389->NXM_NX_TUN_ID[0..31],
        #           load:0xc0a84666->NXM_NX_TUN_IPV4_DST[],
        #           load:0xc0a84665->NXM_NX_TUN_IPV4_SRC[]
        #
        outbound_match = MagmaMatch(direction=Direction.OUT)
        actions = [
            parser.NXActionLearn(
                table_id=self.tunnel_learn_scratch,
                priority=flows.DEFAULT_PRIORITY,
                specs=[
                    parser.NXFlowSpecMatch(
                        src=('eth_src_nxm', 0),
                        dst=('eth_dst_nxm', 0),
                        n_bits=48
                    ),
                    parser.NXFlowSpecMatch(
                        src=Direction.IN,
                        dst=(DIRECTION_REG, 0),
                        n_bits=32
                    ),
                    parser.NXFlowSpecLoad(
                        src=('tun_ipv4_src', 0),
                        dst=('tun_ipv4_dst', 0),
                        n_bits=32
                    ),
                    # TODO This might be getting overwritten by the IP stack,
                    # check if its required
                    parser.NXFlowSpecLoad(
                        src=('tun_ipv4_dst', 0),
                        dst=('tun_ipv4_src', 0),
                        n_bits=32
                    ),
                ]
            )
        ]
        flows.add_resubmit_next_service_flow(dp, self.tbl_num,
                                             outbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

        # The inbound match will first send packets to the scratch table,
        # where global registers will be set and the packet will be dropped
        # Then the final action will send packet down the pipelined(with the
        # necessary tunnel information loaded from the scratch table)
        inbound_match = MagmaMatch(direction=Direction.IN)
        actions = [
            parser.NXActionResubmitTable(table_id=self.tunnel_learn_scratch)]
        flows.add_resubmit_next_service_flow(dp, self.tbl_num,
                                             inbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)
