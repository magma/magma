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

# TODO: Current S1APTester logic is broken. Fix as follows
# 1- Dont send messages automatically on behalf of eNBs
# 2- Handover Notify Msg from target ENB should only be sent after
# successfully receiving the MME Status Transfer Msg
# 3- Source eNB will receive a UE Context Release Command once MME
# receives a Handover Notify Msg from Target eNB

# Once S1APTester logic is fixed, add this back to integ tests


class TestS1HandoverPingPong(unittest.TestCase):
    """Unittest: TestS1HandoverPingPong"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_s1_handover_ping_pong(self):
        """S1 Handover Successful Scenario (Ping-Pong):

        1) Attach UE to ENB 1 (After handover UE should switch to ENB 2)
        2) Trigger handover by sending S1 HO required Message from source ENB
        3) Receive and Handle the S1 HO Request in target ENB
        4) Receive and handle the S1 HO Command in source ENB
        5) Receive and handle the MME Status Transfer message in target ENB
        6) UE is moved to ENB 2. Re-initiate the S1 Handover procedure
        7) After UE is moved back to ENB 1, detach the UE

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

        # Configure the S1 Handover Event type as S1_HO_SUCCESS and
        # Trigger the S1 Handover Procedure from Source ENB by sending S1
        # Handover Required Message
        print(
            "************************* Sending S1 Handover Required for UE Id:",
            req.ue_id,
        )
        s1ho_required = s1ap_types.FwNbS1HoRequired_t()
        s1ho_required.ueId = req.ue_id
        s1ho_required.s1HoEvent = (
            s1ap_types.FwS1HoEvents.FW_S1_HO_SUCCESS.value
        )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.S1_HANDOVER_REQUIRED,
            s1ho_required,
        )

        # After receiving S1 Handover Required from Source ENB, MME sends S1
        # Handover Request to Target ENB.
        # Wait for S1 Handover Request Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.S1_HANDOVER_REQ_IND.value,
        )
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

        # After receiving S1 Handover Request, S1APTester stack sends S1
        # Handover Request Ack from target ENB to MME. MME then sends S1
        # Handover Command to Source ENB.
        # Wait for S1 Handover Command Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.S1_HANDOVER_CMD_IND.value,
        )
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

        # After receiving S1 Handover Command, S1APTester stack sends ENB
        # Status Transfer from source ENB to MME and Handover Notify from
        # target ENB to MME. MME then sends MME Status Transfer to Target ENB.
        # Wait for MME Status Transfer Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.S1_MME_STATUS_TRSFR_IND.value,
        )
        mme_status_trf_ind = response.cast(s1ap_types.FwNbMmeStatusTrnsfrInd_t)
        print(
            "************************* Received MME Status Transfer "
            "Indication (UeId: "
            + str(mme_status_trf_ind.ueId)
            + ", Connected EnbId: "
            + str(mme_status_trf_ind.currEnbId)
            + ") (HO SrcEnbId: "
            + str(mme_status_trf_ind.hoSrcEnbId)
            + ", HO TgtEnbId: "
            + str(mme_status_trf_ind.hoTgtEnbId)
            + ")",
        )

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

        # After S1 Overall Reloc Timer expiry, source ENB sends UE context
        # release request to MME with cause tS1relocoverall-expiry.
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("Received UE Context Release complete indication")

        #####################################
        # Source and Target ENBs are switched
        # Re-initiate the Handover Procedure
        #####################################

        print("Waiting for 2 seconds before triggering ping-pong Handover")
        time.sleep(2)
        # Configure the S1 Handover Event type as S1_HO_SUCCESS and
        # Trigger the S1 Handover Procedure from Source ENB by sending S1
        # Handover Required Message
        print(
            "************************* Sending S1 Handover Required for UE Id:",
            req.ue_id,
        )
        s1ho_required = s1ap_types.FwNbS1HoRequired_t()
        s1ho_required.ueId = req.ue_id
        s1ho_required.s1HoEvent = (
            s1ap_types.FwS1HoEvents.FW_S1_HO_SUCCESS.value
        )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.S1_HANDOVER_REQUIRED,
            s1ho_required,
        )

        # After receiving S1 Handover Required from Source ENB, MME sends S1
        # Handover Request to Target ENB.
        # Wait for S1 Handover Request Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.S1_HANDOVER_REQ_IND.value,
        )
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

        # After receiving S1 Handover Request, S1APTester stack sends S1
        # Handover Request Ack from target ENB to MME. MME then sends S1
        # Handover Command to Source ENB.
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.S1_HANDOVER_CMD_IND.value,
        )
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

        # After receiving S1 Handover Command, S1APTester stack sends ENB
        # Status Transfer from source ENB to MME and Handover Notify from
        # target ENB to MME. MME then sends MME Status Transfer to Target ENB.
        # Wait for MME Status Transfer Indication
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.S1_MME_STATUS_TRSFR_IND.value,
        )
        mme_status_trf_ind = response.cast(s1ap_types.FwNbMmeStatusTrnsfrInd_t)
        print(
            "************************* Received MME Status Transfer "
            "Indication (UeId: "
            + str(mme_status_trf_ind.ueId)
            + ", Connected EnbId: "
            + str(mme_status_trf_ind.currEnbId)
            + ") (HO SrcEnbId: "
            + str(mme_status_trf_ind.hoSrcEnbId)
            + ", HO TgtEnbId: "
            + str(mme_status_trf_ind.hoTgtEnbId)
            + ")",
        )

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

        # After S1 Overall Reloc Timer expiry, source ENB sends UE context
        # release request to MME with cause tS1relocoverall-expiry.
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("Received UE Context Release complete indication")

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
