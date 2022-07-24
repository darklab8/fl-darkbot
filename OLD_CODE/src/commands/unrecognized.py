from discord.ext import commands

from .permissions import connected_to_channel

from .universal import execute_in_storage


def attach_unrecognized(bot, storage):
    """Add commands to regulate unsorted players."""
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def unrecognized(ctx):
        """Set alert level for unrecognized players."""
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @unrecognized.command(name='alert', pass_context=True)
    @execute_in_storage(storage)
    async def unrecognized_alert(ctx, *args):
        pass

    return bot
