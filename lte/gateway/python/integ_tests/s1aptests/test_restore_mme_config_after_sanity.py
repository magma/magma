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


class TestRestoreMMEConfigAfterSanity(unittest.TestCase):
    def test_restore_mme_config_after_sanity(self):
        """
        This test script restores the MME configuration to default values, if
        the config file mme.conf.template has been modified using the test
        script s1aptests/test_modify_mme_config_for_sanity.py
        """
        magmad_client = MagmadServiceGrpc()
        self._magmad_util = MagmadUtil(magmad_client)

        # Replace mme.conf.template with backup of this config file. Backup of
        # this config file with default values is created when running the
        # test script s1aptests/test_modify_mme_config_for_sanity.py
        print(
            "Restoring MME configuration to default values using backup configuration file"
        )
        self._magmad_util.update_mme_config_for_sanity(
            MagmadUtil.config_update_cmds.RESTORE
        )


if __name__ == "__main__":
    unittest.main()
