"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from magma.common.service import MagmaService
from .metrics_collector import MetricsCollector
from orc8r.protos.mconfig import mconfigs_pb2


def main():
    """
    Main co-routine for metricsd
    :return: None
    """
    # Get service config
    service = MagmaService('metricsd', mconfigs_pb2.MetricsD())
    services = service.config['services']
    collect_interval = service.config['collect_interval']
    sync_interval = service.config['sync_interval']
    grpc_timeout = service.config['grpc_timeout']
    queue_length = service.config['queue_length']
    loop = service.loop

    # Create local metrics collector
    collector = MetricsCollector(services, collect_interval, sync_interval,
                                 grpc_timeout, queue_length, loop)

    # Start poll and sync loops
    collector.run()

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
