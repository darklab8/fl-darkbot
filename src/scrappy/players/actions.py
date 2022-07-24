from utils.porto import AbstractAction
from .repository import PlayerRepository
from .schemas import PlayerSchema


class ActionParsePlayers(AbstractAction):
    def __init__(self, data: dict):
        self._data = data

    def run(self) -> list[PlayerSchema]:
        return [
            PlayerSchema(**player, timestamp=self._data["timestamp"])
            for player in self._data["players"]
        ]


class ActionSavePlayersToStorage(AbstractAction):
    def __init__(self, players: list[PlayerSchema], db):
        self._players = players
        self._db = db

    def run(self):
        player_repo = PlayerRepository(self._db)
        for player in self._players:
            player_repo.create_one(**(player.dict()))
        return True
