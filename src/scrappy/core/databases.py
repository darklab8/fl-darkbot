from utils.database.sql import Database, DatabaseFactoryBase

import scrappy.core.settings as settings

class DatabaseFactory(DatabaseFactoryBase):
    settings = settings
