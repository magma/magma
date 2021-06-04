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


class TestAttachDetachMultiUeLooped(unittest.TestCase):

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_multi_ue_looped(self):
        """
        Multiple-attach/detach test with 32 UEs attached at a given time
        """
        num_ues = 16
        num_repeats = 10
        detach_type = [
            s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
            s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
        ]
        self._s1ap_wrapper.configUEDevice(num_ues)
        reqs = tuple(self._s1ap_wrapper.ue_req for _ in range(num_ues))
        for iteration in range(num_repeats):
            for req in reqs:
                print(
                    "************************* Running End to End attach ",
                    "for UE id ", req.ue_id,
                )
                # Now actually complete the attach
                self._s1ap_wrapper._s1_util.attach(
                    req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                    s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                    s1ap_types.ueAttachAccept_t,
                )

                # Wait on EMM Information from MME
                self._s1ap_wrapper._s1_util.receive_emm_info()

            for req in reqs:
                print(
                    "************************* Running UE detach for UE id ",
                    req.ue_id,
                )
                # Now detach the UE
                self._s1ap_wrapper.s1_util.detach(
                    req.ue_id, detach_type[iteration % 2], True,
                )


if __name__ == "__main__":
    unittest.main()
