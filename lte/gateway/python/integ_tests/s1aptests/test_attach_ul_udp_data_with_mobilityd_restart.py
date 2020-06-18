"""
Copyright (c) 2017-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import time
import unittest

import s1ap_types
import s1ap_wrapper


class TestAttachUlUdpDataWithMobilitydRestart(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_ul_udp_data(self):
        """
        Attach, send UL UDP data, restart Mobilityd and
        send UL UDP data again
        """
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

        print("************************* Running UE uplink (UDP) for UE id ",
              req.ue_id)
        with self._s1ap_wrapper.configUplinkTest(
                req, duration=1, is_udp=True) as test:
            test.verify()

        print("************************* Restarting Mobilityd service",
              "on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mobilityd"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        print("************************* Running UE uplink (UDP) for UE id ",
              req.ue_id)
        with self._s1ap_wrapper.configUplinkTest(
                req, duration=1, is_udp=True) as test:
            test.verify()

        print("************************* Running UE detach for UE id ",
              req.ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            True)


if __name__ == "__main__":
    unittest.main()
