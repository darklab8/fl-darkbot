"module for background tasks in the loop"
from discord.ext import tasks, commands
from jinja2 import Template
import datetime
import json
import discord
from storage import render_players
from .views import render_date


class Looper(commands.Cog):
    def __init__(self, bot, storage, chanell_controller):
        self.index = 0
        self.bot = bot
        self.printer.start()
        self.storage = storage
        self.chanell_controller = chanell_controller

    def cog_unload(self):
        self.printer.cancel()
        print('unloading')

    @tasks.loop(seconds=5.0)
    async def printer(self):

        print(self.index)
        self.index += 1

        data = self.storage.get_game_data()
        self.storage.save_channel_settings()

        channel_ids = [int(item) for item in self.storage.channels.keys()]

        for channel_id in channel_ids:
            try:
                # delete expired
                await self.chanell_controller.delete_exp_msgs(channel_id, 40)

                # date stamp
                info = render_date(data.players['timestamp'])

                # bases
                rendered_bases = await self.storage.base.view(
                    channel_id, data.bases)

                # players
                rendered_players = render_players(self.storage, channel_id,
                                                  data.players)
                # send final data update
                try:
                    await self.chanell_controller.update_info(
                        channel_id, info + rendered_players + rendered_bases)
                except discord.errors.HTTPException:
                    await self.chanell_controller.update_info(
                        channel_id, info +
                        '\n**ERROR: you tried to render too much info!**' +
                        '\nremove some of the values from config' +
                        '\nor write them fully instead of tags')
            except discord.errors.Forbidden:
                print("skipping forbidden channel")

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()
