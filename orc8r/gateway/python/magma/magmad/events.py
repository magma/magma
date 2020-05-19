"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import json
import snowflake
from asyncio.events import AbstractEventLoop

from google.protobuf.json_format import MessageToDict
from magma.eventd.eventd_client import log_event
from orc8r.protos.eventd_pb2 import Event
from orc8r.swagger.models.processed_updates import ProcessedUpdates
from orc8r.swagger.models.restarted_services import RestartedServices


def processed_updates(loop: AbstractEventLoop, updates):
    # Convert updates to dicts for JSON serializability
    dict_updates = [MessageToDict(u) for u in updates]
    loop.run_until_complete(
        log_event(
            Event(
                stream_name="magmad",
                event_type="processed_updates",
                tag=snowflake.snowflake(),
                value=json.dumps(ProcessedUpdates(updates=dict_updates).to_dict()),
            )
        )
    )


def restarted_services(loop: AbstractEventLoop, services):
    # Convert to a list for JSON serializability
    services = list(services)
    loop.run_until_complete(
        log_event(
            Event(
                stream_name="magmad",
                event_type="restarted_services",
                tag=snowflake.snowflake(),
                value=json.dumps(RestartedServices(services=services).to_dict()),
            )
        )
    )
