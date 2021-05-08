"starting module"
from commands import attach_commands
from storage import Storage
from looper import Looper
from discord.ext import commands
from channel import ChannelConstroller
# nice settings loading


class CreateApp():
    def __init__(self):
        self.storage = Storage()
        self.bot = commands.Bot(command_prefix='.')
        channel_controller = ChannelConstroller(self.bot,
                                                self.storage.unique_tag)
        bot = attach_commands(self.bot, self.storage, channel_controller)
        _ = Looper(bot, self.storage, channel_controller)

    def run(self):
        self.bot(self.storage.settings.secret_key)


if __name__ == '__main__':
    app = CreateApp()
    app.run()
