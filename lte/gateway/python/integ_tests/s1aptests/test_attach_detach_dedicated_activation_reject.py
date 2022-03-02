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


class TestAttachDetachDedicatedReject(unittest.TestCase):
    """Dedicated bearer activation reject test with a
    single UE
    """

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_dedicated_activation_reject(self):
        """attach/detach + dedicated bearer activation reject test with a
        single UE
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
            time.sleep(5)
            print(
                "********************** Adding dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            # Create default flow list
            flow_list = self._spgw_util.create_default_ipv4_flows()
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach.esmInfo.epsBearerId,
                flow_list,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t,
            )
            ded_bearer_rej = s1ap_types.UeActDedBearCtxtRej_t()
            ded_bearer_rej.ue_Id = req.ue_id
            ded_bearer_rej.bearerId = act_ded_ber_ctxt_req.bearerId
            # Send Bearer Activation Reject
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ACT_DED_BER_REJ, ded_bearer_rej,
            )
            print(
                "Sent UE_ACT_DED_BER_REJ for bearer id",
                ded_bearer_rej.bearerId,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Dedicated bearer creation failed as UE sent UE_ACT_DED_BER_REJ
            # flow rule should not be created for the dedicated bearer
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
