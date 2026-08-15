"""
Copyright 2026 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

from unittest import TestCase, main

import build


class BuildTest(TestCase):
    """Verify the module mapping used by Orc8r image builds."""

    def test_orc8r_deployment_contains_one_module(self) -> None:
        """Represent the single Orc8r module as a one-item tuple."""
        self.assertEqual(
            build.DEPLOYMENT_TO_MODULES['orc8r'],
            ('orc8r',),
        )


if __name__ == '__main__':
    main()
