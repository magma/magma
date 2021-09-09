from contextlib import contextmanager

from sqlalchemy.engine import Engine
from sqlalchemy.orm import Session as sqlalchemy_session
from sqlalchemy.orm import sessionmaker

Session = sqlalchemy_session


class SessionManager:
    def __init__(self, db_engine: Engine):
        self.session_factory = sessionmaker(bind=db_engine)

    @contextmanager
    def session_scope(self) -> Session:
        session = self.session_factory()
        try:
            yield session
        except Exception:
            session.rollback()
            raise
        finally:
            session.close()
