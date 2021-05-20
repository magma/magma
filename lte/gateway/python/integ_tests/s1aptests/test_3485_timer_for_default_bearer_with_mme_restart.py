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

import time
import unittest

import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class Test3485TimerForDefaultBearerWithMmeRestart(unittest.TestCase):
    """Test 3485 timer expiry for default bearer setup while mme restarts"""

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_3485_timer_for_default_bearer_with_mme_restart(self):
        """Test case validates the functionality of 3485 timer for

        default bearer while MME restarts
        Step1: UE attaches to network
        Step2: Send an indication to S1ap stack to drop Activate Default
        Eps Bearer Context Request message, sent as part of secondary PDN
        activation procedure.
        Step3: Initiate activation of secondary PDN
        Step4: Send an indication to S1ap stack to not to Activate Default
        Eps Bearer Context Request message, so that re-transmitted message
        reaches to TFW
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
            'apn_name': 'ims',  # APN-name
            'qci': 5,  # qci
            'priority': 15,  # priority
            'pre_cap': 0,  # preemption-capability
            'pre_vul': 0,  # preemption-vulnerability
            'mbr_ul': 200000000,  # MBR UL
            'mbr_dl': 100000000,  # MBR DL
        }

        # APN list to be configured
        apn_list = [ims]

        self._s1ap_wrapper.configAPN(
            'IMSI' + ''.join([str(idx) for idx in req.imsi]), apn_list,
        )
        print(
            '************************* Running End to End attach for UE id ',
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

        print(
            '*** Sending indication to drop Activate Default EPS bearer Ctxt'
            ' Req ***',
        )
        drop_acctv_dflt_bearer_req = (
            s1ap_types.UeDropActvDefaultEpsBearCtxtReq_t()
        )
        drop_acctv_dflt_bearer_req.ue_id = req.ue_id
        drop_acctv_dflt_bearer_req.dropActDfltEpsBearCtxtReq = True
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DROP_ACTV_DEFAULT_EPS_BEARER_CTXT_REQ,
            drop_acctv_dflt_bearer_req,
        )

        time.sleep(2)
        print('*** Sending PDN connectivity Req ***')
        # Send PDN Connectivity Request
        apn = 'ims'
        self._s1ap_wrapper.sendPdnConnectivityReq(ue_id, apn)
        # Receive PDN CONN RSP/Activate default EPS bearer context request

        print('************************* Restarting MME service on gateway')
        self._s1ap_wrapper.magmad_util.restart_services(['mme'])

        wait_for_restart = 20
        for j in range(wait_for_restart):
            print('Waiting for', j, 'seconds')
            time.sleep(1)

        print(
            '*** Sending indication to drop Activate Default EPS bearer Ctxt'
            ' Req ***',
        )
        drop_acctv_dflt_bearer_req = (
            s1ap_types.UeDropActvDefaultEpsBearCtxtReq_t()
        )
        drop_acctv_dflt_bearer_req.ue_id = req.ue_id
        drop_acctv_dflt_bearer_req.dropActDfltEpsBearCtxtReq = False
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DROP_ACTV_DEFAULT_EPS_BEARER_CTXT_REQ,
            drop_acctv_dflt_bearer_req,
        )

        retransmitted_response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            retransmitted_response.msg_type,
            s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        act_def_bearer_req = retransmitted_response.cast(
            s1ap_types.uePdnConRsp_t,
        )

        # Send PDN Disconnect
        pdn_disconnect_req = s1ap_types.uepdnDisconnectReq_t()
        pdn_disconnect_req.ue_Id = ue_id
        pdn_disconnect_req.epsBearerId = (
            act_def_bearer_req.m.pdnInfo.epsBearerId
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_PDN_DISCONNECT_REQ, pdn_disconnect_req,
        )

        # Receive UE_DEACTIVATE_BER_REQ
        response = self._s1ap_wrapper.s1_util.get_response()
        msg_type = s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value
        while (response.msg_type != msg_type):
            response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
        )

        print(
            '******************* Received deactivate eps bearer context'
            ' request',
        )
        # Send DeactDedicatedBearerAccept
        deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
        self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
            req.ue_id, deactv_bearer_req.bearerId,
        )

        print(
            '************************* Running UE detach (switch-off) for ',
            'UE id ',
            ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            wait_for_s1_ctxt_release=False,
        )


if __name__ == '__main__':
    unittest.main()
