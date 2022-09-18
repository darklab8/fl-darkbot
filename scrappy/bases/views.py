from fastapi import APIRouter
from fastapi import Query, Body

from fastapi import Depends
from scrappy.core.databases import DatabaseFactory
from utils.database.sql import Database
from . import actions as bases_actions
from .schemas import BaseQueryParams, BaseOut

router = APIRouter(
    prefix="/bases",
    tags=["items"],
)


@router.post("")
async def get_bases(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: BaseQueryParams = Body(),
) -> list[BaseOut]:

    bases = bases_actions.ActionGetFilteredBases(
        database=database,
        query=query,
    )
    return bases
