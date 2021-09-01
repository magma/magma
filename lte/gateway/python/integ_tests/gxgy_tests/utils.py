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
from concurrent.futures import Future, ThreadPoolExecutor
from contextlib import ExitStack

import grpc
from lte.protos import session_manager_pb2_grpc
from magma.common.service_registry import ServiceRegistry
from magma.configuration.service_configs import load_service_config
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery as FlowQuery
from magma.pipelined.tests.app.packet_injector import ScapyPacketInjector
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)
from magma.policydb.rule_store import PolicyRuleDict

from .session_manager import MockSessionManager


class GxGyTestUtil(object):
    BRIDGE = 'gtp_br0'
    IFACE = 'gtp_br0'
    CONTROLLER_PORT = 6644

    def __init__(self):
        self.static_rules = PolicyRuleDict()

        # Local sessiond
        self.sessiond = session_manager_pb2_grpc.LocalSessionManagerStub(
            ServiceRegistry.get_rpc_channel("sessiond", ServiceRegistry.LOCAL),
        )

        self.proxy_responder = session_manager_pb2_grpc.SessionProxyResponderStub(
            ServiceRegistry.get_rpc_channel("sessiond", ServiceRegistry.LOCAL),
        )

        # Mock session controller server
        cloud_port = load_service_config("sessiond")["local_controller_port"]
        self.controller = MockSessionManager()
        self.server = grpc.server(ThreadPoolExecutor(max_workers=10))
        session_manager_pb2_grpc.add_CentralSessionControllerServicer_to_server(
            self.controller, self.server,
        )
        self.server.add_insecure_port('127.0.0.1:{}'.format(cloud_port))
        self.server.start()

        # Add new controller to bridge
        BridgeTools.add_controller_to_bridge(self.BRIDGE, self.CONTROLLER_PORT)

        # Start ryu test controller for adding flows
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.Testing],
            references={
                PipelinedController.Testing: testing_controller_reference,
            },
            config={
                'bridge_name': self.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'controller_port': self.CONTROLLER_PORT,
            },
            mconfig=None,
            loop=None,
            integ_test=True,
        )
        self.thread = start_ryu_app_thread(test_setup)
        self.testing_controller = testing_controller_reference.result()

    def cleanup(self):
        # Stop ryu controller
        stop_ryu_app_thread(self.thread)
        # Remove bridge
        BridgeTools.remove_controller_from_bridge(
            self.BRIDGE,
            self.CONTROLLER_PORT,
        )
        # Stop gRPC server
        self.server.stop(0)

    def get_packet_sender(self, subs, packets, count):
        """
        Return a function to call within a greenthread to send packets and
        return the number of packets that went through table 20 (i.e. didn't
        get dropped)
        Args:
            subs ([SubscriberContext]): list of subscribers that may receive
                packets
            packets ([ScapyPacket]): list of packets to send
            count (int): how many of each packet to send
        """
        pkt_sender = ScapyPacketInjector(self.IFACE)

        def packet_sender():
            isolators = [
                RyuDirectTableIsolator(
                    RyuForwardFlowArgsBuilder.from_subscriber(sub)
                                             .build_requests(),
                    self.testing_controller,
                ) for sub in subs
            ]
            flow_query = FlowQuery(20, self.testing_controller)
            pkt_start = sum(flow.packets for flow in flow_query.lookup())
            with ExitStack() as es:
                for iso in isolators:
                    es.enter_context(iso)
                for packet in packets:
                    pkt_sender.send(packet, count=count)
                wait_after_send(self.testing_controller)
                pkt_final = sum(flow.packets for flow in flow_query.lookup())
            return pkt_final - pkt_start
        return packet_sender
