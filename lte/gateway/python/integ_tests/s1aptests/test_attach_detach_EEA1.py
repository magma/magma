"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

import s1ap_types

from integ_tests.s1aptests import s1ap_wrapper


class TestAttachDetachCiphAlgoEea1(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_ciphering_algo_type_eea1(self):
        """ Basic attach/detach test with encrytion algo type Snow3G/EEA1
        Manually configure the preferred encrytion algo type in
        mme.conf.template file for testing this testcase"""
        num_ues = 2
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            req.ueNwCap_pr.pres = 1
            req.ueNwCap_pr.eea2_128 = 0
            req.ueNwCap_pr.eea1_128 = 1
            req.ueNwCap_pr.eea0 = 1
            req.ueNwCap_pr.eia2_128 = 0
            req.ueNwCap_pr.eia1_128 = 1
            req.ueNwCap_pr.eia0 = 1

            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CONFIG, req
            )
            response = self._s1ap_wrapper._s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_CONFIG_COMPLETE_IND.value,
            )

            print(
                "************************* Running End to End attach for ",
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
            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i]
            )


if __name__ == "__main__":
    unittest.main()
