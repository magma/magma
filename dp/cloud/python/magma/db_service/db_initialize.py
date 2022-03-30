"""
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
"""

import importlib
import os

from magma.db_service.config import Config
from magma.db_service.models import DBCbsdState, DBGrantState, DBRequestType
from magma.db_service.session_manager import SessionManager
from magma.mappings.types import CbsdStates, GrantStates, RequestTypes
from sqlalchemy import create_engine


class DBInitializer(object):
    """
    This class is responsible for initializing the database with data.
    """

    def __init__(self, session_manager: SessionManager):
        self.session_manager = session_manager

    def initialize(self) -> None:
        """
        Initialize database with dictionary rows in dict tables
        """
        with self.session_manager.session_scope() as s:
            for request_type in RequestTypes:
                if not s.query(DBRequestType).filter(DBRequestType.name == request_type.value).first():
                    db_request_type = DBRequestType(name=request_type.value)
                    s.add(db_request_type)
            for state in GrantStates:
                if not s.query(DBGrantState).filter(DBGrantState.name == state.value).first():
                    grant_state = DBGrantState(name=state.value)
                    s.add(grant_state)
            for state in CbsdStates:
                if not s.query(DBCbsdState).filter(DBCbsdState.name == state.value).first():
                    cbsd_state = DBCbsdState(name=state.value)
                    s.add(cbsd_state)
            s.commit()


def get_config() -> Config:
    """
    Get configuration for db service
    """
    app_config = os.environ.get('APP_CONFIG', 'ProductionConfig')
    config_module = importlib.import_module(
        '.'.join(f"magma.db_service.config.{app_config}".split('.')[:-1]),
    )
    config_class = getattr(config_module, app_config.split('.')[-1])

    return config_class()


def main():
    config = get_config()
    db_engine = create_engine(
        url=config.SQLALCHEMY_DB_URI,
        encoding=config.SQLALCHEMY_DB_ENCODING,
        echo=config.SQLALCHEMY_ECHO,
        future=config.SQLALCHEMY_FUTURE,
        pool_size=config.SQLALCHEMY_ENGINE_POOL_SIZE,
        max_overflow=config.SQLALCHEMY_ENGINE_MAX_OVERFLOW,
    )
    session_manager = SessionManager(db_engine=db_engine)
    initializer = DBInitializer(session_manager)
    initializer.initialize()


if __name__ == '__main__':
    main()
