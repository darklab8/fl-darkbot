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


router = APIRouter(
    prefix="",
    tags=["items"],
)

query_default_values = actions.ActionRegisterChannel.query_factory(channel_id=0)


@router.post(actions.ActionRegisterChannel.url, response_model=MessageOk)
async def register_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: actions.ActionRegisterChannel.query_factory = Body(),
):

    await actions.ActionRegisterChannel(
        db=database,
        query=query,
    ).run()

    return MessageOk()


@router.delete(actions.ActionRegisterChannel.url, response_model=MessageOk)
async def delete_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: actions.ActionDeleteChannel.query_factory = Body(),
):
    await actions.ActionDeleteChannel(
        db=database,
        query=query,
    ).run()
    return MessageOk()
