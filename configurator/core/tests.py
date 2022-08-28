import pytest
from utils.rest_api.message import MessageOk


def test_read_main(client):
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())
