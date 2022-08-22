from fastapi import APIRouter
from fastapi import Query, Path, Body

from fastapi import Depends
from ..core.databases import DatabaseFactory
from utils.database.sql import Database
from . import actions
from . import storage
from . import schemas
from typing import Union
from pydantic import BaseModel

router = APIRouter(
    prefix="/channels",
    tags=["items"],
)

query_default_values = schemas.ChannelQueryParams(channel_id=0)


class MessageOk(BaseModel):
    result: str = "OK"


@router.post("/{channel_id}", response_model=MessageOk)
async def register_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    channel_id: int = Path(),
    owner_id: Union[int, None] = Body(default=query_default_values.owner_id),
    owner_name: str = Body(default=query_default_values.owner_name),
):

    query = schemas.ChannelQueryParams(
        channel_id=channel_id,
        owner_id=owner_id,
        owner_name=owner_name,
    )
    await actions.ActionRegisterChannel(
        db=database,
        query=query,
    ).run()

    return MessageOk()


@router.delete("/{channel_id}", response_model=MessageOk)
async def delete_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    channel_id: int = Path(),
):
    await actions.ActionDeleteChannel(
        db=database,
        channel_id=channel_id,
    ).run()
    return MessageOk()
