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
from time import sleep

import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestAttachDetachMultipleIpBlocksMobilitydRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )
        self._ip_block = '192.168.125.0/24'

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_multiple_ip_blocks(self):
        """
        Attaching and detaching UE in mobilityd with multiple IP blocks across
        restart
        """
        self._s1ap_wrapper.configUEDevice(1)

        self._s1ap_wrapper.mobility_util.add_ip_block(
            self._ip_block,
        )

        old_blocks = self._s1ap_wrapper.mobility_util.list_ip_blocks()
        assert len(old_blocks) == 2, "2 IP blocks should be allocated on " \
                                     "mobilityd "

        print("************************* Restarting mobilityd")
        self._s1ap_wrapper.magmad_util.restart_services(["mobilityd"])
        for j in range(30):
            print("Waiting for", j, "seconds")
            sleep(1)

        curr_blocks = self._s1ap_wrapper.mobility_util.list_ip_blocks()
        # Check if old_blocks and curr_blocks contain same ip blocks after
        # restart
        self.assertListEqual(old_blocks, curr_blocks)

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for ",
            "UE id ", req.ue_id,
        )

        # Now actually attempt the attach
        self._s1ap_wrapper.s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Detach previously attached UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
        )


if __name__ == "__main__":
    unittest.main()
