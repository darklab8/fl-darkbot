from discorder.msg import queries
from fastapi import APIRouter
from fastapi import Query, Path, Body
from starlette.requests import Request
from fastapi import Depends
from utils.rest_api.message import MessageOk
import discord

# from . import actions
# from . import storage
from . import queries
from . import schemas
from .urls import urls
from typing import Union


router = APIRouter(
    prefix="",
    tags=["msg"],
)

# query_default_values = actions.ActionRegisterChannel.query_factory(channel_id=0)


@router.post(urls.base)
async def send_or_replace_msg(
    request: Request,
    query: queries.CreateOrReplaceMessqgeQueryParams = Body(),
):

    bot: discord.Client = request.app.discord_bot
    channel = bot.get_channel(query.channel_id)
    messages = channel.history(limit=20)

    async for msg in messages:
        if query.id in msg.content:
            print(f"query.id={query.id}, msg.content={msg.content}")
            await msg.edit(content=query.message)
            return MessageOk()

    await channel.send(query.message)
    return MessageOk()


@router.delete(urls.base)
async def delete_msg(
    request: Request,
    query: queries.DeleteMessageQueryParams = Body(),
):

    bot: discord.Client = request.app.discord_bot
    channel = bot.get_channel(query.channel_id)
    messages = channel.history(limit=20)

    async for msg in messages:
        if query.id in msg.content:
            await msg.delete()

    return MessageOk()
