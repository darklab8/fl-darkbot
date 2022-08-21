from utils.database.sql import Database
from .schemas import BaseIn, BaseOut, BaseQueryParams
from scrappy.commons.storage import AbstractStorage
from . import schemas
from sqlalchemy import select
from .models import Base
from sqlalchemy.sql import Select
from sqlalchemy.engine import Row
from sqlalchemy.dialects.postgresql import insert


class BaseQuerySet:
    @classmethod
    def create(cls) -> Select:
        stmt = select(Base)
        return stmt

    @staticmethod
    def from_query_row_to_schema(one_row: Row) -> schemas.BaseOut:
        return schemas.BaseOut(**one_row[0].__dict__)

    @staticmethod
    def from_many_rows_to_schemas(
        many_row: list[Row],
    ) -> list[schemas.BaseOut]:
        return [BaseQuerySet.from_query_row_to_schema(db_row) for db_row in many_row]


class BaseStorage(AbstractStorage):
    def __init__(self, db: Database):
        super().__init__(db=db)

    def _get_all(
        self,
    ) -> list[schemas.BaseOut]:
        with self.db.get_core_session() as session:
            statement = BaseQuerySet.create()
            db_rows = session.execute(statement).all()
            players = BaseQuerySet.from_many_rows_to_schemas(db_rows)
            return players

    def create(
        self,
        *bases: list[BaseIn],
    ) -> None:
        with self.db.get_core_session() as session:
            stmt = insert(Base).values([base.dict() for base in bases])
            stmt = stmt.on_conflict_do_update(
                index_elements=[Base.name],
                set_={k: v for k, v in stmt.excluded.items() if v.primary_key is False},
            )

            result = session.execute(stmt)
            print(result)
            session.commit()

    def get(self, query: BaseQueryParams) -> list[BaseOut]:
        pass
