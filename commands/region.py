from discord.ext import commands
from permissions import connected_to_channel
from .universal import execute_in_storage


def attach_region(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def region(ctx):
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @region.command(name='add', pass_context=True)
    @execute_in_storage(storage)
    async def region_add(ctx, *args):
        pass

    @region.command(name='remove', pass_context=True)
    @execute_in_storage(storage)
    async def region_remove(ctx, *args):
        pass

    return bot
