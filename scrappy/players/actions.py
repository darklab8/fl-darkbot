from .storage import PlayerStorage
from scrappy.commons.actions import (
    ActionGetAndParseAndSaveItems,
    ActionGetFilteredItems,
)
from .subtasks import (  # type: ignore
    SubTaskGetPlayerData,
    SubTaskParsePlayers,
    SubTaskSavePlayersToStorage,
)
from scrappy.core.logger import base_logger
from .schemas import PlayerQueryParams

logger = base_logger.getChild(__name__)


class ActionGetAndParseAndSavePlayers(ActionGetAndParseAndSaveItems):
    task_get = SubTaskGetPlayerData
    task_parse = SubTaskParsePlayers
    task_save = SubTaskSavePlayersToStorage


class ActionGetFilteredPlayers(ActionGetFilteredItems):
    queryparams = PlayerQueryParams
    storage = PlayerStorage
