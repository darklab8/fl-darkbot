import pytest
from app import created_app, storage_builder


def test_build_app():
    STORAGE = storage_builder()
    bot = created_app(STORAGE)
    assert bot is not None
