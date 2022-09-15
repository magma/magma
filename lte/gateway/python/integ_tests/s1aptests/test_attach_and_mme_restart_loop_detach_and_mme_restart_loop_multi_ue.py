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


class TestAttachAndMmeRestartLoopDetachAndMmeRestartLoopMultiUe(
    unittest.TestCase,
):
    """
    Integration Test: TestAttachAndMmeRestartLoopDetachAndMmeRestartLoopMultiUe
    """

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_attach_and_mme_restart_loop_detach_and_mme_restart_loop_multi_ue(
        self,
    ):
        """
        Multi UE attach-detach with MME restart. Steps to be followed:
        1-Attach
        2-MmeRestart + wait for 30 seconds
        3-Repeat step 1 and 2 in loop for all UEs
        4-Detach
        5-MmeRestart + wait for 30 seconds
        6-Repeat step 4 and 5 in loop for all UEs
        """
        num_ues = 32
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        detach_type_str = ["NORMAL", "SWITCHOFF"]
        self._s1ap_wrapper.configUEDevice(num_ues)

        # The inactivity timers for UEs attached in the beginning starts getting
        # expired before all the UEs could be attached. Increasing UE inactivity
        # timer to 15 min (900000 ms) to allow all the UEs to get attached and
        # detached properly
        config_data = s1ap_types.FwNbConfigReq_t()
        config_data.inactvTmrVal_pr.pres = True
        config_data.inactvTmrVal_pr.val = 900000
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.ENB_INACTV_TMR_CFG,
            config_data,
        )
        time.sleep(0.5)

        ue_ids = []
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for "
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

            print(
                "************************* Restarting MME service on gateway",
            )
            wait_for_restart = 30
            self._s1ap_wrapper.magmad_util.restart_services(
                ["mme"], wait_for_restart,
            )

        for ue in ue_ids:
            # Now detach the UE
            random.seed(time.time())
            index = random.randint(0, 1)
            print(
                "************************* Running UE detach for UE id ",
                ue,
                "(Detach type: " + detach_type_str[index] + ")",
            )
            self._s1ap_wrapper.s1_util.detach(ue, detach_type[index])

            print(
                "************************* Restarting MME service on gateway",
            )
            wait_for_restart = 30
            self._s1ap_wrapper.magmad_util.restart_services(
                ["mme"], wait_for_restart,
            )


if __name__ == "__main__":
    unittest.main()
