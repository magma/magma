"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory  this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import asyncio
import subprocess
import unittest
from lte.protos.policydb_pb2 import FlowMatch
from collections import namedtuple
from unittest.mock import patch, MagicMock, call
from magma.pipelined.qos.types import (QosInfo, get_subscriber_key, get_json,
                                       get_key)
from magma.pipelined.qos.utils import IdManager
from magma.pipelined.qos.common import QosManager, QosImplType
from magma.pipelined.qos.qos_meter_impl import MeterManager
from magma.pipelined.qos.qos_tc_impl import TrafficClass


class TestQosCommon(unittest.TestCase):
    def testIdManager(self,):
        id_mgr = IdManager(start_idx=1, max_idx=10)

        # allocate all available ids and verify
        idList = []
        for _ in range(9):
            idList.append(id_mgr.allocate_idx())

        assert(idList == list(range(1, 10)))

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

        self.assertTrue(new_list == list(range(1, 10, 2)) + [8, ])

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
        with self.assertLogs('pipelined.qos.id_manager', level='ERROR') as cm:
            id_mgr.release_idx(99)
        exp_err = "attempting to release invalid idx 99"
        self.assertTrue(cm.output[0].endswith(exp_err))

    def testQosKeyUtils(self,):
        k = get_subscriber_key("imsi1234", 10, 0)
        j = get_json(k)
        self.assertTrue(get_key(j) == k)


