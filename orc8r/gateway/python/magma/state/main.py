"""
Copyright (c) 2018-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from orc8r.protos.mconfig import mconfigs_pb2
from orc8r.protos.state_pb2_grpc import StateServiceStub
from magma.common.grpc_client_manager import GRPCClientManager
from magma.common.service import MagmaService
from .garbage_collector import GarbageCollector
from .state_replicator import StateReplicator


def main():
    """
    main() for gateway state replication service
    """
    service = MagmaService('state', mconfigs_pb2.State())

    # _grpc_client_manager to manage grpc client recycling
    grpc_client_manager = GRPCClientManager(
        service_name="state",
        service_stub=StateServiceStub,
        max_client_reuse=60,
    )

    # Garbage collector propagates state deletions back to Orchestrator
    garbage_collector = GarbageCollector(service, grpc_client_manager)

    # Start state replication loop
    state_manager = StateReplicator(service, garbage_collector,
                                    grpc_client_manager)
    state_manager.start()

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
