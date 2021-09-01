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


class TestAttachIpv4v6PdnType(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_ipv4v6_pdn_type(self):
        """ Test Attach for the UEs that are dual IP stack IPv4v6
            capable """
        # Set PDN TYPE to IPv4V6 i.e. 3. IPV4 is equal to 1
        resp_ipv4_ipv6 = self._create_attach_ipv4v6_pdn_type_req(
            pdn_type_value=3,
        )
        self.assertEqual(
            resp_ipv4_ipv6.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        # IPv6 is equal to 2
        resp_ipv6 = self._create_attach_ipv4v6_pdn_type_req(pdn_type_value=2)
        self.assertEqual(
            resp_ipv6.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )

    def _create_attach_ipv4v6_pdn_type_req(self, pdn_type_value):
        # Ground work.
        self._s1ap_wrapper.configUEDevice(1)
        ue_req = self._s1ap_wrapper.ue_req
        # Trigger Attach Request with PDN_Type = IPv4v6
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        # Set PDN TYPE to IPv4V6 i.e. 3. IPV4 is equal to 1
        # IPV6 is equal to 2 in value
        pdn_type.pdn_type = pdn_type_value
        attach_req.ue_Id = ue_req.ue_id
        attach_req.mIdType = id_type
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type

        print(
            "********Triggering Attach Request with PDN Type IPv4v6 test, "
            "pdn_type_value",
            pdn_type_value,
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_AUTH_REQ_IND.value,
        )

        # Trigger Authentication Response
        auth_res = s1ap_types.ueAuthResp_t()
        auth_res.ue_Id = ue_req.ue_id
        sqnRecvd = s1ap_types.ueSqnRcvd_t()
        sqnRecvd.pres = 0
        auth_res.sqnRcvd = sqnRecvd
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_AUTH_RESP, auth_res,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_SEC_MOD_CMD_IND.value,
        )

        # Trigger Security Mode Complete
        sec_mode_complete = s1ap_types.ueSecModeComplete_t()
        sec_mode_complete.ue_Id = ue_req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_SEC_MOD_COMPLETE, sec_mode_complete,
        )
        # Attach Reject will be sent since IPv6 PDN Type is not configured
        if pdn_type_value == 2:
            response = self._s1ap_wrapper.s1_util.get_response()
            self.assertEqual(
                response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND.value,
            )
            return self._s1ap_wrapper.s1_util.get_response()

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.INT_CTX_SETUP_IND.value,
        )
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND.value,
        )

        # Trigger Attach Complete
        attach_complete = s1ap_types.ueAttachComplete_t()
        attach_complete.ue_Id = ue_req.ue_id
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_COMPLETE, attach_complete,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()
        print("************************* Running UE detach")
        # Now detach the UE
        detach_req = s1ap_types.uedetachReq_t()
        detach_req.ue_Id = ue_req.ue_id
        detach_req.ueDetType = (
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_DETACH_REQUEST, detach_req,
        )

        # Wait for UE context release command
        return self._s1ap_wrapper.s1_util.get_response()


if __name__ == "__main__":
    unittest.main()
