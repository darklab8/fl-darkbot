"starting module"
from commands import attach_commands
from storage import Storage
from looper import Looper

# nice settings loading


def create_app():
    storage = Storage()
    bot = attach_commands(storage)
    _ = Looper(bot, storage)
    return bot, storage.settings.secret_key


if __name__ == '__main__':
    bot, secret_key = create_app()
    bot.run(secret_key)
