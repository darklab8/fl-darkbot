from discord.ext import commands
from permissions import connected_to_channel


def attach_base(bot, storage) -> commands.Bot:
    @bot.group(pass_context=True)
    @commands.check_any(connected_to_channel(storage))
    async def base(ctx):
        if ctx.invoked_subcommand is None:
            await ctx.send('Invalid sub command passed...')

    @base.command(name='add', pass_context=True)
    async def base_add(ctx):
        await ctx.send(f'adding the base, mr {ctx.author.mention}')

    @base.command(name='remove', pass_context=True)
    async def base_remove(ctx):
        await ctx.send(f'removing the base, mr {ctx.author.mention}')

    return bot
