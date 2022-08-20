from utils.porto import AbstractAction
import abc
from scrappy.core.logger import base_logger
from typing import Any, Type
from utils.database.sql import Database
from scrappy.commons.storage import AbstractStorage
from pydantic import BaseModel

logger = base_logger.getChild(__name__)


class ActionGetAndParseAndSaveItems(AbstractAction):
    @abc.abstractmethod
    def task_get(self) -> dict[str, Any]:
        pass

    @abc.abstractmethod
    def task_parse(
        self,
        data: dict[str, Any],
    ) -> list[Any]:
        pass

    @abc.abstractmethod
    def task_save(self, items: list[Any], database: Database) -> bool:
        pass

    def __init__(self, database: Database):
        self._database = database

    def run(self) -> Any:
        item_data = self.task_get()
        items = self.task_parse(item_data)
        self.task_save(items=items, database=self._database)
        logger.debug(f"{self.__class__.__name__} is done")
        return items


class ActionGetFilteredItems(AbstractAction):
    @abc.abstractproperty
    def queryparams(self) -> Type[BaseModel]:
        pass

    @abc.abstractproperty
    def storage(self) -> Type[AbstractStorage]:
        pass

    def __init__(self, database: Database, **kwargs: dict[str, Any]):
        self._database = database
        self.query = self.queryparams(**kwargs)

    def run(self) -> list[Any]:
        storage = self.storage(self._database)
        players = storage.get(self.query)
        return players
