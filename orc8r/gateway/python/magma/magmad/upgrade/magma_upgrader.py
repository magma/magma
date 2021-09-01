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
import re

from magma.common.misc_utils import call_process
from magma.magmad.upgrade.upgrader import Upgrader, UpgraderFactory


class MagmaUpgraderFactory(UpgraderFactory):
    """
    Upgrader factory implementation for Magma.
    """

    def create_upgrader(self, magmad_service, loop):
        return MagmaUpgrader(magmad_service.version, loop)


class MagmaUpgrader(Upgrader):
    """
    Upgrader implementation for Magma.
    """

    def __init__(self, cur_version, loop):
        super().__init__()
        self.cur_version = cur_version
        self.loop = loop

    def perform_upgrade_if_necessary(self, target_version):
        if compare_package_versions(self.cur_version, target_version) > 0:
            logging.info(
                'Upgrading magma to version %s', target_version,
            )
            call_process(
                'apt-get update',
                _get_apt_update_complete_callback(
                    target_version,
                    self.loop,
                ),
                self.loop,
            )
        else:
            logging.info(
                'Service is currently on package version %s, '
                'ignoring upgrade to %s because it is either '
                'equal or behind.', self.cur_version, target_version,
            )


VERSION_RE = re.compile(
    r'(?P<maj>\d+)\.(?P<min>\d+)\.(?P<hotfix>\d+)(-(?P<iter>\d+))?',
)


def compare_package_versions(current_version, target_version):
    """
    Compare 2 package version strings. Returns 1 if target_version is ahead of
    current_version, 0 if they are equal, and -1 if target_version is behind
    current_version.

    Returns:
        1 if target_version is ahead, 0 if both versions are the same, -1 if
        target_version is behind.
    """
    cur_parsed = VERSION_RE.match(current_version)
    target_parsed = VERSION_RE.match(target_version)

    if not cur_parsed:
        raise ValueError(
            'Could not parse current package version '
            '{}'.format(current_version),
        )
    if not target_parsed:
        raise ValueError(
            'Could not parse target package version '
            '{}'.format(target_version),
        )

    for group in ['maj', 'min', 'hotfix', 'iter']:
        target_val = int(target_parsed.group(group) or 0)
        cur_val = int(cur_parsed.group(group) or 0)
        if target_val < cur_val:
            return -1
        elif target_val > cur_val:
            return 1
    return 0


def _get_apt_update_complete_callback(target_version, loop):
    def callback(returncode):
        if returncode != 0:
            logging.error(
                "Apt-get update failed with code: %d", returncode,
            )
            return
        logging.info("apt-get update completed")
        call_process(
            get_autoupgrade_command(target_version, dry_run=True),
            _get_upgrade_dry_run_complete_callback(target_version, loop),
            loop,
        )
    return callback


def _get_upgrade_dry_run_complete_callback(target_version, loop):
    def callback(returncode):
        if returncode != 0:
            logging.error(
                "Magma Upgrade dry-run failed with code: %d", returncode,
            )
            return

        logging.info("Magma upgrade dry-run completed")
        call_process(
            get_autoupgrade_command(target_version, dry_run=False),
            _upgrade_completed,
            loop,
        )
    return callback


def _upgrade_completed(returncode):
    if returncode != 0:
        logging.error("Upgrade magma failed with code: %d", returncode)
    else:
        logging.info("Upgrade completed")


def get_autoupgrade_command(version, *, dry_run=False):
    command = 'apt-get install'
    options = [
        '-o', 'Dpkg::Options::="--force-confnew"',
        '--assume-yes', '--force-yes', '--only-upgrade',
    ]
    if dry_run:
        options.append('--dry-run')
    return '{command} {options_joined} magma={version}'.format(
        command=command,
        options_joined=' '.join(options),
        version=version,
    )
