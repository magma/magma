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
import s1ap_wrapper


class TestSendErrorIndForDlNasWithAuthReq(unittest.TestCase):
    """Test sending of error indication for DL NAS message
    carrying authentication request
    """

    def setUp(self):
        """Initialize"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        """Cleanup"""
        self._s1ap_wrapper.cleanup()

    def test_send_error_ind_for_dl_nas_with_auth_req(self):
        """Send error indication after receiving authentication request"""
        self._s1ap_wrapper.configIpBlock()
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req

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
        print("************************* Sent attach request")
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )
        print("************************* Received authentication request")

        # Send error indication
        error_ind = s1ap_types.fwNbErrIndMsg_t()
        # isUeAssoc flag to include optional MME_UE_S1AP_ID and eNB_UE_S1AP_ID
        error_ind.isUeAssoc = True
        error_ind.ue_Id = req.ue_id
        error_ind.cause.pres = True
        # Radio network causeType = 0
        error_ind.cause.causeType = 0
        # causeVal - Unknown-pair-ue-s1ap-id
        error_ind.cause.causeVal = 15
        print("*** Sending error indication ***")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.ENB_ERR_IND_MSG, error_ind,
        )

        # Context release
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print("************************* Received UE_CTX_REL_IND")


if __name__ == "__main__":
    unittest.main()
