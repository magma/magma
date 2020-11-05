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
from typing import List, Dict  # noqa

import asyncio
import logging
from enum import Enum
from lte.protos.policydb_pb2 import FlowMatch
from magma.pipelined.qos.qos_meter_impl import MeterManager
from magma.pipelined.qos.qos_tc_impl import TCManager, TrafficClass
from magma.pipelined.qos.types import QosInfo, get_json, get_key, get_subscriber_key
from magma.pipelined.qos.utils import QosStore
from magma.configuration.service_configs import load_service_config

LOG = logging.getLogger("pipelined.qos.common")
#LOG.setLevel(logging.DEBUG)


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
    def __init__(self, imsi: str, qos_store : Dict):
        self.imsi = imsi
        self.rules = {}
        self.sessions = {}
        self._qos_store = qos_store

    def check_empty(self,) -> bool:
        return not self.rules and not self.sessions

    def get_or_create_session(self, ip_addr: str):
        session = self.sessions.get(ip_addr)
        if not session:
            session = SubscriberSession(ip_addr)
            self.sessions[ip_addr] = session
        return session

    def remove_session(self, ip_addr: str) -> None:
        del self.sessions[ip_addr]

    def find_session_with_rule(self, rule_num: int) ->SubscriberSession:
        session_with_rule = None
        for session in self.sessions.values():
            if rule_num in session.rules:
                session_with_rule = session
                break
        return session_with_rule

    def update_rule(self, ip_addr: str, rule_num: int, d: FlowMatch.Direction,
                    qos_handle: int) -> None:
        k = get_subscriber_key(self.imsi, ip_addr, rule_num, d)
        self._qos_store[get_json(k)] = qos_handle

        if rule_num not in self.rules:
            self.rules[rule_num] = []

        session = self.get_or_create_session(ip_addr)
        session.rules.add(rule_num)
        self.rules[rule_num].append((d, qos_handle))

    def remove_rule(self, rule_num: int) -> None:
        session_with_rule = self.find_session_with_rule(rule_num)
        if session_with_rule:
            for (d, _) in self.rules[rule_num]:
                k = get_subscriber_key(self.imsi, session_with_rule.ip_addr,
                                        rule_num, d)
                del self._qos_store[get_json(k)]

        del self.rules[rule_num]
        session_with_rule.rules.remove(rule_num)

    def find_rule(self, rule_num: int):
        return self.rules.get(rule_num)

    def get_all_rules(self, ) -> List:
        return self.rules

    def get_all_empty_sessions(self,) -> List:
        return [ s for s in self.sessions.values() if not s.rules]

    def get_qos_handle(self, rule_num: int, direction: FlowMatch.Direction) -> int:
        rule = self.rules.get(rule_num)
        if rule:
            for d, qos_handle in rule:
                if d == direction:
                    return qos_handle
        return 0


