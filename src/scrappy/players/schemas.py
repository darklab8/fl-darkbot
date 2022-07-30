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
