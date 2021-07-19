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


class TestAttachStandaloneActvDfltBearCtxtRej(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_standalone_ActvDfltBearCtxtRej(self):
        """ Test case for sending Activate Default
        EPS Bearer Reject for secondary PDN """

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

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
            "************************* Running End to End attach for UE id ",
            req.ue_id,
        )
        # Attach
        attach = self._s1ap_wrapper.s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        addr = attach.esmInfo.pAddr.addrInfo
        default_ip = ipaddress.ip_address(bytes(addr[:4]))

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Trigger Activate Default EPS Bearer Context Reject indication
        # so that s1ap tester sends Activate default EPS bearer context reject
        # instead of Activate default EPS bearer context accept
        def_ber_rej = s1ap_types.ueActvDfltEpsBearerCtxtRej_t()
        def_ber_rej.ue_Id = req.ue_id
        def_ber_rej.bearerId = attach.esmInfo.epsBearerId
        def_ber_rej.cause = s1ap_types.TFW_EMM_CAUSE_PROT_ERR_UNSP

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.
            UE_STANDALONE_ACTV_DEFAULT_EPS_BEARER_CNTXT_REJECT,
            def_ber_rej,
        )

        print(
            "Sent STANDALONE_ACTV_DEFAULT_EPS_BEARER_CNTXT_REJECT indication",
        )
        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Send PDN Connectivity Request
        apn = "ims"
        # PDN Type 1=IPv4, 2=IPv6, 3=IPv4v6
        self._s1ap_wrapper.sendPdnConnectivityReq(req.ue_id, apn, pdn_type=3)
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        act_def_bearer_req = response.cast(s1ap_types.uePdnConRsp_t)

        print(
            "************************* Sending Activate default EPS bearer "
            "context reject for UE id %d and bearer %d"
            % (req.ue_id, act_def_bearer_req.m.pdnInfo.epsBearerId),
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)

        # Verify that ovs rule is not is created for the secondary pdn
        # as UE rejected the establishment of secondary pdn

        # 1 UL flow for the default bearer
        num_ul_flows = 1
        # No dedicated bearers, so flow list will be empty
        dl_flow_rules = {
            default_ip: [],
        }

        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Now detach the UE
        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ",
            req.ue_id,
        )
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            False,
        )


if __name__ == "__main__":
    unittest.main()
