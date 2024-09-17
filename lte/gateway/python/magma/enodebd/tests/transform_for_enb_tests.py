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

from magma.enodebd.data_models.transform_for_enb import bandwidth


class TransformForMagmaTests(TestCase):
    def test_bandwidth(self) -> None:
        inp = 1.4
        out = bandwidth(inp)
        expected = 'n6'
        self.assertEqual(out, expected, 'Should work with a float')

        inp = 20
        out = bandwidth(inp)
        expected = 'n100'
        self.assertEqual(out, expected, 'Should work with an int')

        inp = 10
        out = bandwidth(inp)
        expected = 'n50'
        self.assertEqual(out, expected, 'Should work with int 10')
