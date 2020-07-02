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


class TestResetMMEConfigAfterSanity(unittest.TestCase):
    def test_reset_mme_config_after_sanity(self):
        """
        This test script resets the MME configuration to default values, if
        the config file mme.conf.template has been updated using the test
        script s1aptests/test_update_mme_config_for_sanity.py
        """
        magmad_client = MagmadServiceGrpc()
        self._magmad_util = MagmadUtil(magmad_client)

        # Replace mme.conf.template with backup of this config file. Backup of
        # this config file with default values is created when running the
        # test script s1aptests/test_update_mme_config_for_sanity.py
        self._magmad_util.reset_mme_config_file()


if __name__ == "__main__":
    unittest.main()
