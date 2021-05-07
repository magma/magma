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


class TestSecondaryPdnDisconnInvalidBearerId(unittest.TestCase):
    """Test secondary pdn disconnect with invalid bearerid"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_secondary_pdn_disconn_invalid_bearer_id(self):
        """Attach a single UE + send standalone PDN Connectivity
        Request + send PDN disconnect with invalid bearer id
        """
        num_ue = 1

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

        # APN list to be configured
        apn_list = [ims]

        self._s1ap_wrapper.configAPN(
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list,
        )

        # Declare an array of len 15 as the bearer id ranges from 5-15
        length = 15
        bearer_idx = [0] * length
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Attach
        attach_accept = self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        addr = attach_accept.esmInfo.pAddr.addrInfo
        default_ip = ipaddress.ip_address(bytes(addr[:4]))

        # Set the bearer index to 1
        bearer_idx[attach_accept.esmInfo.epsBearerId] = 1
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

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

        # Set the bearer index to 1
        bearer_idx[act_def_bearer_req.m.pdnInfo.epsBearerId] = 1
        print(
            "************************* Sending Activate default EPS bearer "
            "context accept for UE id ",
            ue_id,
        )

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify if flow rules are created
        # Flow list is empty as there are no dedicated bearers
        dl_flow_rules = {
            default_ip: [],
            sec_ip: [],
        }
        # 2 UL flows for default and secondary pdns
        num_ul_flows = 2
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print("********************* Sleeping for 5 seconds")
        time.sleep(5)
        # Send PDN Disconnect for a non-existent bearer
        pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
        pdn_disconnect_req.ue_Id = ue_id
        # Find an unassigned bearer id
        # Start from 5th index as the bearer id ranges from 5-15
        for i in range(5, 15):
            if bearer_idx[i] == 0:
                pdn_disconnect_req.epsBearerId = i
                break
        print(
            "****** Sending PDN disconnect for bearer id",
            pdn_disconnect_req.epsBearerId,
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req,
        )

        # Receive PDN Disconnect reject
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REJ.value,
        )

        print("************************* Received PDN disconnect reject")

        print("Sleeping for 5 seconds")
        time.sleep(5)
        # Verify that flow rules are not deleted
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, dl_flow_rules,
        )

        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
        )


if __name__ == "__main__":
    unittest.main()
