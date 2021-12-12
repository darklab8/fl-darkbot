"bot initialization + its commands"
import json

from discord.ext import commands

from jinja2 import Template

from .permissions import (
    all_checks,
    connected_to_channel,
)

from .universal import timedelta
import src.settings as settings


def attach_root(bot, storage, chanell_controller) -> commands.Bot:
    "attaching root commands to application"

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
    async def on_guild_join(guild):
        await guild.system_channel.send(
            "I am the terror that flaps in the night, I am the bed edge that"
            " stumbles your fingers! I am DARKWING DUCK!\n"
            "for a more information, visit https://dd84ai.github.io/darkbot/")

    @bot.event
    async def on_command_error(ctx, error):
        if settings.DEBUG:
            print(f'ERR: {error}')
        await ctx.send(f'ERR: {error}', delete_after=timedelta.medium)

    @bot.command(name='connect')
    @commands.check_any(all_checks())
    async def connect_the_channel(ctx):
        """connects to channel"""
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
    async def get_config(ctx):
        "shows you current config"
        with open('src/templates/json.md') as file_:
            template = Template(file_.read())

            await ctx.send(
                template.render(data=json.dumps(
                    storage.channels[str(ctx.channel.id)], indent=2)))

    @bot.command(name='clear')
    @commands.check_any(connected_to_channel(storage))
    async def clear(ctx):
        "clears all msgs from the current channel"
        await chanell_controller.delete_exp_msgs(ctx.channel.id, 0)

    @bot.command(name='wiki')
    @commands.check_any(all_checks())
    async def more_detailed_info(ctx):
        "gives link to wiki for a more detailed help"
        await ctx.send('**wiki**: https://dd84ai.github.io/darkbot/',
                       delete_after=timedelta.medium)

    @bot.command(name='me')
    @commands.check_any(commands.is_owner())
    async def only_me(ctx, ):
        "secret command"
        await ctx.send('Papa!', delete_after=timedelta.big)

    return bot
