from utils.rest_api.methods import RequestMethod
from utils.database.sql import Database
from utils.rest_api.message import MessageOk
from ..core.logger import base_logger
from utils.porto import AsyncAbstractAction
from . import schemas
from . import storage
from .urls import urls


logger = base_logger.getChild(__name__)


class ActionGetBases(AsyncAbstractAction):
    url = urls.base_get
    method = RequestMethod.post
    query_factory = schemas.BaseGetRequestParams
    response_factory = schemas.BasesManyOut

    def __init__(self, db: Database, query: schemas.BaseGetRequestParams):
        self.db = db
        self.query = query

    async def run(self) -> None:
        logger.debug(f"actions.ActionGetBases.query={self.query}")
        bases = await storage.BaseStorage(self.db).get_bases(self.query)
        logger.debug(f"actions.ActionGetBases.bases={bases}")
        return bases


class ActionRegisterBase(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.post
    query_factory = schemas.BaseRegisterRequestParams
    response_factory = MessageOk

    def __init__(self, db: Database, query: schemas.BaseRegisterRequestParams):
        self.db = db
        self.query = query

    async def run(self) -> None:
        await storage.BaseStorage(self.db).register_base(self.query)
        return MessageOk()


class ActionDeleteBases(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.delete
    query_factory = schemas.BaseDeleteRequestParams
    response_factory = MessageOk

    def __init__(self, db: Database, query: schemas.BaseDeleteRequestParams):
        self.db = db
        self.channel_id = query.channel_id

    async def run(self) -> None:
        await storage.BaseStorage(self.db).delete_bases(self.channel_id)
        return MessageOk()