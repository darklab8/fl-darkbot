from pydantic import BaseModel
from datetime import datetime


class BaseIn(BaseModel):
    pass


class BaseOut(BaseModel):
    pass


class BaseQueryParams(BaseModel):
    pass
