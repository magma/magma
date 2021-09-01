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


class TestMultipleEnbPartialReset(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_multiple_enb_partial_reset(self):
        """ Multi eNB + attach 1 UE + s1ap partial reset + detach """

        """ Note: Before execution of this test case,
        make sure that following steps are correct
        1. Configure same plmn and tac in both MME and s1ap tester
        2. How to configure plmn and tac in MME:
           a. Set mcc and mnc in gateway.mconfig for mme service
           b. Set tac in gateway.mconfig for mme service
           c. Restart MME service
        3. How to configure plmn and tac in s1ap tester,
           a. For multi-eNB test case, configure plmn and tac from test case.
             In each multi-eNB test case, set plmn, plmn length and tac
             in enb_list
           b. For single eNB test case, configure plmn and tac in nbAppCfg.txt
        """

        # column is an enb parameter, row is number of enbs
        """         Cell Id, Tac, EnbType, PLMN Id, PLMN length """
        enb_list = [
            [1, 1, 1, "00101", 5],
            [2, 1, 1, "00101", 5],
            [3, 1, 1, "00101", 5],
            [4, 1, 1, "00101", 5],
            [5, 1, 1, "00101", 5],
        ]

        self._s1ap_wrapper.multiEnbConfig(len(enb_list), enb_list)

        time.sleep(2)
        ue_ids = []
        num_ues = 1

        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Calling attach for UE id ",
                req.ue_id,
            )
            self._s1ap_wrapper.s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            ue_ids.append(req.ue_id)
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
        # Trigger eNB Reset
        # Add delay to ensure S1APTester sends attach partial before sending
        # eNB Reset Request
        time.sleep(0.5)
        print("************************* Sending eNB Partial Reset Request")
        reset_req = s1ap_types.ResetReq()
        reset_req.rstType = s1ap_types.resetType.PARTIAL_RESET.value
        reset_req.cause = s1ap_types.ResetCause()
        reset_req.cause.causeType = \
            s1ap_types.NasNonDelCauseType.TFW_CAUSE_MISC.value
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
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(response.msg_type, s1ap_types.tfwCmd.RESET_ACK.value)
        # Trigger detach request
        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            # self._s1ap_wrapper.s1_util.detach(
            #    ue, detach_type, wait_for_s1)
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
            )


if __name__ == "__main__":
    unittest.main()
