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


class TestAttachIcsDropWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_ics_drop_with_mme_restart(self):
        """Stateless Initial Context Setup Drop Test Case:
        1. Step-by-step UE attach procedure
        2. Set flag in S1APTester to drop ICS request message
        3. Restart MME after getting ICS request dropped indication
        4. Reset flag to not drop next ICS request messages
        5. Handle UE context release after no ICS response to MME
        6. Re-attach UE to verify if UE context was cleared properly
        after handling the NW Initiated detach once MME comes up
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
            "******************** Setting flag to drop Initial Context Setup "
            "Request",
        )
        drop_init_ctxt_setup_req = s1ap_types.UeDropInitCtxtSetup()
        drop_init_ctxt_setup_req.ue_Id = req.ue_id
        drop_init_ctxt_setup_req.flag = 1
        # Timer to release UE context at s1ap tester
        drop_init_ctxt_setup_req.tmrVal = 2000
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_DROP_ICS, drop_init_ctxt_setup_req,
        )

        print("******************** Sending Security Mode Complete")
        # Send Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )

        # enbApp sends UE_ICS_DROPD_IND message to tfwApp after dropping
        # ICS request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ICS_DROPD_IND.value,
        )
        print(
            "******************** Received Initial Context Setup Dropped "
            "Indication",
        )

        print("******************** Restarting MME service on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        print(
            "******************** Resetting flag to not drop next Initial "
            "Context Setup Request messages",
        )
        drop_init_ctxt_setup_req = s1ap_types.UeDropInitCtxtSetup()
        drop_init_ctxt_setup_req.ue_Id = req.ue_id
        drop_init_ctxt_setup_req.flag = 0
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_DROP_ICS, drop_init_ctxt_setup_req,
        )

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        # It has been observed that despite getting the restart command on
        # time, MME sometimes restarts after a delay of 5-6 seconds. If MME
        # restarts before ICS timer expiry of 4 seconds, it will send NW
        # initiated detach request with type re-attach required after coming
        # up, else MME will do implicit detach and send UE context release
        # command to release the S1AP context
        print(
            "******************** Waiting for NW initiated Detach Request or "
            "UE Context Release Indication",
        )
        resp_count = 0
        while True:
            resp_count += 1
            # Receive NW initiated detach request
            response = self._s1ap_wrapper.s1_util.get_response()
            if (
                response.msg_type
                == s1ap_types.tfwCmd.UE_NW_INIT_DETACH_REQUEST.value
            ):
                nw_init_detach_req = response.cast(
                    s1ap_types.ueNwInitdetachReq_t,
                )
                print(
                    "******************** Received NW initiated Detach Request"
                    " with detach type ",
                    nw_init_detach_req.Type,
                )
                self.assertEqual(
                    nw_init_detach_req.Type,
                    s1ap_types.ueNwInitDetType_t.TFW_RE_ATTACH_REQUIRED.value,
                )

                if resp_count == 1:
                    # Send detach accept
                    detach_accept = s1ap_types.ueTrigDetachAcceptInd_t()
                    detach_accept.ue_Id = req.ue_id
                    print(
                        "******************** Sending UE triggered Detach "
                        "Accept message",
                    )
                    self._s1ap_wrapper._s1_util.issue_cmd(
                        s1ap_types.tfwCmd.UE_TRIGGERED_DETACH_ACCEPT,
                        detach_accept,
                    )
                else:
                    print(
                        "******************** Ignoring re-transmitted (",
                        resp_count,
                        ") NW initiated detach request message",
                    )
            else:
                break

        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("******************** Received UE Context Release indication")

        print(
            "******************** Running End to End attach to verify if "
            "UE context was released properly after handling ICS drop for "
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
            "******************** Running UE detach (switch-off) for ",
            "UE id ",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            True,
        )


if __name__ == "__main__":
    unittest.main()
