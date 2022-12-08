import discord
from utils.porto import AsyncAbstractAction
from utils.rest_api.message import MessageOk
from utils.rest_api.methods import RequestMethod
from . import queries
from .urls import urls


class CreateOrReplaceMessage(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.post
    query_factory = queries.CreateOrReplaceMessqgeQueryParams
    response_factory = MessageOk

    def __init__(
        self,
        bot: discord.Client,
        query: queries.CreateOrReplaceMessqgeQueryParams,
    ):
        self.bot = bot
        self.query = query

    async def run(self) -> MessageOk:
        channel = self.bot.get_channel(self.query.channel_id)
        messages = channel.history(limit=20)

        async for msg in messages:
            if self.query.id in msg.content:
                await msg.edit(content=self.query.message)
                return MessageOk()

        await channel.send(self.query.message)
        return MessageOk()


class DeleteMessage(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.delete
    query_factory = queries.DeleteMessageQueryParams
    response_factory = MessageOk

    def __init__(
        self,
        bot: discord.Client,
        query: queries.DeleteMessageQueryParams,
    ):
        self.bot = bot
        self.query = query

    async def run(self) -> MessageOk:
        channel = self.bot.get_channel(self.query.channel_id)
        messages = channel.history(limit=20)

        async for msg in messages:
            if self.query.id in msg.content:
                await msg.delete()

        return MessageOk()
