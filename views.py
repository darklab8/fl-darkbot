from jinja2 import Template
import datetime


class View():
    def __init__(self, data, storage, channel_id):
        self.data = data
        self.storage = storage
        self.channel_id = channel_id

    async def render_date(self, timestamp):
        with open('templates/date.md') as file_:
            template = Template(file_.read())
            return template.render(tag=self.storage.unique_tag,
                                   date=str(datetime.datetime.utcnow()),
                                   timestamp=timestamp)

    async def base_view(self, base, channel_id, bases):
        base_tags = await base.get_data(channel_id)

        if base_tags is None:
            return ''

        rendering_bases = {}
        for base_tag in base_tags:
            adding_bases = {
                key: value
                for key, value in bases.items() if base_tag in key
            }
            rendering_bases = dict(rendering_bases, **adding_bases)

        if not rendering_bases:
            return ''

        with open('templates/base.md') as file_:
            template = Template(file_.read())
            rendered_bases = template.render(data=rendering_bases)
            return rendered_bases

    async def render_players(self, storage, channel_id, players):
        friend_tags, friend_alert = await storage.friend.get_data(channel_id)
        enemy_tags, enemy_alert = await storage.enemy.get_data(channel_id)
        unrecognized_alert = await storage.unrecognized.get_data(channel_id)

        system_tags = await storage.system.get_data(channel_id)
        region_tags = await storage.region.get_data(channel_id)

        players_all_list = {item['name']: item for item in players['players']}

        trackable_players = dict(
            storage.region.process_tag(players_all_list, 'region',
                                       region_tags),
            **(storage.system.process_tag(players_all_list, 'system',
                                          system_tags)))

        friends = storage.friend.process_tag(trackable_players, 'name',
                                             friend_tags)
        enemies = storage.enemy.process_tag(trackable_players, 'name',
                                            enemy_tags)

        unrecognized_tags = set(trackable_players) - set(friends) - set(
            enemies)

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
            rendered_enemies = rendering('Enemies', enemies, enemy_alert)
            rendered_friends = rendering('Friends', friends, friend_alert)

            return (rendered_unrecognized + rendered_enemies +
                    rendered_friends)

    async def render_all(self):
        # date stamp
        info = await self.render_date(self.data.players['timestamp'])

        # bases
        rendered_bases = await self.base_view(self.storage.base,
                                              self.channel_id, self.data.bases)

        # players
        rendered_players = await self.render_players(self.storage,
                                                     self.channel_id,
                                                     self.data.players)

        return info, info + rendered_bases + rendered_players