"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
from typing import Any

from lte.protos.policydb_pb2 import PolicyRule

from magma.common.streamer import StreamerClient
from .rule_store import PolicyRuleDict


class PolicyDBStreamerCallback(StreamerClient.Callback):
    """
    Callback implementation for the PolicyDB StreamerClient instance.
    """

    def __init__(self, loop):
        self._policy_dict = PolicyRuleDict()
        self._loop = loop

    def get_request_args(self, stream_name: str) -> Any:
        return None

    def process_update(self, stream_name, updates, resync):
        logging.info("Processing %d policy updates (resync=%s)",
                     len(updates), resync)
        if resync:
            policy_ids = set()
            for update in updates:
                policy = PolicyRule()
                policy.ParseFromString(update.value)
                self._store_policy_rule(policy)
                policy_ids.add(policy.id)
            logging.debug("Resync with policies: %s", ','.join(policy_ids))
            self._remove_old_policies(policy_ids)
            self._policy_dict.send_update_notification()
        else:
            pass

    def _store_policy_rule(self, policy):
        self._policy_dict[policy.id] = policy

    def _remove_old_policies(self, id_set):
        """
        Scan the set of ids passes in the streaming update to see which have
        been deleted and delete them in the policy dictionary
        """
        missing_rules = set(self._policy_dict.keys()) - id_set
        for rule in missing_rules:
            del self._policy_dict[rule]
