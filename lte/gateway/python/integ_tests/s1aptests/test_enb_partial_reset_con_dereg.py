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
from builtins import range

import s1ap_types
import s1ap_wrapper


class TestEnbPartialResetConDereg(unittest.TestCase):
    """Integration Test: TestEnbPartialResetConDereg"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_enb_partial_reset_con_dereg(self):
        """Test ENB partial reset with 1 UE while UE is connected and
        de-registered
        """
        ue_ids = []
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)

        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Calling attach for UE id",
                req.ue_id,
            )
            # Trigger Attach Request with PDN_Type = IPv4v6
            attach_req = s1ap_types.ueAttachRequest_t()
            sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
            id_type = s1ap_types.TFW_MID_TYPE_IMSI
            eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
            pdn_type = s1ap_types.pdn_Type()
            pdn_type.pres = True
            # Set PDN TYPE to IPv4V6 i.e. 3. IPV4 is equal to 1
            # IPV6 is equal to 2 in value
            pdn_type.pdn_type = 1
            attach_req.ue_Id = req.ue_id
            attach_req.mIdType = id_type
            attach_req.epsAttachType = eps_type
            attach_req.useOldSecCtxt = sec_ctxt
            attach_req.pdnType_pr = pdn_type

            print(
                "************************* Triggering Attach Request with PDN "
                "Type IPv4 test",
            )

            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ATTACH_REQUEST,
                attach_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value

            # Trigger Authentication Response
            auth_res = s1ap_types.ueAuthResp_t()
            auth_res.ue_Id = req.ue_id
            sqn_recvd = s1ap_types.ueSqnRcvd_t()
            sqn_recvd.pres = 0
            auth_res.sqnRcvd = sqn_recvd
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_AUTH_RESP,
                auth_res,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value

            # Trigger Security Mode Complete
            sec_mode_complete = s1ap_types.ueSecModeComplete_t()
            sec_mode_complete.ue_Id = req.ue_id
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE,
                sec_mode_complete,
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
            ue_ids.append(req.ue_id)

        # Trigger eNB Reset
        # Add delay to ensure S1APTester sends attach complete before sending
        # eNB Reset Request
        time.sleep(0.5)
        print("************************* Sending eNB Partial Reset Request")
        reset_req = s1ap_types.ResetReq()
        reset_req.rstType = s1ap_types.resetType.PARTIAL_RESET.value
        reset_req.cause = s1ap_types.ResetCause()
        reset_req.cause.causeType = (
            s1ap_types.NasNonDelCauseType.TFW_CAUSE_MISC.value
        )
        # Set the cause to MISC.hardware-failure
        reset_req.cause.causeVal = 3
        reset_req.r = s1ap_types.R()
        reset_req.r.partialRst = s1ap_types.PartialReset()
        reset_req.r.partialRst.numOfConn = num_ues
        reset_req.r.partialRst.ueS1apIdPairList = (
            (s1ap_types.UeS1apIdPair) * reset_req.r.partialRst.numOfConn
        )()
        for indx in range(reset_req.r.partialRst.numOfConn):
            reset_req.r.partialRst.ueS1apIdPairList[indx].ueId = ue_ids[indx]
            print(
                "Reset_req.r.partialRst.ueS1apIdPairList[",
                indx,
                "].ueId",
                reset_req.r.partialRst.ueS1apIdPairList[indx].ueId,
            )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.RESET_REQ,
            reset_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.RESET_ACK.value


if __name__ == "__main__":
    unittest.main()
