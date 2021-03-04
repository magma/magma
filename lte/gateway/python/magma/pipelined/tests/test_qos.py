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
import asyncio
import subprocess
import unittest
from collections import namedtuple
from unittest.mock import MagicMock, call, patch
from magma.pipelined.bridge_util import BridgeTools

from magma.pipelined.qos.common import QosImplType, QosManager, SubscriberState
from magma.pipelined.qos.qos_meter_impl import MeterManager
from magma.pipelined.qos.qos_tc_impl import TrafficClass, argSplit, run_cmd
from magma.pipelined.qos.types import QosInfo, get_key_json, get_key, get_subscriber_key, \
        get_subscriber_data, get_data, get_data_json
from magma.pipelined.qos.utils import IdManager
from lte.protos.policydb_pb2 import FlowMatch
import logging

class TestQosCommon(unittest.TestCase):
    def testIdManager(self):
        id_mgr = IdManager(start_idx=1, max_idx=10)

        # allocate all available ids and verify
        idList = []
        for _ in range(9):
            idList.append(id_mgr.allocate_idx())

        assert idList == list(range(1, 10))

        # verify if we throw exception when we exceed max ids
        try:
            id_mgr.allocate_idx()
            self.assertTrue(False)
        except ValueError as e:
            self.assertTrue(str(e) == "maximum id allocation exceeded")

        # release odd ids and verify
        for i in range(1, 10, 2):
            id_mgr.release_idx(i)

        # release random id
        id_mgr.release_idx(8)

        new_list = []
        while True:
            try:
                new_list.append(id_mgr.allocate_idx())
            except ValueError:
                break

        self.assertTrue(new_list == list(range(1, 10, 2)) + [8])

        # verify if we recover state properly on reinitialization
        id_mgr = IdManager(start_idx=1, max_idx=10)
        old_id_set = {3, 6, 7, 9}
        id_mgr.restore_state(old_id_set)
        self.assertTrue(id_mgr._counter == 10)

        new_list = []
        while True:
            try:
                new_list.append(id_mgr.allocate_idx())
            except ValueError:
                break
        self.assertTrue(new_list == [1, 2, 4, 5, 8])

        # verify releasing indexes not present
        with self.assertLogs("pipelined.qos.id_manager", level="ERROR") as cm:
            id_mgr.release_idx(99)
        exp_err = "attempting to release invalid idx 99"
        self.assertTrue(cm.output[0].endswith(exp_err))

    def testQosKeyUtils(self):
        k = get_subscriber_key("1234", '1.1.1.1', 10, 0)
        j = get_key_json(k)
        self.assertTrue(get_key(j) == k)


