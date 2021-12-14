# can be run only with `pytest -m integration`
import pytest
import os
import asyncio
from src.app import CreateApp


@pytest.fixture
@pytest.mark.asyncio
@pytest.mark.integration
@pytest.mark.discord
async def app():
    app = CreateApp()
    yield app
    await app.looper.cog_unload()


@pytest.mark.asyncio
@pytest.mark.integration
@pytest.mark.discord
async def test_app_can_run(app):
    loop = asyncio.get_event_loop()
    loop.create_task(app.bot.start(app.storage.settings.discord_bot_key))
    await app.bot.wait_until_ready()
    await asyncio.sleep(int(os.environ.get('testing_time', '5')))
    await app.bot.close()
