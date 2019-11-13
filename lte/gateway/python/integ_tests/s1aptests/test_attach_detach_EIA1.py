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


class TestAttachDetachEIA1(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_eia1(self):
        """ Basic attach with Integrity algo type EIA1(snow3G) """
        """
        Generating this scenario by configuring UE support only eia1 algo
        not eia2 as MME by default configures eia2 if UE support eia2
        """
        num_ues = 2
        detach_type = [s1ap_types.ueDetachType_t.
                       UE_NORMAL_DETACH.value,
                       s1ap_types.ueDetachType_t.
                       UE_SWITCHOFF_DETACH.value]
        wait_for_s1 = [True, False]
        req = self._s1ap_wrapper._sub_util.add_sub(num_ues)

        for i in range(num_ues):
            print("************************* UE device config for ue_id ",
                  req[i].ue_id)
            req[i].ueNwCap_pr.pres = 1
            req[i].ueNwCap_pr.eea2_128 = 0
            req[i].ueNwCap_pr.eea1_128 = 1
            req[i].ueNwCap_pr.eea0 = 1
            req[i].ueNwCap_pr.eia2_128 = 0
            req[i].ueNwCap_pr.eia1_128 = 1
            req[i].ueNwCap_pr.eia0 = 1

            assert (self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CONFIG, req[i]) == 0)
            response = self._s1ap_wrapper._s1_util.get_response()
            assert (s1ap_types.tfwCmd.UE_CONFIG_COMPLETE_IND.value ==
                    response.msg_type)
            self._s1ap_wrapper._configuredUes.append(req[i])
        self._s1ap_wrapper.check_gw_health_after_ue_load()

        # self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            print("************************* Running End to End attach for ",
                  "UE id ", req[i].ue_id)
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req[i].ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t)

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            print("************************* Running UE detach for UE id ",
                  req[i].ue_id)
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req[i].ue_id, detach_type[i], wait_for_s1[i])


if __name__ == "__main__":
    unittest.main()
