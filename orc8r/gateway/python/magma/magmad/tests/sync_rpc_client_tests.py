"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=protected-access

import asyncio
import unittest.mock

from orc8r.protos.sync_rpc_service_pb2 import GatewayRequest, GatewayResponse, \
    SyncRPCResponse

from magma.common.service_registry import ServiceRegistry
from magma.magmad.sync_rpc_client import SyncRPCClient


class MockFuture(object):
    is_error = True

    def __init__(self, is_error, expected_result, expected_err_msg):
        self.is_error = is_error
        self.expected_result = expected_result
        self.expected_err_msg = expected_err_msg

    def exception(self):
        if self.is_error:
            return self.MockException(self.expected_err_msg)

    def result(self):
        if not self.is_error:
            return self.expected_result

    class MockException(object):
        def __init__(self, msg):
            self._msg = msg

        def details(self):
            return self._msg

        def code(self):
            return 0

        def __str__(self):
            return self._msg


class SyncRPCClientTests(unittest.TestCase):
    """
    Tests for the SyncRPCClient
    """
    def setUp(self):
        loop = asyncio.new_event_loop()
        asyncio.set_event_loop(loop)
        self._loop = loop
        self._sync_rpc_client = SyncRPCClient(loop=loop, response_timeout=3)
        ServiceRegistry.add_service('test', '0.0.0.0', 0)
        ServiceRegistry._PROXY_CONFIG = {'local_port': 2345,
                                         'cloud_address': 'test',
                                         'proxy_cloud_connections': True}
        self._req_body = GatewayRequest(gwId="test id", authority='mobility',
                                        path='/magma.MobilityService'
                                             '/ListAddedIPv4Blocks',
                                        headers={'te': 'trailers',
                                                 'content-type':
                                                     'application/grpc',
                                                 'user-agent':
                                                     'grpc-python/1.4.0',
                                                 'grpc-accept-encoding':
                                                     'identity'},
                                        payload=bytes.fromhex('0000000000'))
        self._expected_resp = GatewayResponse(status="400",
                                              headers={"test_key": "test_val"},
                                              payload=b'\x00'
                                                      b'\x00\x00\x00\n\n\x08')
        self._expected_err_msg = "test error"

    def test_send_request_done(self):
        """
        Test if send_request_done puts the right SyncRPCResponse in
        response_queue based on the result of future.
        Returns: None

        """
        self.assertTrue(self._sync_rpc_client._response_queue.empty())
        req_id = 356
        future = MockFuture(is_error=False,
                            expected_result=self._expected_resp,
                            expected_err_msg="")
        # future has result, send_request_done should enqueue a SyncRPCResponse
        # with the req_id, and the expected_resp set in MockFuture
        self._sync_rpc_client.send_request_done(req_id, future)
        self.assertFalse(self._sync_rpc_client._response_queue.empty())
        res = self._sync_rpc_client._response_queue.get(block=False)
        self.assertEqual(res, SyncRPCResponse(reqId=req_id,
                                              respBody=self._expected_resp,
                                              heartBeat=False))

        self.assertTrue(self._sync_rpc_client._response_queue.empty())
        req_id = 234
        expected_err_resp = GatewayResponse(err=self._expected_err_msg)
        future = MockFuture(is_error=True, expected_result=GatewayResponse(),
                            expected_err_msg=self._expected_err_msg)
        self._sync_rpc_client.send_request_done(req_id, future)

        res = self._sync_rpc_client._response_queue.get(block=False)
        self.assertEqual(res, SyncRPCResponse(reqId=req_id,
                                              respBody=expected_err_resp,
                                              heartBeat=False))

    def test_send_sync_rpc_response(self):
        expected = SyncRPCResponse(reqId=123, respBody=self._expected_resp)
        self._sync_rpc_client._response_queue.put(expected)
        res = self._sync_rpc_client.send_sync_rpc_response()
        actual = next(res)
        self.assertEqual(expected, actual)
        expected = SyncRPCResponse(heartBeat=True)
        actual = next(res)
        self.assertEqual(expected, actual)

    def test_retry_connect_sleep(self):
        self._sync_rpc_client._current_delay = 0
        for i in range(5):
            self._sync_rpc_client._set_connect_time()
            self._sync_rpc_client._retry_connect_sleep()
            if i == 4:
                self.assertEqual(self._sync_rpc_client.RETRY_MAX_DELAY_SECS,
                                 self._sync_rpc_client._current_delay)
            else:
                self.assertEqual(2 ** i, self._sync_rpc_client._current_delay)


if __name__ == "__main__":
    unittest.main()
