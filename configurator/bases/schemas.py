from pydantic import BaseModel
from typing import List


class BaseRegisterRequestParams(BaseModel):
    channel_id: int
    base_tags: list[str] = []


class BaseDeleteRequestParams(BaseModel):
    channel_id: int


class BaseGetRequestParams(BaseModel):
    pass


class BaseOut(BaseModel):
    tags: list[str]
    channel_id: int


class BasesManyOut(BaseModel):
    __root__: List[BaseOut]

    def __iter__(self):
        for item in self.__root__:
            yield item
