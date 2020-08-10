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

import unittest
import warnings
import logging

from unittest import TestCase
from unittest.mock import MagicMock
from magma.pipelined.node_state_manager import NodeStateManager 
from lte.protos.mconfig.mconfigs_pb2 import PipelineD

from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)

def mocked_send_node_state_message_success (node_message):
    return True

def mocked_send_node_state_message_failure (node_message):
    return False

class NodeStateManagerTest(unittest.TestCase):

    def setUp(self):
        magma_service_mock = MagicMock()
        magma_service_mock.config = {
            'enodeb_iface': 'eth1'
        }
        self.node_state_manager = NodeStateManager(magma_service_mock)
 
    def _default_settings(self, mock_func):
        self.node_state_manager.assoc_message_version = 0
        self.node_state_manager.assoc_message_count = 0
        self.node_state_manager._send_messsage_wrapper = mock_func

    def test_association_set_message_success (self):
        self._default_settings(mocked_send_node_state_message_success)
        self.node_state_manager.send_association_setup_message()          
        TestCase().assertEqual(self.node_state_manager.assoc_message_version, 1)    
        TestCase().assertEqual(self.node_state_manager.assoc_message_count, 1)    

    def test_association_set_message_failure (self):
        self._default_settings (mocked_send_node_state_message_failure)
        self.node_state_manager.send_association_setup_message()
        TestCase().assertEqual(self.node_state_manager.assoc_message_version, 0)
        TestCase().assertEqual(self.node_state_manager.assoc_message_count, 0)

    def test_association_del_message_success (self):
        self._default_settings (mocked_send_node_state_message_success)
        self.node_state_manager.send_association_release_message()
        TestCase().assertEqual(self.node_state_manager.assoc_message_version, 0)
        TestCase().assertEqual(self.node_state_manager.assoc_message_count, 1)

if __name__ == "__main__":
    unittest.main()
