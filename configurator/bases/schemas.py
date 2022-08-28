from pydantic import BaseModel
import datetime


class BaseRegisterRequestParams(BaseModel):
    channel_id: int
    base_tags: list[str] = []


class BaseDeleteRequestParams(BaseModel):
    channel_id: int


class BaseOut(BaseModel):
    tags: list[str]
    channel_id: int
