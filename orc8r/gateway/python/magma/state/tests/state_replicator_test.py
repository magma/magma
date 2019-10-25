"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
from unittest import TestCase, mock
import grpc
import json
from concurrent import futures
import orc8r.protos.state_pb2_grpc as state_pb2_grpc
from orc8r.protos.state_pb2 import ReportStatesResponse, \
    SyncStatesResponse, IDAndVersion, IDAndError
from unittest.mock import MagicMock
from orc8r.protos.service303_pb2 import LogVerbosity
from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict
from magma.common.redis.serializers import get_proto_deserializer, \
    get_proto_serializer, get_json_deserializer, get_json_serializer, \
    RedisSerde
from magma.state.state_replicator import StateReplicator
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.redis.mocks.mock_redis import MockRedis
from orc8r.protos.state_pb2_grpc import StateServiceStub
from orc8r.protos.common_pb2 import NetworkID, IDList
from google.protobuf.json_format import MessageToDict

NID_TYPE = 'network_id'
IDList_TYPE = 'id_list'
LOG_TYPE = 'log_verbosity'
FOO_TYPE = 'foo'

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
                id_and_error = IDAndError(type=state.type,
                                          deviceID=state.deviceID,
                                          error="mocked_error")
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
                id_and_version = IDAndVersion(id=state.id,
                                              version=state.version)
                unsynced_states.append(id_and_version)
        return SyncStatesResponse(
            unsyncedStates=unsynced_states,
        )

