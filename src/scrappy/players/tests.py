from utils.porto import AbstractAction
from .repository import PlayerRepository
from faker import Faker
from .schemas import PlayerSchema

import pytest
import json
import os
import datetime
from . import actions as player_actions
from celery import shared_task

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

    with open(file_with_data_example, "w") as file_:
        file_.write(json.dumps(data, indent=2))


@pytest.fixture
def mocked_request_url_data():
    with open(file_with_data_example, "r") as file_:
        data = file_.read()

    dict_ = json.loads(data)
    return dict_


from unittest.mock import MagicMock


def test_players_check(db, mocked_request_url_data: dict):

    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get.run = MagicMock(return_value=mocked_request_url_data)
    action(db=db)

    player_repo = PlayerRepository(db)
    assert len(player_repo.get_all()) > 0


def test_repeated_players_override_previous_players(db):
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


@shared_task
def mul(x, y):
    return x * y


@pytest.mark.usefixtures("celery_session_app")
@pytest.mark.usefixtures("celery_session_worker")
def test_try_testing_celery():
    task_handle = mul.delay(2, 3)
    assert task_handle.get() == 6
