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
import asyncio
import logging


class Upgrader(object):
    """
    Interface for software upgraders. Implementation of the actual upgrade
    process is left up to the derived classes.
    """

    def perform_upgrade_if_necessary(self, target_version):
        """
        Perform the software upgrade if it is required, otherwise no-op.

        Args:
            target_version (str): Target software version to upgrade to
        """
        raise NotImplementedError('perform_upgrade must be implemented')


class UpgraderFactory(object):
    """
    An UpgraderFactory instantiates a specific implementation of the Upgrader
    interface for magmad to use.

    The factory is injected into magmad via the magmad.yml configuration.
    magmad will enforce that the injected factory class inherits from this
    base interface.

    Note that factories are not instantiated with any arguments in magmad.
    """

    def create_upgrader(self, magmad_service, loop):
        """
        Instantiate a concrete instance of an Upgrader.

        Args:
            magmad_service (magma.common.service.MagmaService):
                MagmaService for magmad
            loop (asyncio.AbstractEventLoop):
                asyncio event loop which the upgrader will be periodically
                queried in.

        Returns:
            magma.magmad.upgrade.upgrader.Upgrader: upgrader for magmad
        """
        raise NotImplementedError('create_upgrader must be implemented')


@asyncio.coroutine
def start_upgrade_loop(magmad_service, upgrader):
    """
    Check an Upgrader implementation in a loop and upgrade if the Upgrader
    indicates that the software needs to be upgraded.

    Args:
        magmad_service (magma.common.service.MagmaService):
            MagmaService instance for magmad
        upgrader (magma.magmad.upgrade.upgrader.Upgrader):
            Upgrader instance to use
    """
    assert isinstance(upgrader, Upgrader), 'upgrader must implement Upgrader'

    # This loop can happen within seconds of the gateway booting, before
    # even the first checkin. Delay a little bit so the device can
    # record stats/checkin/give someone an opportunity to disable
    logging.info("Waiting before checking for updates for the first time...")
    yield from asyncio.sleep(120)

    while True:
        logging.info('Checking for upgrade...')
        try:
            target_ver = _get_target_version(magmad_service.mconfig)
            upgrader.perform_upgrade_if_necessary(target_ver)
        except Exception:  # pylint: disable=broad-except
            logging.exception(
                'Error encountered while upgrading, will try again after delay',
            )
        poll_interval = max(  # No faster than 1/minute
            60,
            magmad_service.mconfig.autoupgrade_poll_interval,
        )
        yield from asyncio.sleep(poll_interval)


def _get_target_version(magmad_mconfig):
    if magmad_mconfig.package_version is None:
        logging.warning(
            'magmad package_version config not found, '
            'returning 0.0.0-0 as target package version.',
        )
        return '0.0.0-0'

    return magmad_mconfig.package_version
