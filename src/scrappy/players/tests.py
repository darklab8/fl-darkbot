from unicodedata import decimal
import src.scrappy.database as database
import src.scrappy.players.crud as crud
import src.scrappy.players.schemas as schemas
import src.scrappy.players.models as models
import pytest


@pytest.fixture
def db():
    with database.manager_to_get_db() as db:
        yield db


def test_check_db(db):
    player_repo = crud.PlayerRepository()

    players = player_repo.get_all(db)

    assert players == []

    players = player_repo.create_one(db, schemas.PlayerSchema(description="abc"))

    assert player_repo.get_all(db).count(models.Player.id) == 1
