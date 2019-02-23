"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import logging

import snowflake
from orc8r.protos.logging_service_pb2 import LogEntry
from systemd import journal

from magma.common.scribe_client import ScribeClient

SERVICE_EXIT_CATEGORY = "perfpipe_magma_gateway_service_exit"


class JournalEntryParser(object):
    """
    Utility class for parsing journalctl entries into log events for scribe
    """

    @staticmethod
    def entry_to_log_event(entry):
        time = entry['_SOURCE_REALTIME_TIMESTAMP'].timestamp()
        hw_id = "" if snowflake.snowflake() is None else snowflake.snowflake()
        int_map = {'exit_status': entry['EXIT_STATUS']}
        normal_map = {'unit': entry['UNIT'],
                      'exit_code': entry["EXIT_CODE"]}
        return LogEntry(category=SERVICE_EXIT_CATEGORY,
                        time=int(time),
                        hw_id=hw_id,
                        normal_map=normal_map,
                        int_map=int_map)


@asyncio.coroutine
def start_systemd_tailer(magmad_service_config):
    """
    Tail systemd logs for exit statuses and codes, then report them to
    scribe

    Args:
        magmad_service (magma.common.service.MagmaService):
            MagmaService instance for magmad
    """
    loop = asyncio.get_event_loop()
    scribe_client = ScribeClient(loop=loop)
    poll_interval = magmad_service_config.get('systemd_tailer_poll_interval',
                                              10)

    reader = journal.Reader()
    reader.log_level(journal.LOG_INFO)
    # Only include entries since the current box has booted.
    reader.this_boot()
    reader.this_machine()
    reader.add_match(
        SYSLOG_IDENTIFIER=u'systemd',
        CODE_FUNCTION=u'service_sigchld_event'
    )
    # Move to the end of the journal
    reader.seek_tail()
    # Discard old journal entries
    reader.get_previous()
    while True:
        if reader.wait() == journal.APPEND:
            logging.debug("Found systemd exit error codes, reporting to scribe")
            log_events = [JournalEntryParser.entry_to_log_event(e)
                          for e in reader]
            scribe_client.log_to_scribe_with_sampling_rate(log_events)
        yield from asyncio.sleep(poll_interval, loop=loop)
