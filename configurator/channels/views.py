from fastapi import APIRouter
from fastapi import Query, Path, Body

from fastapi import Depends
from ..core.databases import DatabaseFactory
from utils.database.sql import Database
from . import actions
from . import storage
from . import schemas
from typing import Union

router = APIRouter(
    prefix="/channels",
    tags=["items"],
)

query_default_values = schemas.ChannelQueryParams(channel_id=0)


@router.post("/{channel_id}", response_model=schemas.ChannelQueryParams)
async def register_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    channel_id: int = Path(),
    owner_id: Union[int, None] = Body(default=query_default_values.owner_id),
    owner_name: Union[str, None] = Body(default=query_default_values.owner_name),
):

    query_params = schemas.ChannelQueryParams(
        channel_id=channel_id,
        owner_id=owner_id,
        owner_name=owner_name,
    )
    return query_params


@router.delete("/{channel_id}", response_model=schemas.ChannelQueryParams)
async def delete_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
    channel_id: int = Path(),
    owner_id: Union[int, None] = Body(default=query_default_values.owner_id),
    owner_name: Union[str, None] = Body(default=query_default_values.owner_name),
):

    query_params = schemas.ChannelQueryParams(
        channel_id=channel_id,
        owner_id=owner_id,
        owner_name=owner_name,
    )
    return query_params
