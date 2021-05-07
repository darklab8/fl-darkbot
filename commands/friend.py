from discord.ext import commands
from permissions import connected_to_channel
from .universal import execute_in_storage


def attach_friend(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def friend(ctx):
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @friend.command(name='add', pass_context=True)
    @execute_in_storage(storage)
    async def friend_add(ctx, *args):
        pass

    @friend.command(name='remove', pass_context=True)
    @execute_in_storage(storage)
    async def friend_remove(ctx, *args):
        pass

    @friend.command(name='clear', pass_context=True)
    @execute_in_storage(storage)
    async def friend_clear(ctx, *args):
        pass

    @friend.command(name='alert', pass_context=True)
    @execute_in_storage(storage)
    async def friend_alert(ctx, *args):
        pass

    return bot
