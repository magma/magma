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
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestMaximumBearersPerUe(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_max_bearers_per_ue(self):
        """ Attach a single UE and send standalone PDN Connectivity
        Request + add 9 dedicated bearers + detach"""
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

        # 1 oai PDN + 1 dedicated bearer, 1 ims pdn + 8 dedicated bearers
        loop = 8

        for i in range(num_ues):
            ue_id = req.ue_id
            print(
                "********************* Running End to End attach for UE id ",
                ue_id,
            )
            # Attach
            self._s1ap_wrapper.s1_util.attach(
                ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Add dedicated bearer for default bearer 5
            print(
                "********************** Adding dedicated bearer to oai.ipv4"
                " PDN"
            )
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), 5
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            act_ded_ber_req_oai_apn = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_req_oai_apn.bearerId
            )

            time.sleep(5)
            # Send PDN Connectivity Request
            apn = "ims"
            self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
            # Receive PDN CONN RSP/Activate default EPS bearer context req
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
            )
            act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)

            print(
                "********************** Sending Activate default EPS bearer "
                "context accept for UE id ",
                ue_id,
            )

            time.sleep(5)
            for _ in range(loop):
                # Add dedicated bearer to 2nd PDN
                print(
                    "********************** Adding dedicated bearer to ims"
                    " PDN"
                )
                self._spgw_util.create_bearer(
                    "IMSI" + "".join([str(i) for i in req.imsi]),
                    act_def_bearer_req.m.pdnInfo.epsBearerId,
                )

                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertEqual(
                    response.msg_type,
                    s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
                )
                act_ded_ber_req_ims_apn = response.cast(
                    s1ap_types.UeActDedBearCtxtReq_t
                )
                self._s1ap_wrapper.sendActDedicatedBearerAccept(
                    req.ue_id, act_ded_ber_req_ims_apn.bearerId
                )
                print(
                    "************ Added dedicated bearer",
                    act_ded_ber_req_ims_apn.bearerId,
                )
                time.sleep(2)

            print(
                "************************ Running UE detach (switch-off) for ",
                "UE id ",
                ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                False,
            )


if __name__ == "__main__":
    unittest.main()
