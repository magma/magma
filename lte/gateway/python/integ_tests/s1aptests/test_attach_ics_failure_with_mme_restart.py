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
from s1ap_utils import MagmadUtil


class TestAttachIcsFailureWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_ics_failure_with_mme_restart(self):
        """Stateless Initial Context Setup Failure Test Case:
        1. Step-by-step UE attach procedure
        2. Set flag in S1APTester to send ICS failure to MME
        3. Restart MME after sending ICS failure to MME
        4. Reset flag to not send ICS failure for next ICS request
        5. Handle UE context release after ICS failure
        6. Re-attach UE to verify if UE context was cleared properly
        after handling the ICS failure message
        """

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        # Send Attach Request
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        pdn_type.pdn_type = 3
        attach_req.ue_Id = req.ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type

        print("******************** Sending Attach Request")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("******************** Received Authentiction Request Indication")

        # Send Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqnRecvd = s1ap_types.ueSqnRcvd_t()
        sqnRecvd.pres = 0
        auth_res.sqnRcvd = sqnRecvd

        print("******************** Sending Authentiction Response")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        print("******************** Received Security Mode Command Indication")

        print(
            "******************** Setting flag to send Initial Context Setup "
            "Failure",
        )
        init_ctxt_setup_fail = s1ap_types.ueInitCtxtSetupFail()
        init_ctxt_setup_fail.ue_Id = req.ue_id
        init_ctxt_setup_fail.flag = 1
        init_ctxt_setup_fail.causeType = 0
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_INIT_CTXT_SETUP_FAIL, init_ctxt_setup_fail,
        )

        print("******************** Sending Security Mode Complete")
        # Send Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )
        print("******************** Received Initial Context Setup Indication")

        print("******************** Restarting MME service on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        init_ctxt_setup_fail = s1ap_types.ueInitCtxtSetupFail()
        init_ctxt_setup_fail.ue_Id = req.ue_id
        init_ctxt_setup_fail.flag = 0
        init_ctxt_setup_fail.causeType = 0
        print(
            "******************** Resetting flag to not send Initial Context "
            "Setup Failure",
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_INIT_CTXT_SETUP_FAIL, init_ctxt_setup_fail,
        )

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        # Waiting for UE Context Release indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("******************** Received UE Context Release indication")

        print(
            "******************** Running End to End attach to verify if "
            "UE context was released properly after handling ICS failure for "
            "UE id ",
            req.ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait for EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print(
            "******************** Running UE detach for UE id",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            True,
        )


if __name__ == "__main__":
    unittest.main()
