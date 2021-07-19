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


class TestAttachEmergency(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def _test_attach_response_for_id_type(self, id_type, expected_ue_state):
        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running Emergency Attach for ",
            "UE id ", req.ue_id,
        )

        # Attach
        msg = self._s1ap_wrapper._s1_util.attach(
            req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND, s1ap_types.ueAttachFail_t,
            eps_type=s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_EMRG_ATTACH,
            id_type=id_type,
        )

        # Assert cause
        self.assertEqual(msg.ueState, expected_ue_state)
        print(
            "************************* Emergency attach rejection successful",
            "UE id ", req.ue_id,
        )
        # Context release, remove from queue
        self._s1ap_wrapper._s1_util.get_response()

    def test_attach_emergency(self):
        """ Send UE attach emergency for IMSI and IMEI """
        # Initialize
        self._s1ap_wrapper.configIpBlock()
        self._s1ap_wrapper.configUEDevice(2)

        # Test IMSI, should return cause not authorized
        expected_ue_state = 35  # EMM_CAUSE_NOT_AUTHORIZED_IN_PLMN
        self._test_attach_response_for_id_type(
            s1ap_types.TFW_MID_TYPE_IMSI,
            expected_ue_state,
        )

        # Test IMEI, should return IMEI not accepted
        expected_ue_state = 5  # EMM_CAUSE_IMEI_NOT_ACCEPTED
        self._test_attach_response_for_id_type(
            s1ap_types.TFW_MID_TYPE_IMEI,
            expected_ue_state,
        )


if __name__ == "__main__":
    unittest.main()
