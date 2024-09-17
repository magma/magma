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
from concurrent import futures
from unittest import TestCase, mock
from unittest.mock import MagicMock, patch

import fakeredis
import grpc
import orc8r.protos.state_pb2_grpc as state_pb2_grpc
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.redis.containers import RedisFlatDict
from magma.common.redis.serializers import (
    RedisSerde,
    get_json_deserializer,
    get_json_serializer,
    get_proto_deserializer,
    get_proto_serializer,
)
from magma.state.garbage_collector import GarbageCollector
from orc8r.protos.common_pb2 import NetworkID, Void
from orc8r.protos.service303_pb2 import LogVerbosity
from orc8r.protos.state_pb2_grpc import StateServiceStub

NID_TYPE = 'network_id'
LOG_TYPE = 'log_verbosity'
FOO_TYPE = 'foo'

# pylint: disable=protected-access


def get_mock_snowflake():
    return "aaa-bbb"


class DummyStateServer(state_pb2_grpc.StateServiceServicer):
    def __init__(self):
        pass

    def add_to_server(self, server):
        state_pb2_grpc.add_StateServiceServicer_to_server(self, server)

    def DeleteStates(self, request, context):
        for state_id in request.ids:
            if state_id.type == LOG_TYPE:
                raise grpc.RpcError("Test Exception")
        return Void()


class Foo:
    def __init__(self, bar: str, baz: int):
        self.bar = bar
        self.baz = baz


