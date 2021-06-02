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
# pylint: disable=broad-except

import json
import logging

import grpc
import jsonpickle
from google.protobuf.json_format import MessageToDict
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service import MagmaService
from magma.state.garbage_collector import GarbageCollector
from magma.state.keys import make_mem_key, make_scoped_device_id
from magma.state.redis_dicts import (
    PROTO_FORMAT,
    get_json_redis_dicts,
    get_proto_redis_dicts,
)
from orc8r.protos.service303_pb2 import State
from orc8r.protos.state_pb2 import (
    IDAndVersion,
    ReportStatesRequest,
    StateID,
    SyncStatesRequest,
)

# TODO: Make DEFAULT_SYNC_INTERVAL an mconfig parameter
DEFAULT_SYNC_INTERVAL = 60
DEFAULT_GRPC_TIMEOUT = 10
GARBAGE_COLLECTION_ITERATION_INTERVAL = 2


class StateReplicator(SDWatchdogTask):
    """
    StateReplicator periodically fetches all configured state from Redis,
    reporting any updates to the Orchestrator State service.
    """

    def __init__(
        self,
        service: MagmaService,
        garbage_collector: GarbageCollector,
        grpc_client_manager: GRPCClientManager,
    ):
        sync_interval = service.config.get(
            'sync_interval', DEFAULT_SYNC_INTERVAL,
        )
        super().__init__(sync_interval, service.loop)
        self._service = service
        # Garbage collector to propagate deletions back to Orchestrator
        self._garbage_collector = garbage_collector
        # In memory mapping of states to version
        self._state_versions = {}
        # Set of keys from current replication iteration - used to track
        # keys to delete from _state_versions dict
        self._state_keys_from_current_iteration = set()
        # Redis clients for each type of state to replicate
        self._redis_dicts = []
        self._redis_dicts.extend(get_proto_redis_dicts(service.config))
        self._redis_dicts.extend(get_json_redis_dicts(service.config))
        # _grpc_client_manager to manage grpc client recyclings
        self._grpc_client_manager = grpc_client_manager

        # Flag to indicate if resync has completed successfully.
        # Replication cannot proceed until this flag is True
        self._has_resync_completed = False

        # Track replication iteration to track when to trigger garbage
        # collection
        self._replication_iteration = 0

    async def _run(self):
        logging.debug("Check state")
        if not self._has_resync_completed:
            try:
                await self._resync()
            except grpc.RpcError as err:
                logging.error(
                    "GRPC call failed for initial state re-sync: %s",
                    err,
                )
                return
        request = await self._collect_states_to_replicate()
        if request is not None:
            await self._send_to_state_service(request)
        await self._cleanup_deleted_keys()

        self._replication_iteration += 1
        if self._replication_iteration >= \
                GARBAGE_COLLECTION_ITERATION_INTERVAL:
            await self._garbage_collector.run_garbage_collection()
            self._replication_iteration = 0
        logging.debug("")

    async def _resync(self):
        states_to_sync = []
        for redis_dict in self._redis_dicts:
            for key in redis_dict:
                version = redis_dict.get_version(key)
                device_id = make_scoped_device_id(key, redis_dict.state_scope)
                state_id = StateID(
                    type=redis_dict.redis_type,
                    deviceID=device_id,
                )
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
            self._loop,
        )
        unsynced_states = set()
        for id_and_version in response.unsyncedStates:
            unsynced_states.add((
                id_and_version.id.type,
                id_and_version.id.deviceID,
            ))
        # Update in-memory map to add already synced states
        for state in request.states:
            in_mem_key = make_mem_key(state.id.deviceID, state.id.type)
            if (state.id.type, state.id.deviceID) not in unsynced_states:
                self._state_versions[in_mem_key] = state.version

        self._has_resync_completed = True
        logging.info("Successfully resynced state with Orchestrator!")

    async def _collect_states_to_replicate(self):
        states_to_report = []
        for redis_dict in self._redis_dicts:
            for key in redis_dict:
                redis_state = redis_dict.get(key)
                device_id = make_scoped_device_id(key, redis_dict.state_scope)

                in_mem_key = make_mem_key(device_id, redis_dict.redis_type)
                if redis_state is None:
                    logging.debug(
                        "Content of key %s is empty, skipping", in_mem_key,
                    )
                    continue

                redis_version = redis_dict.get_version(key)
                self._state_keys_from_current_iteration.add(in_mem_key)
                if in_mem_key in self._state_versions and \
                        self._state_versions[in_mem_key] == redis_version:
                    logging.debug(
                        "key %s already read on this iteration, skipping", in_mem_key,
                    )
                    continue

                try:
                    if redis_dict.state_format == PROTO_FORMAT:
                        state_to_serialize = MessageToDict(redis_state)
                        serialized_json_state = json.dumps(state_to_serialize)
                    else:
                        serialized_json_state = jsonpickle.encode(redis_state)
                except Exception as e:  # pylint: disable=broad-except
                    logging.error(
                        "Found bad state for %s for %s, not "
                        "replicating this state: %s",
                        key, device_id, e,
                    )
                    continue

                state_proto = State(
                    type=redis_dict.redis_type,
                    deviceID=device_id,
                    value=serialized_json_state.encode(
                        "utf-8",
                    ),
                    version=redis_version,
                )

                logging.debug(
                    "key with version, %s contains: %s", in_mem_key,
                    serialized_json_state,
                )
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
                self._loop,
            )

        except grpc.RpcError as err:
            logging.error("GRPC call failed for state replication: %s", err)
        else:
            unreplicated_states = set()
            for idAndError in response.unreportedStates:
                logging.warning(
                    "Failed to replicate state for (%s,%s): %s",
                    idAndError.type, idAndError.deviceID, idAndError.error,
                )
                unreplicated_states.add((idAndError.type, idAndError.deviceID))
            # Update in-memory map for successfully reported states
            for state in request.states:
                if (state.type, state.deviceID) in unreplicated_states:
                    continue
                in_mem_key = make_mem_key(state.deviceID, state.type)
                self._state_versions[in_mem_key] = state.version

                logging.debug(
                    "Successfully replicated state for: "
                    "deviceID: %s,"
                    "type: %s, "
                    "version: %d",
                    state.deviceID, state.type, state.version,
                )
        finally:
            # reset timeout to config-specified + some buffer
            self.set_timeout(self._interval * 2)

    async def _cleanup_deleted_keys(self):
        deleted_keys = set(self._state_versions) - \
            self._state_keys_from_current_iteration
        for key in deleted_keys:
            del self._state_versions[key]
        self._state_keys_from_current_iteration = set()
