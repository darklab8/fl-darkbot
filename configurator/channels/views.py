from fastapi import APIRouter
from fastapi import Query

from fastapi import Depends
from scrappy.core.databases import DatabaseFactory
from utils.database.sql import Database
from . import actions as player_actions
from .storage import PlayerStorage
from .schemas import PlayerQueryParams

router = APIRouter(
    prefix="/players",
    tags=["items"],
)

query_default_values = PlayerQueryParams()


@router.get("/")
async def get_players(
    database: Database = Depends(DatabaseFactory.get_default_database),
    page: int = Query(default=query_default_values.page),
    player_tag: list[str] = Query(default=query_default_values.player_whitelist_tags),
    region_tag: list[str] = Query(default=query_default_values.region_whitelist_tags),
    system_tag: list[str] = Query(default=query_default_values.system_whitelist_tags),
    is_online: bool = Query(default=query_default_values.is_online),
):

    players = player_actions.ActionGetFilteredPlayers(
        database=database,
        page=page,
        player_whitelist_tags=player_tag,
        region_whitelist_tags=region_tag,
        system_whitelist_tags=system_tag,
        is_online=is_online,
    )
    return players


# purely test to try async
@router.get("/async")
async def register_channel(
    database: Database = Depends(DatabaseFactory.get_default_database),
):
    async with database.get_async_session() as session:

        repo = PlayerStorage(database)
        players = await repo._a_get_all()
        return players
