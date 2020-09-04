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
import time

import gpp_types
import s1ap_types
import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil
from integ_tests.s1aptests.s1ap_utils import SessionManagerUtil
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows


class TestAttachServiceWithMultiPdnsAndBearers(unittest.TestCase):
    SPGW_TABLE = 0
    GTP_PORT = 32768
    LOCAL_PORT = "LOCAL"

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def _verify_flow_rules(self, ueId, num_flows):
        MAX_NUM_RETRIES = 5
        datapath = get_datapath()

        # Check if UL and DL OVS flows are created
        print("************ Verifying flow rules")
        # UPLINK
        print("Checking for uplink flow")
        # try at least 5 times before failing as gateway
        # might take some time to install the flows in ovs
        for i in range(MAX_NUM_RETRIES):
            print("Get uplink flows: attempt ", i)
            uplink_flows = get_flows(
                datapath,
                {
                    "table_id": self.SPGW_TABLE,
                    "match": {"in_port": self.GTP_PORT},
                },
            )
            if len(uplink_flows) > num_flows:
                break
            time.sleep(5)  # sleep for 5 seconds before retrying

        assert len(uplink_flows) > num_flows, "Uplink flow missing for UE"
        self.assertIsNotNone(
            uplink_flows[0]["match"]["tunnel_id"],
            "Uplink flow missing tunnel id match",
        )

        # DOWNLINK
        print("Checking for downlink flow")
        ue_ip = str(self._s1ap_wrapper._s1_util.get_ip(ueId))
        # try at least 5 times before failing as gateway
        # might take some time to install the flows in ovs
        for i in range(MAX_NUM_RETRIES):
            print("Get downlink flows: attempt ", i)
            downlink_flows = get_flows(
                datapath,
                {
                    "table_id": self.SPGW_TABLE,
                    "match": {
                        "nw_dst": ue_ip,
                        "eth_type": 2048,
                        "in_port": self.LOCAL_PORT,
                    },
                },
            )
            if len(downlink_flows) > num_flows:
                break
            time.sleep(5)  # sleep for 5 seconds before retrying
        assert len(downlink_flows) > num_flows, "Downlink flow missing for UE"
        self.assertEqual(
            downlink_flows[0]["match"]["ipv4_dst"],
            ue_ip,
            "UE IP match missing from downlink flow",
        )

        actions = downlink_flows[0]["instructions"][0]["actions"]
        has_tunnel_action = any(
            action
            for action in actions
            if action["field"] == "tunnel_id"
            and action["type"] == "SET_FIELD"
        )
        self.assertTrue(
            has_tunnel_action, "Downlink flow missing set tunnel action"
        )

    def test_attach_service_with_multi_pdns_and_bearers(self):
        """
        Test with a single UE attach + add secondary PDN
        + add 2 dedicated bearers + UE context release + service request
        + detach"""
        self._s1ap_wrapper.configUEDevice(1)
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
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list
        )
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )

        # UL Flow description #1
        ulFlow1 = {
            "ipv4_dst": "192.168.129.42",  # IPv4 destination address
            "tcp_dst_port": 5002,  # TCP dest port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "UL",  # Direction
        }

        # UL Flow description #2
        ulFlow2 = {
            "ipv4_dst": "192.168.129.42",  # IPv4 destination address
            "tcp_dst_port": 5001,  # TCP dest port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "UL",  # Direction
        }

        # UL Flow description #3
        ulFlow3 = {
            "ipv4_dst": "192.168.129.64",  # IPv4 destination address
            "tcp_dst_port": 5003,  # TCP dest port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "UL",  # Direction
        }

        # UL Flow description #4
        ulFlow4 = {
            "ipv4_dst": "192.168.129.42",  # IPv4 destination address
            "tcp_dst_port": 5001,  # TCP dest port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "UL",  # Direction
        }
        # DL Flow description #1
        dlFlow1 = {
            "ipv4_src": "192.168.129.42",  # IPv4 source address
            "tcp_src_port": 5001,  # TCP source port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "DL",  # Direction
        }

        # DL Flow description #2
        dlFlow2 = {
            "ipv4_src": "192.168.129.64",  # IPv4 source address
            "tcp_src_port": 5002,  # TCP source port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "DL",  # Direction
        }

        # DL Flow description #3
        dlFlow3 = {
            "ipv4_src": "192.168.129.64",  # IPv4 source address
            "tcp_src_port": 5003,  # TCP source port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "DL",  # Direction
        }

        # DL Flow description #4
        dlFlow4 = {
            "ipv4_src": "192.168.129.42",  # IPv4 source address
            "tcp_src_port": 5001,  # TCP source port
            "ip_proto": "TCP",  # Protocol Type
            "direction": "DL",  # Direction
        }

        # Flow lists to be configured
        flow_list = [
            ulFlow1,
            ulFlow2,
            ulFlow3,
            ulFlow4,
            dlFlow1,
            dlFlow2,
            dlFlow3,
            dlFlow4,
        ]
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

        policy_id = "internet"

        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        print("Sleeping for 5 seconds")
        time.sleep(5)

        # Add dedicated bearer for default bearer 5
        print(
            "********************** Adding dedicated bearer to magma.ipv4"
            " PDN"
        )
        print(
            "********************** Sending RAR for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        self._sessionManager_util.create_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            policy_id,
            flow_list,
            qos,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        )
        act_ded_ber_req_oai_apn = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t
        )
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_req_oai_apn.bearerId
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Send PDN Connectivity Request
        apn = "ims"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
        )
        sec_pdn = response.cast(s1ap_types.uePdnConRsp_t)
        print(
            "********************** Sending Activate default EPS bearer "
            "context accept for UE id ",
            ue_id,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Add dedicated bearer to 2nd PDN
        print("********************** Adding dedicated bearer to ims PDN")
        self._spgw_util.create_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            sec_pdn.m.pdnInfo.epsBearerId,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        )
        act_ded_ber_req_ims_apn = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t
        )
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            req.ue_id, act_ded_ber_req_ims_apn.bearerId
        )
        print(
            "************* Added dedicated bearer",
            act_ded_ber_req_ims_apn.bearerId,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        print("*********** Moving UE to idle mode")
        print(
            "************* Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Verify if flow rules are created
        num_flows = 3
        self._verify_flow_rules(ue_id, num_flows)
        # Send UE context release request to move UE to idle mode
        req = s1ap_types.ueCntxtRelReq_t()
        req.ue_Id = ue_id
        req.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        )

        print(
            "************************* Sending Service request for UE id ",
            ue_id,
        )
        # Send service request to reconnect UE
        req = s1ap_types.ueserviceReq_t()
        req.ue_Id = ue_id
        req.ueMtmsi = s1ap_types.ueMtmsi_t()
        req.ueMtmsi.pres = False
        req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_SIGNALLING.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SERVICE_REQUEST, req
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)

        # Verify if flow rules are created
        self._verify_flow_rules(ue_id, num_flows)

        print("Sleeping for 5 seconds")
        time.sleep(5)
        print("************************* Running UE detach for UE id ", ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, True
        )


if __name__ == "__main__":
    unittest.main()
