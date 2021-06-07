"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.
This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import ipaddress
import time
import unittest

import s1ap_types
import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SessionManagerUtil
from lte.protos.policydb_pb2 import FlowMatch


class TestIPv4v6SecondaryPdnWithDedBearerMultiUe(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_ipv4v6_secondary_pdn_with_ded_bearer_multi_ue(self):
        """Attach a single UE + add a secondary pdn with
        IPv4v6 + add dedicated bearer + detach
        Repeat for 4 UEs"""
        num_ues = 4
        ue_ids = []
        default_ip = []
        sec_ip_ipv4 = []
        sec_ip_ipv6 = []

        # APN of the secondary PDN
        ims_apn = {
            "apn_name": "ims",  # APN-name
            "qci": 5,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
            "pdn_type": 2,  # PDN Type 0-IPv4,1-IPv6,2-IPv4v6
        }
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
            "ipv6_dst": "5e90:db7b:b18e::1556",  # IPv6 destination address
            "tcp_dst_port": 5003,  # TCP dest port
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
            "ipv6_src": "6e31:1a95:1e7c::df1",  # IPv6 source address
            "tcp_src_port": 5003,  # TCP source port
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.DOWNLINK,  # Direction
        }

        # Flow lists to be configured
        flow_list = [
            ulFlow1,
            ulFlow2,
            ulFlow3,
            dlFlow1,
            dlFlow2,
            dlFlow3,
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

        policy_id = "ims"

        self._s1ap_wrapper.configUEDevice(num_ues)
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id

            apn_list = [ims_apn]
            self._s1ap_wrapper.configAPN(
                "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
            )
            print(
                "*********************** Running End to End attach for UE id ",
                ue_id,
            )

            print("***** Sleeping for 5 seconds")
            time.sleep(5)
            # Attach
            attach = self._s1ap_wrapper.s1_util.attach(
                ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            addr = attach.esmInfo.pAddr.addrInfo
            default_ip.append(ipaddress.ip_address(bytes(addr[:4])))
            ue_ids.append(ue_id)

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

        print("***** Sleeping for 5 seconds")
        time.sleep(5)
        self._s1ap_wrapper._ue_idx = 0
        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_id = req.ue_id

            apn = "ims"
            # PDN Type 2 = IPv6, 3 = IPv4v6
            pdn_type = 3
            # Send PDN Connectivity Request
            self._s1ap_wrapper.sendPdnConnectivityReq(
                ue_id, apn, pdn_type=pdn_type,
            )
            # Receive PDN CONN RSP/Activate default EPS bearer context request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
            )
            act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)

            addr = act_def_bearer_req.m.pdnInfo.pAddr.addrInfo
            sec_ip_ipv4.append(ipaddress.ip_address(bytes(addr[8:12])))

            print(
                "********************** Sending Activate default EPS bearer "
                "context accept for APN-%s, UE id-%d" % (apn, ue_id),
            )
            print(
                "********************** Added default bearer for apn-%s,"
                " bearer id-%d, pdn type-%d"
                % (
                    apn,
                    act_def_bearer_req.m.pdnInfo.epsBearerId,
                    pdn_type,
                ),
            )

            # Receive Router Advertisement message
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ROUTER_ADV_IND.value,
            )
            routerAdv = response.cast(s1ap_types.ueRouterAdv_t)
            print(
                "******************* Received Router Advertisement for APN-%s"
                " ,bearer id-%d" % (apn, routerAdv.bearerId),
            )

            ipv6_addr = "".join([chr(i) for i in routerAdv.ipv6Addr]).rstrip(
                "\x00",
            )
            print("******* UE IPv6 address: ", ipv6_addr)
            sec_ip_ipv6.append(ipaddress.ip_address(ipv6_addr))

            print("***** Sleeping for 5 seconds")
            time.sleep(5)

            # Add dedicated bearer
            print("********************** Adding dedicated bearer to ims PDN")
            print(
                "********************** Sending RAR for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._sessionManager_util.send_ReAuthRequest(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                policy_id,
                flow_list,
                qos,
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

        print("***** Sleeping for 10 seconds")
        time.sleep(10)
        # ipv4 default pdn + ipv4v6(ims) pdn +
        # dedicated bearer for ims pdn for 4 UEs
        num_ul_flows = 12
        for i in range(num_ues):
            dl_flow_rules = {
                default_ip[i]: [],
                sec_ip_ipv4[i]: [flow_list],
                sec_ip_ipv6[i]: [flow_list],
            }
            # Verify if flow rules are created
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

        print("***** Sleeping for 5 seconds")
        time.sleep(5)
        for ue in ue_ids:
            print(
                "******************* Running UE detach (switch-off) for ",
                "UE id ",
                ue,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
            )


if __name__ == "__main__":
    unittest.main()
