from pydantic import BaseModel


class PlayerSchema(BaseModel):
    description: str | None = None

    class Config:
        orm_mode = True
