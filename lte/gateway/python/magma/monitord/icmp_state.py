#  Copyright (c) Facebook, Inc. and its affiliates.
#  All rights reserved.
#
#  This source code is licensed under the BSD-style license found in the
#  LICENSE file in the root directory of this source tree.
import json
from typing import Dict, List, NamedTuple

from orc8r.protos.service303_pb2 import State

ICMP_STATE_TYPE = "icmp_monitoring"

ICMPMonitoringResponse = NamedTuple('ICMPMonitoringResponse',
                                    [('last_reported_time', int),
                                     ('latency_ms', float)])


def serialize_subscriber_states(
        sub_table: Dict[str, ICMPMonitoringResponse]) -> List[State]:
    states = []
    for sub, icmp_resp in sub_table.items():
        serialized = json.dumps(icmp_resp._asdict())
        state = State(
            type=ICMP_STATE_TYPE,
            deviceID=sub,
            value=serialized.encode('utf-8')
        )
        states.append(state)
    return states

