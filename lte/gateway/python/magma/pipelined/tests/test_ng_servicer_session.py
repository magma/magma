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
import warnings
from typing import List
import subprocess
import unittest
from unittest import TestCase
import unittest.mock
from collections import ( 
    OrderedDict,
    namedtuple)
from concurrent.futures import Future
from unittest.mock import MagicMock

from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController)

from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    wait_after_send)

from magma.pipelined.app.ng_services import NGServiceController
from lte.protos import pipelined_pb2_grpc
from lte.protos import pipelined_pb2
from lte.protos import session_manager_pb2_grpc
from lte.protos.session_manager_pb2_grpc import SetInterfaceForUserPlaneStub
from lte.protos.session_manager_pb2 import (
    UPFSessionConfigState)

from lte.protos.session_manager_pb2 import NodeID
from lte.protos.pipelined_pb2 import (
    SessionSet,
    SetGroupFAR,
    FwdParam,
    Action,
    OuterHeaderCreation,
    SetGroupPDR,
    PDI,
    OuterHeadRemoval,
    Fsm_state,
    PdrState,
    CauseIE
)

from unittest.mock import Mock, MagicMock
from magma.pipelined.ng_manager.session_state_manager import SessionStateManager
from magma.pipelined.ng_manager.session_state_manager_util import FARRuleEntry

FAULTY_PDR_SESSION    = 1
FAULTY_FAR_SESSION    = 2
FAULTY_PDRFAR_SESSION = 3

