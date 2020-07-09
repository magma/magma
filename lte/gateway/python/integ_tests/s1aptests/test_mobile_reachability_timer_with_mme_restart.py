"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import time

import gpp_types
import s1ap_types
import s1ap_wrapper


class TestMobileReachabilityTimerWithMmeRestart(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_mobile_reachability_timer_with_mme_restart(self):
        """
        The test case validates Mobile Reachability Timer resumes the
        configured timer value on MME restart
        Step1 : UE attaches to network
        Step2 : UE moves to Idle state
        Step3 : Once MME restarts, MME shall resume
                the Mobile reachability timer for remaining time, on expiry
                MME starts the Implicit Detach Timer. On expiry of
                Implicit Detach Timer, MME implicitly detaches UE.
                MME shall delete the contexts locally
        Step4 : Send Service Request, after Implicit Detach Timer expiry
                expecting Service Reject, as MME has released the UE contexts

        """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print("************************* Running End to End attach for UE id ",
              ue_id)
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t)

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        time.sleep(0.5)

        print("************************* Sending UE context release request ",
              "for UE id ", ue_id)
        # Send UE context release request to move UE to idle mode
        req = s1ap_types.ueCntxtRelReq_t()
        req.ue_Id = ue_id
        req.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req)
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value)

        print("************************* Restarting MME service on",
              "gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        print("Waiting for Mobile Reachability Timer (58 Minutes) and"
              " Implicit Detach Timer (58 minutes) to expire"
              " together timer value is set to 7020 seconds")
        # 58 Minutes + 58 minutes = 116 minutes (6960 seconds)
        # 6960 seconds + 60 seconds, delta(Randomly chosen)
        time.sleep(7020)
        print("************************* Sending Service request for UE id ",
              ue_id)
        # Send service request to reconnect UE
        req = s1ap_types.ueserviceReq_t()
        req.ue_Id = ue_id
        req.ueMtmsi = s1ap_types.ueMtmsi_t()
        req.ueMtmsi.pres = False
        req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SERVICE_REQUEST, req)
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SERVICE_REJECT_IND.value)

        print("************************* Received Service Reject for UE id ",
              ue_id)

        # Wait for UE Context Release command
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value)


if __name__ == "__main__":
    unittest.main()
