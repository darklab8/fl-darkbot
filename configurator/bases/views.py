from fastapi import APIRouter
from fastapi import Query, Path, Body

from fastapi import Depends
from ..core.databases import DatabaseFactory
from utils.database.sql import Database
from utils.rest_api.message import MessageOk
from . import actions
from . import storage
from . import schemas
from typing import Union, List
from pydantic import Field, BaseModel
from ..channels.views import Paths as ChannelPaths


class Paths:
    base = f"{ChannelPaths.base}/base"


router = APIRouter(
    prefix=Paths.base,
    tags=["items"],
)

baseviewinput_defaults = schemas.BaseRegisterRequestParams(channel_id=0)


class BaseBodyInput(BaseModel):
    base_tags: list[str] = Field(default=[])


@router.post("", response_model=MessageOk)
async def register_base(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: schemas.BaseRegisterRequestParams = Body(),
):
    await actions.ActionRegisterBase(db=database, query=query).run()

    return MessageOk()


@router.delete("", response_model=MessageOk)
async def delete_base(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: schemas.BaseDeleteRequestParams = Body(),
):

    await actions.ActionDeleteBases(
        db=database,
        channel_id=query.channel_id,
    ).run()

    return MessageOk()


# TODO add alarm
# @router.post("/alarm", response_model=MessageOk)
# async def register_base_alarms(
#     database: Database = Depends(DatabaseFactory.get_default_database),
#     channel_id: int = Path(),
#     alarm: bool = Query(default=baseviewinput_defaults.base_tags),
# ):

#     query = schemas.BaseAlarmViewInput(
#         channel_id=channel_id,
#         alarm=alarm,
#     )
#     await actions.ActionRegisterBaseAlarms(
#         db=database,
#         query=query,
#     ).run()

#     return MessageOk()
