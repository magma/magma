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
from parameterized import parameterized
from sqlalchemy import select
from sqlalchemy.exc import NoSuchTableError

DOWN_REVISION = '37bd12af762a'
UP_REVISION = '467ad00fbc83'

CBSD_STATES_TABLE = 'cbsd_states'
CBSDS_TABLE = 'cbsds'
CHANNELS_TABLE = 'channels'

CBSDS_CHANNELS_COLUMN = 'channels'

TEST_STATE_ID = 1
CBSD_ID = 1
CHANNEL_ID = 1
CHANNEL_DATA = {
    "low_frequency": 1,
    "high_frequency": 2,
    "max_eirp": 3.0,
}
INCOMPLETE_CHANNEL_DATA = {
    "low_frequency": 1,
    "high_frequency": 2,
}
CHANNEL_DATA_WITH_DEFAULT_MAX_EIRP = {
    "low_frequency": 1,
    "high_frequency": 2,
    "max_eirp": 37,
}


class TestRemoveChannels(AlembicTestCase):
    down_revision = DOWN_REVISION
    up_revision = UP_REVISION

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

    def _given_cbsd_created(self, **data):
        cbsd_states = self.get_table(CBSD_STATES_TABLE)
        self.given_resource_inserted(cbsd_states, id=TEST_STATE_ID, name='some_state')

        cbsds = self.get_table(CBSDS_TABLE)
        self.given_resource_inserted(cbsds, **data)


class TestRemoveChannelsUpgrade(TestRemoveChannels):
    def test_table_removed(self):
        # Given
        self.assertFalse(self.has_column(self.get_table(CBSDS_TABLE), CBSDS_CHANNELS_COLUMN))
        self.get_table(CHANNELS_TABLE)

        # When
        self.upgrade()

        # Then
        self.assertTrue(self.has_column(self.get_table(CBSDS_TABLE), CBSDS_CHANNELS_COLUMN))
        with self.assertRaises(NoSuchTableError):
            self.get_table(CHANNELS_TABLE)

    def test_default_is_set(self):
        # Given
        self._given_cbsd_created(id=CBSD_ID, state_id=TEST_STATE_ID, desired_state_id=TEST_STATE_ID)

        # When
        self.upgrade()

        # Then
        cbsds = self.get_table(CBSDS_TABLE)
        data_after = self.engine.execute(cbsds.select()).mappings().one().get(CBSDS_CHANNELS_COLUMN)
        self.assertEqual(data_after, [])

    @parameterized.expand([
        (CHANNEL_DATA, CHANNEL_DATA),
        (INCOMPLETE_CHANNEL_DATA, CHANNEL_DATA_WITH_DEFAULT_MAX_EIRP),
    ])
    def test_data_migrated(self, channel_data, expected_channel_data):
        # Given
        self._given_cbsd_created(id=CBSD_ID, state_id=TEST_STATE_ID, desired_state_id=TEST_STATE_ID)

        channels = self.get_table(CHANNELS_TABLE)
        self.given_resource_inserted(
            channels, id=CHANNEL_ID, cbsd_id=CBSD_ID, channel_type='channel', rule_applied='rule', **channel_data,
        )

        # When
        self.upgrade()

        # Then
        cbsds = self.get_table(CBSDS_TABLE)
        data_after = self.engine.execute(cbsds.select()).mappings().one().get(CBSDS_CHANNELS_COLUMN)

        self.assertEqual([expected_channel_data], data_after, 'Data was not migrated')


class TestRemoveChannelsDowngrade(TestRemoveChannels):
    down_revision = DOWN_REVISION
    up_revision = UP_REVISION

    def setUp(self) -> None:
        super().setUp()
        self.upgrade()

    def test_table_added(self):
        # Given
        self.assertTrue(self.has_column(self.get_table(CBSDS_TABLE), CBSDS_CHANNELS_COLUMN))
        with self.assertRaises(NoSuchTableError):
            self.get_table(CHANNELS_TABLE)

        # When
        self.downgrade()

        # Then
        self.assertFalse(self.has_column(self.get_table(CBSDS_TABLE), CBSDS_CHANNELS_COLUMN))

        channels = self.get_table(CHANNELS_TABLE)
        self.has_columns(
            channels,
            ['id', 'cbsd_id', 'low_frequency', 'high_frequency', 'max_eirp', 'channel_type', 'rule_applied'],
        )

    def test_data_migrated(self):
        # Given
        self._given_cbsd_created(
            id=CBSD_ID, state_id=TEST_STATE_ID, desired_state_id=TEST_STATE_ID, channels=[CHANNEL_DATA],
        )

        # When
        self.downgrade()

        # Then
        data_before = [dict(id=CHANNEL_ID, cbsd_id=CBSD_ID, **CHANNEL_DATA)]

        channels = self.get_table(CHANNELS_TABLE)
        data_after = [
            dict(r) for r in self.engine.execute(
                select(
                    channels.c.id,
                    channels.c.cbsd_id,
                    channels.c.low_frequency,
                    channels.c.high_frequency,
                    channels.c.max_eirp,
                ),
            ).mappings().all()
        ]
        self.assertEqual(data_before, data_after, 'Data was not migrated')
