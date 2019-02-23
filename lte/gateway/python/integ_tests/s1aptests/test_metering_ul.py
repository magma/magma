"""
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
"""
import unittest

from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows
from integ_tests.s1aptests.workflow.attach_detach import attach_ue, detach_ue
from integ_tests.s1aptests.workflow.data import run_tcp_uplink


class TestMeteringUl(unittest.TestCase):
    METERING_TABLE = 3

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_ul_single_ue(self):
        """
        Send some uplink data from a UE and check that packets are captured by
        an appropriate metering flow.
        """
        self._s1ap_wrapper.configUEDevice(1)
        ue = self._s1ap_wrapper.ue_req

        attach_ue(ue, self._s1ap_wrapper)
        run_tcp_uplink(ue, self._s1ap_wrapper)
        self._verify_ue_metering_flows(ue)
        detach_ue(ue, self._s1ap_wrapper)

    def _verify_ue_metering_flows(self, ue):
        datapath = get_datapath()
        flows = get_flows(datapath, {'table_id': self.METERING_TABLE,
                                     'priority': 10})
        self.assertEqual(2, len(flows), 'There should be 2 UE metering flows, '
                                        'found {}'.format(len(flows)))

        ue_ip = str(self._s1ap_wrapper._s1_util.get_ip(ue.ue_id))
        self._verify_downlink_flow(flows, ue_ip)
        self._verify_uplink_flow(flows, ue_ip)

    def _verify_downlink_flow(self, flows, ue_ip):
        dl_flows_filtered = list(filter(
            lambda flow: 'ipv4_dst' in flow['match'],
            flows,
        ))
        self.assertEqual(1, len(dl_flows_filtered),
                         'Expected 1 UE downlink metering flow, '
                         'found {}'.format(len(dl_flows_filtered)))

        dl_flow = dl_flows_filtered[0]
        self.assertEqual(0, dl_flow['packet_count'])
        self.assertEqual(ue_ip, dl_flow['match']['ipv4_dst'])

    def _verify_uplink_flow(self, flows, ue_ip):
        ul_flows_filtered = list(filter(
            lambda flow: 'ipv4_src' in flow['match'],
            flows,
        ))
        self.assertEqual(1, len(ul_flows_filtered),
                         'Expected 1 UE uplink metering flow, '
                         'found {}'.format(len(ul_flows_filtered)))
        ul_flow = ul_flows_filtered[0]
        self.assertGreater(ul_flow['packet_count'], 0)
        self.assertEqual(ue_ip, ul_flow['match']['ipv4_src'])


if __name__ == '__main__':
    unittest.main()
