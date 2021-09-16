import testing.postgresql
from magma.db_service.tests.db_testcase import DBTestCase

Postgresql = testing.postgresql.PostgresqlFactory(cache_initialized_db=True)


class LocalDBTestCase(DBTestCase):
    postgresql: testing.postgresql.Postgresql

    @classmethod
    def setUpClass(cls) -> None:
        cls.postgresql = Postgresql()

    @classmethod
    def tearDownClass(cls) -> None:
        cls.postgresql.stop()

    def get_config(self):
        config = super().get_config()
        config.SQLALCHEMY_DB_URI = self.postgresql.url()
        return config
