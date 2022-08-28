from utils.rest_api.message import MessageOk
from ..commons.requests import RPCAction
from .requests import config_request


class ActionPingConfig:
    async def run(self) -> MessageOk:
        data = await config_request(path="/", method="get")
        return MessageOk(**data)
