from utils.porto import AbstractAction
from pydantic import BaseModel
import abc
import requests
import scrappy.core.settings as settings
from typing import Type
from scrappy.core.logger import base_logger

logger = base_logger.getChild(__name__)


class StubSchema(BaseModel):
    pass


class SubTaskGetItemsData(AbstractAction):
    @abc.abstractproperty
    def _url(self):
        pass

    def run(self):
        logger.info(f"{self.__class__.__name__} is started")
        response = requests.get(self._url)
        data = response.json()
        logger.debug(f"{self.__class__.__name__} is done")
        return data


class StubStorage(abc.ABC):
    @abc.abstractmethod
    def create(self, *items):
        pass


class SubTaskSaveItemsToStorage(AbstractAction):
    @abc.abstractproperty
    def storage(self) -> Type[StubStorage]:
        pass

    def __init__(self, items: list[StubSchema], database):
        self._items = items
        self._database = database

    def run(self):
        storage = self.storage(self._database)
        storage.create(*(self._items))
        logger.debug(f"{self.__class__.__name__} is done")
        return True