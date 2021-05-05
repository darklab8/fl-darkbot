"rules for command permissions"
from discord.ext import commands

# predicates


def predicate_guild_owner(ctx):
    return ctx.guild is not None and ctx.guild.owner_id == ctx.author.id


def predicate_manage(ctx):
    return (ctx.guild is not None
            and ctx.author.guild_permissions.manage_channels)


def predicate_controller(ctx):
    return ctx.guild is not None and 'bot_controller' in [
        elem.name for elem in ctx.author.roles
    ]
    return True


def predicate_owner(ctx):
    return ctx.bot.owner_id == ctx.author.id


def predicate_connected_to_channels(ctx, storage):
    return (ctx.channel.id) in storage.channels


def predicate_all_permissions(ctx):
    return predicate_owner(ctx) or predicate_guild_owner(
        ctx) or predicate_manage(ctx) or predicate_controller(
            ctx)  # or predicate_bot_owner(ctx)


# command checkers


def connected_to_channel(storage):
    def predicate_connected_and_permissions(ctx):
        return predicate_all_permissions(
            ctx) and predicate_connected_to_channels(ctx, storage)

    return commands.check(predicate_connected_and_permissions)


def all_checks():
    return commands.check(predicate_all_permissions)
