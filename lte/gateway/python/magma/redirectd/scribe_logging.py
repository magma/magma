"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import datetime

import snowflake
from orc8r.protos.logging_service_pb2 import LogEntry

from magma.common.scribe_client import ScribeClient


class RedirectScribeLogger:
    LOGGING_CATEGORY = "perfpipe_magma_redirectd_stats"

    def __init__(self, event_loop):
        self._client = ScribeClient(loop=event_loop)

    def log_to_scribe(self, redirect_info):
        self._client.log_to_scribe_with_sampling_rate(
            [self.generate_redirect_log_entry(redirect_info)]
        )

    def generate_redirect_log_entry(self, redirect_info):
        time = int(datetime.datetime.now().timestamp())
        hw_id = snowflake.snowflake()
        int_map = {'server_response': redirect_info.server_response.http_code}
        normal_map = {
            'subscriber_ip': redirect_info.subscriber_ip,
            'redirect_address': redirect_info.server_response.redirect_address
        }

        return LogEntry(category=self.LOGGING_CATEGORY,
                        time=int(time),
                        hw_id=hw_id,
                        normal_map=normal_map,
                        int_map=int_map)
