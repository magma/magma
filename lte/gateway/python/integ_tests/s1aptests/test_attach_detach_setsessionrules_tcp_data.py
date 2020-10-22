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

from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SessionManagerUtil
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows
from lte.protos.policydb_pb2 import FlowMatch
from s1ap_utils import GTPBridgeUtils

class TestAttachDetachSetSessionRulesTcpData(unittest.TestCase):
    SPGW_TABLE = 0
    GTP_PORT = 32768
    LOCAL_PORT = "LOCAL"

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_setsessionrules_tcp_data(self):
        """ attach/detach + make set session rule calls to"""
        """ session manager twice + run UL and DL tcp traffic"""
        """ with a single UE """
        num_ues = 1
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)
        datapath = get_datapath()
        MAX_NUM_RETRIES = 5

        internet = {
            "apn_name": "internet",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 100000,  # MBR UL
            "mbr_dl": 150000,  # MBR DL
        }

        # APN list to be configured
        apn_list = [internet]

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            self._s1ap_wrapper.configAPN(
                "IMSI" + "".join([str(i) for i in req.imsi]), apn_list, default=False
            )
            print(
                "********************** Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # UL Flow description #1
            ulFlow1 = {
                "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
                "tcp_dst_port": 5001, # TCP UE Port
                "direction": FlowMatch.DOWNLINK,  # Direction
            }

            # DL Flow description #1
            dlFlow1 = {
                "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
                "tcp_dst_port": 7001, # TCP UE Port
                "direction": FlowMatch.DOWNLINK,  # Direction
            }

            ulFlow2 = {
                "ip_proto": FlowMatch.IPPROTO_UDP,  # Protocol Type
                "direction": FlowMatch.UPLINK,  # Direction
            }

            # DL Flow description #1
            dlFlow2 = {
                "ip_proto": FlowMatch.IPPROTO_UDP,  # Protocol Type
                "direction": FlowMatch.DOWNLINK,  # Direction
            }

            # Flow list to be configured
            flow_list = [
                ulFlow1,
                dlFlow1,
                ulFlow2,
                dlFlow2,
            ]

            # QoS
            qos = {
                "qci": 8,  # qci value [1 to 9]
                "priority": 0,  # Range [0-255]
                "max_req_bw_ul": 10000,  # MAX bw Uplink
                "max_req_bw_dl": 25000,  # MAX bw Downlink
                "gbr_ul": 1000,  # GBR Uplink
                "gbr_dl": 2000,  # GBR Downlink
                "arp_prio": 15,  # ARP priority
                "pre_cap": 1,  # pre-emption capability
                "pre_vul": 1,  # pre-emption vulnerability
            }

            policy_id = "tcp_udp_1"

            print("Sleeping for 5 seconds")
            time.sleep(5)
            print(
                "********************** Set Session Rule for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._sessionManager_util.send_SetSessionRules(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                policy_id,
                flow_list,
                qos,
            )

            # Receive Activate dedicated bearer request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )

            # Send Activate dedicated bearer accept
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId
            )

            policy_id = "tcp_udp_2"
            print("Sleeping for 5 seconds")
            time.sleep(5)
            print(
                "********************** Set Session Rule for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._sessionManager_util.send_SetSessionRules(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                policy_id,
                flow_list,
                qos,
            )
            # First rule is replaced by the second rule
            # Triggers a delete bearer followed by a create bearer request
            # Receive Deactivate dedicated bearer request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value
            )
            deactivate_ber_ctxt_req = response.cast(
                s1ap_types.UeDeActvBearCtxtReq_t
            )

            # Send Deactivate dedicated bearer accept
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactivate_ber_ctxt_req.bearerId
            )

            # Receive Activate dedicated bearer request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )

            # Send Activate dedicated bearer accept
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId
            )

            # Check if UL and DL OVS flows are created
            gtp_br_util = GTPBridgeUtils()
            gtp_port_no = gtp_br_util.get_gtp_port_no()
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
                        "match": {"in_port": gtp_port_no},
                    },
                )
                if len(uplink_flows) > 1:
                    break
                time.sleep(5)  # sleep for 5 seconds before retrying

            assert len(uplink_flows) > 1, "Uplink flow missing for UE"
            self.assertIsNotNone(
                uplink_flows[0]["match"]["tunnel_id"],
                "Uplink flow missing tunnel id match",
            )

            # DOWNLINK
            print("Checking for downlink flow")
            ue_ip = str(self._s1ap_wrapper._s1_util.get_ip(req.ue_id))
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
                if len(downlink_flows) > 1:
                    break
                time.sleep(5)  # sleep for 5 seconds before retrying

            assert len(downlink_flows) > 1, "Downlink flow missing for UE"
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

            print("Sleeping for 5 seconds")
            time.sleep(5)
            with self._s1ap_wrapper.configUplinkTest(req, duration=20, is_udp=False) as test:
                test.verify()

            print("Sleeping for 5 seconds")
            time.sleep(5)
            with self._s1ap_wrapper.configDownlinkTest(req, duration=1, is_udp=False) as test:
                test.verify()

            print("Sleeping for 5 seconds")
            time.sleep(5)

            print(
                "********************** Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i]
            )

            print("Checking that uplink/downlink flows were deleted")
            flows = get_flows(
                datapath, {"table_id": self.SPGW_TABLE, "priority": 0}
            )
            self.assertEqual(
                len(flows), 2, "There should only be 2 default table 0 flows"
            )


if __name__ == "__main__":
    unittest.main()
