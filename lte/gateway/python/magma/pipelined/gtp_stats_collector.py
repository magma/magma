"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import asyncio
import logging
import re
from typing import List, NamedTuple, Optional

from magma.common.job import Job
from magma.magmad.check import subprocess_workflow
from magma.pipelined.metrics import (
    GTP_PORT_USER_PLANE_DL_BYTES,
    GTP_PORT_USER_PLANE_UL_BYTES,
)

OVSDBDumpCommandParams = NamedTuple('OVSDBCommandParams',
                                    [('table', str), ('columns', List[str])])
ParsedInterfaceStats = NamedTuple('ParsedInterfaceStats', [
    ('Interface', str),
    ('rx_bytes', str),
    ('tx_bytes', str),
    ('remote_ip', str),
])
OVSDBCommandResult = NamedTuple('OVSDBCommandResult',
                                [('out', List[ParsedInterfaceStats]),
                                 ('err', Optional[str])])

interface_group = r"(?P<Interface>\w+)"
remote_ip_group = r"(?P<remote_ip>.*)"
rx_bytes_group = r"(?P<rx_bytes>\d+)"
tx_bytes_group = r"(?P<tx_bytes>\d+)"
stats_re_str = r'"{}"(.*)(remote_ip="{}")(.*)?rx_bytes={}.+tx_bytes={}'

interface_tx_rx_stats_re = re.compile(
    stats_re_str.format(interface_group, remote_ip_group, rx_bytes_group,
                        tx_bytes_group))

MIN_OVSDB_DUMP_POLLING_INTERVAL = 60
GTP_IP_INTERFACE_PREFIX = 'g_'
GTP_INTERFACE_PREFIX = 'gtp0'


class GTPStatsCollector(Job):
    def __init__(self, polling_interval: int, service_loop):
        self._polling_interval = max(polling_interval,
                                     MIN_OVSDB_DUMP_POLLING_INTERVAL)
        super().__init__(interval=self._polling_interval, loop=service_loop)
        self._loop = service_loop
        logging.info("Running GTP stats collector...")

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
        dump_stats_results = await self._ovsdb_dump_async(
            'Interface', ['name', 'statistics', 'options'])
        for r in list(dump_stats_results)[0].out:
            if GTP_IP_INTERFACE_PREFIX in r.Interface or \
                    r.Interface == GTP_INTERFACE_PREFIX:
                GTP_PORT_USER_PLANE_DL_BYTES.labels(r.remote_ip).set(
                    float(r.tx_bytes))
                GTP_PORT_USER_PLANE_UL_BYTES.labels(r.remote_ip).set(
                    float(r.rx_bytes))


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

    def match_gtp_lines(lines):
        line_matches = []
        for line in lines:
            line_remote_ip_match = interface_tx_rx_stats_re.match(line)
            if line_remote_ip_match:
                match_dict = line_remote_ip_match.groupdict()
                if not 'remote_ip' in match_dict:
                    match_dict['remote_ip'] = ""
                line_matches.append(match_dict)
        return line_matches

    if stderr:
        return create_error_result(stderr)

    try:
        stdout_lines = stdout.decode('ascii').split('\n')
        header_line_idx = find_header_line_idx(stdout_lines)
        gtp_matches = match_gtp_lines(stdout_lines[header_line_idx + 1:])

        results = []
        for m in gtp_matches:
            results.append(ParsedInterfaceStats(**m))

        return OVSDBCommandResult(
            out=results,
            err=None,
        )
    except ValueError as e:
        return create_error_result(str(e.args[0]))
