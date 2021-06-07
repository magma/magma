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

import orc8r.protos.metricsd_pb2 as metricsd
import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestAttachWithMultipleMmeRestarts(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_with_multiple_mme_restarts(self):
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)
        print("************************* Sending Attach Request for ue-id : 1")
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = 1
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("************************* Received UE_AUTH_REQ_IND")

        # Try consecutive mme restarts
        self._s1ap_wrapper.magmad_util.restart_mme_and_wait()
        self._s1ap_wrapper.magmad_util.restart_mme_and_wait()

        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = 1
        sqn_recvd = s1ap_types.ueSqnRcvd_t()
        sqn_recvd.pres = 0
        auth_res.sqnRcvd = sqn_recvd
        print("************************* Sending Auth Response ue-id", auth_res.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )
        print("************************* Received UE_SEC_MOD_CMD_IND")

        self._s1ap_wrapper.magmad_util.restart_mme_and_wait()

        print("************************* Sending UE_SEC_MOD_COMPLETE")
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = 1
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )
        print("************************* Received INT_CTX_SETUP_IND")

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
        )
        print("************************* Received UE_ATTACH_ACCEPT_IND")

        self._s1ap_wrapper.magmad_util.restart_mme_and_wait()

        # Trigger Attach Complete
        print("************************* Sending UE_ATTACH_COMPLETE")
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = 1
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        print("************************* Running UE detach")
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = 1
        detach_req.ueDetType = (
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
        )

        # Wait for UE context release command
        return self._s1ap_wrapper.s1_util.get_response()


if __name__ == "__main__":
    unittest.main()
