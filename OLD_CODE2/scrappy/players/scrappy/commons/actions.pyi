import abc
from _typeshed import Incomplete
from pydantic import BaseModel as BaseModel
from scrappy.commons.storage import AbstractStorage as AbstractStorage
from scrappy.core.logger import base_logger as base_logger
from typing import Any, Type
from utils.database.sql import Database as Database
from utils.porto import AbstractAction

logger: Incomplete

class ActionGetAndParseAndSaveItems(AbstractAction, metaclass=abc.ABCMeta):
    @abc.abstractmethod
    def task_get(self) -> dict[str, Any]: ...
    @abc.abstractmethod
    def task_parse(self, data: dict[str, Any]) -> list[Any]: ...
    @abc.abstractmethod
    def task_save(self, items: list[Any], database: Database) -> bool: ...
    def __init__(self, database: Database) -> None: ...
    def run(self) -> Any: ...

class ActionGetFilteredItems(AbstractAction, metaclass=abc.ABCMeta):
    @property
    @abc.abstractmethod
    def queryparams(self) -> Type[BaseModel]: ...
    @property
    @abc.abstractmethod
    def storage(self) -> Type[AbstractStorage]: ...
    query: Incomplete
    def __init__(self, database: Database, **kwargs: dict[str, Any]) -> None: ...
    def run(self) -> list[Any]: ...
