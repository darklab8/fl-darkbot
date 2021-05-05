import requests
from discord.ext import tasks, commands
from discord import channel
import datetime
from jinja2 import Template
from copy import deepcopy
from universal import (
    delete_messages_older_than_n_seconds,
    deleting_message,
)


class Looper(commands.Cog):
    def __init__(self, bot, storage, settings):
        self.index = 0
        self.bot = bot
        self.printer.start()
        self.settings = settings
        self.storage = storage
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

        # async for elem in channel.history():
        #     print(elem)
        channels = deepcopy(self.storage.channels)
        for channel_id in channels:

            await delete_messages_older_than_n_seconds(self.bot,
                                                       self.unique_tag, 30,
                                                       channel_id)

            content_search = await self.bot.get_channel(channel_id).history(
                limit=200).flatten()
            content = [
                item for item in content_search
                if self.unique_tag in item.content
            ]

            if not content:
                # create first msg
                await self.bot.get_channel(channel_id).send(
                    self.unique_tag + str(datetime.datetime.utcnow()))
            elif len(content) > 1:
                # delete all others
                deleting = content[1:]
                for message in deleting:
                    await deleting_message(message)
            else:
                with open('date.md') as file_:
                    template = Template(file_.read())

                    await content[0].edit(content=str(
                        template.render(date=str(datetime.datetime.utcnow()))))

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()
