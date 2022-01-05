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
import time
from typing import List

import grpc
from magma.common.job import Job
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.service_registry import ServiceRegistry
from magma.magmad.metrics import UNEXPECTED_SERVICE_RESTARTS
from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2_grpc import Service303Stub


class ServiceInfo(object):
    """
    Stores info about the individual services
    """
    # Time buffer for services to restart, in seconds
    SERVICE_RESTART_BUFFER_TIME = 30

    def __init__(self, service_name):
        self.continuous_timeouts = 0
        self._service_name = service_name
        self._expected_start_time = time.time()
        self._status = None
        self._linked_services = []
        # Initialize the counter for each service
        UNEXPECTED_SERVICE_RESTARTS.labels(
            service_name=self._service_name,
        ).inc(0)

    @property
    def status(self):
        return self._status

    @property
    def linked_services(self):
        return self._linked_services

    def add_linked_services(self, service_list):
        for service in service_list:
            if service != self._service_name and \
                    service not in self._linked_services:
                self._linked_services.append(service)

    def update(self, start_time, status):
        self._status = status
        self.continuous_timeouts = 0
        if start_time <= self._expected_start_time:
            # Probably a race in service starts, or magmad restarted
            return
        if (
            start_time - self._expected_start_time
            > self.SERVICE_RESTART_BUFFER_TIME
        ):
            UNEXPECTED_SERVICE_RESTARTS.labels(
                service_name=self._service_name,
            ).inc()
            self._expected_start_time = start_time

    def process_service_restart(self):
        self._expected_start_time = time.time()


class ServicePoller(Job):
    """
    Periodically query the services' Service303 interface
    """
    # Periodicity for getting status from other services, in seconds
    GET_STATUS_INTERVAL = 10
    # Timeout when getting status from other local services, in seconds
    GET_STATUS_TIMEOUT = 8

    def __init__(self, loop, config, dynamic_services: List[str] = None):
        """
        Initialize the ServicePooler

        Args:
            loop: loop
            config: configuration
            dynamic_services: list of dynamic services
        """
        super().__init__(
            interval=self.GET_STATUS_INTERVAL,
            loop=loop,
        )
        self._config = config
        # Holds a map of service name -> ServiceInfo
        self._service_info = {}
        for service in config['magma_services']:
            self._service_info[service] = ServiceInfo(service)
        if dynamic_services is not None:
            for service in dynamic_services:
                self._service_info[service] = ServiceInfo(service)
        for service_list in config.get('linked_services', []):
            for service in service_list:
                self._service_info[service].add_linked_services(service_list)

    def update_dynamic_services(
        self,
        new_services: List[str],
        stopped_services: List[str],
    ):
        """
        Update the service poller when dynamic services are enabled or disabled

        Args:
            new_services: New services which were enabled
            stopped_services: Old services which were disabled
        """
        for service in new_services:
            self._service_info[service] = ServiceInfo(service)
        for service in stopped_services:
            self._service_info.pop(service)

    def get_service_timeouts(self):
        ret = {}
        for service_name, service in self._service_info.items():
            ret[service_name] = service.continuous_timeouts
        return ret

    def reset_timeout_counter(self, service):
        self._service_info[service].continuous_timeouts = 0

    @property
    def service_info(self):
        return self._service_info

    def process_service_restart(self, service_name):
        self._service_info[service_name].process_service_restart()
        for linked_service in self._service_info[service_name].linked_services:
            self._service_info[linked_service].process_service_restart()

    async def _run(self):
        await self._get_service_info()

    async def _get_service_info(self):
        """
        Make RPC calls to 'GetServiceInfo' functions of other services, to
        get current status.
        """
        for service in list(self._service_info):
            # Check whether service provides service303 interface
            if service in self._config['non_service303_services']:
                continue
            try:
                chan = ServiceRegistry.get_rpc_channel(
                    service, ServiceRegistry.LOCAL,
                )
            except ValueError:
                # Service can't be contacted
                logging.error('Cant get RPC channel to %s', service)
                continue
            client = Service303Stub(chan)
            try:
                future = client.GetServiceInfo.future(
                    Void(),
                    self.GET_STATUS_TIMEOUT,
                )
                info = await grpc_async_wrapper(future, self._loop)
                self._service_info[service].update(
                    info.start_time_secs,
                    info.status,
                )
                self._service_info[service].continuous_timeouts = 0
            except grpc.RpcError as err:
                logging.error(
                    "GetServiceInfo Error for %s! [%s] %s",
                    service,
                    err.code(),
                    err.details(),
                )
                self._service_info[service].continuous_timeouts += 1
