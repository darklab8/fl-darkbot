from sqlalchemy import func, or_
from sqlalchemy import select, update
from sqlalchemy.sql import Select, Insert
from sqlalchemy.dialects.postgresql import insert
import scrappy.players.schemas as schemas
from utils.database.sql import Database
from scrappy.players.models import Player
from .schemas import PlayerQueryParams
from scrappy.commons.storage import AbstractStorage
from sqlalchemy.engine import Row
from sqlalchemy import Column, String


class IsOnlineQuery:
    latest_timestamp = select(func.max(Player.timestamp)).scalar_subquery()

    @classmethod
    def create(cls) -> Select:
        stmt = select(
            Player, (cls.latest_timestamp == Player.timestamp).label("is_online")
        )
        return stmt

    @staticmethod
    def from_query_row_to_schema(one_row: Row) -> schemas.PlayerOut:
        return schemas.PlayerOut(**one_row[0].__dict__, is_online=one_row[1])

    @staticmethod
    def from_many_rows_to_schemas(
        many_row: list[Row],
    ) -> list[schemas.PlayerOut]:
        return [IsOnlineQuery.from_query_row_to_schema(db_row) for db_row in many_row]


class PlayerStorage(AbstractStorage):
    def __init__(self, db: Database):
        super().__init__(db=db)

    def _get_all(
        self,
    ) -> list[schemas.PlayerOut]:
        with self.db.get_core_session() as session:
            statement = IsOnlineQuery.create()
            db_rows = session.execute(statement).all()
            players = IsOnlineQuery.from_many_rows_to_schemas(db_rows)
            return players

    async def _a_get_all(
        self,
    ) -> list[schemas.PlayerOut]:
        async with self.db.get_async_session() as session:
            statement: Select = IsOnlineQuery.create()
            db_rows = (await session.execute(statement)).all()
            players = IsOnlineQuery.from_many_rows_to_schemas(db_rows)
            return players

    def create(
        self,
        *players: schemas.PlayerIn,
    ) -> None:

        with self.db.get_core_session() as session:
            stmt = insert(Player).values([item.dict() for item in players])
            stmt = stmt.on_conflict_do_update(
                index_elements=[Player.name],
                set_={k: v for k, v in stmt.excluded.items() if v.primary_key is False},
            )

            result = session.execute(stmt)
            print(result)
            session.commit()

    page_size = 20

    def get(self, query: PlayerQueryParams) -> list[schemas.PlayerOut]:

        with self.db.get_core_session() as session:
            queryset = IsOnlineQuery.create()

            if query.is_online:
                queryset = queryset.where(
                    IsOnlineQuery.latest_timestamp == Player.timestamp
                )

            def contains_any(
                queryset: Select, attribute: Column[String], tags: list[str]
            ) -> Select:
                return queryset.where(
                    or_(*[attribute.like(rf"%{tag}%") for tag in tags])
                )

            if query.player_whitelist_tags:
                queryset = contains_any(
                    queryset, Player.name, query.player_whitelist_tags
                )

            if query.region_whitelist_tags:
                queryset = queryset = contains_any(
                    queryset, Player.region, query.region_whitelist_tags
                )

            if query.system_whitelist_tags:
                queryset = queryset = contains_any(
                    queryset, Player.system, query.system_whitelist_tags
                )

            queryset = queryset.limit(self.page_size).offset(
                query.page * self.page_size
            )

            db_rows = session.execute(queryset).all()

        return IsOnlineQuery.from_many_rows_to_schemas(db_rows)
