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

import asyncio
import logging
from typing import Any, Dict, List, Optional

import grpc
import snowflake
from magma.common.cert_validity import cert_is_invalid
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service import get_service303_client
from magma.common.service_registry import ServiceRegistry
from magma.magmad.bootstrap_manager import BootstrapManager
from magma.magmad.gateway_status import GatewayStatusFactory
from magma.magmad.metrics import CHECKIN_STATUS
from magma.magmad.service_poller import ServiceInfo
from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2 import State
from orc8r.protos.state_pb2 import ReportStatesRequest

States = List[State]


class StateReporterErrorHandler:
    """
    StateReporterErrorHandler holds metadata related to state reporting
    failures. It also defines report_to_cloud_error that has the logic to
    trigger a bootstrap if it finds the certificate is bad or there are
    permission issues.
    """

    def __init__(
        self,
        loop: asyncio.AbstractEventLoop,
        config: Any,
        grpc_client_manager: GRPCClientManager,
        bootstrap_manager: BootstrapManager,
    ):
        self._loop = loop
        # Number of consecutive failed state reporting before we check for an
        # outdated cert
        self.fail_threshold = 1
        # Current number of consecutive failed to report states to cloud
        self.num_failed_state_reporting = 0
        # Current number of consecutive gateway states skipped from sending due
        # to missing fields
        self.num_skipped_gateway_states = 0
        # max_skipped_gw_states specifies number of gateway states that can
        # have an empty/missing service meta before reporting it anyway.
        self.max_skipped_gw_states = int(config.get("max_skipped_checkins", 3))

        # _grpc_client_manager to manage grpc client recyclings
        self._grpc_client_manager = grpc_client_manager
        # bootstrap manager to trigger a bootstrap when state reporting fails
        # because of a bad cert
        self._bootstrap_manager = bootstrap_manager

    def report_to_cloud_error(self, err):
        """
        report_to_cloud_error checks if the number of failed reporting exceeds
        the threshold specified in the config. If it does, it will trigger a
        bootstrap if the certificate is invalid.
        """
        logging.error(
            "Checkin Error! Failed to report states. [%s] %s",
            err.code(), err.details(),
        )
        CHECKIN_STATUS.set(0)
        self.num_failed_state_reporting += 1
        if self.num_failed_state_reporting >= self.fail_threshold:
            logging.info(
                'StateReporting (Checkin) failure threshold met, '
                'remediating...',
            )
            asyncio.ensure_future(
                self._schedule_bootstrap_if_cert_is_invalid(err.code()),
            )
        self._grpc_client_manager.on_grpc_fail(err.code())

    async def _schedule_bootstrap_if_cert_is_invalid(self, err_code):
        """
        Checks for invalid certificate as cause for state reporting failures
        """
        if err_code == grpc.StatusCode.PERMISSION_DENIED:
            await self._bootstrap_manager.schedule_bootstrap_now()
            return
        proxy_config = ServiceRegistry.get_proxy_config()
        host = proxy_config['cloud_address']
        port = proxy_config['cloud_port']
        cert_file = proxy_config['gateway_cert']
        key_file = proxy_config['gateway_key']
        not_valid = await \
            cert_is_invalid(host, port, cert_file, key_file, self._loop)
        if not_valid:
            logging.error('Bootstrapping due to invalid cert')
            await self._bootstrap_manager.schedule_bootstrap_now()
            return
        else:
            logging.error(
                'StateReporting failure likely '
                'not due to invalid cert',
            )


