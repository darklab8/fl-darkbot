"module for background tasks in the loop"
from discord.ext import tasks, commands
from jinja2 import Template
import datetime
import json
import discord


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

            # delete expired
            await self.chanell_controller.delete_exp_msgs(channel_id, 40)

            with open('templates/date.md') as file_:
                template = Template(file_.read())
                info = template.render(date=str(datetime.datetime.utcnow()))

            base_tags = await self.storage.base.get_data(channel_id)

            rendering_bases = {}
            for base_tag in base_tags:
                adding_bases = {
                    key: value
                    for key, value in data.bases.items() if base_tag in key
                }
                rendering_bases = dict(rendering_bases, **adding_bases)
            list_of_bases = [
                json.dumps([value['health'], key])
                for key, value in rendering_bases.items()
            ]

            with open('templates/base.md') as file_:
                template = Template(file_.read())
                rendered_bases = template.render(data=rendering_bases)

            # send final data update
            try:
                await self.chanell_controller.update_info(
                    channel_id, info + rendered_bases)
            except discord.errors.HTTPException:
                await self.chanell_controller.update_info(
                    channel_id,
                    info + '\n**ERROR: you tried to render too much info!**' +
                    '\nremove some of the values from config' +
                    '\nor write them fully instead of tags')

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()
