from pydantic import BaseModel


class ChannelQueryParams(BaseModel):
    channel_id: int
    owner_id: int | None = None
    owner_name: str | None = None
