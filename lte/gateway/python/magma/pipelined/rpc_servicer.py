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
import concurrent.futures
import logging
import os
import queue
from collections import OrderedDict
from concurrent.futures import Future
from typing import List, Tuple

import grpc
from lte.protos import pipelined_pb2_grpc
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import (ActivateFlowsRequest,
                                      ActivateFlowsResult, AllTableAssignments,
                                      CauseIE, DeactivateFlowsRequest,
                                      DeactivateFlowsResult, FlowResponse,
                                      OffendingIE, RequestOriginType,
                                      RuleModResult, SessionSet,
                                      SetupFlowsResult, SetupPolicyRequest,
                                      SetupQuotaRequest, SetupUEMacRequest,
                                      TableAssignment, UESessionSet,
                                      UESessionContextResponse, UPFSessionContextState,
                                      VersionedPolicy)
from lte.protos.session_manager_pb2 import RuleRecordTable
from lte.protos.subscriberdb_pb2 import AggregatedMaximumBitrate
from magma.pipelined.app.check_quota import CheckQuotaController
from magma.pipelined.app.classifier import Classifier
from magma.pipelined.app.dpi import DPIController
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.enforcement_stats import EnforcementStatsController
from magma.pipelined.app.ipfix import IPFIXController
from magma.pipelined.app.ng_services import NGServiceController
from magma.pipelined.app.tunnel_learn import TunnelLearnController
from magma.pipelined.app.ue_mac import UEMacAddressController
from magma.pipelined.app.vlan_learn import VlanLearnController
from magma.pipelined.imsi import encode_imsi
from magma.pipelined.ipv6_prefix_store import (get_ipv6_interface_id,
                                               get_ipv6_prefix)
from magma.pipelined.metrics import (ENFORCEMENT_RULE_INSTALL_FAIL,
                                     ENFORCEMENT_STATS_RULE_INSTALL_FAIL)
from magma.pipelined.ng_manager.session_state_manager_util import PDRRuleEntry
from magma.pipelined.policy_converters import (convert_ipv4_str_to_ip_proto,
                                               convert_ipv6_bytes_to_ip_proto)

grpc_msg_queue = queue.Queue()
DEFAULT_CALL_TIMEOUT = 15


