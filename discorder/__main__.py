from asyncio import gather, get_event_loop
from logging import basicConfig, INFO
from sys import argv
from aiohttp.web import AppRunner, Application, TCPSite
from discorder.core import settings
from discorder.bot import Singleton
from discorder.api import routes

basicConfig(level=INFO)


async def run_bot():

    app = Application()
    app.add_routes(routes)

    runner = AppRunner(app)
    await runner.setup()
    site = TCPSite(runner, "0.0.0.0", 8080)
    await site.start()

    # app["bot"] = Singleton().bot

    try:
        try:
            await Singleton().get_bot().start(settings.DISCORD_TOKEN)
        finally:
            await Singleton().get_bot().close()
    finally:
        await runner.cleanup()


if __name__ == "__main__":
    loop = get_event_loop()
    loop.run_until_complete(run_bot())
