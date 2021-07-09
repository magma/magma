"""
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

===============================================================================

An alternative/replacement interface for the Upgrader interface

The main difference is that this interface provides a specialization that
ties into a cloud component which allows coordinating/pacing upgrades.

In order for that to be possible, upgrading must be broken into stages:
1) Initial state
2) Download resources to upgrade
3) Prepare resources for upgrade
4) Upgrade
"""

import abc
import asyncio
import http.client
import logging
import pathlib
import ssl
import time
import typing
import urllib.request  # TODO: Figure out how to get aiohttp

from magma.common.service import MagmaService
from magma.common.service_registry import ServiceRegistry
from magma.magmad.upgrade.upgrader import Upgrader
from prometheus_client import Gauge

VersionT = typing.NewType("VersionT", str)
ImageNameT = typing.NewType("ImageNameT", str)

# TODO - Add gauges for upgrade actions and failure

UpgraderGauges = typing.NamedTuple(
    'UpgraderGauges',
    [
        ('time_taken', Gauge),
        ('error', Gauge),
        ('downloaded', Gauge),
        ('prepared', Gauge),
        ('canary', Gauge),
        ('stable', Gauge),
        ('idle', Gauge),
    ],
)

_GAUGES = None  # type: typing.Optional[UpgraderGauges]


def get_gauges() -> UpgraderGauges:
    """Lazily create gauges only when the upgrader actually has begun running"""
    global _GAUGES  # pylint: disable=global-statement
    if not _GAUGES:
        _GAUGES = UpgraderGauges(
            time_taken=Gauge(
                "upgrader2_time_taken",
                "how long it took to do a loop of the upgrader",
            ),
            error=Gauge(
                "upgrader2_error",
                "if the last loop threw an exception",
            ),
            prepared=Gauge(
                "upgrader2_prepared",
                "if the last version prepared was successful",
            ),
            downloaded=Gauge(
                "upgrader2_downloaded",
                "if the version to prepare is downloaded",
            ),
            canary=Gauge(
                "upgrader2_canary",
                "if the current version is the canary version",
            ),
            stable=Gauge(
                "upgrader2_stable",
                "if the current version is the stable version",
            ),
            idle=Gauge(
                "upgrader2_idle",
                "if the upgrader doesn't have anything to do",
            ),
        )
    return _GAUGES


class ImageDownloadFailed(Exception):
    pass


class VersionInfo(
    # In 3.6, use the much shorter typing.NamedTuple syntax
    typing.NamedTuple(
        "VersionInfo",
        [
            # Use empty string for "no version"
            ("current_version", VersionT),
            # available_versions are those that have already had the
            # prepare_upgrade() work done on them, and can be kicked off with a
            # upgrade() call.
            #
            # available_versions does not need to contain current_version, but
            # it can if you want.
            ("available_versions", typing.Set[VersionT]),
        ],
    ),
):
    @property
    def all_versions(self) -> typing.Set[VersionT]:
        """All versions available on the device (including the current one)"""
        ret = set(self.available_versions)
        if self.current_version:
            ret.add(self.current_version)
        return ret


class UpgradeIntent(
    typing.NamedTuple(
        "UpgradeIntent",
        [
            # If set, force upgrading to this version if the current version is
            # not in stable or canary
            # use empty string for "no version set"
            ("stable", VersionT),
            # If set, install this version and wait to be told to upgrade
            # by a controller RPC
            # use empty string for "no version set"
            ("canary", VersionT),
        ],
    ),
):
    def version_to_prepare(self, version_info: VersionInfo) -> VersionT:
        """
        Returns the version to prepare for upgrade this loop, if needed.
        Else returns empty string.
        """
        available_versions = version_info.all_versions
        preference = [
            version for version in (
                self.stable, self.canary,
            ) if version
        ]
        for version in preference:
            if version not in available_versions:
                return version
        return VersionT("")

    def version_to_force_upgrade(self, version_info: VersionInfo) -> VersionT:
        """
        Returns the version to be force upgraded, if needed.
        Else returns empty string.
        """
        if not self.stable:  # No stable version = no forcing
            return VersionT("")
        if not version_info.current_version or (
            version_info.current_version not in (self.stable, self.canary)
        ):
            return self.stable
        return VersionT("")

    @property
    def active_version(self) -> VersionT:
        return self.stable or self.canary


