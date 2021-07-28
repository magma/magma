#!/usr/bin/env python3

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

import argparse
import textwrap

from magma.common.rpc_utils import grpc_wrapper
from orc8r.protos import common_pb2
from orc8r.protos.ctraced_pb2 import EndTraceRequest, StartTraceRequest
from orc8r.protos.ctraced_pb2_grpc import CallTraceServiceStub


@grpc_wrapper
def start_call_trace(client, args):
    client.StartCallTrace(StartTraceRequest())


@grpc_wrapper
def end_call_trace(client, args):
    res = client.EndCallTrace(common_pb2.Void())
    print("Result of call trace: ", res)


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        formatter_class=argparse.RawTextHelpFormatter,
        description=textwrap.dedent('''\
            Management CLI for ctraced
            --------------------------
            Use to start and end call traces.
            Options are provided for the type of trace to record.
            Only a single trace can be captured at a time.
        '''),
    )

    # Add subcommands
    subparsers = parser.add_subparsers(title='subcommands', dest='cmd')

    # Add StartCallTrace subcommand
    parser_start_trace = subparsers.add_parser(
        'start_call_trace',
        help='Start a call trace',
    )
    trace_types = list(StartTraceRequest.TraceType.DESCRIPTOR.values_by_name)
    supported_protocols =\
        list(StartTraceRequest.ProtocolName.DESCRIPTOR.values_by_name)
    supported_interfaces =\
        list(StartTraceRequest.InterfaceName.DESCRIPTOR.values_by_name)
    parser_start_trace.add_argument(
        '--type', type=str, choices=trace_types,
        help='Trace type', required=True,
    )
    parser_start_trace.add_argument(
        '--imsi', type=str, choices=trace_types,
        help='Trace type',
    )
    parser_start_trace.add_argument(
        '--protocol', type=str,
        choices=supported_protocols,
    )
    parser_start_trace.add_argument(
        '--interface', type=str,
        choices=supported_interfaces,
    )
    parser_start_trace.set_defaults(func=start_call_trace)

    # Add EndCallTrace subcommand
    parser_end_trace = subparsers.add_parser(
        'end_call_trace',
        help='End a call trace',
    )
    parser_end_trace.set_defaults(func=end_call_trace)

    return parser


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()
    if not args.cmd:
        parser.print_usage()
        exit(1)

    # Execute the subcommand function
    args.func(args, CallTraceServiceStub, 'ctraced')


if __name__ == "__main__":
    main()
