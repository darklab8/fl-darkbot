"module for background tasks in the loop"
import datetime

import discord
from discord.ext import commands, tasks

from views import View


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

    async def cog_unload(self):
        self.printer.cancel()
        print('unloading...')

    @tasks.loop(seconds=5.0, reconnect=True)
    async def printer(self):

        try:
            print(f'{datetime.datetime.utcnow()} OK executing printer loop')

            data = self.storage.get_game_data()

            # calculating previous health about bases
            if self.previous_bases != data.bases:
                self.different_bases = self.previous_bases
            self.previous_bases = data.bases
            data.different_bases = self.different_bases

            def health_diff(a, b):
                if a < b:
                    return '+ '
                elif a > b:
                    return '- '
                return ''

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

                    rendered_date, rendered_all = await View(
                        data, self.storage, channel_id).render_all()
                    # send final data update
                    try:
                        await self.chanell_controller.update_info(
                            channel_id, rendered_all)
                    except discord.errors.HTTPException:
                        await self.chanell_controller.update_info(
                            channel_id, rendered_date +
                            '\n**ERR: you tried to render too much info!**' +
                            '\nremove some of the values from config' +
                            '\nor write them fully instead of tags')
                except discord.errors.DiscordException as error:
                    print(f"{str(datetime.datetime.utcnow())} "
                          f"ERR  {str(error)} for channel: {str(channel_id)}")
                except AttributeError as error:
                    print(f"{str(datetime.datetime.utcnow())} "
                          f"ERR  {str(error)} for channel: {str(channel_id)}")
        except Exception as error:
            print(f"{str(datetime.datetime.utcnow())} "
                  f"ERR massive {str(error)} for loop task")

    @printer.before_loop
    async def before_printer(self):
        print('waiting...')
        await self.bot.wait_until_ready()
