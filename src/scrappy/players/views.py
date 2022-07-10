from fastapi import APIRouter
from typing import Dict
from pydantic.dataclasses import dataclass

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
