"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import time

import s1ap_types
import s1ap_wrapper


class TestSecondaryPdnConnLooped(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_seconday_pdn_conn_looped(self):
        """ Attach a single UE and send standalone PDN Connectivity
        Request + detach. Repeat 3 times """

        num_ues = 1

        # Configure APN before configuring UE device
        # APN of the secondary PDN
        ims = [
            "ims",  # APN-name
            5,  # qci
            15,  # priority
            0,  # preemption-capability
            0,  # preemption-vulnerability
            200000000,  # MBR UL
            100000000,  # MBR DL
        ]

        # APN details to be configured in APN DB
        apn_list = [ims]
        self._s1ap_wrapper.configAPN(apn_list)

        # List of APN names supported by the UE
        apn_supported = ["ims"]
        self._s1ap_wrapper.configUEDevice(num_ues, apn_supported)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        loop = 3
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Attach
        for _ in range(loop):
            time.sleep(5)
            self._s1ap_wrapper.s1_util.attach(
                ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Send PDN Connectivity Request
            apn = "ims"
            self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
            # Receive PDN CONN RSP/Activate default EPS bearer context request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
            )

            print(
                "************************* Sending Activate default EPS "
                "bearer context accept for UE id ",
                ue_id,
            )

            time.sleep(5)

            print(
                "************************* Running UE detach (switch-off)"
                " for UE id ",
                ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
                True,
            )


if __name__ == "__main__":
    unittest.main()
