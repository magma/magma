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


class TestDedicatedBearerActivationIdleMode(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_dedicated_bearer_activation_idle_mode(self):
        """
        Test with a single UE attach + UE context release +
        trigger dedicated bearer activation
        + Page the UE + detach"""

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

        # Flow list to be configured
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
            "qci": 5,  # qci value [1 to 9]
            "priority": 5,  # Range [0-255]
            "max_req_bw_ul": 10000000,  # MAX bw Uplink
            "max_req_bw_dl": 15000000,  # MAX bw Downlink
            "gbr_ul": 1000000,  # GBR Uplink
            "gbr_dl": 2000000,  # GBR Downlink
            "arp_prio": 15,  # ARP priority
            "pre_cap": 1,  # pre-emption capability
            "pre_vul": 1,  # pre-emption vulnerability
        }
        qos2 = {
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

        policy_id1 = "ims-voice"
        policy_id2 = "internet"

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        imsi = req.imsi
        print("*********** Running End to End attach for UE id ", ue_id)
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
        print("*********** Sleeping for 5 seconds")
        time.sleep(5)

        # Move UE to idle mode
        print("*********** Moving UE to idle mode")
        print(
            "*********** Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
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
        ip_list = [default_ip]
        self._s1ap_wrapper.s1_util.verify_paging_flow_rules(ip_list)

        print("*********** Sleeping for 5 seconds")
        time.sleep(5)
        # Add dedicated bearers for default bearer 5

        print(
            "************* Adding dedicated bearers to magma.ipv4"
            " PDN in idle mode for UE",
            ue_id,
        )
        print(
            "********************** Sending 1st RAR for IMSI",
            "".join([str(i) for i in imsi]),
        )
        self._sessionManager_util.send_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in imsi]),
            policy_id1,
            flow_list1,
            qos1,
        )

        # Add 2nd dedicated bearer for default bearer 5
        print(
            "********************** Sending 2nd RAR for IMSI",
            "".join([str(i) for i in imsi]),
        )
        self._sessionManager_util.send_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in imsi]),
            policy_id2,
            flow_list2,
            qos2,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PAGING_IND.value,
        )

        print("*********** Received Paging for UE id ", ue_id)
        print("*********** Sending Service request for UE id ", ue_id)
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

        print("*********** Received ICS Request for UE id ", ue_id)
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
        )
        act_ded_ber_req_oai_apn1 = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t,
        )
        print(
            "*********** Received Activate dedicated EPS bearer"
            " req for bearer ",
            act_ded_ber_req_oai_apn1.bearerId,
        )
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            ue_id, act_ded_ber_req_oai_apn1.bearerId,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
        )
        act_ded_ber_req_oai_apn2 = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t,
        )
        print(
            "*********** Received Activate dedicated EPS bearer"
            " req for bearer ",
            act_ded_ber_req_oai_apn2.bearerId,
        )

        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            ue_id, act_ded_ber_req_oai_apn2.bearerId,
        )

        # Verify if flow rules are created
        dl_flow_rules = {default_ip: [flow_list1, flow_list2]}
        # 1 UL flow is created per bearer
        num_ul_flows = 3
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print("*********** Sleeping for 5 seconds")
        time.sleep(5)
        print("*********** Running UE detach for UE id ", ue_id)

        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
