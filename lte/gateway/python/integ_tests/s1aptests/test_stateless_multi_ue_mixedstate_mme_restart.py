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

import ctypes
import ipaddress
import time
import unittest

import gpp_types
import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestStatelessMultiUeMixedStateMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
           stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )
        self.dl_flow_rules = {}

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def exec_attach_req_step(self, ue_id):
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        # enabling ESM Information transfer flag
        attach_req.eti.pres = 1
        attach_req.eti.esm_info_transfer_flag = 1

        print("Sending Attach Request ue-id", attach_req.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("Received auth req ind ue-id", attach_req.ue_Id)

    def exec_auth_resp_step(self, ue_id):
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = ue_id
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        print("Sending Auth Response ue-id", auth_res.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        print("Received Security Mode Command ue-id", auth_res.ue_Id)

    def exec_sec_mode_complete_step(self, ue_id):
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )
        print(
            "Received Esm Information Request ue-id", sec_mode_complete.ue_Id,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ESM_INFORMATION_REQ.value,
        )
        esm_info_req = response.cast(s1ap_types.ueEsmInformationReq_t)
        return esm_info_req.tId

    def exec_esm_inf_req_step(self, ue_id, tId):
        # Sending Esm Information Response
        print(
            "Sending Esm Information Response ue-id", ue_id,
        )
        esm_info_response = s1ap_types.ueEsmInformationRsp_t()
        esm_info_response.ue_Id = ue_id
        esm_info_response.tId = tId  # esm_info_req.tId
        esm_info_response.pdnAPN_pr.pres = 1
        s = "magma.ipv4"
        esm_info_response.pdnAPN_pr.len = len(s)
        esm_info_response.pdnAPN_pr.pdn_apn = (ctypes.c_ubyte * 100)(
            *[ctypes.c_ubyte(ord(c)) for c in s[:100]],
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ESM_INFORMATION_RSP, esm_info_response,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
        )
        msg = response.cast(s1ap_types.ueAttachAccept_t)
        addr = msg.esmInfo.pAddr.addrInfo
        default_ip = ipaddress.ip_address(bytes(addr[:4]))
        self.dl_flow_rules[default_ip] = []
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

    def test_stateless_multi_ue_mixedstate_mme_restart(self):
        """ Testing of sending Esm Information procedure """
        ue_ids = []
        num_ues_idle = 10
        num_ues_active = 10
        num_attached_ues = num_ues_idle + num_ues_active
        # each list item can be a number in [1,3] and be repeated
        stateof_ues_in_attachproc_before_restart = [
            1, 1, 1, 1, 1,
            2, 2, 2, 2, 2,
            3, 3, 3, 3, 3,
        ]
        num_ues_attaching = len(stateof_ues_in_attachproc_before_restart)

        attach_steps = [
            self.exec_attach_req_step,
            self.exec_auth_resp_step,
            self.exec_sec_mode_complete_step,
            self.exec_esm_inf_req_step,
        ]
        num_of_steps = len(attach_steps)

        tot_num_ues = num_ues_idle + num_ues_active + num_ues_attaching
        self._s1ap_wrapper.configUEDevice(tot_num_ues)

        idle_session_ips = []
        # Prep attached UEs
        for i in range(num_ues_idle + num_ues_active):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* sending Attach Request for "
                "UE id ", req.ue_id,
            )
            attach = self._s1ap_wrapper._s1_util.attach(
                req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            addr = attach.esmInfo.pAddr.addrInfo
            default_ip = ipaddress.ip_address(bytes(addr[:4]))
            if i < num_ues_idle:
                idle_session_ips.append(default_ip)
            else:
                self.dl_flow_rules[default_ip] = []

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)

        # Move first num_ues_idle UEs to idle state
        for i in range(num_ues_idle):
            print(
                "************************* Sending UE context release request ",
                "for UE id ", ue_ids[i],
            )
            # Send UE context release request to move UE to idle mode
            ue_cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
            ue_cntxt_rel_req.ue_Id = ue_ids[i]
            ue_cntxt_rel_req.cause.causeVal = (
                gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
            )
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, ue_cntxt_rel_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )

        tId = {}
        # start attach procedures for the remaining UEs
        for i in range(num_ues_attaching):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Starting Attach procedure "
                "UE id ", req.ue_id,
            )
            # bring each newly attaching UE to the desired point during
            # attach procedure before restarting mme service
            ue_ids.append(req.ue_id)
            for step in range(stateof_ues_in_attachproc_before_restart[i]):
                if attach_steps[step] == self.exec_sec_mode_complete_step:
                    tId[req.ue_id] = attach_steps[step](req.ue_id)
                elif attach_steps[step] == self.exec_esm_inf_req_step:
                    attach_steps[step](req.ue_id, tId[req.ue_id])
                else:
                    attach_steps[step](req.ue_id)

        # Restart mme
        self._s1ap_wrapper.magmad_util.restart_mme_and_wait()

        # Post restart, complete the attach procedures that were cut in between
        for i in range(num_ues_attaching):
            # resume attach for attaching UEs
            print(
                "************************* Resuming Attach procedure "
                "UE id ", ue_ids[i + num_attached_ues],
            )
            for step in range(stateof_ues_in_attachproc_before_restart[i], num_of_steps):
                if attach_steps[step] == self.exec_sec_mode_complete_step:
                    tId[ue_ids[i + num_attached_ues]] = attach_steps[step](ue_ids[i + num_attached_ues])
                elif attach_steps[step] == self.exec_esm_inf_req_step:
                    attach_steps[step](ue_ids[i + num_attached_ues], tId[ue_ids[i + num_attached_ues]])
                else:
                    attach_steps[step](ue_ids[i + num_attached_ues])

        # Verify steady state flows in Table-0
        # Idle users will have paging rules installed
        # Active users will have tunnel rules
        # 1 UL flow is created per active bearer
        num_ul_flows = num_ues_active + num_ues_attaching
        # Verify if flow rules are created
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            num_ul_flows, self.dl_flow_rules,
        )
        # Verify paging flow rules for idle sessions
        self._s1ap_wrapper.s1_util.verify_paging_flow_rules(idle_session_ips)

        # Try to bring idle mode users into active state
        for i in range(num_ues_idle):
            print(
                "************************* Sending Service Request ",
                "for UE id ", ue_ids[i],
            )
            # Send service request to reconnect UE
            ser_req = s1ap_types.ueserviceReq_t()
            ser_req.ue_Id = ue_ids[i]
            ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
            ser_req.ueMtmsi.pres = False
            ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST, ser_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
            )
            self.dl_flow_rules[idle_session_ips[i]] = []

        # Verify default bearer rules
        self._s1ap_wrapper.s1_util.verify_flow_rules(
            tot_num_ues, self.dl_flow_rules,
        )

        # detach everyone
        print("*** Starting Detach Procedure for all UEs ***")
        for ue in ue_ids:
            print(
                "************************* Detaching "
                "UE id ", ue,
            )
            # Now detach the UE
            detach_req = s1ap_types.uedetachReq_t()
            detach_req.ue_Id = ue
            detach_req.ueDetType = (
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
            )
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )


if __name__ == "__main__":
    unittest.main()
