"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import grpc
import logging
from typing import List
from lte.protos.mconfig import mconfigs_pb2
from lte.protos.policydb_pb2 import PolicyRule, FlowDescription, \
    FlowMatch, RatingGroup
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    CreateSessionResponse, UpdateSessionRequest, SessionTerminateResponse, \
    UpdateSessionResponse, StaticRuleInstall, DynamicRuleInstall,\
    CreditUpdateResponse, CreditLimitType
from lte.protos.session_manager_pb2_grpc import \
    CentralSessionControllerServicer, \
    add_CentralSessionControllerServicer_to_server
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.policydb.rating_group_store import RatingGroupsDict
from orc8r.protos.common_pb2 import NetworkID


class SessionRpcServicer(CentralSessionControllerServicer):
    """
    gRPC based server for CentralSessionController service.

    This will act as a bare-bones local PCRF and OCS.
    In current implementation, it is only used for enabling the Captive Portal
    feature.
    """

    def __init__(
        self,
        mconfig: mconfigs_pb2.PolicyDB,
        rating_groups_by_id: RatingGroupsDict,
        subscriberdb_stub: SubscriberDBStub,
    ):
        self._mconfig = mconfig
        self._network_id = NetworkID(id="_")
        self._rating_groups_by_id = rating_groups_by_id
        self._subscriberdb_stub = subscriberdb_stub

    def get_infinite_credit_charging_keys(self) -> List[int]:
        keys = []
        for rating_group in self._rating_groups_by_id.values():
            if rating_group.limit_type == RatingGroup.INFINITE_UNMETERED:
                keys.append(rating_group.id)
        return keys

    def _get_postpay_charging_keys(self) -> List[int]:
        keys = []
        for rating_group in self._rating_groups_by_id.values():
            if rating_group.limit_type == RatingGroup.INFINITE_METERED:
                keys.append(rating_group.id)
        return keys

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

        NOTE: truncate the 'IMSI' prefix
        """
        imsi = request.subscriber.id
        imsi_number = imsi[4:]
        logging.info('Creating a session for subscriber ID: %s', imsi)
        return CreateSessionResponse(
            credits=self._get_credits(imsi),
            static_rules=self._get_rules_for_imsi(imsi_number),
            dynamic_rules=self._get_default_dynamic_rules(imsi_number),
            session_id=request.session_id,
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
                self._get_credits(credit_usage_update.sid),
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

    def _get_default_dynamic_rules(
        self,
        sid: str,
    ) -> List[DynamicRuleInstall]:
        """
        Get a list of dynamic rules to install for whitelisting.
        """
        dynamic_rules = []
        # Build the rule id to be globally unique
        rule_id_info = {'sid': sid}
        rule_id = "whitelist_sid-{sid}".format(**rule_id_info)
        rule = DynamicRuleInstall(
            policy_rule=self._get_allow_all_policy_rule(rule_id),
        )
        dynamic_rules.append(rule)
        return dynamic_rules

    def _get_allow_all_policy_rule(
        self,
        policy_id: str,
    ) -> PolicyRule:
        """
        This builds a PolicyRule used as a default to allow traffic
        for an attached subscriber.
        """
        return PolicyRule(
            # Don't set the rating group
            # Don't set the monitoring key
            # Don't set the hard timeout
            id=policy_id,
            priority=2,
            flow_list=self._get_allow_all_flows(),
            tracking_type=PolicyRule.TrackingType.Value("NO_TRACKING"),
        )

    def _get_allow_all_flows(self) -> List[FlowDescription]:
        """
        Returns:
            Two flows, for outgoing and incoming traffic
        """
        return [
            # Set flow match for ll packets
            # Don't set the app_name field
            FlowDescription(  # uplink flow
                match=FlowMatch(
                    direction=FlowMatch.Direction.Value("UPLINK"),
                ),
                action=FlowDescription.Action.Value("PERMIT"),
            ),
            FlowDescription(  # downlink flow
                match=FlowMatch(
                    direction=FlowMatch.Direction.Value("DOWNLINK"),
                ),
                action=FlowDescription.Action.Value("PERMIT"),
            ),
        ]

    def _get_rules_for_imsi(self, imsi: str) -> List[StaticRuleInstall]:
        """
        Get the list of static rules to be installed for a subscriber
        NOTE: Remove "IMSI" prefix from imsi argument.
        """
        try:
            info = self._subscriberdb_stub.GetSubscriberData(NetworkID(id=imsi))
            return [StaticRuleInstall(rule_id=rule_id)
                    for rule_id in info.lte.assigned_policies]
        except grpc.RpcError:
            logging.error('Unable to find data for subscriber %s', imsi)
            return []

    def _get_credits(self, sid: str) -> List[CreditUpdateResponse]:
        infinite_credit_keys = self.get_infinite_credit_charging_keys()
        postpay_keys = self._get_postpay_charging_keys()
        credit_updates = []
        for charging_key in infinite_credit_keys:
            credit_updates.append(CreditUpdateResponse(
                success=True,
                sid=sid,
                charging_key=charging_key,
                limit_type=CreditLimitType.Value("INFINITE_UNMETERED")
            ))
        for charging_key in postpay_keys:
            credit_updates.append(CreditUpdateResponse(
                success=True,
                sid=sid,
                charging_key=charging_key,
                limit_type=CreditLimitType.Value("INFINITE_METERED")
            ))
        return credit_updates
