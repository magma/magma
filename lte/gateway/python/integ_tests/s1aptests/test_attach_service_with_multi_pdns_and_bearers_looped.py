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
from integ_tests.s1aptests.s1ap_utils import SessionManagerUtil
from lte.protos.policydb_pb2 import FlowMatch


class TestAttachServiceWithMultiPdnsAndBearersLooped(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_service_with_multi_pdns_and_bearers_looped(self):
        """
        Test with a single UE attach + add secondary PDN
        + add 2 dedicated bearers
        + (UE context release + service request) repeat 3 times
        + detach"""
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        loop = 3
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
            "************************* Running End to End attach for UE id ",
            ue_id,
        )

        # UL Flow description #1
        ulFlow1 = {
            "ipv4_dst": "192.168.129.42",  # IPv4 destination address
            "tcp_dst_port": 5002,  # TCP dest port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.UPLINK,  # Direction
        }

        # UL Flow description #2
        ulFlow2 = {
            "ipv4_dst": "192.168.129.42",  # IPv4 destination address
            "tcp_dst_port": 5001,  # TCP dest port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.UPLINK,  # Direction
        }

        # UL Flow description #3
        ulFlow3 = {
            "ipv4_dst": "192.168.129.64",  # IPv4 destination address
            "tcp_dst_port": 5003,  # TCP dest port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.UPLINK,  # Direction
        }

        # UL Flow description #4
        ulFlow4 = {
            "ipv4_dst": "192.168.129.42",  # IPv4 destination address
            "tcp_dst_port": 5001,  # TCP dest port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.UPLINK,  # Direction
        }
        # DL Flow description #1
        dlFlow1 = {
            "ipv4_src": "192.168.129.42",  # IPv4 source address
            "tcp_src_port": 5001,  # TCP source port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.DOWNLINK,  # Direction
        }

        # DL Flow description #2
        dlFlow2 = {
            "ipv4_src": "192.168.129.64",  # IPv4 source address
            "tcp_src_port": 5002,  # TCP source port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.DOWNLINK,  # Direction
        }

        # DL Flow description #3
        dlFlow3 = {
            "ipv4_src": "192.168.129.64",  # IPv4 source address
            "tcp_src_port": 5003,  # TCP source port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.DOWNLINK,  # Direction
        }

        # DL Flow description #4
        dlFlow4 = {
            "ipv4_src": "192.168.129.42",  # IPv4 source address
            "tcp_src_port": 5001,  # TCP source port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.DOWNLINK,  # Direction
        }

        # Flow lists to be configured
        flow_list1 = [
            ulFlow1,
            ulFlow2,
            ulFlow3,
            dlFlow1,
            dlFlow2,
            dlFlow3,
        ]

        flow_list2 = [
            ulFlow4,
            dlFlow4,
        ]

        # QoS
        qos1 = {
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

        qos2 = {
            "qci": 2,  # qci value [1 to 9]
            "priority": 5,  # Range [0-255]
            "max_req_bw_ul": 10000000,  # MAX bw Uplink
            "max_req_bw_dl": 15000000,  # MAX bw Downlink
            "gbr_ul": 1000000,  # GBR Uplink
            "gbr_dl": 2000000,  # GBR Downlink
            "arp_prio": 1,  # ARP priority
            "pre_cap": 1,  # pre-emption capability
            "pre_vul": 1,  # pre-emption vulnerability
        }

        policy_id1 = "internet"
        policy_id2 = "ims"

        # Now actually complete the attach
        attach = self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        addr = attach.esmInfo.pAddr.addrInfo
        default_ip = ipaddress.ip_address(bytes(addr[:4]))

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        print("Sleeping for 5 seconds")
        time.sleep(5)

        # Add dedicated bearer for default bearer 5
        print(
            "********************** Adding dedicated bearer to magma.ipv4"
            " PDN",
        )
        print(
            "********************** Sending RAR for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        self._sessionManager_util.send_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            policy_id1,
            flow_list1,
            qos1,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
        )
        act_ded_ber_req_oai_apn = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t,
        )
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_req_oai_apn.bearerId,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Send PDN Connectivity Request
        apn = "ims"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)
        addr = act_def_bearer_req.m.pdnInfo.pAddr.addrInfo
        sec_ip = ipaddress.ip_address(bytes(addr[:4]))

        print(
            "********************** Sending Activate default EPS bearer "
            "context accept for UE id ",
            ue_id,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Add dedicated bearer to 2nd PDN
        print("********************** Adding dedicated bearer to ims PDN")
        print(
            "********************** Sending RAR for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        self._sessionManager_util.send_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            policy_id2,
            flow_list2,
            qos2,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
        )
        act_ded_ber_req_ims_apn = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t,
        )
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_req_ims_apn.bearerId,
        )
        print(
            "************* Added dedicated bearer",
            act_ded_ber_req_ims_apn.bearerId,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)

        dl_flow_rules = {
            default_ip: [flow_list1],
            sec_ip: [flow_list2],
        }
        # 1 UL flow is created per bearer
        num_ul_flows = 4
        # Verify if flow rules are created
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print("*********** Moving UE to idle mode")
        print(
            "************* Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        for j in range(loop):
            print("********************* Iteration - ", j + 1)
            # Send UE context release request to move UE to idle mode
            rel_req = s1ap_types.ueCntxtRelReq_t()
            rel_req.ue_Id = ue_id
            rel_req.cause.causeVal = (
                gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
            )
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, rel_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )

            # Verify if paging flow rules are created
            ip_list = [default_ip, sec_ip]
            self._s1ap_wrapper.s1_util.verify_paging_flow_rules(ip_list)

            print(
                "************************* Sending Service request for UE id ",
                ue_id,
            )
            # Send service request to reconnect UE
            ser_req = s1ap_types.ueserviceReq_t()
            ser_req.ue_Id = ue_id
            ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
            ser_req.ueMtmsi.pres = False
            ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_SIGNALLING.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST, ser_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)

            # Verify if flow rules are created
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        print("************************* Running UE detach for UE id ", ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
