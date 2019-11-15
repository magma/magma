"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
from typing import Any, Dict
from lte.protos.session_manager_pb2 import CreateSessionRequest, \
    UpdateSessionRequest, CreditUsageUpdate, SessionTerminateRequest
from magma.captive_portal.rpc_servicer import SessionRpcServicer


CSR_CREDITS = '''
  [success: true
  sid: "00101"
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
  result_code: 1]'''

CSR_STATIC_RULES = '[rule_id: "redirect"]'

CSR_DYNAMIC_POLICY_1 = '''
  id: "whitelist_policy_id-192.168.128.1:80"
  priority: 100
  flow_list {
    match {
      ipv4_dst: "192.168.128.1"
      tcp_dst: 80
      ip_proto: IPPROTO_TCP
    }
  }
  flow_list {
    match {
      ipv4_src: "192.168.128.1"
      tcp_src: 80
      ip_proto: IPPROTO_TCP
      direction: DOWNLINK
    }
  }
  qos {
    max_req_bw_ul: 2147483648
    max_req_bw_dl: 2147483648
    gbr_ul: 1048576
    gbr_dl: 1048576
    qci: QCI_3
    arp {
      priority_level: 1
      pre_capability: PRE_CAP_DISABLED
      pre_vulnerability: PRE_VUL_DISABLED
    }
  }
  tracking_type: NO_TRACKING'''

CSR_DYNAMIC_POLICY_2 = '''
  id: "whitelist_policy_id-192.168.128.1:443"
  priority: 100
  flow_list {
    match {
      ipv4_dst: "192.168.128.1"
      tcp_dst: 443
      ip_proto: IPPROTO_TCP
    }
  }
  flow_list {
    match {
      ipv4_src: "192.168.128.1"
      tcp_src: 443
      ip_proto: IPPROTO_TCP
      direction: DOWNLINK
    }
  }
  qos {
    max_req_bw_ul: 2147483648
    max_req_bw_dl: 2147483648
    gbr_ul: 1048576
    gbr_dl: 1048576
    qci: QCI_3
    arp {
      priority_level: 1
      pre_capability: PRE_CAP_DISABLED
      pre_vulnerability: PRE_VUL_DISABLED
    }
  }
  tracking_type: NO_TRACKING'''

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


class SessionRpcServicerTest(unittest.TestCase):
    def setUp(self):
        self.servicer = SessionRpcServicer(self._get_config())

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

        # Check the granted credits
        credits = self._rm_whitespace(str(resp.credits))
        expected = self._rm_whitespace(CSR_CREDITS)
        self.assertEqual(credits, expected, 'Expected that there is a high '
                         'volume of granted credits which will be permissive '
                         'for access to captive portal.')

        # There should be a static rule installed for the redirection
        static_rules = self._rm_whitespace(str(resp.static_rules))
        expected = self._rm_whitespace(CSR_STATIC_RULES)
        self.assertEqual(static_rules, expected, 'There should be one static '
                         'rule installed for redirection.')

        # There should be two dynamic rules:
        #   one for uplink traffic to the captive portal server from the UE
        #   one for downlink traffic from the captive portal server to the UE
        p1 = self._rm_whitespace(str(resp.dynamic_rules[0].policy_rule))
        expected = self._rm_whitespace(CSR_DYNAMIC_POLICY_1)
        self.assertEqual(p1, expected, 'There should be a dynamic rule '
                         'installed for UE traffic to the captive portal '
                         'server')

        p2 = self._rm_whitespace(str(resp.dynamic_rules[1].policy_rule))
        expected = self._rm_whitespace(CSR_DYNAMIC_POLICY_2)
        self.assertEqual(p2, expected, 'There should be a dynamic rule '
                         'installed for captive portal traffic to the UE')

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
