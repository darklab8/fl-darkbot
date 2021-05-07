from discord.ext import commands
from .root import attach_root
from .base import attach_base


def attach_commands(bot, storage) -> commands.Bot:
    bot = attach_root(bot, storage)
    bot = attach_base(bot, storage)

    return bot
