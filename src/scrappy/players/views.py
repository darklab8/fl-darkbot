from fastapi import APIRouter
from typing import Dict
from pydantic.dataclasses import dataclass
from . import repository
import core.databases as databases
from fastapi import Depends
from sqlalchemy.orm import Session


router = APIRouter(
    prefix="/players",
    tags=["items"],
)


pong = {"message": "pong!"}


@dataclass
class Pong:
    message: str = "pong!"


@router.get("/", response_model=Pong)
async def get_ping2():
    return pong


@router.get("/players")
async def get_players(db: Session = Depends(databases.default.get_session)):
    player_storage = repository.PlayerRepository()

    players = player_storage.get_all(db)
    return players
