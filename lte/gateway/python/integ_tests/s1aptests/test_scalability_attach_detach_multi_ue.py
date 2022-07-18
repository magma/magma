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

import time
import unittest

import s1ap_types
import s1ap_wrapper


class TestScalabilityAttachDetachMultiUe(unittest.TestCase):
    """Integration Test: TestScalabilityAttachDetachMultiUe"""

    def setUp(self):
        """Initialize before test case execution"""
        # The default IP pool for UE IP address allocation is configured in
        # s1ap_wrapper.py as 192.168.128.0/24. This mask value of 24 needs to
        # be changed as per the logic explained in s1ap_wrapper.py file to
        # allocate IP address for more than 243 UEs
        # Please follow https://github.com/magma/S1APTester for more details
        self.default_ip_block = s1ap_wrapper.TestWrapper.TEST_IP_BLOCK
        s1ap_wrapper.TestWrapper.TEST_IP_BLOCK = "192.168.128.0/17"
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()
        s1ap_wrapper.TestWrapper.TEST_IP_BLOCK = self.default_ip_block

    def test_scalability_attach_detach_multi_ue(self):
        """Basic attach and detach for 1024 UEs

        This testcase is a reference testcase for handling configurations for
        supporting large number of UEs.
        """
        ue_ids = []
        num_ues = 1024
        self._s1ap_wrapper.configUEDevice(num_ues)

        # The inactivity timers for UEs attached in the beginning starts getting
        # expired before all the UEs could be attached. Increasing UE inactivity
        # timer to 60 min (3600000 ms) to allow all the UEs to get attached and
        # detached properly
        print("Setting the inactivity timer value to 60 mins")
        config_data = s1ap_types.FwNbConfigReq_t()
        config_data.inactvTmrVal_pr.pres = True
        config_data.inactvTmrVal_pr.val = 3600000
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.ENB_INACTV_TMR_CFG,
            config_data,
        )
        time.sleep(0.5)

        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Calling attach for UE id:",
                req.ue_id,
            )
            self._s1ap_wrapper.s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            # Wait for EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        for ue in ue_ids:
            print("************************* Calling detach for UE id:", ue)
            self._s1ap_wrapper.s1_util.detach(
                ue,
                s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            )

        # Reset the inactivity timer value to default 2 mins (120000 ms)
        print("Resetting the inactivity timer value to default value (2 mins)")
        config_data = s1ap_types.FwNbConfigReq_t()
        config_data.inactvTmrVal_pr.pres = True
        config_data.inactvTmrVal_pr.val = 120000
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.ENB_INACTV_TMR_CFG,
            config_data,
        )
        time.sleep(0.5)


if __name__ == "__main__":
    unittest.main()
