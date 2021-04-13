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
from lte.protos.mconfig import mconfigs_pb2
from magma.configuration.mconfig_managers import load_service_mconfig
from magma.pipelined.app.base import ControllerType, MagmaController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch
from magma.pipelined.openflow.registers import Direction
from ryu.lib import hub
from ryu.lib.packet import ether_types


class LIMirrorController(MagmaController):
    """
    LI Mirror controller.

    The LI Mirror controller is responsible for mirroring traffic to the LI
    agent for LEA processing.
    """

    APP_NAME = "li_mirror"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(LIMirrorController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)
        self.next_table = self._service_manager.get_next_table_num(
            self.APP_NAME)
        self._mirror_all = kwargs['config'].get('li_mirror_all', False)
        self._li_local_port = kwargs['config'].get('li_local_iface', None)
        self._li_local_port_num = BridgeTools.get_ofport(self._li_local_port)
        self._li_dst_port = kwargs['config'].get('li_dst_iface', None)
        self._li_dst_port_num = None
        if self._li_dst_port:
            self._li_dst_port_num = BridgeTools.get_ofport(self._li_dst_port)
        else:
            self.logger.warning("LI mirror port not setup, won't mirror pkts")
        self._datapath = None
        self._li_imsis = []

    def initialize_on_connect(self, datapath):
        """
        Install the default flows on datapath connect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)
        self._install_default_flows(datapath)
        self._datapath = datapath
        if self._li_dst_port_num:
            hub.spawn(self._monitor)

    def cleanup_on_disconnect(self, datapath):
        """
        Cleanup flows on datapath disconnect event.

        Args:
            datapath: ryu datapath struct
        """
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def _install_default_flows(self, datapath):
        """
        For each direction set the default flows to just forward to next table.
        If mirror flag set, copy all packets to li mirror port.

        Match traffic from local LI port and redirect it to dst li port

        Args:
            datapath: ryu datapath struct
        """
        parser = datapath.ofproto_parser
        inbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                   direction=Direction.IN)
        outbound_match = MagmaMatch(eth_type=ether_types.ETH_TYPE_IP,
                                    direction=Direction.OUT)
        actions = []
        if self._mirror_all and self._li_dst_port_num:
            self.logger.warning("Mirroring all traffic to LI")
            actions = [parser.OFPActionOutput(self._li_dst_port_num)]
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             inbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)
        flows.add_resubmit_next_service_flow(datapath, self.tbl_num,
                                             outbound_match, actions,
                                             priority=flows.MINIMUM_PRIORITY,
                                             resubmit_table=self.next_table)
        if self._li_dst_port_num:
            li_match = MagmaMatch(in_port=self._li_local_port_num)
            flows.add_output_flow(datapath, self.tbl_num, li_match, [],
                                  output_port=self._li_dst_port_num)

    def _install_mirror_flows(self, imsis):
        parser = self._datapath.ofproto_parser
        for imsi in imsis:
            self.logger.debug("Enabling LI tracking for IMSI %s", imsi)
            match = MagmaMatch(imsi=encode_imsi(imsi))
            actions = [parser.OFPActionOutput(self._li_dst_port_num)]
            flows.add_resubmit_next_service_flow(self._datapath, self.tbl_num,
                match, actions, priority=flows.DEFAULT_PRIORITY,
                resubmit_table=self.next_table)

    def _remove_mirror_flows(self, imsis):
        for imsi in imsis:
            self.logger.error("Disabling LI tracking for IMSI  %s", imsi)
            match = MagmaMatch(imsi=encode_imsi(imsi))
            flows.delete_flow(self._datapath, self.tbl_num, match)

    def _monitor(self, poll_interval=15):
        """
        Main thread that polls config updates and updates LI mirror flows
        """
        while True:
            mconfg_li_imsis = load_service_mconfig(
                'pipelined', mconfigs_pb2.PipelineD()).li_ues.imsis
            
            li_imsis = []
            for imsi in mconfg_li_imsis:
                if any(i.isdigit() for i in imsi):
                    li_imsis.append(imsi)
            imsis_to_add = [imsi for imsi in li_imsis if
                            imsi not in self._li_imsis]
            self._install_mirror_flows(imsis_to_add)
            imsis_to_rm = [imsi for imsi in self._li_imsis if
                           imsi not in li_imsis]
            self._remove_mirror_flows(imsis_to_rm)
            self._li_imsis = li_imsis
            hub.sleep(poll_interval)
