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
from magma.mobilityd.rpc_servicer import MobilityServiceRpcServicer
from lte.protos.mconfig import mconfigs_pb2


def main():
    """ main() for MobilityD """
    service = MagmaService('mobilityd', mconfigs_pb2.MobilityD())

    # Add all servicers to the server
    mobility_service_servicer = MobilityServiceRpcServicer(
        service.mconfig, service.config)
    mobility_service_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
