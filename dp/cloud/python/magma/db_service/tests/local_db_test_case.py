import testing.postgresql
from magma.db_service.models import Base
from magma.db_service.tests.db_testcase import BaseDBTestCase


class LocalDBTestCase(BaseDBTestCase):
    postgresql: testing.postgresql.Postgresql

    @classmethod
    def setUpClass(cls) -> None:
        super().setUpClass()
        cls.postgresql = testing.postgresql.PostgresqlFactory(cache_initialized_db=True)()

    @classmethod
    def tearDownClass(cls) -> None:
        cls.postgresql.stop()

    def setUp(self):
        self.set_up_db_test_case(SQLALCHEMY_DB_URI=self.postgresql.url())
        self.create_all()
