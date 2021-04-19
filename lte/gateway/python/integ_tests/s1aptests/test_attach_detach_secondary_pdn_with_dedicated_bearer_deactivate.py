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


class TestSecondaryPdnConnWithDedBearerDeactivateReq(unittest.TestCase):
    """Test secondary pdn with dedicated bearer deactivation"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_secondary_pdn_conn_ded_bearer_deactivate(self):
        """Attach a single UE and send standalone PDN Connectivity
        Request + add dedicated bearer to each default bearer + deactivate
        dedicated bearers + detach
        """
        num_ues = 1

        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
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
                "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
            )

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
            # Create default flow list
            flow_list2 = self._spgw_util.create_default_ipv4_flows()
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
                flow_list2,
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
            # Verify if flow rules are created
            dl_flow_rules = {
                default_ip: [flow_list1],
                sec_ip: [flow_list2],
            }
            # 2 UL flows for default and secondary pdns +
            # 2 for dedicated bearers
            num_ul_flows = 4
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Delete dedicated bearer of secondary PDN (ims apn)
            print(
                "********************** Deleting dedicated bearer for ims"
                " apn",
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
                act_ded_ber_req_ims_apn.bearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)

            # Send Deactivate dedicated bearer rsp
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId,
            )

            print(
                "********************** Deleted dedicated bearer ",
                deactv_bearer_req.bearerId,
            )
            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Delete dedicated bearer of default PDN (magma.ipv4 apn)
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                attach.esmInfo.epsBearerId,
                act_ded_ber_req_oai_apn.bearerId,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
            # Send Deactivate dedicated bearer rsp
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId,
            )

            print(
                "********************** Deleted dedicated bearer ",
                deactv_bearer_req.bearerId,
            )
            print("Sleeping for 5 seconds")
            time.sleep(5)
            # Verify that flow rules for dedicated bearers are deleted
            dl_flow_rules = {
                default_ip: [],
                sec_ip: [],
            }
            # 2 UL flows for default and secondary pdns
            num_ul_flows = 2
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows, dl_flow_rules,
            )

            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                ue_id,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                False,
            )


if __name__ == "__main__":
    unittest.main()
