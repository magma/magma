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


class TestNoAuthResponseWithMmeRestartReattach(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_no_auth_response_with_mme_restart_reattach(self):
        """ Send attach req + receive auth req + restart mme
        + MME does implicit detach + re-attach """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        print("************************* Running attach ")

        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = req.ue_id
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt

        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        print("************ Sent Attach Request for ue", req.ue_id)

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("************ Received Auth Req for ue", req.ue_id)

        print("************************* Restarting MME service on", "gateway")
        self._s1ap_wrapper.magmad_util.restart_services(["mme"])

        for j in range(30):
            print("Waiting for", j, "seconds")
            time.sleep(1)

        # Wait for UE context release command
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("****** Received UE_CTX_REL_IND for ue", req.ue_id)
        print("****** Triggering end-end attach after mme restart *********")

        self._s1ap_wrapper.s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            False,
        )


if __name__ == "__main__":
    unittest.main()
