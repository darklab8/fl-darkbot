import databases as databases
from . import repository
from . import schemas
import pytest


@pytest.fixture
def db():
    database = databases.Database(
        # url="sqlite:///./test_sql_app.db"
        url="postgresql://postgres:postgres@localhost/default"
    )

    databases.default.Base.metadata.drop_all(bind=database.engine)
    databases.default.Base.metadata.create_all(bind=database.engine)

    with database.manager_to_get_db() as db:
        yield db


def test_check_db(db):
    player_repo = repository.PlayerRepository()

    players = player_repo.get_all(db)

    assert players == []

    players = player_repo.create_one(db, schemas.PlayerSchema(description="abc"))

    assert len(player_repo.get_all(db)) == 1
