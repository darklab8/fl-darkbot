import aiohttp
from ..core import settings
from ..core import exceptions
from typing import Any


def is_succesful_request(code: int):
    if code < 299:
        return True
    return False


async def config_request(path: str, method: str, json: dict[str, Any] = {}):
    async with aiohttp.ClientSession() as session:
        match method:
            case "get":
                async with session.get(settings.CONFIGURATOR_API + path) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return await resp.json()
            case "post":
                async with session.post(
                    settings.CONFIGURATOR_API + path, json=json
                ) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return await resp.json()
            case "delete":
                async with session.delete(
                    settings.CONFIGURATOR_API + path, json=json
                ) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return await resp.json()
            case _:
                raise exceptions.NotImplementedMethod()


class RPCAction:
    async def run(self) -> aiohttp.ClientResponse | None:
        return await config_request(
            path=self.action.url,
            method=self.action.method.name,
            json=dict(self.query),
        )

    def __init__(self, query: Any):
        self.query = query
