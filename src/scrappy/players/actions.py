from utils.porto import AbstractAction
from .storage import PlayerStorage
from pydantic import BaseModel
from scrappy.commons.actions import ActionGetAndParseAndSaveItems
from .subtasks import (
    SubTaskGetPlayerData,
    SubTaskParsePlayers,
    SubTaskSavePlayersToStorage,
)
from scrappy.core.logger import base_logger

logger = base_logger.getChild(__name__)


class ActionGetAndParseAndSavePlayers(ActionGetAndParseAndSaveItems):
    task_get = SubTaskGetPlayerData
    task_parse = SubTaskParsePlayers
    task_save = SubTaskSavePlayersToStorage


class PlayerQueryParams(BaseModel):
    page: int = 0
    player_whitelist_tags: list[str] = []
    region_whitelist_tags: list[str] = []
    system_whitelist_tags: list[str] = []
    is_online: bool = True


class ActionGetFilteredPlayers(AbstractAction):
    def __init__(self, database, **kwargs):
        self._database = database
        self.query = PlayerQueryParams(**kwargs)

    def run(self):
        player_storage = PlayerStorage(self._database)
        players = player_storage.get(self.query)
        return players
