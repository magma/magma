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

from magma.common.redis.client import get_default_client
from magma.common.redis.containers import RedisFlatDict
from magma.common.redis.serializers import (
    RedisSerde,
    get_proto_deserializer,
    get_proto_serializer,
)
from orc8r.protos.service_status_pb2 import ServiceExitStatus


class ServiceStateWrapper:
    """
    Class wraps ServiceState interactions with redis
    """

    # Unique typename for Redis key
    REDIS_VALUE_TYPE = "systemd_status"

    def __init__(self):
        serde = RedisSerde(
            self.REDIS_VALUE_TYPE,
            get_proto_serializer(),
            get_proto_deserializer(ServiceExitStatus),
        )
        self._flat_dict = RedisFlatDict(get_default_client(), serde)

    def update_service_status(
        self, service_name: str,
        service_status: ServiceExitStatus,
    ) -> None:
        """
        Update the service exit status for a given service
        """

        if service_name in self._flat_dict:
            current_service_status = self._flat_dict[service_name]
        else:
            current_service_status = ServiceExitStatus()

        if service_status.latest_service_result == \
                ServiceExitStatus.ServiceResult.Value("SUCCESS"):
            service_status.num_clean_exits = \
                current_service_status.num_clean_exits + 1
            service_status.num_fail_exits = \
                current_service_status.num_fail_exits
        else:
            service_status.num_fail_exits = \
                current_service_status.num_fail_exits + 1
            service_status.num_clean_exits = \
                current_service_status.num_clean_exits
        self._flat_dict[service_name] = service_status

    def get_service_status(self, service_name: str) -> ServiceExitStatus:
        """
        Get the service status protobuf for a given service
        @returns ServiceStatus protobuf object
        """
        return self._flat_dict[service_name]

    def get_all_services_status(self) -> [str, ServiceExitStatus]:
        """
        Get a dict of service name to service status
        @return dict of service_name to service map
        """
        service_status = {}
        for k, v in self._flat_dict.items():
            service_status[k] = v
        return service_status

    def cleanup_service_status(self) -> None:
        """
        Cleanup service status for all services in redis, mostly using for
        testing
        """
        self._flat_dict.clear()
