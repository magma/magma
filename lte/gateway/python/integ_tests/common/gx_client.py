"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from feg.protos.mock_core_pb2_grpc import MockOCSStub, MockPCRFStub
from integ_tests.gateway.rpc import get_ocs_rpc_channel, get_pcrf_rpc_channel
from lte.protos.subscriberdb_pb2 import SubscriberID
from orc8r.protos.common_pb2 import Void


class PCRFGrpcClient(Exception):
    def __init__(self, pcrf_grpc_stub):
        """ Init the gRPC stub.  """
        self._mock_pcrf_stub = pcrf_grpc_stub

    def clear_subscribers(self):
        """ Clearing account in PCRF """
        self._mock_pcrf_stub.ClearSubscribers(Void())

    def add_subscriber(self, imsi):
        imsi = imsi.replace("IMSI", "")
        """ Creating account in PCRF """
        self._mock_pcrf_stub.CreateAccount(SubscriberID(id=imsi))


class PCRFGrpc(PCRFGrpcClient):
    """
    Handle mock PCRF actions by making calls over gRPC directly to the
    pcrf service on the feg.
    """

    def __init__(self):
        """ Init the gRPC stub to connect to mock PCRF. """
        super().__init__(MockPCRFStub(get_pcrf_rpc_channel()))
