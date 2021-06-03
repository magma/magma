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
import json
from concurrent import futures
from unittest import TestCase, mock
from unittest.mock import MagicMock

import fakeredis
import grpc
import jsonpickle
import orc8r.protos.state_pb2_grpc as state_pb2_grpc
from google.protobuf.json_format import MessageToDict
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
from magma.state.keys import make_mem_key
from magma.state.state_replicator import StateReplicator
from orc8r.protos.common_pb2 import IDList, NetworkID
from orc8r.protos.service303_pb2 import LogVerbosity
from orc8r.protos.state_pb2 import (
    IDAndError,
    IDAndVersion,
    ReportStatesResponse,
    SyncStatesResponse,
)
from orc8r.protos.state_pb2_grpc import StateServiceStub

NID_TYPE = 'network_id'
IDList_TYPE = 'id_list'
LOG_TYPE = 'log_verbosity'
FOO_TYPE = 'foo'

# Allow access to protected variables for unit testing
# pylint: disable=protected-access
CS = "magma.state.state_replicator._collect_states_to_replicate"
RS = "magma.state.state_replicator._resync"
SS = "magma.state.state_replicator._send_to_state_service"
SV = "magma.state.state_replicator._state_versions"


def get_mock_snowflake():
    return "aaa-bbb"


class DummyStateServer(state_pb2_grpc.StateServiceServicer):
    def __init__(self):
        pass

    def add_to_server(self, server):
        state_pb2_grpc.add_StateServiceServicer_to_server(self, server)

    def ReportStates(self, request, context):
        unreported_states = []
        for state in request.states:
            # Always 'fail' to report LOG_TYPE states
            if state.type == LOG_TYPE:
                id_and_error = IDAndError(
                    type=state.type,
                    deviceID=state.deviceID,
                    error="mocked_error",
                )
                unreported_states.append(id_and_error)
        return ReportStatesResponse(
            unreportedStates=unreported_states,
        )

    def SyncStates(self, request, context):
        unsynced_states = []
        for state in request.states:
            if state.id.type == LOG_TYPE:
                raise grpc.RpcError("Test Exception")
            elif state.id.type == NID_TYPE:
                id_and_version = IDAndVersion(
                    id=state.id,
                    version=state.version,
                )
                unsynced_states.append(id_and_version)
        return SyncStatesResponse(
            unsyncedStates=unsynced_states,
        )


class Foo:
    def __init__(self, bar: str, baz: int):
        self.bar = bar
        self.baz = baz


