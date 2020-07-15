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


class TestModifyMMEConfigForSanity(unittest.TestCase):
    def test_modify_mme_config_for_sanity(self):
        """
        Some test cases need changes in default mme configuration. This test
        script modifies MME configuration with generic values so that all the
        test cases of sanity suite can be verified
        """
        magmad_client = MagmadServiceGrpc()
        self._magmad_util = MagmadUtil(magmad_client)

        print(
            "Modifying MME configuration for all sanity test cases to pass"
        )
        self._magmad_util.update_mme_config_for_sanity(
            MagmadUtil.config_update_cmds.MODIFY
        )

        print("Restarting services to apply configuration change")
        self._magmad_util.restart_all_services()


if __name__ == "__main__":
    unittest.main()
