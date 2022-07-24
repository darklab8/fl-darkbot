from pydantic import BaseModel


class PlayerSchema(BaseModel):
    description: str = None
    id: int | None = None
