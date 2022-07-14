from fastapi import APIRouter
from typing import Dict
from pydantic.dataclasses import dataclass
from fastapi import Depends
from sqlalchemy.orm import Session

from scrappy.players.repository import PlayerRepository
import scrappy.core.databases as databases


router = APIRouter(
    prefix="/players",
    tags=["items"],
)

@router.get("/")
async def get_players(db: Session = Depends(databases.default.get_session)):
    player_storage = PlayerRepository(db)

    players = player_storage.get_all()
    return players
