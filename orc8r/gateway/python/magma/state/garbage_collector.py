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
import logging

import grpc
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.redis.containers import RedisFlatDict
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.service import MagmaService
from magma.state.keys import make_scoped_device_id
from magma.state.redis_dicts import get_json_redis_dicts, get_proto_redis_dicts
from orc8r.protos.state_pb2 import DeleteStatesRequest, StateID

DEFAULT_GRPC_TIMEOUT = 10


class GarbageCollector:
    """
    GarbageCollector periodically fetches all state in Redis that is marked as
    garbage and deletes that state from the Orchestrator State service. If the
    RPC call succeeds, it then deletes the state from Redis
    """

    def __init__(
        self,
        service: MagmaService,
        grpc_client_manager: GRPCClientManager,
    ):
        self._service = service
        # Redis dicts for each type of state to replicate
        self._redis_dicts = []
        self._redis_dicts.extend(get_proto_redis_dicts(service.config))
        self._redis_dicts.extend(get_json_redis_dicts(service.config))

        # _grpc_client_manager to manage grpc client recyclings
        self._grpc_client_manager = grpc_client_manager

    async def run_garbage_collection(self):
        request = await self._collect_states_to_delete()
        if request is not None:
            await self._send_to_state_service(request)

    async def _collect_states_to_delete(self):
        states_to_delete = []
        for redis_dict in self._redis_dicts:
            for key in redis_dict.garbage_keys():
                state_scope = redis_dict.state_scope
                device_id = make_scoped_device_id(key, state_scope)
                sid = StateID(deviceID=device_id, type=redis_dict.redis_type)
                states_to_delete.append(sid)
        if len(states_to_delete) == 0:
            logging.debug("Not garbage collecting state. No state to delete!")
            return None
        # NetworkID will be filled in by Orchestrator from GW context
        return DeleteStatesRequest(networkID="", ids=states_to_delete)

    async def _send_to_state_service(self, request: DeleteStatesRequest):
        state_client = self._grpc_client_manager.get_client()
        try:
            await grpc_async_wrapper(
                state_client.DeleteStates.future(
                    request,
                    DEFAULT_GRPC_TIMEOUT,
                ),
            )

        except grpc.RpcError as err:
            logging.error("GRPC call failed for state deletion: %s", err)
        else:
            for redis_dict in self._redis_dicts:
                for key in redis_dict.garbage_keys():
                    await self._delete_state_from_redis(redis_dict, key)

    async def _delete_state_from_redis(
        self,
        redis_dict: RedisFlatDict,
        key: str,
    ) -> None:
        # Ensure that the object isn't updated before deletion
        with redis_dict.lock(key):
            deleted = redis_dict.delete_garbage(key)
            if deleted:
                logging.debug(
                    "Successfully garbage collected "
                    "state for key: %s", key,
                )
            else:
                logging.debug(
                    "Successfully garbage collected "
                    "state in cloud for key %s. "
                    "Didn't delete locally as the "
                    "object is no longer garbage", key,
                )
