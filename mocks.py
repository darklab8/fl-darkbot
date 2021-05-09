from types import SimpleNamespace
from channel import IMessageBus


class MockMessage():
    def __init__(self):
        self.author = SimpleNamespace(id=1)
        self.content = 'random content'

    async def delete(self):
        pass


class MockDiscordMessageBus(IMessageBus):
    def __init__(self):
        pass

    def bot_user_id(self):
        return 1

    async def delete(self, message):
        message.delete()

    async def send(self, channel_id, content):
        return [MockMessage(), MockMessage(), MockMessage()]

    async def history(self, channel_id, older_than_n_seconds=0):
        return [MockMessage(), MockMessage(), MockMessage()]
