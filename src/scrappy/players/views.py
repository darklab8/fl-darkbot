from fastapi import APIRouter
from typing import Dict
from fastapi import Query

from fastapi import Depends
from scrappy.core.databases import DatabaseFactory, Database
from . import actions as player_actions


router = APIRouter(
    prefix="/players",
    tags=["items"],
)

query_default_values = player_actions.PlayerQuery()


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
