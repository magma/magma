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


class TestSecurityModeReject(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_security_mode_reject(self):
        """ Testing of security mode reject procedure """
        num_ues = 1

        self._s1ap_wrapper.configUEDevice_ues_same_imsi(num_ues)
        print("************************* sending Attach Request for ue-id : 1")
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = 1
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        pdn_type.pdn_type = 1
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type

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

        sec_mode_reject = s1ap_types.ueSecModeReject_t()
        sec_mode_reject.ue_Id = 1
        sec_mode_reject.cause = s1ap_types.TFW_EMM_CAUSE_SEC_MOD_REJ_UNSP
        print("Sending Security Mode Reject ue-id", sec_mode_reject.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_REJECT, sec_mode_reject,
        )
        time.sleep(2)


if __name__ == "__main__":
    unittest.main()
