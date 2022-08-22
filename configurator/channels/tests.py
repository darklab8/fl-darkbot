import pytest
from . import storage
from . import schemas
from . import actions


@pytest.fixture
def test_query():
    return schemas.ChannelQueryParams(
        channel_id=5,
        owner_id=6,
        owner_name="John",
    )


@pytest.mark.asyncio
async def test_registering_channel_without_view(
    database, async_client, test_query: schemas.ChannelQueryParams
):

    await actions.ActionRegisterChannel(
        db=database,
        query=test_query,
    ).run()

    channels = await storage.ChannelStorage(db=database).get_all()

    assert len(channels) > 0
    assert channels[0].channel_id == test_query.channel_id


@pytest.mark.asyncio
async def test_registering_channel_with_view(
    database, async_client, test_query: schemas.ChannelQueryParams
):

    response = await async_client.post(
        f"/channels/{test_query.channel_id}",
        json={
            "owner_id": test_query.owner_id,
            "owner_name": test_query.owner_name,
        },
    )
    data = response.json()
    print(data)

    assert response.status_code == 200

    channels = await storage.ChannelStorage(db=database).get_all()

    assert len(channels) > 0
    assert channels[0].channel_id == test_query.channel_id
