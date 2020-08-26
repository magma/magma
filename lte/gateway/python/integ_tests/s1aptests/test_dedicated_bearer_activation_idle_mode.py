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
import time

import gpp_types
import s1ap_types
import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil


class TestDedicatedBearerActivationIdleMode(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_dedicated_bearer_activation_idle_mode(self):
        """
        Test with a single UE attach + UE context release +
        trigger dedicated bearer activation
        + Page the UE + send dedicated bearer activation + detach"""
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        imsi = req.imsi
        print("*********** Running End to End attach for UE id ", ue_id)
        # Now actually complete the attach
        attach = self._s1ap_wrapper._s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        print("*********** Sleeping for 5 seconds")
        time.sleep(5)

        print("*********** Moving UE to idle mode")
        print(
            "*********** Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        req = s1ap_types.ueCntxtRelReq_t()
        req.ue_Id = ue_id
        req.cause.causeVal = gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, req
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value
        )

        print("*********** Sleeping for 5 seconds")
        time.sleep(5)
        # Add dedicated bearer for default bearer 5
        print(
            "************* Adding dedicated bearer to magma.ipv4"
            " PDN in idle mode"
        )
        self._spgw_util.create_bearer(
            "IMSI" + "".join([str(i) for i in imsi]),
            attach.esmInfo.epsBearerId,
        )

        # Add 2nd dedicated bearer for default bearer 5
        print(
            "************* Adding dedicated bearer to magma.ipv4"
            " PDN in idle mode"
        )
        self._spgw_util.create_bearer(
            "IMSI" + "".join([str(i) for i in imsi]),
            attach.esmInfo.epsBearerId,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PAGING_IND.value
        )

        print("*********** Received Paging for UE id ", ue_id)
        print("*********** Sending Service request for UE id ", ue_id)
        # Send service request to reconnect UE
        req = s1ap_types.ueserviceReq_t()
        req.ue_Id = ue_id
        req.ueMtmsi = s1ap_types.ueMtmsi_t()
        req.ueMtmsi.pres = False
        req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SERVICE_REQUEST, req
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value
        )

        print("*********** Received ICS Request for UE id ", ue_id)
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        )
        act_ded_ber_req_oai_apn1 = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t
        )
        print(
            "*********** Received Activate dedicated EPS bearer"
            " req for bearer ",
            act_ded_ber_req_oai_apn1.bearerId,
        )
        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            ue_id, act_ded_ber_req_oai_apn1.bearerId
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
        )
        act_ded_ber_req_oai_apn2 = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t
        )
        print(
            "*********** Received Activate dedicated EPS bearer"
            " req for bearer ",
            act_ded_ber_req_oai_apn2.bearerId,
        )

        self._s1ap_wrapper.sendActDedicatedBearerAccept(
            ue_id, act_ded_ber_req_oai_apn2.bearerId
        )
        print("*********** Sleeping for 5 seconds")
        time.sleep(5)
        print("*********** Running UE detach for UE id ", ue_id)

        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, True
        )


if __name__ == "__main__":
    unittest.main()
