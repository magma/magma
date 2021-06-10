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
import json
import logging
import threading
import traceback
from enum import Enum
from typing import Dict, List  # noqa

from lte.protos.policydb_pb2 import FlowMatch
from magma.configuration.service_configs import load_service_config
from magma.pipelined.qos.qos_meter_impl import MeterManager
from magma.pipelined.qos.qos_tc_impl import TCManager, TrafficClass
from magma.pipelined.qos.types import (
    QosInfo,
    get_data,
    get_data_json,
    get_key,
    get_key_json,
    get_subscriber_data,
    get_subscriber_key,
)
from magma.pipelined.qos.utils import QosStore

LOG = logging.getLogger("pipelined.qos.common")
# LOG.setLevel(logging.DEBUG)


def normalize_imsi(imsi: str) -> str:
    imsi = imsi.lower()
    if imsi.startswith("imsi"):
        imsi = imsi[4:]
    return imsi


class QosImplType(Enum):
    LINUX_TC = "linux_tc"
    OVS_METER = "ovs_meter"

    @staticmethod
    def list():
        return list(map(lambda t: t.value, QosImplType))


class SubscriberSession(object):
    def __init__(self, ip_addr: str):
        self.ip_addr = ip_addr
        self.ambr_ul = 0
        self.ambr_ul_leaf = 0
        self.ambr_dl = 0
        self.ambr_dl_leaf = 0

        self.rules = set()

    def set_ambr(self, d: FlowMatch.Direction, root: int, leaf: int) -> None:
        if d == FlowMatch.UPLINK:
            self.ambr_ul = root
            self.ambr_ul_leaf = leaf
        else:
            self.ambr_dl = root
            self.ambr_dl_leaf = leaf

    def get_ambr(self, d: FlowMatch.Direction) -> int:
        if d == FlowMatch.UPLINK:
            return self.ambr_ul
        return self.ambr_dl

    def get_ambr_leaf(self, d: FlowMatch.Direction) -> int:
        if d == FlowMatch.UPLINK:
            return self.ambr_ul_leaf
        return self.ambr_dl_leaf


class SubscriberState(object):
    """
     Map subscriber -> sessions
    """

    def __init__(self, imsi: str, qos_store: Dict):
        self.imsi = imsi
        self.rules = {}  # rule -> qos handle
        self.sessions = {}  # IP -> [sessions(ip, qos_argumets, rule_no), ...]
        self._redis_store = qos_store

    def check_empty(self, ) -> bool:
        return not self.rules and not self.sessions

    def get_or_create_session(self, ip_addr: str):
        session = self.sessions.get(ip_addr)
        if not session:
            session = SubscriberSession(ip_addr)
            self.sessions[ip_addr] = session
        return session

    def remove_session(self, ip_addr: str) -> None:
        del self.sessions[ip_addr]

    def find_session_with_rule(self, rule_num: int) -> SubscriberSession:
        session_with_rule = None
        for session in self.sessions.values():
            if rule_num in session.rules:
                session_with_rule = session
                break
        return session_with_rule

    def _update_rules_map(self, ip_addr: str, rule_num: int, d: FlowMatch.Direction,
                    qos_data) -> None:

        if rule_num not in self.rules:
            self.rules[rule_num] = []

        session = self.get_or_create_session(ip_addr)
        session.rules.add(rule_num)
        self.rules[rule_num].append((d, qos_data))

    def update_rule(self, ip_addr: str, rule_num: int, d: FlowMatch.Direction,
                    qos_handle: int, ambr: int, leaf: int) -> None:
        k = get_key_json(get_subscriber_key(self.imsi, ip_addr, rule_num, d))
        qos_data = get_data_json(get_subscriber_data(qos_handle, ambr, leaf))
        LOG.debug("Update: %s -> %s", k, qos_data)
        self._redis_store[k] = qos_data

        self._update_rules_map(ip_addr, rule_num, d, qos_data)

    def remove_rule(self, rule_num: int) -> None:
        session_with_rule = self.find_session_with_rule(rule_num)
        if session_with_rule:
            for (d, _) in self.rules[rule_num]:
                k = get_subscriber_key(self.imsi, session_with_rule.ip_addr,
                                       rule_num, d)
                if get_key_json(k) in self._redis_store:
                    del self._redis_store[get_key_json(k)]

        del self.rules[rule_num]
        session_with_rule.rules.remove(rule_num)

    def find_rule(self, rule_num: int):
        return self.rules.get(rule_num)

    def get_all_rules(self, ) -> List:
        return self.rules

    def get_all_empty_sessions(self, ) -> List:
        return [s for s in self.sessions.values() if not s.rules]

    def get_qos_handle(self, rule_num: int, direction: FlowMatch.Direction) -> int:
        rule = self.rules.get(rule_num)
        if rule:
            for d, qos_data in rule:
                if d == direction:
                    _, qid, _, _ = get_data(qos_data)
                    return qid
        return 0


