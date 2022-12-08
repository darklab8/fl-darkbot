from utils.porto import AbstractAction
import abc
import requests
from typing import Type
from scrappy.core.logger import base_logger
from typing import Any
from utils.database.sql import Database
from scrappy.commons.storage import AbstractStorage

logger = base_logger.getChild(__name__)


class SubTaskGetItemsData(AbstractAction):
    @abc.abstractproperty
    def _url(self) -> str:
        pass

    def run(self) -> dict[str, Any]:
        logger.info(f"{self.__class__.__name__} is started")
        response = requests.get(self._url)
        data = dict(response.json())
        logger.debug(f"{self.__class__.__name__} is done")
        return data


class SubTaskSaveItemsToStorage(AbstractAction):
    @abc.abstractproperty
    def storage(self) -> Type[AbstractStorage]:
        pass

    def __init__(self, items: list[Any], database: Database):
        self._items = items
        self._database = database

    def run(self) -> None:
        storage = self.storage(self._database)
        storage.create(*(self._items))
        logger.debug(f"{self.__class__.__name__} is done")
