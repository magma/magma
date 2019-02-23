"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from collections import namedtuple
from concurrent.futures import Future
from unittest import TestCase
from unittest.mock import MagicMock

from lte.protos.mconfig.mconfigs_pb2 import PipelineD
from magma.pipelined.service_manager import ServiceManager
from magma.pipelined.tests.app.exceptions import BadConfigError, \
    ServiceRunningError
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


def create_service_manager(services=list):
    """
    Creates a service manager from the given list of services.
    Args:
        services ([int]): Enums of the service from mconfig proto
    Returns:
        A service manager instance from the given config
    """
    mconfig = PipelineD(relay_enabled=True, services=services)
    magma_service = MagicMock()
    magma_service.mconfig = mconfig
    magma_service.config = {
        'static_apps': ['arpd', 'access_control']
    }
    return ServiceManager(magma_service)
