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
COLUMN = 'grant_attempts'


class RemoveGrantAttemptsTestCase(AlembicTestCase):
    down_revision = '530b18568ad9'
    up_revision = '58a1b16ef73c'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

    def test_upgrade(self):
        self.upgrade()
        self.assertFalse(self.has_column(self.get_table(TABLE), COLUMN))

    def test_columns_present_post_upgrade(self):
        self.upgrade()
        self.downgrade()
        self.assertTrue(self.has_column(self.get_table(TABLE), COLUMN))
