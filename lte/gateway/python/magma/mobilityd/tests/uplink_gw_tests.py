"""
Copyright (c) 2020-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""


import logging
import unittest
from collections import defaultdict
import subprocess

from magma.mobilityd.uplink_gw import UplinkGatewayInfo

LOG = logging.getLogger('mobilityd.def_gw.test')
LOG.isEnabledFor(logging.DEBUG)


class DefGwTest(unittest.TestCase):
    """
    Validate default router setting.
    """

    def test_gw_ip_for_DHCP(self):
        gw_store = defaultdict(str)
        dhcp_gw_info = UplinkGatewayInfo(gw_store)
        self.assertEqual(dhcp_gw_info.getIP(), None)
        self.assertEqual(dhcp_gw_info.getMac(), None)

    def test_gw_ip_for_Ip_pool(self):
        gw_store = defaultdict(str)
        dhcp_gw_info = UplinkGatewayInfo(gw_store)
        dhcp_gw_info.read_default_gw()

        def_gw_cmd = "ip route show |grep default| awk '{print $3}'"
        p = subprocess.Popen([def_gw_cmd], stdout=subprocess.PIPE, shell=True)
        def_ip = p.stdout.read().decode("utf-8").strip()
        self.assertEqual(dhcp_gw_info.getIP(), str(def_ip))
        self.assertEqual(dhcp_gw_info.getMac(), None)

