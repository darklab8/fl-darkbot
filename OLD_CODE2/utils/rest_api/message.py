from pydantic import BaseModel


class MessageOk(BaseModel):
    result: str = "OK"
