from .storage import BaseStorage
from pydantic import BaseModel
from scrappy.commons.actions import (
    ActionGetAndParseAndSaveItems,
    ActionGetFilteredItems,
)
from .subtasks import (
    SubTaskGetBaseData,
    SubTaskParseBases,
    SubTaskSaveBasesToStorage,
)
from scrappy.core.logger import base_logger
from scrappy.bases.schemas import BaseQueryParams

logger = base_logger.getChild(__name__)


class ActionGetAndParseAndSaveBases(ActionGetAndParseAndSaveItems):
    task_get = SubTaskGetBaseData
    task_parse = SubTaskParseBases
    task_save = SubTaskSaveBasesToStorage


class ActionGetFilteredBases(ActionGetFilteredItems):
    queryparams = BaseQueryParams
    storage = BaseStorage
