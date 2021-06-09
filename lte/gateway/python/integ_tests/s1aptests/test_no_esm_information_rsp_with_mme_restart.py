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
from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestNoEsmInformationRspWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_no_esm_information_rsp_with_mme_restart(self):
        """ Restart MME after receiving ESM information req during attach +
            after MME recovers perform end-end attach + detach """
        num_ues = 1

        self._s1ap_wrapper.configUEDevice(num_ues)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        # Send Attach
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        # enabling ESM Information transfer flag
        attach_req.eti.pres = 1
        attach_req.eti.esm_info_transfer_flag = 1

        print("Sending Attach Request ue-id", ue_id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("Received Auth Req for ue-id", ue_id)

        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = ue_id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        print("Sending Auth Response for ue-id", ue_id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        print("Received Security Mode Command for ue-id", ue_id)

        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )
        print("Sent Security Mode Complete for ue-id", ue_id)

        # Esm Information Request indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ESM_INFORMATION_REQ.value,
        )
        print("Received Esm Information Request ue-id", ue_id)

        print("************************* Restarting MME service on", "gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        # Receive UE_ESM_INFORMATION_REQ, as sometimes MME retransmits
        # UE_ESM_INFORMATION_REQ message before it restarts
        # as T3489 timer is only 4 secs.
        response = self._s1ap_wrapper.s1_util.get_response()
        if response.msg_type == s1ap_types.tfwCmd.UE_ESM_INFORMATION_REQ.value:
            print(
                "Received retransmitted UE_ESM_INFORMATION_REQ for ue-id",
                ue_id,
            )
            # Wait for UE_CTX_REL_IND
            global response
            response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )

        # Perform end-end attach
        self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print(
            "******************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
        )


if __name__ == "__main__":
    unittest.main()
