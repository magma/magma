"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Module for executing `traceroute` commands via subprocess
"""

import asyncio
from collections import namedtuple

from magma.magmad.check import subprocess_workflow

DEFAULT_TTL = 30
DEFAULT_BYTES_PER_PACKET = 60


TracerouteParams = namedtuple(
    'TracerouteParams',
    ['host_or_ip', 'max_hops', 'bytes_per_packet'],
)

TracerouteResult = namedtuple(
    'TracerouteResult',
    ['error', 'host_or_ip', 'stats'],
)

TracerouteStats = namedtuple('TracerouteStats', ['hops'])

TracerouteHop = namedtuple('TracerouteHop', ['idx', 'probes'])

TracerouteProbe = namedtuple(
    'TracerouteProbe',
    ['hostname', 'ip_addr', 'rtt_ms'],
)


def traceroute(params):
    """
    Execute some `traceroute` commands via subprocess.

    Args:
        params ([TracerouteParams]): params for the `traceroute` commands

    Returns:
        [TracerouteResult]: stats from the executed `traceroute` commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses(
        params,
        _get_traceroute_command_args_list,
        parse_traceroute_output,
    )


@asyncio.coroutine
def traceroute_async(params, loop=None):
    """
    Execute some `traceroute` commands asynchronously and return results.
    Args:
        params ([TracerouteParams]): params for the `traceroute` commands
        loop: event loop to run in (optional)

    Returns:
        [TracerouteResult]: stats from the executed `traceroute` commands
    """
    return subprocess_workflow.exec_and_parse_subprocesses_async(
        params,
        _get_traceroute_command_args_list,
        parse_traceroute_output,
        loop=loop,
    )


def _get_traceroute_command_args_list(param):
    return [
        'traceroute',
        '-m', str(param.max_hops or DEFAULT_TTL),
        param.host_or_ip,
        str(param.bytes_per_packet or DEFAULT_BYTES_PER_PACKET),
    ]


def parse_traceroute_output(stdout, stderr, param):
    def create_error_result(error_msg):
        return TracerouteResult(
            error=error_msg,
            host_or_ip=param.host_or_ip,
            stats=None,
        )

    if stderr:
        return create_error_result(stderr)
    else:
        try:
            stats = TracerouteParser().parse(stdout)
            return TracerouteResult(
                error=None,
                host_or_ip=param.host_or_ip,
                stats=stats,
            )
        except ValueError as e:
            msg = 'Error while parsing output. ' \
                  'Original exception message:\n{}'.format(str(e.args[0]))
            return create_error_result(msg)
        except IndexError as e:
            msg = 'Error while parsing output - an incomplete line ' \
                  'was encountered. Original exception message:\n{}' \
                .format(str(e.args[0]))
            return create_error_result(msg)


class TracerouteParser(object):

    HostnameAndIP = namedtuple('HostnameAndIP', ['hostname', 'ip'])

    DEFAULT_ENDPOINT = HostnameAndIP(hostname=None, ip=None)

    def __init__(self):
        self._probe_endpoint = self.DEFAULT_ENDPOINT

    def parse(self, output):
        """
        Raises:
            ValueError, IndexError
        """
        output_lines = output.decode('ascii').strip().split('\n')
        output_lines.pop(0)     # strip header line

        hops = []
        for line in output_lines:
            self._probe_endpoint = self.DEFAULT_ENDPOINT
            hops.append(self._parse_hop(line))
        return TracerouteStats(hops)

    def _parse_hop(self, line):
        hop_split = line.split()
        hop_idx = int(hop_split.pop(0))

        probes = []
        while hop_split:
            probe = self._parse_next_probe(hop_split)
            if probe:
                probes.append(probe)

        return TracerouteHop(idx=hop_idx, probes=probes)

    def _parse_next_probe(self, tokens):
        head_token = tokens.pop(0)
        if head_token == '*':
            return TracerouteProbe(
                hostname=self._probe_endpoint.hostname,
                ip_addr=self._probe_endpoint.ip,
                rtt_ms=0,
            )

        lookahead_token = tokens.pop(0)
        if lookahead_token == 'ms':
            return TracerouteProbe(
                hostname=self._probe_endpoint.hostname,
                ip_addr=self._probe_endpoint.ip,
                rtt_ms=float(head_token),
            )
        else:
            ip_addr = lookahead_token[1:-1]
            self._probe_endpoint = self.HostnameAndIP(
                hostname=head_token,
                ip=ip_addr,
            )
            return None
