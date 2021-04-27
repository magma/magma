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

import logging
import os
import re
import subprocess
from collections import namedtuple
from concurrent.futures import Future
from difflib import unified_diff
from typing import Dict, List, Optional
from unittest import TestCase, mock
from unittest.mock import MagicMock

import fakeredis
import netifaces
from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from lte.protos.pipelined_pb2 import (
    SetupFlowsResult,
    SetupPolicyRequest,
    SetupUEMacRequest,
    UpdateSubscriberQuotaStateRequest,
)
from magma.pipelined.app.base import global_epoch
from magma.pipelined.bridge_util import BridgeTools
from magma.pipelined.openflow import flows
from magma.pipelined.service_manager import ServiceManager
from magma.pipelined.tests.app.exceptions import (
    BadConfigError,
    ServiceRunningError,
)
from magma.pipelined.tests.app.flow_query import RyuDirectFlowQuery
from magma.pipelined.tests.app.start_pipelined import StartThread
from ryu.lib import hub

"""
Pipelined test util functions can be used for testing pipelined, the usage of
these functions can be seen in pipelined/tests/test_*.py files
"""

SubTest = namedtuple('SubTest', ['context', 'isolator', 'flowtest_list'])
PktsToSend = namedtuple('PacketToSend', ['pkt', 'num'])
QueryMatch = namedtuple('QueryMatch', ['pkts', 'flow_count'])

SNAPSHOT_DIR = 'snapshots/'
SNAPSHOT_EXTENSION = '.snapshot'


# Tuple for FlowVerifier, class wrapper needed becuse of optional flow_count
class FlowTest(namedtuple('FlowTest', ['query', 'match_num', 'flow_count'])):
    __slots__ = ()

    def __new__(cls, query, match_num, flow_count=None):
        return super(FlowTest, cls).__new__(cls, query, match_num, flow_count)


class WaitTimeExceeded(Exception):
    pass


class FlowVerifier:
    """
    FlowVerifier controlls flow pkt matching verification

    When run as a context the verifier will automatically compute deltas of
    packets matched and compare them to the expected pkt differences passed in
    the FlowTest tuples
    """

    def __init__(self, flow_test_list, wait_func):
        """
        Args:
        flow_test_list [FlowTest]: list of FlowTests to query stats for
                                   -> if match_num is number:
                                      verify differences in matched packets
                                   -> if match_num is None:
                                      verify that flows were instantiated
                                      (assert flow query return length == 1)

        wait_func (func): function to run before getting final matched packet
                          num for FlowTests. This function is used to ensure
                          that all pkts were processed by OVS. Currently we
                          have 2 wait functions:
                              'wait_after_send', 'wait_for_enforcement_stats'
        """
        self._initial = []
        self._final = []
        self._flow_tests = flow_test_list
        self._wait_func = wait_func
        self._done = False

    def __enter__(self):
        """
        Runs on entering 'with' (query initial stats)
        """
        self._get_initial_pkt_counts()

    def __exit__(self, type, value, traceback):
        """
        Runs after finishing 'with' (query final stats)
        """
        try:
            self._wait_func()
        except WaitTimeExceeded as e:
            TestCase().fail(e)
        self._get_final_pkt_counts()

    def _get_initial_pkt_counts(self):
        for flow in self._flow_tests:
            matched = flow.query.lookup()
            self._initial.append(QueryMatch(pkt_total(matched), len(matched)))
        self._done = True

    def _get_final_pkt_counts(self):
        for flow in self._flow_tests:
            matched = flow.query.lookup()
            self._final.append(QueryMatch(pkt_total(matched), len(matched)))

    def verify(self):
        """
        Verifies that all queries matches expeceted num of packets
        Assume if the flow.match_num is None we need to verify num of flows
        found by the query, should be one (used when testing if flows
        were added to ovs)
        """
        for f, i, test in zip(self._final, self._initial, self._flow_tests):
            if test.flow_count is not None:
                print(f)
                print(test)
                TestCase().assertEqual(f.flow_count, test.flow_count)
            TestCase().assertEqual(f.pkts, i.pkts + test.match_num)

    def get_query_delts(self):
        """
        If a custom testcase assertion is required
        Return:
            [(pkt_delta, flow_count_delta)]: deltas for all flow test queries

        """
        if self._done is False:
            logging.error("Test didn't finish, can't access final pkt stats")
            return None

        return [(f.pkts - i.pkts, f.flow_count - i.flow_count)
                for f, i in zip(self._final, self._initial)]


