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


class TestAttachDetachNwTriggeredDeleteSecondaryPdn(unittest.TestCase):
    """Test network triggered secondary pdn deletion"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_nw_triggered_delete_secondary_pdn(self):
        """Attach a single UE + add secondary PDN + add dedicated bearer
        to the secondary pdn + delete the secondary pdn + detach
        """
        num_ue = 1

        self._s1ap_wrapper.configUEDevice(num_ue)
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
            "******************* Running End to End attach for UE id ", ue_id,
        )
        # Attach
        attach = self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        addr = attach.esmInfo.pAddr.addrInfo
        default_ip = ipaddress.ip_address(bytes(addr[:4]))

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Send PDN Connectivity Request
        apn = "ims"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)

        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        act_sec_pdn = response.cast(s1ap_types.uePdnConRsp_t)
        addr = act_sec_pdn.m.pdnInfo.pAddr.addrInfo
        sec_ip = ipaddress.ip_address(bytes(addr[:4]))

        print(
            "******************* Sending Activate default EPS bearer "
            "context accept for UE id ",
            ue_id,
        )

        # Add dedicated bearer to IMS PDN
        print(
            "******************* Adding dedicated bearer to IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        # Create default flow list
        flow_list = self._spgw_util.create_default_ipv4_flows()
        self._spgw_util.create_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            act_sec_pdn.m.pdnInfo.epsBearerId,
            flow_list,
        )
        # Receive Activate dedicated EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
        )
        act_ded_ber_ctxt_req = response.cast(s1ap_types.UeActDedBearCtxtReq_t)
        # Send Activate dedicated EPS bearer context accept
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_ctxt_req.bearerId,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify if flow rules are created
        dl_flow_rules = {
            default_ip: [],
            sec_ip: [flow_list],
        }
        # 2 UL flows for default and secondary pdns + 1 for dedicated bearer
        num_ul_flows = 3
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )
        print(
            "******************* Deleting default bearer for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        # Delete secondary pdn
        self._spgw_util.delete_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            act_sec_pdn.m.pdnInfo.epsBearerId,
            act_sec_pdn.m.pdnInfo.epsBearerId,
        )

        # Receive UE_DEACTIVATE_BER_REQ
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
        )

        deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
        print(
            "******************* Received deactivate eps bearer context"
            " request",
        )
        print(
            "******************* Sending deactivate eps bearer context"
            " accept",
        )

        # Send Deactivate EPS bearer context accept
        self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
            ue_id, deactv_bearer_req.bearerId,
        )
        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify if flow rules are deleted for secondary pdn
        dl_flow_rules = {
            default_ip: [],
        }
        # 1 UL flow for default pdn
        num_ul_flows = 1
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )
        print(
            "******************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
        )


if __name__ == "__main__":
    unittest.main()
