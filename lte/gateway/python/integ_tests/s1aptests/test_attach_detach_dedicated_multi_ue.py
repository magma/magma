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


class TestAttachDetachDedicatedMultiUe(unittest.TestCase):
    """Dedicated bearer test with multiple UEs"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_detach(self):
        """attach/detach + dedicated bearer test with 4 UEs"""
        num_ues = 4
        ue_ids = []
        bearer_ids = []
        default_ips = []
        flow_lists = []
        self._s1ap_wrapper.configUEDevice(num_ues)

        for _ in range(num_ues):
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
            default_ips.append(ipaddress.ip_address(bytes(addr[:4])))

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Adding dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            flow_lists.append(self._spgw_util.create_default_ipv4_flows())
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach.esmInfo.epsBearerId,
                flow_lists[i],
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
            )
            print(
                "********************** Received activate dedicated EPS"
                " bearer context request",
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t,
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId,
            )
            bearer_ids.append(act_ded_ber_ctxt_req.bearerId)

        # Verify if flow rules are created
        for i in range(num_ues):
            dl_flow_rules = {
                default_ips[i]: [flow_lists[i]],
            }
            # 4 default + 4 dedicated bearer UL flows for 4 UEs
            num_ul_flows = 8
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

        print("Sleeping for 1 second")
        time.sleep(1)
        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Deleting dedicated bearer for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach.esmInfo.epsBearerId,
                bearer_ids[i],
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print(
                "********************** Received deactivate EPS bearer"
                " context request",
            )

            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, bearer_ids[i],
            )
        print("Sleeping for 2 seconds")
        time.sleep(2)
        # Verify if flow rules are deleted for dedicated bearers
        for i in range(num_ues):
            dl_flow_rules = {
                default_ips[i]: [],
            }
            # 4 default bearer UL flows for 4 UEs
            num_ul_flows = 4
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

        for ue in ue_ids:
            print("********************** Running UE detach for UE id ", ue)
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            )


if __name__ == "__main__":
    unittest.main()
