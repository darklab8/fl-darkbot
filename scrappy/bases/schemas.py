from pydantic import BaseModel, Field
from datetime import datetime


class BaseIn(BaseModel):
    name: str
    affiliation: str
    health: float
    tid: int
    timestamp: datetime = Field(default_factory=datetime.utcnow)


class BaseOut(BaseModel):
    id: int
    name: str
    affiliation: str
    health: float
    tid: int
    timestamp: datetime


class BaseQueryParams(BaseModel):
    page: int = 0
    name_tags: list[str] = []


def check_thing2() -> None:
    thing = BaseQueryParams()
    print(thing.not_existing)
