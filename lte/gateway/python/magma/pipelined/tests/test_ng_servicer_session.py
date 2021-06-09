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
import unittest.mock
import warnings
from collections import OrderedDict
from concurrent.futures import Future
from unittest import TestCase
from unittest.mock import MagicMock

from lte.protos.pipelined_pb2 import CauseIE
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.pipelined_test_util import (
    create_service_manager,
    start_ryu_app_thread,
    stop_ryu_app_thread,
)
from magma.pipelined.ng_set_session_msg import CreateSessionUtil

FAULTY_PDR_SESSION    = 1
FAULTY_FAR_SESSION    = 2
FAULTY_PDRFAR_SESSION = 3

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
                '5G_feature_set': {'enable': True,
                                  'node_identifier': '192.168.220.1'},
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

    # Send message and check response
    def _util_session_message_handler(self, set_session, cause_info=CauseIE.REQUEST_ACCEPTED):
        ng_serv = self.ng_services_controller

        pdr_rules = OrderedDict()

        # Context Response
        context_response = ng_serv.ng_session_message_handler(set_session, pdr_rules)
        TestCase().assertEqual(context_response.cause_info.cause_ie, cause_info)

    # Create generic session create request
    def _util_gen_session_create_request(self, subs_id="IMSI001010000000001",
                                         session_id=1, version=2, pdr_state="ADD",
                                         cause_ie=CauseIE.REQUEST_ACCEPTED):

        cls_sess = CreateSessionUtil(subs_id, session_id, version)

        cls_sess.CreateSession(subs_id, pdr_state, 100, 200,
                               "60.60.60.1", "192.168.10.11")

        #Test case matches values with expected cause_info
        self._util_session_message_handler(cls_sess._set_session, cause_ie)

        return cls_sess._set_session

    def test_smf_session_add_message(self):

        # Create the generic session
        self._util_gen_session_create_request()

        # Delete the generic session
        self._util_gen_session_create_request(pdr_state="REMOVE")

    def test_smf_faulty_session_messages(self):

        #Wrong session version
        cls_sess = CreateSessionUtil("IMSI001010000000001", 100, 0)
        self._util_session_message_handler(cls_sess._set_session, CauseIE.SESSION_CONTEXT_NOT_FOUND)

        #Wrong subscriber id
        cls_sess = CreateSessionUtil("", 100, 10)
        self._util_session_message_handler(cls_sess._set_session, CauseIE.SESSION_CONTEXT_NOT_FOUND)

    def test_smf_faulty_pdr_messages(self):
        #Create Session message with pdr_id=0
        cls_sess = CreateSessionUtil("IMSI001010000000001", 100, 1000)
        cls_sess.CreateSessionWithFaultyPDR()

        #Check if wrong PDR returns error
        self._util_session_message_handler(cls_sess._set_session, CauseIE.MANDATORY_IE_INCORRECT)

    def test_smf_faulty_far_messages(self):

        #Create Session message with PDR="INSTALL" but no FAR
        cls_sess = CreateSessionUtil("IMSI001010000000001", 100, 1000)
        cls_sess.CreateSessionWithFaultyFar()

        self._util_session_message_handler(cls_sess._set_session, CauseIE.INVALID_FORWARDING_POLICY)

if __name__ == "__main__":
    unittest.main()
