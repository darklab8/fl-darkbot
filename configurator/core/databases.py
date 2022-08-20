from utils.database.sql import DatabaseFactoryBase

import scrappy.core.settings as settings


class DatabaseFactory(DatabaseFactoryBase):
    settings = settings
