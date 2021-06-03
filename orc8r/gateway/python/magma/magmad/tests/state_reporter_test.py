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
import asyncio
import unittest.mock
from concurrent import futures
from unittest import TestCase
from unittest.mock import MagicMock

import grpc
import orc8r.protos.state_pb2_grpc as state_pb2_grpc
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.service_registry import ServiceRegistry
from magma.magmad.gateway_status import GatewayStatusFactory, SystemStatus
from magma.magmad.state_reporter import StateReporter
from orc8r.protos.service303_pb2 import GetOperationalStatesResponse, State
from orc8r.protos.state_pb2 import ReportStatesRequest, ReportStatesResponse
from orc8r.protos.state_pb2_grpc import StateServiceStub

# Allow access to protected variables for unit testing
# pylint: disable=protected-access
SR = 'magma.magmad.state_reporter'
GS = 'magma.magmad.gateway_status'
BM = 'magma.magmad._bootstrap_manager'
MS = 'magma.common.service'
SReg = 'magma.common.service_registry'


@staticmethod
def make_awaitable(func, return_value):
    future = asyncio.Future()
    future.set_result(return_value)
    func.return_value = future


class DummpyStateServer(state_pb2_grpc.StateServiceServicer):
    def __init__(self):
        pass

    def add_to_server(self, server):
        state_pb2_grpc.add_StateServiceServicer_to_server(self, server)

    def ReportStates(self, request, context):
        if request.states[0].type == "fail":
            raise grpc.RpcError("Test Exception")
        return ReportStatesResponse(
            unreportedStates=[],
        )


