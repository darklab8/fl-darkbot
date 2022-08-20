from typing import Callable, Any
from typing import Protocol
from utils.database.sql import Database

class ActionGetAndParseAndSaveBasesProtocol(Protocol):
    def __call__(self, database: Database) -> list[Any]: ...

ActionGetAndParseAndSaveBases: ActionGetAndParseAndSaveBasesProtocol

class ActionGetAndParseAndSavePlayersProtocol(Protocol):
    def __call__(
        self, database: Database, page: int, name_tags: list[str]
    ) -> list[Any]: ...

ActionGetFilteredBases: ActionGetAndParseAndSavePlayersProtocol
