from .repository import PlayerRepository
from .schemas import PlayerSchema

def test_check_db(db):
    player_repo = PlayerRepository(db)

    players = player_repo.get_all()

    assert players == []

    player = player_repo.create_one(description="abc")

    assert player.id == 1
    assert player.description == "abc"
