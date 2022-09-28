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

import s1ap_types
import s1ap_wrapper


class TestAttachNoInitialContextResp(unittest.TestCase):
    """Integration Test: TestAttachNoInitialContextResp"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_attach_no_initial_context_resp(self):
        """Attach, no initial context setup response"""
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print("***** Running attach - initial context setup timer expiry test")
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = req.ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST,
            attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqnrecvd = s1ap_types.ueSqnRcvd_t()
        sqnrecvd.pres = 0
        auth_res.sqnRcvd = sqnrecvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP,
            auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()

        assert response.msg_type == s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value

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

        # Send SCTP ABORT to MME
        sctp_abort = s1ap_types.FwSctpAbortReq_t()
        sctp_abort.cause = 0
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.SCTP_ABORT_REQ,
            sctp_abort,
        )
        print(
            "************************* Send Initial context Setup response ",
            "for UE id ",
            ue_id,
        )
        print("***************** Attach Aborted and UE Context released")


if __name__ == "__main__":
    unittest.main()
