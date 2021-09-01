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
            "Modifying MME configuration for all sanity test cases to pass",
        )
        self._magmad_util.update_mme_config_for_sanity(
            MagmadUtil.config_update_cmds.MODIFY,
        )

        print("Restarting services to apply configuration change")
        self._magmad_util.restart_all_services()


if __name__ == "__main__":
    unittest.main()
