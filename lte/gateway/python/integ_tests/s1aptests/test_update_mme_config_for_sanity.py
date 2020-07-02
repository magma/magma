"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

from integ_tests.common.magmad_client import MagmadServiceGrpc
from integ_tests.s1aptests.s1ap_utils import MagmadUtil


class TestUpdateMMEConfigForSanity(unittest.TestCase):
    def test_update_mme_config_for_sanity(self):
        """
        Some test cases need changes in default mme configuration. This test
        script modifies MME configuration with generic values so that all the
        test cases of sanity suite can be verified
        """
        magmad_client = MagmadServiceGrpc()
        self._magmad_util = MagmadUtil(magmad_client)

        # Create a backup of default MME configuration before modifying it,
        # which can be used later to restore the original configuration
        self._magmad_util.create_mme_config_backup()

        # Clear the default PLMN & TAC configuration and update the new values
        self._magmad_util.clear_default_plmn_tac_configuration()
        self._magmad_util.configure_multiple_plmn_tac()

        # Reduce the mobile reachability timer to 1 minute, so that it can
        # quickly be tested as part of Sanity. The current default value of
        # Mobile Reachability Timer is 54 minutes
        # Commented for now because sanity is failing for 7 test cases with
        # this change
        # self._magmad_util.reduce_mobile_reachability_timer_value()

        # Restart the services in order to run MME with updated configuration
        self._magmad_util.restart_all_services()


if __name__ == "__main__":
    unittest.main()
