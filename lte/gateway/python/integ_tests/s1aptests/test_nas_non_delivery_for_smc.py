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

import unittest

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper


class TestNasNonDeliverySmc(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_nas_non_delivery_smc(self):
        """ testing Nas Non Delivery for Security Mode Command message
            for single UE """
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running Nas Non Delivery of Security"
            "Mode Command message for a single UE with UE id",
            req.ue_id,
        )

        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = req.ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        print("Sending Attach Request ue-id", req.ue_id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("Received auth req ind ue-Id", req.ue_id)

        """ The purpose of UE_SET_NAS_NON_DELIVERY command is to prepare
        enbapp to send S1ap-nas non delivery message for next receiving
        downlink nas message
        For testing of nas non delivery of SMC, the NasNonDel flag is set at
        enbApp, before sending Auth Response, to avoid race condition of SMC
        receiving first and then setting NasNonDel flag """

        nas_non_del = s1ap_types.UeNasNonDel()
        nas_non_del.ue_Id = req.ue_id
        nas_non_del.flag = 1
        nas_non_del.causeType = (
            s1ap_types.NasNonDelCauseType.TFW_CAUSE_RADIONW.value
        )
        nas_non_del.causeVal = 3
        print("Sending Set Nas Non Del to enbApp for ue-id", nas_non_del.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_NAS_NON_DELIVERY, nas_non_del,
        )

        print("Send Auth Resp", req.ue_id)
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
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("Received UE_CTX_REL_IND")

        # Reset the nas non delivery flag
        nas_non_del = s1ap_types.UeNasNonDel()
        nas_non_del.ue_Id = req.ue_id
        nas_non_del.flag = 0
        nas_non_del.causeType = (
            s1ap_types.NasNonDelCauseType.TFW_CAUSE_RADIONW.value
        )
        nas_non_del.causeVal = 3
        print("Sending Reset Nas Non Del ind to enb")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_NAS_NON_DELIVERY, nas_non_del,
        )


if __name__ == "__main__":
    unittest.main()
