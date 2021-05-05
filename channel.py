"functions handling work with channels"
import datetime
import discord
from jinja2 import Template


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


async def give_me_last_tagged_channels(bot, unique_tag: str, channel_id: int):
    content_search = await bot.get_channel(channel_id).history(limit=200
                                                               ).flatten()
    return [item for item in content_search if unique_tag in item.content]


async def handle_tagged_messages(bot, unique_tag: str, channel_id: int):
    channels = await give_me_last_tagged_channels(bot, unique_tag, channel_id)

    if not channels:
        # create first msg
        await bot.get_channel(channel_id).send(unique_tag +
                                               str(datetime.datetime.utcnow()))
    elif len(channels) > 1:
        # delete all others
        deleting = channels[1:]
        for message in deleting:
            await deleting_message(message)
    else:
        # edit to apply tag
        with open('date.md') as file_:
            template = Template(file_.read())

            await channels[0].edit(content=str(
                template.render(date=str(datetime.datetime.utcnow()))))
