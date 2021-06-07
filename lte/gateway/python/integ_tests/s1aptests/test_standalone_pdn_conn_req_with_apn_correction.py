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
import threading
import time
import unittest

import gpp_types
import s1ap_types
import s1ap_wrapper
from s1ap_utils import MagmadUtil


class TestStandAlonePdnConnReqWithApnCorrection(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
                apn_correction=MagmadUtil.apn_correction_cmds.ENABLE,
        )

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_standalone_pdn_conn_req_with_apn_correction(self):
        """ Attach a single UE and send standalone PDN Connectivity
        Request with wrong APN and use APN correction feature to
        override it"""
        num_ues = 1

        print('************************* restarting mme')
        self._s1ap_wrapper.magmad_util.restart_services(['mme'])
        for j in range(15):
            print("Waiting mme restart for", j, "seconds")
            time.sleep(1)

        self._s1ap_wrapper.configUEDevice(num_ues)
        req = self._s1ap_wrapper.ue_req
        ue_id = req.ue_id
        print(
            "************************* Running End to End attach for UE id ",
            ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper.s1_util.attach(
            ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print(
            "************************* Sending PDN Connectivity Request ",
            "for UE id ", ue_id,
        )
        req = s1ap_types.uepdnConReq_t()
        req.ue_Id = ue_id
        # Request type = Initial Request
        req.reqType = 1
        req.pdnType_pr.pres = 1
        # PDN Type = IPv4
        req.pdnType_pr.pdn_type = 1
        req.pdnAPN_pr.pres = 1
        s = 'internet.mnc012.mcc345.gprs'
        req.pdnAPN_pr.len = len(s)
        req.pdnAPN_pr.pdn_apn = (ctypes.c_ubyte * 100)(*[ctypes.c_ubyte(ord(c)) for c in s[:100]])
        self._s1ap_wrapper.s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_PDN_CONN_REQ, req,
        )
        # Receive PDN Connectivity Reject
        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_PDN_CONN_RSP_IND.value,
        )
        # self.assertEqual(
        #    response.msg_type, s1ap_types.tfwCmd.UE_ATTACH_FAIL_IND.value)

        print("Received PDN CONNECTIVITY REJECT")
        print(
            "************************* Running UE detach (switch-off) for ",
            "UE id ", ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value, False,
        )

        # Disable APN correction
        self._s1ap_wrapper.magmad_util.config_apn_correction(
                MagmadUtil.apn_correction_cmds.DISABLE,
        )
        self._s1ap_wrapper.magmad_util.restart_services(['mme'])
        for j in range(10):
            print("Waiting mme restart for", j, "seconds")
            time.sleep(1)


if __name__ == "__main__":
    unittest.main()
