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

import unittest
import time

import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class Test3485TimerForDefaultBearerWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_3485_timer_for_default_bearer_with_mme_restart(self):
        """ Test case validates the functionality of 3485 timer for
            default bearer while MME restarts
        Step1: UE attaches to network
        Step2: Send an indication to S1ap stack to drop E-Rab Setup
               Request message, sent as part of secondary PDN activation
               procedure.
        Step3: Initiate activation of secondary PDN
        Step4: Send an indication to S1ap stack to not to drop E-Rab Setup
               Request message, so that re-transmitted message reaches to
               TFW
        Step5: Send command to Magma to restart mme service
        Step6: TFW shall receive the PDN connectivity response
        Step7: TFW shall initiate de-activation of secondary PDN and then
               initiate Detach procedure.
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
            "IMSI" + "".join([str(i) for i in req.imsi]), apn_list
        )
        print(
            "************************* Running End to End attach for UE id ",
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

        print("*** Sending indication to drop Erab Setup Req ***")
        drop_erab_setup_req = s1ap_types.UeDropErabSetupReq_t()
        drop_erab_setup_req.ue_Id = req.ue_id
        drop_erab_setup_req.flag = True
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DROP_ERAB_SETUP_REQ, drop_erab_setup_req
        )

        time.sleep(2)
        print("*** Sending PDN connectivity Req ***")
        # Send PDN Connectivity Request
        apn = "ims"
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        # Receive PDN CONN RSP/Activate default EPS bearer context request
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ERAB_DROP_IND.value
        )
        print(
            "*******************Received ERAB_DROP_IND for last received ",
            "E-Rab Setup Request for secondary PDN's default bearer ",
        )

        print(
            "*** Sending indication to disable dropping of Erab Setup Req",
            " ***",
        )
        drop_erab_setup_req = s1ap_types.UeDropErabSetupReq_t()
        drop_erab_setup_req.ue_Id = req.ue_id
        drop_erab_setup_req.flag = False
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DROP_ERAB_SETUP_REQ, drop_erab_setup_req
        )

        print("************************* Restarting MME service on gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        retransmitted_response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            retransmitted_response.msg_type,
            s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        act_def_bearer_req = retransmitted_response.cast(
            s1ap_types.uePdnConRsp_t
        )

        # Send PDN Disconnect
        pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
        pdn_disconnect_req.ue_Id = ue_id
        pdn_disconnect_req.epsBearerId = (
            act_def_bearer_req.m.pdnInfo.epsBearerId
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req
        )

        # Receive UE_DEACTIVATE_BER_REQ
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value
        )

        print(
            "******************* Received deactivate eps bearer context"
            " request"
        )
        # Send DeactDedicatedBearerAccept
        deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
        self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
            req.ue_id, deactv_bearer_req.bearerId
        )

        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ",
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False
        )


if __name__ == "__main__":
    unittest.main()