class StateReporter(SDWatchdogTask):
    """
    Periodically collects operational states from service303 states and reports
    them to the cloud state service.
    If it fails to send states to the cloud, it will call the failure_cb
    function.
    In this context, check-in refers to the act of connecting and reporting
    states to the cloud.
    """

    def __init__(
        self, config: Any, mconfig: Any,
        loop: asyncio.AbstractEventLoop,
        bootstrap_manager: BootstrapManager,
        gw_status_factory: GatewayStatusFactory,
        grpc_client_manager: GRPCClientManager,
    ):
        super().__init__(
            interval=max(mconfig.checkin_interval, 5),
            loop=loop,
        )
        self._loop = loop
        # keep a pointer to mconfig since config stored can change over time
        self._mconfig = mconfig

        # Manages all metadata and methods on dealing with failures.
        # (invalid gateway status, cloud reporting error)
        self._error_handler = StateReporterErrorHandler(
            loop=loop,
            config=config,
            grpc_client_manager=grpc_client_manager,
            bootstrap_manager=bootstrap_manager,
        )

        # gateway status factory to bundle various information about this
        # gateway into an object.
        self._gw_status_factory = gw_status_factory

        # grpc_client_manager to manage grpc client recycling
        self._grpc_client_manager = grpc_client_manager

        # A dictionary of all services registered with a service303 interface.
        # Holds service name to service info gathered from the config
        self._service_info_by_name = self._construct_service_info_by_name(
            config=config,
        )

        # Initially set status to 1, otherwise on the first round we report a
        # failure. This is particularly an issue if magmad restarts frequenty.
        CHECKIN_STATUS.set(1)
        # set initial timeout to "large" since no reporting can occur until
        # bootstrap succeeds.
        self.set_timeout(60 * 60 * 2)
        # initially set task as alive to wait for bootstrap
        self.heartbeat()

    async def _run(self) -> None:
        request = await self._collect_states()
        if request is not None:
            await self._send_to_state_service(request)
        # reset checkin_interval in case it has changed
        self.set_interval(max(self._mconfig.checkin_interval, 5))

    async def _collect_states(self) -> Optional[ReportStatesRequest]:
        states = []
        for service in self._service_info_by_name:
            result = await self._get_operational_states(service=service)
            states.extend(result)

        gw_state = self._get_gw_state()
        if gw_state is not None:
            states.append(gw_state)

        if len(states) > 0:
            # ReportStates returns error on empty request
            return ReportStatesRequest(states=states)
        return None

    async def _send_to_state_service(
            self,
            request: ReportStatesRequest,
    ) -> None:
        state_client = self._grpc_client_manager.get_client()
        try:
            response = await grpc_async_wrapper(
                state_client.ReportStates.future(
                    request,
                    self._mconfig.checkin_timeout,
                ),
                self._loop,
            )
            for idAndError in response.unreportedStates:
                logging.error(
                    "Failed to report state for (%s,%s): %s",
                    idAndError.type, idAndError.deviceID, idAndError.error,
                )
            # Report that the gateway successfully connected to the cloud
            CHECKIN_STATUS.set(1)
            self._error_handler.num_failed_state_reporting = 0
            logging.info(
                "Checkin Successful! "
                "Successfully sent states to the cloud!",
            )
        except grpc.RpcError as err:
            self._error_handler.report_to_cloud_error(err)
        finally:
            # reset timeout to config-specified + some buffer
            self.set_timeout(self._interval * 2)

    async def _get_operational_states(self, service: str) -> List[States]:
        client = get_service303_client(service, ServiceRegistry.LOCAL)
        if client is None:
            return []
        try:
            states = []
            future = client.GetOperationalStates.future(
                Void(),
                self._mconfig.checkin_timeout,
            )
            result = await grpc_async_wrapper(future, self._loop)
            for i in range(len(result.states)):
                states.append(result.states[i])
            return states
        except Exception as err:
            logging.error(
                "GetOperationalStates Error for %s! [%s] %s",
                service, err.code(), err.details(),
            )
            return []

    def _get_gw_state(self) -> Optional[State]:
        gw_type = "gw_state"
        gw_state, has_all_required_fields = \
            self._gw_status_factory.get_serialized_status()
        if has_all_required_fields:
            self._error_handler.num_skipped_gateway_states = 0
            return self._make_state(gw_type, snowflake.snowflake(), gw_state)

        # check if we have failed to send states too many times in a row
        if 0 < self._error_handler.max_skipped_gw_states < \
                self._error_handler.num_skipped_gateway_states:
            logging.warning(
                "Number of skipped checkins exceeds %d "
                "(cfg: max_skipped_checkins). Checking in anyway.",
                self._error_handler.max_skipped_gw_states,
            )
            # intentionally don't reset num_skipped_gateway_states here
            return self._make_state(gw_type, snowflake.snowflake(), gw_state)

        # skipping reporting gateway state
        self._error_handler.num_skipped_gateway_states += 1
        return None

    @staticmethod
    def _make_state(type_val: str, key: str, value: str) -> State:
        return State(type=type_val, deviceID=key, value=value.encode('utf-8'))

    @staticmethod
    def _construct_service_info_by_name(config) -> Dict[str, ServiceInfo]:
        info = {}
        for service in config['magma_services']:
            # Check whether service provides service303 interface
            if service not in config['non_service303_services']:
                info[service] = ServiceInfo(service)
        return info
