import unittest
from typing import Dict, Optional

import sqlalchemy.engine
from magma.db_service.config import TestConfig
from magma.db_service.models import Base
from magma.db_service.session_manager import Session
from sqlalchemy import MetaData, create_engine


class DBTestCaseBlueprint(unittest.TestCase):
    metadata: MetaData
    engine: sqlalchemy.engine.Engine
    session: Session

    def drop_all(self):
        self.metadata.drop_all()

    def set_up_db_test_case(self, **kwargs: Optional[Dict]):
        self.engine = self.get_test_db_engine(**kwargs)
        self.session = Session(bind=self.engine)
        self.bind_engine()

    @staticmethod
    def get_test_db_engine(**kwargs) -> sqlalchemy.engine.Engine:
        config = TestConfig()
        return create_engine(
            url=kwargs.get("SQLALCHEMY_DB_URI") or config.SQLALCHEMY_DB_URI,
            encoding=kwargs.get("SQLALCHEMY_DB_ENCODING") or config.SQLALCHEMY_DB_ENCODING,
            echo=False,
            future=kwargs.get("SQLALCHEMY_FUTURE") or config.SQLALCHEMY_FUTURE,
        )

    def bind_engine(self):
        self.metadata.bind = self.engine

    def close_session(self):
        self.session.rollback()
        self.session.close()


class BaseDBTestCase(DBTestCaseBlueprint):

    def setUp(self):
        self.metadata = Base.metadata
        self.set_up_db_test_case()
        self.create_all()

    def tearDown(self):
        self.close_session()
        self.drop_all()

    def create_all(self):
        self.metadata.create_all()
