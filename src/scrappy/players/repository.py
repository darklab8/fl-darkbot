from sqlalchemy.orm import Session

from . import models, schemas


class PlayerRepository:
    def get_all(
        self,
        db: Session,
    ):
        return db.query(models.Player).all()

    def create_one(
        self,
        db: Session,
        player: schemas.PlayerSchema,
    ):
        db_user = models.Player(description=player.description)
        db.add(db_user)
        db.commit()
        db.refresh(db_user)
        return db_user
