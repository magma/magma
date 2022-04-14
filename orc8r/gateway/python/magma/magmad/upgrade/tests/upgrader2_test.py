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
import pathlib
import tempfile
import typing
import unittest
from unittest import mock

from magma.magmad.upgrade import upgrader2


async def fake_ensure_downloaded(_config, _name, path, _guage):
    assert isinstance(path, pathlib.Path)
    return None


class FakeUpgrader(upgrader2.Upgrader2):
    """Dummy implementation for testing"""

    def __init__(self, *args: typing.Any, **kwargs: typing.Any) -> None:
        super().__init__(*args, **kwargs)
        self.current_version = upgrader2.VersionT("fake")
        self.installed_versions = set()  # type: typing.Set[upgrader2.VersionT]
        self.stable_version = upgrader2.VersionT("stable_version")
        self.canary_version = upgrader2.VersionT("canary_version")

    async def get_upgrade_intent(self) -> upgrader2.UpgradeIntent:
        return upgrader2.UpgradeIntent(
            stable=self.stable_version, canary=self.canary_version,
        )

    async def get_versions(self) -> upgrader2.VersionInfo:
        return upgrader2.VersionInfo(
            current_version=self.current_version,
            available_versions=self.installed_versions,
        )

    async def prepare_upgrade(
        self, version: upgrader2.VersionT, path_to_image: pathlib.Path,
    ) -> None:
        self.installed_versions.add(version)

    async def upgrade(
        self, version: upgrader2.VersionT, path_to_image: pathlib.Path,
    ) -> None:
        self.installed_versions.remove(version)
        self.current_version = version


