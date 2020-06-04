"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import snowflake
from magma.eventd.eventd_client import log_event
from orc8r.protos.eventd_pb2 import Event


def deleted_stored_mconfig():
    log_event(
        Event(
            stream_name="magmad",
            event_type="deleted_stored_mconfig",
            tag=snowflake.snowflake(),
            value="{}",
        )
    )


def updated_stored_mconfig():
    log_event(
        Event(
            stream_name="magmad",
            event_type="updated_stored_mconfig",
            tag=snowflake.snowflake(),
            value="{}",
        )
    )
