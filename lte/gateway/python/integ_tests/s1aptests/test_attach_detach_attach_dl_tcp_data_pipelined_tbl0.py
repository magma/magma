"""
Copyright 2022 The Magma Authors.

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
from integ_tests.s1aptests.s1ap_utils import MagmadUtil


class TestAttachDetachAttachDlTcpDataPipelinedTbl0(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_attach_dl_tcp_data_pipelined_tbl0(self):
        """ Attach, detach, re-attach, send DL TCP data, and detach """
        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Running 1st End to End attach for ",
            "UE id ", req.ue_id,
        )

        # Enable the enable5g_features flag to pipelined
        self._s1ap_wrapper.magmad_util.config_enable5g_features(
                MagmadUtil.enable5g_features_cmds.ENABLE,
        )

        # Enable the flag to control ovs table=0 using pipelined
        self._s1ap_wrapper.magmad_util.config_pipelined_managed_tbl0(
                MagmadUtil.pipelined_managed_tbl0_cmds.ENABLE,
        )
        # To reflect above config restart pipelined service is mandatory
        self._s1ap_wrapper.magmad_util.restart_services(["pipelined"])

        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print(
            "************************* Running UE detach for UE id ",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            True,
        )

        print(
            "************************* Running 2nd End to End attach for ",
            "UE id ", req.ue_id,
        )
        # Now actually complete the attach
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print(
            "************************* Running UE downlink (TCP) for UE id ",
            req.ue_id,
        )
        with self._s1ap_wrapper.configDownlinkTest(req, duration=1) as test:
            test.verify()

        print(
            "************************* Running UE detach for UE id ",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
            True,
        )

        # Disable the enable5g_features flag to pipelined
        self._s1ap_wrapper.magmad_util.config_enable5g_features(
                MagmadUtil.enable5g_features_cmds.DISABLE,
        )

        # Disable the flag to control ovs table=0 using pipelined
        self._s1ap_wrapper.magmad_util.config_pipelined_managed_tbl0(
                MagmadUtil.pipelined_managed_tbl0_cmds.DISABLE,
        )
        # To reflect above config restart pipelined service is mandatory
        self._s1ap_wrapper.magmad_util.restart_services(["pipelined"])

if __name__ == "__main__":
    unittest.main()
