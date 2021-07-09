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

import time
import unittest

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows
from integ_tests.s1aptests.s1ap_utils import (
    GTPBridgeUtils,
    HeaderEnrichmentUtils,
    SessionManagerUtil,
    SpgwUtil,
)
from lte.protos.policydb_pb2 import FlowMatch, HeaderEnrichment


class TestAttachDetachRarTcpDataWithHE(unittest.TestCase):
    SPGW_TABLE = 0
    LOCAL_PORT = "LOCAL"

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_rar_tcp_data_with_he(self):
        """ Attach/detach + send ReAuth Req to session manager with a
        single UE along with Header enrichment.
        This test validates that same data test wotks with HE policy
        applied.
        Header enrichment should be as transparent as possible.
        """
        num_ues = 1
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)
        datapath = get_datapath()
        MAX_NUM_RETRIES = 5
        gtp_br_util = GTPBridgeUtils()
        GTP_PORT = gtp_br_util.get_gtp_port_no()
        utils = HeaderEnrichmentUtils()

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
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
                "tcp_dst_port": 5004,  # TCP dest port
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
                "tcp_src_port": 5004,  # TCP source port
                "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
                "direction": FlowMatch.DOWNLINK,  # Direction
            }

            # Flow list to be configured
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
                "qci": 5,  # qci value [1 to 9]
                "priority": 0,  # Range [0-255]
                "max_req_bw_ul": 10000000,  # MAX bw Uplink
                "max_req_bw_dl": 15000000,  # MAX bw Downlink
                "gbr_ul": 1000000,  # GBR Uplink
                "gbr_dl": 2000000,  # GBR Downlink
                "arp_prio": 15,  # ARP priority
                "pre_cap": 1,  # pre-emption capability
                "pre_vul": 1,  # pre-emption vulnerability
            }

            policy_id = "ims-voice"

            print("Sleeping for 5 seconds")
            time.sleep(5)
            imsi = "IMSI" + "".join([str(i) for i in req.imsi])

            print(
                "********************** Sending RAR for ", imsi,
            )
            he_domain1 = "192.168.128.1"
            assert utils.he_count_record_of_imsi_to_domain(
                imsi, he_domain1,
            ) == 0

            self._sessionManager_util.send_ReAuthRequest(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                policy_id,
                flow_list,
                qos,
                he_urls=HeaderEnrichment(urls=[he_domain1]),
            )

            # Receive Activate dedicated bearer request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
            )
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Send Activate dedicated bearer accept
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId,
            )

            # Check if UL and DL OVS flows are created
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
                        "match": {"in_port": GTP_PORT},
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
                has_tunnel_action, "Downlink flow missing set tunnel action",
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            with self._s1ap_wrapper.configUplinkTest(req, duration=1) as test:
                test.verify()
            assert utils.he_count_record_of_imsi_to_domain(
                imsi, he_domain1,
            ) == 1

            print(
                "********************** Deleting dedicated bearer for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), 5, 6,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print("******************* Received deactivate eps bearer context")

            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)

            print(
                "********************** Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i],
            )

            time.sleep(20)
            assert utils.he_count_record_of_imsi_to_domain(
                imsi, he_domain1,
            ) == 0

        # Verify that all UL/DL flows are deleted
        self._s1ap_wrapper.s1_util.verify_flow_rules_deletion()


if __name__ == "__main__":
    unittest.main()
