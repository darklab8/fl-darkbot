from pydantic import BaseModel
import abc
from typing import Any
from utils.database.sql import Database


class AbstractStorage(abc.ABC):
    def __init__(self, db: Database):
        self.db = db

    @abc.abstractmethod
    def create(self, *items: Any) -> None:
        pass

    @abc.abstractmethod
    def get(self, query: Any) -> list[Any]:
        pass
