import pytest
from . import queries
from .urls import urls
from utils.rest_api.message import MessageOk
import asyncio


@pytest.mark.asyncio
async def test_create_or_replace_msg(client, channel_id: int):
    await asyncio.sleep(10)
    response = client.post(
        urls.base,
        json=dict(
            queries.CreateOrReplaceMessqgeQueryParams(
                id="123",
                channel_id=channel_id,
                message="create or replace testing message",
            )
        ),
    )
    print(response.json())
    assert response.status_code == 200
    assert response.json() == dict(MessageOk())
