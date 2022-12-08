from fastapi import APIRouter
from fastapi import Query, Body
from fastapi import Depends, Response
from scrappy.core.databases import DatabaseFactory
from utils.database.sql import Database
from . import actions as player_actions
from .storage import PlayerStorage
from .schemas import PlayerQueryParams, PlayerOut

router = APIRouter(
    prefix="/players",
    tags=["items"],
)

@router.get("")
async def get_players(
    database: Database = Depends(DatabaseFactory.get_default_database),
    query: PlayerQueryParams = Body(),
) -> list[PlayerOut]:

    players = player_actions.ActionGetFilteredPlayers(  # type: ignore
        database=database,
        query=query,
    )
    return players


# purely test to try async
@router.get("-async")
async def get_async(
    database: Database = Depends(DatabaseFactory.get_default_database),
) -> list[PlayerOut]:
    async with database.get_async_session() as session:

        repo = PlayerStorage(database)
        players = await repo._a_get_all()
        return players