def start_ryu_app_thread(test_setup):
    """
    Starts the ryu apps, using the information from test_setup config

    Args:
        test_setup (EnforcementStatsController): test_setup config
    Return:
        thread {rule:RuleRecord}: 'imsi|rule_id': RuleRecord dictionary
    """
    launch_successful_future = Future()
    try:
        thread = StartThread(test_setup, launch_successful_future)
    except (BadConfigError, ServiceRunningError) as e:
        logging.error("Failed to start test apps in separate thread: %s" % e)
        exit()

    msg = launch_successful_future.result()
    if msg != "Setup successful":
        logging.error("Failed to start test apps in separate thread: %s" % msg)
        exit()

    return thread


def stop_ryu_app_thread(thread):
    """
    Stops the ryu thread, blocks until finished

    Args:
        test_setup (EnforcementStatsController): test_setup config
    """
    thread.keep_running = False
    while not thread.done:
        hub.sleep(1)


def pkt_total(stats):
    """ Given list of FlowQuery FlowStats tuples return total packet sum """
    return sum(n.packets for n in stats)


def wait_after_send(test_controller, wait_time=1, max_sleep_time=20):
    """
    Wait after sending packets, waits until no packets were received by any
    table in ovs. This will send 2 or more table stat requests to ovs.

    Args:
        test_controller (TestingController): testing controller reference
        wait_time (int): wait time between ovs stat queries
        max_sleep_time (int): max wait time, if exceeded return

    Returns when waiting is done

    Throws a WaitTimeExceeded Exception if max_sleep_time exceeded
    """
    sleep_time = 0
    pkt_cnt_old = -1
    while True:
        hub.sleep(wait_time)

        pkt_cnt_new = sum(
            table.matched_count for table in
            RyuDirectFlowQuery.get_table_stats(test_controller)
        )
        if pkt_cnt_new - pkt_cnt_old == 0:
            return
        else:
            pkt_cnt_old = pkt_cnt_new

        sleep_time = sleep_time + wait_time
        if (sleep_time >= max_sleep_time):
            raise WaitTimeExceeded(
                "Waiting on pkts exceeded the max({}) sleep time".
                    format(max_sleep_time)
            )


def setup_controller(controller, setup_req, sleep_time: float = 1,
                     retries: int = 5):
    for _ in range(0, retries):
        ret = controller.check_setup_request_epoch(setup_req.epoch)
        if ret == SetupFlowsResult.SUCCESS:
            return ret
        else:
            res = controller.handle_restart(setup_req.requests)
            if res.result == SetupFlowsResult.SUCCESS:
                return SetupFlowsResult.SUCCESS
        hub.sleep(sleep_time)
    return res.result


def fake_inout_setup(inout_controller):
    TestCase().assertEqual(setup_controller(
        inout_controller, SetupPolicyRequest(requests=[], epoch=global_epoch)),
        SetupFlowsResult.SUCCESS)


def fake_controller_setup(enf_controller=None,
                          enf_stats_controller=None,
                          startup_flow_controller=None,
                          check_quota_controller=None,
                          setup_flows_request=None,
                          check_quota_request=None):
    """
    Immitate contoller restart. This is done by manually setting contoller init
    fields back to False, and restarting the startup stats controller(optional)

    If no stats controller is given this means a clean restart, if clean restart
    flag is not set fail the test case.
    """
    if setup_flows_request is None:
        setup_flows_request = SetupPolicyRequest(
            requests=[], epoch=global_epoch,
        )
    if startup_flow_controller:
        startup_flow_controller._flows_received = False
        startup_flow_controller._table_flows.clear()
        hub.spawn(startup_flow_controller._poll_startup_flows, 1)
    else:
        TestCase().assertEqual(enf_controller._clean_restart, True)
        if enf_stats_controller:
            TestCase().assertEqual(enf_stats_controller._clean_restart, True)
    enf_controller.init_finished = False
    TestCase().assertEqual(setup_controller(
        enf_controller, setup_flows_request),
        SetupFlowsResult.SUCCESS)
    if enf_stats_controller:
        enf_stats_controller.init_finished = False
        enf_stats_controller.cleanup_state()
        TestCase().assertEqual(setup_controller(
            enf_stats_controller, setup_flows_request),
            SetupFlowsResult.SUCCESS)
    if check_quota_controller:
        check_quota_controller.init_finished = False
        if check_quota_request is None:
            check_quota_request = UpdateSubscriberQuotaStateRequest(
                requests=[], epoch=global_epoch,
            )
        TestCase().assertEqual(setup_controller(
            check_quota_controller, check_quota_request),
            SetupFlowsResult.SUCCESS)


