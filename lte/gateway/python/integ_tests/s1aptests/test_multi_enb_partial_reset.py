"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.
 This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import time
import ctypes
from builtins import range

import s1ap_types
import s1ap_wrapper


class TestMultipleEnbPartialReset(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_multiple_enb_partial_reset(self):
        """ attach 32 UEs """
        num_of_enbs = 5
        # column is a enb parameter,  row is a number of enbs
        """            Cell Id,   Tac, EnbType, PLMN Id, PLMN length """
        enb_list = list([[1,       1,     1,    "001010", 6],
                         [2,       1,     1,    "001010", 6],
                         [3,       1,     1,    "001010", 6],
                         [4,       1,     1,    "001010", 6],
                         [5,       1,     1,    "001010", 6]])

        assert (num_of_enbs == len(enb_list)), "Number of enbs configured"
        "not equal to enbs in the list!!!"

        self._s1ap_wrapper.multiEnbConfig(num_of_enbs, enb_list)

        time.sleep(2)

        ue_ids = []
        num_ues = 1

        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print("************************* Calling attach for UE id ",
                  req.ue_id)
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
        reset_req.r.partialRst.ueIdLst = (
            ctypes.c_ubyte * reset_req.r.partialRst.numOfConn
        )()
        for indx in range(reset_req.r.partialRst.numOfConn):
            reset_req.r.partialRst.ueIdLst[indx] = ue_ids[indx]
            print(
                "Reset_req.r.partialRst.ueIdLst[indx]",
                reset_req.r.partialRst.ueIdLst[indx],
                indx,
            )
        print("ue_ids", ue_ids)
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.RESET_REQ, reset_req)
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(response.msg_type, s1ap_types.tfwCmd.RESET_ACK.value)
        # Trigger detach request
        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            # self._s1ap_wrapper.s1_util.detach(
            #    ue, detach_type, wait_for_s1)
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True
            )


if __name__ == "__main__":
    unittest.main()
