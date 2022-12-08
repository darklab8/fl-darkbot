import asyncio

# from .actions import ActionGetAndParseAndSaveBases
from scrappy.bases import rpc as rpc_bases_data
from scrappy.bases import schemas as rpc_bases_data_schemas
from configurator.bases import rpc as rpc_bases_settings
from configurator.bases import schemas as rpc_bases_settings_schemas
from typing import List
from pydantic import BaseModel
import pathlib
from jinja2 import Template
from ..core.logger import base_logger

logger = base_logger.getChild(__name__)


class BaseInputForRender(BaseModel):
    health: str
    name: str
    affiliation: str


base_template_path = pathlib.Path(__file__).parent / "template.md"
with open(str(base_template_path), "r") as file_:
    base_template = Template(file_.read())


class ActionRenderBase:
    def run(self, bases: List[BaseInputForRender]):
        rendered = base_template.render(data=list([dict(base) for base in bases]))
        return rendered


class RenderedBasesForChannel(BaseModel):
    tag_id: str = "msg_id:||6a8c18f502cd8b67||"
    channel_id: int
    rendered: str


# class BaseOut(BaseModel):
#     id: int
#     name: str
#     affiliation: str
#     health: float
#     tid: int
#     timestamp: datetime

# class BasesOut(BaseModel):
#     __root__: List[BaseOut]
#     def __iter__(self):
#         for item in self.__root__:
#             yield item
class ActionRenderBases:
    def run(self):
        self.base_settings: rpc_bases_settings_schemas.BasesManyOut = asyncio.run(
            rpc_bases_settings.ActionGetBases(
                query=rpc_bases_settings.ActionGetBases.query_factory()
            ).run()
        )
        logger.debug(f"update_all_bases.result={self.base_settings}")

        self.base_data: rpc_bases_data_schemas.BasesOut = asyncio.run(
            rpc_bases_data.ActionGetFilteredBases(
                query=rpc_bases_data.ActionGetFilteredBases.query_factory(
                    page_size=9999999,
                )
            ).run()
        )

        # TODO Render bases per channel
        # Send to channel

    def render_results(self) -> List[RenderedBasesForChannel]:
        pass
