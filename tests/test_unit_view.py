import pytest
from src.storage import Storage
from src.views import View


@pytest.fixture
def storage():
    return Storage()


@pytest.fixture
def api_data(storage):
    return storage.get_load_test_game_data({})


@pytest.mark.asyncio
async def test_rendering_views(storage, api_data):
    channel_ids = [int(item) for item in storage.channels.keys()]
    print(len(channel_ids))
    for channel_id in channel_ids:
        rendered_date, rendered_all, render_forum_records = await View(
            api_data, storage, channel_id).render_all()