class StateReporterTests(TestCase):
    def setUp(self):
        ServiceRegistry.add_service('test1', '0.0.0.0', 0)
        ServiceRegistry.add_service('test2', '0.0.0.0', 0)
        ServiceRegistry.add_service('test3', '0.0.0.0', 0)

        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)

        service = MagicMock()
        service.config = {
            'magma_services': ['test1', 'test2', 'test3'],
            'non_service303_services': ['test2'],
            'skip_checkin_if_missing_meta_services': ['test3'],
        }
        service.mconfig.checkin_interval = 60
        service.mconfig.checkin_timeout = 30
        service.mconfig_metadata.created_at = 0
        service.version = "1.1.1.1"
        service.loop = self.loop

        # Bind the rpc server to a free port
        self._rpc_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
        )
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')
        # Add the servicer
        self._servicer = DummpyStateServer()
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()
        # Create a rpc stub
        self.channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))

        # Set up and start state reporting loop
        self.service_poller = unittest.mock.Mock()
        self.service_poller.service_info = unittest.mock.Mock()
        self.service_poller.service_info.return_value = {}

        # Mock out GatewayStatusFactory
        gateway_status_factory = \
            GatewayStatusFactory(service, self.service_poller, None)
        # Mock _system_status since it tried to read from /data/flash
        gateway_status_factory._system_status = unittest.mock.Mock()
        gateway_status_factory._system_status.return_value = \
            self.default_system_status()

        grpc_client_manager = GRPCClientManager(
            service_name="state",
            service_stub=StateServiceStub,
            max_client_reuse=60,
        )

        # Mock out bootstrap manager
        bootstrap_manager = unittest.mock.Mock()
        bootstrap_manager.schedule_bootstrap_now = unittest.mock.Mock()
        bootstrap_manager.schedule_bootstrap_now.return_value = None

        self.state_reporter = StateReporter(
            config=service.config,
            mconfig=service.mconfig,
            loop=service.loop,
            bootstrap_manager=bootstrap_manager,
            gw_status_factory=gateway_status_factory,
            grpc_client_manager=grpc_client_manager,
        )
        self.state_reporter.FAIL_THRESHOLD = 0

        self.state_reporter.start()

    def tearDown(self):
        self._rpc_server.stop(None)
        self.service_poller.stop()
        self.state_reporter.stop()
        self.loop.close()

    def test__setup(self):
        # Checks that it gets magma_services with service303 interfaces
        # correctly
        service_info_by_name = self.state_reporter._service_info_by_name
        self.assertIsNotNone(service_info_by_name.get("test1"))
        self.assertIsNotNone(service_info_by_name.get("test3"))
        self.assertIsNone(service_info_by_name.get("test2"))

    @staticmethod
    def _construct_operational_state_mock(device_id: str) \
            -> unittest.mock.Mock:
        mock = unittest.mock.Mock()
        future = asyncio.Future()
        future.set_result(
            GetOperationalStatesResponse(
                states=[
                    State(
                        type="test",
                        deviceID=device_id,
                        value="hello!".encode('utf-8'),
                    ),
                ],
            ),
        )
        mock.GetOperationalStates.future.side_effect = [future]
        return mock

    @unittest.mock.patch('%s.Service303Stub' % MS)
    def test__collect_states_missing_meta(self, service_303_mock):
        async def test():
            # Mock out GerOperationalStates.future
            mock_service_1 = self._construct_operational_state_mock("test1")
            mock_service_2 = self._construct_operational_state_mock("test3")
            service_303_mock.side_effect = [mock_service_1, mock_service_2]

            # service info is empty, it should not return a gateway state
            self.service_poller.service_info = {}
            result = await self.state_reporter._collect_states()
            self.assertIsNotNone(result)
            self.assertEqual(len(result.states), 2)
            self.assertEqual(result.states[0].deviceID, "test1")
            self.assertEqual(result.states[1].deviceID, "test3")
            self.assertEqual(
                self.state_reporter._error_handler.num_skipped_gateway_states,
                1,
            )

        # Cancel the reporter's loop so there are no other activities
        self.state_reporter._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @unittest.mock.patch('%s.Service303Stub' % MS)
    def test__collect_states_success(self, service_303_mock):
        async def test():
            # Mock out GerOperationalStates.future
            mock_service_1 = self._construct_operational_state_mock("test1")
            mock_service_2 = self._construct_operational_state_mock("test3")
            service_303_mock.side_effect = [mock_service_1, mock_service_2]

            # populate service info to include the required test3
            service_303_mock.side_effect = [mock_service_1, mock_service_2]
            test3_service_info = unittest.mock.Mock()
            test3_service_info.status = unittest.mock.Mock()
            test3_service_info.status.meta = {"test3": "m e t a"}
            self.service_poller.service_info = {"test3": test3_service_info}
            result = await self.state_reporter._collect_states()
            self.assertIsNotNone(result)

            self.assertIsNotNone(result)
            self.assertEqual(len(result.states), 3)
            self.assertEqual(result.states[0].deviceID, 'test1')
            self.assertEqual(result.states[1].deviceID, 'test3')
            self.assertEqual(result.states[2].type, 'gw_state')

        # Cancel the reporter's loop so there are no other activities
        self.state_reporter._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @unittest.mock.patch('%s.ServiceRegistry.get_proxy_config' % SReg)
    @unittest.mock.patch('%s.cert_is_invalid' % SR)
    @unittest.mock.patch('%s.ServiceRegistry.get_rpc_channel' % SReg)
    def test__report_states_failure(
        self, get_rpc_mock, cert_is_invalid_mock,
        get_proxy_config_mock,
    ):
        bootstrap_manager = \
            self.state_reporter._error_handler._bootstrap_manager
        bootstrap_manager.schedule_bootstrap_now = unittest.mock.Mock()

        async def test():
            # force bootstrap to be called on the first error
            self.state_reporter._error_handler.fail_threshold = 0
            # use dummy state servicer
            get_rpc_mock.return_value = self.channel
            # mock out get_proxy_config
            get_proxy_config_mock.return_value = {
                "cloud_address": 1,
                "cloud_port": 2,
                "gateway_cert": 3,
                "gateway_key": 4,
            }
            # force cert_is_invalid to return false
            future = asyncio.Future()
            future.set_result(True)
            cert_is_invalid_mock.return_value = future
            # make schedule_bootstrap_now awaitable
            future = asyncio.Future()
            future.set_result(None)
            bootstrap_manager.schedule_bootstrap_now.return_value = future

            request = ReportStatesRequest(states=[State(type="fail")])
            await self.state_reporter._send_to_state_service(request)

        # Cancel the reporter's loop so there are no other activities
        self.state_reporter._periodic_task.cancel()
        self.loop.run_until_complete(test())
        # the failure threshold is 0 so assert _schedule_bootstrap_now is
        # called
        bootstrap_manager.schedule_bootstrap_now.assert_has_calls(
            [unittest.mock.call()],
        )

    @staticmethod
    def default_system_status():
        return SystemStatus(
            time=1,
            uptime_secs=2,
            cpu_user=3,
            cpu_system=4,
            cpu_idle=5,
            mem_total=6,
            mem_available=7,
            mem_used=8,
            mem_free=9,
            swap_total=10,
            swap_used=11,
            swap_free=12,
            disk_partitions=[],
        )
