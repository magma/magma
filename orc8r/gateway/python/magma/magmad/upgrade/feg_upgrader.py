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
import os
import shutil

from magma.common.service import MagmaService
from magma.magmad.upgrade.magma_upgrader import compare_package_versions
from magma.magmad.upgrade.upgrader import UpgraderFactory
from magma.magmad.upgrade.upgrader2 import ImageNameT, run_command, \
    UpgradeIntent, Upgrader2, VersionInfo, VersionT

IMAGE_INSTALL_DIR = '/var/cache/magma_feg'
IMAGE_INSTALL_SCRIPT = IMAGE_INSTALL_DIR + '/install.sh'


class FegUpgrader(Upgrader2):
    """
    Downloads and installs the federation gateway images
    """

    def version_to_image_name(self, version: VersionT) -> ImageNameT:
        """
        Returns the image format from the version string.
        (i.e) 0.3.68-1541626353-d1c29db1 -> magma_feg_d1c29db1.zip
        """
        parts = version.split("-")
        if len(parts) != 3:
            raise ValueError("Unknown version format: %s" % version)
        return ImageNameT("magma_feg_%s.zip" % parts[2])

    async def get_upgrade_intent(self) -> UpgradeIntent:
        """
        Returns the desired version for the gateway.
        We don't support downgrading, and so checks are made to update
        only if the target version is higher than the current version.
        """
        tgt_version = self.service.mconfig.package_version
        curr_version = self.service.version
        if (tgt_version == "0.0.0-0" or
                compare_package_versions(curr_version, tgt_version) <= 0):
            tgt_version = curr_version
        return UpgradeIntent(stable=VersionT(tgt_version), canary=VersionT(""))

    async def get_versions(self) -> VersionInfo:
        """ Returns the current version """
        return VersionInfo(
            current_version=self.service.version,
            available_versions=set(),
        )

    async def prepare_upgrade(
        self, version: VersionT, path_to_image: pathlib.Path
    ) -> None:
        """ No-op for the feg upgrader """
        return

    async def upgrade(
        self, version: VersionT, path_to_image: pathlib.Path
    ) -> None:
        """ Time to actually upgrade the Feg using the image """
        # Extract the image to the install directory
        shutil.rmtree(IMAGE_INSTALL_DIR, ignore_errors=True)
        os.mkdir(IMAGE_INSTALL_DIR)
        await run_command("unzip", str(path_to_image), "-d", IMAGE_INSTALL_DIR)
        logging.info("Running image install script: %s", IMAGE_INSTALL_SCRIPT)
        await run_command(IMAGE_INSTALL_SCRIPT)


class FegUpgraderFactory(UpgraderFactory):
    """ Returns an instance of the FegUpgrader """

    def create_upgrader(
        self,
        magmad_service: MagmaService,
        loop: asyncio.AbstractEventLoop,
    ) -> FegUpgrader:
        return FegUpgrader(magmad_service)
