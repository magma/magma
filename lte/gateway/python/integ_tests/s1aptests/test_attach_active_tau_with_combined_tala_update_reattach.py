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
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_active_tau_with_combined_tala_update_reattach(self):
        """This test case validates reattach after active combined TAU reject:
        1. End-to-end attach with attach type COMBINED_EPS_IMSI_ATTACH
        2. Send active TAU request (Combined TALA update)
        3. Receive TAU reject (Combined TALA update not supported)
        4. Retry end-to-end combined EPS IMSI attach to verify if UE context
           was released properly after combined TAU reject
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
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
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

        # Waiting for TAU Reject Indication -Combined TALA update not supported
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_TAU_REJECT_IND.value,
        )
        print(
            "************************* Received Tracking Area Update Reject "
            "Indication",
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print(
            "************************* Received UE context release indication",
        )

        print(
            "************************* Running End to End attach to verify if "
            "UE context was released properly after combined TAU reject for "
            "UE id ",
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
            ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
