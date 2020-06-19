"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import unittest
import time

import s1ap_types
import s1ap_wrapper


class TestAttachDetachSecondaryPdnWithPcscfAddress(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_secondary_pdn_with_pcscf_address(self):
        """ Attach a single UE + add 2 PDN Connections with P-CSCF address """
        """ + disconnect """
        num_pdns = 2
        bearer_ids = []
        num_ue = 1

        self._s1ap_wrapper.configUEDevice(num_ue)
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

        # internet APN
        internet = {
            "apn_name": "internet",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
        }

        # APN list to be configured
        apn_list = [ims, internet]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list
        )
        print(
            "*********************** Running End to End attach for UE id ",
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

        time.sleep(2)
        # APNs of the secondary PDNs
        apn = ["internet", "ims"]
        pcscf_addr_type = ["ipv4", "ipv6"]
        # PDN Type 1-IPv4,2-IPv6,3-IPv4v6
        pdn_type = 1
        for i in range(num_pdns):
            # Send PDN Connectivity Request
            self._s1ap_wrapper.sendPdnConnectivityReq(
                ue_id, apn[i], pdn_type=pdn_type, pcscf_addr_type=pcscf_addr_type[i]
            )
            # Receive PDN CONN RSP/Activate default EPS bearer context request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
            )
            act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)

            print(
                "********************** Sending Activate default EPS bearer "
                "context accept for UE id ",
                ue_id,
            )
            bearer_ids.append(act_def_bearer_req.m.pdnInfo.epsBearerId)
            print(
                "********************** Added default bearer with "
                "bearer id",
                act_def_bearer_req.m.pdnInfo.epsBearerId,
            )

        time.sleep(5)
        for i in range(num_pdns):
            # Send PDN Disconnect
            pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
            pdn_disconnect_req.ue_Id = ue_id
            pdn_disconnect_req.epsBearerId = bearer_ids[i]
            print(
                "******************* Sending PDN Disconnect bearer id\n",
                pdn_disconnect_req.epsBearerId,
            )
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req
            )

            # Receive UE_DEACTIVATE_BER_REQ
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print(
                "******************* Received deactivate eps bearer context"
                " request"
            )
            # Send DeactDedicatedBearerAccept
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                ue_id, bearer_ids[i]
            )

        print(
            "******************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False
        )


if __name__ == "__main__":
    unittest.main()
