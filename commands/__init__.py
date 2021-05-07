from discord.ext import commands
from .root import attach_root
from .base import attach_base
from .system import attach_system


def attach_commands(bot, storage, chanell_controller) -> commands.Bot:
    bot = attach_root(bot, storage, chanell_controller)
    bot = attach_base(bot, storage)
    bot = attach_system(bot, storage)

    return bot
