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

from feg.protos.hello_pb2 import HelloRequest
from feg.protos.hello_pb2_grpc import HelloStub
from magma.common.rpc_utils import grpc_wrapper


@grpc_wrapper
def say_hello(client, args):
    print(client.SayHello(HelloRequest(greeting='world')))


if __name__ == "__main__":
    args = None
    say_hello(args, HelloStub, 'hello')
