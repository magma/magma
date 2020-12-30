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

import threading
import unittest
import time

import gpp_types
import s1ap_types
import s1ap_wrapper


class TestAttachCombinedTauPeriodicActiveReattach(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_combined_tau_periodic_active_reattach(self):
        """Combined attach UE, send active period tracking area updates (TAU)
        and combined reattach after tracking area reject"""

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

        # Wait on EMM Information from MME
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
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        )
        print(
            "************************* Received UE context release indication"
        )

        print(
            "************************* Sending Tracking Area Update ",
            "request for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        req = s1ap_types.ueTauReq_t()
        req.ue_Id = ue_id
        req.type = s1ap_types.Eps_Updt_Type.TFW_COMB_TALA_UPDATING_IMSI.value
        req.Actv_flag = True
        req.ueMtmsi.pres = False
        self._s1ap_wrapper.s1_util.issue_cmd(s1ap_types.tfwCmd.UE_TAU_REQ, req)

        print(
            "************************* Waiting for Tracking Area Update Reject Indication"
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_TAU_REJECT_IND.value
        )
        print(
            "************************* Received Tracking Area Update Reject Indication"
        )

        print(
            "************************* Waiting for UE context release indication"
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        )
        print(
            "************************* Received UE context release indication"
        )

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

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, True
        )


if __name__ == "__main__":
    unittest.main()
