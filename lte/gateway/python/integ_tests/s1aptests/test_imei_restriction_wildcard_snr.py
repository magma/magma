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


import ctypes
import unittest

import s1ap_types
import s1ap_wrapper


class TestImeiRestrictionWildcardSnr(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_imei_restriction_wildcard_snr(self):
        """
        This TC validates imei restriction scenario where only tac
        part of imei is configured i.e snr is wildcarded:
        1. Send security mode complete message with a blocked imeisv
           where tac is same as the one configured in mme.conf.template
        2. Verify that attach reject is received with cause(5)
           IMEI_NOT_ACCEPTED

        If this TC is executed individually run
        test_modify_mme_config_for_sanity.py to add
        { IMEI_TAC="99333821"}
        under the BLOCKED_IMEI_LIST in mme.conf.template.

        After execution of this TC run test_restore_mme_config_after_sanity.py
        to restore the old mme.conf.template.
        """
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)
        req = self._s1ap_wrapper.ue_req
        # Send Attach Request
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
        print(
            "********************** Sent attach req for UE id ", req.ue_id,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        auth_req = response.cast(s1ap_types.ueAuthReqInd_t)
        print(
            "********************** Received auth req for UE id",
            auth_req.ue_Id,
        )

        # Send Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = req.ue_id
        sqnRecvd = s1ap_types.ueSqnRcvd_t()
        sqnRecvd.pres = 0
        auth_res.sqnRcvd = sqnRecvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )
        print("********************** Sent auth rsp for UE id", req.ue_id)

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        sec_mod_cmd = response.cast(s1ap_types.ueSecModeCmdInd_t)
        print(
            "********************** Received security mode cmd for UE id",
            sec_mod_cmd.ue_Id,
        )

        # Send Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = req.ue_id
        sec_mode_complete.imeisv_pres = True
        imeisv = "9933382135103723"
        # Check if the len of imeisv exceeds 16
        self.assertLessEqual(len(imeisv), 16)
        for i in range(0, len(imeisv)):
            sec_mode_complete.imeisv[i] = ctypes.c_ubyte(int(imeisv[i]))
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )
        print(
            "********************** Sent security mode complete for UE id",
            req.ue_id,
        )

        # Receive Attach Reject
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND.value,
        )
        attach_rej = response.cast(s1ap_types.ueAttachRejInd_t)
        print(
            "********************** Received attach reject for UE id %d"
            " with emm cause %d" % (attach_rej.ue_Id, attach_rej.cause),
        )

        # Verify cause
        self.assertEqual(
            attach_rej.cause, s1ap_types.TFW_EMM_CAUSE_IMEI_NOT_ACCEPTED,
        )

        # UE Context release
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )

        ue_context_rel = response.cast(s1ap_types.ueCntxtRelReq_t)
        print(
            "********************** Received UE_CTX_REL_IND for UE id ",
            ue_context_rel.ue_Id,
        )


if __name__ == "__main__":
    unittest.main()
