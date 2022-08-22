from ..core.logger import base_logger
from utils.porto import AbstractAction
from . import schemas
from . import storage
from utils.database.sql import Database

logger = base_logger.getChild(__name__)


class ActionRegisterChannel:
    def __init__(self, db: Database, query: schemas.ChannelQueryParams):
        self.db = db
        self.query = query

    async def run(self) -> None:
        await storage.ChannelStorage(self.db).register(self.query)


class ActionDeleteChannel:
    def __init__(self, db: Database, channel_id: int):
        self.db = db
        self.channel_id = channel_id

    async def run(self) -> None:
        await storage.ChannelStorage(self.db).unregister(self.channel_id)
