from discord.ext import commands

from .permissions import connected_to_channel

from .universal import execute_in_storage


def attach_enemy(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def enemy(ctx):
        "add your enemies to enemy tracking list!"
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @enemy.command(name='add', pass_context=True)
    @execute_in_storage(storage)
    async def enemy_add(ctx, *args):
        pass

    @enemy.command(name='remove', pass_context=True)
    @execute_in_storage(storage)
    async def enemy_remove(ctx, *args):
        pass

    @enemy.command(name='clear', pass_context=True)
    @execute_in_storage(storage)
    async def enemy_clear(ctx, *args):
        pass

    @enemy.command(name='alert', pass_context=True)
    @execute_in_storage(storage)
    async def enemy_alert(ctx, *args):
        pass

    return bot
