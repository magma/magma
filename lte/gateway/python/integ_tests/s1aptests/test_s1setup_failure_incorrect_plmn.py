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

import ctypes
import unittest

import s1ap_types
from integ_tests.common.magmad_client import MagmadServiceGrpc
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import MagmadUtil, S1ApUtil


class TestS1SetupFailureIncorrectPlmn(unittest.TestCase):
    """Integration Test: TestS1SetupFailureIncorrectPlmn"""

    def setUp(self):
        """Initialize before test case execution"""
        if s1ap_wrapper.TestWrapper.TEST_CASE_EXECUTION_COUNT != 0:
            print("\n**Running the test case again to identify flaky behavior")
        s1ap_wrapper.TestWrapper.TEST_CASE_EXECUTION_COUNT += 1
        print(
            "Test Case Execution Count:",
            s1ap_wrapper.TestWrapper.TEST_CASE_EXECUTION_COUNT,
        )
        self._s1_util = S1ApUtil()

    def tearDown(self):
        """Cleanup after test case execution"""
        is_test_successful = s1ap_wrapper.TestWrapper.is_test_successful(self)
        print("************************* Sending SCTP SHUTDOWN")
        self._s1_util.issue_cmd(s1ap_types.tfwCmd.SCTP_SHUTDOWN_REQ, None)

        if not is_test_successful:
            print("************************* Cleaning up TFW")
            self._s1_util.issue_cmd(s1ap_types.tfwCmd.TFW_CLEANUP, None)
            self._s1_util.delete_ovs_flow_rules()

        self._s1_util.cleanup()
        if not is_test_successful:
            print("The test has failed. Restarting Sctpd for cleanup")
            magmad_client = MagmadServiceGrpc()
            magmad_util = MagmadUtil(magmad_client)
            magmad_util.restart_services(['sctpd'], wait_time=30)
            magmad_util.print_redis_state()
            if s1ap_wrapper.TestWrapper.TEST_CASE_EXECUTION_COUNT == 3:
                s1ap_wrapper.TestWrapper.generate_flaky_summary()

        elif s1ap_wrapper.TestWrapper.TEST_CASE_EXECUTION_COUNT > 1:
            s1ap_wrapper.TestWrapper.generate_flaky_summary()

    def test_s1setup_failure_incorrect_plmn(self):
        """S1 Setup with incorrect plmn ID"""

        print("************************* Enb tester configuration")
        req = s1ap_types.FwNbConfigReq_t()
        req.cellId_pr.pres = True
        req.cellId_pr.cell_id = 10
        req.plmnId_pr.pres = True
        # Convert PLMN to ASCII character array of MCC and MNC digits
        # For 5 digit PLMN add \0 in the end, e.g., "00101\0"
        req.plmnId_pr.plmn_id = (ctypes.c_ubyte * 6).from_buffer_copy(
            bytearray(b"333333"),
        )

        print("************************* Sending ENB configuration Request")
        assert self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_CONFIG, req) == 0
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_CONFIG_CONFIRM.value
        res = response.cast(s1ap_types.FwNbConfigCfm_t)
        assert res.status == s1ap_types.CfgStatus.CFG_DONE.value

        print("************************* Sending S1-setup Request")
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
