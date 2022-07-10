from unicodedata import decimal
import src.scrappy.databases as databases
import src.scrappy.players.crud as crud
import src.scrappy.players.schemas as schemas
import src.scrappy.players.models as models
import pytest


@pytest.fixture
def db():
    database = databases.Database(url="sqlite:///./test_sql_app.db")

    models.databases.default.Base.metadata.drop_all(bind=database.engine)
    models.databases.default.Base.metadata.create_all(bind=database.engine)

    with database.manager_to_get_db() as db:
        yield db


def test_check_db(db):
    player_repo = crud.PlayerRepository()

    players = player_repo.get_all(db)

    assert players == []

    players = player_repo.create_one(db, schemas.PlayerSchema(description="abc"))

    assert len(player_repo.get_all(db)) == 1
