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


class TestSctpShutdowniWhileStatelessMmeIsStopped(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_sctp_shutdown_while_stateless_mme_is_stopped(self):
        """
        testing SCTP Shutdown while MME is stopped but Sctpd is running, i.e.:
        1. Attach 1 UE
        2. Stop MME service on AGW
        3. Send SCTP shutdown
        4. Print Redis state
        5. Start MME
        6. Test S1 setup

        """

        self._s1ap_wrapper.configUEDevice(1)

        req = self._s1ap_wrapper.ue_req
        print(
            "************************* Calling attach for UE id ",
            req.ue_id,
        )

        self._s1ap_wrapper.s1_util.attach(
            req.ue_id,
            s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )
        # Wait for EMM Information from MME
        self._s1ap_wrapper._s1_util.receive_emm_info()

        print("Stopping MME service")
        self._s1ap_wrapper.magmad_util.exec_command(
            "sudo service magma@mme stop",
        )

        print("send SCTP SHUTDOWN")
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.SCTP_SHUTDOWN_REQ, None,
        )

        print("Redis state after SCTP shutdown")
        self._s1ap_wrapper.magmad_util.print_redis_state()

        print("Starting MME service and waiting for 20 seconds")
        self._s1ap_wrapper.magmad_util.exec_command(
            "sudo service magma@mobilityd start",
        )
        self._s1ap_wrapper.magmad_util.exec_command(
            "sudo service magma@pipelined start",
        )
        self._s1ap_wrapper.magmad_util.exec_command(
            "sudo service magma@sessiond start",
        )
        self._s1ap_wrapper.magmad_util.exec_command(
            "sudo service magma@mme start",
        )
        time.sleep(30)

        print("Re-establish S1 connection between eNB and MME")
        self._s1ap_wrapper._s1setup()

        # Now detach the UE
        print(
            "************************* Calling detach for UE id ",
            req.ue_id,
        )

        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
