"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=protected-access

import asyncio
import unittest
import unittest.mock

from orc8r.protos.logging_service_pb2 import LogEntry, LogRequest, \
    LoggerDestination

from magma.common.scribe_client import ScribeClient
from magma.common.service_registry import ServiceRegistry


class ScribeClientTests(unittest.TestCase):
    """
    Tests for the ScribeClient
    """
    def setUp(self):
        self._scribe_client = ScribeClient(loop=asyncio.new_event_loop())
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        ServiceRegistry._PROXY_CONFIG = {'local_port': 1234,
                                         'cloud_address': 'test',
                                         'proxy_cloud_connections': True}

    @unittest.mock.patch('magma.common.scribe_client.LoggingServiceStub')
    def test_log_entries_to_dest(self, logging_service_mock_stub):
        """
        Test if the service starts and stops gracefully.
        """
        # mock out Log.future
        mock = unittest.mock.Mock()
        mock.Log.future.side_effect = [unittest.mock.Mock()]
        logging_service_mock_stub.side_effect = [mock]
        data = {}
        data['int'] = {"some_field": 456}
        data['normal'] = {"imsi": "IMSI11111111", "ue_state": "IDLE"}
        entries = [LogEntry(
            category="test_category",
            int_map={"some_field": 456},
            normal_map={"imsi": "IMSI11111111", "ue_state": "IDLE"},
            time=12345,
        )]
        self._scribe_client.log_to_scribe_with_sampling_rate(entries, 0)
        mock.Log.future.assert_not_called()

        self._scribe_client.log_to_scribe_with_sampling_rate(entries, 1)
        mock.Log.future.assert_called_once_with(
            LogRequest(
                Entries=entries,
                Destination=LoggerDestination.Value("SCRIBE")
            )
        )


if __name__ == "__main__":
    unittest.main()
