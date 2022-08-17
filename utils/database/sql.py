from typing import Generator, AsyncGenerator
from contextlib import contextmanager, asynccontextmanager

from sqlalchemy import create_engine
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
import sqlalchemy.orm as orm
from sqlalchemy.orm import Session
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import DeclarativeMeta
from typing import Callable


class ORMBase:
    def __new__(cls) -> Callable[[], DeclarativeMeta]:
        return declarative_base()


class Database:
    def __init__(self, name, url):
        self._name = name
        self._url = url

        self._engine = create_engine(self.full_url, pool_pre_ping=False)

        self._async_engine = create_async_engine(
            self.async_full_url,
            future=True,
            pool_size=20,
            pool_pre_ping=True,
            pool_use_lifo=True,
        )

    @property
    def url(self) -> str:
        return self._url

    @property
    def name(self) -> str:
        return self._name

    @property
    def full_url(self) -> str:
        return "postgresql://" + self._url + self._name

    @property
    def async_full_url(self) -> str:
        return "postgresql+asyncpg://" + self._url + self._name

    @property
    def engine(self):
        return self._engine

    @property
    def async_engine(self):
        return self._async_engine

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
            yield session

    @asynccontextmanager
    async def get_async_session(self) -> AsyncGenerator[AsyncSession, None]:
        try:
            connection = AsyncSession(self.async_engine)
            yield connection
        finally:
            await connection.close()

    def get_self(self):
        return self


class DatabaseFactoryBase:
    def __new__(
        cls,
        url: str | None = None,
        name: str | None = None,
    ) -> Database:
        instance = Database(
            url=cls.settings.DATABASE_URL if url is None else url,
            name=cls.settings.DATABASE_NAME if name is None else name,
        )
        return instance

    @classmethod
    def get_default_database(cls) -> Database:
        database = cls()
        return database
