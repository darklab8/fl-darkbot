import pytest
from . import queries
from .urls import urls
from utils.rest_api.message import MessageOk
import asyncio
from httpx import AsyncClient
from fastapi.testclient import TestClient
from . import rpc


class DummyMessage:
    def __init__(self, id, content):
        self.id = id
        self.content = content

    def __str__(self):
        return f"{self.id}\n```{self.content}```"


@pytest.mark.asyncio
async def test_create_or_replace_msg(client: TestClient, channel_id: int):

    id = "2ca613b64fdc2eb7"
    await asyncio.sleep(3)
    response = client.post(
        urls.base,
        json=dict(
            queries.CreateOrReplaceMessqgeQueryParams(
                id=id,
                channel_id=channel_id,
                message=str(
                    DummyMessage(id=id, content="create or replace testing message")
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
                    DummyMessage(id=id, content="create or replace testing message 2")
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


@pytest.mark.integration
@pytest.mark.asyncio
async def test_rpc(client: AsyncClient, channel_id: int):

    id = "43ru34hri3hr34"
    assert (
        MessageOk()
        == await rpc.CreateOrReplaceMessage(
            query=rpc.CreateOrReplaceMessage.query_factory(
                id=id,
                channel_id=channel_id,
                message=str(
                    DummyMessage(id=id, content="create or replace testing message")
                ),
            )
        ).run()
    )

    assert (
        MessageOk()
        == await rpc.CreateOrReplaceMessage(
            query=rpc.CreateOrReplaceMessage.query_factory(
                id=id,
                channel_id=channel_id,
                message=str(
                    DummyMessage(id=id, content="create or replace testing message 2")
                ),
            )
        ).run()
    )

    assert (
        MessageOk()
        == await rpc.DeleteMessage(
            query=rpc.DeleteMessage.query_factory(
                id=id,
                channel_id=channel_id,
            )
        ).run()
    )
