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
import warnings
from concurrent.futures import Future

from ryu.lib import hub

from magma.pipelined.tests.app.start_pipelined import (
    TestSetup,
    PipelinedController,
)
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.tests.pipelined_test_util import (
    start_ryu_app_thread,
    stop_ryu_app_thread,
    create_service_manager,
    assert_bridge_snapshot_match,
)
from ryu.ofproto.ofproto_v1_4 import OFPP_LOCAL


class LIMirrorTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    LI_LOCAL_IFACE = 'local_li0'
    LI_DST_IFACE = 'dst_li0'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    LI_LOCAL_IP = '1.1.1.1'
    LI_DST_IP = '2.2.3.3'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        super(LIMirrorTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['li_mirror'])

        inout_controller_reference = Future()
        li_mirror_reference = Future()
        testing_controller_reference = Future()
        test_setup = TestSetup(
            apps=[PipelinedController.InOut,
                  PipelinedController.LIMirror,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.InOut:
                    inout_controller_reference,
                PipelinedController.LIMirror:
                    li_mirror_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'CWF',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'internal_ip_subnet': '192.168.0.0/16',
                'ovs_gtp_port_number': 32768,
                'clean_restart': True,
                'li_mirror_all': True,
                'li_local_iface': cls.LI_LOCAL_IFACE,
                'li_dst_iface': cls.LI_DST_IFACE,
                'uplink_port': OFPP_LOCAL
            },
            mconfig=None,
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.LI_LOCAL_IFACE,
                                          cls.LI_LOCAL_IP)
        BridgeTools.create_internal_iface(cls.BRIDGE, cls.LI_DST_IFACE,
                                          cls.LI_DST_IP)

        cls.thread = start_ryu_app_thread(test_setup)
        cls.inout_controller = inout_controller_reference.result()
        cls.li_controller = li_mirror_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def testFlowSnapshotMatch(self):
        hub.sleep(3)
        assert_bridge_snapshot_match(self, self.BRIDGE, self.service_manager)


if __name__ == "__main__":
    unittest.main()
