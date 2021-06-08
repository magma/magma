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


class TestAttachMobileReachabilityTimerExpiry(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_mobile_reachability_timer_expiry(self):
        """ Test Mobile Reachability timer expiry handling """

        """ Note: Before execution of this test case,
        Run the test script s1aptests/test_modify_mme_config_for_sanity.py
        to reduce mobile reachability timer value to 1 minute (default is 54
        minutes) in MME configuration and
        after test case execution, restore the MME configuration by running
        the test script s1aptests/test_restore_mme_config_after_sanity.py

        Or

        Manually update the mme.conf.template file to make sure that the value
        of T3412 timer is set to 1 minute (default is 54 minutes)
        """

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
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

        # Delay by 6 minutes to ensure Mobile reachability timer expires.
        # MOBILE REACHABILITY TIMER VALUE = 1 minute (conf file) + delta value
        print(
            "************************* Waiting for Mobile Reachability Timer"
            " to expire. Sleeping for 360 seconds..",
        )
        timeSlept = 0
        while timeSlept < 360:
            time.sleep(5)
            timeSlept += 5
            print("*********** Slept for " + str(timeSlept) + " seconds")

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
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )

        print("************************* Running UE detach for UE id ", ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
