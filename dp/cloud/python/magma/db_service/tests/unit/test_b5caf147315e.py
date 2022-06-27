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
from magma.db_service.tests.alembic_testcase import AlembicTestCase


class Testeb5caf147315eTestCase(AlembicTestCase):
    down_revision = '00d77c9f7532'
    up_revision = 'b5caf147315e'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

        cbsd_states = self.get_table('cbsd_states')
        self.engine.execute(cbsd_states.insert().values(id=1, name='some_state'))
        self.cbsds_pre_upgrade = self.get_table('cbsds')

    def test_upgrade(self):
        # given
        self.engine.execute(self.cbsds_pre_upgrade.insert().values(id=1, state_id=1, desired_state_id=1))
        self.upgrade()

        # when
        cbsds = self.get_table('cbsds')
        cbsd = self.engine.execute(cbsds.select()).fetchall()[0]

        # then
        self.assertFalse(cbsd.carrier_aggregation_enabled)
        self.assertTrue(cbsd.grant_redundancy)
        self.assertEqual(150, cbsd.max_ibw_mhz)

    def test_downgrade(self):
        # given
        self.upgrade()

        cbsds = self.get_table('cbsds')
        self.given_resource_inserted(cbsds, id=1, state_id=1, desired_state_id=1)

        # when
        self.downgrade()
        cbsds = self.get_table('cbsds')
        cbsd = self.engine.execute(cbsds.select()).fetchall()[0]

        # then
        self.then_column_does_not_exist(cbsd, 'carrier_aggregation_enabled')
        self.then_column_does_not_exist(cbsd, 'max_ibw_mhz')
        self.then_column_does_not_exist(cbsd, 'grant_redundancy')

    def test_downgrade_to_base(self):
        self.verify_downgrade_to_base()
