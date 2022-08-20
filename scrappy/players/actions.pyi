from typing import Callable, Any
from typing import Protocol
from utils.database.sql import Database

ActionGetFilteredPlayers: Callable[[Database, dict[str, Any]], list[Any]]

class ActionGetAndParseAndSavePlayersProtocol(Protocol):
    def __call__(self, database: Database) -> list[Any]: ...

ActionGetAndParseAndSavePlayers: ActionGetAndParseAndSavePlayersProtocol
