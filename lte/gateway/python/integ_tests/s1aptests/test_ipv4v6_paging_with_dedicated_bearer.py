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

import gpp_types
import s1ap_types
import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SessionManagerUtil, SpgwUtil


class TestIpv4v6PagingWithDedicatedBearer(unittest.TestCase):
    """Test ipv4v6 paging"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()
        self._sessionManager_util = SessionManagerUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_ipv4v6_paging_with_dedicated_bearer(self):
        """ IPv4v6 Attach, add dedicated bearer, ue context release,
        paging request """
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        # Default apn over-write
        magma_apn = {
            "apn_name": "magma",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 1,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
            "pdn_type": 2,  # PDN Type 0-IPv4,1-IPv6,2-IPv4v6
        }
        ue_ctxt_rel = False
        apn_list = [magma_apn]
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print(
            "********** Running End to End attach for UE id ", ue_id,
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
            pdn_type=3,
        )
        addr = attach.esmInfo.pAddr.addrInfo
        default_ipv4 = ipaddress.ip_address(bytes(addr[8:12]))
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        time.sleep(5)
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
        default_ipv6 = ipaddress.ip_address(ipv6_addr)
        self._s1ap_wrapper.s1_util.update_ipv6_address(ue_id, ipv6_addr)

        print("********************** Adding dedicated bearer")
        print(
            "********************** Sending RAR for IMSI",
            "".join([str(i) for i in req.imsi]),
        )

        # QoS
        qos = {
            "qci": 1,  # qci value [1 to 9]
            "priority": 1,  # Range [0-255]
            "max_req_bw_ul": 10000000,  # MAX bw Uplink
            "max_req_bw_dl": 15000000,  # MAX bw Downlink
            "gbr_ul": 1000000,  # GBR Uplink
            "gbr_dl": 2000000,  # GBR Downlink
            "arp_prio": 1,  # ARP priority
            "pre_cap": 1,  # pre-emption capability
            "pre_vul": 1,  # pre-emption vulnerability
        }

        flow_list = self._spgw_util.create_default_ipv6_flows()
        policy_id = "magma"
        self._sessionManager_util.send_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            policy_id,
            flow_list,
            qos,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        act_ded_ber_req = response.cast(s1ap_types.UeActDedBearCtxtReq_t)
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_req.bearerId,
        )
        print(
            "************* Added dedicated bearer", act_ded_ber_req.bearerId,
        )

        # Sleep before verifying flows
        print("********** Sleeping for 5 seconds")
        time.sleep(5)
        # Verify flow rules
        num_ul_flows = 2
        dl_flow_rules = {
            default_ipv4: [flow_list],
            default_ipv6: [flow_list],
        }
        # Verify if flow rules are created
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print(
            "********** Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        ue_cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
        ue_cntxt_rel_req.ue_Id = ue_id
        ue_cntxt_rel_req.cause.causeVal = (
            gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, ue_cntxt_rel_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        print("********** UE moved to idle mode")

        print("********** Sleeping for 5 seconds")
        time.sleep(5)

        # Verify paging rules
        ip_list = [default_ipv4, default_ipv6]
        self._s1ap_wrapper.s1_util.verify_paging_flow_rules(ip_list)

        print(
            "********** Running UE downlink (UDP) for UE id ", ue_id,
        )
        self._s1ap_wrapper.s1_util.run_ipv6_data(default_ipv6)
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_PAGING_IND.value
        print("********** Received UE_PAGING_IND")
        # Send service request to reconnect UE
        ser_req = s1ap_types.ueserviceReq_t()
        ser_req.ue_Id = ue_id
        ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
        ser_req.ueMtmsi.pres = False
        ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MT_ACCESS.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SERVICE_REQUEST, ser_req,
        )
        print("********** Sent UE_SERVICE_REQUEST")
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value
        print("********** Received INT_CTX_SETUP_IND")

        print("********** Sleeping for 5 seconds")
        time.sleep(5)
        # Verify flow rules
        num_ul_flows = 2
        dl_flow_rules = {
            default_ipv4: [flow_list],
            default_ipv6: [flow_list],
        }
        # Verify if flow rules are created
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            ue_ctxt_rel,
        )


if __name__ == "__main__":
    unittest.main()
