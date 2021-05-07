"starting module"
from commands import attach_commands
from storage import Storage
from looper import Looper
from discord.ext import commands
from channel import ChannelConstroller
# nice settings loading


def create_app():
    storage = Storage()
    bot = commands.Bot(command_prefix='.')
    chanell_controller = ChannelConstroller(bot, storage.unique_tag)
    bot = attach_commands(bot, storage, chanell_controller)
    _ = Looper(bot, storage, chanell_controller)
    return bot, storage.settings.secret_key


if __name__ == '__main__':
    bot, secret_key = create_app()
    bot.run(secret_key)
