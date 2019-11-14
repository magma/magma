"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=broad-except

import logging
import importlib
import json
import jsonpickle
import snowflake
import grpc

from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict
from magma.common.redis.serializers import get_proto_deserializer, \
    get_proto_serializer, get_json_deserializer, get_json_serializer, \
    RedisSerde
from magma.common.service import MagmaService
from magma.common.sdwatchdog import SDWatchdogTask
from orc8r.protos.state_pb2 import ReportStatesRequest, SyncStatesRequest, \
    IDAndVersion, StateID
from orc8r.protos.service303_pb2 import State
from magma.common.rpc_utils import grpc_async_wrapper
from google.protobuf.json_format import MessageToDict

# TODO: Make DEFAULT_SYNC_INTERVAL an mconfig parameter
DEFAULT_SYNC_INTERVAL = 60
DEFAULT_GRPC_TIMEOUT = 10
MINIMUM_SYNC_INTERVAL = 30
PROTO_FORMAT = 0
JSON_FORMAT = 1


class StateDict(RedisFlatDict):
    """
    StateDict is a RedisFlatDict that holds state metadata and reads/writes
    state to Redis.
    """
    def __init__(self, serde: RedisSerde, state_scope: str, state_format: int):
        super().__init__(get_default_client(), serde)
        # Scope determines the deviceID to report the state with
        self.state_scope = state_scope
        self.state_format = state_format


