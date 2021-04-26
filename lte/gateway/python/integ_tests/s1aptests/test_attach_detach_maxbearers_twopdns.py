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
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestMaximumBearersTwoPdnsPerUe(unittest.TestCase):
    """Test maximum bearers with two pdns"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_maxbearers_twopdns(self):
        """Attach a single UE and send standalone PDN Connectivity
        Request + add 9 dedicated bearers + detach
        """
        num_ues = 1
        flow_lists2 = []
        self._s1ap_wrapper.configUEDevice(num_ues)

        # 1 oai PDN + 1 dedicated bearer, 1 ims pdn + 8 dedicated bearers
        loop = 8

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

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req

            self._s1ap_wrapper.configAPN(
                "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
            )

            ue_id = req.ue_id
            print(
                "********************* Running End to End attach for UE id ",
                ue_id,
            )
            # Attach
            attach = self._s1ap_wrapper.s1_util.attach(
                ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            addr = attach.esmInfo.pAddr.addrInfo
            default_ip = ipaddress.ip_address(bytes(addr[:4]))

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            # Add dedicated bearer for default bearer 5
            print(
                "********************** Adding dedicated bearer to magma.ipv4"
                " PDN",
            )
            # Create default flow list
            flow_list1 = self._spgw_util.create_default_ipv4_flows()
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach.esmInfo.epsBearerId,
                flow_list1,
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
            # Receive PDN CONN RSP/Activate default EPS bearer context req
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
            for i in range(loop):
                # Add dedicated bearer to 2nd PDN
                print(
                    "********************** Adding dedicated bearer to ims"
                    " PDN",
                )
                flow_lists2.append(
                    self._spgw_util.create_default_ipv4_flows(port_idx=i),
                )
                self._spgw_util.create_bearer(
                    "IMSI" + "".join([str(i) for i in req.imsi]),
                    act_def_bearer_req.m.pdnInfo.epsBearerId,
                    flow_lists2[i],
                    qci_val=i + 1,
                )
                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertEqual(
                    response.msg_type,
                    s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
                )
                act_ded_ber_req_ims_apn = response.cast(
                    s1ap_types.UeActDedBearCtxtReq_t,
                )
                self._s1ap_wrapper.sendActDedicatedBearerAccept(
                    req.ue_id, act_ded_ber_req_ims_apn.bearerId,
                )
                print(
                    "************ Added dedicated bearer",
                    act_ded_ber_req_ims_apn.bearerId,
                )
                print("Sleeping for 2 seconds")
                time.sleep(2)

            # Verify if flow rules are created
            # 1 dedicated bearer for default pdn and 8 dedicated bearers
            # for secondary pdn
            dl_flow_rules = {
                default_ip: [flow_list1],
                sec_ip: flow_lists2,
            }
            # 2 UL flows for default and secondray pdn +
            # 9 for dedicated bearers
            num_ul_flows = 11
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

            print(
                "************************ Running UE detach (switch-off) for ",
                "UE id ",
                ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                False,
            )


if __name__ == "__main__":
    unittest.main()
