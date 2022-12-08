from dataclasses import dataclass
from typing import Generator, AsyncGenerator
from contextlib import contextmanager, asynccontextmanager

from sqlalchemy import create_engine
from sqlalchemy.ext.asyncio import AsyncSession, create_async_engine
import sqlalchemy.orm as orm
from sqlalchemy.orm import Session
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import DeclarativeMeta
from typing import Any
from types import ModuleType
from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import AsyncEngine
import abc


class ORMBase:
    def __new__(cls) -> Any:
        return declarative_base()


class Database:
    def __init__(self, name: str, url: str, debug: bool):
        self._name = name
        self._url = url

        self._engine = create_engine(
            self.full_url,
            pool_pre_ping=False,
            echo=debug,
        )

        self._async_engine = create_async_engine(
            self.async_full_url,
            future=True,
            pool_size=20,
            pool_pre_ping=True,
            pool_use_lifo=True,
            echo=debug,
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
    def engine(self) -> Engine:
        return self._engine

    @property
    def async_engine(self) -> AsyncEngine:
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

    def get_self(self) -> "Database":
        return self


@dataclass
class SettingStub:
    DATABASE_URL: str = "example"
    DATABASE_NAME: str = "example"
    DATABASE_DEBUG: bool = False


class DatabaseFactoryBase(Database):
    @property
    def settings(self) -> SettingStub:
        raise NotImplementedError("define settings")

    def __init__(
        cls,
        url: str | None = None,
        name: str | None = None,
        debug: bool | None = None,
    ):
        super().__init__(
            url=cls.settings.DATABASE_URL if url is None else url,
            name=cls.settings.DATABASE_NAME if name is None else name,
            debug=cls.settings.DATABASE_DEBUG if debug is None else debug,
        )

    @classmethod
    def get_default_database(cls) -> Database:
        database: Database = cls()
        return database