class CreateSessionUtil:

    def __init__(self, subscriber_id:str, session_version=2, node_id="192.168.220.1"):
        self._set_session = \
                  SessionSet(subscriber_id=subscriber_id, session_version=session_version,\
                             node_id=NodeID(node_id_type=NodeID.IPv4, node_id=node_id),\
                             state=Fsm_state(state=Fsm_state.CREATED))

    def CreateSessionPDR(self, pdr_id:int, pdr_version:int, pdr_state,
                         precedence:int, local_f_teid:int, ue_ip_addr:str,
                         far_tuple=None):

        far_gr_entry=None

        if far_tuple:
            if far_tuple.o_teid != 0:
                # For pdr_id=2 towards access
                far_gr_entry = SetGroupFAR(far_action_to_apply=[far_tuple.apply_action],\
                                           fwd_parm=FwdParam(dest_iface=0, \
                                           outr_head_cr=OuterHeaderCreation(\
                                             o_teid=far_tuple.o_teid, gnb_ipv4_adr=far_tuple.gnb_ip_addr)))
            else:
                far_gr_entry = SetGroupFAR(far_action_to_apply=[far_tuple.apply_action])


        if local_f_teid != 0:
            self._set_session.set_gr_pdr.extend([\
                         SetGroupPDR(pdr_id=pdr_id, pdr_version=pdr_version,
                                     pdr_state=pdr_state,\
                                     precedence=precedence,\
                                     pdi=PDI(src_interface=0,\
                                              local_f_teid=local_f_teid,\
                                              ue_ip_adr=ue_ip_addr), \
                                      o_h_remo_desc=0, \
                                      set_gr_far=far_gr_entry,
                                      activate_flow_req=None)])

        else:
            self._set_session.set_gr_pdr.extend([\
                         SetGroupPDR(pdr_id=pdr_id, pdr_version=pdr_version,
                                     pdr_state=pdr_state,\
                                     precedence=precedence,\
                                     pdi=PDI(src_interface=1, ue_ip_adr=ue_ip_addr),\
                                     set_gr_far=far_gr_entry,
                                     activate_flow_req=None)])

    def CreateGenSessionTemplate(self, pdr_add=True):
        if pdr_add == True:
            # From Access towards Core pdrid, version, precedence, teid, ue-ip, farid
            self.CreateSessionPDR(1, 1, PdrState.Value('INSTALL'), 32, 98, "60.60.60.1",\
                                  FARRuleEntry(Action.Value('FORW'), 0, ''))

            # From Core towards access, pdrid, version, precedence, teid=0, ue-ip, farid
            self.CreateSessionPDR(2, 1, PdrState.Value('INSTALL'), 32, 0, "60.60.60.1",\
                                  FARRuleEntry(Action.Value('FORW'), 1, "100.200.200.1"))

        else:
            # Delete PDR from Acces toward CORE
            self.CreateSessionPDR(1, 1, PdrState.Value('REMOVE'), 32, 98, "60.60.60.1")
                               
            # Delete PDR from Core toward ACCESS
            self.CreateSessionPDR(2, 1, PdrState.Value('REMOVE'), 32, 0, "60.60.60.1")


    def CreateFaultyPDRFARSessions (self, fault_type):
        if fault_type == FAULTY_PDR_SESSION:
            self.CreateSessionPDR(0, 1, SetGroupPDR.PDR_ADD, 32, 98, "60.60.60.1",\
                                  FARRuleEntry(ApplyAction.FORW, 0, ''))
        elif fault_type == FAULTY_PDRFAR_SESSION:
            self.CreateSessionPDR(1, 1, SetGroupPDR.PDR_ADD, 32, 98, "60.60.60.1")

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

    # Send message and check response
    def _util_session_message_handler(self, set_session, cause_info=CauseIE.REQUEST_ACCEPTED):
        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        pdr_rules = OrderedDict()

        # Context Response
        context_response = ng_serv.ng_session_message_handler(set_session, pdr_rules)
        TestCase().assertEqual(context_response.cause_info.cause_ie, cause_info)

    # Create generic session create request
    def _util_gen_session_create_request(self, pdr_add=True):

        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        cls_sess = CreateSessionUtil("IMSI001010000000001")
        cls_sess.CreateGenSessionTemplate(pdr_add)

        self._util_session_message_handler(cls_sess._set_session, CauseIE.REQUEST_ACCEPTED)

        return cls_sess._set_session

    def mock_validate_report_session_config_state (self, msg_type, upf_session_config_state:UPFSessionConfigState):
        TestCase().assertEqual(self._smf_session_set.subscriber_id, \
                               upf_session_config_state.upf_session_state[0].subscriber_id)

        TestCase().assertEqual(self._smf_session_set.sess_ver_no, \
                               upf_session_config_state.upf_session_state[0].sess_ver_no)

    def test_smf_session_add_message(self):
        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        # Create the generic session 
        set_session = self._util_gen_session_create_request()
       
        # Delete the generic session
        set_session = self._util_gen_session_create_request(False)

    def test_smf_faulty_session_messages(self):
        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        pdr_rules = OrderedDict()

        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        #Wrong session version
        cls_sess = CreateSessionUtil("IMSI001010000000001")
        cls_sess._set_session.session_version = 0
        self._util_session_message_handler(cls_sess._set_session, CauseIE.SESSION_CONTEXT_NOT_FOUND)
 
        #Wrong subscriber id
        cls_sess._set_session.subscriber_id=""
        self._util_session_message_handler(cls_sess._set_session, CauseIE.SESSION_CONTEXT_NOT_FOUND)

    def test_smf_faulty_pdr_messages(self):
        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        pdr_rules = OrderedDict()

        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        cls_sess = CreateSessionUtil("IMSI001010000000001")
        cls_sess.CreateGenSessionTemplate()
        cls_sess._set_session.set_gr_pdr[0].pdr_id = 0

        self._util_session_message_handler(cls_sess._set_session, CauseIE.MANDATORY_IE_INCORRECT)


    def test_smf_faulty_far_messages(self):
        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        pdr_rules = OrderedDict()

        ng_serv = self.ng_services_controller
        sess_mgr = ng_serv._ng_sess_mgr

        cls_sess = CreateSessionUtil("IMSI001010000000001")
        cls_sess.CreateGenSessionTemplate()
        cls_sess._set_session.set_gr_pdr[0].ClearField('set_gr_far')

        self._util_session_message_handler(cls_sess._set_session, CauseIE.INVALID_FORWARDING_POLICY)

if __name__ == "__main__":
    unittest.main()
