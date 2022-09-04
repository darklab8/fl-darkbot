from pydantic import BaseModel


class CreateOrReplaceMessqgeQueryParams(BaseModel):
    identificator = str
    channel_id: int
    message: str


class DeleteMessageQueryParams(BaseModel):
    identificator = str
    channel_id: int
