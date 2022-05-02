import os
from typing import List, Optional

import alembic.config
import sqlalchemy
import testing.postgresql
from alembic.command import downgrade, stamp, upgrade
from alembic.script import ScriptDirectory
from magma.db_service.tests.db_testcase import DBTestCaseBlueprint

HERE = os.path.abspath(os.path.dirname(__file__))
PROJECT_ROOT = os.path.join(HERE, os.pardir)


class AlembicTestCase(DBTestCaseBlueprint):
    postgresql: testing.postgresql.Postgresql
    up_revision: Optional[str] = None
    down_revision: Optional[str] = None
    tables: Optional[List[str]] = None

    def tearDown(self) -> None:
        super().tearDown()
        self.close_session()
        self.reflect_tables()
        self.drop_all()
        stamp(self.alembic_config, ())

    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()
        cls.postgresql = testing.postgresql.PostgresqlFactory(cache_initialized_db=True)()

    @classmethod
    def tearDownClass(cls) -> None:
        cls.postgresql.stop()

    def setUp(self) -> None:
        self.setMetadata(sqlalchemy.MetaData())
        self.set_up_db_test_case(SQLALCHEMY_DB_URI=self.postgresql.url())
        self.up_revision = self.up_revision or "head"
        self.down_revision = self.down_revision or "base"
        self.alembic_config = alembic.config.Config(os.path.join(PROJECT_ROOT, 'migrations/alembic.ini'))
        self.alembic_config.set_section_option('alembic', 'sqlalchemy.url', self.postgresql.url())
        self.alembic_config.set_main_option('script_location', os.path.join(PROJECT_ROOT, 'migrations'))
        self.script = ScriptDirectory.from_config(self.alembic_config)
        # Making sure there are no tables in metadata
        self.assertListEqual(self.metadata.sorted_tables, [])

    def get_table(self, table_name):
        return sqlalchemy.Table(table_name, sqlalchemy.MetaData(), autoload_with=self.engine)

    def upgrade(self, revision=None):
        revision = revision or self.up_revision
        upgrade(self.alembic_config, revision=revision)

    def downgrade(self, revision=None):
        revision = revision or self.down_revision
        downgrade(self.alembic_config, revision=revision)

    def reflect_tables(self):
        self.metadata.reflect()

    def close_session(self):
        self.session.rollback()
        self.session.close()

    def drop_all(self):
        self.metadata.drop_all()

    def then_tables_are(self, tables=None):
        if tables is None:
            tables = []
        self.reflect_tables()
        metadata_tables = dict(self.metadata.tables)
        del metadata_tables["alembic_version"]
        self.assertListEqual(sorted(tables), sorted(metadata_tables.keys()))
