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
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestEpsBearerContextStatusMultipleDedBearerDeact(unittest.TestCase):
    """Test dedicated bearer deactivation with EPS bearer context status"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_eps_bearer_context_status_multiple_ded_bearer_deact(self):
        """Attach a single UE. Add dedicated bearer to the default PDN.
        Create 2 secondary PDNs and add 2
        dedicated bearers to each of the secondary PDNs.
        Send EPS bearer context status
        IE in TAU request with bearer ids
        5(def br),7(def br),8(ded br-LBI 7) and 10(def br) as active
        and bearers 6(ded br LBI-5),9(ded br LBI-7),11(ded br LBI-10)
        and 12(ded br LBI-10) as inactive.
        Set active flag to false.
        """

        num_ue = 1
        num_pdns = 2
        sec_ip = []
        flow_list2 = []
        wait_for_s1_context_rel = False

        self._s1ap_wrapper.configUEDevice(num_ue)
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

        internet = {
            "apn_name": "internet",  # APN-name
            "qci": 9,  # qci
            "priority": 15,  # priority
            "pre_cap": 0,  # preemption-capability
            "pre_vul": 0,  # preemption-vulnerability
            "mbr_ul": 250000000,  # MBR UL
            "mbr_dl": 150000000,  # MBR DL
        }

        # APN list to be configured
        apn_list = [ims, internet]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            apn_list,
        )
        print(
            "************************* Running End to End attach for UE id ",
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

        # Add dedicated bearer to the default bearer
        print(
            "********************** Adding dedicated bearer to IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        # Create default flow list
        flow_list1 = self._spgw_util.create_default_ipv4_flows()
        self._spgw_util.create_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            attach.esmInfo.epsBearerId,
            flow_list1,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        act_ded_ber_ctxt_req = response.cast(s1ap_types.UeActDedBearCtxtReq_t)
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            ue_id,
            act_ded_ber_ctxt_req.bearerId,
        )

        # APNs of the secondary PDNs
        apn = ["ims", "internet"]
        for i in range(num_pdns):
            # Send PDN Connectivity Request
            self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn[i])
            # Receive PDN CONN RSP/Activate default EPS bearer context request
            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value
            act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)

            print(
                "********************** Sending Activate default EPS bearer "
                "context accept for UE id ",
                ue_id,
            )
            print(
                "********************** Added default bearer with bearer id:",
                act_def_bearer_req.m.pdnInfo.epsBearerId,
            )
            addr = act_def_bearer_req.m.pdnInfo.pAddr.addrInfo
            sec_ip.append(ipaddress.ip_address(bytes(addr[:4])))

            # Add dedicated bearer for the default bearers
            print(
                "********************** Adding 1st dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            # Create default flow list
            flow_list2.append(self._spgw_util.create_default_ipv4_flows())
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
                flow_list2[i],
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t,
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                ue_id,
                act_ded_ber_ctxt_req.bearerId,
            )

            # Create 2nd dedicated bearer
            print(
                "********************** Adding 2nd dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )

            flow_list2.append(self._spgw_util.create_default_ipv4_flows())
            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]),
                act_def_bearer_req.m.pdnInfo.epsBearerId,
                flow_list2[i],
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t,
            )
            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                ue_id,
                act_ded_ber_ctxt_req.bearerId,
            )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify if flow rules are created
        for i in range(num_pdns):
            dl_flow_rules = {
                default_ip: [flow_list1],
                sec_ip[i]: [flow_list2[i]],
            }
            # 1 UL flow is created per bearer
            num_ul_flows = 8
            self._s1ap_wrapper.s1_util.verify_flow_rules(
                num_ul_flows,
                dl_flow_rules,
            )

        print(
            "************************* Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
        cntxt_rel_req.ue_Id = ue_id
        cntxt_rel_req.cause.causeVal = (
            gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST,
            cntxt_rel_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        print(" Sleeping for 2 seconds")
        time.sleep(2)
        print(
            "************************* Sending Tracking Area Update ",
            "request for UE id ",
            ue_id,
        )
        tau_req = s1ap_types.ueTauReq_t()
        tau_req.ue_Id = ue_id
        tau_req.type = s1ap_types.Eps_Updt_Type.TFW_TA_UPDATING.value
        tau_req.Actv_flag = True
        # Set default bearers 5,7,8 and 10 as active and
        # dedicated bearers 6,9,11 and 12 as inactive
        # epsBearerCtxSts IE is 16 bits
        # Ref: 3gpp 24.301 sec-9.9.2.1
        tau_req.epsBearerCtxSts = 0x5A0
        tau_req.ueMtmsi.pres = False
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_TAU_REQ,
            tau_req,
        )

        # Receive initial context setup and attach accept indication
        response = (
            self._s1ap_wrapper._s1_util
                .receive_initial_ctxt_setup_and_tau_accept()
        )
        tau_acc = response.cast(s1ap_types.ueTauAccept_t)
        print(
            "************************* Received Tracking Area Update ",
            "accept for UE Id:",
            tau_acc.ue_Id,
        )

        # Verify if flow rules are created
        dl_flow_rules = {
            default_ip: [],
            sec_ip[0]: [flow_list2[0]],
            sec_ip[1]: [],
        }
        # 1 UL flow is created per bearer
        num_ul_flows = 4
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows,
            dl_flow_rules,
        )

        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            wait_for_s1_context_rel,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify that all UL/DL flows are deleted
        self._s1ap_wrapper.s1_util.verify_flow_rules_deletion()


if __name__ == "__main__":
    unittest.main()
