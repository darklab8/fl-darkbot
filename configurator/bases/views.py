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
from ..core.logger import base_logger
import json


logger = base_logger.getChild(__name__)

router = APIRouter(
    prefix="",
    tags=["items"],
)

baseviewinput_defaults = schemas.BaseRegisterRequestParams(channel_id=0)


class BaseBodyInput(BaseModel):
    base_tags: list[str] = Field(default=[])


@router.post(actions.ActionRegisterBase.url, response_model=MessageOk)
async def register_base(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: actions.ActionRegisterBase.query_factory = Body(),
):
    await actions.ActionRegisterBase(db=database, query=query).run()

    return MessageOk()


@router.delete(actions.ActionDeleteBases.url, response_model=MessageOk)
async def delete_base(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: actions.ActionDeleteBases.query_factory = Body(),
):

    await actions.ActionDeleteBases(
        db=database,
        query=query,
    ).run()

    return MessageOk()


@router.post(
    actions.ActionGetBases.url,
    response_model=List[schemas.BaseOut],
)
async def get_bases(
    database: Database = Depends(DatabaseFactory.get_default_database),
    # query: actions.ActionGetBases.query_factory = Body(),
):
    query = actions.ActionGetBases.query_factory()
    logger.debug(f"get_bases: starting view")
    logger.debug(f"get_bases={query}")
    logger.debug(f"get_bases: starting action")
    bases = await actions.ActionGetBases(db=database, query=query).run()
    logger.debug(f"get_bases: finished action")
    logger.debug(f"get_bases: finishing view")
    return list(bases)


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
