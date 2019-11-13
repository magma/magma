"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import logging
from datetime import datetime
from typing import Any, Dict, List
from lte.protos.policydb_pb2 import PolicyRule, FlowMatch, FlowDescription, \
    FlowQos, QosArp
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    CreateSessionResponse, UpdateSessionRequest, SessionTerminateResponse, \
    UpdateSessionResponse, DynamicRuleInstall, StaticRuleInstall, \
    ChargingCredit, GrantedUnits, CreditUnit, CreditUpdateResponse
from lte.protos.session_manager_pb2_grpc import \
    CentralSessionControllerServicer, \
    add_CentralSessionControllerServicer_to_server


class SessionRpcServicer(CentralSessionControllerServicer):
    """
    gRPC based server for CentralSessionController service.

    This will act as a bare-bones local PCRF and OCS.
    In current implementation, it is only used for enabling the Captive Portal
    feature.
    """

    def __init__(self, service_config):
        self._config = service_config

    @property
    def config(self) -> Dict[str, Any]:
        return self._config

    @property
    def redirect_rule_name(self) -> str:
        return self._config['redirect_rule_name']

    @property
    def ip_whitelist(self) -> Dict[str, List[int]]:
        """
        Get whitelisted IP/port combinations which have all traffic allowed.

        Returns: A dict keyed by IP, with values being a list of ports
        """
        return self._config['whitelisted_ips']

    def add_to_server(self, server):
        """ Add the servicer to a gRPC server """
        add_CentralSessionControllerServicer_to_server(
            self, server,
        )

    def CreateSession(
        self,
        request: CreateSessionRequest,
        context,
    ) -> CreateSessionResponse:
        """
        Handles create session request from MME by installing the necessary
        flows in pipelined's enforcement app.
        """
        logging.info('Creating a session for subscriber ID: %s',
                     request.subscriber.id)
        return CreateSessionResponse(
            credits=[self._get_credit_update_response(request.imsi_plmn_id)],
            static_rules=[self._get_static_rule()],
            dynamic_rules=self._get_whitelist_rules(),
        )

    def UpdateSession(
            self,
            request: UpdateSessionRequest,
            context,
    ) -> UpdateSessionResponse:
        """
        On UpdateSession, return an arbitrarily large amount of additional
        credit for the session.
        """
        logging.debug('Updating sessions')
        resp = UpdateSessionResponse()
        for credit_usage_update in request.updates:
            resp.responses.extend(
                [self._get_credit_update_response(credit_usage_update.sid)],
            )
        return resp

    def TerminateSession(
            self,
            request: SessionTerminateResponse,
            context,
    ) -> SessionTerminateResponse:
        logging.info('Terminating a session for session ID: %s',
                     request.session_id)
        return SessionTerminateResponse(
            sid=request.sid,
            session_id=request.session_id,
        )

    def _get_whitelist_rules(self) -> List[DynamicRuleInstall]:
        """
        Get a list of dynamic rules to install for whitelisting.
        These rules will whitelist traffic to/from the captive portal server.
        """
        dynamic_rules = []
        for ip, ports in self.ip_whitelist.items():
            if ip == 'local':
                ip = '192.168.128.1'
            for port in ports:
                # Build the rule id to be globally unique
                rule_id_info = {
                    'ip': ip,
                    'port': port,
                }
                rule_id = "whitelist_policy_id-{ip}:{port}"\
                    .format(**rule_id_info)

                rule = DynamicRuleInstall(
                    policy_rule=self._get_whitelist_policy_rule(
                        rule_id, ip, port
                    ),
                )
                # Activate now, and deactivate long in the future
                t2 = datetime.now()
                t2 = t2.replace(year=t2.year + 1)
                rule.activation_time.FromDatetime(datetime.now())
                rule.deactivation_time.FromDatetime(t2)
                dynamic_rules.append(rule)
        return dynamic_rules

    def _get_whitelist_policy_rule(
        self,
        policy_id: str,
        ip: str,
        port: int,
    ) -> PolicyRule:
        return PolicyRule(
            # Don't set the rating group
            # Don't set the monitoring key
            # Don't set the hard timeout
            id=policy_id,
            priority=100,
            qos=self._get_default_qos(),
            flow_list=self._get_whitelist_flows(ip, port),
            tracking_type=PolicyRule.TrackingType.Value("NO_TRACKING"),
        )

    def _get_whitelist_flows(self, ip: str, port: int) -> List[FlowDescription]:
        """
        Args:
            ip: IP address to allow traffic to. This should be the captive
                portal address
            port:  Port of the captive portal server. Probably 80.

        Returns:
            Two flows, one for traffic towards the captive portal server, and a
            second for traffic from the captive portal server.
        """
        return [
            # Set flow match for outgoing TCP packets to whitelisted IP
            # Don't set the app_name field
            FlowDescription(  # uplink flow
                match=FlowMatch(
                    ipv4_dst=ip,
                    tcp_dst=port,
                    ip_proto=FlowMatch.IPProto.Value("IPPROTO_TCP"),
                    direction=FlowMatch.Direction.Value("UPLINK"),
                ),
                action=FlowDescription.Action.Value("PERMIT"),
            ),
            # Set flow match for incoming TCP packets from whitelisted IP
            # Don't set the app_name field
            FlowDescription(  # downlink flow
                match=FlowMatch(
                    ipv4_src=ip,
                    tcp_src=port,
                    ip_proto=FlowMatch.IPProto.Value("IPPROTO_TCP"),
                    direction=FlowMatch.Direction.Value("DOWNLINK"),
                ),
                action=FlowDescription.Action.Value("PERMIT"),
            ),
        ]

    def _get_default_qos(self) -> FlowQos:
        """
        Get a default QoS, usable for an allow-all flow, and for redirection to
        a captive_portal.
        """
        return FlowQos(
            max_req_bw_ul=2 * 1024 * 1024 * 1024,  # 2G
            max_req_bw_dl=2 * 1024 * 1024 * 1024,  # 2G
            gbr_ul=1 * 1024 * 1024,  # 1 Mb/s
            gbr_dl=1 * 1024 * 1024,  # 1 Mb/s
            qci=FlowQos.Qci.Value('QCI_3'),
            # Allocation and Retention Policy
            # Set to high priority, and disallow pre-emption
            # capability/vulnerability
            arp=QosArp(
                priority_level=1,
                pre_capability=QosArp.PreCap.Value("PRE_CAP_DISABLED"),
                pre_vulnerability=QosArp.PreVul.Value("PRE_VUL_DISABLED"),
            ),
        )

    def _get_static_rule(self) -> StaticRuleInstall:
        """ Return a static rule for redirection to captive portal """
        return StaticRuleInstall(
            rule_id = self.redirect_rule_name,
        )

    def _get_credit_update_response(
        self,
        sid: str,
    ) -> CreditUpdateResponse:
        return CreditUpdateResponse(
            success=True,
            sid=sid,
            charging_key=1,
            credit=self._get_max_charging_credit(),
            type=CreditUpdateResponse.ResponseType.Value('UPDATE'),
            result_code=1,
        )

    def _get_max_charging_credit(self) -> ChargingCredit:
        return ChargingCredit(
            type=ChargingCredit.UnitType.Value('SECONDS'),
            validity_time=86400,  # One day
            is_final=False,
            final_action=ChargingCredit.FinalAction.Value('TERMINATE'),
            granted_units=self._get_max_granted_units(),
        )

    def _get_max_granted_units(self) -> GrantedUnits:
        """ Get an arbitrarily large amount of granted credit """
        return GrantedUnits(
            total=CreditUnit(
                is_valid=True,
                volume=100 * 1024 * 1024 * 1024,  # 100 GiB
            ),
            rx=CreditUnit(
                is_valid=True,
                volume=50 * 1024 * 1024 * 1024,  # 50 GiB
            ),
            tx=CreditUnit(
                is_valid=True,
                volume=50 * 1024 * 1024 * 1024,  # 50 GiB
            ),
        )
