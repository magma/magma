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
import abc
import logging

import grpc
from integ_tests.gateway.rpc import get_rpc_channel
from orc8r.protos.magmad_pb2 import RestartServicesRequest
from orc8r.protos.magmad_pb2_grpc import MagmadStub


class MagmadServiceClient(metaclass=abc.ABCMeta):
    """ Interface for Magmad client """

    @abc.abstractmethod
    def restart_services(self, services):
        """
        Restart magmad services.

        Args:
            services: List of str of services names

        """
        raise NotImplementedError()


class MagmadServiceGrpc(MagmadServiceClient):
    """
    Handle magmad actions by making service calls over gRPC.
    """

    def __init__(self):
        self._magmad_stub = MagmadStub(get_rpc_channel("magmad"))

    def restart_services(self, services):
        request = RestartServicesRequest()
        for s in services:
            request.services.append(s)
        try:
            self._magmad_stub.RestartServices(request)
        except grpc.RpcError as error:
            err_code = error.exception().code()
            if err_code == grpc.StatusCode.FAILED_PRECONDITION:
                logging.info("Ignoring FAILED_PRECONDITION exception")
            else:
                raise
