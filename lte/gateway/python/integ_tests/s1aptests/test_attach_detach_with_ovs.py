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
import s1ap_wrapper
from integ_tests.s1aptests.ovs.rest_api import get_datapath, get_flows
from magma.pipelined.imsi import decode_imsi
from s1ap_utils import GTPBridgeUtils


class TestAttachDetachWithOVS(unittest.TestCase):

    SPGW_TABLE = 0
    LOCAL_PORT = "LOCAL"

    def setUp(self):
        self._s1ap_wrapper = s1ap_wrapper.TestWrapper()

    def tearDown(self):
        self._s1ap_wrapper.cleanup()

    def check_imsi_metadata(self, flow, ue_req):
        """ Checks that IMSI set in the flow metadata matches the one sent """
        sent_imsi = 'IMSI' + ''.join([str(i) for i in ue_req.imsi])
        imsi_action = next(
            (
                a for a in flow["instructions"][0]["actions"]
                if a["field"] == "metadata"
            ), None,
        )
        self.assertIsNotNone(imsi_action)
        imsi64 = imsi_action["value"]
        # Convert between compacted uint IMSI and string
        received_imsi = decode_imsi(imsi64)
        self.assertEqual(
            sent_imsi, received_imsi,
            "IMSI set in metadata field does not match sent IMSI",
        )

    def test_attach_detach_with_ovs(self):
        """
        Basic sanity check of UE downlink/uplink flows during attach and
        detach procedures.
        """
        datapath = get_datapath()
        MAX_NUM_RETRIES = 5

        print("Checking for default table 0 flows")
        flows = get_flows(datapath, {"table_id": self.SPGW_TABLE})
        self.assertEqual(
            len(flows), 2,
            "There should only be 2 default table 0 flows",
        )

        self._s1ap_wrapper.configUEDevice(1)
        req = self._s1ap_wrapper.ue_req

        print("Running End to End attach for UE id ", req.ue_id)
        self._s1ap_wrapper._s1_util.attach(
            req.ue_id, s1ap_types.tfwCmd.UE_END_TO_END_ATTACH_REQUEST,
            s1ap_types.tfwCmd.UE_ATTACH_ACCEPT_IND,
            s1ap_types.ueAttachAccept_t,
        )

        self._s1ap_wrapper._s1_util.receive_emm_info()

        # UPLINK
        gtp_br_util = GTPBridgeUtils()
        gtp_port_no = gtp_br_util.get_gtp_port_no()
        print("Checking for uplink flow in-port %d", gtp_port_no)
        # try at least 5 times before failing as gateway
        # might take some time to install the flows in ovs
        for i in range(MAX_NUM_RETRIES):
            print("Get uplink flows: attempt ", i)
            uplink_flows = get_flows(
                datapath,
                {
                    "table_id": self.SPGW_TABLE,
                    "match": {"in_port": gtp_port_no},
                },
            )
            if len(uplink_flows) > 0:
                break
            time.sleep(5)  # sleep for 5 seconds before retrying

        self.assertEqual(len(uplink_flows), 1, "Uplink flow missing for UE")
        self.assertIsNotNone(
            uplink_flows[0]["match"]["tunnel_id"],
            "Uplink flow missing tunnel id match",
        )
        self.check_imsi_metadata(uplink_flows[0], req)

        # DOWNLINK
        print("Checking for downlink flow")
        ue_ip = str(self._s1ap_wrapper._s1_util.get_ip(req.ue_id))
        # Ryu can't match on ipv4_dst, so match on uplink in port
        # try at least 5 times before failing as gateway
        # might take some time to install the flows in ovs
        for i in range(MAX_NUM_RETRIES):
            print("Get downlink flows: attempt ", i)
            downlink_flows = get_flows(
                datapath,
                {
                    "table_id": self.SPGW_TABLE,
                    "match": {
                        "nw_dst": ue_ip,
                        "eth_type": 2048,
                        "in_port": self.LOCAL_PORT,
                    },
                },
            )
            if len(downlink_flows) > 0:
                break
            time.sleep(5)  # sleep for 5 seconds before retrying

        self.assertEqual(
            len(downlink_flows), 1,
            "Downlink flow missing for UE",
        )
        self.assertEqual(
            downlink_flows[0]["match"]["ipv4_dst"], ue_ip,
            "UE IP match missing from downlink flow",
        )

        actions = downlink_flows[0]["instructions"][0]["actions"]
        has_tunnel_action = any(
            action for action in actions
            if action["field"] == "tunnel_id"
            and action["type"] == "SET_FIELD"
        )
        self.assertTrue(
            has_tunnel_action,
            "Downlink flow missing set tunnel action",
        )
        self.check_imsi_metadata(downlink_flows[0], req)

        print("Running UE detach for UE id ", req.ue_id)
        # Now detach the UE
        self._s1ap_wrapper.s1_util.detach(
            req.ue_id, s1ap_types.ueDetachType_t.UE_NORMAL_DETACH.value, True,
        )

        print("Checking that uplink/downlink flows were deleted")
        flows = get_flows(datapath, {"table_id": self.SPGW_TABLE})
        self.assertEqual(
            len(flows), 2,
            "There should only be 2 default table 0 flows",
        )


if __name__ == "__main__":
    unittest.main()
