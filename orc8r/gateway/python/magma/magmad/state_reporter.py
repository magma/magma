"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=broad-except

import logging
from typing import Optional, List, Dict
import snowflake

from orc8r.protos.common_pb2 import Void
from orc8r.protos.service303_pb2_grpc import Service303Stub
from orc8r.protos.service303_pb2 import State
from orc8r.protos.state_pb2 import ReportStatesRequest
from orc8r.protos.state_pb2_grpc import StateServiceStub
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service_registry import ServiceRegistry
from magma.magmad.service_poller import ServiceInfo
from magma.common.service import MagmaService
from magma.magmad.checkin_manager import CheckinRequest

States = List[State]


class StateReporter(SDWatchdogTask):
    """
    Periodically collects operational states from service303 states and reports
    them to the cloud state service.
    """
    def __init__(self, service: MagmaService, checkin_manager: CheckinRequest):
        super().__init__(
            max(5, service.mconfig.checkin_interval),
            service.loop
        )

        self._service = service
        self._config = service.config
        self._checkin_manager = checkin_manager
        # Holds a map of service name -> ServiceInfo
        self._service_info_by_name = self._construct_service_info()

    async def _run(self) -> None:
        request = await self._collect_states()
        if request is not None:
            await self._send_to_state_service(request)

    async def _get_state(self, service: MagmaService) -> Optional[States]:
        client = self._get_client(service)
        if client is None:
            return None
        try:
            states = []
            future = client.GetOperationalStates.future(
                Void(),
                self._service.mconfig.checkin_timeout,
            )
            result = await grpc_async_wrapper(future, self._loop)
            for i in range(len(result.states)):
                states.append(result.states[i])
            return states
        except Exception as err:
            logging.error("GetOperationalStates Error for %s! [%s]",
                          service, err)
            return None

    async def _collect_states(self) -> Optional[ReportStatesRequest]:
        states = []
        for service in self._service_info_by_name:
            result = await self._get_state(service=service)
            if result is not None:
                states.extend(result)
        gw_state = self._get_gw_state()
        if gw_state is not None:
            states.append(gw_state)
        if len(states) == 0:
            return None
        return ReportStatesRequest(
            states=states,
        )

    async def _send_to_state_service(
            self,
            request: ReportStatesRequest)-> None:
        chan = ServiceRegistry.get_rpc_channel(
            'state',
            ServiceRegistry.CLOUD
        )
        state_client = StateServiceStub(chan)
        try:
            await grpc_async_wrapper(
                state_client.ReportStates.future(
                    request,
                    self._service.mconfig.checkin_timeout,
                ),
                self._loop)
        except Exception as err:
            logging.error("Failed to make a ReportStates request: %s", err)

    def _get_gw_state(self) -> Optional[State]:
        gw_state = self._checkin_manager.get_latest_gw_state()
        if gw_state is not None:
            state = State(type="gw_state",
                          deviceID=snowflake.snowflake(),
                          value=gw_state.encode('utf-8'))
            return state
        return None

    def _construct_service_info(self) -> Dict[str, ServiceInfo]:
        info = {}
        for service in self._config['magma_services']:
            # Check whether service provides service303 interface
            if service not in self._config['non_service303_services']:
                info[service] = ServiceInfo(service)
        return info

    @staticmethod
    def _get_client(service: MagmaService) -> Optional[Service303Stub]:
        try:
            chan = ServiceRegistry.get_rpc_channel(
                service, ServiceRegistry.LOCAL)
            return Service303Stub(chan)
        except ValueError:
            # Service can't be contacted
            logging.error('Failed to get RPC channel to %s', service)
            return None
