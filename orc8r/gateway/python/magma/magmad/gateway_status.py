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
# pylint: disable=broad-except

import json
import logging
import psutil
import platform
import time
import netifaces
from typing import NamedTuple, List, Any, Dict, Optional, Tuple
from collections.abc import KeysView
from magma.common.misc_utils import (
    get_ip_from_if,
    is_interface_up,
    get_all_ips_from_if_cidr,
    get_if_mac_address,
    IpPreference,
)
from magma.common.service import MagmaService
from magma.magmad.check.machine_check.cpu_info import get_cpu_info
from magma.magmad.check.network_check.routing_table import get_routing_table
from magma.magmad.check.kernel_check.kernel_versions import (
    get_kernel_versions_async,
)
from magma.common.job import Job
from magma.magmad.service_poller import ServicePoller

GatewayStatus = NamedTuple(
    'GatewayStatus',
    [('machine_info', Dict[str, Any]), ('meta', Dict[str, str]),
     ('platform_info', Dict[str, Any]), ('system_status', Dict[str, Any])])

SystemStatus = NamedTuple(
    'SystemStatus',
    [('time', int), ('uptime_secs', int), ('cpu_user', int),
     ('cpu_system', int), ('cpu_idle', int), ('mem_total', int),
     ('mem_available', int), ('mem_used', int), ('mem_free', int),
     ('swap_total', int), ('swap_used', int), ('swap_free', int),
     ('disk_partitions', List[Dict[str, Any]])])

PlatformInfo = NamedTuple(
    'PlatformInfo',
    [('vpn_ip', str), ('packages', List[Dict[str, Any]]),
     ('kernel_version', str), ('kernel_versions_installed', List[str]),
     ('config_info',  Dict[str, Any])])

MachineInfo = NamedTuple(
    'MachineInfo',
    [('cpu_info', Dict[str, Any]), ('network_info', Dict[str, Any])])

NetworkInfo = NamedTuple(
    'NetworkInfo',
    [('network_interfaces', List[Dict[str, Any]]),
     ('routing_table', List[Dict[str, Any]])])

DiskPartition = NamedTuple(
    'DiskPartition',
    [('device', str), ('mount_point', str), ('total', int), ('used', int),
     ('free', int)])

ConfigInfo = NamedTuple(
    'ConfigInfo',
    [('mconfig_created_at', int)])

Package = NamedTuple(
    'Package',
    [('name', str), ('version', str)])

CPUInfo = NamedTuple(
    'CPUInfo',
    [('core_count', int), ('threads_per_core', int), ('architecture', str),
     ('model_name', str)])

NetworkInterface = NamedTuple(
    'NetworkInterface',
    [('network_interface_id', str), ('mac_address', str),
     ('ip_addresses', List[str]), ('status', str),
     ('ipv6_addresses', List[str])])


class KernelVersionsPoller(Job):
    """
    KernelVersionsPoller will periodically call get_kernel_versions_async and
    store the result. get_kernel_versions_installed can be called to get the
    latest list.
    """
    def __init__(self, service):
        super().__init__(
            interval=service.mconfig.checkin_interval,
            loop=service.loop
        )
        self._kernel_versions_installed = []

    def get_kernel_versions_installed(self) -> List[str]:
        """ returns the latest list of kernel versions gathered from _run """
        return self._kernel_versions_installed

    async def _run(self):
        try:
            result = await get_kernel_versions_async(loop=self._loop)
            result = list(result)[0].kernel_versions_installed
            self._kernel_versions_installed = result
        except Exception as e:
            logging.error("Error getting kernel versions! %s", e)


