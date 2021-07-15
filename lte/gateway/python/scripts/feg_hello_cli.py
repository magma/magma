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

from feg.protos.hello_pb2 import HelloRequest
from feg.protos.hello_pb2_grpc import HelloStub
from magma.common.rpc_utils import cloud_grpc_wrapper


@cloud_grpc_wrapper
def echo(client, args):
    req = HelloRequest(greeting=args.msg, grpc_err_code=args.err_code)
    print("request:\n", req)
    resp = client.SayHello(req)
    print("resp: ", resp)


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


def main():
    parser = create_parser()

    # Parse the args
    args = parser.parse_args()

    # Execute the subcommand function
    args.func(args, HelloStub, 'feg_hello')


if __name__ == "__main__":
    main()
