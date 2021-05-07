import functools


def execute_in_storage(storage):
    def decorator_repeat(func):
        "useful decorator to preload stuff for unit tests"

        @functools.wraps(func)
        async def wrapper_repeat(*args, **kwargs):
            # print('executing '+func.__name__)

            methods = func.__name__.split('_')
            ctx = args[0]
            names = args[1:]

            await ctx.send(f'executing {methods} operation '
                           f'for objects {names}, mr {ctx.author.mention}')

            # getattr(storage, func.__name__)(args[0])
            return await func(*args, **kwargs)

        return wrapper_repeat

    return decorator_repeat
