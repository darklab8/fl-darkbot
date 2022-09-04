import pytest
from . import queries
from .urls import urls
from utils.rest_api.message import MessageOk
import asyncio
from httpx import AsyncClient


class TestMessage:
    def __init__(self, id, content):
        self.id = id
        self.content = content

    def __str__(self):
        return f"{self.id}\n```{self.content}```"


@pytest.mark.asyncio
async def test_create_or_replace_msg(client: AsyncClient, channel_id: int):

    id = "2ca613b64fdc2eb7"
    await asyncio.sleep(10)
    response = client.post(
        urls.base,
        json=dict(
            queries.CreateOrReplaceMessqgeQueryParams(
                id=id,
                channel_id=channel_id,
                message=str(
                    TestMessage(id=id, content="create or replace testing message")
                ),
            )
        ),
    )
    print(response.json())
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())

    response = client.post(
        urls.base,
        json=dict(
            queries.CreateOrReplaceMessqgeQueryParams(
                id=id,
                channel_id=channel_id,
                message=str(
                    TestMessage(id=id, content="create or replace testing message 2")
                ),
            )
        ),
    )
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())

    response = client.delete(
        urls.base,
        json=dict(
            queries.DeleteMessageQueryParams(
                id=id,
                channel_id=channel_id,
            )
        ),
    )
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())