class StateReplicatorTests(TestCase):
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
                    'proto_file': 'orc8r.protos.common_pb2',
                    'proto_msg': 'IDList',
                    'redis_key': IDList_TYPE,
                    'state_scope': 'gateway',
                },
                {
                    'proto_file': 'orc8r.protos.service303_pb2',
                    'proto_msg': 'LogVerbosity',
                    'redis_key': LOG_TYPE,
                    'state_scope': 'gateway',
                },
            ],
            'json_state': [{'redis_key': FOO_TYPE, 'state_scope': 'network'}],
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
            IDList_TYPE,
            get_proto_serializer(),
            get_proto_deserializer(IDList),
        )
        serde3 = RedisSerde(
            LOG_TYPE,
            get_proto_serializer(),
            get_proto_deserializer(LogVerbosity),
        )
        serde4 = RedisSerde(
            FOO_TYPE,
            get_json_serializer(),
            get_json_deserializer(),
        )

        self.nid_client = RedisFlatDict(self.mock_redis, serde1)
        self.idlist_client = RedisFlatDict(self.mock_redis, serde2)
        self.log_client = RedisFlatDict(self.mock_redis, serde3)
        self.foo_client = RedisFlatDict(self.mock_redis, serde4)

        # Set up and start state replicating loop
        grpc_client_manager = GRPCClientManager(
            service_name="state",
            service_stub=StateServiceStub,
            max_client_reuse=60,
        )

        # mock the get_default_client function used to return the same
        # fakeredis object
        func_mock = mock.MagicMock(return_value=self.mock_redis)
        with mock.patch(
                'magma.state.redis_dicts.get_default_client',
                func_mock,
        ):
            garbage_collector = GarbageCollector(service, grpc_client_manager)

            self.state_replicator = StateReplicator(
                service=service,
                garbage_collector=garbage_collector,
                grpc_client_manager=grpc_client_manager,
            )
        self.state_replicator.start()

    def tearDown(self):
        self._rpc_server.stop(None)
        self.state_replicator.stop()
        self.loop.close()

    def convert_msg_to_state(self, redis_state, is_proto=True):
        if is_proto:
            json_converted_state = MessageToDict(redis_state)
            serialized_json_state = json.dumps(json_converted_state)
        else:
            serialized_json_state = jsonpickle.encode(redis_state)
        return serialized_json_state.encode("utf-8")

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_collect_states_to_replicate(self):
        async def test():
            # Ensure setup is initialized properly
            self.nid_client.clear()
            self.idlist_client.clear()
            self.log_client.clear()
            self.foo_client.clear()

            key = 'id1'

            self.nid_client[key] = NetworkID(id='foo')
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])
            self.foo_client[key] = Foo("boo", 3)

            exp1 = self.convert_msg_to_state(self.nid_client[key])
            exp2 = self.convert_msg_to_state(self.idlist_client[key])
            exp3 = self.convert_msg_to_state(self.foo_client[key], False)

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(3, len(req.states))
            for state in req.states:
                if state.type == NID_TYPE:
                    self.assertEqual('id1', state.deviceID)
                    self.assertEqual(1, state.version)
                    self.assertEqual(exp1, state.value)
                elif state.type == IDList_TYPE:
                    self.assertEqual('aaa-bbb:id1', state.deviceID)
                    self.assertEqual(1, state.version)
                    self.assertEqual(exp2, state.value)
                elif state.type == FOO_TYPE:
                    self.assertEqual('id1', state.deviceID)
                    self.assertEqual(1, state.version)
                    self.assertEqual(exp3, state.value)
                else:
                    self.fail("Unknown state type %s" % state.type)

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_replicate_states_success(self, get_rpc_mock):
        async def test():
            get_rpc_mock.return_value = self.channel

            # Ensure setup is initialized properly
            self.nid_client.clear()
            self.idlist_client.clear()
            self.log_client.clear()
            self.foo_client.clear()

            key = 'id1'
            foo = Foo("boo", 4)
            self.nid_client[key] = NetworkID(id='foo')
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])
            self.foo_client[key] = foo
            # Increment version
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(3, len(req.states))

            # Ensure in-memory map updates properly
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(3, len(self.state_replicator._state_versions))
            mem_key1 = make_mem_key('id1', NID_TYPE)
            mem_key2 = make_mem_key(
                'aaa-bbb:id1',
                IDList_TYPE,
            )
            mem_key3 = make_mem_key('id1', FOO_TYPE)
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key1],
            )
            self.assertEqual(
                2,
                self.state_replicator._state_versions[mem_key2],
            )
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key3],
            )

            # Now add new state and update some existing state
            key2 = 'id2'
            self.nid_client[key2] = NetworkID(id='bar')
            self.idlist_client[key] = IDList(ids=['bar', 'foo'])
            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(2, len(req.states))

            # Ensure in-memory map updates properly
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(4, len(self.state_replicator._state_versions))
            mem_key4 = make_mem_key('id2', NID_TYPE)
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key1],
            )
            self.assertEqual(
                3,
                self.state_replicator._state_versions[mem_key2],
            )
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key3],
            )
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key4],
            )

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_unreplicated_states(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel

            # Add initial state to be replicated
            self.nid_client.clear()
            self.idlist_client.clear()
            self.log_client.clear()
            self.foo_client.clear()

            key = 'id1'
            key2 = 'id2'
            self.nid_client[key] = NetworkID(id='foo')
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])
            # Increment version
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])
            # Set state that will be 'unreplicated'
            self.log_client[key2] = LogVerbosity(verbosity=5)

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(3, len(req.states))

            # Ensure in-memory map updates properly for successful replications
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(2, len(self.state_replicator._state_versions))
            mem_key1 = make_mem_key('id1', NID_TYPE)
            mem_key2 = make_mem_key(
                'aaa-bbb:id1',
                IDList_TYPE,
            )
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key1],
            )
            self.assertEqual(
                2,
                self.state_replicator._state_versions[mem_key2],
            )

            # Now run again, ensuring only the state the wasn't replicated
            # will be sent again
            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(1, len(req.states))
            self.assertEqual('aaa-bbb:id2', req.states[0].deviceID)
            self.assertEqual(LOG_TYPE, req.states[0].type)

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_resync_success(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel
            self.nid_client.clear()
            self.idlist_client.clear()
            self.log_client.clear()
            self.foo_client.clear()

            key = 'id1'
            # Set state that will be 'unsynced'
            self.nid_client[key] = NetworkID(id='foo')
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])
            # Increment state's version
            self.idlist_client[key] = IDList(ids=['bar', 'blah'])

            await self.state_replicator._resync()
            self.assertEqual(True, self.state_replicator._has_resync_completed)
            self.assertEqual(1, len(self.state_replicator._state_versions))
            mem_key = make_mem_key(
                'aaa-bbb:id1',
                IDList_TYPE,
            )
            self.assertEqual(2, self.state_replicator._state_versions[mem_key])

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_resync_failure(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel
            self.nid_client.clear()
            self.idlist_client.clear()
            self.log_client.clear()
            self.foo_client.clear()

            # Set state that will trigger the RpcError
            log_key = 'id1'
            self.log_client[log_key] = LogVerbosity(verbosity=5)

            try:
                await self.state_replicator._resync()
            except grpc.RpcError:
                pass

            self.assertEqual(
                False,
                self.state_replicator._has_resync_completed,
            )
            self.assertEqual(0, len(self.state_replicator._state_versions))

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_deleted_replicated_state(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel
            self.nid_client.clear()
            self.idlist_client.clear()
            self.log_client.clear()
            self.foo_client.clear()

            key = 'id1'
            self.nid_client[key] = NetworkID(id='foo')
            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(1, len(req.states))

            # Ensure in-memory map updates properly
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(1, len(self.state_replicator._state_versions))
            mem_key1 = make_mem_key('id1', NID_TYPE)
            self.assertEqual(
                1,
                self.state_replicator._state_versions[mem_key1],
            )

            # Now delete state and ensure in-memory map gets updated properly
            del self.nid_client[key]
            req = await self.state_replicator._collect_states_to_replicate()
            self.assertIsNone(req)

            await self.state_replicator._cleanup_deleted_keys()
            self.assertFalse(key in self.state_replicator._state_versions)

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())
