"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""
# pylint: disable=broad-except

import logging
import platform
import sys
import time
import asyncio
import netifaces

import grpc
import psutil
import snowflake
from orc8r.protos.magmad_pb2 import (
    CheckinRequest,
    SystemStatus,
    DiskPartition,
    PlatformInfo,
    MachineInfo,
    Package,
    CPUInfo,
    NetworkInfo,
    NetworkInterface,
    Route,
)
from orc8r.protos.magmad_pb2_grpc import CheckindStub

from magma.common.misc_utils import (
    get_ip_from_if,
    is_interface_up,
    get_all_ips_from_if_cidr,
    get_if_mac_address,
    IpPreference,
)
from magma.common.rpc_utils import grpc_async_wrapper
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service_registry import ServiceRegistry
from magma.magmad.check.kernel_check.kernel_versions import (
    get_kernel_versions_async,
)
from magma.magmad.check.machine_check.cpu_info import (
    get_cpu_info,
)
from magma.magmad.check.network_check.routing_table import (
    get_routing_table,
)
from magma.magmad.metrics import CHECKIN_STATUS


class CheckinManager(SDWatchdogTask):
    """
    Periodically sends checkin message to the cloud controller
    """

    def __init__(self, service, service_poller):
        super().__init__(
            max(5, service.mconfig.checkin_interval),
            service.loop
        )

        self._service = service
        self._service_poller = service_poller

        # Number of consecutive failed checkins before we check for an outdated
        # cert
        self.CHECKIN_FAIL_THRESHOLD = 10
        # Current number of consecutive failed checkins
        self.num_failed_checkins = 0
        self._checkin_failure_cb = None

        # cloud controller's client stub
        self._checkin_client = None
        self.MAX_CLIENT_REUSE = 60

        # skip checkin based on missing status meta
        self.num_skipped_checkins = 0

        # One time status info
        self._boot_time = psutil.boot_time()
        self._kernel_name = platform.system()
        self._kernel_version = platform.uname().release
        cpu_info = get_cpu_info()
        if cpu_info.error is not None:
            logging.error('Failed to get cpu info: %s', cpu_info.error)
        self._cpu_info = CPUInfo(
            core_count=cpu_info.core_count,
            threads_per_core=cpu_info.threads_per_core,
            architecture=cpu_info.architecture,
            model_name=cpu_info.model_name,
        )

        self._kernel_versions_installed = []
        self._periodically_check_kernel_versions = \
            service.config.get('enable_kernel_version_checking', False)
        # set initial checkin timeout to "large" since no checkins occur until
        #   bootstrap succeeds.
        self.set_timeout(60 * 60 * 2)
        # initially set task as alive to wait for bootstrap, where try_checkin()
        #   will recheck alive status
        self.heartbeat()

        # Start try_checkin loop
        self.start()

    async def try_checkin(self):
        """
        Attempt to check in. Continue to schedule future checkins

        Uses self.num_skipped_checkins to track skipped checkins
        """
        config = self._service.config

        # specifies number of checkin iterations that can have an empty/missing
        #   service meta before checking in anyway.
        # If 0, then never check in if missing.
        #  Use safe default to make "forever" explicit.
        # (check config early so config is validated)
        max_skipped_checkins = int(config.get("max_skipped_checkins", 3))

        try:
            # gather information required to determine checkin
            service_statusmeta = self._gather_service_statusmeta()

            # use necessary information to determine can_checkin
            can_checkin = self._can_checkin(service_statusmeta)

            if can_checkin:
                # we can checkin! Continue on below to actually _checkin()
                # clear fail count
                self.num_skipped_checkins = 0
            else:
                # we should only not checkin up to a specified limit, at
                #  which time we checkin anyway
                if 0 < max_skipped_checkins < self.num_skipped_checkins:
                    logging.warning(
                        "Number of skipped checkins exceeds %d "
                        "(cfg: max_skipped_checkins). Checking in anyway.",
                        max_skipped_checkins)
                    # intentionally don't reset num_skipped_checkins here
                else:
                    self.num_skipped_checkins += 1
                    return
            await self._checkin(service_statusmeta)

        finally:
            # reset checkin timeout to config plus buffer
            self.set_timeout(self._interval * 2)

    def set_failure_cb(self, checkin_failure_cb):
        self._checkin_failure_cb = checkin_failure_cb

    async def _run(self):
        """
        This functions gets run in a loop in job.py
        """
        if self._periodically_check_kernel_versions:
            await self._check_kernel_versions()
        await self.try_checkin()

    async def _checkin(self, service_statusmeta):
        """
        if previous checkin is successful, create a new channel
        (to make sure the channel does't become stale). Otherwise,
        keep the existing channel.
        """
        if self._checkin_client is None:
            chan = ServiceRegistry.get_rpc_channel(
                    'checkind', ServiceRegistry.CLOUD)
            self._checkin_client = CheckindStub(chan)

        mconfig = self._service.mconfig

        request = CheckinRequest(
            gateway_id=snowflake.snowflake(),
            system_status=self._system_status(),
            platform_info=self._platform_info(),
            machine_info=self._machine_info(),
        )
        logging.debug('Checkin request:\n%s', request)

        for statusmeta in service_statusmeta.values():
            request.status.meta.update(statusmeta)

        try:
            await grpc_async_wrapper(
                self._checkin_client.Checkin.future(
                    request, mconfig.checkin_timeout,
                ),
                self._loop)
            self._checkin_done()
        except grpc.RpcError as err:
            self._checkin_error(err)

    def _system_status(self):
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

        return SystemStatus(
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
            disk_partitions=[partition for partition in partition_gen()],
        )

    def _platform_info(self):
        try:
            gw_ip = get_ip_from_if('tun0')  # look for tun0 interface
        except ValueError:
            gw_ip = 'N/A'

        return PlatformInfo(
            vpn_ip=gw_ip,
            packages=[
                Package(
                    name='magma',
                    version=self._service.version,
                ),
            ],
            kernel_version=self._kernel_version,
            kernel_versions_installed=self._kernel_versions_installed,
        )

    def _network_info(self):
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
                    status=NetworkInterface.UP if is_interface_up(
                        interface) else NetworkInterface.DOWN,
                    mac_address=mac_address,
                    ip_addresses=ip_addresses,
                    ipv6_addresses=ipv6_addresses,
                )

        routing_cmd_result = get_routing_table()
        if routing_cmd_result.error is not None:
            logging.error("Failed to get routing table: %s",
                          routing_cmd_result.error)

        return NetworkInfo(
            network_interfaces=[network_interface for network_interface in
                                network_interface_gen()],
            routing_table=[Route(destination_ip=route.destination,
                                 gateway_ip=route.gateway,
                                 genmask=route.genmask,
                                 network_interface_id=route.interface) for
                           route in routing_cmd_result.routing_table],
        )

    def _machine_info(self):
        return MachineInfo(
            cpu_info=self._cpu_info,
            network_info=self._network_info(),
        )

    def _checkin_error(self, err):
        logging.error("Checkin Error! [%s] %s", err.code(), err.details())
        CHECKIN_STATUS.set(0)
        self.num_failed_checkins += 1
        if self.num_failed_checkins == self.CHECKIN_FAIL_THRESHOLD:
            logging.info('Checkin failure threshold met, remediating...')
            if self._checkin_failure_cb is not None:
                asyncio.ensure_future(
                    self._checkin_failure_cb(err.code()), loop=self._loop)
        self._try_reuse_checkin_client(err.code())

    def _checkin_done(self):
        CHECKIN_STATUS.set(1)
        self._checkin_client = None
        self.num_failed_checkins = 0
        logging.info("Checkin Successful!")

    async def _check_kernel_versions(self):
        try:
            result = await get_kernel_versions_async()
            result = list(result)[0].kernel_versions_installed
            self._kernel_versions_installed = result
        except Exception as e:
            logging.error("Error getting kernel versions! %s", e)

    def _try_reuse_checkin_client(self, err_code):
        """
        Try to reuse the checkin client if possible. We are yet to fix a
        grpc behavior, where if DNS request blackholes then the DNS request
        is retried infinitely even after the channel is deleted. To prevent
        running out of fds, we try to reuse the channel during such failures
        as much as possible.
        """
        if err_code != grpc.StatusCode.DEADLINE_EXCEEDED:
            # Not related to the DNS issue
            self._checkin_client = None
        if (self.num_failed_checkins % self.MAX_CLIENT_REUSE) == 0:
            logging.info('Max client reuse reached. Cleaning up client')
            self._checkin_client = None
            # Sanity check if we are not leaking fds and failing checkin
            proc = psutil.Process()
            max_fds, _ = proc.rlimit(psutil.RLIMIT_NOFILE)
            open_fds = proc.num_fds()
            logging.info('Num open fds: %d', open_fds)
            if open_fds >= (max_fds * 0.8):
                logging.error("Reached 80% of allowed fds. Restarting process")
                sys.exit(1)

    def _gather_service_statusmeta(self):
        """
        returns map of (name: statusmeta) of each service
        """
        service_statusmeta = {}
        for name, info in sorted(self._service_poller.service_info.items()):
            if info.status is not None:
                if len(info.status.meta) == 0:
                    continue
                service_statusmeta[name] = info.status.meta
        return service_statusmeta

    def _can_checkin(self, service_statusmeta):
        """
        Verifies based on status meta pulled from service_poller.

        service_statusmeta contains map of service_name -> statusmeta

        returns True if checkin is allowed
        """

        config = self._service.config

        # track services that are required to have non empty meta in order to checkin
        required_meta = frozenset(
            config.get("skip_checkin_if_missing_meta_services", []))
        got_meta = set(service_statusmeta.keys())

        can_checkin = got_meta.issuperset(required_meta)

        if not can_checkin:
            logging.warning(
                "Missing meta from service: %s "
                "(specified in cfg skip_checkin_if_missing_meta_services)",
                ", ".join(required_meta - got_meta))

        return can_checkin
