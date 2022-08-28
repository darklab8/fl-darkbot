import pytest
from . import storage
from . import schemas
from . import actions
from .views import Paths
from configurator.channels import actions as channel_actions
from configurator.channels import schemas as channel_schemas


@pytest.mark.asyncio
async def test_base_registry(database, async_client):

    test_query = schemas.BaseRegisterRequestParams(
        channel_id=5, base_tags=["abc", "def"]
    )

    await channel_actions.ActionRegisterChannel(
        db=database,
        query=channel_schemas.ChannelCreateQueryParams(
            channel_id=test_query.channel_id
        ),
    ).run()

    response = await async_client.post(
        Paths.base,
        json=dict(test_query),
    )

    assert response.status_code == 200

    bases = await storage.BaseStorage(db=database).get_base(
        channel_id=test_query.channel_id
    )

    assert bases.channel_id == test_query.channel_id
    assert bases.tags == test_query.base_tags