class TestQosManager(unittest.TestCase):
    def setUp(self):
        self.dl_intf = "eth0"
        self.ul_intf = "eth1"

        self.config = {
            "redis_enabled": False,
            "clean_restart": True,
            "enodeb_iface": self.dl_intf,
            "nat_iface": self.ul_intf,
            "qos": {
                "max_rate": 1000000000,
                "enable": True,
                "ovs_meter": {"min_idx": 2, "max_idx": 100000},
                "linux_tc": {"min_idx": 2, "max_idx": 65534},
            },
        }

    def verifyTcAddQos(self, mock_get_action_inst, mock_traffic_cls, d, qid, qos_info,
                       parent_qid=0, skip_filter=False):
        intf = self.ul_intf if d == FlowMatch.UPLINK else self.dl_intf
        mock_get_action_inst.assert_any_call(qid)
        mock_traffic_cls.init_qdisc.assert_any_call(self.ul_intf, enable_pyroute2=False)
        mock_traffic_cls.init_qdisc.assert_any_call(self.dl_intf, enable_pyroute2=False)
        mock_traffic_cls.create_class.assert_any_call(intf, qid, qos_info.mbr,
            rate=qos_info.gbr, parent_qid=parent_qid, skip_filter=skip_filter)

    def verifyTcCleanRestart(self, prior_qids, mock_traffic_cls):
        for qid_tuple in prior_qids[self.ul_intf]:
            qid, _ = qid_tuple
            mock_traffic_cls.delete_class.assert_any_call(
                self.ul_intf, qid)

        for qid_tuple in prior_qids[self.dl_intf]:
            qid, _ = qid_tuple
            mock_traffic_cls.delete_class.assert_any_call(
                self.dl_intf, qid)

    def verifyTcRemoveQos(self, mock_traffic_cls, d, qid, skip_filter=False):
        intf = self.ul_intf if d == FlowMatch.UPLINK else self.dl_intf

        if skip_filter:
            mock_traffic_cls.delete_class.assert_any_call(intf, qid, True)
        else:
            mock_traffic_cls.delete_class.assert_any_call(intf, qid, False)

    def verifyTcRemoveQosBulk(self, mock_traffic_cls, argList):
        call_arg_list = []
        for d, qid in argList:
            intf = self.ul_intf if d == FlowMatch.UPLINK else self.dl_intf
            call_arg_list.append(call(intf, qid, False))
        mock_traffic_cls.delete_class.assert_has_calls(call_arg_list)

    def verifyMeterAddQos(
        self, mock_get_action_inst, mock_meter_cls, d, meter_id, qos_info
    ):
        mbr_kbps = int(qos_info.mbr / 1000)
        mock_get_action_inst.assert_any_call(meter_id)
        mock_meter_cls.add_meter.assert_any_call(
            MagicMock, meter_id, burst_size=0, rate=mbr_kbps
        )

    def verifyMeterRemoveQos(self, mock_meter_cls, d, meter_id):
        mock_meter_cls.del_meter.assert_any_call(MagicMock, meter_id)

    def verifyMeterRemoveQosBulk(self, mock_meter_cls, argList):
        call_arg_list = [call(MagicMock, meter_id) for _, meter_id in argList]
        mock_meter_cls.del_meter.assert_has_calls(call_arg_list)

    def verifyMeterCleanRestart(self, mock_meter_cls):
        mock_meter_cls.del_all_meters.assert_called_with(MagicMock)

    @patch("magma.pipelined.qos.qos_tc_impl.TrafficClass")
    @patch("magma.pipelined.qos.qos_meter_impl.MeterClass")
    @patch("magma.pipelined.qos.qos_tc_impl.TCManager.get_action_instruction")
    @patch(
        "magma.pipelined.qos.qos_meter_impl.MeterManager.\
get_action_instruction"
    )
    def _testSanity(
        self,
        mock_meter_get_action_inst,
        mock_tc_get_action_inst,
        mock_meter_cls,
        mock_traffic_cls,
    ):
        """ This test verifies that qos configuration gets programmed correctly
        for addition and deletion of a single subscriber and a single rule,
        We additionally verify if the clean_restart wipes out everything """

        # mock unclean state in qos
        prior_qids = {self.ul_intf: [(2, 0)], self.dl_intf: [(3, 0)]}
        if self.config["qos"]["impl"] == QosImplType.LINUX_TC:
            mock_traffic_cls.read_all_classes.side_effect = lambda intf: prior_qids[intf]

        qos_mgr = QosManager(MagicMock, asyncio.new_event_loop(), self.config)
        qos_mgr._redis_store = {}
        qos_mgr._setupInternal()
        imsi, ip_addr, rule_num, qos_info = "1234", '1.1.1.1', 0, QosInfo(100000, 100000)

        # add new subscriber qos queue
        qos_mgr.add_subscriber_qos(imsi, ip_addr, 0, rule_num, FlowMatch.UPLINK, qos_info)
        qos_mgr.add_subscriber_qos(imsi, ip_addr, 0, rule_num, FlowMatch.DOWNLINK, qos_info)

        k1 = get_key_json(get_subscriber_key(imsi, ip_addr, rule_num, FlowMatch.UPLINK))
        k2 = get_key_json(get_subscriber_key(imsi, ip_addr, rule_num, FlowMatch.DOWNLINK))
        ul_exp_id = qos_mgr.impl._start_idx
        dl_exp_id = qos_mgr.impl._start_idx + 1
        ul_qid_info = get_data_json(get_subscriber_data(ul_exp_id, 0, 0))
        dl_qid_info = get_data_json(get_subscriber_data(dl_exp_id, 0, 0))

        self.assertTrue(qos_mgr._redis_store[k1] == ul_qid_info)
        self.assertTrue(qos_mgr._redis_store[k2] == dl_qid_info)
        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][0] == (0, ul_qid_info))
        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][1] == (1, dl_qid_info))

        # add the same subscriber and ensure that we didn't create another
        # qos config for the subscriber
        qos_mgr.add_subscriber_qos(imsi, ip_addr, 0, rule_num, FlowMatch.UPLINK, qos_info)
        qos_mgr.add_subscriber_qos(imsi, ip_addr, 0, rule_num, FlowMatch.DOWNLINK, qos_info)

        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][0] == (0, ul_qid_info))
        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][1] == (1, dl_qid_info))

        # verify if traffic class was invoked properly
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterAddQos(
                mock_meter_get_action_inst,
                mock_meter_cls,
                FlowMatch.UPLINK,
                ul_exp_id,
                qos_info,
            )
            self.verifyMeterAddQos(
                mock_meter_get_action_inst,
                mock_meter_cls,
                FlowMatch.DOWNLINK,
                dl_exp_id,
                qos_info,
            )
            self.verifyMeterCleanRestart(mock_meter_cls)
        else:
            self.verifyTcAddQos(
                mock_tc_get_action_inst,
                mock_traffic_cls,
                FlowMatch.UPLINK,
                ul_exp_id,
                qos_info,
            )
            self.verifyTcAddQos(
                mock_tc_get_action_inst,
                mock_traffic_cls,
                FlowMatch.DOWNLINK,
                dl_exp_id,
                qos_info,
            )
            self.verifyTcCleanRestart(prior_qids, mock_traffic_cls)

        # remove the subscriber qos and verify things are cleaned up
        qos_mgr.remove_subscriber_qos(imsi, rule_num)
        self.assertTrue(len(qos_mgr._redis_store) == 0)
        self.assertTrue(imsi not in qos_mgr._subscriber_state)

        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQos(mock_meter_cls, FlowMatch.UPLINK, ul_exp_id)
            self.verifyMeterRemoveQos(mock_meter_cls, FlowMatch.DOWNLINK, dl_exp_id)
        else:
            self.verifyTcRemoveQos(mock_traffic_cls, FlowMatch.UPLINK, ul_exp_id)
            self.verifyTcRemoveQos(mock_traffic_cls, FlowMatch.DOWNLINK, dl_exp_id)

    @patch("magma.pipelined.qos.qos_tc_impl.TrafficClass")
    @patch("magma.pipelined.qos.qos_meter_impl.MeterClass")
    @patch("magma.pipelined.qos.qos_tc_impl.TCManager.get_action_instruction")
    @patch(
        "magma.pipelined.qos.qos_meter_impl.MeterManager.\
get_action_instruction"
    )
    def _testMultipleSubscribers(
        self,
        mock_meter_get_action_inst,
        mock_tc_get_action_inst,
        mock_meter_cls,
        mock_traffic_cls,
    ):
        """ This test verifies that qos configuration gets programmed correctly
        for addition and deletion of a multiple subscribers and rules.
        we additionally run through different scenarios involving
        - deactivating a rule
        - deactivating a subscriber
        - creating gaps in deletion and verifying that new qos configs get
        programmed properly with appropriate qids
        -  additionally we also verify the idempotency of deletion calls
        and ensure that code doesn't behave incorrectly when same items are
        deleted multiple times
        - Finally we delete everything and verify if that behavior is right"""

        if self.config["qos"]["impl"] == QosImplType.LINUX_TC:
            mock_traffic_cls.read_all_classes.side_effect = lambda intf: []
            mock_traffic_cls.delete_class.side_effect = lambda *args: 0

        qos_mgr = QosManager(MagicMock, asyncio.new_event_loop(), self.config)
        qos_mgr._redis_store = {}
        qos_mgr._setupInternal()
        rule_list1 = [
            ("1", 0, 0),
            ("1", 1, 0),
            ("1", 2, 1),
            ("2", 0, 0)
        ]

        rule_list2 = [
            ("2", 1, 0),
            ("3", 0, 0),
            ("4", 0, 0),
            ("5", 0, 0),
            ("5", 1, 1),
            ("6", 0, 0),
        ]

        start_idx, end_idx = 2, 2 + len(rule_list1)
        id_list = list(range(start_idx, end_idx))
        qos_info = QosInfo(100000, 100000)
        exp_id_dict = {}

        # add new subscriber qos queues
        for i, (imsi, rule_num, d) in enumerate(rule_list1):
            qos_mgr.add_subscriber_qos(imsi, '', 0, rule_num, d, qos_info)

            exp_id = id_list[i]
            k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))
            exp_id_dict[k] = exp_id
            # self.assertTrue(qos_mgr._redis_store[k] == exp_id)
            qid_info = get_data_json(get_subscriber_data(exp_id, 0, 0))
            self.assertEqual(qos_mgr._subscriber_state[imsi].rules[rule_num][0], (d, qid_info))

            if self.config["qos"]["impl"] == QosImplType.OVS_METER:
                self.verifyMeterAddQos(
                    mock_meter_get_action_inst, mock_meter_cls, d, exp_id, qos_info
                )
            else:
                self.verifyTcAddQos(
                    mock_tc_get_action_inst, mock_traffic_cls, d, exp_id, qos_info
                )

        # deactivate one rule
        # verify for imsi1 if rule num 0 gets cleaned up
        imsi, rule_num, d = rule_list1[0]
        k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))
        exp_id = exp_id_dict[k]

        qos_mgr.remove_subscriber_qos(imsi, rule_num)
        self.assertTrue(k not in qos_mgr._redis_store)
        self.assertTrue(not qos_mgr._subscriber_state[imsi].find_rule(rule_num))

        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQos(mock_meter_cls, d, exp_id)
        else:
            self.verifyTcRemoveQos(mock_traffic_cls, d, exp_id)

        # deactivate same rule and check if we log properly
        with self.assertLogs("pipelined.qos.common", level="DEBUG") as cm:
            qos_mgr.remove_subscriber_qos(imsi, rule_num)

        error_msg = "unable to find rule_num 0 for imsi 1"
        self.assertTrue(cm.output[1].endswith(error_msg))

        # deactivate imsi
        # verify for imsi1 if rule num 1 and 2 gets cleaned up
        qos_mgr.remove_subscriber_qos(imsi)
        remove_qos_args = []
        for imsi, rule_num, d in rule_list1[1:]:
            if imsi != "1":
                continue

            k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))
            exp_id = exp_id_dict[k]
            self.assertTrue(k not in qos_mgr._redis_store)
            remove_qos_args.append((d, exp_id))

        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQosBulk(mock_meter_cls, remove_qos_args)
        else:
            self.verifyTcRemoveQosBulk(mock_traffic_cls, remove_qos_args)

        self.assertTrue("1" not in qos_mgr._subscriber_state)

        # deactivate same imsi again and ensure nothing bad happens

        logging.debug("removing qos")
        with self.assertLogs("pipelined.qos.common", level="DEBUG") as cm:
            qos_mgr.remove_subscriber_qos("1")
        logging.debug("removing qos: done")

        error_msg = "imsi 1 not found"
        self.assertTrue(error_msg in cm.output[-1])

        # now only imsi2 should remain
        assert(len(qos_mgr._subscriber_state) == 1)
        assert(len(qos_mgr._subscriber_state['2'].rules) == 1)
        assert(len(qos_mgr._subscriber_state['2'].rules[0]) == 1)
        existing_qid = qos_mgr._subscriber_state['2'].rules[0][0][1]
        _, existing_qid, _, _ = get_data(existing_qid)

        # add second rule list and delete and verify if things work
        qos_info = QosInfo(100000, 200000)

        # additional 1 is for accomodating the existing imsi2RuleList1Qid
        start_idx, end_idx = 2, 2 + len(rule_list2) + 1
        id_list = [i for i in range(start_idx, end_idx) if i != existing_qid]

        # add new subscriber qos queues
        for i, (imsi, rule_num, d) in enumerate(rule_list2):
            qos_mgr.add_subscriber_qos(imsi, '', 0, rule_num, d, qos_info)
            exp_id = id_list[i]
            qid_info = get_data_json(get_subscriber_data(exp_id, 0, 0))

            k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))
            self.assertEqual(qos_mgr._redis_store[k], qid_info)
            self.assertEqual(qos_mgr._subscriber_state[imsi].rules[rule_num][0], (d, qid_info))
            if self.config["qos"]["impl"] == QosImplType.OVS_METER:
                self.verifyMeterAddQos(
                    mock_meter_get_action_inst, mock_meter_cls, d, exp_id, qos_info
                )
            else:
                self.verifyTcAddQos(
                    mock_tc_get_action_inst, mock_traffic_cls, d, exp_id, qos_info
                )

        # delete the subscriber qos queues
        for i, (imsi, rule_num, d) in enumerate(rule_list2):
            qos_mgr.remove_subscriber_qos(imsi, rule_num)

            if self.config["qos"]["impl"] == QosImplType.OVS_METER:
                self.verifyMeterRemoveQos(mock_meter_cls, d, id_list[i])
            else:
                self.verifyTcRemoveQos(mock_traffic_cls, d, id_list[i])

        self.assertTrue(len(qos_mgr._subscriber_state) == 1)
        self.assertTrue('2' in qos_mgr._subscriber_state)

        # delete everything
        qos_mgr.remove_subscriber_qos(imsi='2')

        # imsi2 from rule_list1 alone wasn't removed
        imsi, rule_num, d = rule_list1[3]
        k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))
        self.assertTrue(not qos_mgr._redis_store)
        self.assertTrue(not qos_mgr._subscriber_state)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQos(mock_meter_cls, d, existing_qid)
        else:
            self.verifyTcRemoveQos(mock_traffic_cls, d, existing_qid)

    @patch("magma.pipelined.qos.qos_tc_impl.TrafficClass")
    @patch("magma.pipelined.qos.qos_meter_impl.MeterClass")
    def _testUncleanRestart(self, mock_meter_cls, mock_traffic_cls):
        """This test verifies the case when we recover the state upon
        restart. We verify the base case of reconciling differences
        between system qos configs and qos_store configs(real code uses
        redis hash, we simply use dict). Additionally we test cases when
        system qos configs were wiped out and qos store state was wiped out
        and ensure that eventually the system and qos store state remains
        consistent"""
        loop = asyncio.new_event_loop()
        qos_mgr = QosManager(MagicMock, loop, self.config)
        qos_mgr._redis_store = {}

        def populate_db(qid_list, rule_list):
            qos_mgr._redis_store.clear()
            for i, t in enumerate(rule_list):
                k = get_key_json(get_subscriber_key(*t))
                v = get_data_json(get_subscriber_data(qid_list[i], 0, 0))
                qos_mgr._redis_store[k] = v

        MockSt = namedtuple("MockSt", "meter_id")
        dummy_meter_ev_body = [MockSt(11), MockSt(13), MockSt(2), MockSt(15)]

        def tc_read(intf):
            if intf == self.ul_intf:
                return [(2, 65534), (15, 65534)]
            else:
                return [(13, 65534), (11, 65534)]

        # prepopulate qos_store
        old_qid_list = [2, 11, 13, 20]
        old_rule_list = [("1",'1.1.1.1', 0, 0), ("1", '1.1.1.2', 1, 0),
                        ("1", '1.1.1.3', 2, 1), ("2", '1.1.1.4', 0, 0)]
        populate_db(old_qid_list, old_rule_list)

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            qos_mgr.impl.handle_meter_config_stats(dummy_meter_ev_body)
        else:
            mock_traffic_cls.read_all_classes.side_effect = tc_read

        qos_mgr._initialized = False
        qos_mgr._setupInternal()

        # verify that qos_handle 20 not found in system is purged from map
        qid_list = []
        for _, v in qos_mgr._redis_store.items():
            _, qid, _, _ = get_data(v)
            qid_list.append(qid)

        logging.debug("qid_list %s", qid_list)
        self.assertNotIn(20, qid_list)

        # verify that unreferenced qos configs are purged from the system
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_called_with(MagicMock, 15)
        else:
            mock_traffic_cls.delete_class.assert_called_with(self.ul_intf, 15, False)

        # add a new rule to the qos_mgr and check if it is assigned right id
        imsi, rule_num, d, qos_info = "3", 0, 0, QosInfo(100000, 100000)
        qos_mgr.impl.get_action_instruction = MagicMock
        qos_mgr.add_subscriber_qos(imsi, '', 0, rule_num, d, qos_info)
        k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))

        exp_id = 3  # since start_idx 2 is already used
        d3 = get_data_json(get_subscriber_data(exp_id, 0, 0))

        self.assertEqual(qos_mgr._redis_store[k], d3)
        self.assertEqual(qos_mgr._subscriber_state[imsi].rules[rule_num][0], (d, d3))

        # delete the restored rule - ensure that it gets cleaned properly
        purge_imsi = "1"
        purge_rule_num = 0
        purge_qos_handle = 2
        qos_mgr.remove_subscriber_qos(purge_imsi, purge_rule_num)
        self.assertTrue(purge_rule_num not in qos_mgr._subscriber_state[purge_imsi].rules)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_called_with(MagicMock, purge_qos_handle)
        else:
            mock_traffic_cls.delete_class.assert_called_with(
                self.ul_intf, purge_qos_handle, False
            )

        # case 2 - check with empty qos configs, qos_map gets purged
        mock_meter_cls.reset_mock()
        mock_traffic_cls.reset_mock()
        populate_db(old_qid_list, old_rule_list)

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            MockSt = namedtuple("MockSt", "meter_id")
            qos_mgr.impl._fut = loop.create_future()
            qos_mgr.impl.handle_meter_config_stats([])
        else:
            mock_traffic_cls.read_all_classes.side_effect = lambda _: []

        qos_mgr._initialized = False
        qos_mgr._setupInternal()

        self.assertTrue(not qos_mgr._redis_store)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_not_called()
        else:
            mock_traffic_cls.delete_class.assert_not_called()

        # case 3 - check with empty qos_map, all qos configs get purged
        mock_meter_cls.reset_mock()
        mock_traffic_cls.reset_mock()
        qos_mgr._redis_store.clear()

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            qos_mgr.impl._fut = loop.create_future()
            qos_mgr.impl.handle_meter_config_stats(dummy_meter_ev_body)
        else:
            mock_traffic_cls.read_all_classes.side_effect = tc_read

        qos_mgr._initialized = False
        qos_mgr._setupInternal()

        self.assertTrue(not qos_mgr._redis_store)
        # verify that unreferenced qos configs are purged from the system
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 2)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 15)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 13)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 11)
        else:
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 2, False)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 15, False)
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 13, False)
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 11, False)

    @patch("magma.pipelined.qos.qos_tc_impl.TrafficClass")
    @patch("magma.pipelined.qos.qos_meter_impl.MeterClass")
    def _testUncleanRestartWithApnAMBR(self, mock_meter_cls, mock_traffic_cls):
        """This test verifies all tests cases from _testUncleanRestart for APN
        AMBR configs.
        """
        loop = asyncio.new_event_loop()
        qos_mgr = QosManager(MagicMock, loop, self.config)
        qos_mgr._redis_store = {}

        def populate_db(qid_list, old_ambr_list, old_leaf_list, rule_list):
            qos_mgr._redis_store.clear()
            for i, t in enumerate(rule_list):
                k = get_key_json(get_subscriber_key(*t))
                v = get_data_json(get_subscriber_data(qid_list[i],
                                                      old_ambr_list[i],
                                                      old_leaf_list[i]))
                qos_mgr._redis_store[k] = v

        MockSt = namedtuple("MockSt", "meter_id")
        dummy_meter_ev_body = [MockSt(11), MockSt(13), MockSt(2), MockSt(15)]

        def tc_read(intf):
            if intf == self.ul_intf:
                return [(30, 3000), (15, 1500),
                        (300, 3000), (150, 1500),
                        (3000, 65534), (1500, 65534)]
            else:
                return [(13, 1300), (11, 1100),
                        (130, 1300), (110, 1100),
                        (1300, 65534), (1100, 65534)]

        # prepopulate qos_store
        old_qid_list = [2, 11, 13, 30]
        old_leaf_list = [20, 110, 130, 300]
        old_ambr_list = [200, 1100, 1300, 3000]

        old_rule_list = [("1",'1.1.1.1', 0, 0), ("1", '1.1.1.2', 1, 0),
                        ("1", '1.1.1.3', 2, 1), ("2", '1.1.1.4', 0, 0)]
        populate_db(old_qid_list, old_ambr_list, old_leaf_list, old_rule_list)

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            qos_mgr.impl.handle_meter_config_stats(dummy_meter_ev_body)
        else:
            mock_traffic_cls.read_all_classes.side_effect = tc_read

        qos_mgr._initialized = False
        qos_mgr._setupInternal()

        # verify that qos_handle 20 not found in system is purged from map
        qid_list = []
        for _, v in qos_mgr._redis_store.items():
            _, qid, _, _ = get_data(v)
            qid_list.append(qid)

        logging.debug("qid_list %s", qid_list)
        self.assertNotIn(20, qid_list)

        # verify that unreferenced qos configs are purged from the system
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_called_with(MagicMock, 15)
        else:
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 15, False)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 150, False)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 1500, True)

        # add a new rule to the qos_mgr and check if it is assigned right id
        imsi, rule_num, d, qos_info = "3", 0, 0, QosInfo(100000, 100000)
        qos_mgr.impl.get_action_instruction = MagicMock
        qos_mgr.add_subscriber_qos(imsi, '', 10, rule_num, d, qos_info)
        k = get_key_json(get_subscriber_key(imsi, '', rule_num, d))

        exp_id = 3  # since start_idx 2 is already used
        d3 = get_data_json(get_subscriber_data(exp_id+1, exp_id-1, exp_id))

        self.assertEqual(qos_mgr._redis_store[k], d3)
        self.assertEqual(qos_mgr._subscriber_state[imsi].rules[rule_num][0], (d, d3))

        # delete the restored rule - ensure that it gets cleaned properly
        purge_imsi = "1"
        purge_rule_num = 1
        purge_qos_handle = 2
        qos_mgr.remove_subscriber_qos(purge_imsi, purge_rule_num)
        self.assertTrue(purge_rule_num not in qos_mgr._subscriber_state[purge_imsi].rules)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_called_with(MagicMock, purge_qos_handle)
        else:
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 11, False)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 110, False)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 1100, True)

        # case 2 - check with empty qos configs, qos_map gets purged
        mock_meter_cls.reset_mock()
        mock_traffic_cls.reset_mock()
        populate_db(old_qid_list, old_ambr_list, old_leaf_list, old_rule_list)

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            MockSt = namedtuple("MockSt", "meter_id")
            qos_mgr.impl._fut = loop.create_future()
            qos_mgr.impl.handle_meter_config_stats([])
        else:
            mock_traffic_cls.read_all_classes.side_effect = lambda _: []

        qos_mgr._initialized = False
        qos_mgr._setupInternal()

        self.assertTrue(not qos_mgr._redis_store)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_not_called()
        else:
            mock_traffic_cls.delete_class.assert_not_called()

        # case 3 - check with empty qos_map, all qos configs get purged
        mock_meter_cls.reset_mock()
        mock_traffic_cls.reset_mock()
        qos_mgr._redis_store.clear()

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            qos_mgr.impl._fut = loop.create_future()
            qos_mgr.impl.handle_meter_config_stats(dummy_meter_ev_body)
        else:
            mock_traffic_cls.read_all_classes.side_effect = tc_read

        logging.debug("case three")
        qos_mgr._initialized = False
        qos_mgr._setupInternal()

        self.assertTrue(not qos_mgr._redis_store)
        # verify that unreferenced qos configs are purged from the system
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 2)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 15)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 13)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 11)
        else:
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 15, False)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 150, False)

            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 13, False)
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 130, False)

            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 11, False)
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 110, False)

    def testSanity(self):
        for impl_type in (QosImplType.LINUX_TC, QosImplType.OVS_METER):
            self.config["qos"]["impl"] = impl_type
            self._testSanity()

    def testMultipleSubscribers(self):
        for impl_type in (QosImplType.LINUX_TC, QosImplType.OVS_METER):
            self.config["qos"]["impl"] = impl_type
            self._testMultipleSubscribers()

    def testUncleanRestart(self):
        with patch.dict(self.config, {"clean_restart": False}):
            for impl_type in [QosImplType.LINUX_TC]:
                self.config["qos"]["impl"] = impl_type
                self._testUncleanRestart()

    def testUncleanRestartAPN(self):
        with patch.dict(self.config, {"clean_restart": False}):
            for impl_type in [QosImplType.LINUX_TC]:
                self.config["qos"]["impl"] = impl_type
                self._testUncleanRestartWithApnAMBR()

    @patch("magma.pipelined.qos.qos_tc_impl.TrafficClass")
    @patch("magma.pipelined.qos.qos_tc_impl.TCManager.get_action_instruction")
    def testApnAmbrSanity(
        self,
        mock_tc_get_action_inst,
        mock_traffic_cls,
    ):
        """ This test verifies that qos configuration gets programmed correctly
        for addition and deletion of a single subscriber and a single rule,
        We additionally verify if the clean_restart wipes out everything """
        self.config["qos"]["impl"] = QosImplType.LINUX_TC
        # mock unclean state in qos
        prior_qids = {self.ul_intf: [(2, 0)], self.dl_intf: [(3, 0)]}
        mock_traffic_cls.read_all_classes.side_effect = lambda intf: prior_qids[intf]

        qos_mgr = QosManager(MagicMock, asyncio.new_event_loop(), self.config)
        qos_mgr._redis_store = {}
        qos_mgr._setupInternal()
        ambr_ul, ambr_dl = 250000, 500000
        imsi, ip_addr, rule_num, qos_info = ("1234", '1.1.1.1', 1, QosInfo(50000, 100000))

        # add new subscriber qos queue
        qos_mgr.add_subscriber_qos(imsi, ip_addr, ambr_ul, rule_num, FlowMatch.UPLINK, qos_info)
        qos_mgr.add_subscriber_qos(imsi, ip_addr, ambr_dl, rule_num, FlowMatch.DOWNLINK, qos_info)

        k1 = get_key_json(get_subscriber_key(imsi, ip_addr, rule_num, FlowMatch.UPLINK))
        k2 = get_key_json(get_subscriber_key(imsi, ip_addr, rule_num, FlowMatch.DOWNLINK))

        ambr_ul_exp_id = qos_mgr.impl._start_idx
        ul_exp_id = qos_mgr.impl._start_idx + 2
        ambr_dl_exp_id = qos_mgr.impl._start_idx + 3
        dl_exp_id = qos_mgr.impl._start_idx + 5

        qid_info_ul = get_subscriber_data(ul_exp_id, ul_exp_id - 2, ul_exp_id - 1)
        qid_info_dl = get_subscriber_data(dl_exp_id, dl_exp_id - 2, dl_exp_id - 1)
        self.assertEqual(get_data(qos_mgr._redis_store[k1]), qid_info_ul)
        self.assertEqual(get_data(qos_mgr._redis_store[k2]), qid_info_dl)
        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][0] == (0, get_data_json(qid_info_ul)))
        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][1] == (1, get_data_json(qid_info_dl)))

        self.assertEqual(len(qos_mgr._subscriber_state[imsi].sessions), 1)
        self.assertEqual(qos_mgr._subscriber_state[imsi].sessions[ip_addr].ambr_dl, ambr_dl_exp_id)
        self.assertEqual(qos_mgr._subscriber_state[imsi].sessions[ip_addr].ambr_dl_leaf, ambr_dl_exp_id + 1)

        self.assertEqual(qos_mgr._subscriber_state[imsi].sessions[ip_addr].ambr_ul_leaf, ambr_ul_exp_id + 1)

        # add the same subscriber and ensure that we didn't create another
        # qos config for the subscriber
        qos_mgr.add_subscriber_qos(imsi, ip_addr, 0, rule_num, FlowMatch.UPLINK, qos_info)
        qos_mgr.add_subscriber_qos(imsi, ip_addr, 0, rule_num, FlowMatch.DOWNLINK, qos_info)

        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][0] == (0, get_data_json(qid_info_ul)))
        self.assertTrue(qos_mgr._subscriber_state[imsi].rules[rule_num][1] == (1, get_data_json(qid_info_dl)))

        # verify if traffic class was invoked properly
        self.verifyTcAddQos(mock_tc_get_action_inst, mock_traffic_cls, FlowMatch.UPLINK,
                            ul_exp_id, qos_info, parent_qid=ambr_ul_exp_id)
        self.verifyTcAddQos(mock_tc_get_action_inst,mock_traffic_cls, FlowMatch.DOWNLINK,
                            dl_exp_id, qos_info, parent_qid=ambr_dl_exp_id)
        self.verifyTcCleanRestart(prior_qids, mock_traffic_cls)

        # remove the subscriber qos and verify things are cleaned up
        qos_mgr.remove_subscriber_qos(imsi, rule_num)
        self.assertTrue(len(qos_mgr._redis_store) == 0)
        self.assertTrue(imsi not in qos_mgr._subscriber_state)

        self.verifyTcRemoveQos(mock_traffic_cls, FlowMatch.UPLINK, ambr_ul_exp_id, True)
        self.verifyTcRemoveQos(mock_traffic_cls, FlowMatch.UPLINK, ul_exp_id)
        self.verifyTcRemoveQos(mock_traffic_cls, FlowMatch.DOWNLINK, ambr_dl_exp_id, True)
        self.verifyTcRemoveQos(mock_traffic_cls, FlowMatch.DOWNLINK, dl_exp_id)


