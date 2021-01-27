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


class TestSctpShutdownAfterMultiUeAttach(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()
        self._s1ap_wrapper.magmad_util.exec_command(
            "sudo service sctpd restart"
        )
        print(
            "Restart sctpd service to clear Redis state as test case doesn't"
            " intend to initiate detach procedure"
        )
        magtivate_cmd = "source /home/vagrant/build/python/bin/activate"
        state_cli_cmd = "state_cli.py keys IMSI*"
        redis_state = self._s1ap_wrapper.magmad_util.exec_command_output(
            magtivate_cmd + " && " + state_cli_cmd
        )
        print("Redis state is [\n", redis_state, "]")

    def test_sctp_shutdown_after_multi_ue_attach(self):
        """ Attah multiple UEs and send sctp shutdown without detach """
        num_ues = 32
        self._s1ap_wrapper.configUEDevice(num_ues)
        for _ in range(num_ues):
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


if __name__ == "__main__":
    unittest.main()
