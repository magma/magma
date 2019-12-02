"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from typing import Any, Dict, List
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    UpdateSessionRequest, CreditUsageUpdate, SessionTerminateRequest
from lte.protos.subscriberdb_pb2 import SubscriberData, LTESubscription
from orc8r.protos.common_pb2 import NetworkID
from magma.policydb.rpc_servicer import SessionRpcServicer


CSR_STATIC_RULES = '[rule_id: "redirect"]'


USR = '''
  responses {
    success: true
    sid: "abc"
    charging_key: 1
    credit {
      type: SECONDS
      validity_time: 86400
      granted_units {
        total {
          is_valid: true
          volume: 107374182400
        }
        tx {
          is_valid: true
          volume: 53687091200
        }
        rx {
          is_valid: true
          volume: 53687091200
        }
      }
    }
    result_code: 1
  }'''


class MockSubscriberDBStub:
    def __init__(self):
        pass

    def ListSubscribers(self) -> List[str]:
        return ["IMSI001010000000001", "IMSI001010000000002"]

    def GetSubscriberData(self, _: NetworkID) -> SubscriberData:
        return SubscriberData(
            lte=LTESubscription(
                assigned_policies=["redirect"],
            )
        )


class SessionRpcServicerTest(unittest.TestCase):
    def setUp(self):
        self.servicer = SessionRpcServicer(self._get_config(),
                                           MockSubscriberDBStub())

    def _get_config(self) -> Dict[str, Any]:
        return {
            'redirect_rule_name': 'redirect',
            'whitelisted_ips': {
                'local': [80, 443],
            },
        }

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
        msg = CreateSessionRequest()
        msg.session_id = '1234'
        msg.imsi_plmn_id = '00101'
        resp = self.servicer.CreateSession(msg, None)

        # There should be a static rule installed for the redirection
        static_rules = self._rm_whitespace(str(resp.static_rules))
        expected = self._rm_whitespace(CSR_STATIC_RULES)
        self.assertEqual(static_rules, expected, 'There should be one static '
                         'rule installed for redirection.')

    def test_UpdateSession(self):
        """
        Update a session

        Assert:
            An arbitrarily large amount of credit should be granted.
        """
        msg = UpdateSessionRequest()
        credit_update = CreditUsageUpdate()
        credit_update.sid = 'abc'
        msg.updates.extend([credit_update])
        resp = self.servicer.UpdateSession(msg, None)
        session_response = self._rm_whitespace(str(resp))
        expected = self._rm_whitespace(USR)
        self.assertEqual(session_response, expected, 'There should be a large '
                         'amount of additional credit granted')

    def test_TerminateSession(self):
        """
        Terminate a session

        Assert:
            The session can be terminated successfully. Will always succeed.
        """
        msg = SessionTerminateRequest()
        msg.sid = 'abc'
        msg.session_id = 'session_id_123'
        resp = self.servicer.TerminateSession(msg, None)
        self.assertEqual(resp.sid, 'abc', 'SID should be same as request')
        self.assertEqual(resp.session_id, 'session_id_123', 'session ID should '
                         'be same as request')

    def _rm_whitespace(self, inp: str) -> str:
        return inp.replace(' ', '').replace('\n', '')
