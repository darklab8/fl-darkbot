from utils.rest_api.methods import RequestMethod
from utils.database.sql import Database

from ..core.logger import base_logger
from utils.porto import AsyncAbstractAction
from . import schemas
from . import storage
from .urls import urls


logger = base_logger.getChild(__name__)


class ActionRegisterBase(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.post
    query_factory = schemas.BaseRegisterRequestParams

    def __init__(self, db: Database, query: schemas.BaseRegisterRequestParams):
        self.db = db
        self.query = query

    async def run(self) -> None:
        await storage.BaseStorage(self.db).register_base(self.query)


class ActionDeleteBases(AsyncAbstractAction):
    url = urls.base
    method = RequestMethod.delete
    query_factory = schemas.BaseDeleteRequestParams

    def __init__(self, db: Database, query: schemas.BaseDeleteRequestParams):
        self.db = db
        self.channel_id = query.channel_id

    async def run(self) -> None:
        await storage.BaseStorage(self.db).delete_bases(self.channel_id)