class Upgrader2Test(unittest.TestCase):

    version_a = upgrader2.VersionT("a")
    version_b = upgrader2.VersionT("b")
    version_c = upgrader2.VersionT("c")
    version_d = upgrader2.VersionT("d")
    version_none = upgrader2.VersionT("")

    def patch(self, *args, **kwargs):
        patcher = mock.patch(*args, **kwargs)
        ret = patcher.start()
        self.addCleanup(patcher.stop)
        return ret

    def setUp(self):
        asyncio.set_event_loop(asyncio.new_event_loop())
        service = mock.Mock(config={})
        self.upgrader = FakeUpgrader(service)
        self.wrapped_upgrader = mock.Mock(wraps=self.upgrader)

        tmpdir = tempfile.TemporaryDirectory()  # noqa: P201 cleaned up next line
        self.addCleanup(tmpdir.cleanup)
        self.image_directory = tmpdir.name

        self.patch(
            "magma.magmad.upgrade.upgrader2.image_local_path",
            lambda x: pathlib.Path(self.image_directory, x),
        )
        self.ensure_downloaded = self.patch(
            "magma.magmad.upgrade.upgrader2.ensure_downloaded",
            side_effect=fake_ensure_downloaded,
        )

    def run_upgrade_loop(self):
        asyncio.get_event_loop().run_until_complete(
            upgrader2.do_upgrade2(self.wrapped_upgrader),
        )

    def test_version_info(self):
        self.assertEqual(
            set(), upgrader2.VersionInfo(self.version_none, {}).all_versions,
        )
        self.assertEqual(
            {self.version_a}, upgrader2.VersionInfo(
                self.version_a, {},
            ).all_versions,
        )
        self.assertEqual(
            {self.version_a},
            upgrader2.VersionInfo(
                self.version_none, {
                    self.version_a,
                },
            ).all_versions,
        )
        self.assertEqual(
            {self.version_a}, upgrader2.VersionInfo(
                self.version_a, {},
            ).all_versions,
        )
        self.assertEqual(
            {self.version_a, self.version_b},
            upgrader2.VersionInfo(
                self.version_a, {self.version_b},
            ).all_versions,
        )
        self.assertEqual(
            {self.version_a, self.version_b},
            upgrader2.VersionInfo(
                self.version_a, {self.version_a, self.version_b},
            ).all_versions,
        )

    def test_upgrade_intent(self):
        version_info = upgrader2.VersionInfo(
            current_version=self.version_a,
            available_versions={self.version_b, self.version_c},
        )
        intent = upgrader2.UpgradeIntent(
            stable=self.version_none, canary=self.version_none,
        )
        self.assertEqual(
            self.version_none,
            intent.version_to_prepare(version_info),
        )
        self.assertEqual(
            self.version_none, intent.version_to_force_upgrade(version_info),
        )

        # If version is already the one intended, no actions should be taken
        intent = upgrader2.UpgradeIntent(
            stable=version_info.current_version, canary=version_info.current_version,
        )

        self.assertEqual(
            self.version_none,
            intent.version_to_prepare(version_info),
        )
        self.assertEqual(
            self.version_none, intent.version_to_force_upgrade(version_info),
        )

        # If upgrade is already prepared, then no reason to prepare again
        intent = upgrader2.UpgradeIntent(
            stable=self.version_none, canary=self.version_b,
        )

        self.assertEqual(
            self.version_none,
            intent.version_to_prepare(version_info),
        )

        # Unprepared version
        intent = upgrader2.UpgradeIntent(
            stable=self.version_none, canary=self.version_d,
        )
        self.assertEqual(
            self.version_d, intent.version_to_prepare(version_info),
        )

        # Force upgrade needed
        intent = upgrader2.UpgradeIntent(
            stable=self.version_b, canary=self.version_b,
        )
        self.assertEqual(
            self.version_b, intent.version_to_force_upgrade(version_info),
        )

    def test_do_nothing_upgrade(self):
        self.upgrader.stable_version = ""
        self.upgrader.canary_version = ""
        self.upgrader.current_version = ""
        self.upgrader.installed_versions = set()

        self.run_upgrade_loop()

        # Old-style upgrade method not called
        self.assertIs(
            0, self.wrapped_upgrader.perform_upgrade_if_necessary.call_count,
        )
        # New style methods called
        self.assertIs(1, self.wrapped_upgrader.get_upgrade_intent.call_count)
        self.assertIs(1, self.wrapped_upgrader.get_versions.call_count)
        self.assertIs(0, self.wrapped_upgrader.prepare_upgrade.call_count)
        self.assertIs(0, self.wrapped_upgrader.upgrade.call_count)
        self.assertIs(0, self.ensure_downloaded.call_count)

    def test_prepare_upgrade(self):
        self.upgrader.stable_version = ""
        self.upgrader.canary_version = "canary_version"
        self.upgrader.current_version = ""
        self.upgrader.installed_versions = set()

        canary_path = pathlib.Path(self.image_directory, "canary_version")

        self.run_upgrade_loop()

        # Old-style upgrade method not called
        self.assertIs(
            0, self.wrapped_upgrader.perform_upgrade_if_necessary.call_count,
        )
        self.assertIs(1, self.ensure_downloaded.call_count)
        self.assertEqual(
            [mock.call("canary_version", canary_path)],
            self.wrapped_upgrader.prepare_upgrade.call_args_list,
        )
        self.assertIs(0, self.wrapped_upgrader.upgrade.call_count)

    def test_force_upgrade(self):
        self.wrapped_upgrader.stable_version = "stable_version"
        self.wrapped_upgrader.canary_version = ""
        self.wrapped_upgrader.current_version = ""
        self.wrapped_upgrader.installed_versions = set()

        stable_path = pathlib.Path(self.image_directory, "stable_version")

        self.run_upgrade_loop()

        # Old-style upgrade method not called
        self.assertIs(
            0, self.wrapped_upgrader.perform_upgrade_if_necessary.call_count,
        )
        # New style methods called
        self.assertIs(2, self.ensure_downloaded.call_count)
        self.assertIs(1, self.wrapped_upgrader.get_upgrade_intent.call_count)
        self.assertIs(1, self.wrapped_upgrader.get_versions.call_count)
        self.assertEqual(
            [mock.call("stable_version", stable_path)],
            self.wrapped_upgrader.prepare_upgrade.call_args_list,
        )
        self.assertEqual(
            [mock.call("stable_version", stable_path)],
            self.wrapped_upgrader.upgrade.call_args_list,
        )
