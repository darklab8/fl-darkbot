import os
from types import SimpleNamespace
from dotenv import load_dotenv
import discord
# import asyncio
from discord.ext import commands

import random

from permissions import (
    is_guild_owner,
    can_manage_channels,
    is_bot_controller,
    all_checks,
)

from looper import Looper

# nice settings loading
load_dotenv()
settings = SimpleNamespace()
for item, value in os.environ.items():
    setattr(settings, item, value)

# discord.Guild.get_channel

APIDATA = SimpleNamespace()

bot = commands.Bot(command_prefix='$')

cog = Looper(bot, APIDATA, settings)


@bot.event
async def on_ready():
    print('We have logged in as {0.user}'.format(bot))


@bot.command(name='fun')
@commands.check_any(commands.is_owner(), all_checks())
async def nine_nine(ctx):
    "says random message"
    brooklyn_99_quotes = [
        'I\'m the human form of the 💯 emoji.',
        'Bingpot!',
        ('Cool. Cool cool cool cool cool cool cool, '
         'no doubt no doubt no doubt no doubt.'),
    ]

    response = random.choice(brooklyn_99_quotes)
    await ctx.send(response)


@bot.command(name='check')
@commands.check_any(commands.is_owner(), all_checks())
async def check(ctx, number: int):
    await ctx.send(f"{number} is your lucky number!")


@check.error
async def check_error(ctx, error):
    if isinstance(error, commands.CommandError):
        await ctx.send('incorrect command!')


@bot.command(name='author')
@commands.check_any(commands.is_owner())
async def only_me(ctx, ):
    await ctx.send('Only you!')
    # {discord.Guild.get_channel()}


bot.run(settings.secret_key)

# @bot.event
# async def on_message(message):
#     if message.content.startswith('$thumb'):
#         channel = message.channel
#         await channel.send('Send me that 👍 reaction, mate')

#         def check(reaction, user):
#             return user == message.author and str(reaction.emoji) == '👍'

#         try:
#             reaction, user = await bot.wait_for('reaction_add',
#                                                 timeout=60.0,
#                                                 check=check)
#         except asyncio.TimeoutError:
#             await channel.send('👎')
#         else:
#             await channel.send('👍')

# def owner_or_permissions(**perms):
#     original = commands.has_permissions(**perms).predicate

#     async def extended_check(ctx):
#         if ctx.guild is None:
#             return False
#         return ctx.guild.owner_id == ctx.author.id or await original(
#             ctx) or ctx.message.author.id == 370435997974134785

#     return commands.check(extended_check)

# def check_if_it_is_me(ctx):
#     return ctx.message.author.id == 370435997974134785

# @bot.command(name='me')
# @commands.check(check_if_it_is_me)
# async def only_for_me(ctx):
#     await ctx.send('I know you!')
