from sqlalchemy import func, or_
from sqlalchemy.orm.query import Query
from sqlalchemy import select, insert, update
import scrappy.players.schemas as schemas
import scrappy.players.models as models
from scrappy.core.databases import Database
from scrappy.players.models import Player
from contextlib import contextmanager


def filter_by_contains_in_list(queryset: Query, attribute_, list_: list[str]):
    filter_list = [attribute_.contains(x) for x in list_]
    return queryset.filter(or_(*filter_list))


class IsOnlineQuery:
    latest_timestamp = select(func.max(Player.timestamp)).scalar_subquery()

    def __new__(cls):
        stmt = select(
            Player, (cls.latest_timestamp == Player.timestamp).label("is_online")
        )

        return stmt

    @staticmethod
    def from_query_row_to_schema(one_row):
        return schemas.PlayerOut(**one_row[0].__dict__, is_online=one_row[1])

    def from_many_rows_to_schemas(many_row):
        return [IsOnlineQuery.from_query_row_to_schema(db_row) for db_row in many_row]


class PlayerStorage:
    def __init__(self, db: Database):
        self.db: Database = db

    def get_all(
        self,
    ):
        with self.db.get_core_session() as session:
            statement = IsOnlineQuery()
            db_rows = session.execute(statement).all()
            players = IsOnlineQuery.from_many_rows_to_schemas(db_rows)
            return players

    async def a_get_all(
        self,
    ):
        async with self.db.get_async_session() as session:
            statement = IsOnlineQuery()
            db_rows = (await session.execute(statement)).all()
            players = IsOnlineQuery.from_many_rows_to_schemas(db_rows)
            return players

    def create_one(
        self,
        **kwargs: dict,
    ) -> schemas.PlayerOut:
        validated_data = schemas.PlayerIn(**kwargs)

        with self.db.get_core_session() as session:

            get_player = select(Player).where(Player.name == validated_data.name)
            already_present_user = session.execute(get_player).scalar()

            if already_present_user:
                add_or_update_user_query = (
                    update(Player)
                    .where(Player.id == already_present_user.id)
                    .values(**validated_data.dict())
                )
            else:
                add_or_update_user_query = insert(Player).values(
                    **validated_data.dict()
                )

            session.execute(add_or_update_user_query)

            get_refreshed_player = IsOnlineQuery().where(
                IsOnlineQuery.latest_timestamp == Player.timestamp
            )

            db_row = session.execute(get_refreshed_player).first()

            extracted_info = IsOnlineQuery.from_query_row_to_schema(db_row)

            session.commit()

        return extracted_info

    page_size = 20

    def get_players_by_query(self, query: "PlayerQuery") -> list[schemas.PlayerOut]:

        with self.db.get_core_session() as session:
            queryset = IsOnlineQuery()

            if query.is_online:
                queryset = queryset.where(
                    IsOnlineQuery.latest_timestamp == Player.timestamp
                )

            def contains_any(queryset, attribute, tags):
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