class GarbageCollectorTests(TestCase):
    def setUp(self):
        self.mock_redis = fakeredis.FakeStrictRedis()

        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)

        service = MagicMock()
        service.config = {
            # Replicate arbitrary orc8r protos
            'state_protos': [
                {
                    'proto_file': 'orc8r.protos.common_pb2',
                    'proto_msg': 'NetworkID',
                    'redis_key': NID_TYPE,
                    'state_scope': 'network',
                },
                {
                    'proto_file': 'orc8r.protos.service303_pb2',
                    'proto_msg': 'LogVerbosity',
                    'redis_key': LOG_TYPE,
                    'state_scope': 'gateway',
                },
            ],
            'json_state': [{'redis_key': FOO_TYPE, 'state_scope': 'gateway'}],
        }
        service.loop = self.loop

        # Bind the rpc server to a free port
        self._rpc_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10),
        )
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')
        # Add the servicer
        self._servicer = DummyStateServer()
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()
        # Create a rpc stub
        self.channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))

        serde1 = RedisSerde(
            NID_TYPE,
            get_proto_serializer(),
            get_proto_deserializer(NetworkID),
        )
        serde2 = RedisSerde(
            FOO_TYPE,
            get_json_serializer(),
            get_json_deserializer(),
        )
        serde3 = RedisSerde(
            LOG_TYPE,
            get_proto_serializer(),
            get_proto_deserializer(LogVerbosity),
        )

        self.nid_client = RedisFlatDict(self.mock_redis, serde1)
        self.foo_client = RedisFlatDict(self.mock_redis, serde2)
        self.log_client = RedisFlatDict(self.mock_redis, serde3)

        # Set up and start garbage collecting loop
        grpc_client_manager = GRPCClientManager(
            service_name="state",
            service_stub=StateServiceStub,
            max_client_reuse=60,
        )

        # mock the get_default_client function used to return the same
        # fakeredis object
        func_mock = mock.MagicMock(return_value=self.mock_redis)
        with patch('magma.state.redis_dicts.get_default_client', func_mock):
            # Start state garbage collection loop
            self.garbage_collector = GarbageCollector(
                service,
                grpc_client_manager,
            )

    def tearDown(self):
        self._rpc_server.stop(None)
        self.loop.close()

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_collect_states_to_delete(self):
        async def test():
            # Ensure setup is initialized properly
            self.nid_client.clear()
            self.foo_client.clear()
            self.log_client.clear()

            key = 'id1'
            self.nid_client[key] = NetworkID(id='foo')
            self.foo_client[key] = Foo("boo", 3)
            req = await self.garbage_collector._collect_states_to_delete()
            self.assertIsNone(req)

            self.nid_client.mark_as_garbage(key)
            self.foo_client.mark_as_garbage(key)
            req = await self.garbage_collector._collect_states_to_delete()
            self.assertEqual(2, len(req.ids))
            for state_id in req.ids:
                if state_id.type == NID_TYPE:
                    self.assertEqual('id1', state_id.deviceID)
                elif state_id.type == FOO_TYPE:
                    self.assertEqual('aaa-bbb:id1', state_id.deviceID)
                else:
                    self.fail("Unknown state type %s" % state_id.type)

            # Cleanup
            del self.foo_client[key]
            del self.nid_client[key]

        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_garbage_collect_success(self, get_rpc_mock):
        async def test():
            get_rpc_mock.return_value = self.channel
            self.nid_client.clear()
            self.foo_client.clear()
            self.log_client.clear()

            key = 'id1'
            foo = Foo("boo", 4)
            self.nid_client[key] = NetworkID(id='foo')
            self.foo_client[key] = foo
            self.nid_client.mark_as_garbage(key)
            self.foo_client.mark_as_garbage(key)
            req = await self.garbage_collector._collect_states_to_delete()
            self.assertEqual(2, len(req.ids))

            # Ensure all garbage collected objects get deleted from Redis
            await self.garbage_collector._send_to_state_service(req)
            self.assertEqual(0, len(self.nid_client.keys()))
            self.assertEqual(0, len(self.foo_client.keys()))
            self.assertEqual(0, len(self.nid_client.garbage_keys()))
            self.assertEqual(0, len(self.foo_client.garbage_keys()))

        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_garbage_collect_rpc_failure(self, get_rpc_mock):
        async def test():
            get_rpc_mock.return_value = self.channel
            self.nid_client.clear()
            self.foo_client.clear()
            self.log_client.clear()

            key = 'id1'
            self.nid_client[key] = NetworkID(id='foo')
            self.log_client[key] = LogVerbosity(verbosity=3)
            self.nid_client.mark_as_garbage(key)
            self.log_client.mark_as_garbage(key)

            req = await self.garbage_collector._collect_states_to_delete()
            self.assertEqual(2, len(req.ids))

            # Ensure objects on deleted from Redis on RPC failure
            await self.garbage_collector._send_to_state_service(req)
            self.assertEqual(0, len(self.nid_client.keys()))
            self.assertEqual(0, len(self.log_client.keys()))
            self.assertEqual(1, len(self.nid_client.garbage_keys()))
            self.assertEqual(1, len(self.log_client.garbage_keys()))

            # Cleanup
            del self.log_client[key]
            del self.nid_client[key]

        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_garbage_collect_with_state_update(self, get_rpc_mock):
        async def test():
            get_rpc_mock.return_value = self.channel
            self.nid_client.clear()
            self.foo_client.clear()
            self.log_client.clear()

            key = 'id1'
            foo = Foo("boo", 4)
            self.nid_client[key] = NetworkID(id='foo')
            self.foo_client[key] = foo
            self.nid_client.mark_as_garbage(key)
            self.foo_client.mark_as_garbage(key)
            req = await self.garbage_collector._collect_states_to_delete()
            self.assertEqual(2, len(req.ids))

            # Update one of the states, to ensure we don't delete valid state
            # from Redis
            expected = NetworkID(id='bar')
            self.nid_client[key] = expected

            # Ensure all garbage collected objects get deleted from Redis
            await self.garbage_collector._send_to_state_service(req)
            self.assertEqual(1, len(self.nid_client.keys()))
            self.assertEqual(0, len(self.foo_client.keys()))
            self.assertEqual(0, len(self.nid_client.garbage_keys()))
            self.assertEqual(0, len(self.foo_client.garbage_keys()))
            self.assertEqual(expected, self.nid_client[key])

        self.loop.run_until_complete(test())
