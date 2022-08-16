from utils.porto import AbstractAction
import abc
from scrappy.core.logger import base_logger

logger = base_logger.getChild(__name__)


class ActionGetAndParseAndSaveItems(AbstractAction):
    @abc.abstractproperty
    def task_get(self) -> AbstractAction:
        pass

    @abc.abstractproperty
    def task_parse(self) -> AbstractAction:
        pass

    @abc.abstractproperty
    def task_save(self) -> AbstractAction:
        pass

    def __init__(self, database):
        self._database = database

    def run(self):
        item_data = self.task_get()
        items = self.task_parse(item_data)
        self.task_save(items=items, database=self._database)
        logger.debug(f"{self.__class__.__name__} is done")
        return items


class ActionGetFilteredItems(AbstractAction):
    @abc.abstractproperty
    def queryparams(self):
        pass

    @abc.abstractproperty
    def storage(self):
        pass

    def __init__(self, database, **kwargs):
        self._database = database
        self.query = self.queryparams(**kwargs)

    def run(self):
        storage = self.storage(self._database)
        players = storage.get(self.query)
        return players
