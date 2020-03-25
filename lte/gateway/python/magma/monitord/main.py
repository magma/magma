"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

from lte.protos.mconfig import mconfigs_pb2
from magma.common.service import MagmaService
from magma.monitord.icmp_monitoring import ICMPMonitoring
from magma.monitord.icmp_state import serialize_subscriber_states


def main():
    """ main() for monitord service"""
    service = MagmaService('monitord', mconfigs_pb2.MonitorD())

    # Monitoring thread loop
    icmp_monitor = ICMPMonitoring(service.mconfig.polling_interval,
                                  service.loop)
    icmp_monitor.start()

    # Register a callback function for GetOperationalStates
    service.register_operational_states_callback(
        lambda: serialize_subscriber_states(
            icmp_monitor.get_subscriber_state()))

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
