from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker
from contextlib import contextmanager


class Database:
    def __init__(self, url):
        self.SQLALCHEMY_DATABASE_URL = url
        # SQLALCHEMY_DATABASE_URL = "postgresql://user:password@postgresserver/db"

        self.engine = create_engine(
            self.SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False}
        )
        self.SessionLocal = sessionmaker(
            autocommit=False, autoflush=False, bind=self.engine
        )

        self.Base = declarative_base()

    @contextmanager
    def manager_to_get_db(self):
        db = self.SessionLocal()
        try:
            yield db
        finally:
            db.close()

    # Dependency
    def get_db(self):
        db = self.SessionLocal()
        try:
            yield db
        finally:
            db.close()


default = Database(url="sqlite:///./sql_app.db")
