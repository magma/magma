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

import unittest

from lte.protos.mconfig import mconfigs_pb2
from lte.protos.policydb_pb2 import (
    ApnPolicySet,
    AssignedPolicies,
    ChargingRuleNameSet,
    FlowDescription,
    FlowMatch,
    PolicyRule,
    RatingGroup,
    SubscriberPolicySet,
)
from lte.protos.session_manager_pb2 import (
    CommonSessionContext,
    CreateSessionRequest,
    CreditLimitType,
    CreditUpdateResponse,
    CreditUsageUpdate,
    LTESessionContext,
    RatSpecificContext,
    SessionTerminateRequest,
    UpdateSessionRequest,
)
from lte.protos.subscriberdb_pb2 import (
    LTESubscription,
    SubscriberData,
    SubscriberID,
)
from magma.policydb.servicers.session_servicer import SessionRpcServicer

CSR_STATIC_RULES = '[rule_id: "redirect"]'
CSR_STATIC_RULES_2 = '[rule_id: "p6"]'
CSR_STATIC_RULES_3 = '[rule_id: "p5"]'


USR = '''
  responses {
    success: true
    sid: "abc"
    charging_key: 1
    limit_type: INFINITE_UNMETERED
  }
  responses {
    success: true
    sid: "abc"
    charging_key: 2
    limit_type: INFINITE_METERED
  }'''


