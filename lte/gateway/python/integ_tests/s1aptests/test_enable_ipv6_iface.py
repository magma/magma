"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest

from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestEnableIpv6Iface(unittest.TestCase):
    """Unittest: TestEnableIpv6Iface"""

    def test_enable_ipv6_iface(self):
        """Add eth3 interface as nat_iface on pipelined.yml to
        test ipv6 data testcases
        """
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

        print("Enabling ipv6_iface on pipelined.yml")
        cmd = MagmadUtil.config_ipv6_iface_cmds.ENABLE
        self._s1ap_wrapper.enable_disable_ipv6_iface(cmd)


if __name__ == "__main__":
    unittest.main()
