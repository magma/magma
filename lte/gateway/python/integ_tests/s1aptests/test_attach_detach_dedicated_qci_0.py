"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import time
import unittest

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestAttachDetachDedicatedQci0(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_dedicated_qci_0(self):
        """ test attach + create dedicated bearer with QCI, 0 +
        erab_setup_failed_indication + detach, with a single UE """
        num_ues = 1
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
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

            time.sleep(2)
            imsi = "".join([str(i) for i in req.imsi])
            print(
                "********************** Adding dedicated bearer to IMSI",
                imsi,
            )
            self._spgw_util.create_bearer("IMSI" + imsi, 5, 0)

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_FW_ERAB_SETUP_REQ_FAILED_FOR_ERABS.value,
            )
            erab_setup_failed_for_bearers = response.cast(
                s1ap_types.FwErabSetupFailedTosetup
            )
            print(
                "*** Received UE_FW_ERAB_SETUP_REQ_FAILED_FOR_ERABS for "
                "bearer-id:",
                erab_setup_failed_for_bearers.failedErablist[0].erabId,
                end=" ",
            )
            print(
                " with qci:",
                erab_setup_failed_for_bearers.failedErablist[0].qci
            )

            time.sleep(5)
            print(
                "********************** Running UE detach for UE id ",
                req.ue_id
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i]
            )


if __name__ == "__main__":
    unittest.main()
