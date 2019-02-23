"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
import datetime
import unittest.mock
from unittest import TestCase
from uuid import UUID

from orc8r.protos.logging_service_pb2 import LogEntry
from systemd import journal

from magma.magmad.logging.systemd_tailer import JournalEntryParser


class SystemdTailerTest(TestCase):
    """
    Tests for the systemd tailer
    """
    MOCK_JOURNAL_ENTRY = {
        'EXIT_STATUS': 1,
        '__MONOTONIC_TIMESTAMP': journal.Monotonic((
            datetime.timedelta(0, 1281, 913613),
            UUID('7ed2f5e9-e027-4ab8-b0bb-5bd62966b421'))),
        '_TRANSPORT': 'journal',
        'CODE_FILE': '../src/core/service.c',
        '_COMM': 'systemd',
        '_CAP_EFFECTIVE': '3fffffffff',
        '_MACHINE_ID': UUID('7914936b-cb62-428b-b94c-6cb7fefd283a'),
        '_BOOT_ID': UUID('7ed2f5e9-e027-4ab8-b0bb-5bd62966b421'),
        'EXIT_CODE': 'exited',
        'UNIT': 'magma@mobilityd.service',
        '_GID': 0,
        '_PID': 1,
        '_SYSTEMD_CGROUP': '/init.scope',
        '_CMDLINE': '/sbin/init',
        'SYSLOG_IDENTIFIER': 'systemd',
        'SYSLOG_FACILITY': 3, '_UID': 0,
        'CODE_FUNCTION': 'service_sigchld_event',
        '_SYSTEMD_UNIT': 'init.scope',
        'MESSAGE':
            'magma@mobilityd.service:'
            'Main process exited, code=exited, status=1/FAILURE',
        '_EXE': '/lib/systemd/systemd',
        '_SOURCE_REALTIME_TIMESTAMP': datetime.datetime(
            2018, 1, 5, 17, 28, 17, 975170),
        'CODE_LINE': 2681,
        '_HOSTNAME': 'magma-dev',
        '_SYSTEMD_SLICE': '-.slice',
        '__REALTIME_TIMESTAMP': datetime.datetime(
            2018, 1, 5, 17, 28, 17, 975200),
        'PRIORITY': 5}
    GATEWAY_ID = "test ID"
    MOCK_LOG_ENTRY = LogEntry(
        category='perfpipe_magma_gateway_service_exit',
        time=int(datetime.datetime(
            2018, 1, 5, 17, 28, 17, 975170).timestamp()),
        hw_id=GATEWAY_ID,
        int_map={'exit_status': 1},
        normal_map={'unit': 'magma@mobilityd.service', 'exit_code': 'exited'})

    def test_journal_entry_parser(self):
        """
        Test that mconfig updates are handled correctly
        """
        with unittest.mock.patch('snowflake.snowflake') as mock_snowflake:
            mock_snowflake.side_effect = lambda: self.GATEWAY_ID
            parsed_entry = JournalEntryParser.entry_to_log_event(
                self.MOCK_JOURNAL_ENTRY)
            self.assertEqual(parsed_entry, self.MOCK_LOG_ENTRY)
