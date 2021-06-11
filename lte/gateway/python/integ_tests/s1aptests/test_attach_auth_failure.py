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

import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper


class TestAuthFailure(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_auth_failure_proc(self):
        """ Testing of sending authentication failure procedure """
        num_ues = 1

        self._s1ap_wrapper.configUEDevice(num_ues)
        print("************************* sending Attach Request for ue-id : 1")
        attach_req = s1ap_types.ueAttachRequest_t()
        attach_req.ue_Id = 1
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        pdn_type.pdn_type = 3
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type

        print("Sending Attach Request for ue-id", attach_req.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertTrue(response, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value)
        print("Received Authentication Request message ")

        auth_failure = s1ap_types.ueAuthFailure_t()
        auth_failure.ue_Id = 1
        auth_failure.cause = s1ap_types.TFW_EMM_CAUSE_SYNC_FAILURE
        # sending random/zero auts value to simulate failure scenario
        for idx1 in range(14):
            auth_failure.auts[idx1] = 0
        print("Sending Authentication Failure for ue-id", auth_failure.ue_Id)
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_FAILURE, auth_failure,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertTrue(response, s1ap_types.tfwCmd.UE_AUTH_REJ_IND.value)
        print("Received Authentication Reject message")


if __name__ == "__main__":
    unittest.main()
