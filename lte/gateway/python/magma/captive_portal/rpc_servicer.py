"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import logging

from grpc import RpcError
from lte.protos import session_manager_pb2, session_manager_pb2_grpc
from lte.protos.pipelined_pb2 import ActivateFlowsRequest, \
    DeactivateFlowsRequest
from lte.protos.pipelined_pb2_grpc import PipelinedStub
from lte.protos.policydb_pb2 import FlowDescription, FlowMatch, PolicyRule, \
    RedirectInformation
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from orc8r.protos.common_pb2 import Void

from magma.common.misc_utils import get_ip_from_if
from magma.common.service_registry import ServiceRegistry


class SessionRpcServicer(session_manager_pb2_grpc.LocalSessionManagerServicer):
    """
    gRPC based server for LocalSessionManager service
    """
    ALLOW_ALL_PRIORITY = 100
    REDIRECT_PRIORITY = 2000
    RPC_TIMEOUT = 5

    def __init__(self, service):
        chan = ServiceRegistry.get_rpc_channel('pipelined',
                                               ServiceRegistry.LOCAL)
        self._pipelined = PipelinedStub(chan)
        chan = ServiceRegistry.get_rpc_channel('subscriberdb',
                                               ServiceRegistry.LOCAL)
        self._subscriberdb = SubscriberDBStub(chan)
        self._enabled = service.config['captive_portal_enabled']
        self._captive_portal_address = service.config['captive_portal_url']
        self._local_ip = get_ip_from_if(service.config['bridge_interface'])
        self._whitelisted_ips = service.config['whitelisted_ips']
        self._sub_profile_substr = service.config[
            'subscriber_profile_substr_match']

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        session_manager_pb2_grpc.add_LocalSessionManagerServicer_to_server(
            self, server,
        )

    def CreateSession(self, request, context):
        """
        Handles create session request from MME by installing the necessary
        flows in pipelined's enforcement app.
        """
        sid = request.sid
        logging.info('Create session request for sid: %s', sid.id)
        try:
            # Gather the set of policy rules to use
            if self._captive_portal_enabled(sid):
                rules = []
                rules.extend(self._get_whitelisted_policies())
                rules.extend(self._get_redirect_policies())
            else:
                rules = self._get_allow_all_traffic_policies()

            # Activate the flows in the enforcement app in pipelined
            act_request = ActivateFlowsRequest(
                sid=sid, ip_addr=request.ue_ipv4, dynamic_rules=rules)
            act_response = self._pipelined.ActivateFlows(
                act_request, timeout=self.RPC_TIMEOUT)
            for res in act_response.dynamic_rule_results:
                if res.result != res.SUCCESS:
                    # Hmm rolling back partial success is difficult
                    # Let's just log this for now
                    logging.error('Failed to activate rule: %s', res.rule_id)
        except RpcError as err:
            self._set_rpc_error(context, err)
        return session_manager_pb2.LocalCreateSessionResponse()

    def EndSession(self, sid, context):  # pylint: disable=arguments-differ
        """
        Handles end session request from MME by removing all the flows
        for the subscriber in pipelined's enforcement app.
        """
        logging.info('End session request for sid: %s', sid.id)
        try:
            self._pipelined.DeactivateFlows(
                DeactivateFlowsRequest(sid=sid), timeout=self.RPC_TIMEOUT)
        except RpcError as err:
            self._set_rpc_error(context, err)
        return session_manager_pb2.LocalEndSessionResponse()

    def ReportRuleStats(self, request, context):
        """
        Handles stats update from the enforcement app in pipelined. We are
        ignoring this for now, since the applications can poll pipelined for
        the flow stats.
        """
        logging.debug('Ignoring ReportRuleStats rpc')
        return Void()

    def _captive_portal_enabled(self, sid):
        if not self._enabled:
            return False  # Service is disabled

        if self._sub_profile_substr == '':
            return True  # Allow all subscribers

        sub = self._subscriberdb.GetSubscriberData(
            sid, timeout=self.RPC_TIMEOUT)
        return self._sub_profile_substr in sub.sub_profile

    def _get_allow_all_traffic_policies(self):
        """ Policy to allow all traffic to the internet """
        return [PolicyRule(
            id='allow_all_traffic',
            priority=self.ALLOW_ALL_PRIORITY,
            flow_list=[
                FlowDescription(match=FlowMatch(direction=FlowMatch.UPLINK)),
                FlowDescription(match=FlowMatch(direction=FlowMatch.DOWNLINK)),
            ],
        )]

    def _get_whitelisted_policies(self):
        """ Policies to allow http traffic to the whitelisted sites """
        rules = []
        for ip, ports in self._whitelisted_ips.items():
            for port in ports:
                if ip == 'local':
                    ip = self._local_ip
                rules.append(PolicyRule(
                    id='whitelist',
                    priority=self.ALLOW_ALL_PRIORITY,
                    flow_list=[
                        FlowDescription(
                            match=FlowMatch(
                                direction=FlowMatch.UPLINK,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                ipv4_dst=ip,
                                tcp_dst=port)),
                        FlowDescription(
                            match=FlowMatch(
                                direction=FlowMatch.DOWNLINK,
                                ip_proto=FlowMatch.IPPROTO_TCP,
                                ipv4_src=ip,
                                tcp_src=port)),
                    ]))
        return rules

    def _get_redirect_policies(self):
        """ Policy to redirect traffic to the captive portal """
        redirect_info = RedirectInformation(
            support=RedirectInformation.ENABLED,
            address_type=RedirectInformation.URL,
            server_address=self._captive_portal_address)
        return [PolicyRule(
            id='redirect',
            priority=self.REDIRECT_PRIORITY,
            redirect=redirect_info)]

    def _set_rpc_error(self, context, err):
        logging.error(err.details())
        context.set_details(err.details())
        context.set_code(err.code())
