"""
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""
from typing import Any

from magma.db_service.tests.alembic_testcase import AlembicTestCase
from sqlalchemy import Table, select


class Teste793671a19a6TestCase(AlembicTestCase):
    down_revision = '9cd338f28663'
    up_revision = 'e793671a19a6'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

        cbsd_states = self.get_table('cbsd_states')
        self.given_resource_inserted(cbsd_states, id=1, name='unregistered')
        self.given_resource_inserted(cbsd_states, id=2, name='registered')

    def test_upgrade(self):
        cbsds = self.get_table('cbsds')
        amcs = self.get_table('active_mode_configs')
        self.given_resource_inserted(cbsds, id=1, state_id=1)
        self.given_resource_inserted(cbsds, id=2, state_id=1)
        self.given_resource_inserted(cbsds, id=3, state_id=1)
        self.given_resource_inserted(amcs, id=1, cbsd_id=2, desired_state_id=1)
        self.given_resource_inserted(amcs, id=2, cbsd_id=3, desired_state_id=2)

        self.upgrade(self.up_revision)

        cbsds = self.get_table('cbsds')
        actual = self.engine.execute(select(cbsds.c.id, cbsds.c.desired_state_id)).fetchall()
        expected = [(1, 1), (2, 1), (3, 2)]
        self.assertEqual(expected, actual)

    def test_downgrade(self):
        self.upgrade(self.up_revision)

        cbsds = self.get_table('cbsds')
        self.given_resource_inserted(cbsds, id=1, state_id=1, desired_state_id=1)
        self.given_resource_inserted(cbsds, id=2, state_id=1, desired_state_id=2)

        self.downgrade(self.down_revision)

        amcs = self.get_table('active_mode_configs')
        actual = self.engine.execute(select(amcs.c.cbsd_id, amcs.c.desired_state_id)).fetchall()
        expected = [(1, 1), (2, 2)]
        self.assertEqual(expected, actual)
