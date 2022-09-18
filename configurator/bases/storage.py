from . import schemas
from utils.database.sql import Database
from utils.database.queryset import AbstractQuerySet, ModelIn, SchemasOut
from . import models
from sqlalchemy import delete
from sqlalchemy.dialects.postgresql import insert
from sqlalchemy.ext.asyncio import AsyncResult
from collections import defaultdict
from ..core.logger import base_logger

logger = base_logger.getChild(__name__)


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

    async def get_bases(
        self, query: schemas.BaseGetRequestParams
    ) -> schemas.BasesManyOut:
        async with self.db.get_async_session() as session:
            logger.debug("storage.get_bases: started")
            queryset = BaseQuerySet()
            stmt = queryset.select()
            result: AsyncResult = await session.execute(stmt)
            logger.debug(f"storage.get_bases.bases: executed SQL")
            db_rows = result.all()

            logger.debug("storage.get_bases.bases: got results.all")
            rows = list([row for row in db_rows])

            logger.debug(f"storage.get_bases.bases: rows={rows}")

            # group tags for channels
            base_dicts = defaultdict(
                lambda: schemas.BaseOut(
                    channel_id=-1,
                    tags=[],
                )
            )
            for row in rows:
                channel_id = row[0].channel_id
                base_dicts[channel_id].channel_id = channel_id
                base_dicts[channel_id].tags.append(row[0].tag)

            bases = [base for base in base_dicts.values()]

            logger.debug(f"storage.get_bases.bases={bases}")
            data = schemas.BasesManyOut.parse_obj(bases)
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
