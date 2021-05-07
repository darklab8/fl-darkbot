import pytest
from commands import attach_commands
from storage import Storage
from looper import Looper


@pytest.fixture
def storage():
    return Storage()


@pytest.fixture
def bot(storage):
    return attach_commands(storage)


def looper(bot):
    _ = Looper(bot, storage)


def test_request_game_data(storage):
    game_data = storage.get_game_data()

    assert len(game_data.players) > 0
    assert len(game_data.bases) > 0
