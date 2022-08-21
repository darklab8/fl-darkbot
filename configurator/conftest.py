import pytest
from sqlalchemy_utils import database_exists, create_database, drop_database

from fastapi.testclient import TestClient
from httpx import AsyncClient

from configurator.core.databases import DatabaseFactory
from utils.database.sql import Database
from configurator.core.main import app_factory
import configurator.core.settings as settings
from configurator.core.declared_base import Model
import secrets


@pytest.fixture()
def app(database: Database):
    app = app_factory()
    app.dependency_overrides[DatabaseFactory.get_default_database] = database.get_self
    return app


@pytest.fixture()
def client(app):
    client = TestClient(app)
    return client


@pytest.fixture()
async def async_client(app):
    async with AsyncClient(app=app, base_url="http://test") as async_client:
        yield async_client


@pytest.fixture()
def database():
    database = DatabaseFactory(
        url=settings.DATABASE_URL,
        name=f"test_database_{secrets.token_hex(10)}",
    )

    if not database_exists(database.full_url):
        create_database(database.full_url)

    Model.metadata.drop_all(bind=database.engine)

    Model.metadata.create_all(bind=database.engine)
    yield database
    Model.metadata.drop_all(bind=database.engine)
    if database_exists(database.full_url):
        drop_database(database.full_url)


@pytest.fixture
def session(database):
    with database.manager_to_get_session() as session:
        yield session
