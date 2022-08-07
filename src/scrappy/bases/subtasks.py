from utils.porto import AbstractAction
from .storage import BaseStorage
from . import schemas as base_schemas
import scrappy.core.settings as settings
from scrappy.commons.subtasks import SubTaskGetItemsData, SubTaskSaveItemsToStorage
from scrappy.core.logger import base_logger

logger = base_logger.getChild(__name__)


class SubTaskGetBaseData(SubTaskGetItemsData):
    def _url(self):
        return settings.API_BASE_URL


class SubTaskParseBases(AbstractAction):
    def __init__(self, data: dict):
        self._data = data

    def run(self) -> list[base_schemas.BaseIn]:
        # ALL YOUR BASE BELONG TO US
        bases = [
            # player_schemas.PlayerIn(**player, timestamp=self._data["timestamp"])
            # for player in self._data["players"]
        ]
        logger.debug(f"{self.__class__.__name__} is done")
        return bases


class SubTaskSaveBasesToStorage(SubTaskSaveItemsToStorage):
    storage = BaseStorage

    def __init__(self, players: list[base_schemas.BaseIn], database):
        super().__init__(items=players, database=database)
