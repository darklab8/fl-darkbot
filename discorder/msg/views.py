from discorder.msg import queries
from fastapi import APIRouter
from fastapi import Query, Path, Body
from starlette.requests import Request
from fastapi import Depends
from utils.rest_api.message import MessageOk

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

    bot = request.app.discord_bot
    print(repr(bot.get_channel(query.channel_id)))
    await bot.get_channel(query.channel_id).send(query.message)
    return MessageOk()
