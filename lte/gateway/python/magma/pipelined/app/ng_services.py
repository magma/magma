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
from .base import MagmaController, ControllerType
from magma.pipelined.ng_manager.node_state_manager import NodeStateManager
from magma.pipelined.ng_manager.session_state_manager import SessionStateManager

class NGServiceController(MagmaController):
    """
    This class is intended to be a place holder for
    all 5G services which include node and session
    management
    """
    APP_NAME = "ng_services"
    APP_TYPE = ControllerType.LOGICAL

    def __init__(self, *args, **kwargs):
        super(NGServiceController, self).__init__(*args, **kwargs)

        # General Initialization
        self.loop = kwargs['loop']
        self.config = kwargs['config']

        #Get SessionD Channel
        self.sessiond_setinterface = kwargs['rpc_stubs']['sessiond_setinterface']

        # Initialize ng services
        self._ng_node_mgr = NodeStateManager(self.loop, self.sessiond_setinterface,
                                             self.config)
        self._ng_sess_mgr = None

    def initialize_on_connect(self, _):
        """
        Initialize node_state and session_state manager
        with datapath connect event.

        Args:
            None:
        """
        self._ng_node_mgr.send_association_setup_message()
        self._ng_sess_mgr = SessionStateManager(self.loop, self.logger)

    def cleanup_on_disconnect(self, _):
        """
        Send notification to sessiond about association
        release

        Args:
            None
        """
        self._ng_node_mgr.send_association_release_message()

    def delete_all_flows(self, _):
        # in case the stored sesson manager to call delete from tunnel
        pass

    def cleanup_state(self):
        pass

    def ng_services_stop(self):
        self._ng_node_mgr.send_association_release_message()

    def ng_set_node_association(self):
        self._ng_node_mgr.send_association_setup_message()

    # Process the message and send it to SessionStateManager
    def ng_session_message_handler(self, new_session, process_pdr_rules):
        return self._ng_sess_mgr.process_session_message(new_session, process_pdr_rules)
