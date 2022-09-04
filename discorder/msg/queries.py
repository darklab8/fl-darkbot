from pydantic import BaseModel


class CreateOrReplaceMessqgeQueryParams(BaseModel):
    id: str
    channel_id: int
    message: str


class DeleteMessageQueryParams(BaseModel):
    id: str
    channel_id: int
