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
            try:
                # delete expired
                await self.chanell_controller.delete_exp_msgs(channel_id, 40)

                # date stamp
                with open('templates/date.md') as file_:
                    template = Template(file_.read())
                    info = template.render(date=str(
                        datetime.datetime.utcnow()),
                                           timestamp=data.players['timestamp'])

                # bases
                rendered_bases = await self.storage.base.view(
                    channel_id, data.bases)

                friend_tags, friend_alert = await self.storage.friend.get_data(
                    channel_id)
                enemy_tags, enemy_alert = await self.storage.enemy.get_data(
                    channel_id)
                unrecognized_alert = await self.storage.unrecognized.get_data(
                    channel_id)

                system_tags = await self.storage.system.get_data(channel_id)
                region_tags = await self.storage.region.get_data(channel_id)

                players_all_list = {
                    item['name']: item
                    for item in data.players['players']
                }

                trackable_players = dict(
                    self.storage.region.process_tag(players_all_list, 'region',
                                                    region_tags),
                    **(self.storage.system.process_tag(players_all_list,
                                                       'system', system_tags)))

                friends = self.storage.friend.process_tag(
                    trackable_players, 'name', friend_tags)
                enemies = self.storage.enemy.process_tag(
                    trackable_players, 'name', enemy_tags)

                unrecognized_tags = set(trackable_players) - set(
                    friends) - set(enemies)

                unregonizeds = {
                    key: value
                    for key, value in trackable_players.items()
                    if key in unrecognized_tags
                }
                with open('templates/players.md') as file_:
                    template = Template(file_.read())

                    def rendering(title, data, alert_level):
                        if data:
                            alert_needed = False
                            if alert_level is not None:
                                alert_needed = len(data) >= alert_level
                            return template.render(title=title,
                                                   data=data,
                                                   alert=alert_needed)
                        return ''

                    rendered_unrecognized = rendering('Players', unregonizeds,
                                                      unrecognized_alert)
                    rendered_enemies = rendering('Enemies', enemies,
                                                 enemy_alert)
                    rendered_friends = rendering('Friends', friends,
                                                 friend_alert)

                    rendered_all = (rendered_unrecognized + rendered_enemies +
                                    rendered_friends)

                # send final data update
                try:
                    await self.chanell_controller.update_info(
                        channel_id, info + rendered_all + rendered_bases)
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
