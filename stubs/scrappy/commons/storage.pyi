import abc
from pydantic import BaseModel as BaseModel
from typing import Any
from utils.database.sql import Database as Database

class AbstractStorage(abc.ABC, metaclass=abc.ABCMeta):
    def __init__(self, db: Database) -> None: ...
    @abc.abstractmethod
    def create(self, *items: list[Any]) -> list[Any]: ...
    @abc.abstractmethod
    def get(self, query: BaseModel) -> list[Any]: ...
