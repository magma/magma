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

import gpp_types
import s1ap_types
import s1ap_wrapper


class TestAttachServiceMultiUe(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_service_multi_ue(self):
        """
        Test with multi-UE attach, UE context release, service request
        """
        num_ues = 32
        self._s1ap_wrapper.configUEDevice(num_ues)
        reqs = tuple(self._s1ap_wrapper.ue_req for _ in range(num_ues))

        for req in reqs:
            print(
                "************************* Running End to End attach for UE ",
                "id ", req.ue_id,
            )
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        time.sleep(0.5)

        for req in reqs:
            ue_id = req.ue_id
            print(
                "************************* Sending UE context release "
                "request for UE id ", ue_id,
            )
            # Send UE context release request to move UE to idle mode
            req = s1ap_types.ueCntxtRelReq_t()
            req.ue_Id = ue_id
            req.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )

        for req in reqs:
            ue_id = req.ue_id
            print(
                "************************* Sending Service request for UE "
                "id ", ue_id,
            )
            # Send service request to reconnect UE
            req = s1ap_types.ueserviceReq_t()
            req.ue_Id = ue_id
            req.ueMtmsi = s1ap_types.ueMtmsi_t()
            req.ueMtmsi.pres = False
            req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST, req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
            )

        for req in reqs:
            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                True,
            )


if __name__ == "__main__":
    unittest.main()
