import pytest
from src.data_model import DataModel
from src.storage import Storage


@pytest.mark.asyncio
async def test_no_repeated_msgs_mocked():
    storage = Storage(IS_MOCKING_REQUESTS=True)
    api_data = storage.get_game_data({})
    data = DataModel(api_data=api_data)

    updating_api_data = storage.get_game_data(
            data.previous_forum_records)

    for record in updating_api_data.new_forum_records:
        data.previous_forum_records[record.title] = record

    data.update(updating_api_data)

    print("====================")
    print(data.api_data.new_forum_records)

    updating_api_data = storage.get_game_data(
            data.previous_forum_records)

    for record in updating_api_data.new_forum_records:
        data.previous_forum_records[record.title] = record

    data.update(updating_api_data)

    print("====================")
    print(data.api_data.new_forum_records)


# async def test_getting_channel_permissions():
#     storage = Storage(IS_MOCKING_REQUESTS=True)

#     channel_ids = [int(item) for item in storage.channels.keys()]

#     forbidden_channels = []
#     allowed_channels = []

#     for channel_id in channel_ids:
#         channel_info = self.bot.get_channel(channel_id)

#         if channel_info is None:
#             forbidden_channels.append(channel_id)
#         else:
#             allowed_channels.append(channel_id)

# @pytest.mark.asyncio
# async def test_no_repeated_msgs_live_testing():
#     storage = Storage()
#     api_data = storage.get_game_data({})
#     data = DataModel(api_data=api_data)

#     updating_api_data = await storage.a_get_game_data(
#             data.previous_forum_records)

#     for record in updating_api_data.new_forum_records:
#         data.previous_forum_records[record.title] = record

#     data.update(updating_api_data)

#     print("==================2=================")
#     print(data.api_data.new_forum_records)



