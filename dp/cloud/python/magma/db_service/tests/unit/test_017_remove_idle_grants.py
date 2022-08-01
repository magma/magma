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
import sqlalchemy as sa
from magma.db_service.tests.alembic_testcase import AlembicTestCase

GRANTS_TABLE = 'grants'
GRANT_STATES_TABLE = 'grant_states'
IDLE_STATE_NAME = 'idle'


class RemoveIdleGrantsTestCase(AlembicTestCase):
    down_revision = '58a1b16ef73c'
    up_revision = '98f7ccfbd2f8'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

        self._grants_table = self.get_table(GRANTS_TABLE)
        self._grant_states_table = self.get_table(GRANT_STATES_TABLE)

    def test_upgrade_removes_idle_grants_and_state(self) -> None:
        # Given
        self._create_idle_grant_state()
        self._create_idle_grant()

        # When
        self.upgrade()

        # Then
        grant_state = self._get_idle_grant_state()
        self.assertIsNone(grant_state)

        grant = self._get_idle_grant()
        self.assertIsNone(grant)

    def test_downgrade_creates_idle_state(self) -> None:
        # Given
        grant_state = self._get_idle_grant_state()
        self.assertIsNone(grant_state)

        self.upgrade()

        # When
        self.downgrade()

        # Then
        grant_state = self._get_idle_grant_state()
        self.assertEqual(grant_state.name, IDLE_STATE_NAME)

    def _create_idle_grant_state(self) -> None:
        self.engine.execute(
            self._grant_states_table.insert().values(name=IDLE_STATE_NAME),
        )

    def _get_idle_grant_state(self) -> sa.engine.Row:
        return self.engine.execute(self._grant_states_table.select()).first()

    def _create_idle_grant(self) -> None:
        idle_grant_state = self._get_idle_grant_state()
        self.engine.execute(
            self._grants_table.insert().values(
                grant_id="some_grant_id",
                state_id=idle_grant_state.id,
                low_frequency=3550_000_000,
                high_frequency=3570_000_000,
                max_eirp=20,
            ),
        )

    def _get_idle_grant(self) -> sa.engine.Row:
        return self.engine.execute(self._grants_table.select()).first()
