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
from integ_tests.s1aptests.s1ap_utils import SessionManagerUtil
from lte.protos.policydb_pb2 import FlowMatch


class TestAttachDetachRarActivationReject(unittest.TestCase):
    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()
        self._sessionManager_util = SessionManagerUtil()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def test_attach_detach_rar_activation_reject(self):
        """ Attach/detach + rar + dedicated bearer activation reject test
        with a single UE """
        num_ues = 1
        self._s1ap_wrapper.configUEDevice(num_ues)

        req = self._s1ap_wrapper.ue_req
        print(
            "********************** Running End to End attach for ",
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

        # UL Flow description #1
        ulFlow1 = {
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.UPLINK,  # Direction
        }
        # DL Flow description #1
        dlFlow1 = {
            "ip_proto": FlowMatch.IPPROTO_TCP,  # Protocol Type
            "direction": FlowMatch.DOWNLINK,  # Direction
        }

        # Flow list to be configured
        flow_list = [
            ulFlow1,
            dlFlow1,
        ]
        # QoS
        qos = {
            "qci": 5,  # qci value [1 to 9]
            "priority": 15,  # Range [0-255]
            "max_req_bw_ul": 10000000,  # MAX bw Uplink
            "max_req_bw_dl": 15000000,  # MAX bw Downlink
            "gbr_ul": 1000000,  # GBR Uplink
            "gbr_dl": 2000000,  # GBR Downlink
            "arp_prio": 15,  # ARP priority
            "pre_cap": 1,  # pre-emption capability
            "pre_vul": 1,  # pre-emption vulnerability
        }

        policy_id = "ims-voice"

        time.sleep(5)
        print(
            "********************** Sending RAR for IMSI",
            "".join([str(i) for i in req.imsi]),
        )
        self._sessionManager_util.send_ReAuthRequest(
            "IMSI" + "".join([str(i) for i in req.imsi]),
            policy_id,
            flow_list,
            qos,
        )

        response = self._s1ap_wrapper.s1_util.get_response()
        self.assertEqual(
            response.msg_type, s1ap_types.tfwCmd.UE_ACT_DED_BER_REQ.value,
        )
        act_ded_ber_ctxt_req = response.cast(
            s1ap_types.UeActDedBearCtxtReq_t,
        )
        ded_bearer_rej = s1ap_types.UeActDedBearCtxtRej_t()
        ded_bearer_rej.ue_Id = req.ue_id
        ded_bearer_rej.bearerId = act_ded_ber_ctxt_req.bearerId
        time.sleep(15)
        print(
            "********************** Sending activation Reject",
        )
        # Send Bearer Activation Reject
        self._s1ap_wrapper._s1_util.issue_cmd(
            s1ap_types.tfwCmd.UE_ACT_DED_BER_REJ, ded_bearer_rej,
        )

        time.sleep(15)
        print(
            "********************** Running UE detach for UE id ",
            req.ue_id,
        )
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
        )


if __name__ == "__main__":
    unittest.main()