class TestQosManager(unittest.TestCase):
    def setUp(self,):
        self.dl_intf = "eth0"
        self.ul_intf = "eth1"

        self.config = {
            "clean_restart": True,
            "enodeb_iface": self.dl_intf,
            "nat_iface": self.ul_intf,
            'qos': {
                "max_rate": 1000000000,
                'enable': True,
                'ovs_meter': {
                    'min_idx': 2,
                    'max_idx': 100000,
                },
                'linux_tc': {
                    'min_idx': 2,
                    'max_idx': 65535,
                },
            },
        }

    def verifyTcAddQos(self, mock_get_action_inst, mock_traffic_cls, d, qid,
                       qos_info):
        intf = self.ul_intf if d == FlowMatch.UPLINK else self.dl_intf
        mock_get_action_inst.assert_called_with(qid,)
        mock_traffic_cls.init_qdisc.assert_any_call(self.ul_intf)
        mock_traffic_cls.init_qdisc.assert_any_call(self.dl_intf)
        mock_traffic_cls.create_class.assert_called_with(intf, qid,
                                                         qos_info.mbr)

    def verifyTcCleanRestart(self, prior_qids, mock_traffic_cls):
        for qid in prior_qids[self.ul_intf]:
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, qid,
                                                          show_error=False,
                                                          throw_except=False)

        for qid in prior_qids[self.dl_intf]:
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, qid,
                                                          show_error=False,
                                                          throw_except=False)

    def verifyTcRemoveQos(self, mock_traffic_cls, d, qid):
        intf = self.ul_intf if d == FlowMatch.UPLINK else self.dl_intf
        mock_traffic_cls.delete_class.assert_called_with(intf, qid)

    def verifyTcRemoveQosBulk(self, mock_traffic_cls, argList):
        call_arg_list = []
        for d, qid in argList:
            intf = self.ul_intf if d == FlowMatch.UPLINK else self.dl_intf
            call_arg_list.append(call(intf, qid))
        mock_traffic_cls.delete_class.assert_has_calls(call_arg_list)

    def verifyMeterAddQos(self, mock_get_action_inst, mock_meter_cls, d,
                          meter_id, qos_info):
        mbr_kbps = int(qos_info.mbr / 1000)
        mock_get_action_inst.assert_called_with(meter_id,)
        mock_meter_cls.add_meter.assert_called_with(MagicMock, meter_id,
                                                    burst_size=0,
                                                    rate=mbr_kbps)

    def verifyMeterRemoveQos(self, mock_meter_cls, d, meter_id):
        mock_meter_cls.del_meter.assert_called_with(MagicMock, meter_id)

    def verifyMeterRemoveQosBulk(self, mock_meter_cls, argList):
        call_arg_list = [call(MagicMock, meter_id) for _, meter_id in argList]
        mock_meter_cls.del_meter.assert_has_calls(call_arg_list)

    def verifyMeterCleanRestart(self, mock_meter_cls):
        mock_meter_cls.del_all_meters.assert_called_with(MagicMock)

    @patch('magma.pipelined.qos.qos_tc_impl.TrafficClass')
    @patch('magma.pipelined.qos.qos_meter_impl.MeterClass')
    @patch('magma.pipelined.qos.qos_tc_impl.TCManager.get_action_instruction')
    @patch('magma.pipelined.qos.qos_meter_impl.MeterManager.\
get_action_instruction')
    def _testSanity(self,
                    mock_meter_get_action_inst,
                    mock_tc_get_action_inst,
                    mock_meter_cls,
                    mock_traffic_cls):
        ''' This test verifies that qos configuration gets programmed correctly
        for addition and deletion of a single subscriber and a single rule,
        We additionally verify if the clean_restart wipes out everything '''

        # mock unclean state in qos
        prior_qids = {
            self.ul_intf: [2, ],
            self.dl_intf: [3, ]
        }
        if self.config["qos"]["impl"] == QosImplType.LINUX_TC:
            mock_traffic_cls.read_all_classes.side_effect = (lambda intf:
                                                             prior_qids[intf])

        qos_mgr = QosManager(MagicMock, asyncio.new_event_loop(), self.config)
        qos_mgr._qos_store = {}
        qos_mgr._setupInternal()
        imsi, rule_num, d, qos_info = "imsi1234", 0, 0, QosInfo(100000, 100000)

        # add new subscriber qos queue
        qos_mgr.add_subscriber_qos(imsi, rule_num, d, qos_info)
        k = get_json(get_subscriber_key(imsi, rule_num, d))

        exp_id = qos_mgr.qos_impl._start_idx
        self.assertTrue(qos_mgr._qos_store[k] == exp_id)
        self.assertTrue(qos_mgr._subscriber_map[imsi][rule_num] == (exp_id, d))

        # add the same subscriber and ensure that we didn't create another
        # qos config for the subscriber
        qos_mgr.add_subscriber_qos(imsi, rule_num, d, qos_info)
        self.assertTrue(qos_mgr._subscriber_map[imsi][rule_num] == (exp_id, d))

        # verify if traffic class was invoked properly
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterAddQos(mock_meter_get_action_inst, mock_meter_cls,
                                   d, exp_id, qos_info)
            self.verifyMeterCleanRestart(mock_meter_cls)
        else:
            self.verifyTcAddQos(mock_tc_get_action_inst, mock_traffic_cls,
                                d, exp_id, qos_info)
            self.verifyTcCleanRestart(prior_qids, mock_traffic_cls)

        # remove the subscriber qos and verify things are cleaned up
        qos_mgr.remove_subscriber_qos(imsi, rule_num)
        self.assertTrue(len(qos_mgr._qos_store) == 0)
        self.assertTrue(imsi not in qos_mgr._subscriber_map)

        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQos(mock_meter_cls, d, exp_id)
        else:
            self.verifyTcRemoveQos(mock_traffic_cls, d, exp_id)

    @patch('magma.pipelined.qos.qos_tc_impl.TrafficClass')
    @patch('magma.pipelined.qos.qos_meter_impl.MeterClass')
    @patch('magma.pipelined.qos.qos_tc_impl.TCManager.get_action_instruction')
    @patch('magma.pipelined.qos.qos_meter_impl.MeterManager.\
get_action_instruction')
    def _testMultipleSubscribers(self,
                                 mock_meter_get_action_inst,
                                 mock_tc_get_action_inst,
                                 mock_meter_cls,
                                 mock_traffic_cls):
        ''' This test verifies that qos configuration gets programmed correctly
        for addition and deletion of a multiple subscribers and rules.
        we additionally run through different scenarios involving
        - deactivating a rule
        - deactivating a subscriber
        - creating gaps in deletion and verifying that new qos configs get
        programmed properly with appropriate qids
        -  additionally we also verify the idempotency of deletion calls
        and ensure that code doesn't behave incorrectly when same items are
        deleted multiple times
        - Finally we delete everything and verify if that behavior is right'''
        qos_mgr = QosManager(MagicMock, asyncio.new_event_loop(), self.config)
        qos_mgr._qos_store = {}
        qos_mgr._setupInternal()
        rule_list1 = [("imsi1", 0, 0),
                      ("imsi1", 1, 0),
                      ("imsi1", 2, 1),
                      ("imsi2", 0, 0)]

        rule_list2 = [("imsi2", 1, 0),
                      ("imsi3", 0, 0),
                      ("imsi4", 0, 0),
                      ("imsi5", 0, 0),
                      ("imsi5", 1, 1),
                      ("imsi6", 0, 0)]

        start_idx, end_idx = 2, 2 + len(rule_list1)
        id_list = list(range(start_idx, end_idx))
        qos_info = QosInfo(100000, 100000)
        exp_id_dict = {}
        # add new subscriber qos queues
        for i, (imsi, rule_num, d) in enumerate(rule_list1):
            qos_mgr.add_subscriber_qos(imsi, rule_num, d, qos_info)

            exp_id = id_list[i]
            k = get_json(get_subscriber_key(imsi, rule_num, d))
            exp_id_dict[k] = exp_id

            self.assertTrue(qos_mgr._qos_store[k] == exp_id)
            self.assertTrue(qos_mgr._subscriber_map[imsi][rule_num] == (exp_id,
                            d))

            if self.config["qos"]["impl"] == QosImplType.OVS_METER:
                self.verifyMeterAddQos(mock_meter_get_action_inst,
                                       mock_meter_cls, d, exp_id, qos_info)
            else:
                self.verifyTcAddQos(mock_tc_get_action_inst, mock_traffic_cls,
                                    d, exp_id, qos_info)

        # deactivate one rule
        # verify for imsi1 if rule num 0 gets cleaned up
        imsi, rule_num, d = rule_list1[0]
        k = get_json(get_subscriber_key(imsi, rule_num, d))
        exp_id = exp_id_dict[k]
        qos_mgr.remove_subscriber_qos(imsi, rule_num)
        self.assertTrue(k not in qos_mgr._qos_store)
        self.assertTrue(rule_num not in qos_mgr._subscriber_map[imsi])

        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQos(mock_meter_cls, d, exp_id)
        else:
            self.verifyTcRemoveQos(mock_traffic_cls, d, exp_id)

        # deactivate same rule and check if we log properly
        with self.assertLogs('pipelined.qos.common', level='ERROR') as cm:
            qos_mgr.remove_subscriber_qos(imsi, rule_num)
        error_msg = "unable to find rule_num 0 for imsi imsi1"
        self.assertTrue(cm.output[0].endswith(error_msg))

        # deactivate imsi
        # verify for imsi1 if rule num 1 and 2 gets cleaned up
        qos_mgr.remove_subscriber_qos(imsi)
        remove_qos_args = []
        for imsi, rule_num, d in rule_list1[1:]:
            if imsi != "imsi1":
                continue

            k = get_json(get_subscriber_key(imsi, rule_num, d))
            exp_id = exp_id_dict[k]
            self.assertTrue(k not in qos_mgr._qos_store)
            remove_qos_args.append((d, exp_id))

        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQosBulk(mock_meter_cls, remove_qos_args)
        else:
            self.verifyTcRemoveQosBulk(mock_traffic_cls, remove_qos_args)

        self.assertTrue("imsi1" not in qos_mgr._subscriber_map)

        # deactivate same imsi again and ensure nothing bad happens
        with self.assertLogs('pipelined.qos.common', level='ERROR') as cm:
            qos_mgr.remove_subscriber_qos("imsi1")
        error_msg = "unable to find imsi imsi1"
        self.assertTrue(cm.output[0].endswith(error_msg))

        # now only imsi2 should remain
        self.assertTrue(len(qos_mgr._qos_store) == 1)
        existing_qid = list(qos_mgr._qos_store.values())[0]

        # add second rule list and delete and verify if things work
        qos_info = QosInfo(100000, 200000)

        # additional 1 is for accomodating the existing imsi2RuleList1Qid
        start_idx, end_idx = 2, 2 + len(rule_list2) + 1
        id_list = [i for i in range(start_idx, end_idx) if i != existing_qid]

        # add new subscriber qos queues
        for i, (imsi, rule_num, d) in enumerate(rule_list2):
            qos_mgr.add_subscriber_qos(imsi, rule_num, d, qos_info)
            exp_id = id_list[i]
            k = get_json(get_subscriber_key(imsi, rule_num, d))
            self.assertTrue(qos_mgr._qos_store[k] == exp_id)
            self.assertTrue(qos_mgr._subscriber_map[imsi][rule_num] == (exp_id,
                            d))
            if self.config["qos"]["impl"] == QosImplType.OVS_METER:
                self.verifyMeterAddQos(mock_meter_get_action_inst,
                                       mock_meter_cls, d, exp_id, qos_info)
            else:
                self.verifyTcAddQos(mock_tc_get_action_inst, mock_traffic_cls,
                                    d, exp_id, qos_info)

        # delete the subscriber qos queues
        for i, (imsi, rule_num, d) in enumerate(rule_list2):
            qos_mgr.remove_subscriber_qos(imsi, rule_num)
            k = get_json(get_subscriber_key(imsi, rule_num, d))
            self.assertTrue(k not in qos_mgr._qos_store)
            self.assertTrue(rule_num not in qos_mgr._subscriber_map[imsi])
            if self.config["qos"]["impl"] == QosImplType.OVS_METER:
                self.verifyMeterRemoveQos(mock_meter_cls, d, id_list[i])
            else:
                self.verifyTcRemoveQos(mock_traffic_cls, d, id_list[i])

        # delete everything
        qos_mgr.remove_subscriber_qos()

        # imsi2 from rule_list1 alone wasn't removed
        imsi, rule_num, d = rule_list1[3]
        k = get_json(get_subscriber_key(imsi, rule_num, d))
        self.assertTrue(not qos_mgr._qos_store)
        self.assertTrue(not qos_mgr._subscriber_map)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            self.verifyMeterRemoveQos(mock_meter_cls, d, existing_qid)
        else:
            self.verifyTcRemoveQos(mock_traffic_cls, d, existing_qid)

    @patch('magma.pipelined.qos.qos_tc_impl.TrafficClass')
    @patch('magma.pipelined.qos.qos_meter_impl.MeterClass')
    def _testUncleanRestart(self,
                            mock_meter_cls,
                            mock_traffic_cls):
        '''This test verifies the case when we recover the state upon
        restart. We verify the base case of reconciling differences
        between system qos configs and qos_store configs(real code uses
        redis hash, we simply use dict). Additionally we test cases when
        system qos configs were wiped out and qos store state was wiped out
        and ensure that eventually the system and qos store state remains
        consistent'''
        loop = asyncio.new_event_loop()
        qos_mgr = QosManager(MagicMock, loop, self.config)
        qos_mgr._qos_store = {}

        def populate_db(qid_list, rule_list):
            qos_mgr._qos_store.clear()
            for i, t in enumerate(rule_list):
                k = get_json(get_subscriber_key(*t))
                qos_mgr._qos_store[k] = qid_list[i]

        MockSt = namedtuple("MockSt", "meter_id")
        dummy_meter_ev_body = [MockSt(11), MockSt(13), MockSt(2), MockSt(15)]

        def tc_read(intf):
            if intf == self.ul_intf:
                return [2, 15]
            else:
                return [13, 11]

        # prepopulate qos_store
        old_qid_list = [2, 11, 13, 20]
        old_rule_list = [("imsi1", 0, 0),
                         ("imsi1", 1, 0),
                         ("imsi1", 2, 1),
                         ("imsi2", 0, 0)]
        populate_db(old_qid_list, old_rule_list)

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            qos_mgr.qos_impl.handle_meter_config_stats(dummy_meter_ev_body)
        else:
            mock_traffic_cls.read_all_classes.side_effect = tc_read

        qos_mgr._setupInternal()

        # run async loop once to ensure ready items are cleared
        loop._run_once()

        # verify that qos_handle 20 not found in system is purged from map
        self.assertFalse([v for _, v in qos_mgr._qos_store.items() if v == 20])

        # verify that unreferenced qos configs are purged from the system
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_called_with(MagicMock, 15)
        else:
            mock_traffic_cls.delete_class.assert_called_with(self.ul_intf, 15)

        # add a new rule to the qos_mgr and check if it is assigned right id
        imsi, rule_num, d, qos_info = "imsi3", 0, 0, QosInfo(100000, 100000)
        qos_mgr.qos_impl.get_action_instruction = MagicMock
        qos_mgr.add_subscriber_qos(imsi, rule_num, d, qos_info)
        k = get_json(get_subscriber_key(imsi, rule_num, d))

        exp_id = 3  # since start_idx 2 is already used
        self.assertTrue(qos_mgr._qos_store[k] == exp_id)
        self.assertTrue(qos_mgr._subscriber_map[imsi][rule_num] == (exp_id, d))

        # case 2 - check with empty qos configs, qos_map gets purged
        mock_meter_cls.reset_mock()
        mock_traffic_cls.reset_mock()
        populate_db(old_qid_list, old_rule_list)

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            MockSt = namedtuple("MockSt", "meter_id")
            qos_mgr.qos_impl._fut = loop.create_future()
            qos_mgr.qos_impl.handle_meter_config_stats([])
        else:
            mock_traffic_cls.read_all_classes.side_effect = lambda _: []

        qos_mgr._setupInternal()

        # run async loop once to ensure ready items are cleared
        loop._run_once()
        self.assertTrue(not qos_mgr._qos_store)
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_not_called()
        else:
            mock_traffic_cls.delete_class.assert_not_called()

        # case 3 - check with empty qos_map, all qos configs get purged
        mock_meter_cls.reset_mock()
        mock_traffic_cls.reset_mock()
        qos_mgr._qos_store.clear()

        # mock future state
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            qos_mgr.qos_impl._fut = loop.create_future()
            qos_mgr.qos_impl.handle_meter_config_stats(dummy_meter_ev_body)
        else:
            mock_traffic_cls.read_all_classes.side_effect = tc_read

        qos_mgr._setupInternal()

        # run async loop once to ensure ready items are cleared
        loop._run_once()

        self.assertTrue(not qos_mgr._qos_store)
        # verify that unreferenced qos configs are purged from the system
        if self.config["qos"]["impl"] == QosImplType.OVS_METER:
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 2)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 15)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 13)
            mock_meter_cls.del_meter.assert_any_call(MagicMock, 11)
        else:
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 2)
            mock_traffic_cls.delete_class.assert_any_call(self.ul_intf, 15)
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 13)
            mock_traffic_cls.delete_class.assert_any_call(self.dl_intf, 11)

    def testSanity(self,):
        for impl_type in (QosImplType.LINUX_TC, QosImplType.OVS_METER,):
            self.config["qos"]["impl"] = impl_type
            self._testSanity()

    def testMultipleSubscribers(self, ):
        for impl_type in (QosImplType.LINUX_TC, QosImplType.OVS_METER,):
            self.config["qos"]["impl"] = impl_type
            self._testMultipleSubscribers()

    def testRedisConnectionFailure(self,):
        self.config["qos"]["impl"] = QosImplType.LINUX_TC
        qos_mgr = QosManager(MagicMock, asyncio.new_event_loop(), self.config)

        redisConnFailureCount = 5

        def mockRedisAvail(*args, **kw):
            if not hasattr(mockRedisAvail, "count"):
                mockRedisAvail.count = 0
            mockRedisAvail.count += 1
            if mockRedisAvail.count > redisConnFailureCount:
                return True
            return False

        qos_mgr.redisAvailable = mockRedisAvail
        qos_mgr._setupInternal = lambda: True
        with self.assertLogs('pipelined.qos.common', level='INFO') as cm:
            qos_mgr.setup()
        self.assertTrue(len(cm.output), redisConnFailureCount)
        for output in cm.output:
            self.assertTrue("failed to connect to redis" in output)

    def testUncleanRestart(self,):
        with patch.dict(self.config, {"clean_restart": False}):
            for impl_type in (QosImplType.LINUX_TC, QosImplType.OVS_METER):
                self.config["qos"]["impl"] = impl_type
                self._testUncleanRestart()


