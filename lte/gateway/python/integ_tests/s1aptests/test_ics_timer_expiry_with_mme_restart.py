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

import gpp_types
import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestIcsTimerExpiryWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_ics_timer_expiry_with_mme_restart(self):
        """
        The test case validates the functionality of ICS timer expiry while UE
        is in idle mode after mme restart
        Step 1: Attach UE to network
        Step 2: Let the UE go to Idle mode
        Step 3: Send Service Request to move the ECM state to connected mode
        Step 4: Send an indication to s1ap module to drop ICS Request message
                So s1ap shall not send ICS Response message to mme
        Step 5: While ICS timer is running at mme, restart mme service
        Step 6: On mme recovery, it is expected  that ICS timer has lapsed and
                Send Ue Context Release Command to s1ap
        """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        time.sleep(0.5)

        print(
            "************************* Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        req = s1ap_types.ueCntxtRelReq_t()
        req.ue_Id = ue_id
        req.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )

        print("*** Sending indication to drop Initial Context Setup Req ***")
        drop_init_ctxt_setup_req = s1ap_types.UeDropInitCtxtSetup()
        drop_init_ctxt_setup_req.ue_Id = ue_id
        drop_init_ctxt_setup_req.flag = 1
        # Timer to release UE context at s1ap tester
        drop_init_ctxt_setup_req.tmrVal = 2000
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SET_DROP_ICS, drop_init_ctxt_setup_req,
        )

        print(
            "************************* Sending Service request for UE id ",
            ue_id,
        )
        # Send service request to reconnect UE
        req = s1ap_types.ueserviceReq_t()
        req.ue_Id = ue_id
        req.ueMtmsi = s1ap_types.ueMtmsi_t()
        req.ueMtmsi.pres = False
        req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SERVICE_REQUEST, req,
        )

        # enbApp sends UE_ICS_DROPD_IND message to tfwApp after dropping
        # ICS request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ICS_DROPD_IND.value,
        )
        print("************************* Restarting MME service on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )

        print(
            "************************* Received UE_CTX_REL_IND for UE id ",
            ue_id,
        )
        print("************************* Running UE detach for UE id ", ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
        )


if __name__ == "__main__":
    unittest.main()
