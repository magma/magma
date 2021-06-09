"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Util module for executing multiple `dpkg` commands via subprocess.
"""

from typing import Any, Dict, List, NamedTuple, Optional

from magma.magmad.check import subprocess_workflow

RouteCommandParams = NamedTuple('RouteCommandParams', [])
Route = NamedTuple(
    'Route',
    [
        ('destination_ip', str), ('gateway_ip', str),
        ('genmask', str), ('network_interface_id', str),
    ],
)

RouteCommandResult = NamedTuple(
    'RouteCommandResult',
    [
        ('error', Optional[str]),
        ('routing_table', List[Dict[str, Any]]),
    ],
)

# TODO: This relies on the SO language being English. Maybe there is a way to
#  get the info another way.


def get_routing_table() -> RouteCommandResult:
    """
    Execute route command via subprocess. Blocks while waiting for output.
    Returns the routing table in the form of a list of routes.
    """
    return list(
        subprocess_workflow.exec_and_parse_subprocesses(
            [RouteCommandParams()],
            _get_route_command_args_list,
            parse_route_output,
        ),
    )[0]


def _get_route_command_args_list(_):
    return ['route', '-n']


def parse_route_output(stdout, stderr, _):
    """
    Parse stdout output from a route command.
    """
    if stderr:
        return RouteCommandResult(error=stderr, routing_table=[])

    stdout_decoded = stdout.decode().strip()
    heading = stdout_decoded.split('\n')[1]
    if heading.split() != [
        'Destination', 'Gateway', 'Genmask', 'Flags',
        'Metric', 'Ref', 'Use', 'Iface',
    ]:
        return RouteCommandResult(
            error='Unexpected heading: %s' % heading,
            routing_table=[],
        )

    # Ignore the title and heading
    lines = stdout_decoded.split('\n')[2:]
    routes = []
    for line in lines:
        attrs = line.split()
        routes.append(
            Route(
                destination_ip=attrs[0],
                gateway_ip=attrs[1],
                genmask=attrs[2],
                network_interface_id=attrs[7],
            )._asdict(),
        )
    return RouteCommandResult(
        error=None,
        routing_table=routes,
    )
