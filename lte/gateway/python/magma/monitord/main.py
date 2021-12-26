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

import logging
from typing import List

from lte.protos.mconfig import mconfigs_pb2
from lte.protos.mobilityd_pb2 import IPAddress
from magma.common.sentry import sentry_init
from magma.common.service import MagmaService
from magma.configuration import load_service_config
from magma.monitord.cpe_monitoring import CpeMonitoringModule
from magma.monitord.icmp_job import ICMPJob
from magma.monitord.icmp_state import serialize_subscriber_states
from orc8r.protos.service303_pb2 import State


def _get_serialized_subscriber_states(
        cpe_monitor: CpeMonitoringModule,
) -> List[State]:
    return serialize_subscriber_states(cpe_monitor.get_subscriber_state())


def main():
    """Start monitord"""
    manual_ping_targets = {}
    service = MagmaService('monitord', mconfigs_pb2.MonitorD())

    # Optionally pipe errors to Sentry
    sentry_init(service_name=service.name, sentry_mconfig=service.shared_mconfig.sentry_config)

    # Monitoring thread loop
    mtr_interface = load_service_config("monitord")["mtr_interface"]

    # Add manual IP targets from yml file
    try:
        targets = load_service_config("monitord")["ping_targets"]
        for target, data in targets.items():
            ip_string = data.get("ip")
            if ip_string:
                ip = IPAddress(
                    version=IPAddress.IPV4,
                    address=str.encode(ip_string),
                )
                logging.debug(
                    'Adding %s:%s:%s to ping target', target, ip.version,
                    ip.address,
                )
                manual_ping_targets[target] = ip
    except KeyError:
        logging.warning("No ping targets configured")

    cpe_monitor = CpeMonitoringModule()
    cpe_monitor.set_manually_configured_targets(manual_ping_targets)

    icmp_monitor = ICMPJob(
        cpe_monitor,
        service.mconfig.polling_interval,
        service.loop, mtr_interface,
    )
    icmp_monitor.start()

    # Register a callback function for GetOperationalStates
    service.register_operational_states_callback(
        lambda: _get_serialized_subscriber_states(cpe_monitor),
    )

    # Run the service loop
    service.run()

    # Cleanup the service
    service.close()


if __name__ == "__main__":
    main()
