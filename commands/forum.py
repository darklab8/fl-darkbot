from discord.ext import commands

from permissions import connected_to_channel

from .universal import execute_in_storage


def attach_forum(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def forum(ctx):
        """add your threads to forum thread tracking list!\n
        write any tags from the thread for adding"""
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @forum.command(name='add', pass_context=True)
    @execute_in_storage(storage)
    async def forum_add(ctx, *args):
        pass

    @forum.command(name='remove', pass_context=True)
    @execute_in_storage(storage)
    async def forum_remove(ctx, *args):
        pass

    @forum.command(name='clear', pass_context=True)
    @execute_in_storage(storage)
    async def forum_clear(ctx, *args):
        pass

    return bot
