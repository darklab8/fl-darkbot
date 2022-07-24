from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from contextlib import contextmanager

import scrappy.core.settings as settings


class Database:
    DATABASE_NAME = "default"
    initialzed = False

    def _get_database_name(self):
        return self.DATABASE_NAME

    def _get_database_url(self):
        return self.DATABASE_URL + self._get_database_name()

    def __init__(self, url, name=None):
        self.DATABASE_URL = url

        if name is not None:
            self.DATABASE_NAME = name

        self.Base = declarative_base()

    @property
    def engine(self):
        self._delated_init()
        return self._engine

    def _delated_init(self):
        if self.initialzed:
            return

        if "postgresql" in self.DATABASE_URL:
            self._engine = create_engine(self._get_database_url(), pool_pre_ping=False)
        else:
            self._engine = create_engine(
                self._get_database_url(), connect_args={"check_same_thread": False}
            )
        self.SessionLocal = sessionmaker(
            autocommit=False, autoflush=False, bind=self._engine
        )

        self.initialzed = True

    @contextmanager
    def manager_to_get_session(self):
        self._delated_init()
        db = self.SessionLocal()
        try:
            yield db
        finally:
            db.close()

    # Dependency
    def get_session(self):
        self._delated_init()
        db = self.SessionLocal()
        try:
            yield db
        finally:
            db.close()


default = Database(
    # url="sqlite:///./sql_app.db"
    url=settings.DATABASE_URL,
    name=settings.DATABASE_NAME,
)
