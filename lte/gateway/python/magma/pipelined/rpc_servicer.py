"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
from concurrent.futures import Future
from itertools import chain
from typing import List, Tuple

import grpc
from lte.protos import pipelined_pb2_grpc
from lte.protos.pipelined_pb2 import (
    ActivateFlowsResult,
    DeactivateFlowsResult,
    FlowResponse,
    RuleModResult,
    ActivateFlowsRequest)
from lte.protos.policydb_pb2 import PolicyRule
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.meter_stats import MeterStatsController
from magma.pipelined.metrics import (
    ENFORCEMENT_STATS_RULE_INSTALL_FAIL,
    ENFORCEMENT_RULE_INSTALL_FAIL,
)


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

        return self._activate_flows(request)

    def _activate_flows(self, request: ActivateFlowsRequest
                        ) -> ActivateFlowsResult:
        """
        Ensure that the RuleModResult is only successful if the flows are
        successfully added in both the enforcer app and enforcement_stats.
        Install enforcement_stats flows first because even if the enforcement
        flow install fails after, no traffic will be directed to the
        enforcement_stats flows.
        """
        for rule_id in request.rule_ids:
            self._service_manager.session_rule_version_mapper.update_version(
                request.sid.id, rule_id)
        for rule in request.dynamic_rules:
            self._service_manager.session_rule_version_mapper.update_version(
                request.sid.id, rule.id)
        enforcement_stats_res = self._activate_rules_in_enforcement_stats(
            request.sid.id, request.ip_addr, request.rule_ids,
            request.dynamic_rules)

        failed_static_rule_results, failed_dynamic_rule_results = \
            _retrieve_failed_results(enforcement_stats_res)
        # Do not install any rules that failed to install in enforcement_stats.
        static_rule_ids = \
            _filter_failed_static_rule_ids(request, failed_static_rule_results)
        dynamic_rules = \
            _filter_failed_dynamic_rules(request, failed_dynamic_rule_results)
        enforcement_res = self._activate_rules_in_enforcement(
            request.sid.id, request.ip_addr, static_rule_ids, dynamic_rules)

        # Include the failed rules from enforcement_stats in the response.
        enforcement_res.static_rule_results.extend(failed_static_rule_results)
        enforcement_res.dynamic_rule_results.extend(
            failed_dynamic_rule_results)
        return enforcement_res

    def _activate_rules_in_enforcement_stats(self, imsi: str, ip_addr: str,
                                             static_rule_ids: List[str],
                                             dynamic_rules: List[PolicyRule]
                                             ) -> ActivateFlowsResult:
        if not self._service_manager.is_app_enabled(
                EnforcementStatsController.APP_NAME):
            return ActivateFlowsResult()

        fut = Future()  # type: Future[ActivateFlowsResult]
        self._loop.call_soon_threadsafe(
            self._enforcement_stats.activate_rules,
            imsi, ip_addr, static_rule_ids, dynamic_rules, fut)
        enforcement_stats_res = fut.result()
        _report_enforcement_stats_failures(enforcement_stats_res, imsi)
        return enforcement_stats_res

    def _activate_rules_in_enforcement(self, imsi: str, ip_addr: str,
                                       static_rule_ids: List[str],
                                       dynamic_rules: List[PolicyRule]
                                       ) -> ActivateFlowsResult:
        fut = Future()  # type: Future[ActivateFlowsResult]
        self._loop.call_soon_threadsafe(
            self._enforcer_app.activate_rules,
            imsi, ip_addr, static_rule_ids, dynamic_rules, fut)
        enforcement_res = fut.result()
        _report_enforcement_failures(enforcement_res, imsi)
        return enforcement_res

    def DeactivateFlows(self, request, context):
        """
        Deactivate flows for a subscriber
        """
        if not self._service_manager.is_app_enabled(
                EnforcementController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        for rule_id in request.rule_ids:
            self._service_manager.session_rule_version_mapper.update_version(
                request.sid.id, rule_id)
        self._loop.call_soon_threadsafe(
            self._enforcer_app.deactivate_rules,
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


def _retrieve_failed_results(activate_flow_result: ActivateFlowsResult
                             ) -> Tuple[List[RuleModResult],
                                        List[RuleModResult]]:
    failed_static_rule_results = \
        [result for result in activate_flow_result.static_rule_results
         if result.result == RuleModResult.FAILURE]
    failed_dynamic_rule_results = \
        [result for result in
         activate_flow_result.dynamic_rule_results if
         result.result == RuleModResult.FAILURE]
    return failed_static_rule_results, failed_dynamic_rule_results


def _filter_failed_static_rule_ids(request: ActivateFlowsRequest,
                                   failed_results: List[RuleModResult]
                                   ) -> List[str]:
    failed_static_rule_ids = [result.rule_id for result in failed_results]
    return [rule_id for rule_id in request.rule_ids if
            rule_id not in failed_static_rule_ids]


def _filter_failed_dynamic_rules(request: ActivateFlowsRequest,
                                 failed_results: List[RuleModResult]
                                 ) -> List[PolicyRule]:
    failed_dynamic_rule_ids = [result.rule_id for result in failed_results]
    return [rule for rule in request.dynamic_rules if
            rule.id not in failed_dynamic_rule_ids]


def _report_enforcement_failures(activate_flow_result: ActivateFlowsResult,
                                 imsi: str):
    rule_results = chain(activate_flow_result.static_rule_results,
                         activate_flow_result.dynamic_rule_results)
    for result in rule_results:
        if result.result == RuleModResult.SUCCESS:
            continue
        ENFORCEMENT_RULE_INSTALL_FAIL.labels(rule_id=result.rule_id,
                                             imsi=imsi).inc()


def _report_enforcement_stats_failures(
        activate_flow_result: ActivateFlowsResult,
        imsi: str):
    rule_results = chain(activate_flow_result.static_rule_results,
                         activate_flow_result.dynamic_rule_results)
    for result in rule_results:
        if result.result == RuleModResult.SUCCESS:
            continue
        ENFORCEMENT_STATS_RULE_INSTALL_FAIL.labels(rule_id=result.rule_id,
                                                   imsi=imsi).inc()
