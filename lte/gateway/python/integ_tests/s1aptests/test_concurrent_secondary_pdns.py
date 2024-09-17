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


class TestConcurrentSecondaryPdns(unittest.TestCase):
    """Test concurrent secondary pdn session creation"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_concurrent_secondary_pdns(self):
        """Attach a single UE and send concurrent standalone PDN Connectivity
        Requests
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

        internet = {
            "apn_name": "internet",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 250000000,  # MBR UL
            "mbr_dl": 150000000,  # MBR DL
        }

        # APN list to be configured
        apn_list = [ims, internet]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
        )
        print(
            "************************* Running End to End attach for UE id =",
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
        default_ip = ipaddress.ip_address(bytes(addr[:4]))

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Send PDN Connectivity Request
        apn = "ims"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        apn = "internet"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
        act_def_bearer_req1 = response.cast(s1ap_types.uePdnConRsp_t)
        print(
            "************************* Received Activate default EPS bearer "
            "context request for UE id=%d, with bearer id=%d "
            % (
                act_def_bearer_req1.ue_Id,
                act_def_bearer_req1.m.pdnInfo.epsBearerId,
            ),
        )
        addr1 = act_def_bearer_req1.m.pdnInfo.pAddr.addrInfo
        sec_ip1 = ipaddress.ip_address(bytes(addr1[:4]))

        print(
            "************************* Sending Activate default EPS bearer "
            "context accept for UE id =",
            act_def_bearer_req1.ue_Id,
        )

        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
        act_def_bearer_req2 = response.cast(s1ap_types.uePdnConRsp_t)
        addr2 = act_def_bearer_req2.m.pdnInfo.pAddr.addrInfo
        sec_ip2 = ipaddress.ip_address(bytes(addr2[:4]))
        print(
            "************************* Received Activate default EPS bearer "
            "context request for UE id=%d, with bearer id=%d "
            % (
                act_def_bearer_req2.ue_Id,
                act_def_bearer_req2.m.pdnInfo.epsBearerId,
            ),
        )

        print(
            "************************* Sending Activate default EPS bearer "
            "context accept for UE id =",
            act_def_bearer_req2.ue_Id,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify if flow rules are created
        # No dedicated bearers, so flowlist is empty
        dl_flow_rules = {
            default_ip: [],
            sec_ip1: [],
            sec_ip2: [],
        }
        # 1 UL flow is created per bearer
        num_ul_flows = 3
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        # Send PDN Disconnect
        pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
        pdn_disconnect_req.ue_Id = act_def_bearer_req1.ue_Id
        pdn_disconnect_req.epsBearerId = (
            act_def_bearer_req1.m.pdnInfo.epsBearerId
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req,
        )

        # Receive UE_DEACTIVATE_BER_REQ
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value
        deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
        print(
            "******************* Received deactivate eps bearer context"
            " request for UE id=%d with bearer id=%d"
            % (deactv_bearer_req.ue_Id, deactv_bearer_req.bearerId),
        )
        # Send DeactDedicatedBearerAccept
        self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
            deactv_bearer_req.ue_Id, deactv_bearer_req.bearerId,
        )
        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify that flow rule is deleted for ims secondary pdn
        dl_flow_rules = {
            default_ip: [],
            sec_ip2: [],
        }
        # 1 UL flow is created per bearer
        num_ul_flows = 2
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id =",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            wait_for_s1_ctxt_release=False,
        )


if __name__ == "__main__":
    unittest.main()
