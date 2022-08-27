from . import schemas
from utils.database.sql import Database
from utils.database.queryset import AbstractQuerySet, ModelIn, SchemasOut
from . import models
from sqlalchemy import delete
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.ext.asyncio import AsyncResult


class BaseQuerySet(AbstractQuerySet):
    def __init__(self):
        super().__init__(
            model=models.Base,
            out=schemas.BaseOut,
        )


class BaseStorage:
    def __init__(self, db: Database):
        self.db = db

    async def get_base(self, channel_id: int) -> schemas.BaseOut:
        async with self.db.get_async_session() as session:
            queryset = BaseQuerySet()
            stmt = queryset.select().where(models.Base.channel_id == channel_id)
            result: AsyncResult = await session.execute(stmt)
            db_rows = result.all()

            data = schemas.BaseOut(
                channel_id=channel_id,
                tags=[row[0].tag for row in db_rows],
            )
            return data

    async def register_base(self, query: schemas.BaseRegisterRequestParams):
        async with self.db.get_async_session() as session:

            input_data = list(
                [
                    {"channel_id": query.channel_id, "tag": tag}
                    for tag in query.base_tags
                ]
            )
            stmt = insert(models.Base).values(input_data)
            stmt = stmt.on_conflict_do_nothing(
                constraint=models.BaseConstraints.channel_tag
            )

            await session.execute(stmt)
            await session.commit()

    async def delete_bases(self, channel_id: int):
        async with self.db.get_async_session() as session:
            stmt = delete(models.Base).where(models.Base.channel_id == channel_id)
            await session.execute(stmt)
            await session.commit()
