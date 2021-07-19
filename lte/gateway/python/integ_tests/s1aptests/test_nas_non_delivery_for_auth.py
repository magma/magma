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


class TestNasNonDeliveryAuthReq(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_nas_non_delivery_auth_req(self):
        """ testing Nas Non Delivery functionality for Auth Req for a
             single UE """
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running Nas Non Delivery of Auth Req"
            "for a single UE UE id ",
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
        print("Sending Attach Request for ue-id", req.ue_id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        """ The purpose of UE_SET_NAS_NON_DELIVERY command is to prepare
        enbapp to send S1ap-nas non delivery message for next receiving
        downlink nas message """

        nas_non_del = s1ap_types.UeNasNonDel()
        nas_non_del.ue_Id = req.ue_id
        nas_non_del.flag = 1
        nas_non_del.causeType = (
            s1ap_types.NasNonDelCauseType.TFW_CAUSE_RADIONW.value
        )
        nas_non_del.causeVal = 3
        print("Sending Set Nas Non Del to enb for ue-id ", nas_non_del.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_NAS_NON_DELIVERY, nas_non_del,
        )

        """ Waiting for UE Context Release from MME """
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
        print("Sending Reset Nas Non Del ind to enbapp")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_NAS_NON_DELIVERY, nas_non_del,
        )


if __name__ == "__main__":
    unittest.main()
