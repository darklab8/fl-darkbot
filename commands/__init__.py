from discord.ext import commands
from .root import attach_root
from .base import attach_base
from .system import attach_system
from .region import attach_region
from .friend import attach_friend
from .enemy import attach_enemy
from .unrecognized import attach_unrecognized
from .forum import attach_forum


def attach_commands(bot, storage, chanell_controller) -> commands.Bot:
    bot = attach_root(bot, storage, chanell_controller)
    bot = attach_base(bot, storage)
    bot = attach_system(bot, storage)
    bot = attach_region(bot, storage)
    bot = attach_friend(bot, storage)
    bot = attach_enemy(bot, storage)
    bot = attach_unrecognized(bot, storage)
    bot = attach_forum(bot, storage)

    return bot
