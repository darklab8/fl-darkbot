from .repository import PlayerRepository
from faker import Faker
from .schemas import PlayerSchema
import requests
import scrappy.core.settings as settings
import pytest
import json

fake = Faker()


class PlayerTestFactory:
    repo_model = PlayerRepository

    def __new__(cls, db, **kwargs: dict) -> PlayerSchema:
        repo = cls.repo_model(db)
        return repo.create_one(description=kwargs.get("description", fake.name()))


def test_check_db(db):
    player_repo = PlayerRepository(db)

    players = player_repo.get_all()

    assert players == []

    player = player_repo.create_one(description="abc")

    assert player.id == 1
    assert player.description == "abc"


def test_check_test_factory(db):

    player = PlayerTestFactory(db)
    assert player.id == 1
    assert isinstance(player.description, str)


def test_check_endpoint_to_get_players(db, client):
    assert client.get("/players/").json() == []

@pytest.mark.integration
def test_get_player_data():
    response = requests.get(settings.API_PLAYER_URL)
    data = response.json()
    with open("scrappy/players/test_example/players.json", "w") as file_:
        file_.write(json.dumps(data, indent=2))