class TestMeters(unittest.TestCase):
    @patch("magma.pipelined.qos.qos_meter_impl.MeterClass")
    def testBrokenMeter(self, meter_cls):
        config = {
            "redis_enabled": False,
            "qos": {
                "enable": True,
                "max_rate": 10000000,
                "ovs_meter": {"min_idx": 2, "max_idx": 100000},
            }
        }
        m = MeterManager(None, asyncio.new_event_loop(), config)
        self.assertTrue(meter_cls.dump_meter_features.called)
        MockSt = namedtuple("MockSt", "max_meter")
        m.handle_meter_feature_stats([MockSt(0)])
        try:
            m.add_qos(0, QosInfo(100000, 100000))
            self.fail("unexpectedly add_qos succeeded")
        except RuntimeError:
            pass


class TestTrafficClass(unittest.TestCase):

    @patch("subprocess.check_call")
    @patch("os.geteuid", return_value=1)
    def testSudoUser(self, _, mock_check_call):
        intf = 'qt'
        BRIDGE = 'qtbr0'
        BridgeTools.create_bridge(BRIDGE, BRIDGE)
        BridgeTools.create_internal_iface(BRIDGE, intf, None)

        TrafficClass.init_qdisc(intf)
        mock_check_call.assert_any_call(
            ["sudo", "tc", "qdisc", "add", "dev", intf, "root", "handle", "1:", "htb"]
        )

    @patch("subprocess.check_call")
    def testError(self, mock_check_call):
        def dummy_check_call(*args):
            raise subprocess.CalledProcessError(returncode=1, cmd="tc")

        mock_check_call.side_effect = dummy_check_call
        with self.assertLogs("pipelined.qos.tc_cmd", level="ERROR") as cm:
            TrafficClass.init_qdisc("eth0", show_error=True)
        self.assertTrue("error: 1 running: tc qdisc add dev " in cm.output[0])

    def testSanityTrafficClass(self, ):
        intf = 'qt'
        BRIDGE = 'qtbr0'
        BridgeTools.create_bridge(BRIDGE, BRIDGE)
        BridgeTools.create_internal_iface(BRIDGE, intf, None)

        parent_qid = 2
        qid = 3
        apn_ambr = 1000000
        bearer_mbr = 500000
        bearer_gbr = 250000
        TrafficClass.init_qdisc(intf, show_error=False)

        # create APN level ambr
        TrafficClass.create_class(intf, qid=parent_qid, max_bw=apn_ambr)

        # create child queue
        TrafficClass.create_class(intf, qid=qid, rate=bearer_gbr, max_bw=bearer_mbr,
                                  parent_qid=parent_qid)

        # check if the filters installed for leaf class only
        filter_output = subprocess.check_output(['tc', 'filter', 'show', 'dev', intf])
        filter_list = filter_output.decode('utf-8').split("\n")
        filter_list = [ln for ln in filter_list if 'classid' in ln]
        assert('classid 1:{qid}'.format(qid=parent_qid) in filter_list[0])
        assert('classid 1:{qid}'.format(qid=qid) in filter_list[1])

        # check if classes are installed with appropriate bandwidth limits
        class_output = subprocess.check_output(['tc', 'class', 'show', 'dev', intf])
        class_list = class_output.decode('utf-8').split("\n")
        for info in class_list:
            if 'class htb 1:{qid}'.format(qid=qid) in info:
                child_class = info

            if 'class htb 1:{qid}'.format(qid=parent_qid) in info:
                parent_class = info

        assert(parent_class and 'ceil 1Mbit' in parent_class)
        assert(child_class and 'rate 250Kbit ceil 500Kbit' in child_class)

        # check if fq_codel is associated only with the leaf class
        qdisc_output = subprocess.check_output(['tc', 'qdisc', 'show', 'dev', intf])

        # check if read_all_classes work
        qid_list = TrafficClass.read_all_classes(intf)
        assert((qid, parent_qid) in qid_list)

        # delete leaf class
        TrafficClass.delete_class(intf, 3)

        # check class for qid 3 removed
        class_output = subprocess.check_output(['tc', 'class', 'show', 'dev', intf])
        class_list = class_output.decode('utf-8').split("\n")
        assert( not [info for info in class_list  if 'class htb 1:{qid}'.format(
            qid=qid) in info])

        # delete APN AMBR class
        TrafficClass.delete_class(intf, 2)

        # verify that parent class is removed
        class_output = subprocess.check_output(['tc', 'class', 'show', 'dev', intf])
        class_list = class_output.decode('utf-8').split("\n")
        assert( not [info for info in class_list  if 'class htb 1:{qid}'.format(
            qid=parent_qid) in info])

        # check if no fq_codel nor filter exists
        qdisc_output = subprocess.check_output(['tc', 'qdisc', 'show', 'dev', intf])
        filter_output = subprocess.check_output(['tc', 'filter', 'show', 'dev', intf])
        filter_list = filter_output.decode('utf-8').split("\n")
        filter_list = [ln for ln in filter_list if 'classid' in ln]
        qdisc_list = qdisc_output.decode('utf-8').split("\n")
        qdisc_list = [ln for ln in qdisc_list if 'fq_codel' in ln]
        assert(not filter_list and not qdisc_list)

        # destroy all qos on intf
        run_cmd(['tc qdisc del dev {intf} root'.format(intf=intf)])


