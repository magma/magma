import logging

import asyncio
from typing import List

import re
from collections import namedtuple

from magma.common.job import Job
from magma.magmad.check import subprocess_workflow

from magma.pipelined.metrics import GTP_PORT_USER_PLANE_DL_BYTES, \
    GTP_PORT_USER_PLANE_UL_BYTES

OVSDBDumpCommandParams = namedtuple('OVSDBCommandParams', ['table', 'columns'])
ParsedInterfaceStats = namedtuple('ParsedInterfaceStats', [
    'Interface',
    'rx_bytes',
    'tx_bytes',
    'remote_ip',
])
OVSDBCommandResult = namedtuple('OVSDBCommandResult', ['out', 'err'])

interface_tx_rx_stats_re = re.compile(
    r'"(?P<Interface>\w+)"(.*)remote_ip="(?P<remote_ip>.*)"(.*)rx_bytes=(?P<rx_bytes>\d).+tx_bytes=(?P<tx_bytes>\d)')


class GTPStatsCollector(Job):
    def __init__(self, polling_interval: int, service_loop):
        self._polling_interval = max(polling_interval, 60)
        super().__init__(interval=self._polling_interval, loop=service_loop)
        self._loop = service_loop

    @asyncio.coroutine
    def _ovsdb_dump_async(self, table: str, columns: List[str]):
        """
        Execute ovsdb-client dump command asynchronously and parse stdout
        results.

        """
        params = [OVSDBDumpCommandParams(table=table, columns=columns)]
        return subprocess_workflow.exec_and_parse_subprocesses_async(
            params,
            _get_ovsdb_dump_params,
            _parse_ovsdb_dump_output,
            self._loop,
        )

    async def _run(self) -> None:
        logging.info("Running GTP stats collector...")
        while True:
            dump_stats_results = await self._ovsdb_dump_async(
                'Interface', ['name', 'statistics', 'options'])
            for r in list(dump_stats_results)[0].out:
                if 'g_' in r.Interface:
                    GTP_PORT_USER_PLANE_DL_BYTES.labels(r.remote_ip).inc(
                        float(r.rx_bytes))
                    GTP_PORT_USER_PLANE_UL_BYTES.labels(r.remote_ip).inc(
                        float(r.tx_bytes))
            await asyncio.sleep(self._polling_interval, self._loop)


def _get_ovsdb_dump_params(params: OVSDBDumpCommandParams) -> List[str]:
    params_list = ['ovsdb-client', 'dump', params.table]
    params_list.extend(params.columns)
    return params_list


def _parse_ovsdb_dump_output(stdout: str, stderr: str,
                             _) -> OVSDBCommandResult:
    """
    Parse stdout output from ovsdb-client dump command.

    Raises:
        ValueError: If any errors are encountered while parsing output.
    """

    def create_error_result(error_msg):
        return OVSDBCommandResult(
            out='',
            err=error_msg,
        )

    def find_header_line_idx(lines):
        line_re = re.compile('^---.+$')
        for i, line in enumerate(lines):
            if line_re.match(line):
                return i
        raise ValueError('Could not find header in ovsdb-client output')

    def match_gtp_lines(lines, line_re):
        line_matches = []
        for line in lines:
            line_match = line_re.match(line)
            if line_match:
                line_matches.append(line_match.groupdict())
        return line_matches

    if stderr:
        return create_error_result(stderr)
    else:
        try:
            stdout_lines = stdout.decode('ascii').split('\n')
            header_line_idx = find_header_line_idx(stdout_lines)
            gtp_matches = match_gtp_lines(
                stdout_lines[header_line_idx + 1:],
                interface_tx_rx_stats_re)

            results = []
            for m in gtp_matches:
                logging.info(m)
                results.append(ParsedInterfaceStats(**m))

            return OVSDBCommandResult(
                out=results,
                err=None,
            )
        except ValueError as e:
            return create_error_result(str(e.args[0]))
