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


class TestAttachDetachNoUeContextReleseComp(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_no_ueContRelComp(self):
        """ Basic attach/detach test with a single UE - SCTP Abort"""

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print(
            "************************ Running End to End attach for ",
            "UE id ", ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        print(
            "************************* Running UE detach for ",
            "UE id", ue_id,
        )
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = req.ue_id
        detach_req.ueDetType = s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
        )
        response = self._s1ap_wrapper._s1_util.get_response()
        assert (
            s1ap_types.tfwCmd.UE_DETACH_ACCEPT_IND.value
            == response.msg_type
        )
        time.sleep(0.1)
        # Send SCTP ABORT to MME
        sctp_abort = s1ap_types.FwSctpAbortReq_t()
        sctp_abort.cause = 0
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.SCTP_ABORT_REQ, sctp_abort,
        )


if __name__ == "__main__":
    unittest.main()
