"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
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
