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

    @classmethod
    def drop_all(cls):
        cls.metadata.drop_all()

    @classmethod
    def create_all(cls):
        cls.metadata.create_all()

    @classmethod
    def setMetadata(cls, metadata: MetaData = Base.metadata):
        cls.metadata = metadata

    @classmethod
    def setUpClass(cls) -> None:
        cls.setMetadata(metadata=Base.metadata)

    @classmethod
    def set_up_db_test_case(cls, **kwargs: Optional[Dict]):
        cls.engine = cls.get_test_db_engine(**kwargs)
        cls.session = Session(bind=cls.engine)
        cls.bind_engine()

    @staticmethod
    def get_test_db_engine(**kwargs) -> sqlalchemy.engine.Engine:
        config = TestConfig()
        return create_engine(
            url=kwargs.get("SQLALCHEMY_DB_URI", config.SQLALCHEMY_DB_URI),
            encoding=kwargs.get("SQLALCHEMY_DB_ENCODING", config.SQLALCHEMY_DB_ENCODING),
            echo=False,
            future=kwargs.get("SQLALCHEMY_FUTURE", config.SQLALCHEMY_FUTURE),
        )

    @classmethod
    def bind_engine(cls):
        cls.metadata.bind = cls.engine

    @classmethod
    def close_session(cls):
        cls.session.rollback()
        cls.session.close()


class BaseDBTestCase(DBTestCaseBlueprint):

    def setUp(self):
        self.set_up_db_test_case()
        self.create_all()

    def tearDown(self):
        self.close_session()
        self.drop_all()
