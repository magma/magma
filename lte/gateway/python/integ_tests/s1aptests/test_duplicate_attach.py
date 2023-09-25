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


class TestDuplicateAttach(unittest.TestCase):
    """Integration Test: TestDuplicateAttach"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_duplicate_attach(self):
        """Start the attach procedure for an UE, proceed until attach accept is
        recvd. from the network. Without sending an attach accept, resend an
        attach request.
        """
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        print("************************* Running duplicate attach test")
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_AUTH_REQ_IND,
            s1ap_types.ueAttachComplete_t,
        )
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = False
        auth_res.sqnRcvd = sqn_recvd

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
            "************************ Dropping received attach accept for UE Id:",
            attach_acc.ue_Id,
        )

        while True:
            response = self._s1ap_wrapper.s1_util.get_response()
            if (
                response.msg_type
                == s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value
            ):
                attach_acc = response.cast(s1ap_types.ueAttachAccept_t)
                print(
                    "************************ Dropping received attach accept for UE Id:",
                    attach_acc.ue_Id,
                )
                continue

            assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value
            print(
                "************************ Received UE context release indication",
            )
            break

        print("************************ Running RE-Attach")
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        print("************************* Running UE detach")
        # Now detach the UE
        self._s1ap_wrapper._s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
        )


if __name__ == "__main__":
    unittest.main()