class SessionRpcServicerTest(unittest.TestCase):
    def setUp(self):
        rating_groups_by_id = {
            1: RatingGroup(
                id=1,
                limit_type=RatingGroup.INFINITE_UNMETERED,
            ),
            2: RatingGroup(
                id=2,
                limit_type=RatingGroup.INFINITE_METERED,
            ),
        }
        basenames_dict = {
            'bn1': ChargingRuleNameSet(RuleNames=['p5']),
            'bn2': ChargingRuleNameSet(RuleNames=['p6']),
        }
        apn_rules_by_sid = {
            "IMSI1234": SubscriberPolicySet(
                rules_per_apn=[
                    ApnPolicySet(
                        apn="apn1",
                        assigned_base_names=[],
                        assigned_policies=["redirect"],
                    ),
                ],
            ),
            "IMSI2345": SubscriberPolicySet(
                rules_per_apn=[
                    ApnPolicySet(
                        apn="apn1",
                        assigned_base_names=["bn1"],
                        assigned_policies=[],
                    ),
                    ApnPolicySet(
                        apn="apn2",
                        assigned_base_names=["bn2"],
                        assigned_policies=[],
                    ),
                ],
            ),
            "IMSI3456": SubscriberPolicySet(
                global_base_names=["bn1"],
                global_policies=[],
                rules_per_apn=[
                    ApnPolicySet(
                        apn="apn1",
                        assigned_base_names=[],
                        assigned_policies=[],
                    ),
                ],
            ),
        }
        self.servicer = SessionRpcServicer(
            self._get_mconfig(),
            rating_groups_by_id,
            basenames_dict,
            apn_rules_by_sid,
        )

    def _get_mconfig(self) -> mconfigs_pb2.PolicyDB:
        return mconfigs_pb2.PolicyDB(
            log_level=1,
        )

    def tearDown(self):
        # TODO: not sure if this is actually needed
        self.servicer = None

    def test_CreateSession(self):
        """
        Create a session

        Assert:
            There is a high volume of granted credits
            A static rule is installed for redirection to the captive portal
            server
            Two dynamic rules are installed for traffic from UE to captive
            portal and vice versa
        """
        msg = CreateSessionRequest(
            session_id='1234',
            common_context=CommonSessionContext(
                sid=SubscriberID(
                    id='IMSI1234',
                ),
                apn='apn1',
            ),
            rat_specific_context=RatSpecificContext(
                lte_context=LTESessionContext(
                    imsi_plmn_id='00101',
                ),
            ),
        )
        resp = self.servicer.CreateSession(msg, None)

        # There should be a static rule installed for the redirection
        static_rules = self._rm_whitespace(str(resp.static_rules))
        expected = self._rm_whitespace(CSR_STATIC_RULES)
        self.assertEqual(
            static_rules, expected, 'There should be one static '
            'rule installed for redirection.',
        )

        # Credit granted should be unlimited and un-metered
        credit_limit_type = resp.credits[0].limit_type
        expected = CreditLimitType.Value("INFINITE_UNMETERED")
        self.assertEqual(
            credit_limit_type, expected, 'There should be an '
            'infinite, unmetered credit grant',
        )

        msg_2 = CreateSessionRequest(
            session_id='2345',
            common_context=CommonSessionContext(
                sid=SubscriberID(
                    id='IMSI2345',
                ),
                apn='apn2',
            ),
            rat_specific_context=RatSpecificContext(
                lte_context=LTESessionContext(
                    imsi_plmn_id='00101',
                ),
            ),
        )
        resp = self.servicer.CreateSession(msg_2, None)

        # There should be a static rule installed
        static_rules = self._rm_whitespace(str(resp.static_rules))
        expected = self._rm_whitespace(CSR_STATIC_RULES_2)
        self.assertEqual(
            static_rules, expected, 'There should be one static '
            'rule installed.',
        )

        # Credit granted should be unlimited and un-metered
        credit_limit_type = resp.credits[0].limit_type
        expected = CreditLimitType.Value("INFINITE_UNMETERED")
        self.assertEqual(
            credit_limit_type, expected,
            'There should be an infinite, unmetered credit grant',
        )

        msg_3 = CreateSessionRequest(
            session_id='3456',
            common_context=CommonSessionContext(
                sid=SubscriberID(
                    id='IMSI3456',
                ),
                apn='apn1',
            ),
            rat_specific_context=RatSpecificContext(
                lte_context=LTESessionContext(
                    imsi_plmn_id='00101',
                ),
            ),
        )
        resp = self.servicer.CreateSession(msg_3, None)

        # There should be a static rule installed
        static_rules = self._rm_whitespace(str(resp.static_rules))
        expected = self._rm_whitespace(CSR_STATIC_RULES_3)
        self.assertEqual(
            static_rules, expected, 'There should be one static '
            'rule installed.',
        )

        # Credit granted should be unlimited and un-metered
        credit_limit_type = resp.credits[0].limit_type
        expected = CreditLimitType.Value("INFINITE_UNMETERED")
        self.assertEqual(
            credit_limit_type, expected,
            'There should be an infinite, unmetered credit grant',
        )

    def test_UpdateSession(self):
        """
        Update a session

        Assert:
            An arbitrarily large amount of credit should be granted.
        """
        msg = UpdateSessionRequest()
        credit_update = CreditUsageUpdate()
        credit_update.common_context.sid.id = 'abc'
        msg.updates.extend([credit_update])
        resp = self.servicer.UpdateSession(msg, None)
        session_response = self._rm_whitespace(str(resp))
        expected = self._rm_whitespace(USR)
        self.assertEqual(
            session_response, expected, 'There should be an '
            'infinite, unmetered credit grant and an infinite, '
            'metered credit grant',
        )

    def test_TerminateSession(self):
        """
        Terminate a session

        Assert:
            The session can be terminated successfully. Will always succeed.
        """
        msg = SessionTerminateRequest()
        msg.common_context.sid.id = 'abc'
        msg.session_id = 'session_id_123'
        resp = self.servicer.TerminateSession(msg, None)
        self.assertEqual(resp.sid, 'abc', 'SID should be same as request')
        self.assertEqual(
            resp.session_id, 'session_id_123', 'session ID should '
            'be same as request',
        )

    def _rm_whitespace(self, inp: str) -> str:
        return inp.replace(' ', '').replace('\n', '')
