from discorder.msg import queries
from fastapi import APIRouter
from fastapi import Body
from starlette.requests import Request
from utils.rest_api.message import MessageOk
import discord

from . import actions

# from . import storage
from . import queries
from . import schemas
from .urls import urls
from typing import Union


router = APIRouter(
    prefix="",
    tags=["msg"],
)


@router.post(actions.CreateOrReplaceMessage.url)
async def send_or_replace_msg(
    request: Request,
    query: queries.CreateOrReplaceMessqgeQueryParams = Body(),
):

    bot: discord.Client = request.app.discord_bot
    response = await actions.CreateOrReplaceMessage(
        bot=bot,
        query=query,
    ).run()
    return response


@router.delete(actions.DeleteMessage.url)
async def delete_msg(
    request: Request,
    query: queries.DeleteMessageQueryParams = Body(),
):

    bot: discord.Client = request.app.discord_bot
    response = await actions.DeleteMessage(
        bot=bot,
        query=query,
    ).run()
    return response
