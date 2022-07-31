from sqlalchemy.orm import Session
from sqlalchemy import func, or_
from sqlalchemy.orm.query import Query
import scrappy.players.schemas as schemas
import scrappy.players.models as models


def filter_by_contains_in_list(queryset: Query, attribute_, list_: list[str]):
    filter_list = [attribute_.contains(x) for x in list_]
    return queryset.filter(or_(*filter_list))


class IsOnlineQuery:
    def __new__(cls, db):
        statement = db.query(func.max(models.Player.timestamp)).subquery()
        return db.query(
            models.Player, (statement == models.Player.timestamp).label("is_online")
        )

    @staticmethod
    def from_query_row_to_schema(one_row):
        return schemas.PlayerOut(**one_row[0].__dict__, is_online=one_row[1])

    def from_many_rows_to_schemas(many_row):
        return [IsOnlineQuery.from_query_row_to_schema(db_row) for db_row in many_row]


class PlayerRepository:
    def __init__(self, db: Session):
        self.db: Session = db

    def get_all(
        self,
    ):
        return IsOnlineQuery.from_many_rows_to_schemas(IsOnlineQuery(self.db).all())

    def create_one(
        self,
        **kwargs: dict,
    ) -> schemas.PlayerIn:
        validated_data = schemas.PlayerIn(**kwargs)

        db_user = (
            self.db.query(models.Player)
            .filter(models.Player.name == validated_data.name)
            .first()
        )

        if db_user:
            for key, value in validated_data:
                setattr(db_user, key, value)

        if db_user is None:
            db_user = models.Player(**validated_data.dict())
            self.db.add(db_user)

        self.db.commit()
        self.db.refresh(db_user)

        db_user = (
            IsOnlineQuery(self.db)
            .filter(models.Player.name == validated_data.name)
            .first()
        )

        return IsOnlineQuery.from_query_row_to_schema(db_user)

    page_size = 20

    def get_players_by_query(self, query: "PlayerQuery") -> list[schemas.PlayerOut]:

        queryset = IsOnlineQuery(self.db)

        if query.is_online:
            maximum = self.db.query(func.max(models.Player.timestamp)).subquery()
            queryset = queryset.filter(models.Player.timestamp == maximum)

        if query.player_whitelist_tags:
            queryset = filter_by_contains_in_list(
                queryset, models.Player.name, query.player_whitelist_tags
            )

        if query.region_whitelist_tags:
            queryset = filter_by_contains_in_list(
                queryset, models.Player.region, query.region_whitelist_tags
            )

        if query.system_whitelist_tags:
            queryset = filter_by_contains_in_list(
                queryset, models.Player.system, query.system_whitelist_tags
            )

        queryset = queryset.limit(self.page_size).offset(query.page * self.page_size)
        players = queryset.all()

        return IsOnlineQuery.from_many_rows_to_schemas(players)
