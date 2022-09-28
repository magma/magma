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


class TestAttachActiveTauWithCombinedTalaUpdateReattach(unittest.TestCase):
    """Test combined TAU with TA/LA updating and active flag set to true"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_active_tau_with_combined_tala_update_reattach(self):
        """This test case validates reattach after active combined TAU with TA/LA updating:
        1. End-to-end attach with attach type COMBINED_EPS_IMSI_ATTACH
        2. Send active TAU request (Combined TALA update)
        3. Receive TAU accept
        4. Perform end-to-end combined EPS IMSI attach
        """

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id

        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
            eps_type=s1ap_types.TFW_EPS_ATTACH_TYPE_COMB_EPS_IMSI_ATTACH,
        )

        # Wait for EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        time.sleep(0.5)

        print(
            "************************* Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        req = s1ap_types.ueCntxtRelReq_t()
        req.ue_Id = ue_id
        req.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST,
            req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        print(
            "************************* Received UE context release indication",
        )

        print(
            "************************* Sending active TAU request (Combined "
            "TALA update) for UE id ",
            ue_id,
        )
        # Send active TAU request with combined TALA update as update type
        req = s1ap_types.ueTauReq_t()
        req.ue_Id = ue_id
        req.type = s1ap_types.Eps_Updt_Type.TFW_COMB_TALA_UPDATING.value
        req.Actv_flag = True
        req.ueMtmsi.pres = False
        self._s1ap_wrapper.s1_util.issue_cmd(s1ap_types.tfwCmd.UE_TAU_REQ, req)

        response = (
            self._s1ap_wrapper._s1_util.receive_initial_ctxt_setup_and_tau_accept()
        )
        tau_acc = response.cast(s1ap_types.ueTauAccept_t)
        print(
            "************************* Received Tracking Area Update",
            "accept for UE Id:",
            tau_acc.ue_Id,
        )

        print(
            "************************* Running End to End attach for UE Id:",
            ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
            eps_type=s1ap_types.TFW_EPS_ATTACH_TYPE_COMB_EPS_IMSI_ATTACH,
        )

        # Wait for EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print("************************* Running UE detach for UE id", ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            wait_for_s1_ctxt_release=True,
        )


if __name__ == "__main__":
    unittest.main()