def fake_cwf_setup(ue_mac_controller, setup_ue_mac_request=None):
    if setup_ue_mac_request is None:
        setup_ue_mac_request = SetupUEMacRequest(
            requests=[], epoch=global_epoch,
        )
    ue_mac_controller.init_finished = False
    TestCase().assertEqual(setup_controller(
        ue_mac_controller, setup_ue_mac_request),
        SetupFlowsResult.SUCCESS)


def wait_for_enforcement_stats(controller, rule_list, wait_time=1,
                               max_sleep_time=25):
    """
    Wait until all rules from rule_list appear in reports from
    enforcement_controller to sessiond. This is done by checking the mocked
    EnforcementStatsController '_report_usage' method call arguments.

    Args:
        controller (EnforcementStatsController): EnfStats controller reference
        rule_list ([string]): list of strings such as 'IMSI0088888|nice_rule'
        wait_time (int): wait time between checking _report_usage calls
        max_sleep_time (int): max wait time, if exceeded return

    Returns when waiting is done or max_sleep_time exceeded

    Throws a WaitTimeExceeded Exception if max_sleep_time exceeded
    """
    sleep_time = 0
    stats_reported = {rule: False for rule in rule_list}
    while not all(stats_reported[rule] for rule in rule_list):
        hub.sleep(wait_time)
        for reported_stats in controller._report_usage.call_args_list:
            stats = reported_stats[0][0]
            for rule in rule_list:
                if rule in stats:
                    stats_reported[rule] = True

        sleep_time = sleep_time + wait_time
        if (sleep_time >= max_sleep_time):
            raise WaitTimeExceeded(
                "Waiting on enforcement stats exceeded the max({}) sleep time".
                    format(max_sleep_time)
            )


def get_enforcement_stats(enforcement_stats):
    """
    Parses multiple _report_usage(delta_usage) from EnforcementStatsController
    delta_usage structs into one dict with the max byte totals.

    Args:
        controller (EnforcementStatsController): EnfStats controller reference
    Return:
        stats {rule:RuleRecord}: 'imsi|rule_id': RuleRecord dictionary
    """
    stats = {}
    for reported_stats in enforcement_stats:
        for (rule, info) in reported_stats[0][0].items():
            if rule not in stats:
                stats[rule] = info
            else:
                stats[rule].bytes_rx = max(stats[rule].bytes_rx, info.bytes_rx)
                stats[rule].bytes_tx = max(stats[rule].bytes_tx, info.bytes_tx)
    return stats


def create_service_manager(services: List[int],
                           static_services: List[str] = None):
    """
    Creates a service manager from the given list of services.
    Args:
        services ([int]): Enums of the service from mconfig proto
    Returns:
        A service manager instance from the given config
    """
    mconfig = PipelineD(services=services)
    magma_service = MagicMock()
    magma_service.mconfig = mconfig
    if static_services is None:
        static_services = []
    magma_service.config = {
        'static_services': static_services,
        '5G_feature_set': {'enable': False}
    }
    # mock the get_default_client function used to return a fakeredis object
    func_mock = MagicMock(return_value=fakeredis.FakeStrictRedis())
    with mock.patch(
            'magma.pipelined.rule_mappers.get_default_client',
            func_mock):
        service_manager = ServiceManager(magma_service)

    # Workaround as we don't use redis in unit tests
    service_manager.rule_id_mapper._rule_nums_by_rule = {}
    service_manager.rule_id_mapper._rules_by_rule_num = {}
    service_manager.session_rule_version_mapper._version_by_imsi_and_rule = {}
    service_manager.interface_to_prefix_mapper._prefix_by_interface = {}
    service_manager.tunnel_id_mapper._tunnel_map = {}

    return service_manager


