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


class TestNoAttachCompleteWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )
        self.gateway_services = self._s1ap_wrapper.get_gateway_services_util()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_no_attach_complete_with_mme_restart(self):
        """
        Step 1: UE sends Attach Request, receives Attach Accept and UE shall
                not respond to mme
        Step 2: After sending Attach Accept, mme runs 3450 timer, while timer
                is running, mme restarts.
        Step 3: On mme recovery, mme sends Detach Request with re-attach
                required, S1ap shall send Detach Accept and release the UE
                contexts
        Step 4: Once again attach UE to network
        """

        self._s1ap_wrapper.configIpBlock()
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running attach setup timer expiry test",
        )

        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = req.ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
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

        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
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

        print(
            "************************* Restarting MME service on", "gateway",
        )
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(20):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        # Receive NW initiated detach request
        response = self._s1ap_wrapper.s1_util.get_response()

        while response.msg_type == s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND:
            print(
                "Received Attach Accept retransmission from before restart",
                "Ignoring...",
            )
            response = self._s1ap_wrapper.s1_util.get_response()

        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.UE_NW_INIT_DETACH_REQUEST.value,
        )
        nw_init_detach_req = response.cast(s1ap_types.ueNwInitdetachReq_t)
        print(
            "**************** Received NW initiated Detach Req with detach "
            "type set to ",
            nw_init_detach_req.Type,
        )
        self.assertEqual(
            nw_init_detach_req.Type,
            s1ap_types.ueNwInitDetType_t.TFW_RE_ATTACH_REQUIRED.value,
        )
        # Receive NW initiated detach request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.UE_NW_INIT_DETACH_REQUEST.value,
        )
        nw_init_detach_req = response.cast(s1ap_types.ueNwInitdetachReq_t)
        print(
            "**************** Received NW initiated Detach Req with detach "
            "type set to ",
            nw_init_detach_req.Type,
        )
        self.assertEqual(
            nw_init_detach_req.Type,
            s1ap_types.ueNwInitDetType_t.TFW_RE_ATTACH_REQUIRED.value,
        )

        print("**************** Sending Detach Accept")
        # Send detach accept
        detach_accept = s1ap_types.ueTrigDetachAcceptInd_t()
        detach_accept.ue_Id = req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_TRIGGERED_DETACH_ACCEPT, detach_accept,
        )

        # Wait for UE context release command
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("****** Received Ue context release command *********")

        print("****** Triggering end-end attach after mme restart *********")
        self._s1ap_wrapper.s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            False,
        )


if __name__ == "__main__":
    unittest.main()
