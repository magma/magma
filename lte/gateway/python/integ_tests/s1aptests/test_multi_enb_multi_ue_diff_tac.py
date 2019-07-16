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
import s1ap_wrapper


class TestMultiEnbWithDifferentTac(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_multienb_different_tac(self):
        """ Multi Enb attach with different TAC values """
        req = s1ap_types.multiEnbConfigReq_t()
        # Num of Enbs (Max Number of ENBS is 5)
        num_of_enbs = 5
        plmn_length = 6
        # column is a enb parameter,  row is a number of enbs
        """            Cell Id,   Tac, EnbType, PLMN Id """
        enb_list = list([[1,       1,     1,    "001010"],
                         [2,       2,     1,    "001010"],
                         [3,       3,     1,    "001010"],
                         [4,       4,     1,    "001010"],
                         [5,       5,     1,    "001010"]])
        # Maximum 5 Enbs can be configured
        req.numOfEnbs = num_of_enbs

        # ENB Parameter column index initialization
        cellid_col_idx = 0
        tac_col_idx = 1
        enbtype_col_idx = 2
        plmnid_col_idx = 3

        for idx1 in range(num_of_enbs):
            req.multiEnbCfgParam[idx1].cell_id = \
                    (enb_list[idx1][cellid_col_idx])

        for idx1 in range(num_of_enbs):
            req.multiEnbCfgParam[idx1].tac = \
                    (enb_list[idx1][tac_col_idx])

        for idx1 in range(num_of_enbs):
            req.multiEnbCfgParam[idx1].enbType = \
                    (enb_list[idx1][enbtype_col_idx])

        for idx1 in range(num_of_enbs):
            for idx3 in range(plmn_length):
                val = enb_list[idx1][plmnid_col_idx][idx3]
                req.multiEnbCfgParam[idx1].plmn_id[idx3] = int(val)

        print("***************** Sending Multiple Enb Config Request\n")
        assert (self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.MULTIPLE_ENB_CONFIG_REQ, req) == 0)

        time.sleep(2)

        ue_ids = []
        # UEs will attach to the ENBs in a round-robin fashion
        # each ENBs will be connected with 32UEs
        num_ues = 5
        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print("******************** Calling attach for UE id ",
                  req.ue_id)
            self._s1ap_wrapper.s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t)
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            self._s1ap_wrapper.s1_util.detach(
                ue,
                s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value)


if __name__ == "__main__":
    unittest.main()
