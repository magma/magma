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


class TestX2HandOver(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_x2_handover(self):
        """ Multi Enb Multi UE attach detach """

        """ Note: Before execution of this test case,
        Run the test script s1aptests/test_modify_mme_config_for_sanity.py
        to update multiple PLMN/TAC configuration in MME and
        after test case execution, restore the MME configuration by running
        the test script s1aptests/test_restore_mme_config_after_sanity.py

        Or

        Make sure that following steps are correct
        1. Configure same plmn and tac in both MME and s1ap tester
        2. How to configure plmn and tac in MME:
           a. Set mcc and mnc in gateway.mconfig for mme service
           b. Set tac in gateway.mconfig for mme service
           c. Restart MME service
        3. How to configure plmn and tac in s1ap tester,
           a. For multi-eNB test case, configure plmn and tac from test case.
             In each multi-eNB test case, set plmn, plmn length and tac
             in enb_list
           b. For single eNB test case, configure plmn and tac in nbAppCfg.txt
        """

        # column is an enb parameter, row is a number of enb
        """         Cell Id, Tac, EnbType, PLMN Id, PLMN length """
        enb_list = [[1, 1, 1, "00101", 5], [2, 2, 1, "00101", 5]]

        self._s1ap_wrapper.multiEnbConfig(len(enb_list), enb_list)

        time.sleep(2)
        """ Attach to Src eNB.HO to TeNB  """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for UE id ",
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

        time.sleep(3)
        print("************************* Sending ENB_CONFIGURATION_TRANSFER")
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.ENB_CONFIGURATION_TRANSFER, req,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.MME_CONFIGURATION_TRANSFER.value,
        )

        print("************************* Received MME_CONFIGURATION_TRANSFER")
        print("************************* Sending ENB_CONFIGURATION_TRANSFER")
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type,
            s1ap_types.tfwCmd.MME_CONFIGURATION_TRANSFER.value,
        )

        print("************************* Received MME_CONFIGURATION_TRANSFER")
        print("************************* Sending X2_HO_TRIGGER_REQ")

        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.X2_HO_TRIGGER_REQ, req,
        )
        # Receive Path Switch Request Ack
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.PATH_SW_REQ_ACK.value,
        )

        print("************************* Received Path Switch Request Ack")

        print(
            "************************* Running UE detach for UE id ", req.ue_id,
        )
        # Now detach the UE
        time.sleep(3)
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            False,
        )


if __name__ == "__main__":
    unittest.main()
