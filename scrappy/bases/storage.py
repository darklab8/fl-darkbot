from utils.database.sql import Database
from .schemas import BaseIn, BaseOut, BaseQueryParams
from scrappy.commons.storage import AbstractStorage


class BaseStorage(AbstractStorage):
    def __init__(self, db: Database):
        super().__init__(db=db)

    def create(
        self,
        *bases: list[BaseIn],
    ) -> list[BaseOut]:
        pass

    def get(self, query: BaseQueryParams) -> list[BaseOut]:
        pass
