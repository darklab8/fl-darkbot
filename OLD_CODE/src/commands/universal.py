import functools
from types import SimpleNamespace

timedelta = SimpleNamespace(small=5, medium=10, big=20, super_big=40)


def execute_in_storage(storage):
    """Move main code being executed from discord.commands
    to my own code sections, which can be much easily tested
    """
    def decorator_repeat(func):
        @functools.wraps(func)
        async def wrapper_repeat(*args, **kwargs):
            # print('executing '+func.__name__)

            methods = func.__name__.split('_')
            category = methods[0]
            operation = methods[1]
            ctx = args[0]
            names = args[1:]

            await ctx.send(f'executing {methods} operation '
                           f'for objects {names}, {ctx.author.mention}')

            category_controller = getattr(storage, category)
            await getattr(category_controller, operation)(ctx.channel.id,
                                                          names)
            return await func(*args, **kwargs)

        return wrapper_repeat

    return decorator_repeat
