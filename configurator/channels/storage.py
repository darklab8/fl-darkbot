from . import schemas
from utils.database.sql import Database
from . import models
from sqlalchemy import delete, select
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.sql import Select
from sqlalchemy.engine import Row


class ChannelQuerySet:
    @classmethod
    def select(cls) -> Select:
        stmt = select(models.Channel)
        return stmt

    @staticmethod
    def from_query_row_to_schema(one_row: tuple[models.Channel]) -> schemas.ChannelOut:
        return schemas.ChannelOut(channel_id=one_row[0].id)

    @classmethod
    def from_many_rows_to_schemas(
        cls,
        many_row: list[Row],
    ) -> list[schemas.ChannelOut]:
        return [cls.from_query_row_to_schema(db_row) for db_row in many_row]


class ChannelStorage:
    def __init__(self, db: Database):
        self.db = db

    async def get_all(self):
        async with self.db.get_async_session() as session:
            stmt = ChannelQuerySet.select()
            db_rows = await session.execute(stmt)
            return ChannelQuerySet.from_many_rows_to_schemas(db_rows)

    async def register(self, query: schemas.ChannelQueryParams):
        async with self.db.get_async_session() as session:

            stmt = insert(models.Channel).values(id=query.channel_id)
            stmt = stmt.on_conflict_do_nothing(index_elements=[models.Channel.id])

            await session.execute(stmt)
            await session.commit()

        await self.record_owner(query=query)

    async def record_owner(self, query: schemas.ChannelQueryParams):
        if query.owner_id is None:
            return
        async with self.db.get_async_session() as session:
            stmt = insert(models.ChannelOwner).values(
                owner_id=query.owner_id,
                channel=query.channel_id,
                owner_name=query.owner_name,
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
