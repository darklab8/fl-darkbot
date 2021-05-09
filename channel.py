"""functions handling work with channels"""
import datetime
from abc import ABC, abstractmethod

import discord


class IMessageBus(ABC):
    @abstractmethod
    def __init__(self):
        pass

    @abstractmethod
    def bot_user_id(self, test: str) -> int:
        pass

    @abstractmethod
    async def delete(self, message):
        pass

    @abstractmethod
    async def send(self, channel_id, content):
        pass

    @abstractmethod
    async def history(self, channel_id, older_than_n_seconds=0):
        pass


class DiscordMessageBus(IMessageBus):
    def __init__(self, bot):
        self.bot = bot

    def bot_user_id(self):
        return self.bot.user.id

    async def delete(self, message):
        try:
            await message.delete()
        except discord.errors.DiscordException as error:
            print(f"{str(datetime.datetime.utcnow())} "
                  f"ERR  {str(error)} for channel: {str(message.channel.id)}"
                  f"can't delete msg {str(message.content)}")

    async def send(self, channel_id, content):
        return await self.bot.get_channel(channel_id).send(content)

    async def history(self, channel_id, older_than_n_seconds=0):
        return await self.bot.get_channel(channel_id).history(
            limit=200,
            before=datetime.datetime.utcnow() -
            datetime.timedelta(seconds=older_than_n_seconds)).flatten()


class ChannelConstroller():
    def __init__(self, message_bus: DiscordMessageBus, unique_tag: str):
        self.message_bus = message_bus
        self.unique_tag = unique_tag

    async def delete_exp_msgs(self, channel_id: int, time_msg_expiration: int):
        messages = await self.message_bus.history(channel_id,
                                                  time_msg_expiration)

        for message in messages:
            if self.unique_tag not in message.content:
                await self.message_bus.delete(message)

    async def get_tagged_msgs(self, channel_id: int):
        content_search = await self.message_bus.history(channel_id)
        return [
            item for item in content_search if self.unique_tag in item.content
            and self.message_bus.bot_user_id() == item.author.id
        ]

    async def update_info(self, channel_id: int, info: str):
        messages = await self.get_tagged_msgs(channel_id)

        if not messages:
            # create first msg
            await self.message_bus.history(
                channel_id, self.unique_tag + ' forming the message')
        elif len(messages) > 1:
            # delete all others
            deleting = messages[1:]
            for message in deleting:
                await self.message_bus.delete(message)
        else:
            # edit to apply tag
            await messages[0].edit(content=str(info))
