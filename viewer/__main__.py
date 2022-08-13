import discord
from discord.ext import commands, tasks
from threading import Thread
import asyncio
import time
import logging
from utils.logger import Logger
from viewer.core import settings

logger = Logger(console_level="DEBUG")


class Looper(commands.Cog):
    def __init__(self, bot):
        logger.debug("executing __init__ started")
        self.bot = bot
        self.printer.start()
        logger.debug("executing __init__ finished")

    @tasks.loop(seconds=5.0, reconnect=True, count=1)
    async def printer(self, msg="default"):
        logger.debug(f"executing printer, msg={msg}")
        await self.bot.get_channel(840251517398548521).send(f"msg={msg}")

    def task(self, loop):
        logger.debug("executing task started")
        asyncio.run_coroutine_threadsafe(self.printer(msg="task"), loop)
        logger.debug("executing task finished")

    def task_creator(self, loop, delay=5):
        logger.debug("executing task_creator started")
        while True:
            thread = Thread(
                target=self.task,
                args=(loop,),
                daemon=True,
            )
            thread.start()
            time.sleep(delay)

    def create_task_creator(self, loop):
        logger.debug("executing create_task_creator started")
        thread = Thread(
            target=self.task_creator,
            args=(loop,),
            daemon=True,
        )
        thread.start()
        logger.debug("executing create_task_creator finished")

    @printer.before_loop
    async def before_printer(self):
        logger.debug("executing before_printer started")
        self.create_task_creator(asyncio.get_running_loop())
        await self.bot.wait_until_ready()
        logger.debug("executing before_printer finished")


class CreateApp:
    def __init__(self):
        self.bot = commands.Bot(command_prefix=".")
        self.looper = Looper(self.bot)

    def run(self):
        self.bot.run(settings.DISCORD_TOKEN)


if __name__ == "__main__":
    app = CreateApp()
    app.run()
