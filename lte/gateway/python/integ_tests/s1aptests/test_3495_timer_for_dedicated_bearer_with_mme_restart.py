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
from s1ap_utils import MagmadUtil


class Test3495TimerForDedicatedBearerWithMmeRestart(unittest.TestCase):
    """Test case validates the functionality of 3495 timer for
    dedicated bearer while MME restarts
    """

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_3495_timer_for_dedicated_bearer_with_mme_restart(self):
        """Test case validates the functionality of 3495 timer for
        dedicated bearer while MME restarts
        Step1: UE attaches to network
        Step2: Creates a dedicated bearer
        Step3: Initiates dedicated bearer deletion, as part of which mme sends
        deactivate EPS bearer context request and starts 3495 timer
        Step4: TFW shall not respond to first Deactivate EPS bearer context
        request message
        Step5: Send command to Magma to restart mme service
        Step6: TFW shall receive re-transmitted Deactivate EPS bearer context
        request message and send Deactivate EPS bearer Context Accept
        Step7: TFW shall initiate Detach procedure.
        """
        num_ues = 1
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
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Verify if flow rules are created
            dl_flow_rules = {
                default_ip: [flow_list],
            }
            # 1 UL flow for default bearer + 1 for dedicated bearer
            num_ul_flows = 2
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

            print(
                "********************** Deleting dedicated bearer for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach.esmInfo.epsBearerId,
                act_ded_ber_ctxt_req.bearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print("******************* Received deactivate eps bearer context")
            # Do not send deactivate eps bearer context accept
            print(
                "************************* Restarting MME service on",
                "gateway",
            )
            self._s1ap_wrapper.magmad_util.restart_services(["mme"])

            for j in range(30):
                print("Waiting for", j, "seconds")
                time.sleep(1)

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )
            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId,
            )
            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Verify if flow rules are deleted for dedicated bearer
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
                req.ue_id,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                False,
            )


if __name__ == "__main__":
    unittest.main()