class QosManager(object):
    # TODO: convert QosManager to singleton class.
    qos_mgr = None
    # protect QoS object create and delete across all QoSManager Objects.
    lock = threading.Lock()

    @staticmethod
    def get_qos_manager(datapath, loop, config):
        if QosManager.qos_mgr:
            LOG.debug("Got QosManager instance")
            return QosManager.qos_mgr
        QosManager.qos_mgr = QosManager(datapath, loop, config)
        QosManager.qos_mgr.setup()
        return QosManager.qos_mgr

    @staticmethod
    def get_impl(datapath, loop, config):
        try:
            impl_type = QosImplType(config["qos"]["impl"])
        except ValueError:
            LOG.error("%s is not a valid qos impl type", impl_type)
            raise

        if impl_type == QosImplType.OVS_METER:
            return MeterManager(datapath, loop, config)
        else:
            return TCManager(datapath, loop, config)

    @classmethod
    def debug(cls, _, __, ___):
        config = load_service_config('pipelined')
        qos_impl_type = QosImplType(config["qos"]["impl"])
        qos_store = QosStore(cls.__name__)
        for k, v in qos_store.items():
            _, imsi, ip_addr, rule_num, d = get_key(k)
            _, qid, ambr, leaf = get_data(v)
            print('imsi :', imsi)
            print('ip_addr :', ip_addr)
            print('rule_num :', rule_num)
            print('direction :', d)
            print('qos_handle:', qid)
            print('qos_handle_ambr:', ambr)
            print('qos_handle_ambr_leaf:', leaf)
            if qos_impl_type == QosImplType.OVS_METER:
                MeterManager.dump_meter_state(v)
            else:
                intf = 'nat_iface' if d == FlowMatch.UPLINK else 'enodeb_iface'
                print("Dev: ", config[intf])
                TrafficClass.dump_class_state(config[intf], qid)
                if leaf and leaf != qid:
                    print("Leaf:")
                    TrafficClass.dump_class_state(config[intf], leaf)
                if ambr:
                    print("AMBR (parent):")
                    TrafficClass.dump_class_state(config[intf], ambr)

        if qos_impl_type == QosImplType.LINUX_TC:
            dev = config['nat_iface']
            print("Root stats for: ", dev)
            TrafficClass.dump_root_class_stats(dev)
            dev = config['enodeb_iface']
            print("Root stats for: ", dev)
            TrafficClass.dump_root_class_stats(dev)

    def _is_redis_available(self):
        try:
            self._redis_store.client.ping()
        except ConnectionError:
            return False
        return True

    def __init__(self, datapath, loop, config):
        self._qos_enabled = config["qos"]["enable"]
        if not self._qos_enabled:
            return
        self._apn_ambr_enabled = config["qos"].get("apn_ambr_enabled", True)
        LOG.info("QoS: apn_ambr_enabled: %s", self._apn_ambr_enabled)
        self._clean_restart = config["clean_restart"]
        self._subscriber_state = {}
        self._loop = loop
        self.impl = QosManager.get_impl(datapath, loop, config)
        self._redis_store = QosStore(self.__class__.__name__)
        self._initialized = False
        self._redis_conn_retry_secs = 1

    def setup(self):
        with QosManager.lock:
            if not self._qos_enabled:
                return

            if self._is_redis_available():
                return self._setupInternal()
            else:
                LOG.info("failed to connect to redis..retrying in %d secs",
                         self._redis_conn_retry_secs)
                self._loop.call_later(self._redis_conn_retry_secs, self.setup)

    def _setupInternal(self):
        if self._initialized:
            return
        if self._clean_restart:
            LOG.info("Qos Setup: clean start")
            self.impl.destroy()
            self._redis_store.clear()
            self.impl.setup()
            self._initialized = True
        else:
            # read existing state from qos_impl
            LOG.info("Qos Setup: recovering existing state")
            self.impl.setup()

            cur_qos_state, apn_qid_list = self.impl.read_all_state()
            LOG.debug("Initial qos_state -> %s", json.dumps(cur_qos_state, indent=1))
            LOG.debug("apn_qid_list -> %s", apn_qid_list)
            LOG.debug("Redis state: %s", self._redis_store)
            try:
                # populate state from db
                in_store_qid = set()
                in_store_ambr_qid = set()
                purge_store_set = set()
                for rule, sub_data in self._redis_store.items():
                    _, qid, ambr, leaf = get_data(sub_data)
                    if qid not in cur_qos_state:
                        LOG.warning("missing qid: %s in TC", qid)
                        purge_store_set.add(rule)
                        continue
                    if ambr and ambr != 0 and ambr not in cur_qos_state:
                        purge_store_set.add(rule)
                        LOG.warning("missing ambr class: %s of qid %d",ambr, qid)
                        continue
                    if leaf and leaf != 0 and leaf not in cur_qos_state:
                        purge_store_set.add(rule)
                        LOG.warning("missing leaf class: %s of qid %d",leaf, qid)
                        continue

                    if ambr:
                        qid_state = cur_qos_state[qid]
                        if qid_state['ambr_qid'] != ambr:
                            purge_store_set.add(rule)
                            LOG.warning("Inconsistent amber class: %s of qid %d", qid_state['ambr_qid'], ambr)
                            continue

                    in_store_qid.add(qid)
                    if ambr:
                        in_store_qid.add(ambr)
                        in_store_ambr_qid.add(ambr)
                    in_store_qid.add(leaf)

                    _, imsi, ip_addr, rule_num, direction = get_key(rule)

                    subscriber = self._get_or_create_subscriber(imsi)
                    subscriber.update_rule(ip_addr, rule_num, direction, qid, ambr, leaf)
                    session = subscriber.get_or_create_session(ip_addr)
                    session.set_ambr(direction, ambr, leaf)

                # purge entries from qos_store
                for rule in purge_store_set:
                    LOG.debug("purging qos_store entry %s", rule)
                    del self._redis_store[rule]

                # purge unreferenced qos configs from system
                # Step 1. Delete child nodes
                lost_and_found_apn_list = set()
                for qos_handle in cur_qos_state:
                    if qos_handle not in in_store_qid:
                        if qos_handle in apn_qid_list:
                            lost_and_found_apn_list.add(qos_handle)
                        else:
                            LOG.debug("removing qos_handle %d", qos_handle)
                            self.impl.remove_qos(qos_handle,
                                                 cur_qos_state[qos_handle]['direction'],
                                                 recovery_mode=True)

                if len(lost_and_found_apn_list) > 0:
                    # Step 2. delete qos ambr without any leaf nodes
                    for qos_handle in lost_and_found_apn_list:
                        if qos_handle not in in_store_ambr_qid:
                            LOG.debug("removing apn qos_handle %d", qos_handle)
                            self.impl.remove_qos(qos_handle,
                                                 cur_qos_state[qos_handle]['direction'],
                                                 recovery_mode=True,
                                                 skip_filter=True)
                final_qos_state, _ = self.impl.read_all_state()
                LOG.info("final_qos_state -> %s", json.dumps(final_qos_state, indent=1))
                LOG.info("final_redis state -> %s", self._redis_store)
            except Exception as e:  # pylint: disable=broad-except
                # in case of any exception start clean slate

                LOG.error("error %s. restarting clean %s", e, traceback.format_exc())
                self._clean_restart = True

            self._initialized = True

    def _get_or_create_subscriber(self, imsi):
        subscriber_state = self._subscriber_state.get(imsi)
        if not subscriber_state:
            subscriber_state = SubscriberState(imsi, self._redis_store)
            self._subscriber_state[imsi] = subscriber_state
        return subscriber_state

    def add_subscriber_qos(
            self,
            imsi: str,
            ip_addr: str,
            apn_ambr: int,
            rule_num: int,
            direction: FlowMatch.Direction,
            qos_info: QosInfo,
    ):
        with QosManager.lock:
            if not self._qos_enabled or not self._initialized:
                LOG.debug("add_subscriber_qos: not enabled or initialized")
                return None, None

            LOG.debug("adding qos for imsi %s rule_num %d direction %d apn_ambr %d, %s",
                      imsi, rule_num, direction, apn_ambr, qos_info)

            imsi = normalize_imsi(imsi)

            # ip_addr identifies a specific subscriber session, each subscriber session
            # must be associated with a default bearer and can be associated with dedicated
            # bearers. APN AMBR specifies the aggregate max bit rate for a specific
            # subscriber across all the bearers. Queues for dedicated bearers will be
            # children of default bearer Queues. In case the dedicated bearers exceed the
            # rate, then they borrow from the default bearer queue
            subscriber_state = self._get_or_create_subscriber(imsi)

            qos_handle = subscriber_state.get_qos_handle(rule_num, direction)
            if qos_handle:
                LOG.debug("qos exists for imsi %s rule_num %d direction %d",
                          imsi, rule_num, direction)

                return self.impl.get_action_instruction(qos_handle)

            ambr_qos_handle_root = 0
            ambr_qos_handle_leaf = 0
            if self._apn_ambr_enabled and apn_ambr > 0:
                session = subscriber_state.get_or_create_session(ip_addr)
                ambr_qos_handle_root = session.get_ambr(direction)
                LOG.debug("existing root rec: ambr_qos_handle_root %d", ambr_qos_handle_root)

                if not ambr_qos_handle_root:
                    ambr_qos_handle_root = self.impl.add_qos(direction, QosInfo(gbr=None, mbr=apn_ambr), skip_filter=True)
                    if not ambr_qos_handle_root:
                        LOG.error('Failed adding root ambr qos mbr %u direction %d',
                                  apn_ambr, direction)
                        return None, None
                    else:
                        LOG.debug('Added root ambr qos mbr %u direction %d qos_handle %d ',
                                  apn_ambr, direction, ambr_qos_handle_root)

                ambr_qos_handle_leaf = session.get_ambr_leaf(direction)
                LOG.debug("existing leaf rec: ambr_qos_handle_leaf %d", ambr_qos_handle_leaf)

                if not ambr_qos_handle_leaf:
                    ambr_qos_handle_leaf = self.impl.add_qos(direction,
                                                             QosInfo(gbr=None, mbr=apn_ambr),
                                                             parent=ambr_qos_handle_root)
                    if ambr_qos_handle_leaf:
                        session.set_ambr(direction, ambr_qos_handle_root, ambr_qos_handle_leaf)
                        LOG.debug('Added ambr qos mbr %u direction %d qos_handle %d/%d ',
                                  apn_ambr, direction, ambr_qos_handle_root, ambr_qos_handle_leaf)
                    else:
                        LOG.error('Failed adding leaf ambr qos mbr %u direction %d',
                                  apn_ambr, direction)
                        self.impl.remove_qos(ambr_qos_handle_root, direction, skip_filter=True)
                        return None, None
                qos_handle = ambr_qos_handle_leaf

            if qos_info:
                qos_handle = self.impl.add_qos(direction, qos_info, parent=ambr_qos_handle_root)
                LOG.debug("Added ded brr handle: %d", qos_handle)
                if qos_handle:
                    LOG.debug('Adding qos %s direction %d qos_handle %d ',
                              qos_info, direction, qos_handle)
                else:
                    LOG.error('Failed adding qos %s direction %d', qos_info, direction)
                    return None, None

            if qos_handle:
                subscriber_state.update_rule(ip_addr, rule_num, direction,
                                             qos_handle, ambr_qos_handle_root, ambr_qos_handle_leaf)
                return self.impl.get_action_instruction(qos_handle)
            return None, None

    def remove_subscriber_qos(self, imsi: str = "", del_rule_num: int = -1):
        with QosManager.lock:
            if not self._qos_enabled or not self._initialized:
                LOG.debug("remove_subscriber_qos: not enabled or initialized")
                return

            LOG.debug("removing Qos for imsi %s del_rule_num %d", imsi, del_rule_num)
            if not imsi:
                LOG.error('imsi %s invalid, failed removing', imsi)
                return

            imsi = normalize_imsi(imsi)
            subscriber_state = self._subscriber_state.get(imsi)
            if not subscriber_state:
                LOG.debug('imsi %s not found, nothing to remove ', imsi)
                return

            to_be_deleted_rules = set()
            qid_to_remove = {}
            qid_in_use = set()
            if del_rule_num == -1:
                # deleting all rules for the subscriber
                rules = subscriber_state.get_all_rules()
                for (rule_num, rule) in rules.items():
                    for (d, qos_data) in rule:
                        _, qid, ambr, leaf = get_data(qos_data)
                        if ambr != qid:
                            qid_to_remove[qid] = d
                            if leaf and leaf != qid:
                                qid_to_remove[leaf] = d

                    to_be_deleted_rules.add(rule_num)
            else:
                rule = subscriber_state.find_rule(del_rule_num)
                if rule:
                    rules = subscriber_state.get_all_rules()
                    for (rule_num, rule) in rules.items():
                        for (d, qos_data) in rule:
                            _, qid, ambr, leaf = get_data(qos_data)
                            if rule_num == del_rule_num:
                                if ambr != qid:
                                    qid_to_remove[qid] = d
                                    if leaf and leaf != qid:
                                        qid_to_remove[leaf] = d
                            else:
                                qid_in_use.add(qid)

                    LOG.debug("removing rule %s %s ", imsi, del_rule_num)
                    to_be_deleted_rules.add(del_rule_num)
                else:
                    LOG.debug("unable to find rule_num %d for imsi %s", del_rule_num, imsi)

            for (qid, d) in qid_to_remove.items():
                if qid not in qid_in_use:
                    self.impl.remove_qos(qid, d)

            for rule_num in to_be_deleted_rules:
                subscriber_state.remove_rule(rule_num)

            # purge sessions with no rules
            for session in subscriber_state.get_all_empty_sessions():
                for d in (FlowMatch.UPLINK, FlowMatch.DOWNLINK):
                    ambr_qos_handle = session.get_ambr(d)
                    if ambr_qos_handle:
                        LOG.debug("removing root ambr qos handle %d direction %d", ambr_qos_handle, d)
                        self.impl.remove_qos(ambr_qos_handle, d, skip_filter=True)
                LOG.debug("purging session %s %s ", imsi, session.ip_addr)
                subscriber_state.remove_session(session.ip_addr)

            # purge subscriber state with no rules
            if subscriber_state.check_empty():
                LOG.debug("purging subscriber state for %s, empty rules and sessions", imsi)
                del self._subscriber_state[imsi]
