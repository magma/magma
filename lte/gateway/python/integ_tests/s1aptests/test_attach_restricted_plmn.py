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


import ctypes
import time
import unittest

import s1ap_types
import s1ap_wrapper


class TestAttachRestrictedPlmn(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_restricted_plmn(self):
        """
        If this TC is executed individually run
        test_modify_mme_config_for_sanity.py to add
        PLMN - { MCC="123" ; MNC="450";}
        under the RESTRICTED_PLMN_LIST
        in mme.conf.template.

        This TC does the following:
        1. Send attach request with IMSI containing the PLMN
           configured in mme.conf.template
        2. Verify that MME sends attach reject with cause(11) PLMN NOT ALLOWED
        3. Attach a 2nd UE with an allowed PLMN by invoking attach utility
           function
        4. Detach the UE

        After execution of this TC run test_restore_mme_config_after_sanity.py
        to restore the old mme.conf.template.
        """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        # Trigger Attach Request with restricted PLMN
        attach_req = s1ap_types.ueAttachRequest_t()
        sec_ctxt = s1ap_types.TFW_CREATE_NEW_SECURITY_CONTEXT
        id_type = s1ap_types.TFW_MID_TYPE_IMSI
        eps_type = s1ap_types.TFW_EPS_ATTACH_TYPE_EPS_ATTACH
        pdn_type = s1ap_types.pdn_Type()
        pdn_type.pres = True
        pdn_type.pdn_type = 1
        attach_req.ue_Id = req.ue_id
        attach_req.mIdType = id_type
        # Generate IMSI with prefix IMSI12345
        imsi = self._s1ap_wrapper._s1_util.generate_imsi(prefix="IMSI12345")
        for i in range(0, 15):
            attach_req.imsi[i] = ctypes.c_ubyte(int(imsi[4 + i]))
        attach_req.imsi_len = 15
        attach_req.epsAttachType = eps_type
        attach_req.useOldSecCtxt = sec_ctxt
        attach_req.pdnType_pr = pdn_type

        print(
            "************************* Sending attach for ",
            "UE id ",
            req.ue_id,
        )
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ATTACH_REQUEST, attach_req,
        )

        # Attach Reject
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_REJECT_IND.value,
        )

        attach_rej = response.cast(s1ap_types.ueAttachRejInd_t)
        print(
            "************************* Received attach reject for "
            "UE id %d with cause %d" % (req.ue_id, attach_rej.cause),
        )

        # Verify cause
        self.assertEqual(attach_rej.cause, s1ap_types.TFW_EMM_CAUSE_PLMN_NA)

        # Context release
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_CTX_REL_IND.value,
        )
        print(
            "************************* Received ue context release cmd for "
            "UE id ",
            req.ue_id,
        )

        # Attach the 2nd UE with allowed PLMN
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running End to End attach for "
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
        print("********************** Running UE detach for UE id ", req.ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
        )


if __name__ == "__main__":
    unittest.main()
