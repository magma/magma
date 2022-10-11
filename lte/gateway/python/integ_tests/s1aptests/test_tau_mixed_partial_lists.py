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

import gpp_types
import s1ap_types
import s1ap_wrapper


class TestTauMixedPartialLists(unittest.TestCase):
    """Test TAU with TACs belonging to different partial TAI lists"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_tau_mixed_partial_lists(self):
        """Attach 3 UEs. Move the UEs to idle mode.
        Send tracking area update(TAU)request with TACs
        belonging to different partial TAI lists
        """
        num_ues = 3
        ue_ids = []
        tac = [1, 20, 45]
        wait_for_s1_context_rel = False
        self._s1ap_wrapper.configUEDevice(num_ues)
        # Attach
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            self._s1ap_wrapper.s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            ue_ids.append(req.ue_id)
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Add delay to ensure that S1APTester sends Attach Complete
            time.sleep(0.5)

        for i in range(num_ues):
            # Configure tacs in s1aptester
            config_tai = s1ap_types.nbConfigTai_t()
            config_tai.ue_Id = ue_ids[i]
            config_tai.tac = tac[i]
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.ENB_CONFIG_TAI, config_tai,
            )

            print(
                "************************* Sending ENB_CONFIG_TAI ",
                "for UE id ",
                ue_ids[i],
            )
            # Move the UE to idle state
            print(
                "************************* Sending UE context release request ",
                "for UE id ",
                ue_ids[i],
            )
            cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
            cntxt_rel_req.ue_Id = ue_ids[i]
            cntxt_rel_req.cause.causeVal = (
                gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
            )
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, cntxt_rel_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value

            print(
                "************************* Sending Tracking Area Update ",
                "request for UE id ",
                ue_ids[i],
            )
            tau_req = s1ap_types.ueTauReq_t()
            tau_req.ue_Id = ue_ids[i]
            tau_req.type = s1ap_types.Eps_Updt_Type.TFW_TA_UPDATING.value
            tau_req.Actv_flag = False
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

            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value

        print("************************* Sleeping for 2 seconds")
        time.sleep(2)
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
