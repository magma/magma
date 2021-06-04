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

import s1ap_types
import s1ap_wrapper


class TestAttachWithoutIpsAvailable(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._blocks = []

    def tearDown(self):
        for block in self._blocks:
            self._s1ap_wrapper.mobility_util.add_ip_block(block)
        self._s1ap_wrapper.cleanup()

    def test_attach_without_ips_available(self):
        """ Attaching without available IPs in mobilityd """
        self._s1ap_wrapper.configUEDevice(1)

        # Clear blocks
        self._s1ap_wrapper.mobility_util.cleanup()

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for ",
            "UE id ", req.ue_id,
        )
        # Now actually attempt the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND, s1ap_types.ueAttachFail_t,
        )


if __name__ == "__main__":
    unittest.main()
