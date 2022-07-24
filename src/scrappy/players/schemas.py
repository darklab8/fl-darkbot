from pydantic import BaseModel


class PlayerSchema(BaseModel):
    id: int | None = None
    name: str
    region: str
    system: str
    time: str