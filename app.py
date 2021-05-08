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
        bot = commands.Bot(command_prefix='.')
        channel_controller = ChannelConstroller(bot, self.storage.unique_tag)
        self.bot = attach_commands(bot, self.storage, channel_controller)
        self.looper = Looper(bot, self.storage, channel_controller)

    def run(self):
        self.bot.run(self.storage.settings.secret_key)


if __name__ == '__main__':
    app = CreateApp()
    app.run()
