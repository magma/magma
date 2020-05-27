"""
Copyright (c) 2016-present, Facebook, Inc.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree. An additional grant
of patent rights can be found in the PATENTS file in the same directory.
"""

import asyncio
import logging
import pathlib

from magma.common.service import MagmaService
from magma.configuration.service_configs import load_service_config
from magma.magmad.upgrade.upgrader import UpgraderFactory
from magma.magmad.upgrade.upgrader2 import ImageNameT, UpgradeIntent, \
    Upgrader2, VersionInfo, VersionT, run_command

MAGMA_GITHUB_PATH = "/tmp/magma_upgrade"
MAGMA_GITHUB_URL = "https://github.com/facebookincubator/magma.git"


class DockerUpgrader(Upgrader2):
    """
    Downloads and installs images
    """

    def perform_upgrade_if_necessary(self, target_version: str) -> None:
        """
        Target version comes from tier configuration currently.
        """

        if self.upgrade_task and not self.upgrade_task.done():
            logging.info("Not starting another upgrade, upgrade in progress")
            return
        self.upgrade_task = self.loop.create_task(self.do_docker_upgrade())

    def version_to_image_name(self, version: VersionT) -> ImageNameT:
        split_version = version.split('|')
        if len(split_version) > 2:
            raise ValueError(
                'Expected version formatted as '
                '<image tag>|<git hash>, got {}'.format(version),
            )
        return ImageNameT(version)

    async def get_upgrade_intent(self) -> UpgradeIntent:
        """
        Returns the desired version tag for the gateway.
        """
        version_info = await asyncio.gather(self.get_versions())
        current_version = version_info[0].current_version
        tgt_version = self.service.mconfig.package_version
        if tgt_version is None or tgt_version == "":
            logging.warning('magmad package_version not found, '
                            'using current tag: %s as target tag.',
                            current_version)
            return UpgradeIntent(stable=VersionT(current_version),
                                 canary=VersionT(""))

        tgt_tag = self.version_to_image_name(tgt_version)
        return UpgradeIntent(stable=VersionT(tgt_tag), canary=VersionT(""))

    async def get_versions(self) -> VersionInfo:
        """ Returns the current version by parsing the IMAGE_VERSION in the
        .env file
        """
        with open('/var/opt/magma/docker/.env', 'r') as env:
            for line in env:
                if line.startswith("IMAGE_VERSION="):
                    current_version = line.split("=")[1].strip()
                    break

        return VersionInfo(
            current_version=current_version,
            available_versions=set(),
        )

    async def prepare_upgrade(
        self, version: VersionT, path_to_image: pathlib.Path
    ) -> None:
        """Install the new docker-compose file"""
        gw_module = self.service.config["upgrader_factory"]\
            .get("gateway_module")

        # Update any mounted static configs
        await run_command("cp -TR {}/magma/{}/gateway/configs /etc/magma".
                          format(MAGMA_GITHUB_PATH, gw_module),
                          shell=True, check=True)
        # Update any mounted template configs
        await run_command("cp -TR {}/magma/orc8r/gateway/configs/templates "
                          "/etc/magma/templates".format(MAGMA_GITHUB_PATH),
                          shell=True, check=True)
        # Copy updated docker-compose
        await run_command("cp {}/magma/{}/gateway/docker/docker-compose.yml "
                          "/var/opt/magma/docker".format(MAGMA_GITHUB_PATH,
                                                         gw_module),
                          shell=True, check=True)

    async def upgrade(
            self, version: VersionT, path_to_image: pathlib.Path,
    ) -> None:
        """Upgrade is a no-op as an external process (e.g. cron) must
        trigger it
        """
        pass

    async def do_docker_upgrade(self) -> None:
        """
           Do a single loop of the upgrade process, using the given upgrader.
        """
        try:
            await self._do_docker_upgrade()
        except Exception as exp:  # pylint: disable=broad-except
            logging.error("Upgrade loop failed! Error: %s", exp)

    async def _do_docker_upgrade(self) -> None:
        upgrade_intent, version_info = await asyncio.gather(
            self.get_upgrade_intent(), self.get_versions()
        )
        current_version = version_info.current_version

        # For back-compat, checkout from master if the version doesn't have a
        # git hash appended
        version_parts = upgrade_intent.stable.split('|')
        if len(version_parts) == 2:
            target_image, git_hash = version_parts
        else:
            logging.info('No target git hash was found, will pull configs '
                         'from master')
            target_image, git_hash = version_parts[0], 'master'

        if target_image != current_version:
            logging.info(
                "There is work to be done:\n"
                "  current: %s\n"
                "  to_upgrade: %s",
                current_version,
                target_image,
            )

            use_proxy = self.service.config["upgrader_factory"]\
                .get("use_proxy", True)

            await download_update(target_image, git_hash, use_proxy)
            await self.prepare_upgrade(
                current_version,
                pathlib.Path(MAGMA_GITHUB_PATH, "magma"),
            )

            # As a last step, update the IMAGE_VERSION in .env
            sed_args = "sed -i s/IMAGE_VERSION={}/IMAGE_VERSION={}/g " \
                       "var/opt/magma/docker/.env".format(current_version,
                                                          target_image)
            logging.info("Successfully downloaded version %s! Awaiting docker "
                         "container recreation...", target_image)
            await run_command(sed_args, shell=True, check=True)
        else:
            logging.info(
                'Service is currently on image tag %s, '
                'ignoring upgrade to tag %s, since they\'re equal.',
                current_version, target_image
            )


async def download_update(
    target_image: str,
    git_hash: str,
    use_proxy: bool,
) -> None:
    """
    Download the images for the given tag and clones the github repo.
    """
    await run_command("rm -rf {}".format(MAGMA_GITHUB_PATH), shell=True,
                      check=True)
    await run_command("mkdir -p {}".format(MAGMA_GITHUB_PATH), shell=True,
                      check=True)

    control_proxy_config = load_service_config('control_proxy')
    await run_command("cp {} /usr/local/share/ca-certificates/rootCA.crt".
                      format(control_proxy_config['rootca_cert']), shell=True,
                      check=True)
    await run_command("update-ca-certificates", shell=True, check=True)

    if use_proxy:
        git_clone_cmd = "git -c http.proxy=https://{}:{} -C {} clone {}".format(
            control_proxy_config['bootstrap_address'],
            control_proxy_config['bootstrap_port'], MAGMA_GITHUB_PATH,
            MAGMA_GITHUB_URL)
    else:
        git_clone_cmd = "git -C {} clone {}".format(MAGMA_GITHUB_PATH,
            MAGMA_GITHUB_URL)

    await run_command(git_clone_cmd, shell=True, check=True)

    git_checkout_cmd = "git -C {}/magma checkout {}".format(
        MAGMA_GITHUB_PATH, git_hash,
    )
    await run_command(git_checkout_cmd, shell=True, check=True)
    docker_login_cmd = "docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD " \
                       "$DOCKER_REGISTRY"
    await run_command(docker_login_cmd, shell=True, check=True)
    docker_pull_cmd = "IMAGE_VERSION={} docker-compose --project-directory " \
                      "/var/opt/magma/docker -f " \
                      "/var/opt/magma/docker/docker-compose.yml pull -q".\
        format(target_image)
    await run_command(docker_pull_cmd, shell=True, check=True)


class DockerUpgraderFactory(UpgraderFactory):
    """ Returns an instance of the DockerUpgrader """

    def create_upgrader(
        self,
        magmad_service: MagmaService,
        loop: asyncio.AbstractEventLoop,
    ) -> DockerUpgrader:
        return DockerUpgrader(magmad_service)
