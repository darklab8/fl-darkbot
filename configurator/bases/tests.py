import pytest
from . import storage
from . import schemas
from . import actions
from configurator.channels import actions as channel_actions
from configurator.channels import schemas as channel_schemas


@pytest.mark.asyncio
async def test_base_registry(database, async_client):

    test_query = schemas.BaseRegisterRequestParams(
        channel_id=5, base_tags=["abc", "def"]
    )

    await channel_actions.ActionRegisterChannel(
        db=database,
        query=channel_schemas.ChannelQueryParams(channel_id=test_query.channel_id),
    ).run()

    response = await async_client.post(
        f"/channels/{test_query.channel_id}/base",
        json={
            "base_tags": test_query.base_tags,
        },
    )

    assert response.status_code == 200

    bases = await storage.BaseStorage(db=database).get_base(
        channel_id=test_query.channel_id
    )

    assert bases.channel_id == test_query.channel_id
    assert bases.tags == test_query.base_tags
