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


class TestTauTaUpdatingConnectedMode(unittest.TestCase):
    """Test TAU with TA updating in connected mode"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_tau_ta_updating_connected_mode(self):
        """Attach 2 UEs.
        1. For the 1st UE, send tracking area update(TAU) request
        in UL NAS transport with EPS Update Type=TA Updating and
        active flag set to false
        2. For the 2nd UE, send tracking area update(TAU) request
        in UL NAS transport with EPS Update Type=TA Updating and
        active flag set to true
        """
        num_ues = 2
        wait_for_s1_context_rel = False
        self._s1ap_wrapper.configUEDevice(num_ues)
        ue_ids = []
        # Attach
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id
            self._s1ap_wrapper.s1_util.attach(
                ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
                id_type=s1ap_types.TFW_MID_TYPE_GUTI,
            )
            ue_ids.append(req.ue_id)
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

        # Add delay to ensure that S1APTester sends Attach Complete
        time.sleep(0.5)

        active_flag = [False, True]
        for i in range(num_ues):
            print(
                "************************* Sending Tracking Area Update ",
                "request for UE id ",
                ue_ids[i],
            )
            tau_req = s1ap_types.ueTauReq_t()
            tau_req.ue_Id = ue_ids[i]
            tau_req.type = s1ap_types.Eps_Updt_Type.TFW_TA_UPDATING.value
            tau_req.Actv_flag = active_flag[i]
            tau_req.ueMtmsi.pres = False
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_TAU_REQ, tau_req,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_TAU_ACCEPT_IND.value
            tau_acc = response.cast(s1ap_types.ueTauAccept_t)
            print(
                "************************* Received Tracking Area Update ",
                "accept for UE id ",
                tau_acc.ue_Id,
            )

        # Detach
        for ue in ue_ids:
            print(
                "************************* Running UE detach (switch-off) for ",
                "UE id ",
                ue,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                ue,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                wait_for_s1_context_rel,
            )


if __name__ == "__main__":
    unittest.main()
