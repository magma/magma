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
from integ_tests.s1aptests import s1ap_wrapper
from integ_tests.s1aptests.s1ap_utils import MagmadUtil, SpgwUtil


class TestAttachNwInitiatedDetachFail(unittest.TestCase):
    """
    S1AP Integration test for Failed Network Initiated Detach
    """

    def setUp(self):
        """Initialize s1ap wrapper and spgw utility
        """
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper(
            stateless_mode=MagmadUtil.stateless_cmds.ENABLE,
        )
        self._spgw_util = SpgwUtil()

    def tearDown(self):
        """Clean up utilities and sctp connection
        """
        self._s1ap_wrapper.cleanup()

    def test_attach_nw_initiated_detach_fail(self):
        """
        The test case validates retransmission of Detach Request after MME
        restarts
        Step 1: UE attaches to network
        Step 2: Send request to delete default bearer, since deletion is
                invoked for default bearer, MME initiates detach procedure
        Step 3: MME starts 3422 timer to receive Detach Accept message
        Step 4: S1AP tester ignores and does not send Detach Accept
        """
        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "********************** Running End to End attach for ",
            "UE id ",
            req.ue_id,
        )
        # Now actually complete the attach
        attach = self._s1ap_wrapper._s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        # Wait on EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print("Sleeping for 5 seconds")
        time.sleep(5)
        print(
            "********************** Deleting default bearer for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        # Delete default bearer
        self._spgw_util.delete_bearer(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            attach.esmInfo.epsBearerId,
            attach.esmInfo.epsBearerId,
        )
        # Receive NW initiated detach request
        # Wait for timer 3422 expiry 5 times
        for _ in range(5):
            response = self._s1ap_wrapper.s1_util.get_response()
            assert response.msg_type == s1ap_types.tfwCmd.UE_NW_INIT_DETACH_REQUEST.value
            print("**************** Received NW initiated Detach Req")

        time.sleep(6)


if __name__ == "__main__":
    unittest.main()
