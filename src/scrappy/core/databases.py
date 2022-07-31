from typing import Generator
from sqlalchemy import create_engine
import sqlalchemy.orm as orm
from contextlib import contextmanager
from sqlalchemy.orm import Session

import scrappy.core.settings as settings


class SessionWrapper:
    def __init__(self, session: Session):
        self._session = session

    def execute(self, stmt):
        return self._session.execute(stmt)


class Database:
    def __init__(self, name, url):
        self._name = name
        self._url = url

        if "postgresql" in self.url:
            self._engine = create_engine(self.full_url, pool_pre_ping=False)
        else:
            self._engine = create_engine(
                self.full_url,
                connect_args={"check_same_thread": False},
            )

    @property
    def url(self) -> str:
        return self._url

    @property
    def name(self) -> str:
        return self._name

    @property
    def full_url(self) -> str:
        return self._url + self._name

    @property
    def engine(self):
        return self._engine

    @contextmanager
    def get_orm_sessiom(self) -> Generator[orm.Session, None, None]:

        session_maker = orm.sessionmaker(
            autocommit=False, autoflush=False, bind=self._engine
        )

        session = session_maker()
        try:
            yield session
        finally:
            session.close()

    @contextmanager
    def get_core_session(self) -> Generator[Session, None, None]:
        with Session(self.engine, future=True) as session:
            yield SessionWrapper(session)

    def get_self(self):
        return self


class DatabaseFactory:
    def __new__(
        cls,
        url: str = settings.DATABASE_URL,
        name: str = settings.DATABASE_NAME,
    ) -> Database:
        instance = Database(
            url=url,
            name=name,
        )
        return instance

    @staticmethod
    def get_default_database() -> Database:
        database = DatabaseFactory()
        return database
