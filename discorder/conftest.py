import pytest
from fastapi.testclient import TestClient
from httpx import AsyncClient

from .app import create_app


@pytest.fixture(scope="session")
def app():
    app = create_app()

    # example how to override stuff
    # app.dependency_overrides[DatabaseFactory.get_default_database] = database.get_self

    return app


@pytest.fixture(scope="session")
def client(app):
    # with-contextmanager is used in order to await `startup` event creating discord bot
    with TestClient(app) as client:
        yield client


# @pytest_asyncio.fixture(scope="module")
@pytest.fixture(scope="session")
async def async_client(app):
    async with AsyncClient(app=app, base_url="http://test") as async_client:
        yield async_client


@pytest.fixture()
def channel_id():
    return 840251517398548521
