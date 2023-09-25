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
from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestIpv6NonNatDpDlTcp(unittest.TestCase):
    """Integration Test: TestIpv6NonNatDpDlTcp"""

    def __init__(self, method_name: str) -> None:
        """Initialize unittest class"""
        super().__init__(methodName=method_name)
        self.magma_utils = MagmadUtil(None)

    def setUp(self):
        """Initialize before test case execution"""
        self.magma_utils.disable_nat(ip_version=6)
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(ip_version=6)

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()
        self.magma_utils.enable_nat(ip_version=6)

    def test_ipv6_non_nat_dp_dl_tcp(self):
        """Basic attach/detach and DL TCP ipv6 data test with a single UE"""
        num_ues = 1

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

        wait_for_s1 = True
        ue_ips = ["fdee::"]
        apn_list = [magma_apn]

        self._s1ap_wrapper.configUEDevice(num_ues, [], ue_ips)

        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
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
        self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
            pdn_type=2,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Receive Router Advertisement message
        apn = "magma"
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_ROUTER_ADV_IND.value
        router_adv = response.cast(s1ap_types.ueRouterAdv_t)
        print(
            "********** Received Router Advertisement for APN-%s"
            " bearer id-%d" % (apn, router_adv.bearerId),
        )
        ipv6_addr = "".join([chr(i) for i in router_adv.ipv6Addr]).rstrip(
            "\x00",
        )
        print("********** UE IPv6 address: ", ipv6_addr)
        self._s1ap_wrapper.s1_util.update_ipv6_address(ue_id, ipv6_addr)

        print("Sleeping for 5 secs")
        time.sleep(5)
        print(
            "************************* Running UE downlink (TCP) for UE id ",
            req.ue_id,
        )
        self._s1ap_wrapper.configMtuSize(True)
        with self._s1ap_wrapper.configDownlinkTest(req, duration=1) as test:
            test.verify()
        self._s1ap_wrapper.configMtuSize(False)

        print(
            "************************* Running UE detach for UE id ",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            wait_for_s1,
        )


if __name__ == "__main__":
    unittest.main()
