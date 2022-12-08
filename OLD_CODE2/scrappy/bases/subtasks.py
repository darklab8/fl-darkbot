from utils.porto import AbstractAction
from .storage import BaseStorage
from . import schemas as base_schemas
import scrappy.core.settings as settings
from scrappy.commons.subtasks import SubTaskGetItemsData, SubTaskSaveItemsToStorage
from scrappy.core.logger import base_logger
from utils.database.sql import Database

logger = base_logger.getChild(__name__)


class SubTaskGetBaseData(SubTaskGetItemsData):
    @property
    def _url(self) -> str:
        return settings.API_BASE_URL


class SubTaskParseBases(AbstractAction):
    def __init__(self, data: dict):
        self._data = data

    def run(self) -> list[base_schemas.BaseIn]:
        # ALL YOUR BASE BELONG TO US
        bases = [
            base_schemas.BaseIn(
                name=base_name,
                affiliation=base_data.get("affiliation", "unknown"),
                health=float(base_data.get("health", -1.0)),
                tid=int(base_data.get("tid", "-1")),
            )
            for base_name, base_data in self._data.items()
        ]
        logger.debug(f"{self.__class__.__name__} is done")
        return bases


class SubTaskSaveBasesToStorage(SubTaskSaveItemsToStorage):
    storage = BaseStorage

    def __init__(self, items: list[base_schemas.BaseIn], database: Database):
        super().__init__(items=items, database=database)