class GatewayStatusFactory:
    """
    GatewayStatusFactory is used to generate an object with information about
    the gateway. get_serialized_status is the public interface used to generate
    the gateway status object. The object mimics the swagger spec for
    GatewayStatus defined in the orc8r.
    """
    def __init__(self, service: MagmaService,
                 service_poller: ServicePoller,
                 kernel_version_poller: Optional[KernelVersionsPoller]):
        self._service = service
        self._service_poller = service_poller

        # Get one time status info
        self._kernel_version = platform.uname().release
        self._boot_time = psutil.boot_time()
        self._cpu_info = self._get_cpu_info()

        self._kernel_version_poller = kernel_version_poller

        # track services that are required to have non empty meta in order to
        # check-in
        self._required_service_metas = frozenset(
            service.config.get("skip_checkin_if_missing_meta_services", []))

    def get_serialized_status(self) -> Tuple[str, bool]:
        """
        get_serialized_status returns a tuple of the serialized gateway status
        and a boolean on whether or not its meta fields has all the required
        services specified by the service config.
        """
        system_status = self._system_status()._asdict()
        platform_info = \
            self._get_platform_info()._asdict()
        machine_info = self._get_machine_info()._asdict()

        gw_status = GatewayStatus(
            machine_info=machine_info,
            platform_info=platform_info,
            system_status=system_status,
            meta={},
        )
        gw_status, meta_services = self._fill_in_meta(gw_status)

        has_required_service_meta = \
            self._meta_has_required_services(meta_services)
        return json.dumps(gw_status._asdict(), default=str), \
               has_required_service_meta

    def _fill_in_meta(
        self, gw_status: GatewayStatus
    ) -> Tuple[GatewayStatus, KeysView]:
        service_status_meta = self._gather_service_status_meta()
        for statusmeta in service_status_meta.values():
            gw_status.meta.update(statusmeta)
        return gw_status, service_status_meta.keys()

    def _gather_service_status_meta(self):
        """
        returns map of (name: statusmeta) of each service
        """
        status_meta_by_name = {}
        for name, info in sorted(self._service_poller.service_info.items()):
            if info.status is not None:
                if len(info.status.meta) == 0:
                    continue
                status_meta_by_name[name] = info.status.meta
        return status_meta_by_name

    @staticmethod
    def _get_cpu_info() -> CPUInfo:
        cpu_info = get_cpu_info()
        if cpu_info.error is not None:
            logging.error('Failed to get cpu info: %s', cpu_info.error)
        return CPUInfo(
            core_count=cpu_info.core_count,
            threads_per_core=cpu_info.threads_per_core,
            architecture=cpu_info.architecture,
            model_name=cpu_info.model_name,
        )

    def _system_status(self) -> SystemStatus:
        cpu = psutil.cpu_times()
        mem = psutil.virtual_memory()
        swap = psutil.swap_memory()

        def partition_gen():
            for partition in psutil.disk_partitions():
                usage = psutil.disk_usage(partition.mountpoint)
                yield DiskPartition(
                    device=partition.device,
                    mount_point=partition.mountpoint,
                    total=usage.total,
                    used=usage.used,
                    free=usage.free,
                )

        system_status = SystemStatus(
            time=int(time.time()),
            uptime_secs=int(time.time() - self._boot_time),
            cpu_user=int(cpu.user * 1000),  # convert second to millisecond
            cpu_system=int(cpu.system * 1000),
            cpu_idle=int(cpu.idle * 1000),
            mem_total=mem.total,
            mem_available=mem.available,
            mem_used=mem.used,
            mem_free=mem.free,
            swap_total=swap.total,
            swap_used=swap.used,
            swap_free=swap.free,
            disk_partitions=[partition._asdict() for partition in
                             partition_gen()],
        )
        return system_status

    def _get_platform_info(self) -> PlatformInfo:
        try:
            gw_ip = get_ip_from_if('tun0')  # look for tun0 interface
        except ValueError:
            gw_ip = 'N/A'

        kernel_versions_installed = []
        if self._kernel_version_poller is not None:
            kernel_versions_installed = \
                self._kernel_version_poller.get_kernel_versions_installed()

        platform_info = PlatformInfo(
            vpn_ip=gw_ip,
            packages=[
                Package(
                    name='magma',
                    version=self._service.version,
                )._asdict(),
            ],
            kernel_version=self._kernel_version,
            kernel_versions_installed=kernel_versions_installed,
            config_info=ConfigInfo(
                mconfig_created_at=self._service.mconfig_metadata.created_at,
            )._asdict(),
        )
        return platform_info

    def _get_machine_info(self) -> MachineInfo:
        machine_info = MachineInfo(
            cpu_info=self._cpu_info._asdict(),
            network_info=self._get_network_info()._asdict(),
        )
        return machine_info

    @staticmethod
    def _get_network_info() -> NetworkInfo:
        def network_interface_gen():
            for interface in netifaces.interfaces():
                try:
                    mac_address = get_if_mac_address(interface)
                except ValueError:
                    mac_address = None

                try:
                    ip_addresses = get_all_ips_from_if_cidr(
                        interface, IpPreference.IPV4_ONLY)
                except ValueError:
                    ip_addresses = []

                try:
                    ipv6_addresses = get_all_ips_from_if_cidr(
                        interface, IpPreference.IPV6_ONLY)
                except ValueError:
                    ipv6_addresses = []

                yield NetworkInterface(
                    network_interface_id=interface,
                    status="UP" if is_interface_up(interface) else "DOWN",
                    mac_address=mac_address,
                    ip_addresses=ip_addresses,
                    ipv6_addresses=ipv6_addresses,
                )

        routing_cmd_result = get_routing_table()
        if routing_cmd_result.error is not None:
            logging.error("Failed to get routing table: %s",
                          routing_cmd_result.error)

        network_info = NetworkInfo(
            network_interfaces=[
                network_interface._asdict() for network_interface in
                network_interface_gen()],
            routing_table=routing_cmd_result.routing_table,
        )
        return network_info

    def _meta_has_required_services(self, meta_services: List[str]) -> bool:
        """
        Verifies based on status meta pulled from service_poller.
        service_statusmeta contains map of service_name -> statusmeta
        returns True if gateway state reporting is allowed
        """
        got_meta_services = set(meta_services)

        has_required_services = \
            got_meta_services.issuperset(self._required_service_metas)

        if not has_required_services:
            logging.warning(
                "Missing meta from services: %s "
                "(specified in cfg skip_checkin_if_missing_meta_services)",
                ", ".join(self._required_service_metas - got_meta_services))

        return has_required_services
