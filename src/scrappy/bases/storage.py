from scrappy.core.databases import Database
from .schemas import BaseIn, BaseOut, BaseQueryParams


class BaseStorage:
    def __init__(self, db: Database):
        self.db: Database = db

    def create(
        self,
        *bases: list[BaseIn],
    ) -> list[BaseOut]:
        pass

    def get(self, query: BaseQueryParams) -> list[BaseOut]:
        pass
