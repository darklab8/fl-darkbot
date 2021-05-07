from jinja2 import Template
from .info_controller import InfoController
import datetime


def render_date(timestamp):
    with open('templates/date.md') as file_:
        template = Template(file_.read())
        return template.render(date=str(datetime.datetime.utcnow()),
                               timestamp=timestamp)


class BaseViewer(InfoController):
    async def view(self, channel_id, bases):
        base_tags = await self.get_data(channel_id)

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


async def render_players(storage, channel_id, players):
    friend_tags, friend_alert = await storage.friend.get_data(channel_id)
    enemy_tags, enemy_alert = await storage.enemy.get_data(channel_id)
    unrecognized_alert = await storage.unrecognized.get_data(channel_id)

    system_tags = await storage.system.get_data(channel_id)
    region_tags = await storage.region.get_data(channel_id)

    players_all_list = {item['name']: item for item in players['players']}

    trackable_players = dict(
        storage.region.process_tag(players_all_list, 'region', region_tags),
        **(storage.system.process_tag(players_all_list, 'system',
                                      system_tags)))

    friends = storage.friend.process_tag(trackable_players, 'name',
                                         friend_tags)
    enemies = storage.enemy.process_tag(trackable_players, 'name', enemy_tags)

    unrecognized_tags = set(trackable_players) - set(friends) - set(enemies)

    unregonizeds = {
        key: value
        for key, value in trackable_players.items() if key in unrecognized_tags
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

        return (rendered_unrecognized + rendered_enemies + rendered_friends)


async def render_all(data, storage, channel_id):
    # date stamp
    info = render_date(data.players['timestamp'])

    # bases
    rendered_bases = await storage.base.view(channel_id, data.bases)

    # players
    rendered_players = render_players(storage, channel_id, data.players)

    return info, info + rendered_bases + rendered_players
