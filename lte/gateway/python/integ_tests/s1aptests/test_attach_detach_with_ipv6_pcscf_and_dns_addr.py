"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest

import s1ap_types

from integ_tests.s1aptests import s1ap_wrapper


class TestAttachDetachWithIpv6PcscfAndDnsAddress(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_with_ipv6_pcscf_and_dns_address(self):
        """ Basic attach/detach test with IPv6 P-CSCF and DNS IPv6 address """
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)
        req = self._s1ap_wrapper.ue_req
        pcscf_addr_type = "ipv6"
        # PDN Type 1-IPv4,2-IPv6,3-IPv4v6 as per 3gpp 24.301/29.274
        pdn_type = 1

        # APN of the secondary PDN
        ims = {
            "apn_name": "ims",  # APN-name
            "qci": 5,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
            "pdn_type": 1,  # PDN Type 0-IPv4,1-IPv6,2-IPv4v6
            # as per 3gpp 29.272
        }

        # APN list to be configured
        apn_list = [ims]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list
        )

        print(
            "************************* Running End to End attach for ",
            "UE id ",
            req.ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
            pdn_type=pdn_type,
            pcscf_addr_type=pcscf_addr_type,
            dns_ipv6_addr=True,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Send PDN Connectivity Request
        print("***************** Sending secondary PDN request for IMS APN")
        apn = "ims"
        pcscf_addr_type = "ipv6"
        # PDN Type 1-IPv4,2-IPv6,3-IPv4v6 as per 29.274
        pdn_type = 2

        self._s1ap_wrapper.sendPdnConnectivityReq(
            req.ue_id,
            apn,
            pdn_type=pdn_type,
            pcscf_addr_type=pcscf_addr_type,
            dns_ipv6_addr=True,
        )
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
        )
        print(
            "************************* Sending Activate default EPS bearer "
            "context accept for UE id ",
            req.ue_id,
        )

        print(
            "************************* Running UE detach for UE id ",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
