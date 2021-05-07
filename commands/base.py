from discord.ext import commands
from permissions import connected_to_channel
from .universal import execute_in_storage


def attach_base(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def base(ctx):
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @base.command(name='add', pass_context=True)
    @execute_in_storage(storage)
    async def base_add(ctx, *args):
        pass

    @base.command(name='remove', pass_context=True)
    @execute_in_storage(storage)
    async def base_remove(ctx, *args):
        pass

    @base.command(name='clear', pass_context=True)
    @execute_in_storage(storage)
    async def base_clear(ctx, *args):
        pass

    @base.command(name='list', pass_context=True)
    @execute_in_storage(storage)
    async def base_lst(ctx, *args):
        pass

    return bot
