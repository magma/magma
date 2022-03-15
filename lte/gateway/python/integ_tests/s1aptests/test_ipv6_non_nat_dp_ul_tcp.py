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
import time 
import unittest

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestAttachDetachNonNatDpUlTcp(unittest.TestCase):
    """Integration Test: TestAttachDetachNonNatDpUlTcp"""

    def __init__(self, method_name: str = ...) -> None:
        """Initialize unittest class"""
        super().__init__(methodName=method_name)
        self.magma_utils = MagmadUtil(None)

    def setUp(self):
        """Initialize before test case execution"""
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
        magma_apn = {
            "apn_name": "magma",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 1,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
            "pdn_type": 1,  # PDN Type 0-IPv4,1-IPv6,2-IPv4v6
        }

        wait_for_s1 = [True, True, False]
        # third IP validates handling of invalid IP address.
        ue_ips = ["fdee::"]
        apn_list = [magma_apn]

        self._s1ap_wrapper.configUEDevice(num_ues, [], ue_ips)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id
            print("******* Iteration ******", i)
            print(
                "************************* Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )

            self._s1ap_wrapper.configAPN(
                "IMSI" + "".join([str(j) for j in req.imsi]),
                apn_list,
                default=False,
            )

            # Now actually complete the attach
            attach = self._s1ap_wrapper.s1_util.attach(
             ue_id,
             s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
             s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
             s1ap_types.ueAttachAccept_t,
             pdn_type=2,
            )

            """addr = attach.esmInfo.pAddr.addrInfo
            ue_ipv4 = ipaddress.ip_address(bytes(addr[:4]))
            self.assertEqual(ue_ipv4, ipaddress.IPv4Address(ue_ips[i]))"""

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Delay to ensure S1APTester sends attach complete before sending UE
            # context release
            time.sleep(5)
            # Receive Router Advertisement message
            apn = "magma"
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ROUTER_ADV_IND.value,
            )
            router_adv = response.cast(s1ap_types.ueRouterAdv_t)
            print(
                "********** Received Router Advertisement for APN-%s"
                " bearer id-%d" % (apn, router_adv.bearerId),
            )
            ipv6_addr = "".join([chr(i) for i in router_adv.ipv6Addr]).rstrip(
                "\x00",
            )
            print("********** UE IPv6 address: ", ipv6_addr)
            default_ipv6 = ipaddress.ip_address(ipv6_addr)
            self._s1ap_wrapper.s1_util.update_ipv6_address(ue_id, ipv6_addr)

            print("Sleeping for 20 secs")
            time.sleep(20)
            print(
                "************************* Running UE uplink (TCP) for UE id ",
                req.ue_id,
            )
            with self._s1ap_wrapper.configUplinkTest(req, duration=10) as test:
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
