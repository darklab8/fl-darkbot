import pytest
from src.storage import Storage
from src.views import View
from src.data_model import DataModel
import copy


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


@pytest.fixture
@pytest.mark.asyncio
async def storage_with_tracked_base(storage, api_data):
    channel_id = 1
    base_to_track = ["Emlenton"]

    await storage.base.add(channel_id, base_to_track)
    base_data = await storage.base.get_data(channel_id)
    assert base_data == base_to_track
    print(base_data)

    return storage


@pytest.mark.asyncio
async def test_check_base_statuses_dynamically(storage_with_tracked_base,
                                               api_data):
    def request_api_data():
        return copy.deepcopy(api_data)

    channel_id = 1
    base_to_track = ["Emlenton"]

    def find_tracked_base(api_data_to_search):
        return ([
            value for key, value in api_data_to_search.bases.items()
            if key.startswith(base_to_track[0])
        ])[0]

    data_model = DataModel(request_api_data())
    data_model.update(request_api_data())
    # arrange
    rendered_date, rendered_base_data, _ = await View(
        data_model.api_data, storage_with_tracked_base,
        channel_id).render_all()

    assert "Static" in rendered_base_data

    print("=====================")
    new_api_data = request_api_data()
    find_tracked_base(new_api_data)["health"] += 0.1
    data_model.update(new_api_data)

    rendered_date, rendered_base_data, _ = await View(
        data_model.api_data, storage_with_tracked_base,
        channel_id).render_all()

    assert "Repairing" in rendered_base_data

    new_api_data = request_api_data()
    find_tracked_base(new_api_data)["health"] -= 0.1
    data_model.update(new_api_data)

    rendered_date, rendered_base_data, _ = await View(
        data_model.api_data, storage_with_tracked_base,
        channel_id).render_all()

    assert "Degrading" in rendered_base_data

    print(rendered_base_data)
