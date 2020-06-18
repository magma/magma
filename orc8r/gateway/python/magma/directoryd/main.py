"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.common.service import MagmaService
from magma.directoryd.rpc_servicer import GatewayDirectoryServiceRpcServicer
from orc8r.protos.mconfig import mconfigs_pb2


def main():
    """ main() for Directoryd """
    service = MagmaService('directoryd', mconfigs_pb2.DirectoryD())

    # Add servicer to the server
    gateway_directory_servicer = GatewayDirectoryServiceRpcServicer()
    gateway_directory_servicer.add_to_server(service.rpc_server)

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
