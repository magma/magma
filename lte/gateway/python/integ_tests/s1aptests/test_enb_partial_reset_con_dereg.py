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
from builtins import range

import s1ap_types
import s1ap_wrapper


class TestEnbPartialReset(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_enb_partial_reset(self):
        """ attach 32 UEs """
        ue_ids = []
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Calling attach for UE id ",
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

            print("********Triggering Attach Request with PND Type IPv4 test")

            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
            )

            # Trigger Authentication Response
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

            # Trigger Security Mode Complete
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
            print("response message type for ATTTACH ACC", response.msg_type)
            print("acc", s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value)
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
            )
            ue_ids.append(req.ue_id)

        # Trigger eNB Reset
        # Add delay to ensure S1APTester sends attach partial before sending
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
                "Reset_req.r.partialRst.ueS1apIdPairList[indx].ueId",
                reset_req.r.partialRst.ueS1apIdPairList[indx].ueId,
                indx,
            )
        print("ue_ids", ue_ids)
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.RESET_REQ, reset_req,
        )
        response1 = self._s1ap_wrapper.s1_util.get_response()
        print("response1 message type", response1.msg_type)
        self.assertEqual(response1.msg_type, s1ap_types.tfwCmd.RESET_ACK.value)
        # Trigger detach request
        """time.sleep(0.5)
        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            #self._s1ap_wrapper.s1_util.detach(
            #    ue, detach_type, wait_for_s1)
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True)
        """


if __name__ == "__main__":
    unittest.main()
