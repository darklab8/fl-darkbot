from . import repository
from . import schemas


def test_check_db(db):
    player_repo = repository.PlayerRepository()

    players = player_repo.get_all(db)

    assert players == []

    players = player_repo.create_one(db, schemas.PlayerSchema(description="abc"))

    assert len(player_repo.get_all(db)) == 1
