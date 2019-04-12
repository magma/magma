"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import s1ap_types
from integ_tests.s1aptests import s1ap_utils
from integ_tests.s1aptests.s1ap_utils import S1ApUtil


class TestS1SetupFailureIncorrectTac(unittest.TestCase):
    """S1 Setup Request with incorrect TAC value """
    s1ap_utils._s1_util = S1ApUtil()
    print("************************* Enb tester config")
    req = s1ap_types.FwNbConfigReq_t()
    req.cellId_pr.pres = True
    req.cellId_pr.cell_id = 10
    req.tac_pr.pres = True
    req.tac_pr.tac = 0

    assert (
        s1ap_utils._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_CONFIG, req) == 0)
    response = s1ap_utils._s1_util.get_response()
    assert (response.msg_type ==
            s1ap_types.tfwCmd.ENB_CONFIG_CONFIRM.value)
    res = response.cast(s1ap_types.FwNbConfigCfm_t)
    assert (res.status == s1ap_types.CfgStatus.CFG_DONE.value)

    req = None
    assert (s1ap_utils._s1_util.issue_cmd(
        s1ap_types.tfwCmd.ENB_S1_SETUP_REQ, req) == 0)
    response = s1ap_utils._s1_util.get_response()
    assert (response.msg_type ==
            s1ap_types.tfwCmd.ENB_S1_SETUP_RESP.value)


if __name__ == "__main__":
    unittest.main()
