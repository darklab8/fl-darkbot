"bot initialization + its commands"
from discord.ext import commands
from permissions import (
    all_checks,
    connected_to_channel,
)
from channel import delete_messages_older_than_n_seconds
import random

from .consts import timedelta
from jinja2 import Template
import json


def attach_root(bot, storage) -> commands.Bot:
    class MyHelpCommand(commands.DefaultHelpCommand):
        async def send_pages(self):
            """A helper utility to send the page output
            from :attr:`paginator` to the destination."""
            destination = self.get_destination()
            for page in self.paginator.pages:
                await destination.send(page, delete_after=timedelta.super_big)

            await destination.send(
                (f'the message is going to be auto destroyed'
                 f' in {timedelta.super_big} seconds'),
                delete_after=timedelta.super_big)

    bot.help_command = MyHelpCommand()

    @bot.event
    async def on_ready():
        print('We have logged in as {0.user}'.format(bot))

    @bot.event
    async def on_command_error(ctx, error):
        print(f'ERR: {error}')
        await ctx.send(f'ERR: {error}', delete_after=timedelta.medium)

    @bot.command(name='connect')
    @commands.check_any(all_checks())
    async def connect_the_channel(ctx):
        "connects to channel"
        if str(ctx.channel.id) not in storage.channels:
            storage.channels[str(ctx.channel.id)] = {}
            await ctx.send('connected', delete_after=timedelta.medium)
        else:
            await ctx.send('we are already connected',
                           delete_after=timedelta.medium)

    @bot.command(name='disconnect')
    @commands.check_any(all_checks())
    async def diconnect_the_channel(ctx):
        "disconnects from channel"
        if str(ctx.channel.id) in storage.channels:
            storage.channels.pop(str(ctx.channel.id))
            await ctx.send('disconnected', delete_after=timedelta.medium)
        else:
            await ctx.send('we are already disconnected',
                           delete_after=timedelta.medium)

    @bot.command(name='config')
    @commands.check_any(connected_to_channel(storage))
    async def check_number(ctx):
        with open('templates/json.md') as file_:
            template = Template(file_.read())

            await ctx.send(
                template.render(data=json.dumps(
                    storage.channels[str(ctx.channel.id)], indent=2)))

    @check_number.error
    async def check_error(ctx, error):
        if isinstance(error, commands.CommandError):
            await ctx.send('incorrect number!', delete_after=timedelta.medium)

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
        await ctx.send(response, delete_after=timedelta.medium)

    @bot.command(name='me')
    @commands.check_any(commands.is_owner())
    async def only_me(ctx, ):
        "secret command"
        await ctx.send('Papa!', delete_after=timedelta.big)

    @bot.command(name='info')
    @commands.check_any(all_checks())
    async def more_detailed_info(ctx):
        "more detailed help"
        with open('templates/info.md') as file_:
            await ctx.send(file_.read(), delete_after=timedelta.super_big)

    return bot
