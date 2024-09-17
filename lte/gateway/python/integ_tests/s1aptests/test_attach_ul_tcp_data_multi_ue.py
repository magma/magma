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


class TestAttachUlTcpDataMultiUe(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_ul_tcp_data_multi_ue(self):
        """ Attach and send UL TCP data with multiple UEs """
        reqs = []
        num_ues = 32
        self._s1ap_wrapper.configUEDevice(num_ues)

        for _ in range(num_ues):
            req = self._s1ap_wrapper.ue_req
            print(
                "************************* Running End to End attach for ",
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

            # Wait for EMM Information from MME
            self._s1ap_wrapper._s1_util.receive_emm_info()
            reqs.append(req)

        print(
            "************************* Running UE uplink (TCP) for UE ids ",
            [req.ue_id for req in reqs],
        )
        with self._s1ap_wrapper.configUplinkTest(*reqs, duration=1) as test:
            test.verify()

        for req in reqs:
            print(
                "************************* Running UE detach for UE id ",
                req.ue_id,
            )
            # Now detach the UE
            self._s1ap_wrapper.s1_util.detach(
                req.ue_id,
                s1ap_types.ueDetachType_t.UE_SWITCHOFF_DETACH.value,
                True,
            )


if __name__ == "__main__":
    unittest.main()
