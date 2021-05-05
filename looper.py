"module for background tasks in the loop"
import requests
from discord.ext import tasks, commands
from copy import deepcopy
from channel import (
    delete_messages_older_than_n_seconds,
    handle_tagged_messages,
)


class Looper(commands.Cog):
    def __init__(self, bot, storage):
        self.index = 0
        self.bot = bot
        self.printer.start()
        self.storage = storage
        self.settings = storage.settings
        self.unique_tag = storage.unique_tag

    def cog_unload(self):
        self.printer.cancel()
        print('unloading')

    @tasks.loop(seconds=5.0)
    async def printer(self):

        self.storage.players = requests.get(
            self.settings.player_request_url).json()
        self.storage.bases = requests.get(
            self.settings.base_request_url).json()

        print(self.index)
        self.index += 1

        channels = deepcopy(self.storage.channels)
        for channel_id in channels:

            await delete_messages_older_than_n_seconds(self.bot,
                                                       self.unique_tag, 10,
                                                       channel_id)
            await handle_tagged_messages(self.bot, self.unique_tag, channel_id)

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()
