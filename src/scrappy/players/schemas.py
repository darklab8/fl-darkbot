from pydantic import BaseModel
from datetime import datetime


class PlayerIn(BaseModel):
    name: str
    region: str
    system: str
    time: str
    timestamp: datetime


class PlayerOut(BaseModel):
    id: int
    name: str
    region: str
    system: str
    time: str
    timestamp: datetime
    is_online: bool


class PlayerQueryParams(BaseModel):
    page: int = 0
    player_whitelist_tags: list[str] = []
    region_whitelist_tags: list[str] = []
    system_whitelist_tags: list[str] = []
    is_online: bool = True
