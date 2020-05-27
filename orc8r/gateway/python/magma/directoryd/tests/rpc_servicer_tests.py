"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from concurrent import futures
import grpc
from unittest import TestCase, mock
from magma.common.redis.mocks.mock_redis import MockRedis, MockUnavailableRedis
from magma.directoryd.rpc_servicer import GatewayDirectoryServiceRpcServicer
from orc8r.protos.common_pb2 import Void
from orc8r.protos.directoryd_pb2 import UpdateRecordRequest, \
    DeleteRecordRequest, GetDirectoryFieldRequest
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub

# Allow access to protected variables for unit testing
# pylint: disable=protected-access


def get_mock_snowflake():
    return "aaa-bbb"


class DirectorydRpcServiceTests(TestCase):
    @mock.patch("redis.Redis", MockRedis)
    def setUp(self):
        # Bind the rpc server to a free port
        thread_pool = futures.ThreadPoolExecutor(max_workers=10)
        self._rpc_server = grpc.server(thread_pool)
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')

        # Add the servicer
        self._servicer = GatewayDirectoryServiceRpcServicer()
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()

        # Create a rpc stub
        channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))
        self._stub = GatewayDirectoryServiceStub(channel)

    @mock.patch("redis.Redis", MockRedis)
    def tearDown(self):
        self._rpc_server.stop(None)

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_update_record(self):
        self._servicer._redis_dict.clear()

        req = UpdateRecordRequest()
        req.id = "IMSI555"
        self._stub.UpdateRecord(req)
        actual_record = self._servicer._redis_dict[req.id]
        self.assertEqual(actual_record.location_history, ['aaa-bbb'])
        self.assertEqual(actual_record.identifiers, {})

        req.fields["mac_addr"] = "aa:aa:bb:bb:cc:cc"
        req.fields["ipv4_addr"] = "192.168.172.12"

        self._stub.UpdateRecord(req)
        actual_record2 = self._servicer._redis_dict[req.id]
        self.assertEqual(actual_record2.location_history, ["aaa-bbb"])
        self.assertEqual(actual_record2.identifiers['mac_addr'],
                         "aa:aa:bb:bb:cc:cc")
        self.assertEqual(actual_record2.identifiers['ipv4_addr'],
                         "192.168.172.12")

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_update_record_bad_location(self):
        self._servicer._redis_dict.clear()

        req = UpdateRecordRequest()
        req.id = "IMSI556"
        req.location = "bbb-ccc"

        self._stub.UpdateRecord(req)
        actual_record = self._servicer._redis_dict[req.id]
        self.assertEqual(actual_record.location_history, ['aaa-bbb'])
        self.assertEqual(actual_record.identifiers, {})

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_delete_record(self):
        self._servicer._redis_dict.clear()

        req = UpdateRecordRequest()
        req.id = "IMSI557"
        self._stub.UpdateRecord(req)
        self.assertTrue(req.id in self._servicer._redis_dict)

        del_req = DeleteRecordRequest()
        del_req.id = "IMSI557"
        self._stub.DeleteRecord(del_req)
        self.assertFalse(req.id in self._servicer._redis_dict)

        with self.assertRaises(grpc.RpcError) as err:
            self._stub.DeleteRecord(del_req)
        self.assertEqual(err.exception.code(), grpc.StatusCode.NOT_FOUND)

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_get_field(self):
        self._servicer._redis_dict.clear()

        req = UpdateRecordRequest()
        req.id = "IMSI557"
        req.fields["mac_addr"] = "aa:bb:aa:bb:aa:bb"
        self._stub.UpdateRecord(req)
        self.assertTrue(req.id in self._servicer._redis_dict)

        get_req = GetDirectoryFieldRequest()
        get_req.id = "IMSI557"
        get_req.field_key = "mac_addr"
        ret = self._stub.GetDirectoryField(get_req)
        self.assertEqual("aa:bb:aa:bb:aa:bb", ret.value)

        with self.assertRaises(grpc.RpcError) as err:
            get_req.field_key = "ipv4_addr"
            self._stub.GetDirectoryField(get_req)
        self.assertEqual(err.exception.code(), grpc.StatusCode.NOT_FOUND)

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_get_all(self):
        self._servicer._redis_dict.clear()

        req = UpdateRecordRequest()
        req.id = "IMSI557"
        req.fields["mac_addr"] = "aa:bb:aa:bb:aa:bb"
        self._stub.UpdateRecord(req)
        self.assertTrue(req.id in self._servicer._redis_dict)

        req2 = UpdateRecordRequest()
        req2.id = "IMSI556"
        req2.fields["ipv4_addr"] = "192.168.127.11"
        self._stub.UpdateRecord(req2)
        self.assertTrue(req2.id in self._servicer._redis_dict)

        void_req = Void()
        ret = self._stub.GetAllDirectoryRecords(void_req)
        self.assertEqual(2, len(ret.records))
        for record in ret.records:
            if record.id == "IMSI556":
                self.assertEqual(record.fields["ipv4_addr"], "192.168.127.11")
            elif record.id == "IMSI557":
                self.assertEqual(record.fields["mac_addr"],
                                 "aa:bb:aa:bb:aa:bb")
            else:
                raise AssertionError()

    @mock.patch("redis.Redis", MockUnavailableRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_redis_unavailable(self):
        self._servicer._redis_dict = MockUnavailableRedis("localhost", 6380)
        req = UpdateRecordRequest()
        req.id = "IMSI557"
        req.fields["mac_addr"] = "aa:bb:aa:bb:aa:bb"

        with self.assertRaises(grpc.RpcError) as err:
            self._stub.UpdateRecord(req)
        self.assertEqual(err.exception.code(), grpc.StatusCode.UNAVAILABLE)

        get_req = GetDirectoryFieldRequest()
        get_req.id = "IMSI557"
        get_req.field_key = "mac_addr"
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.GetDirectoryField(get_req)
        self.assertEqual(err.exception.code(), grpc.StatusCode.UNAVAILABLE)

        del_req = DeleteRecordRequest()
        del_req.id = "IMSI557"
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.DeleteRecord(del_req)
        self.assertEqual(err.exception.code(), grpc.StatusCode.UNAVAILABLE)

        void_req = Void()
        with self.assertRaises(grpc.RpcError) as err:
            self._stub.GetAllDirectoryRecords(void_req)
        self.assertEqual(err.exception.code(), grpc.StatusCode.UNAVAILABLE)
