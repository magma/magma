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
