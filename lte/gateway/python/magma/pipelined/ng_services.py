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
from magma.pipelined.node_state_manager import NodeStateManager 
from magma.common.service import MagmaService

class NGServices:
    """
    This class is intended to be a place holder for
    all 5G services which include node and session
    management
    """
    def __init__(self, magma_service: MagmaService):
        self._magma_service = magma_service
        self._ng_node_state_manager = 0

    # Initilize the class for node and session management
    def ng_services_start(self):
        self._ng_node_state_manager = NodeStateManager(self._magma_service)
        #self._ng_sessiong_mgmt = SessionManagement(self._magma_service)
 
        self.ng_set_node_association ()

    def ng_services_stop(self):
        self._ng_node_state_manager.send_association_release_message()

    def ng_set_node_association(self):
        self._ng_node_state_manager.send_association_setup_message()
