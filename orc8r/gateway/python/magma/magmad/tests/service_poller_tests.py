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
import unittest
import unittest.mock

import grpc
from magma.common.service_registry import ServiceRegistry
from magma.magmad.service_poller import ServicePoller
from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2 import ServiceInfo

# Allow access to protected variables for unit testing
# pylint: disable=protected-access
SP = "magma.magmad.service_poller"


class ServicePollerTests(unittest.TestCase):
    """
    Tests for the ServicePoller
    """

    def setUp(self):
        ServiceRegistry.add_service('test1', '0.0.0.0', 0)
        ServiceRegistry.add_service('test2', '0.0.0.0', 0)
        config = {
            'magma_services': ['test1', 'test2'],
            'non_service303_services': ['test2'],
        }
        self._loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self._loop)
        self._service_poller = ServicePoller(self._loop, config)

    @unittest.mock.patch('%s.Service303Stub' % SP)
    def test_poll(self, service303_mock):
        """
        Test if the query to Service303 succeeds.
        """
        async def test():
            # Mock out GetServiceInfo.future
            mock = unittest.mock.Mock()
            service_info_future = asyncio.Future()
            service_info_future.set_result(ServiceInfo())
            mock.GetServiceInfo.future.side_effect = [service_info_future]
            service303_mock.side_effect = [mock]

            await self._service_poller._get_service_info()

            mock.GetServiceInfo.future.assert_called_once_with(
                Void(), self._service_poller.GET_STATUS_TIMEOUT,
            )
        self._loop.run_until_complete(test())

    @unittest.mock.patch('%s.Service303Stub' % SP)
    def test_poll_exception(self, service303_mock):
        """
        Test if the query to Service303 fails and handled gracefully.
        """
        def fake_add_done(_):
            grpc_err = grpc.RpcError()
            grpc_err.code = lambda: grpc.StatusCode.UNKNOWN
            grpc_err.details = lambda: "Test Exception"
            raise grpc_err

        async def test():
            # Mock out GetServiceInfo.future
            mock = unittest.mock.Mock()
            service_info_future = asyncio.Future()
            # Force an exception to happen
            service_info_future.add_done_callback = fake_add_done
            mock.GetServiceInfo.future.side_effect = [service_info_future]
            service303_mock.side_effect = [mock]

            await self._service_poller._get_service_info()

            mock.GetServiceInfo.future.assert_called_once_with(
                Void(), self._service_poller.GET_STATUS_TIMEOUT,
            )
        self._loop.run_until_complete(test())


if __name__ == "__main__":
    unittest.main()
