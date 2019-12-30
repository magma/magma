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


class TestSecondaryPdnWithDedBearerMultiUe(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_seconday_pdn_with_dedicated_bearer_multi_ue(self):
        """ attach/detach + PDN Connectivity Requests + dedicated bearer for 4
            UEs """
        num_ues = 4
        ue_ids = []
        bearer_ids = []
        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id
            print(
                "******************* Running End to End attach for UE id ",
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
            ue_ids.append(ue_id)

        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            # Send PDN Connectivity Request
            apn = "ims"
            ue_id = req.ue_id
            self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
            # Receive PDN CONN RSP/Activate default EPS bearer context request
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
            print(
                "********************** Added IMS PDN with bearer id",
                act_def_bearer_req.m.pdnInfo.epsBearerId,
            )
            bearer_ids.append(act_def_bearer_req.m.pdnInfo.epsBearerId)

            time.sleep(2)
            # Add dedicated bearer to IMS PDN
            print("********************** Adding dedicated bearer to IMS PDN")
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId
            )

            print(
                "********************** Added dedicated bearer",
                act_ded_ber_ctxt_req.bearerId,
            )

        time.sleep(5)
        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id
            print("******************* Deleting IMS PDN for ue", ue_id)
            # Send PDN Disconnect
            pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
            pdn_disconnect_req.ue_Id = ue_id
            pdn_disconnect_req.epsBearerId = bearer_ids[i]
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req
            )

            # Receive UE_DEACTIVATE_BER_REQ
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(
                response, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value
            )

            print(
                "******************* Received deactivate eps bearer context"
                " request"
            )
            # Send DeactDedicatedBearerAccept
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                ue_id, bearer_ids[i]
            )

            print(
                "******************* Deletd IMS PDN with bearer ID",
                bearer_ids[i],
            )
        time.sleep(2)
        # Now detach the UE
        for ue in ue_ids:
            print(
                "******************** Running UE detach (switch-off) for ",
                "UE id ",
                ue,
            )
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False
            )


if __name__ == "__main__":
    unittest.main()
