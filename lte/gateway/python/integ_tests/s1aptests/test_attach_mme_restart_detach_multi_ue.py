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

import random
import time
import unittest

import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestAttachMmeRestartDetachMultiUe(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_mme_restart_detach_multi_ue(self):
        """
        Multi UE attach-detach with MME restart. Steps to be followed:
        1-Attach
        2-Repeat step 1 in loop for all UEs
        3-MmeRestart + wait for 30 seconds
        4-Detach
        5-Repeat step 4 in loop for all UEs
        """
        num_ues = 32
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        detach_type_str = ["NORMAL", "SWITCHOFF"]
        self._s1ap_wrapper.configUEDevice(num_ues)

        ue_ids = []
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        print("************************* Restarting MME service on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        for ue in ue_ids:
            # Now detach the UE
            random.seed(time.clock())
            index = random.randint(0, 1)
            print(
                "************************* Running UE detach for UE id ",
                ue,
                "(Detach type: " + detach_type_str[index] + ")",
            )
            self._s1ap_wrapper.s1_util.detach(ue, detach_type[index])


if __name__ == "__main__":
    unittest.main()
