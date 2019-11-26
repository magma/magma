"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import logging
from typing import Any, Dict, List
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    CreateSessionResponse, UpdateSessionRequest, SessionTerminateResponse, \
    UpdateSessionResponse, StaticRuleInstall, \
    ChargingCredit, GrantedUnits, CreditUnit, CreditUpdateResponse
from lte.protos.session_manager_pb2_grpc import \
    CentralSessionControllerServicer, \
    add_CentralSessionControllerServicer_to_server
from magma.policydb.rule_store import PolicyRuleDict


class SessionRpcServicer(CentralSessionControllerServicer):
    """
    gRPC based server for CentralSessionController service.

    This will act as a bare-bones local PCRF and OCS.
    In current implementation, it is only used for enabling the Captive Portal
    feature.
    """

    def __init__(self, service_config):
        self._config = service_config
        self._rules = PolicyRuleDict()

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
            credits=[],
            static_rules=self._get_static_rules(),
            dynamic_rules=[],
        )

    def UpdateSession(
        self,
        request: UpdateSessionRequest,
        context,
    ) -> UpdateSessionResponse:
        """
        On UpdateSession, return an arbitrarily large amount of additional
        credit for the session.

        NOTE: This really shouldn't be called, as no credit should have been
        granted on CreateSession.
        """
        logging.info('UpdateSession called')
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

    def _get_static_rules(self) -> List[StaticRuleInstall]:
        """ Return a static rule for redirection to captive portal """
        if self.redirect_rule_name in self._rules:
            return [StaticRuleInstall(rule_id=self.redirect_rule_name)]
        return []

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
