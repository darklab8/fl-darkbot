import pytest
from fastapi.testclient import TestClient

from .app import create_app


@pytest.fixture()
def app():
    app = create_app()

    # example how to override stuff
    # app.dependency_overrides[DatabaseFactory.get_default_database] = database.get_self

    return app


@pytest.fixture()
def client(app):
    # with-contextmanager is used in order to await `startup` event creating discord bot
    with TestClient(app) as client:
        yield client
