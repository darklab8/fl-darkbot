import requests
from discord.ext import tasks, commands


class Looper(commands.Cog):
    def __init__(self, bot, apidata, settings):
        self.index = 0
        self.bot = bot
        self.printer.start()
        self.settings = settings
        self.apidata = apidata

    def cog_unload(self):
        self.printer.cancel()
        print('unloading')

    @tasks.loop(seconds=5.0)
    async def printer(self):
        self.apidata.players = requests.get(
            self.settings.player_request_url).json()
        self.apidata.bases = requests.get(
            self.settings.base_request_url).json()
        print(self.index)
        self.index += 1

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()
