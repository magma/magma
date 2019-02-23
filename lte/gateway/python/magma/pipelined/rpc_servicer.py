"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from concurrent.futures import Future

import grpc
from lte.protos import pipelined_pb2_grpc
from lte.protos.pipelined_pb2 import DeactivateFlowsResult, FlowResponse
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.meter_stats import MeterStatsController


class PipelinedRpcServicer(pipelined_pb2_grpc.PipelinedServicer):
    """
    gRPC based server for Pipelined.
    """

    def __init__(self, loop, metering_stats, enforcer_app, enforcement_stats,
                 dpi_app, service_manager):
        self._loop = loop
        self._metering_stats = metering_stats
        self._enforcer_app = enforcer_app
        self._enforcement_stats = enforcement_stats
        self._dpi_app = dpi_app
        self._service_manager = service_manager

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        pipelined_pb2_grpc.add_PipelinedServicer_to_server(self, server)

    # --------------------------
    # Metering App
    # --------------------------

    def GetSubscriberMeteringFlows(self, request, context):
        """
        Returns all subscriber metering flows
        """
        if not self._service_manager.is_app_enabled(
                MeterStatsController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        fut = Future()
        self._loop.call_soon_threadsafe(
            self._metering_stats.get_subscriber_metering_flows, fut)
        return fut.result()

    # --------------------------
    # Enforcement App
    # --------------------------

    def ActivateFlows(self, request, context):
        """
        Activate flows for a subscriber based on the pre-defined rules
        """
        if not self._service_manager.is_app_enabled(
                EnforcementController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        fut = Future()
        self._loop.call_soon_threadsafe(
            self._enforcer_app.activate_flows,
            request.sid.id, request.ip_addr, request.rule_ids,
            request.dynamic_rules, fut)
        return fut.result()

    def DeactivateFlows(self, request, context):
        """
        Deactivate flows for a subscriber
        """
        if not self._service_manager.is_app_enabled(
                EnforcementController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        self._loop.call_soon_threadsafe(
            self._enforcer_app.deactivate_flows,
            request.sid.id, request.rule_ids)
        if self._service_manager.is_app_enabled(
                EnforcementStatsController.APP_NAME):
            self._loop.call_soon_threadsafe(
                self._enforcement_stats.delete_stats,
                request.sid.id, request.rule_ids)
        return DeactivateFlowsResult()

    # --------------------------
    # DPI App
    # --------------------------

    def CreateFlow(self, request, context):
        """
        Add dpi flow
        """
        if not self._service_manager.is_app_enabled(
                DPIController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        resp = FlowResponse()
        self._loop.call_soon_threadsafe(
            self._dpi_app.classify_flow,
            request.match, request.app_name)
        return resp

    def UpdateFlowStats(self, request, context):
        """
        Update stats for a flow
        """
        if not self._service_manager.is_app_enabled(
                DPIController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        resp = FlowResponse()
        return resp
