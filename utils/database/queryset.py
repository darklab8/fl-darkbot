from typing import TypeVar
from sqlalchemy import select
from sqlalchemy.engine import Row
from sqlalchemy.sql import Select

ModelIn = TypeVar("ModelIn")
SchemasOut = TypeVar("SchemasOut")


class AbstractQuerySet:
    def __init__(self, model: ModelIn, out: SchemasOut):
        self.model = model
        self.schema_out = out

    def select(self) -> "Select":
        self.stmt = select(self.model)
        return self.stmt

    def from_query_row_to_schema(self, one_row: tuple[ModelIn]) -> SchemasOut:
        return self.schema_out(**one_row[0].__dict__)

    def from_many_rows_to_schemas(
        self,
        many_row: list[Row],
    ) -> list[SchemasOut]:
        return [self.from_query_row_to_schema(db_row) for db_row in many_row]
