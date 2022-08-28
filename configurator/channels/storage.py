from . import schemas
from utils.database.sql import Database
from utils.database.queryset import AbstractQuerySet, ModelIn, SchemasOut
from . import models
from sqlalchemy import delete
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.ext.asyncio import AsyncResult


class ChanellQuerySet(AbstractQuerySet):
    def __init__(self):
        super().__init__(
            model=models.Channel,
            out=schemas.ChannelOut,
        )

    def from_query_row_to_schema(self, one_row: tuple[ModelIn]) -> SchemasOut:
        return self.schema_out(channel_id=one_row[0].id)


class ChannelOwnerQuerySet(AbstractQuerySet):
    def __init__(self):
        super().__init__(
            model=models.ChannelOwner,
            out=schemas.ChannelOwnerOut,
        )


class ChannelStorage:
    def __init__(self, db: Database):
        self.db = db

    async def get_all(self):
        async with self.db.get_async_session() as session:
            queryset = ChanellQuerySet()
            stmt = queryset.select()
            result: AsyncResult = await session.execute(stmt)
            db_rows = result.all()
            return queryset.from_many_rows_to_schemas(db_rows)

    async def get_owner_by_channel_id(self, channel_id: int):
        async with self.db.get_async_session() as session:
            queryset = ChannelOwnerQuerySet()
            stmt = queryset.select().where(models.ChannelOwner.channel_id == channel_id)
            result: AsyncResult = await session.execute(stmt)
            db_row = result.first()
            return queryset.from_query_row_to_schema(db_row)

    async def register(self, query: schemas.ChannelCreateQueryParams):
        async with self.db.get_async_session() as session:

            stmt = insert(models.Channel).values(id=query.channel_id)
            stmt = stmt.on_conflict_do_nothing(index_elements=[models.Channel.id])

            await session.execute(stmt)
            await session.commit()

        await self.record_owner(query=query)

    async def record_owner(self, query: schemas.ChannelCreateQueryParams):
        if query.owner_id is None:
            return
        async with self.db.get_async_session() as session:
            stmt = insert(models.ChannelOwner).values(
                id=query.owner_id,
                name=query.owner_name,
                channel_id=query.channel_id,
            )
            stmt = stmt.on_conflict_do_update(
                constraint=models.ChannelOwnerConstraints.owner_channel,
                set_={k: v for k, v in stmt.excluded.items() if v.primary_key is False},
            )

            await session.execute(stmt)
            await session.commit()

    async def unregister(self, channel_id: int):
        async with self.db.get_async_session() as session:
            stmt = delete(models.Channel).where(models.Channel.id == channel_id)
            await session.execute(stmt)
            await session.commit()