class QosManager(object):
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
            print('imsi :', imsi)
            print('ip_addr :', ip_addr)
            print('rule_num :', rule_num)
            print('direction :', d)
            print('qos_handle:', v)
            if qos_impl_type == QosImplType.OVS_METER:
                MeterManager.dump_meter_state(v)
            else:
                intf = 'nat_iface' if d == FlowMatch.UPLINK else 'enodeb_iface'
                TrafficClass.dump_class_state(config[intf], v)

    def redisAvailable(self):
        try:
            self._qos_store.client.ping()
        except ConnectionError:
            return False
        return True

    def __init__(self, datapath, loop, config):
        self._qos_enabled = config["qos"]["enable"]
        if not self._qos_enabled:
            return

        self._clean_restart = config["clean_restart"]
        self._subscriber_state = {}
        self._loop = loop
        self.impl = QosManager.get_impl(datapath, loop, config)
        self._qos_store = QosStore(self.__class__.__name__)
        self._initialized = False
        self._redis_conn_retry_secs = 1

    def setup(self):
        if not self._qos_enabled:
            return

        if self.redisAvailable():
            return self._setupInternal()
        else:
            LOG.info("failed to connect to redis..retrying in %d secs",
                     self._redis_conn_retry_secs)
            self._loop.call_later(self._redis_conn_retry_secs, self.setup)

    def _setupInternal(self):
        if self._clean_restart:
            LOG.info("Qos Setup: clean start")
            self.impl.destroy()
            self._qos_store.clear()
            self.impl.setup()
            self._initialized = True
        else:
            # read existing state from qos_impl
            LOG.info("Qos Setup: recovering existing state")

            def callback(fut):
                LOG.debug("read_all_state complete => \n%s", fut.result())
                qos_state = fut.result()
                try:
                    # populate state from db
                    in_store_qid = set()
                    purge_store_set = set()
                    for k, v in self._qos_store.items():
                        if v not in qos_state:
                            purge_store_set.add(k)
                            continue
                        in_store_qid.add(v)
                        _, imsi, ip_addr, rule_num, d = get_key(k)
                        subscriber = self.get_or_create_subscriber(imsi)
                        subscriber.update_rule(ip_addr, rule_num, d, v)

                        qid_state = qos_state[v]
                        if qid_state['ambr_qid'] != 0:
                            session = subscriber.get_or_create_session(ip_addr)
                            ambr_qid = qid_state['ambr_qid']
                            leaf = 0
                            # its AMBR QoS handle if its config matches with parent
                            # pylint: disable=too-many-function-args
                            if self.impl.same_qos_config(d, ambr_qid, v):
                                leaf = v
                            session.set_ambr(d, qid_state['ambr_qid'], leaf)

                    # purge entries from qos_store
                    for k in purge_store_set:
                        LOG.debug("purging qos_store entry %s qos_handle", k)
                        del self._qos_store[k]

                    # purge unreferenced qos configs from system
                    for qos_handle in qos_state:
                        if qos_handle not in in_store_qid:
                            LOG.debug("removing qos_handle %d", qos_handle)
                            self.impl.remove_qos(qos_handle,
                                qos_state[qos_handle]['direction'], recovery_mode=True)

                    self._initialized = True
                    LOG.info("init complete with state recovered successfully")
                except Exception as e:  # pylint: disable=broad-except
                    # in case of any exception start clean slate
                    LOG.error("error %s. restarting clean", str(e))
                    self._clean_restart = True
                    self.setup()

            asyncio.ensure_future(self.impl.read_all_state(), loop=self._loop).add_done_callback(callback)

    def get_or_create_subscriber(self, imsi):
        subscriber_state = self._subscriber_state.get(imsi)
        if not subscriber_state:
            subscriber_state = SubscriberState(imsi, self._qos_store)
            self._subscriber_state[imsi] = subscriber_state
        return subscriber_state

    def add_subscriber_qos(
        self,
        imsi: str,
        ip_addr: str,
        apn_ambr : int,
        rule_num: int,
        direction: FlowMatch.Direction,
        qos_info: QosInfo,
    ):
        if not self._qos_enabled or not self._initialized:
            LOG.error("add_subscriber_qos: not enabled or initialized")
            return None, None

        LOG.debug("adding qos for imsi %s rule_num %d direction %d apn_ambr %d, qos_info %s",
                   imsi, rule_num, direction, apn_ambr, qos_info)

        imsi = normalize_imsi(imsi)

        # ip_addr identifies a specific subscriber session, each subscriber session
        # must be associated with a default bearer and can be associated with dedicated
        # bearers. APN AMBR specifies the aggregate max bit rate for a specific
        # subscriber across all the bearers. Queues for dedicated bearers will be
        # children of default bearer Queues. In case the dedicated bearers exceed the
        # rate, then they borrow from the default bearer queue
        subscriber_state = self.get_or_create_subscriber(imsi)

        qos_handle = subscriber_state.get_qos_handle(rule_num, direction)
        LOG.debug("existing rec: qos_handle %d", qos_handle)
        if qos_handle:
            LOG.debug("qos exists for imsi %s rule_num %d direction %d",
                      imsi, rule_num, direction)
            return self.impl.get_action_instruction(qos_handle)

        ambr_qos_handle_root = None
        if apn_ambr > 0:
            session = subscriber_state.get_or_create_session(ip_addr)
            ambr_qos_handle_root = session.get_ambr(direction)
            LOG.debug("existing root rec: ambr_qos_handle_root %d", ambr_qos_handle_root)

            if not ambr_qos_handle_root:
                ambr_qos_handle_root = self.impl.add_qos(direction, QosInfo(gbr=None, mbr=apn_ambr))
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
                    self.impl.remove_qos(ambr_qos_handle_root, direction)
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

        LOG.debug("qos_handle %d", qos_handle)
        subscriber_state.update_rule(ip_addr, rule_num, direction, qos_handle)
        return self.impl.get_action_instruction(qos_handle)

    def remove_subscriber_qos(self, imsi: str = "", rule_num: int = -1):
        if not self._qos_enabled or not self._initialized:
            LOG.error("remove_subscriber_qos: not enabled or initialized")
            return

        LOG.debug("removing Qos for imsi %s rule_num %d", imsi, rule_num)
        if not imsi:
            LOG.error('imsi %s invalid, failed removing', imsi)
            return

        imsi = normalize_imsi(imsi)
        subscriber_state = self._subscriber_state.get(imsi)
        if not subscriber_state:
            LOG.debug('imsi %s not found, nothing to remove ', imsi)
            return

        to_be_deleted_rules = []
        if rule_num == -1:
            # deleting all rules for the subscriber
            rules = subscriber_state.get_all_rules()
            for (rule_num, rule) in rules.items():
                LOG.debug("removing rule %s %s ", imsi, rule_num)
                session_with_rule = subscriber_state.find_session_with_rule(rule_num)
                for (d, qos_handle) in rule:
                    if session_with_rule.get_ambr(d) != qos_handle:
                        self.impl.remove_qos(qos_handle, d)
                to_be_deleted_rules.append(rule_num)
        else:
            rule = subscriber_state.find_rule(rule_num)
            if rule is None:
                LOG.error("unable to find rule_num %d for imsi %s", rule_num, imsi)
                return

            session_with_rule = subscriber_state.find_session_with_rule(rule_num)
            for (d, qos_handle) in rule:
                if session_with_rule.get_ambr(d) != qos_handle:
                    self.impl.remove_qos(qos_handle, d)
            LOG.debug("removing rule %s %s ", imsi, rule_num)
            to_be_deleted_rules.append(rule_num)

        for rule_num in to_be_deleted_rules:
            subscriber_state.remove_rule(rule_num)

        # purge sessions with no rules
        for session in subscriber_state.get_all_empty_sessions():
            for d in (FlowMatch.UPLINK, FlowMatch.DOWNLINK):
                ambr_qos_handle = session.get_ambr(d)
                if ambr_qos_handle:
                    LOG.debug("removing root ambr qos handle %d direction %d", ambr_qos_handle, d)
                    self.impl.remove_qos(ambr_qos_handle, d)
            LOG.debug("purging session %s %s ", imsi, session.ip_addr)
            subscriber_state.remove_session(session.ip_addr)

        # purge subscriber state with no rules
        if subscriber_state.check_empty():
            LOG.debug("purging subscriber state for %s, empty rules and sessions", imsi)
            del self._subscriber_state[imsi]