class StateReplicator(SDWatchdogTask):
    """
    StateReplicator periodically fetches all configured state from Redis,
    reporting any updates to the Orchestrator State service.
    """
    def __init__(self,
                 service: MagmaService,
                 grpc_client_manager: GRPCClientManager):
        super().__init__(DEFAULT_SYNC_INTERVAL, service.loop)
        self._service = service
        # In memory mapping of states to version
        self._state_versions = {}
        # Redis clients for each type of state to replicate
        self._redis_clients = []
        self._redis_clients.extend(self._get_proto_redis_clients())
        self._redis_clients.extend(self._get_json_redis_clients())
        # _grpc_client_manager to manage grpc client recyclings
        self._grpc_client_manager = grpc_client_manager

        # Flag to indicate if resync has completed successfully.
        # Replication cannot proceed until this flag is True
        self._has_resync_completed = False

    def _get_proto_redis_clients(self):
        clients = []
        state_protos = self._service.config.get('state_protos', []) or []
        for proto_cfg in state_protos:
            is_invalid_cfg = 'proto_msg' not in proto_cfg or \
                             'proto_file' not in proto_cfg or \
                             'redis_key' not in proto_cfg or \
                             'state_scope' not in proto_cfg
            if is_invalid_cfg:
                logging.warning("Invalid proto config found in state_protos "
                                "configuration: %s", proto_cfg)
                continue
            try:
                proto_module = importlib.import_module(proto_cfg['proto_file'])
                msg = getattr(proto_module, proto_cfg['proto_msg'])
                redis_key = proto_cfg['redis_key']
                logging.info('Initializing RedisSerde for proto state %s',
                             proto_cfg['redis_key'])
                serde = RedisSerde(redis_key,
                                   get_proto_serializer(),
                                   get_proto_deserializer(msg))
                client = StateDict(serde,
                                   proto_cfg['state_scope'],
                                   PROTO_FORMAT)
                clients.append(client)

            except (ImportError, AttributeError) as err:
                logging.error(err)

        return clients

    def _get_json_redis_clients(self):
        clients = []
        json_state = self._service.config.get('json_state', []) or []
        for json_cfg in json_state:
            is_invalid_cfg = 'redis_key' not in json_cfg or \
                             'state_scope' not in json_cfg
            if is_invalid_cfg:
                logging.warning("Invalid json state config found in json_state"
                                "configuration: %s", json_cfg)
                continue

            logging.info('Initializing RedisSerde for json state %s',
                         json_cfg['redis_key'])
            redis_key = json_cfg['redis_key']
            serde = RedisSerde(redis_key,
                               get_json_serializer(),
                               get_json_deserializer())
            client = StateDict(serde,
                           json_cfg['state_scope'],
                           JSON_FORMAT)
            clients.append(client)

        return clients

    async def _run(self):
        if not self._has_resync_completed:
            try:
                await self._resync()
            except grpc.RpcError as err:
                logging.error("GRPC call failed for initial state re-sync: %s",
                              err)
                return
        request = await self._collect_states_to_replicate()
        if request is not None:
            await self._send_to_state_service(request)

    async def _resync(self):
        states_to_sync = []
        for client in self._redis_clients:
            for key in client:
                version = client.get_version(key)
                device_id = self.make_scoped_device_id(key, client.state_scope)
                state_id = StateID(type=client.redis_type, deviceID=device_id)
                id_and_version = IDAndVersion(id=state_id, version=version)
                states_to_sync.append(id_and_version)

        if len(states_to_sync) == 0:
            logging.debug("Not re-syncing state. No local state found.")
            return
        state_client = self._grpc_client_manager.get_client()
        request = SyncStatesRequest(states=states_to_sync)
        response = await grpc_async_wrapper(
            state_client.SyncStates.future(
                request,
                DEFAULT_GRPC_TIMEOUT,
            ),
            self._loop)
        unsynced_states = set()
        for id_and_version in response.unsyncedStates:
            unsynced_states.add((id_and_version.id.type,
                                 id_and_version.id.deviceID))
        # Update in-memory map to add already synced states
        for state in request.states:
            in_mem_key = self.make_mem_key(state.id.deviceID, state.id.type)
            if (state.id.type, state.id.deviceID) not in unsynced_states:
                self._state_versions[in_mem_key] = state.version

        self._has_resync_completed = True
        logging.info("Successfully resynced state with Orchestrator!")

    async def _collect_states_to_replicate(self):
        states_to_report = []
        for client in self._redis_clients:
            for key in client:
                device_id = self.make_scoped_device_id(key, client.state_scope)
                in_mem_key = self.make_mem_key(device_id, client.redis_type)
                redis_version = client.get_version(key)

                if in_mem_key in self._state_versions and \
                        self._state_versions[in_mem_key] == redis_version:
                    continue

                redis_state = client.get(key)
                if client.state_format == PROTO_FORMAT:
                    state_to_serialize = MessageToDict(redis_state)
                    serialized_json_state = json.dumps(state_to_serialize)
                else:
                    serialized_json_state = jsonpickle.encode(redis_state)
                state_proto = State(type=client.redis_type,
                      deviceID=device_id,
                      value=serialized_json_state.encode("utf-8"),
                      version=redis_version)

                states_to_report.append(state_proto)

        if len(states_to_report) == 0:
            logging.debug("Not replicating state. No state has changed!")
            return None
        return ReportStatesRequest(states=states_to_report)

    async def _send_to_state_service(self, request: ReportStatesRequest):
        state_client = self._grpc_client_manager.get_client()
        try:
            response = await grpc_async_wrapper(
                state_client.ReportStates.future(
                    request,
                    DEFAULT_GRPC_TIMEOUT,
                ),
                self._loop)

        except grpc.RpcError as err:
            logging.error("GRPC call failed for state replication: %s", err)
        else:
            unreplicated_states = set()
            for idAndError in response.unreportedStates:
                logging.warning(
                    "Failed to replicate state for (%s,%s): %s",
                    idAndError.type, idAndError.deviceID, idAndError.error)
                unreplicated_states.add((idAndError.type, idAndError.deviceID))
            # Update in-memory map for successfully reported states
            for state in request.states:
                if (state.type, state.deviceID) in unreplicated_states:
                    continue
                in_mem_key = self.make_mem_key(state.deviceID, state.type)
                self._state_versions[in_mem_key] = state.version

                logging.debug("Successfully replicated state for: "
                              "deviceID: %s,"
                              "type: %s, "
                              "version: %d",
                              state.deviceID, state.type, state.version)
        finally:
            # reset timeout to config-specified + some buffer
            self.set_timeout(self._interval * 2)

    @staticmethod
    def make_mem_key(device_id, state_type):
        """
        Create a key of the format <id>:<type>
        """
        return device_id + ":" + state_type

    @staticmethod
    def make_scoped_device_id(idval, scope):
        """
        Create a deviceID of the format <id> for scope 'network'
        Otherwise create a key of the format <hwid>:<id> for 'gateway' or
        unrecognized scope.
        """
        if scope == "network":
            return idval
        else:
            return snowflake.snowflake() + ":" + idval
