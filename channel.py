"""functions handling work with channels"""
import datetime

import discord


async def msg_deleter(message):
    try:
        await message.delete()
        return True
    except discord.errors.DiscordException as error:
        print(f"{str(datetime.datetime.utcnow())} "
              f"ERR  {str(error)} for channel: {str(message.channel.id)}"
              f"can't delete msg {str(message.content)}")
        return False


class ChannelConstroller():
    def __init__(self, bot, unique_tag: str, deleting_message=msg_deleter):
        self.bot = bot
        self.unique_tag = unique_tag
        self.deleting_message = deleting_message

    async def delete_exp_msgs(self, channel_id: int, time_msg_expiration: int):
        messages = await self.bot.get_channel(channel_id).history(
            limit=200,
            before=datetime.datetime.utcnow() -
            datetime.timedelta(seconds=time_msg_expiration)).flatten()

        for message in messages:
            # if message.author.id == self.bot.user.id:
            if self.unique_tag not in message.content:
                await self.deleting_message(message)

    async def get_tagged_msgs(self, channel_id: int):
        content_search = await self.bot.get_channel(channel_id).history(
            limit=200).flatten()
        return [
            item for item in content_search if self.unique_tag in item.content
            and self.bot.user.id == item.author.id
        ]

    async def update_info(self, channel_id: int, info: str):
        messages = await self.get_tagged_msgs(channel_id)

        if not messages:
            # create first msg
            await self.bot.get_channel(channel_id).send(self.unique_tag +
                                                        ' forming the message')
        elif len(messages) > 1:
            # delete all others
            deleting = messages[1:]
            for message in deleting:
                await self.deleting_message(message)
        else:
            # edit to apply tag
            await messages[0].edit(content=str(info))
