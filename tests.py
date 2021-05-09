import pytest
import os
from app import CreateApp
from views import View
from channel import ChannelConstroller
from mocks import MockDiscordMessageBus
import asyncio


@pytest.fixture
@pytest.mark.asyncio
async def app():
    app = CreateApp()
    yield app
    await app.looper.cog_unload()


@pytest.mark.skipif(not bool(os.environ.get("pipline")), reason="long test")
@pytest.mark.asyncio
async def test_app_can_run(app):
    loop = asyncio.get_event_loop()
    loop.create_task(app.bot.start(app.storage.settings.secret_key))
    await app.bot.wait_until_ready()
    await asyncio.sleep(int(os.environ.get('testing_time', '0')))
    await app.bot.close()


@pytest.fixture
def data(app):
    game_data = app.storage.get_game_data()

    assert len(game_data.players) > 0
    assert len(game_data.bases) > 0

    return game_data


def test_saving_correctly_and_loading_back(app):
    app.storage.save_channel_settings()
    app.storage.get_game_data()


@pytest.fixture
@pytest.mark.asyncio
async def filled_data(app, data):
    for i, item in enumerate(data.bases.keys()):
        await app.storage.base.add(1, (item, ))
        if i > 10:
            break
    # bases = await storage.base.get_data(1)

    for i, item in enumerate(data.players['players']):
        await app.storage.system.add(1, (item['system'], ))
        if i > 10:
            break

    return data


@pytest.mark.asyncio
async def test_render_all(app, filled_data):

    rendered_date, rendered_all = await View(filled_data, app.storage,
                                             1).render_all()
    assert len(rendered_date) > 50
    assert len(rendered_all) > 50


@pytest.mark.asyncio
async def test_check_channel_controller(mocker):
    controller = ChannelConstroller(MockDiscordMessageBus(),
                                    'self.storage.unique_tag')
    await controller.delete_exp_msgs(1, 40)

    messages = await controller.get_tagged_msgs(1)
    for message in messages:
        controller.message_bus.delete(message)

    await controller.update_info(1, 'info')
