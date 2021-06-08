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


class TestAttachDetachDuplicateNASRespMessages(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_duplicate_nas_resp_messages(self):
        """ Duplicate NAS Response Messages Test Case"""
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        maxNasMsgRetransmission = 4

        # Trigger Attach Request
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.ue_Id = req.ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        print("*** Triggering Attach Request ***")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        # Waiting for Authentication Request
        # Wait for last Timer T3460 Expiry
        for i in range(maxNasMsgRetransmission):
            print(
                "*** Waiting for Authentication Request Message (",
                str(i + 1),
                ") ***",
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
            )
            print(
                "*** Authentication Request Message Received (",
                str(i + 1),
                ") ***",
            )

        # Trigger Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqnRecvd = s1ap_types.ueSqnRcvd_t()
        sqnRecvd.pres = 0
        auth_res.sqnRcvd = sqnRecvd
        for i in range(maxNasMsgRetransmission):
            print(
                "*** Sending Authentication Response Message (",
                str(i + 1),
                ") ***",
            )
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
            )

        # Waiting for Security mode command
        # Wait for last Timer T3460 Expiry
        for i in range(maxNasMsgRetransmission):
            print(
                "*** Waiting for Security Mode Command Message (",
                str(i + 1),
                ") ***",
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
            )
            print(
                "*** Security Mode Command Message Received (",
                str(i + 1),
                ") ***",
            )

        # Trigger Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
        for i in range(maxNasMsgRetransmission):
            print(
                "*** Sending Security Mode Complete Message (",
                str(i + 1),
                ") ***",
            )
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
            )

        # Waiting for Attach accept
        # Wait for last Timer T3450 Expiry
        for i in range(maxNasMsgRetransmission):
            print(
                "*** Waiting for Attach Accept Message (", str(i + 1), ") ***",
            )
            # Attach accept will be sent in ICSR only for the first time
            # Re-transmitted Attach Accept will be sent in DL NAS transport
            if i < 1:
                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertEqual(
                    response.msg_type,
                    s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
                )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
            )
            print(
                "*** Attach Accept Message Received (", str(i + 1), ") ***",
            )

        # Trigger Attach Complete
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = req.ue_id
        for i in range(maxNasMsgRetransmission):
            print(
                "*** Sending Attach Complete Message (", str(i + 1), ") ***",
            )
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete,
            )

        print("*** Running UE detach ***")
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = req.ue_id
        detach_req.ueDetType = s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
        )


if __name__ == "__main__":
    unittest.main()
