from sqlalchemy.orm import Session

import scrappy.players.schemas as schemas
import scrappy.players.models as models


class PlayerRepository:
    def __init__(self,db: Session):
        self.db: Session = db

    def get_all(
        self,
    ):
        return self.db.query(models.Player).all()

    def create_one(
        self,
        player: schemas.PlayerSchema,
    ):
        db_user = models.Player(description=player.description)
        self.db.add(db_user)
        self.db.commit()
        self.db.refresh(db_user)
        return db_user
