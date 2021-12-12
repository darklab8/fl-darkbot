"starting module"
from discord.ext import commands

from src.commands import attach_commands
from src.channel import ChannelConstroller
from src.looper import Looper
from src.storage import Storage
from src.channel import DiscordMessageBus


class CreateApp():
    def __init__(self):
        bot = commands.Bot(command_prefix='.')
        self.storage = Storage()
        self.channel_controller = ChannelConstroller(DiscordMessageBus(bot),
                                                     self.storage.unique_tag)
        self.bot = attach_commands(bot, self.storage, self.channel_controller)
        self.looper = Looper(bot, self.storage, self.channel_controller)

    def run(self):
        self.bot.run(self.storage.settings.secret_key)

        # self.bot


if __name__ == '__main__':
    app = CreateApp()
    app.run()