class TestMeters(unittest.TestCase):
    @patch('magma.pipelined.qos.qos_meter_impl.MeterClass')
    def testBrokenMeter(self, meter_cls):
        config = {
            'qos': {
                'enable': True,
                'max_rate': 10000000,
                'ovs_meter': {
                    'min_idx': 2,
                    'max_idx': 100000,
                },
            },
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
    @patch('subprocess.check_call')
    @patch('magma.pipelined.qos.qos_tc_impl.TrafficClass.delete_class')
    def testSanity(self, mock_del_cls, mock_check_call):
        TrafficClass.create_class("en0", 2, 10)
        mock_del_cls.assert_called_with("en0", 2, show_error=False,
                                        throw_except=False)
        mock_check_call.assert_any_call(['tc', 'class', 'add', 'dev', 'en0',
                                         'parent', '1:fffe', 'classid', '1:2',
                                         'htb', 'rate', '12000', 'ceil', '10'])
        mock_check_call.assert_any_call(['tc', 'qdisc', 'add', 'dev', 'en0',
                                         'parent', '1:2', 'fq_codel'])
        mock_check_call.assert_any_call(['tc', 'filter', 'add', 'dev', 'en0',
                                         'protocol', 'ip', 'parent', '1:',
                                         'prio', '1', 'handle', '2', 'fw',
                                         'flowid', '1:2'])

    @patch('subprocess.check_call')
    @patch('os.geteuid', return_value=1)
    def testSudoUser(self, _, mock_check_call):
        TrafficClass.init_qdisc("en0")
        mock_check_call.assert_any_call(['sudo', 'tc', 'qdisc', 'add',
                                         'dev', 'en0', 'root', 'handle',
                                         '1:', 'htb'])

    @patch('subprocess.check_output')
    def testReadAllClasses(self, mock_check_output):
        tc_output = '''
class htb 1:1 parent 1:fffe prio 0 rate 12Kbit ceil 1Gbit burst 1599b\n
class htb 1:fffe root rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b\n
class htb 1:2 root rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b\n
class htb 1:5 root rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b\n
class htb 1:7 root rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b\n
class htb 1:8 root rate 1Gbit ceil 1Gbit burst 1375b cburst 1375b\n
class fq_codel 8005:23b parent 8005: \n
class fq_codel 8005:383 parent 8005: \n
'''
        mock_check_output.return_value = bytes(tc_output, 'utf-8')
        qid_list = TrafficClass.read_all_classes("testIntf")
        self.assertTrue(qid_list == [1, 65534, 2, 5, 7, 8])

    @patch('subprocess.check_call')
    def testError(self, mock_check_call):
        def dummy_check_call(*args):
            raise subprocess.CalledProcessError(returncode=1, cmd="tc")

        mock_check_call.side_effect = dummy_check_call
        with self.assertLogs('pipelined.qos.qos_tc_impl', level='ERROR') as cm:
            try:
                TrafficClass.init_qdisc("en0", throw_except=True,
                                        show_error=True)
                self.fail("init_qdisc didn't raise an exception")
            except subprocess.CalledProcessError:
                pass
        self.assertTrue("error running tc qdisc add dev" in cm.output[0])
        args = ['tc', 'qdisc', 'add', 'dev', 'en0', 'root', 'handle',
                '1:', 'htb']
        mock_check_call.assert_called_with(args)