class TestSubscriberState(unittest.TestCase):
    def testSingleRuleWithNoApnAmbr(self, ):
        rule_num, d = 10, FlowMatch.UPLINK
        ip_addr = '1.1.1.1'
        qos_handle = 10

        # add rule
        subscriber_state = SubscriberState('IMSI101', {})
        assert(subscriber_state.check_empty())
        assert(not subscriber_state.get_qos_handle(rule_num, d))
        subscriber_state.update_rule(ip_addr, rule_num, d, qos_handle, 0, 0)

        assert(subscriber_state.find_rule(rule_num))
        session_with_rule = subscriber_state.find_session_with_rule(rule_num)
        assert(session_with_rule)

        # remove rule
        subscriber_state.remove_rule(rule_num)
        empty_sessions = subscriber_state.get_all_empty_sessions()
        assert(len(empty_sessions) == 1)
        assert(session_with_rule == empty_sessions[0])

        subscriber_state.remove_session(ip_addr)
        assert(subscriber_state.check_empty())

    def testSingleRuleWithApnAmbr(self,):
        rule_num, d = 10, FlowMatch.UPLINK
        ip_addr = '1.1.1.1'
        ambr_qos_handle = 5
        qos_handle = 10
        subscriber_state = SubscriberState('IMSI101', {})
        assert(subscriber_state.check_empty())
        assert(not subscriber_state.get_qos_handle(rule_num, d))

        session = subscriber_state.get_or_create_session(ip_addr)
        session.set_ambr(d, ambr_qos_handle, 0)
        subscriber_state.update_rule(ip_addr, rule_num, d, qos_handle, 0, 0)

        assert(subscriber_state.find_rule(rule_num))
        session_with_rule = subscriber_state.find_session_with_rule(rule_num)
        assert(session_with_rule)

        # remove rule
        subscriber_state.remove_rule(rule_num)
        empty_sessions = subscriber_state.get_all_empty_sessions()
        assert(len(empty_sessions) == 1)
        assert(session_with_rule == empty_sessions[0])

        subscriber_state.remove_session(ip_addr)
        assert(subscriber_state.check_empty())

    def testMultipleRuleWithApnAmbr(self,):
        rule_num1, d1 = 10, FlowMatch.UPLINK
        rule_num2, d2 = 20, FlowMatch.UPLINK
        rule_num3, d3 = 20, FlowMatch.DOWNLINK
        ip_addr = '1.1.1.1'
        ambr_qos_handle = 5
        qos_handle1 = 10
        qos_handle2 = 20
        qos_handle3 = 30

        subscriber_state = SubscriberState('IMSI101', {})
        assert(subscriber_state.check_empty())
        assert(not subscriber_state.get_qos_handle(rule_num1, d1))

        subscriber_state.get_or_create_session(ip_addr)
        subscriber_state.update_rule(ip_addr, rule_num1, d1, qos_handle1, 0, 0)

        # add rule_num2
        subscriber_state.update_rule(ip_addr, rule_num2, d2, qos_handle2, 0, 0)

        # add rule_num3 with apn_ambr in downlink direction
        session = subscriber_state.get_or_create_session(ip_addr)
        assert(session)
        session.set_ambr(d3, ambr_qos_handle, 0)
        subscriber_state.update_rule(ip_addr, rule_num3, d3, qos_handle3, 0, 0)

        assert(len(subscriber_state.rules) == 2)
        assert(rule_num1 in subscriber_state.rules)
        assert(rule_num2 in subscriber_state.rules)
        assert(len(subscriber_state.sessions) == 1)

        subscriber_state.remove_rule(rule_num1)
        assert(not subscriber_state.get_all_empty_sessions())

        subscriber_state.remove_rule(rule_num2)
        session = subscriber_state.get_all_empty_sessions()
        assert(len(session) == 1)
        subscriber_state.remove_session(session[0].ip_addr)
        assert(subscriber_state.check_empty())

    def testMultipleSessionsWithMultipleRulesAmbr(self, ):
        # session1 information
        rule_num1, d1 = 10, FlowMatch.UPLINK
        rule_num2, d2 = 20, FlowMatch.UPLINK
        rule_num3, d3 = 20, FlowMatch.DOWNLINK
        ip_addr = '1.1.1.1'
        ambr_qos_handle = 5
        qos_handle1 = 10
        qos_handle2 = 20
        qos_handle3 = 30

        # session2 information
        new_session_ip_addr = '2.2.2.2'
        new_rule_num4, new_d1 = 110, FlowMatch.UPLINK
        new_ambr_qos_handle = 105
        new_qos_handle1 = 110

        subscriber_state = SubscriberState('IMSI101', {})
        assert(subscriber_state.check_empty())
        assert(not subscriber_state.get_qos_handle(rule_num1, d1))

        subscriber_state.get_or_create_session(ip_addr)
        subscriber_state.update_rule(ip_addr, rule_num1, d1, qos_handle1, 0, 0)

        # add rule_num2
        subscriber_state.update_rule(ip_addr, rule_num2, d2, qos_handle2, 0, 0)

        # add rule_num3 with apn_ambr in downlink direction
        session = subscriber_state.get_or_create_session(ip_addr)
        assert(session)
        session.set_ambr(d3, ambr_qos_handle, 0)
        subscriber_state.update_rule(ip_addr, rule_num3, d3, qos_handle3, 0, 0)

        # add rule_num4 with apn_ambr in uplink direction
        session = subscriber_state.get_or_create_session(new_session_ip_addr)
        assert(session)
        session.set_ambr(new_d1, new_ambr_qos_handle, 0)
        subscriber_state.update_rule(new_session_ip_addr, new_rule_num4, new_d1,
            new_qos_handle1, 0, 0)

        assert(len(subscriber_state.rules) == 3)
        assert(rule_num1 in subscriber_state.rules)
        assert(rule_num2 in subscriber_state.rules)
        assert(new_rule_num4 in subscriber_state.rules)
        assert(len(subscriber_state.sessions) == 2)

        # remove the rules
        subscriber_state.remove_rule(rule_num1)
        assert(not subscriber_state.get_all_empty_sessions())

        subscriber_state.remove_rule(rule_num2)
        session = subscriber_state.get_all_empty_sessions()
        assert(len(session) == 1)
        subscriber_state.remove_session(session[0].ip_addr)
        assert(not subscriber_state.check_empty())

        # remove the new session
        subscriber_state.remove_rule(new_rule_num4)
        session = subscriber_state.get_all_empty_sessions()
        assert(len(session) == 1)
        subscriber_state.remove_session(session[0].ip_addr)
        assert(subscriber_state.check_empty())
