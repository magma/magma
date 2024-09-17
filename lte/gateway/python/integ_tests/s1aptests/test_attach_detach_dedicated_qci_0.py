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
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestAttachDetachDedicatedQci0(unittest.TestCase):
    """Dedicated bearer test with qci 0"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_dedicated_qci_0(self):
        """Test attach + create dedicated bearer with QCI, 0 +
        erab_setup_failed_indication + detach, with a single UE
        """
        num_ues = 1
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Running End to End attach for ",
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
            addr = attach.esmInfo.pAddr.addrInfo
            default_ip = ipaddress.ip_address(bytes(addr[:4]))

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            print("Sleeping for 5 seconds")
            time.sleep(2)
            imsi = "".join([str(i) for i in req.imsi])
            print(
                "********************** Adding dedicated bearer to IMSI", imsi,
            )
            # Create default flow list
            flow_list = self._spgw_util.create_default_ipv4_flows()
            self._spgw_util.create_bearer(
                "IMSI" + imsi,
                attach.esmInfo.epsBearerId,
                flow_list,
                qci_val=0,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_FW_ERAB_SETUP_REQ_FAILED_FOR_ERABS.value
            erab_setup_failed_for_bearers = response.cast(
                s1ap_types.FwErabSetupFailedTosetup,
            )
            print(
                "*** Received UE_FW_ERAB_SETUP_REQ_FAILED_FOR_ERABS for "
                "bearer-id:",
                erab_setup_failed_for_bearers.failedErablist[0].erabId,
                end=" ",
            )
            print(
                " with qci:",
                erab_setup_failed_for_bearers.failedErablist[0].qci,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Verify that flow rules are created only for default bearer
            dl_flow_rules = {
                default_ip: [],
            }
            # 1 UL flow for default bearer
            num_ul_flows = 1
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

            print(
                "********************** Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i],
            )


if __name__ == "__main__":
    unittest.main()
