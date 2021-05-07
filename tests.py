import pytest
from commands import attach_commands
from storage import Storage
from looper import Looper
from views import View


@pytest.fixture
def storage():
    return Storage()


@pytest.fixture
def bot(storage):
    return attach_commands(storage)


def looper(bot):
    _ = Looper(bot, storage)


@pytest.fixture
def data(storage):
    game_data = storage.get_game_data()

    assert len(game_data.players) > 0
    assert len(game_data.bases) > 0

    return game_data


def test_saving_correctly_and_loading_back(storage, data):
    storage.save_channel_settings()
    storage.get_game_data()


@pytest.mark.asyncio
async def test_render_all(storage, data):

    for i, item in enumerate(data.bases.keys()):
        await storage.base.add(1, (item, ))
        if i > 10:
            break
    # bases = await storage.base.get_data(1)

    for i, item in enumerate(data.players['players']):
        await storage.system.add(1, (item['system'], ))
        if i > 10:
            break

    rendered_date, rendered_all = await View(data, storage, 1).render_all()
    assert len(rendered_date) > 50
    assert len(rendered_all) > 50
