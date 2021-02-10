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
import s1ap_wrapper
import time


class TestDuplicateAttachWithoutAttachAcceptOrReject(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_duplicate_attach_without_attach_accept_or_reject(self):
        """
        UE sends Attach Request, it goes through auth and smc procedures
        But UE doesn't receive neither Attach Accept nor Reject, so
        UE re-transmits the Attach Request after 15 seconds, for whcih
        UE is expected to attach successfully to network
        """
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)

        # Trigger Attach Request
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.ue_Id = 1
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        print("***Triggering Attach Request ***")

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value
        )

        # Trigger Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = 1
        sqnRecvd = s1ap_types.ueSqnRcvd_t()
        sqnRecvd.pres = 0
        auth_res.sqnRcvd = sqnRecvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value
        )

        # Trigger Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = 1
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete
        )

        print("***Triggering Re-Attach Request after 15 seconds ***")
        time.sleep(15)

        req = self._s1ap_wrapper.ue_req
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        print(
            "************************* Running UE detach for UE id ", req.ue_id
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            True,
        )

        time.sleep(1)


if __name__ == "__main__":
    unittest.main()
