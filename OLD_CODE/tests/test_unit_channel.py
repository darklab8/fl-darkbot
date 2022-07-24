import pytest
from src.channel import ChannelConstroller
from src.mocks import MockDiscordMessageBus


@pytest.mark.asyncio
async def test_check_channel_controller():
    controller = ChannelConstroller(MockDiscordMessageBus(),
                                    'self.storage.unique_tag')
    await controller.delete_exp_msgs(1, 40)

    messages = await controller.get_tagged_msgs(1)
    for message in messages:
        controller.message_bus.delete(message)

    await controller.update_info(1, 'info')
