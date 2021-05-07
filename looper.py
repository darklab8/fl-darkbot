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
                async def bases_to_view():
                    base_tags, base_alert = await self.storage.base.get_data(
                        channel_id)

                    if base_tags is None:
                        return ''

                    rendering_bases = {}
                    for base_tag in base_tags:
                        adding_bases = {
                            key: value
                            for key, value in data.bases.items()
                            if base_tag in key
                        }
                        rendering_bases = dict(rendering_bases, **adding_bases)

                    if not rendering_bases:
                        return ''

                    with open('templates/base.md') as file_:
                        template = Template(file_.read())
                        rendered_bases = template.render(data=rendering_bases)
                        return rendered_bases

                rendered_bases = await bases_to_view()

                friend_tags, friend_alert = await self.storage.friend.get_data(
                    channel_id)
                enemy_tags, enemy_alert = await self.storage.enemy.get_data(
                    channel_id)
                unrecognized_tags, unrecognized_alert = await self.storage.unrecognized.get_data(
                    channel_id)
                system_tags, _ = await self.storage.system.get_data(channel_id)
                region_tags, _ = await self.storage.region.get_data(channel_id)

                players_all_list = {
                    item['name']: item
                    for item in data.players['players']
                }

                def process_tag(output, from_where, access_key, tags):
                    if tags is not None:
                        for tag in tags:
                            found = {
                                key: value
                                for key, value in from_where.items()
                                if tag in value[access_key]
                            }
                            output = dict(output, **found)
                    return output

                trackable_players = {}
                trackable_players = process_tag(trackable_players,
                                                players_all_list, 'region',
                                                region_tags)
                trackable_players = process_tag(trackable_players,
                                                players_all_list, 'system',
                                                system_tags)

                friends = {}
                friends = process_tag(friends, trackable_players, 'name',
                                      friend_tags)

                enemies = {}
                enemies = process_tag(enemies, trackable_players, 'name',
                                      enemy_tags)

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
