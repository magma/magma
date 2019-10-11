"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
import time
from collections import namedtuple
from concurrent.futures import Future

import abc
import grpc
from lte.protos.pipelined_pb2 import ActivateFlowsRequest, \
    DeactivateFlowsRequest
from ryu.lib import hub

from magma.subscriberdb.sid import SIDUtils

SubContextConfig = namedtuple('ContextConfig', ['imsi', 'ip', 'table_id'])


def try_grpc_call_with_retries(grpc_call, retry_count=5, retry_interval=1):
    """ Attempt a grpc call and retry if unavailable """
    for i in range(retry_count):
        try:
            return grpc_call()
        except grpc.RpcError as error:
            err_code = error.exception().code()
            # Retry if unavailable
            if (err_code == grpc.StatusCode.UNAVAILABLE):
                logging.warning("Pipelined unavailable, retrying...")
                time.sleep(retry_interval * (2 ** i))
                continue
            logging.error("Pipelined grpc call failed with error : %s",
                          error)
            raise


class SubscriberContext(abc.ABC):
    """
    Interface for SubscriberContext

    SubscriberContext handles adding new subscribers to pipelined:
        - Stores all subscriber information
        - Communicates with the Enforcement Table to activate flows
    """

    @abc.abstractmethod
    def add_static_rule(self, id):
        """
        Adds a new static rule to the subscriber
        Args:
            id (String): PolicyRule id
        Returns:
            Self
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def add_dynamic_rule(self, policy_rule):
        """
        Adds new dynamic rule to subcriber
        Args:
            policy_rule (PolicyRule): PolicyRule value
        Returns:
            Self
        """
        raise NotImplementedError()

    def __enter__(self):
        """
        Used for running 'with'
        """
        self._activate_subscriber_rules()

    def __exit__(self, type, value, traceback):
        """
        Clean up after using 'with'
        """
        self._deactivate_subscriber_rules()

    @abc.abstractmethod
    def _activate_subscriber_rules(self):
        """
        Activates all subscriber rules

        Adds flows for subscriber rules to the Enforcement table
        """
        raise NotImplementedError()

    @abc.abstractmethod
    def _deactivate_subscriber_rules(self):
        """
        Deactivate the rules that were added

        Removes subscriber flows from Enforcement table
        """
        raise NotImplementedError()


class RyuRPCSubscriberContext(SubscriberContext):
    """
    RyuRestSubscriberContext uses grpc calls to enforcement_controller for
    testing subscriber rules
    """

    def __init__(self, imsi, ip, pipelined_stub, table_id=5):
        self.cfg = SubContextConfig(imsi, ip, table_id)
        self._dynamic_rules = []
        self._static_rule_names = []
        self._pipelined_stub = pipelined_stub

    def add_static_rule(self, id):
        self._static_rule_names.append(id)
        return self

    def add_dynamic_rule(self, policy_rule):
        self._dynamic_rules.append(policy_rule)
        return self

    def _activate_subscriber_rules(self):
        try_grpc_call_with_retries(
            lambda: self._pipelined_stub.ActivateFlows(
                ActivateFlowsRequest(sid=SIDUtils.to_pb(self.cfg.imsi),
                                     rule_ids=self._static_rule_names,
                                     dynamic_rules=self._dynamic_rules))
        )

    def _deactivate_subscriber_rules(self):
        try_grpc_call_with_retries(
            lambda: self._pipelined_stub.DeactivateFlows(
                DeactivateFlowsRequest(sid=SIDUtils.to_pb(self.cfg.imsi)))
        )


class RyuDirectSubscriberContext(SubscriberContext):
    """
    RyuDirectSubscriberContext uses ryu.hub and enforcement_controller to
    directly manage subscriber flows
    """

    def __init__(self, imsi, ip, enforcement_controller, table_id=5,
                 enforcement_stats_controller=None, nuke_flows_on_exit=True):
        self.cfg = SubContextConfig(imsi, ip, table_id)
        self._dynamic_rules = []
        self._static_rule_names = []
        self._ec = enforcement_controller
        self._esc = enforcement_stats_controller
        self._nuke_flows_on_exit = nuke_flows_on_exit

    def add_static_rule(self, id):
        self._static_rule_names.append(id)
        return self

    def add_dynamic_rule(self, policy_rule):
        self._dynamic_rules.append(policy_rule)
        return self

    def _activate_subscriber_rules(self):
        def activate_flows():
            self._ec.activate_rules(imsi=self.cfg.imsi,
                                    ip_addr=self.cfg.ip,
                                    static_rule_ids=self._static_rule_names,
                                    dynamic_rules=self._dynamic_rules)
            if self._esc:
                self._esc.activate_rules(
                    imsi=self.cfg.imsi,
                    ip_addr=self.cfg.ip,
                    static_rule_ids=self._static_rule_names,
                    dynamic_rules=self._dynamic_rules)
        hub.joinall([hub.spawn(activate_flows)])

    def _deactivate_subscriber_rules(self):
        if self._nuke_flows_on_exit:
            def deactivate_flows():
                self._ec.deactivate_rules(imsi=self.cfg.imsi, rule_ids=None)
            hub.joinall([hub.spawn(deactivate_flows)])
