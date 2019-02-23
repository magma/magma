"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import unittest
import unittest.mock

from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2 import ServiceInfo

from magma.common.service_registry import ServiceRegistry
from magma.magmad.service_poller import ServicePoller


class MockFuture(object):
    def __init__(self, is_error):
        self._is_error = is_error

    def exception(self):
        if self._is_error:
            return self.MockException()
        return None

    def result(self):
        return ServiceInfo()

    class MockException(object):
        def details(self):
            return ''

        def code(self):
            return 0


class ServicePollerTests(unittest.TestCase):
    """
    Tests for the ServicePoller
    """
    def setUp(self):
        ServiceRegistry.add_service('test1', '0.0.0.0', 0)
        ServiceRegistry.add_service('test2', '0.0.0.0', 0)
        config = {
            'magma_services': ['test1', 'test2'],
            'non_service303_services': ['test2']
        }
        self._loop = asyncio.new_event_loop()
        self._service_poller = ServicePoller(self._loop, config)

    @unittest.mock.patch('magma.magmad.service_poller.Service303Stub')
    @unittest.mock.patch('magma.configuration.service_configs')
    def test_poll(self, _service_configs_mock, service303_mock):
        """
        Test if the query to Service303 succeeds.
        """
        # Mock out GetServiceInfo.future
        mock = unittest.mock.Mock()
        mock.GetServiceInfo.future.side_effect = [unittest.mock.Mock()]
        service303_mock.side_effect = [mock]

        self._service_poller.start()
        mock.GetServiceInfo.future.assert_called_once_with(
            Void(), self._service_poller.GET_STATUS_TIMEOUT)
        # pylint: disable=protected-access
        self._service_poller._get_service_info_done('test1', MockFuture(False))

    @unittest.mock.patch('magma.magmad.service_poller.Service303Stub')
    @unittest.mock.patch('magma.configuration.service_configs')
    def test_poll_exception(self, _service_configs_mock, service303_mock):
        """
        Test if the query to Service303 fails and handled gracefully.
        """
        # Mock out GetServiceInfo.future
        mock = unittest.mock.Mock()
        mock.GetServiceInfo.future.side_effect = [unittest.mock.Mock()]
        service303_mock.side_effect = [mock]

        self._service_poller.start()
        mock.GetServiceInfo.future.assert_called_once_with(
            Void(), self._service_poller.GET_STATUS_TIMEOUT)
        # pylint: disable=protected-access
        self._service_poller._get_service_info_done('test1', MockFuture(True))


if __name__ == "__main__":
    unittest.main()
