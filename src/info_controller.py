class InfoController():
    def __init__(self, source, category):
        self.source = source
        self.category = category

    def create_if_none(self, channel_id):
        if str(channel_id) not in self.source:
            self.source[str(channel_id)] = {}

        if self.category not in self.source[str(channel_id)]:
            self.source[str(channel_id)][self.category] = {'list': []}

    async def add(self, channel_id, *args):
        self.create_if_none(channel_id)
        for item in args[0]:
            self.source[str(channel_id)][self.category]['list'].append(item)

    async def remove(self, channel_id, *args):
        self.create_if_none(channel_id)

        for item in args[0]:
            self.source[str(channel_id)][self.category]['list'].remove(item)

    async def clear(self, channel_id, *args):
        self.create_if_none(channel_id)

        self.source[str(channel_id)][self.category]['list'].clear()

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
    def create_if_none(self, channel_id):
        if self.category not in self.source[str(channel_id)]:
            self.source[str(channel_id)][self.category] = {
                'list': [],
                'alert': 999
            }

    async def alert(self, channel_id, *args):
        self.create_if_none(channel_id)

        if args:
            self.source[str(channel_id)][self.category]['alert'] = int(
                args[0][0])

    async def get_data(self, channel_id):
        if self.category in self.source[str(channel_id)]:
            return self.source[str(channel_id)][self.category][
                'list'], self.source[str(channel_id)][self.category]['alert']
        return None, 999


class AlertOnlyController(InfoWithAlertController):
    def create_if_none(self, channel_id):
        if self.category not in self.source[str(channel_id)]:
            self.source[str(channel_id)][self.category] = {'alert': 999}

    async def get_data(self, channel_id):
        if self.category in self.source[str(channel_id)]:
            return self.source[str(channel_id)][self.category]['alert']
        return 999
