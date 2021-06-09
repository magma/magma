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


class TestMultiEnbPagingRequest(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_multi_enb_paging_request(self):
        """ Multi Enb + attach 1 UE + UE enters into idle mode + paging +
            service req + detach
        """

        """ Note: Before execution of this test case,
        make sure that following steps are correct
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

        # column is an enb parameter, row is number of enbs
        """ Cell Id, Tac, EnbType, PLMN Id, PLMN length """
        enb_list = [
            (1, 1, 1, "00101", 5),
            (2, 1, 1, "00101", 5),
            (3, 1, 1, "00101", 5),
            (4, 1, 1, "00101", 5),
            (5, 1, 1, "00101", 5),
        ]

        self._s1ap_wrapper.multiEnbConfig(len(enb_list), enb_list)

        time.sleep(2)

        ue_ids = []
        # UEs will attach to the ENBs in a round-robin fashion
        # each ENBs will be connected with 32UEs
        num_ues = 1
        # Ground work.
        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print("******************** Calling attach for UE id ", req.ue_id)
            self._s1ap_wrapper.s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
            )
            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            ue_ids.append(req.ue_id)
            # Delay to ensure S1APTester sends attach complete before
            # sending UE context release
            time.sleep(0.5)
            print(
                "*********************  Sending UE context release request ",
                "for UE id ",
                req.ue_id,
            )
            # Send UE context release request to move UE to idle mode
            ue_cntxt_rel_req = s1ap_types.ueCntxtRelReq_t()
            ue_cntxt_rel_req.ue_Id = req.ue_id
            ue_cntxt_rel_req.cause.causeVal = (
                gpp_types.CauseRadioNetwork.USER_INACTIVITY.value
            )
            self._s1ap_wrapper.s1_util.issue_cmd(
                s1ap_types.tfwCmd.UE_CNTXT_REL_REQUEST, ue_cntxt_rel_req,
            )
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
            )
            time.sleep(0.5)
            print(
                "********************* Running UE downlink (UDP) for UE id ",
                req.ue_id,
            )
            with self._s1ap_wrapper.configDownlinkTest(
                req, duration=1, is_udp=True,
            ) as test:
                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertTrue(
                    response, s1ap_types.tfwCmd.UE_PAGING_IND.value,
                )

                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertTrue(
                    response, s1ap_types.tfwCmd.UE_PAGING_IND.value,
                )

                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertTrue(
                    response, s1ap_types.tfwCmd.UE_PAGING_IND.value,
                )

                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertTrue(
                    response, s1ap_types.tfwCmd.UE_PAGING_IND.value,
                )

                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertTrue(
                    response, s1ap_types.tfwCmd.UE_PAGING_IND.value,
                )
                # Send service request to reconnect UE
                ser_req = s1ap_types.ueserviceReq_t()
                ser_req.ue_Id = req.ue_id
                ser_req.ueMtmsi = s1ap_types.ueMtmsi_t()
                ser_req.ueMtmsi.pres = False
                ser_req.rrcCause = s1ap_types.Rrc_Cause.TFW_MT_ACCESS.value
                self._s1ap_wrapper.s1_util.issue_cmd(
                    s1ap_types.tfwCmd.UE_SERVICE_REQUEST, ser_req,
                )
                response = self._s1ap_wrapper.s1_util.get_response()
                self.assertEqual(
                    response.msg_type,
                    s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
                )
                test.verify()
        time.sleep(0.5)
        # Now detach the UE
        for ue in ue_ids:
            print("************************* Calling detach for UE id ", ue)
            self._s1ap_wrapper.s1_util.detach(
                ue, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
            )

        time.sleep(1)


if __name__ == "__main__":
    unittest.main()
