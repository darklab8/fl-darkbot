from pydantic import BaseModel
from datetime import datetime


class PlayerSchema(BaseModel):
    id: int | None = None
    name: str
    region: str
    system: str
    time: str
    timestamp: datetime
