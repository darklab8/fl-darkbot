from typing import Generator
from sqlalchemy import create_engine
import sqlalchemy.orm as orm
from contextlib import contextmanager

import scrappy.core.settings as settings


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

        self._session_maker = orm.sessionmaker(
            autocommit=False, autoflush=False, bind=self._engine
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
    def manager_to_get_session(self) -> Generator[orm.Session, None, None]:
        session = self._session_maker()
        try:
            yield session
        finally:
            session.close()

    # Dependency
    def get_session(self) -> Generator[orm.Session, None, None]:
        session = self._session_maker()
        try:
            yield session
        finally:
            session.close()


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
    def get_default_session() -> Generator[orm.Session, None, None]:
        database = DatabaseFactory()
        yield database.get_session()
