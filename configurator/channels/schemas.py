from pydantic import BaseModel


class ChannelOut(BaseModel):
    channel_id: int


class ChannelQueryParams(BaseModel):
    channel_id: int
    owner_id: int | None = None
    owner_name: str = ""
