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

from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.state.garbage_collector import GarbageCollector
from magma.state.state_replicator import StateReplicator
from orc8r.protos.mconfig import mconfigs_pb2
from orc8r.protos.state_pb2_grpc import StateServiceStub


def main():
    """
    main() for gateway state replication service
    """
    service = MagmaService('state', mconfigs_pb2.State())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    # _grpc_client_manager to manage grpc client recycling
    grpc_client_manager = GRPCClientManager(
        service_name="state",
        service_stub=StateServiceStub,
        max_client_reuse=60,
    )

    config = service.config
    print_grpc_payload = config.get('print_grpc_payload', False)

    # Garbage collector propagates state deletions back to Orchestrator
    garbage_collector = GarbageCollector(
        service, grpc_client_manager, print_grpc_payload,
    )

    # Start state replication loop
    state_manager = StateReplicator(
        service, garbage_collector,
        grpc_client_manager,
        print_grpc_payload,
    )
    state_manager.start()

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
