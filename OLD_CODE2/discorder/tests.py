import pytest
from utils.rest_api.message import MessageOk


def test_root_ping(client):
    response = client.get("/")
    assert response.status_code == 200
    # assert response.json() == dict(MessageOk())
    assert "bot.Bot" in response.json()
