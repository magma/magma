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
import time
import unittest

import gpp_types
import s1ap_types
import s1ap_wrapper


class TestGutiAttachWithZeroMtmsi(unittest.TestCase):
    """GUTI attach test with 0 M-TMSI for single UE"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_guti_attach_with_zero_mtmsi(self):
        """1. Perform IMSI attach
        2. Move UE to idle mode
        3. Send GUTI attach request for the same UE with M-TMSI value 0
        4. Detach the UE
        """
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)
        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for UE id:",
            req.ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        time.sleep(0.5)

        # Send UE context release request to move UE to idle mode
        print(
            "************************* Sending UE context release request ",
            "for UE id ",
            req.ue_id,
        )
        ureq = s1ap_types.ueCntxtRelReq_t()
        ureq.ue_Id = req.ue_id
        ureq.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST,
            ureq,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value

        time.sleep(5)

        # Send GUTI attach request with M-TMSI value as 0
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_GUTI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        pdn_type.pdn_type = 1
        mcc = "001"
        mnc = "01"
        mcc_len = len(mcc)
        mnc_len = len(mnc)
        attach_req.ue_Id = req.ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type
        attach_req.guti_mi.pres = True
        for i in range(0, mcc_len):
            attach_req.guti_mi.guti.mcc[i] = ctypes.c_ubyte(int(mcc[i]))
        for i in range(0, mnc_len):
            attach_req.guti_mi.guti.mnc[i] = ctypes.c_ubyte(int(mnc[i]))
        attach_req.guti_mi.guti.mmeGrdId = 314
        attach_req.guti_mi.guti.mmeCode = 30
        attach_req.guti_mi.guti.mTmsi = 0

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST,
            attach_req,
        )
        print(
            "********************** Sent attach req for UE id ",
            attach_req.ue_Id,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_IDENTITY_REQ_IND.value
        id_req = response.cast(s1ap_types.ueIdentityReqInd_t)
        print(
            "********************** Received identity req for UE id",
            id_req.ue_Id,
        )

        identity_resp = s1ap_types.ueIdentityResp_t()
        identity_resp.ue_Id = id_req.ue_Id
        identity_resp.idType = s1ap_types.TFW_MID_TYPE_IMSI
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_IDENTITY_RESP,
            identity_resp,
        )
        print(
            "********************** Sent identity rsp for UE id",
            id_req.ue_Id,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value
        auth_req = response.cast(s1ap_types.ueAuthReqInd_t)
        print(
            "********************** Received auth req for UE id",
            auth_req.ue_Id,
        )
        # Send Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = auth_req.ue_Id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP,
            auth_res,
        )
        print("********************** Sent auth rsp for UE id", auth_req.ue_Id)

        response = self._s1ap_wrapper.s1_util.get_response()
        sec_mode_cmd = response.cast(s1ap_types.ueSecModeCmdInd_t)

        assert response.msg_type == s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value
        print(
            "********************** Received security mode cmd for UE id",
            sec_mode_cmd.ue_Id,
        )
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = sec_mode_cmd.ue_Id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE,
            sec_mode_complete,
        )
        print(
            "********************** Sent security mode complete for UE id",
            sec_mode_cmd.ue_Id,
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

        # Trigger Attach Complete
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = attach_acc.ue_Id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE,
            attach_complete,
        )
        print(
            "********************** Sent attach complete for UE id",
            attach_complete.ue_Id,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        print(
            "************************* Running UE detach for UE id ",
            attach_complete.ue_Id,
        )
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = attach_complete.ue_Id
        detach_req.ueDetType = (
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST,
            detach_req,
        )
        # Wait for UE context release command
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value


if __name__ == "__main__":
    unittest.main()
