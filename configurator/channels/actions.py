from utils.database.sql import Database
from utils.rest_api.methods import RequestMethod
from utils.porto import AsyncAbstractAction
from utils.rest_api.message import MessageOk
from . import schemas
from . import storage
from .urls import urls
from ..core.logger import base_logger


logger = base_logger.getChild(__name__)


class ActionRegisterChannel(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.post
    query_factory = schemas.ChannelCreateQueryParams
    response_factory = MessageOk

    def __init__(self, db: Database, query: schemas.ChannelCreateQueryParams):
        self.db = db
        self.query = query

    async def run(self) -> None:
        await storage.ChannelStorage(self.db).register(self.query)
        return MessageOk()


class ActionDeleteChannel(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.delete
    query_factory = schemas.ChannelDeleteQueryParams
    response_factory = MessageOk

    def __init__(self, db: Database, query: schemas.ChannelDeleteQueryParams):
        self.db = db
        self.channel_id = query.channel_id

    async def run(self) -> None:
        await storage.ChannelStorage(self.db).unregister(self.channel_id)
        return MessageOk()
