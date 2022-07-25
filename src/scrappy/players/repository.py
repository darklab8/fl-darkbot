from sqlalchemy.orm import Session
from sqlalchemy import func, or_
from sqlalchemy.orm.query import Query
import scrappy.players.schemas as schemas
import scrappy.players.models as models


def filter_by_contains_in_list(queryset: Query, attribute_, list_: list[str]):
    filter_list = [attribute_.contains(x) for x in list_]
    return queryset.filter(or_(*filter_list))


class PlayerRepository:
    def __init__(self, db: Session):
        self.db: Session = db

    def get_all(
        self,
    ):
        return self.db.query(models.Player).all()

    def create_one(
        self,
        **kwargs: dict,
    ) -> schemas.PlayerSchema:
        validated_data = schemas.PlayerSchema(**kwargs)

        db_user = (
            self.db.query(models.Player)
            .filter(models.Player.name == validated_data.name)
            .first()
        )

        if db_user:
            for key, value in validated_data.fields:
                setattr(db_user, key, value)

        if db_user is None:
            db_user = models.Player(**validated_data.dict())
            self.db.add(db_user)

        self.db.commit()
        self.db.refresh(db_user)
        return schemas.PlayerSchema(**db_user.__dict__)

    page_size = 20

    def get_players_by_query(self, query: "PlayerQuery") -> list[schemas.PlayerSchema]:
        queryset = self.db.query(models.Player)

        if query.is_online:
            timedate_when_online = self.db.query(func.max(models.Player.timestamp))
            queryset = queryset.filter(models.Player.timestamp == timedate_when_online)

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

        return list([schemas.PlayerSchema(**player.__dict__) for player in players])
