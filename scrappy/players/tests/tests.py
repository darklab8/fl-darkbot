from ..storage import PlayerStorage
from faker import Faker
from .. import schemas as player_schemas

import pytest
import json
import datetime
from .. import actions as player_actions
from celery import shared_task
from ..tasks import update_players
from unittest.mock import MagicMock, patch
from fastapi.testclient import TestClient
from scrappy.core.logger import base_logger
from pathlib import Path

logger = base_logger.getChild(__name__)

fake = Faker()


def test_check_endpoint_to_get_players(database, client):
    assert client.get("/players/").json() == []


file_with_data_example = Path(__file__).parent / "data" / "players.json"


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


def test_players_check(database, mocked_request_url_data: dict):

    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get = MagicMock(return_value=mocked_request_url_data)
    action(database)

    items = PlayerStorage(database)._get_all()
    assert len(items) > 0

    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get = MagicMock(return_value=mocked_request_url_data)
    action(database)

    items2 = PlayerStorage(database)._get_all()
    assert len(items2) > 0
    assert len(items) == len(items2)


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
def test_trying_players_update_with_celery_integration(
    database, mocked_request_url_data
):

    with patch.object(
        player_actions.ActionGetAndParseAndSavePlayers,
        "task_get",
        return_value=mocked_request_url_data,
    ) as mock_method:
        task_handle = update_players.delay(database_name=database.name)
        task_handle.get()

    player_storage = PlayerStorage(database)
    players_amount = len(player_storage._get_all())

    assert players_amount > 0


@pytest.fixture
def loaded_players(database, mocked_request_url_data):
    action = player_actions.ActionGetAndParseAndSavePlayers
    action.task_get.run = MagicMock(return_value=mocked_request_url_data)
    return action(database=database)


def test_filter_players(loaded_players, database):
    assert len(player_actions.ActionGetFilteredPlayers(database)) == 19

    assert (
        len(
            player_actions.ActionGetFilteredPlayers(
                database, player_whitelist_tags=["AWES", "Aiv"]
            )
        )
        > 0
    )


def test_get_players_from_endpoint(
    database, mocked_request_url_data: dict, client: TestClient, loaded_players
):
    assert len(client.get("/players?player_tag=AWES&player_tag=Aiv").json()) == 2


@pytest.mark.asyncio
async def test_async_try(
    database, mocked_request_url_data: dict, app, loaded_players, async_client
):
    response = await async_client.get("/players-async")
    players = response.json()
    assert len(players) > 0


@pytest.mark.asyncio
async def test_asyncio_stuff():
    pass
