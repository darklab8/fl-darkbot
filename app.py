import os
from types import SimpleNamespace
from dotenv import load_dotenv
import discord
# import asyncio
from discord.ext import commands
import json
import random
from permissions import connected_to_channel
from permissions import (
    is_guild_owner,
    can_manage_channels,
    is_bot_controller,
    all_checks,
)

from universal import delete_messages_older_than_n_seconds
from looper import Looper

# nice settings loading
load_dotenv()

settings = SimpleNamespace()
for item, value in os.environ.items():
    setattr(settings, item, value)

# discord.Guild.get_channel

STORAGE = SimpleNamespace()
with open('channels.json') as file_:
    STORAGE.channels = json.loads(file_.read())

STORAGE.unique_tag = 'dark_info:'
bot = commands.Bot(command_prefix='$')

cog = Looper(bot, STORAGE, settings)


@bot.command(name='connect')
@commands.check_any(all_checks())
async def connect_the_channel(ctx):
    "connects to channel"
    if (ctx.channel.id) not in STORAGE.channels:
        STORAGE.channels[(ctx.channel.id)] = {}
        await ctx.send('connected')
    else:
        await ctx.send('we are already connected')


@bot.command(name='disconnect')
@commands.check_any(all_checks())
async def diconnect_the_channel(ctx):
    "disconnects from channel"
    if (ctx.channel.id) in STORAGE.channels:
        STORAGE.channels.pop((ctx.channel.id))
        await ctx.send('disconnected')
    else:
        await ctx.send('we are already disconnected')


@bot.command(name='check')
@commands.check_any(connected_to_channel(STORAGE))
async def check(ctx, number: int):
    await ctx.send(f"{number} is your lucky number!")


@bot.command(name='clear')
@commands.check_any(connected_to_channel(STORAGE))
async def clear(ctx):
    await delete_messages_older_than_n_seconds(ctx.bot, STORAGE.unique_tag, 0,
                                               ctx.channel.id)


@bot.event
async def on_ready():
    print('We have logged in as {0.user}'.format(bot))
    # breakpoint()


@bot.command(name='fun')
@commands.check_any(connected_to_channel(STORAGE))
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


@check.error
async def check_error(ctx, error):
    if isinstance(error, commands.CommandError):
        await ctx.send('incorrect command!')


@bot.command(name='me')
@commands.check_any(commands.is_owner())
async def only_me(ctx, ):
    "secret command"
    await ctx.send('Papa!')
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
