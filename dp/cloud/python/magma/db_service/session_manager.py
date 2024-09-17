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

from contextlib import contextmanager

from sqlalchemy.engine import Engine
from sqlalchemy.orm import Session as SQLAlchemy_Session
from sqlalchemy.orm import sessionmaker

Session = SQLAlchemy_Session


class SessionManager(object):
    """
    Database session manager class
    """

    def __init__(self, db_engine: Engine):
        self.session_factory = sessionmaker(bind=db_engine)

    @contextmanager
    def session_scope(self) -> Session:
        """
        Get database session

        Yields:
            Session: database session object

        Raises:
            Exception: generic exception
        """
        session = self.session_factory()
        try:
            yield session
        except Exception:
            session.rollback()
            raise
        finally:
            session.close()
