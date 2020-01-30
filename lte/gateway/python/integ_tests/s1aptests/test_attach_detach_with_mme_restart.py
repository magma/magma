"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

import s1ap_types
import time

from integ_tests.s1aptests import s1ap_wrapper


class TestAttachDetachWithMmeRestart(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach(self):
        """
        Basic attach/detach test with two UEs,
        where MME restarts between each attach and detach
        """
        num_ues = 2
        detach_type = [s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
                       s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
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

            print("************************* Restarting MME service on",
                  "gateway")
            self._s1ap_wrapper.magmad_util.restart_services(["mme"])

            for j in range(30):
                print("Waiting for", j, "seconds")
                time.sleep(1)

            # Now detach the UE
            print("************************* Running UE detach for UE id ",
                  req.ue_id)
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i])

            if i == 1:
                break

            for j in range(15):
                print("Connecting next UE in", 15 - j, "seconds")
                time.sleep(1)


if __name__ == "__main__":
    unittest.main()
