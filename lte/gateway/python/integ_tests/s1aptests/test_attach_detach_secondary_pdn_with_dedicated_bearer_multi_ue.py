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
import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestSecondaryPdnWithDedBearerMultiUe(unittest.TestCase):
    """Test secondary pdn creation with dedicated bearer for 4 UEs"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_secondary_pdn_with_dedicated_bearer_multi_ue(self):
        """Attach/detach + PDN Connectivity Requests + dedicated bearer for 4
        UEs
        """
        num_ues = 4
        ue_ids = []
        bearer_ids = []
        default_ips = []
        sec_ips = []
        flow_list = []

        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id
            # APN of the secondary PDN
            ims = {
                "apn_name": "ims",  # APN-name
                "qci": 5,  # qci
                "priority": 15,  # priority
                "pre_cap": 0,  # preemption-capability
                "pre_vul": 0,  # preemption-vulnerability
                "mbr_ul": 200000000,  # MBR UL
                "mbr_dl": 100000000,  # MBR DL
            }

            # APN list to be configured
            apn_list = [ims]

            self._s1ap_wrapper.configAPN(
                "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
            )

            print(
                "******************* Running End to End attach for UE id ",
                ue_id,
            )
            # Attach
            attach = self._s1ap_wrapper.s1_util.attach(
                ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            addr = attach.esmInfo.pAddr.addrInfo
            default_ips.append(ipaddress.ip_address(bytes(addr[:4])))

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(ue_id)

        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            # Send PDN Connectivity Request
            apn = "ims"
            ue_id = req.ue_id
            self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
            # Receive PDN CONN RSP/Activate default EPS bearer context request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
            )
            act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)
            addr = act_def_bearer_req.m.pdnInfo.pAddr.addrInfo
            sec_ips.append(ipaddress.ip_address(bytes(addr[:4])))

            print(
                "********************** Sending Activate default EPS bearer "
                "context accept for UE id ",
                ue_id,
            )
            print(
                "********************** Added IMS PDN with bearer id",
                act_def_bearer_req.m.pdnInfo.epsBearerId,
            )
            bearer_ids.append(act_def_bearer_req.m.pdnInfo.epsBearerId)

            print("********************* Sleeping for 2 seconds")
            time.sleep(2)
            # Add dedicated bearer to IMS PDN
            print("********************** Adding dedicated bearer to IMS PDN")
            flow_list.append(self._spgw_util.create_default_ipv4_flows())
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
                flow_list[i],
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

            print(
                "********************** Added dedicated bearer",
                act_ded_ber_ctxt_req.bearerId,
            )

        print("********************* Sleeping for 5 seconds")
        time.sleep(5)
        for i in range(num_ues):
            # Verify if flow rules are created
            dl_flow_rules = {
                default_ips[i]: [],
                sec_ips[i]: [flow_list[i]],
            }
            # 8 UL flows for default and secondary pdns +
            # 4 UL flows for dedicated bearers
            num_ul_flows = 12
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id
            print("******************* Deleting IMS PDN for ue", ue_id)
            # Send PDN Disconnect
            pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
            pdn_disconnect_req.ue_Id = ue_id
            pdn_disconnect_req.epsBearerId = bearer_ids[i]
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req,
            )

            # Receive UE_DEACTIVATE_BER_REQ
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print(
                "******************* Received deactivate eps bearer context"
                " request",
            )
            # Send DeactDedicatedBearerAccept
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                ue_id, bearer_ids[i],
            )

            print(
                "******************* Deleted IMS PDN with bearer ID",
                bearer_ids[i],
            )
        print("********************* Sleeping for 5 seconds")
        time.sleep(2)
        for i in range(num_ues):
            # Verify that flow rules for secondary pdn are deleted
            dl_flow_rules = {
                default_ips[i]: [],
            }
            # 4 UL flows for default bearers
            num_ul_flows = 4
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

        # Now detach the UE
        for ue in ue_ids:
            print(
                "******************** Running UE detach (switch-off) for ",
                "UE id ",
                ue,
            )
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
            )


if __name__ == "__main__":
    unittest.main()
