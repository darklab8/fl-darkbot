import pytest
from app import created_app
from storage import storage_builder


@pytest.fixture
def storage():
    return storage_builder()


def test_build_app(storage):
    bot = created_app(storage)
    assert bot is not None
