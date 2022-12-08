import pytest
from scrappy.bases import rpc as rpc_bases_data
from configurator.bases import rpc as rpc_bases_settings
from . import actions


def test_check():
    assert True


@pytest.mark.asyncio
async def test_get_bases_data():
    size = 10
    result = await rpc_bases_data.ActionGetFilteredBases(
        query=rpc_bases_data.ActionGetFilteredBases.query_factory(
            page_size=size,
        )
    ).run()

    list_of_items = list(result)
    assert len(list_of_items) == size
    print(list_of_items)


@pytest.mark.asyncio
async def test_get_bases_settings():
    print("started test")
    result = await rpc_bases_settings.ActionGetBases(
        query=rpc_bases_settings.ActionGetBases.query_factory()
    ).run()

    print("result:")
    print(result)


def test_render_bases():

    bases = list(
        [
            actions.BaseInputForRender(
                health=i,
                name=i,
                affiliation=i,
            )
            for i in range(10)
        ]
    )

    print(actions.ActionRenderBase().run(bases=bases))
