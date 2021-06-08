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

import random
import time
import unittest
from builtins import range

import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestEnbPartialResetMultiUeWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_enb_partial_reset_multi_ue_with_mme_restart(self):
        """ENB Partial Reset with MME restart for multiple UEs:
        1) Attach 32 UEs
        2) Send partial reset for a random subset of the attached UEs
        3) Restart MME
        4) Send service request for all the UEs which were reset
        5) Detach all the 32 UEs
        """
        ue_ids = []
        num_ues = 32
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

        # Add delay to ensure S1APTester sends attach complete before sending
        # eNB Reset Request
        time.sleep(0.5)

        # Set the reset UEs list
        random.seed(time.clock())
        reset_ue_count = random.randint(1, num_ues)
        random.seed(time.clock())
        reset_ue_list = random.sample(range(num_ues), reset_ue_count)

        print(
            "************************* Sending eNB Partial Reset Request for",
            reset_ue_count,
            "UEs",
        )
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
        reset_req.r.partialRst.numOfConn = reset_ue_count
        reset_req.r.partialRst.ueS1apIdPairList = (
            (s1ap_types.UeS1apIdPair) * reset_req.r.partialRst.numOfConn
        )()
        for indx in range(reset_req.r.partialRst.numOfConn):
            reset_req.r.partialRst.ueS1apIdPairList[indx].ueId = ue_ids[
                reset_ue_list[indx]
            ]
            print(
                "Reset_req.r.partialRst.ueS1apIdPairList[",
                indx,
                "].ueId",
                reset_req.r.partialRst.ueS1apIdPairList[indx].ueId,
            )

        # Send eNB Partial Reset
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.RESET_REQ, reset_req,
        )

        print("************************* Restarting MME service on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(response.msg_type, s1ap_types.tfwCmd.RESET_ACK.value)

        # Send service request for all the Reset UE Ids
        for indx in range(reset_req.r.partialRst.numOfConn):
            print(
                "************************* Sending Service request for UE id",
                reset_req.r.partialRst.ueS1apIdPairList[indx].ueId,
            )
            req = s1ap_types.ueserviceReq_t()
            req.ue_Id = reset_req.r.partialRst.ueS1apIdPairList[indx].ueId
            req.ueMtmsi = s1ap_types.ueMtmsi_t()
            req.ueMtmsi.pres = False
            req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST, req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
            )

        # Trigger detach request
        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
            )


if __name__ == "__main__":
    unittest.main()
