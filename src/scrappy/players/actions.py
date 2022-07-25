from utils.porto import AbstractAction
from .repository import PlayerRepository
from .schemas import PlayerSchema
import requests
import scrappy.core.settings as settings
from pydantic import BaseModel


class SubTaskParsePlayers(AbstractAction):
    def __init__(self, data: dict):
        self._data = data

    def run(self) -> list[PlayerSchema]:
        return [
            PlayerSchema(**player, timestamp=self._data["timestamp"])
            for player in self._data["players"]
        ]


class SubTaskSavePlayersToStorage(AbstractAction):
    def __init__(self, players: list[PlayerSchema], session):
        self._players = players
        self._session = session

    def run(self):
        player_repo = PlayerRepository(self._session)
        for player in self._players:
            player_repo.create_one(**(player.dict()))
        return True


class SubTaskGetPlayerData(AbstractAction):
    def __init__(self):
        self._url = settings.API_PLAYER_URL

    def run(self):
        response = requests.get(settings.API_PLAYER_URL)
        data = response.json()
        print("CALL IS MADE")
        return data


class ActionGetAndParseAndSavePlayers(AbstractAction):
    task_get = SubTaskGetPlayerData
    task_parse = SubTaskParsePlayers
    task_save = SubTaskSavePlayersToStorage

    def __init__(self, session):
        self._session = session

    def run(self):
        player_data = self.task_get()
        players = self.task_parse(player_data)
        self.task_save(players=players, session=self._session)
        return players


class PlayerQuery(BaseModel):
    page: int = 0
    player_whitelist_tags: list[str] = []
    region_whitelist_tags: list[str] = []
    system_whitelist_tags: list[str] = []
    is_online: bool = True


class ActionGetFilteredPlayers(AbstractAction):
    def __init__(self, session, **kwargs):
        self._session = session
        self.query = PlayerQuery(**kwargs)

    def run(self):
        player_storage = PlayerRepository(self._session)
        players = player_storage.get_players_by_query(self.query)
        return players
