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

import math
import time
import unittest

import gpp_types
import s1ap_types
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import HaUtil


class TestAgwOffloadMixedIdleActiveMultiUe(unittest.TestCase):
    """Unittest: TestAgwOffloadMixedIdleActiveMultiUe"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._ha_util = HaUtil()

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_agw_offload_mixed_idle_active_multiue(self):
        """Test case to offload 100 UEs in both active and idle states

        NOTE: HA service must be enabled for running this test case. Set the
        parameter 'use_ha' in configuration file /etc/magma/mme.yml to 'true'
        on magma-dev VM and restart MME to enable the HA service
        """
        num_ues = 100
        # column is a enb parameter,  row is a number of enbs
        # column description:
        #     1.Cell Id, 2.Tac, 3.EnbType, 4.PLMN Id 5. PLMN length
        enb_list = [(1, 1, 1, "00101", 5)]
        self._s1ap_wrapper.multiEnbConfig(len(enb_list), enb_list)

        time.sleep(2)
        self._s1ap_wrapper.configUEDevice(num_ues)

        ue_ids = []
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for ",
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
            ue_ids.append(req.ue_id)

        # Send UE context release request for half of UEs to move them
        # to ECM-IDLE state
        for i in range(math.floor(num_ues / 2)):
            ue_cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
            ue_cntxt_rel_req.ue_Id = ue_ids[i]
            ue_cntxt_rel_req.cause.causeVal = (
                gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
            )
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST,
                ue_cntxt_rel_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )

        print("*************************  Send Offload Request to AGW")
        # Send offloading request
        self.assertTrue(self._ha_util.offload_agw(None, enb_list[0][0]))

        # All UEs should eventually receive Context Release Request
        # The first half should get it immediately
        # The second half should first get paging
        for _ in range(num_ues):
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertIn(
                response.msg_type,
                [
                    s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
                    s1ap_types.tfwCmd.UE_PAGING_IND.value,
                ],
                "Not a paging or ue context release message",
            )

        # Send service request as paging response
        for i in range(math.floor(num_ues / 2)):
            # Send service request to reconnect UE
            # Auto-release should happen
            ser_req = s1ap_types.ueserviceReq_t()
            ser_req.ue_Id = ue_ids[i]
            ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
            ser_req.ueMtmsi.pres = False
            ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST,
                ser_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
            )

            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )

        for i in range(num_ues):
            # Send service request again:
            ser_req = s1ap_types.ueserviceReq_t()
            ser_req.ue_Id = ue_ids[i]
            ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
            ser_req.ueMtmsi.pres = False
            ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MO_DATA.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST,
                ser_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type,
                s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
            )

        # Now detach the UEs normally
        for ue in ue_ids:
            print("************************* Running UE detach for UE id ", ue)
            self._s1ap_wrapper.s1_util.detach(
                ue,
                s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
                wait_for_s1_ctxt_release=True,
            )


if __name__ == "__main__":
    unittest.main()
