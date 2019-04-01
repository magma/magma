"""
Copyright (c) 2019-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.

Util module for executing multiple `dpkg` commands via subprocess.
"""

from typing import List, NamedTuple, Optional

from magma.magmad.check import subprocess_workflow

RouteCommandParams = NamedTuple('RouteCommandParams', [])
Route = NamedTuple('Route',
                   [('destination', str), ('gateway', str), ('genmask', str),
                    ('flags', str), ('metric', str), ('ref', str),
                    ('use', str), ('interface', str)])
RouteCommandResult = NamedTuple('RouteCommandResult',
                                [('error', Optional[str]),
                                 ('routing_table', Optional[List[Route]])])


def get_routing_table() -> RouteCommandResult:
    """
    Execute route command via subprocess. Blocks while waiting for output.
    Returns the routing table in the form of a list of routes.
    """
    return list(subprocess_workflow.exec_and_parse_subprocesses(
        [RouteCommandParams()],
        _get_route_command_args_list,
        parse_route_output,
    ))[0]


def _get_route_command_args_list(_):
    return ['route', '-n']


def parse_route_output(stdout, stderr, _):
    """
    Parse stdout output from a route command.
    """
    if stderr:
        return RouteCommandResult(error=stderr, routing_table=None)

    stdout_decoded = stdout.decode().strip()
    heading = stdout_decoded.split('\n')[1]
    if heading.split() != ['Destination', 'Gateway', 'Genmask', 'Flags',
                           'Metric', 'Ref', 'Use', 'Iface']:
        return RouteCommandResult(error='Unexpected heading: %s' % heading,
                                  routing_table=None)

    # Ignore the title and heading
    lines = stdout_decoded.split('\n')[2:]
    return RouteCommandResult(
        error=None,
        routing_table=[Route(*line.split()) for line in lines]
    )
