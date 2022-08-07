from utils.porto import AbstractAction
from .storage import PlayerStorage
from . import schemas as player_schemas
import requests
import scrappy.core.settings as settings
from pydantic import BaseModel
from scrappy.commons.subtasks import SubTaskGetItemsData, SubTaskSaveItemsToStorage
from scrappy.commons.actions import ActionGetAndParseAndSaveItems
from scrappy.core.logger import base_logger

logger = base_logger.getChild(__name__)


class SubTaskGetPlayerData(SubTaskGetItemsData):
    def _url(self):
        return settings.API_PLAYER_URL


class SubTaskParsePlayers(AbstractAction):
    def __init__(self, data: dict):
        self._data = data

    def run(self) -> list[player_schemas.PlayerIn]:
        players = [
            player_schemas.PlayerIn(**player, timestamp=self._data["timestamp"])
            for player in self._data["players"]
        ]
        logger.debug(f"{self.__class__.__name__} is done")
        return players


class SubTaskSavePlayersToStorage(SubTaskSaveItemsToStorage):
    storage = PlayerStorage

    def __init__(self, players: list[player_schemas.PlayerIn], database):
        super().__init__(items=players, database=database)


class ActionGetAndParseAndSavePlayers(ActionGetAndParseAndSaveItems):
    task_get = SubTaskGetPlayerData
    task_parse = SubTaskParsePlayers
    task_save = SubTaskSavePlayersToStorage


class PlayerQuery(BaseModel):
    page: int = 0
    player_whitelist_tags: list[str] = []
    region_whitelist_tags: list[str] = []
    system_whitelist_tags: list[str] = []
    is_online: bool = True


class ActionGetFilteredPlayers(AbstractAction):
    def __init__(self, database, **kwargs):
        self._database = database
        self.query = PlayerQuery(**kwargs)

    def run(self):
        player_storage = PlayerStorage(self._database)
        players = player_storage.get_players_by_query(self.query)
        return players
