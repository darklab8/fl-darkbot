import pytest
import json
from scrappy.bases import subtasks
from typing import Any
import pathlib
from .. import actions
from .. import subtasks
from unittest.mock import MagicMock
from .. import storage
from fastapi.testclient import TestClient

file_with_data_example = pathlib.Path(__file__).parent / "data" / "bases.json"


@pytest.mark.integration
def test_request_base_url() -> None:
    data = subtasks.SubTaskGetBaseData()
    from pprint import pprint as print

    print(data)
    with open(file_with_data_example, "w") as file_:
        file_.write(json.dumps(data, indent=2))


@pytest.fixture
def mocked_request_url_data() -> dict[str, Any]:
    with open(file_with_data_example, "r") as file_:
        data = file_.read()

    dict_: dict = json.loads(data)  # type: ignore
    return dict_


def test_base_action(mocked_request_url_data, database):

    mocked_request_url_data = {
        key: value
        for key, value in mocked_request_url_data.items()
        if "Aergia" in key or "Junktown" in key
    }
    action = actions.ActionGetAndParseAndSaveBases
    action.task_get = MagicMock(return_value=mocked_request_url_data)
    action(database)

    items = storage.BaseStorage(database)._get_all()
    assert len(items) > 0

    action = actions.ActionGetAndParseAndSaveBases
    action.task_get = MagicMock(return_value=mocked_request_url_data)
    action(database)

    items2 = storage.BaseStorage(database)._get_all()
    assert len(items2) > 0
    assert len(items) == len(items2)

    print(items)


@pytest.fixture
def loaded_items(database, mocked_request_url_data):
    action = actions.ActionGetAndParseAndSaveBases
    action.task_get = MagicMock(return_value=mocked_request_url_data)
    return action(database=database)


def test_get_bases_from_action(
    database, mocked_request_url_data: dict, client: TestClient, loaded_items
):
    page_size = 10
    bases = actions.ActionGetFilteredBases(
        database=database,
        page=0,
        page_size=page_size,
        name_tags=[],
    )
    assert len(bases) == page_size


def test_get_players_from_endpoint(
    database, mocked_request_url_data: dict, client: TestClient, loaded_items
):
    size = 10

    result = client.get(f"/bases/?page_size={size}").json()
    assert len(result) == size
