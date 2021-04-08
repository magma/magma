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

# pylint: disable=protected-access
from unittest import TestCase

from magma.enodebd.data_models.transform_for_magma import bandwidth, gps_tr181
from magma.enodebd.exceptions import ConfigurationError


class TransformForMagmaTests(TestCase):
    def test_gps_tr181(self) -> None:
        # Negative longitude
        inp = '-122150583'
        out = gps_tr181(inp)
        expected = '-122.150583'
        self.assertEqual(out, expected, 'Should convert negative longitude')

        inp = '122150583'
        out = gps_tr181(inp)
        expected = '122.150583'
        self.assertEqual(out, expected, 'Should convert positive longitude')

        inp = '0'
        out = gps_tr181(inp)
        expected = '0.0'
        self.assertEqual(out, expected, 'Should leave zero as zero')

    def test_bandwidth(self) -> None:
        inp = 'n6'
        out = bandwidth(inp)
        expected = 1.4
        self.assertEqual(out, expected, 'Should convert RBs')

        inp = 1.4
        out = bandwidth(inp)
        expected = 1.4
        self.assertEqual(out, expected, 'Should accept MHz')

        with self.assertRaises(ConfigurationError):
            inp = 'asdf'
            bandwidth(inp)

        with self.assertRaises(ConfigurationError):
            inp = 1234
            bandwidth(inp)
