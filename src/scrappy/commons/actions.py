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
        player_data = self.task_get()
        players = self.task_parse(player_data)
        self.task_save(players=players, database=self._database)
        logger.debug(f"{self.__class__.__name__} is done")
        return players
