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

CBSDS = 'cbsds'
TEST_STATE_ID = 1
NEW_COLUMNS = [
    'single_step_enabled',
    'cbsd_category',
    'latitude_deg',
    'longitude_deg',
    'height_m',
    'height_type',
    'horizontal_accuracy_m',
    'antenna_azimuth_deg',
    'antenna_downtilt_deg',
    'antenna_beamwidth_deg',
    'antenna_model',
    'eirp_capability_dbm_mhz',
    'cpi_digital_signature',
    'indoor_deployment',
]


class Test00d77c9f7532TestCase(AlembicTestCase):
    down_revision = 'e793671a19a6'
    up_revision = '00d77c9f7532'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)
        cbsd_states = self.get_table('cbsd_states')
        self.engine.execute(cbsd_states.insert().values(id=TEST_STATE_ID, name='some_state'))
        self.cbsds_pre_upgrade = self.get_table(CBSDS)
        self.engine.execute(
            self.cbsds_pre_upgrade.insert().values(
                id=1, state_id=TEST_STATE_ID, desired_state_id=TEST_STATE_ID,
            ),
        )

    def test_columns_not_present_pre_upgrade(self):
        for col in NEW_COLUMNS:
            self.assertFalse(self.has_column(self.cbsds_pre_upgrade, col))

    def test_columns_present_post_upgrade(self):
        # given
        self.upgrade()

        # when
        cbsds = self.get_table(CBSDS)

        # then
        self.assertTrue(self.has_columns(cbsds, NEW_COLUMNS))

    def test_default_values_post_upgrade(self):
        # given
        self.upgrade()

        # when
        cbsds = self.get_table(CBSDS)
        cbsd = self.engine.execute(cbsds.select()).fetchall()[0]

        # then
        self.assertFalse(cbsd.single_step_enabled)
        self.assertFalse(cbsd.indoor_deployment)
        self.assertEqual('b', cbsd.cbsd_category)

    def test_downgrade(self):
        # given
        self.upgrade()

        # when
        self.downgrade()
        for col in NEW_COLUMNS:
            self.assertFalse(self.has_column(self.cbsds_pre_upgrade, col))
