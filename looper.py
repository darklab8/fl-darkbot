"module for background tasks in the loop"
import datetime

import discord
from discord.ext import commands, tasks
from threading import Thread
import asyncio
import time

from views import View
from forum_parser import get_forum_threads


class Looper(commands.Cog):
    def __init__(self, bot, storage, chanell_controller):
        self.bot = bot
        self.printer.start()
        self.storage = storage
        self.chanell_controller = chanell_controller

        bases = self.storage.get_game_data().bases
        # last data about bases different in values
        self.different_bases = bases
        # value from different loop
        self.previous_bases = bases

        self.previous_forum_records = {}
        # add on start all records
        self.get_new_forum_records()
        self.new_forum_records = []

    def get_new_forum_records(self):
        forum_records = get_forum_threads(
            forum_acc=self.storage.settings.forum_acc,
            forum_pass=self.storage.settings.forum_pass,
        )

        new_records = []
        for record in forum_records:
            if record.title not in self.previous_forum_records:
                self.previous_forum_records[record.title] = record
                new_records.append(record)
            else:
                if record.date != self.previous_forum_records[
                        record.title].date:
                    self.previous_forum_records[record.title] = record
                    new_records.append(record)
        return new_records

    async def cog_unload(self):
        self.printer.cancel()
        print('unloading...')

    @tasks.loop(seconds=5.0, reconnect=True, count=1)
    async def printer(self):

        try:
            print(f'{datetime.datetime.utcnow()} OK executing printer loop')

            data = self.storage.get_game_data()
            self.new_forum_records = await asyncio.to_thread(
                self.get_new_forum_records)

            print(f"to_thread={self.new_forum_records}")
            # for new_record in new_forum_records:
            #     print(new_record)

            # calculating previous health about bases
            if self.previous_bases != data.bases:
                self.different_bases = self.previous_bases
            self.previous_bases = data.bases
            data.different_bases = self.different_bases

            def health_diff(a, b):
                if a < b:
                    return 'Repairing'
                elif a > b:
                    return 'Degrading'
                return 'Static'

            data.bases = {
                key: dict(
                    value, **{
                        "diff":
                        health_diff(data.different_bases[key]['health'],
                                    value['health'])
                    })
                for key, value in data.bases.items()
            }

            self.storage.save_channel_settings()

            channel_ids = [int(item) for item in self.storage.channels.keys()]

            for channel_id in channel_ids:
                try:
                    print(
                        f'channel {channel_id} in {self.bot.get_channel(channel_id).guild}'
                    )

                    # delete expired
                    await self.chanell_controller.delete_exp_msgs(
                        channel_id, 40)

                    rendered_date, rendered_all, render_forum_records = await View(
                        data, self.storage, channel_id,
                        self.new_forum_records).render_all()

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
                    print(f"{str(datetime.datetime.utcnow())} "
                          f"ERR  {str(error)} for channel: {str(channel_id)}")
                    if isinstance(error, KeyboardInterrupt):
                        raise KeyboardInterrupt(
                            "time to exit, KeyboardInterrupt")
            self.new_forum_records = []
        except Exception as error:
            print(f"{str(datetime.datetime.utcnow())} "
                  f"ERR massive {str(error)} for loop task")
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
