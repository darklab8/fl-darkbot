from ..core.logger import base_logger
from configurator.commons.actions import AbstractAction
from . import schemas
from . import storage
from utils.database.sql import Database

logger = base_logger.getChild(__name__)


class ActionRegisterBase:
    def __init__(self, db: Database, query: schemas.BaseRegisterRequestParams):
        self.db = db
        self.query = query

    async def run(self) -> None:
        await storage.BaseStorage(self.db).register_base(self.query)


class ActionDeleteBases:
    def __init__(self, db: Database, channel_id: int):
        self.db = db
        self.channel_id = channel_id

    async def run(self) -> None:
        await storage.BaseStorage(self.db).delete_bases(self.channel_id)
