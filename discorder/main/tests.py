import pytest
from utils.rest_api.message import MessageOk


def test_main_ping(client):
    response = client.get("/msg/ping")
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())
