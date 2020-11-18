import warnings
from typing import List
import subprocess
import unittest
from unittest import TestCase
import unittest.mock
from collections import OrderedDict
from concurrent.futures import Future
from unittest.mock import MagicMock
from concurrent import futures
import grpc

from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.start_pipelined import (
                TestSetup,
                PipelinedController,
            )

#from magma.pipelined.tests.pipelined_test_util import (
#    start_ryu_app_thread,
#    stop_ryu_app_thread,
    #create_service_manager,
#    wait_after_send
#)
#from magma.pipelined.app.ng_services import NGServiceController
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos import pipelined_pb2_grpc
from lte.protos import pipelined_pb2
from lte.protos.pipelined_pb2 import (
            NodeID,
            SessionSet,
            UPFSessionContextState,
            FailureRuleInformation,
            OffendingIE,
            UPFSessionContextState,
            CauseIE,
            Fsm_state,
            ActionType)
from lte.protos import session_manager_pb2_grpc
from lte.protos.session_manager_pb2_grpc import SetInterfaceForUserPlaneStub
from lte.protos.session_manager_pb2 import (
                                            UPFSessionState,
                                            UPFSessionConfigState)

from lte.protos.pipelined_pb2 import (
            SessionSet,
            #SetGroupFAR,
            FwdParam,
            #ApplyAction,
            OuterHeaderCreation,
            SetGroupPDR,
            PDI,
            #LocalFTeid,
            OuterHeadRemoval,
            CauseIE)

from unittest.mock import Mock, MagicMock


class UpfRpcSessionServer(pipelined_pb2_grpc.PipelinedServicer):
  
     #def __init__ (self, loop, service_manager):
     def __init__ (self, loop):
         self._loop = loop
         #self._ng_servicer_app = ng_servicer_app
         #self._service_manager = service_manager
         
     def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        pipelined_pb2_grpc.add_PipelinedServicer_to_server(self, server)

     def SetSMFSessions(self, request, context):
         context_response=\
            UPFSessionContextState(cause_info=CauseIE(cause_ie=CauseIE.REQUEST_ACCEPTED),
                                   session_snapshot=UPFSessionState(subscriber_id=request.subscriber_id,
                                                                   session_version=request.session_version))
         return context_response                                                          
   
class RpcTests(unittest.TestCase):
    """
    Tests for the IPAllocator rpc servicer and stub
    """

    def setUp(self):
        # Bind the rpc server to a free port
        thread_pool = futures.ThreadPoolExecutor(max_workers=10)
        self._rpc_server = grpc.server(thread_pool)
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')

        loop_mock = MagicMock()
        #self.service_manager = create_service_manager([])

        self._servicer = UpfRpcSessionServer(loop_mock)
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        self._stub = PipelinedStub(channel)
        self._set_session = \
            SessionSet(subscriber_id="IMSI001222333", session_version=2,\
                       node_id=NodeID(node_id_type=NodeID.IPv4, identifier="192.168.220.1"),\
                       state=Fsm_state(state=Fsm_state.SESSION_ACTIVE))
       
        self._set_session.set_gr_pdr.extend ([\
                         SetGroupPDR (pdr_id=50, pdr_version=34,
                                      precedence=32,\
                                      pdr_state=SetGroupPDR.ADD,\
                                      pdi=PDI(src_interface=0, \
                                              local_f_teid=98,\
                                              ue_ip_adr="1.1.1.1"), \
                                      outer_head_desc=0, \
                                      far_action_to_apply=[ActionType(act_type=ActionType.FORW),\
                                                           ActionType(act_type=ActionType.BUFF)],\
                                      far_fwd_parm=FwdParam(dest_iface=0, \
                                                        outer_head_cr=\
                                                          OuterHeaderCreation(o_teid=68,ipv4_adr="192.168.76.100")))])



    def tearDown(self):
        self._rpc_server.stop(0)

    def test_smf_session_set_message (self):
        self._stub.SetSMFSessions(self._set_session)
