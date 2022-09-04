import discord
from discord.ext import commands
from .core import settings
import secrets


async def create_bot():
    intents = discord.Intents.default()
    intents.message_content = True
    intents.members = True

    bot = commands.Bot(
        command_prefix=secrets.token_hex(10),
        case_insensitive=True,
        intents=intents,
    )

    return bot


async def run_bot(bot: commands.Bot):
    async with bot:
        await bot.start(token=settings.DISCORD_TOKEN)
