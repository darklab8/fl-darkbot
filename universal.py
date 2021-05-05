import datetime
import discord


async def deleting_message(message):
    try:
        await message.delete()
    except discord.errors.NotFound:
        print('channel was already deleted')


async def delete_messages_older_than_n_seconds(bot, unique_tag: str, n: int,
                                               channel_id: int):
    messages = await bot.get_channel(channel_id).history(
        limit=200,
        before=datetime.datetime.utcnow() -
        datetime.timedelta(seconds=n)).flatten()

    for message in messages:
        # if message.author.id == self.bot.user.id:
        if unique_tag not in message.content:
            await deleting_message(message)
