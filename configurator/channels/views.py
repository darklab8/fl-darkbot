from fastapi import APIRouter
from fastapi import Query, Path, Body

from fastapi import Depends
from ..core.databases import DatabaseFactory
from utils.database.sql import Database
from utils.rest_api.message import MessageOk
from . import actions
from . import storage
from . import schemas
from typing import Union


class Paths:
    base = "/channel"


router = APIRouter(
    prefix=Paths.base,
    tags=["items"],
)

query_default_values = schemas.ChannelCreateQueryParams(channel_id=0)


@router.post("", response_model=MessageOk)
async def register_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: schemas.ChannelCreateQueryParams = Body(),
):

    await actions.ActionRegisterChannel(
        db=database,
        query=query,
    ).run()

    return MessageOk()


@router.delete("", response_model=MessageOk)
async def delete_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: schemas.ChannelDeleteQueryParams = Body(),
):
    await actions.ActionDeleteChannel(
        db=database,
        channel_id=query.channel_id,
    ).run()
    return MessageOk()
