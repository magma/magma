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
from lte.protos.mconfig import mconfigs_pb2
from magma.common.service import MagmaService
from magma.configuration import load_service_config
from magma.monitord.icmp_monitoring import ICMPMonitoring
from magma.monitord.icmp_state import serialize_subscriber_states
from magma.monitord.cpe_monitoring import CpeMonitoringModule
import logging
from lte.protos.mobilityd_pb2 import IPAddress


def main():
    """ main() for monitord service"""
    manual_ping_targets = {}
    service = MagmaService('monitord', mconfigs_pb2.MonitorD())

    # Monitoring thread loop
    mtr_interface = load_service_config("monitord")["mtr_interface"]

    try:
        targets = load_service_config("monitord")["ping_targets"]
        for target, data in targets.items():
            if "ip" in data:
                ip = IPAddress(version=IPAddress.IPV4, address=str.encode(data["ip"]))
                logging.debug('Adding {}:{}:{} to ping target'.format(target, ip.version, ip.address))
                manual_ping_targets[target] = ip
    except KeyError:
        logging.error("No ping targets configured")

    obj = CpeMonitoringModule()
    obj.set_manually_configured_targets(manual_ping_targets)

    icmp_monitor = ICMPMonitoring(obj, service.mconfig.polling_interval,
                                  service.loop, mtr_interface)
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
