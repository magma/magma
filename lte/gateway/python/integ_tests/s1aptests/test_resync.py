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


class TestResync(unittest.TestCase):
    """Integration Test: TestResync"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_resync(self):
        """Attach logic where we bump up the seq number stored in UE SQN
        array in order to trigger authentication failure from UE to
        verify the re-sync of seq number logic in EPC"""

        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        print("************************* Running sync failure test")

        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = req.ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        # Trigger Attach
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST,
            attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        max_sqn_recvd = s1ap_types.ueSqnRcvd_t()
        # Use SQN received from EPC as is
        sqn_recvd.pres = False
        # Set UE side max SQN value to bigger than value received from EPC
        # This is to trigger auth failure with cause = Sync Failure
        max_sqn_recvd.pres = True
        max_sqn_recvd.sqn = (0, 0, 0, 0, 0, 96)
        auth_res.sqnRcvd = sqn_recvd
        auth_res.maxSqnRcvd = max_sqn_recvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP,
            auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()

        # Verify that EPC re-sends Authentication Request
        assert response.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value
        print("************************* Retrying authentication")
        # Use SQN received from EPC as is
        sqn_recvd.pres = False
        auth_res.sqnRcvd = sqn_recvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP,
            auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value

        # Trigger Security mode complete message
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

        # Trigger Attach complete message
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE,
            attach_complete,
        )
        print("************************* UE Attach Complete")

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print("************************* Running UE detach")
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = req.ue_id
        detach_req.ueDetType = (
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST,
            detach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value


if __name__ == "__main__":
    unittest.main()
