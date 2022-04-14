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
import random
import threading
import unittest

import s1ap_types
import s1ap_wrapper


class TestContinuousRandomAttach(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def handle_msg(self, msg):
        if msg.msg_type == s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value:
            self.auth_req_ind_count += 1
            m = msg.cast(s1ap_types.ueAuthReqInd_t)
            print("====> Received UE_AUTH_REQ_IND ue-id", m.ue_Id)
            auth_res = s1ap_types.ueAuthResp_t()
            auth_res.ue_Id = m.ue_Id
            sqn_recvd = s1ap_types.ueSqnRcvd_t()
            sqn_recvd.pres = 0
            auth_res.sqnRcvd = sqn_recvd
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
            )
        elif msg.msg_type == s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value:
            self.sec_mod_cmd_ind_count += 1
            m = msg.cast(s1ap_types.ueSecModeCmdInd_t)
            print("============>  Received UE_SEC_MOD_CMD_IND ue-id", m.ue_Id)
            sec_mode_complete = s1ap_types.ueSecModeComplete_t()
            sec_mode_complete.ue_Id = m.ue_Id
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
            )
        elif msg.msg_type == s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value:
            self.attach_accept_ind_count += 1
            m = msg.cast(s1ap_types.ueAttachAccept_t)
            print("====================> Received UE_ATTACH_ACCEPT_IND ue-id", m.ue_Id)
            pdn_type = m.esmInfo.pAddr.pdnType
            addr = m.esmInfo.pAddr.addrInfo
            if self._s1ap_wrapper._s1_util.CM_ESM_PDN_IPV4 == pdn_type:
                # Cast and cache the IPv4 address
                ip = ipaddress.ip_address(bytes(addr[:4]))
                with self._s1ap_wrapper._s1_util._lock:
                    self._s1ap_wrapper._s1_util._ue_ip_map[m.ue_Id] = ip
            attach_complete = s1ap_types.ueAttachComplete_t()
            attach_complete.ue_Id = m.ue_Id
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete,
            )
        elif msg.msg_type == s1ap_types.tfwCmd.UE_IDENTITY_REQ_IND.value:
            self.identity_req_ind_count += 1
            m = msg.cast(s1ap_types.ueIdentityReqInd_t)
            print("=> Received UE_IDENTITY_REQ_IND ue-id", m.ue_Id)
            us_identity_resp = s1ap_types.ueIdentityResp_t()
            us_identity_resp.ue_Id = m.ue_Id
            us_identity_resp.idType = m.idType
            self._s1ap_wrapper._s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_IDENTITY_RESP, us_identity_resp,
            )
        elif msg.msg_type == s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value:
            self.int_ctx_setup_ind_count += 1
        elif msg.msg_type == s1ap_types.tfwCmd.UE_EMM_INFORMATION.value:
            self.ue_emm_information_count += 1
        elif msg.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value:
            self.ue_ctx_rel_ind_count += 1
        else:
            print("Unhandled msg type", msg.msg_type)

    def send_attach_req(self, ue_id):
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        self.attach_req_sent_count += 1

    def send_ue_detach(self, ue_id):
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = ue_id
        detach_req.ueDetType = (
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
        )
        self.detach_req_sent_count += 1

    def handle_detach_timer(self, ue_state):
        print("Detaching ue_id", ue_state.ue_id)
        self.send_ue_detach(ue_state.ue_id)
        attachTime = random.uniform(self.attach_delay_t0, self.attach_delay_t1)
        ue_state.attachTimer = threading.Timer(
            attachTime, self.handle_attach_timer, args=(ue_state,),
        )
        ue_state.attachTimer.start()

    def handle_attach_timer(self, ue_state):
        print("Attaching ue_id", ue_state.ue_id)
        attachDuration = random.uniform(
            self.attach_duration_t0, self.attach_duration_t1,
        )
        ue_state.detachTimer = threading.Timer(
            attachDuration, self.handle_detach_timer, args=(ue_state,),
        )
        ue_state.detachTimer.start()
        self.send_attach_req(ue_state.ue_id)

    def start_ue(self, ue_state):
        attachTime = random.uniform(self.attach_delay_t0, self.attach_delay_t1)
        ue_state.attachTimer = threading.Timer(
            attachTime, self.handle_attach_timer, args=(ue_state,),
        )
        ue_state.attachTimer.start()

    def hadle_end_timer(self):
        self.test_ended = True

    class UeState(object):
        def __init__(self, ue_id):
            self.ue_id = ue_id
            self.attachTimer = threading.Timer(1, None)
            self.detachTimer = threading.Timer(1, None)

    def test_continuous_random_attach(self):
        test_duration = 30
        num_ues = 100

        # These specify the attach rate as well as the duration of each attach
        # Actual value is uniformly distributed between t0 and t1
        self.attach_delay_t0 = 1
        self.attach_delay_t1 = 10
        self.attach_duration_t0 = 10
        self.attach_duration_t1 = 15

        self.ue_state_store = []
        self._s1ap_wrapper.configUEDevice(num_ues)
        self.test_ended = False

        # Collect some stats
        self.attach_req_sent_count = 0
        self.detach_req_sent_count = 0
        self.auth_req_ind_count = 0
        self.sec_mod_cmd_ind_count = 0
        self.attach_accept_ind_count = 0
        self.identity_req_ind_count = 0
        self.int_ctx_setup_ind_count = 0
        self.ue_emm_information_count = 0
        self.ue_ctx_rel_ind_count = 0

        for index in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            ue_state = self.UeState(req.ue_id)
            self.ue_state_store.append(ue_state)
            self.start_ue(ue_state)

        # Schedule test end
        endTimer = threading.Timer(test_duration, self.hadle_end_timer)
        endTimer.start()

        while True:
            if self.test_ended:
                break
            response = self._s1ap_wrapper.s1_util.get_response()
            self.handle_msg(response)

        # stop all active UE timers
        for index in range(num_ues):
            ue_state = self.ue_state_store[index]
            ue_state.attachTimer.cancel()
            ue_state.detachTimer.cancel()

        print("==============    stats   ======================")
        print("attach_req_sent_count: ", self.attach_req_sent_count)
        print("auth_req_ind_count: ", self.auth_req_ind_count)
        print("sec_mod_cmd_ind_count: ", self.sec_mod_cmd_ind_count)
        print("attach_accept_ind_count: ", self.attach_accept_ind_count)
        print("identity_req_ind_count: ", self.identity_req_ind_count)
        print("int_ctx_setup_ind_count: ", self.int_ctx_setup_ind_count)
        print("ue_emm_information_count: ", self.ue_emm_information_count)
        print("detach_req_sent_count: ", self.detach_req_sent_count)
        print("ue_ctx_rel_ind_count: ", self.ue_ctx_rel_ind_count)


if __name__ == "__main__":
    unittest.main()
