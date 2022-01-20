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
import datetime

from feg.protos.hello_pb2 import HelloRequest
from feg.protos.hello_pb2_grpc import HelloStub
from magma.common.rpc_utils import cloud_grpc_wrapper


@cloud_grpc_wrapper
def echo(client, args):
    start_time = datetime.datetime.utcnow()
    req = HelloRequest(greeting=args.msg, grpc_err_code=args.err_code)
    print("- Request:\n", req)
    resp = client.SayHello(req)
    end_time = datetime.datetime.utcnow()
    print(f'- Response: {resp.greeting}')
    print_stats(resp, start_time, end_time)


def print_stats(resp, start_time, end_time):
    times = resp.timestamps
    delta1 = times.agw_to_feg_relay_timestamp.ToDatetime() - start_time
    delta2 = times.feg_timestamp.ToDatetime() - start_time
    delta3 = times.feg_relay_to_agw_timestamp.ToDatetime() - start_time
    delta4 = end_time - start_time

    a = to_ms_string(delta1.total_seconds())
    b = to_ms_string((delta2 - delta1).total_seconds())
    c = to_ms_string((delta3 - delta2).total_seconds(), left=True)
    d = to_ms_string((delta4 - delta3).total_seconds(), left=True)
    total_time = to_ms_string((end_time - start_time).total_seconds())

    print('\n- Stats:')
    print(f'  * Total time: {total_time} ms')
    print(f'  * Approximate path (ms):')
    print(
        f'    ┌─────┐─> {   a} ─>┌───────┐-> {   b} ─>┌─────┐\n'
        f'    │ AGW │            │ Orc8r │            │ Feg │\n'
        f'    └─────┘<─ {d   } <─└───────┘<─ {c   } <─└─────┘\n',
    )


def create_parser():
    """
    Creates the argparse parser with all the arguments.
    """
    parser = argparse.ArgumentParser(
        description='CLI to send echo requests to feg_echo',
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    # usage: ./python/scripts/feg_echo_cli.py echo "hello world" 0
    parser.add_argument('msg', type=str, help='echo message')
    parser.add_argument('err_code', type=int, help='echo err code')

    # Add function callbacks
    parser.set_defaults(func=echo)
    return parser


def to_ms_string(seconds, left=False):
    """
    Convert seconds into a string of milliseconds. If the string is smaller
    than 6 positions it will add leading spaces to complete a total of 6 chars
    Args:
        seconds: time in seconds expressed as float
        left: reverse justification and fills blanks at the right
    Returns:
        String with with at least 6 positions.

    """
    ms = str(int(1000 * seconds))
    if left:
        return ms.ljust(6)
    return ms.rjust(6)


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()

    # Execute the subcommand function
    args.func(args, HelloStub, 'feg_hello')


if __name__ == "__main__":
    main()
