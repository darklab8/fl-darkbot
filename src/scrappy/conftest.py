import pytest
from sqlalchemy import create_engine
from sqlalchemy_utils import database_exists, create_database

from fastapi.testclient import TestClient

import scrappy.core.databases as databases
from scrappy.core.main import app_factory
import scrappy.core.settings as settings

from unittest.mock import patch


@pytest.fixture()
def client():
    app = app_factory()
    client = TestClient(app)
    return client


@pytest.fixture
def database():
    test_database_name = "test_database"

    database_url = settings.DATABASE_URL + test_database_name

    engine = create_engine(database_url)
    if not database_exists(engine.url):
        create_database(engine.url)

    database = databases.Database(
        # url="sqlite:///./test_sql_app.db"
        url=settings.DATABASE_URL,
        name=test_database_name,
    )

    databases.default.Base.metadata.drop_all(bind=database.engine)
    databases.default.Base.metadata.create_all(bind=database.engine)

    with patch.object(
        databases.default, "_get_database_name", return_value=test_database_name
    ):
        yield database

    databases.default.Base.metadata.drop_all(bind=database.engine)


@pytest.fixture
def session(database):
    with database.manager_to_get_session() as session:
        yield session


@pytest.fixture(scope="session")
def celery_config():
    return {"broker_url": "memory://", "result_backend": "redis://redis:6379/0"}
