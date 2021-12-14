"module for background tasks in the loop"
import datetime

import discord
from discord.ext import commands, tasks
from threading import Thread
import asyncio
import time

from src.views import View
from src.data_model import DataModel
import src.settings as settings
from src.storage import Storage


class Looper(commands.Cog):
    def __init__(self, bot, storage: Storage, chanell_controller):
        self.bot = bot
        self.printer.start()
        self.storage = storage
        self.chanell_controller = chanell_controller

        api_data = self.storage.get_game_data({})
        self.data = DataModel(api_data=api_data)

    async def cog_unload(self):
        self.printer.cancel()
        print('unloading...')

    @tasks.loop(seconds=5.0, reconnect=True, count=1)
    async def printer(self):

        try:
            if settings.DEBUG:
                print(
                    f'{datetime.datetime.utcnow()} OK executing printer loop')

            updating_api_data = await self.storage.a_get_game_data(
                self.data.previous_forum_records)

            for record in updating_api_data.new_forum_records:
                self.data.previous_forum_records[record.title] = record

            self.data.update(updating_api_data)

            self.storage.save_channel_settings()

            channel_ids = [int(item) for item in self.storage.channels.keys()]

            for channel_id in channel_ids:
                try:
                    if settings.DEBUG:
                        print(
                            f'channel {channel_id} in {self.bot.get_channel(channel_id).guild}'
                        )

                    # delete expired
                    await self.chanell_controller.delete_exp_msgs(
                        channel_id, 40)

                    rendered_date, rendered_all, render_forum_records = await View(
                        self.data.api_data, self.storage, channel_id,
                        self.data.api_data.new_forum_records).render_all()

                    if settings.DEBUG:
                        print(f"view_new_records={render_forum_records}")
                    # send final data update
                    try:
                        await self.chanell_controller.update_info(
                            channel_id,
                            rendered_all,
                            render_forum_records=render_forum_records)
                    except discord.errors.HTTPException:
                        await self.chanell_controller.update_info(
                            channel_id, rendered_date +
                            '\n**ERR: you tried to render too much info!**' +
                            '\nremove some of the values from config' +
                            '\nor write them fully instead of tags')
                except (discord.errors.DiscordException, AttributeError,
                        Exception) as error:
                    if settings.DEBUG:
                        print(
                            f"{str(datetime.datetime.utcnow())} "
                            f"ERR  {str(error)} for channel: {str(channel_id)}"
                        )
                    if isinstance(error, KeyboardInterrupt):
                        raise KeyboardInterrupt(
                            "time to exit, KeyboardInterrupt")
        except Exception as error:
            if settings.DEBUG:
                print(f"{str(datetime.datetime.utcnow())} "
                      f"ERR massive {str(error)} for loop task")
                raise Exception("detailed_exception") from error
            if isinstance(error, KeyboardInterrupt):
                print("gracefully exiting")

    def task(self, loop):
        asyncio.run_coroutine_threadsafe(self.printer(), loop)

    def task_creator(self, loop, delay=5):
        print("starting task creator")
        while True:
            thread = Thread(
                target=self.task,
                args=(loop, ),
                daemon=True,
            )
            thread.start()
            time.sleep(delay)

    def create_task_creator(self, loop):
        "launch background daemon process"
        thread = Thread(
            target=self.task_creator,
            args=(loop, ),
            daemon=True,
        )
        thread.start()

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        self.create_task_creator(asyncio.get_running_loop())
        await self.bot.wait_until_ready()
