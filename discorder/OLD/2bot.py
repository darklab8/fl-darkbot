from discord.ext.commands import Bot


class Singleton:
    @classmethod
    def get_bot(cls):
        if not hasattr(cls, "bot"):
            cls.bot = Bot(command_prefix="$")

        return cls.bot
