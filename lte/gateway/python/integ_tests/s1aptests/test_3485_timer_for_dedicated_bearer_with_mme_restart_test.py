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
import time

from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import SpgwUtil
from s1ap_utils import MagmadUtil


class Test3485TimerForDedicatedBearerWithMmeRestart(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE
        )
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_3485_timer_for_dedicated_bearer_with_mme_restart(self):
        """ Test case validates the functionality of 3485 timer for
            Dedicated bearer while MME restarts
        Step1: UE attaches to network
        Step2: Send an indication to S1ap stack to drop E-Rab Setup
               Request message, sent as part of secondary PDN activation
               procedure.
        Step3: Send an indication to initiate Dedicated bearer activation
        Step4: Send an indication to S1ap stack to not to drop E-Rab Setup
               Request message, so that re-transmitted message reaches to
               TFW
        Step5: Send command to Magma to restart mme service
        Step6: TFW shall receive re-transmitted Actiavte Dedicated EPS
               Bearer Context Request and send Actiavte Dedicated EPS Bearer
               Context Accept
        Step7: TFW shall initaite de-activation of deidicated bearer and then
               initiate Detach procedure.
        """

        num_ues = 1
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        wait_for_s1 = [True, False]
        self._s1ap_wrapper.configUEDevice(num_ues)

        for i in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "********************** Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()

            time.sleep(2)
            print(
                "********************** Adding dedicated bearer to IMSI",
                "".join([str(i) for i in req.imsi]),
            )

            self._spgw_util.create_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), 5
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            print(
                "*******************Received first Activate dedicated bearer "
                "request"
            )

            print("***** Restarting MME service on gateway")
            self._s1ap_wrapper.magmad_util.restart_services(["mme"])

            for j in range(20):
                print("Waiting for", j, "seconds")
                time.sleep(1)

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            )
            print(
                "***** Received re-transmitted Activate dedicated bearer req"
            )

            act_ded_ber_ctxt_req = response.cast(
                s1ap_types.UeActDedBearCtxtReq_t
            )

            self._s1ap_wrapper.sendActDedicatedBearerAccept(
                req.ue_id, act_ded_ber_ctxt_req.bearerId
            )
            print("*******************Send Activate dedicated bearer accept")
            time.sleep(4)

            print(
                "********************** Deleting dedicated bearer for IMSI",
                "".join([str(i) for i in req.imsi]),
            )
            self._spgw_util.delete_bearer(
                "IMSI" + "".join([str(i) for i in req.imsi]), 5, 6
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            while (
                response.msg_type == s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value
            ):
                print("Ignore re-transmitted Activate dedicate bearer message")
                response = self._s1ap_wrapper.s1_util.get_response()

            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_DEACTIVATE_BER_REQ.value,
            )

            print("******************* Received deactivate eps bearer context")

            deactv_bearer_req = response.cast(s1ap_types.UeDeActvBearCtxtReq_t)
            self._s1ap_wrapper.sendDeactDedicatedBearerAccept(
                req.ue_id, deactv_bearer_req.bearerId
            )

            time.sleep(5)
            print(
                "********************** Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id, detach_type[i], wait_for_s1[i]
            )


if __name__ == "__main__":
    unittest.main()
