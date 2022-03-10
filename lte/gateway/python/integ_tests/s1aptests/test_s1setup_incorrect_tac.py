"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import unittest

import s1ap_types
from integ_tests.common.magmad_client import MagmadServiceGrpc
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import MagmadUtil, S1ApUtil


class TestS1SetupFailureIncorrectTac(unittest.TestCase):
    """Integration Test: TestS1SetupFailureIncorrectTac"""

    def setUp(self):
        """Initialize before test case execution"""
        if not s1ap_wrapper.TestWrapper.IS_TEST_RUNNING_FIRST_TIME:
            print("\n**Running the test case again to identify flaky behavior")
        s1ap_wrapper.TestWrapper.IS_TEST_RUNNING_FIRST_TIME = False
        self._s1_util = S1ApUtil()

    def tearDown(self):
        """Cleanup after test case execution"""
        is_test_passed = s1ap_wrapper.TestWrapper.get_test_status(self)
        print("************************* Sending SCTP SHUTDOWN")
        self._s1_util.issue_cmd(s1ap_types.tfwCmd.SCTP_SHUTDOWN_REQ, None)

        if not is_test_passed:
            print("************************* Cleaning up TFW")
            self._s1_util.issue_cmd(s1ap_types.tfwCmd.TFW_CLEANUP, None)

        self._s1_util.cleanup()
        if not is_test_passed:
            print("The test has failed. Restarting Sctpd for cleanup")
            magmad_client = MagmadServiceGrpc()
            magmad_util = MagmadUtil(magmad_client)
            magmad_util.restart_sctpd()
            magmad_util.print_redis_state()

    def test_s1setup_failure_incorrect_tac(self):
        """S1 Setup Request with incorrect TAC value"""

        print("************************* Enb tester configuration")
        req = s1ap_types.FwNbConfigReq_t()
        req.cellId_pr.pres = True
        req.cellId_pr.cell_id = 10
        req.tac_pr.pres = True
        req.tac_pr.tac = 0

        print("************************* Sending ENB configuration Request")
        assert self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_CONFIG, req) == 0
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_CONFIG_CONFIRM.value
        res = response.cast(s1ap_types.FwNbConfigCfm_t)
        assert res.status == s1ap_types.CfgStatus.CFG_DONE.value

        print("************************* Sending S1-Setup Request")
        req = None
        assert (
            self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_S1_SETUP_REQ, req)
            == 0
        )
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_S1_SETUP_RESP.value
        res = response.cast(s1ap_types.FwNbS1setupRsp_t)
        assert res.res == s1ap_types.S1_setp_Result.S1_SETUP_FAILED.value


if __name__ == "__main__":
    unittest.main()
