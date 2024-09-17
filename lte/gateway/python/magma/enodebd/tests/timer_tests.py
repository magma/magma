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

from magma.enodebd.state_machines.timer import StateMachineTimer


class StateMachineTimerTests(TestCase):
    def test_is_done(self):
        timer_a = StateMachineTimer(0)
        self.assertTrue(timer_a.is_done(), 'Timer should be done')

        timer_b = StateMachineTimer(600)
        self.assertFalse(timer_b.is_done(), 'Timer should not be done')
