from utils.database.sql import DatabaseFactoryBase, SettingStub

import scrappy.core.settings as settings


class DatabaseFactory(DatabaseFactoryBase):
    settings = settings  # type: ignore
