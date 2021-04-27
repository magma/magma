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
from typing import List

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.mobilityd_pb2 import IPAddress
from lte.protos.pipelined_pb2 import VersionedPolicy
from lte.protos.policydb_pb2 import (
    FlowDescription,
    FlowMatch,
    HeaderEnrichment,
    PolicyRule,
)
from magma.pipelined.app import he
from magma.pipelined.app.enforcement import EnforcementController
from magma.pipelined.app.he import HeaderEnrichmentController
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow.messages import MessageHub, MsgChannel
from magma.pipelined.openflow.registers import Direction
from magma.pipelined.policy_converters import (
    convert_ip_str_to_ip_proto,
    convert_ipv4_str_to_ip_proto,
)
from magma.pipelined.tests.app.start_pipelined import (
    PipelinedController,
    TestSetup,
)
from magma.pipelined.tests.app.subscriber import RyuDirectSubscriberContext
from magma.pipelined.tests.app.table_isolation import (
    RyuDirectTableIsolator,
    RyuForwardFlowArgsBuilder,
)
from magma.pipelined.tests.pipelined_test_util import (
    SnapshotVerifier,
    create_service_manager,
    fake_controller_setup,
    start_ryu_app_thread,
    stop_ryu_app_thread,
    wait_after_send,
)


def mocked_activate_he_urls_for_ue(ip: IPAddress, rule_id: str, urls: List[str], imsi: str, msisdn: str):
    return True


def mocked_deactivate_he_urls_for_ue(ip: IPAddress, rule_id: str):
    return True


class HeTableTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    BRIDGE_IP = '192.168.128.1'
    UE_BLOCK = '192.168.128.0/24'
    UE_MAC = '5e:cc:cc:b1:49:4b'
    UE_IP = '192.168.128.22'
    OTHER_MAC = '0a:00:27:00:00:02'
    OTHER_IP = '1.2.3.4'
    VETH = 'tveth'
    VETH_NS = 'tveth_ns'
    PROXY_PORT = '15'

    @classmethod
    @unittest.mock.patch('netifaces.ifaddresses',
                         return_value=[[{'addr': '00:11:22:33:44:55'}]])
    @unittest.mock.patch('netifaces.AF_LINK', 0)
    def setUpClass(cls, *_):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures.
        """
        he.activate_he_urls_for_ue = mocked_activate_he_urls_for_ue
        he.deactivate_he_urls_for_ue = mocked_deactivate_he_urls_for_ue

        super(HeTableTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([], ['proxy'])
        cls._tbl_num = cls.service_manager.get_table_num(HeaderEnrichmentController.APP_NAME)

        BridgeTools.create_veth_pair(cls.VETH, cls.VETH_NS)
        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)
        BridgeTools.add_ovs_port(cls.BRIDGE, cls.VETH, cls.PROXY_PORT)

        he_controller_reference = Future()
        testing_controller_reference = Future()

        test_setup = TestSetup(
            apps=[
                PipelinedController.HeaderEnrichment,
                PipelinedController.Testing,
                PipelinedController.StartupFlows
            ],
            references={
                PipelinedController.HeaderEnrichment:
                    he_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'setup_type': 'LTE',
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': cls.BRIDGE_IP,
                'uplink_port': 20,
                'proxy_port_name': cls.VETH,
                'clean_restart': True,
                'enable_nat': True,
                'ovs_gtp_port_number': 10,
            },
            mconfig=PipelineD(
                ue_ip_block=cls.UE_BLOCK,
            ),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False,
        )

        cls.thread = start_ryu_app_thread(test_setup)
        cls.he_controller = he_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    def _wait_for_responses(self, chan, response_count, logger):
        def fail(err):
            logger.error("Failed to install rule for subscriber: %s", err)

        for _ in range(response_count):
            try:
                result = chan.get()

            except MsgChannel.Timeout:
                return fail("No response from OVS policy mixin")
            if not result.ok():
                return fail(result.exception())

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def tearDown(self):
        cls = self.__class__
        dp = HeTableTest.he_controller._datapath
        cls.he_controller.delete_all_flows(dp)

    def test_default_flows(self):
        """
        Verify that a proxy flows are setup
        """

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass

    def test_ue_flows_add(self):
        """
        Verify that a proxy flows are setup
        """
        cls = self.__class__
        self._msg_hub = MessageHub(HeTableTest.he_controller.logger)

        ue_ip = '1.1.1.1'
        tun_id = 1
        dest_server = '2.2.2.2'
        flow_msg = cls.he_controller.get_subscriber_he_flows("rule1", Direction.OUT, ue_ip, tun_id, dest_server, 123,
                                                             ['abc.com'], 'IMSI01', b'1')
        chan = self._msg_hub.send(flow_msg,
                                  HeTableTest.he_controller._datapath, )
        self._wait_for_responses(chan, len(flow_msg), HeTableTest.he_controller.logger)

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass

    def test_ue_flows_add_direction_in(self):
        """
        Verify that a proxy flows are setup
        """
        cls = self.__class__
        self._msg_hub = MessageHub(HeTableTest.he_controller.logger)

        ue_ip = '1.1.1.1'
        tun_id = 1
        dest_server = '2.2.2.2'
        flow_msg = cls.he_controller.get_subscriber_he_flows("rule1", Direction.IN, ue_ip, tun_id, dest_server, 123,
                                                             ['abc.com'], 'IMSI01', b'1')
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip), 0)
        chan = self._msg_hub.send(flow_msg,
                                  HeTableTest.he_controller._datapath, )
        self._wait_for_responses(chan, len(flow_msg), HeTableTest.he_controller.logger)

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass

    def test_ue_flows_add2(self):
        """
        Verify that a proxy flows are setup
        """
        cls = self.__class__
        self._msg_hub = MessageHub(HeTableTest.he_controller.logger)
        dp = HeTableTest.he_controller._datapath
        ue_ip1 = '1.1.1.200'
        tun_id1 = 1
        dest_server1 = '2.2.2.4'
        rule1 = 123
        flow_msg = cls.he_controller.get_subscriber_he_flows("rule1", Direction.OUT, ue_ip1, tun_id1, dest_server1, rule1,
                                                             ['abc.com'], 'IMSI01', b'1')

        ue_ip2 = '10.10.10.20'
        tun_id2 = 2
        dest_server2 = '20.20.20.40'
        rule2 = 1230
        flow_msg.extend(cls.he_controller.get_subscriber_he_flows("rule2", Direction.OUT, ue_ip2, tun_id2, dest_server2, rule2,
                                                                  ['abc.com'], 'IMSI01', b'1'))
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 1)
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip2), 1)

        chan = self._msg_hub.send(flow_msg, dp)
        self._wait_for_responses(chan, len(flow_msg), HeTableTest.he_controller.logger)

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass

    def test_ue_flows_del(self):
        """
        Verify that a proxy flows are setup
        """
        cls = self.__class__
        self._msg_hub = MessageHub(HeTableTest.he_controller.logger)
        dp = HeTableTest.he_controller._datapath
        ue_ip1 = '1.1.1.200'
        tun_id1 = 1

        dest_server1 = '2.2.2.4'
        rule1 = 123
        flow_msg = cls.he_controller.get_subscriber_he_flows('rule1', Direction.OUT, ue_ip1, tun_id1, dest_server1, rule1,
                                                             ['abc.com'], 'IMSI01', b'1')

        ue_ip2 = '10.10.10.20'
        tun_id2 = 2
        dest_server2 = '20.20.20.40'
        rule2 = 1230
        flow_msg2 = cls.he_controller.get_subscriber_he_flows('rule2', Direction.OUT, ue_ip2, tun_id2, dest_server2, rule2,
                                                              ['abc.com'], 'IMSI01', b'1')
        flow_msg.extend(flow_msg2)
        chan = self._msg_hub.send(flow_msg, dp)
        self._wait_for_responses(chan, len(flow_msg), HeTableTest.he_controller.logger)

        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip2), 'rule2', rule2)

        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip2), 'rule_random', 3223)

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass

    def test_ue_flows_del2(self):
        """
        Verify that a proxy flows are setup
        """
        cls = self.__class__
        self._msg_hub = MessageHub(HeTableTest.he_controller.logger)
        dp = HeTableTest.he_controller._datapath
        ue_ip1 = '1.1.1.200'
        tun_id1 = 1
        dest_server1 = '2.2.2.4'
        rule1 = 123
        flow_msg = cls.he_controller.get_subscriber_he_flows('rule1', Direction.OUT, ue_ip1, tun_id1, dest_server1, rule1,
                                                             ['abc.com'], 'IMSI01', b'1')

        ue_ip2 = '10.10.10.20'
        tun_id2 = 2
        dest_server2 = '20.20.20.40'
        rule2 = 1230
        flow_msg.extend(cls.he_controller.get_subscriber_he_flows('rule2', Direction.OUT, ue_ip2, tun_id2, dest_server2, rule2,
                                                                  ['abc.com'], 'IMSI01', b'1'))
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 1)
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip2), 1)

        ue_ip2 = '10.10.10.20'
        dest_server2 = '20.20.40.40'
        rule2 = 1230
        flow_msg.extend(cls.he_controller.get_subscriber_he_flows('rule2', Direction.OUT, ue_ip2, tun_id2, dest_server2, rule2,
                                                                  ['abc.com'], 'IMSI01', None))

        chan = self._msg_hub.send(flow_msg, dp)
        self._wait_for_responses(chan, len(flow_msg), HeTableTest.he_controller.logger)

        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip2))

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass
        # verify multiple remove works.
        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip2))
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip2), 0)

    def test_ue_flows_multi_rule(self):
        """
        Verify that a proxy flows are setup
        """
        cls = self.__class__
        self._msg_hub = MessageHub(HeTableTest.he_controller.logger)
        dp = HeTableTest.he_controller._datapath
        ue_ip1 = '1.1.1.200'
        tun_id1 = 1
        dest_server1 = '2.2.2.4'
        rule1 = 123
        flow_msg = cls.he_controller.get_subscriber_he_flows('rule1', Direction.OUT, ue_ip1, tun_id1, dest_server1,
                                                             rule1, ['abc.com'], 'IMSI01', b'1')

        tun_id2 = 2
        dest_server2 = '20.20.20.40'
        rule2 = 1230
        flow_msg.extend(cls.he_controller.get_subscriber_he_flows('rule2', Direction.OUT, ue_ip1, tun_id2, dest_server2,
                                                                  rule2, ['abc1.com'], 'IMSI01', b'1'))
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 2)

        dest_server2 = '20.20.40.40'
        rule3 = 1230
        flow_msg.extend(cls.he_controller.get_subscriber_he_flows('rule3', Direction.OUT, ue_ip1, tun_id2, dest_server2,
                                                                  rule3, ['abc2.com'], 'IMSI01', None))

        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 3)

        dest_server2 = '20.20.50.50'
        rule4 = 22
        flow_msg.extend(cls.he_controller.get_subscriber_he_flows('rule4', Direction.OUT, ue_ip1, tun_id2, dest_server2,
                                                                  rule4, ['abc2.com'], 'IMSI01', None))

        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 4)

        chan = self._msg_hub.send(flow_msg, dp)
        self._wait_for_responses(chan, len(flow_msg), HeTableTest.he_controller.logger)

        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip1), "rule1", rule1)

        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager,
                                             max_sleep_time=20,
                                             datapath=HeTableTest.he_controller._datapath)

        with snapshot_verifier:
            pass
        # verify multiple remove works.
        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip1), "rule2", rule2)
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 2)
        cls.he_controller.remove_subscriber_he_flows(convert_ip_str_to_ip_proto(ue_ip1))
        self.assertEqual(cls.he_controller._ue_rule_counter.get(ue_ip1), 0)


class EnforcementTableHeTest(unittest.TestCase):
    BRIDGE = 'testing_br'
    IFACE = 'testing_br'
    MAC_DEST = "5e:cc:cc:b1:49:4b"
    he_controller_reference = Future()
    VETH = 'tveth1'
    VETH_NS = 'tveth1_ns'
    PROXY_PORT = '16'

    @classmethod
    def setUpClass(cls):
        """
        Starts the thread which launches ryu apps

        Create a testing bridge, add a port, setup the port interfaces. Then
        launch the ryu apps for testing pipelined. Gets the references
        to apps launched by using futures, mocks the redis policy_dictionary
        of enforcement_controller
        """
        super(EnforcementTableHeTest, cls).setUpClass()
        warnings.simplefilter('ignore')
        cls.service_manager = create_service_manager([PipelineD.ENFORCEMENT], ['proxy'])
        cls._tbl_num = cls.service_manager.get_table_num(
            EnforcementController.APP_NAME)
        BridgeTools.create_bridge(cls.BRIDGE, cls.IFACE)

        BridgeTools.create_veth_pair(cls.VETH, cls.VETH_NS)
        BridgeTools.add_ovs_port(cls.BRIDGE, cls.VETH, cls.PROXY_PORT)

        enforcement_controller_reference = Future()
        testing_controller_reference = Future()
        he.activate_he_urls_for_ue = mocked_activate_he_urls_for_ue
        he.deactivate_he_urls_for_ue = mocked_deactivate_he_urls_for_ue

        test_setup = TestSetup(
            apps=[PipelinedController.Enforcement,
                  PipelinedController.HeaderEnrichment,
                  PipelinedController.Testing,
                  PipelinedController.StartupFlows],
            references={
                PipelinedController.Enforcement:
                    enforcement_controller_reference,
                PipelinedController.HeaderEnrichment:
                    cls.he_controller_reference,
                PipelinedController.Testing:
                    testing_controller_reference,
                PipelinedController.StartupFlows:
                    Future(),
            },
            config={
                'bridge_name': cls.BRIDGE,
                'bridge_ip_address': '192.168.128.1',
                'nat_iface': 'eth2',
                'enodeb_iface': 'eth1',
                'qos': {'enable': False},
                'clean_restart': True,
                'uplink_port': 20,
                'proxy_port_name': cls.VETH,
                'enable_nat': True,
                'ovs_gtp_port_number': 10,
                'setup_type': 'LTE',
            },
            mconfig=PipelineD(),
            loop=None,
            service_manager=cls.service_manager,
            integ_test=False
        )

        cls.thread = start_ryu_app_thread(test_setup)

        cls.enforcement_controller = enforcement_controller_reference.result()
        cls.testing_controller = testing_controller_reference.result()

    @classmethod
    def tearDownClass(cls):
        stop_ryu_app_thread(cls.thread)
        BridgeTools.destroy_bridge(cls.BRIDGE)

    def test_subscriber_policy_with_he(self):
        """
        Add policy to subscriber with HE config

        """
        cls = self.__class__

        fake_controller_setup(self.enforcement_controller)

        imsi = 'IMSI010000000088888'
        sub_ip = '192.168.128.74'
        uplink_tunnel = 0x1234
        flow_list1 = [FlowDescription(
            match=FlowMatch(
                ip_dst=convert_ipv4_str_to_ip_proto('45.10.0.0/24'),
                direction=FlowMatch.UPLINK),
            action=FlowDescription.PERMIT)
        ]
        he = HeaderEnrichment(urls=['abc.com'])
        policies = [
            VersionedPolicy(
                rule=PolicyRule(id='simple_match', priority=2, flow_list=flow_list1, he=he),
                version=1,
            )
        ]

        # ============================ Subscriber ============================
        sub_context = RyuDirectSubscriberContext(
            imsi, sub_ip, uplink_tunnel, self.enforcement_controller,
            self._tbl_num
        ).add_policy(policies[0])

        isolator = RyuDirectTableIsolator(
            RyuForwardFlowArgsBuilder.from_subscriber(sub_context.cfg)
                .build_requests(),
            self.testing_controller
        )
        snapshot_verifier = SnapshotVerifier(self,
                                             self.BRIDGE,
                                             self.service_manager)
        with isolator, sub_context, snapshot_verifier:
            pass


if __name__ == "__main__":
    unittest.main()
