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


class TestAttachUeCtxtReleaseCmpDelay(unittest.TestCase):
    """Integration Test: TestAttachUeCtxtReleaseCmpDelay"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_attach_ue_ctxt_release_cmp_delay(self):
        """Attach, Delay Ue context release complete"""
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        # Trigger Attach Request
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.ue_Id = req.ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        print("********Triggering Attach Request ")

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST,
            attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value

        # Trigger Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP,
            auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value

        # Trigger Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE,
            sec_mode_complete,
        )

        # Receive initial context setup and attach accept indication
        response = (
            self._s1ap_wrapper._s1_util
                .receive_initial_ctxt_setup_and_attach_accept()
        )
        attach_acc = response.cast(s1ap_types.ueAttachAccept_t)
        print(
            "********************** Received attach accept for UE Id:",
            attach_acc.ue_Id,
        )

        delay_ue_ctxt_rel_cmp = s1ap_types.UeDelayUeCtxtRelCmp()
        delay_ue_ctxt_rel_cmp.ue_Id = req.ue_id
        delay_ue_ctxt_rel_cmp.flag = 1
        delay_ue_ctxt_rel_cmp.tmrVal = 1000

        print("*** Setting Delay for Ue context release complete ***")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_DELAY_UE_CTXT_REL_CMP,
            delay_ue_ctxt_rel_cmp,
        )

        # Trigger Attach Complete
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE,
            attach_complete,
        )
        time.sleep(0.5)
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_EMM_INFORMATION.value

        time.sleep(0.5)
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = req.ue_id
        detach_req.ueDetType = s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST,
            detach_req,
        )
        response = self._s1ap_wrapper._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_DETACH_ACCEPT_IND.value

        print(
            "*** Sending UE context release request ",
            "for UE id ***",
            req.ue_id,
        )

        # Send UE context release request to move UE to idle mode
        uectxtrel_req = s1ap_types.ueCntxtRelReq_t()
        uectxtrel_req.ue_Id = req.ue_id
        uectxtrel_req.cause.causeVal = (
            gpp_types.CauseRadioNetwork.RELEASE_DUE_TO_EUTRAN_GENERATED_REASON.value
        )

        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST,
            uectxtrel_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value

        time.sleep(10)


if __name__ == "__main__":
    unittest.main()
