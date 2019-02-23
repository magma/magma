"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
from unittest import TestCase, main, mock

from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2 import ServiceInfo
from orc8r.protos.service303_pb2_grpc import Service303Stub

from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry


class Service303Tests(TestCase):
    """
    Tests for the MagmaService and the Service303 interface
    """

    @mock.patch('time.time', mock.MagicMock(return_value=12345))
    def setUp(self):
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        self._stub = None

        # Use a new event loop to ensure isolated tests
        self._service = MagmaService('test', loop=asyncio.new_event_loop())
        # Clear the global event loop so tests rely only on the event loop that
        # was manually set
        asyncio.set_event_loop(None)

    def test_service_run(self):
        """
        Test if the service starts and stops gracefully.
        """
        self.assertEqual(self._service.state, ServiceInfo.STARTING)

        # Start the service and pause the loop
        self._service.loop.stop()
        self._service.run()
        self.assertEqual(self._service.state, ServiceInfo.ALIVE)

        # Create a rpc stub and query the Service303 interface
        ServiceRegistry.add_service('test', '0.0.0.0', self._service.port)
        channel = ServiceRegistry.get_rpc_channel('test', ServiceRegistry.LOCAL)
        self._stub = Service303Stub(channel)

        info = ServiceInfo(name='test',
                           version='0.0.0',
                           state=ServiceInfo.ALIVE,
                           health=ServiceInfo.APP_HEALTHY,
                           start_time_secs=12345)
        self.assertEqual(self._stub.GetServiceInfo(Void()), info)

        # Stop the service
        self._stub.StopService(Void())
        self._service.loop.run_forever()
        self.assertEqual(self._service.state, ServiceInfo.STOPPED)


if __name__ == "__main__":
    main()
