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


class TestSecondaryPdnConnWithDedBearerDeactivateReq(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_seconday_pdn_conn_ded_bearer_deactivate(self):
        """ Attach a single UE and send standalone PDN Connectivity
        Request + add dedicated bearer to each default bearer + deactivate
        dedicated bearers + detach"""
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
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

            time.sleep(5)
            # Add dedicated bearer to 2nd PDN
            print("********************** Adding dedicated bearer to ims PDN")
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            act_ded_ber_req_ims_apn = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_req_ims_apn.bearerId
            )
            print(
                "************* Added dedicated bearer",
                act_ded_ber_req_ims_apn.bearerId,
            )

            time.sleep(5)
            # Delete dedicated bearer of secondary PDN (ims apn)
            print(
                "********************** Deleting dedicated bearer for ims"
                " apn"
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
                act_ded_ber_req_ims_apn.bearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)

            # Send Deactivate dedicated bearer rsp
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId
            )

            print(
                "********************** Deleted dedicated bearer ",
                deactv_bearer_req.bearerId,
            )
            time.sleep(5)
            # Delete dedicated bearer of secondary PDN (oai.ipv4 apn)
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                5,
                act_ded_ber_req_oai_apn.bearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
            # Send Deactivate dedicated bearer rsp
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId
            )

            print(
                "********************** Deleted dedicated bearer ",
                deactv_bearer_req.bearerId,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                ue_id,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                False,
            )


if __name__ == "__main__":
    unittest.main()