def _parse_flow(flow):
    fields_to_remove = [
        r'duration=[\d\w\.]*, ',
        r'idle_age=[\d]*, ',
    ]
    for field in fields_to_remove:
        flow = re.sub(field, '', flow)
    return flow


def _get_current_bridge_snapshot(bridge_name, service_manager,
                                 include_stats=True) -> List[str]:
    table_assignments = service_manager.get_all_table_assignments()
    # Currently, the unit test setup library does not set up the ryu api app.
    # For now, snapshots are created from the flow dump output using ovs and
    # parsed using regex. Once the ryu api works for unit tests, we can
    # directly parse the api response and avoid the regex.
    flows = BridgeTools.get_annotated_flows_for_bridge(bridge_name,
                                                       table_assignments,
                                                       include_stats=include_stats)
    return [_parse_flow(flow) for flow in flows]


def fail(test_case: TestCase, err_msg: str, _bridge_name: str,
         snapshot_file, current_snapshot):
    ofctl_cmd = "sudo ovs-ofctl dump-flows %s" % _bridge_name
    p = subprocess.Popen([ofctl_cmd],
                         stdout=subprocess.PIPE,
                         shell=True)
    ofctl_dump = p.stdout.read().decode("utf-8").strip()
    logging.error("cmd ofctl_dump: %s", ofctl_dump)

    msg = 'Snapshot mismatch with error:\n' \
          '{}\n' \
          'To fix the error, update "{}" to the current snapshot:\n' \
          '{}'.format(err_msg, snapshot_file,
                      '\n'.join(current_snapshot))

    # For fixing test cases after major changes
    #file = open(snapshot_file,'w+')
    #file.write("\n".join(current_snapshot))
    #file.close()

    return test_case.fail(msg)


def expected_snapshot(test_case: TestCase,
                      bridge_name: str,
                      current_snapshot,
                      snapshot_name: Optional[str] = None) -> bool:
    if snapshot_name is not None:
        combined_name = '{}.{}{}'.format(test_case.id(), snapshot_name,
                                         SNAPSHOT_EXTENSION)
    else:
        combined_name = '{}{}'.format(test_case.id(), SNAPSHOT_EXTENSION)
    snapshot_file = os.path.join(
        os.path.dirname(os.path.realpath(__file__)),
        SNAPSHOT_DIR,
        combined_name)

    try:
        with open(snapshot_file, 'r') as file:
            prev_snapshot = []
            for line in file:
                prev_snapshot.append(line.rstrip('\n'))
    except OSError as e:
        fail(test_case, str(e), bridge_name, combined_name, current_snapshot)

    return snapshot_file, prev_snapshot


def assert_bridge_snapshot_match(test_case: TestCase, bridge_name: str,
                                 service_manager: ServiceManager,
                                 snapshot_name: Optional[str] = None,
                                 include_stats: bool = True):
    """
    Verifies the current bridge snapshot matches the snapshot saved in file for
    the given test case. Fails the test case if the snapshots differ.

    Args:
        test_case: Test case instance of the current test
        bridge_name: Name of the bridge
        service_manager: Service manager instance used to obtain the app to
            table number mapping
        snapshot_name: Name of the snapshot. For tests with multiple snapshots,
            this is used to distinguish the snapshots
    """

    current_snapshot = _get_current_bridge_snapshot(bridge_name,
                                                    service_manager,
                                                    include_stats)

    snapshot_file, expected = expected_snapshot(test_case,
                                                bridge_name,
                                                current_snapshot,
                                                snapshot_name)
    if set(current_snapshot) != set(expected):
        fail(test_case,
             '\n'.join(list(unified_diff(expected, current_snapshot,
                                         fromfile='previous snapshot',
                                         tofile='current snapshot'))),
             bridge_name,
             snapshot_file,
             current_snapshot)


