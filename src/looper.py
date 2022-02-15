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
import logging
from .message_sent_history import message_history
from .shuffler import shuffled_dict

class Looper(commands.Cog):
    def __init__(self, bot, storage: Storage, chanell_controller):
        self.bot = bot
        self.printer.start()
        self.storage = storage
        self.chanell_controller = chanell_controller

        api_data = self.storage.get_game_data()

        channels_ids = self.storage.get_channels_id()

        for record in api_data.new_forum_records:
            for channel_id in channels_ids:
                message_history.add_message(channel_id=channel_id, record=record)

        self.data = DataModel(api_data=api_data)

    async def cog_unload(self):
        self.printer.cancel()
        print('unloading...')

    @tasks.loop(seconds=5.0, reconnect=True, count=1)
    async def printer(self):

        try:
            logging.info('OK executing printer loop')

            updating_api_data = await self.storage.a_get_game_data()

            self.data.update(updating_api_data)

            # logging.info(f"context=new_forum_records type=looper, data={self.data.api_data.new_forum_records}")

            self.storage.save_channel_settings()

            channel_ids = self.storage.get_channels_id()

            forbidden_channels = []
            allowed_channels = {}

            for channel_id in channel_ids:
                channel_info = self.bot.get_channel(channel_id)

                if channel_info is None:
                    forbidden_channels.append(channel_id)
                else:
                    allowed_channels[channel_id] = channel_info

            logging.info(f'context=allowed_channels, allowed_channels={allowed_channels.keys()}')
            logging.info(f'context=forbidden_channels, forbidden_channels={forbidden_channels}')
            
            shuffled_allowed_channels = shuffled_dict(allowed_channels)
            for channel_id, channel_info in shuffled_allowed_channels.items():
                try:
                    logging.info(f'context=loop_begins_for_channel channel={channel_id} in guild={self.bot.get_channel(channel_id).guild}')

                    # delete expired messages
                    await self.chanell_controller.delete_exp_msgs(
                        channel_id, 40)

                    logging.info(f'context=channel_loop, channel={channel_id}, msg=deleting_old_msgs')

                    rendered_date, rendered_all, render_forum_records = await View(
                        self.data.api_data, self.storage,
                        channel_id).render_all()

                    logging.info(f'context=channel_loop, channel={channel_id}, msg=rendered_all')

                    # send final data update
                    try:
                        await self.chanell_controller.update_info(
                            channel_id,
                            rendered_all,
                            render_forum_records=render_forum_records)

                        logging.info(f'context=channel_loop, channel={channel_id}, msg=update_info_is_done')
                    except discord.errors.HTTPException:
                        await self.chanell_controller.update_info(
                            channel_id, rendered_date +
                            '\n**ERR: you tried to render too much info!**' +
                            '\nremove some of the values from config' +
                            '\nor write them fully instead of tags')
                except (discord.errors.DiscordException, AttributeError,
                        Exception) as error:
                    error_msg = f"ERR, loop_cycle, channel_id={str(channel_id)}, error={str(error)}"
                    
                    searched_error = "error=403 Forbidden (error code: 50001): Missing Access"
                    is_access_to_channel_forbidden_and_bot_removed = searched_error in error_msg
                    if is_access_to_channel_forbidden_and_bot_removed:
                        self.storage.channels.pop(channel_id)
                        logging.info(f"channel_id={str(channel_id)}, msg=removed channel which had {searched_error}")

                    logging.info(error_msg)
                    if isinstance(error, KeyboardInterrupt):
                        raise KeyboardInterrupt(
                            "time to exit, KeyboardInterrupt")
        except Exception as error:
            logging.error(f"ERR, context=whole_loop, error={str(error)}")
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
