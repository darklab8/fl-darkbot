from . import repository
from . import schemas


def test_check_db(db):
    player_repo = repository.PlayerRepository(db)

    players = player_repo.get_all()

    assert players == []

    players = player_repo.create_one(schemas.PlayerSchema(description="abc"))

    assert len(player_repo.get_all()) == 1
