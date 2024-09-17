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

from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.directoryd.rpc_servicer import GatewayDirectoryServiceRpcServicer
from orc8r.protos.mconfig import mconfigs_pb2


def main():
    """ main() for Directoryd """
    service = MagmaService('directoryd', mconfigs_pb2.DirectoryD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    service_config = service.config

    # Add servicer to the server
    gateway_directory_servicer = GatewayDirectoryServiceRpcServicer(
        service_config.get('print_grpc_payload', False),
    )
    gateway_directory_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
