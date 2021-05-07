from discord.ext import commands
from permissions import connected_to_channel
from .universal import execute_in_storage


def attach_system(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def system(ctx):
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @system.command(name='add', pass_context=True)
    @execute_in_storage(storage)
    async def system_add(ctx, *args):
        pass

        # getattr(storage, )(ctx.channel.id, name)

    @system.command(name='remove', pass_context=True)
    @execute_in_storage(storage)
    async def system_remove(ctx, *args):
        pass

    return bot
