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

import s1ap_types
import s1ap_wrapper


class TestAttachCompleteAfterIcsTimerExpiry(unittest.TestCase):
    """Send attach complete message after ICS timer expiry"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_complete_after_ics_timer_expiry(self):
        """Simulating sending of attach complete msg after ICS timer expiry.

        MME should discard attach complete message as it is already in the
        process of implicit detach

        ### Do not add this TC to sanity suite ###
        """
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        # Trigger Attach Request
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        pdn_type.pdn_type = 1
        attach_req.ue_Id = req.ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        print("****** Sent ATTACH_REQUEST")
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("****** Received AUTH_REQ_IND")
        delay_init_ctxt_setup_resp = s1ap_types.UeDelayInitCtxtSetupRsp()
        delay_init_ctxt_setup_resp.ue_Id = req.ue_id
        delay_init_ctxt_setup_resp.flag = 1
        delay_init_ctxt_setup_resp.tmrVal = 4500
        print("*** Setting Initial Context Setup Resp Delay ***")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_DELAY_ICS_RSP, delay_init_ctxt_setup_resp,
        )

        # Trigger Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )
        print("****** Sent UE_AUTH_RESP")
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        print("****** Received UE_SEC_MOD_CMD_IND")

        # Trigger Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )
        print("****** Sent UE_SEC_MOD_COMPLETE")

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )
        print("****** Received INT_CTX_SETUP_IND")

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
        )
        print("****** Received UE_ATTACH_ACCEPT_IND")

        print("****** Sleeping for 4 seconds. Waiting for ICS timer to expire")
        time.sleep(4)
        # Trigger Attach Complete
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete,
        )
        print("****** Sent ATTACH_COMPLETE")

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("****** Received UE_CTX_REL_IND")


if __name__ == "__main__":
    unittest.main()
