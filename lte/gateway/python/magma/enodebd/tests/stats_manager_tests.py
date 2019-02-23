"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import pkg_resources
from unittest import TestCase
from xml.etree import ElementTree

from magma.enodebd import metrics
from magma.enodebd.stats_manager import StatsManager


class StatsManagerTest(TestCase):
    """
    Tests for eNodeB statistics manager
    """

    def test_parse_stats(self):
        """ Test that example statistics from eNodeB can be parsed, and metrics
            updated """
        # Example performance metrics structure, sent by eNodeB
        pm_file_example = pkg_resources.resource_string(__name__,
                                                        'pm_file_example.xml')

        root = ElementTree.fromstring(pm_file_example)

        mgr = StatsManager()
        mgr.parse_pm_xml(root)

        # Check that metrics were correctly populated
        # See '<V i="5">123</V>' in pm_file_example
        rrc_estab_attempts = metrics.STAT_RRC_ESTAB_ATT.collect()
        self.assertEqual(rrc_estab_attempts[0].samples[0][2], 123)
        # See '<V i="7">99</V>' in pm_file_example
        rrc_estab_successes = metrics.STAT_RRC_ESTAB_SUCC.collect()
        self.assertEqual(rrc_estab_successes[0].samples[0][2], 99)
        # See '<SV>654</SV>' in pm_file_example
        rrc_reestab_att_reconf_fail = \
            metrics.STAT_RRC_REESTAB_ATT_RECONF_FAIL.collect()
        self.assertEqual(rrc_reestab_att_reconf_fail[0].samples[0][2], 654)
        # See '<SV>65537</SV>' in pm_file_example
        erab_rel_req_radio_conn_lost = \
            metrics.STAT_ERAB_REL_REQ_RADIO_CONN_LOST.collect()
        self.assertEqual(erab_rel_req_radio_conn_lost[0].samples[0][2], 65537)
