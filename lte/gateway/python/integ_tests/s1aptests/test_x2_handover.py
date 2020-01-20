"""
Copyright (c) 2017-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

import s1ap_types
import s1ap_wrapper
import time


class TestX2HandOver(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_x2_handover(self):
        """ Multi Enb Multi UE attach detach """

        # column is an enb parameter, row is a number of enb
        """            Cell Id,   Tac, EnbType, PLMN Id """
        enb_list = list([[1,       1,     1,    "001010"],
                         [2,       2,     1,    "001010"]])

        self._s1ap_wrapper.multiEnbConfig(len(enb_list), enb_list)

        time.sleep(2)
        """ Attach to Src eNB.HO to TeNB  """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        print("************************* Running End to End attach for ",
              "UE id ", req.ue_id)
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t)

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        time.sleep(3)
        print("************************* Sending ENB_CONFIGURATION_TRANSFER")
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.ENB_CONFIGURATION_TRANSFER, req)

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.MME_CONFIGURATION_TRANSFER.value)

        print("************************* Received MME_CONFIGURATION_TRANSFER")
        print("************************* Sending ENB_CONFIGURATION_TRANSFER")
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.MME_CONFIGURATION_TRANSFER.value)

        print("************************* Received MME_CONFIGURATION_TRANSFER")
        print("************************* Sending X2_HO_TRIGGER_REQ")

        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.X2_HO_TRIGGER_REQ, req)
        # Receive Path Switch Request Ack
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.PATH_SW_REQ_ACK.value)

        print("************************* Received Path Switch Request Ack")

        print("************************* Running UE detach for UE id ",
              req.ue_id)
        # Now detach the UE
        time.sleep(3)
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            False)


if __name__ == "__main__":
    unittest.main()
