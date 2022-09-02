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
import unittest

import s1ap_types
from integ_tests.s1aptests.s1ap_utils import S1ApUtil


class TestS1SetupSuccessSecondaryPLMN(unittest.TestCase):
    def setUp(self):
        self._s1_util = S1ApUtil()

    def tearDown(self):
        print("************************* Sending SCTP SHUTDOWN")
        self._s1_util.issue_cmd(s1ap_types.tfwCmd.SCTP_SHUTDOWN_REQ, None)
        self._s1_util.cleanup()

    def test_s1setup_success_secondary_plmn(self):
        """ S1 Setup with multiple bPLMN IDs, the second being valid. """

        print("************************* Enb tester configuration")
        req = s1ap_types.FwNbConfigReq_t()

        req.cellId_pr.pres = True
        req.cellId_pr.cell_id = 1

        req.plmnId_pr.pres = True
        req.plmnId_pr.plmn_id = (ctypes.c_ubyte * 6).from_buffer_copy(
            bytearray(b"333333"),
        )

        req.suppTAs.pres = True
        req.suppTAs.numTAs = 1
        req.suppTAs.suppTA[0].tac = 1
        req.suppTAs.suppTA[0].bPlmnList.numBPlmns = 2

        # 333333 - invalid MCC/MNC
        req.suppTAs.suppTA[0].bPlmnList.bPlmn[0].numMncDigits = 3
        req.suppTAs.suppTA[0].bPlmnList.bPlmn[0].mcc = (ctypes.c_ubyte * 3).from_buffer_copy(
            bytearray(b"\x03\x03\x03"),
        )
        req.suppTAs.suppTA[0].bPlmnList.bPlmn[0].mnc = (ctypes.c_ubyte * 3).from_buffer_copy(
            bytearray(b"\x03\x03\x03"),
        )

        # 00101 - Valid MCC/MNC
        req.suppTAs.suppTA[0].bPlmnList.bPlmn[1].numMncDigits = 2
        req.suppTAs.suppTA[0].bPlmnList.bPlmn[1].mcc = (ctypes.c_ubyte * 3).from_buffer_copy(
            bytearray(b"\x00\x00\x01"),
        )
        req.suppTAs.suppTA[0].bPlmnList.bPlmn[1].mnc = (ctypes.c_ubyte * 3).from_buffer_copy(
            bytearray(b"\x00\x01\x00"),
        )

        print("************************* Sending ENB configuration Request")
        assert self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_CONFIG, req) == 0
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_CONFIG_CONFIRM.value
        res = response.cast(s1ap_types.FwNbConfigCfm_t)
        assert res.status == s1ap_types.CfgStatus.CFG_DONE.value

        print("************************* Sending S1-setup Request")
        req = None
        assert (
            self._s1_util.issue_cmd(s1ap_types.tfwCmd.ENB_S1_SETUP_REQ, req)
            == 0
        )
        response = self._s1_util.get_response()
        assert response.msg_type == s1ap_types.tfwCmd.ENB_S1_SETUP_RESP.value
        res = response.cast(s1ap_types.FwNbS1setupRsp_t)
        assert res.res == s1ap_types.S1_setp_Result.S1_SETUP_SUCCESS.value


if __name__ == "__main__":
    unittest.main()
