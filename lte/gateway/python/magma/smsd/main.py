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

from magma.common.service import MagmaService
from lte.protos.mconfig import mconfigs_pb2
from lte.protos.sms_orc8r_pb2_grpc import SMSOrc8rServiceStub, SMSOrc8rGatewayServiceStub, SmsDStub
from magma.common.service_registry import ServiceRegistry
from orc8r.protos.directoryd_pb2_grpc import GatewayDirectoryServiceStub
from .relay import SmsRelay


def main():
    """ main() for smsd """
    service = MagmaService('smsd', None)

    directoryd_chan = ServiceRegistry.get_rpc_channel('directoryd',
                                                      ServiceRegistry.LOCAL)
    mme_chan = ServiceRegistry.get_rpc_channel('sms_mme_service',
                                               ServiceRegistry.LOCAL)
    smsd_chan = ServiceRegistry.get_rpc_channel('smsd', ServiceRegistry.CLOUD)

    # Add all servicers to the server
    smsd_relay = SmsRelay(service.loop,
                          GatewayDirectoryServiceStub(directoryd_chan),
                          SMSOrc8rGatewayServiceStub(mme_chan),
                          SmsDStub(smsd_chan))
    smsd_relay.add_to_server(service.rpc_server)
    smsd_relay.start()

    # Run the service loop
    service.run()
    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
