from src.storage import Storage
import pytest


@pytest.fixture
def storage():
    return Storage()


@pytest.mark.asyncio
async def test_storage_1(storage):
    updating_api_data = await storage.a_get_game_data()
    updating_new_forum_records = await storage.a_get_new_forum_records({})