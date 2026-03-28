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

from typing import Dict, List

from magma.db_service.tests.alembic_testcase import AlembicTestCase


class TestAddIndicesForRelations(AlembicTestCase):
    down_revision = '98f7ccfbd2f8'
    up_revision = '48e8b58fcc24'

    def setUp(self) -> None:
        super().setUp()
        self.upgrade(self.down_revision)

    def test_upgrade(self):
        self.then_indexed_columns_are({
            'channels': [],
            'requests': [],
            'grants': [],
        })
        self.upgrade()
        self.then_indexed_columns_are({
            'channels': ['cbsd_id'],
            'requests': ['cbsd_id'],
            'grants': ['cbsd_id'],
        })

    def test_downgrade(self):
        self.upgrade()
        self.downgrade()
        self.then_indexed_columns_are({
            'channels': [],
            'requests': [],
            'grants': [],
        })

    def then_indexed_columns_are(self, indexes_dict: Dict[str, List[str]]):
        for tab, indexes in indexes_dict.items():
            table = self.get_table(tab)

            indexed_columns = [i.expressions[0].name for i in table.indexes]
            self.assertCountEqual(indexed_columns, indexes)