class Upgrader2(Upgrader, metaclass=abc.ABCMeta):
    """Electric bugaloo"""

    def __init__(self, service: MagmaService) -> None:
        self.service = service
        self.upgrade_task = None  # type: typing.Optional[asyncio.Task]

    @property
    def loop(self) -> asyncio.AbstractEventLoop:
        return self.service.loop

    def perform_upgrade_if_necessary(self, target_version: str) -> None:
        """
        Target version comes from tier configuration currently.

        Additionally, if we have loop from the constructor, we don't even
        need this function - we could create our own loop, with blackjack
        and courtesans.
        """

        if self.upgrade_task and not self.upgrade_task.done():
            logging.info("Not starting another upgrade, upgrade in progress")
            return
        self.upgrade_task = self.loop.create_task(do_upgrade2(self))

    def version_to_image_name(self, version: VersionT) -> ImageNameT:
        """
        Return the image name for this device given the version string

        For deployments with heterogeneous machine types, the image will likely
        differ by machine type.
        """
        return ImageNameT(version)  # Assume name == path

    @abc.abstractmethod
    async def get_upgrade_intent(self) -> UpgradeIntent:
        """
        Returns the instructions for what this device should be installing.
        """
        # TODO - default to returning tier information
        raise NotImplementedError

    @abc.abstractmethod
    async def get_versions(self) -> VersionInfo:
        """
        Return version info for exporting to the status.

        Return the version string running on the device, as well as any others
        that are available to be upgraded.
        """
        raise NotImplementedError

    @abc.abstractmethod
    async def prepare_upgrade(
        self, version: VersionT, path_to_image: pathlib.Path,
    ) -> None:
        """
        Prepare the device for upgrade, i.e. by unpacking the image on disk.

        After this step, it is required the version appears in the
        VersionInfo.available_versions returned by get_versions().
        """
        raise NotImplementedError

    @abc.abstractmethod
    async def upgrade(self, version: VersionT, path_to_image: pathlib.Path) -> None:
        """
        Apply the upgrade to the device, potentially restarting it.

        After this is called, if the upgrade is successful, the version should
        appear in the VersionInfo.current_version returned by get_versions()

        if the device needs to reboot to finish this installation, it should be
        triggered here.
        """
        raise NotImplementedError


async def do_upgrade2(upgrader: Upgrader2) -> None:
    """
    Do a single loop of the upgrade process, using the given upgrader.

    Run the update loop, which executes the interface methods on the Upgrader2
    object in the sequence needed to orchestrate an update.
    """
    gauges = get_gauges()
    start_time = time.time()
    try:
        await _do_upgrade2(upgrader)
        gauges.error.set(0)
    except Exception as exp:  # pylint: disable=broad-except
        logging.error("Upgrade loop failed! Error: %s", exp)
        gauges.error.set(1)
    gauges.time_taken.set(time.time() - start_time)


async def _do_upgrade2(upgrader: Upgrader) -> None:
    gauges = get_gauges()
    upgrade_intent, version_info = await asyncio.gather(
        upgrader.get_upgrade_intent(), upgrader.get_versions(),
    )

    current_version = version_info.current_version or object()
    gauges.canary.set(current_version == upgrade_intent.canary)
    gauges.stable.set(current_version == upgrade_intent.stable)

    # TODO - export current version information using
    #        self.upgrader.service.register_get_status_callback
    #        which will eventually be used by cloud upgrade orchestrator

    to_prepare = upgrade_intent.version_to_prepare(version_info)
    to_upgrade_to = upgrade_intent.version_to_force_upgrade(version_info)

    if to_prepare or to_upgrade_to:
        logging.info(
            "There is work to be done:\n"
            "  canary: %s\n"
            "  stable: %s\n"
            "  current: %s\n"
            "  available: %s\n"
            "  to_prepare: %s\n"
            "  to_upgrade: %s",
            upgrade_intent.canary,
            upgrade_intent.stable,
            version_info.current_version,
            version_info.available_versions,
            to_prepare,
            to_upgrade_to,
        )
        gauges.idle.set(0)
    else:
        gauges.idle.set(1)

    # If need to do an upgrade, skip preparing a new image
    if to_upgrade_to and to_prepare not in (VersionT(""), to_upgrade_to):
        logging.warning(
            "Need to upgrade, so skipping preparation of %s",
            to_prepare,
        )
        to_prepare = VersionT("")

    if to_prepare:
        gauges.prepared.set(0)
        assert to_prepare != version_info.current_version
        image_name = upgrader.version_to_image_name(to_prepare)
        logging.info("Preparing %r, image is %r", to_prepare, image_name)
        image_path = image_local_path(image_name)
        await ensure_downloaded(
            upgrader.service.config, image_name, image_path, gauges.downloaded,
        )
        await upgrader.prepare_upgrade(to_prepare, image_path)
        logging.info("%r is prepared", to_prepare)
        gauges.prepared.set(1)
    elif upgrade_intent.active_version not in (
        VersionT(""), version_info.current_version,
    ):
        # We're not the version to be installed, but we don't need to prepare,
        # so we must be prepared
        gauges.downloaded.set(1)
        gauges.prepared.set(1)

    if to_upgrade_to:
        assert to_upgrade_to != version_info.current_version
        logging.warning(
            "Version %r is out of date! Force upgrading to %r",
            version_info.current_version,
            to_upgrade_to,
        )
        image_name = upgrader.version_to_image_name(to_upgrade_to)
        image_path = image_local_path(image_name)
        await ensure_downloaded(
            upgrader.service.config, image_name, image_path, gauges.downloaded,
        )
        logging.warning(
            "Performing upgrade to %r, device may reboot",
            to_upgrade_to,
        )
        await upgrader.upgrade(to_upgrade_to, image_path)
        logging.info(
            "Upgrade to %r is complete (I guess we didn't restart!)",
            to_upgrade_to,
        )
        gauges.downloaded.set(0)
        gauges.prepared.set(0)
    logging.info("Upgrade2 loop complete")


