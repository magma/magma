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


class TestAttachDetachTwoPDNsWithTcpTraffic(unittest.TestCase):
    """Test two secondary pdn connections with tcp data"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_two_pdns_with_tcp_traffic(self):
        """Attach a single UE, send standalone PDN Connectivity
        Request, generate TCP traffic for each PDN session
        """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id

        # ims APN
        ims = {
            "apn_name": "ims",  # APN-name
            "qci": 5,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
        }
        apn_list = [ims]
        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
        )

        print(
            "************************* Running End to End"
            " attach for UE id ",
            ue_id,
        )
        # Attach
        self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        default_apn_ip = self._s1ap_wrapper._s1_util.get_ip(ue_id)

        time.sleep(2)
        # Send PDN Connectivity Request
        apn = "ims"
        print(
            "************************* Sending Standalone PDN "
            "CONNECTIVITY REQUEST for UE id ",
            ue_id,
            " APN ",
            apn,
        )
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        pdn_conn_rsp = response.cast(s1ap_types.uePdnConRsp_t)
        ims_addr = pdn_conn_rsp.m.pdnInfo.pAddr.addrInfo
        ims_ip = ipaddress.ip_address(bytes(ims_addr[:4]))

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify if flow rules are created
        # No dedicated bearers, so flowlist is empty
        dl_flow_rules = {
            default_apn_ip: [],
            ims_ip: [],
        }
        # 1 UL flow is created per bearer
        num_ul_flows = 2
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print(
            "************************* Running UE uplink (TCP) for UE id ",
            req.ue_id,
            " UE IP addr: ",
            default_apn_ip,
            " APN: oai_IPv4",
        )
        with self._s1ap_wrapper._trf_util.generate_traffic_test(
            [default_apn_ip], is_uplink=True, duration=5, is_udp=False,
        ) as test:
            test.verify()

        print("Sleeping for 5 seconds...")
        time.sleep(5)
        print(
            "************************* Running UE uplink (TCP) for UE id ",
            req.ue_id,
            " ue ip addr: ",
            ims_ip,
            " APN: ",
            apn,
        )
        with self._s1ap_wrapper._trf_util.generate_traffic_test(
            [ims_ip], is_uplink=True, duration=5, is_udp=False,
        ) as test:
            test.verify()

        print("Sleeping for 5 seconds...")
        time.sleep(5)
        # Send PDN Disconnect
        print("******************* Sending PDN Disconnect" " request")
        pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
        pdn_disconnect_req.ue_Id = ue_id
        pdn_disconnect_req.epsBearerId = pdn_conn_rsp.m.pdnInfo.epsBearerId
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req,
        )

        # Receive UE_DEACTIVATE_BER_REQ
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
        )

        print(
            "******************* Received deactivate EPS bearer context",
            " request",
        )
        deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
        self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
            req.ue_id, deactv_bearer_req.bearerId,
        )
        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify that flow rule is deleted for the secondary pdn
        dl_flow_rules = {
            default_apn_ip: [],
        }
        # 1 UL flow the default bearer
        num_ul_flows = 1
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
        )


if __name__ == "__main__":
    unittest.main()
