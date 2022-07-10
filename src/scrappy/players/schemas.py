from pydantic import BaseModel


class PlayerSchema(BaseModel):
    title: str
    description: str | None = None

    class Config:
        orm_mode = True