async def ensure_downloaded(
    config: dict,
    image_name: ImageNameT,
    image_path: pathlib.Path,
    gauge: Gauge,
) -> None:
    """
    Download the image for the given version and returns where it was put.

    Will garbage collect if needed before downloading.
    """
    # TODO - Checksum? Validate? Implementation of install could be allowed
    #        to delete if no good
    if not image_path.is_file():
        gauge.set(0)
        logging.info("Downloading %r to %s from s3", image_name, image_path)
        base_url = config["upgrader_factory"].get(
            "http_base_url", "https://api.magma.test/s3",
        )
        size = await asyncio.get_event_loop().run_in_executor(
            None,
            download,
            base_url,
            image_name,
            image_path,
        )
        gauge.set(1)
        logging.info(
            "Download of %r to %s is complete (%.2f MB)",
            image_name,
            image_path,
            size / 1024 / 1024,
        )
    else:
        gauge.set(1)
        logging.info(
            "Skipping download, image %r is already in %s", image_name, image_path,
        )


def image_local_path(image_name: ImageNameT) -> pathlib.Path:
    """Where an image to be downloaded should be stored"""
    return pathlib.Path("/var/cache", image_name)


def get_ssl_context() -> ssl.SSLContext:
    proxy_config = ServiceRegistry.get_proxy_config()
    ret = ssl.SSLContext(ssl.PROTOCOL_TLSv1_2)
    try:
        ret.load_cert_chain(
            certfile=proxy_config["gateway_cert"],
            keyfile=proxy_config["gateway_key"],
        )
    except FileNotFoundError:
        raise RuntimeError("Gateway cert or key file not found")
    return ret


def download(
    base_url: str,
    image_name: ImageNameT,
    dst: pathlib.Path,
) -> int:
    """
    Download the image from the MagmaC s3 storage

    Raises an ImageDownloadFailed exception if it fails.

    TODO: Blocking call! Could be replaced with aiohttp
    """
    assert not dst.is_dir()
    tmp_path = dst.parent.joinpath(".%s" % dst.name)
    url = "{}/{}".format(base_url, image_name)
    headers = {"Host": "download.cloud"}
    resume = False
    if tmp_path.exists():
        # Already downloaded some, maybe the destination supports
        # resume download
        headers["Range"] = "bytes={}-".format(tmp_path.stat().st_size)
        resume = True
    req = urllib.request.Request(url, headers=headers)

    with urllib.request.urlopen(
        req,
        context=get_ssl_context(),
        timeout=30,  # How long without getting data, not total time
    ) as response:
        assert isinstance(response, http.client.HTTPResponse)  # for types
        # For some reason, talking to any url on api.magma always returns
        # 200, but it sends back empty responses if the file doesn't exist
        assert response.status == 200
        size = 0
        write_mode = "w+b"
        if resume and response.getheader("Content-Range"):
            write_mode = "a+b"
            size += tmp_path.stat().st_size
            logging.warning(
                "Resuming download of %s at %.2f Mb",
                image_name,
                size / 1024 / 1024,
            )
        else:
            logging.warning(
                "Attempted to resume download of %s, range unsupported",
                image_name,
            )
            resume = False

        with tmp_path.open(write_mode) as f:
            while True:
                chunk = response.read(4 * 1024)  # 4kb
                if not chunk:
                    break
                size += len(chunk)
                f.write(chunk)
        if size:
            tmp_path.rename(dst)
        else:
            tmp_path.unlink()
            raise ImageDownloadFailed("%s is not available yet" % image_name)
        return size


async def run_command(*args, **kwargs):
    """Shortcut for asyncronous running of a command"""
    fn = asyncio.subprocess.create_subprocess_exec
    if kwargs.pop("shell", False):
        fn = asyncio.subprocess.create_subprocess_shell
    check = kwargs.pop("check", False)
    process = await fn(*args, **kwargs)
    stdout, stderr = await process.communicate()
    if check:
        if process.returncode != 0:
            raise Exception("Command failed: %s" % args)
    return process.returncode, stdout, stderr
