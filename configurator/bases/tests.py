import pytest
from . import storage
from . import schemas
from . import actions
from configurator.channels import actions as channel_actions


@pytest.mark.asyncio
async def test_base_registry(database, async_client):

    base_registry_query = schemas.BaseRegisterRequestParams(
        channel_id=5, base_tags=["abc", "def"]
    )

    registry_action = channel_actions.ActionRegisterChannel
    await registry_action(
        db=database,
        query=registry_action.query_factory(channel_id=base_registry_query.channel_id),
    ).run()

    response = await async_client.post(
        actions.ActionRegisterBase.url,
        json=dict(base_registry_query),
    )

    assert response.status_code == 200

    bases = await storage.BaseStorage(db=database).get_base(
        channel_id=base_registry_query.channel_id
    )

    assert bases.channel_id == base_registry_query.channel_id
    assert bases.tags == base_registry_query.base_tags

    # get bases

    response = await async_client.post(
        actions.ActionGetBases.url,
        json={},
    )
    bases = response.json()
    assert response.status_code == 200
    assert len(bases) > 0
    print(bases)
