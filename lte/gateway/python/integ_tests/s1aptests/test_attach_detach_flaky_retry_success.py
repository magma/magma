"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import ipaddress
import time
import unittest

import s1ap_types
from flaky import flaky
from integ_tests.s1aptests import s1ap_wrapper


class TestAttachDetachFlakyRetrySuccess(unittest.TestCase):
    """Integration Test: TestAttachDetachFlakyRetrySuccess"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    @flaky(max_runs=3, min_passes=3)
    def test_attach_detach_flaky_retry_success(self):
        """Basic attach/detach test for a single UE with flaky retry

        This testcase runs the basic attach-detach scenario multiple times even
        for successful case to validate TFW and S1APTester cleanup
        """
        num_ues = 2
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for ",
                "UE Id",
                req.ue_id,
            )
            # Now actually complete the attach
            attach = self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            addr = attach.esmInfo.pAddr.addrInfo
            default_ip = ipaddress.ip_address(bytes(addr[:4]))

            # Wait for EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            print("Waiting for 3 seconds for the flow rules creation")
            time.sleep(3)
            # Verify if flow rules are created
            # 1 UL flow for default bearer
            num_ul_flows = 1
            dl_flow_rules = {default_ip: []}
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows,
                dl_flow_rules,
            )

            # Now detach the UE
            print(
                "************************* Running UE detach for UE Id",
                req.ue_id,
            )
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                detach_type[i],
                wait_for_s1[i],
            )

            print("Waiting for 5 seconds for the flow rules deletion")
            time.sleep(5)
            # Verify that all UL/DL flows are deleted
            self._s1ap_wrapper.s1_util.verify_flow_rules_deletion()


if __name__ == "__main__":
    unittest.main()
