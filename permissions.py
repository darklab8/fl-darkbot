from discord.ext import commands


def predicate_guild_owner(ctx):
    return ctx.guild is not None and ctx.guild.owner_id == ctx.author.id


def is_guild_owner():
    return commands.check(predicate_guild_owner)


def predicate_manage(ctx):
    return ctx.guild is not None and ctx.author.guild_permissions.manage_channels


def can_manage_channels():
    return commands.check(predicate_manage)


def controller_predicate(ctx):
    return ctx.guild is not None and 'bot_controller' in [
        elem.name for elem in ctx.author.roles
    ]
    return True


def is_bot_controller():
    return commands.check(controller_predicate)


# async def predicate_bot_owner(ctx):
#     # TODO it should be awaited
#     return await ctx.bot.is_owner(ctx.author)

# def bow_owner():
#     return commands.check(predicate_bot_owner)


def all_checks():
    def predicate(ctx):
        return predicate_guild_owner(ctx) or predicate_manage(
            ctx) or controller_predicate(ctx)  # or predicate_bot_owner(ctx)

    return commands.check(predicate)
