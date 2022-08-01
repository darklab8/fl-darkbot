import pytest
from sqlalchemy_utils import database_exists, create_database, drop_database

from fastapi.testclient import TestClient
from httpx import AsyncClient

from scrappy.core.databases import DatabaseFactory, Database
from scrappy.core.main import app_factory
import scrappy.core.settings as settings
from scrappy.core.declared_base import Model
from unittest.mock import patch
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
async def async_client(async_app):
    async with AsyncClient(app=async_app, base_url="http://test") as async_client:
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


@pytest.fixture(scope="session")
def celery_config():
    return {"broker_url": "memory://", "result_backend": "redis://redis:6379/0"}
