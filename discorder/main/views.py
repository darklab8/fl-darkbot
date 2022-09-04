from fastapi import APIRouter
from starlette.requests import Request
from utils.rest_api.message import MessageOk

from typing import Union


router = APIRouter(
    prefix="",
    tags=["main"],
)


@router.get("/ping")
async def ping(request: Request):
    return MessageOk()


@router.get("/guilds")
async def get_guilds(request: Request):
    bot = request.app.discord_bot
    return {"guilds": [guild.id for guild in bot.guilds]}
