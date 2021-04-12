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
import time
import ipaddress

from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestAttachDetachMultipleDedicated(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach(self):
        """ attach/detach + multiple dedicated bearer test with a single UE """
        num_dedicated_bearers = 3
        bearer_ids = []
        flow_lists = []
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "********************** Running End to End attach for UE id",
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

        print("Sleeping for 2 seconds")
        time.sleep(2)
        for i in range(num_dedicated_bearers):
            print(
                "********************** Adding dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            flow_lists.append(self._spgw_util.create_default_flows())
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), attach.esmInfo.epsBearerId,
                flow_lists[i]
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId
            )
            bearer_ids.append(act_ded_ber_ctxt_req.bearerId)
            print(
                "********************** Added dedicated bearer with ",
                "bearer id",
                act_ded_ber_ctxt_req.bearerId,
            )

        print("Sleeping for 2 seconds")
        time.sleep(2)
        # flow_lists for 3 dedicated bearers
        dl_flow_rules = {
            default_ip: [flow_lists[0],flow_lists[1],flow_lists[2]],
        }
        # 1 default UL flow + 3 dedicated bearer UL flows
        num_ul_flows = 4
        # Verify if flow rules are created
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules
        )


        for i in range(num_dedicated_bearers):
            print(
                "********************** Deleting dedicated bearer for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), attach.esmInfo.epsBearerId, bearer_ids[i]
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print("******************* Received deactivate eps bearer context")

            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, bearer_ids[i]
            )

            print(
                "********************** Deleted dedicated bearer with "
                "bearer id",
                bearer_ids[i],
            )

        print("Sleeping for 2 seconds")
        time.sleep(2)
        print("********************** Running UE detach for UE id ", req.ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value
        )


if __name__ == "__main__":
    unittest.main()
