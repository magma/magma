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

TABLE = 'cbsds'
COLUMNS = [
    'antenna_beamwidth_deg', 'cpi_digital_signature',
    'horizontal_accuracy_m', 'antenna_model', 'eirp_capability_dbm_mhz',
    'antenna_azimuth_deg', 'antenna_downtilt_deg',
]


class TestRemoveCpiRelatedFields(AlembicTestCase):
    down_revision = 'fa12c537244a'
    up_revision = '37bd12af762a'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

    def test_upgrade(self):
        self.upgrade()
        table = self.get_table(TABLE)
        has = any(self.has_column(table, c) for c in COLUMNS)
        self.assertFalse(has)

    def test_columns_present_post_upgrade(self):
        self.upgrade()
        self.downgrade()
        table = self.get_table(TABLE)
        has = self.has_columns(table, COLUMNS)
        self.assertTrue(has)
