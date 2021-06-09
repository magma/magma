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


class TestAttachEnbRlf(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()
        print(
            "Restart sctpd service to clear Redis state as test case doesn't"
            " intend to initiate detach procedure",
        )
        self._s1ap_wrapper.magmad_util.restart_sctpd()
        self._s1ap_wrapper.magmad_util.print_redis_state()

    def test_attach_enb_rlf(self):
        """
        Attach, after attach procedure is completed, eNB detects the RLF
        i.e., releasing the context by sending UE context release request
        with cause Radio Link Failure by that time we sending one more
        attach request
        """
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)

        # Trigger Attach Request
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.ue_Id = 1
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        print("***Triggering Attach Request ***")

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )

        # Trigger Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = 1
        sqnRecvd = s1ap_types.ueSqnRcvd_t()
        sqnRecvd.pres = 0
        auth_res.sqnRcvd = sqnRecvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )

        # Trigger Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = 1
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
        )

        # Trigger Attach Complete
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = 1
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete,
        )

        print("***Triggering Re-Attach Request ***")

        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.ue_Id = 1
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        # Send UE context release request with cause Radio Link Failure
        ue_ctxt_rel_req = s1ap_types.ueCntxtRelReq_t()
        ue_ctxt_rel_req.ue_Id = 1
        ue_ctxt_rel_req.cause.causeVal = (
            gpp_types.CauseRadioNetwork.RADIO_CONNECTION_WITH_UE_LOST.value
        )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, ue_ctxt_rel_req,
        )
        time.sleep(1)


if __name__ == "__main__":
    unittest.main()
