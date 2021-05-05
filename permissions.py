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


def predicate_owner(ctx):
    return ctx.bot.owner_id == ctx.author.id


# async def predicate_bot_owner(ctx):
#     # TODO it should be awaited
#     return await ctx.bot.is_owner(ctx.author)

# def bow_owner():
#     return commands.check(predicate_bot_owner)


def predicate_connected_to_channels(ctx, storage):
    return (ctx.channel.id) in storage.channels


def predicate_all_permissions(ctx):
    return predicate_owner(ctx) or predicate_guild_owner(
        ctx) or predicate_manage(ctx) or controller_predicate(
            ctx)  # or predicate_bot_owner(ctx)


def connected_to_channel(storage):
    def predicate_connected_and_permissions(ctx):
        return predicate_all_permissions(
            ctx) and predicate_connected_to_channels(ctx, storage)

    return commands.check(predicate_connected_and_permissions)


def all_checks():
    return commands.check(predicate_all_permissions)