class StateReplicatorTests(TestCase):
    @mock.patch("redis.Redis", MockRedis)
    def setUp(self):

        self.loop = asyncio.new_event_loop()
        asyncio.set_event_loop(self.loop)

        service = MagicMock()
        service.config = {
            # Replicate arbitrary orc8r protos
            'state_protos': [{'proto_file': 'orc8r.protos.common_pb2',
                              'proto_msg': 'NetworkID',
                              'redis_key': NID_TYPE,
                              'state_scope': 'network'},
                             {'proto_file': 'orc8r.protos.common_pb2',
                              'proto_msg': 'IDList',
                              'redis_key': IDList_TYPE,
                              'state_scope': 'gateway'},
                             {'proto_file': 'orc8r.protos.service303_pb2',
                              'proto_msg': 'LogVerbosity',
                              'redis_key': LOG_TYPE,
                              'state_scope': 'gateway'}],
            'json_state': [{'redis_key': FOO_TYPE,
                              'state_scope': 'network'}]
        }
        service.loop = self.loop

        # Bind the rpc server to a free port
        self._rpc_server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=10)
        )
        port = self._rpc_server.add_insecure_port('0.0.0.0:0')
        # Add the servicer
        self._servicer = DummyStateServer()
        self._servicer.add_to_server(self._rpc_server)
        self._rpc_server.start()
        # Create a rpc stub
        self.channel = grpc.insecure_channel('0.0.0.0:{}'.format(port))

        serde1 = RedisSerde(NID_TYPE,
                            get_proto_serializer(),
                            get_proto_deserializer(NetworkID))
        serde2 = RedisSerde(IDList_TYPE,
                            get_proto_serializer(),
                            get_proto_deserializer(IDList))
        serde3 = RedisSerde(LOG_TYPE,
                            get_proto_serializer(),
                            get_proto_deserializer(LogVerbosity))
        serde4 = RedisSerde(FOO_TYPE,
                            get_json_serializer(),
                            get_json_deserializer())
        serde_map = {}
        serde_map[NID_TYPE] = serde1
        serde_map[IDList_TYPE] = serde2
        serde_map[LOG_TYPE] = serde3
        serde_map[FOO_TYPE] = serde4

        self.mock_client = RedisFlatDict(get_default_client(), serde_map)

        # Set up and start state replicating loop
        grpc_client_manager = GRPCClientManager(
            service_name="state",
            service_stub=StateServiceStub,
            max_client_reuse=60,
        )

        self.state_replicator = StateReplicator(
            service=service,
            grpc_client_manager=grpc_client_manager,
        )
        self.state_replicator.start()

    @mock.patch("redis.Redis", MockRedis)
    def tearDown(self):
        self._rpc_server.stop(None)
        self.state_replicator.stop()
        self.loop.close()

    def convert_msg_to_state(self, redis_state, is_proto=True):
        if is_proto:
            json_converted_state = MessageToDict(redis_state)
        else:
            json_converted_state = redis_state
        serialized_json_state = json.dumps(json_converted_state)
        return serialized_json_state.encode("utf-8")

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_collect_states_to_replicate(self):
        async def test():
            # Ensure setup is initialized properly
            self.mock_client.clear()

            nid_key = 'id1:' + NID_TYPE
            idlist_key = 'id1:' + IDList_TYPE
            foo_key = 'id1:' + FOO_TYPE
            foo = {'blah': 'bar', 'foo_id': 3}

            self.mock_client[nid_key] = NetworkID(id='foo')
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])
            self.mock_client[foo_key] = foo

            exp1 = self.convert_msg_to_state(self.mock_client[nid_key])
            exp2 = self.convert_msg_to_state(self.mock_client[idlist_key])
            exp3 = self.convert_msg_to_state(self.mock_client[foo_key], False)

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

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_replicate_states_success(self, get_rpc_mock):
        async def test():
            get_rpc_mock.return_value = self.channel

            # Ensure setup is initialized properly
            self.mock_client.clear()

            nid_key = 'id1:' + NID_TYPE
            idlist_key = 'id1:' + IDList_TYPE
            foo_key = 'id1:' + FOO_TYPE
            foo = {'blah': 'bar', 'foo_id': 3}
            self.mock_client[nid_key] = NetworkID(id='foo')
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])
            self.mock_client[foo_key] = foo
            # Increment version
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(3, len(req.states))

            # Ensure in-memory map updates properly
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(3, len(self.state_replicator._state_versions))
            mem_key1 = self.state_replicator.make_mem_key('id1', NID_TYPE)
            mem_key2 = self.state_replicator.make_mem_key('aaa-bbb:id1',
                                                          IDList_TYPE)
            mem_key3 = self.state_replicator.make_mem_key('id1', FOO_TYPE)
            self.assertEqual(1,
                             self.state_replicator._state_versions[mem_key1])
            self.assertEqual(2,
                             self.state_replicator._state_versions[mem_key2])
            self.assertEqual(1,
                             self.state_replicator._state_versions[mem_key3])

            # Now add new state and update some existing state
            nid_key2 = 'id2:' + NID_TYPE
            self.mock_client[nid_key2] = NetworkID(id='bar')
            self.mock_client[idlist_key] = IDList(ids=['bar', 'foo'])
            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(2, len(req.states))

            # Ensure in-memory map updates properly
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(4, len(self.state_replicator._state_versions))
            mem_key4 = self.state_replicator.make_mem_key('id2', NID_TYPE)
            self.assertEqual(1,
                             self.state_replicator._state_versions[mem_key1])
            self.assertEqual(3,
                             self.state_replicator._state_versions[mem_key2])
            self.assertEqual(1,
                             self.state_replicator._state_versions[mem_key3])
            self.assertEqual(1,
                             self.state_replicator._state_versions[mem_key4])

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_unreplicated_states(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel

            # Add initial state to be replicated
            self.mock_client.clear()

            nid_key = 'id1:' + NID_TYPE
            idlist_key = 'id1:' + IDList_TYPE
            log_key = 'id2:' + LOG_TYPE

            self.mock_client[nid_key] = NetworkID(id='foo')
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])
            # Increment version
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])
            # Set state that will be 'unreplicated'
            self.mock_client[log_key] = LogVerbosity(verbosity=5)

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(3, len(req.states))

            # Ensure in-memory map updates properly for successful replications
            await self.state_replicator._send_to_state_service(req)
            self.assertEqual(2, len(self.state_replicator._state_versions))
            mem_key1 = self.state_replicator.make_mem_key('id1', NID_TYPE)
            mem_key2 = self.state_replicator.make_mem_key('aaa-bbb:id1',
                                                          IDList_TYPE)
            self.assertEqual(1,
                             self.state_replicator._state_versions[mem_key1])
            self.assertEqual(2,
                             self.state_replicator._state_versions[mem_key2])

            # Now run again, ensuring only the state the wasn't replicated
            # will be sent again
            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(1, len(req.states))
            self.assertEqual('aaa-bbb:id2', req.states[0].deviceID)
            self.assertEqual(LOG_TYPE, req.states[0].type)

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_resync_success(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel
            self.mock_client.clear()

            nid_key = 'id1:' + NID_TYPE
            idlist_key = 'id1:' + IDList_TYPE

            # Set state that will be 'unsynced'
            self.mock_client[nid_key] = NetworkID(id='foo')
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])
            # Increment state's version
            self.mock_client[idlist_key] = IDList(ids=['bar', 'blah'])

            await self.state_replicator._resync()
            self.assertEqual(True, self.state_replicator._has_resync_completed)
            self.assertEqual(1, len(self.state_replicator._state_versions))
            mem_key = self.state_replicator.make_mem_key('aaa-bbb:id1',
                                                         IDList_TYPE)
            self.assertEqual(2, self.state_replicator._state_versions[mem_key])

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    @mock.patch('magma.magmad.state_reporter.ServiceRegistry.get_rpc_channel')
    def test_resync_failure(self, get_grpc_mock):
        async def test():
            get_grpc_mock.return_value = self.channel
            self.mock_client.clear()

            # Set state that will trigger the RpcError
            log_key = 'id1:' + LOG_TYPE
            self.mock_client[log_key] = LogVerbosity(verbosity=5)

            try:
                await self.state_replicator._resync()
            except grpc.RpcError:
                pass

            self.assertEqual(False,
                             self.state_replicator._has_resync_completed)
            self.assertEqual(0, len(self.state_replicator._state_versions))

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_invalid_key(self):
        async def test():
            self.mock_client.clear()

            bad_key = 'id1'
            # Force bad state format into mocked redis server
            self.mock_client.redis.set(bad_key, NetworkID(id='foo'))

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(None, req)

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())

    @mock.patch("redis.Redis", MockRedis)
    @mock.patch('snowflake.snowflake', get_mock_snowflake)
    def test_invalid_serde(self):
        async def test():
            self.mock_client.clear()

            bad_serde_key = 'id1:missing_serde'
            good_serde_key = 'id1:' + NID_TYPE
            # Force bad serde into mocked redis server
            self.mock_client.redis.set(bad_serde_key, NetworkID(id='foo'))
            self.mock_client[good_serde_key] = NetworkID(id='bar')

            req = await self.state_replicator._collect_states_to_replicate()
            self.assertEqual(1, len(req.states))
            self.assertEqual(NID_TYPE, req.states[0].type)

        # Cancel the replicator's loop so there are no other activities
        self.state_replicator._periodic_task.cancel()
        self.loop.run_until_complete(test())
