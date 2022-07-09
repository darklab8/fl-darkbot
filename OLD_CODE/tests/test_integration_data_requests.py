from src.storage import Storage
from src.forum_parser import forum_record
import pytest
import json
import dataclasses
from types import SimpleNamespace


class EnhancedJSONEncoder(json.JSONEncoder):
    def default(self, o):
        if dataclasses.is_dataclass(o):
            return dataclasses.asdict(o)
        return super().default(o)


@pytest.fixture
def storage():
    return Storage()


@pytest.mark.integration
@pytest.mark.requests
def test_request_players(storage: Storage):
    data = storage.get_players_data()

    with open('tests/examples/players.json', 'w') as file_:
        file_.write(json.dumps(data))


@pytest.mark.integration
@pytest.mark.requests
def test_request_bases(storage: Storage):
    data = storage.get_base_data()

    with open('tests/examples/bases.json', 'w') as file_:
        file_.write(json.dumps(data))


@pytest.mark.integration
@pytest.mark.requests
def test_request_forum_records(storage: Storage):
    data = storage.get_new_forum_records({})

    with open('tests/examples/forum_records.json', 'w') as file_:
        file_.write(json.dumps(data, cls=EnhancedJSONEncoder))


@pytest.mark.integration
@pytest.mark.requests
def test_read_forum_record_back(storage: Storage):
    with open('tests/examples/forum_records.json', 'r') as file_:
        data = file_.read()

    loaded = json.loads(data)
    print(type(loaded))
    assert len(loaded) > 3

    morfed = [forum_record(**record) for record in loaded]

    assert len(morfed) > 3
    print(morfed[0])


def test_load_tested_data(storage):

    output = storage.get_load_test_game_data({})

    assert hasattr(output, "new_forum_records")
    assert hasattr(output, "players")
    assert hasattr(output, "bases")
