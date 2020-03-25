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
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestAttachDetachDedicatedMultiUe(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach(self):
        """ attach/detach + dedicated bearer test with 4 UEs """
        num_ues = 4
        ue_ids = []
        bearer_ids = []
        self._s1ap_wrapper.configUEDevice(num_ues)

        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Adding dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), 5
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            print(
                "********************** Received activate dedicated EPS"
                " bearer context request"
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId
            )
            bearer_ids.append(act_ded_ber_ctxt_req.bearerId)

        time.sleep(1)
        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Deleting dedicated bearer for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), 5, bearer_ids[i]
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print(
                "********************** Received deactivate EPS bearer"
                " context request"
            )

            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, bearer_ids[i]
            )

        time.sleep(2)
        for ue in ue_ids:
            print("********************** Running UE detach for UE id ", ue)
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value
            )


if __name__ == "__main__":
    unittest.main()
