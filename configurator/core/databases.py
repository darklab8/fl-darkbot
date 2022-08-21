from utils.database.sql import DatabaseFactoryBase

from . import settings


class DatabaseFactory(DatabaseFactoryBase):
    settings = settings
