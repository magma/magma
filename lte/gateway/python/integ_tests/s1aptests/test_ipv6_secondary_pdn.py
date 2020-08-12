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


class TestIPv6SecondaryPdn(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_ipv6_secondary_pdn(self):
        """ Attach a single UE + add 2 secondary pdns with """
        """ IPv4v6 and IPv6 address respectively + detach """
        num_ue = 1

        self._s1ap_wrapper.configUEDevice(num_ue)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id

        # ims and internet PDNs
        num_pdns = 2
        # APN of the secondary PDNs
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

        internet_apn = {
            "apn_name": "internet",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 200000000,  # MBR UL
            "mbr_dl": 100000000,  # MBR DL
            "pdn_type": 1,  # PDN Type 0-IPv4,1-IPv6,2-IPv4v6
        }

        apn_list = [ims_apn, internet_apn]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list
        )
        print(
            "*********************** Running End to End attach for UE id ",
            ue_id,
        )

        print("***** Sleeping for 5 seconds")
        time.sleep(5)
        # Attach
        self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        apn = ["ims", "internet"]
        # PDN Type 2 = IPv6, 3 = IPv4v6
        pdn_type = [3, 2]
        # Send PDN Connectivity Request for ims and internet apns
        for i in range(num_pdns):
            # Send PDN Connectivity Request
            self._s1ap_wrapper.sendPdnConnectivityReq(
                ue_id, apn[i], pdn_type=pdn_type[i]
            )
            # Receive PDN CONN RSP/Activate default EPS bearer context request
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
            )
            act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)

            print(
                "********************** Sending Activate default EPS bearer "
                "context accept for APN-%s, UE id-%d" % (apn[i], ue_id),
            )
            print(
                "********************** Added default bearer for apn-%s,"
                " bearer id-%d, pdn type-%d"
                % (
                    apn[i],
                    act_def_bearer_req.m.pdnInfo.epsBearerId,
                    pdn_type[i],
                )
            )

            time.sleep(5)
            # TODO- Uncomment once its imlemeted in magma
            # Receive Router Advertisement message
            # response = self._s1ap_wrapper.s1_util.get_response()
            # self.assertEqual(
            #    response.msg_type, s1ap_types.tfwCmd.UE_ROUTER_ADV_IND.value
            # )
            # routerAdv = response.cast(s1ap_types.ueRouterAdv_t)
            # print(
            #    "******************* Received Router Advertisement for APN-%s
            #    with IPv6 address-%s ,bearer id-%d",
            #    apn[i], routerAdv.ipv6Addr, routerAdv.bearerId
            # )

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
