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
import netifaces
from collections import namedtuple

from magma.common.misc_utils import cidr_to_ip_netmask_tuple
from magma.pipelined.app.base import MagmaController, ControllerType
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction, load_passthrough
from magma.pipelined.directoryd_client import get_all_records

from ryu.controller import dpset
from ryu.lib.packet import ether_types
from ryu.ofproto.inet import IPPROTO_ICMPV6


class IPV6RouterSolicitationController(MagmaController):
    """
    IPV6RouterSolicitationController responds to ipv6 router solicitation
    messages

    (1) Listens to flows with IPv6 src address prefixed with ""fe80".
    (2) Extracts interface ID (lower 64 bits) from the Router Solicitation
        message.
    (3) Performs a look up to find the IPv6 prefix that corresponds to the
        interface ID. The look up can be done using a local look up table that
        is updated during session creation where the full 128 bit IPv6 address
        assigned to UE is provided.
    (4) Generates a router advertisement message targeting the GTP tunnel.

    """
    APP_NAME = 'ipv6_router_solicitation'
    APP_TYPE = ControllerType.PHYSICAL

    # Inherited from app_manager.RyuApp
    _CONTEXTS = {
        'dpset': dpset.DPSet,
    }

    def __init__(self, *args, **kwargs):
        super(IPV6RouterSolicitationController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self.setup_type = kwargs['config']['setup_type']
        self._datapath = None

    def initialize_on_connect(self, datapath):
        self._datapath = datapath
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._install_default_ipv6_flows(datapath)

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next app.
        """
        match = MagmaMatch()

        flows.add_resubmit_next_service_flow(datapath, self.tbl_num, match, [],
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

    def _install_default_ipv6_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next app.
        """
        match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IPV6,
                           ipv6_dst='fe80::/10',
                           ip_proto=IPPROTO_ICMPV6,
                           direction=Direction.IN)

        flows.add_resubmit_next_service_flow(datapath, self.tbl_num, match, [],
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)

        flows.add_output_flow(datapath, self.tbl_num,
                              match=MagmaMatch(), actions=[],
                              priority=flows.PASSTHROUGH_PRIORITY,
                              output_port=ofproto.OFPP_CONTROLLER,
                              copy_table=self.next_table,
                              max_len=ofproto.OFPCML_NO_BUFFER)


    def handle_restart(self):
        pass

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)