def wait_for_snapshots(test_case: TestCase,
                       bridge_name: str,
                       service_manager: ServiceManager,
                       snapshot_name: Optional[str] = None,
                       wait_time: int = 1, max_sleep_time: int = 20,
                       datapath=None,
                       try_snapshot=False):
    """
    Wait after checking ovs snapshot as new changes might still come in,

    Args:
        wait_time (int): wait time between ovs stat queries
        max_sleep_time (int): max wait time, if exceeded return

    Returns when waiting is done

    Throws a WaitTimeExceeded Exception if max_sleep_time exceeded
    """
    sleep_time = 0
    old_snapshot = _get_current_bridge_snapshot(bridge_name, service_manager)
    while True:
        if datapath:
            flows.set_barrier(datapath)
        hub.sleep(wait_time)

        new_snapshot = _get_current_bridge_snapshot(bridge_name, service_manager)
        if try_snapshot:
            snapshot_file, expected_ = expected_snapshot(test_case,
                                                         bridge_name,
                                                         snapshot_name)
            if new_snapshot == expected_:
                return
        else:
            if new_snapshot == old_snapshot:
                return
            else:
                old_snapshot = new_snapshot

        sleep_time = sleep_time + wait_time
        if sleep_time >= max_sleep_time:
            raise WaitTimeExceeded(
                "Waiting on pkts exceeded the max({}) sleep time".
                    format(max_sleep_time)
            )


class SnapshotVerifier:
    """
    SnapshotVerifier is a context wrapper for verifying bridge snapshots.
    """

    def __init__(self, test_case: TestCase, bridge_name: str,
                 service_manager: ServiceManager,
                 snapshot_name: Optional[str] = None,
                 include_stats: bool = True,
                 max_sleep_time: int = 20,
                 datapath=None,
                 try_snapshot=False):
        """
        These arguments are used to call assert_bridge_snapshot_match on exit.

        Args:
            test_case: Test case instance of the current test
            bridge_name: Name of the bridge
            service_manager: Service manager instance used to obtain the app to
                table number mapping
            snapshot_name: Name of the snapshot. For tests with multiple snapshots,
                this is used to distinguish the snapshots
        """
        self._test_case = test_case
        self._bridge_name = bridge_name
        self._service_manager = service_manager
        self._snapshot_name = snapshot_name
        self._include_stats = include_stats
        self._max_sleep_time = max_sleep_time
        self._datapath = datapath
        self._try_snapshot = try_snapshot

    def __enter__(self):
        pass

    def __exit__(self, type, value, traceback):
        """
        Runs after finishing 'with' (Verify snapshot)
        """
        try:
            wait_for_snapshots(self._test_case,
                               self._bridge_name,
                               self._service_manager,
                               self._snapshot_name,
                               max_sleep_time=self._max_sleep_time,
                               datapath=self._datapath,
                               try_snapshot=self._try_snapshot)
        except WaitTimeExceeded as e:
            ofctl_cmd = "sudo ovs-ofctl dump-flows %s".format(self._bridge_name)
            p = subprocess.Popen([ofctl_cmd],
                                 stdout=subprocess.PIPE,
                                 shell=True)
            ofctl_dump = p.stdout.read().decode("utf-8").strip()
            logging.error("ofctl_dump: [%s]", ofctl_dump)
            TestCase().fail(e)

        assert_bridge_snapshot_match(self._test_case, self._bridge_name,
                                     self._service_manager,
                                     self._snapshot_name, self._include_stats)


def get_ovsdb_port_tag(port_name: str) -> str:
    dump1 = subprocess.Popen(["ovsdb-client", "dump", "Port", "name", "tag"],
                             stdout=subprocess.PIPE)
    for port in dump1.stdout.readlines():
        if port_name not in str(port):
            continue
        try:
            tokens = str(port.decode("utf-8")).strip('\"').split()
            return tokens[1]
        except ValueError:
            pass


def get_iface_ipv4(iface: str) -> List[str]:
    virt_ifaddresses = netifaces.ifaddresses(iface)
    ip_addr_list = []
    for ip_rec in virt_ifaddresses[netifaces.AF_INET]:
        ip_addr_list.append(ip_rec['addr'])

    return ip_addr_list


def get_iface_gw_ipv4(iface: str) -> List[str]:
    gateways = netifaces.gateways()
    gateway_ip_addr_list = []
    for gw_ip, gw_iface, _ in gateways[netifaces.AF_INET]:
        if gw_iface != iface:
            continue
        gateway_ip_addr_list.append(gw_ip)

    return gateway_ip_addr_list
