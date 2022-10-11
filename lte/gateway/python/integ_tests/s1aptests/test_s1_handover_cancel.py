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


class TestS1HandoverCancel(unittest.TestCase):
    """Integration Test: TestS1HandoverCancel"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_s1_handover_cancel(self):
        """S1 Handover Unsuccessful Scenario (S1 HO Cancel):

        1) Attach UE to ENB 1 (After handover UE should switch to ENB 2)
        2) Trigger handover by sending S1 HO required Message from source ENB
        3) Receive and Handle the S1 HO Request by Target ENB
        4) Send S1 HO Cancel in response to S1 HO Command from source ENB
        5) Receive and handle the S1 HO Cancel Ack Message in source ENB
        6) S1 HO has failed. UE is still attached to ENB 1. Detach the UE

        Note: Before execution of this test case,
        Run the test script s1aptests/test_modify_mme_config_for_sanity.py
        to update multiple PLMN/TAC configuration in MME and
        after test case execution, restore the MME configuration by running
        the test script s1aptests/test_restore_mme_config_after_sanity.py

        Or

        Make sure that following steps are correct
        1. Configure same plmn and tac in both MME and S1APTester
        2. How to configure plmn and tac in MME:
        2a. Set mcc and mnc in gateway.mconfig for mme service
        2b. Set tac in gateway.mconfig for mme service
        2c. Restart MME service
        3. How to configure plmn and tac in S1APTester,
        3a. For multi-eNB test case, configure plmn and tac from test case.
        In each multi-eNB test case, set plmn, plmn length and tac in enb_list
        3b. For single eNB test case, configure plmn and tac in nbAppCfg.txt
        """
        # Column is an ENB parameter, Row is number of ENB
        # Columns: Cell Id, Tac, EnbType, PLMN Id, PLMN length
        enb_list = [
            [1, 1, 1, "00101", 5],
            [2, 2, 1, "00101", 5],
        ]

        self._s1ap_wrapper.multiEnbConfig(len(enb_list), enb_list)

        print("Waiting for 2 seconds for multiple ENBs to get configured")
        time.sleep(2)
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for UE Id:",
            req.ue_id,
        )
        # Now actually complete the attach
        attach = self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        addr = attach.esmInfo.pAddr.addrInfo
        default_ip = ipaddress.ip_address(bytes(addr[:4]))

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print("Waiting for 3 seconds for the flow rules creation")
        time.sleep(3)
        # Verify if flow rules are created
        # 1 UL flow for default bearer
        num_ul_flows = 1
        dl_flow_rules = {default_ip: []}
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows,
            dl_flow_rules,
        )

        # Trigger the S1 Handover Procedure from Source ENB by sending S1
        # Handover Required Message to MME
        print(
            "************************* Sending S1 Handover Required for UE Id:",
            req.ue_id,
        )
        s1ho_required = s1ap_types.FwNbS1HoRequired_t()
        s1ho_required.ueId = req.ue_id
        s1ho_required.s1HoEvent = s1ap_types.FwS1HoEvents.FW_S1_HO_CANCEL.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.S1_HANDOVER_REQUIRED,
            s1ho_required,
        )

        # After receiving S1 Handover Required from Source ENB, MME sends S1
        # Handover Request to Target ENB.
        # Wait for S1 Handover Request Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.S1_HANDOVER_REQ_IND.value
        s1ho_req_ind = response.cast(s1ap_types.FwNbS1HoReqInd_t)
        print(
            "************************* Received S1 Handover Request "
            "Indication (UeId: "
            + str(s1ho_req_ind.ueId)
            + ", Connected EnbId: "
            + str(s1ho_req_ind.currEnbId)
            + ") (HO SrcEnbId: "
            + str(s1ho_req_ind.hoSrcEnbId)
            + ", HO TgtEnbId: "
            + str(s1ho_req_ind.hoTgtEnbId)
            + ")",
        )

        # Send the S1 Handover Request Ack message from Target ENB to MME
        print(
            "************************* Sending S1 Handover Req Ack for UE Id:",
            req.ue_id,
        )
        s1ho_req_ack = s1ap_types.FwNbS1HoReqAck_t()
        s1ho_req_ack.ueId = req.ue_id
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.S1_HANDOVER_REQ_ACK,
            s1ho_req_ack,
        )

        # After receiving S1 Handover Req Ack from Target ENB, MME sends S1
        # Handover Command to Source ENB.
        # Wait for S1 Handover Command Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.S1_HANDOVER_CMD_IND.value
        s1ho_cmd_ind = response.cast(s1ap_types.FwNbS1HoCmdInd_t)
        print(
            "************************* Received S1 Handover Command "
            "Indication (UeId: "
            + str(s1ho_cmd_ind.ueId)
            + ", Connected EnbId: "
            + str(s1ho_cmd_ind.currEnbId)
            + ") (HO SrcEnbId: "
            + str(s1ho_cmd_ind.hoSrcEnbId)
            + ", HO TgtEnbId: "
            + str(s1ho_cmd_ind.hoTgtEnbId)
            + ")",
        )

        # Send the S1 Handover Cancel message from source ENB to MME
        print(
            "************************* Sending S1 Handover Cancel for UE Id:",
            req.ue_id,
        )
        s1ho_cancel = s1ap_types.FwNbS1HoCancel_t()
        s1ho_cancel.ueId = req.ue_id
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.S1_HANDOVER_CANCEL,
            s1ho_cancel,
        )

        # After receiving S1 Handover Cancel from Source ENB, MME sends S1
        # Handover Cancel Ack to Source ENB.
        # Wait for S1 Handover Cancel Ack Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.S1_HANDOVER_CANCEL_ACK_IND.value
        s1ho_cancl_ack_ind = response.cast(s1ap_types.FwNbS1HoCancelAckInd_t)
        print(
            "************************* Received S1 Handover Cancel Ack "
            "Indication (UeId: "
            + str(s1ho_cancl_ack_ind.ueId)
            + ", Connected EnbId: "
            + str(s1ho_cancl_ack_ind.currEnbId)
            + ")",
        )

        print(
            "************************* Running UE detach for UE Id:",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        )

        print("Waiting for 5 seconds for the flow rules deletion")
        time.sleep(5)
        # Verify that all UL/DL flows are deleted
        self._s1ap_wrapper.s1_util.verify_flow_rules_deletion()


if __name__ == "__main__":
    unittest.main()
