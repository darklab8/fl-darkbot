import pytest
from utils.rest_api.message import MessageOk


def test_main_ping(client):
    response = client.get("/ping")
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())


def test_get_guilds(client):
    # wait, i can eliminate deleted servers? :thinking:
    response = client.get("/guilds")
    assert response.status_code == 200
    print(response.json())
