import json
from jinja2 import Template


class InfoController():
    def __init__(self, source, category):
        self.source = source
        self.category = category

    def create_if_none(self, ctx):
        if self.category not in self.source[str(ctx.channel.id)]:
            self.source[str(ctx.channel.id)][self.category] = {'list': []}

    async def add(self, ctx, *args):
        self.create_if_none(ctx)
        for item in args[0]:
            self.source[str(
                ctx.channel.id)][self.category]['list'].append(item)

    async def remove(self, ctx, *args):
        self.create_if_none(ctx)

        for item in args[0]:
            self.source[str(
                ctx.channel.id)][self.category]['list'].remove(item)

    async def clear(self, ctx, *args):
        self.create_if_none(ctx)

        self.source[str(ctx.channel.id)][self.category]['list'].clear()

    async def lst(self, ctx, *args):
        self.create_if_none(ctx)

        with open('templates/json.md') as file_:
            template = Template(file_.read())

            await ctx.send(
                template.render(data=json.dumps(
                    self.source[str(ctx.channel.id)][self.category], indent=2))
            )

    async def get_data(self, channel_id):
        if self.category in self.source[str(channel_id)]:
            return self.source[str(channel_id)][self.category]['list']
        return None

    def process_tag(self, from_where, access_key, tags):
        output = {}
        if tags is not None:
            for tag in tags:
                found = {
                    key: value
                    for key, value in from_where.items()
                    if tag in value[access_key]
                }
                output = dict(output, **found)
        return output


class InfoWithAlertController(InfoController):
    def create_if_none(self, ctx):
        if self.category not in self.source[str(ctx.channel.id)]:
            self.source[str(ctx.channel.id)][self.category] = {
                'list': [],
                'alert': 999
            }

    async def alert(self, ctx, *args):
        self.create_if_none(ctx)

        if args:
            self.source[str(ctx.channel.id)][self.category]['alert'] = int(
                args[0][0])

    async def get_data(self, channel_id):
        if self.category in self.source[str(channel_id)]:
            return self.source[str(channel_id)][self.category][
                'list'], self.source[str(channel_id)][self.category]['alert']
        return None, 999


class AlertOnlyController(InfoWithAlertController):
    def create_if_none(self, ctx):
        if self.category not in self.source[str(ctx.channel.id)]:
            self.source[str(ctx.channel.id)][self.category] = {'alert': 999}

    async def get_data(self, channel_id):
        if self.category in self.source[str(channel_id)]:
            return self.source[str(channel_id)][self.category]['alert']
        return 999
