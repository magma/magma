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

import ipaddress
import unittest

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestAttachDetachNonNatDpUlTcp(unittest.TestCase):
    """Integration Test: TestAttachDetachNonNatDpUlTcp"""

    def setUp(self):
        """Initialize before test case execution"""
        self.magma_utils = MagmadUtil(None)
        self.magma_utils.disable_nat()
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()
        self.magma_utils.enable_nat()

    def test_attach_detach_non_nat_dp_ul_tcp(self):
        """Basic attach/detach test with a single UE"""
        num_ues = 1

        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, True, False]
        # third IP validates handling of invalid IP address.
        ue_ips = ["192.168.129.100"]
        self._s1ap_wrapper.configUEDevice(num_ues, [], ue_ips)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )

            # Now actually complete the attach
            attach = self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Validate assigned IP address.
            addr = attach.esmInfo.pAddr.addrInfo
            ue_ipv4 = ipaddress.ip_address(bytes(addr[:4]))
            assert ue_ipv4 == ipaddress.IPv4Address(ue_ips[i])

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            print(
                "************************* Running UE uplink (TCP) for UE id ",
                req.ue_id,
            )
            with self._s1ap_wrapper.configUplinkTest(req, duration=1) as test:
                test.verify()

            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                detach_type[i],
                wait_for_s1[i],
            )


if __name__ == "__main__":
    unittest.main()
