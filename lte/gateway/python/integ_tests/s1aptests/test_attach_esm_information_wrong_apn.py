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

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper


class TestEsmInformationWrongApn(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_esm_information_wrong_apn(self):
        """ Testing of sending wrong APN in Esm Information procedure """
        num_ues = 1

        self._s1ap_wrapper.configUEDevice_ues_same_imsi(num_ues)
        print("************************* sending Attach Request for ue-id : 1")
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = 1
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        # enabling ESM Information transfer flag
        attach_req.eti.pres = 1
        attach_req.eti.esm_info_transfer_flag = 1

        print("Sending Attach Request ue-id", attach_req.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("Received auth req ind ")

        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = 1
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        print("Sending Auth Response ue-id", auth_res.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        print("Received Security Mode Command ue-id", auth_res.ue_Id)

        time.sleep(1)

        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = 1
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )

        # Esm Information Request indication
        print(
            "Received Esm Information Request ue-id", sec_mode_complete.ue_Id,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ESM_INFORMATION_REQ.value,
        )
        esm_info_req = response.cast(s1ap_types.ueEsmInformationReq_t)
        # Sending Esm Information Response
        print(
            "Sending Esm Information Response ue-id", sec_mode_complete.ue_Id,
        )
        esm_info_response = s1ap_types.ueEsmInformationRsp_t()
        esm_info_response.ue_Id = 1
        esm_info_response.tId = esm_info_req.tId
        esm_info_response.pdnAPN_pr.pres = 1
        s = "oai"
        esm_info_response.pdnAPN_pr.len = len(s)
        esm_info_response.pdnAPN_pr.pdn_apn = (ctypes.c_ubyte * 100)(
            *[ctypes.c_ubyte(ord(c)) for c in s[:100]],
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ESM_INFORMATION_RSP, esm_info_response,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND.value,
        )

        print("*** Running UE detach ***")
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = 1
        detach_req.ueDetType = (
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
        )
        # Wait for UE context release command
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )


if __name__ == "__main__":
    unittest.main()
