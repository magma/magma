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

import unittest
from collections import namedtuple
from unittest import mock
from unittest.mock import patch

# Allow unused arguments in DIY mock functions
# pylint: disable=unused-argument
from magma.magmad.upgrade import magma_upgrader


class VersionComparisonTests(unittest.TestCase):

    def test_function(self):
        TestCase = namedtuple('TestCase', ['v1', 'v2'])
        greater_test_cases = [
            TestCase('1.0.0-0', '2.0.0-0'),
            TestCase('1.0.0-0', '1.1.0-0'),
            TestCase('1.0.0-0', '1.0.1-0'),
            TestCase('1.0.0-0', '1.0.0-1'),
            TestCase('1.0.0', '1.0.0-1'),
        ]
        equal_test_case = TestCase('1.1.1-0', '1.1.1-0')
        less_test_cases = [
            TestCase('2.0.0-0', '1.0.0-0'),
            TestCase('1.1.0-0', '1.0.0-0'),
            TestCase('1.0.1-0', '1.0.0-0'),
            TestCase('1.0.0-1', '1.0.0-0'),
            TestCase('1.0.0-1', '1.0.0'),
        ]

        for test_case in greater_test_cases:
            self.assertEqual(
                1,
                magma_upgrader.compare_package_versions(
                    test_case.v1,
                    test_case.v2,
                ),
            )
        for test_case in less_test_cases:
            self.assertEqual(
                -1,
                magma_upgrader.compare_package_versions(
                    test_case.v1,
                    test_case.v2,
                ),
            )
        self.assertEqual(
            0,
            magma_upgrader.compare_package_versions(
                equal_test_case.v1,
                equal_test_case.v2,
            ),
        )

    def test_bad_version_strings(self):
        test_cases = [
            '1.0-0',
            '1-0',
            '1',
            'a.b.c-d',
        ]
        for ver in test_cases:
            try:
                magma_upgrader.compare_package_versions('1.0.0-0', ver)
                self.fail('An exception should have been raised')
            except ValueError as e:
                self.assertEqual(
                    'Could not parse target package version {}'.format(ver),
                    str(e),
                )

            try:
                magma_upgrader.compare_package_versions(ver, '1.0.0-0')
                self.fail('An exception should have been raised')
            except ValueError as e:
                self.assertEqual(
                    'Could not parse current package version {}'.format(ver),
                    str(e),
                )


class MagmaUpgraderTests(unittest.TestCase):

    @staticmethod
    def _mock_call_process_success(cmd, callback, loop):
        callback(0)

    @staticmethod
    def _mock_call_process_fail(cmd, callback, loop):
        callback(1)

    @staticmethod
    def _mock_call_process_fail_on_install(cmd, callback, loop):
        if 'install' in cmd:
            callback(1)
        else:
            callback(0)

    @staticmethod
    def _get_upgrader(cur_ver):
        return magma_upgrader.MagmaUpgrader(cur_ver, mock.Mock())

    def test_full_upgrade_workflow(self):
        with patch(
                'magma.magmad.upgrade.magma_upgrader.call_process',
                mock.Mock(wraps=self._mock_call_process_success),
        ) as m:
            upgrader = self._get_upgrader('1.0.0-0')
            upgrader.perform_upgrade_if_necessary('1.1.0-0')

            self.assertEqual(3, m.call_count)
            m.assert_has_calls([
                mock.call('apt-get update', mock.ANY, mock.ANY),
                mock.call(
                    'apt-get install '
                    '-o Dpkg::Options::="--force-confnew" '
                    '--assume-yes --force-yes --only-upgrade --dry-run '
                    'magma=1.1.0-0',
                    mock.ANY, mock.ANY,
                ),
                mock.call(
                    'apt-get install '
                    '-o Dpkg::Options::="--force-confnew" '
                    '--assume-yes --force-yes --only-upgrade '
                    'magma=1.1.0-0',
                    mock.ANY, mock.ANY,
                ),
            ])

    def test_apt_update_failure(self):
        with patch(
                'magma.magmad.upgrade.magma_upgrader.call_process',
                mock.Mock(wraps=self._mock_call_process_fail),
        ) as m:
            upgrader = self._get_upgrader('1.0.0-0')
            upgrader.perform_upgrade_if_necessary('1.1.0-0')

            self.assertEqual(1, m.call_count)
            m.assert_has_calls([
                mock.call('apt-get update', mock.ANY, mock.ANY),
            ])

    def test_upgrade_dry_run_failure(self):
        with patch(
                'magma.magmad.upgrade.magma_upgrader.call_process',
                mock.Mock(wraps=self._mock_call_process_fail_on_install),
        ) as m:
            upgrader = self._get_upgrader('1.0.0-0')
            upgrader.perform_upgrade_if_necessary('1.1.0-0')

            self.assertEqual(2, m.call_count)
            m.assert_has_calls([
                mock.call('apt-get update', mock.ANY, mock.ANY),
                mock.call(
                    'apt-get install '
                    '-o Dpkg::Options::="--force-confnew" '
                    '--assume-yes --force-yes --only-upgrade --dry-run '
                    'magma=1.1.0-0',
                    mock.ANY, mock.ANY,
                ),
            ])

    def test_already_in_newest_version(self):
        with patch(
                'magma.magmad.upgrade.magma_upgrader.call_process',
                mock.Mock(wraps=self._mock_call_process_success),
        ) as m:
            upgrader = self._get_upgrader('1.1.0-0')
            upgrader.perform_upgrade_if_necessary('1.1.0-0')
            self.assertEqual(0, m.call_count)

    def test_downgrade_is_ignored(self):
        with patch(
                'magma.magmad.upgrade.magma_upgrader.call_process',
                mock.Mock(wraps=self._mock_call_process_success),
        ) as m:
            upgrader = self._get_upgrader('1.2.0-0')
            upgrader.perform_upgrade_if_necessary('1.1.0-0')

            self.assertEqual(0, m.call_count)
