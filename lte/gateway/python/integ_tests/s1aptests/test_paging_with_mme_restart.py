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
from s1ap_utils import MagmadUtil


class TestPagingWithMmeRestart(unittest.TestCase):
    """Unittest: TestPagingWithMmeRestart"""

    def setUp(self):
        """Initialize before test case execution"""
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )

    def tearDown(self):
        """Cleanup after test case execution"""
        self._s1ap_wrapper.cleanup()

    def test_paging_with_mme_restart(self):
        """
        The test case validates resumption of Paging Response Timer
        on mme restart
        Step1 : UE attaches to network
        Step2 : UE moves to Idle state
        Step3 : Send DL data, on arrival of DL data mme sends Paging message
                and starts Paging Response timer, while timer is running
                mme restarts, once mme restarts, mme shall be able to resume
                Paging response timer
        Step4 : On expiry of Paging Response timer, mme shall re-send Paging
                message
        Step5 : In response to Paging message, UE sends Service Request message
        Step6 : Expecting normal flow of DL data
        """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper.s1_util.attach(
            ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        # Delay to ensure S1APTester sends attach complete before sending UE
        # context release
        time.sleep(0.5)

        print(
            "************************* Sending UE context release request ",
            "for UE id ",
            ue_id,
        )
        # Send UE context release request to move UE to idle mode
        ue_cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
        ue_cntxt_rel_req.ue_Id = ue_id
        ue_cntxt_rel_req.cause.causeVal = (
            gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
        )
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST,
            ue_cntxt_rel_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.UE_CTX_REL_IND.value

        time.sleep(0.3)
        print(
            "************************* Running UE downlink (UDP) for UE id ",
            ue_id,
        )
        with self._s1ap_wrapper.configDownlinkTest(
            req,
            duration=1,
            is_udp=True,
        ) as test:
            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_PAGING_IND.value
            print("************************ Received Paging Indication")

            print(
                "************************* Restarting MME service on",
                "gateway",
            )
            wait_for_restart = 30
            self._s1ap_wrapper.magmad_util.restart_services(
                ["mme"], wait_for_restart,
            )

            # Send service request to reconnect UE
            ser_req = s1ap_types.ueserviceReq_t()
            ser_req.ue_Id = ue_id
            ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
            ser_req.ueMtmsi.pres = False
            ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MT_ACCESS.value
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_SERVICE_REQUEST,
                ser_req,
            )

            # It has been observed that despite getting the restart command on
            # time, MME sometimes restarts after a delay of 5-6 seconds.
            # Currently Paging response timer is set to 4 seconds defined by
            # macro, MME_APP_PAGING_RESPONSE_TIMER_VALUE in the file
            # mme_app_ue_context.h. If MME restarts after expiry of paging
            # response timer of 4 sec, it will again send the Paging indication
            resp_count = 0
            while True:
                response = self._s1ap_wrapper.s1_util.get_response()

                if response.msg_type == s1ap_types.tfwCmd.UE_PAGING_IND.value:
                    resp_count += 1
                    print(
                        "******************** Ignoring re-transmitted (",
                        resp_count,
                        ") Paging indication",
                    )
                else:
                    break

            assert response.msg_type == s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value
            test.verify()

        time.sleep(0.5)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id,
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
        )


if __name__ == "__main__":
    unittest.main()
