from pydantic import BaseModel


class PlayerSchema(BaseModel):
    description: str | None = None
    id: int | None = None
