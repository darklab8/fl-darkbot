from .repository import PlayerRepository
from faker import Faker
from . import schemas as player_schemas

import pytest
import json
import os
import datetime
from . import actions as player_actions
from celery import shared_task
from .tasks import update_players
from unittest.mock import MagicMock, patch
from fastapi.testclient import TestClient

fake = Faker()


class PlayerTestFactory:
    repo_model = PlayerRepository

    def __new__(cls, session, **kwargs: dict) -> player_schemas.PlayerOut:
        repo = cls.repo_model(session)
        return repo.create_one(
            name=kwargs.get("name", fake.name()),
            region=kwargs.get("region", fake.name()),
            system=kwargs.get("system", fake.name()),
            time=kwargs.get("time", fake.name()),
            timestamp=kwargs.get("timestamp", datetime.datetime.utcnow()),
        )


def test_check_test_factory(session):

    player = PlayerTestFactory(session)
    assert player.id == 1
    assert isinstance(player.name, str)
    assert isinstance(player.is_online, bool)
    print(player)


def test_check_endpoint_to_get_players(session, client):
    assert client.get("/players/").json() == []


file_with_data_example = os.path.join(
    os.path.dirname(__file__), "test_example", "players.json"
)


@pytest.mark.integration
def test_get_player_data():
    data = player_actions.SubTaskGetPlayerData()
    with open(file_with_data_example, "w") as file_:
        file_.write(json.dumps(data, indent=2))


@pytest.fixture
def mocked_request_url_data():
    with open(file_with_data_example, "r") as file_:
        data = file_.read()

    dict_ = json.loads(data)
    return dict_


def test_players_check(session, mocked_request_url_data: dict):

    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get.run = MagicMock(return_value=mocked_request_url_data)
    action(session=session)

    player_repo = PlayerRepository(session)
    players = player_repo.get_all()
    assert len(players) > 0
    print(players)


def test_repeated_players_override_previous_players(session):
    fixed_player_name = "Alpha"
    player = PlayerTestFactory(session, name=fixed_player_name)

    player_repo = PlayerRepository(session)
    players_amount = len(player_repo.get_all())

    player = PlayerTestFactory(session, name=fixed_player_name)

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


@pytest.mark.usefixtures("celery_session_app")
@pytest.mark.usefixtures("celery_session_worker")
def test_trying_players_update(session, mocked_request_url_data):

    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get = lambda self: mocked_request_url_data

    action(session)

    player_repo = PlayerRepository(session)
    players_amount = len(player_repo.get_all())

    assert players_amount > 0


@pytest.mark.usefixtures("celery_session_app")
@pytest.mark.usefixtures("celery_session_worker")
def test_trying_players_update_with_celery_integration(
    session, database, mocked_request_url_data
):

    with patch.object(
        player_actions.ActionGetAndParseAndSavePlayers,
        "task_get",
        return_value=mocked_request_url_data,
    ) as mock_method:
        task_handle = update_players.delay(database_name=database.name)
        task_handle.get()

    player_repo = PlayerRepository(session)
    players_amount = len(player_repo.get_all())

    assert players_amount > 0


@pytest.fixture
def loaded_players(session, mocked_request_url_data):
    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get.run = MagicMock(return_value=mocked_request_url_data)
    return action(session=session)


def test_filter_players(loaded_players, session):
    assert len(player_actions.ActionGetFilteredPlayers(session)) == 19

    assert (
        len(
            player_actions.ActionGetFilteredPlayers(
                session, player_whitelist_tags=["AWES", "Aiv"]
            )
        )
        > 0
    )


def test_get_players_from_endpoint(
    session, mocked_request_url_data: dict, client: TestClient, loaded_players
):
    assert len(client.get("/players/?player_tag=AWES&player_tag=Aiv").json()) == 2
