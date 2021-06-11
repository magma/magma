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

import time
import unittest

import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestAttachDetachWithCorruptStatelessMME(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
            health_service=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_with_corrupt_stateless_mme(self):
        """
        Basic attach/detach test with two UEs,
        with purpose of validating corruption of MME state and recovery
        """
        detach_type = s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value
        wait_for_s1 = True
        self._s1ap_wrapper.configUEDevice(1)

        services_state_dict = {
            'mme': 'mme_nas_state',
        }

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for ",
            "UE id ", req.ue_id,
        )

        for s in services_state_dict:
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            print("************************* Corrupting %s state" % s)
            self._s1ap_wrapper.magmad_util.corrupt_agw_state(
                services_state_dict[s],
            )

            print("************************* Restarting %s service" % s)
            self._s1ap_wrapper.magmad_util.restart_services([s])

            for j in range(100):
                print("Waiting for", j, "seconds")
                time.sleep(1)

            # Re-establish S1 connection between eNB and MME
            self._s1ap_wrapper._s1setup()

            print(
                "************************* Re-running End to End attach for ",
                "UE id ",
                req.ue_id,
            )

            # Repeat the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Now detach the UE
            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type, wait_for_s1,
            )


if __name__ == "__main__":
    unittest.main()
