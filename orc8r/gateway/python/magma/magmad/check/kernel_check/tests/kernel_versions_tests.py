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

import textwrap
import unittest

from magma.magmad.check.kernel_check import kernel_versions

# Allow access to protected variables for unit testing
# pylint: disable=protected-access


class DpkgArgFactoryTests(unittest.TestCase):
    def test_function(self):
        actual = kernel_versions._get_dpkg_command_args_list(
            kernel_versions.DpkgCommandParams(),
        )
        self.assertEqual(
            ['dpkg', '--list'],
            actual,
        )


class DpkgParseTests(unittest.TestCase):

    def setUp(self):
        self.param = kernel_versions._get_dpkg_command_args_list(
            kernel_versions.DpkgCommandParams(),
        )

    def test_parse_with_errors(self):
        param = kernel_versions.DpkgCommandParams()
        actual = kernel_versions.parse_dpkg_output('test', 'test', param)
        expected = kernel_versions.DpkgCommandResult(
            error='test',
            kernel_versions_installed=None,
        )
        self.assertEqual(expected, actual)

    def test_parse_good_output(self):
        output = textwrap.dedent('''
            ii  console-setup-linux           1.164             all
            ii  firmware-linux-free           3.4               all
            ii  libselinux1:amd64             2.6-3+b3          amd64
            ii  linux-base                    4.5               all
            ii  linux-compiler-gcc-6-x86      4.9.110-1         amd64
            ii  linux-headers-4.9.0-4-amd64   4.9.65-3+deb9u1   amd64
            ii  linux-headers-4.9.0-4-common  4.9.65-3+deb9u1   all
            ii  linux-headers-4.9.0-6-amd64   4.9.88-1+deb9u1   amd64
            ii  linux-headers-4.9.0-6-common  4.9.88-1+deb9u1   all
            ii  linux-headers-4.9.0-7-amd64   4.9.110-1         amd64
            ii  linux-headers-4.9.0-7-common  4.9.110-1         all
            ii  linux-headers-amd64           4.9+80+deb9u5     amd64
            ii  linux-image-4.9.0-4-amd64     4.9.65-3+deb9u1   amd64
            ii  linux-image-4.9.0-6-amd64     4.9.88-1+deb9u1   amd64
            ii  linux-image-4.9.0-7-amd64     4.9.110-3+deb9u2  amd64
            ii  linux-image-amd64             4.9+80+deb9u2     amd64
            ii  linux-kbuild-4.9              4.9.110-1         amd64
            ii  linux-libc-dev:amd64          4.9.110-1         amd64
            ii  util-linux                    2.29.2-1+deb9u1   amd64
            ii  util-linux-locales            2.29.2-1+deb9u1   all
        ''').strip().encode('ascii')

        expected = [
            'linux-image-4.9.0-4-amd64',
            'linux-image-4.9.0-6-amd64',
            'linux-image-4.9.0-7-amd64',
            'linux-image-amd64',
        ]
        actual = kernel_versions.parse_dpkg_output(output, '', self.param)
        self.assertEqual(expected, actual.kernel_versions_installed)


if __name__ == '__main__':
    unittest.main()
