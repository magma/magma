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

from magmaexample.hello import metrics
from protos import hello_pb2, hello_pb2_grpc


class HelloRpcServicer(hello_pb2_grpc.HelloServicer):
    """
    gRPC based server for Hello service
    """

    def __init__(self):
        pass

    def add_to_server(self, server):
        """
        Add the servicer to a gRPC server
        """
        hello_pb2_grpc.add_HelloServicer_to_server(self, server)

    def SayHello(self, request, context):
        """
        Echo the message as the response
        """
        metrics.NUM_REQUESTS.inc()
        return hello_pb2.HelloReply(greeting=request.greeting)