class PipelinedRpcServicer(pipelined_pb2_grpc.PipelinedServicer):
    """
    gRPC based server for Pipelined.
    """

    def __init__(self, loop, gy_app, enforcer_app, enforcement_stats, dpi_app,
                 ue_mac_app, check_quota_app, ipfix_app, vlan_learn_app,
                 tunnel_learn_app, classifier_app, inout_app, ng_servicer_app,
                 service_config, service_manager):
        self._loop = loop
        self._gy_app = gy_app
        self._enforcer_app = enforcer_app
        self._enforcement_stats = enforcement_stats
        self._dpi_app = dpi_app
        self._ue_mac_app = ue_mac_app
        self._check_quota_app = check_quota_app
        self._ipfix_app = ipfix_app
        self._vlan_learn_app = vlan_learn_app
        self._tunnel_learn_app = tunnel_learn_app
        self._service_config = service_config
        self._classifier_app = classifier_app
        self._inout_app = inout_app
        self._ng_servicer_app = ng_servicer_app
        self._service_manager = service_manager

        self._call_timeout = service_config.get('call_timeout',
                                                DEFAULT_CALL_TIMEOUT)
        self._print_grpc_payload = os.environ.get('MAGMA_PRINT_GRPC_PAYLOAD')
        if self._print_grpc_payload is None:
            self._print_grpc_payload = \
                service_config.get('magma_print_grpc_payload', False)

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        pipelined_pb2_grpc.add_PipelinedServicer_to_server(self, server)

    # --------------------------
    # General setup rpc
    # --------------------------

    def SetupDefaultControllers(self, request, _) -> SetupFlowsResult:
        """
        Setup default controllers, used on pipelined restarts
        """
        self._log_grpc_payload(request)
        ret = self._inout_app.check_setup_request_epoch(request.epoch)
        if ret is not None:
            return SetupFlowsResult(result=ret)

        fut = Future()
        self._loop.call_soon_threadsafe(self._setup_default_controllers, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("SetupDefaultControllers processing timed out")
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

    def _setup_default_controllers(self, fut: 'Future(SetupFlowsResult)'):
        res = self._inout_app.handle_restart(None)
        fut.set_result(res)

    # --------------------------
    # Enforcement App
    # --------------------------

    def SetupPolicyFlows(self, request, context) -> SetupFlowsResult:
        """
        Setup flows for all subscribers, used on pipelined restarts
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                EnforcementController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        for controller in [self._gy_app, self._enforcer_app,
                           self._enforcement_stats]:
            ret = controller.check_setup_request_epoch(request.epoch)
            if ret is not None:
                return SetupFlowsResult(result=ret)

        fut = Future()
        self._loop.call_soon_threadsafe(self._setup_flows, request, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("SetupPolicyFlows processing timed out")
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

    def _setup_flows(self, request: SetupPolicyRequest,
                     fut: 'Future[List[SetupFlowsResult]]'
                     ) -> SetupFlowsResult:
        gx_reqs = [req for req in request.requests
                   if req.request_origin.type == RequestOriginType.GX]
        gy_reqs = [req for req in request.requests
                   if req.request_origin.type == RequestOriginType.GY]
        enforcement_res = self._enforcer_app.handle_restart(gx_reqs)
        # TODO check these results and aggregate
        self._gy_app.handle_restart(gy_reqs)
        self._enforcement_stats.handle_restart(gx_reqs)
        fut.set_result(enforcement_res)

    def ActivateFlows(self, request, context):
        """
        Activate flows for a subscriber based on the pre-defined rules
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                EnforcementController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        for controller in [self._gy_app, self._enforcer_app,
                           self._enforcement_stats]:
            if not controller.is_controller_ready():
                context.set_code(grpc.StatusCode.UNAVAILABLE)
                context.set_details('Enforcement service not initialized!')
                return ActivateFlowsResult()

        fut = Future()  # type: Future[ActivateFlowsResult]
        self._loop.call_soon_threadsafe(self._activate_flows, request, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("ActivateFlows request processing timed out")
            return ActivateFlowsResult()

    def _update_ipv6_prefix_store(self, ipv6_addr: bytes):
        ipv6_str = ipv6_addr.decode('utf-8')
        interface = get_ipv6_interface_id(ipv6_str)
        prefix = get_ipv6_prefix(ipv6_str)
        self._service_manager.interface_to_prefix_mapper.save_prefix(
            interface, prefix)

    def _update_tunnel_map_store(self, uplink_tunnel: int,
                                 downlink_tunnel: int):
        self._service_manager.tunnel_id_mapper.save_tunnels(uplink_tunnel,
                                                            downlink_tunnel)

    def _update_version(self, request: ActivateFlowsRequest):
        """
        Update version for a given subscriber and rule.
        """
        for policy in request.policies:
            self._service_manager.session_rule_version_mapper.save_version(
                request.sid.id, request.uplink_tunnel, policy.rule.id,
                policy.version)

    def _remove_version(self, request: DeactivateFlowsRequest):
        def cleanup_dict(imsi, teid, rule_id, version):
            self._service_manager.session_rule_version_mapper \
                .remove(imsi, teid, rule_id, version)

        if not request.policies:
            self._service_manager.session_rule_version_mapper\
                .update_all_ue_versions(request.sid.id, request.uplink_tunnel)
            return

        for policy in request.policies:
            self._service_manager.session_rule_version_mapper \
                .save_version(request.sid.id, request.uplink_tunnel,
                              policy.rule_id, policy.version)
            cleanup_dict(request.sid.id, request.uplink_tunnel, policy.rule_id,
                         policy.version)

    def _activate_flows(self, request: ActivateFlowsRequest,
                        fut: 'Future[ActivateFlowsResult]'
                        ) -> None:
        """
        Activate flows for ipv4 / ipv6 or both

        CWF won't have an ip_addr passed
        """
        ret = ActivateFlowsResult()
        if self._service_config['setup_type'] == 'CWF' or request.ip_addr:
            ipv4 = convert_ipv4_str_to_ip_proto(request.ip_addr)
            if request.request_origin.type == RequestOriginType.GX:
                ret_ipv4 = self._install_flows_gx(request, ipv4)
            else:
                ret_ipv4 = self._install_flows_gy(request, ipv4)
            ret.policy_results.extend(ret_ipv4.policy_results)
        if request.ipv6_addr:
            ipv6 = convert_ipv6_bytes_to_ip_proto(request.ipv6_addr)
            self._update_ipv6_prefix_store(request.ipv6_addr)
            if request.request_origin.type == RequestOriginType.GX:
                ret_ipv6 = self._install_flows_gx(request, ipv6)
            else:
                ret_ipv6 = self._install_flows_gy(request, ipv6)
            ret.policy_results.extend(ret_ipv6.policy_results)
        if request.uplink_tunnel and request.downlink_tunnel:
            self._update_tunnel_map_store(request.uplink_tunnel,
                                          request.downlink_tunnel)

        fut.set_result(ret)

    def _install_flows_gx(self, request: ActivateFlowsRequest,
                         ip_address: IPAddress
                         ) -> ActivateFlowsResult:
        """
        Ensure that the RuleModResult is only successful if the flows are
        successfully added in both the enforcer app and enforcement_stats.
        Install enforcement_stats flows first because even if the enforcement
        flow install fails after, no traffic will be directed to the
        enforcement_stats flows.
        """
        logging.debug('Activating GX flows for %s', request.sid.id)
        self._update_version(request)
        # Install rules in enforcement stats
        enforcement_stats_res = self._activate_rules_in_enforcement_stats(
            request.sid.id, request.msisdn, request.uplink_tunnel,
            ip_address, request.apn_ambr, request.policies)

        failed_policies_results = \
            _retrieve_failed_results(enforcement_stats_res)
        # Do not install any rules that failed to install in enforcement_stats.
        policies = \
            _filter_failed_policies(request, failed_policies_results)

        enforcement_res = self._activate_rules_in_enforcement(
            request.sid.id, request.msisdn, request.uplink_tunnel,
            ip_address, request.apn_ambr, policies)

        # Include the failed rules from enforcement_stats in the response.
        enforcement_res.policy_results.extend(
            failed_policies_results)
        return enforcement_res

    def _install_flows_gy(self, request: ActivateFlowsRequest,
                          ip_address: IPAddress
                          ) -> ActivateFlowsResult:
        """
        Ensure that the RuleModResult is only successful if the flows are
        successfully added in both the enforcer app and enforcement_stats.
        Install enforcement_stats flows first because even if the enforcement
        flow install fails after, no traffic will be directed to the
        enforcement_stats flows.
        """
        logging.debug('Activating GY flows for %s', request.sid.id)
        self._update_version(request)
        # Install rules in enforcement stats
        enforcement_stats_res = self._activate_rules_in_enforcement_stats(
            request.sid.id, request.msisdn, request.uplink_tunnel,
            ip_address, request.apn_ambr, request.policies)

        failed_policies_results = \
            _retrieve_failed_results(enforcement_stats_res)
        # Do not install any rules that failed to install in enforcement_stats.
        policies = \
            _filter_failed_policies(request, failed_policies_results)

        gy_res = self._activate_rules_in_gy(request.sid.id, request.msisdn, request.uplink_tunnel,
                                            ip_address, request.apn_ambr,
                                            policies)

        # Include the failed rules from enforcement_stats in the response.
        gy_res.policy_results.extend(failed_policies_results)
        return gy_res

    def _activate_rules_in_enforcement_stats(self, imsi: str,
                                             msisdn: bytes,
                                             uplink_tunnel: int,
                                             ip_addr: IPAddress,
                                             apn_ambr: AggregatedMaximumBitrate,
                                             policies: List[VersionedPolicy]
                                             ) -> ActivateFlowsResult:
        if not self._service_manager.is_app_enabled(
                EnforcementStatsController.APP_NAME):
            return ActivateFlowsResult()

        enforcement_stats_res = self._enforcement_stats.activate_rules(
            imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, policies)
        _report_enforcement_stats_failures(enforcement_stats_res, imsi)
        return enforcement_stats_res

    def _activate_rules_in_enforcement(self, imsi: str, msisdn: bytes,
                                       uplink_tunnel: int,
                                       ip_addr: IPAddress,
                                       apn_ambr: AggregatedMaximumBitrate,
                                       policies: List[VersionedPolicy]
                                       ) -> ActivateFlowsResult:
        # TODO: this will crash pipelined if called with both static rules
        # and dynamic rules at the same time
        enforcement_res = self._enforcer_app.activate_rules(
            imsi, msisdn, uplink_tunnel, ip_addr, apn_ambr, policies)
        # TODO ?? Should the enforcement failure be reported per imsi session
        _report_enforcement_failures(enforcement_res, imsi)
        return enforcement_res

    def _activate_rules_in_gy(self, imsi: str, msisdn: bytes,
                              uplink_tunnel: int,
                              ip_addr: IPAddress,
                              apn_ambr: AggregatedMaximumBitrate,
                              policies: List[VersionedPolicy]
                              ) -> ActivateFlowsResult:
        gy_res = self._gy_app.activate_rules(imsi, msisdn, uplink_tunnel,
                                             ip_addr, apn_ambr, policies)
        # TODO: add metrics
        return gy_res

    def DeactivateFlows(self, request, context):
        """
        Deactivate flows for a subscriber
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                EnforcementController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        for controller in [self._gy_app, self._enforcer_app,
                           self._enforcement_stats]:
            if not controller.is_controller_ready():
                context.set_code(grpc.StatusCode.UNAVAILABLE)
                context.set_details('Enforcement service not initialized!')
                return ActivateFlowsResult()

        self._loop.call_soon_threadsafe(self._deactivate_flows, request)
        return DeactivateFlowsResult()

    def _deactivate_flows(self, request):
        """
        Deactivate flows for ipv4 / ipv6 or both

        CWF won't have an ip_addr passed
        """
        if self._service_config['setup_type'] == 'CWF' or request.ip_addr:
            ipv4 = convert_ipv4_str_to_ip_proto(request.ip_addr)
            if self._should_remove_from_gx(request):
                self._deactivate_flows_gx(request, ipv4)
            if self._should_remove_from_gy(request):
                self._deactivate_flows_gy(request, ipv4)
        if request.ipv6_addr:
            ipv6 = convert_ipv6_bytes_to_ip_proto(request.ipv6_addr)
            self._update_ipv6_prefix_store(request.ipv6_addr)
            if self._should_remove_from_gx(request):
                self._deactivate_flows_gx(request, ipv6)
            if self._should_remove_from_gy(request):
                self._deactivate_flows_gy(request, ipv6)

    def _should_remove_from_gy(self, request: DeactivateFlowsRequest) -> bool:
        is_gy = request.request_origin.type == RequestOriginType.GY
        is_wildcard = request.request_origin.type == RequestOriginType.WILDCARD
        return is_gy or is_wildcard

    def _should_remove_from_gx(self, request: DeactivateFlowsRequest) -> bool:
        is_gx = request.request_origin.type == RequestOriginType.GX
        is_wildcard = request.request_origin.type == RequestOriginType.WILDCARD
        return is_gx or is_wildcard

    def _deactivate_flows_gx(self, request, ip_address: IPAddress):
        logging.debug('Deactivating GX flows for %s', request.sid.id)
        self._remove_version(request)
        if request.remove_default_drop_flows:
            self._enforcement_stats.deactivate_default_flow(request.sid.id,
                                                            ip_address,
                                                            request.uplink_tunnel)
        rule_ids = [policy.rule_id for policy in request.policies]
        self._enforcer_app.deactivate_rules(request.sid.id, ip_address,
                                            request.uplink_tunnel, rule_ids)

    def _deactivate_flows_gy(self, request, ip_address: IPAddress):
        logging.debug('Deactivating GY flows for %s', request.sid.id)
        # Only deactivate requested rules here to not affect GX
        self._remove_version(request)
        rule_ids = [policy.rule_id for policy in request.policies]
        self._gy_app.deactivate_rules(request.sid.id, ip_address,
                                      request.uplink_tunnel, rule_ids)

    def GetPolicyUsage(self, request, context):
        """
        Get policy usage stats
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                EnforcementStatsController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        fut = Future()
        self._loop.call_soon_threadsafe(
            self._enforcement_stats.get_policy_usage, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("GetPolicyUsage processing timed out")
            return RuleRecordTable()

    # -------------------------
    # GRPC messages from MME
    #--------------------------
    def UpdateUEState(self, request, context):
        
        self._log_grpc_payload(request)

        if not self._service_manager.is_app_enabled(
              Classifier.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        fut = Future()
        self._loop.call_soon_threadsafe(\
                      self._setup_pg_tunnel_update, request, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("UpdateUEState processing timed out")
            return UESessionContextResponse(operation_type=request.ue_session_state.ue_config_state,
                                            cause_info=CauseIE(cause_ie=CauseIE.REQUEST_REJECTED_NO_REASON))

    def _setup_pg_tunnel_update(self, request: UESessionSet, fut: 'Future(UESessionContextResponse)'):
        res = self._classifier_app.process_mme_tunnel_request(request)
        fut.set_result(res)

    # --------------------------
    # IPFIX App
    # --------------------------

    def UpdateIPFIXFlow(self, request, context):
        """
        Update IPFIX sampling record
        """
        self._log_grpc_payload(request)
        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            # Install trace flow
            self._loop.call_soon_threadsafe(
                self._ipfix_app.add_ue_sample_flow, request.sid.id,
                request.msisdn, request.ap_mac_addr, request.ap_name,
                request.pdp_start_time)

        resp = FlowResponse()
        return resp

    # --------------------------
    # DPI App
    # --------------------------

    def CreateFlow(self, request, context):
        """
        Add dpi flow
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                DPIController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        resp = FlowResponse()
        self._loop.call_soon_threadsafe(self._dpi_app.add_classify_flow,
                                        request.match, request.state,
                                        request.app_name, request.service_type)
        return resp

    def RemoveFlow(self, request, context):
        """
        Add dpi flow
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                DPIController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        resp = FlowResponse()
        self._loop.call_soon_threadsafe(self._dpi_app.remove_classify_flow,
                                        request.match)
        return resp

    def UpdateFlowStats(self, request, context):
        """
        Update stats for a flow
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                DPIController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None
        resp = FlowResponse()
        return resp

    # --------------------------
    # UE MAC App
    # --------------------------

    def SetupUEMacFlows(self, request, context) -> SetupFlowsResult:
        """
        Activate a list of attached UEs
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                UEMacAddressController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        ret = self._ue_mac_app.check_setup_request_epoch(request.epoch)
        if ret is not None:
            return SetupFlowsResult(result=ret)

        fut = Future()
        self._loop.call_soon_threadsafe(self._setup_ue_mac,
                                        request, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("SetupUEMacFlows processing timed out")
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

    def _setup_ue_mac(self, request: SetupUEMacRequest,
                      fut: 'Future(SetupFlowsResult)'
                      ) -> SetupFlowsResult:
        res = self._ue_mac_app.handle_restart(request.requests)

        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            for req in request.requests:
                self._ipfix_app.add_ue_sample_flow(req.sid.id, req.msisdn,
                                                   req.ap_mac_addr,
                                                   req.ap_name,
                                                   req.pdp_start_time)

        fut.set_result(res)

    def AddUEMacFlow(self, request, context):
        """
        Associate UE MAC address to subscriber
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                UEMacAddressController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        # 12 hex characters + 5 colons
        if len(request.mac_addr) != 17:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details('Invalid UE MAC address provided')
            return None

        fut = Future()
        self._loop.call_soon_threadsafe(self._add_ue_mac_flow, request, fut)

        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("AddUEMacFlow processing timed out")
            return FlowResponse()

    def _add_ue_mac_flow(self, request, fut: 'Future(FlowResponse)'):
        res = self._ue_mac_app.add_ue_mac_flow(request.sid.id, request.mac_addr)

        fut.set_result(res)

    def DeleteUEMacFlow(self, request, context):
        """
        Delete UE MAC address to subscriber association
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                UEMacAddressController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        if not self._ue_mac_app.is_controller_ready():
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('UE MAC service not initialized!')
            return FlowResponse()

        # 12 hex characters + 5 colons
        if len(request.mac_addr) != 17:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            context.set_details('Invalid UE MAC address provided')
            return None

        self._loop.call_soon_threadsafe(
            self._ue_mac_app.delete_ue_mac_flow,
            request.sid.id, request.mac_addr)

        if self._service_manager.is_app_enabled(CheckQuotaController.APP_NAME):
            self._loop.call_soon_threadsafe(
                self._check_quota_app.remove_subscriber_flow, request.sid.id)

        if self._service_manager.is_app_enabled(VlanLearnController.APP_NAME):
            self._loop.call_soon_threadsafe(
                self._vlan_learn_app.remove_subscriber_flow, request.sid.id)

        if self._service_manager.is_app_enabled(TunnelLearnController.APP_NAME):
            self._loop.call_soon_threadsafe(
                self._tunnel_learn_app.remove_subscriber_flow, request.mac_addr)

        if self._service_manager.is_app_enabled(IPFIXController.APP_NAME):
            # Delete trace flow
            self._loop.call_soon_threadsafe(
                self._ipfix_app.delete_ue_sample_flow, request.sid.id)

        resp = FlowResponse()
        return resp

    # --------------------------
    # Check Quota App
    # --------------------------

    def SetupQuotaFlows(self, request, context) -> SetupFlowsResult:
        """
        Activate a list of quota rules
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                CheckQuotaController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        ret = self._check_quota_app.check_setup_request_epoch(request.epoch)
        if ret is not None:
            return SetupFlowsResult(result=ret)

        fut = Future()
        self._loop.call_soon_threadsafe(self._setup_quota,
                                        request, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("SetupQuotaFlows processing timed out")
            return SetupFlowsResult(result=SetupFlowsResult.FAILURE)

    def _setup_quota(self, request: SetupQuotaRequest,
                     fut: 'Future(SetupFlowsResult)'
                     ) -> SetupFlowsResult:
        res = self._check_quota_app.handle_restart(request.requests)
        fut.set_result(res)

    def UpdateSubscriberQuotaState(self, request, context):
        """
        Updates the subcsciber quota state
        """
        self._log_grpc_payload(request)
        if not self._service_manager.is_app_enabled(
                CheckQuotaController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return None

        if not self._check_quota_app.is_controller_ready():
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Check Quota service not initialized!')
            return FlowResponse()

        resp = FlowResponse()
        self._loop.call_soon_threadsafe(
            self._check_quota_app.update_subscriber_quota_state, request.updates)
        return resp

    # --------------------------
    # Debugging
    # --------------------------

    def GetAllTableAssignments(self, request, context):
        """
        Get the flow table assignment for all apps ordered by main table number
        and name
        """
        self._log_grpc_payload(request)
        table_assignments = self._service_manager.get_all_table_assignments()
        return AllTableAssignments(table_assignments=[
            TableAssignment(app_name=app_name, main_table=tables.main_table,
                            scratch_tables=tables.scratch_tables) for
            app_name, tables in table_assignments.items()])

    # --------------------------
    # Internal
    # --------------------------

    def _log_grpc_payload(self, grpc_request):
        if not grpc_request:
            return
        indent = '  '
        dbl_indent = indent + indent
        indented_text = dbl_indent + \
            str(grpc_request).replace('\n', '\n' + dbl_indent)
        log_msg = 'Got RPC payload:\n{0}{1} {{\n{2}\n{0}}}'.format(indent,
            grpc_request.DESCRIPTOR.name, indented_text.rstrip())

        grpc_msg_queue.put(log_msg)
        if grpc_msg_queue.qsize() > 100:
            grpc_msg_queue.get()

        if not self._print_grpc_payload:
            return
        logging.info(log_msg)

    def SetSMFSessions(self, request, context):
        """
        Setup the 5G Session flows for the subscriber
        """
        #if 5G Services are not enabled return UNAVAILABLE
        if not self._service_manager.is_ng_app_enabled(
                NGServiceController.APP_NAME):
            context.set_code(grpc.StatusCode.UNAVAILABLE)
            context.set_details('Service not enabled!')
            return UPFSessionContextState()

        fut = Future()
        self._log_grpc_payload(request)
        self._loop.call_soon_threadsafe(\
                      self.ng_update_session_flows, request, fut)
        try:
            return fut.result(timeout=self._call_timeout)
        except concurrent.futures.TimeoutError:
            logging.error("SetSMFSessions processing timed out")
            return UPFSessionContextState()

    def ng_update_session_flows(self, request: SessionSet,
                                fut: 'Future(UPFSessionContextState)') -> UPFSessionContextState:
        """
        Install PDR, FAR and QER flows for the 5G Session send by SMF
        """
        logging.debug('Update 5G Session Flow for SessionID:%s, SessionVersion:%d',
                      request.subscriber_id, request.session_version)

        # Convert message containing PDR to Named Tuple Rules.
        process_pdr_rules = OrderedDict()
        response = self._ng_servicer_app.ng_session_message_handler(request, process_pdr_rules)

        # Failure in message processing return failure
        if response.cause_info.cause_ie == CauseIE.REQUEST_ACCEPTED:
            for _, pdr_entries in process_pdr_rules.items():
                # Create the Tunnel
                ret = self._ng_tunnel_update(pdr_entries, request.subscriber_id)
                if ret == False:
                    offending_ie = OffendingIE(identifier=pdr_entries.pdr_id,
                                               version=pdr_entries.pdr_version)

                    #Session information is filled already
                    response.cause_info.cause_ie = CauseIE.RULE_CREATION_OR_MODIFICATION_FAILURE
                    response.failure_rule_id.pdr.extend([offending_ie])
                    break

        fut.set_result(response)

    def _ng_tunnel_update(self, pdr_entry: PDRRuleEntry, subscriber_id: str) -> bool:

        ret = self._classifier_app.gtp_handler(pdr_entry.pdr_state,
                                                pdr_entry.precedence,
                                                pdr_entry.local_f_teid,
                                                pdr_entry.far_action.o_teid,
                                                pdr_entry.ue_ip_addr,
                                                pdr_entry.far_action.gnb_ip_addr,
                                                encode_imsi(subscriber_id),
                                                True)

        return ret

def _retrieve_failed_results(activate_flow_result: ActivateFlowsResult
                             ) -> Tuple[List[RuleModResult],
                                        List[RuleModResult]]:
    failed_policies_results = \
        [result for result in
         activate_flow_result.policy_results if
         result.result == RuleModResult.FAILURE]
    return failed_policies_results


def _filter_failed_policies(request: ActivateFlowsRequest,
                            failed_results: List[RuleModResult]
                            ) -> List[VersionedPolicy]:
    failed_policies = [result.rule_id for result in failed_results]
    return [policy for policy in request.policies if
            policy.rule.id not in failed_policies]


def _report_enforcement_failures(activate_flow_result: ActivateFlowsResult,
                                 imsi: str):
    for result in activate_flow_result.policy_results:
        if result.result == RuleModResult.SUCCESS:
            continue
        ENFORCEMENT_RULE_INSTALL_FAIL.labels(rule_id=result.rule_id,
                                             imsi=imsi).inc()


def _report_enforcement_stats_failures(
        activate_flow_result: ActivateFlowsResult,
        imsi: str):
    for result in activate_flow_result.policy_results:
        if result.result == RuleModResult.SUCCESS:
            continue
        ENFORCEMENT_STATS_RULE_INSTALL_FAIL.labels(rule_id=result.rule_id,
                                                   imsi=imsi).inc()
