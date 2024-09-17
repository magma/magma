"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Util module for executing multiple `ping` commands via subprocess.
"""

import asyncio
import re
from collections import namedtuple

from magma.magmad.check import subprocess_workflow

DEFAULT_NUM_PACKETS = 4
DEFAULT_TIMEOUT_SECS = 20


PingCommandParams = namedtuple(
    'PingCommandParams',
    ['host_or_ip', 'num_packets', 'timeout_secs'],
)

PingInterfaceCommandParams = namedtuple(
    'PingInterfaceCommandParams',
    ['host_or_ip', 'num_packets', 'interface', 'timeout_secs'],
)

PingCommandResult = namedtuple(
    'PingCommandResult',
    ['error', 'host_or_ip', 'num_packets', 'stats'],
)

ParsedPingStats = namedtuple(
    'ParsedPingStats', [
        'packets_transmitted',
        'packets_received',
        'packet_loss_pct',
        'rtt_min',
        'rtt_avg',
        'rtt_max',
        'rtt_mdev',
    ],
)

# regexp's for parsing
dec_re = r'\d+(\.\d+)?'
packet_line_re = re.compile(
    r'^(?P<packets_transmitted>\d+) packets transmitted, '
    + r'(?P<packets_received>\d+) received, '
    + r'(?P<packet_loss_pct>{d})% packet loss, '.format(d=dec_re)
    + r'time .+$',
)
rtt_line_re = re.compile(
    r'^rtt min/avg/max/mdev = '
    + r'(?P<rtt_min>{d})/(?P<rtt_avg>{d})/'.format(d=dec_re)
    + r'(?P<rtt_max>{d})/(?P<rtt_mdev>{d}) ms$'.format(d=dec_re),
)


def ping(ping_params):
    """
    Execute ping commands via subprocess. Blocks while waiting for output.

    Args:
        ping_params ([PingCommandParams]): params for the pings to execute
    Returns:
        [PingCommandResult]: stats from the executed ping commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses(
        ping_params,
        _get_ping_command_args_list,
        parse_ping_output,
    )


@asyncio.coroutine
def ping_async(ping_params, loop=None):
    """
    Execute ping commands asynchronously.

    Args:
        ping_params ([PingCommandParams]): params for the pings to execute
        loop: asyncio event loop (optional)

    Returns:
        [PingCommandResult]: stats from the executed ping commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses_async(
        ping_params,
        _get_ping_command_args_list,
        parse_ping_output,
        loop,
    )


@asyncio.coroutine
def ping_interface_async(ping_params, loop=None):
    """
    Execute ping commands asynchronously through specified interface.

    Args:
        ping_params ([PingCommandParams]): params for the pings to execute
        loop: asyncio event loop (optional)

    Returns:
        [PingCommandResult]: stats from the executed ping commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses_async(
        ping_params,
        _get_ping_command_interface_args_list,
        parse_ping_output,
        loop,
    )


def _get_ping_command_args_list(ping_param):
    return [
        'ping', ping_param.host_or_ip,
        '-c', str(ping_param.num_packets or DEFAULT_NUM_PACKETS),
        '-w', str(ping_param.timeout_secs or DEFAULT_TIMEOUT_SECS),
    ]


def _get_ping_command_interface_args_list(ping_param):
    return [
        'ping', ping_param.host_or_ip,
        '-c', str(ping_param.num_packets or DEFAULT_NUM_PACKETS),
        '-I', str(ping_param.interface),
        '-w', str(ping_param.timeout_secs or DEFAULT_TIMEOUT_SECS),
    ]


def parse_ping_output(stdout, stderr, param):
    """
    Parse stdout output from a ping command.

    Raises:
        ValueError: If any errors are encountered while parsing ping output.
    """
    def create_error_result(error_msg):
        return PingCommandResult(
            error=error_msg,
            host_or_ip=param.host_or_ip,
            num_packets=param.num_packets or DEFAULT_NUM_PACKETS,
            stats=None,
        )

    def find_statistic_line_idx(ping_lines):
        line_re = re.compile('^--- .+ statistics ---$')
        for i, line in enumerate(ping_lines):
            if line_re.match(line):
                return i
        raise ValueError('Could not find statistics header in ping output')

    def match_ping_line(line, line_re, line_name='ping'):
        line_match = line_re.match(line)
        if not line_match:
            raise ValueError(
                'Could not parse {name} line:\n{line}'.format(
                    name=line_name,
                    line=line,
                ),
            )
        return line_match

    def str_to_num(s_in):
        try:
            return int(s_in)
        except ValueError:
            return float(s_in)

    if stderr:
        return create_error_result(stderr)
    else:
        try:
            stdout_lines = stdout.decode('ascii').strip().split('\n')
            stat_header_line_idx = find_statistic_line_idx(stdout_lines)
            if len(stdout_lines) <= stat_header_line_idx + 2:
                raise ValueError(
                    'Not enough output lines in ping output. '
                    'The ping may have timed out.',
                )

            packet_match = match_ping_line(
                stdout_lines[stat_header_line_idx + 1],
                packet_line_re,
                line_name='packet',
            )
            rtt_match = match_ping_line(
                stdout_lines[stat_header_line_idx + 2],
                rtt_line_re,
                line_name='rtt',
            )

            match_dict = {}
            match_dict.update(packet_match.groupdict())
            match_dict.update(rtt_match.groupdict())
            match_dict = {k: str_to_num(v) for k, v in match_dict.items()}

            return PingCommandResult(
                error=None,
                host_or_ip=param.host_or_ip,
                num_packets=param.num_packets or DEFAULT_NUM_PACKETS,
                stats=ParsedPingStats(**match_dict),
            )
        except ValueError as e:
            return create_error_result(str(e.args[0]))
