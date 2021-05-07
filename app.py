"starting module"
from commands import attach_commands
from storage import Storage
from looper import Looper
from discord.ext import commands

# nice settings loading


def create_app():
    storage = Storage()
    bot = commands.Bot(command_prefix='.')
    bot = attach_commands(bot, storage)
    _ = Looper(bot, storage)
    return bot, storage.settings.secret_key


if __name__ == '__main__':
    bot, secret_key = create_app()
    bot.run(secret_key)
