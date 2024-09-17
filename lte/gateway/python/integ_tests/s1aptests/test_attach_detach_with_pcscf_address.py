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
from integ_tests.s1aptests import s1ap_wrapper


class TestAttachDetachWithPcscfAddress(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_with_pcscf_address(self):
        """ Basic VoLTE attach/detach test with P-CSCF addresses """
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)
        req = self._s1ap_wrapper.ue_req
        pcscf_addr_type = ["ipv4", "ipv6", "ipv4v6"]
        # PDN Type 1-IPv4,2-IPv6,3-IPv4v6
        pdn_type = 1

        for idx in pcscf_addr_type:
            print(
                "************************* Running End to End attach for ",
                "UE id ",
                req.ue_id,
            )
            print("***************** P-CSCF address type", idx)
            # Now actually complete the attach
            self._s1ap_wrapper._s1_util.attach(
                req.ue_id,
                s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
                s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
                s1ap_types.ueAttachAccept_t,
                pdn_type=pdn_type,
                pcscf_addr_type=idx,
            )

            # Wait on EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value,
                True,
            )


if __name__ == "__main__":
    unittest.main()
