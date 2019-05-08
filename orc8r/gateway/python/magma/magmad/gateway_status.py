"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import json
import logging
import psutil
import platform
import time
import netifaces
from typing import (
    NamedTuple,
    List,
    Any,
    Dict,
)
from magma.common.misc_utils import (
    get_ip_from_if,
    is_interface_up,
    get_all_ips_from_if_cidr,
    get_if_mac_address,
    IpPreference,
)
from magma.common.service import MagmaService
from magma.magmad.check.machine_check.cpu_info import (
    get_cpu_info,
)
from magma.magmad.check.network_check.routing_table import (
    get_routing_table,
)


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
     ('ipv6_addresses', List[str])]
)

Route = NamedTuple(
    'Route',
    [('destination_ip', str), ('gateway_ip', str),
     ('genmask', str), ('network_interface_id', str)])


class GatewayStatusNative:
    def __init__(self, service: MagmaService):
        self._service = service
        self._kernel_version = platform.uname().release
        self._boot_time = psutil.boot_time()
        cpu_info = get_cpu_info()
        if cpu_info.error is not None:
            logging.error('Failed to get cpu info: %s', cpu_info.error)
        self._cpu_info = CPUInfo(
            core_count=cpu_info.core_count,
            threads_per_core=cpu_info.threads_per_core,
            architecture=cpu_info.architecture,
            model_name=cpu_info.model_name,
        )

    def make_status(
            self,
            service_statusmeta: Dict[str, Any],
            kernel_versions_installed: List[str]) -> str:
        system_status = self._system_status_tuple()._asdict()
        platform_info = \
            self._get_platform_info_tuple(kernel_versions_installed)._asdict()
        machine_info = self._get_machine_info_tuple()._asdict()

        gw_status = GatewayStatus(
            machine_info=machine_info,
            platform_info=platform_info,
            system_status=system_status,
            meta={},
        )
        for statusmeta in service_statusmeta.values():
            gw_status.meta.update(statusmeta)

        return json.dumps(gw_status._asdict())

    def _system_status_tuple(self) -> SystemStatus:
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

    def _get_platform_info_tuple(
            self, kernel_versions: List[str]) -> PlatformInfo:
        try:
            gw_ip = get_ip_from_if('tun0')  # look for tun0 interface
        except ValueError:
            gw_ip = 'N/A'

        mconfig_metadata = self._service.mconfig_metadata

        platform_info = PlatformInfo(
            vpn_ip=gw_ip,
            packages=[
                Package(
                    name='magma',
                    version=self._service.version,
                )._asdict(),
            ],
            kernel_version=self._kernel_version,
            kernel_versions_installed=kernel_versions,
            config_info=ConfigInfo(
                mconfig_created_at=mconfig_metadata.created_at,
            )._asdict(),
        )
        return platform_info

    def _get_machine_info_tuple(self) -> MachineInfo:
        machine_info = MachineInfo(
            cpu_info=self._cpu_info._asdict(),
            network_info=self._get_network_info_tuple()._asdict(),
        )
        return machine_info

    def _get_network_info_tuple(self) -> NetworkInfo:
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

        def make_route_tuple(route) -> Route:
            return Route(
                destination_ip=route.destination,
                gateway_ip=route.gateway,
                genmask=route.genmask,
                network_interface_id=route.interface,
            )

        routing_cmd_result = get_routing_table()
        if routing_cmd_result.error is not None:
            logging.error("Failed to get routing table: %s",
                          routing_cmd_result.error)

        network_info = NetworkInfo(
            network_interfaces=[
                network_interface._asdict() for network_interface in
                network_interface_gen()],
            routing_table=[
                make_route_tuple(route)._asdict() for
                route in routing_cmd_result.routing_table],
        )
        return network_info
