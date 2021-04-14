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

import logging
import subprocess
import unittest
import unittest.mock
import warnings
from collections import OrderedDict
from concurrent.futures import Future
from typing import List
from unittest import TestCase
from unittest.mock import MagicMock, Mock

from lte.protos import (
    pipelined_pb2,
    pipelined_pb2_grpc,
    session_manager_pb2_grpc,
)
from lte.protos.session_manager_pb2 import UPFAssociationState
from magma.pipelined.app.ng_services import NGServiceController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)


def mocked_send_node_state_message_success (node_message):
    return True

def mocked_send_node_state_message_failure (node_message):
    return False


class NGServiceControllerTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'

    def setUp(self):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(NGServiceControllerTest, self).setUpClass()
        warnings.simplefilter('ignore')
        self.service_manager = create_service_manager([])

        ng_services_controller_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.NGServiceController,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.NGServiceController:
                    ng_services_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'enodeb_iface': 'eth1',
                'clean_restart': True,
                '5G_feature_set': {'enable': True},
                '5G_feature_set': {'node_identifier': '192.168.220.1'},
                'bridge_name': self.BRIDGE,
            },
            mconfig=None,
            loop=None,
            service_manager=self.service_manager,
            integ_test=False,
            rpc_stubs={'sessiond_setinterface': MagicMock()}
        )

        BridgeTools.create_bridge(self.BRIDGE, self.IFACE)

        self.thread = start_ryu_app_thread(test_setup)
        self.ng_services_controller = \
            ng_services_controller_reference.result()
        self.testing_controller = testing_controller_reference.result()

    def tearDown(self):
        stop_ryu_app_thread(self.thread)
        BridgeTools.destroy_bridge(self.BRIDGE)

    def _default_settings(self, mock_func, version=0, msg_count=0):
        ng_serv = self.ng_services_controller
        node_mgr = ng_serv._ng_node_mgr

        node_mgr._smf_assoc_version = version
        node_mgr._assoc_message_count = msg_count
        node_mgr._send_messsage_wrapper = mock_func
        node_message = node_mgr.get_node_assoc_message()

        return (node_message.associaton_state)

    def test_association_setup_message_request (self):
        ng_serv = self.ng_services_controller
        node_mgr = ng_serv._ng_node_mgr

        assoc_message = self._default_settings(mocked_send_node_state_message_success)
        node_mgr._send_association_request_message(assoc_message)
        TestCase().assertEqual(node_mgr._smf_assoc_version, 1)
        TestCase().assertEqual(node_mgr._assoc_message_count, 1)

    def test_association_release_message (self):
        ng_serv = self.ng_services_controller
        node_mgr = ng_serv._ng_node_mgr

        # Assume that the Association is established
        node_mgr._smf_assoc_state = UPFAssociationState.ESTABLISHED

        self._default_settings(mocked_send_node_state_message_success, 1, 1)
        node_mgr.send_association_release_message()
        TestCase().assertEqual(node_mgr._smf_assoc_version, 0)
        TestCase().assertEqual(node_mgr._assoc_message_count, 2)

    def test_association_setup_message_request_failure (self):
        ng_serv = self.ng_services_controller
        node_mgr = ng_serv._ng_node_mgr

        # Assume that the Association is established
        node_mgr._smf_assoc_state = UPFAssociationState.ESTABLISHED

        assoc_message = self._default_settings(mocked_send_node_state_message_failure)
        node_mgr.send_association_release_message()

        TestCase().assertEqual(node_mgr._smf_assoc_version, 0)
        TestCase().assertEqual(node_mgr._assoc_message_count, 0)
        TestCase().assertEqual(node_mgr._smf_assoc_state, UPFAssociationState.RELEASE)

if __name__ == "__main__":
    unittest.main()
