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


class TestAttachDetachNwTriggeredDeleteSecondaryPdn(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_nw_triggered_delete_secondary_pdn(self):
        """ Attach a single UE + add secondary PDN + add dedicated bearer
        to the secondary pdn + delete the secondary pdn + detach"""
        num_ue = 1

        self._s1ap_wrapper.configUEDevice(num_ue)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id

        # APN of the secondary PDN
        ims = {
            "apn_name": "ims",  # APN-name
            "qci": 5,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
        }

        # APN list to be configured
        apn_list = [ims]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list
        )
        print(
            "******************* Running End to End attach for UE id ", ue_id,
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

        # Send PDN Connectivity Request
        apn = "ims"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)

        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
        )
        act_sec_pdn = response.cast(s1ap_types.uePdnConRsp_t)

        print(
            "******************* Sending Activate default EPS bearer "
            "context accept for UE id ",
            ue_id,
        )

        # Add dedicated bearer to IMS PDN
        print(
            "******************* Adding dedicated bearer to IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        self._spgw_util.create_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            act_sec_pdn.m.pdnInfo.epsBearerId,
        )

        # Receive Activate dedicated EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        )
        act_ded_ber_ctxt_req = response.cast(s1ap_types.UeActDedBearCtxtReq_t)
        # Send Activate dedicated EPS bearer context accept
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_ctxt_req.bearerId
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        print(
            "******************* Deleting default bearer for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        # Delete secondary pdn
        self._spgw_util.delete_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            act_sec_pdn.m.pdnInfo.epsBearerId,
            act_sec_pdn.m.pdnInfo.epsBearerId,
        )

        # Receive UE_DEACTIVATE_BER_REQ
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
        )

        deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
        print(
            "******************* Received deactivate eps bearer context"
            " request"
        )
        print(
            "******************* Sending deactivate eps bearer context"
            " accept"
        )

        # Send Deactivate EPS bearer context accept
        self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
            ue_id, deactv_bearer_req.bearerId
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        print(
            "******************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False
        )


if __name__ == "__main__":
    unittest.main()
