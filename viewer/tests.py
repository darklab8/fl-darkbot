import pytest
from scrappy.bases import rpc as rpc_bases


def test_check():
    assert True


@pytest.mark.asyncio
async def test_get_bases_data():
    size = 10
    result = await rpc_bases.ActionGetFilteredBases(
        query=rpc_bases.ActionGetFilteredBases.query_factory(
            page_size=size,
        )
    ).run()

    list_of_items = list(result)
    assert len(list_of_items) == size
    print(list_of_items)


@pytest.mark.asyncio
async def test_get_bases_settings():
    # TODO write correct test here :)
    pass
