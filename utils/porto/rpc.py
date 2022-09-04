import abc
import aiohttp
from typing import Any
import enum
from pydantic import BaseModel
from utils.porto import AsyncAbstractAction


def is_succesful_request(code: int):
    if code < 299:
        return True
    return False


async def config_request(
    api_url: str, path: str, method: str, json: dict[str, Any] = {}
):
    async with aiohttp.ClientSession() as session:
        match method:
            case "get":
                async with session.get(api_url + path) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return await resp.json()
            case "post":
                async with session.post(api_url + path, json=json) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return await resp.json()
            case "delete":
                async with session.delete(api_url + path, json=json) as resp:
                    if is_succesful_request(resp.status):
                        print("successful request")
                    return await resp.json()
            case _:
                raise exceptions.NotImplementedMethod()


class AbstractRPCAction(abc.ABC):
    @abc.abstractproperty
    def action(self) -> AsyncAbstractAction:
        pass

    @abc.abstractproperty
    def api_url(self) -> str:
        pass

    @property
    def url(self) -> str:
        return self.action.url

    @property
    def method(self) -> enum.Enum:
        return self.action.method

    @property
    def response_factory(self) -> BaseModel:
        return self.action.response_factory

    @property
    def query_factory(self) -> BaseModel:
        return self.action.query_factory

    async def run(self):
        response = await config_request(
            api_url=self.api_url,
            path=self.url,
            method=self.method.name,
            json=dict(self.query),
        )

        return self.response_factory(**response)

    def __init__(self, query: Any):
        self.query = query
