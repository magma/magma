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
from collections import OrderedDict
from concurrent.futures import Future
from unittest.mock import MagicMock
import grpc
from concurrent import futures

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

from lte.protos.pipelined_pb2 import (
    NodeID,
    SessionSet,
    SetGroupFAR,
    FwdParam,
    ApplyAction,
    OuterHeaderCreation,
    SetGroupPDR,
    PDI,
    OuterHeadRemoval,
    CauseIE,
    Fsm_state)

from unittest.mock import Mock, MagicMock
from magma.pipelined.ng_manager.session_state_manager import SessionStateManager
from orc8r.protos.common_pb2 import Void

class CreateSessionUtil:
    def __init__(self, subscriber_id:str, session_version=2, node_id="192.168.220.1"):
        self._set_session = \
                  SessionSet(subscriber_id=subscriber_id, session_version=session_version,\
                             node_id=NodeID(node_id_type=NodeID.IPv4, node_id=node_id),\
                             state=Fsm_state(state=Fsm_state.CREATED))

    def CreateSessionPDR(self, pdr_id:int, pdr_version:int, pdr_state,
                         precedence:int, local_f_teid:int, ue_ip_addr:str,
                         qos_enforcer=None, far_tuple=None):

        far_gr_entry=None

        if far_tuple:
            if far_tuple.o_teid != 0:
                # For pdr_id=2 towards access
                far_gr_entry = SetGroupFAR(far_action_to_apply=[far_tuple.apply_action],\
                                           fwd_parm=FwdParam(dest_iface=0, \
                                           outr_head_cr=OuterHeaderCreation(\
                                             o_teid=far_tuple.o_teid, ipv4_adr=far_tuple.enodeb_ip_addr)))
            else:
                far_gr_entry = SetGroupFAR(far_action_to_apply=[far_tuple.apply_action])

        qos_enforce_rule = self.CreateQERinPDR(qos_enforcer)

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
                                      activate_flow_req=qos_enforce_rule)])

    def CreateGenSessionTemplate(self, to_core=True):

        if to_core:
            self.PdrTowardsCore("100")
        else:    
            self.PdrTowardsAccess("200")
        

class SMFSessionConfigTest(session_manager_pb2_grpc.SetInterfaceForUserPlaneServicer):

     def __init__ (self, loop):
         self._loop = loop

     def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        session_manager_pb2_grpc.add_SetInterfaceForUserPlaneServicer_to_server(self, server)

     def SetUPFSessionsConfig(self, request, context):
         return (Void())

class RpcTests(unittest.TestCase):
    """
    Tests NG Node related servicers
    """
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    ASSOCIATED = 0

    def setUp(self):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        warnings.simplefilter('ignore')

        loop_mock = MagicMock()

        # Bind the rpc server to a free port
        thread_pool = futures.ThreadPoolExecutor(max_workers=10)
        self._rpc_server = grpc.server(thread_pool)
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')

        self._servicer = SMFSessionConfigTest(loop_mock)
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))

        self._channel = channel

        config_mock ={
                   'enodeb_iface': 'eth1',
                   'clean_restart': True,
                   '5G_feature_set': {'enable': True},
                   '5G_feature_set': {'node_identifier': '192.168.220.1'},
                   'bridge_name': self.BRIDGE,
               }

        self._ng_sess_mgr = SessionStateManager(loop_mock, channel)

    def tearDown(self):
        self._rpc_server.stop(0)

    def test_session_config_thread(self):
        cls_sess = CreateSessionUtil();
        cls_sess.CreateGenSessionTemplate()
        
        sess_mgr = self._ng_sess_mgr
        
        pdr_rules = OrderedDict()

        # Create session 100 with core PDR        
        context_response = sess_mgr.process_session_message(cls_sess._set_session, pdr_rules) 
        TestCase().assertEqual(context_response.cause_info.cause_ie, CauseIE.REQUEST_ACCEPTED)
       
        sess_mgr.update_session_registry(cls_sess._set_session)
        
        pdr_rules.clear()
        
        # Create session 200 with access PDR
        cls_sess.CreateGenSessionTemplate(False)
        TestCase().assertEqual(context_response.cause_info.cause_ie, CauseIE.REQUEST_ACCEPTED)
      
        # As its a created case
        sess_mgr.update_session_registry(cls_sess._set_session)
        
        sess_mgr._report_session_config_state()
        
        TestCase().assertEqual(sess_mgr.periodic_config_msg_count, 1)
