"bot initialization + its commands"
from discord.ext import commands
from permissions import (
    all_checks,
    connected_to_channel,
)
from channel import delete_messages_older_than_n_seconds
import random


def attach_commands(storage, command_prefix='.'):
    bot = commands.Bot(command_prefix)

    @bot.event
    async def on_ready():
        print('We have logged in as {0.user}'.format(bot))

    @bot.command(name='connect')
    @commands.check_any(all_checks())
    async def connect_the_channel(ctx):
        "connects to channel"
        if (ctx.channel.id) not in storage.channels:
            storage.channels[(ctx.channel.id)] = {}
            await ctx.send('connected')
        else:
            await ctx.send('we are already connected')

    @bot.command(name='disconnect')
    @commands.check_any(all_checks())
    async def diconnect_the_channel(ctx):
        "disconnects from channel"
        if (ctx.channel.id) in storage.channels:
            storage.channels.pop((ctx.channel.id))
            await ctx.send('disconnected')
        else:
            await ctx.send('we are already disconnected')

    @bot.command(name='check')
    @commands.check_any(connected_to_channel(storage))
    async def check(ctx, number: int):
        await ctx.send(f"{number} is your lucky number!")

    @bot.command(name='clear')
    @commands.check_any(connected_to_channel(storage))
    async def clear(ctx):
        await delete_messages_older_than_n_seconds(ctx.bot, storage.unique_tag,
                                                   0, ctx.channel.id)

    @bot.command(name='fun')
    @commands.check_any(connected_to_channel(storage))
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

    return bot
