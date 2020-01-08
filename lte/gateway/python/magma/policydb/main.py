"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging

from lte.protos.mconfig import mconfigs_pb2
from lte.protos.subscriberdb_pb2_grpc import SubscriberDBStub
from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.common.streamer import StreamerClient
from magma.policydb.rpc_servicer import SessionRpcServicer

from .streamer_callback import PolicyDBStreamerCallback, \
    RuleMappingsStreamerCallback, BaseNamesStreamerCallback
from .basename_store import BaseNameDict


def main():
    """ main() for subscriberdb """
    service = MagmaService('policydb', mconfigs_pb2.PolicyDB())

    # Add all servicers to the server
    chan = ServiceRegistry.get_rpc_channel('subscriberdb',
                                           ServiceRegistry.LOCAL)
    subscriberdb_stub = SubscriberDBStub(chan)
    session_servicer = SessionRpcServicer(service.mconfig, subscriberdb_stub)
    session_servicer.add_to_server(service.rpc_server)

    # Start a background thread to stream updates from the cloud
    if service.config['enable_streaming']:
        stream = StreamerClient(
            {
                'policydb': PolicyDBStreamerCallback(),
                'rule_mappings': RuleMappingsStreamerCallback(),
                'base_names': BaseNamesStreamerCallback(BaseNameDict()),
            },
            service.loop,
        )
        stream.start()
    else:
        logging.info('enable_streaming set to False. Streamer disabled!')

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
