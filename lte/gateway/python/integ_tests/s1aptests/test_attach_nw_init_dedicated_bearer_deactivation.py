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


class TestDedicatedBearer(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_different_enb_s1ap_id_same_ue(self):
        """ Testing of Dedicated Bearer Request from Network """
        num_ues = 1

        self._s1ap_wrapper.configUEDevice(num_ues)
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print("**************** sending Attach Request for ue-id : ",
                  req.ue_id)
            attach_req = s1ap_types.ueAttachRequest_t()

            attach_req.ue_Id = req.ue_id
            sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
            id_type = s1ap_types.TFW_MID_TYPE_IMSI
            eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
            attach_req.mIdType = id_type
            attach_req.epsAttachType = eps_type
            attach_req.useOldSecCtxt = sec_ctxt

            print("Sending Attach Request ue-id", attach_req.ue_Id)
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req)

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value)
            print("Received auth req ind ")

            auth_res = s1ap_types.ueAuthResp_t()
            auth_res.ue_Id = req.ue_id
            sqn_recvd = s1ap_types.ueSqnRcvd_t()
            sqn_recvd.pres = 0
            auth_res.sqnRcvd = sqn_recvd
            print("Sending Auth Response ue-id", auth_res.ue_Id)
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res)

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response,
                            s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value)
            print("Received Security Mode Command ue-id", auth_res.ue_Id)

            time.sleep(1)

            sec_mode_complete = s1ap_types.ueSecModeComplete_t()
            sec_mode_complete.ue_Id = req.ue_id
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete)
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response,
                            s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value)
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response,
                            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value)

            # Trigger Attach Complete
            attach_complete = s1ap_types.ueAttachComplete_t()
            attach_complete.ue_Id = req.ue_id
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete)
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            time.sleep(5)

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response,
                            s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value)

            act_ded_ber_ctxt_req = response.cast(
                                   s1ap_types.UeActDedBearCtxtReq_t)
            ded_bearer_acc = s1ap_types.UeActDedBearCtxtAcc_t()
            ded_bearer_acc.ue_Id = req.ue_id
            ded_bearer_acc.bearerId = act_ded_ber_ctxt_req.bearerId

            print("** Bearer ID received in Activate Ded Bearer Context",
                  "Request message **", ded_bearer_acc.bearerId)

            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ACT_DED_BER_ACC, ded_bearer_acc)

            time.sleep(3)

            print("************************* De-Activate EPS Bearer Context "
                  "Request Indication")
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response,
                            s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value)

            print("**********Received De activate eps bearer context")

            deactv_bearer_req = response.cast(
                               s1ap_types.UeDeActvBearCtxtReq_t)
            print("*************************Sending De-Activate EPS Bearer"
                  " Context Accept")
            deactv_bearer_acc = s1ap_types.UeDeActvBearCtxtAcc_t()
            deactv_bearer_acc.ue_Id = req.ue_id
            deactv_bearer_acc.bearerId = deactv_bearer_req.bearerId
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_ACC, deactv_bearer_acc)

            time.sleep(5)
            print("************************* Running UE detach")
            # Now detach the UE
            detach_req = s1ap_types.uedetachReq_t()
            detach_req.ue_Id = req.ue_id
            detach_req.ueDetType = s1ap_types.ueDetachType_t.\
                UE_NORMAL_DETACH.value
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req)
            response = self._s1ap_wrapper._s1_util.get_response()
            self.assertTrue(
                response, s1ap_types.tfwCmd.UE_DETACH_ACCEPT_IND.value)
            # Wait for UE context release command
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertTrue(response, s1ap_types.tfwCmd.UE_CTX_REL_IND.value)


if __name__ == "__main__":
    unittest.main()
