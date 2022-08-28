from pydantic import BaseModel
import datetime


class ChannelOut(BaseModel):
    channel_id: int


class ChannelOwnerOut(BaseModel):
    id: int
    name: str
    channel_id: int
    created: datetime.datetime


class ChannelCreateQueryParams(BaseModel):
    channel_id: int
    owner_id: int | None = None
    owner_name: str = ""


class ChannelDeleteQueryParams(BaseModel):
    channel_id: int
