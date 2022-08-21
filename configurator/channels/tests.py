import pytest


@pytest.mark.asyncio
async def test_try_getting_channel(database, async_client):
    response = await async_client.post(
        "/channels/5",
        json={
            "owner_id": 6,
            "owner_name": "Json",
        },
    )
    data = response.json()
    print(data)
