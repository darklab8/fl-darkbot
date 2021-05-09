"starting module"
from discord.ext import commands

from commands import attach_commands
from channel import ChannelConstroller
from looper import Looper
from storage import Storage


class CreateApp():
    def __init__(self):
        self.storage = Storage()
        bot = commands.Bot(command_prefix='.')
        channel_controller = ChannelConstroller(bot, self.storage.unique_tag)
        self.bot = attach_commands(bot, self.storage, channel_controller)
        self.looper = Looper(bot, self.storage, channel_controller)

    def run(self):
        self.bot.run(self.storage.settings.secret_key)


if __name__ == '__main__':
    app = CreateApp()
    app.run()
