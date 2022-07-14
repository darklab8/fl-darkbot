import pytest
from sqlalchemy import create_engine
from sqlalchemy_utils import database_exists, create_database, drop_database

from fastapi.testclient import TestClient

import scrappy.core.databases as databases
from scrappy.core.main import app_factory
from scrappy.core.settings import get_database_url

@pytest.fixture()
def client():
    app = app_factory()
    client = TestClient(app)
    return client


@pytest.fixture
def db():
    database_url = get_database_url("test_database")
    engine = create_engine(database_url)
    if not database_exists(engine.url):
        create_database(engine.url)

    database = databases.Database(
        # url="sqlite:///./test_sql_app.db"
        url=database_url
    )

    databases.default.Base.metadata.drop_all(bind=database.engine)
    databases.default.Base.metadata.create_all(bind=database.engine)

    with database.manager_to_get_session() as session:
        yield session

    databases.default.Base.metadata.drop_all(bind=database.engine)


