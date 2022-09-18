from pydantic import BaseModel, Field
from datetime import datetime
from typing import List

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

class BasesOut(BaseModel):
    __root__: List[BaseOut]
    def __iter__(self):
        for item in self.__root__:
            yield item

class BaseQueryParams(BaseModel):
    page: int = 0
    name_tags: list[str] = []
    page_size: int = 20
