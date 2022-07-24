from .repository import PlayerRepository
from faker import Faker
from .schemas import PlayerSchema
import requests
import scrappy.core.settings as settings
import pytest
import json
import os
import datetime
from . import actions as player_actions

fake = Faker()


class PlayerTestFactory:
    repo_model = PlayerRepository

    def __new__(cls, db, **kwargs: dict) -> PlayerSchema:
        repo = cls.repo_model(db)
        return repo.create_one(
            name=kwargs.get("name", fake.name()),
            region=kwargs.get("region", fake.name()),
            system=kwargs.get("system", fake.name()),
            time=kwargs.get("time", fake.name()),
            timestamp=kwargs.get("timestamp", datetime.datetime.utcnow()),
        )


def test_check_test_factory(db):

    player = PlayerTestFactory(db)
    assert player.id == 1
    assert isinstance(player.name, str)


def test_check_endpoint_to_get_players(db, client):
    assert client.get("/players/").json() == []


file_with_data_example = os.path.join(
    os.path.dirname(__file__), "test_example", "players.json"
)


@pytest.mark.integration
def test_get_player_data():
    response = requests.get(settings.API_PLAYER_URL)
    data = response.json()
    with open(file_with_data_example, "w") as file_:
        file_.write(json.dumps(data, indent=2))


@pytest.fixture
def mocked_request_url_data():
    with open(file_with_data_example, "r") as file_:
        data = file_.read()

    dict_ = json.loads(data)
    return dict_


def test_players_check(db, mocked_request_url_data: dict):

    players = player_actions.ActionParsePlayers(mocked_request_url_data)
    player_actions.ActionSavePlayersToStorage(players=players, db=db)

    player_repo = PlayerRepository(db)
    assert len(player_repo.get_all()) > 0


def test_repeated_players_override_previous_players(db, mocked_request_url_data: dict):
    fixed_player_name = "Alpha"
    player = PlayerTestFactory(db, name=fixed_player_name)

    player_repo = PlayerRepository(db)
    players_amount = len(player_repo.get_all())

    player = PlayerTestFactory(db, name=fixed_player_name)

    players_amount2 = len(player_repo.get_all())
    player_in_db = player_repo.get_all()[0]

    assert players_amount > 0
    assert players_amount == players_amount2
    assert player.name == player_in_db.name
    assert player.region == player_in_db.region
