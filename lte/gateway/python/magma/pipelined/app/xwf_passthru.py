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
from magma.pipelined.app.inout import INGRESS
from magma.pipelined.openflow import flows
from magma.pipelined.openflow.magma_match import MagmaMatch


class XWFPassthruController(MagmaController):

    APP_NAME = "xwf_passthru"
    APP_TYPE = ControllerType.SPECIAL

    def __init__(self, *args, **kwargs):
        super(XWFPassthruController, self).__init__(*args, **kwargs)
        self.tbl_num = self._service_manager.get_table_num(self.APP_NAME)

    def initialize_on_connect(self, datapath):
        self.delete_all_flows(datapath)
        self._install_passthrough_flow(datapath)

    def cleanup_on_disconnect(self, datapath):
        self.delete_all_flows(datapath)

    def delete_all_flows(self, datapath):
        flows.delete_all_flows_from_table(datapath, self.tbl_num)

    def _install_passthrough_flow(self, datapath):
        flows.add_resubmit_next_service_flow(
            datapath,
            self.tbl_num,
            MagmaMatch(),
            actions=[],
            priority=flows.DEFAULT_PRIORITY,
            resubmit_table=self._service_manager.get_table_num(INGRESS),
        )
