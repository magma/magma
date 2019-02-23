"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import logging
import platform
import sys
import time

import grpc
import psutil
import snowflake
from orc8r.protos.magmad_pb2 import CheckinRequest, SystemStatus
from orc8r.protos.magmad_pb2_grpc import CheckindStub

from magma.common.misc_utils import get_ip_from_if
from magma.common.sdwatchdog import SDWatchdogTask
from magma.common.service_registry import ServiceRegistry
from magma.magmad.kernel_check.kernel_versions import get_kernel_versions_async
from magma.magmad.metrics import CHECKIN_STATUS


class CheckinManager(SDWatchdogTask):
    """
    Periodically sends checkin message to the cloud controller
    """

    def __init__(self, service, service_poller):
        super().__init__()  # runs SDWatchdogTask.__init__()

        self._loop = service.loop
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
        self._kernel_version = platform.uname().release

        self._kernel_versions_installed = []
        if service.config.get('enable_kernel_version_checking', False):
            self._periodic_check_kernel_versions()

        # set initial checkin timeout to "large" since no checkins occur until
        #   bootstrap succeeds.
        self.SetSDWatchdogTimeout(60 * 60 * 2)
        # initially set task as alive to wait for bootstrap, where try_checkin()
        #   will recheck alive status
        self.SetSDWatchdogAlive()

    def set_failure_cb(self, checkin_failure_cb):
        self._checkin_failure_cb = checkin_failure_cb

    def _checkin(self, service_statusmeta):
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
        cpu = psutil.cpu_times()
        mem = psutil.virtual_memory()
        try:
            gw_ip = get_ip_from_if('tun0')  # look for tun0 interface
        except ValueError:
            gw_ip = 'N/A'

        request = CheckinRequest(
            gateway_id=snowflake.snowflake(),
            magma_pkg_version=self._service.version,
            system_status=SystemStatus(
                cpu_user=int(cpu.user * 1000),  # convert second to millisecond
                cpu_system=int(cpu.system * 1000),
                cpu_idle=int(cpu.idle * 1000),
                mem_total=mem.total,
                mem_available=mem.available,
                mem_used=mem.used,
                mem_free=mem.free,
                uptime_secs=int(time.time() - self._boot_time),
            ),
            vpn_ip=gw_ip,
            kernel_version=self._kernel_version,
            kernel_versions_installed=self._kernel_versions_installed,
        )

        for statusmeta in service_statusmeta.values():
            request.status.meta.update(statusmeta)

        future = self._checkin_client.Checkin.future(
            request, mconfig.checkin_timeout,
        )
        future.add_done_callback(
            lambda f: self._loop.call_soon_threadsafe(self._checkin_done, f),
        )

    def _checkin_done(self, future):
        err = future.exception()
        if err:
            logging.error("Checkin Error! [%s] %s", err.code(), err.details())
            CHECKIN_STATUS.set(0)
            self.num_failed_checkins += 1
            if self.num_failed_checkins == self.CHECKIN_FAIL_THRESHOLD:
                logging.info('Checkin failure threshold met, remediating...')
                if self._checkin_failure_cb is not None:
                    self._checkin_failure_cb(err.code())
            self._try_reuse_checkin_client(err.code())
        else:
            CHECKIN_STATUS.set(1)
            self._checkin_client = None
            self.num_failed_checkins = 0
            logging.info("Checkin Successful!")

    def _periodic_check_kernel_versions(self):
        mconfig = self._service.mconfig
        self._check_kernel_versions()
        self._loop.call_later(
            max(mconfig.checkin_interval, 5),
            self._periodic_check_kernel_versions
        )

    def _check_kernel_versions(self):
        get_kernel_versions_future = self._loop.create_task(
            get_kernel_versions_async()
        )
        get_kernel_versions_future.add_done_callback(
            lambda future:
            self._loop.call_soon(self._get_kernel_versions_done, future)
        )

    def _get_kernel_versions_done(self, future):
        err = future.exception()
        if err:
            logging.error(
                "Error getting kernel versions! %s",
                str(err))
        else:
            result = list(future.result())[0].kernel_versions_installed
            self._kernel_versions_installed = result


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

    def try_checkin(self):
        """
        Attempt to check in. Continue to schedule future checkins

        Uses self.num_skipped_checkins to track skipped checkins
        """
        mconfig = self._service.mconfig
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

            self._checkin(service_statusmeta)

        finally:
            # always schedule the next checkin, don't allow interval < 5 sec
            checkin_interval = max(5, mconfig.checkin_interval)
            self._loop.call_later(checkin_interval, self.try_checkin)

            # flag to ensure the loop is still running, successfully or not
            self.SetSDWatchdogAlive()
            # reset checkin timeout to config plus buffer
            self.SetSDWatchdogTimeout(checkin_interval * 2)
