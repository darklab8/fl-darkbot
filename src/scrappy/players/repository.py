from sqlalchemy.orm import Session

import scrappy.players.schemas as schemas
import scrappy.players.models as models


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
